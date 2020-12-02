/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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

 * File name : ClientTest1.${EXT}
 * Created on: 03/27/2018
 * Created by: vincent
 * SVN Id: $Id: ClientTest1.java 3996 2020-05-16 01:09:49Z vchung $
 */

package com.tibco.tgdb.test;

import java.util.Arrays;
import java.util.Collection;
import java.util.Iterator;
import java.util.List;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.utils.EntityUtils;
import com.tibco.tgdb.query.TGQueryOption;

/*
 * This test uses the generic initdb.conf that comes with the product.
 * However, please add the following line to initdb.conf before running the test.
 * nameontestidx   = @attrs:name @unique:true @ontype:testnode
 */
public class ClientTest1 {
	private String url = "tcp://scott@localhost:8222/{dbname=demodb}";
    private String passwd = "scott";
    private TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    private int depth = 5;
    private int printDepth = 5;
    private int resultCount = 100;
    private int edgeLimit = 0;
    private boolean initDB = false;


    private String getStringValue(Iterator<String> argIter) {
    	while (argIter.hasNext()) {
    		String s = argIter.next();
    		return s;
    	}
    	return null;
    }
    
    private String getStringValue(Iterator<String> argIter, String defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
    		return s;
    	}
    }

    private int getIntValue(Iterator<String> argIter, int defaultValue) {
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

    private void getArgs(String[] args) {
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
    		} else if (s.equalsIgnoreCase("-initdb") || s.equalsIgnoreCase("-i")) {
    			initDB = true;
    		} else {
    			System.out.printf("Skip argument %s\n", s);
    		}
    	}
    }

    private TGNode createNode(TGGraphObjectFactory gof, TGNodeType nodeType) {
    	if (nodeType != null) {
    		return gof.createNode(nodeType);
    	} else {
    		return gof.createNode();
    	}
    }
    
    private void testGetByKey() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

       	boolean exceptionThrown = false;
       	try {
       		TGKey key = gof.createCompositeKey("No good desc");
       		//TGKey key = gof.createCompositeKey("testnode");
       	} catch (TGException e) {
       		exceptionThrown = true;
        	System.out.printf("Exception : %s\n", e.getMessage());
       	}
       	if (exceptionThrown) {
        	System.out.println("Exception received - good");
       	} else {
        	System.out.println("Exception not received - bad");
       	}
        conn.disconnect();
    }

    private void testGetNodeTypes() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

       	TGGraphMetadata gmd = conn.getGraphMetadata(false);
       	Collection<TGNodeType> types = gmd.getNodeTypes();
        System.out.println("Show all node types");
       	for (TGNodeType type : types) {
       		System.out.printf("Type : %s\n", type.getName());
       	}
        System.out.println("Show all node types end");
        //Thread.sleep(30000);
        conn.disconnect();
    }

    //This method is replaced by setupData
    private void test() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
       	TGNodeType nullNodeType = null;
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

       	TGGraphMetadata gmd = conn.getGraphMetadata(false);
        nullNodeType = gmd.getNodeType("nulltestnode");
        if (nullNodeType != null) {
        	System.out.printf("'nulltestnode' is found with %d attributes\n", nullNodeType.getAttributeDescriptors().size());
        } else {
        	System.out.println("'nulltestnode' is not found from meta data fetch");
        }

      	System.out.println("Start transaction 1");
      	System.out.println("Create node1");
      	//createNode/Edge by itself does not include the it in the transaction.
        //Explicit call to insertEntity to add it to the transaction.
        TGNode node1 = createNode(gof, nullNodeType);
        node1.setAttribute("name", "john doe");
        node1.setAttribute("multiple", 7);
        node1.setAttribute("rate", 3.3);
        node1.setAttribute("nickname", "美麗");
        conn.insertEntity(node1);
      	System.out.println("Commit transaction 1");
        conn.commit(); //----- write data to database ----------. Everything is create
      	System.out.println("Commit transaction 1 completed");
        conn.disconnect();
        System.out.println("Connection test connection disconnected.");
    }

    private void setupData() throws Exception {
    	if (initDB == false) {
    		return;
    	}
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
       	TGNodeType testnodetype = null;
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

       	TGGraphMetadata gmd = conn.getGraphMetadata(false);
        testnodetype = gmd.getNodeType("testnode");
        if (testnodetype != null) {
        	System.out.printf("'testnode' is found with %d attributes\n", testnodetype.getAttributeDescriptors().size());
        } else {
        	System.out.println("'testnode' is not found from meta data fetch");
        }

      	System.out.println("Start transaction 1");
      	System.out.println("Create node1");
      	//createNode/Edge by itself does not include the it in the transaction.
        //Explicit call to insertEntity to add it to the transaction.
        TGNode node1 = createNode(gof, testnodetype);
        node1.setAttribute("name", "john doe");
        node1.setAttribute("multiple", 7);
        node1.setAttribute("rate", 3.3);
        node1.setAttribute("nickname", "美麗");
        conn.insertEntity(node1);
      	System.out.println("Create node2");
        TGNode node2 = createNode(gof, testnodetype);
        node2.setAttribute("name", "Jane doe");
        node2.setAttribute("multiple", 5);
        node2.setAttribute("rate", 8.3);
        node2.setAttribute("nickname", "black widow");
        conn.insertEntity(node2);
      	System.out.println("Create edge1");
      	TGEdgeType edgeType = gmd.getEdgeType("basicedge");
      	if (edgeType != null) {
      		System.out.println("Edge desc 'basicedge' found");
      	} else {
      		System.out.println("Edge desc 'basicedge' not found - bad");
      	}
      	TGEdge edge1 = gof.createEdge(node1, node2, edgeType);
        edge1.setAttribute("name", "spouse");
        conn.insertEntity(edge1);
      	TGEdge edge2 = gof.createEdge(node2, node1, edgeType);
        edge2.setAttribute("name", "partner");
        conn.insertEntity(edge2);
      	TGEdge edge = gof.createEdge(node1, node1, edgeType);
        edge.setAttribute("name", "toSelf1");
        edge.setAttribute("age", 100);
        conn.insertEntity(edge);
      	edge = gof.createEdge(node1, node1, edgeType);
        edge.setAttribute("name", "toSelf2");
        edge.setAttribute("age", 200);
        conn.insertEntity(edge);
      	edge = gof.createEdge(node1, node1, edgeType);
        edge.setAttribute("name", "toSelf3");
        edge.setAttribute("age", 300);
        conn.insertEntity(edge);
      	edge = gof.createEdge(node1, node1, edgeType);
        edge.setAttribute("name", "toSelf4");
        edge.setAttribute("age", 400);
        conn.insertEntity(edge);
      	edge = gof.createEdge(node1, node1, edgeType);
        edge.setAttribute("name", "toSelf5");
        edge.setAttribute("age", 500);
        conn.insertEntity(edge);
      	System.out.println("Commit transaction 1");
        conn.commit(); //----- write data to database ----------. Everything is create
      	System.out.println("Commit transaction 1 completed");
        conn.disconnect();
        System.out.println("Connection test connection disconnected.");
    }

    private void testGetEdges() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

       	TGGraphMetadata gmd = conn.getGraphMetadata(false);
       	TGEdgeType edgeType = gmd.getEdgeType("basicedge");

       	try {
       		TGQueryOption option = TGQueryOption.DEFAULT_QUERY_OPTION;
       		TGKey key = gof.createCompositeKey("testnode");
       		key.setAttribute("name", "john doe");
       		TGNode node = (TGNode) conn.getEntity(key, option);
       		//get all the edges
       		getEdges(node, null, TGEdge.Direction.Any);
       		getEdges(node, null, TGEdge.Direction.Outbound);
       		getEdges(node, null, TGEdge.Direction.Inbound);
       		getEdges(node, edgeType, TGEdge.Direction.Any);
       		getEdges(node, edgeType, TGEdge.Direction.Outbound);
       		getEdges(node, edgeType, TGEdge.Direction.Inbound);
       	} catch (TGException e) {
        	System.out.printf("Exception : %s\n", e.getMessage());
       	}
        conn.disconnect();
    }

    private void testSelfTraversal() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        String query = "@nodetype = 'testnode' and name = 'john doe';";
        String traversal = "@edgetype = 'basicedge' and @edge.name = 'toSelf3';";
        String end = "@nodetype = 'testnode' and name = 'john doe';";

       	try {
       		TGQueryOption option = TGQueryOption.DEFAULT_QUERY_OPTION;
       		TGResultSet result = conn.executeQuery(query,  null, traversal, end, TGQueryOption.DEFAULT_QUERY_OPTION);

       		if (result != null) {
       			while (result.hasNext()) {
       				TGNode node = (TGNode) result.next();
       				EntityUtils.printEntitiesBreadth(node, 3);
       			}
       		}
       	} catch (TGException e) {
        	System.out.printf("Exception : %s\n", e.getMessage());
       	}
        conn.disconnect();
    }

    private void getEdges(TGNode node, TGEdgeType edgeType, TGEdge.Direction direction) {
       if (node == null) {
    	   System.out.println("node is null");
    	   return;
       }
       Collection<TGEdge> edges = node.getEdges(edgeType, direction);
       System.out.printf("Get edge desc : '%s' with direction : '%s'\n", edgeType == null ? "All edge types" : edgeType.getName(), direction.name());
       for (TGEdge edge : edges) {
    	   System.out.printf("Edge name : %s\n", edge.getAttribute("name").getAsString());
       }
    }

    public static void main(String[] args) throws Exception {
    	ClientTest1 nt = new ClientTest1();
    	nt.getArgs(args);
    	nt.setupData();
    	nt.testGetNodeTypes();
    	nt.testGetByKey();
    	nt.testGetEdges();
    	//nt.testSelfTraversal();
    }
}

