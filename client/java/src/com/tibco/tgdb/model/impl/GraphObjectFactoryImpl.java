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
 * File name : GraphModelFactoryImpl.${EXT}
 * Created on: 1/28/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: GraphObjectFactoryImpl.java 2191 2018-03-31 00:47:05Z vchung $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.*;


public class GraphObjectFactoryImpl implements TGGraphObjectFactory {

    protected TGGraphMetadata graphMetadata;
    protected TGChangeListener lsnr;

    public GraphObjectFactoryImpl(TGGraphMetadata graphMetadata, TGChangeListener lsnr) {
    	//FIXME: Graph meta data cannot be passed in.
    	//There will be one meta data object per object factory
    	//And one object factory per connection even though connections can share 
    	//the same channel. How do we handle update notifications from other clients?
        //this.graphMetadata = graphMetadata;
        if (graphMetadata == null) {
            this.graphMetadata = new GraphMetadataImpl();
        } else {
            this.graphMetadata = graphMetadata;
        }
        this.lsnr = lsnr;
    }

    public TGGraphMetadata getGraphMetaData() { return graphMetadata;}

    @Override
    public TGNode createNode()  {
        TGNode node = new NodeImpl(graphMetadata);

        return node;
    }

    /**
     * Create Entity based on kind. This is used for deserialization purpose only. Does not notify  the listener.
     * @param kind
     * @return
     */
    public TGEntity createEntity(TGEntity.TGEntityKind kind) {
        switch (kind) {
            case Node:
                return new NodeImpl(graphMetadata);
            case Edge:
                return new EdgeImpl(graphMetadata);
            case Graph:
                return new GraphImpl(graphMetadata);
        }
        return null;
    }

    @Override
    public TGNode createNode(TGNodeType nodeType) {
        TGNode node = new NodeImpl(graphMetadata, nodeType);

        return node;
    }

    @Override
    public TGEdge createEdge(TGNode fromNode, TGNode toNode, TGEdge.DirectionType directionType) {
        TGEdge edge = new EdgeImpl(graphMetadata, fromNode, toNode, directionType);
        ((NodeImpl)fromNode).addEdge(edge);
        ((NodeImpl)toNode).addEdge(edge);

        return edge;
    }

    @Override
    public TGEdge createEdge(TGNode fromNode, TGNode toNode, TGEdgeType edgeType)  {
        TGEdge edge =  new EdgeImpl(graphMetadata, fromNode, toNode, edgeType);
        ((NodeImpl)fromNode).addEdge(edge);
        ((NodeImpl)toNode).addEdge(edge);

        return edge;
    }


    public TGGraph createGraph(String name)  {
        TGGraph graph = new GraphImpl(graphMetadata, name);

        return graph;
    }

    public TGEntityId createEntityId(byte[] buf) throws TGException {
        return new ByteArrayEntityId(buf);
    }

    @Override
    public TGKey createCompositeKey(String nodeTypeName) throws TGException {
    	TGNodeType nodeType = graphMetadata.getNodeType(nodeTypeName);
    	if (nodeType == null) {
            throw TGException.buildException("Node type not found", null, null);
    	}
        return new CompositeKeyImpl(graphMetadata, nodeTypeName);
    }
}
