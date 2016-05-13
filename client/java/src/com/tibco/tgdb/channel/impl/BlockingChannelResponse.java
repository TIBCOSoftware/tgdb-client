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
 * File name : BlockingChannelResponse.${EXT}
 * Created on: 2/5/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: BlockingChannelResponse.java 583 2016-03-15 02:02:39Z vchung $
 */


package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.channel.TGChannelResponse;
import com.tibco.tgdb.pdu.TGMessage;

import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class BlockingChannelResponse implements TGChannelResponse {

    final Lock lock = new ReentrantLock();
    final Condition cond = lock.newCondition();
    Status status;
    long requestId;
    long timeout;
    TGMessage reply = null;


    public BlockingChannelResponse(long requestId) {
        this.requestId = requestId;
        this.status = Status.Waiting;
        this.timeout = -1;
    }

    public BlockingChannelResponse(long requestId, long timeout) {
        this.requestId = requestId;
        this.status = Status.Waiting;
        this.timeout = timeout;
    }


    @Override
    public boolean isBlocking() {
        return true;
    }

    @Override
    public Status getStatus() {
        lock.lock();
        try {
            return status;
        }
        finally {
            lock.unlock();
        }
    }

    @Override
    public long getRequestId() {
        return requestId;
    }

    @Override
    public void await(StatusTester tester) throws InterruptedException {
        lock.lock();
        try {
            while (tester.test(this.status)) {
                cond.await(timeout, TimeUnit.MILLISECONDS);
            }
        }
        catch(InterruptedException ie) {
            throw ie;
        }
        finally {
            lock.unlock();
        }

    }

    @Override
    public void signal(Status status) {
        lock.lock();
        try {
            this.status = status;
            cond.signalAll();

        }
        finally {
            lock.unlock();
        }
    }

    @Override
    public TGMessage getReply() {
        return reply;
    }

    @Override
    public void setReply(TGMessage msg) {
        lock.lock();
        try {
            this.reply = msg;
            this.status = Status.Ok;
            cond.signalAll();
        }
        finally {
            lock.unlock();
        }
    }

    @Override
    public Callback getCallback() {
        throw new RuntimeException("Not available for BlockingChannelResponse");
    }

    @Override
    public void setRequestId(long requestId) {
    	this.requestId = requestId;
    }

    public void reset()
    {
        lock.lock();
        try {
            this.status = Status.Waiting;
            this.reply = null;
        }
        finally {
            lock.unlock();
        }
    }
}
