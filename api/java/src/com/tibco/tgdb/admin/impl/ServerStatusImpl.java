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
 *  File name : ServerStatusImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: ServerStatusImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.time.Duration;

import com.tibco.tgdb.TGVersion;
import com.tibco.tgdb.admin.TGServerStatus;

public class ServerStatusImpl implements TGServerStatus {
	
	

	@Override
	public String toString() {
		return "ServerStatusImpl [name=" + name + ", version=" + version + ", status=" + status + ", processId="
				+ processId + ", uptime=" + uptime + "]";
	}

	protected String name;
	//protected ServerVersionInfo version;
	protected TGVersion version;
	protected ServerStates status;
	protected String processId;
	protected Duration uptime; 
	
	public ServerStatusImpl (
		String _name,
		TGVersion _version,
		ServerStates _status,
		String _processId,
		Duration _uptime)
	{
		name = _name;
		version = _version;
		status = _status;
		processId = _processId;
		uptime = _uptime; 
	}

	@Override
	public String getName() {
		return this.name;
	}

	@Override
	/*
	public ServerVersionInfo getVersion() {
		return this.version;
	}
	*/
	public TGVersion getVersion ()
	{
		return this.version;
	}
	

	@Override
	public ServerStates getStates() {
		return this.status;
	}

	@Override
	public String getProcessId() {
		return this.processId;
	}

	@Override
	public Duration getUptime() {
		return this.uptime;
	}

}
