package com.tibco.tgdb.test;

import java.util.Calendar;
import java.util.GregorianCalendar;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;

import com.tibco.tgdb.exception.TGAuthenticationException;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

public class JustTry {

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
			
			conn.getGraphMetadata(true);
			while (true) {;}
		}
		
		finally {
			if (conn != null)
				conn.disconnect();
		}
	}	
}

