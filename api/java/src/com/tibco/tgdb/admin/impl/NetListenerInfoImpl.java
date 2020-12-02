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
 *  File name : NetListenerInfoImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: NetListenerInfoImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.admin.TGNetListenerInfo;

public class NetListenerInfoImpl implements TGNetListenerInfo {
	
	

	@Override
	public String toString() {
		return "NetListenerInfoImpl [listenerName=" + listenerName + ", currentConnections=" + currentConnections
				+ ", maxConnections=" + maxConnections + ", portNumber=" + portNumber + "]";
	}

	protected String listenerName;
	protected int currentConnections;
	protected int maxConnections;
	protected String portNumber;

	
	public NetListenerInfoImpl (
		String _listenerName,
		int _currentConnections,
		int _maxConnections,
		String _portNumber
	) 
	{	
		this.listenerName = _listenerName;
		this.currentConnections = _currentConnections;
		this.maxConnections = _maxConnections;
		this.portNumber = _portNumber;
	}
	
	@Override
	public String getListenerName() {
		return this.listenerName;
	}

	@Override
	public int getCurrentConnections() {
		return this.currentConnections;
	}

	@Override
	public int getMaxConnections() {
		return this.maxConnections;
	}

	@Override
	public String getPortNumber() {
		return this.portNumber;
	}

}
