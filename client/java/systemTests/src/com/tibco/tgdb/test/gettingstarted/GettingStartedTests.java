package com.tibco.tgdb.test.gettingstarted;

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.PrintStream;

import org.testng.Assert;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.Test;

import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.utils.ClasspathResource;

import com.tibco.tgdb.exception.TGException;

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

public class GettingStartedTests {

	private static TGServer tgServer;
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");
	
	private static String buildGraphSuccessMsg = "House of Bonaparte graph completed successfully";
	private static String searchGraphMember = "Napoleon Bonaparte"; 
	private static String searchGraphChild = "Francois Bonaparte";
	private static String searchGraphSuccessMsg = "House member '" + searchGraphMember + "' found";
	private static String searchGraphFailureMsg = "House member '" + searchGraphMember + "' not found";
	private static String updateGraphHead = "true";
	private static String updateGraphBorn = "1990";
	private static String updateGraphReignEnd = "31 Jan 2016";
	private static String updateGraphCrown = "Napoleon XVIII";
	private static String updateGraphSuccessMsg = "House member '" + searchGraphMember + "' updated successfully";
	private static String deleteGraphSuccessMsg = "House member '" + searchGraphMember + "' deleted successfully";
	private static String queryGraphMember = "Napoleon VIII Jean-Christophe";
	private static String queryGraphSuccessMsg = "House member '" + queryGraphMember + "' found";
	
	/**
	 * Init TG server before test suite
	 * @throws Exception
	 */
	@BeforeSuite(description = "Init TG Server")
	public void initServer() throws Exception  {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/inithousedb.conf", tgWorkingDir + "/inithousedb.conf");
		File confFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/tgdb.conf", tgWorkingDir + "/tgdb.conf");
		tgServer = new TGServer(tgHome);
		tgServer.setConfigFile(confFile);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 60000);
		}
		catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
		//File confFile = ClasspathResource.getResourceAsFile(
		//		this.getClass().getPackage().getName().replace('.', '/') + "/tgdb.conf", tgWorkingDir + "/tgdb.conf");
		//tgServer.setConfigFile(confFile);
		//tgServer.start(10000);
		System.out.println(tgServer.getBanner());
	}
	
	/**
	 * Start TG server before each test method
	 * @throws Exception
	 */
	@BeforeMethod
	public void startServer() throws Exception {
		tgServer.start(10000);
	}

	/**
	 * Stop TG server after each test method
	 * @throws Exception
	 */
	@AfterMethod
	public void stopServer() throws Exception {
		Thread.sleep(1000); // avoid corrupted shm file if stopping too early (1.1.1 release)
		TGAdmin.stopServer(tgServer, tgServer.getNetListeners()[0].getName(), null, null, 60000);
	}
	
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/	
	
	
  /**
	 * testBuildGraph - Run the BuildGraph program from Getting Started Guide
	 * @throws Exception
	 */
	@Test(description = "Run the BuildGraph program from Getting Started Guide",
		  timeOut = 30000)
	public void testBuildGraph() throws Exception {
		PrintStream old = System.out; // Save old System.out
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		PrintStream ps = new PrintStream(baos);
		System.setOut(ps);
		try {
			BuildGraph.main(null);
			if (!baos.toString().contains(buildGraphSuccessMsg))
				Assert.fail("Could not find output '" + buildGraphSuccessMsg + "'");
		}
		finally {
			System.out.flush();
		    System.setOut(old); // Restore regular System.out
		    System.out.println(baos.toString());
		}
	}
	
	/**
	 * testRebuildGraph - Run the BuildGraph program again. Expect TGException UniqueContraintViolation
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testBuildGraph" },
		  expectedExceptions = {TGException.class},
		  description = "Run the BuildGraph program again - Expect TGException UniqueContraintViolation",
		  timeOut = 30000)
	public void testRebuildGraph() throws Exception {
		BuildGraph.main(null);
	}
	
	/**
	 * testSearch1Graph - Run the SearchGraph program from Getting Started Guide
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testBuildGraph" },
		  description = "Run the SearchGraph program from Getting Started Guide")
	public void testSearch1Graph() throws Exception {
		PrintStream old = System.out; // Save old System.out
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		PrintStream ps = new PrintStream(baos);
		System.setOut(ps);
		try {
			SearchGraph.main(new String[]{"-memberName",searchGraphMember});
			if (!baos.toString().contains(searchGraphSuccessMsg))
				Assert.fail("Could not find output '" + searchGraphSuccessMsg + "'");
			if (!baos.toString().contains("child: " + searchGraphChild))
				Assert.fail("Could not find output 'child: " + searchGraphChild + "'");
		}
		finally {
			System.out.flush();
		    System.setOut(old); // Restore regular System.out
		    System.out.println(baos.toString());
		}
	}
	
	/**
	 * testUpdateGraph - Run the UpdateGraph program from Getting Started Guide
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testBuildGraph" },
		  description = "Run the UpdateGraph program from Getting Started Guide",
		  timeOut = 30000)
	public void testUpdateGraph() throws Exception {
		PrintStream old = System.out; // Save old System.out
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		PrintStream ps = new PrintStream(baos);
		System.setOut(ps);
		try {
			UpdateGraph.main(new String[]{"-memberName",searchGraphMember,"-crownName",updateGraphCrown,"-houseHead",updateGraphHead,"-yearBorn",updateGraphBorn,"-reignEnd",updateGraphReignEnd});
			if (!baos.toString().contains(updateGraphSuccessMsg))
				Assert.fail("Could not find output '" + updateGraphSuccessMsg + "'");
		}
		finally {
			System.out.flush();
		    System.setOut(old); // Restore regular System.out
		    System.out.println(baos.toString());
		}
	}
	
	/**
	 * testSearch2Graph - Run the SearchGraph program again to validate the node update
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testUpdateGraph" },
		  description = "Run the SearchGraph program again to validate the node update",
		  timeOut = 30000)
	public void testSearch2Graph() throws Exception {
		PrintStream old = System.out; // Save old System.out
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		PrintStream ps = new PrintStream(baos);
		System.setOut(ps);
		try {
			SearchGraph.main(new String[]{"-memberName",searchGraphMember});
			if (!baos.toString().contains(searchGraphSuccessMsg))
				Assert.fail("Could not find output '" + searchGraphSuccessMsg + "'");
			if (!baos.toString().contains("houseHead: "+updateGraphHead))
				Assert.fail("Could not find output 'houseHead: " + updateGraphHead + "'");
			if (!baos.toString().contains("yearBorn: "+updateGraphBorn))
				Assert.fail("Could not find output 'yearBorn: " + updateGraphBorn + "'");
			if (!baos.toString().contains("reignEnd: "+updateGraphReignEnd))
				Assert.fail("Could not find output 'reignEnd: " + updateGraphReignEnd + "'");
			if (!baos.toString().contains("crownName: "+updateGraphCrown))
				Assert.fail("Could not find output 'crownName: " + updateGraphCrown + "'");
		}
		finally {
			System.out.flush();
		    System.setOut(old); // Restore regular System.out
		    System.out.println(baos.toString());
		}
	}
	
	/**
	 * testDeleteGraph - Run the DeleteGraph program from Getting Started Guide
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testSearch2Graph" },
		  description = "Run the DeleteGraph program from Getting Started Guide",
		  timeOut = 30000)
	public void testDeleteGraph() throws Exception {
		PrintStream old = System.out; // Save old System.out
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		PrintStream ps = new PrintStream(baos);
		System.setOut(ps);
		try {
			DeleteGraph.main(new String[]{"-memberName",searchGraphMember});
			if (!baos.toString().contains(deleteGraphSuccessMsg))
				Assert.fail("Could not find output '" + deleteGraphSuccessMsg + "'");
		}
		finally {
			System.out.flush();
		    System.setOut(old); // Restore regular System.out
		    System.out.println(baos.toString());
		}
	}
	
	/**
	 * testSearch3Graph - Run the SearchGraph program again to validate the node deletion
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testDeleteGraph" },
		  description = "Run the SearchGraph program again to validate the node deletion",
		  timeOut = 30000)
	public void testSearch3Graph() throws Exception {
		PrintStream old = System.out; // Save old System.out
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		PrintStream ps = new PrintStream(baos);
		System.setOut(ps);
		try {
			SearchGraph.main(new String[]{"-memberName",searchGraphMember});
			if (!baos.toString().contains(searchGraphFailureMsg))
				Assert.fail("House member '" + searchGraphMember + "' not found");
		}
		finally {
			System.out.flush();
		    System.setOut(old); // Restore regular System.out
		    System.out.println(baos.toString());
		}
	}
	
	/**
	 * testQueryGraph - Run the QueryGraph program from Getting Started Guide
	 * @throws Exception
	 */
	@Test(dependsOnMethods = { "testBuildGraph" },
		  description = "Run the QueryGraph program from Getting Started Guide",
		  timeOut = 30000)
	public void testQueryGraph() throws Exception {
		PrintStream old = System.out; // Save old System.out
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		PrintStream ps = new PrintStream(baos);
		System.setOut(ps);
		try {
			QueryGraph.main(new String[]{"-startyear", "1900", "-endyear", "2000"});
			if (!baos.toString().contains(queryGraphSuccessMsg))
				Assert.fail("Could not find output '" + queryGraphSuccessMsg + "'");
		}
		finally {
			System.out.flush();
		    System.setOut(old); // Restore regular System.out
		    System.out.println(baos.toString());
		}
	}
}
