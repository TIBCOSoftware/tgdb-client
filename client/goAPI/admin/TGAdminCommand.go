package admin

import (
	"bytes"
)

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGAdminCommand.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= Admin Command Types =======
type AdminCommand int

const (
	AdminCommandInvalid AdminCommand = iota
	AdminCommandCreateUser
	AdminCommandCreateAttrDesc
	AdminCommandCreateIndex
	AdminCommandCreateNodeType
	AdminCommandCreateEdgeType
	AdminCommandShowUsers
	AdminCommandShowAttrDescs
	AdminCommandShowIndices
	AdminCommandShowTypes
	AdminCommandShowInfo
	AdminCommandShowConnections
	AdminCommandDescribe
	AdminCommandSetLogLevel
	AdminCommandStopServer
	AdminCommandCheckpointServer
	AdminCommandDisconnectClient
	AdminCommandKillConnection
)

func (command AdminCommand) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if command&AdminCommandInvalid == AdminCommandInvalid {
		buffer.WriteString("Admin Command Invalid")
	} else if command&AdminCommandCreateUser == AdminCommandCreateUser {
		buffer.WriteString("Admin Command Create User")
	} else if command&AdminCommandCreateAttrDesc == AdminCommandCreateAttrDesc {
		buffer.WriteString("Admin Command Create AttrDesc")
	} else if command&AdminCommandCreateIndex == AdminCommandCreateIndex {
		buffer.WriteString("Admin Command Create Index")
	} else if command&AdminCommandCreateNodeType == AdminCommandCreateNodeType {
		buffer.WriteString("Admin Command Create NodeType")
	} else if command&AdminCommandCreateEdgeType == AdminCommandCreateEdgeType {
		buffer.WriteString("Admin Command Create EdgeType")
	} else if command&AdminCommandShowUsers == AdminCommandShowUsers {
		buffer.WriteString("Admin Command Show Users")
	} else if command&AdminCommandShowAttrDescs == AdminCommandShowAttrDescs {
		buffer.WriteString("Admin Command Show AttrDesc")
	} else if command&AdminCommandShowIndices == AdminCommandShowIndices {
		buffer.WriteString("Admin Command Show Indices")
	} else if command&AdminCommandShowTypes == AdminCommandShowTypes {
		buffer.WriteString("Admin Command Show Types")
	} else if command&AdminCommandShowInfo == AdminCommandShowInfo {
		buffer.WriteString("Admin Command Show Info")
	} else if command&AdminCommandShowConnections == AdminCommandShowConnections {
		buffer.WriteString("Admin Command Show Connections")
	} else if command&AdminCommandDescribe == AdminCommandDescribe {
		buffer.WriteString("Admin Command Describe")
	} else if command&AdminCommandSetLogLevel == AdminCommandSetLogLevel {
		buffer.WriteString("Admin Command Set LogLevel")
	} else if command&AdminCommandStopServer == AdminCommandStopServer {
		buffer.WriteString("Admin Command Stop Server")
	} else if command&AdminCommandCheckpointServer == AdminCommandCheckpointServer {
		buffer.WriteString("Admin Command Checkpoint Server")
	} else if command&AdminCommandDisconnectClient == AdminCommandDisconnectClient {
		buffer.WriteString("Admin Command Disconnect Client")
	} else if command&AdminCommandKillConnection == AdminCommandKillConnection {
		buffer.WriteString("Admin Command Kill Connection")
	}
	return buffer.String()
}
