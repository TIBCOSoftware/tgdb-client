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
 * File name : ConnectionPoolImpl.${EXT}
 * Created on: 1/10/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: ConnectionPoolImpl.java 583 2016-03-15 02:02:39Z vchung $
 */


package com.tibco.tgdb.connection.impl;

import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionExceptionListener;
import com.tibco.tgdb.connection.TGConnectionPool;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.utils.TGProperties;

import java.util.ArrayList;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;


public class ConnectionPoolImpl implements TGConnectionPool, TGChannel.LinkEventHandler {

    Queue<TGConnection> connPool;
    List<TGConnection> connlist;

    boolean bUseDedicateChannel;
    int poolSize;
    TGProperties<String,String> properties;
    TGConnectionExceptionListener lsnr = null;
    private ReadWriteLock adminLock = new ReentrantReadWriteLock();

    public ConnectionPoolImpl(List<TGChannel> channels, boolean bUseDedicatedChannel, int poolSize, TGProperties<String, String> properties) throws TGException
    {
        this.connPool = new ConcurrentLinkedQueue<>();
        this.connlist = new ArrayList<>();
        this.bUseDedicateChannel = bUseDedicatedChannel;
        this.poolSize            = poolSize;
        this.properties          = properties;



        for (int i=0; i<poolSize; i++)
        {
            TGChannel channel = !bUseDedicatedChannel ? channels.get(0) : channels.get(i);
            TGConnection connection = new ConnectionImpl(this, channel, properties);
            connPool.add(connection);
            connlist.add(connection);
            channel.setLinkEventHandler(this);
        }
    }

    @Override
    public void connect() throws Exception {

        for(TGConnection connection : connPool) {
            connection.connect();
        }

    }

    @Override
    public synchronized void setExceptionListener(TGConnectionExceptionListener lsnr)  {
        this.lsnr = lsnr;
    }

    @Override
    public void disconnect() throws Exception {

        for (TGConnection conn : connPool) {
            conn.disconnect();
        }
    }



    @Override
    public TGConnection get() {
        adminLock.readLock().lock();
        try {
            return connPool.remove();
        }
        finally {
            adminLock.readLock().unlock();
        }
    }

    @Override
    public void release(TGConnection conn) {
         connPool.offer(conn);  //No need to lock. Let it comeback to the pool
    }

    @Override
    public int getPoolSize() {
        return poolSize;
    }


    @Override
    public void onException(Exception ex, boolean duringClose) {

        adminLock.writeLock().lock();
        try {
            for (TGConnection conn : connlist) {
                conn.disconnect();
            }

        }
        finally {
            adminLock.writeLock().unlock();
        }

        if (lsnr != null) {
            lsnr.onException(ex);
        }
    }

    @Override
    public boolean onReconnect() {
        return false;
    }

    @Override
    public String getTerminatedText() {
        return "ConnectionPool terminated.";
    }

    void adminlock() {
        adminLock.readLock().lock();
    }

    void adminUnlock() {
        adminLock.readLock().unlock();
    }
}
