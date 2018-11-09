package com.tibco.tgdb.test.datatype.index;

import java.io.IOException;

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
 * CRUD tests for integer data type index
 */
public class IntIndexTests extends LifecycleServer {

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateIntIndex - Insert nodes with int index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with int index")
	public void testCreateIntIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeIntIdxType = gmd.getNodeType("nodeIntIdx");
		if (nodeIntIdxType == null)
			throw new Exception("Node type not found");
		
		Object[][] data = this.getIntData();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodeIntIdxType);
			node.setAttribute("intAttr", data[i][0]);
			node.setAttribute("key", i);
			conn.insertEntity(node);
		}
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testReadIntIndex - Retrieve nodes with int index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with int index",
		  dependsOnMethods = { "testCreateIntIndex" })
	public void testReadIntIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeIntIdx");
		
		Object[][] data = this.getIntData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("intAttr", data[i][0]);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			// System.out.println("READ ATTR :" + data[i][0]);
			// Assert on Node attribute
			Assert.assertEquals(entity.getAttribute("key").getValue(), i);
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateIntIndex - Update int index
	 * @throws Exception
	 */
	
	@Test(description = "Update int index",
		  dependsOnMethods = { "testReadIntIndex" },
		  enabled = false)
	public void testUpdateIntIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeIntIdx");
		
		Object[][] data = this.getIntData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("intAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedIntIndex - Retrieve nodes with updated int index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated int index",
		  dependsOnMethods = { "testUpdateIntIndex" },
		  enabled = false)
	public void testReadUpdatedIntIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeIntIdx");
		
		Object[][] data = this.getIntData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			// Assert.assertFalse(entity.getAttribute("intAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("intAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteIntIndex - Delete int index
	 * @throws Exception
	 */
	
	@Test(description = "Delete int index",
		  dependsOnMethods = { "testReadUpdatedIntIndex" },
		  enabled = false)
	public void testDeleteIntIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeIntIdx");
		
		Object[][] data = this.getIntData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("intAttr", null); // delete the boolean value
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedIntIndex - Retrieve nodes with deleted int index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted int index",
		  dependsOnMethods = { "testDeleteIntIndex" },
		  enabled = false)
	public void testReadDeletedIntIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeIntIdx");
		
		Object[][] data = this.getIntData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("intAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of int data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "IntData")
	public Object[][] getIntData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/int.data"));
		return data;
	}	
}
