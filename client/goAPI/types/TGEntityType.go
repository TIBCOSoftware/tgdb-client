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
 * File name: TGEntityType.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGEntityType interface {
	TGSystemObject
	// AddAttributeDescriptor add an attribute descriptor to the map
	AddAttributeDescriptor(attrName string, attrDesc TGAttributeDescriptor)
	// GetEntityTypeId gets Entity Type id
	GetEntityTypeId() int
	// DerivedFrom gets the parent Entity Type
	DerivedFrom() TGEntityType
	// GetAttributeDescriptor gets the attribute descriptor for the specified name
	GetAttributeDescriptor(attrName string) TGAttributeDescriptor
	// GetAttributeDescriptors returns a collection of attribute descriptors associated with this Entity Type
	GetAttributeDescriptors() []TGAttributeDescriptor
	// SetEntityTypeId sets Entity Type id
	SetEntityTypeId(eTypeId int)
	// SetName sets the system object's name
	SetName(eTypeName string)
	// SetSystemType sets system object's type
	SetSystemType(eSysType TGSystemType)
	// Additional Method to help debugging
	String() string
}
