package com.tibco.tgdb.utils;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGEntity;


/**
 * Copyright 2017 TIBCO Software Inc. All rights reserved.
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

 * File name :EntityUtils
 * Created on: 10/5/17
 * Created by: chung

 * SVN Id: $Id: $
 */

public class EntityUtils {
    public static void printEntities(TGEntity ent, int maxDepth, int currDepth, String indent, boolean showAllPath, Map<Integer, TGEntity> traverseMap) {
        if (currDepth == maxDepth) {
        	return;
        }
    	if (ent == null) {
    		return;
    	}
        TGEntity.TGEntityKind kind = ent.getEntityKind();
        if (kind == TGEntity.TGEntityKind.Node || (kind == TGEntity.TGEntityKind.Edge && !showAllPath)) {
    	    if (traverseMap.get(ent.hashCode()) != null) {
                /*
                System.out.printf("%sKind : %s, hashcode : %d visited\n", indent, 
                    (ent.getEntityKind() == TGEntity.TGEntityKind.Node ? "Node" : "Edge"), ent.hashCode());
                return;
                */
                return;
            }
            traverseMap.put(ent.hashCode(), ent);
        }
       	System.out.printf("%sKind : %s, hashcode : %d\n", indent, (kind == TGEntity.TGEntityKind.Node ? "Node" : "Edge"), ent.hashCode());
        for (TGAttribute attrib : ent.getAttributes()) {
        	System.out.printf("%s Attr : %s\n", indent, attrib.getValue());
        }
        if (kind == TGEntity.TGEntityKind.Node) {
        	int edgeCount = ((TGNode) ent).getEdges().size();
        	System.out.printf("%s Has %d edges\n", indent, edgeCount);
        	String newIndent = new String(indent).concat(" ");
        	for (TGEdge edge : ((TGNode) ent).getEdges()) {
        		printEntities(edge, maxDepth, currDepth, newIndent, showAllPath, traverseMap);
        	}
        } else if (kind == TGEntity.TGEntityKind.Edge) {
        	TGNode[] nodes = ((TGEdge) ent).getVertices();
        	if (nodes.length > 0) {
        		System.out.printf("%s Has end nodes\n", indent);
        	}
        	String newIndent = new String(indent).concat("  ");
        	++currDepth;
        	for (int j=0; j<nodes.length; j++) {
        		printEntities(nodes[j], maxDepth, currDepth, newIndent, showAllPath, traverseMap);
        	}
        }
        if (showAllPath && kind == TGEntity.TGEntityKind.Node) {
            traverseMap.remove(ent.hashCode());
        }
    }

    public static void printEntitiesBreadth(TGNode node, int maxDepth) {
        String indent = "";
        Map<Integer, TGEntity> traverseMap = new HashMap<Integer, TGEntity>();
        List<TGNode> nodeList = new ArrayList<TGNode>();
        nodeList.add(node);
        int level = 1;

        while (nodeList.size() > 0) {
            int listSize = nodeList.size();
            while (listSize > 0) {
                listSize--;
                node = nodeList.get(0);
                nodeList.remove(0);
                if (traverseMap.get(node.hashCode()) != null) {
                    System.out.printf("%s%d:Node(%d) visited\n", indent, level, node.hashCode());
                    continue;
                } 
                System.out.printf("%s%d:Node(%d)\n", indent, level, node.hashCode());
                for (TGAttribute attrib : node.getAttributes()) {
                    System.out.printf("%s %d:Attr : %s\n", indent, level, attrib.getValue());
                }
                traverseMap.put(node.hashCode(), node);
        	    System.out.printf("%s %d:has %d edges\n", indent, level, node.getEdges().size());
                for (TGEdge edge : node.getEdges()) {
                    TGNode[] nodes = edge.getVertices();
                    if (nodes.length > 0) {
                        if (nodes[0] == null) {
                           	//System.out.printf("%s   %d:passed last edge\n", indent, level);
                            //break;
                        	continue;
                        }
                        for (int j=0; j<nodes.length; j++) {
                            if (nodes[j].hashCode() == node.hashCode()) {
                                continue;
                            } else {
                            	Iterator<TGAttribute> itr = edge.getAttributes().iterator();
                            	if (itr.hasNext()) {
                            		TGAttribute attrib = itr.next();
                            		System.out.printf("%s   %d:edge(%d)(%s) with node(%d)\n", indent, level, edge.hashCode(), attrib.getValue(), nodes[j].hashCode());
                            	} else {
                            		System.out.printf("%s   %d:edge(%d) with node(%d)\n", indent, level, edge.hashCode(), nodes[j].hashCode());
                            	}
                                nodeList.add(nodes[j]);
                            }
                        }
                    }
                }
            }
            indent += "    ";
            level++;
        }
    }
}

