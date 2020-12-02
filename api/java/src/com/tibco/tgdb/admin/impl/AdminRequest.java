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
 *
 *  File name : AdminRequest.java
 *  Created on: 03/29/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: AdminRequest.java 4070 2020-06-11 00:48:16Z sbangar $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.io.IOException;
import java.util.Collection;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

import com.tibco.tgdb.admin.TGServerLogDetails;
import com.tibco.tgdb.admin.impl.CreateUserInfo;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.impl.EdgeTypeImpl;
import com.tibco.tgdb.model.impl.MutableNodeTypeImpl;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.pdu.impl.AbstractProtocolMessage;

public class AdminRequest extends AbstractProtocolMessage {
	
	protected TGAdminCommand command;
	protected long sessionId;
	TGServerLogDetails logDetails;
	MutableNodeTypeImpl createNodeTypeInfo;
	EdgeTypeImpl edgeType;
	TGAttributeDescriptor attrDesc;	
	CreateUserInfo createUserInfo;
	Map<TGEdge.DirectionType, String> directionType2String = new HashMap <TGEdge.DirectionType, String> ();
	{
		directionType2String.put(TGEdge.DirectionType.Directed, "directed");
		directionType2String.put(TGEdge.DirectionType.UnDirected, "undirected");
		directionType2String.put(TGEdge.DirectionType.BiDirectional, "bidirected");
	}

    public AdminRequest() {
    	super();
    }

    public AdminRequest(long authToken, long sessionId) {
    	super(authToken, sessionId);
    }

	public void setCommand (TGAdminCommand _command)
	{
		this.command = _command;
	}

	public void setSessionId (long _sessionId)
	{
		this.sessionId = _sessionId;
	}
	
	public void setLogLevel (TGServerLogDetails _logDetails)
	{
		logDetails = _logDetails;
	}
	
	public void setCreateNodeTypeInfo (MutableNodeTypeImpl createNodeTypeInfo)
	{
		this.createNodeTypeInfo = createNodeTypeInfo;
	}
	public void setCreateUserInfo (CreateUserInfo createUserInfo)
	{
		this.createUserInfo = createUserInfo;
	}
	
	public TGAdminCommand getCommand ()
	{
		return this.command;
	}


	@Override
	protected void writePayload(TGOutputStream os) throws TGException, IOException {
		
		int dataLength = 0;
	    int checksum = 0;
	    
	    switch (this.command)
	    {
		    case ShowInfo:
		    case ShowUsers:
		    case ShowConnections:
		    case ShowAttrDescs:
		    case ShowIndices:
		    case StopServer:
		    case CheckpointServer:
		    {
			    int command = this.command.ordinal() + 1;
			    os.writeInt(dataLength);
			    os.writeInt(checksum);
			    os.writeInt(command);
			    break;
		    }
		    case KillConnection:
		    {
		    	int command = this.command.ordinal() + 1;
			    os.writeInt(dataLength);
			    os.writeInt(checksum);
			    os.writeInt(command);
			    os.writeLong(this.sessionId);
			    os.writeBoolean(true);
		    	break;
		    }
		    case SetLogLevel:
		    {
		    	int command = this.command.ordinal() + 1;
		    	os.writeInt(dataLength);
			    os.writeInt(checksum);
			    os.writeInt(command);
			    int logLevel = logDetails.getLogLevel().getLogLevel();
			    os.writeShort(logLevel);
			    long component = logDetails.getLogComponent().getLogComponent();
				os.writeLong(component);
		    	break;
		    }
		    case CreateNodeType:
		    {
		    	int command = this.command.ordinal() + 1;
		    	os.writeInt(dataLength);
			    os.writeInt(checksum);
			    os.writeInt(command);
			    os.writeUTF(this.createNodeTypeInfo.getName());
			    os.writeInt(512); // this is &request->pageSize.
			    int attrCount = 0;
			    Collection<TGAttributeDescriptor> listOfAttributes = this.createNodeTypeInfo.getAttributeDescriptors();

			    if (listOfAttributes != null)
			    {
			    	attrCount = listOfAttributes.size();
			    }
			    os.writeInt(attrCount);

			    if (attrCount != 0)
			    {
				    Iterator<TGAttributeDescriptor> itAttributes = listOfAttributes.iterator();
				    for (;itAttributes.hasNext();)
				    {
				    	TGAttributeDescriptor currentAttribute = itAttributes.next();
				    	os.writeBytes(currentAttribute.getName());
				    }
			    }

			    int keyCount = 0;
			    TGAttributeDescriptor[] pKeyAttributeDescriptors = this.createNodeTypeInfo.getPKeyAttributeDescriptors();

			    if (pKeyAttributeDescriptors != null)
			    {
			    	keyCount = pKeyAttributeDescriptors.length;
			    }
			    os.writeInt(keyCount);

			    if (keyCount != 0)
			    {
			    	for (int keyIndex = 0; keyIndex < keyCount; ++keyIndex)
			    	{
			    		TGAttributeDescriptor currentPKeyAttribute = pKeyAttributeDescriptors[keyIndex];
			    		os.writeBytes(currentPKeyAttribute.getName());
			    	}
			    }
		    	break;
		    }
		    case CreateEdgeType:
		    {
		    	int command = this.command.ordinal() + 1;

		    	os.writeInt(dataLength);
			    os.writeInt(checksum);
			    os.writeInt(command);

			    String nameOfEdgeType = this.edgeType.getName();
			    String nameOfFromNode = this.edgeType.getFromNodeType().getName();
			    String nameOfToNode = this.edgeType.getToNodeType().getName();

			    String direction = this.directionType2String.get(this.edgeType.getDirectionType());
			    int attrCount = this.edgeType.getAttributeDescriptors().size();

			    os.writeUTF(nameOfEdgeType);
			    os.writeUTF(nameOfFromNode);
			    os.writeUTF(nameOfToNode);
			    os.writeUTF(direction);

			    os.writeInt(attrCount);

			    Iterator<TGAttributeDescriptor> itAttributes = this.edgeType.getAttributeDescriptors().iterator();

			    for (; itAttributes.hasNext(); )
			    {
			    	TGAttributeDescriptor currentAttribute = itAttributes.next();
			    	os.writeChars(currentAttribute.getName());//os.writeUTF(currentAttribute.getName());
			    }

				int pkeyCount = 0;
			    os.writeInt(pkeyCount);
		    	break;
		    }
		    case CreateAttrDesc:
		    {
				int command = this.command.ordinal() + 1;
		    	
		    	os.writeInt(dataLength);
			    os.writeInt(checksum);
			    os.writeInt(command);
			    
			    os.writeBoolean(attrDesc.isArray());
			    os.writeBoolean(attrDesc.isEncrypted());
			    os.writeUTF(attrDesc.getName());
			    os.writeInt(attrDesc.getType().ordinal());
			    
			    if (attrDesc.getType() == TGAttributeType.Number)
			    {
			    	os.writeShort(attrDesc.getPrecision());
			    	os.writeShort(attrDesc.getScale());
			    }

				break;
		    }
		    case CreateUser:
		    {
		    	int command = this.command.ordinal() + 1;
		    	
		    	os.writeInt(dataLength);
			    os.writeInt(checksum);
			    os.writeInt(command);
			    
			    String userName = this.createUserInfo.getName();
			    String password = this.createUserInfo.getPasswd();
			    List<String> roleList = this.createUserInfo.getRoles();
			    
			    os.writeUTF(userName);
			    os.writeUTF(password);			    
			    os.writeInt(roleList.size());
			    for(String roleName: roleList)
			    	os.writeUTF(roleName);		
			    break;
		    }
		    default: {
		    	break;
		    }
	    }
	}
	
	@Override
	protected void readPayload(TGInputStream is) throws TGException, IOException {

	}

	@Override
	public boolean isUpdateable() {
		return true;
	}

	@Override
	public VerbId getVerbId() {
		return VerbId.AdminRequest;
	}


	public void setCreateEdgeTypeInfo(EdgeTypeImpl parameters) {
		this.edgeType = parameters;
	}

	public void setAttrDesc (TGAttributeDescriptor parameters) {
		this.attrDesc = parameters;
	}
}
