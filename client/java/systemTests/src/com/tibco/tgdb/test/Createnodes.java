package com.tibco.tgdb.test;

import java.util.ArrayList;
import java.util.Calendar;
import java.util.List;
import java.util.UUID;

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
			
			Double[] arrayd = {0.0D,1.0D,Double.POSITIVE_INFINITY,Double.NEGATIVE_INFINITY,Double.NaN,Double.MIN_VALUE,Double.MAX_VALUE,Double.MIN_NORMAL};
			
			for (Double d : arrayd) {
				TGNode node = gof.createNode(nodeType);
				node.setAttribute("name", UUID.randomUUID().toString());
				node.setAttribute("rate", d);
				conn.insertEntity(node);
			
			}
			conn.commit();
			System.out.println("Entities created\n");
			
			for (Double d : arrayd) {
				TGKey key = gof.createCompositeKey("ratenode");
				key.setAttribute("rate", d);
				TGEntity entity = conn.getEntity(key, null);
				if (entity != null) {
					System.out.println("rate = " + entity.getAttribute("rate").getValue());
				}	
				else 
					System.out.println("Could not retrieve entity with rate="+d);
			}
					
		} finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}