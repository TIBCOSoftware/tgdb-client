package com.tibco.tgdb.test.transactions.parallel;

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

import java.util.UUID;

import java.io.File;
import java.io.IOException;
import java.math.BigDecimal;
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.List;
import java.util.TimeZone;

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


public class TransactionTests {
	private static TGServer tgServer;
	private static String tgUrl = "tcp://127.0.0.1:8222"; // /{connectTimeout=1000}";
	private static String tgUser = "scott";
	private static String tgPwd = "scott";
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");
	
	static int connectionSuccessCount = 0;
	static int connectionFailedCount = 0;
	
	
	/**
	 * Init TG server before test suite
	 * @throws Exception
	 */
	@BeforeSuite
	public void initServer() throws Exception  {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/initdb.conf", tgWorkingDir + "/initdb.conf");
		tgServer = new TGServer(tgHome);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 300000);
			System.out.println(tgServer.getBanner());
		}
		catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
	}
	
	
	/**
	 * Start TG server with maximum 10 connections
    */
	@BeforeGroups("MaxConnect10")
	public void startMaxConnect10Server() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/tgdb10.conf", tgWorkingDir + "/tgdb10.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	/**
	 * Stop TG server started with maximum 10 connections
    */
	@AfterGroups("MaxConnect10")
	public void stopMaxConnect10Server() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".maxconnect10");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}

	/**
	 * Start TG server with maximum 100 connections
    */
	@BeforeGroups("MaxConnect100")
	public void startMaxConnect100Server() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/tgdb100.conf", tgWorkingDir + "/tgdb100.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}

	/**
	 * Stop TG server started with maximum 100 connections
    */
	@AfterGroups("MaxConnect100")
	public void stopMaxConnect100Server() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".maxconnect100");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	
	/**
	 * Start TG server with maximum 20 connections
    */
	@BeforeGroups("MaxConnect20")
	public void startMaxConnect20Server() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/tgdb10.conf", tgWorkingDir + "/tgdb10.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	/**
	 * Stop TG server started with maximum 90 connections
    */
	@AfterGroups("MaxConnect20")
	public void stopMaxConnect20Server() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".maxconnect90");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	/**
	 * Start TG server with maximum 500 connections
    */
	@BeforeGroups("MaxConnect500")
	public void startMaxConnect500Server() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/tgdb500.conf", tgWorkingDir + "/tgdb500.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	
	/**
	 * Stop TG server started with maximum 500 connections
    */
	@AfterGroups("MaxConnect500")
	public void stopMaxConnec500Server() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".maxconnect500");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	/**
	 * Start TG server with maximum 1000 connections
    */
	@BeforeGroups("MaxConnect1000")
	public void startMaxConnect1000Server() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/tgdb1000.conf", tgWorkingDir + "/tgdb1000.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	
	/**
	 * Stop TG server started with maximum 1000 connections
    */
	@AfterGroups("MaxConnect1000")
	public void stopMaxConnec1000Server() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".maxconnect1000");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
	/**
	 * Connect with 10 parallel threads. Each thread insert one node
    */
	@Test( 	groups = "MaxConnect10",
			description ="Connect with 10 parallel threads. Each thread insert one node",
			invocationCount = 10,
			threadPoolSize = 10,
			timeOut=120000)
	public void testConnectWith10Threads() throws Exception {
		connectAndInsertNode(1000);
	}	
	
	/**
	 * Connect with 100 parallel threads. Each thread insert one node
    */
	@Test( 	groups = "MaxConnect100",
			description ="Connect with 100 parallel threads. Each thread insert one node",
			invocationCount = 100,
			threadPoolSize = 100,
			timeOut=120000)
	public void testConnectWith100Threads() throws Exception {
		connectAndInsertNode(1000);
	}
	
	/**
	 * Connect with 500 parallel threads. Each thread insert one node
    */
	@Test( 	groups = "MaxConnect500",
			description ="Connect with 500 parallel threads. Each thread insert one node",
			invocationCount = 500,
			threadPoolSize = 500,
			timeOut=180000)
	public void testConnectWith500Threads() throws Exception {
		connectAndInsertNode(1000);
	}
	
	/**
	 * Connect with 1000 parallel threads. Each thread insert one node
    */
	@Test( 	groups = "MaxConnect1000",
			description ="Connect with 1000 parallel threads. Each thread insert one node",
			invocationCount = 1000,
			threadPoolSize = 1000,
			timeOut=180000)
	public void testConnectWith1000Threads() throws Exception {
		connectAndInsertNode(3000);
	}
	
	/**
	 * Connect with 20 parallel threads while the max connections is 10. Expect 10 TGException for Max Exceeded Connections
    */
	@Test( 	groups = "MaxConnect20",
			description ="Connect with 20 parallel threads while the max connections is 10. Expect 10 TGException for Max Exceeded Connections",
			invocationCount = 20,
			threadPoolSize = 20,
			timeOut=120000)
	public void testConnectWith20Threads() throws Exception {
		int expectedConnectionFailure = 10;
		int expectedConnectionSuccess = 10;
		
		try {
		     connectAndInsertNode(11000); // hold on to connection for 11 sec. It seems connectTimeout=10 sec and configuring it does not work.
		     synchronized (this) {
		    	 connectionSuccessCount  ++;
		     }
		     Assert.assertTrue(connectionSuccessCount <= expectedConnectionSuccess, "The number of successful connection is higher than " + expectedConnectionSuccess + " -");
		}
		catch(TGException tge) { // Expected 10 of this
			//System.out.println("TGException caught : "  + tge.getMessage());
			synchronized(this) {
				connectionFailedCount++;
			}
			// Assert that we get the correct TGException message (Exceeded max connections)
			Assert.assertEquals("TGException: "+tge.getMessage(), "TGException: Exceeded max connections of " + expectedConnectionFailure + " for " + tgServer.getNetListeners()[1].getName() + " listener");
			// Assert that nb failed connections is lower or equal to 10
			Assert.assertTrue(connectionFailedCount <= expectedConnectionFailure, "The number of ExceededMaxConnection Exception is higher than " + expectedConnectionFailure + " -");
		}
		finally { 
			;
		}
	}

	/**
	 * Connect to server and insert a node
	 * @throws Exception
	 */
	public void connectAndInsertNode(long sleep) throws Exception {
		
		TGConnection conn = null;
		try {
			conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
			conn.connect();
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null) {
				throw new org.testng.TestException("TG object factory is null");
			}
			UUID uuid = UUID.randomUUID();
			String randomUUIDString = uuid.toString();
			TGGraphMetadata gmd = conn.getGraphMetadata(true);
			TGNodeType basicnode = gmd.getNodeType("basicnode");
			if (basicnode == null)
				throw new Exception("Node type not found");
		
			TGNode basic1 = gof.createNode(basicnode);
			//System.out.println("Unique Id : "+randomUUIDString);
			basic1.setAttribute("name", randomUUIDString);
			basic1.setAttribute("networth", new java.math.BigDecimal("1.5E+6"));
			basic1.setAttribute("age", 73);
			basic1.setAttribute("createtm", Calendar.getInstance());
			conn.insertEntity(basic1);
			conn.commit();
			//Adding this sleep to make sure all connections are consumed in parallel
			Thread.sleep(sleep);
			//System.out.println("Entity created");	
		} 
		finally {
			if (conn != null) {
				conn.disconnect();
				//System.out.println("Disconnected !");
			}
		}
	}
}
