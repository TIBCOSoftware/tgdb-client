package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;

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
 * File name :HandshakeRequest
 * Created on: 12/24/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: HandshakeRequest.java 723 2016-04-16 19:21:18Z vchung $
 */
public class HandshakeRequest extends AbstractProtocolMessage {

    private boolean sslMode = false;
    private int challenge = 0;
    private RequestType requestType = RequestType.Invalid;

    @Override
    public VerbId getVerbId() {
        return VerbId.HandShakeRequest;
    }

    @Override
    public boolean isUpdateable() {
        return true;
    }

    public boolean getSslMode() {
        return sslMode;
    }

    public void setSslMode(boolean sslMode)
    {
        this.sslMode = sslMode;
    }

    public int getChallenge() {
        return challenge;
    }

    public void setChallenge(int challenge) {
        this.challenge = challenge;
    }

    public RequestType getRequestType() {
        return this.requestType;
    }

    public void setRequestType(RequestType type) {
        this.requestType = type;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {

    	//FIXME: Need to change the enum definition not to use ordinal
        os.writeByte((byte) this.requestType.ordinal());
        os.writeBoolean(sslMode);
        os.writeInt(challenge);

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        //Purely for testing.
        int b = is.readByte();
        this.requestType = RequestType.values()[b];
        this.sslMode = is.readBoolean();
        this.challenge = is.readInt();

    }

    public enum RequestType {
        Invalid,
        Initiate,
        ChallengeAccepted
    }
}
