package admin

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"time"
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
 * File name: ServerStatusImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type ServerStatusImpl struct {
	name      string
	processId string
	status    ServerStates
	uptime    time.Duration
	version   *utils.TGServerVersion
}

// Make sure that the ServerStatusImpl implements the TGServerStatus interface
var _ TGServerStatus = (*ServerStatusImpl)(nil)

func DefaultServerStatusImpl() *ServerStatusImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ServerStatusImpl{})

	return &ServerStatusImpl{}
}

func NewServerStatusImpl(_name string, _version *utils.TGServerVersion, _processId string, _status ServerStates, _uptime time.Duration) *ServerStatusImpl {
	newServerStatus := DefaultServerStatusImpl()
	newServerStatus.name = _name
	newServerStatus.processId = _processId
	newServerStatus.status = _status
	newServerStatus.uptime = _uptime
	newServerStatus.version = _version
	return newServerStatus
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGServerStatusImpl
/////////////////////////////////////////////////////////////////

func (obj *ServerStatusImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ServerStatusImpl:{")
	buffer.WriteString(fmt.Sprintf("Name: '%s'", obj.name))
	buffer.WriteString(fmt.Sprintf(", ProcessId: '%s'", obj.processId))
	buffer.WriteString(fmt.Sprintf(", Status: '%+v'", obj.status))
	buffer.WriteString(fmt.Sprintf(", Uptime: '%+v'", obj.uptime))
	buffer.WriteString(fmt.Sprintf(", Version: '%+v'", obj.version))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGServerStatus
/////////////////////////////////////////////////////////////////

// GetName returns the name of the server instance
func (obj *ServerStatusImpl) GetName() string {
	return obj.name
}

// GetProcessId returns the process ID of server
func (obj *ServerStatusImpl) GetProcessId() string {
	return obj.processId
}

// GetServerStatus returns the state information of server
func (obj *ServerStatusImpl) GetServerStatus() ServerStates {
	return obj.status
}

// GetUptime returns the uptime information of server
func (obj *ServerStatusImpl) GetUptime() time.Duration {
	return obj.uptime
}

// GetServerVersion returns the server version information
func (obj *ServerStatusImpl) GetServerVersion() *utils.TGServerVersion {
	return obj.version
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerStatusImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.name, obj.processId, obj.status, obj.uptime, obj.version)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerStatusImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerStatusImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.name, &obj.processId, &obj.status, &obj.uptime, &obj.version)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerStatusImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
