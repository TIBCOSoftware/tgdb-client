
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
 *  File name :TGServerStatus.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  <p>This interface allows users to retrieve the status of server.
 *  
 *  SVN Id: $Id: TGServerStatus.java 3120 2019-04-25 21:21:48Z nimish $  
 */

package com.tibco.tgdb.admin;

import java.time.Duration;

import com.tibco.tgdb.TGVersion;


public interface TGServerStatus {
	
	enum ServerStates {
	    Created,
	    Initialized,
	    Started,
	    Suspended,
	    Interrupted,
	    RequestStop,
	    Stopped,
	    ShutDown
	}

	/**
	 * Retrieves the name of the server instance
	 * @return server name
	 */
	public String getName ();
	
	
	/**
	 * Retrieves the server version information
	 * @return {@link TGVersion} information of server
	 */
	public TGVersion getVersion ();
	
	
	/**
	 * Retrieves the state information of server
	 * @return server state information
	 */
	public ServerStates getStates ();
	
	
	/**
	 * Retrieves the process ID of server
	 * @return process ID
	 */
	public String getProcessId ();
	
	
	/**
	 * Retrieves the uptime information of server
	 * @return uptime information
	 */
	public Duration getUptime (); 

}
