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
 *  File name : AdminHelper.java
 *  Created on: 03/29/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: AdminHelper.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.io.IOException;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Collection;

import com.tibco.tgdb.TGVersion;
import com.tibco.tgdb.admin.TGCacheStatistics;
import com.tibco.tgdb.admin.TGConnectionInfo;
import com.tibco.tgdb.admin.TGIndexInfo;
import com.tibco.tgdb.admin.TGNetListenerInfo;
import com.tibco.tgdb.admin.TGServerMemoryInfo;
import com.tibco.tgdb.admin.TGServerStatus;
import com.tibco.tgdb.admin.TGTransactionStatistics;
import com.tibco.tgdb.admin.TGUserInfo;
import com.tibco.tgdb.admin.TGServerStatus.ServerStates;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.TGSystemObject;
import com.tibco.tgdb.model.impl.AttributeDescriptorImpl;
import com.tibco.tgdb.pdu.TGInputStream;

class AdminHelper {
	
	
	private static TGServerStatus convertFromStreamToAdminServerInfo (TGInputStream is) throws IOException {
		String name = is.readUTF();
		ServerStates status = ServerStates.values()[is.readByte()];
		String processId = "" + is.readInt();
		long duration = is.readLong();
		long versionOnWire = is.readLong();
		
		//ServerVersionInfo serverVersionInfo = new ServerVersionInfo(versionOnWire);
		TGVersion serverVersionInfo = TGVersion.getInstanceFromLong(versionOnWire);
		Duration durationObject = Duration.ofMillis((System.currentTimeMillis() - (duration/1000000)));
		
		ServerStatusImpl serverInfo = new ServerStatusImpl (name, serverVersionInfo, status, processId, durationObject);
		return serverInfo;
	}

	private static ServerMemoryInfoImpl convertFromStreamToAdminGetMemoryInfo(TGInputStream is) throws IOException {
		long freeMemoryProcess = is.readLong();
		int memUsagePercentProcess = is.readInt();
		long maxMemoryProcess = is.readLong();
		long usedMemoryProcess = maxMemoryProcess - freeMemoryProcess;
		
		String sharedMemoryFileLocation = is.readUTF();
		
		long freeMemoryShared = is.readLong();
		int memUsagePercentShared = is.readInt();
		long maxMemoryShared = is.readLong();
		long usedMemoryShared = maxMemoryShared - freeMemoryShared;
		
		MemoryInfoImpl processMemory = new MemoryInfoImpl(usedMemoryProcess, freeMemoryProcess, maxMemoryProcess, null);
		MemoryInfoImpl sharedMemory = new MemoryInfoImpl(usedMemoryShared, freeMemoryShared, maxMemoryShared, sharedMemoryFileLocation);
		
		ServerMemoryInfoImpl memoryInfo = new ServerMemoryInfoImpl(processMemory, sharedMemory);
		return memoryInfo;
	}

	private static Collection<TGNetListenerInfo> convertFromStreamToAdminNetListenersInfo(TGInputStream is) throws IOException {
		long no_Of_listeners = is.readLong();
		ArrayList<TGNetListenerInfo> listenersInfo = new ArrayList<TGNetListenerInfo>();
		
		for (int i = 0; i < no_Of_listeners; ++i)
		{
			String listenerName = "";
			int bufferLength = is.readInt();
			for (int ii = 0; ii < bufferLength; ++ii)
			{
				byte[] bytes = new byte[1];
				bytes[0] = is.readByte();
				listenerName = listenerName + new String(bytes);
			}
			
			int currentConnections = is.readInt();
			int maxConnections = is.readInt();
			String portNumber = "";
			int bufferLength4Port = is.readInt();
			for (int jj = 0; jj < bufferLength4Port; ++jj)
			{
				byte[] bytes = new byte[1];
				bytes[0] = is.readByte();
				portNumber = portNumber + new String(bytes);
			}
			
			
			NetListenerInfoImpl listenerInfo = new NetListenerInfoImpl(listenerName, currentConnections, maxConnections, portNumber);
			listenersInfo.add(listenerInfo);
		}
		
		/*
		AdminNetListenersInfoImpl netListenersInfo = new AdminNetListenersInfoImpl(listenersInfo);
		return netListenersInfo;
		*/
		return listenersInfo;
	}

	private static TransactionStatisticsImpl convertFromStreamToAdminTransactionsInfo(TGInputStream is) throws IOException {
		long transactionProcessorsCount = is.readShort();
		long transactionProcessedCount = is.readLong();
		long transactionSuccessfulCount = is.readLong();
		double averageProcessingTime = is.readDouble();
		long pendingTransactionsCount = is.readLong();
		int transactionLoggerQueueDepth = is.readInt();
		//
		// Fill the trasactionsInfo
		//
		TransactionStatisticsImpl transactionsInfo = new TransactionStatisticsImpl(transactionProcessorsCount, transactionProcessedCount, transactionSuccessfulCount, averageProcessingTime, pendingTransactionsCount, transactionLoggerQueueDepth);
		return transactionsInfo;
	}

	private static CacheStatisticsImpl convertFromStreamToAdminCacheInfo(TGInputStream is) throws IOException {
		int dataCacheMaxEntries = is.readInt();
		int dataCacheEntries = is.readInt();
		long dataCacheHits = is.readLong();
		long dataCacheMisses = is.readLong();
		long dataCacheMaxMemory = is.readLong();
		int indexCacheMaxEntries = is.readInt();
		int indexCacheEntries = is.readInt();
		long indexCacheHits = is.readLong();
		long indexCacheMisses = is.readLong();
		long indexCacheMaxMemory = is.readLong();
		
		CacheStatisticsImpl cacheInfoImpl = new CacheStatisticsImpl(dataCacheMaxEntries, dataCacheEntries, dataCacheHits, dataCacheMisses, dataCacheMaxMemory, indexCacheMaxEntries, indexCacheEntries, indexCacheHits, indexCacheMisses, indexCacheMaxMemory);
		return cacheInfoImpl;
	}

	private static DatabaseStatisticsImpl convertFromStreamToAdminDatabaseInfo(TGInputStream is) throws IOException {
		long dbSize = is.readLong();
		int numDataSegments = is.readInt();
		long dataSize = is.readLong();
		long dataUsed = is.readLong();
		long dataFree = is.readLong();
		int dataBlockSize = is.readInt();
		
		int numIndexSegments = is.readInt();
		long indexSize = is.readLong();
		long indexUsed = is.readLong();
		long indexFree = is.readLong();
		int blockSize = is.readInt();
					
		DatabaseStatisticsImpl databaseInfoImpl = new DatabaseStatisticsImpl(dbSize, numDataSegments, dataSize, dataUsed, dataFree, dataBlockSize, numIndexSegments, indexSize, indexUsed, indexFree, blockSize);			
		
		
		return databaseInfoImpl;

	}

	public static ServerInfoImpl convertFromStreamToAdminCommandInfoResult(TGInputStream is) throws IOException {
		//
		// Fill TGAdminServerInfo
		//
		TGServerStatus serverInfo = AdminHelper.convertFromStreamToAdminServerInfo(is);
		
		//
		// Fill TGAdminGetMemoryInfo
		//
		TGServerMemoryInfo memoryInfo = AdminHelper.convertFromStreamToAdminGetMemoryInfo(is);
		
		//
		// Fill TGAdminNetListenersInfo
		//
		//TGAdminNetListenersInfo netListenersInfo = TGAdminStreamHelper.convertFromStreamToAdminNetListenersInfo(is);
		Collection<TGNetListenerInfo> netListenersInfo = AdminHelper.convertFromStreamToAdminNetListenersInfo(is);
		
		//
		// Fill TransactionsInfo
		//
		TGTransactionStatistics transactionsInfo = AdminHelper.convertFromStreamToAdminTransactionsInfo(is);
		
		//
		// Fill CacheInfo
		//
		TGCacheStatistics cacheInfoImpl = AdminHelper.convertFromStreamToAdminCacheInfo(is);
		
		//
		// Fill database info
		//
		DatabaseStatisticsImpl databaseInfoImpl = AdminHelper.convertFromStreamToAdminDatabaseInfo(is);
		
		ServerInfoImpl adminCommandInfoResult = new ServerInfoImpl (serverInfo, netListenersInfo, memoryInfo, transactionsInfo, cacheInfoImpl, databaseInfoImpl);
		return adminCommandInfoResult;

	}

	public static Collection<TGUserInfo> convertFromStreamToAdminCommandShowUsers(TGInputStream is) throws IOException {
	//public static AdminShowUsersImpl convertFromStreamToAdminCommandShowUsers(TGInputStream is) throws IOException {
		
		//System.out.println("Show Users Command");
		int count = is.readInt();
		
		ArrayList<TGUserInfo> listOfUsers = new ArrayList<TGUserInfo>(); 
		for (int i = 0; i < count; ++i)
		{
			byte type = is.readByte();
			int id = is.readInt();
			String name = is.readUTF();
			
			if (type == TGSystemObject.TGSystemType.AttributeDescriptor.type())
			{
				//TODO
			}
			else if (type == TGSystemObject.TGSystemType.NodeType.type())
			{
				//TODO
			}
			else if (type == TGSystemObject.TGSystemType.EdgeType.type())
			{
				//TODO
			}
			else if (type == TGSystemObject.TGSystemType.Principle.type())
			{
				int bufferLen = is.readInt();
				byte[] buffer = new byte[bufferLen];
				is.read(buffer);
				String principleRole = is.readUTF();
				//System.out.println("Buffer Value = " + new String(buffer) + " PrincipleRole = " + principleRole);
			}
			else if (type == TGSystemObject.TGSystemType.Index.type())
			{
				//TODO
			}
			UserInfoImpl user = new UserInfoImpl (type, id, name); 
			listOfUsers.add(user);
		}
		
		//return new AdminShowUsersImpl(listOfUsers);
		return listOfUsers;

	}

	public static Collection<TGConnectionInfo> convertFromStreamToAdminCommandShowConnections(TGInputStream is) throws IOException {
		
		long countOfConnections = is.readLong();
		
		ArrayList<TGConnectionInfo> listOfConnections = new ArrayList<TGConnectionInfo>(); 
		
		for (int i = 0; i < countOfConnections; ++i)
		{
			String listnerName = is.readUTF();
			String clientID = is.readUTF();
			long sessionID = is.readLong();
			String userName = is.readUTF();
			String remoteAddress = is.readUTF();
			long createdTimeInSeconds = (System.currentTimeMillis() / 1000) - (is.readLong()/1000000000);
			
			TGConnectionInfo connection = new ConnectionInfoImpl (listnerName, clientID, sessionID, userName, remoteAddress, createdTimeInSeconds); 
			listOfConnections.add(connection);
		}
		
		return listOfConnections;
	}

	public static Collection<TGAttributeDescriptor> convertFromStreamToAdminCommandShowAttrDescs(TGInputStream is) throws IOException {
		ArrayList<TGAttributeDescriptor> attrDescs = new ArrayList<TGAttributeDescriptor>();
		
		int countOfAttrDescs = is.readInt();

		for (int i = 0; i < countOfAttrDescs; ++i)
		{
			byte type = is.readByte();
			int sysid = is.readInt();
			
			String name = is.readUTF();
			
			byte attrType = is.readByte();
			boolean isArray = is.readBoolean();
			boolean isEncrypted = is.readBoolean();
			
		    AttributeDescriptorImpl desc = new AttributeDescriptorImpl(name, TGAttributeType.fromTypeId(attrType), isArray);
		    desc.setEncrypted(isEncrypted);
		    
		    desc.setAttributeId(sysid);

			
		    if (attrType == TGAttributeType.Number.ordinal())
			{
				short precision = is.readShort();
				short scale = is.readShort();
				desc.setPrecision(precision);
			    desc.setScale(scale);
			}
		    
		    attrDescs.add(desc);
		}
		
		return attrDescs;
	}

	/*
	public static Collection<TGEntityType> convertFromStreamToAdminCommandShowTypes(TGInputStream is) throws IOException {
		ArrayList<TGEntityType> attrDescs = new ArrayList<TGEntityType>();
		
		int countOfTypes = is.readInt();
		for (int i = 0; i < countOfTypes; ++i)
		{
			byte type = is.readByte();
			int sysid = is.readInt();
			String name = is.readUTF();
			
			if (type == 1) // NodeType
			{
				int pageSize = is.readInt();
				
				short attributeCount = is.readShort();
				for (int j=0; j < attributeCount;++j)
				{
					String attrName = is.readUTF();
					System.out.println("AttributeName = " + attrName);
				}
				
				short countKeyAttributes = is.readShort();
				for (int k=0; k < countKeyAttributes; ++k)
				{
					String attrKeyName = is.readUTF();
					System.out.println("KeyAttributeName = " + attrKeyName);
				}
				
				////
				short countIndexIDs = is.readShort();
				for (int l=0; l < countIndexIDs; ++l)
				{
					int indexId = is.readInt();
				}
				
				long nCountEntries = is.readLong();
				int ii = 0;
				////
				
			}
			else if (type == 2) // EdgeType
			{
				
			}
		}
		
		return attrDescs;
	}
	*/

	public static Collection<TGIndexInfo> convertFromStreamToAdminCommandShowIndices(TGInputStream is) throws IOException {
		
		ArrayList<TGIndexInfo> listOfIndices = new ArrayList<TGIndexInfo>();
		
		int count = is.readInt();
		for (int i = 0; i < count; ++i)
		{
			byte type = is.readByte();
			int sysid = is.readInt();
			String name = is.readUTF();
			
			//System.out.println(name);
			
			boolean isUnique = is.readBoolean();
			int attrCount = is.readInt();

			ArrayList<String> attributes = new ArrayList<String>();
			for (int j = 0; j < attrCount; ++j)
			{
				String name1 = is.readUTF();
				//System.out.println(name1);
				attributes.add(name1);
			}
			
			int nodeTypesCount = is.readInt();
			
			ArrayList<String> nodeTypes = new ArrayList<String>();			
			for (int k = 0; k < nodeTypesCount; ++k)
			{
				String name2 = is.readUTF();
				//System.out.println(name2);
				nodeTypes.add(name2);
			}
			
			int blockSize = is.readInt();
			//System.out.println("BlockSize = " + blockSize);
			
			long numEntries = is.readLong();
			
			String status = new String (is.readBytes());
			
			IndexInfoImpl showIndexInfo = new IndexInfoImpl(sysid, type, name, isUnique, attributes, nodeTypes, numEntries, status);
			listOfIndices.add(showIndexInfo);
		}
		return listOfIndices;
	} 

}
