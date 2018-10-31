package com.tibco.tgdb.test.admin.index;

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
import org.testng.annotations.Test;

import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.utils.ClasspathResource;

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

public class IndexTests {

	private static TGServer tgServer;
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");

	final private String indexCreationSuccessMsg = "Successfully created index on server.";
	final private String indexCreationDuplicateMsg = "A duplicate entity exists on the server.";
	final private String indexShowMsg = "index(es) returned";
	final private String indexDropMsg = "????"; // Not yet defined in the product
	int expectedNbIndex = 0;

	/**
	 * Init TG server before test suite
	 * 
	 * @throws Exception
	 */
	@BeforeSuite(description = "Init TG Admin")
	public void initServer() throws Exception {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replaceFirst("\\.[a-z]*$", "").replace('.', '/') + "/initdb.conf",
				tgWorkingDir + "/initdb.conf");
		tgServer = new TGServer(tgHome);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 100000);
		} catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
	}

	/**
	 * Start TG server before each test method
	 * 
	 * @throws Exception
	 */
	@BeforeMethod
	public void startServer() throws Exception {
		File confFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replaceFirst("\\.[a-z]*$", "").replace('.', '/') + "/tgdb.conf",
				tgWorkingDir + "/tgdb.conf");
		tgServer.setConfigFile(confFile);
		tgServer.start(10000);
	}

	/**
	 * Kill TG server after each test method
	 * 
	 * @throws Exception
	 */
	@AfterMethod
	public void killServer() throws Exception {
		TGAdmin.stopServer(tgServer, tgServer.getNetListeners()[0].getName(), null, null, 60000);
		// tgServer.kill();
		// Backup log file before moving to next test
		// File logFile = tgServer.getLogFile();
		// File backLogFile = new File(logFile + ".attr");
		// Files.copy(logFile.toPath(), backLogFile.toPath(),
		// StandardCopyOption.REPLACE_EXISTING);
	}

	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/

	/**
	 * testCreateIndex - Create various index via TG Admin
	 * 
	 * @throws Exception
	 */
	@Test(description = "Create various index via TG Admin")
	public void testCreateIndex() throws Exception {

		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/CreateIndex.data",
				tgWorkingDir + "/CreateIndex.data");

		// Create attributes via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminCreateIndex.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);

		// Get expected number of index creations
		expectedNbIndex = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.matches("^create (unique |)index .*"))
				expectedNbIndex++;
		}
		br.close();

		// Check actual index creation
		Scanner scanner = new Scanner(console);
		int indexCreation = 0;
		while (scanner.hasNextLine()) {
			if (scanner.nextLine().contains(indexCreationSuccessMsg))
				indexCreation++;
		}
		//System.out.println(indexCreation);
		scanner.close();
		Assert.assertEquals(indexCreation, expectedNbIndex, "Index creation does not match -");
	}

	/**
	 * testDuplicateIndexTypes - Re-create the same index. Should get duplicate
	 * message in Admin
	 * 
	 * @throws Exception
	 */
	@Test(dependsOnMethods = {
			"testCreateIndex" }, description = "Re-create same index via TG Admin and check for duplicates")
	public void testDuplicateIndex() throws Exception {
		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DuplicateIndex.data",
				tgWorkingDir + "/DuplicateIndex.data");

		// Re-Create index via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminDuplicateIndex.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);
		//System.out.println(console);

		// Check index duplication
		Scanner scanner = new Scanner(console);
		int indexCreation = 0;
		while (scanner.hasNextLine()) {
			if (scanner.nextLine().contains(indexCreationDuplicateMsg))
				indexCreation++;
			// System.out.println("indexcreate number"+indexCreation+"\n");

		}
		scanner.close();
		Assert.assertEquals(indexCreation, expectedNbIndex, "Index duplication does not match -");
	}

	@Test(dependsOnMethods = { "testDuplicateIndex" }, description = "Show index in TG Admin")
	public void testShowIndex() throws Exception {

		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/ShowIndex.data",
				tgWorkingDir + "/ShowIndex.data");

		//System.out.println("in command file");

		// Show index via Admin
		String console = TGAdmin.invoke(tgServer, tgServer.getNetListeners()[0].getName(), tgWorkingDir + "/adminShowIndex.log", null,
				cmdFile.getAbsolutePath(), -1, 10000);

		//System.out.println(console);

		// Check show index - +1 since we have 1 pkey that become indices
		Assert.assertTrue(console.contains((expectedNbIndex+1) + " " + indexShowMsg), "Expected " + expectedNbIndex + "+1=16 " + indexShowMsg + " but did not get that -");
	}

	/**
	 * testDropIndex - Drop index previously created via TG Admin
	 * 
	 * @throws Exception
	 */
	/*
	@Test(dependsOnMethods = { "testShowIndex" }, description = "Drop index previously created via TG Admin")
	public void testDropIndex() throws Exception {

		File cmdFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/DropIndex.data",
				tgWorkingDir + "/DropIndex.data");

		// Drop index via Admin
		String console = TGAdmin.invoke(tgHome, null, null, null, tgWorkingDir + "/adminDropIndex.log",
				cmdFile.getAbsolutePath(), 10000);
		// System.out.println(console);

		// Get expected number of index deleted
		expectedNbIndex = 0;
		BufferedReader br = new BufferedReader(new FileReader(cmdFile));
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.contains("drop index"))
				expectedNbIndex++;
		}
		br.close();

		// Check actual index deletion
		Scanner scanner = new Scanner(console);
		int indexDrop = 0;
		while (scanner.hasNextLine()) {
			if (scanner.nextLine().contains(indexDropMsg))
				indexDrop++;
		}
		scanner.close();
		Assert.assertEquals(indexDrop, expectedNbIndex, "Index deletion does not match -");
	}
	*/

}