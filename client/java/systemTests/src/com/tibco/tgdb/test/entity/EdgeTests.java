package com.tibco.tgdb.test.entity;
import java.io.File;
	import java.io.IOException;
import java.text.SimpleDateFormat;
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
import com.tibco.tgdb.model.TGAttribute;
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

import bsh.EvalError;

import java.util.Calendar;

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


public class EdgeTests
{ 
	private static TGServer tgServer;
	private static String tgUrl;
	private static String tgUser = "scott";
	private static String tgPwd = "scott";
	private static String tgHome = System.getProperty("TGDB_HOME");
	private static String tgWorkingDir = System.getProperty("TGDB_WORKING", tgHome + "/test");	
	
	Object[][] data;
	
	public EdgeTests() throws IOException, EvalError {
		this.data = this.getEdgeData();
	}
	
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
	public void testCreateEdgeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType edgeAllAttrsType = gmd.getNodeType("nodeAllAttrs");
		if (edgeAllAttrsType == null)
			throw new Exception("Node type not found");
		
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(edgeAllAttrsType);
			
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
			if (i>0) {
			TGEdge edge = gof.createEdge(nodes.get(i-1), nodes.get(i), TGEdge.DirectionType.UnDirected);
				edge.setAttribute("boolAttr", data[i-1][1]);
				edge.setAttribute("intAttr", data[i-1][2]);
				edge.setAttribute("charAttr", data[i-1][3]);
				edge.setAttribute("byteAttr", data[i-1][4]);
				edge.setAttribute("longAttr", data[i-1][5]);
				edge.setAttribute("stringAttr", data[i-1][6]);
				edge.setAttribute("shortAttr", data[i-1][7]);
				edge.setAttribute("floatAttr", data[i-1][8]);
				edge.setAttribute("doubleAttr", data[i-1][9]);
				edge.setAttribute("dateAttr", data[i-1][10]);
			edge.setAttribute("timeAttr", data[i-1][11]);
			edge.setAttribute("timestampAttr", data[i-1][12]);
				edge.setAttribute("numberAttr", data[i-1][13]);
				
				conn.insertEntity(edge);
				//System.out.println("in create edge");
			}
		}
		
		TGEdge edge = gof.createEdge(nodes.get(data.length-1), nodes.get(0), TGEdge.DirectionType.UnDirected);
		edge.setAttribute("boolAttr", data[data.length-1][1]);
		edge.setAttribute("intAttr", data[data.length-1][2]);
		edge.setAttribute("charAttr", data[data.length-1][3]);
		edge.setAttribute("byteAttr", data[data.length-1][4]);
		edge.setAttribute("longAttr", data[data.length-1][5]);
		edge.setAttribute("stringAttr", data[data.length-1][6]);
		edge.setAttribute("shortAttr", data[data.length-1][7]);
		edge.setAttribute("floatAttr", data[data.length-1][8]);
		edge.setAttribute("doubleAttr", data[data.length-1][9]);
		edge.setAttribute("dateAttr", data[data.length-1][10]);
		edge.setAttribute("timeAttr", data[data.length-1][11]);
		edge.setAttribute("timestampAttr", data[data.length-1][12]);
		edge.setAttribute("numberAttr", data[data.length-1][13]);
		conn.insertEntity(edge);
		
		conn.commit();
		//Assert.assertEquals(conn.commit().count(),2*booleanData.length,"Expected " + booleanData.length + " nodes + " + (booleanData.length-1) + " edges inserts -");
	
		conn.disconnect();
	}
	
	@Test(description = "Retrieve edge with Node attribute",
			  dependsOnMethods = { "testCreateEdgeData" })
		public void testReadEdgeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getEdgeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
				//System.out.println("Edge count: " + ((TGNode)entity).getEdges().size());
      			for (TGEdge edge : ((TGNode)entity).getEdges()) {
      				if (edge.getVertices()[0].equals((TGNode)entity)) {
      					System.out.println("Edge attr " + edge.getAttribute("boolAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("intAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("charAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("byteAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("longAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("stringAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("shortAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("floatAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("doubleAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("dateAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("timeAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("timestampAttr").getValue());
      					System.out.println("Edge attr " + edge.getAttribute("numberAttr").getValue());
      					
      					Assert.assertEquals(edge.getAttribute("boolAttr").getValue(), data[i][1]);
      					Assert.assertEquals(edge.getAttribute("intAttr").getValue(), data[i][2]);
      					Assert.assertEquals(edge.getAttribute("charAttr").getValue(), data[i][3]);
      					Assert.assertEquals(edge.getAttribute("byteAttr").getValue(), data[i][4]);
      					Assert.assertEquals(edge.getAttribute("longAttr").getValue(), data[i][5]);
      					Assert.assertEquals(edge.getAttribute("stringAttr").getValue(), data[i][6]);
      					Assert.assertEquals(edge.getAttribute("shortAttr").getValue(), data[i][7]);
      					Assert.assertEquals(edge.getAttribute("floatAttr").getValue(), data[i][8]);
      					Assert.assertEquals(edge.getAttribute("doubleAttr").getValue(), data[i][9]);
      					Assert.assertEquals(edge.getAttribute("dateAttr").getValue(), data[i][10]);
      					Calendar timeAttr = (Calendar)edge.getAttribute("timeAttr").getValue();
      					Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) data[i][11]).get(Calendar.HOUR_OF_DAY));
      					Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) data[i][11]).get(Calendar.MINUTE));
      					Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) data[i][11]).get(Calendar.SECOND));
      					Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) data[i][11]).get(Calendar.MILLISECOND));
      					Assert.assertEquals(edge.getAttribute("timestampAttr").getValue(), data[i][12]);
      					Assert.assertEquals(edge.getAttribute("numberAttr").getValue(), data[i][13]);
      				}
      			}
		}	
		conn.disconnect();
	}
	@Test(description = "Update Edge attribute",
			  dependsOnMethods = { "testReadEdgeData" })
		public void testUpdateEdgeData() throws Exception {
	TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
			
			conn.connect();
			
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null) {
				throw new org.testng.TestException("TG object factory is null");
			}
			
			conn.getGraphMetadata(true);
			TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
			
			Object[][] updatedata = this.getUpdatedEdgeData();
			for (int i=0; i<updatedata.length; i++) {
				tgKey.setAttribute("key", i);
				TGEntity entity = conn.getEntity(tgKey, null);
				if (entity == null) {
					throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
				}
				
					
				entity.setAttribute("boolAttr", updatedata[i][1]); 
				entity.setAttribute("intAttr", updatedata[i][2]); 
				entity.setAttribute("charAttr", updatedata[i][3]); 
				entity.setAttribute("byteAttr", updatedata[i][4]); 
				entity.setAttribute("longAttr", updatedata[i][5]);
				entity.setAttribute("stringAttr", updatedata[i][6]); 
				entity.setAttribute("shortAttr", updatedata[i][7]);
				entity.setAttribute("floatAttr", updatedata[i][8]); 
				entity.setAttribute("doubleAttr", updatedata[i][9]); 

				entity.setAttribute("dateAttr", updatedata[i][10]); 
				entity.setAttribute("timeAttr", updatedata[i][11]); 
				entity.setAttribute("timestampAttr", data[i][12]); 
				entity.setAttribute("numberAttr", updatedata[i][13]); 
				
				
				
				
				
				conn.updateEntity(entity);
				conn.commit();
			}
			conn.disconnect();
		}
	
	/**
	 * testRead2BooleanData - Retrieve nodes with updated boolean attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve edge with updated edge attribute",
		  dependsOnMethods = { "testUpdateEdgeData" })
	public void testReadUpdatedEdgeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] updatedata = this.getUpdatedEdgeData();
		for (int i=0; i<updatedata.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			Assert.assertFalse(entity.getAttribute("boolAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("intAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("charAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
		    Assert.assertFalse(entity.getAttribute("byteAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("longAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("stringAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("shortAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("floatAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("doubleAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("dateAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("timeAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertFalse(entity.getAttribute("numberAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			
			Assert.assertEquals(entity.getAttribute("boolAttr").getAsBoolean(), updatedata[i][1]);
			Assert.assertEquals(entity.getAttribute("intAttr").getAsInt(), updatedata[i][2]);
			Assert.assertEquals(entity.getAttribute("charAttr").getAsChar(), updatedata[i][3]);
			Assert.assertEquals(entity.getAttribute("byteAttr").getAsByte(), updatedata[i][4]);
			Assert.assertEquals(entity.getAttribute("longAttr").getAsLong(), updatedata[i][5]);
			Assert.assertEquals(entity.getAttribute("stringAttr").getAsString(), updatedata[i][6]);
			Assert.assertEquals(entity.getAttribute("shortAttr").getAsShort(), updatedata[i][7]);
			Assert.assertEquals(entity.getAttribute("floatAttr").getAsFloat(), updatedata[i][8]);
			Assert.assertEquals(entity.getAttribute("doubleAttr").getAsDouble(), updatedata[i][9]);
			Assert.assertEquals(entity.getAttribute("dateAttr").getValue(), updatedata[i][10]);
			//Assert.assertEquals(entity.getAttribute("timeAttr").getValue(), updatedata[i][11]);
			Calendar timeAttr = (Calendar)entity.getAttribute("timeAttr").getValue();
			Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) updatedata[i][11]).get(Calendar.HOUR_OF_DAY));
			Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) updatedata[i][11]).get(Calendar.MINUTE));
			Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) updatedata[i][11]).get(Calendar.SECOND));
			Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) updatedata[i][11]).get(Calendar.MILLISECOND));
			Assert.assertEquals(entity.getAttribute("timestampAttr").getValue(), data[i][12]);
			Assert.assertEquals(entity.getAttribute("numberAttr").getValue(), updatedata[i][13]);
		}
		conn.disconnect();
	}
	/**
	 * testDeleteBooleanData - Delete boolean attribute
	 * @throws Exception
	 */
	
	@Test(description = "Delete node attribute",
		  dependsOnMethods = { "testReadUpdatedEdgeData" })
	public void testDeleteEdgeData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getEdgeData();
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
			entity.setAttribute("timestampAttr", null);
			entity.setAttribute("numberAttr", null);
			conn.deleteEntity(entity);
			//conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	/**
	 * testRead3BooleanData - Retrieve nodes with updated boolean attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve edge with deleted edge attribute",
		  dependsOnMethods = { "testDeleteEdgeData" })
	public void testReadDeletedEdgeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getEdgeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
					Assert.assertNull(entity);
		
		}
		conn.disconnect();
	}
	
	
	@DataProvider(name = "EdgeData")
		public Object[][] getEdgeData() throws IOException, EvalError {
			Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/edge.data"));
			return data;
		}
		@DataProvider(name = "EdgeData")
		public Object[][] getUpdatedEdgeData() throws IOException, EvalError {
			Object[][] updatedata =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/updatEdge.data"));
			return updatedata;
		}
		

}

//>>>>>>> branch 'master' of https://git.tibco.com/git/product/sgdb.git
