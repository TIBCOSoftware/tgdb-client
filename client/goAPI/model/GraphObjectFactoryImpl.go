package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
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
 * File name: TGGraphObjectFactory.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type GraphObjectFactory struct {
	graphMData *GraphMetadata
	connection types.TGConnection
}

func DefaultGraphObjectFactory() *GraphObjectFactory {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GraphObjectFactory{})

	newGraphObjectFactory := GraphObjectFactory{}
	return &newGraphObjectFactory
}

func NewGraphObjectFactory(conn types.TGConnection) *GraphObjectFactory {
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

func (obj *GraphObjectFactory) GetConnection() types.TGConnection {
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
func (obj *GraphObjectFactory) CreateCompositeKey(nodeTypeName string) (types.TGKey, types.TGError) {
	_, err := obj.graphMData.GetNodeType(nodeTypeName)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateCompositeKey as node NOT found"))
		errMsg := fmt.Sprintf("Node desc with name %s not found", nodeTypeName)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return NewCompositeKey(obj.graphMData, nodeTypeName), nil
}

// CreateEdgeWithEdgeType creates an Edge
func (obj *GraphObjectFactory) CreateEdgeWithEdgeType(fromNode types.TGNode, toNode types.TGNode, edgeType types.TGEdgeType) (types.TGEdge, types.TGError) {
	newEdge := NewEdgeWithEdgeType(obj.graphMData, fromNode, toNode, edgeType)
	if newEdge.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateEdgeWithEdgeType as edge in NOT initialized"))
		errMsg := fmt.Sprintf("Unable to create an edge with type %s", edgeType.(*EdgeType).name)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	fromNode.AddEdge(newEdge)
	toNode.AddEdge(newEdge)
	return newEdge, nil
}

// CreateEdgeWithDirection creates an Edge with a direction
func (obj *GraphObjectFactory) CreateEdgeWithDirection(fromNode types.TGNode, toNode types.TGNode, directionType types.TGDirectionType) (types.TGEdge, types.TGError) {
	newEdge := NewEdgeWithDirection(obj.graphMData, fromNode, toNode, directionType)
	if newEdge.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateEdgeWithDirection as edge in NOT initialized"))
		errMsg := fmt.Sprintf("Unable to create an edge with direction %s", directionType)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	fromNode.AddEdge(newEdge)
	toNode.AddEdge(newEdge)
	return newEdge, nil
}

// CreateEntity creates entity based on the entity kind specified
func (obj *GraphObjectFactory) CreateEntity(entityKind types.TGEntityKind) (types.TGEntity, types.TGError) {
	switch entityKind {
	case types.EntityKindNode:
		return NewNode(obj.graphMData), nil
	case types.EntityKindEdge:
		return NewEdge(obj.graphMData), nil
	case types.EntityKindGraph:
		return NewGraph(obj.graphMData), nil
	}
	logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateEntity as entity kind specified is INVALID"))
	errMsg := fmt.Sprintf("Invalid entity kind %d specified", entityKind)
	return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
}

// CreateEntityId creates entity id from input buffer
func (obj *GraphObjectFactory) CreateEntityId(buf []byte) (types.TGEntityId, types.TGError) {
	return NewByteArrayEntity(buf), nil
}

// CreateGraph creates a Graph
func (obj *GraphObjectFactory) CreateGraph(name string) (types.TGGraph, types.TGError) {
	graph := NewGraphWithName(obj.graphMData, name)
	if graph.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateGraph as graph in NOT initialized"))
		errMsg := fmt.Sprint("Unable to create a graph with this Graph Object Factory")
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return graph, nil
}

// CreateNode creates a Node
func (obj *GraphObjectFactory) CreateNode() (types.TGNode, types.TGError) {
	node := NewNode(obj.graphMData)
	if node.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateNode as node in NOT initialized"))
		errMsg := fmt.Sprint("Unable to create a node with this Graph Object Factory")
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return node, nil
}

// CreateNodeInGraph creates Node within this Graph. There is a default Root Graph.
func (obj *GraphObjectFactory) CreateNodeInGraph(nodeType types.TGNodeType) (types.TGNode, types.TGError) {
	node := NewNodeWithType(obj.graphMData, nodeType)
	if node.isInitialized != true {
		logger.Error(fmt.Sprint("ERROR: Returning GraphObjectFactory:CreateNodeInGraph as node in NOT initialized"))
		errMsg := fmt.Sprint("Unable to create a node with this Graph Object Factory")
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
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
