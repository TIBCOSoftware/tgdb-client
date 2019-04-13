package admin

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGServerInfo.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// TGServerInfo allows users to retrieve the Server Information; this includes
// the server status, collection of net-listener objects, information on server-memory,
// information on transaction-statistics, cache-statistics, and database information
type TGServerInfo interface {
	// GetCacheInfo returns cache statistics information from server
	GetCacheInfo() TGCacheStatistics
	// GetDatabaseInfo returns database statistics information from server
	GetDatabaseInfo() TGDatabaseStatistics
	// GetMemoryInfo returns object corresponding to specific memory type
	GetMemoryInfo(memType MemType) TGMemoryInfo
	// GetNetListenersInfo returns a collection of information on NetListeners
	GetNetListenersInfo() []TGNetListenerInfo
	// GetServerStatus returns the information on Server Status including name, version etc.
	GetServerStatus() TGServerStatus
	// GetTransactionsInfo returns transaction statistics from server including processed transaction count, successful transaction count, average processing time etc.
	GetTransactionsInfo() TGTransactionStatistics
}
