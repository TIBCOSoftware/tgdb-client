package com.tibco.tgdb.test.gettingstarted;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;

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

/**
 * Delete a given member of the House
 * 
 * Usage : java DeleteGraph -memberName <memberName>
 *    
 *  For instance to delete the house member named "Napoleon Bonaparte" :
 *  java DeleteGraph -memberName "Napoleon Bonaparte"
 *
 */
public class DeleteGraph {
	
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
      			conn.deleteEntity(houseMember);
      			conn.commit();
      			System.out.printf("House member '%s' deleted successfully", memberName);
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
