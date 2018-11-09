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
 * CRUD tests for float data type index
 */
public class FloatIndexTests extends LifecycleServer {
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateFloatIndex - Insert nodes with float index
	 * @throws Exception
	 */
	@Test(description = "create Index Insert nodes with float index")
	public void testCreateFloatIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeFloatIdxType = gmd.getNodeType("nodeFloatIdx");
		if (nodeFloatIdxType == null)
			throw new Exception("Node type not found");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodeFloatIdxType);
			node.setAttribute("floatAttr", data[i][0]);
			node.setAttribute("key", i);
			conn.insertEntity(node);
		}
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testReadFloatIndex - Retrieve nodes by float index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes by float index",
		  dependsOnMethods = { "testCreateFloatIndex" })
	public void testReadFloatIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeFloatIdx");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("floatAttr", data[i][0]);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			//System.out.println("READ ATTR :" + entity.getAttribute("floatAttr").getValue());
			// Assert on Node attribute
			Assert.assertEquals(entity.getAttribute("key").getValue(), i);
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateFloatIndex - Update float index
	 * @throws Exception
	 */
	
	@Test(description = "Update float index",
		  dependsOnMethods = { "testReadFloatIndex" },
		  enabled = false)
	public void testUpdateFloatIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeFloatIdx");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("floatAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedFloatIndex - Retrieve nodes with updated float index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated float index",
		  dependsOnMethods = { "testUpdateFloatIndex" },
		  enabled = false)
	public void testReadUpdatedFloatIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeFloatIdx");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			// Assert.assertFalse(entity.getAttribute("floatAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("floatAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteFloatIndex - Delete float index
	 * @throws Exception
	 */
	
	@Test(description = "Delete float index",
		  dependsOnMethods = { "testReadUpdatedFloatIndex" },
		  enabled = false)
	public void testDeleteFloatIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeFloatIdx");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("floatAttr", null); // delete the boolean value
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedFloatIndex - Retrieve nodes with deleted float index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted float index",
		  dependsOnMethods = { "testDeleteFloatIndex" },
		  enabled = false)
	public void testReadDeletedFloatIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeFloatIdx");
		
		Object[][] data = this.getFloatData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
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
