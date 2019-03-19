package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"reflect"
	"strings"
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
 * File name: BlobAttribute.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

//var gUniqueId = NewAtomicLong(0)
var UniqueId int64

type BlobAttribute struct {
	*AbstractAttribute
	EntityId int64
	IsCached bool
}

// Create NewTGDecimal Attribute Instance
func DefaultBlobAttribute() *BlobAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(BlobAttribute{})

	newAttribute := BlobAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
		IsCached:          false,
	}
	newAttribute.EntityId = atomic.AddInt64(&UniqueId, 1)
	newAttribute.AttrValue = []byte{}
	return &newAttribute
}

func NewBlobAttributeWithOwner(ownerEntity types.TGEntity) *BlobAttribute {
	newAttribute := DefaultBlobAttribute()
	newAttribute.Owner = ownerEntity
	return newAttribute
}

func NewBlobAttribute(attrDesc *AttributeDescriptor) *BlobAttribute {
	newAttribute := DefaultBlobAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewBlobAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *BlobAttribute {
	newAttribute := NewBlobAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for BlobAttribute
/////////////////////////////////////////////////////////////////

func (obj *BlobAttribute) getValueAsBytes() ([]byte, types.TGError) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(obj.AttrValue)
	if err != nil {
		errMsg := "BlobAttribute::getValueAsBytes - Unable to encode attribute value"
		return nil, exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	dec := gob.NewDecoder(&network)
	var v []byte
	err = dec.Decode(&v)
	if err != nil {
		errMsg := "BlobAttribute::getValueAsBytes - Unable to decode attribute value"
		return nil, exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	return v, nil
}

func (obj *BlobAttribute) SetBlob(b []byte) {
	obj.AttrValue = b
	obj.setIsModified(true)
}

func (obj *BlobAttribute) GetAsBytes() []byte {
	if obj.EntityId < 0 || obj.IsCached {
		ba, _ := obj.getValueAsBytes()
		return ba
	}
	// TODO: Revisit later once connection.GetLargeObjectAsBytes() is implemented
	gmd := obj.GetOwner().GetGraphMetadata()
	conn := gmd.(*GraphMetadata).GetConnection()
	v, _ := conn.GetLargeObjectAsBytes(obj.EntityId)
	obj.AttrValue = v
	obj.IsCached = true
	return obj.AttrValue.([]byte)
}

func (obj *BlobAttribute) GetAsByteBuffer() *bytes.Buffer {
	buf := obj.GetAsBytes()
	return bytes.NewBuffer(buf)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *BlobAttribute) GetAttributeDescriptor() types.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *BlobAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the name for this attribute as the most generic form
func (obj *BlobAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *BlobAttribute) GetOwner() types.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *BlobAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *BlobAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *BlobAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *BlobAttribute) SetOwner(ownerEntity types.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *BlobAttribute) SetValue(value interface{}) types.TGError {
	if value == nil {
		//errMsg := fmt.Sprintf("Attribute value is required")
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() {
		return nil
	}

	if reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to BlobAttribute")
		return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.Float32 {
		v, err := utils.BigDecimalToByteArray(value.(float32))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:SetValue - unable to extract attribute value in float format/type"))
			errMsg := fmt.Sprintf("Failure to covert float to BlobAttribute")
			return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetBlob(v)
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		v, err := utils.DoubleToByteArray(value.(float64))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:SetValue - unable to extract attribute value in double format/type"))
			errMsg := fmt.Sprintf("Failure to covert double to BlobAttribute")
			return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetBlob(v)
	} else if reflect.TypeOf(value).Kind() == reflect.String {
		v := []byte(value.(string))
		obj.SetBlob(v)
	} else {
		obj.AttrValue = value
		obj.setIsModified(true)
	}
	// TODO: Revisit later - use retrospection via go/constant package
	//if (value instanceof byte[]) {
	//    this.value = value;
	//}
	//else if (value instanceof String) {
	//    this.value = ((String)value).getBytes();
	//}
	//else if (value instanceof BigDecimal) {
	//    this.value = ConversionUtils.bigDecimal2ByteArray(BigDecimal.class.cast(value));
	//}
	//else if (value instanceof ByteBuffer) {
	//    this.value = ((ByteBuffer)value).array();
	//}
	//else if (value instanceof InputStream) {
	//    this.value = ConversionUtils.inputStream2ByteArray((InputStream) value);
	//}
	//else {
	//    throw new TGErrorTypeCoercionNotSupported(TGAttributeType.AttributeTypeBlob, value.getClass().getSimpleName());
	//}

	return nil
}

// ReadValue reads the value from input stream
func (obj *BlobAttribute) ReadValue(is types.TGInputStream) types.TGError {
	entityId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:ReadValue w/ Error in reading entityId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("BlobAttribute::ReadValue - read entityId: '%+v'", entityId))
	obj.EntityId = entityId
	obj.IsCached = false
	return nil
}

// WriteValue writes the value to output stream
func (obj *BlobAttribute) WriteValue(os types.TGOutputStream) types.TGError {
	os.(*iostream.ProtocolDataOutputStream).WriteLong(obj.EntityId)
	if obj.AttrValue == nil {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(false)
	} else {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(true)
		v, err := obj.getValueAsBytes()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning BlobAttribute:WriteValue - Unable to decode attribute value w/ Error: '%s'", err.Error()))
			errMsg := "BlobAttribute::WriteValue - Unable to decode attribute value"
			return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
		err = os.(*iostream.ProtocolDataOutputStream).WriteBytes(v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning BlobAttribute:WriteValue - Unable to write attribute value w/ Error: '%s'", err.Error()))
			errMsg := "BlobAttribute::WriteValue - Unable to write attribute value"
			return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	}
	return nil
}

func (obj *BlobAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("BlobAttribute:{")
	buffer.WriteString(fmt.Sprintf("EntityId: %+v", obj.EntityId))
	buffer.WriteString(fmt.Sprintf(", IsCached: %+v", obj.IsCached))
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *BlobAttribute) ReadExternal(is types.TGInputStream) types.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *BlobAttribute) WriteExternal(os types.TGOutputStream) types.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *BlobAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.Owner, obj.AttrDesc, obj.AttrValue, obj.IsModified, obj.EntityId, obj.IsCached)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BlobAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *BlobAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.Owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified, &obj.EntityId, &obj.IsCached)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BlobAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
