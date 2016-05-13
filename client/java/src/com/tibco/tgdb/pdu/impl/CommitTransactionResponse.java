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
 * File name : CommitTransactionResponse.${EXT}
 * Created on: 2/4/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: CommitTransactionResponse.java 771 2016-05-05 11:40:52Z vchung $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGEntity.TGEntityKind;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.impl.AbstractEntity;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;

public class CommitTransactionResponse extends AbstractProtocolMessage {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private ArrayList<Long> addedIdList = null; 
    private ArrayList<Long> updatedIdList = null; 
    private ArrayList<Long> removedIdList = null; 
    private ArrayList<Integer> attrDescIdList = null; 
    private int attrDescCount = 0;
    private int entityCount = 0;
    private TGGraphObjectFactory gof = null;
    private TGInputStream entityStream = null;

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
    	gLogger.log(TGLevel.Debug, "Entering commit transaction response readPayload");
    	is.readInt(); // buf length
    	is.readInt(); // checksum
    	is.readInt();// status code - currently zero
        while (is.available() > 0) {
    		short opCode = is.readShort();
    		switch (opCode) {
    		case 0x1010: 
                attrDescIdList = new ArrayList<Integer>();
                attrDescCount = is.readInt();
                for (int i=0; i<attrDescCount; i++) {
                    int tempId = is.readInt();
                    attrDescIdList.add(tempId);
                    int realId = is.readInt();
                    attrDescIdList.add(realId);
                }
                gLogger.log(TGLevel.Debug, "Received %d attr desc", attrDescCount);
    			break;
    		case 0x1011:
                addedIdList = new ArrayList<Long>();
                entityCount = is.readInt();
                for (int i=0; i<entityCount; i++) {
                    Long tempId = is.readLong();
                    addedIdList.add(tempId);
                    Long realId = is.readLong();
                    addedIdList.add(realId);
                }
                gLogger.log(TGLevel.Debug, "Received %d entity", entityCount);
    			break;
            case 0x1012:
                gLogger.log(TGLevel.Debug, "Received update results");
                break;
            case 0x1013:
                removedIdList = new ArrayList<Long>();
                entityCount = is.readInt();
                for (int i=0; i<entityCount; i++) {
                    Long id = is.readLong();
                    removedIdList.add(id);
                }
                gLogger.log(TGLevel.Debug, "Received %d delete results", entityCount);
                break;
            case 0x6789:
            	entityStream = is;
            	long pos = is.getPosition();
            	int count = is.readInt();
            	is.setPosition(pos);
                gLogger.log(TGLevel.Debug, "Received %d debug entities", count);
                return;
            default:
                break;
    		}
        }
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.CommitTransactionResponse;
    }

    public int getAttrDescCount() {
        return attrDescCount;
    }

    public int getAddedEntityCount() {
        return entityCount;
    } 

    public List<Integer> getAttrDescIdList() {
        return attrDescIdList;
    }

    public List<Long> getAddedIdList() {
        return addedIdList;
    }
    
    public TGInputStream getEntityStream() {
    	return entityStream;
    }

}
