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
 * File name: TGAttributeDescriptor.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// An AttributeDescriptor is Basic definition of an Attribute
type TGAttributeDescriptor interface {
	TGSystemObject
	// GetAttributeId returns the attributeId
	GetAttributeId() int64
	// GetAttrType returns the type of Attribute Descriptor
	GetAttrType() int
	// GetPrecision returns the precision for Attribute Descriptor of type Number. The default precision is 20
	GetPrecision() int16
	// GetScale returns the scale for Attribute Descriptor of type Number. The default scale is 5
	GetScale() int16
	// IsAttributeArray checks whether the AttributeType an array desc or not
	IsAttributeArray() bool
	// IsEncrypted checks whether this attribute is Encrypted or not
	IsEncrypted() bool
	// SetPrecision sets the prevision for Attribute Descriptor of type Number
	SetPrecision(precision int16)
	// SetScale sets the scale for Attribute Descriptor of type Number
	SetScale(scale int16)
}
