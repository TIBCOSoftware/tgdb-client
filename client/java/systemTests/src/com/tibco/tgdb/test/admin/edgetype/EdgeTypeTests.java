package com.tibco.tgdb.test.admin.edgetype;

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
 * Edgetype-related test cases from Admin
 * After each admin operation (create, update ...) the server and admin are re-started
 * @author Serge
 *
 */
public class EdgeTypeTests {

	private static TGServer tgServer;
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");

	final private String edgeCreationSuccessMsg = "Successfully created edgetype on server.";
	final private String edgeCreationDuplicateMsg = "A duplicate entity exists on the server.";
	final private String edgeDescribeMsg = "???";
	final private String edgeAlterMsg = "???";
	final private String edgeDropMsg = "???";
	final private String edgeShowMsg = "type(s) returned";
	
	int expectedNbEdge = 0;
	
	// Init TG server
	@BeforeSuite(description = "Init TG Admin")
	public void initServer() throws Exception  {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/../Initdb.conf", tgWorkingDir + "/InitDB.conf");
		tgServer = new TGServer(tgHome);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 60000);
		}
		catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
	}
	
	// Start TG server
	@BeforeMethod
	public void startServer() throws Exception {
		File confFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/../tgdb.conf", tgWorkingDir + "/tgdb.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}

	// Kill TG server
	@AfterMethod
	public void killServer() throws Exception {
		TGAdmin.stopServer(tgServer, tgServer.getNetListeners()[0].getName(), null, null, 60000);
		//tgServer.kill();
		// Backup log file before moving to next test
		//File logFile = tgServer.getLogFile();
		//File backLogFile = new File(logFile + ".edge");
		//Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
	/**
	 * testCreateEdgeTypes - Create various edgetypes via TG Admin
	 * @throws Exception
	 */
	@Test(description = "Create various edgetypes via TG Admin")
	public void testCreateEdgeTypes() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/CreateEdge.data", tgWorkingDir + "/CreateEdge.data");
		
		// Create edgetypes via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminCreateEdge.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Get expected number of edgetype creations
		expectedNbEdge = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.matches("^create .* edgetype.*"))
				expectedNbEdge++;
		}
		br.close();
		
		// Check actual edgetype creation
		Scanner scanner = new Scanner(console);
		int edgeCreation = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(edgeCreationSuccessMsg))
				edgeCreation++;
		}
		scanner.close();
		Assert.assertEquals(edgeCreation, expectedNbEdge, "Edgetype creation does not match -");
	}
	
	/**
	 * testDuplicateEdgeTypes - Re-create the same edgetypes. Should get duplicate message in Admin
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testCreateEdgeTypes" },
		  description = "Re-create same edgetypes via TG Admin and check for duplicates")
	public void testDuplicateEdgeTypes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DuplicateEdge.data", tgWorkingDir + "/DuplicateEdge.data");
		
		// Re-Create edgetypes via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminDuplicateEdge.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Check edgetype duplication 
		Scanner scanner = new Scanner(console);
		int edgeDuplicate = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(edgeCreationDuplicateMsg))
				edgeDuplicate++;
		}
		scanner.close();
		Assert.assertEquals(edgeDuplicate, expectedNbEdge, "Edgetype duplication does not match -");
	}
	
	
	 
	
	/**
	 * testShowEdgeTypes - Show edgetypes previously created
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testDuplicateEdgeTypes" },
		  description = "Show edgetypes in TG Admin")
	public void testShowEdgeTypes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/ShowEdge.data", tgWorkingDir + "/Showedge.data");
		
		// Show edgetype via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminShowEdge.log", null, 
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Check show edgetype. We expect expectedNbEdge (12 edgetypes we created) + 3 default edgetypes + 2 nodetypes we created for the test + 1 default nodetype => 18 types
		Assert.assertTrue(console.contains((expectedNbEdge+3+2+1) + " " + edgeShowMsg), "Expected " + (expectedNbEdge+3+2+1) + " " + edgeShowMsg + " but did not get that -");
	}
	
	/**
	 * testAlterEdgeTypes - Alter edgetypes previously created via TG Admin
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testShowEdgeTypes" },
		  description = "Alter edgetypes previously created via TG Admin")
	public void testAlterEdgeTypes() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/AlterEdge.data", tgWorkingDir + "/AlterEdge.data");
		
		// Alter edgetype via Admin
		String console = TGAdmin.invoke(tgHome, null, null, null, tgWorkingDir + "/adminAlterEdge.log",
				cmdFile.getAbsolutePath(), 10000);
		//System.out.println(console);
		
		// Get expected number of edgetype altered
		expectedNbEdge = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("alter edgetype"))
				expectedNbEdge++;
		}
		br.close();
		
		// Check actual edgetype alteration
		Scanner scanner = new Scanner(console);
		int edgeAlter = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(edgeAlterMsg))
				edgeAlter++;
		}
		scanner.close();
		Assert.assertEquals(edgeAlter, expectedNbEdge, "Edgetype alteration does not match -");
	}
	*/
	/**
	 * testDescribeedgeTypes - Describe edgetypes previously created
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testAlterEdgeTypes" },
		  description = "Describe edgetypes in TG Admin")
	public void testDescribeEdgeTypes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DescribeEdge.data", tgWorkingDir + "/Describeedge.data");
		
		// Describe edgetypes via Admin
		String console = TGAdmin.invoke(tgHome, null, null, null, tgWorkingDir + "/adminDescribeEdge.log",
				cmdFile.getAbsolutePath(), 10000);
		// System.out.println(console);
		
		// Get expected number of edgetype description
		expectedNbEdge = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("describe edgetype"))
				expectedNbEdge++;
			}
		br.close();
		
		// Check describe edgetype
		Scanner scanner = new Scanner(console);
		int edgeDescribe = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(edgeDescribeMsg))
				edgeDescribe++;
		}
		scanner.close();
		Assert.assertEquals(edgeDescribe, expectedNbEdge, "Describe edgetype does not match -");
	}
	*/
	/**
	 * testDropEdgeTypes - Drop edgetypes previously created via TG Admin
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testDescribeEdgeTypes" },
		  description = "Drop edgetypes previously created via TG Admin")
	public void testDropEdgeTypes() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DropEdge.data", tgWorkingDir + "/DropEdge.data");
		
		// Drop edgetype via Admin
		String console = TGAdmin.invoke(tgHome, null, null, null, tgWorkingDir + "/adminDropEdge.log",
				cmdFile.getAbsolutePath(), 10000);
		//System.out.println(console);
		
		// Get expected number of edgetype deleted
		expectedNbEdge = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("drop edgetype"))
				expectedNbEdge++;
		}
		br.close();
		
		// Check actual edgetype deletion
		Scanner scanner = new Scanner(console);
		int edgeDrop = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(edgeDropMsg))
				edgeDrop++;
		}
		scanner.close();
		Assert.assertEquals(edgeDrop, expectedNbEdge, "Edgetype deletion does not match -");
	}
	*/
}
