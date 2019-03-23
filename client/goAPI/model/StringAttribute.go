package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
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
 * File name: StringAttribute.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const MaxStringAttrLength = 1000 - 2 - 1 // 2 ==> size of int16 or short

type StringAttribute struct {
	*AbstractAttribute
}

// Create New Attribute Instance
func DefaultStringAttribute() *StringAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(StringAttribute{})

	newAttribute := StringAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewStringAttributeWithOwner(ownerEntity types.TGEntity) *StringAttribute {
	newAttribute := DefaultStringAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewStringAttribute(attrDesc *AttributeDescriptor) *StringAttribute {
	newAttribute := DefaultStringAttribute()
	newAttribute.attrDesc = attrDesc
	return newAttribute
}

func NewStringAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *StringAttribute {
	newAttribute := NewStringAttributeWithOwner(ownerEntity)
	newAttribute.attrDesc = attrDesc
	newAttribute.attrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for StringAttribute
/////////////////////////////////////////////////////////////////

func UtfLength(str string) int {
	strLen := len(str)
	utfLen := 0
	for i := 0; i < strLen; i++ {
		c := str[i]
		if c >= 0x0001 && c <= 0x007F {
			utfLen++
		} else if int(c) > 0x07FF {
			utfLen += 3
		} else {
			utfLen += 2
		}
	}
	return utfLen
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *StringAttribute) GetAttributeDescriptor() types.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *StringAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the name for this attribute as the most generic form
func (obj *StringAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *StringAttribute) GetOwner() types.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *StringAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *StringAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *StringAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *StringAttribute) SetOwner(ownerEntity types.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *StringAttribute) SetValue(value interface{}) types.TGError {
	if value == nil {
		obj.attrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.attrValue == value {
		return nil
	}

	s := value.(string)
	strLen := UtfLength(s)
	if strLen > MaxStringAttrLength {
		logger.Error(fmt.Sprint("ERROR: Returning StringAttribute:SetValue as strLen > MaxStringAttrLength"))
		errMsg := fmt.Sprintf("UTF length of String exceed the maximum string length supported for a String Attribute (%d > %d", strLen, MaxStringAttrLength)
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	obj.attrValue = value.(string)
	obj.setIsModified(true)
	return nil
}

// ReadValue reads the value from input stream
func (obj *StringAttribute) ReadValue(is types.TGInputStream) types.TGError {
	value, err := is.(*iostream.ProtocolDataInputStream).ReadUTF()
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprintf("StringAttribute::ReadValue - read value: '%+v'", value))
	obj.attrValue = value
	return nil
}

// WriteValue writes the value to output stream
func (obj *StringAttribute) WriteValue(os types.TGOutputStream) types.TGError {
	return os.(*iostream.ProtocolDataOutputStream).WriteUTF(obj.attrValue.(string))
}

func (obj *StringAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("StringAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *StringAttribute) ReadExternal(is types.TGInputStream) types.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *StringAttribute) WriteExternal(os types.TGOutputStream) types.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaler
/////////////////////////////////////////////////////////////////

func (obj *StringAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.attrDesc, obj.attrValue, obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning StringAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *StringAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.attrDesc, &obj.attrValue, &obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning StringAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
