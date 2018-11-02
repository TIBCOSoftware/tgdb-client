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

public class TestEdge {

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
			TGNodeType basicnode = gmd.getNodeType("basicnode");
			if (basicnode == null)
				throw new Exception("Node type not found");
			
			TGNode basic1 = gof.createNode(basicnode);
			TGNode basic2 = gof.createNode(basicnode);
			TGEdge edge1, edge2, edge3;
			basic1.setAttribute("name", "Mike3");
			basic2.setAttribute("name", "John3");
			edge1 = gof.createEdge(basic1, basic2, TGEdge.DirectionType.Directed);
			edge1.setAttribute("type", "friend");
			edge2 = gof.createEdge(basic1, basic2, TGEdge.DirectionType.BiDirectional);
			edge2.setAttribute("type", "enemy");
			edge3 = gof.createEdge(basic1, basic2, TGEdge.DirectionType.UnDirected);
			edge3.setAttribute("type", "family");
			conn.insertEntity(basic1);
			conn.insertEntity(basic2);
			conn.insertEntity(edge1);
			conn.insertEntity(edge2);
			conn.insertEntity(edge3);
			conn.commit();
			System.out.println("Entities created successfully");
			
			conn.getGraphMetadata(true);
			TGKey key = gof.createCompositeKey("basicnode");
			
			key.setAttribute("name", "John3");
      		TGEntity entity = conn.getEntity(key, null);
      		if (entity != null) {
      			System.out.println("From Node " + entity.getAttribute("name").getAsString() + " we get:");
      			for (TGEdge edge : ((TGNode)entity).getEdges(TGEdge.DirectionType.Directed)) {
      				System.out.println("  - Edge from " + edge.getVertices()[0].getAttribute("name").getAsString() + " to " + edge.getVertices()[1].getAttribute("name").getAsString());
      				System.out.println("  - Edge direction : " + edge.getDirectionType());
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
