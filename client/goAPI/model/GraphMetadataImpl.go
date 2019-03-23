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
	initialized     bool
	descriptors     map[string]types.TGAttributeDescriptor
	descriptorsById map[int64]types.TGAttributeDescriptor
	edgeTypes       map[string]types.TGEdgeType
	edgeTypesById   map[int]types.TGEdgeType
	nodeTypes       map[string]types.TGNodeType
	nodeTypesById   map[int]types.TGNodeType
	graphObjFactory *GraphObjectFactory
}

func DefaultGraphMetadata() *GraphMetadata {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GraphMetadata{})

	newGraphMetadata := GraphMetadata{
		initialized:     false,
		descriptors:     make(map[string]types.TGAttributeDescriptor, 0),
		descriptorsById: make(map[int64]types.TGAttributeDescriptor, 0),
		edgeTypes:       make(map[string]types.TGEdgeType, 0),
		edgeTypesById:   make(map[int]types.TGEdgeType, 0),
		nodeTypes:       make(map[string]types.TGNodeType, 0),
		nodeTypesById:   make(map[int]types.TGNodeType, 0),
	}
	return &newGraphMetadata
}

func NewGraphMetadata(gof *GraphObjectFactory) *GraphMetadata {
	newGraphMetadata := DefaultGraphMetadata()
	newGraphMetadata.graphObjFactory = gof
	return newGraphMetadata
}

/////////////////////////////////////////////////////////////////
// Helper functions for types.TGGraphMetadata
/////////////////////////////////////////////////////////////////

func (obj *GraphMetadata) GetNewAttributeDescriptors() ([]types.TGAttributeDescriptor, types.TGError) {
	if obj.descriptors == nil || len(obj.descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNewAttributeDescriptors as there are NO new attrDesc`"))
		return nil, nil
	}
	attrDesc := make([]types.TGAttributeDescriptor, 0)
	for _, desc := range obj.descriptors {
		if desc.(*AttributeDescriptor).GetAttributeId() < 0 {
			attrDesc = append(attrDesc, desc)
		}
	}
	return attrDesc, nil
}

func (obj *GraphMetadata) GetConnection() types.TGConnection {
	return obj.graphObjFactory.GetConnection()
}

func (obj *GraphMetadata) GetGraphObjectFactory() *GraphObjectFactory {
	return obj.graphObjFactory
}

func (obj *GraphMetadata) GetAttributeDescriptorsById() map[int64]types.TGAttributeDescriptor {
	return obj.descriptorsById
}

func (obj *GraphMetadata) GetEdgeTypesById() map[int]types.TGEdgeType {
	return obj.edgeTypesById
}

func (obj *GraphMetadata) GetNodeTypesById() map[int]types.TGNodeType {
	return obj.nodeTypesById
}

func (obj *GraphMetadata) IsInitialized() bool {
	return obj.initialized
}

func (obj *GraphMetadata) SetInitialized(flag bool) {
	obj.initialized = flag
}

func (obj *GraphMetadata) SetAttributeDescriptors(attrDesc map[string]types.TGAttributeDescriptor) {
	obj.descriptors = attrDesc
}

func (obj *GraphMetadata) SetAttributeDescriptorsById(attrDescId map[int64]types.TGAttributeDescriptor) {
	obj.descriptorsById = attrDescId
}

func (obj *GraphMetadata) SetEdgeTypes(edgeTypes map[string]types.TGEdgeType) {
	obj.edgeTypes = edgeTypes
}

func (obj *GraphMetadata) SetEdgeTypesById(edgeTypesId map[int]types.TGEdgeType) {
	obj.edgeTypesById = edgeTypesId
}

func (obj *GraphMetadata) SetNodeTypes(nodeTypes map[string]types.TGNodeType) {
	obj.nodeTypes = nodeTypes
}

func (obj *GraphMetadata) SetNodeTypesById(nodeTypes map[int]types.TGNodeType) {
	obj.nodeTypesById = nodeTypes
}

func (obj *GraphMetadata) UpdateMetadata(attrDescList []types.TGAttributeDescriptor, nodeTypeList []types.TGNodeType, edgeTypeList []types.TGEdgeType) types.TGError {
	if attrDescList != nil {
		for _, attrDesc := range attrDescList {
			obj.descriptors[attrDesc.GetName()] = attrDesc.(*AttributeDescriptor)
			obj.descriptorsById[attrDesc.GetAttributeId()] = attrDesc.(*AttributeDescriptor)
		}
	}
	if nodeTypeList != nil {
		for _, nodeType := range nodeTypeList {
			obj.nodeTypes[nodeType.GetName()] = nodeType.(*NodeType)
			obj.nodeTypesById[nodeType.GetEntityTypeId()] = nodeType.(*NodeType)
		}
	}
	if edgeTypeList != nil {
		for _, edgeType := range edgeTypeList {
			obj.edgeTypes[edgeType.GetName()] = edgeType.(*EdgeType)
			obj.edgeTypesById[edgeType.GetEntityTypeId()] = edgeType.(*EdgeType)
		}
	}
	obj.initialized = true
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGGraphMetadata
/////////////////////////////////////////////////////////////////

// CreateAttributeDescriptor creates Attribute Descriptor
func (obj *GraphMetadata) CreateAttributeDescriptor(attrName string, attrType int, isArray bool) types.TGAttributeDescriptor {
	newAttrDesc := NewAttributeDescriptorAsArray(attrName, attrType, isArray)
	obj.descriptors[attrName] = newAttrDesc
	return newAttrDesc
}

// CreateAttributeDescriptorForDataType creates Attribute Descriptor for data/attribute type - New in GO Lang
func (obj *GraphMetadata) CreateAttributeDescriptorForDataType(attrName string, dataTypeClassName string) types.TGAttributeDescriptor {
	attrType := types.GetAttributeTypeFromName(dataTypeClassName)
	logger.Log(fmt.Sprintf("GraphMetadata CreateAttributeDescriptorForDataType creating attribute descriptor for '%+v' w/ type '%+v'", attrName, attrType))
	newAttrDesc := NewAttributeDescriptorAsArray(attrName, attrType.GetTypeId(), false)
	obj.descriptors[attrName] = newAttrDesc
	return newAttrDesc
}

// GetAttributeDescriptor gets the Attribute Descriptor by name
func (obj *GraphMetadata) GetAttributeDescriptor(attrName string) (types.TGAttributeDescriptor, types.TGError) {
	if obj.descriptors == nil || len(obj.descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptor as there are NO attrDesc`"))
		return nil, nil
	}
	desc := obj.descriptors[attrName]
	//if desc == nil || desc.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Attribute Descriptors")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return desc, nil
}

// GetAttributeDescriptorById gets the Attribute Descriptor by name
func (obj *GraphMetadata) GetAttributeDescriptorById(id int64) (types.TGAttributeDescriptor, types.TGError) {
	if obj.descriptorsById == nil || len(obj.descriptorsById) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptorById as there are NO attrDesc`"))
		return nil, nil
	}
	desc := obj.descriptorsById[id]
	//if desc.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Attribute Descriptors")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return desc, nil
}

// GetAttributeDescriptors gets a list of Attribute Descriptors
func (obj *GraphMetadata) GetAttributeDescriptors() ([]types.TGAttributeDescriptor, types.TGError) {
	if obj.descriptors == nil || len(obj.descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptors as there are NO attrDesc`"))
		return nil, nil
	}
	attrDesc := make([]types.TGAttributeDescriptor, 0)
	for _, desc := range obj.descriptors {
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
	if obj.edgeTypes == nil || len(obj.edgeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetEdgeType as there are NO edges`"))
		return nil, nil
	}
	edge := obj.edgeTypes[typeName]
	//if edge.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Edges")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return edge, nil
}

// GetEdgeTypeById returns the Edge Type by id
func (obj *GraphMetadata) GetEdgeTypeById(id int) (types.TGEdgeType, types.TGError) {
	if obj.edgeTypesById == nil || len(obj.edgeTypesById) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetEdgeTypeById as there are NO edges`"))
		return nil, nil
	}
	edge := obj.edgeTypesById[id]
	//if edge.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Edges")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return edge, nil
}

// GetEdgeTypes returns a set of know edge Type
func (obj *GraphMetadata) GetEdgeTypes() ([]types.TGEdgeType, types.TGError) {
	if obj.edgeTypes == nil || len(obj.edgeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetEdgeTypes as there are NO edges`"))
		return nil, nil
	}
	edgeTypes := make([]types.TGEdgeType, 0)
	for _, edgeType := range obj.edgeTypes {
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
	if obj.nodeTypes == nil || len(obj.nodeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNodeType as there are NO nodes"))
		return nil, nil
	}
	node := obj.nodeTypes[typeName]
	//if node.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Nodes")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return node, nil
}

// GetNodeTypeById returns the Node types by id
func (obj *GraphMetadata) GetNodeTypeById(id int) (types.TGNodeType, types.TGError) {
	if obj.nodeTypesById == nil || len(obj.nodeTypesById) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNodeTypeById as there are NO nodes"))
		return nil, nil
	}
	node := obj.nodeTypesById[id]
	logger.Log(fmt.Sprintf("Inside GraphMetadata:GetNodeTypeById read node as '%+v'", node))
	//if node.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Nodes")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return node, nil
}

// GetNodeTypes returns a set of Node Type defined in the System
func (obj *GraphMetadata) GetNodeTypes() ([]types.TGNodeType, types.TGError) {
	if obj.nodeTypes == nil || len(obj.nodeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNodeTypes as there are NO nodes"))
		return nil, nil
	}
	nodeTypes := make([]types.TGNodeType, 0)
	for _, nodeType := range obj.nodeTypes {
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
	_, err := fmt.Fprintln(&b, obj.initialized, obj.descriptors, obj.descriptorsById, obj.edgeTypes, obj.edgeTypesById,
		obj.nodeTypes, obj.nodeTypesById, obj.graphObjFactory)
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
	_, err := fmt.Fscanln(b, &obj.initialized, &obj.descriptors, &obj.descriptorsById, &obj.edgeTypes, &obj.edgeTypesById,
		&obj.nodeTypes, &obj.nodeTypesById, &obj.graphObjFactory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GraphMetadata:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
