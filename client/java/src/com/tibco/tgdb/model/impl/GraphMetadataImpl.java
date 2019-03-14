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
 * File name : GraphMetadataImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: GraphMetadataImpl.java 2348 2018-06-22 16:34:26Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

//TODO: There is no deletion of meta data ??
public class GraphMetadataImpl implements TGGraphMetadata {
	
	private boolean isInitialized = false;
	private HashMap<String, TGAttributeDescriptor> descriptorMap;
	private HashMap<String, TGNodeType> nodeTypeMap;
	private HashMap<String, TGEdgeType> edgeTypeMap;
	private HashMap<Integer, TGAttributeDescriptor> descriptorMapById;
	private HashMap<Integer, TGNodeType> nodeTypeMapById;
	private HashMap<Integer, TGEdgeType> edgeTypeMapById;
	private GraphObjectFactoryImpl gof;
	
	public GraphMetadataImpl(GraphObjectFactoryImpl gof) {
		descriptorMap = new HashMap<String, TGAttributeDescriptor>();
		nodeTypeMap = new HashMap<String, TGNodeType>();
		edgeTypeMap = new HashMap<String, TGEdgeType>();
		descriptorMapById = new HashMap<Integer, TGAttributeDescriptor>();
		nodeTypeMapById = new HashMap<Integer, TGNodeType>();
		edgeTypeMapById = new HashMap<Integer, TGEdgeType>();
		this.gof = gof;
	}

    @Override
    public Set<TGNodeType> getNodeTypes() {
    	return nodeTypeMap.values().stream().collect(Collectors.toSet());
    }

    @Override
    public TGNodeType getNodeType(String typeName) {
    	return nodeTypeMap.get(typeName);
    }

    @Override
    public Set<TGEdgeType> getEdgeTypes() {
    	return edgeTypeMap.values().stream().collect(Collectors.toSet());
    }

    @Override
    public TGEdgeType getEdgeType(String typeName) {
        return edgeTypeMap.get(typeName);
    }

    @Override
    public Set<TGAttributeDescriptor> getAttributeDescriptors() {
    	return descriptorMap.values().stream().collect(Collectors.toSet());
    }


    public Set<TGAttributeDescriptor> getNewAttributeDescriptors() {
    	return descriptorMap.values().stream().filter(a -> a.getAttributeId() < 0).collect(Collectors.toSet());
    }

    @Override
    public TGAttributeDescriptor getAttributeDescriptor(String attributeName) {
        return descriptorMap.get(attributeName);
    }

    @Override
    public TGAttributeDescriptor createAttributeDescriptor(String attrName, TGAttributeType attrType, boolean isArray) throws TGException {
    	TGAttributeDescriptor desc = new AttributeDescriptorImpl(attrName, attrType, isArray);
    	descriptorMap.put(attrName, desc);
        return desc;
    }


    public TGAttributeDescriptor createAttributeDescriptor(String attrName, Class klazz) throws TGException {
    	TGAttributeDescriptor desc = new AttributeDescriptorImpl(attrName, TGAttributeType.fromClass(klazz));
    	descriptorMap.put(attrName, desc);
    	return desc;
    }


    public TGNodeType createNodeType(String typeName, TGNodeType parentNodeType) {
    	TGNodeType nodeType = new NodeTypeImpl(typeName, parentNodeType);
    	return nodeType;
    }


    public TGEdgeType createEdgeType(String typeName, TGEdgeType parentEdgeType) {
    	TGEdgeType edgeType = new EdgeTypeImpl(typeName, parentEdgeType.getDirectionType(), parentEdgeType);
        return edgeType;
    }

    public TGKey createCompositeKey(String typeName) throws TGException {

        TGKey key = new CompositeKeyImpl(this, typeName);
        return key;
    }
    
    public void updateMetadata(List<TGAttributeDescriptor> attrDescList, 
    		List<TGNodeType> nodeTypeList, List<TGEdgeType> edgeTypeList) throws TGException {
    	if (attrDescList != null) {
    		for (TGAttributeDescriptor desc : attrDescList) {
    			descriptorMapById.put(((AttributeDescriptorImpl) desc).getAttributeId(), desc);
    			descriptorMap.put(desc.getName(), desc);
    		}
    	}
    	if (nodeTypeList != null) {
    		for (TGNodeType nt : nodeTypeList) {
                ((NodeTypeImpl) nt).updateMetadata(this);
    			nodeTypeMap.put(nt.getName(), nt);
    			nodeTypeMapById.put(((NodeTypeImpl)nt).getId(), nt);
    		}
    	}
    	if (edgeTypeList != null) {
    		for (TGEdgeType et : edgeTypeList) {
                ((EdgeTypeImpl) et).updateMetadata(this);
    			edgeTypeMap.put(et.getName(), et);
    			edgeTypeMapById.put(((EdgeTypeImpl)et).getId(), et);
    		}
    	}
    	isInitialized = true;
    }

    public TGAttributeDescriptor getAttributeDescriptor(int id) {
    	return descriptorMapById.get(id);
    }

    TGNodeType getNodeType(int id) {
    	return nodeTypeMapById.get(id);
    }

    TGEdgeType getEdgeType(int id) {
    	return edgeTypeMapById.get(id);
    }

    public boolean isInitialized() {
    	return isInitialized;
    }

    @Override
    public void writeExternal(TGOutputStream os) throws TGException, IOException {

    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {

    }
    
    public GraphObjectFactoryImpl getGraphObjectFactory() {
	    return this.gof;
    }

    public ConnectionImpl getConnection() { return this.gof.connection; }
}
