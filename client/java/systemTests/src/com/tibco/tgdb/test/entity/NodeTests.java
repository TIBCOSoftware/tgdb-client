package com.tibco.tgdb.test.entity;


import java.io.File;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import org.testng.Assert;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.DataProvider;
import org.testng.annotations.Test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.utils.ClasspathResource;
import com.tibco.tgdb.test.utils.PipedData;
import java.util.Calendar;
import java.math.BigDecimal;




import bsh.EvalError;

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

public class NodeTests
{ 
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
	@BeforeSuite(description = "Init TG Server")
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
	 * testCreateBooleanData - Insert nodes and edge with boolean attribute
	 * @throws Exception
	 */
	@Test(description = "Create nodes with all possible datatypes")
	public void testCreateNodeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeAllAttrsType = gmd.getNodeType("nodeAllAttrs");
		if (nodeAllAttrsType == null)
			throw new Exception("Node type not found");
		
		Object[][] data = this.getNodeData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodeAllAttrsType);
			
		    node.setAttribute("boolAttr", data[i][1]);
			node.setAttribute("intAttr", data[i][2]);
			node.setAttribute("charAttr", data[i][3]);
			node.setAttribute("byteAttr", data[i][4]);
			node.setAttribute("longAttr", data[i][5]);
			node.setAttribute("stringAttr", data[i][6]);
			node.setAttribute("shortAttr", data[i][7]);
			node.setAttribute("floatAttr", data[i][8]);
			node.setAttribute("doubleAttr", data[i][9]);
			node.setAttribute("dateAttr", data[i][10]);
			node.setAttribute("timeAttr", data[i][11]);
			node.setAttribute("timestampAttr", data[i][12]);
			node.setAttribute("numberAttr", data[i][13]);

			
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
		
			if (i>0) {
			TGEdge edge = gof.createEdge(nodes.get(i-1), nodes.get(i), TGEdge.DirectionType.UnDirected);
				edge.setAttribute("boolAttr", data[i-1][1]);
				
				conn.insertEntity(edge);
				System.out.println("edge created");
			}
		}
		// complete the circle - FIX TGDB-176
		TGEdge edge = gof.createEdge(nodes.get(data.length-1), nodes.get(0), TGEdge.DirectionType.UnDirected);
		edge.setAttribute("boolAttr", data[data.length-1][1]);
		conn.insertEntity(edge);
		conn.commit();
		//Assert.assertEquals(conn.commit().count(),2*booleanData.length,"Expected " + booleanData.length + " nodes + " + (booleanData.length-1) + " edges inserts -");
	
		conn.disconnect();
	}
	
	@Test(description = "Retrieve nodes with Node attribute",
			  dependsOnMethods = { "testCreateNodeData" })
		public void testReadNodeData() throws Exception {
			TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
			
			conn.connect();
			
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null) {
				throw new org.testng.TestException("TG object factory is null");
			}
			
			conn.getGraphMetadata(true);
			TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
			
			Object[][] data = this.getNodeData();
			for (int i=0; i<data.length; i++) {
				tgKey.setAttribute("key", i);
				TGEntity entity = conn.getEntity(tgKey, null);
				if (entity == null) {
					throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
				}
				System.out.println("READ ATTR :" + data[i][0]);
				// Assert on Node attribute
				
				System.out.println("READ ATTR :" + data[i][11]);
				Assert.assertEquals(entity.getAttribute("boolAttr").getValue(),  data[i][1]);
				Assert.assertEquals(entity.getAttribute("intAttr").getValue(),  data[i][2]);
				Assert.assertEquals(entity.getAttribute("charAttr").getValue(),  data[i][3]);
				Assert.assertEquals(entity.getAttribute("byteAttr").getValue(),  data[i][4]);
				Assert.assertEquals(entity.getAttribute("longAttr").getValue(),  data[i][5]);
				Assert.assertEquals(entity.getAttribute("stringAttr").getValue(),  data[i][6]);
				Assert.assertEquals(entity.getAttribute("shortAttr").getValue(),  data[i][7]);
				Assert.assertEquals(entity.getAttribute("floatAttr").getValue(),  data[i][8]);
				Assert.assertEquals(entity.getAttribute("doubleAttr").getValue(),  data[i][9]);
				Assert.assertEquals(entity.getAttribute("dateAttr").getValue(),  data[i][10]);
				Calendar timeAttr = (Calendar)entity.getAttribute("timeAttr").getValue();
				Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) data[i][11]).get(Calendar.HOUR_OF_DAY));
				Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) data[i][11]).get(Calendar.MINUTE));
				Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) data[i][11]).get(Calendar.SECOND));
				Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) data[i][11]).get(Calendar.MILLISECOND));
				Assert.assertEquals(entity.getAttribute("timestampAttr").getValue(),  data[i][12]);
				Assert.assertEquals(entity.getAttribute("numberAttr").getValue(),  data[i][13]);
				System.out.println("READ ATTR :" + data[i][1]);
				System.out.println("READ ATTR :" + data[i][2]);
				System.out.println("READ ATTR :" + data[i][3]);
				System.out.println("READ ATTR :" + data[i][4]);
				System.out.println("READ ATTR :" + data[i][5]);
				System.out.println("READ ATTR :" + data[i][6]);
				System.out.println("READ ATTR :" + data[i][7]);
				System.out.println("READ ATTR :" + data[i][8]);
				System.out.println("READ ATTR :" + data[i][9]);
				System.out.println("READ ATTR :" + data[i][10]);
				System.out.println("READ ATTR :" + data[i][11]);
				System.out.println("READ ATTR :" + data[i][12]);
				System.out.println("READ ATTR :" + data[i][13]);
				
				
				/*for (TGEdge edge : ((TGNode)entity).getEdges()) {
					if (edge.getVertices()[0].equals(entity))  {
						// Assert on Edge attribute
						Assert.assertEquals(edge.getAttribute("boolAttr").getAsBoolean(), data[i][0]);
					}
				}*/
			}
			conn.disconnect();
		}
	
	
	/**
	 * testUpdateBooleanData - Update boolean attribute
	 * @throws Exception
	 */
	
	@Test(description = "Update node attribute",
		  dependsOnMethods = { "testReadNodeData" })
	public void testUpdateNodeData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] UpdateData = this.getUpdateNodeData();
		for (int i=0; i<UpdateData.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("boolAttr", UpdateData[i][1]); 
			entity.setAttribute("intAttr", UpdateData[i][2]); 
			entity.setAttribute("charAttr", UpdateData[i][3]); 
			entity.setAttribute("byteAttr", UpdateData[i][4]); 
			entity.setAttribute("longAttr", UpdateData[i][5]);
			entity.setAttribute("stringAttr", UpdateData[i][6]); 
			entity.setAttribute("shortAttr", UpdateData[i][7]);
			entity.setAttribute("floatAttr", UpdateData[i][8]); 
			entity.setAttribute("doubleAttr", UpdateData[i][9]); 

			entity.setAttribute("dateAttr", UpdateData[i][10]); 
			entity.setAttribute("timeAttr", UpdateData[i][11]);
			entity.setAttribute("timestampAttr", UpdateData[i][12]); 
			entity.setAttribute("numberAttr", UpdateData[i][13]); 
			
			
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testRead2BooleanData - Retrieve nodes with updated boolean attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated node attribute",
		  dependsOnMethods = { "testUpdateNodeData" })
	public void testReadUpdatedNodeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] UpdateData = this.getUpdateNodeData();
		for (int i=0; i<UpdateData.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
		
			Assert.assertEquals(entity.getAttribute("boolAttr").getAsBoolean(), UpdateData[i][1]);
			Assert.assertEquals(entity.getAttribute("intAttr").getAsInt(), UpdateData[i][2]);
			Assert.assertEquals(entity.getAttribute("charAttr").getAsChar(), UpdateData[i][3]);
			Assert.assertEquals(entity.getAttribute("byteAttr").getAsByte(), UpdateData[i][4]);
			Assert.assertEquals(entity.getAttribute("longAttr").getAsLong(), UpdateData[i][5]);
			Assert.assertEquals(entity.getAttribute("stringAttr").getAsString(), UpdateData[i][6]);
			Assert.assertEquals(entity.getAttribute("shortAttr").getAsShort(), UpdateData[i][7]);
			Assert.assertEquals(entity.getAttribute("floatAttr").getAsFloat(), UpdateData[i][8]);
			Assert.assertEquals(entity.getAttribute("doubleAttr").getAsDouble(), UpdateData[i][9]);
			Assert.assertEquals(entity.getAttribute("dateAttr").getValue(), UpdateData[i][10]);
			Calendar timeAttr = (Calendar)entity.getAttribute("timeAttr").getValue();
			Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) UpdateData[i][11]).get(Calendar.HOUR_OF_DAY));
			Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) UpdateData[i][11]).get(Calendar.MINUTE));
			Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) UpdateData[i][11]).get(Calendar.SECOND));
			Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) UpdateData[i][11]).get(Calendar.MILLISECOND));
			
			Assert.assertEquals(entity.getAttribute("timestampAttr").getValue(), UpdateData[i][12]);
			Assert.assertEquals(entity.getAttribute("numberAttr").getValue(), UpdateData[i][13]);
			System.out.println("READ ATTR :" + UpdateData[i][1]);
			System.out.println("READ ATTR :" + UpdateData[i][2]);
			System.out.println("READ ATTR :" + UpdateData[i][3]);
			System.out.println("READ ATTR :" + UpdateData[i][4]);
			System.out.println("READ ATTR :" + UpdateData[i][5]);
			System.out.println("READ ATTR :" + UpdateData[i][6]);
			System.out.println("READ ATTR :" + UpdateData[i][7]);
			System.out.println("READ ATTR :" + UpdateData[i][8]);
			System.out.println("READ ATTR :" + UpdateData[i][9]);
			//System.out.println("READ ATTR :" + UpdateData[i][10]);
			System.out.println("READ ATTR :" + UpdateData[i][11]);
			System.out.println("READ ATTR :" + UpdateData[i][12]);
			System.out.println("READ ATTR :" + UpdateData[i][13]);
		}
		
		conn.disconnect();
	}
	
	/**
	 * testDeleteBooleanData - Delete boolean attribute
	 * @throws Exception
	 */
	
	@Test(description = "Delete node attribute",
		  dependsOnMethods = { "testReadUpdatedNodeData" })
	public void testDeleteNodeData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
				
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getNodeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("boolAttr", null);// delete the boolean value
			entity.setAttribute("intAttr", null);
			entity.setAttribute("charAttr", null);
			entity.setAttribute("byteAttr", null);
			entity.setAttribute("longAttr", null);
			entity.setAttribute("stringAttr", null);
			entity.setAttribute("shortAttr", null);
			entity.setAttribute("floatAttr", null);
			entity.setAttribute("doubleAttr", null);
			entity.setAttribute("dateAttr", null);
			entity.setAttribute("timeAttr", null);
			entity.setAttribute("timestampAttr",null);
			entity.setAttribute("numberAttr",null);
			
			conn.deleteEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	/**
	 * testRead3BooleanData - Retrieve nodes with updated boolean attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted node attribute",
		  dependsOnMethods = { "testDeleteNodeData" })
	public void testReadDeletedNodeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getNodeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
		
			/*Assert.assertTrue(entity.getAttribute("boolAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("intAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("charAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("byteAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("longAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("stringAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("shortAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("floatAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("doubleAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("dateAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
//			Assert.assertTrue(entity.getAttribute("timeAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
			Assert.assertTrue(entity.getAttribute("timestampAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");*/
			
	Assert.assertNull(entity);
		
		}
		conn.disconnect();
	}
	
	
	
	
	
	@DataProvider(name = "NodeData")
	public Object[][] getNodeData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/node.data"));
		return data;
	}
	
	@DataProvider(name = "NodeData")
	public Object[][] getUpdateNodeData() throws IOException, EvalError {
		Object[][] UpdateData =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/updateNode.data"));
		return UpdateData;
	}
	
	
	
}