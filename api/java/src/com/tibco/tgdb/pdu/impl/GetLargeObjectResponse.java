
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
 * File name : GetLargeObjectResponse.${EXT}
 * Created on: 10/15/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: GetLargeObjectResponse.java 3136 2019-04-25 23:50:36Z nimish $
 */

package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.ByteArrayOutputStream;
import java.io.IOException;

public class GetLargeObjectResponse extends AbstractProtocolMessage {

    ByteArrayOutputStream bos;
    long entityId;

    public GetLargeObjectResponse() {
        super();
    }

    public GetLargeObjectResponse(long authToken, long sessionId)
    {
        super(authToken, sessionId);
    }

    public byte[] getBuffer() {
        if (bos == null) return new byte[0];
        return bos.toByteArray(); //ReWrite ByteArrayOutputStream
    }

    public int getBufferLength() {
        if (bos == null) return 0;
        return bos.size();
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        int status = is.readInt();
        if (status > 0) throw new TGException(String.format("Read Large Object failed with status : %d", status));
        //Read the chunks.
        entityId = is.readLong();
        boolean bHasData = is.readBoolean();
        if (bHasData) {
            int numChunks = is.readInt();
            bos = new ByteArrayOutputStream();
            for (int i=0; i<numChunks; i++) {
                byte[] buf = is.readBytes();
                bos.write(buf);
            }
        }

    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.GetLargeObjectResponse;
    }
}
