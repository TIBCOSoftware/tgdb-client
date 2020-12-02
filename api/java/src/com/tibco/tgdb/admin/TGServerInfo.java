
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
 *  File name :TGServerInfo.java
 *  Created on: 03/28/2019
 *  Created by: nimish

 *  <p>This interface allows users to retrieve the Server Information; this includes 
 *  the server status, collection of net-listener objects, information on server-memory, 
 *  information on transaction-statistics, cache-statistics, and database information  
 * 
 *  SVN Id: $Id: TGServerInfo.java 3120 2019-04-25 21:21:48Z nimish $ 
 */

package com.tibco.tgdb.admin;

import java.util.Collection;

import com.tibco.tgdb.admin.TGServerMemoryInfo.MEMORY_TYPE;



public interface TGServerInfo {

	/**
	 * Retrieves the information on Server Status
	 * @return {@link TGServerStatus} gives status of the server including name, version etc.
	 */
	public TGServerStatus getServerStatus ();
	

	
	/**
	 * Retrieves a collection of information on NetListeners
	 * @return a collection of {@link TGNetListenerInfo}
	 */
	public Collection<TGNetListenerInfo> getNetListenersInfo ();
	
	
	
	/**
	 * Retrieves {@link TGMemoryInfo} object corresponding to specific {@link MEMORY_TYPE}
	 * @param type specific {@link MEMORY_TYPE} for which the memory information is needed
	 * @return {@link TGMemoryInfo} object for a specific {@link MEMORY_TYPE}
	 */
	public TGMemoryInfo getMemoryInfo (MEMORY_TYPE type);
	
	
	/**
	 * Retrieves transaction statistics from server
	 * @return {@link TGTransactionStatistics} object that gives various pieces of information including processed transaction count, successful transaction count, average processing time etc.
	 */
	public TGTransactionStatistics getTransactionsInfo ();
	
	
	
	/**
	 * Retrieves cache statistics information from server
	 * @return {@link TGCacheStatistics} object
	 */
	public TGCacheStatistics getCacheInfo ();

	
	/**
	 * Retrieves database statistics information from server
	 * @return {@link TGDatabaseStatistics} object
	 */
	public TGDatabaseStatistics getDatabaseInfo ();
	
}
