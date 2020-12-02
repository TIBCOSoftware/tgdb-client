/**
 * Copyright 2019 TIBCO Software Inc.
 * All rights reserved.
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
 * <p/>
 * File name: TGAdminConnection.java
 * Created on: 2019-03-01
 * Created by: nimish
 * <p/>
 * SVN Id: $Id: TGAdminConnection.java 3157 2019-04-26 20:28:37Z kattaylo $
 */

package com.tibco.tgdb.admin;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;

import java.util.Collection;

public interface TGAdminConnection extends TGConnection {

	/**
	 * Retrieves the server information (including the server status, 
	 * listener information, memory information, transaction statistics, 
	 * cache statistics, database statistics)
	 * 
	 * @return {@link TGServerInfo}
	 * @throws TGException if any error occurs while retrieving the server information
	 */
	public TGServerInfo getInfo() throws TGException;
	


	/**
	 * Retrieves the collection of {@link TGUserInfo} objects from server
	 * 
	 * @return A Collection of {@link TGUserInfo} objects
	 * @throws TGException if any error occurs while retrieving user-information
	 * Information from server  
	 */
	public Collection<TGUserInfo> getUsers() throws TGException;



	/**
	 * Retrieves the collection of {@link TGConnectionInfo} objects from server
	 *
	 * @return A Collection of {@link TGConnectionInfo} objects from server
	 * @throws TGException if any error occurs while getting the connection information
	 */
	public Collection<TGConnectionInfo> getConnections() throws TGException;



	/**
	 * Allows the programmatic-stop of the server execution
	 *
	 * @throws TGException if any error occurs while stopping the server
	 */
	public void stopServer() throws TGException;



	/**
	 * Allows the programmatic control to dump the stack trace on the server console
	 *
	 * @throws TGException if any error occurs while communicating with the server
	 */
	public void dumpServerStacktrace() throws TGException;



	/**
	 * Allows the programmatic control to take a checkpoint on server
	 *
	 * @throws TGException if any error occurs while communicating with the server
	 */
	public void checkpointServer() throws TGException;



	/**
	 * Allows the programmatic control to stop a particular connection instance
	 *
	 * @param sessionId  session ID for a particular connection
	 * @throws TGException if any error occurs while performing the operation on server
	 */
	public void killConnection(long sessionId) throws TGException;




	// TODO:
	// public void createAttributeDescriptor (TGAttributeDescriptor attrDesc) throws TGException;



	/**
	 * Retrieves a collection of {@link TGAttributeDescriptor} from the server
	 * @return A Collection of {@link TGAttributeDescriptor}
	 * @throws TGException if any error occurs while retrieving information from the server
	 */
	public Collection<TGAttributeDescriptor> getAttributeDescriptors() throws TGException;




	// TODO:
	//getNodeTypes()
	//getNodeType(String name)
	//getEdgeTypes()
	//getEdgeType(String name)



	/**
	 * Retrieves a collection of {@link TGIndexInfo} from the server
	 *
	 * @return A Collection of {@link TGIndexInfo}
	 * @throws TGException if any error occurs while retrieving information from the server
	 */
	public Collection<TGIndexInfo> getIndices() throws TGException;


	//
	// TODO:
	// Create AttributeDescriptor, NodeType, EdgeType
	//


	//
	// TODO:
	// Drop of NodeType, EdgeType, Indices (post 2.0.1)
	//



	/**
	 * Sets the appropriate log level on server
	 * @param logComponent Specific {@link TGServerLogDetails.LogComponent} to be set
	 * @param logLevel Specific {@link TGServerLogDetails.LogLevel} to be set
	 * @throws TGException if any error occurs while setting the log-level on the server
	 */
	public void setServerLogLevel(TGServerLogDetails.LogComponent logComponent, TGServerLogDetails.LogLevel logLevel) throws TGException;

	//
	// TODO:
	// may not be available on server yet
	// getServerLogLevel()
	
	/**
	 * Creates a User on Server and assigns it the specified roles if any
	 * @param userName for the user to be created
	 * @param passwd password for the user to be created
	 * @param roles optional roles to be assigned to the user
	 * @throws TGException if any error occurs while creating the user
	 */
	public void createUser(String userName, String passwd, String ...roles) throws TGException;
}
