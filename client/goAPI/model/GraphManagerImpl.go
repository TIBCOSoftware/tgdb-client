package model

import (
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
 * File name: TGGraphManager.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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

// GetName gets Graph Manager's name
func (obj *GraphManager) GetName() string {
	return obj.name
}

///////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGGraphManager //
///////////////////////////////////////////////////////////

// CreateNode creates Node within this Graph. There is a default Root Graph.
func (obj *GraphManager) CreateNode() (types.TGNode, types.TGError) {
	return nil, nil
}

// CreateNodeForNodeType creates Node of particular Type
func (obj *GraphManager) CreateNodeForNodeType(nodeType types.TGNodeType) (types.TGNode, types.TGError) {
	return nil, nil
}

// CreateEdge creates an Edge
func (obj *GraphManager) CreateEdge(fromNode types.TGNode, toNode types.TGNode, edgeType int) (types.TGEdge, types.TGError) {
	return nil, nil
}

// CreateEdgeWithDirection creates an Edge with direction
func (obj *GraphManager) CreateEdgeWithDirection(fromNode types.TGNode, toNode types.TGNode, directionType types.TGDirectionType) (types.TGEdge, types.TGError) {
	return nil, nil
}

// CreateGraph creates a SubGraph at the Root level.
func (obj *GraphManager) CreateGraph(name string) (types.TGGraph, types.TGError) {
	return nil, nil
}

// DeleteNode removes this node from the graph
func (obj *GraphManager) DeleteNode(filter types.TGFilter) (types.TGGraphManager, types.TGError) {
	return nil, nil
}

// DeleteNodes removes the nodes from this graph that match the filter
func (obj *GraphManager) DeleteNodes(filter types.TGFilter) (types.TGGraphManager, types.TGError) {
	return nil, nil
}

// CreateQuery creates a Reusable Query
func (obj *GraphManager) CreateQuery(filter types.TGFilter) types.TGQuery {
	return nil
}

// QueryNodes gets Nodes based on the Filter condition with a set of Arguments
func (obj *GraphManager) QueryNodes(filter types.TGFilter, args ...interface{}) types.TGResultSet {
	return nil
}

// Traverse follows the graph using the traversal descriptor
func (obj *GraphManager) Traverse(descriptor types.TGTraversalDescriptor, startingPoints []types.TGNode) types.TGResultSet {
	return nil
}

// GetGraphMetadata gets the Graph Metadata
func (obj *GraphManager) GetGraphMetadata() types.TGGraphMetadata {
	return nil
}
