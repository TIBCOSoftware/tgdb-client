package com.tibco.tgdb.test.datatype.attribute;

import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
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
 * Create Read Update Delete (CRUD) tests for number data type attribute
 */
public class NumberAttrTests extends LifecycleServer {	
	
	Object[][] data;
	
	public NumberAttrTests() throws IOException, EvalError {
		
		this.data = this.getNumberData();
	}
	
	/************************
	 * 
	 * Test Cases
	 * 
	 ************************/
	
  /**
	 * testCreateNumberData - Insert nodes and edge with number attribute
	 * @throws Exception
	 */
	@Test(description = "Insert nodes and edge with number attribute")
	public void testCreateNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		TGGraphMetadata gmd = conn.getGraphMetadata(true);
		TGNodeType nodeNumberAttrsType = gmd.getNodeType("nodeNumberAttrs");
		if (nodeNumberAttrsType == null)
			throw new Exception("Node type not found");
		
		List<TGNode> nodes = new ArrayList<TGNode>();
		for (int i=0; i<data.length; i++) {
			//if (data[i][0] != null)
			//	System.out.println("CREATE BIGDECIMAL = " + data[i][0] + " - UnscaledValue = " + ((BigDecimal)data[i][0]).unscaledValue() + " - Precision = " + ((BigDecimal)data[i][0]).precision() + " - Scale = " + ((BigDecimal)data[i][0]).scale());
			//else
			//	System.out.println("CREATE BIGDECIMAL = " + data[i][0]);
			TGNode node = gof.createNode(nodeNumberAttrsType);
			node.setAttribute("numberAttr"+i, data[i][0]);
			node.setAttribute("key", i);
			nodes.add(node);
			conn.insertEntity(node);
			if (i>0) {
				TGEdge edge = gof.createEdge(nodes.get(i-1), nodes.get(i), TGEdge.DirectionType.UnDirected);
				edge.setAttribute("numberAttr"+(i-1), data[i-1][0]);
				conn.insertEntity(edge);
			}
		}
		TGEdge edge = gof.createEdge(nodes.get(data.length-1), nodes.get(0), TGEdge.DirectionType.UnDirected);
		edge.setAttribute("numberAttr"+(data.length-1), data[data.length-1][0]);
		conn.insertEntity(edge);
		conn.commit();
	
		conn.disconnect();
		// Expect 0 Errors in log file
		int nbErrorsInLog = tgServer.getErrorsInLog().size();
		Assert.assertEquals(nbErrorsInLog, 0, "Found " + nbErrorsInLog + " Error(s) in server log file -");
	}
	
	/**
	 * testReadNumberData - Retrieve nodes and edge with number attribute
	 * @throws Exception
	 */
	
	@Test(description = "Retrieve nodes and edge with number attribute",
		  dependsOnMethods = { "testCreateNumberData" },enabled = true)
	public void testReadNumberData() throws Exception {
		
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		try {
			conn.connect();
		
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null) {
				throw new org.testng.TestException("TG object factory is null");
			}
		
			conn.getGraphMetadata(true);
			TGKey tgKey = gof.createCompositeKey("nodeNumberAttrs");
		
			for (int i=0; i<data.length; i++) {
				tgKey.setAttribute("key", i);
				TGEntity entity = conn.getEntity(tgKey, null);
				if (entity == null) {
					throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
				}
				// Assert on Node attribute
				if (data[i][0] != null) { 
					//System.out.println("READ BIGDECIMAL = " + entity.getAttribute("numberAttr"+i).getValue() + " - UnscaledValue = " + ((BigDecimal)entity.getAttribute("numberAttr"+i).getValue()).unscaledValue() + " - Precision = " + ((BigDecimal)entity.getAttribute("numberAttr"+i).getValue()).precision() + " - Scale = " + ((BigDecimal)entity.getAttribute("numberAttr"+i).getValue()).scale());	
					Assert.assertEquals(((BigDecimal) entity.getAttribute("numberAttr"+i).getValue()).compareTo((BigDecimal)data[i][0]), 0, "Actual and Expected BigDecimal are not the same");
				}
				else {
					//System.out.println("READ BIGDECIMAL = " + entity.getAttribute("numberAttr"+i).getValue());
					Assert.assertEquals(entity.getAttribute("numberAttr"+i).getValue(), data[i][0]);
				}
				for (TGEdge edge : ((TGNode)entity).getEdges()) {
					if (edge.getVertices()[0].equals(entity))  {
						// Assert on Edge attribute
						if (data[i][0] != null) 
							Assert.assertEquals(((BigDecimal) edge.getAttribute("numberAttr"+i).getValue()).compareTo((BigDecimal) data[i][0]), 0, "Actual and Expected BigDecimal are not the same");
						else
							Assert.assertEquals(edge.getAttribute("numberAttr"+i).getValue(), data[i][0]);
					}
				}
			}
			// Expect 0 Errors in log file
			int nbErrorsInLog = tgServer.getErrorsInLog().size();
			Assert.assertEquals(nbErrorsInLog, 0, "Found " + nbErrorsInLog + " Error(s) in server log file -");
		}
		finally {
			conn.disconnect();
		}
	}
	
	/**
	 * testUpdateNumberData - Update number attribute
	 * @throws Exception
	 */
	
	@Test(description = "Update number attribute",
		  dependsOnMethods = { "testReadNumberData" },
		  enabled = true)
	public void testUpdateNumberData() throws Exception {
TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberAttrs");
		
		for (int i=0; i<data.length; i++) {
			//if (data[i][1] != null)
			//	System.out.println("UPDATE BIGDECIMAL = " + data[i][1] + " - UnscaledValue = " + ((BigDecimal)data[i][1]).unscaledValue() + " - Precision = " + ((BigDecimal)data[i][1]).precision() + " - Scale = " + ((BigDecimal)data[i][1]).scale());
			//else 
			//	System.out.println("UPDATE BIGDECIMAL = " + data[i][1]);
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("numberAttr"+i, data[i][1]); // 
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
		// Expect 0 Errors in log file
		int nbErrorsInLog = tgServer.getErrorsInLog().size();
		Assert.assertEquals(nbErrorsInLog, 0, "Found " + nbErrorsInLog + " Error(s) in server log file -");
	}
	
	/**
	 * testReadUpdatedNumberData - Retrieve nodes with updated number attribute
	 * @throws Exception
	 */
	@Test(description = "Retrieve nodes with updated number attribute",
		  dependsOnMethods = { "testUpdateNumberData" }, 
		  enabled = true)
	public void testReadUpdatedNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberAttrs");
		
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
  		
			// Assert on Node attribute
			if (data[i][1] != null) {
				//System.out.println("READ UPDATED BIGDECIMAL = " + entity.getAttribute("numberAttr"+i).getValue() + " - UnscaledValue = " + ((BigDecimal)entity.getAttribute("numberAttr"+i).getValue()).unscaledValue() + " - Precision = " + ((BigDecimal)entity.getAttribute("numberAttr"+i).getValue()).precision() + " - Scale = " + ((BigDecimal)entity.getAttribute("numberAttr"+i).getValue()).scale());
				Assert.assertEquals(((BigDecimal) entity.getAttribute("numberAttr"+i).getValue()).compareTo((BigDecimal) data[i][1]), 0, "Actual and Expected BigDecimal are not the same");
			} else {
				//System.out.println("READ UPDATED BIGDECIMAL = " + entity.getAttribute("numberAttr"+i).getValue());
				Assert.assertEquals(entity.getAttribute("numberAttr"+i).getValue(), data[i][1]);
			}
		}
		conn.disconnect();
	}
	
	/**
	 * testDeleteNumberData - Delete number attribute
	 * @throws Exception
	 */
	@Test(description = "Delete number attribute",
		  dependsOnMethods = { "testReadUpdatedNumberData" }, 
		  enabled = true)
	public void testDeleteNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberAttrs");
		
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null);
			if (entity == null) {
				throw new org.testng.TestException("TG entity #"+i+" was not retrieved");
			}
			entity.setAttribute("numberAttr"+i, null); // delete the number value by setting it up to null
			conn.updateEntity(entity);
			conn.commit();
		}
		conn.disconnect();
		// Expect 0 Errors in log file
		int nbErrorsInLog = tgServer.getErrorsInLog().size();
		Assert.assertEquals(nbErrorsInLog, 0, "Found " + nbErrorsInLog + " Error(s) in server log file -");
	}
	
	/**
	 * testReadDeletedNumberData - Retrieve nodes with updated number attribute
	 * @throws Exception
	 */
	@Test(description = "Retrieve nodes with deleted number attribute",
		  dependsOnMethods = { "testDeleteNumberData" }, 
		  enabled = true)
	public void testReadDeletedNumberData() throws Exception {
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(tgUrl, tgUser, tgPwd, null);
		
		conn.connect();
		
		TGGraphObjectFactory gof = conn.getGraphObjectFactory();
		if (gof == null) {
			throw new org.testng.TestException("TG object factory is null");
		}
		
		conn.getGraphMetadata(true);
		TGKey tgKey = gof.createCompositeKey("nodeNumberAttrs");
		
		for (int i=0; i<data.length; i++) {
			tgKey.setAttribute("key", i);
			TGEntity entity = conn.getEntity(tgKey, null); 
 
			// Assert on Node attribute
			Assert.assertTrue(entity.getAttribute("numberAttr"+i).isNull(), "Expected attribute #"+i+" null but found it non null -");
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
	
	public static BigDecimal getRandomBigDecimal(int precision, int scale) {
		String chars = "0123456789";
		int length = chars.length();
		String number = "";
		for (int i=0; i<precision; i++){
			long j = Math.round((Math.random() * (length-1)));
			number += chars.split("")[(int)j];
		}
		return new BigDecimal(new BigInteger(number), scale);
	}
	
}
