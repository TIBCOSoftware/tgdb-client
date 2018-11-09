package com.tibco.tgdb.test.datatype.pkey;

import java.io.File;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import org.testng.Assert;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeClass;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.DataProvider;
import org.testng.annotations.Test;

import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.utils.ClasspathResource;
import com.tibco.tgdb.test.utils.PipedData;

import bsh.EvalError;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

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
 * CRUD tests for boolean data type primary key
 */
public class BooleanPKeyTests {

	private static TGServer tgServer;
	private static String tgUrl;
	private static String tgUser = "scott";
	private static String tgPwd = "scott";
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");	
	
	/**
	 * Init TG server before test suite
	 * @throws Exception
	 */
	@BeforeClass(description = "Init TG Server")
	public void initServer() throws Exception  {
		TGServer.killAll(); // Clean up everything first
		File initFile = ClasspathResource.getResourceAsFile(this.getClass().getPackage().getName().replace('.', '/') + "/initdb.conf", tgWorkingDir + "/inidb.conf");
		File confFile = ClasspathResource.getResourceAsFile(
				this.getClass().getPackage().getName().replace('.', '/') + "/tgdb.conf", tgWorkingDir + "/tgdb.conf");
		tgServer = new TGServer(tgHome);
		tgServer.setConfigFile(confFile);
		try {
			tgServer.init(initFile.getAbsolutePath(), true, 60000);
			System.out.println(tgServer.getBanner());
		}
		catch (TGInitException ie) {
			System.out.println(ie.getOutput());
			throw ie;
		}
		tgUrl = "tcp://" + tgServer.getNetListeners()[0].getHost() + ":" + tgServer.getNetListeners()[0].getPort();
		//File confFile = ClasspathResource.getResourceAsFile(
		//		this.getClass().getPackage().getName().replace('.', '/') + "/tgdb.conf", tgWorkingDir + "/tgdb.conf");
		//tgServer.setConfigFile(confFile);
		//tgServer.start(10000);
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
		TGAdmin.stopServer(tgServer, tgServer.getNetListeners()[0].getName(), null, null, 60000);
	}
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testDefinePKey - Insert nodes with boolean primary key
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with boolean primary key")
	public void testDefinePKey() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeAllAttrsType = gmd.getNodeType("nodeBooleanKey");
		if (nodeAllAttrsType == null)
			throw new Exception("Node type not found");
		
		Assert.assertEquals(nodeAllAttrsType.getPKeyAttributeDescriptors()[0].getName(), "boolAttr");
		Object[][] data = this.getBooleanData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodeAllAttrsType);
			node.setAttribute("boolAttr",data[i][0]);
			//node.setAttribute("boolAttr", i);
			nodes.add(node);
			conn.insertEntity(node);
		}
		
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testRetrievePKey - Retrieve nodes with boolean primary key
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with string primary key",
			  dependsOnMethods = { "testDefinePKey" })
		public void testRetrievePKey() throws Exception {
			TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
			
			conn.connect();
			
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null) {
				throw new org.testng.TestException("TG object factory is null");
			}
			
			conn.getGraphMetadata(true);
			TGKey tgKey = gof.createCompositeKey("nodeBooleanKey");
			
			Object[][] data = this.getBooleanData();
			for (int i=0; i<data.length; i++) {
				tgKey.setAttribute("boolAttr", data[i][0]);
				TGEntity entity = conn.getEntity(tgKey, null);
				System.out.println(data[i][0]);
				if (entity == null) {
					throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
				}
				//System.out.println("READ ATTR:" + entity.getAttribute("boolAttr").getAsString());
				// Assert on Node attribute
				Assert.assertEquals(entity.getAttribute("boolAttr").getValue(),data[i][0]);
			}
			conn.disconnect();
		}
		
	/**
	 * testUpdatePKey - Update boolean primary key
	 * @throws Exception
	 */
	
	@Test(description = "Update boolean primary key",
		  dependsOnMethods = { "testRetrievePKey" },
		  enabled = false)
	public void testUpdatePKey() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeBooleanKey");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("boolAttr", data[i][0]);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("boolAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testRetrieveUpdatedPKey - Retrieve nodes with updated boolean primary key
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated boolean primary key",
		  dependsOnMethods = { "testUpdatePKey" },
		  enabled = false)
	public void testRetrieveUpdatedPKey() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeBooleanKey");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("boolAttr", data[i][0]);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			Assert.assertFalse(entity.getAttribute("boolAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("boolAttr").getAsBoolean(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeletePKey - Delete boolean primary key
	 * @throws Exception
	 */
	
	@Test(description = "Delete boolean primary key",
		  dependsOnMethods = { "testRetrieveUpdatedPKey" },
		  enabled = false)
	public void testDeletePKey() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeBooleanKey");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("boolAttr", (Boolean)data[i][0]);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("boolAttr", null); // delete the value by setting it up to null
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testRetrieveDeletedPKey - Retrieve nodes with updated boolean primary key
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted boolean primary key",
		  dependsOnMethods = { "testDeletePKey" },
		  enabled = false)
	public void testRetrieveDeletedPKey() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeBooleanKey");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("boolAttr", data[i][0]);
			TGEntity entity = conn.getEntity(tgKey, null); 
 
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("boolAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}

	
	/**
	 * Provide a set of boolean data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "BoolData")
	public Object[][] getBooleanData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/bool.data"));
		return data;
	}	
}
