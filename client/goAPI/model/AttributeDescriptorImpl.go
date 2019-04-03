package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"sync/atomic"
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
 * File name: TGAttributeDescriptor.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

//var gLogger = TGLogManager.getInstance().getLogger()
var LocalAttributeId int64

type AttributeDescriptor struct {
	sysType     types.TGSystemType
	attributeId int64
	name        string
	attrType    int
	isArray     bool
	isEncrypted bool
	precision   int16
	scale       int16
}

func DefaultAttributeDescriptor() *AttributeDescriptor {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AttributeDescriptor{})

	newAttributeDescriptor := AttributeDescriptor{
		sysType:     types.SystemTypeAttributeDescriptor,
		name:        "",
		attrType:    types.AttributeTypeInvalid,
		isArray:     false,
		isEncrypted: false,
		precision:   0,
		scale:       0,
	}
	newAttributeDescriptor.attributeId = atomic.AddInt64(&LocalAttributeId, 1)
	return &newAttributeDescriptor
}

func NewAttributeDescriptor(id int64) *AttributeDescriptor {
	newAttributeDescriptor := DefaultAttributeDescriptor()
	newAttributeDescriptor.attributeId = id
	return newAttributeDescriptor
}

func NewAttributeDescriptorWithType(name string, attrType int) *AttributeDescriptor {
	newAttributeDescriptor := DefaultAttributeDescriptor()
	newAttributeDescriptor.name = name
	newAttributeDescriptor.attrType = attrType
	if attrType == types.AttributeTypeNumber {
		newAttributeDescriptor.precision = 20
		newAttributeDescriptor.scale = 5
	}
	return newAttributeDescriptor
}

func NewAttributeDescriptorAsArray(name string, attrType int, isArray bool) *AttributeDescriptor {
	newAttributeDescriptor := NewAttributeDescriptorWithType(name, attrType)
	newAttributeDescriptor.isArray = isArray
	return newAttributeDescriptor
}

// TODO: To be used when created from server side data
func NewAttributeDescriptorOnServer(name string, attrType int, isArray bool, attributeId int64) *AttributeDescriptor {
	newAttributeDescriptor := NewAttributeDescriptorAsArray(name, attrType, isArray)
	newAttributeDescriptor.attributeId = attributeId
	return newAttributeDescriptor
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGAbstractEntity
/////////////////////////////////////////////////////////////////

func (obj *AttributeDescriptor) SetAttributeId(attrId int64) {
	obj.attributeId = attrId
}

func (obj *AttributeDescriptor) SetAttrType(attrType int) {
	obj.attrType = attrType
}

func (obj *AttributeDescriptor) SetIsArray(arrayFlag bool) {
	obj.isArray = arrayFlag
}

func (obj *AttributeDescriptor) SetIsEncrypted(encryptedFlag bool) {
	obj.isEncrypted = encryptedFlag
}

func (obj *AttributeDescriptor) SetName(attrName string) {
	obj.name = attrName
}

func (obj *AttributeDescriptor) SetSystemType(sysType types.TGSystemType) {
	obj.sysType = sysType
}

func (obj *AttributeDescriptor) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AttributeDescriptor:{")
	buffer.WriteString(fmt.Sprintf("SysType: %d", obj.sysType))
	buffer.WriteString(fmt.Sprintf(", AttributeId: %d", obj.attributeId))
	buffer.WriteString(fmt.Sprintf(", Name: %s", obj.name))
	buffer.WriteString(fmt.Sprintf(", AttrType: %s", types.GetAttributeTypeFromId(obj.attrType).GetTypeName()))
	buffer.WriteString(fmt.Sprintf(", IsArray: %+v", obj.isArray))
	buffer.WriteString(fmt.Sprintf(", Precision: %+v", obj.precision))
	buffer.WriteString(fmt.Sprintf(", Scale: %+v", obj.scale))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGAttributeDescriptor
/////////////////////////////////////////////////////////////////

// GetAttributeId returns the attributeId
func (obj *AttributeDescriptor) GetAttributeId() int64 {
	return obj.attributeId
}

// GetAttrType returns the type of Attribute Descriptor
func (obj *AttributeDescriptor) GetAttrType() int {
	return obj.attrType
}

// GetPrecision returns the precision for Attribute Descriptor of type Number. The default precision is 20
func (obj *AttributeDescriptor) GetPrecision() int16 {
	return obj.precision
}

// GetScale returns the scale for Attribute Descriptor of type Number. The default scale is 5
func (obj *AttributeDescriptor) GetScale() int16 {
	return obj.scale
}

// IsAttributeArray checks whether the AttributeType an array desc or not
func (obj *AttributeDescriptor) IsAttributeArray() bool {
	return obj.isArray
}

// IsEncrypted checks whether this attribute is Encrypted or not
func (obj *AttributeDescriptor) IsEncrypted() bool {
	return obj.isEncrypted
}

// SetPrecision sets the prevision for Attribute Descriptor of type Number
func (obj *AttributeDescriptor) SetPrecision(precision int16) {
	if obj.attrType == types.AttributeTypeNumber {
		obj.precision = precision
	}
}

// SetScale sets the scale for Attribute Descriptor of type Number
func (obj *AttributeDescriptor) SetScale(scale int16) {
	if obj.attrType == types.AttributeTypeNumber {
		obj.scale = scale
	}
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSystemObject
/////////////////////////////////////////////////////////////////

// GetName gets the system object's name
func (obj *AttributeDescriptor) GetName() string {
	return obj.name
}

// GetSystemType gets system object's type
func (obj *AttributeDescriptor) GetSystemType() types.TGSystemType {
	return types.SystemTypeAttributeDescriptor
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *AttributeDescriptor) ReadExternal(is types.TGInputStream) types.TGError {
	sType, err := is.(*iostream.ProtocolDataInputStream).ReadByte() // Read the sysobject desc field which should be 0 for attribute descriptor
	if err != nil {
		return err
	}
	if types.TGSystemType(sType) != types.SystemTypeAttributeDescriptor {
		// TODO: Revisit later - Do we need to throw exception is needed
		//errMsg := fmt.Sprintf("Attribute descriptor has invalid input stream value: %d", sType)
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		logger.Warning(fmt.Sprint("WARNING: AttributeDescriptor:ReadExternal - types.TGSystemType(sType) != types.SystemTypeAttributeDescriptor"))
	}
	attributeId, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read attributeId w/ Error: '%+v'", err.Error()))
		return err
	}
	attrName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read attrName w/ Error: '%+v'", err.Error()))
		return err
	}
	attrType, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read attrType w/ Error: '%+v'", err.Error()))
		return err
	}
	isArray, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read isArray w/ Error: '%+v'", err.Error()))
		return err
	}
	isEncrypted, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read isEncrypted w/ Error: '%+v'", err.Error()))
		return err
	}
	var precision, scale int16
	if attrType == types.AttributeTypeNumber {
		precision, err = is.(*iostream.ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read precision w/ Error: '%+v'", err.Error()))
			return err
		}
		scale, err = is.(*iostream.ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read scale w/ Error: '%+v'", err.Error()))
			return err
		}
	}
	obj.SetSystemType(types.TGSystemType(sType))
	obj.SetAttributeId(int64(attributeId))
	obj.SetName(attrName)
	obj.SetAttrType(int(attrType))
	obj.SetIsArray(isArray)
	obj.SetIsEncrypted(isEncrypted)
	obj.SetPrecision(precision)
	obj.SetScale(scale)
	logger.Log(fmt.Sprintf("Returning AttributeDescriptor:ReadExternal w/ NO error, for attrDesc: '%+v'", obj))
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *AttributeDescriptor) WriteExternal(os types.TGOutputStream) types.TGError {
	os.(*iostream.ProtocolDataOutputStream).WriteByte(int(types.SystemTypeAttributeDescriptor)) // SysObject desc attribute descriptor
	os.(*iostream.ProtocolDataOutputStream).WriteInt(int(obj.GetAttributeId()))
	err := os.(*iostream.ProtocolDataOutputStream).WriteUTF(obj.GetName())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:WriteExternal - unable to write attrDesc name w/ Error: '%+v'", err.Error()))
		return err
	}
	os.(*iostream.ProtocolDataOutputStream).WriteByte(obj.GetAttrType())
	os.(*iostream.ProtocolDataOutputStream).WriteBoolean(obj.IsAttributeArray())
	os.(*iostream.ProtocolDataOutputStream).WriteBoolean(obj.IsEncrypted())
	if obj.attrType == types.AttributeTypeNumber {
		os.(*iostream.ProtocolDataOutputStream).WriteShort(int(obj.GetPrecision()))
		os.(*iostream.ProtocolDataOutputStream).WriteShort(int(obj.GetScale()))
	}
	//logger.Log(fmt.Sprintf("Exported Attribute Descriptor object as '%+v' from byte format", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *AttributeDescriptor) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.sysType, obj.attributeId, obj.name, obj.attrType, obj.isArray, obj.precision, obj.scale)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *AttributeDescriptor) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.sysType, &obj.attributeId, &obj.name, &obj.attrType, &obj.isArray, &obj.precision, &obj.scale)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
