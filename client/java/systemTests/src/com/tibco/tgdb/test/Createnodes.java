package com.tibco.tgdb.test;

import java.util.ArrayList;
import java.util.Calendar;
import java.util.List;
import java.util.Random;
import java.util.UUID;

import org.testng.Assert;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

public class Createnodes {

	public static void main(String[] args) throws Exception {

		String url = "tcp://127.0.0.1:8222";
		String user = "scott";
		String pwd = "scott";
		String pkey = "Tom"+new Random().nextInt();
		String pkey2 = "Jim"+new Random().nextInt();
		TGConnection conn = null;
		try {
			conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
			conn.connect();

			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null) {
				throw new Exception("Graph object not found");
			}
			TGGraphMetadata gmd = conn.getGraphMetadata(true);
			TGNodeType nodeType = gmd.getNodeType("ratenode");
			if (nodeType == null)
				throw new Exception("Node type not found");
			
			
			TGNode node1 = gof.createNode(nodeType);
			node1.setAttribute("name", pkey);
			node1.setAttribute("extra", true); // true works fine
			conn.insertEntity(node1);
			
			conn.commit();
			System.out.println("Entity created\n");
			
			
			TGKey key = gof.createCompositeKey("ratenode");
			key.setAttribute("name", pkey);
			TGEntity entity = conn.getEntity(key, null);
			System.out.println("TYPE = " + entity.getEntityType().getName());
			if (entity != null) {
				System.out.println("boolean attribute on Node = " + entity.getAttribute("extra").getAsString());
			}	
			else 
				System.out.println("Could not retrieve entity");
					
		} finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}