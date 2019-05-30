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
 *  File name : ServerMemoryInfoImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: ServerMemoryInfoImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.admin.TGMemoryInfo;
import com.tibco.tgdb.admin.TGServerMemoryInfo;

public class ServerMemoryInfoImpl implements TGServerMemoryInfo {
	
	

	

	@Override
	public String toString() {
		return "ServerMemoryInfoImpl [processMemory=" + processMemory + ", sharedMemory=" + sharedMemory + "]";
	}

	protected TGMemoryInfo processMemory;
	protected TGMemoryInfo sharedMemory;
	
	
	public ServerMemoryInfoImpl (
		TGMemoryInfo _processMemory,
		TGMemoryInfo _sharedMemory
	)
	{
		this.processMemory = _processMemory;
		this.sharedMemory = _sharedMemory;
	}
	
	public TGMemoryInfo getMemoryInfo (MEMORY_TYPE type)
	{
		if (type == MEMORY_TYPE.PROCESS)
		{
			return getProcessMemory();
		}
		else if (type == MEMORY_TYPE.SHARED)
		{
			return getSharedMemory();
		}
		else {
			return null;
		}
	}


	//@Override
	public TGMemoryInfo getProcessMemory() {
		return this.processMemory;
	}

	//@Override
	public TGMemoryInfo getSharedMemory() {
		return this.sharedMemory;
	}

}
