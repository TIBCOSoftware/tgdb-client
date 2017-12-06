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
 * File name : GetEntityResponse.${EXT}
 * Created on: 2/4/15
 * Created by: chung 
 * <p/>
 * SVN Id: $Id: GetEntityResponse.java 583 2016-03-15 02:02:39Z vchung $
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

public class GetEntityResponse extends AbstractProtocolMessage {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private TGInputStream entityStream;
    private boolean hasResult = false;
    private int resultId = 0;
    private int totalCount = 0;
    private int resultCount = 0;

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
    	gLogger.log(TGLevel.Debug, "Entering get entity response readPayload");
    	if (is.available() == 0) {
    		gLogger.log(TGLevel.Debug, "Get entity response has no data");
    		return;
    	}
    	entityStream = is;
        resultId = is.readInt();
    	long pos = is.getPosition();
    	int totalCount = is.readInt();
    	if (totalCount > 0) {
    		hasResult = true;
    	}
    	is.setPosition(pos);
    	gLogger.log(TGLevel.Debug, "Received %d get entities", totalCount);
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.GetEntityResponse;
    }

    public TGInputStream getEntityStream() {
    	return entityStream;
    }

    public boolean hasResult() {
    	return hasResult;
    }

    public int getResultId() {
        return resultId;
    }

    public int getTotalCount() {
        return totalCount;
    }

    public int getResultCount() {
        return resultCount;
    }
}

