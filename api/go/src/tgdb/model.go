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
 * File name: model.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: model.go 3626 2019-12-09 19:35:03Z nimish $
 */

package tgdb

import (
	"bytes"
)

// ======= Various Entity Kind =======
type TGEntityKind int
// ======= System Types =======
type TGSystemType int

const (
	EntityKindInvalid   TGEntityKind = iota
	EntityKindEntity    //TGEntityKind = 1
	EntityKindNode      //TGEntityKind = 2
	EntityKindEdge      //TGEntityKind = 3
	EntityKindGraph     //TGEntityKind = 4
	EntityKindHyperEdge //TGEntityKind = 5
)

func (kind TGEntityKind) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if kind&EntityKindInvalid == EntityKindInvalid {
		buffer.WriteString("InvalidKind")
	} else if kind&EntityKindEntity == EntityKindEntity {
		buffer.WriteString("EntityKind")
	} else if kind&EntityKindNode == EntityKindNode {
		buffer.WriteString("NodeKind")
	} else if kind&EntityKindEdge == EntityKindEdge {
		buffer.WriteString("EdgeKind")
	} else if kind&EntityKindGraph == EntityKindGraph {
		buffer.WriteString("GraphKind")
	} else if kind&EntityKindHyperEdge == EntityKindHyperEdge {
		buffer.WriteString("HyperEdgeKind")
	}
	if buffer.Len() == 0 {
		return ""
	}

	return buffer.String()
}



type TGEntity interface {
	TGSerializable
	// GetAttribute gets the attribute for the name specified
	GetAttribute(name string) TGAttribute
	// GetAttributes lists of all the attributes set
	GetAttributes() ([]TGAttribute, TGError)
	// GetEntityKind returns the EntityKind as a constant
	GetEntityKind() TGEntityKind
	// GetEntityType returns the EntityType
	GetEntityType() TGEntityType
	// GetGraphMetadata returns the Graph Meta Data	- New in GO Lang
	GetGraphMetadata() TGGraphMetadata
	// GetIsDeleted checks whether this entity is already deleted in the system or not
	GetIsDeleted() bool
	// GetIsNew checks whether this entity that is currently being added to the system is new or not
	GetIsNew() bool
	// GetVersion gets the version of the Entity
	GetVersion() int
	// GetVirtualId gets Entity identifier
	// At the time of creation before reaching the server, it is the virtual id
	// Upon successful creation, server returns a valid entity id that gets set in place of virtual id
	GetVirtualId() int64
	// IsAttributeSet checks whether this entity is an Attribute set or not
	IsAttributeSet(attrName string) bool
	// ResetModifiedAttributes resets the dirty flag on attributes
	ResetModifiedAttributes()
	// SetAttribute associates the specified Attribute to this Entity
	SetAttribute(attr TGAttribute) TGError
	// SetOrCreateAttribute dynamically associates the attribute to this entity
	// If the AttributeDescriptor doesn't exist in the database, create a new one
	SetOrCreateAttribute(name string, value interface{}) TGError
	// SetEntityId sets Entity id and reset Virtual id after creation
	SetEntityId(id int64)
	// SetIsDeleted set the deleted flag
	SetIsDeleted(flag bool)
	// SetIsInitialized set the initialized flag
	SetIsInitialized(flag bool)
	// SetIsNew sets the flag that this is a new entity
	SetIsNew(flag bool)
	// SetVersion sets the version of the Entity
	SetVersion(version int)
	// Additional Method to help debugging
	String() string
}


// An Attribute is simple scalar value that is associated with an Entity.
type TGAttribute interface {
	TGSerializable
	// GetAttributeDescriptor returns the TGAttributeDescriptor for this attribute
	GetAttributeDescriptor() TGAttributeDescriptor
	// GetIsModified checks whether the attribute modified or not
	GetIsModified() bool
	// GetName gets the name for this attribute as the most generic form
	GetName() string
	// GetOwner gets owner Entity of this attribute
	GetOwner() TGEntity
	// GetValue gets the value for this attribute as the most generic form
	GetValue() interface{}
	// IsNull checks whether the attribute value is null or not
	IsNull() bool
	// ResetIsModified resets the IsModified flag - recursively, if needed
	ResetIsModified()
	// SetOwner sets the owner entity - Need this indirection to traverse the chain
	SetOwner(attrOwner TGEntity)
	// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
	// If the object is Null, then the object is explicitly set, but no value is provided.
	SetValue(value interface{}) TGError
	// ReadValue reads the attribute value from input stream
	ReadValue(is TGInputStream) TGError
	// WriteValue writes the attribute value to output stream
	WriteValue(os TGOutputStream) TGError
	// Additional Method to help debugging
	String() string
}


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
	// Is_Encrypted checks whether this attribute is Encrypted or not
	Is_Encrypted() bool
	// SetPrecision sets the prevision for Attribute Descriptor of type Number
	SetPrecision(precision int16)
	// SetScale sets the scale for Attribute Descriptor of type Number
	SetScale(scale int16)
}

type TGSystemObject interface {
	TGSerializable
	// GetName gets the system object's name
	GetName() string
	// GetSystemType gets the system object's type
	GetSystemType() TGSystemType
}



const (
	SystemTypeInvalid TGSystemType = -1
	SystemTypeEntity  TGSystemType = -2
)
const (
	SystemTypeAttributeDescriptor TGSystemType = iota
	SystemTypeNode
	SystemTypeEdge
	SystemTypeIndex
	SystemTypePrincipal
	SystemTypeRole
	SystemTypeSequence
	SystemTypeMaxSysObject
)

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

type TGKey interface {
	TGSerializable
	// Dynamically set the attribute to this entity. If the AttributeDescriptor doesn't exist in the database, create a new one.
	SetOrCreateAttribute(name string, value interface{}) TGError
}

type TGEdgeType interface {
	TGEntityType
	// GetDirectionType gets direction type as one of the constants
	GetDirectionType() TGDirectionType
	// GetFromNodeType gets From-Node Type
	GetFromNodeType() TGNodeType
	// GetFromTypeId gets From-Node ID
	GetFromTypeId() int
	// GetToNodeType gets To-Node Type
	GetToNodeType() TGNodeType
	// GetToTypeId gets To-Node ID
	GetToTypeId() int
	// SetFromNodeType sets From-Node Type
	SetFromNodeType(fromNode TGNodeType)
	// SetFromTypeId sets From-Node ID
	SetFromTypeId(fromTypeId int)
	// SetToNodeType sets From-Node Type
	SetToNodeType(toNode TGNodeType)
	// SetToTypeId sets To-Node ID
	SetToTypeId(toTypeId int)
}

type TGNodeType interface {
	TGEntityType
	// GetPKeyAttributeDescriptors returns a set of primary key descriptors
	GetPKeyAttributeDescriptors() []TGAttributeDescriptor
}


type TGDirectionType int

const (
	DirectionTypeUnDirected TGDirectionType = iota
	DirectionTypeDirected
	DirectionTypeBiDirectional
)

type TGGraphObjectFactory interface {
	// CreateCompositeKey creates a CompositeKey for a SystemTypeNode. The composite key can also be a single key
	CreateCompositeKey(nodeTypeName string) (TGKey, TGError)
	// CreateEdgeWithEdgeType creates an Edge
	CreateEdgeWithEdgeType(fromNode TGNode, toNode TGNode, edgeType TGEdgeType) (TGEdge, TGError)
	// CreateEdgeWithDirection creates an Edge with a direction
	CreateEdgeWithDirection(fromNode TGNode, toNode TGNode, directionType TGDirectionType) (TGEdge, TGError)
	// CreateEntity creates entity based on the entity kind specified
	CreateEntity(entityKind TGEntityKind) (TGEntity, TGError)
	// CreateEntityId creates entity id from input buffer
	CreateEntityId(buf []byte) (TGEntityId, TGError)
	// CreateGraph creates a Graph
	CreateGraph(name string) (TGGraph, TGError)
	// CreateNode creates a Node
	CreateNode() (TGNode, TGError)
	// CreateNodeInGraph creates Node within this Graph. There is a default Root Graph.
	CreateNodeInGraph(nodeType TGNodeType) (TGNode, TGError)
}

type TGNode interface {
	TGEntity
	// Add another Edge to Node
	AddEdge(edge TGEdge)
	// Add another Edge with direction to Node
	AddEdgeWithDirectionType(node TGNode, edgeType TGEdgeType, directionType TGDirectionType) TGEdge
	// Return entire collection edges
	GetEdges() []TGEdge
	// Return collection edges for the direction desc
	GetEdgesForDirectionType(directionType TGDirectionType) []TGEdge
	// Return collection filtered edges
	GetEdgesForEdgeType(edgeType TGEdgeType, direction TGDirection) []TGEdge
}

type TGEdge interface {
	TGEntity
	// GetDirectionType gets direction type as one of the constants
	GetDirectionType() TGDirectionType
	// GetVertices gets array of NODE (Entity) types for this EDGE (Entity) type
	GetVertices() []TGNode
}

type TGDirection int

const (
	DirectionInbound TGDirection = iota
	DirectionOutbound
	DirectionAny
)

type TGEntityId interface {
	TGSerializable
	// ToBytes converts the Entity Id in binary format
	ToBytes() ([]byte, TGError)
}


type TGGraph interface {
	TGNode
	// AddNode adds an EntityKindEdge to another EntityKindNode
	AddNode(node TGNode) (TGGraph, TGError)
	// AddEdges adds a collection of edges for this node
	AddEdges(edges []TGEdge) (TGGraph, TGError)
	// GetNode gets a unique node which matches the unique constraint
	GetNode(filter TGFilter) (TGNode, TGError)
	// ListNodes lists all the nodes that match the filter and recurse All sub graphs
	ListNodes(filter TGFilter, recurseAllSubGraphs bool) (TGNode, TGError)
	// CreateGraph creates a sub graph within this graph
	CreateGraph(name string) (TGGraph, TGError)
	// RemoveGraph removes the graph
	RemoveGraph(name string) (TGGraph, TGError)
	// RemoveNode removes this node from the graph
	RemoveNode(node TGNode) (TGGraph, TGError)
	// RemoveNodes removes a set of nodes from this graph that match the filter
	RemoveNodes(filter TGFilter) int
}

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


// ChangeListener is an event listener that gets triggered when that event occurs
type TGChangeListener interface {
	// AttributeAdded gets called when an attribute is Added to an entity.
	AttributeAdded(attr TGAttribute, owner TGEntity)
	// AttributeChanged gets called when an attribute is set.
	AttributeChanged(attr TGAttribute, oldValue, newValue interface{})
	// AttributeRemoved gets called when an attribute is removed from the entity.
	AttributeRemoved(attr TGAttribute, owner TGEntity)
	// EntityCreated gets called when an entity is Added
	EntityCreated(entity TGEntity)
	// EntityDeleted gets called when the entity is deleted
	EntityDeleted(entity TGEntity)
	// NodeAdded gets called when a node is Added
	NodeAdded(graph TGGraph, node TGNode)
	// NodeRemoved gets called when a node is removed
	NodeRemoved(graph TGGraph, node TGNode)
}

type TGGraphManager interface {
	// CreateNode creates Node within this Graph. There is a default Root Graph.
	CreateNode() (TGNode, TGError)
	// CreateNodeForNodeType creates Node of particular Type
	CreateNodeForNodeType(nodeType TGNodeType) (TGNode, TGError)
	// CreateEdge creates an Edge
	CreateEdge(fromNode TGNode, toNode TGNode, edgeType int) (TGEdge, TGError)
	// CreateEdgeWithDirection creates an Edge with direction
	CreateEdgeWithDirection(fromNode TGNode, toNode TGNode, directionType TGDirectionType) (TGEdge, TGError)
	// CreateGraph creates a SubGraph at the Root level.
	CreateGraph(name string) (TGGraph, TGError)
	// DeleteNode removes this node from the graph
	DeleteNode(filter TGFilter) (TGGraphManager, TGError)
	// DeleteNodes removes the nodes from this graph that match the filter
	DeleteNodes(filter TGFilter) (TGGraphManager, TGError)
	// CreateQuery creates a Reusable Query
	CreateQuery(filter TGFilter) TGQuery
	// QueryNodes gets Nodes based on the Filter condition with a set of Arguments
	QueryNodes(filter TGFilter, args ...interface{}) TGResultSet
	// Traverse follows the graph using the traversal descriptor
	Traverse(descriptor TGTraversalDescriptor, startingPoints []TGNode) TGResultSet
	// GetGraphMetadata gets the Graph Metadata
	GetGraphMetadata() TGGraphMetadata
}

func (systemType TGSystemType) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if systemType&SystemTypeInvalid == SystemTypeInvalid {
		buffer.WriteString("SystemTypeInvalid")
	} else if systemType&SystemTypeAttributeDescriptor == SystemTypeAttributeDescriptor {
		buffer.WriteString("SystemTypeAttributeDescriptor")
	} else if systemType&SystemTypeEntity == SystemTypeEntity {
		buffer.WriteString("SystemTypeEntity")
	} else if systemType&SystemTypeNode == SystemTypeNode {
		buffer.WriteString("SystemTypeNode")
	} else if systemType&SystemTypeEdge == SystemTypeEdge {
		buffer.WriteString("SystemTypeEdge")
	} else if systemType&SystemTypeIndex == SystemTypeIndex {
		buffer.WriteString("SystemTypeIndex")
	} else if systemType&SystemTypePrincipal == SystemTypePrincipal {
		buffer.WriteString("SystemTypePrincipal")
	} else if systemType&SystemTypeRole == SystemTypeRole {
		buffer.WriteString("SystemTypeRole")
	} else if systemType&SystemTypeSequence == SystemTypeSequence {
		buffer.WriteString("SystemTypeSequence")
	} else if systemType&SystemTypeMaxSysObject == SystemTypeMaxSysObject {
		buffer.WriteString("SystemTypeMaxSysObject")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}


// ======= Various Direction Types for EDGE (Entity) type =======

func (directionType TGDirectionType) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if directionType&DirectionTypeUnDirected == DirectionTypeUnDirected {
		buffer.WriteString("UnDirected")
	} else if directionType&DirectionTypeDirected == DirectionTypeDirected {
		buffer.WriteString("Directed")
	} else if directionType&DirectionTypeBiDirectional == DirectionTypeBiDirectional {
		buffer.WriteString("BiDirectional")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}

// ======= Various Directions associated with EDGE (Entity) type =======

func (direction TGDirection) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if direction&DirectionInbound == DirectionInbound {
		buffer.WriteString("Inbound")
	}
	if direction&DirectionOutbound == DirectionOutbound {
		buffer.WriteString("Outbound")
	}
	if direction&DirectionAny == DirectionAny {
		buffer.WriteString("Any")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}

