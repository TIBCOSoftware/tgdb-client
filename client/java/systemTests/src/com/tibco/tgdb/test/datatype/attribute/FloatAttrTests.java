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
 * CRUD tests for float data type attribute
 */
public class FloatAttrTests extends LifecycleServer {
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateFloatData - Insert nodes and edge with float attribute
	 * @throws Exception
	 */
	@Test(description = "Insert nodes and edge with float attribute")
	public void testCreateFloatData() throws Exception {
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
		
		Object[][] data = this.getFloatData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			System.out.println("CREATE FLOAT = " + data[i][0]);
			TGNode node = gof.createNode(nodeAllAttrsType);
			node.setAttribute("floatAttr", data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
			if (i>0) {
				TGEdge edge = gof.createEdge(nodes.get(i-1), nodes.get(i), TGEdge.DirectionType.UnDirected);
				edge.setAttribute("floatAttr", data[i-1][0]);
				conn.insertEntity(edge);
			}
		}
		TGEdge edge = gof.createEdge(nodes.get(data.length-1), nodes.get(0), TGEdge.DirectionType.UnDirected);
		edge.setAttribute("floatAttr", data[data.length-1][0]);
		conn.insertEntity(edge);
		conn.commit();
	
		conn.disconnect();
	}
	
	/**
	 * testReadFloatData - Retrieve nodes and edge with float attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes and edge with float attribute",
		  dependsOnMethods = { "testCreateFloatData" })
	public void testReadFloatData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
  			TGEntity entity = conn.getEntity(tgKey, null);
  			if (entity == null) {
  				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
  			}
  			// Assert on Node attribute
  			System.out.println("READ FLOAT = " + entity.getAttribute("floatAttr").getValue());
  			Assert.assertEquals(entity.getAttribute("floatAttr").getValue(), data[i][0]);
  			for (TGEdge edge : ((TGNode)entity).getEdges()) {
  				if (edge.getVertices()[0].equals(entity))  {
  					// Assert on Edge attribute
  					Assert.assertEquals(edge.getAttribute("floatAttr").getValue(), data[i][0]);
  				}
  			}
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateFloatData - Update float attribute
	 * @throws Exception
	 */
	
	@Test(description = "Update float attribute",
		  dependsOnMethods = { "testReadFloatData" })
	public void testUpdateFloatData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("floatAttr", data[i][1]); // 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedFloatData - Retrieve nodes with updated float attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated float attribute",
		  dependsOnMethods = { "testUpdateFloatData" })
	public void testReadUpdatedFloatData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			//System.out.println("READ UPDATED FLOAT = " + entity.getAttribute("floatAttr").getValue());
			//Assert.assertFalse(entity.getAttribute("floatAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("floatAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteFloatData - Set float attribute to null on nodes
	 * @throws Exception
	 */
	
	@Test(description = "Set float attribute to null on nodes",
		  dependsOnMethods = { "testReadUpdatedFloatData" })
	public void testDeleteFloatData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("floatAttr", null); // delete the float value by setting it up to null
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedFloatData - Retrieve nodes with null float attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with null float attribute",
		  dependsOnMethods = { "testDeleteFloatData" })
	public void testReadDeletedFloatData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
 
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("floatAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of float data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "FloatData")
	public Object[][] getFloatData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/float.data"));
		return data;
	}
	
}
