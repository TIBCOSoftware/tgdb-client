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
 * File name : EntityTypeImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: EntityTypeImpl.java 622 2016-03-19 20:51:12Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.util.*;

public abstract class EntityTypeImpl implements TGEntityType {

    Map<String, TGAttributeDescriptor> attributes = new LinkedHashMap<String, TGAttributeDescriptor>();

    private int id; //issued only for creation and not valid later
    private String name;
    private TGEntityType parent;

    static TGLogger gLogger = TGLogManager.getInstance().getLogger();

    /*
    protected EntityTypeImpl(TGGraphMetadata gmd) {
    }
    */

    public EntityTypeImpl() {
    	
    }

    @Override
    public Collection<TGAttributeDescriptor> getAttributeDescriptors() {
        return attributes.values();
    }

    @Override
    public TGAttributeDescriptor getAttributeDescriptor(String attrName) {
        return attributes.get(attrName);
    }

    public int getId() {
        return id;
    }

    void updateMetadata(TGGraphMetadata gmd) {
        for (String attrName : attributes.keySet()) {
            TGAttributeDescriptor desc = gmd.getAttributeDescriptor(attrName);
            if (desc == null) {
    	        gLogger.log(TGLevel.Warning, "Cannot find '%s' attribute descriptor", attrName);
                continue;
            }
            attributes.put(attrName, desc);
        }
    }

    @Override
    public String getName() {
    	return name;
    }

    @Override
    public TGEntityType derivedFrom() {
        return parent;
    }

    @Override
    public void writeExternal(TGOutputStream os) throws TGException, IOException {
    	gLogger.log(TGLevel.Warning, "writeExternal for entity desc is not implemented");
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
    	//FIXME: Do we save the desc value??
    	int typeValue = is.readByte();
    	TGSystemType type = TGSystemType.fromValue(typeValue);
    	if (type == TGSystemType.InvalidType) {
    		gLogger.log(TGLevel.Warning, "Entity desc input stream has invalid desc value : %d", typeValue);
    		//FIXME: Need to throw Exception
    	}

    	id = is.readInt();
    	name = is.readUTF();

    	int pageSize = is.readInt(); // pagesize

    	int attrCount = is.readShort();
    	for (int i=0; i<attrCount; i++) {
    		String name = is.readUTF();
    		//FIXME: The stream only contains name of the descriptor.
    		//Need to lookup the attribute descriptor from GraphMetaData object
    		TGAttributeDescriptor attrDesc = new AttributeDescriptorImpl(name, TGAttributeType.String);
    		attributes.put(name, attrDesc);
    	}
    }
}

