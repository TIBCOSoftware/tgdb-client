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
 * File name : QueryRequest.${EXT}
 * Created on: 2/4/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: QueryRequest.java 687 2016-04-08 01:02:44Z cltran $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.query.TGQuery;

import java.io.IOException;

public class QueryRequest extends AbstractProtocolMessage {

    private String queryExpr;
    private long queryHashId;
    private int command;
    private TGQuery queryObject;
    
    static TGLogger gLogger = TGLogManager.getInstance().getLogger();

    public QueryRequest() {
    	super();
    }
    
    public QueryRequest(long authToken, long sessionId) {
    	super(authToken, sessionId);
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
        int startPos = os.getPosition();
        os.writeInt(1);
        os.writeInt(1);
        gLogger.log(TGLevel.Debug, "Entering query request writePayload at output buffer position at : %d", startPos);
        os.writeInt(command);
        // CREATE, EXECUTE.
        if (command == 1 || command == 2) {
            os.writeUTF(this.queryExpr);
        }
        // EXECUTEID, CLOSE
        else if (command == 3 || command == 4){
        	os.writeLong(this.queryHashId);
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
        return VerbId.QueryRequest;
    }

    public void setQuery(String expr) {
        this.queryExpr = expr;
    }
    
    public void setQueryHashId(long queryHashId) {
        this.queryHashId = queryHashId;
    }
    
    public void setCommand(int command) {
        this.command = command;
    }
}
