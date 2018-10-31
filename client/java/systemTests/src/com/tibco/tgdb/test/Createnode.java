package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.query.TGResultSet;

public class Createnode {

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
			
			System.out.printf("Querying...");
            String queryString = "@nodetype = 'houseMemberType' and yearBorn > 1700 and yearBorn < 1800;";
            TGResultSet resultSet = conn.executeQuery(queryString, null);
            System.out.println("ResultSet = " + resultSet);
      		
		} finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}
