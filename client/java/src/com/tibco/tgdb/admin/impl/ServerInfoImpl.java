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
 *  File name : ServerInfoImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: ServerInfoImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.util.Collection;

import com.tibco.tgdb.admin.TGCacheStatistics;
import com.tibco.tgdb.admin.TGDatabaseStatistics;
import com.tibco.tgdb.admin.TGMemoryInfo;
import com.tibco.tgdb.admin.TGNetListenerInfo;
import com.tibco.tgdb.admin.TGServerInfo;
import com.tibco.tgdb.admin.TGServerMemoryInfo;
import com.tibco.tgdb.admin.TGServerStatus;
import com.tibco.tgdb.admin.TGTransactionStatistics;
import com.tibco.tgdb.admin.TGServerMemoryInfo.MEMORY_TYPE;

public class ServerInfoImpl implements TGServerInfo {
	
	

	@Override
	public String toString() {
		return "ServerInfoImpl [serverInfo=" + serverInfo + ", netListenersInfo=" + netListenersInfo + ", memoryInfo="
				+ memoryInfo + ", transactionsInfo=" + transactionsInfo + ", cacheInfo=" + cacheInfo + ", databaseInfo="
				+ databaseInfo + "]";
	}

	protected TGServerStatus serverInfo;
	
	//protected TGAdminNetListenersInfo netListenersInfo;
	protected Collection<TGNetListenerInfo> netListenersInfo;
	
	
	protected TGServerMemoryInfo memoryInfo;
	protected TGTransactionStatistics transactionsInfo;
	protected TGCacheStatistics cacheInfo;
	protected TGDatabaseStatistics databaseInfo;

	
	public ServerInfoImpl (
		TGServerStatus _serverInfo, 
		Collection<TGNetListenerInfo> _netListenersInfo, 
		TGServerMemoryInfo _memoryInfo, 
		TGTransactionStatistics _transactionsInfo,
		TGCacheStatistics _cacheInfo,
		TGDatabaseStatistics _databaseInfo
	)
	{
		serverInfo = _serverInfo;
		netListenersInfo = _netListenersInfo; 
		memoryInfo = _memoryInfo;
		transactionsInfo = _transactionsInfo;
		cacheInfo = _cacheInfo;
		databaseInfo = _databaseInfo;
	}

	@Override
	public TGServerStatus getServerStatus() {
		return this.serverInfo;
	}

	@Override
	public Collection<TGNetListenerInfo> getNetListenersInfo () {
		return this.netListenersInfo;
	}
	
	
	@Override
	public TGMemoryInfo getMemoryInfo (MEMORY_TYPE type)
	{
		return memoryInfo.getMemoryInfo(type);
	}
	
//	@Override
//	public TGServerMemoryInfo getMemoryInfo() {
//		return this.memoryInfo;
//	}

	@Override
	public TGTransactionStatistics getTransactionsInfo() {
		return this.transactionsInfo;
	}

	@Override
	public TGCacheStatistics getCacheInfo() {
		return this.cacheInfo;
	}

	@Override
	public TGDatabaseStatistics getDatabaseInfo() {
		return this.databaseInfo;
	}

}
