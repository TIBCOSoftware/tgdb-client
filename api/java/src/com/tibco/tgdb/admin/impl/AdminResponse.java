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
 *  File name : AdminResponse.java
 *  Created on: 03/29/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: AdminResponse.java 4070 2020-06-11 00:48:16Z sbangar $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.io.IOException;
import java.util.Collection;

import com.tibco.tgdb.admin.TGConnectionInfo;
import com.tibco.tgdb.admin.TGIndexInfo;
import com.tibco.tgdb.admin.TGServerInfo;
import com.tibco.tgdb.admin.TGUserInfo;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.pdu.impl.AbstractProtocolMessage;
import java.io.EOFException;
import com.tibco.tgdb.utils.TGConstants;

public class AdminResponse extends AbstractProtocolMessage {
	
	protected ServerInfoImpl adminCommandInfoResult = null;
	protected Collection<TGUserInfo> users = null;
	protected Collection<TGConnectionInfo> connections = null;
	
	protected Collection<TGAttributeDescriptor> attrDescs = null;
	protected Collection<TGIndexInfo> indices = null;
	
	//SSB : TODO: Can we rename StatusForCreateEntityType to make it generic for all the commands returning some status and error message?
	protected StatusForCreateEntityType createNodeTypeStatus = new StatusForCreateEntityType();
	
	public boolean isUpdateable() {
		return false;
	}

	public VerbId getVerbId() {
		return VerbId.AdminResponse;
	}
	
    protected void readPayload(TGInputStream is) throws TGException, IOException {
   	
    	int resultId = is.readInt();
    	int adminCommand = is.readInt();
    	--adminCommand;
    	TGAdminCommand command = TGAdminCommand.values[adminCommand];
    	
    	switch (command)
    	{
    		case ShowInfo: {
    			adminCommandInfoResult = AdminHelper.convertFromStreamToAdminCommandInfoResult(is);
    			break;
    		}
    		case ShowUsers: {
    			users = AdminHelper.convertFromStreamToAdminCommandShowUsers(is);
    			break;
    		}
    		case ShowConnections: {
    			connections = AdminHelper.convertFromStreamToAdminCommandShowConnections(is);
    			break;
    		}
    		case ShowAttrDescs: {
    			attrDescs = AdminHelper.convertFromStreamToAdminCommandShowAttrDescs(is);
    			break;
    		}
    		case ShowIndices: {
    			indices = AdminHelper.convertFromStreamToAdminCommandShowIndices(is);
    			break;
    		}
    		case CreateNodeType:
    		case CreateEdgeType:
    		case CreateAttrDesc:
    		case CreateUser: {
        		if (resultId == 0)
        		{
        			createNodeTypeStatus.setResultId(resultId);
        		}
        		else {
        			createNodeTypeStatus.setResultId(resultId);
        			String errorMessage = TGConstants.EmptyString;
        			try{
        			    errorMessage = is.readUTF();
        			}
        			catch(EOFException eof){
        			    //if server did not send en error message
        			    errorMessage = "Error Processing create user request.";
        			}
                    createNodeTypeStatus.setErrorMessage(errorMessage);
        		}
    			break;
    		}
    		default: {
    			break;
    		}
    	}
    }

	protected void writePayload(TGOutputStream os) throws TGException {
	}
	
	public TGServerInfo getAdminCommandInfoResult () {
		return this.adminCommandInfoResult;
	}
	
	public Collection<TGUserInfo> getAdminCommandUsersResult () {
		return this.users;
	}
	
	public Collection<TGConnectionInfo> getAdminCommandConnectionsResult () {
		return this.connections;
	}

	public Collection<TGAttributeDescriptor> getAdminCommandAttrDescsResult() {
		return this.attrDescs;
	}

	public Collection<TGIndexInfo> getIndices () {
		return this.indices;
	}
}