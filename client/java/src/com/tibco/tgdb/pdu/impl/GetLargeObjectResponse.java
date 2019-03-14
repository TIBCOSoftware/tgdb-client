/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : GetLargeObjectResponse.${EXT}
 * Created on: 10/15/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
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
