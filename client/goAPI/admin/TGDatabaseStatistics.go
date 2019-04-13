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
 * File name: TGDatabaseStatistics.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
