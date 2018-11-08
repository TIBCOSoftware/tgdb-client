package com.tibco.tgdb.test.connect;

import org.testng.annotations.Test;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.DataProvider;
import org.testng.annotations.BeforeGroups;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.AfterGroups;
import org.testng.Assert;

import static org.testng.Assert.fail;

import java.io.File;
import java.io.IOException;
import java.lang.reflect.Method;
import java.net.Inet4Address;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.connection.TGConnectionPool;
import com.tibco.tgdb.exception.TGAuthenticationException;
import com.tibco.tgdb.exception.TGChannelDisconnectedException;
import com.tibco.tgdb.exception.TGConnectionTimeoutException;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGGraphObjectFactory;

import com.tibco.tgdb.test.lib.*;
import com.tibco.tgdb.test.utils.*;

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

public class ConnectionTests {
	
	private static TGServer tgServer;
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");
	
	static String tgUser = "scott";
	static String tgPwd = "scott";
	
	static TGConnection tgConn = null;
	static TGConnectionPool tgPool = null;

	static int nbPoolTimeout = 0;
	static int nbPoolSuccess = 0;
	
	@BeforeSuite
	public void initServer() throws Exception  {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/Initdb.conf", tgWorkingDir + "/Initdb.conf");
		tgServer = new TGServer(tgHome);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 150000);
			System.out.println(tgServer.getBanner());
		}
		catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
	}
	
	@BeforeGroups("ipv4Grp")
	public void startIPv4Server() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/IPv4.conf", tgWorkingDir + "/IPv4.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	@AfterGroups("ipv4Grp")
	public void stopIPv4Server() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".ipv4");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
		
	}

	@BeforeGroups("ipv6Grp")
	public void startIPv6Server() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/IPv6.conf", tgWorkingDir + "/IPv6.conf");

		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	@AfterGroups("ipv6Grp")
	public void stopIPv6Server() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".ipv6");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}

	@BeforeGroups("authGrp")
	public void startAuthServer() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/Auth.conf", tgWorkingDir + "/Auth.conf");

		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	@AfterGroups("authGrp")
	public void stopAuthServer() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".auth");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	@AfterGroups("maxConnectGrp")
	public void stopMaxConnectionServer() throws Exception {
		
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".maxConnect");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	@BeforeGroups("massConnectGrp")
	public void startMassConnectionServer() throws Exception {
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/MassConnect.conf", tgWorkingDir + "/MassConnect.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}
	
	@AfterGroups("massConnectGrp")
	public void stopMassConnectionServer() throws Exception {
		
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".massConnect");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	@BeforeGroups("connectionPoolGrp")
	public void startConnectionPoolServer() throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/ConnectionPool.conf", tgWorkingDir + "/ConnectionPool.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
		
		String url = "tcp://"+ tgServer.getNetListeners()[0].getHost() + ":" + tgServer.getNetListeners()[0].getPort() + "/{connectionReserveTimeoutSeconds=5}";
		tgPool = TGConnectionFactory.getInstance().createConnectionPool(url, tgUser, tgPwd, 10, null);
		tgPool.connect();
	}
	
	@AfterGroups("connectionPoolGrp")
	public void stopConnectionPoolServer() throws Exception {
		
		tgPool.disconnect();
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".connectionPool");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
	/**
	 * testIPv4Connect - Connect to server with various IPv4 addresses
	 * @throws Exception
	 */
	@Test(	dataProvider = "ipv4Data", 
			groups = "ipv4Grp",
			description = "Server binds on various IPv4-type addresses - Connect to server with various IPv4 addresses")
	public void testIPv4Connect(String host, int port) throws Exception {
		
		String url = "tcp://scott@" + host + ":" + port;
		String pwd = "scott";
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, pwd, null);
		conn.connect();
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();

		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		conn.disconnect();
	}

	/**
	 * testIPv6Connect - Connect to server with various IPv6 addresses
	 * @throws Exception
	 */
	@Test(	dataProvider = "ipv6Data", 
			groups = "ipv6Grp", 
			description = "Server binds on various IPv6-type addresses - Connect to server with various IPv6 addresses")
	public void testIPv6Connect(String host, int port) throws Exception {
		
		String url = "tcp://[" + host + ":" + port+"]";
//		String url = "tcp://[fe80:0:0:0:e118:55ae:70a3:7369:8223]";
		
		
		System.out.println("tcp://[fe80:0:0:0:e118:55ae:70a3:7369:8223]\n\n\n");
		String user = "scott";
		String pwd = "scott";
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
		conn.connect();
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();

		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		conn.disconnect();
	}
	
	/**
	 * testAuthConnect - Authenticate to server with various user/pwd/role
	 * @throws Exception
	 */
	@Test(	dataProvider = "authData", 
			groups = "authGrp",
			description = "Authenticate to the server with various user/passwd/role")
	public void testAuthConnect(String user, String pwd, String role, String expectedLogin) throws Exception {

		String actualLogin = "";
		String url = "tcp://localhost:" + tgServer.getNetListeners()[0].getPort();
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);

		try {
			conn.connect();
			actualLogin = "loginSuccess";
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();

			if (gof == null) {
				throw new org.testng.TestException("TG object factory is null");
			}
		}
		catch (TGAuthenticationException ae) {
			actualLogin = "loginFailure";
		}
		finally {
			if (conn != null)
				conn.disconnect();
			Assert.assertEquals(actualLogin, expectedLogin);
		}
		
	}

	/**
	 * testMaxConnectWithNoPool - Try to connect with more than the max number of connections allowed
	 * @throws Exception
	 */
	@Test(	dataProvider = "maxConnectData",
			groups = "maxConnectGrp",
			description = "Connect more than the max number of connections allowed")
	public void testMaxConnectWithNoPool(String configFile, String user, String pwd) throws Exception {
			
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/" + configFile, tgWorkingDir + "/" + configFile);
		
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);	
		
		String url = "tcp://localhost:" + tgServer.getNetListeners()[0].getPort();
		int maxConnectionIPv4 = tgServer.getNetListeners()[0].getMaxConnections();
		int currentConnectionIPv4 = 0;
		try {
			// connect maxconnection+1 - Last connection should fail
			for (currentConnectionIPv4=0; currentConnectionIPv4 < maxConnectionIPv4+1; currentConnectionIPv4++) { 
				TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
				conn.connect();
				TGGraphObjectFactory gof = conn.getGraphObjectFactory();
				if (gof == null)	
					Assert.fail("TG object factory is null for connection #" + (currentConnectionIPv4+1));
			}
			Assert.fail("Expected a TGChannelDisconnectedException since maxConnection is " + maxConnectionIPv4 + " and current connection is already #" + currentConnectionIPv4);
		}
		catch(Exception e) {
			if (!(e instanceof TGChannelDisconnectedException))
				Assert.fail("Expected a TGChannelDisconnectedException upon connection but got a " + e.getClass().getName() + " instead");
			Assert.assertEquals(currentConnectionIPv4+1, maxConnectionIPv4+1, "Expected TGChannelDisconnectedException on connection #" + (maxConnectionIPv4+1) + " but got it on connection #" + (currentConnectionIPv4+1));
		}
		finally { 
			tgServer.kill();
		}
	}
	
	/**
	 * testMaxConnectWithPool - Try to connect with more than the max number of connections allowed
	 * @throws Exception
	 */
	@Test(	dataProvider = "maxConnectData",
			groups = "maxConnectGrp",
			description = "Connect pools more than the max number of connections allowed")
	public void testMaxConnectWithPool(String configFile, String user, String pwd) throws Exception {
		
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/" + configFile, tgWorkingDir + "/" + configFile);
		
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);	
		
		String url = "tcp://localhost:" + tgServer.getNetListeners()[0].getPort();
		int maxConnectionIPv4 = tgServer.getNetListeners()[0].getMaxConnections();
		int currentConnectionIPv4 = 0;
		try {
			for (currentConnectionIPv4=0; currentConnectionIPv4 < maxConnectionIPv4+1; currentConnectionIPv4++) {
				TGConnectionPool pool = TGConnectionFactory.getInstance().createConnectionPool(url, user, pwd, 10, null);
				pool.connect();
				TGGraphObjectFactory gof = pool.get().getGraphObjectFactory();
				if (gof == null) {
					Assert.fail("TG object factory is null for connection pool #" + (currentConnectionIPv4+1));
				}
			}
			Assert.fail("Expected a TGChannelDisconnectedException since maxConnection is " + maxConnectionIPv4 + " and current connection is already #" + currentConnectionIPv4);
		}
		catch(Exception e) {
			e.printStackTrace();
			if (!(e instanceof TGChannelDisconnectedException))
				Assert.fail("Expected a TGChannelDisconnectedException upon connection but got a " + e.getClass().getName() + " instead");
			Assert.assertEquals(currentConnectionIPv4+1, maxConnectionIPv4+1, "Expected TGChannelDisconnectedException on connection #" + (maxConnectionIPv4+1) + " but got it on connection #" + (currentConnectionIPv4+1));
		}
		finally { 
			tgServer.kill();
		}
	}
	
	/**
	 * testMassConnectWithNoPool - Connect and disconnect 5,000 times
	 * @throws Exception
	 */
	@Test(	groups = "massConnectGrp",
			description = "Connect and disconnect 5,000 times",
			timeOut = 500000)
	public void testMassConnectWithNoPool() throws Exception {
		
		String url = "tcp://"+ tgServer.getNetListeners()[0].getHost() + ":" + tgServer.getNetListeners()[0].getPort();
		int nbConnection = 5000;
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(url,"scott", "scott", null);
		try {
			for (int i=0; i < nbConnection; i++) { 
				
				conn.connect();
				TGGraphObjectFactory gof = conn.getGraphObjectFactory();
				if (gof == null)	
					Assert.fail("TG object factory is null for connection #" + nbConnection);
				conn.disconnect();
			}
		}
		finally { 
			try {conn.disconnect();}
			catch(Exception e) {;}
		}
	}
	
	/**
	 * testMassConnectWithPool - Get and release a connection from the pool 10,000 times
	 * @throws Exception
	 */
	@Test(	groups = "massConnectGrp",
			description = "Get and release a connection from the pool 10,000 times",
			timeOut = 120000)
	public void testMassConnectWithPool() throws Exception {
		TGConnection conn = null;
		int nbConnection = 10000;
		String url = "tcp://"+ tgServer.getNetListeners()[0].getHost() + ":" + tgServer.getNetListeners()[0].getPort();
		TGConnectionPool pool = TGConnectionFactory.getInstance().createConnectionPool(url, tgUser, tgPwd, 1, null); // consume 1 connection
		pool.connect();
		try {
			for (int i=0; i < nbConnection; i++) { 
				conn = pool.get();
				TGGraphObjectFactory gof = conn.getGraphObjectFactory();
				if (gof == null)	
					Assert.fail("TG object factory is null for connection #" + nbConnection);
				pool.release(conn);
			}
		}
		finally { 
			if (pool != null)
				pool.disconnect();
		}
	}
	
	/**
	 * testMassConnectWithPool2 - Connect and disconnect a connection pool 5,000 times
	 * @throws Exception
	 */
	@Test(	groups = "massConnectGrp",
			description = "Connect and disconnect a connection pool 5,000 times",
			timeOut = 600000)
	public void testMassConnectWithPool2() throws Exception {
		TGConnection conn = null;
		TGConnectionPool pool = null;
		int nbConnection = 5000;
		String url = "tcp://"+ tgServer.getNetListeners()[0].getHost() + ":" + tgServer.getNetListeners()[0].getPort();
		pool = TGConnectionFactory.getInstance().createConnectionPool(url, tgUser, tgPwd, 10, null); // consume 1 connection
		try {
			for (int i=0; i < nbConnection; i++) {
				long startTime = System.currentTimeMillis();
				pool.connect();
				long endTime = System.currentTimeMillis();
				conn = pool.get();
				TGGraphObjectFactory gof = conn.getGraphObjectFactory();
				if (gof == null)	
					Assert.fail("TG object factory is null for connection #" + nbConnection);
				//pool.release(conn); // skip releasing connection
				System.out.println("LOOP ITER #" + i + " - Connect Time = " + (endTime - startTime) + " msec");
				pool.disconnect();
			}
		}
		finally { 
			try { pool.disconnect();}
			catch (Exception e) {;}
		}
	}
	
	/**
	 * testMassConnectWithPool3 - Connect and disconnect a dedicated-channel connection pool 1,000 times
	 * @throws Exception
	 */
	@Test(	groups = "massConnectGrp",
			description = "Connect and disconnect a dedicated-channel connection pool 1,000 times",
			timeOut = 600000)
	public void testMassConnectWithPool3() throws Exception {
		TGConnection conn = null;
		TGConnectionPool pool = null;
		int nbConnection = 1000;
		String url = "tcp://"+ tgServer.getNetListeners()[0].getHost() + ":" + tgServer.getNetListeners()[0].getPort() + "/{useDedicatedChannelPerConnection=true}";
		pool = TGConnectionFactory.getInstance().createConnectionPool(url, tgUser, tgPwd, 10, null); // consume 10 connections
		try {
			for (int i=0; i < nbConnection; i++) {
				
				long startTime = System.currentTimeMillis();
				pool.connect();
				long endTime = System.currentTimeMillis();
				conn = pool.get();
				TGGraphObjectFactory gof = conn.getGraphObjectFactory();
				if (gof == null)	
					Assert.fail("TG object factory is null for connection #" + nbConnection);
				// pool.release(conn); // skip releasing connection
				System.out.println("LOOP ITER #" + i + " - Connect Time = " + (endTime - startTime) + " msec");
				pool.disconnect();
			}
		}
		finally { 
			try { pool.disconnect();}
			catch (Exception e) {;}
		}
	}
	
	/**
	 * testConnectionPool - Connect 20 clients to a pool of size 10. All clients should connect successfully since connectionReserveTimeoutSeconds allows sufficient time
	 * @throws Exception
	 */
	@Test(	groups = "connectionPoolGrp",
			description = "Connect 20 clients to a pool of size 10. All clients should connect successfully since connectionReserveTimeoutSeconds allows sufficient time",
			invocationCount = 20,
			threadPoolSize = 20,
			timeOut = 60000)
	public void testConnectionPool() throws Exception {
		
		tgConn = tgPool.get();
		TGGraphObjectFactory gof = tgConn.getGraphObjectFactory();
		if (gof == null)	
			Assert.fail("TG object factory is null");
		Thread.sleep(3000); // connectionReserveTimeoutSeconds is longer than this sleep so all clients should connect
		tgPool.release(tgConn);
				
	}
	
	/**
	 * testConnectionPool - Connect 20 clients to a pool of size 10. Half clients should fail to connect since connectionReserveTimeoutSeconds is short
	 * @throws Exception
	 */
	@Test(	groups = "connectionPoolGrp",
			description = "Connect 20 clients to a pool of size 10. Half clients should fail to connect since connectionReserveTimeoutSeconds is short",
			invocationCount = 20,
			threadPoolSize = 20,
			timeOut = 60000)
	public void testConnectionPool2() throws Exception {
		int expectedPoolTimeout = 10;
		int expectedPoolSuccess = 10;
		try {
			tgConn = tgPool.get();
			synchronized (this) {
				nbPoolSuccess  ++;
				//System.out.println("Number of Pool Get : " + nbPoolSuccess);
			}
			Assert.assertTrue(nbPoolSuccess <= expectedPoolSuccess, "The number of successful pool get is higher than " + expectedPoolSuccess);
			TGGraphObjectFactory gof = tgConn.getGraphObjectFactory();
			if (gof == null)	
				Assert.fail("TG object factory is null");
			Thread.sleep(6000); // sleep longer than connectionReserveTimeoutSeconds so half clients should fail to connect
			tgPool.release(tgConn);
		}
		catch (TGConnectionTimeoutException tge) {
			synchronized(this) {
				nbPoolTimeout  ++;
				//System.out.println("Number of TimeoutException : " + nbPoolTimeout);
			}
			Assert.assertTrue(nbPoolTimeout <= expectedPoolTimeout, "The number of TimeoutException is higher than " + expectedPoolTimeout);
		}
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	@DataProvider(name = "ipv4Data")
	public Object[][] getIPv4() throws IOException, EvalError {
		String[] host = new String[2];
		host[0] = InetAddress.getLocalHost().getHostName();
		host[1] = "localhost"; // get ipv4 loopback address as well
		int port = 8222;
		
		List<Object[]> urlParams = new ArrayList<Object[]>();
		System.setProperty("java.net.preferIPv6Addresses", "false");
		
		// Get all the IPv6 addresses on the local machine
		for (int i=0; i<host.length; i++) {
			urlParams.add(new Object[] {host[i],port});
			InetAddress[] addr = InetAddress.getAllByName(host[i]);
	    	for (InetAddress address : addr) {
	    		if (address instanceof Inet4Address) {
	    			String tmpAddr = address.getHostAddress();
	    			urlParams.add(new Object[] {tmpAddr,port});
	    		}
	    	}
		}
		return (Object[][])urlParams.toArray(new Object[urlParams.size()][2]);
	}

	/**
	 * Get all IPv6 addresses available on the current machine
	 * @throws UnknownHostException 
	 */
	@DataProvider(name = "ipv6Data")
	public Object[][] getIPv6() throws UnknownHostException {
		String[] host = new String[2];
		host[0] = InetAddress.getLocalHost().getHostName();
		host[1] = "localhost"; // get ipv6 loopback address as well
		int port = 8223;
		
		List<Object[]> urlParams = new ArrayList<Object[]>();
		System.setProperty("java.net.preferIPv6Addresses", "true");
		
		// Get all the IPv6 addresses on the local machine
		for (int i=0; i<host.length; i++) {
			urlParams.add(new Object[] {host[i],port});
			InetAddress[] addr = InetAddress.getAllByName(host[i]);
	    	for (InetAddress address : addr) {
	    		if (address instanceof Inet6Address) {
	    			String tmpAddr = address.getHostAddress().substring(0, (address.getHostAddress().contains("%")?address.getHostAddress().indexOf('%'):address.getHostAddress().length()));
	    			urlParams.add(new Object[] {tmpAddr,port});
	    		}
	    	}
		}
		return (Object[][])urlParams.toArray(new Object[urlParams.size()][2]);
	}
	
	@DataProvider(name = "authData")
	public Object[][] getUsers() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/Auth.data"));
		return data;
	}
	
	@DataProvider(name = "maxConnectData")
	public Object[][] getMaxConnection() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/MaxConnect.data"));
		return data;
	}
}
