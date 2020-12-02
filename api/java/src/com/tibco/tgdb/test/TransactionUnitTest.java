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
 *
 * File name : TransactionUnitTest.${EXT}
 * Created on: 10/02/2016
 * Created by: suresh
 * SVN Id: $Id: TransactionUnitTest.java 3998 2020-05-17 02:31:57Z vchung $
 */

package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTransactionException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.query.TGResultSet;

import java.math.BigDecimal;
import java.util.Calendar;
import java.util.GregorianCalendar;
import java.util.TimeZone;

import static java.lang.System.out;

public class TransactionUnitTest {

    public String url = "tcp://scott@localhost:8222/{dbName=demodb}";
    //public String url = "tcp://scott@10.98.201.111:8228/{connectTimeout=30}";
    //public String url = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
    //public String url = "tcp://scott@localhost6:8228";
    public String passwd = "scott";
    public TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    TGGraphObjectFactory gof;
    TGGraphMetadata gmd;
    TGConnection conn;
    TGNodeType basicNodeType, rateNodeType, testNodeType, nodeAllAttrs;
    TGNode john, smith, kelly;
    TGEdge brother, wife;
    
    
    public TransactionUnitTest(String args[])
    {
        TGLogger logger = TGLogManager.getInstance().getLogger();
        logger.setLevel(logLevel);
        parseArgs(args);
    }
    
    public void connect() throws TGException {
        out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        gof = conn.getGraphObjectFactory();
        gmd = conn.getGraphMetadata(true);
        basicNodeType = gmd.getNodeType("basicnode");
        if (basicNodeType == null) throw new TGException("Node desc basicnode not found");

        rateNodeType = gmd.getNodeType("ratenode");
        if (rateNodeType == null) throw new TGException("Node desc ratenode not found");

        testNodeType = gmd.getNodeType("testnode");
        if (testNodeType == null) throw new TGException("Node desc testnode not found");
    }

    void disconnect() {
        conn.disconnect();
    }
    
    public void runTestCases() throws TGException, Exception
    {
        //jira_testCase_157();
        //jira_testCase_182();
        //testCase0();
        //testCase1_1();

        testCase1();
//        testCase1_1();

        testCase1_2();
        testCase2();

        testCase3();
        testCase3_1();
        testCase4();
        testCase4_0();
        testCase4_2();
        testCase5();
        testCase5_1();
        testCase6();
        testCase6_1();
        testCase1();
        //testCase7();
        out.println("Done all the cases...");


    }

    TGNode createNode(TGNodeType nodeType) {
        if (nodeType != null) {
            return gof.createNode(nodeType);
        } else {
            return gof.createNode();
        }
    }

    public void jira_testCase_157() throws TGException {

        out.println("Begin test2");
        nodeAllAttrs = gmd.getNodeType("nodeAllAttrs");
        if (nodeAllAttrs == null) throw new TGException("Node desc nodeAllAttrs not found");

        TGNode node = gof.createNode(nodeAllAttrs);

        node.setAttribute("boolAttr", false);
        node.setAttribute("byteAttr", (byte) 0xba);
        node.setAttribute("charAttr", '*');
        node.setAttribute("shortAttr", (short) 6385);
        node.setAttribute("numberAttr", new BigDecimal("907323.070"));
        node.setAttribute("intAttr", 73741825);
        node.setAttribute("longAttr", (long) 1342177281);
        node.setAttribute("floatAttr", (float) 2.23);
        node.setAttribute("doubleAttr", 2336.32424);
        node.setAttribute("stringAttr", "betterStringKey");
        node.setAttribute("dateAttr", new Calendar
                .Builder()
                .setDate(2016, 10, 31)
                .build());
        node.setAttribute("timeAttr", new Calendar
                .Builder()
                .setTimeOfDay(21, 32, 12, 845)
                .setTimeZone(TimeZone.getDefault())
                .build());
        node.setAttribute("timestampAttr", new Calendar
                .Builder()
                .setDate(2016, 10, 25)
                .setTimeOfDay(8,9,30,999)
                .build());


        conn.insertEntity(node);
        out.println("before commit");
        conn.commit();
        out.println("after commit");

        TGKey key = gof.createCompositeKey("nodeAllAttrs");
        key.setAttribute("stringAttr", "betterStringKey");
        TGQueryOption option = TGQueryOption.createQueryOption();

        TGEntity entity = conn.getEntity(key, option);

        // FIXME the commented lines below indicate attrdesc that do not have a get() function
        if (entity instanceof TGNode) {
            out.println("Found node.");
            out.println("stringAttr = " + entity.getAttribute("stringAttr").getAsString());
            out.println("boolAttr = " + entity.getAttribute("boolAttr").getAsBoolean());
            out.println("charAttr = " + entity.getAttribute("charAttr").getAsChar());
            out.println("shortAttr = " + entity.getAttribute("shortAttr").getAsShort());
            out.println("intAttr = " + entity.getAttribute("intAttr").getAsInt());
            out.println("longAttr = " + entity.getAttribute("longAttr").getAsLong());
            out.println("floatAttr = " + entity.getAttribute("floatAttr").getAsFloat());
            out.println("doubleAttr = " + entity.getAttribute("doubleAttr").getAsDouble());
            //    out.println("numberAttr = " + entity.getAttribute("numberAttr").getAsString());
            //    out.println("dateAttr = " + entity.getAttribute("dateAttr").getAsString());
            //    out.println("timeAttr = " + entity.getAttribute("timeAttr").getAsLong());
            //    out.println("timestampAttr = " + entity.getAttribute("timestampAttr").getAsLong());
        }
    }

    private void jira_testCase_182() throws TGException
    {
        TGNode basic1 = gof.createNode(basicNodeType);
        TGNode basic2 = gof.createNode(basicNodeType);
        TGEdge edge1;

        edge1 = gof.createEdge(basic1, basic2, TGEdge.DirectionType.UnDirected);
        edge1.setAttribute("ratedate", new Calendar
                .Builder()
                .setDate(2016, 12, 1)
                .set(Calendar.ERA, GregorianCalendar.BC)
                .build());
        basic1.setAttribute("name", "Mike");
        basic2.setAttribute("name", "Kevin");
        conn.insertEntity(basic1);
        conn.insertEntity(basic2);
        conn.insertEntity(edge1);

        conn.commit();
        out.println("Entities created");

        conn.getGraphMetadata(true);
        TGKey key = gof.createCompositeKey("basicnode");

        key.setAttribute("name", "Mike");
        TGEntity entity = conn.getEntity(key, null);
        if (entity != null) {
            out.println("Name = " + entity.getAttribute("name").getValue());
        }
    }


    String keyName1 = new String ("Gabe");
    String keyName2 = new String ("Georgia");

    private void testCase0() throws TGException
    {
        try {
            out.println("Test Case 0: Insert Simple Node(John) of basicnode with a few properties");

            TGNode j1 = createNode(basicNodeType);

            j1.setAttribute("name", keyName1); //name is the primary key
            j1.setAttribute("age", 30);
            //j1.setAttribute("nickname", "美麗");
            j1.setAttribute("createtm", new Calendar
                    .Builder()
                    .setDate(2016, 10, 25)
                    .setTimeOfDay(15, 9, 30, 999)
                    .build());

            j1.setAttribute("networth", new BigDecimal("2378989.567"));
            j1.setAttribute("flag", 'D');
            j1.setAttribute("desc", "Hi TIBCO Team!\n" +
                    "\n" +
                    "The second stop on the TIBCO NOW Global Tour is just days away. We saw extreme value from Singapore and the excitement and now we do it again in Berlin this time with 545 registered attendees! We have reached and exceeded our target and will be closing registration before we run into any capacity issues. We are very excited about this event and to see what is coming with some game changing product updates shown for the first time at TIBCO NOW Berlin. (There will be a Sharpen the Saw on this Friday)\n" +
                    "\n");

            conn.insertEntity(j1);

            TGNode j2 = createNode(basicNodeType);

            j2.setAttribute("name", keyName2); //name is the primary key
            j2.setAttribute("age", 30);
            //j2.setAttribute("nickname", "美麗");
            j2.setAttribute("createtm", new Calendar
                    .Builder()
                    .setDate(2016, 10, 25)
                    .setTimeOfDay(15, 9, 30, 999)
                    .build());

            j2.setAttribute("networth", new BigDecimal("2378989.567"));
            j2.setAttribute("flag", 'D');
            j2.setAttribute("desc", "Hi TIBCO Team!\n" +
                    "\n" +
                    "The second stop on the TIBCO NOW Global Tour is just days away. We saw extreme value from Singapore and the excitement and now we do it again in Berlin this time with 545 registered attendees! We have reached and exceeded our target and will be closing registration before we run into any capacity issues. We are very excited about this event and to see what is coming with some game changing product updates shown for the first time at TIBCO NOW Berlin. (There will be a Sharpen the Saw on this Friday)\n" +
                    "\n");

            conn.insertEntity(j2);

            for (int i=0;i<1000;i++) {
                TGEdge edge = gof.createEdge(j1, j2, TGEdge.DirectionType.Directed);
                edge.setAttribute("name", "spouse");
                edge.setAttribute("desc", "This is test...");
                conn.insertEntity(edge);
            }
            conn.commit();
            john = j1;


        }
        catch (TGException e) {
            e.printStackTrace();
        }
        out.printf("Test Case 0 end\n\n");
    }

    private void testCase1() throws TGException
    {
        try {
            out.println("Test Case 1: Insert Simple Node(John) of basicnode with a few properties");

            TGNode node = createNode(basicNodeType);

            node.setAttribute("name", "john"); //name is the primary key
            node.setAttribute("age", 40);
            node.setAttribute("nickname", "美麗");
            node.setAttribute("createtm", new Calendar
                    .Builder()
                    .setDate(2016, 10, 25)
                    .setTimeOfDay(15, 9, 30, 999)
                    .build());

            node.setAttribute("networth", new BigDecimal("2378989.567"));
            node.setAttribute("flag", 'D');
            node.setAttribute("ssn", "123-456-7890");

            conn.insertEntity(node);
            conn.commit();
            john = node;
        }
        catch (TGException e) {
            e.printStackTrace();
        }
        out.printf("Test Case 1 end\n\n");
    }

    private void testCase1_1() throws TGException
    {
        out.println("Test Case 1_0: Get the Entity that we inserted");
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
            out.println("John's age is :" + entity.getAttribute("age").getAsInt());
            out.println("John's createtm:" + entity.getAttribute("createtm").getValue().toString());
            out.println("John's networth:" + entity.getAttribute("networth").getValue().toString());
            out.println("John's ssn:" + entity.getAttribute("ssn").getValue().toString());
        }
        john = TGNode.class.cast(entity);

    }

    private void testCase1_2() throws TGException
    {
        try {
            out.println("Test Case 1_2: Again insert John. This should raise Unique Key Constraint violation.");
            TGNode node = createNode(basicNodeType);
            node.setAttribute("name", "john");
            node.setAttribute("age", 30);
            node.setAttribute("nickname", "美麗");
            conn.insertEntity(node);
            conn.commit();
        }
        catch (TGException e) {
            e.printStackTrace();
            out.printf("Expected exception for TestCase 1_2: %s\n", e.getExceptionType());
        }
        out.printf("Test Case 1_2 end\n\n");
        return;

    }

    private void testCase2() throws TGException
    {
        out.println("Test Case 2: Update Node John's attribute.");
        john.setAttribute("age", 35);
        //john.setAttribute("nickname", "麗美"); //swapped the character
        john.setAttribute("nickname", "This is a long nickname"); //swapped the character
        conn.updateEntity(john);
        conn.commit(); //Should be successful.
        out.printf("Test Case 2 end\n\n");
    }
    
    private void testCase3() throws TGException
    {
        try {
            out.println("Test Case 3: Insert 2 nodes, and set a relation between them");
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
        out.printf("Test Case 3 end\n\n");
    }

    private void testCase3_1() throws TGException
    {
        out.println("Test Case 3_1: Get the Entity that we inserted");
        TGKey key = gof.createCompositeKey("basicnode");
        key.setAttribute("name", "smith");;
        TGEntity entity = conn.getEntity(key, TGQueryOption.DEFAULT_QUERY_OPTION);
        if (entity instanceof TGNode) {
            out.println("John's age is :" + entity.getAttribute("age").getAsInt());
        }
        out.printf("Test Case 3_1 end\n\n");
    }
    
    private void testCase4() throws TGException
    {
        try {
            out.println("Test Case 4: Add an edge between 2 existing nodes - In case between john and kelly");
            wife = gof.createEdge(kelly, john, TGEdge.DirectionType.Directed);
            wife.setAttribute("name", "wife");
            conn.insertEntity(wife);

            conn.commit();
        }
        catch (TGException e) {
            e.printStackTrace();
            wife = null;
        }
        out.printf("Test Case 4 end\n\n");
    }

    private void testCase4_0() throws TGException
    {
        out.println("Test Case 4_0: Get the Entity that we inserted");
        TGKey key = gof.createCompositeKey("basicnode");
        key.setAttribute("name", "kelly");;
        TGEntity entity = conn.getEntity(key, TGQueryOption.DEFAULT_QUERY_OPTION);
        if (entity instanceof TGNode) {
            out.println("Kelly's age is :" + entity.getAttribute("age").getAsInt());
        }
        kelly = TGNode.class.cast(entity);
        out.printf("Test Case 4_0 end\n\n");
    }

    private void testCase5() throws TGException
    {
        out.println("Test Case 5: Update an existing Edge");
        if (wife != null) {
            wife.setAttribute("name", "Ex-wife");
            //wife.setAttribute("dom", "10/2/2016");  //Adding date of marriage
            conn.updateEntity(wife);
            conn.commit();
        }
        out.printf("Test Case 5 end\n\n");
    }

    //This test case needs to be self contained because once a new attribute descriptor is created locally, all
    //the subsequent transactions will failed because attribute descriptor cannot be created on the fly on the server
    //side
    private void testCase5_1() throws TGException
    {
        out.println("Test Case 5_1: Update a node with undefined attribute");
        TGConnection c = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);
        c.connect();

        try {
            TGGraphObjectFactory g = c.getGraphObjectFactory();
            TGKey key = g.createCompositeKey("basicnode");
            key.setAttribute("name", "smith");
            TGQueryOption option = TGQueryOption.createQueryOption();

            TGEntity entity = c.getEntity(key, option);
            if (entity != null) {
                entity.setAttribute("nickname", "joe");
                entity.setAttribute("dom", "10/2/2016");  //Adding date of marriage
                c.updateEntity(entity);
                c.commit();
            }

        } catch (TGTransactionException te) {
            out.printf("Expected exception for TestCase 5_1: %s(%s)\n", te, te.getExceptionType());
        }
        c.disconnect();
        out.printf("Test Case 5_1 end\n\n");
    }

    private void testCase6() throws TGException
    {
        out.println("Test Case 6: Deleting Node 1");
        conn.deleteEntity(john);
        conn.commit();
        out.printf("Test Case 6 end\n\n");
    }

    private void testCase6_1() throws TGException
    {
        try {
            out.println("Test Case 6_1: Updating Node 1 Again - Should throw mismatch of ERA. or deleted");
            john.setAttribute("age", 40);
            john.setAttribute("nickname", "美麗"); //swapped the character
            conn.updateEntity(john);
            conn.commit();
        }
        catch (TGException e) {
            e.printStackTrace();
        }
        out.printf("Test Case 6_1 end\n\n");
    }

    private void testCase7() throws TGException
    {
        out.println("Test Case 7: Deleting Egde");
        conn.deleteEntity(brother);
        conn.commit();

    }

    private void testCase4_2() throws TGException
    {
        out.printf("Test Case 4_2 start\n");
        String startFilter = "@nodetype = 'basicnode' and name = 'smith';";
        String traverserFilter = "@degree = 1;";
        String endFilter = "@nodetype = 'basicnode' and name = 'kelly';";
        TGResultSet result = conn.executeQuery(startFilter, null, traverserFilter, endFilter, TGQueryOption.DEFAULT_QUERY_OPTION);
        while (result.hasNext()) {
            out.printf("result.next = %s\n", ((TGEntity) result.next()).getAttribute("name").getAsString());
        }
        out.printf("Test Case 4_2 end\n\n");
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
        TransactionUnitTest tut = new TransactionUnitTest(args);
        try {

            tut.connect();
            tut.runTestCases();

        }
        catch (Exception e) {
            e.printStackTrace();
        }
        finally {
            tut.disconnect();
        }
    }

       
}
