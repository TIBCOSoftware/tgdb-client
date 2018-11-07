package com.tibco.tgdb.test.admin.connect;

import org.testng.annotations.Test;
import org.testng.Assert;
import org.testng.annotations.AfterSuite;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.DataProvider;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGAdminException;
import com.tibco.tgdb.test.lib.TGGeneralException;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.utils.ClasspathResource;
import com.tibco.tgdb.test.utils.PipedData;

import bsh.EvalError;

import java.io.File;
import java.io.IOException;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.nio.file.StandardOpenOption;
import java.util.ArrayList;
import java.util.List;

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
	
	final private String adminConnectSuccessMsg = "Successfully connected to server";	
	
	// Config file is used 2 times in this class
	private File getConfigFile() throws IOException {
		File confFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replaceFirst("\\.[a-z]*$", "").replace('.', '/') + "/tgdb.conf",tgWorkingDir + "/tgdb.conf");
		return confFile;
	}

	/**
	 * Init TG server before test suite
	 * @throws Exception
	 */
	@BeforeSuite(description = "Init TG Admin")
	public void initServer() throws Exception {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replaceFirst("\\.[a-z]*$", "").replace('.', '/') + "/initdb.conf",tgWorkingDir + "/initdb.conf");
		tgServer = new TGServer(tgHome);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 60000);
		} catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
		System.out.println(tgServer.getBanner());
		tgServer.setConfigFile(getConfigFile());
		tgServer.start(15000);
	}

	/**
	 * Kill TG server after suite
	 * @throws Exception
	 */
	@AfterSuite
	public void killServer() throws Exception {
		tgServer.kill();
		// Backup log file before moving to next test
		File logFile = tgServer.getLogFile();
		File backLogFile = new File(logFile + ".adminconnection");
		Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}
	
	/**
	 * Start TG server before test
	 * @throws Exception
	 */
	/*@BeforeTest
	public void startServer() throws Exception {
		File confFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/../tgdb.conf",
				tgWorkingDir + "/tgdb.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}*/

	/**
	 * Kill TG server after test
	 * @throws Exception
	 */
	//@AfterTest
	//public void killServer() throws Exception {
		//tgServer.kill();
		// Backup log file before moving to next test
		//File logFile = tgServer.getLogFile();
		//File backLogFile = new File(logFile + ".adminconnection");
		//Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	//}

	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/

	/**
	 * testIPv6Connect - Connect TG Admin to TG Server via IPv6
	 * 
	 * @throws Exception
	 */
	@Test(dataProvider = "ipv6Data",
		  description = "Connect TG Admin to TG Server via IPv6")
	public void testIPv6Connect(String host, int port) throws Exception {

		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/Connection.cmd",
				tgWorkingDir + "/Connection.cmd");

		// Start admin console and connect via IPv6
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[1].getName(), tgWorkingDir + "/admin.ipv6.log", null, cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);

		Assert.assertTrue(console.contains(adminConnectSuccessMsg), "Admin did not connect to server");
	}
	
	/**
	 * testWrongUserPwd - Try connecting TG Admin to TG Server via IPv6 with wrong user/pwd
	 * 
	 * @throws Exception
	 */
	@Test(dataProvider = "wrongUserData",
		  description = "Try connecting TG Admin to TG Server via IPv6 with wrong user/pwd")
	public void testWrongUserPwd(String user, String pwd) throws Exception {

		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/Connection.cmd",
				tgWorkingDir + "/Connection.cmd");
		String console = "";
		String url = "tcp://[" + tgServer.getNetListeners()[1].getHost() + ":" + tgServer.getNetListeners()[1].getPort() + "]";
		try {
			// Start admin console and connect via IPv6 with wrong user/pwd
			console = TGAdmin.invoke(tgHome, url, user, pwd, tgWorkingDir + "/admin.wronguserpwd.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
			System.out.println(console);
			Assert.fail("Expected a TGAdminException due to wrong user/pwd but did not get it");
		}
		catch(TGAdminException e) { // Expected since wrong user/pwd
			// Even though we got the exception, make sure it is for the good reason
			Assert.assertFalse(console.contains(adminConnectSuccessMsg), "Admin connected to server even though user/pwd was wrong");
		}
	}
	
	/**
	 * testWrongUrl - Try connecting TG Admin to TG Server with wrong url
	 * 
	 * @throws Exception
	 */
	@Test(description = "Try connecting TG Admin to TG Server with wrong url")
	public void testWrongUrl() throws Exception {

		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/Connection.cmd",
				tgWorkingDir + "/Connection.cmd");

		String host = "my.machine.com"; // random host
		int port = 1234; // random port
		String console = "";
		try {
			// Start admin console and connect via IPv6 with wrong url
			console = TGAdmin.invoke(tgHome, "tcp://"+host+":"+port, tgServer.getSystemUser(), tgServer.getSystemPwd(), tgWorkingDir + "/admin.wronguserpwd.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
			//System.out.println(console);
			Assert.fail("Expected a TGAdminException due to wrong url but did not get it");
		}
		catch(TGAdminException e) { // Expected since wrong url
			// Even though we got the exception, make sure it is for the good reason
			Assert.assertFalse(console.contains(adminConnectSuccessMsg), "Admin connected to server even though url was wrong");
		}
	}
	
	/**
	 * testShowConnections - Connect 5 times via API and show connections in TG admin"
	 * 
	 * @throws Exception
	 */
	@Test(description = "Connect 5 times via API and show connections in TG admin")
	public void testShowConnections() throws Exception {
		
		String user = "user1";
		String pwd = "pass1";
		int nbConnections = 5;
		// create a new user
		if (!TGAdmin.createUser(tgServer, tgServer.getNetListeners()[0].getName(), user, pwd, null, null, 5000))
			throw new org.testng.TestException("TGAdmin could not create user " + user);
		
		// connect that user 5 times - do not disconnect
		for (int i=0; i<nbConnections; i++) {
			String url = "tcp://" + tgServer.getNetListeners()[0].getHost() + ":" + tgServer.getNetListeners()[0].getPort();
			TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
			conn.connect();
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();

			if (gof == null) {
				throw new org.testng.TestException("TG object factory is null on connection #" + i);
			}
		}
		String[] connection = TGAdmin.getConnectionsByUser(tgServer, tgServer.getNetListeners()[0].getName(), user, null, null, 5000);
		Assert.assertEquals(connection.length, nbConnections, "Expected " + nbConnections + " connections but got " + connection.length + " -");
	}
	
	/**
	 * testKillConnection - Retrieve connections and kill one. 
	 * Show connections in TG admin should display one less connection"
	 * 
	 * @throws Exception
	 */
	@Test(description = "Retrieve connections and kill one. Show connections in TG admin should display one less connection",
		  dependsOnMethods = { "testShowConnections" })
	public void testKillConnection() throws Exception {
		String user = "user1";
		int nbConnections = 5;
		// Retrieve connection for the user
		String[] connection = TGAdmin.getConnectionsByUser(tgServer, tgServer.getNetListeners()[0].getName(), user, null, null, 5000);
		if (connection.length != nbConnections)
			throw new org.testng.TestException("Expected " + nbConnections + " but got " + connection.length);
		// Kill one of them
		TGAdmin.killConnection(tgServer, tgServer.getNetListeners()[0].getName(), connection[0], null, null, 5000);
		// Now check nb connections again. It should be 1 less
		connection = TGAdmin.getConnectionsByUser(tgServer, tgServer.getNetListeners()[0].getName(), user, null, null, 5000);
		Assert.assertEquals(connection.length, (nbConnections-1), "Expected " + (nbConnections-1) + " connections after 1 killed but got " + connection.length + " -");
	}
	
	/**
	 * testMassConnections - Connect and disconnect TG admin 10,000 times"
	 * 
	 * Note: this test is disabled because tgdb-admin crashes in this test,
	 * when it is started by the Apache common-exec library (we use it internally),
	 * but it is working fine when we take the same parameters and run it 
	 * manually from a terminal
	 * 
	 * @throws Exception
	 */
	@Test(description = "Connect and disconnect TG admin 10,000 times",
		  enabled = false)
	public void testMassConnections() throws Exception {
		
		int nbConnections = 10000;
		String user = tgServer.getSystemUser();
		String pwd = tgServer.getSystemPwd();
		String host = tgServer.getNetListeners()[0].getHost();
		int port = tgServer.getNetListeners()[0].getPort();
		String url = "tcp://" + (tgServer.getNetListeners()[0].isIPv6()?"["+host+":"+port+"]":host+":"+port);
		
		// create admin script file
		File cmdFile = new File(tgWorkingDir + "/connAdminScript.txt");
		Files.write(Paths.get(cmdFile.toURI()), "".getBytes(StandardCharsets.UTF_8));
		String adminCmd = "connect " + url + " " + user + " " + pwd + "\ndisconnect\n";
		for (int i=0; i<nbConnections; i++) {
			Files.write(Paths.get(cmdFile.toURI()), adminCmd.getBytes(StandardCharsets.UTF_8), StandardOpenOption.APPEND);
		}
		Files.write(Paths.get(cmdFile.toURI()), "exit".getBytes(StandardCharsets.UTF_8), StandardOpenOption.APPEND);
		
		TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), null, null, cmdFile.getAbsolutePath(), -1, 400000);
	}
	
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/

	/**
	 * Get all IPv6 addresses available on the current machine
	 * @throws TGGeneralException 
	 * @throws IOException 
	 */
	@DataProvider(name = "ipv6Data")
	public Object[][] getIPv6() throws TGGeneralException, IOException {
		
		String[] host = new String[2];
		host[0] = InetAddress.getLocalHost().getHostName();
		host[1] = "localhost"; // get ipv6 loopback address as well
		
		// We need to get a new server here to get the port
		// since @DataProvider might run before @BeforeSuite tgServer might not exist yet
		TGServer tgTempServer = new TGServer(tgHome);
		tgTempServer.setConfigFile(getConfigFile());
		int port = tgTempServer.getNetListeners()[1].getPort(); // get port of ipv6 listener
		
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
	
	/**
	 * Get several combinations of wrong user and pwd 
	 */
	@DataProvider(name = "wrongUserData")
	public Object[][] getUsers() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/WrongUsers.data"));
		return data;
	}

}
