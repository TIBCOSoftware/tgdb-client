
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
 * File name : GetLargeObjectRequest.${EXT}
 * Created on: 10/15/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: GetLargeObjectRequest.java 3631 2019-12-11 01:12:03Z ssubrama $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;

public class GetLargeObjectRequest extends AbstractProtocolMessage {

    private long entityId;
    private boolean bDecrypt;

    public GetLargeObjectRequest() { super();}

    public GetLargeObjectRequest(long authToken, long sessionId)
    {
        super(authToken, sessionId);
    }

    public void setEntityId(long entityId) {
        this.entityId = entityId;
    }

    public void setDecryptionEnable(boolean bDecrypt) {
        this.bDecrypt = bDecrypt;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
        os.writeLong(entityId);
        os.writeBoolean(bDecrypt);

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        entityId = is.readLong();
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.GetLargeObjectRequest;
    }
}
