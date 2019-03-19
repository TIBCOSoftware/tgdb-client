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
 * File name: TGGraphObjectFactory.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
