package com.tibco.tgdb.test.datatype.index;

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
 * CRUD tests for char data type index
 */
public class CharIndexTests extends LifecycleServer {

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateCharIndex - Insert nodes with char index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with char index")
	public void testCreateCharIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeCharIdxType = gmd.getNodeType("nodeCharIdx");
		if (nodeCharIdxType == null)
			throw new Exception("Node type not found");
		
		Object[][] data = this.getCharData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodeCharIdxType);
			node.setAttribute("charAttr", data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
		}
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testReadCharIndex - Retrieve nodes by char index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes by char index",
		  dependsOnMethods = { "testCreateCharIndex" })
	public void testReadCharIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeCharIdx");
		
		Object[][] data = this.getCharData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("charAttr", data[i][0]);
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
	 * testUpdateCharIndex - Update char index
	 * @throws Exception
	 */
	
	@Test(description = "Update char index",
		  dependsOnMethods = { "testReadCharIndex" },
		  enabled = false)
	public void testUpdateCharIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeCharIdx");
		
		Object[][] data = this.getCharData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("charAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedCharIndex - Retrieve nodes with updated char index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated char index",
		  dependsOnMethods = { "testUpdateCharIndex" },
		  enabled = false)
	public void testReadUpdatedCharIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeCharIdx");
		
		Object[][] data = this.getCharData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			// Assert.assertFalse(entity.getAttribute("charAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("charAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteCharIndex - Delete char index
	 * @throws Exception
	 */
	
	@Test(description = "Delete char index",
		  dependsOnMethods = { "testReadUpdatedCharIndex" },
		  enabled = false)
	public void testDeleteCharIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeCharIdx");
		
		Object[][] data = this.getCharData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("charAttr", null); // delete the boolean value
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedCharIndex - Retrieve nodes with updated char index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted boolean index",
		  dependsOnMethods = { "testDeleteCharIndex" },
		  enabled = false)
	public void testReadDeletedCharIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeCharIdx");
		
		Object[][] data = this.getCharData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("charAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of char data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "CharData")
	public Object[][] getCharData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/char.data"), Charset.forName("UTF-8"));
		return data;
	}	
}
