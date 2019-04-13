package admin

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
 * File name: ConnectionInfoImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type ConnectionInfoImpl struct {
	listenerName         string
	clientID             string
	sessionID            int64
	userName             string
	remoteAddress        string
	createdTimeInSeconds int64
}

// Make sure that the ConnectionInfoImpl implements the TGConnectionInfo interface
var _ TGConnectionInfo = (*ConnectionInfoImpl)(nil)

func DefaultConnectionInfoImpl() *ConnectionInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ConnectionInfoImpl{})

	return &ConnectionInfoImpl{}
}

func NewConnectionInfoImpl(_listnerName, _clientID string, _sessionID int64,
	_userName, _remoteAddress string, _createdTimeInSeconds int64) *ConnectionInfoImpl {
	newConnectionInfo := DefaultConnectionInfoImpl()
	newConnectionInfo.listenerName = _listnerName
	newConnectionInfo.clientID = _clientID
	newConnectionInfo.sessionID = _sessionID
	newConnectionInfo.userName = _userName
	newConnectionInfo.remoteAddress = _remoteAddress
	newConnectionInfo.createdTimeInSeconds = _createdTimeInSeconds
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGConnectionInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *ConnectionInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ConnectionInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("ListenerName: '%s'", obj.listenerName))
	buffer.WriteString(fmt.Sprintf(", ClientID: '%s'", obj.clientID))
	buffer.WriteString(fmt.Sprintf(", SessionID: '%d'", obj.sessionID))
	buffer.WriteString(fmt.Sprintf(", UserName: '%s'", obj.userName))
	buffer.WriteString(fmt.Sprintf(", RemoteAddress: '%s'", obj.remoteAddress))
	buffer.WriteString(fmt.Sprintf(", CreatedTimeInSeconds: '%d'", obj.createdTimeInSeconds))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConnectionInfo
/////////////////////////////////////////////////////////////////

// GetClientID returns a client ID of listener
func (obj *ConnectionInfoImpl) GetClientID() string {
	return obj.clientID
}

// GetCreatedTimeInSeconds returns a time when the listener was created
func (obj *ConnectionInfoImpl) GetCreatedTimeInSeconds() int64 {
	return obj.createdTimeInSeconds
}

// GetListenerName returns a name of a particular listener
func (obj *ConnectionInfoImpl) GetListenerName() string {
	return obj.listenerName
}

// GetRemoteAddress returns a remote address of listener
func (obj *ConnectionInfoImpl) GetRemoteAddress() string {
	return obj.remoteAddress
}

// GetSessionID returns a session ID of listener
func (obj *ConnectionInfoImpl) GetSessionID() int64 {
	return obj.sessionID
}

// GetUserName returns a user-name associated with listener
func (obj *ConnectionInfoImpl) GetUserName() string {
	return obj.userName
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ConnectionInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.listenerName, obj.clientID, obj.sessionID,
		obj.userName, obj.remoteAddress, obj.createdTimeInSeconds)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ConnectionInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ConnectionInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.listenerName, &obj.clientID, &obj.sessionID,
		&obj.userName, &obj.remoteAddress, &obj.createdTimeInSeconds)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ConnectionInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
