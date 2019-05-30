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
 *  SVN Id: $Id: AdminRequest.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.io.IOException;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.pdu.impl.AbstractProtocolMessage;

public class AdminRequest extends /*AuthenticatedMessage*/ AbstractProtocolMessage {
	
	protected TGAdminCommand command;
	
	protected long sessionId;
	
	TGServerLogDetails logDetails;

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
	
	/*
	public void setAttrDesc(TGAttributeDescriptor _attrDesc) {
		this.attrDesc = _attrDesc;
	}
	*/
	
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
		    //case SHOW_TYPES:
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
		    /*
		    case CREATE_ATTRDESC:
		    {
		    	int command = this.command.ordinal() + 1;
		    	
		    	os.writeInt(dataLength);
			    os.writeInt(checksum);
			    os.writeInt(command);
			    
			    os.writeBoolean(attrDesc.isArray());
			    os.writeUTF(attrDesc.getName());
			    os.writeInt(attrDesc.getType().ordinal());
			    
			    if (attrDesc.getType() == TGAttributeType.Number)
			    {
			    	os.writeShort(attrDesc.getPrecision());
			    	os.writeShort(attrDesc.getScale());
			    }
			    
		    	break;
		    }
		    */
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


}
