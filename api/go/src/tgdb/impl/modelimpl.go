/*
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File Name: modelimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: modelimpl.go 4270 2020-08-19 22:18:58Z nimish $
 */

package impl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"tgdb"
)

var logger = DefaultTGLogManager().GetLogger()


type GraphMetadata struct {
	initialized     bool
	descriptors     map[string]tgdb.TGAttributeDescriptor
	descriptorsById map[int64]tgdb.TGAttributeDescriptor
	edgeTypes       map[string]tgdb.TGEdgeType
	edgeTypesById   map[int]tgdb.TGEdgeType
	nodeTypes       map[string]tgdb.TGNodeType
	nodeTypesById   map[int]tgdb.TGNodeType
	graphObjFactory *GraphObjectFactory
}

type GraphObjectFactory struct {
	graphMData *GraphMetadata
	connection tgdb.TGConnection
}





// TODO: Revisit later - Once SetAttributeViaDescriptor is properly implemented after discussing with TGDB Engineering Team
func setAttributeViaDescriptor(obj tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) tgdb.TGError {
	if attrDesc == nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractEntity:setAttributeViaDescriptor as AttrDesc is EMPTY"))
		errMsg := fmt.Sprintf("Attribute Descriptor cannot be null")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if value == nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractEntity:setAttributeViaDescriptor as value is EMPTY"))
		errMsg := fmt.Sprintf("Attribute value is required")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// TODO: Do we need to validate if this descriptor exists as part of Graph Meta Data???
	// If attribute is not present in the set, create a new one
	attrDescName := attrDesc.GetName()
	attr := obj.(*AbstractEntity).Attributes[attrDescName]
	if attr == nil {
		if attrDesc.GetAttrType() == AttributeTypeInvalid {
			logger.Error(fmt.Sprint("ERROR: Returning AbstractEntity:setAttributeViaDescriptor as AttrDesc.GetAttrType() == types.AttributeTypeInvalid"))
			errMsg := fmt.Sprintf("Attribute descriptor is of incorrect type")
			return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
		}
		// TODO: Revisit later - For some reason, it goes in infinite loop, alternative is to create attribute and assign owner and value later
		//newAttr, aErr := CreateAttribute(obj, AttrDesc, value)
		newAttr, aErr := CreateAttributeWithDesc(nil, attrDesc, nil)
		if aErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:setAttributeViaDescriptor unable to create attribute '%s' w/ descriptor and value '%+v'", attrDesc, value))
			return aErr
		}
		newAttr.SetOwner(obj)
		attr = newAttr
	}
	// Value can be null here
	if !attr.GetIsModified() {
		obj.(*AbstractEntity).ModifiedAttributes = append(obj.(*AbstractEntity).ModifiedAttributes, attr)
	}
	// Set the attribute value
	err := attr.SetValue(value)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:setAttributeViaDescriptor unable to set attribute value w/ error '%+v'", err.Error()))
		return err
	}
	// Add it to the set
	obj.(*AbstractEntity).Attributes[attrDesc.Name] = attr
	return nil
}





func DefaultGraphMetadata() *GraphMetadata {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GraphMetadata{})

	newGraphMetadata := GraphMetadata{
		initialized:     false,
		descriptors:     make(map[string]tgdb.TGAttributeDescriptor, 0),
		descriptorsById: make(map[int64]tgdb.TGAttributeDescriptor, 0),
		edgeTypes:       make(map[string]tgdb.TGEdgeType, 0),
		edgeTypesById:   make(map[int]tgdb.TGEdgeType, 0),
		nodeTypes:       make(map[string]tgdb.TGNodeType, 0),
		nodeTypesById:   make(map[int]tgdb.TGNodeType, 0),
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

func (obj *GraphMetadata) GetNewAttributeDescriptors() ([]tgdb.TGAttributeDescriptor, tgdb.TGError) {
	if obj.descriptors == nil || len(obj.descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNewAttributeDescriptors as there are NO new AttrDesc`"))
		return nil, nil
	}
	attrDesc := make([]tgdb.TGAttributeDescriptor, 0)
	for _, desc := range obj.descriptors {
		if desc.(*AttributeDescriptor).GetAttributeId() < 0 {
			attrDesc = append(attrDesc, desc)
		}
	}
	return attrDesc, nil
}

func (obj *GraphMetadata) GetConnection() tgdb.TGConnection {
	return obj.graphObjFactory.GetConnection()
}

func (obj *GraphMetadata) GetGraphObjectFactory() *GraphObjectFactory {
	return obj.graphObjFactory
}

func (obj *GraphMetadata) GetAttributeDescriptorsById() map[int64]tgdb.TGAttributeDescriptor {
	return obj.descriptorsById
}

func (obj *GraphMetadata) GetEdgeTypesById() map[int]tgdb.TGEdgeType {
	return obj.edgeTypesById
}

func (obj *GraphMetadata) GetNodeTypesById() map[int]tgdb.TGNodeType {
	return obj.nodeTypesById
}

func (obj *GraphMetadata) IsInitialized() bool {
	return obj.initialized
}

func (obj *GraphMetadata) SetInitialized(flag bool) {
	obj.initialized = flag
}

func (obj *GraphMetadata) SetAttributeDescriptors(attrDesc map[string]tgdb.TGAttributeDescriptor) {
	obj.descriptors = attrDesc
}

func (obj *GraphMetadata) SetAttributeDescriptorsById(attrDescId map[int64]tgdb.TGAttributeDescriptor) {
	obj.descriptorsById = attrDescId
}

func (obj *GraphMetadata) SetEdgeTypes(edgeTypes map[string]tgdb.TGEdgeType) {
	obj.edgeTypes = edgeTypes
}

func (obj *GraphMetadata) SetEdgeTypesById(edgeTypesId map[int]tgdb.TGEdgeType) {
	obj.edgeTypesById = edgeTypesId
}

func (obj *GraphMetadata) SetNodeTypes(nodeTypes map[string]tgdb.TGNodeType) {
	obj.nodeTypes = nodeTypes
}

func (obj *GraphMetadata) SetNodeTypesById(nodeTypes map[int]tgdb.TGNodeType) {
	obj.nodeTypesById = nodeTypes
}

func (obj *GraphMetadata) UpdateMetadata(attrDescList []tgdb.TGAttributeDescriptor, nodeTypeList []tgdb.TGNodeType, edgeTypeList []tgdb.TGEdgeType) tgdb.TGError {
	if attrDescList != nil {
		for _, attrDesc := range attrDescList {
			obj.descriptors[attrDesc.GetName()] = attrDesc.(*AttributeDescriptor)
			obj.descriptorsById[attrDesc.GetAttributeId()] = attrDesc.(*AttributeDescriptor)
		}
	}
	if nodeTypeList != nil {
		for _, nodeType := range nodeTypeList {
			nodeType.(*NodeType).UpdateMetadata(obj)
			obj.nodeTypes[nodeType.GetName()] = nodeType.(*NodeType)
			obj.nodeTypesById[nodeType.GetEntityTypeId()] = nodeType.(*NodeType)
		}
	}
	if edgeTypeList != nil {
		for _, edgeType := range edgeTypeList {
			edgeType.(*EdgeType).UpdateMetadata(obj)
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
func (obj *GraphMetadata) CreateAttributeDescriptor(attrName string, attrType int, isArray bool) tgdb.TGAttributeDescriptor {
	newAttrDesc := NewAttributeDescriptorAsArray(attrName, attrType, isArray)
	obj.descriptors[attrName] = newAttrDesc
	return newAttrDesc
}

// CreateAttributeDescriptorForDataType creates Attribute Descriptor for data/attribute type - New in GO Lang
func (obj *GraphMetadata) CreateAttributeDescriptorForDataType(attrName string, dataTypeClassName string) tgdb.TGAttributeDescriptor {
	attrType := GetAttributeTypeFromName(dataTypeClassName)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("GraphMetadata CreateAttributeDescriptorForDataType creating attribute descriptor for '%+v' w/ type '%+v'", attrName, attrType))
	}
	newAttrDesc := NewAttributeDescriptorAsArray(attrName, attrType.GetTypeId(), false)
	obj.descriptors[attrName] = newAttrDesc
	return newAttrDesc
}

// GetAttributeDescriptor gets the Attribute Descriptor by Name
func (obj *GraphMetadata) GetAttributeDescriptor(attrName string) (tgdb.TGAttributeDescriptor, tgdb.TGError) {
	if obj.descriptors == nil || len(obj.descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptor as there are NO AttrDesc`"))
		return nil, nil
	}
	desc := obj.descriptors[attrName]
	//if desc == nil || desc.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Attribute Descriptors")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return desc, nil
}

// GetAttributeDescriptorById gets the Attribute Descriptor by Name
func (obj *GraphMetadata) GetAttributeDescriptorById(id int64) (tgdb.TGAttributeDescriptor, tgdb.TGError) {
	if obj.descriptorsById == nil || len(obj.descriptorsById) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptorById as there are NO AttrDesc`"))
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
func (obj *GraphMetadata) GetAttributeDescriptors() ([]tgdb.TGAttributeDescriptor, tgdb.TGError) {
	if obj.descriptors == nil || len(obj.descriptors) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetAttributeDescriptors as there are NO AttrDesc`"))
		return nil, nil
	}
	attrDesc := make([]tgdb.TGAttributeDescriptor, 0)
	for _, desc := range obj.descriptors {
		attrDesc = append(attrDesc, desc)
	}
	return attrDesc, nil
}

// CreateCompositeKey creates composite key
func (obj *GraphMetadata) CreateCompositeKey(nodeTypeName string) tgdb.TGKey {
	compKey := NewCompositeKey(obj, nodeTypeName)
	return compKey
}

// CreateEdgeType creates Edge Type
func (obj *GraphMetadata) CreateEdgeType(typeName string, parentEdgeType tgdb.TGEdgeType) tgdb.TGEdgeType {
	newEdgeType := NewEdgeType(typeName, parentEdgeType.GetDirectionType(), parentEdgeType)
	return newEdgeType
}

// GetEdgeType returns the Edge by Name
func (obj *GraphMetadata) GetEdgeType(typeName string) (tgdb.TGEdgeType, tgdb.TGError) {
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
func (obj *GraphMetadata) GetEdgeTypeById(id int) (tgdb.TGEdgeType, tgdb.TGError) {
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
func (obj *GraphMetadata) GetEdgeTypes() ([]tgdb.TGEdgeType, tgdb.TGError) {
	if obj.edgeTypes == nil || len(obj.edgeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetEdgeTypes as there are NO edges`"))
		return nil, nil
	}
	edgeTypes := make([]tgdb.TGEdgeType, 0)
	for _, edgeType := range obj.edgeTypes {
		edgeTypes = append(edgeTypes, edgeType)
	}
	return edgeTypes, nil
}

// CreateNodeType creates a Node Type
func (obj *GraphMetadata) CreateNodeType(typeName string, parentNodeType tgdb.TGNodeType) tgdb.TGNodeType {
	newNodeType := NewNodeType(typeName, parentNodeType)
	return newNodeType
}

// GetNodeType gets Node Type by Name
func (obj *GraphMetadata) GetNodeType(typeName string) (tgdb.TGNodeType, tgdb.TGError) {
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
func (obj *GraphMetadata) GetNodeTypeById(id int) (tgdb.TGNodeType, tgdb.TGError) {
	if obj.nodeTypesById == nil || len(obj.nodeTypesById) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNodeTypeById as there are NO nodes"))
		return nil, nil
	}
	node := obj.nodeTypesById[id]
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning GraphMetadata:GetNodeTypeById read node as '%+v'", node))
	}
	//if node.GetSystemType() == types.SystemTypeInvalid {
	//	errMsg := fmt.Sprintf("There are no Nodes")
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	return node, nil
}

// GetNodeTypes returns a set of Node Type defined in the System
func (obj *GraphMetadata) GetNodeTypes() ([]tgdb.TGNodeType, tgdb.TGError) {
	if obj.nodeTypes == nil || len(obj.nodeTypes) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning GraphMetadata:GetNodeTypes as there are NO nodes"))
		return nil, nil
	}
	nodeTypes := make([]tgdb.TGNodeType, 0)
	for _, nodeType := range obj.nodeTypes {
		nodeTypes = append(nodeTypes, nodeType)
	}
	return nodeTypes, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *GraphMetadata) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *GraphMetadata) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
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


func DefaultGraphObjectFactory() *GraphObjectFactory {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GraphObjectFactory{})

	newGraphObjectFactory := GraphObjectFactory{}
	return &newGraphObjectFactory
}

func NewGraphObjectFactory(conn tgdb.TGConnection) *GraphObjectFactory {
	newGraphObjectFactory := DefaultGraphObjectFactory()
	// TODO: Graph meta data cannot be passed in.
	// There will be one meta data object per object factory and one object factory per connection even though
	// connections can share the same channel. How do we handle update notifications from other clients?
	newGraphObjectFactory.graphMData = NewGraphMetadata(newGraphObjectFactory)
	newGraphObjectFactory.connection = conn
	return newGraphObjectFactory
}

/////////////////////////////////////////////////////////////////
// Helper functions for GraphObjectFactory
/////////////////////////////////////////////////////////////////

func (obj *GraphObjectFactory) GetConnection() tgdb.TGConnection {
	return obj.connection
}

func (obj *GraphObjectFactory) GetGraphMetaData() *GraphMetadata {
	return obj.graphMData
}

// TODO: Revisit later to optimize and consolidate common functions here instead of each individual structure implementation
/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGGraphObjectFactory
/////////////////////////////////////////////////////////////////

// CreateCompositeKey creates a CompositeKey for a SystemTypeNode. The composite key can also be a single key
func (obj *GraphObjectFactory) CreateCompositeKey(nodeTypeName string) (tgdb.TGKey, tgdb.TGError) {
	_, err := obj.graphMData.GetNodeType(nodeTypeName)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateCompositeKey as node NOT found"))
		errMsg := fmt.Sprintf("Node desc with Name %s not found", nodeTypeName)
		return nil, GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return NewCompositeKey(obj.graphMData, nodeTypeName), nil
}

// CreateEdgeWithEdgeType creates an Edge
func (obj *GraphObjectFactory) CreateEdgeWithEdgeType(fromNode tgdb.TGNode, toNode tgdb.TGNode, edgeType tgdb.TGEdgeType) (tgdb.TGEdge, tgdb.TGError) {
	newEdge := NewEdgeWithEdgeType(obj.graphMData, fromNode, toNode, edgeType)
	if newEdge.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateEdgeWithEdgeType as edge in NOT initialized"))
		errMsg := fmt.Sprintf("Unable to create an edge with type %s", edgeType.(*EdgeType).Name)
		return nil, GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	fromNode.AddEdge(newEdge)
	toNode.AddEdge(newEdge)
	return newEdge, nil
}

// CreateEdgeWithDirection creates an Edge with a direction
func (obj *GraphObjectFactory) CreateEdgeWithDirection(fromNode tgdb.TGNode, toNode tgdb.TGNode, directionType tgdb.TGDirectionType) (tgdb.TGEdge, tgdb.TGError) {
	newEdge := NewEdgeWithDirection(obj.graphMData, fromNode, toNode, directionType)
	if newEdge.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateEdgeWithDirection as edge in NOT initialized"))
		errMsg := fmt.Sprintf("Unable to create an edge with direction %s", directionType)
		return nil, GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	fromNode.AddEdge(newEdge)
	toNode.AddEdge(newEdge)
	return newEdge, nil
}

// CreateEntity creates entity based on the entity kind specified
func (obj *GraphObjectFactory) CreateEntity(entityKind tgdb.TGEntityKind) (tgdb.TGEntity, tgdb.TGError) {
	switch entityKind {
	case tgdb.EntityKindNode:
		return NewNode(obj.graphMData), nil
	case tgdb.EntityKindEdge:
		return NewEdge(obj.graphMData), nil
	case tgdb.EntityKindGraph:
		return NewGraph(obj.graphMData), nil
	}
	logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateEntity as entity kind specified is INVALID"))
	errMsg := fmt.Sprintf("Invalid entity kind %d specified", entityKind)
	return nil, GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
}

// CreateEntityId creates entity id from input buffer
func (obj *GraphObjectFactory) CreateEntityId(buf []byte) (tgdb.TGEntityId, tgdb.TGError) {
	return NewByteArrayEntity(buf), nil
}

// CreateGraph creates a Graph
func (obj *GraphObjectFactory) CreateGraph(name string) (tgdb.TGGraph, tgdb.TGError) {
	graph := NewGraphWithName(obj.graphMData, name)
	if graph.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateGraph as graph in NOT initialized"))
		errMsg := fmt.Sprint("Unable to create a graph with this Graph Object Factory")
		return nil, GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return graph, nil
}

// CreateNode creates a Node
func (obj *GraphObjectFactory) CreateNode() (tgdb.TGNode, tgdb.TGError) {
	node := NewNode(obj.graphMData)
	if node.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateNode as node in NOT initialized"))
		errMsg := fmt.Sprint("Unable to create a node with this Graph Object Factory")
		return nil, GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return node, nil
}

// CreateNodeInGraph creates Node within this Graph. There is a default Root Graph.
func (obj *GraphObjectFactory) CreateNodeInGraph(nodeType tgdb.TGNodeType) (tgdb.TGNode, tgdb.TGError) {
	node := NewNodeWithType(obj.graphMData, nodeType)
	if node.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateNodeInGraph as node in NOT initialized"))
		errMsg := fmt.Sprint("Unable to create a node with this Graph Object Factory")
		return nil, GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return node, nil
}

// The following methods are indirectly used by encoding/gob methods that are used above

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *GraphObjectFactory) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.graphMData, obj.connection)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GraphObjectFactory:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *GraphObjectFactory) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.graphMData, &obj.connection)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GraphObjectFactory:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}




type Graph struct {
	*Node
	name string
}

func DefaultGraph() *Graph {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(Graph{})

	newGraph := Graph{
		Node: DefaultNode(),
	}
	newGraph.EntityKind = tgdb.EntityKindGraph
	return &newGraph
}

func NewGraph(gmd *GraphMetadata) *Graph {
	newGraph := DefaultGraph()
	newGraph.graphMetadata = gmd
	return newGraph
}

func NewGraphWithName(gmd *GraphMetadata, name string) *Graph {
	newGraph := NewGraph(gmd)
	newGraph.name = name
	return newGraph
}

/////////////////////////////////////////////////////////////////
// Helper functions for Graph
/////////////////////////////////////////////////////////////////

func (obj *Graph) GetModifiedAttributes() []tgdb.TGAttribute {
	return obj.getModifiedAttributes()
}

func (obj *Graph) GetName() string {
	return obj.name
}

func (obj *Graph) SetName(name string) {
	obj.name = name
}

// TODO: Revisit later - Ask TGDB Engineering Team as to if-n-how implement these methods
/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGGraph
/////////////////////////////////////////////////////////////////

func (obj *Graph) AddNode(node tgdb.TGNode) (tgdb.TGGraph, tgdb.TGError) {
	return obj, nil
}

func (obj *Graph) AddEdges(edges []tgdb.TGEdge) (tgdb.TGGraph, tgdb.TGError) {
	return obj, nil
}

func (obj *Graph) GetNode(filter tgdb.TGFilter) (tgdb.TGNode, tgdb.TGError) {
	return nil, nil
}

func (obj *Graph) ListNodes(filter tgdb.TGFilter, recurseAllSubGraphs bool) (tgdb.TGNode, tgdb.TGError) {
	return nil, nil
}

func (obj *Graph) CreateGraph(name string) (tgdb.TGGraph, tgdb.TGError) {
	return nil, nil
}

func (obj *Graph) RemoveGraph(name string) (tgdb.TGGraph, tgdb.TGError) {
	return nil, nil
}

func (obj *Graph) RemoveNode(node tgdb.TGNode) (tgdb.TGGraph, tgdb.TGError) {
	return nil, nil
}

func (obj *Graph) RemoveNodes(filter tgdb.TGFilter) int {
	return 0
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNode
/////////////////////////////////////////////////////////////////

func (obj *Graph) AddEdge(edge tgdb.TGEdge) {
	obj.Edges = append(obj.Edges, edge)
}

func (obj *Graph) AddEdgeWithDirectionType(node tgdb.TGNode, edgeType tgdb.TGEdgeType, directionType tgdb.TGDirectionType) tgdb.TGEdge {
	newEdge := NewEdgeWithDirection(obj.graphMetadata, obj, node, directionType)
	obj.AddEdge(newEdge)
	return newEdge
}

func (obj *Graph) GetEdges() []tgdb.TGEdge {
	return obj.Edges
}

func (obj *Graph) GetEdgesForDirectionType(directionType tgdb.TGDirectionType) []tgdb.TGEdge {
	edgesWithDirections := make([]tgdb.TGEdge, 0)
	if len(obj.Edges) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning Graph:GetEdgesForDirectionType as there are NO edges"))
		return edgesWithDirections
	}

	for _, edge := range obj.Edges {
		if edge.(*Edge).directionType == directionType {
			edgesWithDirections = append(edgesWithDirections, edge)
		}
	}
	return edgesWithDirections
}

func (obj *Graph) GetEdgesForEdgeType(edgeType tgdb.TGEdgeType, direction tgdb.TGDirection) []tgdb.TGEdge {
	edgesWithDirections := make([]tgdb.TGEdge, 0)
	if len(obj.Edges) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning Graph:GetEdgesForEdgeType as there are NO edges"))
		return edgesWithDirections
	}

	if edgeType == nil && direction == tgdb.DirectionAny {
		for _, edge := range obj.Edges {
			if edge.(*Edge).GetIsInitialized() {
				edgesWithDirections = append(edgesWithDirections, edge)
			}
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning Graph:GetEdgesForEdgeType w/ all edges of ANY directions"))
		}
		return obj.Edges
	}

	for _, edge := range obj.Edges {
		if !edge.(*Edge).GetIsInitialized() {
			logger.Warning(fmt.Sprintf("WARNING: Returning Graph:GetEdgesForEdgeType - skipping uninitialized edge '%+v'", edge))
			continue
		}
		eType := edge.GetEntityType()
		if edgeType != nil && eType != nil && eType.GetName() != edgeType.GetName() {
			logger.Warning(fmt.Sprintf("WARNING: Returning Graph:GetEdgesForEdgeType - skipping (entity type NOT matching) edge '%+v'", edge))
			continue
		}
		if direction == tgdb.DirectionAny {
			edgesWithDirections = append(edgesWithDirections, edge)
		} else if direction == tgdb.DirectionOutbound {
			edgesForThisNode := edge.GetVertices()
			if obj.GetVirtualId() == edgesForThisNode[0].GetVirtualId() {
				edgesWithDirections = append(edgesWithDirections, edge)
			}
		} else {
			edgesForThisNode := edge.GetVertices()
			if obj.GetVirtualId() == edgesForThisNode[1].GetVirtualId() {
				edgesWithDirections = append(edgesWithDirections, edge)
			}
		}
	}
	return edgesWithDirections
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntity
/////////////////////////////////////////////////////////////////

// GetAttribute gets the attribute for the Name specified
func (obj *Graph) GetAttribute(attrName string) tgdb.TGAttribute {
	return obj.getAttribute(attrName)
}

// GetAttributes lists of all the Attributes set
func (obj *Graph) GetAttributes() ([]tgdb.TGAttribute, tgdb.TGError) {
	return obj.getAttributes()
}

// GetEntityKind returns the EntityKind as a constant
func (obj *Graph) GetEntityKind() tgdb.TGEntityKind {
	return obj.getEntityKind()
}

// GetEntityType returns the EntityType
func (obj *Graph) GetEntityType() tgdb.TGEntityType {
	return obj.getEntityType()
}

// GetGraphMetadata returns the Graph Meta Data	- New in GO Lang
func (obj *Graph) GetGraphMetadata() tgdb.TGGraphMetadata {
	return obj.getGraphMetadata()
}

// GetIsDeleted checks whether this entity is already deleted in the system or not
func (obj *Graph) GetIsDeleted() bool {
	return obj.getIsDeleted()
}

// GetIsNew checks whether this entity that is currently being added to the system is new or not
func (obj *Graph) GetIsNew() bool {
	return obj.getIsNew()
}

// GetVersion gets the version of the Entity
func (obj *Graph) GetVersion() int {
	return obj.getVersion()
}

// GetVirtualId gets Entity identifier
// At the time of creation before reaching the server, it is the virtual id
// Upon successful creation, server returns a valid entity id that gets set in place of virtual id
func (obj *Graph) GetVirtualId() int64 {
	return obj.getVirtualId()
}

// IsAttributeSet checks whether this entity is an Attribute set or not
func (obj *Graph) IsAttributeSet(attrName string) bool {
	return obj.isAttributeSet(attrName)
}

// ResetModifiedAttributes resets the dirty flag on Attributes
func (obj *Graph) ResetModifiedAttributes() {
	obj.resetModifiedAttributes()
}

// SetAttribute associates the specified Attribute to this Entity
func (obj *Graph) SetAttribute(attr tgdb.TGAttribute) tgdb.TGError {
	return obj.setAttribute(attr)
}

// SetOrCreateAttribute dynamically associates the attribute to this entity
// If the AttributeDescriptor doesn't exist in the database, create a new one
func (obj *Graph) SetOrCreateAttribute(name string, value interface{}) tgdb.TGError {
	return obj.setOrCreateAttribute(name, value)
}

// SetEntityId sets Entity id and reset Virtual id after creation
func (obj *Graph) SetEntityId(id int64) {
	obj.setEntityId(id)
}

// SetIsDeleted set the deleted flag
func (obj *Graph) SetIsDeleted(flag bool) {
	obj.setIsDeleted(flag)
}

// SetIsInitialized set the initialized flag
func (obj *Graph) SetIsInitialized(flag bool) {
	obj.setIsInitialized(flag)
}

// SetIsNew sets the flag that this is a new entity
func (obj *Graph) SetIsNew(flag bool) {
	obj.setIsNew(flag)
}

// SetVersion sets the version of the Entity
func (obj *Graph) SetVersion(version int) {
	obj.setVersion(version)
}

func (obj *Graph) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Graph:{")
	buffer.WriteString(fmt.Sprintf("Name: %+v", obj.name))
	//buffer.WriteString(fmt.Sprintf(", Edges: %+v", obj.Edges))
	strArray := []string{buffer.String(), obj.entityToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *Graph) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering Graph:ReadExternal"))
	}
	// TODO: Revisit later - Do we need to validate length?
	nodeBufLen, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside Graph:ReadExternal read nodeBufLen as '%+v'", nodeBufLen))
	}

	err = obj.AbstractEntityReadExternal(is)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside Graph:ReadExternal read abstractEntity"))
	}

	edgeCount, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Graph:ReadExternal - unable to read edgeCount w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside Graph:ReadExternal read edgeCount as '%d'", edgeCount))
	}
	for i := 0; i < edgeCount; i++ {
		edgeId, err := is.(*ProtocolDataInputStream).ReadLong()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning Graph:ReadExternal - unable to read entId w/ Error: '%+v'", err.Error()))
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Graph:ReadExternal read edgeId as '%d'", edgeId))
		}
		var edge *Edge
		var entity tgdb.TGEntity
		refMap := is.(*ProtocolDataInputStream).GetReferenceMap()
		if refMap != nil {
			entity = refMap[edgeId]
		}
		if entity == nil {
			edge1 := NewEdge(obj.graphMetadata)
			edge1.SetEntityId(edgeId)
			edge1.SetIsInitialized(false)
			if refMap != nil {
				refMap[edgeId] = edge1
			}
			edge = edge1
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside Graph:ReadExternal created new edge: '%+v'", edge))
			}
		} else {
			edge = entity.(*Edge)
		}
		obj.Edges = append(obj.Edges, edge)
	}

	obj.SetIsInitialized(true)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning Graph:ReadExternal w/ NO error, for graph: '%+v'", obj))
	}
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *Graph) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering Graph:WriteExternal"))
	}
	startPos := os.(*ProtocolDataOutputStream).GetPosition()
	os.(*ProtocolDataOutputStream).WriteInt(0)
	// Write Attributes from the base class
	err := obj.AbstractEntityWriteExternal(os)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside Graph:WriteExternal - exported base entity Attributes"))
	}
	newCount := 0
	for _, edge := range obj.Edges {
		if edge.GetIsNew() {
			newCount++
		}
	}
	os.(*ProtocolDataOutputStream).WriteInt(newCount)
	// Write the edges ids - ONLY include new edges
	for _, edge := range obj.Edges {
		if ! edge.GetIsNew() {
			continue
		}
		os.(*ProtocolDataOutputStream).WriteLong(obj.GetVirtualId())
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Graph:WriteExternal - exported a new edge: '%+v'", edge))
		}
	}
	currPos := os.(*ProtocolDataOutputStream).GetPosition()
	length := currPos - startPos
	_, err = os.(*ProtocolDataOutputStream).WriteIntAt(startPos, length)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Graph:WriteExternal - unable to update data length in the buffer w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning Graph:WriteExternal w/ NO error, for graph: '%+v'", obj))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *Graph) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.isNew, obj.EntityKind, obj.VirtualId, obj.Version, obj.EntityId, obj.EntityType,
		obj.isDeleted, obj.isInitialized, obj.graphMetadata, obj.Attributes, obj.ModifiedAttributes, obj.Edges, obj.name)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Graph:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *Graph) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.isNew, &obj.EntityKind, &obj.VirtualId, &obj.Version, &obj.EntityId, &obj.EntityType,
		&obj.isDeleted, &obj.isInitialized, &obj.graphMetadata, &obj.Attributes, &obj.ModifiedAttributes, &obj.Edges, &obj.name)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Graph:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}




type GraphManager struct {
	name string
}

func NewGraphManager(gmd GraphMetadata) GraphManager {
	newGraphManager := GraphManager{
		name: "TGDB Graph Manager",
	}
	return newGraphManager
}

///////////////////////////////////////
// Helper functions for GraphManager //
///////////////////////////////////////

// GetName gets Graph Manager's Name
func (obj *GraphManager) GetName() string {
	return obj.name
}

///////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGGraphManager //
///////////////////////////////////////////////////////////

// CreateNode creates Node within this Graph. There is a default Root Graph.
func (obj *GraphManager) CreateNode() (tgdb.TGNode, tgdb.TGError) {
	return nil, nil
}

// CreateNodeForNodeType creates Node of particular Type
func (obj *GraphManager) CreateNodeForNodeType(nodeType tgdb.TGNodeType) (tgdb.TGNode, tgdb.TGError) {
	return nil, nil
}

// CreateEdge creates an Edge
func (obj *GraphManager) CreateEdge(fromNode tgdb.TGNode, toNode tgdb.TGNode, edgeType int) (tgdb.TGEdge, tgdb.TGError) {
	return nil, nil
}

// CreateEdgeWithDirection creates an Edge with direction
func (obj *GraphManager) CreateEdgeWithDirection(fromNode tgdb.TGNode, toNode tgdb.TGNode, directionType tgdb.TGDirectionType) (tgdb.TGEdge, tgdb.TGError) {
	return nil, nil
}

// CreateGraph creates a SubGraph at the Root level.
func (obj *GraphManager) CreateGraph(name string) (tgdb.TGGraph, tgdb.TGError) {
	return nil, nil
}

// DeleteNode removes this node from the graph
func (obj *GraphManager) DeleteNode(filter tgdb.TGFilter) (tgdb.TGGraphManager, tgdb.TGError) {
	return nil, nil
}

// DeleteNodes removes the nodes from this graph that match the filter
func (obj *GraphManager) DeleteNodes(filter tgdb.TGFilter) (tgdb.TGGraphManager, tgdb.TGError) {
	return nil, nil
}

// CreateQuery creates a Reusable Query
func (obj *GraphManager) CreateQuery(filter tgdb.TGFilter) tgdb.TGQuery {
	return nil
}

// QueryNodes gets Nodes based on the Filter condition with a set of Arguments
func (obj *GraphManager) QueryNodes(filter tgdb.TGFilter, args ...interface{}) tgdb.TGResultSet {
	return nil
}

// Traverse follows the graph using the traversal descriptor
func (obj *GraphManager) Traverse(descriptor tgdb.TGTraversalDescriptor, startingPoints []tgdb.TGNode) tgdb.TGResultSet {
	return nil
}

// GetGraphMetadata gets the Graph Metadata
func (obj *GraphManager) GetGraphMetadata() tgdb.TGGraphMetadata {
	return nil
}



