package com.tibco.tgdb.test.datatype.index;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.List;

import org.testng.Assert;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeClass;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.DataProvider;
import org.testng.annotations.Ignore;
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
 * CRUD tests for timestamp data type index
 */
@Ignore
public class TimestampIndexTests {

	private static TGServer tgServer;
	private static String tgUrl;
	private static String tgUser = "scott";
	private static String tgPwd = "scott";
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");	
	
	Object[][] data;
	
	public TimestampIndexTests() throws IOException, EvalError {
		this.data = this.getTimestampData();
	}
	
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
	 * testCreateTimestampData - Insert nodes with timestamp index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with timestamp index")
	public void testCreateTimestampData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeTimestampIdxType = gmd.getNodeType("nodeTimestampIdx");
		if (nodeTimestampIdxType == null)
			throw new Exception("Node type not found");
		
		//Object[][] data = this.getTimestampData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			//System.out.println("CREATE ATTR:" + data[i][0]);
			TGNode node = gof.createNode(nodeTimestampIdxType);
			node.setAttribute("timestampAttr", data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
			if (i>0) {
				TGEdge edge = gof.createEdge(nodes.get(i-1), nodes.get(i), TGEdge.DirectionType.UnDirected);
				edge.setAttribute("timestampAttr", data[i-1][0]);
				conn.insertEntity(edge);
			}
		}
		// complete the circle - FIX TGDB-176
		//TGEdge edge = gof.createEdge(nodes.get(booleanData.length-1), nodes.get(0), TGEdge.DirectionType.UnDirected);
		//edge.setAttribute("timestampAttr2", booleanData[booleanData.length-1][0]);
		//conn.insertEntity(edge);
		conn.commit();
		//Assert.assertEquals(conn.commit().count(),2*booleanData.length,"Expected " + booleanData.length + " nodes + " + (booleanData.length-1) + " edges inserts -");
	
		conn.disconnect();
	}
	
	/**
	 * testReadTimestampData - Retrieve nodes and edge with timestamp index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes and edge with timestamp index",
		  dependsOnMethods = { "testCreateTimestampData" })
	public void testReadTimestampData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimestampIdx");
		
		//Object[][] data = this.getTimestampData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			//System.out.println("READ ATTR:" + entity.getAttribute("timestampAttr").getValue());
			// Assert on Node attribute
			Assert.assertEquals(entity.getAttribute("timestampAttr").getValue(), data[i][0]);
			/*for (TGEdge edge : ((TGNode)entity).getEdges()) {
				if (edge.getVertices()[0].equals(entity))  {
					// Assert on Edge attribute
					Assert.assertEquals(edge.getAttribute("timestampAttr").getValue(), data[i][0]);
				}
			}*/
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateTimestampData - Update timestamp index
	 * @throws Exception
	 */
	
	@Test(description = "Update timestamp index",
		  dependsOnMethods = { "testReadTimestampData" })
	public void testUpdateTimestampData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimestampIdx");
		
		//Object[][] data = this.getTimestampData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #" + i + " was not retrieved");
			}
			//System.out.println("UPDATE ATTR:" + data[i][1] + " - Length:" + ((Timestamp) data[i][1]).length());
			entity.setAttribute("timestampAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedTimestampData - Retrieve nodes with updated timestamp index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated timestamp index",
		  dependsOnMethods = { "testUpdateTimestampData" })
	public void testReadUpdatedTimestampData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimestampIdx");
		
		//Object[][] data = this.getTimestampData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			//System.out.println("READ UPDATED ATTR:" + entity.getAttribute("timestampAttr").getValue());
			// Assert on Node attribute
			// Assert.assertFalse(entity.getAttribute("timestampAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("timestampAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteTimestampData - Delete timestamp index
	 * @throws Exception
	 */
	
	@Test(description = "Delete timestamp index",
		  dependsOnMethods = { "testReadUpdatedTimestampData" })
	public void testDeleteTimestampData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimestampIdx");
		
		//Object[][] data = this.getTimestampData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("timestampAttr", null); // delete the timestamp value by setting it up to null
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedTimestampData - Retrieve nodes with deleted timestamp index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted timestamp index",
		  dependsOnMethods = { "testDeleteTimestampData" })
	public void testReadDeletedTimestampData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimestampIdx");
		
		//Object[][] data = this.getTimestampData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
 
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("timestampAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of timestamp data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "TimestampData")
	public Object[][] getTimestampData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/timestamp.data"));
		return data;
	}
}
