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
 * File name : CommitTransactionMessage.${EXT}
 * Created on: 2/4/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: CommitTransactionRequest.java 2344 2018-06-11 23:21:45Z ssubrama $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGEntityId;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;
import java.util.LinkedHashMap;
import java.util.Set;

public class CommitTransactionRequest extends AbstractProtocolMessage {

    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private LinkedHashMap<Long, TGEntity> addedList; 
    private LinkedHashMap<Long, TGEntity> updatedList;
    private LinkedHashMap<Long, TGEntity> removedList;
    private Set<TGAttributeDescriptor> attrDescSet;

    LinkedHashMap<Long, TGEntity> changedList;
    
    public CommitTransactionRequest() {
    	super();
    }

    public CommitTransactionRequest(long authToken, long sessionId) {
    	super(authToken, sessionId);
    }

    public void addCommitLists(LinkedHashMap<Long, TGEntity> addedList, LinkedHashMap<Long, TGEntity> updatedList,
    		LinkedHashMap<Long, TGEntity> removedList, Set<TGAttributeDescriptor> attrDescSet) {
    	this.addedList = addedList;
    	this.updatedList =updatedList;
    	this.removedList = removedList;
    	this.attrDescSet = attrDescSet;
    }

    //Need to make each attribute as modified when changed especially when it's brand new so that
    //we can send only modified values to the server
    //How to handle exception inside stream.
    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
    	int startPos = os.getPosition();
    	os.writeInt(0); // This is for the commit buffer length
    	os.writeInt(0); // This is for the checksum for the commit buffer to be added later.  Currently not used
    	gLogger.log(TGLevel.Debug, "Entering commit transaction request writePayload at output buffer position at : %d", startPos);
    	//<A> for attribute descriptor, <N> for node desc definitions, <E> for edge desc definitions
    	//meta should be sent before the instance data
    	if (!attrDescSet.isEmpty()) {
    		os.writeShort(0x1010); // for attribute descriptor
    		//There should be nothing after the marker due to no new attribute descriptor
    		//Need to check for new descriptor only with attribute id as negative number
    		// check for size overrun
    		os.writeInt((int) attrDescSet.stream().filter(e -> e.getAttributeId() < 0).count());
    		attrDescSet.stream().forEach(e -> {try{e.writeExternal(os);} catch(TGException tex){} catch(IOException iex){}});
    	}
    	//FIXME: Need to deal with stream related exception handling such as IOException below
    	if (!addedList.isEmpty()) {
    		os.writeShort(0x1011); // for entity creation
    		os.writeInt(addedList.size()); // check for size 
            addedList.entrySet().stream().forEach(e -> {try{e.getValue().writeExternal(os);} catch(TGException ex) {} catch(IOException iex){}});
    	}
    	//FIXME: Need to write only the modified attributes
    	if (!updatedList.isEmpty()) {
    		os.writeShort(0x1012); // for entity update
    		os.writeInt(updatedList.size()); // check for size 
            updatedList.entrySet().stream().forEach(e -> {try{e.getValue().writeExternal(os);} catch(TGException ex) {} catch(IOException iex){}});
    	}
    	//FIXME: Need to write the Id only
    	if (!removedList.isEmpty()) {
    		os.writeShort(0x1013); // for deleted entities
    		os.writeInt(removedList.size());
//            removedList.entrySet().stream().forEach(e -> {try{os.writeLong(e.getKey());} catch(IOException iex){}});
            removedList.entrySet().stream().forEach(e -> {try{e.getValue().writeExternal(os);} catch(TGException ex) {} catch(IOException iex){}});
    	}
    	int currPos = os.getPosition();
    	int length = currPos - startPos;
    	os.writeIntAt(startPos, length);
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
    	gLogger.log(TGLevel.Debug, "Commit transaction request readPayload is called");
    	//Commit response need to send back real id for all entities and descriptors.
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.CommitTransactionRequest;
    }
}
