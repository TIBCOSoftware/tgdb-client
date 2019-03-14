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
 * File name : TGGraphModelFactory.java
 * Created on: 1/28/15
 * Created by: suresh 

 * SVN Id: $Id: TGGraphObjectFactory.java 2348 2018-06-22 16:34:26Z ssubrama $
 */


package com.tibco.tgdb.model;

import com.tibco.tgdb.exception.TGException;

/**
 * A model Factory for creating Graph Objects
 */
public interface TGGraphObjectFactory {

    /**
     * Create a Node
     * @return a node created
     */
     TGNode createNode() ;

    /**
     * Create Node within this Graph. There is a default Root Graph
     * @param nodeType - The typenode that this node is instanceof
     * @return the node created from this nodeType
     */
     TGNode createNode(TGNodeType nodeType);


    /**
     * Create an Edge
     * @param fromNode the starting node
     * @param toNode to end node
     * @param directionType - BiDirection/UniDirection or None
     * @return a TGEdge
     */
     TGEdge createEdge(TGNode fromNode, TGNode toNode, TGEdge.DirectionType directionType);


    /**
     * Create an Edge
     * @param fromNode the starting node
     * @param toNode to end node
     * @param edgeType the edgeType associated to this edge
     * @return a TGEdge
     */
     TGEdge createEdge(TGNode fromNode, TGNode toNode, TGEdgeType edgeType);


    /**
     * Create a CompositeKey for a NodeType. The composite key can also be a single key
     * @param nodeTypeName - An optional Nodetype name
     * @return returns a TGKey that can be used to query the Database.
     * @throws TGException if a NodeType is not found.
     */
     TGKey createCompositeKey(String nodeTypeName) throws TGException;


}
