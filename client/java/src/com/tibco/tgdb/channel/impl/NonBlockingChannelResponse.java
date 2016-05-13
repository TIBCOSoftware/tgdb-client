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
 * File name : NonBlockingChannelResponse.${EXT}
 * Created on: 2/5/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: NonBlockingChannelResponse.java 583 2016-03-15 02:02:39Z vchung $
 */


package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.channel.TGChannelResponse;
import com.tibco.tgdb.pdu.TGMessage;

import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class NonBlockingChannelResponse implements TGChannelResponse {

    long requestId;
    Callback callback;
    Status status = Status.Waiting;
    Lock lock = new ReentrantLock();


    public NonBlockingChannelResponse(long requestId, Callback callback)
    {
        this.requestId = requestId;
        this.callback = callback;
    }
    @Override
    public boolean isBlocking() {
        return false;
    }

    @Override
    public Status getStatus() {
        return status;
    }

    @Override
    public long getRequestId() {
        return requestId;
    }

    @Override
    public void setRequestId(long requestId) {
        this.requestId = requestId;
    }

    @Override
    public void await(StatusTester tester) throws InterruptedException {
        //Nothing to do
    }

    @Override
    public void signal(Status status) {
        lock.lock();
        try {
            this.status = status;
        }
        finally {
            lock.unlock();
        }
    }

    @Override
    public TGMessage getReply() {
        return null;
    }

    @Override
    public void setReply(TGMessage msg) {
        lock.lock();
        try {
            this.status = Status.Ok;
            if (callback != null) {
                callback.onResponse(msg);
            }
        }
        finally {
            lock.unlock();
        }
    }

    @Override
    public Callback getCallback() {
        return callback;
    }

    @Override
    public void reset() {
        lock.lock();
        try {
            this.status = Status.Waiting;

        }
        finally {
            lock.unlock();
        }
    }
}
