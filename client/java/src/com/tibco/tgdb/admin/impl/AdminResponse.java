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
 *  SVN Id: $Id: AdminResponse.java 3122 2019-04-25 21:38:58Z nimish $
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

public class AdminResponse extends /*AuthenticatedMessage*/ AbstractProtocolMessage {
	
	protected ServerInfoImpl adminCommandInfoResult = null;
	protected Collection<TGUserInfo> users = null;
	protected Collection<TGConnectionInfo> connections = null;
	
	protected Collection<TGAttributeDescriptor> attrDescs = null;
	protected Collection<TGIndexInfo> indices = null;
	
	public boolean isUpdateable() {
		// TODO Auto-generated method stub
		return false;
	}

	public VerbId getVerbId() {
		return VerbId.AdminResponse;
	}
	
    protected void readPayload(TGInputStream is) throws TGException, IOException {
    	//System.out.println("readPayload called..");
    	//super.readPayload(is);
    	
    	int resultId = is.readInt();
    	int command = is.readInt();
    	
    	//if (command == 10)
    	if (command == (TGAdminCommand.ShowInfo.ordinal() + 1))
    	{    		    		
    		adminCommandInfoResult = AdminHelper.convertFromStreamToAdminCommandInfoResult(is);
    	}
    	else if (command == (TGAdminCommand.ShowUsers.ordinal() + 1))
    	{
    		users = AdminHelper.convertFromStreamToAdminCommandShowUsers(is);
    	}
    	else if (command == (TGAdminCommand.ShowConnections.ordinal() + 1))
    	{
    		connections = AdminHelper.convertFromStreamToAdminCommandShowConnections(is);
    	}
    	else if (command == (TGAdminCommand.ShowAttrDescs.ordinal() + 1))
    	{
    		attrDescs = AdminHelper.convertFromStreamToAdminCommandShowAttrDescs(is);
    	}
    	/*
    	else if (command == (TGAdminCommand.SHOW_TYPES.ordinal() + 1))
    	{
    		entityTypes = TGAdminStreamHelper.convertFromStreamToAdminCommandShowTypes(is);
    	}
    	*/
    	else if (command == (TGAdminCommand.ShowIndices.ordinal() + 1))
    	{
    		indices = AdminHelper.convertFromStreamToAdminCommandShowIndices(is);
    	}
    }

	protected void writePayload(TGOutputStream os) throws TGException {
		//super.writePayload(os);
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

	/*
	public Collection<TGEntityType> getAdminCommandEntityTypesResult() {
		// TODO Auto-generated method stub
		return this.entityTypes;
	}
	*/
	
	public Collection<TGIndexInfo> getIndices () {
		return this.indices;
	}

}