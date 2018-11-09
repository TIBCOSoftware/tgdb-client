package com.tibco.tgdb.test.datatype.index;

import java.io.IOException;
import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.List;

import org.testng.Assert;
import org.testng.annotations.DataProvider;
import org.testng.annotations.Ignore;
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
 * CRUD tests for number data type index
 */
@Ignore
public class NumberIndexTests extends LifecycleServer {

	Object[][] data;
	
	public NumberIndexTests() throws IOException, EvalError {
		this.data = this.getNumberData();
	}
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateNumberData - Insert nodes with number index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with number index")
	public void testCreateNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeNumberIdxType = gmd.getNodeType("nodeNumberIdx");
		if (nodeNumberIdxType == null)
			throw new Exception("Node type not found");
		
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodeNumberIdxType);
			node.setAttribute("numberAttr", data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
		}
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testReadNumberData - Retrieve nodes with number index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with number index",
		  dependsOnMethods = { "testCreateNumberData" })
	public void testReadNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberIdx");
		
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			// System.out.println("READ ATTR :" + data[i][0]);
			// Assert on Node attribute
			if (data[i][0] != null)
				Assert.assertEquals(((BigDecimal) entity.getAttribute("numberAttr").getValue()).compareTo((BigDecimal)data[i][0]), 0, "Actual and Expected BigDecimal are not the same");
			else
				Assert.assertEquals(entity.getAttribute("numberAttr").getValue(), data[i][0]);
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateNumberData - Update number index
	 * @throws Exception
	 */
	
	@Test(description = "Update number index",
		  dependsOnMethods = { "testReadNumberData" },
		  enabled = false)
	public void testUpdateNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberIdx");
		
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("numberAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedNumberData - Retrieve nodes with updated number index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated number index",
		  dependsOnMethods = { "testUpdateNumberData" },
		  enabled = false)
	public void testReadUpdatedNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberIdx");
		
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			if (data[i][1] != null)
				Assert.assertEquals(((BigDecimal) entity.getAttribute("numberAttr").getValue()).compareTo((BigDecimal)data[i][1]), 0, "Actual and Expected BigDecimal are not the same");
			else
				Assert.assertEquals(entity.getAttribute("numberAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteNumberData - Delete number index
	 * @throws Exception
	 */
	
	@Test(description = "Delete number index",
		  dependsOnMethods = { "testReadUpdatedNumberData" },
		  enabled = false)
	public void testDeleteNumberData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberIdx");
		
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("numberAttr", null); // delete the boolean value
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedNumberData - Retrieve nodes with deleted number index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted number index",
		  dependsOnMethods = { "testDeleteNumberData" },
		  enabled = false)
	public void testReadDeletedNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberIdx");
		
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("numberAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of number data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "NumberData")
	public Object[][] getNumberData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/number.data"));
		return data;
	}	
}
