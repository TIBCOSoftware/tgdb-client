/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * File name : TGGraph.java
 * Created on: 1/22/15
 * Created by: suresh 
 *
 * SVN Id: $Id: TGGraph.java 623 2016-03-19 21:41:13Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.query.TGFilter;

import java.util.List;

/**
 * @todo - This is a placeholder interface for future use.
 */
public interface TGGraph extends TGNode {


    /**
     * Add a node to this graph, but not the edges
     * @param node
     * @return TGNode
     * @throws TGException
     */
    TGNode addNode(TGNode node) throws TGException;

    /**
     * A Subgraph covers a list of edges. This implies that the two nodes that edge connects
     * are also included
     * @param edges
     */
    void addEdges(List<TGEdge> edges);

    /**
     * Return a unique node which matches the unique constraint
     * @param filter
     * @return TGNode
     * @throws TGException
     */
    TGNode getNode(TGFilter filter) throws TGException;

    /**
     * List all the nodes that match the filter and recurse All subgraphs
     * @param filter
     * @return TGNode
     * @throws TGException
     */
    TGNode listNodes(TGFilter filter, boolean recurseAllSubgraphs) throws TGException;


    /**
     * Create a subgraph within this graph
     * @param name graph name
     * @return TGGraph
     */
    TGGraph createGraph(String name);

    /**
     * Remove the nodes from this graph
     * @param filter
     * @return nos of nodes removed
     */
    int removeNodes(TGFilter filter);

    /**
     * Remove this node from the graph.
     * The Edge is also removed.
     * @param node to remove
     */
    void removeNode(TGNode node);

    /**
     * Remove the Graph
     * @param name graph name
     */
    void removeGraph(String name);

}
