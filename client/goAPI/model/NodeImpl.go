package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"strings"
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
 * File name: TGNode.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type Node struct {
	*AbstractEntity
	edges []types.TGEdge
}

func DefaultNode() *Node {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(Node{})

	newNode := Node{
		AbstractEntity: DefaultAbstractEntity(),
	}
	newNode.EntityKind = types.EntityKindNode
	newNode.EntityType = DefaultNodeType()
	newNode.edges = make([]types.TGEdge, 0)
	return &newNode
}

func NewNode(gmd *GraphMetadata) *Node {
	newNode := DefaultNode()
	newNode.graphMetadata = gmd
	return newNode
}

func NewNodeWithType(gmd *GraphMetadata, nodeType types.TGNodeType) *Node {
	newNode := NewNode(gmd)
	newNode.EntityType = nodeType
	return newNode
}

/////////////////////////////////////////////////////////////////
// Helper functions for Node
/////////////////////////////////////////////////////////////////

func (obj *Node) GetIsInitialized() bool {
	return obj.isInitialized
}

func (obj *Node) GetModifiedAttributes() []types.TGAttribute {
	return obj.getModifiedAttributes()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNode
/////////////////////////////////////////////////////////////////

func (obj *Node) AddEdge(edge types.TGEdge) {
	obj.edges = append(obj.edges, edge)
}

func (obj *Node) AddEdgeWithDirectionType(node types.TGNode, edgeType types.TGEdgeType, directionType types.TGDirectionType) types.TGEdge {
	newEdge := NewEdgeWithDirection(obj.graphMetadata, obj, node, directionType)
	obj.AddEdge(newEdge)
	return newEdge
}

func (obj *Node) GetEdges() []types.TGEdge {
	return obj.edges
}

func (obj *Node) GetEdgesForDirectionType(directionType types.TGDirectionType) []types.TGEdge {
	edgesWithDirections := make([]types.TGEdge, 0)
	if len(obj.edges) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning Node:GetEdgesForDirectionType as there are NO edges"))
		return edgesWithDirections
	}

	for _, edge := range obj.edges {
		if edge.(*Edge).directionType == directionType {
			edgesWithDirections = append(edgesWithDirections, edge)
		}
	}
	return edgesWithDirections
}

func (obj *Node) GetEdgesForEdgeType(edgeType types.TGEdgeType, direction types.TGDirection) []types.TGEdge {
	edgesWithDirections := make([]types.TGEdge, 0)
	if len(obj.edges) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning Node:GetEdgesForEdgeType as there are NO edges"))
		return edgesWithDirections
	}

	if edgeType == nil && direction == types.DirectionAny {
		for _, edge := range obj.edges {
			if edge.(*Edge).GetIsInitialized() {
				edgesWithDirections = append(edgesWithDirections, edge)
			}
		}
		return obj.edges
	}

	for _, edge := range obj.edges {
		if !edge.(*Edge).GetIsInitialized() {
			logger.Warning(fmt.Sprintf("WARNING: Continuing loop Node:GetEdgesForEdgeType - skipping uninitialized edge '%+v'", edge))
			continue
		}
		eType := edge.GetEntityType()
		if edgeType != nil && eType != nil && eType.GetName() != edgeType.GetName() {
			logger.Warning(fmt.Sprintf("WARNING: Continuing loop Node:GetEdgesForEdgeType - skipping (entity type NOT matching) edge '%+v'", edge))
			continue
		}
		if direction == types.DirectionAny {
			edgesWithDirections = append(edgesWithDirections, edge)
		} else if direction == types.DirectionOutbound {
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

// GetAttribute gets the attribute for the name specified
func (obj *Node) GetAttribute(attrName string) types.TGAttribute {
	return obj.getAttribute(attrName)
}

// GetAttributes lists of all the attributes set
func (obj *Node) GetAttributes() ([]types.TGAttribute, types.TGError) {
	return obj.getAttributes()
}

// GetEntityKind returns the EntityKind as a constant
func (obj *Node) GetEntityKind() types.TGEntityKind {
	return obj.getEntityKind()
}

// GetEntityType returns the EntityType
func (obj *Node) GetEntityType() types.TGEntityType {
	return obj.getEntityType()
}

// GetGraphMetadata returns the Graph Meta Data	- New in GO Lang
func (obj *Node) GetGraphMetadata() types.TGGraphMetadata {
	return obj.getGraphMetadata()
}

// GetIsDeleted checks whether this entity is already deleted in the system or not
func (obj *Node) GetIsDeleted() bool {
	return obj.getIsDeleted()
}

// GetIsNew checks whether this entity that is currently being added to the system is new or not
func (obj *Node) GetIsNew() bool {
	return obj.getIsNew()
}

// GetVersion gets the version of the Entity
func (obj *Node) GetVersion() int {
	return obj.getVersion()
}

// GetVirtualId gets Entity identifier
// At the time of creation before reaching the server, it is the virtual id
// Upon successful creation, server returns a valid entity id that gets set in place of virtual id
func (obj *Node) GetVirtualId() int64 {
	return obj.getVirtualId()
}

// IsAttributeSet checks whether this entity is an Attribute set or not
func (obj *Node) IsAttributeSet(attrName string) bool {
	return obj.isAttributeSet(attrName)
}

// ResetModifiedAttributes resets the dirty flag on attributes
func (obj *Node) ResetModifiedAttributes() {
	obj.resetModifiedAttributes()
}

// SetAttribute associates the specified Attribute to this Entity
func (obj *Node) SetAttribute(attr types.TGAttribute) types.TGError {
	return obj.setAttribute(attr)
}

// SetOrCreateAttribute dynamically associates the attribute to this entity
// If the AttributeDescriptor doesn't exist in the database, create a new one
func (obj *Node) SetOrCreateAttribute(name string, value interface{}) types.TGError {
	return obj.setOrCreateAttribute(name, value)
}

// SetEntityId sets Entity id and reset Virtual id after creation
func (obj *Node) SetEntityId(id int64) {
	obj.setEntityId(id)
}

// SetIsDeleted set the deleted flag
func (obj *Node) SetIsDeleted(flag bool) {
	obj.setIsDeleted(flag)
}

// SetIsInitialized set the initialized flag
func (obj *Node) SetIsInitialized(flag bool) {
	obj.setIsInitialized(flag)
}

// SetIsNew sets the flag that this is a new entity
func (obj *Node) SetIsNew(flag bool) {
	obj.setIsNew(flag)
}

// SetVersion sets the version of the Entity
func (obj *Node) SetVersion(version int) {
	obj.setVersion(version)
}

func (obj *Node) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Node:{")
	buffer.WriteString(fmt.Sprintf("Edges: %+v", obj.edges))
	strArray := []string{buffer.String(), obj.entityToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *Node) ReadExternal(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering Node:ReadExternal"))
	// TODO: Revisit later - Do we need to validate length?
	nodeBufLen, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:ReadExternal - unable to read length w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside Node:ReadExternal read nodeBufLen as '%+v'", nodeBufLen))

	err = obj.AbstractEntityReadExternal(is)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:ReadExternal - unable to obj.AbstractEntityReadExternal(is) w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprint("Inside Node:ReadExternal read abstractEntity"))

	edgeCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:ReadExternal - unable to read edgeCount w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside Node:ReadExternal read edgeCount as '%d'", edgeCount))
	for i := 0; i < edgeCount; i++ {
		edgeId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning Node:ReadExternal - unable to read entId w/ Error: '%+v'", err.Error()))
			return err
		}
		logger.Log(fmt.Sprintf("Inside Node:ReadExternal read edgeId as '%d'", edgeId))
		var edge *Edge
		var entity types.TGEntity
		refMap := is.(*iostream.ProtocolDataInputStream).GetReferenceMap()
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
			logger.Log(fmt.Sprintf("Inside Node:ReadExternal created new edge: '%+v'", edge))
		} else {
			edge = entity.(*Edge)
		}
		obj.edges = append(obj.edges, edge)
		logger.Log(fmt.Sprintf("Inside Node:ReadExternal Node has '%d' edges & StreamEntityCount is '%d'", len(obj.edges), len(is.(*iostream.ProtocolDataInputStream).GetReferenceMap())))
	}

	obj.SetIsInitialized(true)
	logger.Log(fmt.Sprintf("Returning Node:ReadExternal w/ NO error, for node: '%+v'", obj))
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *Node) WriteExternal(os types.TGOutputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering Node:WriteExternal"))
	startPos := os.(*iostream.ProtocolDataOutputStream).GetPosition()
	os.(*iostream.ProtocolDataOutputStream).WriteInt(0)
	// Write attributes from the base class
	err := obj.AbstractEntityWriteExternal(os)
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprint("Inside Node:WriteExternal - exported base entity attributes"))
	newCount := 0
	for _, edge := range obj.edges {
		if edge.GetIsNew() {
			newCount++
		}
	}
	os.(*iostream.ProtocolDataOutputStream).WriteInt(newCount)
	logger.Log(fmt.Sprintf("Inside Node:WriteExternal - exported new edge count '%d'", newCount))
	// Write the edges ids - ONLY include new edges
	for _, edge := range obj.edges {
		if ! edge.GetIsNew() {
			continue
		}
		os.(*iostream.ProtocolDataOutputStream).WriteLong(obj.GetVirtualId())
		logger.Log(fmt.Sprintf("Inside Node:WriteExternal - exported a new edge: '%+v'", edge))
	}
	currPos := os.(*iostream.ProtocolDataOutputStream).GetPosition()
	length := currPos - startPos
	_, err = os.(*iostream.ProtocolDataOutputStream).WriteIntAt(startPos, length)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:WriteExternal - unable to update data length in the buffer w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Returning Node:WriteExternal w/ NO error, for node: '%+v'", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *Node) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.isNew, obj.EntityKind, obj.virtualId, obj.version, obj.entityId, obj.EntityType,
		obj.isDeleted, obj.isInitialized, obj.graphMetadata, obj.attributes, obj.modifiedAttributes, obj.edges)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *Node) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.isNew, &obj.EntityKind, &obj.virtualId, &obj.version, &obj.entityId, &obj.EntityType,
		&obj.isDeleted, &obj.isInitialized, &obj.graphMetadata, &obj.attributes, &obj.modifiedAttributes, &obj.edges)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
