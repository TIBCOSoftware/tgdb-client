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
	owner      types.TGEntity
	attrDesc   *AttributeDescriptor
	attrValue  interface{}
	isModified bool
}

// Create New Attribute Instance
func defaultNewAbstractAttribute() *AbstractAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AbstractAttribute{})

	newAttribute := AbstractAttribute{
		attrDesc:   DefaultAttributeDescriptor(),
		isModified: false,
	}
	return &newAttribute
}

func NewAbstractAttributeWithOwner(ownerEntity types.TGEntity) *AbstractAttribute {
	newAttribute := defaultNewAbstractAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewAbstractAttribute(attrDesc *AttributeDescriptor) *AbstractAttribute {
	newAttribute := defaultNewAbstractAttribute()
	newAttribute.attrDesc = attrDesc
	return newAttribute
}

func NewAbstractAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *AbstractAttribute {
	newAttribute := NewAbstractAttributeWithOwner(ownerEntity)
	newAttribute.attrDesc = attrDesc
	newAttribute.attrValue = value
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
	return obj.attrDesc
}

func (obj *AbstractAttribute) getIsModified() bool {
	return obj.isModified
}

func (obj *AbstractAttribute) getName() string {
	return obj.GetAttributeDescriptor().GetName()
}

func (obj *AbstractAttribute) getOwner() types.TGEntity {
	return obj.owner
}

func (obj *AbstractAttribute) getValue() interface{} {
	return obj.attrValue
}

func (obj *AbstractAttribute) isNull() bool {
	return obj.attrValue == nil
}

func (obj *AbstractAttribute) resetIsModified() {
	obj.isModified = false
}

func (obj *AbstractAttribute) setIsModified(flag bool) {
	obj.isModified = flag
}

func (obj *AbstractAttribute) setNull() {
	obj.attrValue = nil
}

func (obj *AbstractAttribute) setOwner(ownerEntity types.TGEntity) {
	obj.owner = ownerEntity
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

	obj.attrValue = value
	obj.isModified = true
	return nil
}

func (obj *AbstractAttribute) attributeToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("AbstractAttribute:{")
	//buffer.WriteString(fmt.Sprintf("Owner: %+v ", obj.owner))
	buffer.WriteString(fmt.Sprintf("AttrDesc: %s", obj.attrDesc.String()))
	buffer.WriteString(fmt.Sprintf(", AttrValue: %+v", obj.attrValue))
	buffer.WriteString(fmt.Sprintf(", IsModified: %+v", obj.isModified))
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

func AbstractAttributeReadDecrypted(obj types.TGAttribute, is types.TGInputStream) types.TGError {
	conn := obj.GetOwner().GetGraphMetadata().GetConnection()
	sysType := obj.GetOwner().GetEntityType().GetSystemType()
	decryptBuf := make([]byte, 0)
	switch sysType {
	case types.SystemTypeNode:
		entityId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:AbstractAttributeReadDecrypted w/ Error in reading entityId from message buffer: %s", err.Error()))
			return err
		}
		decryptBuf, err = conn.DecryptEntity(entityId)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:AbstractAttributeReadDecrypted w/ Error in conn.DecryptEntity(entityId): %s", err.Error()))
			return err
		}
	case types.SystemTypeEdge:
		encryptedBuf, err := is.(*iostream.ProtocolDataInputStream).ReadBytes()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:AbstractAttributeReadDecrypted w/ Error in reading encryptedBuf from message buffer: %s", err.Error()))
			return err
		}
		decryptBuf, err = conn.DecryptBuffer(encryptedBuf)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:AbstractAttributeReadDecrypted w/ Error in conn.DecryptBuffer(encryptedBuf): %s", err.Error()))
			return err
		}
	default:
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:AbstractAttributeReadDecrypted - Decryption not supported for entity type:'%d'", sysType))
		errMsg := fmt.Sprintf("Decryption not supported for entity type:'%d'", sysType)
		return exception.GetErrorByType(types.TGErrorIOException, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	value, err := ObjectFromByteArray(decryptBuf, obj.GetAttributeDescriptor().GetAttrType())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:AbstractAttributeReadDecrypted w/ Error in ObjectFromByteArray(): %s", err.Error()))
		return err
	}

	return obj.SetValue(value)
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
		obj.(*AbstractAttribute).attrValue = nil
		return nil
	}
	if obj.GetAttributeDescriptor().IsEncrypted() {
		return AbstractAttributeReadDecrypted(obj, is)
	}
	logger.Log(fmt.Sprint("Returning AbstractAttribute::AbstractAttributeReadExternal after reading the attribute value"))
	return obj.ReadValue(is)
}

func AbstractAttributeWriteEncrypted(obj types.TGAttribute, os types.TGOutputStream) types.TGError {
	buff, err := ObjectToByteArray(obj.GetValue(), obj.GetAttributeDescriptor().GetAttrType())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractAttribute:AbstractAttributeWriteEncrypted w/ Error in ObjectToByteArray"))
		return err
	}
	conn := obj.GetOwner().GetGraphMetadata().GetConnection()
	encryptedBuf, err := conn.EncryptEntity(buff)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractAttribute:AbstractAttributeWriteEncrypted w/ Error in conn.EncryptEntity(buff)"))
		return err
	}
	return os.(*iostream.ProtocolDataOutputStream).WriteBytes(encryptedBuf)
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
	if obj.GetAttributeDescriptor().IsEncrypted() {
		return AbstractAttributeWriteEncrypted(obj, os)
	}
	logger.Log(fmt.Sprint("Returning AbstractAttribute:AbstractAttributeWriteExternal after writing the attribute value"))
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
	_, err := fmt.Fprintln(&b, obj.owner, obj.attrDesc, obj.attrValue, obj.isModified)
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
	_, err := fmt.Fscanln(b, &obj.owner, &obj.attrDesc, &obj.attrValue, &obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
