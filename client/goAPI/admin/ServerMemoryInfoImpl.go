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
 * File name: ServerMemoryInfoImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type ServerMemoryInfoImpl struct {
	processMemory *MemoryInfoImpl
	sharedMemory  *MemoryInfoImpl
}

// Make sure that the ServerMemoryInfoImpl implements the TGServerMemoryInfo interface
var _ TGServerMemoryInfo = (*ServerMemoryInfoImpl)(nil)

func DefaultServerMemoryInfoImpl() *ServerMemoryInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ServerMemoryInfoImpl{})

	return &ServerMemoryInfoImpl{}
}

func NewServerMemoryInfoImpl(_processMemory, _sharedMemory *MemoryInfoImpl) *ServerMemoryInfoImpl {
	newConnectionInfo := DefaultServerMemoryInfoImpl()
	newConnectionInfo.processMemory = _processMemory
	newConnectionInfo.sharedMemory = _sharedMemory
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGServerMemoryInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *ServerMemoryInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ServerMemoryInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("ProcessMemory: '%+v'", obj.processMemory))
	buffer.WriteString(fmt.Sprintf(", SharedMemory: '%+v'", obj.sharedMemory))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGServerMemoryInfo
/////////////////////////////////////////////////////////////////

// GetProcessMemory returns server process memory
func (obj *ServerMemoryInfoImpl) GetProcessMemory() TGMemoryInfo {
	return obj.processMemory
}

// GetSharedMemory returns server shared memory
func (obj *ServerMemoryInfoImpl) GetSharedMemory() TGMemoryInfo {
	return obj.sharedMemory
}

// GetMemoryInfo returns the memory info for the specified type
func (obj *ServerMemoryInfoImpl) GetServerMemoryInfo(memType MemType) TGMemoryInfo {
	if memType == MemoryProcess {
		return obj.GetProcessMemory()
	} else if memType == MemoryShared {
		return obj.GetSharedMemory()
	} else {
		return nil
	}
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerMemoryInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.processMemory, obj.sharedMemory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerMemoryInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerMemoryInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.processMemory, &obj.sharedMemory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerMemoryInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
