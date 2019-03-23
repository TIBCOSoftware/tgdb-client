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
 * File name: NumberAttribute.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type NumberAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultNumberAttribute() *NumberAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(NumberAttribute{})

	newAttribute := NumberAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	newAttribute.GetAttributeDescriptor().SetPrecision(20)
	newAttribute.GetAttributeDescriptor().SetScale(5)
	return &newAttribute
}

func NewNumberAttributeWithOwner(ownerEntity types.TGEntity) *NumberAttribute {
	newAttribute := DefaultNumberAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewNumberAttribute(attrDesc *AttributeDescriptor) *NumberAttribute {
	newAttribute := DefaultNumberAttribute()
	newAttribute.attrDesc = attrDesc
	return newAttribute
}

func NewNumberAttributeWithDesc(ownerEntity types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *NumberAttribute {
	newAttribute := NewNumberAttributeWithOwner(ownerEntity)
	newAttribute.attrDesc = attrDesc
	newAttribute.attrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for NumberAttribute
/////////////////////////////////////////////////////////////////

func (obj *NumberAttribute) SetDecimal(b utils.TGDecimal, precision, scale int) {
	logger.Log(fmt.Sprintf("Inside NumberAttribute::SetDecimal about to set value '%+v' w/ precision:scale %d/%d", b, precision, scale))
	if !obj.IsNull() && obj.attrValue == b {
		return
	}
	obj.GetAttributeDescriptor().SetPrecision(int16(precision))
	obj.GetAttributeDescriptor().SetScale(int16(scale))
	obj.attrValue = b.String() //strings.Replace(b.String(), ".", "", -1)
	logger.Log(fmt.Sprintf("Inside NumberAttribute::SetDecimal attrValue is '%+v'", obj.attrValue))
	obj.setIsModified(true)
}

//func roundU(val float64) int {
//	if val > 0 {
//		return int(val + 1.0)
//	}
//	return int(val)
//}
//
//func roundD(val float64) int {
//	if val < 0 {
//		return int(val - 1.0)
//	}
//	return int(val)
//}
//
//func round(val float64) int {
//	if val < 0 {
//		return int(val - 0.5)
//	}
//	return int(val + 0.5)
//}
//
//func round1(num float64) int {
//	return int(num + math.Copysign(0.5, num))
//}
//
//func toFixed(num float64, precision int) float64 {
//	output := math.Power(10, float64(precision))
//	return float64(round(num*output)) / output
//}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *NumberAttribute) GetAttributeDescriptor() types.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *NumberAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the name for this attribute as the most generic form
func (obj *NumberAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *NumberAttribute) GetOwner() types.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *NumberAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *NumberAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *NumberAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *NumberAttribute) SetOwner(ownerEntity types.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *NumberAttribute) SetValue(value interface{}) types.TGError {
	logger.Log(fmt.Sprintf("Inside NumberAttribute::SetValue about to set value '%+v' of kind '%+v'", value, reflect.TypeOf(value).Kind()))
	if value == nil {
		obj.attrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.attrValue == value {
		return nil
	}

	precision := 20 	// Default Precision
	scale := 5			// Default Scale

	if 	reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Int64 &&
		reflect.TypeOf(value).Kind() != reflect.Struct &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to NumberAttribute")
		return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v1 := value.(string)
		parts := strings.Split(v1, ".")
		if len(parts) == 1 {
			// There is no decimal point, we can just parse the original string as an int
			scale = 0
			precision = len(parts[0])
		} else if len(parts) == 2 {
			scale = len(parts[1])
			precision = len(parts[0]) + len(parts[1])
		}
		//v, err := strconv.ParseFloat(v1, 64)
		v, err := utils.NewTGDecimalFromString(v1)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to NumberAttribute")
			return exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetDecimal(v, precision, scale)
	} else if reflect.TypeOf(value).Kind() == reflect.Int32 ||
		reflect.TypeOf(value).Kind() == reflect.Int64 {
		v1 := value.(int64)
		scale = 0
		precision = len(strconv.Itoa(int(v1)))
		obj.SetDecimal(utils.NewTGDecimal(value.(int64),0), precision, scale)
	} else if reflect.TypeOf(value).Kind() == reflect.Float32 ||
		reflect.TypeOf(value).Kind() == reflect.Float64 {
		v1 := value.(float64)
		v2 := fmt.Sprintf("%f", v1)
		logger.Log(fmt.Sprintf("Inside NumberAttribute::SetValue - read v1: '%f' and v2: '%s'", v1, v2))
		parts := strings.Split(v2, ".")
		if len(parts) == 1 {
			// There is no decimal point, we can just parse the original string as an int
			scale = 0
			precision = len(parts[0])
		} else if len(parts) == 2 {
			scale = len(parts[1])
			precision = len(parts[0]) + len(parts[1])
		}
		obj.SetDecimal(utils.NewTGDecimalFromFloatWithExponent(v1, 0), precision, scale)
	} else {
		v1 := value.(utils.TGDecimal)
		parts := strings.Split(v1.String(), ".")
		if len(parts) == 1 {
			// There is no decimal point, we can just parse the original string as an int
			scale = 0
			precision = len(parts[0])
		} else if len(parts) == 2 {
			scale = len(parts[1])
			precision = len(parts[0]) + len(parts[1])
		}
		obj.SetDecimal(value.(utils.TGDecimal), precision, scale)
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *NumberAttribute) ReadValue(is types.TGInputStream) types.TGError {
	precision, err := is.(*iostream.ProtocolDataInputStream).ReadShort() // precision
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:ReadValue w/ Error in reading precision from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside NumberAttribute::ReadValue - read precision: '%+v'", precision))

	scale, err := is.(*iostream.ProtocolDataInputStream).ReadShort() // scale
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:ReadValue w/ Error in reading scale from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside NumberAttribute::ReadValue - read scale: '%+v'", scale))

	bdStr, err := is.(*iostream.ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:ReadValue w/ Error in reading bdStr from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside NumberAttribute::ReadValue - read bdStr: '%+v'", bdStr))
	obj.attrValue, _ = utils.NewTGDecimalFromString(bdStr)
	return nil
}

// WriteValue writes the value to output stream
func (obj *NumberAttribute) WriteValue(os types.TGOutputStream) types.TGError {
	os.(*iostream.ProtocolDataOutputStream).WriteShort(int(obj.GetAttributeDescriptor().GetPrecision()))
	os.(*iostream.ProtocolDataOutputStream).WriteShort(int(obj.GetAttributeDescriptor().GetScale()))
	logger.Log(fmt.Sprintf("Inside NumberAttribute::WriteValue attrValue is '%+v'", obj.attrValue))
	//dValue := reflect.ValueOf(obj.AttrValue).Float()
	//strValue := strconv.FormatFloat(dValue, 'f', int(obj.GetAttributeDescriptor().GetPrecision()), 64)
	//newStr := strings.Replace(strValue, ".", "", -1)
	dValue := obj.attrValue.(string)
	newStr := dValue	//strings.Replace(dValue, ".", "", -1)
	logger.Log(fmt.Sprintf("Inside NumberAttribute::WriteValue newStr is '%+v'", newStr))
	return os.(*iostream.ProtocolDataOutputStream).WriteUTF(newStr)
}

func (obj *NumberAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("NumberAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *NumberAttribute) ReadExternal(is types.TGInputStream) types.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *NumberAttribute) WriteExternal(os types.TGOutputStream) types.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *NumberAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.attrDesc, obj.attrValue, obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NumberAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *NumberAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.attrDesc, &obj.attrValue, &obj.isModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NumberAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
