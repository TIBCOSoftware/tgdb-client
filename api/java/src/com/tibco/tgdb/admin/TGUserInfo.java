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
 *  File name :TGUserInfo.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  <p>This interface exposes the getters to retrieve the user information
 *  
 *  SVN Id: $Id: TGUserInfo.java 3120 2019-04-25 21:21:48Z nimish $  
 */

package com.tibco.tgdb.admin;

public interface TGUserInfo {

	/**
	 * 
	 * @return user type??
	 */
	public byte getType();
	
	
	/**
	 * Get the system ID for this user
	 * @return system ID
	 */
	public int getId();
	
	
	/**
	 * Get the user name
	 * @return user name
	 */
	public String getName();

	
}