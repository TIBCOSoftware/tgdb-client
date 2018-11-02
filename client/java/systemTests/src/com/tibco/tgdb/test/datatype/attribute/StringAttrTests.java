package com.tibco.tgdb.test.datatype.attribute;

import java.io.IOException;
import java.nio.charset.Charset;
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
 * CRUD tests for string data type attribute
 */
public class StringAttrTests extends LifecycleServer {

	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateStringData - Insert nodes and edge with string attribute
	 * @throws Exception
	 */
	@Test(description = "Insert nodes and edge with string attribute")
	public void testCreateStringData() throws Exception {
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
		
		Object[][] data = this.getStringData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			//System.out.println("CREATE ATTR:" + data[i][0]);
			TGNode node = gof.createNode(nodeAllAttrsType);
			node.setAttribute("stringAttr", data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
			if (i>0) {
				TGEdge edge = gof.createEdge(nodes.get(i-1), nodes.get(i), TGEdge.DirectionType.UnDirected);
				edge.setAttribute("stringAttr", data[i-1][0]);
				conn.insertEntity(edge);
			}
		}
		TGEdge edge = gof.createEdge(nodes.get(data.length-1), nodes.get(0), TGEdge.DirectionType.UnDirected);
		edge.setAttribute("stringAttr", data[data.length-1][0]);
		conn.insertEntity(edge);
		conn.commit();
	
		conn.disconnect();
	}
	
	/**
	 * testReadStringData - Retrieve nodes and edge with string attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes and edge with string attribute",
		  dependsOnMethods = { "testCreateStringData" })
	public void testReadStringData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			//System.out.println("READ ATTR:" + entity.getAttribute("stringAttr").getAsString());
			// Assert on Node attribute
			Assert.assertEquals(entity.getAttribute("stringAttr").getValue(), (String)data[i][0]);
			for (TGEdge edge : ((TGNode)entity).getEdges()) {
				if (edge.getVertices()[0].equals(entity))  {
					// Assert on Edge attribute
					Assert.assertEquals(edge.getAttribute("stringAttr").getValue(), (String)data[i][0]);
				}
			}
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateStringData - Update string attribute on nodes
	 * @throws Exception
	 */
	
	@Test(description = "Update string attribute on nodes",
		  dependsOnMethods = { "testReadStringData" })
	public void testUpdateStringData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #" + i + " was not retrieved");
			}
			//System.out.println("UPDATE ATTR:" + data[i][1] + " - Length:" + ((String) data[i][1]).length());
			entity.setAttribute("stringAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedStringData - Retrieve nodes with updated string attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated string attribute",
		  dependsOnMethods = { "testUpdateStringData" })
	public void testReadUpdatedStringData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			//System.out.println("READ UPDATED ATTR:" + entity.getAttribute("stringAttr").getAsString());
			// Assert on Node attribute
			// Assert.assertFalse(entity.getAttribute("stringAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("stringAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteStringData - Set string attribute to null on nodes
	 * @throws Exception
	 */
	
	@Test(description = "Set string attribute to null on nodes",
		  dependsOnMethods = { "testReadUpdatedStringData" })
	public void testDeleteStringData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("stringAttr", null); // delete the string value by setting it up to null
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedStringData - Retrieve nodes with null string attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with null string attribute",
		  dependsOnMethods = { "testDeleteStringData" })
	public void testReadDeletedStringData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeAllAttrs");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
 
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("stringAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of string data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "StringData")
	public Object[][] getStringData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/string.data"), Charset.forName("UTF-8"));
		return data;
	}
}
