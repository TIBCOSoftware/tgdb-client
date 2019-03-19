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
 * File name: TGGraphManager.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
