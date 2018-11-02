package com.tibco.tgdb.test;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.GregorianCalendar;
import java.util.List;
import java.util.TimeZone;

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

public class RetrieveNodes {

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
			TGNodeType nodeType = gmd.getNodeType("nodeType1");
			if (nodeType == null)
				throw new Exception("Node type not found");
			
			int nbNodes = 9;

			conn.getGraphMetadata(true);
			TGKey key = gof.createCompositeKey("nodeType1");
			
			for (int i=0; i<nbNodes; i++) {
				key.setAttribute("key", i);
				TGEntity entity = conn.getEntity(key, null);
				if (entity != null) {
					System.out.print("key = " + entity.getAttribute("key").getValue() + " - ");
					System.out.println("stringAttr = " + entity.getAttribute("stringAttr").getValue());
				}	
				else 
					System.out.println("Could not retrieve entity #" + i);
			}
		} finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}