package com.tibco.tgdb.test.datatype.index;

import java.io.IOException;

import org.testng.Assert;
import org.testng.annotations.DataProvider;
import org.testng.annotations.Test;

import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.utils.ClasspathResource;
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
 * CRUD tests for boolean data type index
 */
public class BooleanIndexTests extends LifecycleServer {

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateBooleanIndex - Insert nodes with boolean index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with boolean index")
	public void testCreateBooleanIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodetype = gmd.getNodeType("boolNodetype");
		if (nodetype == null)
			throw new Exception("Node type not found");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodetype);
			node.setAttribute("key", i);
			node.setAttribute("boolAttr", data[i][0]);
			
			conn.insertEntity(node);
			//System.out.println("data " + data[i][0]);
		}
		
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testReadBooleanIndex - Retrieve nodes by boolean index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes by boolean index",
		  dependsOnMethods = { "testCreateBooleanIndex" })
	public void testReadBooleanIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("boolNodetype");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("boolAttr", data[i][0]);
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
	 * testUpdateBooleanData - Update boolean index
	 * @throws Exception
	 */
	
	@Test(description = "Update boolean index",
		  dependsOnMethods = { "testReadBooleanIndex" },
		  enabled = false)
	public void testUpdateBooleanIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("boolNodetype");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("boolAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedBooleanIndex - Retrieve nodes with updated boolean index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated boolean index",
		  dependsOnMethods = { "testUpdateBooleanIndex" },
		  enabled = false)
	public void testReadUpdatedBooleanIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("boolNodetype");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("boolAttr", data[i][1]);
			TGEntity entity = conn.getEntity(tgKey, null); 
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
  		
			// Assert on Node attribute
			//Assert.assertFalse(entity.getAttribute("boolAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("key").getValue(), i);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteBooleanIndex - Delete boolean index
	 * @throws Exception
	 */
	
	@Test(description = "Delete boolean index",
		  dependsOnMethods = { "testReadUpdatedBooleanIndex" },
		  enabled = false)
	public void testDeleteBooleanIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("boolNodetype");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("boolAttr", null); // nullify the boolean value. Null should fail for index
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedBooleanIndex - Retrieve nodes with updated boolean index
	 * @throws Exception
	 */
	@Test(description = "Retrieve nodes with deleted boolean attribute",
		  dependsOnMethods = { "testDeleteBooleanIndex" },
	      enabled = false) 
	public void testReadDeletedBooleanIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("boolNodetype");
		
		Object[][] data = this.getBooleanData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("boolAttr", null);
			TGEntity entity = conn.getEntity(tgKey, null); 
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
  		
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("boolAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of boolean data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "BoolData")
	public Object[][] getBooleanData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/bool.data"));
		return data;
	}	
}
