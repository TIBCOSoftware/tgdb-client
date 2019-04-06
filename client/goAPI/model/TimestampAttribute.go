package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"reflect"
	"strings"
	"time"
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
 * File name: TimestampAttribute.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TimestampAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultTimestampAttribute() *TimestampAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TimestampAttribute{})

	newAttribute := TimestampAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewTimestampAttributeWithOwner(ownerEntity types.TGEntity) *TimestampAttribute {
	newAttribute := DefaultTimestampAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewTimestampAttribute(attrDesc *AttributeDescriptor) *TimestampAttribute {
	newAttribute := DefaultTimestampAttribute()
	newAttribute.attrDesc = attrDesc
	return newAttribute
}

func NewTimestampAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *TimestampAttribute {
	newAttribute := NewTimestampAttributeWithOwner(ownerEntity)
	newAttribute.attrDesc = attrDesc
	newAttribute.attrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for TimestampAttribute
/////////////////////////////////////////////////////////////////

func (obj *TimestampAttribute) SetCalendar(b time.Time) {
	if ! obj.IsNull() {
		return
	}
	obj.attrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *TimestampAttribute) GetAttributeDescriptor() types.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *TimestampAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the name for this attribute as the most generic form
func (obj *TimestampAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *TimestampAttribute) GetOwner() types.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *TimestampAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *TimestampAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *TimestampAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *TimestampAttribute) SetOwner(ownerEntity types.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *TimestampAttribute) SetValue(value interface{}) types.TGError {
	if value == nil {
		obj.attrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.attrValue == value {
		return nil
	}

	logger.Log(fmt.Sprintf("Input value '%+v' is of type: '%+v'\n", value, reflect.TypeOf(value).Kind()))
	if reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Int64 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprint("Failure to cast the attribute value to TimestampAttribute")
		return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := StringToCalendar(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprint("Failure to covert string to TimestampAttribute")
			return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
		}
		logger.Log(fmt.Sprintf("Transformed value '%+v' is of type: '%+v'\n", v, reflect.TypeOf(v).Kind()))
		obj.SetCalendar(v)
		return nil
	} else if reflect.TypeOf(value).Kind() == reflect.Int32 ||
		reflect.TypeOf(value).Kind() == reflect.Int64 {
		v := LongToCalendar(value.(int64))
		logger.Log(fmt.Sprintf("Transformed value '%+v' is of type: '%+v'\n", v, reflect.TypeOf(v).Kind()))
		obj.SetCalendar(v)
		return nil
	} else {
		obj.attrValue = value
		obj.setIsModified(true)
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *TimestampAttribute) ReadValue(is types.TGInputStream) types.TGError {
	var v time.Time
	var year, mon, dom, hr, min, sec, ms, tzType int
	era, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading era from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read era: '%+v'", era))

	yr, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading yr from message buffer"))
		return err
	}
	year = int(yr)
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read year: '%+v'", year))

	mth, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading mth from message buffer"))
		return err
	}
	mon = int(mth)
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read mth: '%+v'", mth))

	day, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading day from message buffer"))
		return err
	}
	dom = int(day)
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read dom: '%+v'", dom))

	hour, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading hour from message buffer"))
		return err
	}
	hr = int(hour)
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read hr: '%+v'", hr))

	mts, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading mts from message buffer"))
		return err
	}
	min = int(mts)
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read min: '%+v'", min))

	secs, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading secs from message buffer"))
		return err
	}
	sec = int(secs)
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read sec: '%+v'", sec))

	mSec, err := is.(*iostream.ProtocolDataInputStream).ReadUnsignedShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading mSec from message buffer"))
		return err
	}
	ms = int(mSec)
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read ms: '%+v'", ms))

	tz, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading tz from message buffer"))
		return err
	}
	tzType = int(tz)
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read tzType: '%+v'", tzType))

	if tzType != -1 {
		tzId, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading tzId from message buffer"))
			return err
		}
		logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read tzId: '%+v'", tzId))
	}

	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - attribute type: '%+v'", obj.attrDesc.GetAttrType()))
	switch obj.attrDesc.GetAttrType() {
	case types.AttributeTypeDate:
		v = time.Date(year, time.Month(mon), dom, 0, 0, 0, 0, time.Local)
		break
	case types.AttributeTypeTime:
		v = time.Date(1970, time.January, 01, hr, min, sec, ms*1000, time.Local)
		break
	case types.AttributeTypeTimeStamp:
		v = time.Date(year, time.Month(mon), dom, hr, min, sec, ms*1000, time.Local)
		break
	default:
		errMsg := fmt.Sprintf("Bad Descriptor: %s", string(obj.attrDesc.GetAttrType()))
		return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, "")
	}
	logger.Log(fmt.Sprintf("TimestampAttribute::ReadValue - read v: '%+v'", v))

	//obj.AttrValue = v.In(loc)		// TODO: Revisit later to use this once location/zone information is available
	obj.attrValue = v
	return nil
}

// WriteValue writes the value to output stream
func (obj *TimestampAttribute) WriteValue(os types.TGOutputStream) types.TGError {
	era := true // Corresponding to GregorianCalendar.AD = 1
	logger.Log(fmt.Sprintf("Object value '%+v' is of type: '%+v'\n", obj.GetValue(), reflect.TypeOf(obj.GetValue()).Kind()))
	v := obj.GetValue().(time.Time)
	yr, mth, day := v.Date()
	hr := v.Hour()
	min := v.Minute()
	sec := v.Second()
	msec := v.Nanosecond() / 1000
	switch obj.attrDesc.GetAttrType() {
	case types.AttributeTypeDate:
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(era)
		os.(*iostream.ProtocolDataOutputStream).WriteShort(yr)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(int(mth))
		os.(*iostream.ProtocolDataOutputStream).WriteByte(day)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(0)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(0)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(0)
		os.(*iostream.ProtocolDataOutputStream).WriteShort(0)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(TGNoZone)
		break
	case types.AttributeTypeTime:
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(era)
		os.(*iostream.ProtocolDataOutputStream).WriteShort(0)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(0)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(0)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(hr)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(min)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(sec)
		os.(*iostream.ProtocolDataOutputStream).WriteShort(msec)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(TGNoZone)
		break
	case types.AttributeTypeTimeStamp:
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(era)
		os.(*iostream.ProtocolDataOutputStream).WriteShort(yr)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(int(mth))
		os.(*iostream.ProtocolDataOutputStream).WriteByte(day)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(hr)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(min)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(sec)
		os.(*iostream.ProtocolDataOutputStream).WriteShort(msec)
		os.(*iostream.ProtocolDataOutputStream).WriteByte(TGNoZone)
		break
	default:
		errMsg := fmt.Sprintf("Bad Descriptor: %s", string(obj.attrDesc.GetAttrType()))
		return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, "")
	}
	return nil
}

func (obj *TimestampAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TimestampAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *TimestampAttribute) ReadExternal(is types.TGInputStream) types.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *TimestampAttribute) WriteExternal(os types.TGOutputStream) types.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TimestampAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.attrDesc, obj.attrValue, obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TimestampAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TimestampAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.attrDesc, &obj.attrValue, &obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TimestampAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
