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

package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.query.TGQueryOption;

import java.math.BigDecimal;
import java.util.Calendar;
import java.util.TimeZone;

public class MetadataTest {
    public String url = "tcp://admin@localhost:8228";
    public String passwd = "admin";
    TGGraphObjectFactory gof;
    TGGraphMetadata gmd;
    TGConnection conn;
    TGNodeType nodeAllAttrs, nodeOneAttr;

    public void connect() throws TGException {
        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);
        conn.connect();

        gof = conn.getGraphObjectFactory();
        gmd = conn.getGraphMetadata(true);

        nodeAllAttrs = gmd.getNodeType("nodeAllAttrs");
        if (nodeAllAttrs == null) throw new TGException("Node desc nodeAllAttrs not found");

        nodeOneAttr = gmd.getNodeType("nodeOneAttr");
        if (nodeOneAttr == null) throw new TGException("Node desc nodeOneAttr not found");

        System.out.println("Leaving connect()...");
    }

    public void test1() throws Exception {
        System.out.println("Begin test1");
        TGNode node = gof.createNode(nodeOneAttr);

        node.setAttribute("stringAttr", "StringKey");

        conn.insertEntity(node);
        conn.commit(); // FIXME this is where the test fails with result TGCatalogIndexNotLoaded.

        TGKey key = gof.createCompositeKey("nodeOneAttr");
        key.setAttribute("stringAttr", "StringKey");
        TGQueryOption option = TGQueryOption.createQueryOption();

        TGEntity entity = conn.getEntity(key, option);

        if (entity instanceof TGNode) {
            System.out.println("Found node. stringAttr = " + entity.getAttribute("stringAttr").getAsString());
        }

    }

    public void test2() throws Exception {
        System.out.println("Begin test2");
        TGNode node = gof.createNode(nodeAllAttrs);

        node.setAttribute("boolAttr", false);
        node.setAttribute("byteAttr", (byte) 0xba);
//        node.setAttribute("charAttr", '*'); // FIXME this line creates a malformed transaction.
        node.setAttribute("shortAttr", (short) 6385);
        node.setAttribute("intAttr", 73741825);
        node.setAttribute("longAttr", (long) 1342177281);
        node.setAttribute("floatAttr", (float) 2.23);
        node.setAttribute("doubleAttr", 2336.32424);
        node.setAttribute("numberAttr", new BigDecimal("234235732234235723590735124523.89813275891735070"));
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
        conn.commit(); // fixme this is where this test crashes, eventually causing a core dump.

        TGKey key = gof.createCompositeKey("nodeAllAttrs");
        key.setAttribute("stringAttr", "betterStringKey");
        TGQueryOption option = TGQueryOption.createQueryOption();

        TGEntity entity = conn.getEntity(key, option);

        if (entity instanceof TGNode) {
            System.out.println("Found node. stringAttr = " + entity.getAttribute("stringAttr").getAsString());
        }
    }

    public static void main(String[] args) throws Exception {
        try {
            MetadataTest test = new MetadataTest();
            test.connect();
            test.test1();
            test.test2();
        }
        catch (Exception e) {
            e.printStackTrace();
        }
    }
}
