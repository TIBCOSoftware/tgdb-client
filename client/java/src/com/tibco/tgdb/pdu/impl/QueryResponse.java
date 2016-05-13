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
 * SVN Id: $Id: QueryResponse.java 687 2016-04-08 01:02:44Z cltran $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;

public class QueryResponse extends AbstractProtocolMessage {

	public int result;
	public long queryHashId;
    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        is.readInt(); // buf length
        is.readInt(); // checksum
        this.result = is.readInt(); // query result
        this.queryHashId = is.readLong();  // query hash Id
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
}
