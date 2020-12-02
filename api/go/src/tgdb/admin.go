/*
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
 * File name: admin.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: admin.go 3516 2019-11-13 19:54:15Z nimish $
 */

package tgdb

import "time"

type TGAdminConnection interface {
	TGConnection
	// CheckpointServer allows the programmatic control to do a checkpoint on server
	CheckpointServer() TGError

	// DumpServerStackTrace allows the programmatic control to dump the stack trace on the server console
	DumpServerStackTrace() TGError

	// GetAttributeDescriptors gets the list of attribute descriptors
	GetAttributeDescriptors() ([]TGAttributeDescriptor, TGError)

	// GetConnections gets the list of all socket connections using this connection type
	GetConnections() ([]TGConnectionInfo, TGError)

	// GetIndices gets the list of all indices
	GetIndices() ([]TGIndexInfo, TGError)

	// GetInfo retrieves the server information (including the server status, listener information,
	// memory information, transaction statistics, cache statistics, database statistics)
	GetInfo() (TGServerInfo, TGError)

	// GetUsers gets the list of users
	GetUsers() ([]TGUserInfo, TGError)

	// KillConnection allows the programmatic control to stop a particular connection instance
	KillConnection(sessionId int64) TGError

	// SetServerLogLevel sets the appropriate log level on server
	SetServerLogLevel(logLevel int, logComponent int64) TGError

	// StopServer allows the programmatic-stop of the server execution
	StopServer() TGError
}

// TGCacheStatistics allows users to retrieve the Cache Statistics from server
type TGCacheStatistics interface {
	// GetDataCacheEntries returns the data-cache entries
	GetDataCacheEntries() int
	// GetDataCacheHits returns the data-cache hits
	GetDataCacheHits() int64
	// GetDataCacheMisses returns the data-cache misses
	GetDataCacheMisses() int64
	// GetDataCacheMaxEntries returns the data-cache max entries
	GetDataCacheMaxEntries() int
	// GetDataCacheMaxMemory returns the data-cache max memory
	GetDataCacheMaxMemory() int64
	// GetIndexCacheEntries returns the index-cache entries
	GetIndexCacheEntries() int
	// GetIndexCacheHits returns the index-cache hits
	GetIndexCacheHits() int64
	// GetIndexCacheMisses returns the index-cache misses
	GetIndexCacheMisses() int64
	// GetIndexCacheMaxMemory returns the index-cache max memory
	GetIndexCacheMaxMemory() int64
	// GetIndexCacheMaxEntries returns the index-cache max entries
	GetIndexCacheMaxEntries() int
}

// TGConnectionInfo allows users to retrieve the individual Connection Information from server
type TGConnectionInfo interface {
	// GetClientID returns a client ID of listener
	GetClientID() string
	// GetCreatedTimeInSeconds returns a time when the listener was created
	GetCreatedTimeInSeconds() int64
	// GetListenerName returns a name of a particular listener
	GetListenerName() string
	// GetRemoteAddress returns a remote address of listener
	GetRemoteAddress() string
	// GetSessionID returns a session ID of listener
	GetSessionID() int64
	// GetUserName returns a user-name associated with listener
	GetUserName() string
}

// TGDatabaseStatistics allows users to retrieve the database statistics from server
type TGDatabaseStatistics interface {
	// GetBlockSize returns the block size
	GetBlockSize() int
	// GetDataBlockSize returns the block size of data
	GetDataBlockSize() int
	// GetDataFree returns the free data size
	GetDataFree() int64
	// GetDataSize returns data size
	GetDataSize() int64
	// GetDataUsed returns the size of data used
	GetDataUsed() int64
	// GetDbSize returns the size of database
	GetDbSize() int64
	// GetIndexFree returns the free index size
	GetIndexFree() int64
	// GetIndexSize returns the index size
	GetIndexSize() int64
	// GetIndexUsed returns the size of index used
	GetIndexUsed() int64
	// GetNumDataSegments returns the number of data segments
	GetNumDataSegments() int
	// GetNumIndexSegments returns the number of index segments
	GetNumIndexSegments() int
}

// TGIndexInfo users to retrieve the index information from server
type TGIndexInfo interface {
	// GetAttributes returns a collection of attribute names
	GetAttributeNames() []string
	// GetName returns the index name
	GetName() string
	// GetNumEntries returns the number of entries for the index
	GetNumEntries() int64
	// GetType returns the index type
	GetType() byte
	// GetStatus returns the status of the index
	GetStatus() string
	// GetSystemId returns the system ID
	GetSystemId() int
	// GetNodeTypes returns a collection of node types
	GetNodeTypes() []string
	// IsUnique returns the information whether the index is unique
	IsUnique() bool
}

// TGMemoryInfo allows users to retrieve the memory information from server
type TGMemoryInfo interface {
	// GetFreeMemory returns the free memory size from server
	GetFreeMemory() int64
	// GetMaxMemory returns the max memory size from server
	GetMaxMemory() int64
	// GetSharedMemoryFileLocation returns the shared memory file location
	GetSharedMemoryFileLocation() string
	// GetUsedMemory returns the used memory size from server
	GetUsedMemory() int64
}

// TGNetListenerInfo allows users to retrieve the Net-Listener information from server
type TGNetListenerInfo interface {
	// GetCurrentConnections returns the count of current connections
	GetCurrentConnections() int
	// GetMaxConnections returns the count of max connections
	GetMaxConnections() int
	// GetListenerName returns the listener name
	GetListenerName() string
	// GetPortNumber returns the port detail of this listener
	GetPortNumber() string
}


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


// ======= Link State Types =======
type MemType int

const (
	MemoryProcess MemType = iota
	MemoryShared
)

// TGServerMemoryInfo allows users to retrieve the Server-Process-Memory or Server-Shared-Memory Information
type TGServerMemoryInfo interface {
	// GetMemoryInfo returns the memory info for the specified type
	GetServerMemoryInfo(memType MemType) TGMemoryInfo
}


// ======= Link State Types =======
type ServerStates int

const (
	ServerStateCreated ServerStates = iota
	ServerStateInitialized
	ServerStateStarted
	ServerStateSuspended
	ServerStateInterrupted
	ServerStateRequestStop
	ServerStateStopped
	ServerStateShutDown
)

// TGServerStatus allows users to retrieve the status of server
type TGServerStatus interface {
	// GetName returns the name of the server instance
	GetName() string
	// GetProcessId returns the process ID of server
	GetProcessId() string
	// GetServerStatus returns the state information of server
	GetServerStatus() ServerStates
	// GetUptime returns the uptime information of server
	GetUptime() time.Duration
	// GetVersion returns the server version information
	//GetVersion() TGServerVersion
}


// TGTransactionStatistics allows users to retrieve the transaction statistics from server
type TGTransactionStatistics interface {
	// GetAverageProcessingTime returns the average processing time for the transactions
	GetAverageProcessingTime() float64
	// GetPendingTransactionsCount returns the pending transactions count
	GetPendingTransactionsCount() int64
	// GetTransactionLoggerQueueDepth returns the queue depth of transactionLogger
	GetTransactionLoggerQueueDepth() int
	// GetTransactionProcessorsCount returns the transaction processors count
	GetTransactionProcessorsCount() int64
	// GetTransactionProcessedCount returns the processed transaction count
	GetTransactionProcessedCount() int64
	// GetTransactionSuccessfulCount returns the successful transactions count
	GetTransactionSuccessfulCount() int64
}



// TGUserInfo allows users to retrieve the memory information from server
type TGUserInfo interface {
	// GetName returns the user name
	GetName() string
	// GetSystemId returns the system ID for this user
	GetSystemId() int
	// GetType returns the user type
	GetType() byte
}
