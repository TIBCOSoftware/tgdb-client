
/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 * File name :HandshakeResponse
 * Created on: 12/24/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: HandshakeResponse.java 3138 2019-04-25 23:53:21Z nimish $
 */

package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;

public class HandshakeResponse extends AbstractProtocolMessage {

    long challenge;
    ResponseStatus responseStatus;
    long version;
    String errorMessage;

    public String getErrorMessage() {
		return errorMessage;
	}
    
	public long getVersion() {
		return version;
	}

	public void setVersion(long version) {
		this.version = version;
	}

	@Override
    public VerbId getVerbId() {
        return VerbId.HandShakeResponse;
    }

    public long getChallenge() {
        return challenge;
    }

    public void setChallenge(long c) {
        challenge = c;
    } //For testing purpose

    public ResponseStatus getResponseStatus() {
        return responseStatus;
    }

    public void setResponseStatus(ResponseStatus status) {
        this.responseStatus = status;
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
        //This is purely for testing. Client never writes out the response.
        os.writeByte(responseStatus.ordinal());
        os.writeLong(challenge);
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
    	byte byteRead = is.readByte();
        this.responseStatus = ResponseStatus.values()[(int) byteRead];
        this.challenge = is.readLong();
    	if (this.responseStatus == ResponseStatus.ChallengeFailed) 
        {
    		byte[] readBytes = is.readBytes();
    		errorMessage = new String (readBytes);
    		return;
        }
    }

    public enum ResponseStatus {
        Invalid,
        AcceptChallenge,
        ProceedWithAuthentication,
        ChallengeFailed
    }
}
