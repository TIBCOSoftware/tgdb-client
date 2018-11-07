package com.tibco.tgdb.test.admin.nodetype;

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
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;
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
 * nodetype-related test cases from Admin
 * After each admin operation (create, update ...) the server and admin are re-started
 * @author Serge
 *
 */
public class NodeTypeTests {

	private static TGServer tgServer;
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");

	final private String nodeCreationSuccessMsg = "Successfully created nodetype on server.";
	final private String nodeCreationDuplicateMsg = "A duplicate entity exists on the server.";
	final private String nodeDescribeSuccessMsg = "???"; // Not yet defined in the product
	final private String nodeAlterSuccessMsg = "???"; // Not yet defined in the product
	final private String nodeDropSuccessMsg = "???"; // Not yet defined in the product
	final private String nodeShowMsg = "type(s) returned"; 
	
	int expectedNbNode = 0;
	
	/**
	 * Init TG server before test suite
	 * @throws Exception
	 */
	@BeforeSuite(description = "Init TG Admin")
	public void initServer() throws Exception  {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replaceFirst("\\.[a-z]*$", "").replace('.', '/') + "/initdb.conf", tgWorkingDir + "/initDB.conf");
		tgServer = new TGServer(tgHome);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 60000);
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
		tgServer.start(10000);
	}

	/**
	 * Kill TG server after each test method
	 * @throws Exception
	 */
	@AfterMethod
	public void killServer() throws Exception {
		TGAdmin.stopServer(tgServer, tgServer.getNetListeners()[0].getName(), null, null, 60000);
		//tgServer.kill();
		// Backup log file before moving to next test
		//File logFile = tgServer.getLogFile();
		//File backLogFile = new File(logFile + ".node");
		//Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
	/**
	 * testCreateNodeTypes - Create various node types via TG Admin
	 * @throws Exception
	 */
	@Test(description = "Create various node types via TG Admin")
	public void testCreateNodeTypes() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/CreateNode.data", tgWorkingDir + "/CreateNode.data");
		
		// Create nodetypes via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminCreateNode.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Get expected number of nodetype creations
		expectedNbNode = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.matches("^create nodetype .*"))
				expectedNbNode++;
		}
		br.close();
		
		// Check actual nodetype creation
		Scanner scanner = new Scanner(console);
		int nodeCreation = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(nodeCreationSuccessMsg))
				nodeCreation++;
		}
		scanner.close();
		Assert.assertEquals(nodeCreation, expectedNbNode, "Nodetype creation does not match -");
	}
	
	/**
	 * testDuplicateNodeTypes - Re-create the same nodetypes. Should get duplicate message in Admin
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testCreateNodeTypes" },
		  description = "Re-create same nodetypes via TG Admin and check for duplicates")
	public void testDuplicateNodeTypes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DuplicateNode.data", tgWorkingDir + "/DuplicateNode.data");
		
		// Re-Create nodetypes via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminDuplicateNode.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Check nodetype duplication 
		Scanner scanner = new Scanner(console);
		int nodeDuplicate = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(nodeCreationDuplicateMsg))
				nodeDuplicate++;
		}
		scanner.close();
		Assert.assertEquals(nodeDuplicate, expectedNbNode, "Nodetype duplication does not match -");
	}
	
	/**
	 * testShowNodeTypes - Show nodetypes previously created
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testDuplicateNodeTypes" },
		  description = "Show nodetypes in TG Admin")
	public void testShowNodeTypes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/ShowNode.data", tgWorkingDir + "/ShowNode.data");
		
		// Show nodetype via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminShowNode.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Check show nodetype. We expect expectedNbNode nodetypes (16 nodetypes) + 1 default nodetype + 3 default edgetypes ==> 20 types
		Assert.assertTrue(console.contains((expectedNbNode+1+3) + " " + nodeShowMsg), "Expected " + (expectedNbNode+1+3) + " " + nodeShowMsg + " but did not get that.");
	}
	
	/**
	 * testAlterNodeTypes - Alter nodetypes previously created via TG Admin
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testShowNodeTypes" },
		  description = "Alter nodetypes previously created via TG Admin")
	public void testAlterNodeTypes() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/AlterNode.data", tgWorkingDir + "/AlterNode.data");
		
		// Alter nodetype via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminAlterNode.log",
				cmdFile.getAbsolutePath(), 10000);
		//System.out.println(console);
		
		// Get expected number of nodetype altered
		expectedNbNode = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("alter nodetype"))
				expectedNbNode++;
		}
		br.close();
		
		// Check actual nodetype alteration
		Scanner scanner = new Scanner(console);
		int nodeAlter = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(nodeAlterSuccessMsg))
				nodeAlter++;
		}
		scanner.close();
		Assert.assertEquals(nodeAlter, expectedNbNode, "Nodetype alteration does not match -");
	}
	*/
	
	/**
	 * testDescribeNodeTypes - Describe nodetypes previously created
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testAlterNodeTypes" },
		  description = "Describe node types in TG Admin")
	public void testDescribeNodeTypes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DescribeNode.data", tgWorkingDir + "/DescribeNode.data");
		
		// Describe nodetypes via Admin
		String console = TGAdmin.invoke(tgHome, null, null, null, tgWorkingDir + "/adminDescribeNode.log",
				cmdFile.getAbsolutePath(), 10000);
		// System.out.println(console);
		
		// Get expected number of nodetype description
		expectedNbNode = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("describe nodetype"))
				expectedNbNode++;
			}
		br.close();
		
		// Check describe nodetype
		Scanner scanner = new Scanner(console);
		int nodeDescribe = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(nodeDescribeSuccessMsg))
				nodeDescribe++;
		}
		scanner.close();
		Assert.assertEquals(nodeDescribe, expectedNbNode, "Describe nodetype does not match -");
	}
	*/
	
	/**
	 * testDropNodeTypes - Drop nodetypes previously created via TG Admin
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testShowNodeTypes" },
		  description = "Drop nodetypes previously created via TG Admin")
	public void testDropNodeTypes() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DropNode.data", tgWorkingDir + "/DropNode.data");
		
		// Drop nodetype via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminDropNode.log",
				cmdFile.getAbsolutePath(), 10000);
		//System.out.println(console);
		
		// Get expected number of nodetype deleted
		expectedNbNode = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("drop nodetype"))
				expectedNbNode++;
		}
		br.close();
		
		// Check actual nodetype deletion
		Scanner scanner = new Scanner(console);
		int nodeDrop = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(nodeDropSuccessMsg))
				nodeDrop++;
		}
		scanner.close();
		Assert.assertEquals(nodeDrop, expectedNbNode, "Nodetype deletion does not match -");
	}
	*/
}
