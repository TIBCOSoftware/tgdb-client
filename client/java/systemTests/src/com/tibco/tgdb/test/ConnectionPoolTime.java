package com.tibco.tgdb.test;

import java.util.ArrayList;
import java.util.Calendar;
import java.util.Iterator;
import java.util.List;
import java.util.UUID;

import org.testng.Assert;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.connection.TGConnectionPool;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

public class ConnectionPoolTime {

	public static void main(String[] args) throws Exception {
		TGConnection conn = null;
		TGConnectionPool pool = null;
		int nbConnection = 1000;
		String url = "tcp://127.0.0.1:8222/{useDedicatedChannelPerConnection=false}";
		pool = TGConnectionFactory.getInstance().createConnectionPool(url, "scott", "scott", 10, null); 
		try {
			for (int i=0; i < nbConnection; i++) {
				
				long startTime = System.currentTimeMillis();
				pool.connect();
				long endTime = System.currentTimeMillis();
				conn = pool.get();
				TGGraphObjectFactory gof = conn.getGraphObjectFactory();
				if (gof == null)	
					Assert.fail("TG object factory is null for connection #" + nbConnection);
				// pool.release(conn); // skip releasing connection
				System.out.println("LOOP ITER #" + i + " - Connect Time = " + (endTime - startTime) + " msec");
				pool.disconnect();
			}
		}
		finally { 
			try { pool.disconnect();}
			catch (Exception e) {;}
		}
	}
}
