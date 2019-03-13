package com.tibco.tgdb.test;

import java.util.Calendar;
import java.util.GregorianCalendar;
import java.util.List;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;

import com.tibco.tgdb.exception.TGAuthenticationException;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.test.lib.TGServer;

public class JustTry {

	public static void main(String[] args) throws Exception {
		String url = "tcp://127.0.0.1:8222";
		String user = "scott";
		String pwd = "scott";
		TGConnection conn = null;

		TGServer tgServer = new TGServer("C:/tgdb/2.0");
		List list = tgServer.getErrorsInLog();
		System.out.println(list.size());
		System.out.println(list.get(0).toString());
		
	}	
}

