package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"log"
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
 * File name: AbstractEntity.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	DATE_ONLY    = 0
	TIME_ONLY    = 1
	TIMESTAMP    = 2
	TGNoZone     = -1
	TGZoneOffset = 0
	TGZoneId     = 1
	TGZoneName   = 2
)

//var gLogger = TGLogManager.getInstance().getLogger()

type AbstractAttribute struct {
	Owner      types.TGEntity
	AttrDesc   *AttributeDescriptor
	AttrValue  interface{}
	IsModified bool
}

// Create New Attribute Instance
func defaultNewAbstractAttribute() *AbstractAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AbstractAttribute{})

	newAttribute := AbstractAttribute{
		AttrDesc:   DefaultAttributeDescriptor(),
		IsModified: false,
	}
	return &newAttribute
}

func NewAbstractAttributeWithOwner(ownerEntity types.TGEntity) *AbstractAttribute {
	newAttribute := defaultNewAbstractAttribute()
	newAttribute.Owner = ownerEntity
	return newAttribute
}

func NewAbstractAttribute(attrDesc *AttributeDescriptor) *AbstractAttribute {
	newAttribute := defaultNewAbstractAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewAbstractAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *AbstractAttribute {
	newAttribute := NewAbstractAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Private functions for AbstractAttribute / TGAttribute
/////////////////////////////////////////////////////////////////

// interfaceEncode encodes the interface value into the encoder.
func interfaceEncode(enc *gob.Encoder, p types.TGAttribute) types.TGError {
	// The encode will fail unless the concrete type has been
	// registered. We registered it in the calling function.

	// Pass pointer to interface so Encode sees (and hence sends) a value of
	// interface type. If we passed p directly it would see the concrete type instead.
	// See the blog post, "The Laws of Reflection" for background.
	err := enc.Encode(&p)
	if err != nil {
		log.Fatal("encode:", err)
		errMsg := "Unable to encode interface"
		return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, "")
	}
	return nil
}

// interfaceDecode decodes the next interface value from the stream and returns it.
func interfaceDecode(dec *gob.Decoder) types.TGAttribute {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.
	var p types.TGAttribute
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal("decode:", err)
	}
	return p
}

func (obj *AbstractAttribute) getAttributeDescriptor() *AttributeDescriptor {
	return obj.AttrDesc
}

func (obj *AbstractAttribute) getIsModified() bool {
	return obj.IsModified
}

func (obj *AbstractAttribute) getName() string {
	return obj.GetAttributeDescriptor().GetName()
}

func (obj *AbstractAttribute) getOwner() types.TGEntity {
	return obj.Owner
}

func (obj *AbstractAttribute) getValue() interface{} {
	return obj.AttrValue
}

func (obj *AbstractAttribute) isNull() bool {
	return obj.AttrValue == nil
}

func (obj *AbstractAttribute) resetIsModified() {
	obj.IsModified = false
}

func (obj *AbstractAttribute) setIsModified(flag bool) {
	obj.IsModified = flag
}

func (obj *AbstractAttribute) setOwner(ownerEntity types.TGEntity) {
	obj.Owner = ownerEntity
}

func (obj *AbstractAttribute) setValue(value interface{}) types.TGError {
	if value == nil {
		errMsg := fmt.Sprintf("Attribute value is required")
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if !obj.isNull() && obj.getValue() == value {
		return nil
	}

	//var precision, scale int16
	//if obj.AttrDesc.GetAttrType() == types.AttributeTypeNumber {
	//	precision = obj.AttrDesc.GetPrecision()
	//	scale = obj.AttrDesc.GetScale()
	//}
	//obj.setValueWithPrecisionAndScale(value, precision, scale)

	obj.AttrValue = value
	obj.IsModified = true
	return nil
}

func (obj *AbstractAttribute) attributeToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("AbstractAttribute:{")
	//buffer.WriteString(fmt.Sprintf("Owner: %+v ", obj.Owner))
	buffer.WriteString(fmt.Sprintf("AttrDesc: %s", obj.AttrDesc.String()))
	buffer.WriteString(fmt.Sprintf(", AttrValue: %+v", obj.AttrValue))
	buffer.WriteString(fmt.Sprintf(", IsModified: %+v", obj.IsModified))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Helper functions for AbstractAttribute
/////////////////////////////////////////////////////////////////

func (obj *AbstractAttribute) SetIsModified(flag bool) {
	obj.setIsModified(flag)
}

func ReadExternalForEntity(owner types.TGEntity, is types.TGInputStream) (types.TGAttribute, types.TGError) {
	logger.Log(fmt.Sprint("Entering AbstractAttribute:ReadExternalForEntity"))
	attrId, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		return nil, err
	}
	logger.Log(fmt.Sprintf("AbstractAttribute::ReadExternalForEntity - read attrId: '%+v'", attrId))
	gmd := owner.(*AbstractEntity).GetGraphMetadata()
	if gmd == nil {
		errMsg := fmt.Sprintf("Invalid graph meta data associated with owner: '%+v'", owner.(*AbstractEntity))
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	//logger.Log(fmt.Sprintf("AbstractAttribute::ReadExternalForEntity - read gmd: '%+v'", gmd))
	attrDesc, err := gmd.GetAttributeDescriptorById(int64(attrId))
	if err != nil {
		return nil, err
	}
	logger.Log(fmt.Sprintf("AbstractAttribute::ReadExternalForEntity - read attrDesc: '%+v'", attrDesc))
	if attrDesc == nil {
		errMsg := fmt.Sprintf("Invalid attributeId:'%d' encountered while deserialized", attrId)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	newAttr, err := CreateAttributeWithDesc(owner, attrDesc.(*AttributeDescriptor), nil)
	if err != nil {
		return nil, err
	}
	logger.Log(fmt.Sprintf("AbstractAttribute::ReadExternalForEntity - created new attribute newAttr: '%+v'", newAttr))
	err = newAttr.ReadExternal(is)
	if err != nil {
		return nil, err
	}
	logger.Log(fmt.Sprintf("AbstractAttribute::ReadExternalForEntity - updated newAttr from stream: '%+v'", newAttr))
	return newAttr, nil
}

func AbstractAttributeReadExternal(obj types.TGAttribute, is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering AbstractAttribute:EntityTypeReadExternal"))
	// We have already read the AttributeId, so no need to read it.
	isNull, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractAttribute:AbstractAttributeReadExternal w/ Error in reading isNull from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("AbstractAttribute::AbstractAttributeReadExternal - read isnull: '%+v'", isNull))
	if isNull {
		obj.(*AbstractAttribute).AttrValue = nil
		return nil
	}
	logger.Log(fmt.Sprint("Returning AbstractAttribute::AbstractAttributeReadExternal after reading the attribute value"))
	return obj.ReadValue(is)
}

func AbstractAttributeWriteExternal(obj types.TGAttribute, os types.TGOutputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering AbstractAttribute:AbstractAttributeWriteExternal"))
	attrId := obj.GetAttributeDescriptor().GetAttributeId()
	// Null attribute is not allowed during entity creation
	os.(*iostream.ProtocolDataOutputStream).WriteInt(int(attrId))
	os.(*iostream.ProtocolDataOutputStream).WriteBoolean(obj.IsNull())
	if obj.IsNull() {
		return nil
	}

	logger.Log(fmt.Sprint("Entering AbstractAttribute:AbstractAttributeWriteExternal after writing the attribute value"))
	return obj.WriteValue(os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *AbstractAttribute) GetAttributeDescriptor() types.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *AbstractAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the name for this attribute as the most generic form
func (obj *AbstractAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *AbstractAttribute) GetOwner() types.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *AbstractAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *AbstractAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *AbstractAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *AbstractAttribute) SetOwner(ownerEntity types.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *AbstractAttribute) SetValue(value interface{}) types.TGError {
	return obj.setValue(value)
}

// ReadValue reads the value from input stream
func (obj *AbstractAttribute) ReadValue(is types.TGInputStream) types.TGError {
	return nil
}

// WriteValue writes the value to output stream
func (obj *AbstractAttribute) WriteValue(os types.TGOutputStream) types.TGError {
	return nil
}

func (obj *AbstractAttribute) String() string {
	return obj.attributeToString()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *AbstractAttribute) ReadExternal(is types.TGInputStream) types.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *AbstractAttribute) WriteExternal(os types.TGOutputStream) types.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *AbstractAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.Owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *AbstractAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.Owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
