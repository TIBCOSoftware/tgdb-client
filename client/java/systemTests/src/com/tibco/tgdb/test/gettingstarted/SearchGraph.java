package com.tibco.tgdb.test.gettingstarted;

/**
 * Copyright 2018 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except 
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import java.text.SimpleDateFormat;
import java.util.Calendar;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;

/**
 * Search for a member in the House of Bonaparte graph 
 * and display the member attributes and children
 * 
 * Usage : java SearchGraph -memberName "Carlo Bonaparte" 
 * 
 */
public class SearchGraph {
	
	static String url = "tcp://127.0.0.1:8222";
	static String user = "napoleon";
	static String pwd = "bonaparte";
	
	static String memberName = "Napoleon Bonaparte";
	
	public static void main(String[] args) throws Exception {
		
		for (int i=0; i<args.length; i++) {
			if (args[i].equals("-memberName"))
				memberName = args[i+1];
		}
		
		TGConnection conn = null;
		try {
			conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
			conn.connect();

			TGGraphObjectFactory gof = conn.getGraphObjectFactory();
			if (gof == null) {
				throw new Exception("Graph object not found");
			}

			conn.getGraphMetadata(true);
			
			TGKey houseKey = gof.createCompositeKey("houseMemberType");
			
			houseKey.setAttribute("memberName", memberName);
			System.out.printf("Searching for member '%s'...\n",memberName);
      		TGEntity houseMember = conn.getEntity(houseKey, null);
      		if (houseMember != null) {
      			SimpleDateFormat simpleFormat = new SimpleDateFormat("dd MMM yyyy");
      			System.out.printf("House member '%s' found\n",houseMember.getAttribute("memberName").getAsString());
      			for (TGAttribute attr : houseMember.getAttributes()) {
      				if (attr.getValue() == null)
      					System.out.printf("\t%s: %s\n", attr.getAttributeDescriptor().getName(), "");
      				else
      					System.out.printf("\t%s: %s\n", attr.getAttributeDescriptor().getName(), (attr.getValue() instanceof Calendar)?(simpleFormat.format(((Calendar)attr.getValue()).getTime())):attr.getValue());
      			}
      			for (TGEdge relation : ((TGNode)houseMember).getEdges(TGEdge.DirectionType.Directed)) { // Directed == child
      				TGNode[] vertices = relation.getVertices();
      				TGNode fromMember = vertices[0];
      				TGNode toMember = vertices[1];
      				if (fromMember == houseMember) {
      					System.out.printf("\tchild: %s\n", toMember.getAttribute("memberName").getAsString());
      				}
      			}
      		} else {
      			System.out.printf("House member '%s' not found", memberName);
      		}
		}
		finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}
