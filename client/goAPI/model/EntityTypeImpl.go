package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * File name: TGEntityType.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type EntityType struct {
	sysType    types.TGSystemType
	id         int // Issued only for creation and not valid later
	name       string
	parent     types.TGEntityType
	attributes map[string]*AttributeDescriptor
}

func DefaultEntityType() *EntityType {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(EntityType{})

	newEntityType := EntityType{
		//SysType:    types.SystemTypeEntity,
		attributes: make(map[string]*AttributeDescriptor, 0),
	}
	// TODO: Check with TGDB Engineering Team how and when Id will be set - It is supposed to be set at the time of creation
	return &newEntityType
}

func NewEntityType(name string, parent *EntityType) *EntityType {
	newEntityType := DefaultEntityType()
	newEntityType.name = name
	newEntityType.parent = parent
	return newEntityType
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGEntityType
/////////////////////////////////////////////////////////////////

func (obj *EntityType) UpdateMetadata(gmd *GraphMetadata) types.TGError {
	return EntityTypeUpdateMetadata(obj, gmd)
}

func (obj *EntityType) SetAttributeMap(attrMap map[string]*AttributeDescriptor) {
	obj.attributes = attrMap
}

func (obj *EntityType) SetAttributeDesc(attrName string, attrDesc *AttributeDescriptor) {
	obj.attributes[attrName] = attrDesc
}

func (obj *EntityType) SetParent(parentEntity types.TGEntityType) {
	obj.parent = parentEntity
}

func EntityTypeReadExternal(obj types.TGEntityType, is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering EntityType:EntityTypeReadExternal"))
	// TODO: Revisit later - Do we save the desc value?
	sType, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read sType w/ Error: '%+v'", err.Error()))
		return err
	}
	if types.TGSystemType(sType) == types.SystemTypeInvalid {
		logger.Warning(fmt.Sprint("WARNING: EntityType:EntityTypeReadExternal - types.TGSystemType(sType) == types.SystemTypeInvalid"))
		// TODO: Revisit later - Do we need to throw Exception?
		//errMsg := fmt.Sprintf("Entity desc input stream has invalid desc value: %d", sType)
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	eId, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read eId w/ Error: '%+v'", err.Error()))
		return err
	}

	eName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read eName w/ Error: '%+v'", err.Error()))
		return err
	}

	_, err = is.(*iostream.ProtocolDataInputStream).ReadInt() // pagesize
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read pageSize w/ Error: '%+v'", err.Error()))
		return err
	}

	// TODO: Check with TGDB Engineering Team why parent is not being sent over

	attrCount, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read attrCount w/ Error: '%+v'", err.Error()))
		return err
	}

	for i := 0; i < int(attrCount); i++ {
		attrName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read aName w/ Error: '%+v'", err.Error()))
			return err
		}
		// TODO: The stream only contains name of the descriptor. Do we need to lookup the descriptor from GraphMetaData?
		attrDesc := NewAttributeDescriptorWithType(attrName, types.AttributeTypeString)
		obj.AddAttributeDescriptor(attrName, attrDesc)
	}
	obj.SetEntityTypeId(eId)
	obj.SetName(eName)
	obj.SetSystemType(types.TGSystemType(sType))
	logger.Log(fmt.Sprintf("Returning EntityType:EntityTypeReadExternal w/ NO error, for entityType: '%+v'", obj))
	return nil
}

func EntityTypeUpdateMetadata(obj types.TGEntityType, gmd *GraphMetadata) types.TGError {
	logger.Log(fmt.Sprint("Entering EntityType:EntityTypeUpdateMetadata"))
	for attrName, _ := range obj.(*EntityType).attributes {
		attrDesc, err := gmd.GetAttributeDescriptor(attrName)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeUpdateMetadata - unable to get attrDesc w/ Error: '%+v'", err.Error()))
			return err
		}
		// TODO: Revisit later - Something not correct - should we continue or throw an error
		if attrDesc == nil {
			logger.Warning(fmt.Sprintf("WARNING: Continuing loop EntityType:EntityTypeUpdateMetadata - cannot find '%s' attribute descriptor", attrName))
			//gLogger.log(TGLevel.Warning, "Cannot find '%s' attribute descriptor", attrName)
			continue
		}
		if attrDesc.GetAttrType() == types.AttributeTypeInvalid {
			logger.Warning(fmt.Sprint("WARNING: Continuing loop EntityType:EntityTypeUpdateMetadata as attrDesc.GetAttrType() == types.AttributeTypeInvalid"))
			continue
		}
		obj.AddAttributeDescriptor(attrName, attrDesc)
	}
	logger.Log(fmt.Sprintf("Returning EntityType:EntityTypeUpdateMetadata w/ NO error, for entityType: '%+v'", obj))
	return nil
}

func (obj *EntityType) entityTypeToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("EntityType:{")
	buffer.WriteString(fmt.Sprintf("SysType: %d", obj.sysType))
	buffer.WriteString(fmt.Sprintf(", Id: %+v", obj.id))
	buffer.WriteString(fmt.Sprintf(", Name: %+v", obj.name))
	//buffer.WriteString(fmt.Sprintf(", Parent: %+v", obj.parent))
	//buffer.WriteString(fmt.Sprintf(", EntityTypeAttributes: %+v", obj.attributes))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntityType
/////////////////////////////////////////////////////////////////

// AddAttributeDescriptor add an attribute descriptor to the map
func (obj *EntityType) AddAttributeDescriptor(attrName string, attrDesc types.TGAttributeDescriptor) {
	if attrName != "" && attrDesc != nil {
		obj.attributes[attrName] = attrDesc.(*AttributeDescriptor)
	}
}

// GetEntityTypeId gets Entity Type id
func (obj *EntityType) GetEntityTypeId() int {
	return obj.id
}

// DerivedFrom gets the parent Entity Type
func (obj *EntityType) DerivedFrom() types.TGEntityType {
	return obj.parent
}

// GetAttributeDescriptor gets the attribute descriptor for the specified name
func (obj *EntityType) GetAttributeDescriptor(attrName string) types.TGAttributeDescriptor {
	attrDesc := obj.attributes[attrName]
	return attrDesc
}

// GetAttributeDescriptors returns a collection of attribute descriptors associated with this Entity Type
func (obj *EntityType) GetAttributeDescriptors() []types.TGAttributeDescriptor {
	attrDescriptors := make([]types.TGAttributeDescriptor, 0)
	for _, attrDesc := range obj.attributes {
		attrDescriptors = append(attrDescriptors, attrDesc)
	}
	return attrDescriptors
}

// SetEntityTypeId sets Entity Type id
func (obj *EntityType) SetEntityTypeId(eTypeId int) {
	obj.id = eTypeId
}

// SetName sets the system object's name
func (obj *EntityType) SetName(eTypeName string) {
	obj.name = eTypeName
}

// SetSystemType sets system object's type
func (obj *EntityType) SetSystemType(eSysType types.TGSystemType) {
	obj.sysType = eSysType
}

func (obj *EntityType) String() string {
    return obj.entityTypeToString()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSystemObject
/////////////////////////////////////////////////////////////////

// GetName gets the name for this entity type as the most generic form
func (obj *EntityType) GetName() string {
	return obj.name
}

func (obj *EntityType) GetSystemType() types.TGSystemType {
	return obj.sysType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *EntityType) ReadExternal(is types.TGInputStream) types.TGError {
	return EntityTypeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *EntityType) WriteExternal(os types.TGOutputStream) types.TGError {
	logger.Warning(fmt.Sprint("WARNING: Returning EntityType:WriteExternal is not implemented"))
	//errMsg := fmt.Sprint("EntityType WriteExternal message is not implemented")
	//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *EntityType) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.sysType, obj.id, obj.name, obj.parent, obj.attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *EntityType) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.sysType, &obj.id, &obj.name, &obj.parent, &obj.attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
