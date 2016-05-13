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
 * File name : GraphImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: GraphImpl.java 623 2016-03-19 21:41:13Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.query.TGFilter;

import java.util.List;

//TODO:  Need to add implementation
public class GraphImpl extends NodeImpl implements TGGraph {

    String name;

    GraphImpl(TGGraphMetadata gmd) {
        super(gmd);
    }

    GraphImpl(TGGraphMetadata gmd, String name) {

        super(gmd);
        this.name = name;
    }

    @Override
    public TGNode addNode(TGNode node) throws TGException {
        return null;
    }

    @Override
    public void addEdges(List<TGEdge> edges) {

    }

    @Override
    public TGNode getNode(TGFilter filter) throws TGException {
        return null;
    }

    @Override
    public TGNode listNodes(TGFilter filter, boolean recurseAllSubgraphs) throws TGException {
        return null;
    }

    @Override
    public TGGraph createGraph(String name) {
        return null;
    }

    @Override
    public int removeNodes(TGFilter filter) {
        return 0;
    }

    @Override
    public void removeNode(TGNode node) {

    }

    @Override
    public void removeGraph(String name) {

    }
}
