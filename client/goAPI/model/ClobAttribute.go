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
 * File name: ClobAttribute.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

//var gUniqueId = NewAtomicLong(0)

type ClobAttribute struct {
	*BlobAttribute
}

// Create New Attribute Instance
func DefaultClobAttribute() *ClobAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ClobAttribute{})

	newAttribute := ClobAttribute{
		BlobAttribute: DefaultBlobAttribute(),
	}
	return &newAttribute
}

func NewClobAttributeWithOwner(ownerEntity types.TGEntity) *ClobAttribute {
	newAttribute := DefaultClobAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewClobAttribute(attrDesc *AttributeDescriptor) *ClobAttribute {
	newAttribute := DefaultClobAttribute()
	newAttribute.attrDesc = attrDesc
	return newAttribute
}

func NewClobAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *ClobAttribute {
	newAttribute := NewClobAttributeWithOwner(ownerEntity)
	newAttribute.attrDesc = attrDesc
	newAttribute.attrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for ClobAttribute
/////////////////////////////////////////////////////////////////

func (obj *ClobAttribute) getValueAsBytes() ([]byte, types.TGError) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(obj.attrValue)
	if err != nil {
		errMsg := "ClobAttribute::getValueAsBytes - Unable to encode attribute value"
		return nil, exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	dec := gob.NewDecoder(&network)
	var v []byte
	err = dec.Decode(&v)
	if err != nil {
		errMsg := "ClobAttribute::getValueAsBytes - Unable to decode attribute value"
		return nil, exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	return v, nil
}

func (obj *ClobAttribute) SetCharBuffer(b string) {
	if !obj.IsNull() && obj.attrValue == b {
		return
	}
	obj.attrValue = []byte(b)
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *ClobAttribute) GetAttributeDescriptor() types.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *ClobAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the name for this attribute as the most generic form
func (obj *ClobAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *ClobAttribute) GetOwner() types.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *ClobAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *ClobAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *ClobAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *ClobAttribute) SetOwner(ownerEntity types.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *ClobAttribute) SetValue(value interface{}) types.TGError {
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
	// TODO: Revisit later
	//if (value == null)
	//{
	//	this.value = value;
	//	setModified();
	//	return;
	//}
	//else if (value instanceof char[]) {
	//	setCharBuffer(CharBuffer.wrap((char[])value));
	//}
	//else if (value instanceof CharBuffer) {
	//	setCharBuffer((CharBuffer) value);
	//}
	//else if (value instanceof CharSequence) {
	//	setCharBuffer(CharBuffer.wrap((CharSequence)value));
	//}
	//else {
	//	super.setValue(value);
	//}

	obj.SetCharBuffer(value.(string))
	return nil
}

// ReadValue reads the value from input stream
func (obj *ClobAttribute) ReadValue(is types.TGInputStream) types.TGError {
	entityId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ClobAttribute:ReadValue w/ Error in reading entityId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("BlobAttribute::ReadValue - read entityId: '%+v'", entityId))
	obj.entityId = entityId
	obj.isCached = false
	return nil
}

// WriteValue writes the value to output stream
func (obj *ClobAttribute) WriteValue(os types.TGOutputStream) types.TGError {
	if obj.attrValue == nil {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(false)
	} else {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(true)
		v, err := obj.getValueAsBytes()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ClobAttribute:WriteValue - Unable to decode attribute value w/ Error: '%s'", err.Error()))
			errMsg := "ClobAttribute::WriteValue - Unable to decode attribute value"
			return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
		err = os.(*iostream.ProtocolDataOutputStream).WriteBytes(v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ClobAttribute:WriteValue - Unable to write attribute value w/ Error: '%s'", err.Error()))
			errMsg := "ClobAttribute::WriteValue - Unable to write attribute value"
			return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	}
	return nil
}

func (obj *ClobAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ClobAttribute:{")
	buffer.WriteString(fmt.Sprintf("EntityId: %+v", obj.entityId))
	buffer.WriteString(fmt.Sprintf(", IsCached: %+v", obj.isCached))
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

//@Override
//public char[] getAsChars() throws TGException {
//	return getAsChars("UTF-8");
//}
//
//@Override
//public CharBuffer getAsCharBuffer() throws TGException {
//	return getAsCharBuffer("UTF-8");
//}
//
//public char[] getAsChars(String encoding) throws TGException
//{
//	CharBuffer cb = getAsCharBuffer(encoding);
//	return cb.array();
//}
//
//public CharBuffer getAsCharBuffer(String encoding) throws TGException {
//	ByteBuffer bb = getAsByteBuffer();
//	Charset cs = Charset.forName(encoding);
//	CharBuffer cb = cs.decode(bb);
//	return cb;
//}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *ClobAttribute) ReadExternal(is types.TGInputStream) types.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *ClobAttribute) WriteExternal(os types.TGOutputStream) types.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ClobAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.attrDesc, obj.attrValue, obj.isModified, obj.entityId, obj.isCached)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ClobAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ClobAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.attrDesc, &obj.attrValue, &obj.isModified, &obj.entityId, &obj.isCached)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ClobAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
