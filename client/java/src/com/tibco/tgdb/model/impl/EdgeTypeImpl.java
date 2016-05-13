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
 * SVN Id: $Id: EdgeTypeImpl.java 723 2016-04-16 19:21:18Z vchung $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGEntityType;
import com.tibco.tgdb.model.TGGraphMetadata;

public class EdgeTypeImpl extends EntityTypeImpl implements TGEdgeType {
    TGEdge.DirectionType directionType;

    //FIXME: Can directionType different from parent direction type?
    public EdgeTypeImpl(String name, TGEdge.DirectionType directionType, TGEdgeType parent) {
    	super();
    	this.directionType = directionType;
    }

    @Override
    public TGSystemType getSystemType() {
    	return TGSystemType.EdgeType;
    }

    @Override
    public TGEdge.DirectionType getDirectionType() {
        return directionType;
    }
}
