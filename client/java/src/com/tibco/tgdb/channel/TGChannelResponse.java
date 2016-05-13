package com.tibco.tgdb.channel;

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
 * <p>
 * File name : TGChannelResponse.${EXT}
 * Created on: 2/5/15
 * Created by: suresh
 * <p>
 * SVN Id: $Id: TGChannelResponse.java 583 2016-03-15 02:02:39Z vchung $
 */

import com.tibco.tgdb.pdu.TGMessage;

/**
 * An interface for response object
 */
public interface TGChannelResponse {

    public enum Status {
        Waiting,
        Ok,
        Pushed,
        Resend,
        Disconnected,
        Closed
    }

    public interface Callback {
        void onResponse(TGMessage msg);
    }

    public interface StatusTester {
        boolean test(Status status);
    }

    boolean isBlocking();

    Status getStatus();

    long getRequestId();

    void setRequestId(long requestId);
    
    /**
     * Await while the response is in this state
     * @param tester
     * @throws InterruptedException
     */
    void await(StatusTester tester) throws InterruptedException;

    /**
     * Signal the new status
     * @param status
     */
    void signal(Status status);

    /**
     * Get Reply object
     * @return
     */
    TGMessage getReply();

    /**
     * Set the Reply
     * @param msg
     */
    void setReply(TGMessage msg);

    /**
     * Get a Callback object
     * @return
     */
    Callback getCallback();

    void reset();

}
