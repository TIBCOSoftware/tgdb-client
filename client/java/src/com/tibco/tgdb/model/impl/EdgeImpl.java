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
 * File name : EdgeImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: EdgeImpl.java 1633 2017-08-26 04:44:08Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import java.io.IOException;
import java.util.Map;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

public class EdgeImpl extends AbstractEntity implements TGEdge {

	private TGNode fromNode;
	private TGNode toNode;
	private DirectionType directionType;
	
    EdgeImpl(TGGraphMetadata gmd) {
        super(gmd);
    }

    EdgeImpl(TGGraphMetadata gmd, TGNode fromNode, TGNode toNode, DirectionType directionType) {
        super(gmd);
        this.fromNode = fromNode;
        this.toNode = toNode;
        this.directionType = directionType;
    }

    EdgeImpl(TGGraphMetadata gmd, TGNode fromNode, TGNode toNode, TGEdgeType edgeType) {
        super(gmd);
        this.fromNode = fromNode;
        this.toNode = toNode;
        this.entityType = edgeType;
        this.directionType = edgeType.getDirectionType();
    }

    @Override
    public TGNode[] getVertices() {
        return new TGNode[] {fromNode, toNode};
    }

    @Override
    public DirectionType getDirectionType() {
    	if (entityType != null) {
    		return ((TGEdgeType) entityType).getDirectionType();
    	} else {
    		return directionType;
        }
    }

    @Override
    public TGEntityKind getEntityKind() {
        return TGEntityKind.Edge;
    }

    @Override
    public void writeExternal(TGOutputStream os) throws TGException, IOException {
    	int startPos = os.getPosition();
    	os.writeInt(0);
    	//write attributes from the based class
        super.writeExternal(os);
        //write the edges ids
        if (isNew == true) {
        	os.writeByte(directionType.ordinal());
        	os.writeLong(((AbstractEntity) fromNode).getVirtualId());
        	os.writeLong(((AbstractEntity) toNode).getVirtualId());
        } else {
        	//FIXME:  Not sending the direction, is it ok?
            os.writeByte(directionType.ordinal()); //The Server expects it - so better send it.
        	os.writeLong(((AbstractEntity) fromNode).getVirtualId());
        	os.writeLong(((AbstractEntity) toNode).getVirtualId());
        }
        int currPos = os.getPosition();
        int length = currPos - startPos;
        os.writeIntAt(startPos, length);
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
    	//FIXME: Need to validate length
    	int buflen = is.readInt();
        super.readExternal(is);
        byte dir = is.readByte();
        if (dir == 0) {
        	this.directionType = DirectionType.UnDirected;
        } else if (dir == 1) {
        	this.directionType = DirectionType.Directed;
        } else {
        	this.directionType = DirectionType.BiDirectional;
        }
        NodeImpl fromNode = null;
        NodeImpl toNode = null;
        long id = is.readLong();
        Map refMap = is.getReferenceMap();
        if (refMap != null) {
        	fromNode = (NodeImpl) refMap.get(id);
        }
        if (fromNode == null) {
        	fromNode = new NodeImpl(graphMetadata);
        	fromNode.setEntityId(id);
        	fromNode.setInitialized(false);
        	if (refMap != null) {
        		refMap.put(id, fromNode);
        	}
        }
        this.fromNode = fromNode;

        id = is.readLong();
        if (refMap != null) {
        	toNode = (NodeImpl) refMap.get(id);
        }
        if (toNode == null) {
        	toNode = new NodeImpl(graphMetadata);
        	toNode.setEntityId(id);
        	toNode.setInitialized(false);
        	if (refMap != null) {
        		refMap.put(id, toNode);
        	}
        }
        this.toNode = toNode;
    	isInitialized = true;
    }
}
