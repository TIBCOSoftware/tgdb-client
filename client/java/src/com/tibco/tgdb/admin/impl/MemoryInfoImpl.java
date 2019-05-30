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
 *  File name : MemoryInfoImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: MemoryInfoImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.admin.TGMemoryInfo;

public class MemoryInfoImpl implements TGMemoryInfo {

	

	@Override
	public String toString() {
		return "MemoryInfoImpl [usedMemory=" + usedMemory + ", freeMemory=" + freeMemory + ", maxMemory=" + maxMemory
				+ ", sharedMemoryFileLocation=" + sharedMemoryFileLocation + "]";
	}

	protected long usedMemory;
	protected long freeMemory;
	protected long maxMemory;
	protected String sharedMemoryFileLocation;
	
	public String getSharedMemoryFileLocation() {
		return sharedMemoryFileLocation;
	}

	public MemoryInfoImpl (
		long _usedMemory,
		long _freeMemory,
		long _maxMemory,
		String _sharedMemoryFileLocation)
	{
		usedMemory = _usedMemory;
		freeMemory = _freeMemory;
		maxMemory = _maxMemory;
		sharedMemoryFileLocation = _sharedMemoryFileLocation;
	}
	
	@Override
	public long getUsedMemory() {
		return usedMemory;
	}

	@Override
	public long getFreeMemory() {
		return freeMemory;
	}

	@Override
	public long getMaxMemory() {
		return maxMemory;
	}

}
