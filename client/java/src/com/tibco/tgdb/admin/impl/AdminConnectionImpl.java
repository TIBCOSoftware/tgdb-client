
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
 *  File name : AdminConnectionImpl.java
 *  Created on: 04/24/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: AdminConnectionImpl.java 3121 2019-04-25 21:36:20Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.util.Collection;

import com.tibco.tgdb.admin.TGConnectionInfo;
import com.tibco.tgdb.admin.TGIndexInfo;
import com.tibco.tgdb.admin.TGServerInfo;
import com.tibco.tgdb.admin.TGUserInfo;
import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.channel.TGChannelResponse;
import com.tibco.tgdb.channel.impl.BlockingChannelResponse;
import com.tibco.tgdb.connection.TGAdminConnection;
import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.connection.impl.ConnectionPoolImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogComponent;
import com.tibco.tgdb.log.TGLogLevel;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.pdu.TGMessageFactory;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.TGProperties;



public class AdminConnectionImpl extends ConnectionImpl implements TGAdminConnection {
	
	public AdminConnectionImpl(ConnectionPoolImpl connPool, TGChannel channel,
			TGProperties<String, String> properties) {
		super(connPool, channel, properties);
		
		
	}
	
	public TGServerInfo getInfo() throws TGException {
		return ((TGServerInfo) getCommandWithoutParameters (TGAdminCommand.ShowInfo));
	}
	
	public Collection<TGUserInfo> getUsers () throws TGException {
		return ((Collection<TGUserInfo>) getCommandWithoutParameters (TGAdminCommand.ShowUsers));
	}
	
	public Collection<TGConnectionInfo> getConnections () throws TGException {
		return ((Collection<TGConnectionInfo>) getCommandWithoutParameters (TGAdminCommand.ShowConnections));
	}
	
	private Object getCommandWithoutParameters (TGAdminCommand command) throws TGException {
		connPool.adminlock();
		
		TGChannelResponse channelResponse;
		long timeout = Long.parseLong(properties.getProperty(ConfigName.ConnectionOperationTimeoutSeconds, "-1"));
	
		long requestId  = requestIds.getAndIncrement();
		channelResponse = new BlockingChannelResponse(requestId, timeout);
		
		AdminResponse response = null;
		try {
			AdminRequest request = (AdminRequest) TGMessageFactory.getInstance().createMessage(VerbId.AdminRequest);
			request.setCommand(command);
			response = (AdminResponse) this.channel.sendRequest(request, channelResponse);
		}
		finally {
			connPool.adminUnlock();
		}
		
		Object result = null;
		switch (command)
		{
			case ShowInfo: {
				result = response.getAdminCommandInfoResult();
				break;
			}
			case ShowUsers: {
				result = response.getAdminCommandUsersResult();
				break;
			}
			case ShowConnections: {
				result = response.getAdminCommandConnectionsResult();
				break;
			}
			case ShowAttrDescs: {
				result = response.getAdminCommandAttrDescsResult();
				break;
			}
			case ShowIndices: {
				result = response.getIndices();
				break;
			}
			/*
			case SHOW_TYPES: {
				result = response.getAdminCommandEntityTypesResult();
				break;
			}
			*/
			case StopServer: {
				break;
			}
			case SetLogLevel: {
				break;
			}
		}
		return result;
	}
	
	private Object getCommandWithParameters (TGAdminCommand command, Object parameters) throws TGException {
		connPool.adminlock();
		
		TGChannelResponse channelResponse;
		long timeout = Long.parseLong(properties.getProperty(ConfigName.ConnectionOperationTimeoutSeconds, "-1"));
	
		long requestId  = requestIds.getAndIncrement();
		channelResponse = new BlockingChannelResponse(requestId, timeout);
		
		AdminResponse response = null;
		try {
			AdminRequest request = (AdminRequest) TGMessageFactory.getInstance().createMessage(VerbId.AdminRequest);
			request.setCommand(command);
			
			switch (command)
			{
				case KillConnection: {
					//request.setKillConnectionInfo((TGAdminKillConnectionInfo)parameters);
					request.setSessionId((Long)parameters);
					//request.setKillConnectionInfo(parameters);
					break;
				}
				case SetLogLevel: {
					request.setLogLevel((TGServerLogDetails) parameters);
				}
				/*
				case CREATE_ATTRDESC: {
					request.setAttrDesc ((TGAttributeDescriptor)parameters);
					break;
				}
				*/
			}
			
			response = (AdminResponse) this.channel.sendRequest(request, channelResponse);
		}
		finally {
			connPool.adminUnlock();
		}
		
		
		Object result = null;
		
		/*
		 
		//
		// The switch case may be needed for other commands later
		//
		 
		switch (command)
		{
			case SHOW_INFO: {
				result = response.getAdminCommandInfoResult();
				break;
			}
			case SHOW_USERS: {
				result = response.getAdminCommandUsersResult();
				break;
			}
			case SHOW_CONNECTIONS: {
				result = response.getAdminCommandConnectionsResult();
				break;
			}
			case STOP_SERVER: {
				break;
			}
		}
		*/
		
		return result;
	}
	

	@Override
	public void stopServer() throws TGException {
		getCommandWithoutParameters (TGAdminCommand.StopServer);
		
	}

	@Override
	public void dumpServerStacktrace() throws TGException {
		
/*		
		connPool.adminlock();
		
		TGChannelResponse channelResponse;
		long timeout = Long.parseLong(properties.getProperty(ConfigName.ConnectionOperationTimeoutSeconds, "-1"));
	
		long requestId  = requestIds.getAndIncrement();
		channelResponse = new BlockingChannelResponse(requestId, timeout);
		
		Object response = null;
		try {
			DumpStacktraceRequest request = (DumpStacktraceRequest) TGMessageFactory.getInstance().createMessage(VerbId.DumpStacktraceRequest);
			//response = this.channel.sendRequest(request, channelResponse);
			this.channel.sendMessage(request);
		}
		finally {
			connPool.adminUnlock();
		}
*/
		
		connPool.adminlock();
		
		try {
			DumpStacktraceRequest request = (DumpStacktraceRequest) TGMessageFactory.getInstance().createMessage(VerbId.DumpStacktraceRequest);
			this.channel.sendMessage(request);
		}
		finally {
			connPool.adminUnlock();
		}
		
	}

	@Override
	public void checkpointServer() throws TGException {
		Object obj = getCommandWithoutParameters (TGAdminCommand.CheckpointServer);
	}

	/*
	@Override
	public void killConnection(TGAdminKillConnectionInfo killConnectionInfo) throws TGException {
		getCommandWithParameters(TGAdminCommand.KILL_CONNECTION, killConnectionInfo);
		
	}
	*/
	
	@Override
	public void killConnection (long sessionId) throws TGException {
		getCommandWithParameters(TGAdminCommand.KillConnection, new Long (sessionId));
	}

	

	/*
	@Override
	public void createAttributeDescriptor(TGAttributeDescriptor attrDesc) throws TGException {
		getCommandWithParameters(TGAdminCommand.CREATE_ATTRDESC, attrDesc);
	}
	*/

	@Override
	public Collection<TGAttributeDescriptor> getAttributeDescriptors () throws TGException {
		return ((Collection<TGAttributeDescriptor>) getCommandWithoutParameters (TGAdminCommand.ShowAttrDescs));
	}

	/*
	@Override
	public Collection<TGEntityType> getEntityTypes() throws TGException {
		return ((Collection<TGEntityType>) getCommandWithoutParameters (TGAdminCommand.SHOW_TYPES));
	}
	*/

	@Override
	public Collection<TGIndexInfo> getIndices() throws TGException {
		return ((Collection<TGIndexInfo>) getCommandWithoutParameters (TGAdminCommand.ShowIndices));
	}

	@Override
	//public void setServerLogLevel (long logComponent, int logLevel) throws TGException {
	public void setServerLogLevel (TGLogComponent logComponent, TGLogLevel logLevel) throws TGException {
	//public void setLogLevel(int logLevel, long logComponent) throws TGException {
		
		TGServerLogDetails logDetails = new TGServerLogDetails(logComponent, logLevel);
		getCommandWithParameters (TGAdminCommand.SetLogLevel, logDetails);
		
	}

	
	

}
