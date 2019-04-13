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
 * File name: TGCacheStatistics.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
