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
 * CRUD tests for time data type index
 */
@Ignore
public class TimeIndexTests {

	private static TGServer tgServer;
	private static String tgUrl;
	private static String tgUser = "scott";
	private static String tgPwd = "scott";
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");	
	
	Object[][] data;
	
	public TimeIndexTests() throws IOException, EvalError {
		this.data = this.getTimeData();
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
	 * testCreateTimeData - Insert nodes with time index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with time index")
	public void testCreateTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeTimeIdxType = gmd.getNodeType("nodeTimeIdx");
		if (nodeTimeIdxType == null)
			throw new Exception("Node type not found");
		
		//Object[][] data = this.getTimeData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			//System.out.println("CREATE ATTR:" + data[i][0]);
			TGNode node = gof.createNode(nodeTimeIdxType);
			node.setAttribute("timeAttr", data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
			/*if (i>0) {
				TGEdge edge = gof.createEdge(nodes.get(i-1), nodes.get(i), TGEdge.DirectionType.UnDirected);
				edge.setAttribute("timeAttr", data[i-1][0]);
				conn.insertEntity(edge);
			}*/
		}
		// complete the circle - FIX TGDB-176
		//TGEdge edge = gof.createEdge(nodes.get(booleanData.length-1), nodes.get(0), TGEdge.DirectionType.UnDirected);
		//edge.setAttribute("timeAttr2", booleanData[booleanData.length-1][0]);
		//conn.insertEntity(edge);
		conn.commit();
		//Assert.assertEquals(conn.commit().count(),2*booleanData.length,"Expected " + booleanData.length + " nodes + " + (booleanData.length-1) + " edges inserts -");
	
		conn.disconnect();
	}
	
	/**
	 * testReadTimeData - Retrieve nodes with time index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with time index",
		  dependsOnMethods = { "testCreateTimeData" })
	public void testReadTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			System.out.println("READ ATTR:" + entity.getAttribute("timeAttr").getValue());
			// Assert on Node attribute
			// Assert only on time (HOUR, MIN, SEC, MILLISEC) since the Date part that comes back from DB is junk
			Calendar timeAttr = (Calendar)entity.getAttribute("timeAttr").getValue();
			if (timeAttr != null) {
				Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) data[i][0]).get(Calendar.HOUR_OF_DAY));
				Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) data[i][0]).get(Calendar.MINUTE));
				Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) data[i][0]).get(Calendar.SECOND));
				Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) data[i][0]).get(Calendar.MILLISECOND));
			}
			else { // Attribute value is Null. Make sure the original value was Null too
				Assert.assertEquals(timeAttr, data[i][0]);
			}
			/*for (TGEdge edge : ((TGNode)entity).getEdges()) {
				if (edge.getVertices()[0].equals(entity))  {
					// Assert on Edge attribute
					// Assert only on time (HOUR, MIN, SEC, MILLISEC) since the Date part that comes back from DB is junk
					timeAttr = (Calendar)edge.getAttribute("timeAttr").getValue();
					Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) data[i][0]).get(Calendar.HOUR_OF_DAY));
					Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) data[i][0]).get(Calendar.MINUTE));
					Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) data[i][0]).get(Calendar.SECOND));
					Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) data[i][0]).get(Calendar.MILLISECOND));
				}
			}*/
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateTimeData - Update time index
	 * @throws Exception
	 */
	
	@Test(description = "Update time index",
		  dependsOnMethods = { "testReadTimeData" })
	public void testUpdateTimeData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #" + i + " was not retrieved");
			}
			//System.out.println("UPDATE ATTR:" + data[i][1] + " - Length:" + ((Time) data[i][1]).length());
			entity.setAttribute("timeAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedTimeData - Retrieve nodes with updated time index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated time index",
		  dependsOnMethods = { "testUpdateTimeData" })
	public void testReadUpdatedTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			//System.out.println("READ UPDATED ATTR:" + entity.getAttribute("timeAttr").getValue());
			// Assert on Node attribute
			// Assert only on time (HOUR, MIN, SEC, MILLISEC) since the Date part that comes back from DB is junk
			// Assert.assertFalse(entity.getAttribute("timeAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Calendar timeAttr = (Calendar)entity.getAttribute("timeAttr").getValue();
			if (timeAttr != null) {
				Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) data[i][1]).get(Calendar.HOUR_OF_DAY));
				Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) data[i][1]).get(Calendar.MINUTE));
				Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) data[i][1]).get(Calendar.SECOND));
				Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) data[i][1]).get(Calendar.MILLISECOND));
			}
			else {// Attribute value is Null. Make sure the original value was Null too
				Assert.assertEquals(timeAttr, data[i][1]);
			}
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteTimeData - Delete time index
	 * @throws Exception
	 */
	
	@Test(description = "Delete time index",
		  dependsOnMethods = { "testReadUpdatedTimeData" })
	public void testDeleteTimeData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("timeAttr", null); // delete the time value by setting it up to null
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedTimeData - Retrieve nodes with deleted time index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted time index",
		  dependsOnMethods = { "testDeleteTimeData" })
	public void testReadDeletedTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
 
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("timeAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of time data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "TimeData")
	public Object[][] getTimeData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/time.data"));
		return data;
	}
}
