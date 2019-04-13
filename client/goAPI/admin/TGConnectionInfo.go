package admin

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
 * File name: TGConnectionInfo.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// TGConnectionInfo allows users to retrieve the individual Connection Information from server
type TGConnectionInfo interface {
	// GetClientID returns a client ID of listener
	GetClientID() string
	// GetCreatedTimeInSeconds returns a time when the listener was created
	GetCreatedTimeInSeconds() int64
	// GetListenerName returns a name of a particular listener
	GetListenerName() string
	// GetRemoteAddress returns a remote address of listener
	GetRemoteAddress() string
	// GetSessionID returns a session ID of listener
	GetSessionID() int64
	// GetUserName returns a user-name associated with listener
	GetUserName() string
}
