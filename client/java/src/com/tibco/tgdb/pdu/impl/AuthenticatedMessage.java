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
 * File name : AuthenticatedMessage.${EXT}
 * Created on: 2/4/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: AuthenticatedMessage.java 583 2016-03-15 02:02:39Z vchung $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGAuthenticatedMessage;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;

public abstract class AuthenticatedMessage extends AbstractProtocolMessage implements TGAuthenticatedMessage {

//    long authToken;
//    long sessionId;
    int connectionId;
    String clientId;

    /*
    public long getAuthToken() {
        return authToken;
    }
    */

    /*
    public int getConnectionId() {
        return connectionId;
    }

    public void setConnectionId(int connectionId) {
        this.connectionId = connectionId;
    }

    public void setAuthToken(long authToken) {
        this.authToken = authToken;
    }

    public long getSessionId() {
        return sessionId;
    }

    public void setSessionId(long sessionId) {
        this.sessionId = sessionId;
    }
    */

    public String getClientId() {
        return clientId;
    }

    public void setClientId(String clientId) {
        this.clientId = clientId;
    }

    //FIXME: Is it ok to expose this?
    //Expose this so that it can be set
    //by the channel during sendRequest
    public void setRequestId(long requestId) {
    	this.requestId = requestId;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
        if ((authToken == -1) || (sessionId == -1))
            throw new TGException("Message not authenticated");

        os.writeLong(authToken);
        os.writeLong(sessionId);

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        this.authToken = is.readLong();
        this.sessionId = is.readLong();
    }
}
