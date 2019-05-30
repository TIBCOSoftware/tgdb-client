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
 *  File name :TGIndexInfo.java
 *  Created on: 03/28/2019
 *  Created by: nimish
  
 *  <p>This interface allows users to retrieve the index information from server  
 * 
 *  SVN Id: $Id: TGIndexInfo.java 3120 2019-04-25 21:21:48Z nimish $ 
 */

package com.tibco.tgdb.admin;

import java.util.Collection;


public interface TGIndexInfo {
	
	/**
	 * Retrieves the system ID
	 * 
	 * @return system ID
	 */
	public int getSysid();

	

	/**
	 * Retrieves the index type
	 * @return type
	 */
	public byte getType();

	

	/**
	 * Retrieves the index name
	 * @return index name
	 */
	public String getName();
	
	

	/**
	 * Retrieves the information whether the index is unique
	 * @return result whether the index is unique  
	 */
	public boolean isUnique();
	

	
	/**
	 * Retrieves a collection of attribute names
	 * @return a collection of attribute names
	 */
	public Collection<String> getAttributes();
	
	
	
	/**
	 * Retrieves a collection of node types
	 * @return a collection of node types
	 */
	public Collection<String> getNodeTypes();
	
	/**
	 * Retrieves number of entries for the index
	 * @return number of entries for the index
	 */
	public long getNumEntries ();
	
	
	/**
	 * Retrieves the index status
	 * @return the index status
	 */
	public String getStatus ();


}
