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
 * File name: NetListenerInfoImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type NetListenerInfoImpl struct {
	currentConnections int
	maxConnections     int
	listenerName       string
	portNumber         string
}

// Make sure that the NetListenerInfoImpl implements the TGNetListenerInfo interface
var _ TGNetListenerInfo = (*NetListenerInfoImpl)(nil)

func DefaultNetListenerInfoImpl() *NetListenerInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(NetListenerInfoImpl{})

	return &NetListenerInfoImpl{}
}

func NewNetListenerInfoImpl(_currentConnections, _maxConnections int, _listenerName, _portNumber string) *NetListenerInfoImpl {
	newConnectionInfo := DefaultNetListenerInfoImpl()
	newConnectionInfo.currentConnections = _currentConnections
	newConnectionInfo.maxConnections = _maxConnections
	newConnectionInfo.listenerName = _listenerName
	newConnectionInfo.portNumber = _portNumber
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGNetListenerInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *NetListenerInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("NetListenerInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("CurrentConnections: '%d'", obj.currentConnections))
	buffer.WriteString(fmt.Sprintf(", MaxConnections: '%d'", obj.maxConnections))
	buffer.WriteString(fmt.Sprintf(", ListenerName: '%s'", obj.listenerName))
	buffer.WriteString(fmt.Sprintf(", PortNumber: '%s'", obj.portNumber))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNetListenerInfo
/////////////////////////////////////////////////////////////////

// GetCurrentConnections returns the count of current connections
func (obj *NetListenerInfoImpl) GetCurrentConnections() int {
	return obj.currentConnections
}

// GetMaxConnections returns the count of max connections
func (obj *NetListenerInfoImpl) GetMaxConnections() int {
	return obj.maxConnections
}

// GetListenerName returns the listener name
func (obj *NetListenerInfoImpl) GetListenerName() string {
	return obj.listenerName
}

// GetPortNumber returns the port detail of this listener
func (obj *NetListenerInfoImpl) GetPortNumber() string {
	return obj.portNumber
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *NetListenerInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.currentConnections, obj.maxConnections, obj.listenerName, obj.portNumber)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NetListenerInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *NetListenerInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.currentConnections, &obj.maxConnections, &obj.listenerName, &obj.portNumber)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NetListenerInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
