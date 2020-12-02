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
 *  File name : UserInfoImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: UserInfoImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.admin.TGUserInfo;

public class UserInfoImpl implements TGUserInfo {

	@Override
	public String toString() {
		return "UserInfoImpl [type=" + type + ", id=" + id + ", name=" + name + "]";
	}

	protected byte type;
	protected int id;
	protected String name;
	
	public UserInfoImpl(byte _type, int _id, String _name) {
		this.type = _type;
		this.id = _id;
		this.name = _name;
	}
	
	public byte getType() {
		return type;
	}
	public int getId() {
		return id;
	}
	public String getName() {
		return name;
	}
	
	
	
}
