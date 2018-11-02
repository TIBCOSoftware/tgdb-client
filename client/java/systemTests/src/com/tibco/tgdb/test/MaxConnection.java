package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGGraphObjectFactory;


public class MaxConnection {

	public static void main(String[] args) throws Exception {
		
		String url = "tcp://localhost:8222";

		TGConnection conn = null;
		try {
			conn = TGConnectionFactory.getInstance().createConnection(url, "scott", "scott", null);
			conn.connect();
			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null)	
				System.out.println("GOF is null");
		}
		finally { 
			System.out.println("THE END");
			//conn.disconnect();
		}
	}
}
