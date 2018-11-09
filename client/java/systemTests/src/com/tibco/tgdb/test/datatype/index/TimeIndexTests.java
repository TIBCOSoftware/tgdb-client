package com.tibco.tgdb.test.datatype.index;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Calendar;
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
 * CRUD tests for time data type index
 */
@Ignore
public class TimeIndexTests extends LifecycleServer {
	
	Object[][] data;
	
	public TimeIndexTests() throws IOException, EvalError {
		this.data = this.getTimeData();
	}
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateTimeData - Insert nodes with time index
	 * @throws Exception
	 */
	@Test(description = "Insert nodes with time index")
	public void testCreateTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeTimeIdxType = gmd.getNodeType("nodeTimeIdx");
		if (nodeTimeIdxType == null)
			throw new Exception("Node type not found");
		
		//Object[][] data = this.getTimeData();
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			//System.out.println("CREATE ATTR:" + data[i][0]);
			TGNode node = gof.createNode(nodeTimeIdxType);
			node.setAttribute("timeAttr", data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
		}
		conn.commit();
		conn.disconnect();
	}
	
	/**
	 * testReadTimeData - Retrieve nodes with time index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with time index",
		  dependsOnMethods = { "testCreateTimeData" })
	public void testReadTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			System.out.println("READ ATTR:" + entity.getAttribute("timeAttr").getValue());
			// Assert on Node attribute
			// Assert only on time (HOUR, MIN, SEC, MILLISEC) since the Date part that comes back from DB is junk
			Calendar timeAttr = (Calendar)entity.getAttribute("timeAttr").getValue();
			if (timeAttr != null) {
				Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) data[i][0]).get(Calendar.HOUR_OF_DAY));
				Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) data[i][0]).get(Calendar.MINUTE));
				Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) data[i][0]).get(Calendar.SECOND));
				Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) data[i][0]).get(Calendar.MILLISECOND));
			}
			else { // Attribute value is Null. Make sure the original value was Null too
				Assert.assertEquals(timeAttr, data[i][0]);
			}
		}
		conn.disconnect();
	}
	
	/**
	 * testUpdateTimeData - Update time index
	 * @throws Exception
	 */
	
	@Test(description = "Update time index",
		  dependsOnMethods = { "testReadTimeData" },
		  enabled = false)
	public void testUpdateTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #" + i + " was not retrieved");
			}
			//System.out.println("UPDATE ATTR:" + data[i][1] + " - Length:" + ((Time) data[i][1]).length());
			entity.setAttribute("timeAttr", data[i][1]); 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadUpdatedTimeData - Retrieve nodes with updated time index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with updated time index",
		  dependsOnMethods = { "testUpdateTimeData" },
		  enabled = false)
	public void testReadUpdatedTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			//System.out.println("READ UPDATED ATTR:" + entity.getAttribute("timeAttr").getValue());
			// Assert on Node attribute
			// Assert only on time (HOUR, MIN, SEC, MILLISEC) since the Date part that comes back from DB is junk
			// Assert.assertFalse(entity.getAttribute("timeAttr").isNull(), "Expected attribute #"+i+" non null but found it null -");
			Calendar timeAttr = (Calendar)entity.getAttribute("timeAttr").getValue();
			if (timeAttr != null) {
				Assert.assertEquals(timeAttr.get(Calendar.HOUR_OF_DAY), ((Calendar) data[i][1]).get(Calendar.HOUR_OF_DAY));
				Assert.assertEquals(timeAttr.get(Calendar.MINUTE), ((Calendar) data[i][1]).get(Calendar.MINUTE));
				Assert.assertEquals(timeAttr.get(Calendar.SECOND), ((Calendar) data[i][1]).get(Calendar.SECOND));
				Assert.assertEquals(timeAttr.get(Calendar.MILLISECOND), ((Calendar) data[i][1]).get(Calendar.MILLISECOND));
			}
			else {// Attribute value is Null. Make sure the original value was Null too
				Assert.assertEquals(timeAttr, data[i][1]);
			}
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteTimeData - Delete time index
	 * @throws Exception
	 */
	
	@Test(description = "Delete time index",
		  dependsOnMethods = { "testReadUpdatedTimeData" },
		  enabled = false)
	public void testDeleteTimeData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("timeAttr", null); // delete the time value by setting it up to null
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
	}
	
	/**
	 * testReadDeletedTimeData - Retrieve nodes with deleted time index
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes with deleted time index",
		  dependsOnMethods = { "testDeleteTimeData" },
		  enabled = false)
	public void testReadDeletedTimeData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeTimeIdx");
		
		//Object[][] data = this.getTimeData();
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
 
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("timeAttr").isNull(), "Expected attribute #"+i+" null but found it non null -");
		}
		conn.disconnect();
	}
	
	/************************
	 * 
	 * Data Providers 
	 * 
	 ************************/
	
	
	/**
	 * Provide a set of time data
	 * @return Object[][] of data
	 * @throws IOException
	 * @throws EvalError
	 */
	@DataProvider(name = "TimeData")
	public Object[][] getTimeData() throws IOException, EvalError {
		Object[][] data =  PipedData.read(this.getClass().getResourceAsStream("/"+this.getClass().getPackage().getName().replace('.', '/') + "/time.data"));
		return data;
	}
}
