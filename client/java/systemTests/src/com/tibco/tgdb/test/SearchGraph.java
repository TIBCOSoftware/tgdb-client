package com.tibco.tgdb.test;

import java.text.SimpleDateFormat;
import java.util.Calendar;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.utils.SortedProperties;
import com.tibco.tgdb.utils.TGProperties;

public class SearchGraph {

	public static void main(String[] args) throws Exception {
		String url = "tcp://127.0.0.1:8222";
		String user = "napoleon";
		String pwd = "bonaparte";
		TGConnection conn = null;
		try {
			conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
			conn.connect();

			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null) {
				throw new Exception("Graph object not found");
			}

			TGGraphMetadata gmd = conn.getGraphMetadata(true);
			TGNodeType houseMemberType = gmd.getNodeType("houseMemberType");
			TGNode houseMemberNode = gof.createNode(houseMemberType);
			
			houseMemberNode.setAttribute("memberName", "Napoleon Bonaparte");
			houseMemberNode.setAttribute("crownName", "Napoleon I");
			houseMemberNode.setAttribute("houseHead", false);
			houseMemberNode.setAttribute("yearBorn", 1769);
			houseMemberNode.setAttribute("yearDied", 1821);
			houseMemberNode.setAttribute("crownTitle", "Emperor of the French");
			Calendar date = Calendar.getInstance();
			date.setTime((new SimpleDateFormat("dd MMM yyyy").parse("18 May 1804")));
			houseMemberNode.setAttribute("reignStart", date);
			date.setTime((new SimpleDateFormat("dd MMM yyyy").parse("22 Jun 1815")));
			houseMemberNode.setAttribute("reignEnd", date);
			conn.insertEntity(houseMemberNode);
			conn.commit(); // Write data to database
			
			TGKey houseKey = gof.createCompositeKey("houseMemberType");
			
			houseKey.setAttribute("memberName", "Napoleon Bonaparte");
      		TGProperties<String, String> houseProps = new SortedProperties<String, String>();
      		houseProps.put("fetchsize", "0");
      		houseProps.put("traversaldepth", "0");
      		houseProps.put("edgelimit", "0");
      		TGQueryOption options = TGQueryOption.createQueryOption();
      		options.setEdgeLimit(0);
      		options.setPrefetchSize(0);
      		options.setTraversalDepth(0);
      		TGEntity houseEntity = conn.getEntity(houseKey, options);
      		if (houseEntity != null) {
      			System.out.println("Entity found : " + houseEntity.getAttribute("crownTitle").getAsString());
      		} else {
      			System.out.println("Entity not found");
      		}
      		
		}
		finally {
			if (conn != null)
				conn.disconnect();
		}
	}

}
