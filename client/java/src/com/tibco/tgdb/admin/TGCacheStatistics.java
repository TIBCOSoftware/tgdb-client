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
 *  File name :TGCacheStatistics.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  <p>This interface allows users to retrieve the Cache Statistics from server
 *  
 *  SVN Id: $Id: TGCacheStatistics.java 3120 2019-04-25 21:21:48Z nimish $
 * 
 */

package com.tibco.tgdb.admin;

public interface TGCacheStatistics {
	
	/**
	 * Retrieves the data-cache max entries
	 * @return data-cache max entries
	 */
	
	public int getDataCacheMaxEntries();
	
	

	/**
	 * Retrieves the data-cache entries
	 * @return data-cache entries
	 */
	public int getDataCacheEntries();
	
	
	
	/**
	 * Retrieves the data-cache hits
	 * @return data-cache hits value
	 */
	public long getDataCacheHits();
	
	
	
	/**
	 * Retrieves the data-cache misses
	 * @return data-cache miss values
	 */
	public long getDataCacheMisses();
	
	
	
	/**
	 * Retrieves the data-cache max memory
	 * @return data-cache max memory value
	 */
	public long getDataCacheMaxMemory();
	
	
	
	/**
	 * Retrieves the index-cache max entries
	 * @return index-cache max entries
	 */
	public int getIndexCacheMaxEntries();
	
	
	
	
	/**
	 * Retrieves the index-cache entries
	 * @return index-cache entries
	 */
	public int getIndexCacheEntries();
	
	
	
	/**
	 * Retrieves the index-cache hits
	 * @return index-cache hits value
	 */
	public long getIndexCacheHits();
	
	
	
	/**
	 * Retrieves the index-cache misses
	 * @return index-cache misses value
	 */
	public long getIndexCacheMisses();
	

	
	/**
	 * Retrieves the index-cache max memory
	 * @return index-cache max memory value
	 */
	public long getIndexCacheMaxMemory();


}
