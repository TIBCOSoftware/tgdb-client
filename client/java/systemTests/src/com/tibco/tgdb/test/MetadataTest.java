package com.tibco.tgdb.test;


/**
 * Copyright (c) 2016 TIBCO Software Inc.

 * All rights reserved.
 * <p/>
 * File name : MetadataTest.${EXT}
 * Created on: 11/7/16
 * Created by: Katie
 * <p/>
 * SVN Id: $Id$
 */



import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.query.TGQueryOption;

import java.math.BigDecimal;
import java.util.Calendar;
import java.util.TimeZone;

import static java.lang.Thread.sleep;

public class MetadataTest {
    public String url = "tcp://admin@localhost:8223";
    public String passwd = "admin";
    TGGraphObjectFactory gof;
    TGGraphMetadata gmd;
    TGConnection conn;
    TGNodeType nodeAllAttrs;

    public void connect() throws TGException {
        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);
        conn.connect();

        gof = conn.getGraphObjectFactory();
        gmd = conn.getGraphMetadata(true);

        nodeAllAttrs = gmd.getNodeType("nodeAllAttrs");
        if (nodeAllAttrs == null) throw new TGException("Node type nodeAllAttrs not found");

        System.out.println("Leaving connect()...");
    }

    public void disconnect() {

        conn.disconnect();
    }

    //test a node with all attrdescs
    public void test1() throws Exception {
        System.out.println("Begin test1...");
        TGNode node = gof.createNode(nodeAllAttrs);

        node.setAttribute("numberAttr", new BigDecimal("907323.070"));
        node.setAttribute("boolAttr", false);
        node.setAttribute("byteAttr", (byte) 0xba);
        node.setAttribute("charAttr", '*');
        node.setAttribute("shortAttr", (short) 6385);
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
        System.out.println("before commit");
        conn.commit();
        System.out.println("after commit");

        TGKey key = gof.createCompositeKey("nodeAllAttrs");
        key.setAttribute("stringAttr", "betterStringKey");
        TGQueryOption option = TGQueryOption.createQueryOption();

        TGEntity entity = conn.getEntity(key, option);

        // FIXME the commented lines below indicate attrdesc that do not have a get() function (jira ticket tgdb-157)
        if (entity instanceof TGNode) {
            System.out.println("Found node.");
            System.out.println("stringAttr = " + entity.getAttribute("stringAttr").getAsString());
            System.out.println("boolAttr = " + entity.getAttribute("boolAttr").getAsBoolean());
            System.out.println("charAttr = " + entity.getAttribute("charAttr").getAsChar());
            System.out.println("shortAttr = " + entity.getAttribute("shortAttr").getAsShort());
            System.out.println("intAttr = " + entity.getAttribute("intAttr").getAsInt());
            System.out.println("longAttr = " + entity.getAttribute("longAttr").getAsLong());
            System.out.println("floatAttr = " + entity.getAttribute("floatAttr").getAsFloat());
            System.out.println("doubleAttr = " + entity.getAttribute("doubleAttr").getAsDouble());
        //    System.out.println("numberAttr = " + entity.getAttribute("numberAttr").getAsString());
        //    System.out.println("dateAttr = " + entity.getAttribute("dateAttr").getAsString());
        //    System.out.println("timeAttr = " + entity.getAttribute("timeAttr").getAsLong());
        //    System.out.println("timestampAttr = " + entity.getAttribute("timestampAttr").getAsLong());
        }

        System.out.println("Leaving test1.");
    }

    // test creating indices with different attributes
    public void test2() throws Exception {
        System.out.println("Begin test2...");
        TGNode node = gof.createNode(nodeAllAttrs);

        node.setAttribute("numberAttr", new BigDecimal("907323.070"));
        
        node.setAttribute("boolAttr", false);
        node.setAttribute("byteAttr", (byte) 0xba);
        node.setAttribute("charAttr", '*');
        node.setAttribute("shortAttr", (short) 6385);
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
        conn.commit();
        System.out.println("After commit..:"+node);
        System.out.printf("print nodeAllAtributes...%s:\n",node.getAttribute("stringAttr").getAsString());
        //houseMember.getAttribute("memberName").getAsString()
        TGKey key = gof.createCompositeKey("nodeAllAttrs");
        key.setAttribute("boolAttr", false);
        TGQueryOption option = TGQueryOption.createQueryOption();
        TGEntity entity = conn.getEntity(key, option);
        System.out.println("entity...:"+entity);
        
        
        

    }


    public static void main(String[] args) throws Exception {
        try {
            MetadataTest test = new MetadataTest();
            test.connect();
           // test.test1();
            test.test2();
            test.disconnect();
        }
        catch (Exception e) {
            e.printStackTrace();
        }

    }
}
