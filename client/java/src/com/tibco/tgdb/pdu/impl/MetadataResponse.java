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
 * File name : MetadataResponse.${EXT}
 * Created on: 2/4/15
 * Created by: chung 
 * <p/>
 * SVN Id: $Id: MetadataResponse.java 583 2016-03-15 02:02:39Z vchung $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGEntityId;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.model.TGSystemObject.TGSystemType;
import com.tibco.tgdb.model.impl.AttributeDescriptorImpl;
import com.tibco.tgdb.model.impl.ByteArrayEntityId;
import com.tibco.tgdb.model.impl.EdgeTypeImpl;
import com.tibco.tgdb.model.impl.NodeTypeImpl;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class MetadataResponse extends AbstractProtocolMessage {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private ArrayList<TGAttributeDescriptor> attrDescList = new ArrayList();
    private ArrayList<TGNodeType> nodeTypeList = new ArrayList(); 
    private ArrayList<TGEdgeType> edgeTypeList = new ArrayList(); 

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
    	gLogger.log(TGLevel.Debug, "Entering metadata response readPayload");
    	if (is.available() == 0) {
    		gLogger.log(TGLevel.Debug, "Entering metadata response has no data");
    		return;
    	}
    	int count  = is.readInt();
    	while (count > 0) {
    		int sysType = is.readByte();
    		int typeCount = is.readInt();
    		if (sysType == TGSystemType.AttributeDescriptor.type()) {
    			for (int i=0; i<typeCount; i++) {
    				TGAttributeDescriptor attrDesc = new AttributeDescriptorImpl("temp", TGAttributeType.String);
    				attrDesc.readExternal(is);
    				attrDescList.add(attrDesc);
    			}
    		} else if (sysType == TGSystemType.NodeType.type()) {
    			for (int i=0; i<typeCount; i++) {
    				TGNodeType nodeType = new NodeTypeImpl("temp", null);
    				nodeType.readExternal(is);
    				String name = nodeType.getName();
    				if (name.startsWith("@") || name.startsWith("$")) {
    					continue;
    				}
    				nodeTypeList.add(nodeType);
    			}
    		} else if (sysType == TGSystemType.EdgeType.type()) {
    			for (int i=0; i<typeCount; i++) {
    				TGEdgeType edgeType = new EdgeTypeImpl("temp", TGEdge.DirectionType.BiDirectional, null);
    				edgeType.readExternal(is);
    				edgeTypeList.add(edgeType);
    			}
    		} else {
    			gLogger.log(TGLevel.Warning, "Invalid meta data desc received %d\n", sysType);
    			//FIXME: Need to throw exception
    		}
    		count -= typeCount;
    	}
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.MetadataResponse;
    }

    public List<TGAttributeDescriptor> getAttrDescList() {
        return attrDescList;
    }

    public List<TGNodeType> getNodeTypeList() {
        return nodeTypeList;
    }

    public List<TGEdgeType> getEdgeTypeList() {
        return edgeTypeList;
    }
}
