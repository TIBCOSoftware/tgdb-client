/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except 
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * <p/>
 * File name : ChannelPinger.${EXT}
 * Created on: 1/6/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: MultiChannelPinger.java 723 2016-04-16 19:21:18Z vchung $
 */


package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.pdu.impl.PingMessage;
import com.tibco.tgdb.utils.TGEnvironment;

import java.util.HashSet;
import java.util.Iterator;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicBoolean;

public class MultiChannelPinger {

    private static MultiChannelPinger gInstance = new MultiChannelPinger();
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    HashSet<AbstractChannel> channelSet = new HashSet<>();
    ExecutorService executorService;
    Pinger pinger = new Pinger();



    private MultiChannelPinger() { }

    public static MultiChannelPinger getInstance() { return gInstance;}

    public void addChannel(AbstractChannel channel) {
        synchronized(channelSet) {
            channelSet.add(channel);
        }
    }

    public void removeChannel(AbstractChannel channel) {
        synchronized(channelSet) {
            channelSet.remove(channel);
        }
    }

    public synchronized void start() throws TGException {
         if (executorService == null) {
            executorService = Executors.newSingleThreadExecutor();
            executorService.execute(pinger);
        }

        return;
    }

    public void stop() {
        if (channelSet.size() == 0) {  //If there are no more channels to ping
            pinger.stop();
            executorService.shutdown();
            executorService = null;
        }
    }

    class Pinger implements Runnable {

        AtomicBoolean isRunning = new AtomicBoolean(true);
        @Override
        public void run() {

            long pingInterval = TGEnvironment.getInstance().getChannelPingInterval() * 1000;  //seconds
            Thread.currentThread().setName("TGChannelPinger");

            gLogger.log(TGLogger.TGLevel.Debug, "Pinger thread %s is using ping interval %d(ms)", 
            		Thread.currentThread().getName(), pingInterval);
            while (isRunning.get()) {
                try {
                    Thread.sleep(pingInterval);

                    Iterator<AbstractChannel> itr = MultiChannelPinger.this.channelSet.iterator();
                    while (itr.hasNext()) {

                        AbstractChannel channel = itr.next();

                        if (channel.needsPing()) {
                            try {
                                PingMessage ping = new PingMessage();
                                channel.send(ping);
                            }
                            catch (Exception ioe) {
                            	if (isRunning.get()) {
                            		gLogger.logException("Pinger invoke channel exception callbacks", ioe);
                            		channel.handleException(ioe);
                            	} else {
                            		gLogger.logException("Pinger is preparing to stop", ioe);
                            	}
                                channel.disablePing();
                            }
                        }
                        else {
                            channel.enablePing();
                        }
                    }
                }
                catch (InterruptedException ie) { //Should happen only when the thread is interrupted
                	gLogger.logException("Pinger thread stopped with exception", ie);
                    break;
                }
            }
            gLogger.log(TGLogger.TGLevel.Debug, "Pinger thread is exiting");
        }

        void stop() {
            isRunning.set(false);
        }
    }
}
