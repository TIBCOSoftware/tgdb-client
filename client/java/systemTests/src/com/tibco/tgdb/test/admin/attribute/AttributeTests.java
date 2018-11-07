package com.tibco.tgdb.test.admin.attribute;

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
 * Attribute-related test cases from Admin
 * After each admin operation (create, update ...) the server and admin are re-started
 * @author Serge
 *
 */
public class AttributeTests {

	private static TGServer tgServer;
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");

	final private String attrCreationSuccessMsg = "Successfully created attrdesc on server.";
	final private String attrCreationDuplicateMsg = "A duplicate entity exists on the server.";
	final private String attrDescribeMsg = "???";  // Not yet defined in the product
	final private String attrDropMsg = "???"; // Not yet defined in the product
	final private String attrShowMsg = "attribute descriptor(s) returned"; 
	
	int expectedNbAttr = 0;
	
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
			tgServer.init(initFile.getAbsolutePath(), true, 100000);
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
		//File backLogFile = new File(logFile + ".attr");
		//Files.copy(logFile.toPath(), backLogFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
	}

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
	/**
	 * testCreateAttributes - Create various attributes via TG Admin
	 * @throws Exception
	 */
	@Test(description = "Create various attributes via TG Admin")
	public void testCreateAttributes() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/CreateAttr.data", tgWorkingDir + "/CreateAttr.data");
		
		// Create attributes via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminCreateAttr.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Get expected number of attr creations
		expectedNbAttr = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("create"))
				expectedNbAttr++;
		}
		br.close();
		
		// Check actual attr creation
		Scanner scanner = new Scanner(console);
		int attrCreation = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(attrCreationSuccessMsg))
				attrCreation++;
		}
		scanner.close();
		Assert.assertEquals(attrCreation, expectedNbAttr, "Attribute creation does not match -");
	}
	
	/**
	 * testDuplicateAttributes - Re-create the same attributes. Should get duplicate message in Admin
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testCreateAttributes" },
		  description = "Re-create same attributes via TG Admin and check for duplicates")
	public void testDuplicateAttributes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/CreateAttr.data", tgWorkingDir + "/CreateAttr.data");
		
		// Re-Create attribute via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminDuplicateAttr.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		// System.out.println(console);
		
		// Check attr duplication 
		Scanner scanner = new Scanner(console);
		int attrDuplicate = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(attrCreationDuplicateMsg))
				attrDuplicate++;
		}
		scanner.close();
		Assert.assertEquals(attrDuplicate, expectedNbAttr, "Attribute duplication does not match -");
	}
	
	/**
	 * testShowAttributes - Show attributes previously created
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testDuplicateAttributes" },
		  description = "Show attributes in TG Admin")
	public void testShowAttributes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/ShowAttr.data", tgWorkingDir + "/ShowAttr.data");
		
		// Show attr via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminShowAttr.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);
		
		// Check show attr
		Assert.assertTrue(console.contains(expectedNbAttr + " " + attrShowMsg), "Expected " + expectedNbAttr + " " + attrShowMsg + " but did not get that -"); 
	}
	
	/**
	 * testDescribeAttributes - Describe attributes previously created
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testShowAttributes" },
		  description = "Describe attributes in TG Admin")
	public void testDescribeAttributes() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DescribeAttr.data", tgWorkingDir + "/DescribeAttr.data");
		
		// Describe attr via Admin
		String console = TGAdmin.invoke(tgHome, null, null, null, tgWorkingDir + "/adminDescribeAttr.log",
				cmdFile.getAbsolutePath(), 10000);
		// System.out.println(console);
		
		// Get expected number of attr description
		expectedNbAttr = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("describe attrdesc"))
				expectedNbAttr++;
			}
		br.close();
		
		// Check describe attr
		Scanner scanner = new Scanner(console);
		int attrDescribe = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(attrDescribeMsg))
				attrDescribe++;
		}
		scanner.close();
		Assert.assertEquals(attrDescribe, expectedNbAttr, "Describe attributes does not match -");
	}
	*/
	
	/**
	 * testDropAttributes - Drop attributes previously created via TG Admin
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testShowAttributes" },
		  description = "Drop attributes previously created via TG Admin")
	public void testDropAttributes() throws Exception {
		
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DropAttr.data", tgWorkingDir + "/DropAttr.data");
		
		// Drop attr via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminDropAttr.log",
				cmdFile.getAbsolutePath(), 10000);
		//System.out.println(console);
		
		// Get expected number of attr deleted
		expectedNbAttr = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("drop attrdesc"))
				expectedNbAttr++;
		}
		br.close();
		
		// Check actual attr deletion
		Scanner scanner = new Scanner(console);
		int attrDrop = 0;
		while(scanner.hasNextLine()) {
			if (scanner.nextLine().contains(attrDropMsg))
				attrDrop++;
		}
		scanner.close();
		Assert.assertEquals(attrDrop, expectedNbAttr, "Attribute deletion does not match -");
	}
	*/
}
