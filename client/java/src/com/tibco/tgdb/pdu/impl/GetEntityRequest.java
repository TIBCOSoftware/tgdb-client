package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;

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
 * File name :GetEntityRequest
 * Created on: 12/24/14
 * Created by: chung
 * <p/>
 * SVN Id: $Id: GetEntityRequest.java 583 2016-03-15 02:02:39Z vchung $
 */
public class GetEntityRequest extends AbstractProtocolMessage {
	private TGKey key;

    //0 - get, 1 - getbyid, 2 - get multiples, 10 - continue, 20 - close
    private short getCommand = 0;
	private int fetchSize = 1000;
	private short batchSize = 50;
	private short traversalDepth = 3;
	private short edgeFetchSize = 0; // zero means no limitation
    private int resultId = 0;
	
    public GetEntityRequest() {
    	super();
    }

    public GetEntityRequest(long authToken, long sessionId) {
    	super(authToken, sessionId);
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.GetEntityRequest;
    }

    @Override
    public boolean isUpdateable() {
        return true;
    }

    public void setKey(TGKey key) {
    	this.key = key;
    }

    public void setBatchSize(short size) {
    	if (size < 10 || size > 32767) {
    		batchSize = 50;
    	} else {
    		batchSize = size;
    	}
    }
    
    public short getBatchSize() {
    	return batchSize;
    }
    
    public void setFetchSize(int size) {
    	if (size < 0) {
    		fetchSize = 1000;
    	} else {
    		fetchSize = size;
    	}
    }
    
    public int getFetchSize() {
    	return fetchSize;
    }
    
    public void setEdgeFetchSize(short size) {
    	if (size < 0 || size > 32767) {
    		edgeFetchSize = 1000;
    	} else {
    		edgeFetchSize = size;
    	}
    }
    
    public int getEdgeFetchSize() {
    	return edgeFetchSize;
    }
    
    public void setTraversalDepth(short depth) {
    	if (depth < 1 || depth > 1000) {
    		this.traversalDepth = 3;
    	} else {
    		this.traversalDepth = depth;
    	}
    }
    
    public short getTraversalDepth() {
    	return traversalDepth;
    }

    public void setResultId(int id) {
        resultId = id;
    }

    public int getResultId() {
        return resultId;
    }

    public void setCommand(short cmd) {
        this.getCommand = cmd;
    }

    public short getCommand() {
        return getCommand;
    }

	@Override
	protected void writePayload(TGOutputStream os) throws TGException, IOException {
        os.writeShort(getCommand);
        os.writeInt(resultId);
        if (getCommand == 0 || getCommand == 1 || getCommand == 2) {
            os.writeInt(fetchSize);
            os.writeShort(batchSize);
            os.writeShort(traversalDepth);
            os.writeShort(edgeFetchSize);
		    key.writeExternal(os);
        }
	}

	@Override
	protected void readPayload(TGInputStream is) throws TGException, IOException {
		// TODO Auto-generated method stub
		
	}
}
