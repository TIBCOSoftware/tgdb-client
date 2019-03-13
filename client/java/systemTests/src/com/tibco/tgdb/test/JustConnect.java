package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;

import com.tibco.tgdb.exception.TGAuthenticationException;
import com.tibco.tgdb.exception.TGChannelDisconnectedException;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGGraphObjectFactory;

public class JustConnect {

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
			Thread.sleep(15000);
		}
		
		/*
		catch (TGAuthenticationException e) {
			System.out.println("Caught TGAuth : " + e.getClass() + " - "  + e.getMessage());
		}
		catch (Exception e) {
			System.out.println("Caught something else : " + e.getClass() + " - " + e.getMessage());
		}
		
		catch (TGException e) {
			System.out.println("Caught something : " + e.getMessage());
		}
		*/
		finally {
			if (conn != null)
				conn.disconnect();
		}
	}	
}

