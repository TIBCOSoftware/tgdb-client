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
 * File name: IndexInfoImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type IndexInfoImpl struct {
	sysId      int
	indexType  byte
	name       string
	uniqueFlag bool
	attributes []string
	nodeTypes  []string
}

// Make sure that the IndexInfoImpl implements the TGIndexInfo interface
var _ TGIndexInfo = (*IndexInfoImpl)(nil)

func DefaultIndexInfoImpl() *IndexInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(IndexInfoImpl{})

	return &IndexInfoImpl{}
}

func NewIndexInfoImpl(_sysId int, _name string, _indexType byte,
	_isUnique bool, _attributes, _nodeTypes []string) *IndexInfoImpl {
	newConnectionInfo := DefaultIndexInfoImpl()
	newConnectionInfo.sysId = _sysId
	newConnectionInfo.indexType = _indexType
	newConnectionInfo.name = _name
	newConnectionInfo.uniqueFlag = _isUnique
	newConnectionInfo.attributes = _attributes
	newConnectionInfo.nodeTypes = _nodeTypes
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGIndexInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *IndexInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("IndexInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("SysId: '%d'", obj.sysId))
	buffer.WriteString(fmt.Sprintf(", IndexType: '%+v'", obj.indexType))
	buffer.WriteString(fmt.Sprintf(", Name: '%s'", obj.name))
	buffer.WriteString(fmt.Sprintf(", IsUnique: '%+v'", obj.uniqueFlag))
	buffer.WriteString(fmt.Sprintf(", Attributes: '%+v'", obj.attributes))
	buffer.WriteString(fmt.Sprintf(", NodeTypes: '%+v'", obj.nodeTypes))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGIndexInfo
/////////////////////////////////////////////////////////////////

// GetAttributes returns a collection of attribute names
func (obj *IndexInfoImpl) GetAttributeNames() []string {
	return obj.attributes
}

// GetName returns the index name
func (obj *IndexInfoImpl) GetName() string {
	return obj.name
}

// GetType returns the index type
func (obj *IndexInfoImpl) GetType() byte {
	return obj.indexType
}

// GetSystemId returns the system ID
func (obj *IndexInfoImpl) GetSystemId() int {
	return obj.sysId
}

// GetNodeTypes returns a collection of node types
func (obj *IndexInfoImpl) GetNodeTypes() []string {
	return obj.nodeTypes
}

// IsUnique returns the information whether the index is unique
func (obj *IndexInfoImpl) IsUnique() bool {
	return obj.uniqueFlag
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *IndexInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.sysId, obj.indexType, obj.name, obj.uniqueFlag, obj.attributes, obj.nodeTypes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning IndexInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *IndexInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.sysId, &obj.indexType, &obj.name, &obj.uniqueFlag, &obj.attributes, &obj.nodeTypes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning IndexInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
