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
 *  File name : TGAdminCommand.java
 *  Created on: 04/22/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: TGAdminCommand.java 3813 2020-03-19 20:25:34Z dhudson $
 * 
 */

package com.tibco.tgdb.admin.impl;

public enum TGAdminCommand {
	
	CreateUser,
	CreateRole,
	CreateAttrDesc,
	CreateIndex,
	CreateNodeType,
	CreateEdgeType,
	ShowUsers,
	ShowRoles,
	ShowAttrDescs,
	ShowIndices,
	ShowTypes,
	ShowInfo,
	ShowConnections,
	Describe,
	SetLogLevel,
	SetSPDirectory,
	UpdateRole,
	StopServer,
	CheckpointServer,
	DisconnectClient,
	KillConnection;

	public static final TGAdminCommand values[] = values();
	
}
