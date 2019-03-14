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
 * File name : QueryResponse.${EXT}
 * Created on: 2/4/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: QueryResponse.java 2659 2018-11-07 16:04:58Z vchung $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;

public class QueryResponse extends AbstractProtocolMessage {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private TGInputStream entityStream;
    private boolean hasResult = false;
    private int totalCount = 0;
    private int resultCount = 0;
	public int result;
	public long queryHashId;
    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
    	gLogger.log(TGLevel.Debug, "Entering query response readPayload");
    	if (is.available() == 0) {
    		gLogger.log(TGLevel.Debug, "Query response has no data");
    		return;
    	}
    	entityStream = is;
        is.readInt(); // buf length
        is.readInt(); // checksum
        this.result = is.readInt(); // query result
        this.queryHashId = is.readLong();  // query hash Id
        int syntax = is.readByte();
        resultCount = is.readInt();
        if (resultCount > 0) {
        	hasResult = true;
        }
        if (syntax == 1) {
    	   	totalCount = is.readInt();
    	   	gLogger.log(TGLevel.Debug, "Query has %d result entities and %d total entities", resultCount, totalCount);
        } else {
    	   	gLogger.log(TGLevel.Debug, "Query has %d result count", resultCount);
        }
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.QueryResponse;
    }
    
    public int getResult() {
        return this.result;
    }
    
    public long getQueryHashId() {
        return this.queryHashId;
    }

    public TGInputStream getEntityStream() {
    	return entityStream;
    }

    public boolean hasResult() {
    	return hasResult;
    }

    public int getTotalCount() {
        return totalCount;
    }

    public int getResultCount() {
        return resultCount;
    }
}
