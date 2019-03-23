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
 * File name: TGEdgeType.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

//static TGLogger gLogger = TGLogManager.getInstance().getLogger();

type EdgeType struct {
	*EntityType
	directionType types.TGDirectionType
	fromTypeId    int
	fromNodeType  types.TGNodeType
	toTypeId      int
	toNodeType    types.TGNodeType
	numEntries    int64
}

func DefaultEdgeType() *EdgeType {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(EdgeType{})

	newEdgeType := &EdgeType{
		EntityType: DefaultEntityType(),
		numEntries: 0,
	}
	newEdgeType.sysType = types.SystemTypeEdge
	return newEdgeType
}

func NewEdgeType(name string, directionType types.TGDirectionType, parent types.TGEntityType) *EdgeType {
	newEdgeType := DefaultEdgeType()
	newEdgeType.name = name
	newEdgeType.directionType = directionType
	newEdgeType.parent = parent
	return newEdgeType
}

/////////////////////////////////////////////////////////////////
// Helper functions for EdgeType
/////////////////////////////////////////////////////////////////

func (obj *EdgeType) SetAttributeMap(attrMap map[string]*AttributeDescriptor) {
	obj.attributes = attrMap
}

func (obj *EdgeType) SetAttributeDesc(attrName string, attrDesc *AttributeDescriptor) {
	obj.attributes[attrName] = attrDesc
}

func (obj *EdgeType) SetParent(parentEntity types.TGEntityType) {
	obj.parent = parentEntity
}

func (obj *EdgeType) SetDirectionType(dirType types.TGDirectionType) {
	obj.directionType = dirType
}

func (obj *EdgeType) SetNumEntries(num int64) {
	obj.numEntries = num
}

func (obj *EdgeType) UpdateMetadata(gmd *GraphMetadata) types.TGError {
	logger.Log(fmt.Sprint("Entering EdgeType:UpdateMetadata"))
	// Base Class EntityType::UpdateMetadata()
	err := EntityTypeUpdateMetadata(obj, gmd)
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprint("Inside EdgeType:UpdateMetadata, updated base entity type's attributes"))
	nType, nErr := gmd.GetNodeTypeById(obj.fromTypeId)
	if nErr == nil {
		obj.fromNodeType = nType.(*NodeType)
	}
	nType, nErr = gmd.GetNodeTypeById(obj.toTypeId)
	if nErr == nil {
		obj.toNodeType = nType.(*NodeType)
	}
	logger.Log(fmt.Sprintf("Returning EdgeType:UpdateMetadata w/ NO error, for EdgeType: '%+v'", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNodeType
/////////////////////////////////////////////////////////////////

// GetDirectionType gets direction type as one of the constants
func (obj *EdgeType) GetDirectionType() types.TGDirectionType {
	return obj.directionType
}

// GetFromNodeType gets From-Node Type
func (obj *EdgeType) GetFromNodeType() types.TGNodeType {
	return obj.fromNodeType
}

// GetFromTypeId gets From-Node ID
func (obj *EdgeType) GetFromTypeId() int {
	return obj.fromTypeId
}

// GetToNodeType gets To-Node Type
func (obj *EdgeType) GetToNodeType() types.TGNodeType {
	return obj.toNodeType
}

// GetToTypeId gets To-Node ID
func (obj *EdgeType) GetToTypeId() int {
	return obj.toTypeId
}

// SetFromNodeType sets From-Node Type
func (obj *EdgeType) SetFromNodeType(fromNode types.TGNodeType) {
	obj.fromNodeType = fromNode
}

// SetFromTypeId sets From-Node ID
func (obj *EdgeType) SetFromTypeId(fromTypeId int) {
	obj.fromTypeId = fromTypeId
}

// SetToNodeType sets From-Node Type
func (obj *EdgeType) SetToNodeType(toNode types.TGNodeType) {
	obj.toNodeType = toNode
}

// SetToTypeId sets To-Node ID
func (obj *EdgeType) SetToTypeId(toTypeId int) {
	obj.toTypeId = toTypeId
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntityType
/////////////////////////////////////////////////////////////////

// AddAttributeDescriptor add an attribute descriptor to the map
func (obj *EdgeType) AddAttributeDescriptor(attrName string, attrDesc types.TGAttributeDescriptor) {
	if attrName != "" && attrDesc != nil {
		obj.attributes[attrName] = attrDesc.(*AttributeDescriptor)
	}
}

// GetEntityTypeId gets Entity Type id
func (obj *EdgeType) GetEntityTypeId() int {
	return obj.id
}

// DerivedFrom gets the parent Entity Type
func (obj *EdgeType) DerivedFrom() types.TGEntityType {
	return obj.parent
}

// GetAttributeDescriptor gets the attribute descriptor for the specified name
func (obj *EdgeType) GetAttributeDescriptor(attrName string) types.TGAttributeDescriptor {
	attrDesc := obj.attributes[attrName]
	return attrDesc
}

// GetAttributeDescriptors returns a collection of attribute descriptors associated with this Entity Type
func (obj *EdgeType) GetAttributeDescriptors() []types.TGAttributeDescriptor {
	attrDescriptors := make([]types.TGAttributeDescriptor, 0)
	for _, attrDesc := range obj.attributes {
		attrDescriptors = append(attrDescriptors, attrDesc)
	}
	return attrDescriptors
}

// SetEntityTypeId sets Entity Type id
func (obj *EdgeType) SetEntityTypeId(eTypeId int) {
	obj.id = eTypeId
}

// SetName sets the system object's name
func (obj *EdgeType) SetName(eTypeName string) {
	obj.name = eTypeName
}

// SetSystemType sets system object's type
func (obj *EdgeType) SetSystemType(eSysType types.TGSystemType) {
	obj.sysType = eSysType
}

func (obj *EdgeType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("EdgeType:{")
	buffer.WriteString(fmt.Sprintf("DirectionType: %+v", obj.directionType))
	buffer.WriteString(fmt.Sprintf(", FromTypeId: %+v", obj.fromTypeId))
	buffer.WriteString(fmt.Sprintf(", FromNodeType: %+v", obj.fromNodeType))
	buffer.WriteString(fmt.Sprintf(", ToTypeId: %+v", obj.toTypeId))
	buffer.WriteString(fmt.Sprintf(", ToNodeType: %+v", obj.toNodeType))
	buffer.WriteString(fmt.Sprintf(", NumEntries: %+v", obj.numEntries))
	strArray := []string{buffer.String(), obj.entityTypeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSystemObject
/////////////////////////////////////////////////////////////////

// GetName gets the system object's name
func (obj *EdgeType) GetName() string {
	return obj.name
}

// GetSystemType gets system object's type
func (obj *EdgeType) GetSystemType() types.TGSystemType {
	return obj.sysType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *EdgeType) ReadExternal(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering EdgeType:ReadExternal"))
	// Base Class EntityType's ReadExternal()
	err := EntityTypeReadExternal(obj, is)
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprint("Inside EdgeType:ReadExternal, read base entity type's attributes"))

	fromTypeId, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:ReadExternal - unable to read fromTypeId w/ Error: '%+v'", err.Error()))
		return err
	}

	toTypeId, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:ReadExternal - unable to read toTypeId w/ Error: '%+v'", err.Error()))
		return err
	}

	direction, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:ReadExternal - unable to read direction w/ Error: '%+v'", err.Error()))
		return err
	}
	if direction == 0 {
		obj.SetDirectionType(types.DirectionTypeUnDirected)
	} else if direction == 1 {
		obj.SetDirectionType(types.DirectionTypeDirected)
	} else {
		obj.SetDirectionType(types.DirectionTypeBiDirectional)
	}

	numEntries, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:ReadExternal - unable to read numEntries w/ Error: '%+v'", err.Error()))
		return err
	}

	obj.SetFromTypeId(fromTypeId)
	obj.SetToTypeId(toTypeId)
	obj.SetNumEntries(numEntries)
	logger.Log(fmt.Sprintf("Returning EdgeType:ReadExternal w/ NO error, for entityType: '%+v'", obj))
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *EdgeType) WriteExternal(os types.TGOutputStream) types.TGError {
	logger.Warning(fmt.Sprint("WARNING: Returning EdgeType:WriteExternal is not implemented"))
	//errMsg := fmt.Sprint("EdgeType WriteExternal message is not implemented")
	//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *EdgeType) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.sysType, obj.id, obj.name, obj.parent, obj.attributes, obj.directionType,
		obj.fromTypeId, obj.fromNodeType, obj.toTypeId, obj.toNodeType, obj.numEntries)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *EdgeType) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.sysType, &obj.id, &obj.name, &obj.parent, &obj.attributes, &obj.directionType,
		&obj.fromTypeId, &obj.fromNodeType, &obj.toTypeId, &obj.toNodeType, &obj.numEntries)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
