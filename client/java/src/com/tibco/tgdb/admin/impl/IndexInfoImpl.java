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
 *  File name : IndexInfoImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: IndexInfoImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.util.Collection;

import com.tibco.tgdb.admin.TGIndexInfo;

public class IndexInfoImpl implements TGIndexInfo {

	protected int sysid;
	protected byte type;
	protected String name;
	protected boolean isUnique;
	protected Collection<String> attributes;
	protected Collection<String> nodeTypes;
	protected long numEntries;
	protected String status;

	
	public IndexInfoImpl(int sysid, byte type, String name, boolean isUnique, Collection<String> attributes,
			Collection<String> nodeTypes, long numEntries, String status) {
		super();
		this.sysid = sysid;
		this.type = type;
		this.name = name;
		this.isUnique = isUnique;
		this.attributes = attributes;
		this.nodeTypes = nodeTypes;
		this.numEntries = numEntries;
		this.status = status;
	}
	
	public int getSysid() {
		return sysid;
	}
	public byte getType() {
		return type;
	}
	public String getName() {
		return name;
	}
	public boolean isUnique() {
		return isUnique;
	}
	public Collection<String> getAttributes() {
		return attributes;
	}
	public Collection<String> getNodeTypes() {
		return nodeTypes;
	}
	

	@Override
	public long getNumEntries() {
		return numEntries;
	}

	@Override
	public String getStatus() {
		return status;
	}

	@Override
	public String toString() {
		return "IndexInfoImpl [sysid=" + sysid + ", type=" + type + ", name=" + name + ", isUnique=" + isUnique
				+ ", attributes=" + attributes + ", nodeTypes=" + nodeTypes + ", numEntries=" + numEntries + ", status="
				+ status + "]";
	}
	
	
}
