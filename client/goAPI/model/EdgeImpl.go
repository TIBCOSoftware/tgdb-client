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
 * File name: TEEdge.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type Edge struct {
	*AbstractEntity
	directionType types.TGDirectionType
	fromNode      types.TGNode
	toNode        types.TGNode
}

func DefaultEdge() *Edge {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(Edge{})

	newEdge := Edge{
		AbstractEntity: DefaultAbstractEntity(),
	}
	newEdge.EntityKind = types.EntityKindEdge
	newEdge.EntityType = DefaultEdgeType()
	return &newEdge
}

func NewEdge(gmd *GraphMetadata) *Edge {
	newEdge := DefaultEdge()
	newEdge.graphMetadata = gmd
	return newEdge
}

func NewEdgeWithDirection(gmd *GraphMetadata, fromNode types.TGNode, toNode types.TGNode, directionType types.TGDirectionType) *Edge {
	newEdge := NewEdge(gmd)
	newEdge.directionType = directionType
	newEdge.fromNode = fromNode
	newEdge.toNode = toNode
	return newEdge
}

func NewEdgeWithEdgeType(gmd *GraphMetadata, fromNode types.TGNode, toNode types.TGNode, edgeType types.TGEdgeType) *Edge {
	newEdge := NewEdge(gmd)
	newEdge.fromNode = fromNode
	newEdge.toNode = toNode
	newEdge.EntityType = edgeType
	newEdge.directionType = edgeType.GetDirectionType()
	return newEdge
}

/////////////////////////////////////////////////////////////////
// Helper functions for Edge
/////////////////////////////////////////////////////////////////

func (obj *Edge) GetFromNode() types.TGNode {
	return obj.fromNode
}

func (obj *Edge) GetIsInitialized() bool {
	return obj.isInitialized
}

func (obj *Edge) GetModifiedAttributes() []types.TGAttribute {
	return obj.getModifiedAttributes()
}

func (obj *Edge) GetToNode() types.TGNode {
	return obj.toNode
}

func (obj *Edge) SetDirectionType(dirType types.TGDirectionType) {
	obj.directionType = dirType
}

func (obj *Edge) SetFromNode(node types.TGNode) {
	obj.fromNode = node
}

func (obj *Edge) SetToNode(node types.TGNode) {
	obj.toNode = node
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEdge
/////////////////////////////////////////////////////////////////

// GetDirectionType gets direction type as one of the constants
func (obj *Edge) GetDirectionType() types.TGDirectionType {
	if obj.EntityType != nil {
		return obj.EntityType.(*EdgeType).GetDirectionType()
	} else {
		return obj.directionType
	}
}

// GetVertices gets array of NODE (Entity) types for this EDGE (Entity) type
func (obj *Edge) GetVertices() []types.TGNode {
	return []types.TGNode{obj.fromNode, obj.toNode}
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntity
/////////////////////////////////////////////////////////////////

// GetAttribute gets the attribute for the name specified
func (obj *Edge) GetAttribute(attrName string) types.TGAttribute {
	return obj.getAttribute(attrName)
}

// GetAttributes lists of all the attributes set
func (obj *Edge) GetAttributes() ([]types.TGAttribute, types.TGError) {
	return obj.getAttributes()
}

// GetEntityKind returns the EntityKind as a constant
func (obj *Edge) GetEntityKind() types.TGEntityKind {
	return obj.getEntityKind()
}

// GetEntityType returns the EntityType
func (obj *Edge) GetEntityType() types.TGEntityType {
	return obj.getEntityType()
}

// GetGraphMetadata returns the Graph Meta Data	- New in GO Lang
func (obj *Edge) GetGraphMetadata() types.TGGraphMetadata {
	return obj.getGraphMetadata()
}

// GetIsDeleted checks whether this entity is already deleted in the system or not
func (obj *Edge) GetIsDeleted() bool {
	return obj.getIsDeleted()
}

// GetIsNew checks whether this entity that is currently being added to the system is new or not
func (obj *Edge) GetIsNew() bool {
	return obj.getIsNew()
}

// GetVersion gets the version of the Entity
func (obj *Edge) GetVersion() int {
	return obj.getVersion()
}

// GetVirtualId gets Entity identifier
// At the time of creation before reaching the server, it is the virtual id
// Upon successful creation, server returns a valid entity id that gets set in place of virtual id
func (obj *Edge) GetVirtualId() int64 {
	return obj.getVirtualId()
}

// IsAttributeSet checks whether this entity is an Attribute set or not
func (obj *Edge) IsAttributeSet(attrName string) bool {
	return obj.isAttributeSet(attrName)
}

// ResetModifiedAttributes resets the dirty flag on attributes
func (obj *Edge) ResetModifiedAttributes() {
	obj.resetModifiedAttributes()
}

// SetAttribute associates the specified Attribute to this Entity
func (obj *Edge) SetAttribute(attr types.TGAttribute) types.TGError {
	return obj.setAttribute(attr)
}

// SetOrCreateAttribute dynamically associates the attribute to this entity
// If the AttributeDescriptor doesn't exist in the database, create a new one
func (obj *Edge) SetOrCreateAttribute(name string, value interface{}) types.TGError {
	return obj.setOrCreateAttribute(name, value)
}

// SetEntityId sets Entity id and reset Virtual id after creation
func (obj *Edge) SetEntityId(id int64) {
	obj.setEntityId(id)
}

// SetIsDeleted set the deleted flag
func (obj *Edge) SetIsDeleted(flag bool) {
	obj.setIsDeleted(flag)
}

// SetIsInitialized set the initialized flag
func (obj *Edge) SetIsInitialized(flag bool) {
	obj.setIsInitialized(flag)
}

// SetIsNew sets the flag that this is a new entity
func (obj *Edge) SetIsNew(flag bool) {
	obj.setIsNew(flag)
}

// SetVersion sets the version of the Entity
func (obj *Edge) SetVersion(version int) {
	obj.setVersion(version)
}

func (obj *Edge) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Edge:{")
	buffer.WriteString(fmt.Sprintf("DirectionType: %+v", obj.directionType))
	if obj.fromNode != nil {
		buffer.WriteString(fmt.Sprintf(", FromNode: %+v", obj.fromNode.GetVirtualId()))
	}
	if obj.toNode != nil {
		buffer.WriteString(fmt.Sprintf(", ToNode: %+v", obj.toNode.GetVirtualId()))
	}
	strArray := []string{buffer.String(), obj.entityToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *Edge) ReadExternal(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering Edge:ReadExternal"))
	// TODO: Revisit later - Do we need to validate length?
	edgeBufLen, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprintf("Inside Edge:ReadExternal read edgeBufLen as '%+v'", edgeBufLen))

	err = obj.AbstractEntityReadExternal(is)
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprint("Inside Edge:ReadExternal read abstractEntity"))

	direction, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:ReadExternal - unable to read direction w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside Edge:ReadExternal read direction as '%+v'", direction))
	if direction == 0 {
		obj.SetDirectionType(types.DirectionTypeUnDirected)
	} else if direction == 1 {
		obj.SetDirectionType(types.DirectionTypeDirected)
	} else {
		obj.SetDirectionType(types.DirectionTypeBiDirectional)
	}

	var fromEntity, toEntity types.TGEntity
	var fromNode, toNode *Node

	fromNodeId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:ReadExternal - unable to read fromNodeId w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside Edge:ReadExternal read fromNodeId as '%d'", fromNodeId))
	refMap := is.(*iostream.ProtocolDataInputStream).GetReferenceMap()
	if refMap != nil {
		fromEntity = refMap[fromNodeId]
	}
	if fromEntity == nil {
		fNode := NewNode(obj.graphMetadata)
		fNode.SetEntityId(fromNodeId)
		fNode.SetIsInitialized(false)
		if refMap != nil {
			refMap[fromNodeId] = fNode
		}
		fromNode = fNode
		logger.Log(fmt.Sprintf("Inside Edge:ReadExternal created new fromNode: '%+v'", fromNode))
	} else {
		fromNode = fromEntity.(*Node)
	}
	obj.SetFromNode(fromNode)
	logger.Log(fmt.Sprintf("Inside Edge:ReadExternal Edge has fromNode & StreamEntityCount is '%d'", len(is.(*iostream.ProtocolDataInputStream).GetReferenceMap())))

	toNodeId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:ReadExternal - unable to read toNodeId w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside Edge:ReadExternal read toNodeId as '%d'", toNodeId))
	if refMap != nil {
		toEntity = refMap[toNodeId]
	}
	if toEntity == nil {
		tNode := NewNode(obj.graphMetadata)
		tNode.SetEntityId(toNodeId)
		tNode.SetIsInitialized(false)
		if refMap != nil {
			refMap[toNodeId] = tNode
		}
		toNode = tNode
		logger.Log(fmt.Sprintf("Inside Edge:ReadExternal created new toNode: '%+v'", toNode))
	} else {
		toNode = toEntity.(*Node)
	}
	obj.SetToNode(toNode)
	logger.Log(fmt.Sprintf("Inside Edge:ReadExternal Edge has toNode & StreamEntityCount is '%d'", len(is.(*iostream.ProtocolDataInputStream).GetReferenceMap())))

	obj.SetIsInitialized(true)
	logger.Log(fmt.Sprintf("Returning Edge:ReadExternal w/ NO error, for edge: '%+v'", obj))
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *Edge) WriteExternal(os types.TGOutputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering Edge:WriteExternal"))
	startPos := os.(*iostream.ProtocolDataOutputStream).GetPosition()
	os.(*iostream.ProtocolDataOutputStream).WriteInt(0)
	// Write attributes from the base class
	err := obj.AbstractEntityWriteExternal(os)
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprint("Inside Edge:WriteExternal - exported base entity attributes"))
	// TODO: Revisit later - Check w/ TGDB Engineering team as what the difference should be for if-n-else conditions
	// Write the edges ids
	if obj.GetIsNew() {
		os.(*iostream.ProtocolDataOutputStream).WriteByte(int(obj.GetDirectionType()))
		os.(*iostream.ProtocolDataOutputStream).WriteLong(obj.GetFromNode().GetVirtualId())
		os.(*iostream.ProtocolDataOutputStream).WriteLong(obj.GetToNode().GetVirtualId())
	} else {
		os.(*iostream.ProtocolDataOutputStream).WriteByte(int(obj.GetDirectionType())) // The Server expects it - so better send it.
		os.(*iostream.ProtocolDataOutputStream).WriteLong(obj.GetFromNode().GetVirtualId())
		os.(*iostream.ProtocolDataOutputStream).WriteLong(obj.GetToNode().GetVirtualId())
	}
	currPos := os.(*iostream.ProtocolDataOutputStream).GetPosition()
	length := currPos - startPos
	_, err = os.(*iostream.ProtocolDataOutputStream).WriteIntAt(startPos, length)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:WriteExternal - unable to update data length in the buffer w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Returning Edge:WriteExternal w/ NO error, for edge: '%+v'", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *Edge) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.isNew, obj.EntityKind, obj.virtualId, obj.version, obj.entityId, obj.EntityType,
		obj.isDeleted, obj.isInitialized, obj.graphMetadata, obj.attributes, obj.modifiedAttributes,
		obj.directionType, obj.fromNode, obj.toNode)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *Edge) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.isNew, &obj.EntityKind, &obj.virtualId, &obj.version, &obj.entityId, &obj.EntityType,
		&obj.isDeleted, &obj.isInitialized, &obj.graphMetadata, &obj.attributes, &obj.modifiedAttributes,
		&obj.directionType, &obj.fromNode, &obj.toNode)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
