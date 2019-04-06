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
 * File name: TGGraphMetadata.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGGraphMetadata interface {
	TGSerializable
	// CreateAttributeDescriptor creates Attribute Descriptor
	CreateAttributeDescriptor(attrName string, attrType int, isArray bool) TGAttributeDescriptor
	// CreateAttributeDescriptorForDataType creates Attribute Descriptor for data/attribute type - New in GO Lang
	CreateAttributeDescriptorForDataType(attrName string, dataTypeClassName string) TGAttributeDescriptor
	// GetAttributeDescriptor gets the Attribute Descriptor by name
	GetAttributeDescriptor(attributeName string) (TGAttributeDescriptor, TGError)
	// GetAttributeDescriptorById gets the Attribute Descriptor by name
	GetAttributeDescriptorById(id int64) (TGAttributeDescriptor, TGError)
	// GetAttributeDescriptors gets a list of Attribute Descriptors
	GetAttributeDescriptors() ([]TGAttributeDescriptor, TGError)
	// CreateCompositeKey creates composite key
	CreateCompositeKey(nodeTypeName string) TGKey
	// CreateEdgeType creates Edge Type
	CreateEdgeType(typeName string, parentEdgeType TGEdgeType) TGEdgeType
	// GetConnection returns the connection from its graph object factory
	GetConnection() TGConnection
	// GetEdgeType returns the Edge by name
	GetEdgeType(typeName string) (TGEdgeType, TGError)
	// GetEdgeTypeById returns the Edge Type by id
	GetEdgeTypeById(id int) (TGEdgeType, TGError)
	// GetEdgeTypes returns a set of know edge Type
	GetEdgeTypes() ([]TGEdgeType, TGError)
	// CreateNodeType creates a Node Type
	CreateNodeType(typeName string, parentNodeType TGNodeType) TGNodeType
	// GetNodeType gets Node Type by name
	GetNodeType(typeName string) (TGNodeType, TGError)
	// GetNodeTypeById returns the Node types by id
	GetNodeTypeById(id int) (TGNodeType, TGError)
	// GetNodeTypes returns a set of Node Type defined in the System
	GetNodeTypes() ([]TGNodeType, TGError)
}
