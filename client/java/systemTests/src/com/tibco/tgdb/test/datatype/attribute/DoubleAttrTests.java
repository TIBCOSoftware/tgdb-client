package com.tibco.tgdb.test.datatype.attribute;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import org.testng.Assert;
import org.testng.annotations.DataProvider;
import org.testng.annotations.Test;

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
 * CRUD tests for double data type attribute
 */
public class DoubleAttrTests extends LifecycleServer {

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateDoubleData - Insert nodes and edge with double attribute
	 * @throws Exception
	 */
	@Test(description = "Insert nodes and edge with double attribute")
	public void testCreateDoubleData() throws Exception {
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
		
		Object[][] data = this.getDoubleData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			//System.out.println("CREATE DOUBLE = " + data[i][0]);
			TGNode node = gof.createNode(nodeAllAttrsType);
			node.setAttribute("doubleAttr", data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
			if (i>0) {
				TGEdge edge = gof.createEdge(nodes.get(i-1), nodes.get(i), TGEdge.DirectionType.UnDirected);
				edge.setAttribute("doubleAttr", data[i-1][0]);
				conn.insertEntity(edge);
			}
		}
		TGEdge edge = gof.createEdge(nodes.get(data.length-1), nodes.get(0), TGEdge.DirectionType.UnDirected);
		edge.setAttribute("doubleAttr", data[data.length-1][0]);
		conn.insertEntity(edge);
		conn.commit();
	
		conn.disconnect();
	}
	
	/**
	 * testReadDoubleData - Retrieve nodes and edge with double attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes and edge with double attribute",
		  dependsOnMethods = { "testCreateDoubleData" })
	public void testReadDoubleData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getDoubleData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
  			TGEntity entity = conn.getEntity(tgKey, null);
  			if (entity == null) {
  				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
  			}
  			// Assert on Node attribute
  			//System.out.println("READ DOUBLE = " + entity.getAttribute("doubleAttr").getValue());
  			Assert.assertEquals(entity.getAttribute("doubleAttr").getValue(), data[i][0]);
  			for (TGEdge edge : ((TGNode)entity).getEdges()) {
  				if (edge.getVertices()[0].equals(entity))  {
  					// Assert on Edge attribute
  					Assert.assertEquals(edge.getAttribute("doubleAttr").getValue(), data[i][0]);
  				}
  			}
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateDoubleData - Update double attribute
	 * @throws Exception
	 */
	
	@Test(description = "Update double attribute",
		  dependsOnMethods = { "testReadDoubleData" })
	public void testUpdateDoubleData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getDoubleData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("doubleAttr", data[i][1]); // 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedDoubleData - Retrieve nodes with updated double attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated double attribute",
		  dependsOnMethods = { "testUpdateDoubleData" })
	public void testReadUpdatedDoubleData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getDoubleData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			//System.out.println("READ UPDATED DOUBLE = " + entity.getAttribute("doubleAttr").getValue());
			//Assert.assertFalse(entity.getAttribute("doubleAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("doubleAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteDoubleData - Set double attribute to null on nodes
	 * @throws Exception
	 */
	
	@Test(description = "Set double attribute to null on nodes",
		  dependsOnMethods = { "testReadUpdatedDoubleData" })
	public void testDeleteDoubleData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getDoubleData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("doubleAttr", null); // delete the double value by setting it up to null
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedDoubleData - Retrieve nodes with null double attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with null double attribute",
		  dependsOnMethods = { "testDeleteDoubleData" })
	public void testReadDeletedDoubleData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getDoubleData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
 
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("doubleAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of double data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "DoubleData")
	public Object[][] getDoubleData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/double.data"));
		return data;
	}
	
}
