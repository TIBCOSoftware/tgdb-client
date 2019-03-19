package types

import "bytes"

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
 * File name: TGEntity.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= Various Entity Kind =======
type TGEntityKind int

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
