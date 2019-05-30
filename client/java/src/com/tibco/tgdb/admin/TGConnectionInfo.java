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
 *  File name :TGConnectionInfo.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  <p>This interface allows users to retrieve the individual Connection Information from server
 *  SVN Id: $Id: TGConnectionInfo.java 3120 2019-04-25 21:21:48Z nimish $
 * 
 */

package com.tibco.tgdb.admin;

public interface TGConnectionInfo {

	/**
	 * Retrieves a name of a particular listener
	 * @return listener name
	 */
	public String getListenerName();

	
	
	/**
	 * Retrieves a client ID of listener
	 * @return client ID
	 */
	public String getClientID();

	
	/**
	 * Retrieves a session ID of listener
	 * @return session ID
	 */
	public long getSessionID();

	
	/**
	 * Retrieves a user-name associated with listener
	 * @return user name
	 */
	public String getUserName();

	
	/**
	 * Retrieves a remote address of listener
	 * @return remote address
	 */
	public String getRemoteAddress();

	
	/**
	 * Retrieves a time when the listener was created 
	 * @return long value corresponds to the time when listener was created
	 */
	public long getCreatedTimeInSeconds();
	
}
