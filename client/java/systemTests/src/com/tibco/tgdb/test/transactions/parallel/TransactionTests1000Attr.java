package com.tibco.tgdb.test.transactions.parallel;
import java.util.UUID;

import java.io.File;
import java.io.IOException;
import java.math.BigDecimal;
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.TimeZone;
import java.util.Timer;

import org.testng.Assert;
import org.testng.annotations.AfterGroups;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeGroups;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.DataProvider;
import org.testng.annotations.Test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTransactionException;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.utils.ClasspathResource;
import com.tibco.tgdb.test.utils.PipedData;

import bsh.EvalError;

/**
 * Copyright 2018 TIBCO Software Inc. All rights reserved.
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
 */

public class TransactionTests1000Attr {
	private static TGServer tgServer;
	private static String tgUrl;
	private static String tgUser = "scott";
	private static String tgPwd = "scott";
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");
	int cnt = 0;
	
	/**
	 * Init TG server before test suite
	 * @throws Exception
	 */
	@BeforeSuite
	public void initServer() throws Exception  {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/initdb1000Attr.conf", tgWorkingDir + "/initdb1000Attr.conf");
		tgServer = new TGServer(tgHome);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 300000);
		}
		catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
	}
	
	
	/**
	 * Start TG server with maximum 10 connections
    */
	
	@BeforeGroups("AttributeTest")
	public void startMaxConnect10Server() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/tgdb10.conf", tgWorkingDir + "/tgdb10.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	/**
	 * Stop TG server started with maximum 10 connections
    */
	
	@AfterGroups("AttributeTest")
	public void stopMaxConnect10Server() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".maxconnect10");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}

		
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
	/**
	 * Test to verify 10 threads connecting in parallel,tgserver configured with max 10 connections
    */
	
	@Test( groups = "AttributeTest",description ="Test to verify creating 1000 attributes by each thread,while 10 threads connecting in parallel")
	public void testtoverify1000attributes() throws Exception {
		     ConnectServer(101);
		
	        }
	
	/**
	 * Test to verify 10 threads connecting in parallel,tgserver configured with max 10 connections
    */
	
	@Test( groups = "AttributeTest",description ="Test to verify exception creating 3000 attributes by each thread,while 10 threads connecting in parallel")
	public void testtoverify2000attributes() throws Exception {
		     

				try {
					ConnectServer(301);
				}
				catch(Exception e) { 
					
					//System.out.println(e.toString());
					if(!e.toString().contains("TGTransactionResourceExceededException"))
					Assert.fail("Expected a TGTransactionResourceExceededException upon connection but got a "+ e.getClass().getName() +" instead");
			
				}
				finally { 
					
				}
	        }
	

	/**
	 * Start TG server with custom configuration and connects to it.
	 * @throws Exception
	 */
	
	public void ConnectServer(int attrcount) throws Exception {
		String url = "tcp://127.0.0.1:8222";
		String user = "scott";
		String pwd = "scott";
		
		 
		TGConnection conn = null;
		try {
			conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		 UUID uuid = UUID.randomUUID();
	    String randomUUIDString = uuid.toString();
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType basicnode = gmd.getNodeType("nodeAllAttrs");
		if (basicnode == null)
			throw new Exception("Node type not found");
		
		TGNode basic1 = gof.createNode(basicnode);
		Object[][] data = this.getNodeData();

		
        for(int i=1;i<attrcount;i++)
		 {

		     // basic1.setAttribute("key", i);
		      
			  basic1.setAttribute("boolAttr"+i, data[1][1]);
			  basic1.setAttribute("intAttr"+i,data[1][2]);
			  basic1.setAttribute("charAttr"+i, data[1][3]);
			  basic1.setAttribute("byteAttr"+i, data[1][4]);
			  basic1.setAttribute("longAttr"+i, data[1][5]);
			  basic1.setAttribute("stringAttr"+i,data[1][6]);	 
			  basic1.setAttribute("shortAttr"+i, data[1][7]);
			  basic1.setAttribute("floatAttr"+i,data[1][8]);
			  basic1.setAttribute("doubleAttr"+i, data[1][9]);
			 // basic1.setAttribute("dateAttr"+i, data[1][10]);
			  basic1.setAttribute("number"+i, data[1][13]);		
	   }
	   conn.insertEntity(basic1);
	   conn.commit();
		System.out.println("Entity 1 created");
		
	} finally {
		if (conn != null)
			conn.disconnect();
	}
	}


	
	@DataProvider(name = "NodeData")
	public Object[][] getNodeData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/node.data"));
		return data;
	}


	/**
	 * Test to verify 100 threads connecting in parallel,tgserver configured with max 100 connections
	*/
	@Test( 	groups = "MaxConnect100",
			description ="Test to verify 100 threads connecting in parallel,tgserver configured with max 100 connections",
			invocationCount = 30,
			threadPoolSize = 30,
			timeOut=600000)
	public void testConnectwithMaxConnect100() throws Exception {
		connectAndInsertNode(1000);
	}


	private void connectAndInsertNode(int i) {
		// TODO Auto-generated method stub
		
	}
}
