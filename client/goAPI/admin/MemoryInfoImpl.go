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
 * File name: MemoryInfoImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type MemoryInfoImpl struct {
	freeMemory               int64
	maxMemory                int64
	usedMemory               int64
	sharedMemoryFileLocation string
}

// Make sure that the MemoryInfoImpl implements the TGMemoryInfo interface
var _ TGMemoryInfo = (*MemoryInfoImpl)(nil)

func DefaultMemoryInfoImpl() *MemoryInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(MemoryInfoImpl{})

	return &MemoryInfoImpl{}
}

func NewMemoryInfoImpl(_freeMemory, _maxMemory, _usedMemory int64, _sharedMemoryFileLocation string) *MemoryInfoImpl {
	newConnectionInfo := DefaultMemoryInfoImpl()
	newConnectionInfo.freeMemory = _freeMemory
	newConnectionInfo.maxMemory = _maxMemory
	newConnectionInfo.usedMemory = _usedMemory
	newConnectionInfo.sharedMemoryFileLocation = _sharedMemoryFileLocation
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGMemoryInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *MemoryInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("MemoryInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("FreeMemory: '%d'", obj.freeMemory))
	buffer.WriteString(fmt.Sprintf(", MaxMemory: '%d'", obj.maxMemory))
	buffer.WriteString(fmt.Sprintf(", UsedMemory: '%d'", obj.usedMemory))
	buffer.WriteString(fmt.Sprintf(", SharedMemoryFileLocation: '%s'", obj.sharedMemoryFileLocation))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMemoryInfo
/////////////////////////////////////////////////////////////////

// GetFreeMemory returns the free memory size from server
func (obj *MemoryInfoImpl) GetFreeMemory() int64 {
	return obj.freeMemory
}

// GetMaxMemory returns the max memory size from server
func (obj *MemoryInfoImpl) GetMaxMemory() int64 {
	return obj.maxMemory
}

// GetSharedMemoryFileLocation returns the shared memory file location
func (obj *MemoryInfoImpl) GetSharedMemoryFileLocation() string {
	return obj.sharedMemoryFileLocation
}

// GetUsedMemory returns the used memory size from server
func (obj *MemoryInfoImpl) GetUsedMemory() int64 {
	return obj.usedMemory
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *MemoryInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.freeMemory, obj.maxMemory, obj.usedMemory, obj.sharedMemoryFileLocation)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MemoryInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *MemoryInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.freeMemory, &obj.maxMemory, &obj.usedMemory, &obj.sharedMemoryFileLocation)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MemoryInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
