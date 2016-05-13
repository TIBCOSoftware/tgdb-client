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
 * File name : ConnectionTest1.${EXT}
 * Created on: 1/13/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: ConnectionTest1.java 823 2016-05-12 12:47:02Z vchung $
 */


package com.tibco.tgdb.test;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.utils.SortedProperties;
import com.tibco.tgdb.utils.TGProperties;

public class ConnectionTest1 {
	public String url = "tcp://scott@localhost:8222";
    public String passwd = "scott";
    public TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    public boolean waitForExit = false;
    public boolean prefetchMetaData = false;
    public boolean testGetOperation = false;
    public boolean interactiveGet = false;
    public int depth = 5;
    public int printDepth = 5;
    public int resultCount = 100;
    public int edgeLimit = 0;
    
    String getStringValue(Iterator<String> argIter) {
    	while (argIter.hasNext()) {
    		String s = argIter.next();
    		return s;
    	}
    	return null;
    }
    
    String getStringValue(Iterator<String> argIter, String defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
    		return s;
    	}
    }

    int getIntValue(Iterator<String> argIter, int defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
    		try {
    			int i = Integer.valueOf(s);
    			return i;
    		} catch (NumberFormatException e) {
    			System.out.printf("Invalid number : %s\n", s);
    		}
    		return defaultValue;
    	}
    }

    void getArgs(String[] args) {
    	List<String> argList = Arrays.asList(args);
    	Iterator<String> argIter = argList.iterator();
    	while (argIter.hasNext()) {
    		String s = argIter.next();
    		System.out.printf("Arg : \"%s\"\n", s);
    		if (s.equalsIgnoreCase("-url")) {
    			url = getStringValue(argIter, "tcp://scott@localhost:8222");
    		} else if (s.equalsIgnoreCase("-password") || s.equalsIgnoreCase("-pw")) {
    			passwd = getStringValue(argIter, "scott");
    		} else if (s.equalsIgnoreCase("-loglevel") || s.equalsIgnoreCase("-ll")) {
    			String ll = getStringValue(argIter, "Debug");
    			try {
    				logLevel = TGLogger.TGLevel.valueOf(ll);
    			} catch(IllegalArgumentException e) {
    				System.out.printf("Invalid log level value '%s'...ignored\n", ll);
    			}
    		} else if (s.equalsIgnoreCase("-waitforexit") || s.equalsIgnoreCase("-we")) {
    			waitForExit = true;
    		} else if (s.equalsIgnoreCase("-prefetchmeta") || s.equalsIgnoreCase("-pm")) {
    			prefetchMetaData = true;
    		} else if (s.equalsIgnoreCase("-getentity") || s.equalsIgnoreCase("-ge")) {
    			testGetOperation = true;
    		} else if (s.equalsIgnoreCase("-getentityconsole") || s.equalsIgnoreCase("-gecon")) {
    			interactiveGet = true;
    		} else if (s.equalsIgnoreCase("-getentitydepth") || s.equalsIgnoreCase("-ged")) {
    			depth = getIntValue(argIter, 5);
    		} else if (s.equalsIgnoreCase("-getentitycount") || s.equalsIgnoreCase("-gec")) {
    			resultCount = getIntValue(argIter, 1000);
    		} else if (s.equalsIgnoreCase("-getedgelimit") || s.equalsIgnoreCase("-gel")) {
    			edgeLimit = getIntValue(argIter, 50);
    		} else if (s.equalsIgnoreCase("-printdepth") || s.equalsIgnoreCase("-pd")) {
    			printDepth = getIntValue(argIter, 5);
    		} else {
    			System.out.printf("Skip argument %s\n", s);
    		}
    	}
    }

    void printEntities(TGEntity ent, int maxDepth, int currDepth, String indent, Map<Integer, TGEntity> traverseMap) {
        if (currDepth == maxDepth) {
        	return;
        }
    	if (ent == null) {
    		return;
    	}
    	if (traverseMap.get(ent.hashCode()) != null) {
    		System.out.printf("%sKind : %s, hashcode : %d visited\n", indent, 
    			(ent.getEntityKind() == TGEntity.TGEntityKind.Node ? "Node" : "Edge"), ent.hashCode());
    		return;
    	} 
        traverseMap.put(ent.hashCode(), ent);
       	System.out.printf("%sKind : %s, hashcode : %d\n", indent, 
       		(ent.getEntityKind() == TGEntity.TGEntityKind.Node ? "Node" : "Edge"), ent.hashCode());
        for (TGAttribute attrib : ent.getAttributes()) {
        	System.out.printf("%s Attr : %s\n", indent, attrib.getValue());
        }
        if (ent.getEntityKind() == TGEntity.TGEntityKind.Node) {
        	int edgeCount = ((TGNode) ent).getEdges().size();
        	System.out.printf("%s Has %d edges\n", indent, edgeCount);
        	String newIndent = new String(indent).concat(" ");
        	for (TGEdge edge : ((TGNode) ent).getEdges()) {
        		printEntities(edge, maxDepth, currDepth, newIndent, traverseMap);
        	}
        } else if (ent.getEntityKind() == TGEntity.TGEntityKind.Edge) {
        	TGNode[] nodes = ((TGEdge) ent).getVertices();
        	if (nodes.length > 0) {
        		System.out.printf("%s Has end nodes\n", indent);
        	}
        	String newIndent = new String(indent).concat("  ");
        	++currDepth;
        	for (int j=0; j<nodes.length; j++) {
        		printEntities(nodes[j], maxDepth, currDepth, newIndent, traverseMap);
        	}
        }
    }

    void printEntitiesBreadth(TGNode node, int maxDepth) {
        String indent = "";
        Map<Integer, TGEntity> traverseMap = new HashMap<Integer, TGEntity>();
        List<TGNode> nodeList = new ArrayList<TGNode>();
        nodeList.add(node);

        while (nodeList.size() > 0) {
            int listSize = nodeList.size();
            while (listSize > 0) {
                listSize--;
                node = nodeList.get(0);
                nodeList.remove(0);
                if (traverseMap.get(node.hashCode()) != null) {
                    System.out.printf("%sNode(%d) visited\n", indent, node.hashCode());
                    continue;
                } 
                System.out.printf("%sNode(%d)\n", indent, node.hashCode());
                for (TGAttribute attrib : node.getAttributes()) {
                    System.out.printf("%s Attr : %s\n", indent, attrib.getValue());
                }
                traverseMap.put(node.hashCode(), node);
        	    System.out.printf("%s has %d edges\n", indent, node.getEdges().size());
                for (TGEdge edge : node.getEdges()) {
                    TGNode[] nodes = edge.getVertices();
                    if (nodes.length > 0) {
                        if (nodes[0] == null) {
                           	System.out.printf("%s   passed last edge\n", indent);
                            break;
                        }
                        for (int j=0; j<nodes.length; j++) {
                            if (nodes[j].hashCode() == node.hashCode()) {
                                continue;
                            } else {
                            	Iterator<TGAttribute> itr = edge.getAttributes().iterator();
                            	if (itr.hasNext()) {
                            		TGAttribute attrib = itr.next();
                            		System.out.printf("%s   edge(%d)(%s) with node(%d)\n", indent, edge.hashCode(), attrib.getValue(), nodes[j].hashCode());
                            	} else {
                            		System.out.printf("%s   edge(%d) with node(%d)\n", indent, edge.hashCode(), nodes[j].hashCode());
                            	}
                                nodeList.add(nodes[j]);
                            }
                        }
                    }
                }
            }
            indent += "    ";
        }
    }

    TGNode createNode(TGGraphObjectFactory gof, TGNodeType nodeType) {
    	if (nodeType != null) {
    		return gof.createNode(nodeType);
    	} else {
    		return gof.createNode();
    	}
    }

    void test() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
       	TGNodeType rateNodeType = null;
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

        if (prefetchMetaData) {
        	TGGraphMetadata gmd = conn.getGraphMetadata(true);
        	TGAttributeDescriptor attrDesc = gmd.getAttributeDescriptor("factor");
        	if (attrDesc != null) {
        		System.out.printf("'factor' has id %d and data type : %d\n", attrDesc.getAttributeId(), attrDesc.getType().typeId());
        	} else {
        		System.out.println("'factor' is not found from meta data fetch");
        	}
        	attrDesc = gmd.getAttributeDescriptor("level");
        	if (attrDesc != null) {
        		System.out.printf("'level' has id %d and data type : %d\n", attrDesc.getAttributeId(), attrDesc.getType().typeId());
        	} else {
        		System.out.println("'level' is not found from meta data fetch");
        	}
        	rateNodeType = gmd.getNodeType("testnode");
        	if (rateNodeType != null) {
        		System.out.printf("'testnode' is found with %d attributes\n", rateNodeType.getAttributeDescriptors().size());
        	} else {
        		System.out.println("'testnode' is not found from meta data fetch");
        	}
        }

      	System.out.println("Start transaction 1");
      	System.out.println("Create node1");
      	//createNode/Edge by itself does not include the it in the transaction.
        //Explicit call to insertEntity to add it to the transaction.
        TGNode node1 = createNode(gof, rateNodeType);
        node1.setAttribute("name", "john doe");
        node1.setAttribute("multiple", 7);
        node1.setAttribute("rate", 3.3);
        node1.setAttribute("nickname", "美麗");
        conn.insertEntity(node1);

      	System.out.println("Create node2");
        TGNode node2 = createNode(gof, rateNodeType);
        node2.setAttribute("name", "julie");
        node2.setAttribute("factor", 3.3);
        conn.insertEntity(node2);

      	System.out.println("Create node3");
        TGNode node3 = createNode(gof, rateNodeType);
        node3.setAttribute("name", "margo");
        node3.setAttribute("factor", 2.3);
        node3.setAttribute("是超人嗎", false);
        conn.insertEntity(node3);

      	System.out.println("Create edge1");
        TGEdge edge1 = gof.createEdge(node1, node2, TGEdge.DirectionType.BiDirectional);
        edge1.setAttribute("name", "spouse");
        conn.insertEntity(edge1);

      	System.out.println("Create edge2");
        TGEdge edge2 = gof.createEdge(node1, node3, TGEdge.DirectionType.Directed);
        edge2.setAttribute("name", "daughter");
        conn.insertEntity(edge2);

      	System.out.println("Commit transaction 1");
        conn.commit(); //----- write data to database ----------. Everything is create
      	System.out.println("Commit transaction 1 completed");

      	System.out.println("Start transaction 2");
        // updates
      	System.out.println("Update node1");
        node1.setAttribute("age", 40);
        node1.setAttribute("multiple", 8);
        conn.updateEntity(node1);
        
      	System.out.println("Delete edge1");
        conn.deleteEntity(edge1);
        
      	System.out.println("Update edge2");
        //add a new attribute
        edge2.setAttribute("extra", true);
        //update existing one
        edge2.setAttribute("name", "kid");
        conn.updateEntity(edge2);
        
      	System.out.println("Create node4");
        TGNode node4 = createNode(gof, rateNodeType);
        node4.setAttribute("name", "Smith");
        node4.setAttribute("level", 3.0);
        conn.insertEntity(node4);

      	System.out.println("Create edge3");
        TGEdge edge3 = gof.createEdge(node1, node4, TGEdge.DirectionType.BiDirectional);
        edge3.setAttribute("name", "Tennis Partner");
        conn.insertEntity(edge3);

      	System.out.println("Commit transaction 2");
        conn.commit(); // update John doe's value --------.
      	System.out.println("Commit transaction 2 completed");

      	if (testGetOperation) {
      		if (prefetchMetaData) {
      			conn.getGraphMetadata(true);
      		}
      		System.out.println("Start unique get operation");
      		TGKey key = gof.createCompositeKey("testnode");
      		key.setAttribute("name", "foo");
      		TGProperties<String, String> props = new SortedProperties<String, String>();
      		props.put("fetchsize", String.valueOf(resultCount));
      		props.put("traversaldepth", String.valueOf(depth));
      		props.put("edgelimit", String.valueOf(edgeLimit));
      		TGEntity ent = conn.getEntity(key, props);
      		if (ent != null) {
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for 'foo' returns nothing");
      		}
      		key.setAttribute("name", "margo");
      		ent = conn.getEntity(key, null);
      		if (ent != null) {
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for 'margo' returns nothing");
      		}
      		System.out.println("End unique get operation");
      	}

        if (waitForExit) {
        	System.out.println("Wait for any key to exit");
        	System.in.read();
        }

        conn.disconnect();
        System.out.println("Connection test connection disconnected.");
    }

    public void interactiveGet() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
      	conn.getGraphMetadata(true);
        
      	System.out.println("Start get console");
      	TGKey hkey = gof.createCompositeKey("hESCNodetype");
      	TGKey mkey = gof.createCompositeKey("mESCNodetype");
      	BufferedReader in = new BufferedReader(new InputStreamReader(System.in));
        String s;
        while ((s = in.readLine()) != null && s.length() != 0) {
      		System.out.printf("Received key : '%s'\n", s);
        	if (s.equalsIgnoreCase("quit")) {
        		break;
        	}
        	hkey.setAttribute("symbol", s);
      	 	TGProperties<String, String> props = new SortedProperties<String, String>();
      	 	props.put("fetchsize", String.valueOf(resultCount));
      	 	props.put("traversaldepth", String.valueOf(depth));
      		props.put("edgelimit", String.valueOf(edgeLimit));
      	 	TGEntity ent = conn.getEntity(hkey, props);
      	 	if (ent != null) {
      		 	//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		 	printEntitiesBreadth((TGNode) ent, printDepth);
      		 	System.out.printf("End of get %s\n", s);
      	 	} else {
      		 	System.out.printf("getEntity for %s returns nothing\n", s);
      	 	}
        }
        conn.disconnect();
        System.out.println("Exiting get entity console");
    }

    public static void main(String[] args) throws Exception {
    	ConnectionTest1 ct1 = new ConnectionTest1();

    	ct1.getArgs(args);
    	if (ct1.interactiveGet) {
    		ct1.interactiveGet();
    	} else {
    		ct1.test();
    	}
    }
}
