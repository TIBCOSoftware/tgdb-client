package types

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
 * File name: TGAdminConnection.go
 * Created on: Mar 03, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGAdminConnection interface {
	TGConnection
	//// GetAttributeDescriptors gets the list of attribute descriptors
	//GetAttributeDescriptors() ([]TGAttributeDescriptor, TGError)
	//
	//// GetConnections gets the list of all socket connections using this connection type
	//GetConnections() (TGAdminShowConnections, TGError)
	//
	//// GetIndices gets the list of all indices
	//GetIndices() ([]TGAdminShowIndex, TGError)
	//
	//// GetInfo gets the information about this connection type
	//GetInfo() (TGAdminCommandInfoResult, TGError)
	//
	//// GetUsers gets the list of users
	//GetUsers() (TGAdminShowUsers, TGError)
	//
	//// StopServer stops the admin connection
	//StopServer() TGError
	//
	//// CheckpointServer allows the programmatic control to do a checkpoint on server
	//CheckpointServer() TGError
	//
	//// DumpServerStackTrace prints the stack trace
	//DumpServerStackTrace() TGError
	//
	//// KillConnection terminates the connection forcefully
	//KillConnection(killConnectionInfo TGAdminKillConnectionInfo ) TGError
	//
	//// setServerLogLevel set the log level
	//setServerLogLevel(logLevel int, logComponent int64) TGError
}
