package com.tibco.tgdb.test.datatype.index;

import java.io.IOException;
import java.nio.charset.Charset;

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
 * CRUD tests for string data type index
 */
public class StringIndexTests extends LifecycleServer {

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateStringIndex - Insert nodes with string index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with string index")
	public void testCreateStringIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeStringIdxType = gmd.getNodeType("nodeStringIdx");
		if (nodeStringIdxType == null)
			throw new Exception("Node type not found");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodeStringIdxType);
			node.setAttribute("stringAttr", data[i][0]);
			node.setAttribute("key", i);
			conn.insertEntity(node);
			System.out.println("Key : " + i + " strIdx : " + data[i][0]);
		}
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testReadStringIndex - Retrieve nodes by string index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes by string index",
		  dependsOnMethods = { "testCreateStringIndex" })
	public void testReadStringIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeStringIdx");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("stringAttr", data[i][0]);
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
	 * testUpdateStringIndex - Update string index
	 * @throws Exception
	 */
	
	@Test(description = "Update string index",
		  dependsOnMethods = { "testCreateStringIndex" },
		  enabled = false)
	public void testUpdateStringIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeStringIdx");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("stringAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedStringIndex - Retrieve nodes with updated string index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated string index",
		  dependsOnMethods = { "testUpdateStringIndex" },
		  enabled = false)
	public void testReadUpdatedStringIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeStringIdx");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("stringAttr", data[i][1]);
			TGEntity entity = conn.getEntity(tgKey, null); 
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
  		
			// Assert on Node attribute
			// Assert.assertFalse(entity.getAttribute("stringAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("key").getValue(), i);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteStringIndex - Delete string index
	 * @throws Exception
	 */
	
	@Test(description = "Delete string index",
		  dependsOnMethods = { "testReadUpdatedStringIndex" },
		  enabled = false)
	public void testDeleteStringIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeStringIdx");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("stringAttr", null); // delete the boolean value
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedStringIndex - Retrieve nodes with deleted string index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted string index",
		  dependsOnMethods = { "testDeleteStringIndex" },
		  enabled = false)
	public void testReadDeletedStringIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeStringIdx");
		
		Object[][] data = this.getStringData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
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
