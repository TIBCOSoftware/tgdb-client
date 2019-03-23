package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"reflect"
	"strconv"
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
 * File name: ByteAttribute.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type ByteAttribute struct {
	*AbstractAttribute
}

// Create New Attribute Instance
func DefaultByteAttribute() *ByteAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ByteAttribute{})

	newAttribute := ByteAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewByteAttributeWithOwner(ownerEntity types.TGEntity) *ByteAttribute {
	newAttribute := DefaultByteAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewByteAttribute(attrDesc *AttributeDescriptor) *ByteAttribute {
	newAttribute := DefaultByteAttribute()
	newAttribute.attrDesc = attrDesc
	return newAttribute
}

func NewByteAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *ByteAttribute {
	newAttribute := NewByteAttributeWithOwner(ownerEntity)
	newAttribute.attrDesc = attrDesc
	newAttribute.attrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for ByteAttribute
/////////////////////////////////////////////////////////////////

func (obj *ByteAttribute) SetByte(b uint8) {
	if !obj.IsNull() {
		return
	}
	obj.attrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *ByteAttribute) GetAttributeDescriptor() types.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *ByteAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the name for this attribute as the most generic form
func (obj *ByteAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *ByteAttribute) GetOwner() types.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *ByteAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *ByteAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *ByteAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *ByteAttribute) SetOwner(ownerEntity types.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *ByteAttribute) SetValue(value interface{}) types.TGError {
	logger.Log(fmt.Sprintf("ByteAttribute::SetValue trying to set attribute value '%+v' of type '%+v'", value, reflect.TypeOf(value).Kind()))
	if value == nil {
		//errMsg := fmt.Sprintf("Attribute value is required")
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		obj.attrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.attrValue == value {
		return nil
	}

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(value)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:SetValue - unable to encode attribute value"))
		errMsg := "ByteAttribute::SetValue - Unable to encode attribute value"
		return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	dec := gob.NewDecoder(&network)

	if reflect.TypeOf(value).Kind() != reflect.Bool &&
		reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.Uint &&
		reflect.TypeOf(value).Kind() != reflect.Uint8 &&
		reflect.TypeOf(value).Kind() != reflect.Uint16 &&
		reflect.TypeOf(value).Kind() != reflect.Uint32 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to ByteAttribute")
		return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.Bool {
		var v bool
		err = dec.Decode(&v)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:SetValue - unable to decode attribute value"))
			errMsg := "ByteAttribute::SetValue - Unable to decode attribute value"
			return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
		}
		if v {
			obj.SetByte(1)
		} else {
			obj.SetByte(0)
		}
	} else if reflect.TypeOf(value).Kind() == reflect.String ||
	   reflect.TypeOf(value).Kind() == reflect.Float32 ||
		reflect.TypeOf(value).Kind() == reflect.Float64 {
		v := value.(string)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to ByteAttribute")
			return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		v1, _ := strconv.Atoi(v)
		obj.SetByte(uint8(v1))
	} else {
		v := uint8(reflect.ValueOf(value).Uint())
		logger.Log(fmt.Sprintf("CharAttribute::SetValue finally trying to set attribute value '%+v' of type '%+v'", v, reflect.TypeOf(v).Kind()))
		obj.SetByte(v)
	}

	return nil
}

// ReadValue reads the value from input stream
func (obj *ByteAttribute) ReadValue(is types.TGInputStream) types.TGError {
	value, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:ReadValue w/ Error in reading value from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("ByteAttribute::ReadValue - read value: '%+v'", value))
	obj.attrValue = value
	return nil
}

// WriteValue writes the value to output stream
func (obj *ByteAttribute) WriteValue(os types.TGOutputStream) types.TGError {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(obj.attrValue)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:WriteValue - unable to encode attribute value"))
		errMsg := "AbstractAttribute::WriteExternal - Unable to encode attribute value"
		return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	dec := gob.NewDecoder(&network)
	var v byte
	err = dec.Decode(&v)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:WriteValue - unable to decode attribute value"))
		errMsg := "AbstractAttribute::WriteExternal - Unable to decode attribute value"
		return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	os.(*iostream.ProtocolDataOutputStream).WriteByte(int(v))
	return nil
}

func (obj *ByteAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ByteAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *ByteAttribute) ReadExternal(is types.TGInputStream) types.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *ByteAttribute) WriteExternal(os types.TGOutputStream) types.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ByteAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.attrDesc, obj.attrValue, obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ByteAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ByteAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.attrDesc, &obj.attrValue, &obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ByteAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
