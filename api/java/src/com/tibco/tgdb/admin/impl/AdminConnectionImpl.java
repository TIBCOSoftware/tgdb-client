
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
 *  SVN Id: $Id: AdminConnectionImpl.java 4070 2020-06-11 00:48:16Z sbangar $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.util.Collection;

import com.tibco.tgdb.admin.TGAdminConnection;
import com.tibco.tgdb.admin.TGConnectionInfo;
import com.tibco.tgdb.admin.TGIndexInfo;
import com.tibco.tgdb.admin.TGServerInfo;
import com.tibco.tgdb.admin.TGServerLogDetails;
import com.tibco.tgdb.admin.TGUserInfo;
import com.tibco.tgdb.admin.impl.CreateUserInfo;
import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.channel.TGChannelResponse;
import com.tibco.tgdb.channel.impl.BlockingChannelResponse;
import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.connection.impl.ConnectionPoolImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGException.TGExceptionType;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.impl.EdgeTypeImpl;
import com.tibco.tgdb.model.impl.MutableNodeTypeImpl;
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
			case StopServer: {
				break;
			}
			case SetLogLevel: {
				break;
			}
			default: {
				break;
			}
		}
		return result;
	}
	
	Object getCommandWithParameters (TGAdminCommand command, Object parameters) throws TGException {
		connPool.adminlock();
		
		TGChannelResponse channelResponse;
		long timeout = Long.parseLong(properties.getProperty(ConfigName.ConnectionOperationTimeoutSeconds, "-1"));
	
		long requestId  = requestIds.getAndIncrement();
		channelResponse = new BlockingChannelResponse(requestId, timeout);
		
		AdminResponse response = null;
		try {
			AdminRequest request = (AdminRequest) TGMessageFactory.getInstance().createMessage(VerbId.AdminRequest,channel.getAuthToken(),channel.getSessionId());
			request.setCommand(command);
			switch (command)
			{
				case KillConnection: {
					request.setSessionId((Long)parameters);
					break;
				}
				case SetLogLevel: {
					request.setLogLevel((TGServerLogDetails) parameters);
					break;
				}
				case CreateNodeType: {
					request.setCreateNodeTypeInfo ((MutableNodeTypeImpl) parameters);
					break;
				}
				case CreateEdgeType: {
					request.setCreateEdgeTypeInfo ((EdgeTypeImpl) parameters);
					break;
				}
				case CreateAttrDesc: {
					request.setAttrDesc ((TGAttributeDescriptor)parameters);
					break;
				}
				case CreateUser:
					request.setCreateUserInfo((CreateUserInfo)parameters);
				default: {
					break;
				}
			}
			response = (AdminResponse) this.channel.sendRequest(request, channelResponse);
		}
		finally {
			connPool.adminUnlock();
		}
		
		return response;
	}
	

	@Override
	public void stopServer() throws TGException {
		getCommandWithoutParameters (TGAdminCommand.StopServer);
		
	}

	@Override
	public void dumpServerStacktrace() throws TGException {
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

	@Override
	public void killConnection (long sessionId) throws TGException {
		getCommandWithParameters(TGAdminCommand.KillConnection, new Long (sessionId));
	}

	
	@Override
	public Collection<TGAttributeDescriptor> getAttributeDescriptors () throws TGException {
		return ((Collection<TGAttributeDescriptor>) getCommandWithoutParameters (TGAdminCommand.ShowAttrDescs));
	}

	@Override
	public Collection<TGIndexInfo> getIndices() throws TGException {
		return ((Collection<TGIndexInfo>) getCommandWithoutParameters (TGAdminCommand.ShowIndices));
	}

	@Override
	public void setServerLogLevel (TGServerLogDetails.LogComponent logComponent, TGServerLogDetails.LogLevel logLevel) throws TGException {
		TGServerLogDetails logDetails = new TGServerLogDetails(logComponent, logLevel);
		getCommandWithParameters (TGAdminCommand.SetLogLevel, logDetails);
	}

	@Override
	public void createUser(String userName, String passwd, String ...roles) throws TGException {
		CreateUserInfo userCreateInfo = new CreateUserInfo(userName,passwd,roles);
		AdminResponse adminResponse = (AdminResponse) getCommandWithParameters (TGAdminCommand.CreateUser,userCreateInfo);
        if (adminResponse.createNodeTypeStatus.getResultId() != 0)
        {
            if(adminResponse.createNodeTypeStatus.getResultId() == 369){
                TGException tgException = new TGException("Duplicated system object. Username cannot be same as any other user, nor same as a nodetype, edgetype, or index name." + adminResponse.createNodeTypeStatus.getErrorMessage(),TGExceptionType.DuplicateSystemObject,adminResponse.createNodeTypeStatus.getResultId());
                throw tgException;
            }
            else if(adminResponse.createNodeTypeStatus.getResultId() == 328){
                 //for now making it a general exception
                 TGException tgException = new TGException("Invalid role. " + adminResponse.createNodeTypeStatus.getErrorMessage(),null,adminResponse.createNodeTypeStatus.getResultId());
                 throw tgException;
            }else{
                 TGException tgException = new TGException(adminResponse.createNodeTypeStatus.getErrorMessage(),null,adminResponse.createNodeTypeStatus.getResultId());
                 throw tgException;
            }

        }
	}
}
