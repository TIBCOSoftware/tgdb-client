/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : SessionForcefullyTermincated.${EXT}
 * Created on: 3/18/18
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

public class SessionForcefullyTerminated extends ExceptionMessage {


    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        msg = is.readUTF();
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.SessionForcefullyTerminated;
    }

    public String getKillString() { return msg; }


}
