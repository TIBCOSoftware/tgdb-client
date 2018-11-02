/**
 * Copyright (c) 2016 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : TransactionUnitTest.${EXT}
 * Created on: 10/2/16
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.utils.SortedProperties;
import com.tibco.tgdb.utils.TGProperties;

import java.io.DataOutputStream;
import java.math.BigDecimal;
import java.util.Calendar;

public class TransactionUnitTest {

    public String url = "tcp://scott@localhost:8222";
    //public String url = "tcp://scott@10.98.201.111:8228/{connectTimeout=30}";
    //public String url = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
    //public String url = "tcp://scott@localhost6:8228";
    public String passwd = "scott";
    public TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    TGGraphObjectFactory gof;
    TGGraphMetadata gmd;
    TGConnection conn;
    TGNodeType basicNodeType, rateNodeType;
    TGNode john, smith, kelly;
    TGEdge brother, wife;
    
    
    public TransactionUnitTest(String args[])
    {
        TGLogger logger = TGLogManager.getInstance().getLogger();
        logger.setLevel(logLevel);
        parseArgs(args);
    }
    
    public void connect() throws TGException {
        System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        gof = conn.getGraphObjectFactory();
        gmd = conn.getGraphMetadata(true);
        basicNodeType = gmd.getNodeType("basicnode");
        if (basicNodeType == null) throw new TGException("Node type basicnode not found");

        rateNodeType = gmd.getNodeType("ratenode");
        if (rateNodeType == null) throw new TGException("Node type ratenode not found");
        
        
    }
    
    public void runTestCases() throws TGException
    {
        testCase1();

        testCase1_1();
        testCase1_2();
        testCase2();

        testCase3();
        testCase3_1();
        testCase4();
        testCase4_0();
        testCase5();

        testCase6();
        testCase6_1();
        //testCase7();

    }

    TGNode createNode(TGNodeType nodeType) {
        if (nodeType != null) {
            return gof.createNode(nodeType);
        } else {
            return gof.createNode();
        }
    }
    
    private void testCase1() throws TGException
    {
        try {
            System.out.println("Test Case 1: Insert Simple Node(John) of basicnode with a few properties");

            TGNode node = createNode(basicNodeType);

            node.setAttribute("name", "john"); //name is the primary key
            node.setAttribute("age", 30);
            node.setAttribute("nickname", "美麗");
            node.setAttribute("createtm", new Calendar
                    .Builder()
                    .setDate(2016, 10, 25)
                    .setTimeOfDay(8, 9, 30, 999)
                    .build());

            node.setAttribute("networth", new BigDecimal(2378989.567));
            node.setAttribute("flag", 'D');

            conn.insertEntity(node);
            conn.commit();
            john = node;
        }
        catch (TGException e) {
            e.printStackTrace();
        }
        
    }

    private void testCase1_1() throws TGException
    {
        System.out.println("Test Case 1_0: Get the Entity that we inserted");
        TGKey key = gof.createCompositeKey("basicnode");
        key.setAttribute("name", "john");
        TGQueryOption option = TGQueryOption.createQueryOption();
        //option.setPrefetchSize(0); //Test for Server Crash
        /*
        TGProperties props = new SortedProperties();
        props.put("fetchsize", "-1");
        props.put("traversaldepth", "-1");
        props.put("edgelimit", "-1");
        */
        TGEntity entity = conn.getEntity(key, option);
        if (entity instanceof TGNode) {
            System.out.println("John's age is :" + entity.getAttribute("age").getAsInt());
            System.out.println("John's createtm:" + entity.getAttribute("createtm").getValue().toString());
            System.out.println("John's networth:" + entity.getAttribute("networth").getValue().toString());
        }
        john = TGNode.class.cast(entity);

    }

    private void testCase1_2() throws TGException
    {
        try {
            System.out.println("Test Case 1_1: Again insert John. This should raise Unique Key Constraint violation.");
            TGNode node = createNode(basicNodeType);
            node.setAttribute("name", "john");
            node.setAttribute("age", 30);
            node.setAttribute("nickname", "美麗");
            conn.insertEntity(node);
            conn.commit();
        }
        catch (TGException e) {
            e.printStackTrace();
            System.out.printf("Expected exception for TestCase 1_1: %s\n", e.getErrorCode());
        }
        return;

    }

    private void testCase2() throws TGException
    {
        System.out.println("Test Case 2: Update Node John's attribute.");
        john.setAttribute("age", 35);
        //john.setAttribute("nickname", "麗美"); //swapped the character
        john.setAttribute("nickname", "This is a long nickname"); //swapped the character
        conn.updateEntity(john);
        conn.commit(); //Should be successful.
    }
    
    private void testCase3() throws TGException
    {
        try {
            System.out.println("Test Case 3: Insert 2 nodes, and set a relation between them");
            smith = createNode(basicNodeType);
            smith.setAttribute("name", "smith"); //name is the primary key
            smith.setAttribute("age", 30);
            smith.setAttribute("nickname", "will");
            conn.insertEntity(smith);

            kelly = createNode(basicNodeType);
            kelly.setAttribute("name", "kelly"); //name is the primary key
            kelly.setAttribute("age", 28);
            kelly.setAttribute("nickname", "Ki");
            conn.insertEntity(kelly);

            brother = gof.createEdge(smith, kelly, TGEdge.DirectionType.Directed);
            brother.setAttribute("name", "Sister");
            conn.insertEntity(brother);

            conn.commit();
        }
        catch (TGException e) {
            e.printStackTrace();
        }

    }

    private void testCase3_1() throws TGException
    {
        System.out.println("Test Case 3_1: Get the Entity that we inserted");
        TGKey key = gof.createCompositeKey("basicnode");
        key.setAttribute("name", "smith");;
        TGEntity entity = conn.getEntity(key, TGQueryOption.DEFAULT_QUERY_OPTION);
        if (entity instanceof TGNode) {
            System.out.println("John's age is :" + entity.getAttribute("age").getAsInt());
        }

    }
    
    private void testCase4() throws TGException
    {
        try {
            System.out.println("Test Case 4: Add an edge between 2 existing nodes - In case between john and kelly");
            wife = gof.createEdge(kelly, john, TGEdge.DirectionType.Directed);
            wife.setAttribute("name", "wife");
            conn.insertEntity(wife);

            conn.commit();
        }
        catch (TGException e) {
            e.printStackTrace();
            wife = null;
        }
    }

    private void testCase4_0() throws TGException
    {
        System.out.println("Test Case 4_0: Get the Entity that we inserted");
        TGKey key = gof.createCompositeKey("basicnode");
        key.setAttribute("name", "kelly");;
        TGEntity entity = conn.getEntity(key, TGQueryOption.DEFAULT_QUERY_OPTION);
        if (entity instanceof TGNode) {
            System.out.println("Kelly's age is :" + entity.getAttribute("age").getAsInt());
        }
        kelly = TGNode.class.cast(entity);

    }

    private void testCase5() throws TGException
    {
        System.out.println("Test Case 5: Update an existing Edge");
        if (wife != null) {
            wife.setAttribute("name", "wife");
            wife.setAttribute("dom", "10/2/2016");  //Adding date of marriage
            conn.updateEntity(wife);
            conn.commit();
        }
    }

    private void testCase6() throws TGException
    {
        System.out.println("Test Case 6: Deleting Node 1");
        conn.deleteEntity(john);
        conn.commit();

    }

    private void testCase6_1() throws TGException
    {
        try {
            System.out.println("Test Case 6_1: Updating Node 1 Again - Should throw mismatch of ERA. or deleted");
            john.setAttribute("age", 40);
            john.setAttribute("nickname", "美麗"); //swapped the character
            conn.updateEntity(john);
            conn.commit();
        }
        catch (TGException e) {
            e.printStackTrace();
        }

    }

    private void testCase7() throws TGException
    {
        System.out.println("Test Case 7: Deleting Egde");
        conn.deleteEntity(brother);
        conn.commit();

    }


    private void parseArgs(String args[])
    {
        if (args.length == 0) return;
        for (int i=0; i < args.length; i++)
        {
            String arg = args[i];
            if (arg.equalsIgnoreCase("-url")) {
                url = args[++i];
            }
            else if (arg.equalsIgnoreCase("-passwd")) {
                passwd = args[++i];
            }
            else if (arg.equalsIgnoreCase("-loglevel")) {
                logLevel = TGLogger.TGLevel.valueOf(args[++i]);
            }
        }
    }
    

    public static void main(String[] args) {
        try {
            TransactionUnitTest tut = new TransactionUnitTest(args);
            tut.connect();
            tut.runTestCases();

        }
        catch (Exception e) {
            e.printStackTrace();
        }
    }

       
}
