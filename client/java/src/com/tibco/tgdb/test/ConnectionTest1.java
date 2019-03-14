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
 * SVN Id: $Id: ConnectionTest1.java 2344 2018-06-11 23:21:45Z ssubrama $
 */


package com.tibco.tgdb.test;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.sql.Timestamp;
import java.time.LocalDateTime;
import java.util.*;

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
import com.tibco.tgdb.query.TGQueryOption;

public class ConnectionTest1 {
	public String url = "tcp://scott@localhost:8222";
    public String passwd = "scott";
    public TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    public boolean waitForExit = false;
    public boolean prefetchMetaData = false;
    public boolean testGetOperation = false;
    public boolean interactiveGet = false;
    public boolean testTransaction = false;
    public boolean testQuery = false;
    public int depth = 5;
    public int printDepth = 5;
    public int resultCount = 100;
    public int edgeLimit = 0;
    public String typeName = "hESCNodetype";
    public String keyName = "symbol";
    
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
    		} else if (s.equalsIgnoreCase("-testtransaction") || s.equalsIgnoreCase("-txn")) {
                testTransaction = true;
    		} else if (s.equalsIgnoreCase("-testquery") || s.equalsIgnoreCase("-qry")) {
                testQuery = true;
    		} else if (s.equalsIgnoreCase("-getentitydepth") || s.equalsIgnoreCase("-ged")) {
    			depth = getIntValue(argIter, 5);
    		} else if (s.equalsIgnoreCase("-getentitycount") || s.equalsIgnoreCase("-gec")) {
    			resultCount = getIntValue(argIter, 1000);
    		} else if (s.equalsIgnoreCase("-getedgelimit") || s.equalsIgnoreCase("-gel")) {
    			edgeLimit = getIntValue(argIter, 50);
    		} else if (s.equalsIgnoreCase("-gettypename") || s.equalsIgnoreCase("-getn")) {
    			typeName = getStringValue(argIter, "hESCNodeType");
    		} else if (s.equalsIgnoreCase("-getkeyname") || s.equalsIgnoreCase("-gek")) {
    			keyName = getStringValue(argIter, "symbol");
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
       	TGNodeType testNodeType = null;
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

        if (prefetchMetaData) {
        	TGGraphMetadata gmd = conn.getGraphMetadata(true);
        	TGAttributeDescriptor attrDesc = gmd.getAttributeDescriptor("factor");
        	if (attrDesc != null) {
        		System.out.printf("'factor' has id %d and data desc : %d\n", attrDesc.getAttributeId(), attrDesc.getType().typeId());
        	} else {
        		System.out.println("'factor' is not found from meta data fetch");
        	}
        	attrDesc = gmd.getAttributeDescriptor("level");
        	if (attrDesc != null) {
        		System.out.printf("'level' has id %d and data desc : %d\n", attrDesc.getAttributeId(), attrDesc.getType().typeId());
        	} else {
        		System.out.println("'level' is not found from meta data fetch");
        	}
        	testNodeType = gmd.getNodeType("testnode");
        	if (testNodeType != null) {
        		System.out.printf("'testnode' is found with %d attributes\n", testNodeType.getAttributeDescriptors().size());
        	} else {
        		System.out.println("'testnode' is not found from meta data fetch");
        	}
        }

      	System.out.println("Start transaction 1");
      	System.out.println("Create node1");
      	//createNode/Edge by itself does not include the it in the transaction.
        //Explicit call to insertEntity to add it to the transaction.
        TGNode node1 = createNode(gof, testNodeType);
        node1.setAttribute("name", "john doe");
        node1.setAttribute("multiple", 7);
        node1.setAttribute("rate", 3.3);
        node1.setAttribute("nickname", "美麗");
        conn.insertEntity(node1);

      	System.out.println("Create node2");
        TGNode node2 = createNode(gof, testNodeType);
        node2.setAttribute("name", "julie");
        node2.setAttribute("factor", 3.3);
        conn.insertEntity(node2);

      	System.out.println("Create node3");
        TGNode node3 = createNode(gof, testNodeType);
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
        TGNode node4 = createNode(gof, testNodeType);
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
            TGQueryOption option = TGQueryOption.createQueryOption();
            option.setPrefetchSize(resultCount);
            option.setTraversalDepth(depth);
            option.setEdgeLimit(edgeLimit);
            /*
      		TGProperties<String, String> props = new SortedProperties<String, String>();
      		props.put("fetchsize", String.valueOf(resultCount));
      		props.put("traversaldepth", String.valueOf(depth));
      		props.put("edgelimit", String.valueOf(edgeLimit));
      		*/
      		TGEntity ent = conn.getEntity(key, option);
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

      		key = gof.createCompositeKey("testnode");
      		key.setAttribute("nickname", "margo");
      		ent = conn.getEntity(key, null);
      		if (ent != null) {
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for nickname 'margo' returns nothing");
      		}
      		
      		key.setAttribute("nickname", "margo");
      		key.setAttribute("rate",  12.34);
      		ent = conn.getEntity(key, null);
      		if (ent != null) {
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for nickname 'margo' and rate - 12.34 returns nothing");
      		}
      		
      		key = gof.createCompositeKey("testnode");
      		key.setAttribute("multiple", 3);
      		ent = conn.getEntity(key, null);
      		if (ent != null) {
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for multiple - 3");
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

    public void testTransaction() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
       	TGNodeType testNodeType = null;
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

        if (prefetchMetaData) {
        	TGGraphMetadata gmd = conn.getGraphMetadata(true);
        	TGAttributeDescriptor attrDesc = gmd.getAttributeDescriptor("factor");
        	if (attrDesc != null) {
        		System.out.printf("'factor' has id %d and data desc : %d\n", attrDesc.getAttributeId(), attrDesc.getType().typeId());
        	} else {
        		System.out.println("'factor' is not found from meta data fetch");
        	}
        	attrDesc = gmd.getAttributeDescriptor("level");
        	if (attrDesc != null) {
        		System.out.printf("'level' has id %d and data desc : %d\n", attrDesc.getAttributeId(), attrDesc.getType().typeId());
        	} else {
        		System.out.println("'level' is not found from meta data fetch");
        	}
        	testNodeType = gmd.getNodeType("testnode");
        	if (testNodeType != null) {
        		System.out.printf("'testnode' is found with %d attributes\n", testNodeType.getAttributeDescriptors().size());
        	} else {
        		System.out.println("'testnode' is not found from meta data fetch");
        	}
        }

        System.out.println("Test Transaction : Insert Simple Node(John) of testnode desc with a few properties");
      	System.out.println("Create node");
        TGNode node1 = createNode(gof, testNodeType);
        node1.setAttribute("name", "Bruce Wayne");
        node1.setAttribute("multiple", 7);
        node1.setAttribute("rate", 5.5);
        node1.setAttribute("nickname", "超級英雄");
        node1.setAttribute("level", 4.0);
        node1.setAttribute("age", 38);
        conn.insertEntity(node1);
        TGNode node11 = createNode(gof, testNodeType);
        node11.setAttribute("name", "Peter Parker");
        node11.setAttribute("multiple", 7);
        node11.setAttribute("rate", 3.3);
        node11.setAttribute("nickname", "超級英雄");
        node11.setAttribute("level", 4.0);
        node11.setAttribute("age", 24);
        conn.insertEntity(node11);
        TGNode node12 = createNode(gof, testNodeType);
        node12.setAttribute("name", "Clark Kent");
        node12.setAttribute("multiple", 7);
        node12.setAttribute("rate", 6.6);
        node12.setAttribute("nickname", "超級英雄");
        node12.setAttribute("level", 4.0);
        node12.setAttribute("age", 32);
        conn.insertEntity(node12);
        TGNode node13 = createNode(gof, testNodeType);
        node13.setAttribute("name", "James Logan Howlett");
        node13.setAttribute("multiple", 7);
        node13.setAttribute("rate", 4.4);
        node13.setAttribute("nickname", "超級英雄");
        node13.setAttribute("level", 4.0);
        node13.setAttribute("age", 40);
        conn.insertEntity(node13);

        TGNode node2 = createNode(gof, testNodeType);
        node2.setAttribute("name", "Mary Jane Watson");
        node2.setAttribute("multiple", 14);
        node2.setAttribute("rate", 6.3);
        node2.setAttribute("nickname", "美麗");
        node2.setAttribute("level", 3.0);
        node2.setAttribute("age", 22);
        conn.insertEntity(node2);
        TGNode node21 = createNode(gof, testNodeType);
        node21.setAttribute("name", "Lois Lane");
        node21.setAttribute("multiple", 14);
        node21.setAttribute("rate", 6.4);
        node21.setAttribute("nickname", "美麗");
        node21.setAttribute("level", 3.0);
        node21.setAttribute("age", 30);
        conn.insertEntity(node21);
        TGNode node22 = createNode(gof, testNodeType);
        node22.setAttribute("name", "Jean Grey");
        node22.setAttribute("multiple", 14);
        node22.setAttribute("rate", 6.2);
        node22.setAttribute("nickname", "美麗");
        node22.setAttribute("level", 3.0);
        node22.setAttribute("age", 30);
        conn.insertEntity(node22);
        TGNode node23 = createNode(gof, testNodeType);
        node23.setAttribute("name", "Selina Kyle");
        node23.setAttribute("multiple", 14);
        node23.setAttribute("rate", 6.5);
        node23.setAttribute("nickname", "Criminal");
        node23.setAttribute("level", 3.0);
        node23.setAttribute("age", 30);
        conn.insertEntity(node23);
        TGNode node24 = createNode(gof, testNodeType);
        node24.setAttribute("name", "Harley Quinn");
        node24.setAttribute("multiple", 14);
        node24.setAttribute("rate", 5.8);
        node24.setAttribute("nickname", "Criminal");
        node24.setAttribute("level", 3.0);
        node24.setAttribute("age", 30);
        conn.insertEntity(node24);

        TGNode node3 = createNode(gof, testNodeType);
        node3.setAttribute("name", "Lex Luthor");
        node3.setAttribute("multiple", 52);
        node3.setAttribute("rate", 7.1);
        node3.setAttribute("nickname", "壞人");
        node3.setAttribute("level", 11.0);
        node3.setAttribute("age", 26);
        conn.insertEntity(node3);
        TGNode node31 = createNode(gof, testNodeType);
        node31.setAttribute("name", "Harvey Dent");
        node31.setAttribute("multiple", 52);
        node31.setAttribute("rate", 4.6);
        node31.setAttribute("nickname", "壞人");
        node31.setAttribute("level", 11.0);
        node31.setAttribute("age", 40);
        conn.insertEntity(node31);
        TGNode node32 = createNode(gof, testNodeType);
        node32.setAttribute("name", "Victor Creed");
        node32.setAttribute("multiple", 52);
        node32.setAttribute("rate", 6.4);
        node32.setAttribute("nickname", "壞人");
        node32.setAttribute("level", 11.0);
        node32.setAttribute("age", 40);
        conn.insertEntity(node32);
        TGNode node33 = createNode(gof, testNodeType);
        node33.setAttribute("name", "Norman Osborn");
        node33.setAttribute("multiple", 52);
        node33.setAttribute("rate", 6.3);
        node33.setAttribute("nickname", "壞人");
        node33.setAttribute("level", 11.0);
        node33.setAttribute("age", 50);
        conn.insertEntity(node33);

        TGEdge edge1 = gof.createEdge(node1, node31, TGEdge.DirectionType.BiDirectional);
        edge1.setAttribute("name", "Nemesis");
        conn.insertEntity(edge1);
        TGEdge edge2 = gof.createEdge(node1, node23, TGEdge.DirectionType.BiDirectional);
        edge2.setAttribute("name", "Frenemy");
        conn.insertEntity(edge2);
        TGEdge edge3 = gof.createEdge(node1, node24, TGEdge.DirectionType.BiDirectional);
        edge3.setAttribute("name", "Nemesis");
        conn.insertEntity(edge3);
        TGEdge edge4 = gof.createEdge(node1, node12, TGEdge.DirectionType.BiDirectional);
        edge4.setAttribute("name", "Teammate");
        conn.insertEntity(edge4);

        conn.commit();

      	System.out.println("Commit transaction completed");

      	if (testGetOperation == true) {
      		if (prefetchMetaData) {
      			conn.getGraphMetadata(true);
      		}
      		System.out.println("Start unique get operation");

      		TGKey key = gof.createCompositeKey("testnode");
      		key.setAttribute("name", "John Doe");
            TGQueryOption option = TGQueryOption.createQueryOption();
            option.setPrefetchSize(resultCount);
            option.setTraversalDepth(depth);
            option.setEdgeLimit(edgeLimit);
            /*
      		TGProperties<String, String> props = new SortedProperties<String, String>();
      		props.put("fetchsize", String.valueOf(resultCount));
      		props.put("traversaldepth", String.valueOf(depth));
      		props.put("edgelimit", String.valueOf(edgeLimit));
      		*/
      		TGEntity ent = conn.getEntity(key, option);
      		if (ent != null) {
      			System.out.println("getEntity for 'John Doe' - found - wrong");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for 'John Doe' returns nothing - good");
      		}

      		key.setAttribute("name", "Bruce Wayne");
      		ent = conn.getEntity(key, null);
      		if (ent != null) {
      			System.out.println("getEntity for 'Bruce Wayne' - found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for 'Bruce Wayne' returns nothing - wrong");
      		}
      		key.setAttribute("name", "Peter Parker");
      		ent = conn.getEntity(key, null);
      		if (ent != null) {
      			System.out.println("getEntity for 'Peter Parker' - found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for 'Peter Parker' returns nothing - wrong");
      		}
      		key.setAttribute("name", "Mary Jane Watson");
      		ent = conn.getEntity(key, null);
      		if (ent != null) {
      			System.out.println("getEntity for 'Mary Jane Watson' - found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for 'Mary Jane' returns nothing - wrong");
      		}
      		key.setAttribute("name", "Super Jane");
      		ent = conn.getEntity(key, null);
      		if (ent != null) {
      			System.out.println("getEntity for 'Super Jane' - found - wrong");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for 'Super Jane' returns nothing - good");
      		}

      		TGKey key1 = gof.createCompositeKey("testnode");
      		key1.setAttribute("nickname", "Stupid");
      		ent = conn.getEntity(key1, null);
      		if (ent != null) {
      			System.out.println("getEntity for 'Stupid' - found - wrong");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for nickname 'Stupid' returns nothing - good");
      		}
      		key1.setAttribute("nickname", "美麗");
      		ent = conn.getEntity(key1, null);
      		if (ent != null) {
      			System.out.println("getEntity for '美麗' - found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for nickname '美麗' returns nothing - wrong");
      		}
      		key1.setAttribute("nickname", "Criminal");
      		ent = conn.getEntity(key1, null);
      		if (ent != null) {
      			System.out.println("getEntity for 'Criminal' - found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for nickname 'Criminal' returns nothing - wrong");
      		}
      		key1.setAttribute("nickname", "超級英雄");
      		ent = conn.getEntity(key1, null);
      		if (ent != null) {
      			System.out.println("getEntity for '超級英雄' - found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for nickname '超級英雄' returns nothing - wrong");
      		}

      		TGKey key2 = gof.createCompositeKey("testnode");
      		key2.setAttribute("nickname", "壞人");
      		key2.setAttribute("level", 11.0);
      		key2.setAttribute("rate", 7.1);
      		ent = conn.getEntity(key2, null);
      		if (ent != null) {
      			System.out.println("getEntity for '壞人'(11.0)(7.1) found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for '壞人'(11.0)(7.1) returns nothing - wrong");
      		}

      		TGKey key3 = gof.createCompositeKey("testnode");
      		key3.setAttribute("multiple", 52);
      		ent = conn.getEntity(key3, null);
      		if (ent != null) {
      			System.out.println("getEntity for multiple 52 found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for multiple 52 returns nothing - wrong");
      		}
      		System.out.println("End unique get operation");

      		System.out.println("Start update test");
            node13.setAttribute("age", 43);
            //node13.setAttribute("nickname", "Wolverine");
            node13.setAttribute("nickname", "超級");
            conn.updateEntity(node13);

            //node32.setAttribute("nickname", "Sabertooth");
            node32.setAttribute("nickname", "害怕");
            node32.setAttribute("level", 11.0);
            node32.setAttribute("age", 38);
            conn.updateEntity(node32);
            conn.commit();
      		System.out.println("End update test");

      		System.out.println("Get update entity");
      		key2.setAttribute("nickname", "害怕");
      		key2.setAttribute("level", 11.0);
      		key2.setAttribute("rate", 6.4);
      		ent = conn.getEntity(key2, null);
      		if (ent != null) {
      			System.out.println("getEntity for 'Sabertooth'(11.0)(6.4) found - good");
      			printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
      		} else {
      			System.out.println("getEntity for 'Sabertooth'(11.0)(6.4) returns nothing - wrong");
      		}
      		System.out.println("Get update entity end");
      	}

        if (testQuery == true) {
            System.out.println("Start test query");
            //TGQuery Query1 = conn.createQuery("testquery < X '5ef';");
            //Query1.execute();
            //Query1.close();
            conn.executeQuery("@nodetype = 'testnode' and ((nickname = '壞人' and level = 11.0) or (level = 3.0));",
                TGQueryOption.DEFAULT_QUERY_OPTION);
            System.out.println("End test query");
        }


        conn.disconnect();
        System.out.println("Transaction test connection disconnected.");
    }

    public void interactiveGet() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
        System.out.printf("Depth : %d, Fetch : %d, Edge : %d, desc : %s, key : %s\n", depth, resultCount,
            edgeLimit, typeName, keyName);
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
      	conn.getGraphMetadata(true);
        
      	System.out.println("Start get console");
      	TGKey hkey = gof.createCompositeKey(typeName);
      	BufferedReader in = new BufferedReader(new InputStreamReader(System.in));
        String s;
        while ((s = in.readLine()) != null && s.length() != 0) {
      		System.out.printf("Received key : '%s'\n", s);
        	if (s.equalsIgnoreCase("quit")) {
        		break;
        	}
        	hkey.setAttribute(keyName, s);
            TGQueryOption option = TGQueryOption.createQueryOption();
            option.setPrefetchSize(resultCount);
            option.setTraversalDepth(depth);
            option.setEdgeLimit(edgeLimit);
            /*
      	 	TGProperties<String, String> props = new SortedProperties<String, String>();
      	 	props.put("fetchsize", String.valueOf(resultCount));
      	 	props.put("traversaldepth", String.valueOf(depth));
      		props.put("edgelimit", String.valueOf(edgeLimit));
      		*/
      	 	TGEntity ent = conn.getEntity(hkey, option);
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
    	} else if (ct1.testTransaction) {
            ct1.testTransaction();
        }else {
    		ct1.test();
    	}
    }
}
