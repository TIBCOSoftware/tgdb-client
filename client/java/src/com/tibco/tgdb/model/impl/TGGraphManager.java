package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.query.TGFilter;
import com.tibco.tgdb.query.TGQuery;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.query.TGTraversalDescriptor;

import java.util.Set;

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

 * File name : TGGraphManager.java
 * Created on: 1/21/15
 * Created by: suresh

 * SVN Id: $Id: TGGraphManager.java 623 2016-03-19 21:41:13Z ssubrama $
 */

/**
 * @todo For Future Development. Currently a Marker. For Object creation use GraphObjectFactory
 * @see com.tibco.tgdb.model.TGGraphObjectFactory
 */
public interface TGGraphManager  {



    /**
     * Create Node within this Graph. There is a default Root Graph
     * @return
     * @throws TGException
     */
    TGNode createNode() throws TGException;

    /**
     * Create Node of particular Type
     * @param nodeType
     * @return
     * @throws TGException
     */
    TGNode createNode(TGNodeType nodeType) throws TGException;

    /**
     * Create an Edge
     * @param fromNode
     * @param toNode
     * @return
     * @throws TGException
     */
    TGEdge createEdge(TGNode fromNode, TGNode toNode, TGEdge.DirectionType directionType) throws TGException;


    /**
     * Create an Edge
     * @param fromNode
     * @param toNode
     * @param edgeType
     * @return
     * @throws TGException
     */
    TGEdge createEdge(TGNode fromNode, TGNode toNode, TGEdgeType edgeType) throws TGException;


    /**
     * Create a SubGraph at the Root level.
     * @param name
     * @return
     * @throws TGException
     */
    TGGraph createGraph(String name) throws TGException;



    /**
     * Query Nodes based on the Filter condition with a set of Arguments
     * @param filter
     * @param args
     * @return
     */
    TGResultSet queryNodes(TGFilter filter, Object... args);

    /**
     * Create a Resuable Query
     * @param filter
     * @return
     */
    TGQuery createQuery(TGFilter filter);

    /**
     * Traverse the graph using the traversal descriptor
     * @param descriptor
     * @param startingPoints
     * @return
     */
    TGResultSet traverse(TGTraversalDescriptor descriptor, Set<TGNode> startingPoints);

    /**
     * RemoveNode using the filter constraint.
     * @param filter
     * @throws Exception
     */
    void deleteNode(TGFilter filter) throws Exception;


    void deleteNodes(TGFilter filter) throws Exception;

    /**
     * Get the Graph Metadata
     * @return
     */
    TGGraphMetadata getGraphMetadata();


}
