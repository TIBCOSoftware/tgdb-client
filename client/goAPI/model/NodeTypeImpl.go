package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"strings"
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
 * File name: TGNodeType.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type NodeType struct {
	*EntityType
	pKeys      []*AttributeDescriptor
	idxIds     []int
	numEntries int64
}

func DefaultNodeType() *NodeType {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(NodeType{})

	newNodeType := NodeType{
		EntityType: DefaultEntityType(),
		pKeys:      make([]*AttributeDescriptor, 0),
		idxIds:     make([]int, 0),
		numEntries: 0,
	}
	newNodeType.sysType = types.SystemTypeNode
	return &newNodeType
}

func NewNodeType(name string, parent types.TGEntityType) *NodeType {
	newNodeType := DefaultNodeType()
	newNodeType.name = name
	newNodeType.parent = parent
	return newNodeType
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGNodeType
/////////////////////////////////////////////////////////////////

func (obj *NodeType) GetIndexIds() []int {
	return obj.idxIds
}

func (obj *NodeType) GetNumEntries() int64 {
	return obj.numEntries
}

func (obj *NodeType) SetAttributeMap(attrMap map[string]*AttributeDescriptor) {
	obj.attributes = attrMap
}

func (obj *NodeType) SetAttributeDesc(attrName string, attrDesc *AttributeDescriptor) {
	obj.attributes[attrName] = attrDesc
}

func (obj *NodeType) SetParent(parentEntity types.TGEntityType) {
	obj.parent = parentEntity
}

func (obj *NodeType) SetNumEntries(num int64) {
	obj.numEntries = num
}

func (obj *NodeType) UpdateMetadata(gmd *GraphMetadata) types.TGError {
	logger.Log(fmt.Sprint("Entering NodeType:UpdateMetadata"))
	// Base Class EntityType::UpdateMetadata()
	err := EntityTypeUpdateMetadata(obj, gmd)
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprint("Inside NodeType:UpdateMetadata, updated base entity type's attributes"))
	for id, key := range obj.pKeys {
		attrDesc, err := gmd.GetAttributeDescriptor(key.GetName())
		if err == nil {
			logger.Warning(fmt.Sprintf("WARNING: Continuing loop NodeType:UpdateMetadata - cannot find '%s' attribute descriptor", key.GetName()))
			continue
		}
		if attrDesc.GetAttrType() == types.AttributeTypeInvalid {
			logger.Warning(fmt.Sprint("WARNING: Continuing loop NodeType:UpdateMetadata as attrDesc.GetAttrType() == types.AttributeTypeInvalid"))
			//gLogger.log(TGLevel.Warning, "Cannot find '%s' attribute descriptor", attrName)
			continue
		}
		obj.pKeys[id] = attrDesc.(*AttributeDescriptor)
	}
	logger.Log(fmt.Sprintf("Returning NodeType:UpdateMetadata w/ NO error, for entityType: '%+v'", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNodeType
/////////////////////////////////////////////////////////////////

// GetPKeyAttributeDescriptors returns a set of primary key descriptors
func (obj *NodeType) GetPKeyAttributeDescriptors() []types.TGAttributeDescriptor {
	pkDesc := make([]types.TGAttributeDescriptor, 0)
	for _, pk := range obj.pKeys {
		pkDesc = append(pkDesc, pk)
	}
	return pkDesc
}

// SetPKeyAttributeDescriptors sets primary key descriptors
func (obj *NodeType) SetPKeyAttributeDescriptors(keys []*AttributeDescriptor) {
	obj.pKeys = keys
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntityType
/////////////////////////////////////////////////////////////////

// AddAttributeDescriptor add an attribute descriptor to the map
func (obj *NodeType) AddAttributeDescriptor(attrName string, attrDesc types.TGAttributeDescriptor) {
	if attrName != "" && attrDesc != nil {
		obj.attributes[attrName] = attrDesc.(*AttributeDescriptor)
	}
}

// GetEntityTypeId gets Entity Type id
func (obj *NodeType) GetEntityTypeId() int {
	return obj.id
}

// DerivedFrom gets the parent Entity Type
func (obj *NodeType) DerivedFrom() types.TGEntityType {
	return obj.parent
}

// GetAttributeDescriptor gets the attribute descriptor for the specified name
func (obj *NodeType) GetAttributeDescriptor(attrName string) types.TGAttributeDescriptor {
	attrDesc := obj.attributes[attrName]
	return attrDesc
}

// GetAttributeDescriptors returns a collection of attribute descriptors associated with this Entity Type
func (obj *NodeType) GetAttributeDescriptors() []types.TGAttributeDescriptor {
	attrDescriptors := make([]types.TGAttributeDescriptor, 0)
	for _, attrDesc := range obj.attributes {
		attrDescriptors = append(attrDescriptors, attrDesc)
	}
	return attrDescriptors
}

// SetEntityTypeId sets Entity Type id
func (obj *NodeType) SetEntityTypeId(eTypeId int) {
	obj.id = eTypeId
}

// SetName sets the system object's name
func (obj *NodeType) SetName(eTypeName string) {
	obj.name = eTypeName
}

// SetSystemType sets system object's type
func (obj *NodeType) SetSystemType(eSysType types.TGSystemType) {
	obj.sysType = eSysType
}

func (obj *NodeType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("NodeType:{")
	//buffer.WriteString(fmt.Sprintf("PKeys: %+v", obj.PKeys))
	buffer.WriteString(fmt.Sprintf(", IdxIds: %+v", obj.idxIds))
	buffer.WriteString(fmt.Sprintf(", NumEntries: %+v", obj.numEntries))
	strArray := []string{buffer.String(), obj.entityTypeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSystemObject
/////////////////////////////////////////////////////////////////

// GetName gets the system object's name
func (obj *NodeType) GetName() string {
	return obj.name
}

// GetSystemType gets system object's type
func (obj *NodeType) GetSystemType() types.TGSystemType {
	return obj.sysType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *NodeType) ReadExternal(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering NodeType:ReadExternal"))
	// Base Class EntityType's ReadExternal()
	err := EntityTypeReadExternal(obj, is)
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprint("Inside NodeType:ReadExternal, read base entity type's attributes"))

	attrCount, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read attrCount w/ Error: '%+v'", err.Error()))
		return err
	}
	for i := 0; i < int(attrCount); i++ {
		attrName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read attrName w/ Error: '%+v'", err.Error()))
			return err
		}
		attrDesc := NewAttributeDescriptorWithType(attrName, types.AttributeTypeString)
		obj.pKeys = append(obj.pKeys, attrDesc)
	}

	idxCount, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read idxCount w/ Error: '%+v'", err.Error()))
		return err
	}
	for i := 0; i < int(idxCount); i++ {
		// TODO: Revisit later to get meta data needs to return index definitions
		indexId, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read indexId w/ Error: '%+v'", err.Error()))
			return err
		}
		obj.idxIds = append(obj.idxIds, indexId)
	}

	numEntries, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read numEntries w/ Error: '%+v'", err.Error()))
		return err
	}

	obj.SetNumEntries(numEntries)
	logger.Log(fmt.Sprintf("Returning NodeType:ReadExternal w/ NO error, for NodeType: '%+v'", obj))
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *NodeType) WriteExternal(os types.TGOutputStream) types.TGError {
	logger.Warning(fmt.Sprint("WARNING: Returning NodeType:WriteExternal is not implemented"))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *NodeType) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.sysType, obj.id, obj.name, obj.parent, obj.attributes, obj.pKeys, obj.idxIds, obj.numEntries)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *NodeType) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.sysType, &obj.id, &obj.name, &obj.parent, &obj.attributes, &obj.pKeys,
		&obj.idxIds, &obj.numEntries)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
