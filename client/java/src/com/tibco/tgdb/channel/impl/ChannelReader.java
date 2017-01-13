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
 * File name : ChannelReader.${EXT}
 * Created on: 1/6/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: ChannelReader.java 955 2016-09-02 16:32:55Z vchung $
 */


package com.tibco.tgdb.channel.impl;


import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.VerbId;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicLong;

public class ChannelReader implements Runnable {

    AbstractChannel channel;
    AtomicBoolean isRunning = new AtomicBoolean(false);

    ExecutorService executorService = null;

    static AtomicLong gReaders = new AtomicLong(0);
    long readerNum = gReaders.getAndIncrement();

    ChannelReader (AbstractChannel channel) {
        this.channel = channel;
    }

    public synchronized void start() {
        if (!isRunning.get()) {
            executorService = Executors.newSingleThreadExecutor();
            executorService.execute(this);
            isRunning.set(true);
        }
    }

    public synchronized  void stop() {
        if (isRunning.get()) {
            isRunning.set(false);
            executorService.shutdown();
            executorService = null;
        }
    }

    @Override
    public void run() {
        Thread.currentThread().setName(String.format("TGLinkReader@%s-%d", channel.clientId, readerNum));

        while(isRunning.get()) {
            try {
                if (Thread.currentThread().isInterrupted())
                    return;

                TGMessage msg = channel.readWireMsg();

                if (channel.isClosed())
                    return;

                if ((Thread.currentThread().isInterrupted() || (!isRunning.get()))) {
                    if (!isRunning.get()) {
                        AbstractChannel.gLogger.log(TGLogger.TGLevel.Info, String.format("Thread %s interrupted. Stopping reader", Thread.currentThread().getName()));
                    }
                    return;
                }


                if ((msg == null) || (msg.getVerbId() == VerbId.PingMessage))
                    continue;

                channel.processMessage(msg);

                if (channel.isClosed())
                    return;
            }

            catch(Exception e) {
                //AbstractChannel.gLogger.logException("Channel reader thread interrupted.  Stopping reader...", e);
//                AbstractChannel.gLogger.log(TGLogger.TGLevel.Info, String.format("Thread %s interrupted. Stopping reader", Thread.currentThread().getName()));
                if ((Thread.currentThread().isInterrupted() || (!isRunning.get()))) {
                    AbstractChannel.gLogger.log(TGLogger.TGLevel.Debug, "Exiting channel reader thread");
                    return;
                }

                if ( channel.handleException(e) == 1) continue;  //SS:TODO instead of 1, specify an enum
                return;
            }
        }
    }
}
