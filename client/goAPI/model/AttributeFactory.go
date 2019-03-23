package model

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/logging"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * File name: AttributeFactory.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

var logger = logging.DefaultTGLogManager().GetLogger()

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttributeFactory
/////////////////////////////////////////////////////////////////

// CreateAttributeByType creates a new attribute based on the type specified
func CreateAttributeByType(attrTypeId int) (types.TGAttribute, types.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputAttrTypeId := attrTypeId

	// Use a switch case to switch between attribute types, if a type exist then error is nil (null)
	// Whenever new attribute type gets into the mix, just add a case below
	switch inputAttrTypeId {
	case types.AttributeTypeBoolean:
		return DefaultBooleanAttribute(), nil
	case types.AttributeTypeByte:
		return DefaultByteAttribute(), nil
	case types.AttributeTypeChar:
		return DefaultCharAttribute(), nil
	case types.AttributeTypeShort:
		return DefaultShortAttribute(), nil
	case types.AttributeTypeInteger:
		return DefaultIntegerAttribute(), nil
	case types.AttributeTypeLong:
		return DefaultLongAttribute(), nil
	case types.AttributeTypeFloat:
		return DefaultFloatAttribute(), nil
	case types.AttributeTypeDouble:
		return DefaultDoubleAttribute(), nil
	case types.AttributeTypeNumber:
		return DefaultNumberAttribute(), nil
	case types.AttributeTypeString:
		return DefaultStringAttribute(), nil
	case types.AttributeTypeDate:
		return DefaultTimestampAttribute(), nil
	case types.AttributeTypeTime:
		return DefaultTimestampAttribute(), nil
	case types.AttributeTypeTimeStamp:
		return DefaultTimestampAttribute(), nil
	case types.AttributeTypeBlob:
		return DefaultBlobAttribute(), nil
	case types.AttributeTypeClob:
		return DefaultClobAttribute(), nil
	case types.AttributeTypeInvalid:
		fallthrough
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Attribute Type '%s'", types.GetAttributeTypeFromId(inputAttrTypeId).GetTypeName())
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return nil, nil
}

// CreateAttribute creates a new attribute based on AttributeDescriptor
func CreateAttribute(attrDesc *AttributeDescriptor) (types.TGAttribute, types.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	attrTypeId := attrDesc.GetAttrType()
	inputAttrTypeId := attrTypeId

	var newAttribute types.TGAttribute
	// Use a switch case to switch between attribute types, if a type exist then error is nil (null)
	// Whenever new attribute type gets into the mix, just add a case below
	switch inputAttrTypeId {
	case types.AttributeTypeBoolean:
		// Execute Individual Attribute's method
		newAttribute = NewBooleanAttribute(attrDesc)
	case types.AttributeTypeByte:
		// Execute Individual Attribute's method
		newAttribute = NewByteAttribute(attrDesc)
	case types.AttributeTypeChar:
		// Execute Individual Attribute's method
		newAttribute = NewCharAttribute(attrDesc)
	case types.AttributeTypeShort:
		// Execute Individual Attribute's method
		newAttribute = NewShortAttribute(attrDesc)
	case types.AttributeTypeInteger:
		// Execute Individual Attribute's method
		newAttribute = NewIntegerAttribute(attrDesc)
	case types.AttributeTypeLong:
		// Execute Individual Attribute's method
		newAttribute = NewLongAttribute(attrDesc)
	case types.AttributeTypeFloat:
		// Execute Individual Attribute's method
		newAttribute = NewFloatAttribute(attrDesc)
	case types.AttributeTypeDouble:
		// Execute Individual Attribute's method
		newAttribute = NewDoubleAttribute(attrDesc)
	case types.AttributeTypeNumber:
		// Execute Individual Attribute's method
		newAttribute = NewNumberAttribute(attrDesc)
	case types.AttributeTypeString:
		// Execute Individual Attribute's method
		newAttribute = NewStringAttribute(attrDesc)
	case types.AttributeTypeDate:
		newAttribute = NewTimestampAttribute(attrDesc)
	case types.AttributeTypeTime:
		newAttribute = NewTimestampAttribute(attrDesc)
	case types.AttributeTypeTimeStamp:
		// Execute Individual Attribute's method
		newAttribute = NewTimestampAttribute(attrDesc)
	case types.AttributeTypeBlob:
		// Execute Individual Attribute's method
		newAttribute = NewBlobAttribute(attrDesc)
	case types.AttributeTypeClob:
		// Execute Individual Attribute's method
		newAttribute = NewClobAttribute(attrDesc)
	case types.AttributeTypeInvalid:
		fallthrough
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Attribute Type '%s'", types.GetAttributeTypeFromId(inputAttrTypeId).GetTypeName())
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return newAttribute, nil
}

// CreateAttributeWithDesc creates new attribute based on the owner and AttributeDescriptor
func CreateAttributeWithDesc(attrOwner types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) (types.TGAttribute, types.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	attrTypeId := attrDesc.GetAttrType()
	inputAttrTypeId := attrTypeId

	var newAttribute types.TGAttribute
	// Use a switch case to switch between attribute types, if a type exist then error is nil (null)
	// Whenever new attribute type gets into the mix, just add a case below
	switch inputAttrTypeId {
	case types.AttributeTypeBoolean:
		// Execute Individual Attribute's method
		newAttribute = NewBooleanAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeByte:
		// Execute Individual Attribute's method
		newAttribute = NewByteAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeChar:
		// Execute Individual Attribute's method
		newAttribute = NewCharAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeShort:
		// Execute Individual Attribute's method
		newAttribute = NewShortAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeInteger:
		// Execute Individual Attribute's method
		newAttribute = NewIntegerAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeLong:
		// Execute Individual Attribute's method
		newAttribute = NewLongAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeFloat:
		// Execute Individual Attribute's method
		newAttribute = NewFloatAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeDouble:
		// Execute Individual Attribute's method
		newAttribute = NewDoubleAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeNumber:
		// Execute Individual Attribute's method
		newAttribute = NewNumberAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeString:
		// Execute Individual Attribute's method
		newAttribute = NewStringAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeDate:
		newAttribute = NewTimestampAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeTime:
		newAttribute = NewTimestampAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeTimeStamp:
		// Execute Individual Attribute's method
		newAttribute = NewTimestampAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeBlob:
		// Execute Individual Attribute's method
		newAttribute = NewBlobAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeClob:
		// Execute Individual Attribute's method
		newAttribute = NewClobAttributeWithDesc(attrOwner, attrDesc, value)
	case types.AttributeTypeInvalid:
		fallthrough
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Attribute Type '%s'", types.GetAttributeTypeFromId(inputAttrTypeId).GetTypeName())
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return newAttribute, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute type
func GetAttributeDescriptor(attrTypeId int) types.TGAttributeDescriptor {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return nil
	}
	// Execute Individual Attribute's method
	return attr.GetAttributeDescriptor()
}

// GetIsModified checks whether the attribute of this type is modified or not
func GetIsModified(attrTypeId int) bool {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return false
	}
	// Execute Individual Attribute's method
	return attr.GetIsModified()
}

// ResetIsModified resets the IsModified flag of this attribute type - recursively, if needed
func ResetIsModified(attrTypeId int) {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return
	}
	// Execute Individual Attribute's method
	attr.ResetIsModified()
}

// GetName gets the name for this attribute type as the most generic form
func GetName(attrTypeId int) interface{} {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return false
	}
	// Execute Individual Attribute's method
	return attr.GetName()
}

// GetOwner gets owner Entity of this attribute type
func GetOwner(attrTypeId int) interface{} {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return nil
	}
	// Execute Individual Attribute's method
	return attr.GetOwner()
}

// GetValue gets the value for this attribute type as the most generic form
func GetValue(attrTypeId int) interface{} {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return false
	}
	// Execute Individual Attribute's method
	return attr.GetValue()
}

// IsNull checks whether the value of this attribute type is null or not
func IsNull(attrTypeId int) bool {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return false
	}
	// Execute Individual Attribute's method
	return attr.IsNull()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func SetOwner(attrTypeId int, attrOwner types.TGEntity) {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return
	}
	// Execute Individual Attribute's method
	attr.SetOwner(attrOwner)
}

// SetValue sets the value for this attribute type. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func SetValue(attrTypeId int, value interface{}) types.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.SetValue(value)
}

// ReadValue reads the value of this attribute type from input stream
func ReadValue(attrTypeId int, is types.TGInputStream) types.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.ReadValue(is)
}

// WriteValue writes the value of this attribute type to output stream
func WriteValue(attrTypeId int, os types.TGOutputStream) types.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.WriteValue(os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func ReadExternal(attrTypeId int, is types.TGInputStream) types.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.ReadExternal(is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func WriteExternal(attrTypeId int, os types.TGOutputStream) types.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.WriteExternal(os)
}
