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
 *  File name :TGDatabaseStatistics.java
 *  Created on: 03/28/2019
 *  Created by: nimish
  
 *  <p>This interface allows users to retrieve the database statistics from server
 *    
 *  SVN Id: $Id: TGDatabaseStatistics.java 3120 2019-04-25 21:21:48Z nimish $ 
 */

package com.tibco.tgdb.admin;

public interface TGDatabaseStatistics {

	
	/**
	 * Retrieves the size of database
	 * 
	 * @return database size
	 */
	public long getDbSize();
	
	
	/**
	 * Retrieves the number of data segments
	 * 
	 * @return number of data segments
	 */
	public int getNumDataSegments();
	
	
	/**
	 * Retrieves the datasize
	 * @return datasize
	 */
	public long getDataSize();
	
	
	/**
	 * Retrieves the size of data used
	 * 
	 * @return data used size
	 */
	public long getDataUsed();
	
	
	
	/**
	 * Retrieves the free data size 
	 * 
	 * @return free data size
	 */
	public long getDataFree();
	
	
	/**
	 * Retrieves the block size of data
	 * 
	 * @return data block size
	 */
	public int getDataBlockSize();
	
	
	/**
	 * Retrieves the number of index segments
	 * 
	 * @return index segment count
	 */
	public int getNumIndexSegments();
	
	
	/**
	 * Retrieves the index size
	 * 
	 * @return index size
	 */
	public long getIndexSize();
	
	
	/**
	 * Retrieves the used index size
	 * 
	 * @return used index size
	 */
	public long getIndexUsed();
	
	
	/**
	 * Retrieves the free index size
	 * 
	 * @return free index size
	 */
	public long getIndexFree();
	
	
	/**
	 * Retrieves the block size
	 * 
	 * @return block size
	 */
	public int getBlockSize();	
	
}