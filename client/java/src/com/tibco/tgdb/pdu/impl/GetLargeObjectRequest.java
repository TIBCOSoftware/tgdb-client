/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : GetLargeObjectRequest.${EXT}
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

import java.io.IOException;

public class GetLargeObjectRequest extends AbstractProtocolMessage {

    private long entityId;

    public GetLargeObjectRequest() { super();}

    public GetLargeObjectRequest(long authToken, long sessionId)
    {
        super(authToken, sessionId);
    }

    public void setEntityId(long entityId) {
        this.entityId = entityId;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
        os.writeLong(entityId);

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
