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

public class DeleteRecreateNode {

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
			
			TGKey key = gof.createCompositeKey("basicnode");
			
			key.setAttribute("name", "Serge");
      		TGEntity entity = conn.getEntity(key, null);
      		if (entity != null) {
      			conn.deleteEntity(entity);
      			// conn.commit();  // If we do commit(), re-insert works
      			//TGGraphMetadata gmd = conn.getGraphMetadata(true);
      			TGNodeType basicnode = gmd.getNodeType("basicnode");
      			entity = gof.createNode(basicnode);
    			entity.setAttribute("name", "Serge");
    			conn.insertEntity(entity);
    			conn.commit();
    			System.out.println("Entity deleted and re-created");
      		}
      		else
      			System.out.println("Entity not found");
      		
		} finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}
