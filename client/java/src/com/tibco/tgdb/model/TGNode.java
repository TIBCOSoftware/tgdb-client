package com.tibco.tgdb.model;


import java.util.Collection;

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
 *  File name :TGNode.java
 *  Created by: suresh
 *
 *		SVN Id: $Id: TGNode.java 2344 2018-06-11 23:21:45Z ssubrama $
 *
 */

public interface TGNode extends TGEntity {

    /**
     * Get the List of Edges
     * @return a collection of edges for this node
     */
    Collection<TGEdge> getEdges();

    /**
     *
     * @param directionType the direction desc associated
     * @return collection edges for the direction desc
     */
    Collection<TGEdge> getEdges(TGEdge.DirectionType directionType);

    /**
     *
     * @param edgeType the edge desc and it can be null
     * @param direction the edge direction relative to the node(inbound, outbound or any)
     * @return collection filtered edges
     */
    Collection<TGEdge> getEdges(TGEdgeType edgeType, TGEdge.Direction direction);

    /**
     * Add an Edge to another Node
     * @param node The destination node
     * @param edgeType Optional EdgeType
     * @param directionType the direction of the edge
     * @return the newly created Edge.
     */
    TGEdge addEdge(TGNode node, TGEdgeType edgeType, TGEdge.DirectionType directionType );

}
