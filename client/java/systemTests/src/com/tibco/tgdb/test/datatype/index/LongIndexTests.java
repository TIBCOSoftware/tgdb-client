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
 * CRUD tests for long data type index
 */
public class LongIndexTests extends LifecycleServer {

	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateLongIndex - Insert nodes with long index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with long index")
	public void testCreateLongIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeLongIdxType = gmd.getNodeType("nodeLongIdx");
		if (nodeLongIdxType == null)
			throw new Exception("Node type not found");
		
		Object[][] data = this.getLongData();
		for (int i=0; i<data.length; i++) {
			TGNode node = gof.createNode(nodeLongIdxType);
			node.setAttribute("longAttr", data[i][0]);
			node.setAttribute("key", i);
			conn.insertEntity(node);
		}
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testReadLongIndex - Retrieve nodes by long index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes by long index",
		  dependsOnMethods = { "testCreateLongIndex" })
	public void testReadLongIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeLongIdx");
		
		Object[][] data = this.getLongData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("longAttr", data[i][0]);
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
	 * testUpdateLongIndex - Update long index
	 * @throws Exception
	 */
	
	@Test(description = "Update long index",
		  dependsOnMethods = { "testReadLongIndex" },
		  enabled = false)
	public void testUpdateLongIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeLongIdx");
		
		Object[][] data = this.getLongData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("longAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedLongIndex - Retrieve nodes with updated long index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated long index",
		  dependsOnMethods = { "testUpdateLongIndex" },
		  enabled = false)
	public void testReadUpdatedLongIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeLongIdx");
		
		Object[][] data = this.getLongData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			// Assert.assertFalse(entity.getAttribute("longAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Assert.assertEquals(entity.getAttribute("longAttr").getValue(), data[i][1]);
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteLongIndex - Delete long index
	 * @throws Exception
	 */
	
	@Test(description = "Delete long index",
		  dependsOnMethods = { "testReadUpdatedLongIndex" },
		  enabled = false)
	public void testDeleteLongIndex() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeLongIdx");
		
		Object[][] data = this.getLongData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("longAttr", null); // delete the boolean value
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedLongIndex - Retrieve nodes with deleted long index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted long index",
		  dependsOnMethods = { "testDeleteLongIndex" },
		  enabled = false)
	public void testReadDeletedLongIndex() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeLongIdx");
		
		Object[][] data = this.getLongData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", 0);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("longAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of long data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "LongData")
	public Object[][] getLongData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/long.data"));
		return data;
	}	
}
