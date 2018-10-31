package com.tibco.tgdb.test.admin.user;

import org.testng.annotations.Test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.utils.ClasspathResource;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.util.Scanner;

import org.testng.Assert;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.BeforeSuite;

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

/**
 * User-related test cases from Admin console.
 * After each admin operation (create, update ...) the server and admin are re-started.
 * @author Serge
 *
 */
public class UserTests {

	private static TGServer tgServer;
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");

	final private String userCreationSuccessMsg = "Successfully created user on server.";
	final private String userCreationDuplicateMsg = "A duplicate entity exists on the server.";
	final private String userUpdateSuccessMsg = "???"; // Not yet defined in the product
	final private String userDropSuccessMsg = "???"; // Not yet defined in the product
	final private String userShowMsg = "user(s) returned";
	
	int expectedNbUsers = 0;
	
	/**
	 * Init TG server before test suite
	 * @throws Exception
	 */
	@BeforeSuite(description = "Init TG Admin")
	public void initServer() throws Exception  {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replaceFirst("\\.[a-z]*$", "").replace('.', '/') + "/initdb.conf", tgWorkingDir + "/initdb.conf");
		tgServer = new TGServer(tgHome);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 15000);
		}
		catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
	}
	
	/**
	 * Start TG server before each test method
	 * @throws Exception
	 */
	@BeforeMethod
	public void startServer() throws Exception {
		File confFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replaceFirst("\\.[a-z]*$", "").replace('.', '/') + "/tgdb.conf", tgWorkingDir + "/tgdb.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(15000);
	}

	/**
	 *  Kill TG server after each test method
	 * @throws Exception
	 */
	@AfterMethod
	public void killServer() throws Exception {
		TGAdmin.stopServer(tgServer, tgServer.getNetListeners()[0].getName(), null, null, 300000);
		//tgServer.kill();
		// Backup log file before moving to next test
		//File logFile = tgServer.getLogFile();
		//File backLogFile = new File(logFile + ".user");
		//Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
	/**
	 * testCreateUsers - Create 10,000 users via TG Admin
	 * @throws Exception
	 */
	@Test(description = "Create 10,000 users via TG Admin")
	public void testCreateUsers() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/CreateUsers.data", tgWorkingDir + "/CreateUsers.data");
		
		// Create users via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminCreateUsers.log", null, 
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Get expected number of user creations
		expectedNbUsers = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("create user"))
				expectedNbUsers++;
		}
		br.close();
		
		// Check actual user creation
		Scanner scanner = new Scanner(console);
		int userCreation = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(userCreationSuccessMsg))
				userCreation++;
		}
		scanner.close();
		Assert.assertEquals(userCreation, expectedNbUsers, "User creation does not match -");
	}
	
	/**
	 * testDuplicateUsers - Re-create the same users. Should get duplicate message in Admin
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testCreateUsers" },
		  description = "Re-create same users via TG Admin and check for duplicates")
	public void testDuplicateUsers() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/CreateUsers.data", tgWorkingDir + "/CreateUsers.data");
		
		// Re-Create users via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminDuplicateUsers.log", null, 
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Get expected number of user creations
		expectedNbUsers = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("create user"))
				expectedNbUsers++;
			}
		br.close();
		
		// Check user duplication 
		Scanner scanner = new Scanner(console);
		int userCreation = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(userCreationDuplicateMsg))
				userCreation++;
		}
		scanner.close();
		Assert.assertEquals(userCreation, expectedNbUsers, "User duplication does not match -");
	}
	
	/**
	 * testUpdateUsers - Update user passwords
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testCreateUsers" },
		  description = "Update users via TG Admin")
	public void testUpdateUsers() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/UpdateUsers.data", tgWorkingDir + "/UpdateUsers.data");
		
		// Update users via Admin
		String console = TGAdmin.invoke(tgHome, null, null, null, tgWorkingDir + "/adminUpdateUsers.log",
				cmdFile.getAbsolutePath(), 10000);
		//System.out.println(console);
		
		// Get expected number of user updates
		int expectedUserUpdate = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("update user"))
				expectedUserUpdate++;
			}
		br.close();
		
		// Check user update 
		Scanner scanner = new Scanner(console);
		int userUpdate = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(userUpdateSuccessMsg))
				userUpdate++;
		}
		scanner.close();
		Assert.assertEquals(userUpdate, expectedUserUpdate, "User update does not match -");
	}
	*/
	
	/**
	 * testShowUsers - Show users previously created
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testCreateUsers" },
		  description = "Show users in TG Admin")
	public void testShowUsers() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/ShowUsers.data", tgWorkingDir + "/ShowUsers.data");
		
		// Show users via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminShowUsers.log", null, cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Assert show users. +1 to count admin user
		Assert.assertTrue(console.contains((expectedNbUsers+1) + " " + userShowMsg), "Expected " + expectedNbUsers + "+1 " + userShowMsg + " but did not get that.");
	}
	
	/**
	 * testConnectUsers - Connect to server with previously 10,000 created users
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testCreateUsers" },
		  description = "Connect to TG server with previously 10,000 created users")
	public void testConnectUsers() throws Exception {
		
		String url = "tcp://" + tgServer.getNetListeners()[0].getHost() + ":" + tgServer.getNetListeners()[0].getPort();
		
		for (int i=0; i<expectedNbUsers; i++) { 
			TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, "user"+i, "pass"+i, null);
			conn.connect();
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			conn.disconnect(); 
			if (gof == null)	
				Assert.fail("TG object factory is null for user" + i);
		}
	}
	
	/**
	 * testConnectUpdatedUsers - Connect to server with previously created users
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testUpdateUsers" },
		  description = "Connect to TG server with previously updated users")
	public void testConnectUpdatedUsers() throws Exception {
		
		String url = "tcp://localhost:" + tgServer.getNetListeners()[0].getPort();
		String pwd = "pass"; // all pwd are "pass" after user update
		
		for (int i=0; i<expectedNbUsers; i++) { 
			TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, "user"+i, pwd, null);
			conn.connect();
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null)	
				Assert.fail("TG object factory is null for user" + i);
		}
	}
	*/
	
	/**
	 * testDropUsers - Drop users previously created via TG Admin
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testShowUsers" },
		  description = "Drop users previously created via TG Admin")
	public void testDropUsers() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DropUsers.data", tgWorkingDir + "/DropUsers.data");
		
		// Drop users via Admin
		String console = TGAdmin.invoke(tgHome, null, null, null, tgWorkingDir + "/adminDropUsers.log",
				cmdFile.getAbsolutePath(), 10000);
		//System.out.println(console);
		
		// Get expected number of user deleted
		expectedNbUsers = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("drop user"))
				expectedNbUsers++;
		}
		br.close();
		
		// Check actual user deletion
		Scanner scanner = new Scanner(console);
		int userDrop = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(userDropSuccessMsg))
				userDrop++;
		}
		scanner.close();
		Assert.assertEquals(userDrop, expectedNbUsers, "User deletion does not match -");
	}
	*/
	/**
	 * testConnectDroppedUsers - Connect to server with previously dropped users - Connection should fail
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testDropUsers" },
		  description = "Connect to TG server with previously dropped users - Connection should fail")
	public void testConnectDroppedUsers() throws Exception {
		
		String url = "tcp://localhost:" + tgServer.getNetListeners()[0].getPort();
		String pwd = "pass"; // all pwd are "pass" after user update
		
		for (int i=0; i<expectedNbUsers; i++) { 
			try {
				TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, "user"+i, pwd, null);
				conn.connect();
				TGGraphObjectFactory gof = conn.getGraphObjectFactory();
				conn.disconnect();
				if (gof == null)	
					Assert.fail("TG object factory is null for user" + i);
				Assert.fail("Expected a TGException upon connection since user was dropped but it seems connection was successful !");
			}
			catch(Exception e) {
				if (!(e instanceof TGException))
					Assert.fail("Expected a TGException upon connection but got a " + e.getClass().getName() + " instead");
			}
		}
	}
	*/
	
}
