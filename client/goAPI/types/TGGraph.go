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
 * File name: TGGraph.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
