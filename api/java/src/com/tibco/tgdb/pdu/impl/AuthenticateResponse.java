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
 * File name :AuthenticateResponse
 * Created on: 12/24/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: AuthenticateResponse.java 3631 2019-12-11 01:12:03Z ssubrama $
 */

package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;


public class AuthenticateResponse extends AbstractProtocolMessage {

    private boolean bSuccess = false;
    long authToken = -1;
    long sessionId = -1;
    byte certBuffer[];
    int  errorStatus;

    @Override
    public VerbId getVerbId() {
        return VerbId.AuthenticateResponse;
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
        os.writeBoolean(bSuccess);
        os.writeLong(authToken);
        os.writeLong(sessionId);
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {

        bSuccess = is.readBoolean();
        if (!bSuccess) {
            errorStatus = is.readInt();
            return;
        }
        authToken = is.readLong();
        sessionId = is.readLong();
        certBuffer = is.readBytes();

    }

    public boolean isSuccess() { return bSuccess;}
    public void setSuccess(boolean b) { this.bSuccess = b;}

    public long getAuthToken() { return authToken;}
    public void setAuthToken(long l) { this.authToken = l;}

    public long getSessionId() { return sessionId;}
    public void setSessionId(long l) { this.sessionId = l;}

    public byte[] getServerCertificateBuffer() { return this.certBuffer;}


}
