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
 * <p/>
 * File name : GraphManager.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: GraphManagerImpl.java 622 2016-03-19 20:51:12Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.query.TGFilter;
import com.tibco.tgdb.query.TGQuery;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.query.TGTraversalDescriptor;

import java.util.Set;

public class GraphManagerImpl implements TGGraphManager {

    @Override
    public TGNode createNode() throws TGException {
        return null;
    }

    @Override
    public TGNode createNode(TGNodeType nodeType) throws TGException {
        return null;
    }

    @Override
    public TGEdge createEdge(TGNode fromNode, TGNode toNode, TGEdge.DirectionType directionType) throws TGException {
        return null;
    }

    @Override
    public TGEdge createEdge(TGNode fromNode, TGNode toNode, TGEdgeType edgeType) throws TGException {
        return null;
    }

    @Override
    public TGGraph createGraph(String name) throws TGException {
        return null;
    }

    @Override
    public TGResultSet queryNodes(TGFilter filter, Object... args) {
        return null;
    }

    @Override
    public TGQuery createQuery(TGFilter filter) {
        return null;
    }

    @Override
    public TGResultSet traverse(TGTraversalDescriptor descriptor, Set<TGNode> startingPoints) {
        return null;
    }

    @Override
    public void deleteNode(TGFilter filter) throws Exception {

    }

    @Override
    public void deleteNodes(TGFilter filter) throws Exception {

    }

    @Override
    public TGGraphMetadata getGraphMetadata() {
        return null;
    }
}
