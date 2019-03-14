package com.tibco.tgdb.test;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.query.TGQuery;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.query.TGResultSet;

public class QueryTest {

    public TGConnection conn = null;
    public String url = "tcp://scott@localhost:8222";
    public String passwd = "scott";

    TGNode createNode(TGGraphObjectFactory gof, TGNodeType nodeType) {
    	if (nodeType != null) {
    		return gof.createNode(nodeType);
    	} else {
    		return gof.createNode();
    	}
    }

    public void createData() throws Exception {
        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
       	TGNodeType testNodeType = null;
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        }

    	TGGraphMetadata gmd = conn.getGraphMetadata(true);
        testNodeType = gmd.getNodeType("testnode");
        if (testNodeType != null) {
            System.out.printf("'testnode' is found with %d attributes\n", testNodeType.getAttributeDescriptors().size());
        } else {
            System.out.println("'testnode' is not found from meta data fetch");
            return;
        }

      	System.out.println("Create test nodes");
        TGNode node1 = createNode(gof, testNodeType);
        node1.setAttribute("name", "Bruce Wayne");
        node1.setAttribute("multiple", 7);
        node1.setAttribute("nickname", "Superhero");
        node1.setAttribute("level", 4.0);
        node1.setAttribute("rate", 5.5);
        node1.setAttribute("age", 40);
        conn.insertEntity(node1);
        TGNode node11 = createNode(gof, testNodeType);
        node11.setAttribute("name", "Peter Parker");
        node11.setAttribute("multiple", 7);
        node11.setAttribute("nickname", "Superhero");
        node11.setAttribute("level", 4.0);
        node11.setAttribute("rate", 3.3);
        node11.setAttribute("age", 24);
        conn.insertEntity(node11);
        TGNode node12 = createNode(gof, testNodeType);
        node12.setAttribute("name", "Clark Kent");
        node12.setAttribute("multiple", 7);
        node12.setAttribute("nickname", "Superhero");
        node12.setAttribute("level", 4.0);
        node12.setAttribute("rate", 6.6);
        node12.setAttribute("age", 32);
        conn.insertEntity(node12);
        TGNode node13 = createNode(gof, testNodeType);
        node13.setAttribute("name", "James Logan Howlett");
        node13.setAttribute("multiple", 7);
        node13.setAttribute("nickname", "Superhero");
        node13.setAttribute("level", 4.0);
        node13.setAttribute("rate", 4.4);
        node13.setAttribute("age", 40);
        conn.insertEntity(node13);
        TGNode node14 = createNode(gof, testNodeType);
        node14.setAttribute("name", "Diana Prince");
        node14.setAttribute("multiple", 7);
        node14.setAttribute("nickname", "Superheroine");
        node14.setAttribute("level", 4.0);
        node14.setAttribute("rate", 5.9);
        node14.setAttribute("age", 35);
        conn.insertEntity(node14);
        TGNode node15 = createNode(gof, testNodeType);
        node15.setAttribute("name", "Jean Grey");
        node15.setAttribute("multiple", 7);
        node15.setAttribute("nickname", "Superheroine");
        node15.setAttribute("level", 4.0);
        node15.setAttribute("rate", 6.2);
        node15.setAttribute("age", 35);
        conn.insertEntity(node15);

        TGNode node2 = createNode(gof, testNodeType);
        node2.setAttribute("name", "Mary Jane Watson");
        node2.setAttribute("multiple", 14);
        node2.setAttribute("nickname", "Girlfriend");
        node2.setAttribute("level", 3.0);
        node2.setAttribute("rate", 6.3);
        node2.setAttribute("age", 22);
        conn.insertEntity(node2);
        TGNode node21 = createNode(gof, testNodeType);
        node21.setAttribute("name", "Lois Lane");
        node21.setAttribute("multiple", 14);
        node21.setAttribute("nickname", "Girlfriend");
        node21.setAttribute("level", 3.0);
        node21.setAttribute("rate", 6.4);
        node21.setAttribute("age", 30);
        conn.insertEntity(node21);
        TGNode node22 = createNode(gof, testNodeType);
        node22.setAttribute("name", "Raven Darkhölme");
        node22.setAttribute("multiple", 14);
        node22.setAttribute("nickname", "Criminal");
        node22.setAttribute("level", 3.0);
        node22.setAttribute("rate", 7.2);
        node22.setAttribute("age", 36);
        conn.insertEntity(node22);
        TGNode node23 = createNode(gof, testNodeType);
        node23.setAttribute("name", "Selina Kyle");
        node23.setAttribute("multiple", 14);
        node23.setAttribute("nickname", "Criminal");
        node23.setAttribute("level", 3.0);
        node23.setAttribute("rate", 6.5);
        node23.setAttribute("age", 30);
        conn.insertEntity(node23);
        TGNode node24 = createNode(gof, testNodeType);
        node24.setAttribute("name", "Harley Quinn");
        node24.setAttribute("multiple", 14);
        node24.setAttribute("nickname", "Criminal");
        node24.setAttribute("level", 3.0);
        node24.setAttribute("rate", 5.8);
        node24.setAttribute("age", 30);
        conn.insertEntity(node24);

        TGNode node3 = createNode(gof, testNodeType);
        node3.setAttribute("name", "Lex Luthor");
        node3.setAttribute("multiple", 52);
        node3.setAttribute("nickname", "Bad guy");
        node3.setAttribute("level", 11.0);
        node3.setAttribute("rate", 7.1);
        node3.setAttribute("age", 26);
        conn.insertEntity(node3);
        TGNode node31 = createNode(gof, testNodeType);
        node31.setAttribute("name", "Harvey Dent");
        node31.setAttribute("multiple", 52);
        node31.setAttribute("nickname", "Bad guy");
        node31.setAttribute("level", 11.0);
        node31.setAttribute("rate", 4.6);
        node31.setAttribute("age", 40);
        conn.insertEntity(node31);
        TGNode node32 = createNode(gof, testNodeType);
        node32.setAttribute("name", "Victor Creed");
        node32.setAttribute("multiple", 52);
        node32.setAttribute("nickname", "Bad guy");
        node32.setAttribute("level", 11.0);
        node32.setAttribute("rate", 6.4);
        node32.setAttribute("age", 40);
        conn.insertEntity(node32);
        TGNode node33 = createNode(gof, testNodeType);
        node33.setAttribute("name", "Norman Osborn");
        node33.setAttribute("multiple", 52);
        node33.setAttribute("nickname", "Bad guy");
        node33.setAttribute("level", 11.0);
        node33.setAttribute("rate", 6.3);
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

        TGEdge edge5 = gof.createEdge(node13, node15, TGEdge.DirectionType.BiDirectional);
        edge5.setAttribute("name", "Friend");
        conn.insertEntity(edge5);
        TGEdge edge6 = gof.createEdge(node13, node22, TGEdge.DirectionType.BiDirectional);
        edge6.setAttribute("name", "Enemy");
        conn.insertEntity(edge6);
        TGEdge edge7 = gof.createEdge(node13, node32, TGEdge.DirectionType.BiDirectional);
        edge7.setAttribute("name", "Enemy");
        conn.insertEntity(edge7);

        conn.commit();

      	System.out.println("Commit transaction completed");
    }

    // name is primary key
    public void testPKey() throws Exception {
        System.out.println("Start unique query");
        // use primary key
        TGResultSet resultSet = conn.executeQuery("@nodetype = 'testnode' and name = 'Lex Luthor';", TGQueryOption.DEFAULT_QUERY_OPTION);
        if (resultSet != null) {
            int i=1;
            while (resultSet.hasNext()) {
                TGNode node = (TGNode) resultSet.next();
                System.out.printf("Node : %d\n", i++);
                printEntitiesBreadth(node, 5);
            }
        }
        System.out.println("End unique query");
    }

    public void testPartialUnique() throws Exception {
        System.out.println("Start partial unique query");
        // use unique key testidx1
        TGResultSet resultSet = conn.executeQuery("@nodetype = 'testnode' and nickname = 'Bad guy' and level = 11.0;", TGQueryOption.DEFAULT_QUERY_OPTION);
        if (resultSet != null) {
            int i=1;
            while (resultSet.hasNext()) {
                TGNode node = (TGNode) resultSet.next();
                System.out.printf("Node : %d\n", i++);
                printEntitiesBreadth(node, 5);
            }
        }
        System.out.println("End partial unique query");
    }

    public void testNonunique() throws Exception {
        System.out.println("Start non-unique query");
        // use nonunique key testidx2
        TGResultSet resultSet = conn.executeQuery("@nodetype = 'testnode' and nickname = 'Superhero';", TGQueryOption.DEFAULT_QUERY_OPTION);
        if (resultSet != null) {
            int i=1;
            while (resultSet.hasNext()) {
                TGNode node = (TGNode) resultSet.next();
                System.out.printf("Node : %d\n", i++);
                printEntitiesBreadth(node, 5);
            }
        }
        System.out.println("End non-unique query");
    }

    public void testGreaterThan() throws Exception {
        System.out.println("Start GT query");
        // use testidx3 partially
        TGResultSet resultSet = conn.executeQuery("@nodetype = 'testnode' and age > 32;", TGQueryOption.DEFAULT_QUERY_OPTION);
        if (resultSet != null) {
            int i=1;
            while (resultSet.hasNext()) {
                TGNode node = (TGNode) resultSet.next();
                System.out.printf("Node : %d\n", i++);
                printEntitiesBreadth(node, 5);
            }
        }
        System.out.println("End GT query");
    }

    public void testLessThan() throws Exception {
        System.out.println("Start LT query");
        // use testidx3 partially
        TGResultSet resultSet = conn.executeQuery("@nodetype = 'testnode' and age < 28;", TGQueryOption.DEFAULT_QUERY_OPTION);
        if (resultSet != null) {
            int i=1;
            while (resultSet.hasNext()) {
                TGNode node = (TGNode) resultSet.next();
                System.out.printf("Node : %d\n", i++);
                printEntitiesBreadth(node, 5);
            }
        }
        System.out.println("End LT query");
    }

    public void testRange() throws Exception {
        System.out.println("Start range query");
        // use testidx3 
        TGResultSet resultSet = conn.executeQuery("@nodetype = 'testnode' and age > 28 and age < 33 and level > 2.9 and level < 4.1;", TGQueryOption.DEFAULT_QUERY_OPTION);
        if (resultSet != null) {
            int i=1;
            while (resultSet.hasNext()) {
                TGNode node = (TGNode) resultSet.next();
                System.out.printf("Node : %d\n", i++);
                printEntitiesBreadth(node, 5);
            }
        }
        System.out.println("End range query");
    }

    public void testComplex() throws Exception {
        System.out.println("Start complex query");
//        TGQuery Query1 = conn.createQuery("testquery < X '5ef';");
//        Query1.execute();
        TGResultSet resultSet = conn.executeQuery("@nodetype =     'testnode'    and nickname = 'Criminal' and level = 3.0 and (rate = 5.8 or rate = 6.5);", TGQueryOption.DEFAULT_QUERY_OPTION);
        resultSet = conn.executeQuery("@nodetype = 'testnode' and (nickname = 'Criminal' or nickname = 'foo') and level = 3.0 and (rate = 5.8 or rate = 6.5);", TGQueryOption.DEFAULT_QUERY_OPTION);
        resultSet = conn.executeQuery("@nodetype = 'testnode' and (nickname = 'Criminal' and level = 3.0 and rate = 5.8) or (nickname = '壞人' and level = 11.0 and rate = 6.3);", TGQueryOption.DEFAULT_QUERY_OPTION);
        resultSet = conn.executeQuery("@nodetype = 'testnode' and ((nickname = 'Criminal' and level = 3.0 and rate = 5.8) or (nickname = '壞人' and level = 11.0 and rate = 6.3));", TGQueryOption.DEFAULT_QUERY_OPTION);
        //resultSet = conn.executeQuery("@nodetype = 'testnode' and level < 5.2 and rate < 4.5;", TGQueryOption.DEFAULT_QUERY_OPTION);
        //resultSet = conn.executeQuery("@nodetype = 'testnode' and level < 5.2 and rate < 4.5;", TGQueryOption.DEFAULT_QUERY_OPTION);
//        Query1.close();
//        resultSet = conn.executeQuery("testquery < X '5ef';", TGQueryOption.DEFAULT_QUERY_OPTION);

        System.out.println("End complex query");
    }

    public void testNodeTypeOnly() throws Exception {
        System.out.println("Start node desc only query");
        TGResultSet resultSet = conn.executeQuery("@nodetype = 'testnode';", TGQueryOption.DEFAULT_QUERY_OPTION);
        if (resultSet != null) {
            int i=1;
            while (resultSet.hasNext()) {
                TGNode node = (TGNode) resultSet.next();
                System.out.printf("Node : %d\n", i++);
                printEntitiesBreadth(node, 5);
            }
        }
        System.out.println("End node desc only query");
    }

    public void connect() throws Exception {
        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);
        conn.connect();
    }

    public void disconnect() throws Exception {
        conn.disconnect();
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

    public static void main(String[] args) throws Exception {
        QueryTest qt = new QueryTest();
        qt.connect();

        qt.createData();
        qt.testPKey();
        qt.testPartialUnique();
        qt.testNonunique();
        qt.testGreaterThan();
        qt.testLessThan();
        qt.testRange();
        qt.testNodeTypeOnly();
        //testComplex();

        qt.disconnect();
        System.out.println("Disconnected.");
    }
}
