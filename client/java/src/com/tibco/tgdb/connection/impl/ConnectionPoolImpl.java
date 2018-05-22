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
 * SVN Id: $Id: ConnectionPoolImpl.java 2179 2018-03-29 21:49:54Z ssubrama $
 */


package com.tibco.tgdb.connection.impl;

import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.channel.impl.ChannelFactoryImpl;
import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionExceptionListener;
import com.tibco.tgdb.connection.TGConnectionPool;
import com.tibco.tgdb.exception.TGConnectionTimeoutException;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.TGProperties;

import java.util.*;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;


public class ConnectionPoolImpl implements TGConnectionPool{

    static enum ConnectionPoolState {
        Initialized,
        Connected,
        Disconnected,
        Disconnecting, Connecting, Stopped
    }

    BlockingQueue<TGConnection> connPool;
    List<TGConnection> connlist;

    boolean bUseDedicateChannel;
    int poolSize;
    int connectReserveTimeOut;
    TGProperties<String,String> properties;
    TGConnectionExceptionListener lsnr = null;
    Map<Thread, TGConnection> consumers;
    ConnectionPoolState state;
    private ReadWriteLock adminLock = new ReentrantReadWriteLock();

    public ConnectionPoolImpl(TGChannelUrl url, int poolSize, TGProperties<String, String> properties) throws TGException
    {
        this.connPool   = new ArrayBlockingQueue<TGConnection>(poolSize, true);
        this.connlist   = new ArrayList<>();
        this.poolSize   = poolSize;
        this.properties = properties;
        this.consumers  = new HashMap<>();
        this.connectReserveTimeOut = Integer.parseInt(properties.getProperty(ConfigName.ConnectionReserveTimeoutSeconds));
        bUseDedicateChannel = Boolean.parseBoolean(properties.getProperty(ConfigName.ConnectionPoolUseDedicatedChannelPerConnection, "false"));
        TGChannel channel = null;
        ChannelFactoryImpl channelFactory = (ChannelFactoryImpl) ChannelFactoryImpl.getInstance();
        for (int i=0; i<poolSize; i++)
        {
            if ((channel == null) || (bUseDedicateChannel)) {
                channel = channelFactory.createChannel(url, properties);
            }
            TGConnection connection = new ConnectionImpl(this, channel, properties);
            connPool.add(connection);
            connlist.add(connection);
        }
        state = ConnectionPoolState.Initialized;
    }

    @Override
    public void connect() throws Exception {
        adminLock.readLock().lock();
        try {
            if (this.state == ConnectionPoolState.Connected) {
                throw new TGException("ConnectionPool is already connected. Disconnect and then reconnect");
            }
            this.state = ConnectionPoolState.Connecting;
            for (TGConnection connection : connPool) {
                connection.connect();
            }
            state = ConnectionPoolState.Connected;
        }
        finally {
            adminLock.readLock().unlock();
        }

    }

    @Override
    public synchronized void setExceptionListener(TGConnectionExceptionListener lsnr)  {
        this.lsnr = lsnr;
    }

    @Override
    public void disconnect() throws Exception
    {
        adminLock.readLock().lock();
        try {
            if (this.state != ConnectionPoolState.Connected) {
                throw new TGException(String.format("ConnectionPool is not connected. State is :%s", this.state));
            }
            this.state = ConnectionPoolState.Disconnecting;
            for (TGConnection conn : connlist) {
                conn.disconnect();
            }
        }
        finally {
            this.state = ConnectionPoolState.Disconnected;
            adminLock.readLock().unlock();
        }

    }

    @Override
    public TGConnection get() throws TGException {

        adminLock.readLock().lock();
        try {
            if (this.state != ConnectionPoolState.Connected) {
                throw new TGException(String.format("ConnectionPool is not connected. State is :%s", this.state));
            }
            TGConnection connection = consumers.get(Thread.currentThread());
            if (connection != null) return connection;

            connection =  connPool.poll(this.connectReserveTimeOut, TimeUnit.SECONDS);
            if (connection == null) {
                throw new TGConnectionTimeoutException(String.format("Failed to get connection within :%d seconds", this.connectReserveTimeOut));
            }
            consumers.put(Thread.currentThread(), connection);
            return connection;
        }
        catch (InterruptedException ie) {
              throw  TGException.buildException("ConnectionPool interrupted", "ThreadIntr", ie);
        }
        finally {
            adminLock.readLock().unlock();
        }
    }

    TGConnection getConnection() throws TGException {

        adminLock.readLock().lock();
        try {
            TGConnection connection = consumers.get(Thread.currentThread());
            if (connection != null) return connection;

            connection =  connPool.poll(this.connectReserveTimeOut, TimeUnit.SECONDS);
            if (connection == null) {
                throw new TGConnectionTimeoutException(String.format("Failed to get connection within :%d seconds", this.connectReserveTimeOut));
            }
            consumers.put(Thread.currentThread(), connection);
            return connection;
        }
        catch (InterruptedException ie) {
              throw  TGException.buildException("ConnectionPool interrupted", "ThreadIntr", ie);
        }
        finally {
            adminLock.readLock().unlock();
        }
    }

    @Override
    public void release(TGConnection conn) {
        adminLock.readLock().lock();
        try {
            consumers.remove(Thread.currentThread());
            connPool.offer(conn);
        }
        finally {
            adminLock.readLock().unlock();
        }
    }

    @Override
    public int getPoolSize() {
        return poolSize;
    }


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


    void adminlock() {
        adminLock.readLock().lock();
    }

    void adminUnlock() {
        adminLock.readLock().unlock();
    }
}
