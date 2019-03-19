package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
 * File name: TGGraphMetadata.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type GraphMetadata struct {
	Initialized     bool
	Descriptors     map[string]types.TGAttributeDescriptor
	DescriptorsById map[int64]types.TGAttributeDescriptor
	EdgeTypes       map[string]types.TGEdgeType
	EdgeTypesById   map[int]types.TGEdgeType
	NodeTypes       map[string]types.TGNodeType
	NodeTypesById   map[int]types.TGNodeType
	GraphObjFactory *GraphObjectFactory
}

func DefaultGraphMetadata() *GraphMetadata {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GraphMetadata{})

	newGraphMetadata := GraphMetadata{
		Initialized:     false,
		Descriptors:     make(map[string]types.TGAttributeDescriptor, 0),
		DescriptorsById: make(map[int64]types.TGAttributeDescriptor, 0),
		EdgeTypes:       make(map[string]types.TGEdgeType, 0),
		EdgeTypesById:   make(map[int]types.TGEdgeType, 0),
		NodeTypes:       make(map[string]types.TGNodeType, 0),
		NodeTypesById:   make(map[int]types.TGNodeType, 0),
	}
	return &newGraphMetadata
}

func NewGraphMetadata(gof *GraphObjectFactory) *GraphMetadata {
	newGraphMetadata := DefaultGraphMetadata()
	newGraphMetadata.GraphObjFactory = gof
	return newGraphMetadata
}

/////////////////////////////////////////////////////////////////
// Helper functions for types.TGGraphMetadata
/////////////////////////////////////////////////////////////////

func (obj *GraphMetadata) GetNewAttributeDescriptors() ([]types.TGAttributeDescriptor, types.TGError) {
	if obj.Descriptors == nil || len(obj.Descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNewAttributeDescriptors as there are NO new attrDesc`"))
		return nil, nil
	}
	attrDesc := make([]types.TGAttributeDescriptor, 0)
	for _, desc := range obj.Descriptors {
		if desc.(*AttributeDescriptor).GetAttributeId() < 0 {
			attrDesc = append(attrDesc, desc)
		}
	}
	return attrDesc, nil
}

func (obj *GraphMetadata) GetConnection() types.TGConnection {
	return obj.GraphObjFactory.GetConnection()
}

func (obj *GraphMetadata) GetGraphObjectFactory() *GraphObjectFactory {
	return obj.GraphObjFactory
}

func (obj *GraphMetadata) IsInitialized() bool {
	return obj.Initialized
}

func (obj *GraphMetadata) SetInitialized(flag bool) {
	obj.Initialized = flag
}

func (obj *GraphMetadata) UpdateMetadata(attrDescList []types.TGAttributeDescriptor, nodeTypeList []types.TGNodeType, edgeTypeList []types.TGEdgeType) types.TGError {
	if attrDescList != nil {
		for _, attrDesc := range attrDescList {
			obj.Descriptors[attrDesc.GetName()] = attrDesc.(*AttributeDescriptor)
			obj.DescriptorsById[attrDesc.GetAttributeId()] = attrDesc.(*AttributeDescriptor)
		}
	}
	if nodeTypeList != nil {
		for _, nodeType := range nodeTypeList {
			obj.NodeTypes[nodeType.GetName()] = nodeType.(*NodeType)
			obj.NodeTypesById[nodeType.GetEntityTypeId()] = nodeType.(*NodeType)
		}
	}
	if edgeTypeList != nil {
		for _, edgeType := range edgeTypeList {
			obj.EdgeTypes[edgeType.GetName()] = edgeType.(*EdgeType)
			obj.EdgeTypesById[edgeType.GetEntityTypeId()] = edgeType.(*EdgeType)
		}
	}
	obj.Initialized = true
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGGraphMetadata
/////////////////////////////////////////////////////////////////

// CreateAttributeDescriptor creates Attribute Descriptor
func (obj *GraphMetadata) CreateAttributeDescriptor(attrName string, attrType int, isArray bool) types.TGAttributeDescriptor {
	newAttrDesc := NewAttributeDescriptorAsArray(attrName, attrType, isArray)
	obj.Descriptors[attrName] = newAttrDesc
	return newAttrDesc
}

// CreateAttributeDescriptorForDataType creates Attribute Descriptor for data/attribute type - New in GO Lang
func (obj *GraphMetadata) CreateAttributeDescriptorForDataType(attrName string, dataTypeClassName string) types.TGAttributeDescriptor {
	attrType := types.GetAttributeTypeFromName(dataTypeClassName)
	logger.Log(fmt.Sprintf("GraphMetadata CreateAttributeDescriptorForDataType creating attribute descriptor for '%+v' w/ type '%+v'", attrName, attrType))
	newAttrDesc := NewAttributeDescriptorAsArray(attrName, attrType.TypeId, false)
	obj.Descriptors[attrName] = newAttrDesc
	return newAttrDesc
}

// GetAttributeDescriptor gets the Attribute Descriptor by name
func (obj *GraphMetadata) GetAttributeDescriptor(attrName string) (types.TGAttributeDescriptor, types.TGError) {
	if obj.Descriptors == nil || len(obj.Descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptor as there are NO attrDesc`"))
		return nil, nil
	}
	desc := obj.Descriptors[attrName]
	//if desc == nil || desc.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Attribute Descriptors")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return desc, nil
}

// GetAttributeDescriptorById gets the Attribute Descriptor by name
func (obj *GraphMetadata) GetAttributeDescriptorById(id int64) (types.TGAttributeDescriptor, types.TGError) {
	if obj.DescriptorsById == nil || len(obj.DescriptorsById) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptorById as there are NO attrDesc`"))
		return nil, nil
	}
	desc := obj.DescriptorsById[id]
	//if desc.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Attribute Descriptors")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return desc, nil
}

// GetAttributeDescriptors gets a list of Attribute Descriptors
func (obj *GraphMetadata) GetAttributeDescriptors() ([]types.TGAttributeDescriptor, types.TGError) {
	if obj.Descriptors == nil || len(obj.Descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptors as there are NO attrDesc`"))
		return nil, nil
	}
	attrDesc := make([]types.TGAttributeDescriptor, 0)
	for _, desc := range obj.Descriptors {
		attrDesc = append(attrDesc, desc)
	}
	return attrDesc, nil
}

// CreateCompositeKey creates composite key
func (obj *GraphMetadata) CreateCompositeKey(nodeTypeName string) types.TGKey {
	compKey := NewCompositeKey(obj, nodeTypeName)
	return compKey
}

// CreateEdgeType creates Edge Type
func (obj *GraphMetadata) CreateEdgeType(typeName string, parentEdgeType types.TGEdgeType) types.TGEdgeType {
	newEdgeType := NewEdgeType(typeName, parentEdgeType.GetDirectionType(), parentEdgeType)
	return newEdgeType
}

// GetEdgeType returns the Edge by name
func (obj *GraphMetadata) GetEdgeType(typeName string) (types.TGEdgeType, types.TGError) {
	if obj.EdgeTypes == nil || len(obj.EdgeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetEdgeType as there are NO edges`"))
		return nil, nil
	}
	edge := obj.EdgeTypes[typeName]
	//if edge.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Edges")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return edge, nil
}

// GetEdgeTypeById returns the Edge Type by id
func (obj *GraphMetadata) GetEdgeTypeById(id int) (types.TGEdgeType, types.TGError) {
	if obj.EdgeTypesById == nil || len(obj.EdgeTypesById) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetEdgeTypeById as there are NO edges`"))
		return nil, nil
	}
	edge := obj.EdgeTypesById[id]
	//if edge.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Edges")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return edge, nil
}

// GetEdgeTypes returns a set of know edge Type
func (obj *GraphMetadata) GetEdgeTypes() ([]types.TGEdgeType, types.TGError) {
	if obj.EdgeTypes == nil || len(obj.EdgeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetEdgeTypes as there are NO edges`"))
		return nil, nil
	}
	edgeTypes := make([]types.TGEdgeType, 0)
	for _, edgeType := range obj.EdgeTypes {
		edgeTypes = append(edgeTypes, edgeType)
	}
	return edgeTypes, nil
}

// CreateNodeType creates a Node Type
func (obj *GraphMetadata) CreateNodeType(typeName string, parentNodeType types.TGNodeType) types.TGNodeType {
	newNodeType := NewNodeType(typeName, parentNodeType)
	return newNodeType
}

// GetNodeType gets Node Type by name
func (obj *GraphMetadata) GetNodeType(typeName string) (types.TGNodeType, types.TGError) {
	if obj.NodeTypes == nil || len(obj.NodeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNodeType as there are NO nodes"))
		return nil, nil
	}
	node := obj.NodeTypes[typeName]
	//if node.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Nodes")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return node, nil
}

// GetNodeTypeById returns the Node types by id
func (obj *GraphMetadata) GetNodeTypeById(id int) (types.TGNodeType, types.TGError) {
	if obj.NodeTypesById == nil || len(obj.NodeTypesById) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNodeTypeById as there are NO nodes"))
		return nil, nil
	}
	node := obj.NodeTypesById[id]
	logger.Log(fmt.Sprintf("Inside GraphMetadata:GetNodeTypeById read node as '%+v'", node))
	//if node.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Nodes")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return node, nil
}

// GetNodeTypes returns a set of Node Type defined in the System
func (obj *GraphMetadata) GetNodeTypes() ([]types.TGNodeType, types.TGError) {
	if obj.NodeTypes == nil || len(obj.NodeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNodeTypes as there are NO nodes"))
		return nil, nil
	}
	nodeTypes := make([]types.TGNodeType, 0)
	for _, nodeType := range obj.NodeTypes {
		nodeTypes = append(nodeTypes, nodeType)
	}
	return nodeTypes, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *GraphMetadata) ReadExternal(is types.TGInputStream) types.TGError {
	// No-op for Now
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *GraphMetadata) WriteExternal(os types.TGOutputStream) types.TGError {
	// No-op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *GraphMetadata) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.Initialized, obj.Descriptors, obj.DescriptorsById, obj.EdgeTypes, obj.EdgeTypesById,
		obj.NodeTypes, obj.NodeTypesById, obj.GraphObjFactory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GraphMetadata:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *GraphMetadata) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.Initialized, &obj.Descriptors, &obj.DescriptorsById, &obj.EdgeTypes, &obj.EdgeTypesById,
		&obj.NodeTypes, &obj.NodeTypesById, &obj.GraphObjFactory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GraphMetadata:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
