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
 *  File name : ConnectionInfoImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: ConnectionInfoImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.admin.TGConnectionInfo;

public class ConnectionInfoImpl implements TGConnectionInfo {
	
	

	@Override
	public String toString() {
		return "ConnectionInfoImpl [listnerName=" + listenerName + ", clientID=" + clientID + ", sessionID=" + sessionID
				+ ", userName=" + userName + ", remoteAddress=" + remoteAddress + ", createdTimeInSeconds="
				+ createdTimeInSeconds + "]";
	}

	public String getListenerName() {
		return listenerName;
	}

	public String getClientID() {
		return clientID;
	}

	public long getSessionID() {
		return sessionID;
	}

	public String getUserName() {
		return userName;
	}

	public String getRemoteAddress() {
		return remoteAddress;
	}

	public long getCreatedTimeInSeconds() {
		return createdTimeInSeconds;
	}

	protected String listenerName;
	protected String clientID;
	protected long sessionID;
	protected String userName;
	protected String remoteAddress;
	protected long createdTimeInSeconds;

	public ConnectionInfoImpl(String _listnerName, String _clientID, long _sessionID, String _userName,
			String _remoteAddress, long _createdTimeInSeconds) {
		
		listenerName = _listnerName;
		clientID = _clientID;
		sessionID = _sessionID;
		userName = _userName;
		remoteAddress = _remoteAddress;
		createdTimeInSeconds = _createdTimeInSeconds;
		
	}
	
	

}
