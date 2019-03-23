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
 * File name: BooleanAttribute.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type BooleanAttribute struct {
	*AbstractAttribute
}

// Create New Attribute Instance
func DefaultBooleanAttribute() *BooleanAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(BooleanAttribute{})

	newAttribute := BooleanAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	newAttribute.attrValue = false
	return &newAttribute
}

func NewBooleanAttributeWithOwner(ownerEntity types.TGEntity) *BooleanAttribute {
	newAttribute := DefaultBooleanAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewBooleanAttribute(attrDesc *AttributeDescriptor) *BooleanAttribute {
	newAttribute := DefaultBooleanAttribute()
	newAttribute.attrDesc = attrDesc
	return newAttribute
}

func NewBooleanAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *BooleanAttribute {
	newAttribute := NewBooleanAttributeWithOwner(ownerEntity)
	newAttribute.attrDesc = attrDesc
	newAttribute.attrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for BooleanAttribute
/////////////////////////////////////////////////////////////////

func (obj *BooleanAttribute) SetBoolean(b bool) {
	if !obj.IsNull() && obj.attrValue.(bool) == b {
		return
	}
	obj.attrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *BooleanAttribute) GetAttributeDescriptor() types.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *BooleanAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the name for this attribute as the most generic form
func (obj *BooleanAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *BooleanAttribute) GetOwner() types.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *BooleanAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *BooleanAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *BooleanAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *BooleanAttribute) SetOwner(ownerEntity types.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *BooleanAttribute) SetValue(value interface{}) types.TGError {
	if value == nil {
		obj.attrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.attrValue == value {
		return nil
	}
	if reflect.TypeOf(value).Kind() != reflect.Bool &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning BooleanAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to BooleanAttribute")
		return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := strconv.ParseBool(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning BooleanAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to BooleanAttribute")
			return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetBoolean(v)
	} else {
		obj.SetBoolean(value.(bool))
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *BooleanAttribute) ReadValue(is types.TGInputStream) types.TGError {
	value, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading value from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("BooleanAttribute::ReadValue - read value: '%+v'", value))
	obj.attrValue = value
	return nil
}

// WriteValue writes the value to output stream
func (obj *BooleanAttribute) WriteValue(os types.TGOutputStream) types.TGError {
	if obj.attrValue == nil {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(false)
	} else {
		var network bytes.Buffer
		enc := gob.NewEncoder(&network)
		err := enc.Encode(obj.attrValue)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning BooleanAttribute:WriteValue - unable to encode attribute value w/ '%s'", err.Error()))
			errMsg := "AbstractAttribute::WriteExternal - Unable to encode attribute value"
			return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
		}
		dec := gob.NewDecoder(&network)
		var v bool
		err = dec.Decode(&v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning BooleanAttribute:WriteValue - unable to decode attribute value w/ '%s'", err.Error()))
			errMsg := "AbstractAttribute::WriteExternal - Unable to decode attribute value"
			return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
		}
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(v)
	}
	return nil
}

func (obj *BooleanAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("BooleanAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *BooleanAttribute) ReadExternal(is types.TGInputStream) types.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *BooleanAttribute) WriteExternal(os types.TGOutputStream) types.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *BooleanAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.attrDesc, obj.attrValue, obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BooleanAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *BooleanAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.attrDesc, &obj.attrValue, &obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BooleanAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
