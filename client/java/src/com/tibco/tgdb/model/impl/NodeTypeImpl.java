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
 * File name : NodeTypeImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: NodeTypeImpl.java 978 2016-09-18 21:43:27Z vchung $
 */


package com.tibco.tgdb.model.impl;

import java.io.IOException;
import java.util.ArrayList;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEntityType;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.pdu.TGInputStream;

public class NodeTypeImpl extends EntityTypeImpl implements TGNodeType {

    ArrayList<String> pKeys = new ArrayList<String>();

    public NodeTypeImpl(String name, TGNodeType parent) {
    }

    @Override
    public TGSystemType getSystemType() {
    	return TGSystemType.NodeType;
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
    	super.readExternal(is);

    	int attrCount = is.readShort();
    	for (int i=0; i<attrCount; i++) {
    		pKeys.add(is.readUTF());
    	}

    	int idxCount = is.readShort();
    	for (int i=0; i<idxCount; i++) {
            //FIXME: Get meta data needs to return index definitions
    		is.readInt();
    	}
    }
}
