
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
 *  File name :TGMemoryInfo.java
 *  Created on: 03/28/2019
 *  Created by: nimish
  
 *  <p>This interface allows users to retrieve the memory information from server
 * 
 *  SVN Id: $Id: TGMemoryInfo.java 3120 2019-04-25 21:21:48Z nimish $ 
 */

package com.tibco.tgdb.admin;

import com.tibco.tgdb.admin.TGServerMemoryInfo.MEMORY_TYPE;


public interface TGMemoryInfo {
	
	/**
	 * Retrieve the used memory size from server  
	 * @return used memory size
	 */
	public long getUsedMemory();
	
	
	/**
	 * Retrieve the free memory size from server  
	 * @return free memory size
	 */
	public long getFreeMemory();
	
	
	/**
	 * Retrieve the max memory size from server  
	 * @return max memory size
	 */
	public long getMaxMemory();
	
	
	/**
	 * Retrieve the shared memory file location   
	 * @return shared memory file location; for {@link MEMORY_TYPE} {@code PROCESS}, this 
	 * method will return null value
	 */

	public String getSharedMemoryFileLocation ();
	
}
