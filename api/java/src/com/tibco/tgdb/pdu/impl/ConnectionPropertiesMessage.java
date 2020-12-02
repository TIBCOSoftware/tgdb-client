/**
 * Copyright (c) 2020 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : SetConnectionPropertiesRequest.
 * Created on: 7/17/20
 * Created by: suresh
 * <p/>
 * SVN Id: $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.utils.TGProperties;

import java.io.IOException;
import java.util.Map;

public class ConnectionPropertiesMessage extends AbstractProtocolMessage {

    TGProperties<String, String> properties = null;

    public ConnectionPropertiesMessage() {

    }

    public ConnectionPropertiesMessage(long authToken, long sessionId) {
        super(authToken, sessionId);
    }

    public void setProperties(TGProperties<String, String> properties) {
        this.properties = properties;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
        if (properties == null) {
            throw new TGException("Can't write Null Properties");
        }
        os.writeInt(this.properties.size());
        for(Map.Entry<String,String> entry : properties.entrySet()) {
            String key = entry.getKey();
            String value = entry.getValue();
            System.out.printf("%s = %s\n", key, value);
            os.writeUTF(key);
            if (value != null) {
                os.writeBoolean(false);
                os.writeUTF(value);
            }
            else {
                os.writeBoolean(true);
            }
        }
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {

    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.ConnectionPropertiesMessage;
    }
}
