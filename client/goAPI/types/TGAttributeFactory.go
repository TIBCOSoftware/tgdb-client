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
 * File name: TGAttributeFactory.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// An AttributeFactory is a convenient mechanism to create Attribute(s) of various types
type TGAttributeFactory interface {
	TGSerializable
	// CreateAttributeByType creates a new attribute based on the type specified
	CreateAttributeByType(attrTypeId int) (TGAttribute, TGError)
	// CreateAttribute creates a new attribute based on AttributeDescriptor
	CreateAttribute(attrDesc TGAttributeDescriptor) (TGAttribute, TGError)
	// CreateAttributeWithDesc creates new attribute based on the owner and AttributeDescriptor
	CreateAttributeWithDesc(attrOwner TGEntity, attrDesc TGAttributeDescriptor, value interface{}) (TGAttribute, TGError)
}
