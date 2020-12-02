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
 *  File name : CacheStatisticsImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: CacheStatisticsImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.admin.TGCacheStatistics;

public class CacheStatisticsImpl implements TGCacheStatistics{
	

	
	
	@Override
	public String toString() {
		return "CacheStatisticsImpl [dataCacheMaxEntries=" + dataCacheMaxEntries + ", dataCacheEntries="
				+ dataCacheEntries + ", dataCacheHits=" + dataCacheHits + ", dataCacheMisses=" + dataCacheMisses
				+ ", dataCacheMaxMemory=" + dataCacheMaxMemory + ", indexCacheMaxEntries=" + indexCacheMaxEntries
				+ ", indexCacheEntries=" + indexCacheEntries + ", indexCacheHits=" + indexCacheHits
				+ ", indexCacheMisses=" + indexCacheMisses + ", indexCacheMaxMemory=" + indexCacheMaxMemory + "]";
	}

	int dataCacheMaxEntries;
	int dataCacheEntries;
	long dataCacheHits;
	long dataCacheMisses;
	long dataCacheMaxMemory;
	int indexCacheMaxEntries;
	int indexCacheEntries;
	long indexCacheHits;
	long indexCacheMisses;
	long indexCacheMaxMemory;
	
	public CacheStatisticsImpl (
		int _dataCacheMaxEntries,
		int _dataCacheEntries,
		long _dataCacheHits,
		long _dataCacheMisses,
		long _dataCacheMaxMemory,
		int _indexCacheMaxEntries,
		int _indexCacheEntries,
		long _indexCacheHits,
		long _indexCacheMisses,
		long _indexCacheMaxMemory
	)
	{
		dataCacheMaxEntries = _dataCacheMaxEntries;
		dataCacheEntries = _dataCacheEntries;
		dataCacheHits = _dataCacheHits;
		dataCacheMisses = _dataCacheMisses;
		dataCacheMaxMemory = _dataCacheMaxMemory;
		indexCacheMaxEntries = _indexCacheMaxEntries;
		indexCacheEntries = _indexCacheEntries;
		indexCacheHits = _indexCacheHits;
		indexCacheMisses = _indexCacheMisses;
		indexCacheMaxMemory = _indexCacheMaxMemory;
	}

	@Override
	public int getDataCacheMaxEntries() {
		return dataCacheMaxEntries;
	}

	@Override
	public int getDataCacheEntries() {
		return dataCacheEntries;
	}

	@Override
	public long getDataCacheHits() {
		return dataCacheHits;
	}

	@Override
	public long getDataCacheMisses() {
		return dataCacheMisses;
	}

	@Override
	public long getDataCacheMaxMemory() {
		return dataCacheMaxMemory;
	}

	@Override
	public int getIndexCacheMaxEntries() {
		return indexCacheMaxEntries;
	}

	@Override
	public int getIndexCacheEntries() {
		return indexCacheEntries;
	}

	@Override
	public long getIndexCacheHits() {
		return indexCacheHits;
	}

	@Override
	public long getIndexCacheMisses() {
		return indexCacheMisses;
	}

	@Override
	public long getIndexCacheMaxMemory() {
		return indexCacheMaxMemory;
	}

	
}
