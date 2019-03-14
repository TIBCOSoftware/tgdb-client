package com.tibco.tgdb.query.impl;
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
 *
 * File name : TGQueryOption.java
 * SVN Id: $Id$
 */

import java.io.IOException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.impl.GraphObjectFactoryImpl;
import com.tibco.tgdb.model.impl.attribute.AbstractAttribute;
import com.tibco.tgdb.pdu.TGInputStream;

public class GremlinResult {

   	private enum ElementType {
       	Invalid(0),
       	List(1),
       	Attr(2),
       	AttrValue(3),
       	Entity(4),
       	Map(5);

    	private int id;
    
    	ElementType (int id) {
    		this.id = (short) id;
    	}

    	public static ElementType fromValue(int id) {
    		for (ElementType et : ElementType.values()) {
    			if (id == et.id) return et;
    		}
    		return Invalid;
    	}

    	public int getId() { 
    		return id; 
    	}
   	}

    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    //Only handle entity right now
    static public void fillCollection(TGInputStream entityStream, GraphObjectFactoryImpl gof, Collection col) {
    	try {
    		//These types will have an enum to represent them
    		ElementType retType = ElementType.fromValue(entityStream.readByte());
    		if (retType == ElementType.List) {
    			constructList(entityStream, gof, col);
    		} else {
    			gLogger.log(TGLogger.TGLevel.Error, "Invalid gremlin response collection type : %d", retType.getId());
    		}
    	} catch (IOException ioe) {
            gLogger.log(TGLogger.TGLevel.Error, "Failed to read gremlin query response stream");
    	} catch (TGException tge) {
            gLogger.log(TGLogger.TGLevel.Error, "Failed to read gremlin query response stream");
    	}
    }
    
    static public void constructList(TGInputStream entityStream, GraphObjectFactoryImpl gof, Collection col) throws IOException, TGException {
   		int size = entityStream.readInt();
   		int et = entityStream.readByte();
   		ElementType elemType = ElementType.fromValue(et);
   		TGNode dummyNode = null;
   		if (elemType == ElementType.Attr || elemType == ElementType.AttrValue) {
   			dummyNode = gof.createNode();
   		}
   		for (int i=0; i<size; i++) {
   			if (elemType == ElementType.Entity) {
   				TGEntity.TGEntityKind kind = TGEntity.TGEntityKind.fromValue(entityStream.readByte());
   				if (kind == TGEntity.TGEntityKind.InvalidKind) {
   					gLogger.log(TGLogger.TGLevel.Error, "Invalid entity kind from gremlin response stream");
   					break;
   				}
   				if (kind == TGEntity.TGEntityKind.Node) {
   					TGNode node = gof.createNode();
   					node.readExternal(entityStream);
   					col.add(node);
   				}
   			} else if (elemType == ElementType.Attr) {
   				TGAttribute attr = AbstractAttribute.readExternal(dummyNode, entityStream);
   				col.add(attr);
   			} else if (elemType == ElementType.AttrValue) {
   				TGAttribute attr = AbstractAttribute.readExternal(dummyNode, entityStream);
   				col.add(attr.getValue());
   			} else if (elemType == ElementType.List) {
   				Collection colElem = new ArrayList();
   				constructList(entityStream, gof, colElem);
   				col.add(colElem);
   			} else if (elemType == ElementType.Map) {
   				Map<String, Object> mapElem = new HashMap<String, Object>();
   				constructMap(entityStream, gof, mapElem);
   				col.add(mapElem);
   			} else {
				gLogger.log(TGLogger.TGLevel.Error, "Invalid element type %d from gremlin response stream", et);
   			}
   		}
    }

    static public void constructMap(TGInputStream entityStream, GraphObjectFactoryImpl gof, Map<String, Object> map) throws IOException, TGException {
   		int size = entityStream.readInt();
   		for (int i=0; i<size; i++) {
   			String key = entityStream.readUTF();
			int et = entityStream.readByte();
			ElementType elemType = ElementType.fromValue(et);
			TGNode dummyNode = null;
			if (elemType == ElementType.Attr || elemType == ElementType.AttrValue) {
				dummyNode = gof.createNode();
			}
   			if (elemType == ElementType.Entity) {
   				TGEntity.TGEntityKind kind = TGEntity.TGEntityKind.fromValue(entityStream.readByte());
   				if (kind == TGEntity.TGEntityKind.InvalidKind) {
   					gLogger.log(TGLogger.TGLevel.Error, "Invalid entity kind from gremlin response stream");
   					break;
   				}
   				if (kind == TGEntity.TGEntityKind.Node) {
   					TGNode node = gof.createNode();
   					node.readExternal(entityStream);
   					map.put(key, node);
   				} else if (kind == TGEntity.TGEntityKind.Node) {
   					//It has no from/or node right now
   					TGEntity edge = gof.createEntity(TGEntity.TGEntityKind.Edge);
   					edge.readExternal(entityStream);
   					map.put(key, edge);
   				}
   			} else if (elemType == ElementType.Attr) {
   				TGAttribute attr = AbstractAttribute.readExternal(dummyNode, entityStream);
   				map.put(key, attr);
   			} else if (elemType == ElementType.AttrValue) {
   				TGAttribute attr = AbstractAttribute.readExternal(dummyNode, entityStream);
   				map.put(key, attr.getValue());
   			} else if (elemType == ElementType.List) {
   				Collection colElem = new ArrayList();
   				constructList(entityStream, gof, colElem);
   				map.put(key, colElem);
   			} else if (elemType == ElementType.Map) {
   				Map<String, Object> mapElem = new HashMap<String, Object>();
   				constructMap(entityStream, gof, mapElem);
   				map.put(key, mapElem);
   			} else {
				gLogger.log(TGLogger.TGLevel.Error, "Invalid element type %d from gremlin response stream", et);
   			}
   		}
    }
}
