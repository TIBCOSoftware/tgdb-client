package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

public class Searchnodes {

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
			
			conn.getGraphMetadata(false);   //<===== no refresh = server crash
			
			TGKey key = gof.createCompositeKey("basicnode");
			
			key.setAttribute("name", "John");
      		TGEntity entity = conn.getEntity(key, null);
      		if (entity != null) {
      			System.out.println("From Node " + entity.getAttribute("name").getAsString() + " we get:");
      			for (TGEdge edge : ((TGNode)entity).getEdges()) {
      				System.out.println("  - Edge from " + edge.getVertices()[0].getAttribute("name").getAsString() + " to " + edge.getVertices()[1].getAttribute("name").getAsString());
      				System.out.println("    -- Attribute 'type' = " + (edge.getAttribute("type")==null?"null!":edge.getAttribute("type").getAsString()));	
      			}
      		}
      		else
      			System.out.println("Entity not found");
      		
		} finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}
