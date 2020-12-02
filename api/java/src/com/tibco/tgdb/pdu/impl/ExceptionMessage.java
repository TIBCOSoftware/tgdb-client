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
 * File name : ExceptionMessage.${EXT}
 * Created on: 1/17/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: ExceptionMessage.java 4309 2020-09-08 16:32:46Z ssubrama $
 */


package com.tibco.tgdb.pdu.impl;

import static com.tibco.tgdb.exception.TGException.TGExceptionType.UnexpectedOperation;
import com.tibco.tgdb.exception.TGException.TGExceptionType;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;
import java.io.PrintWriter;
import java.io.StringWriter;

public class ExceptionMessage extends AbstractProtocolMessage {

    int servercode;
    String msg;
    TGExceptionType exceptionType;

    ExceptionMessage() {
        exceptionType = TGExceptionType.GeneralException;
        servercode = -1;
        msg = "Server General Exception";
    }

    public ExceptionMessage(TGExceptionType type, String msg) {
        this.msg = msg;
        this.exceptionType = type;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.ExceptionMessage;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        boolean isNull;
        servercode = is.readInt();
        if (servercode == 0) return;
        isNull = is.readBoolean();
        if (!isNull) msg = is.readUTF();
        mapServerCodeToExceptionType();

    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    public TGExceptionType getExceptionType() {
        return exceptionType;
    }

    public String getMessage() {
        return msg;
    }

    private void mapServerCodeToExceptionType() {
        //TODO

    }
}
