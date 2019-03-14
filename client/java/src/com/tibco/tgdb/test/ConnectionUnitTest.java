/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : ConnectionUnitTest.${EXT}
 * Created on: 3/15/18
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

public class ConnectionUnitTest {

    //public String url = "ssl://scott@192.168.1.15:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
    public String url = "ssl://scott@10.108.16.93:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
    //public String url = "ssl://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
    //public String url = "tcp://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
    //public String url = "tcp://scott@10.98.201.111:8228/{connectTimeout=30}";
    //public String url = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
    //public String url = "tcp://scott@localhost:8222";
    public String passwd = "scott";
    public TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    TGGraphObjectFactory gof;
    TGGraphMetadata gmd;
    TGConnection conn;
    TGNodeType basicNodeType, rateNodeType, testNodeType, nodeAllAttrs;
    TGNode john, smith, kelly;
    TGEdge brother, wife;


    public ConnectionUnitTest(String args[])
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

//        TGBlob blob = gof.createBlob();
//        TGClob clob = gof.createClob();
//        byte[] buf = new byte[100000];
//        blob.setBytes(buf);
//        char[] cbuf = new char[10000];
//        clob.setChars(cbuf);

        gmd = conn.getGraphMetadata(true);
        basicNodeType = gmd.getNodeType("basicnode");
        if (basicNodeType == null) throw new TGException("Node desc basicnode not found");

        rateNodeType = gmd.getNodeType("ratenode");
        if (rateNodeType == null) throw new TGException("Node desc ratenode not found");

        testNodeType = gmd.getNodeType("testnode");
        if (testNodeType == null) throw new TGException("Node desc testnode not found");


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

    private static void case1(String[] args) {
        ConnectionUnitTest cut = new ConnectionUnitTest(args);
        try {

            cut.connect();


        }
        catch (Exception e) {
            e.printStackTrace();
        }
        finally {
            cut.conn.disconnect();
        }
    }

    private static void case2(String[] args) {
        ConnectionUnitTest cut = new ConnectionUnitTest(args);
        try {

            cut.connect();
            Thread.sleep(Integer.MAX_VALUE);
        }
        catch (Exception e) {
            e.printStackTrace();
        }
        finally {
            cut.conn.disconnect();
        }
    }



    public static void main(String[] args) {
        //case1(args);
        case2(args);
    }
}
