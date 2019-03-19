package types

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
 * File name: TGAttribute.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// An Attribute is simple scalar value that is associated with an Entity.
type TGAttribute interface {
	TGSerializable
	// GetAttributeDescriptor returns the TGAttributeDescriptor for this attribute
	GetAttributeDescriptor() TGAttributeDescriptor
	// GetIsModified checks whether the attribute modified or not
	GetIsModified() bool
	// GetName gets the name for this attribute as the most generic form
	GetName() string
	// GetOwner gets owner Entity of this attribute
	GetOwner() TGEntity
	// GetValue gets the value for this attribute as the most generic form
	GetValue() interface{}
	// IsNull checks whether the attribute value is null or not
	IsNull() bool
	// ResetIsModified resets the IsModified flag - recursively, if needed
	ResetIsModified()
	// SetOwner sets the owner entity - Need this indirection to traverse the chain
	SetOwner(attrOwner TGEntity)
	// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
	// If the object is Null, then the object is explicitly set, but no value is provided.
	SetValue(value interface{}) TGError
	// ReadValue reads the attribute value from input stream
	ReadValue(is TGInputStream) TGError
	// WriteValue writes the attribute value to output stream
	WriteValue(os TGOutputStream) TGError
	// Additional Method to help debugging
	String() string
}
