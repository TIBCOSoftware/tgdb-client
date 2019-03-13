package com.tibco.tgdb.test;

import java.util.ArrayList;
import java.util.Calendar;
import java.util.Iterator;
import java.util.List;
import java.util.UUID;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

public class ParallelConnections {

	public static void main(String[] args) throws Exception {
		String url = "tcp://127.0.0.1:8222"; // /{connectTimeout=15000}";
		String user = "scott";
		String pwd = "scott";

		List<Thread> threadList = new ArrayList<Thread>();

		for (int i = 0; i < 50; i++) {
			final int j = i;
			Thread t = new Thread() {
				public void run() {
					TGConnection conn = null;
					try {
						System.out.println("Thread-" + j + " connecting...");
						conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
						conn.connect();
						System.out.println("Thread-" + j + " connected");
						TGGraphObjectFactory gof = conn.getGraphObjectFactory();
						if (gof == null) {
							throw new org.testng.TestException("thread-" + j + " TG object factory is null");
						}

						 // sleep a little to make sure we use all connections on server at the same time
						
						
						UUID uuid = UUID.randomUUID();
						String randomUUIDString = uuid.toString();
						TGGraphMetadata gmd = conn.getGraphMetadata(true);
						TGNodeType basicnode = gmd.getNodeType("basicnode");
						if (basicnode == null)
							throw new Exception("Thread-" + j + " Node type not found");
						TGNode basic1 = gof.createNode(basicnode);
						basic1.setAttribute("name", randomUUIDString);
						basic1.setAttribute("age", 73);
						//basic1.setAttribute("createtm", new Calendar.Builder().setDate(2016, 10, 30).build());
						conn.insertEntity(basic1);
						conn.commit();
						
						Thread.sleep(3000);
						
						System.out.println("Thread-" + j + " disconnected");
					} catch (Exception v) {
						System.out.println("Thread-" + j + " --> " + v);
						v.printStackTrace();
					} 
				finally {
						if (conn != null)
							conn.disconnect();
					}
				}
			};
			threadList.add(t);
		}

		Iterator<Thread> iter = threadList.iterator();
		while (iter.hasNext()) {
			Thread t = iter.next();
			t.start();
		}
		Iterator<Thread> iter2 = threadList.iterator();
		while (iter2.hasNext()) {
			Thread t = iter2.next();
			t.join();
		}
	}
}
