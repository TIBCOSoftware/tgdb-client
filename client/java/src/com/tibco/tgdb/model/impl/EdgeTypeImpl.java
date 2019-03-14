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
 * File name : EdgeTypeImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: EdgeTypeImpl.java 2344 2018-06-11 23:21:45Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import java.io.IOException;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.pdu.TGInputStream;

public class EdgeTypeImpl extends EntityTypeImpl implements TGEdgeType {
    private TGEdge.DirectionType directionType;
    private int fromTypeId;
    private int toTypeId;
    private TGNodeType fromNodeType;
    private TGNodeType toNodeType;
    private long numEntries;

    //FIXME: Can directionType different from parent direction desc?
    public EdgeTypeImpl(String name, TGEdge.DirectionType directionType, TGEdgeType parent) {
    	super();
    	this.directionType = directionType;
    }

    int getFromId() {
        return fromTypeId;
    }

    int getToId() {
        return toTypeId;
    }

    void updateMetadata(TGGraphMetadata gmd) {
        fromNodeType = ((GraphMetadataImpl) gmd).getNodeType(fromTypeId);
        toNodeType = ((GraphMetadataImpl) gmd).getNodeType(toTypeId);
    }

    @Override
    public TGSystemType getSystemType() {
    	return TGSystemType.EdgeType;
    }

    @Override
    public TGEdge.DirectionType getDirectionType() {
        return directionType;
    }

    @Override
    public TGNodeType getFromNodeType() {
        return fromNodeType;
    }

    @Override
    public TGNodeType getToNodeType() {
        return toNodeType;
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
    	super.readExternal(is);

        fromTypeId = is.readInt();
        toTypeId = is.readInt();
        byte dir = is.readByte();
        if (dir == 0) {
        	this.directionType = TGEdge.DirectionType.UnDirected;
        } else if (dir == 1) {
        	this.directionType = TGEdge.DirectionType.Directed;
        } else {
        	this.directionType = TGEdge.DirectionType.BiDirectional;
        }
        numEntries = is.readLong();
    }
}
