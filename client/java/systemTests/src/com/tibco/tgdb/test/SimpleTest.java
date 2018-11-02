package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.test.lib.TGAdmin;
import com.tibco.tgdb.test.lib.TGAdminException;
import com.tibco.tgdb.test.lib.TGInitException;
import com.tibco.tgdb.test.lib.TGServer;
import com.tibco.tgdb.test.lib.TGStartException;

public class SimpleTest {

	private String tgHome = null;
	private TGServer tgServer = null;
	
	public SimpleTest(String tgHome) throws Exception {
		TGServer.killAll();
		this.tgHome = tgHome;
		tgServer = new TGServer(tgHome, tgHome+"/bin/tgdb.conf");
	}
	
	public void initServer() throws TGInitException {
		tgServer.init(tgHome+"/bin/initdb.conf", true, 60000);
	}
	
	public void startServer() throws TGStartException   {
		tgServer.start(10000);
		System.out.println(tgServer.getBanner());
	}
	
	public void stopServer() throws TGAdminException {
		// tgServer.kill();
		TGAdmin.stopServer(tgServer, null, null, null, 10000);
	}
	
	public void test1() throws TGException {
		
		System.out.println("Test1 - Start");
		// Connect via IPv6
		String url = tgServer.getNetListeners()[0].getUrl();
		String user = tgServer.getSystemUser();
		String pwd = tgServer.getSystemPwd();
		TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
        TGGraphMetadata gmd = conn.getGraphMetadata(true);

        TGNodeType basicnode = gmd.getNodeType("basicnode");
        if (basicnode == null) throw new TGException("Node type not found");
        
        TGNode node1 = gof.createNode(basicnode);

        node1.setAttribute("name", "Mike");
        conn.insertEntity(node1);
        conn.commit();
        System.out.println("Test1 - End Commit successful");
        conn.disconnect();

	}
	
public void test2() throws TGException {
		
	System.out.println("Test2 - Start");
	// Connect via IPv4
	String url = tgServer.getNetListeners()[1].getUrl();
	String user = tgServer.getSystemUser();
	String pwd = tgServer.getSystemPwd();
	TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
    conn.connect();

    TGGraphObjectFactory gof = conn.getGraphObjectFactory();
    TGGraphMetadata gmd = conn.getGraphMetadata(true);

    TGKey key = gof.createCompositeKey("basicnode");
    key.setAttribute("name", "Mike");
    TGEntity entity = conn.getEntity(key, null);
    System.out.println("Test2 - End - Entity : "+entity);
    conn.disconnect();

	}
	
	public static void main(String[] args) throws Exception {
		
		SimpleTest simple = new SimpleTest("C:/tgdb/1.0");
		simple.initServer();
		simple.startServer();
		simple.test1();
		simple.stopServer();
		simple.startServer();
		simple.test2();
		simple.stopServer();
		
	}	
}

