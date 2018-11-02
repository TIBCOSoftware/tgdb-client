package com.tibco.tgdb.test.gettingstarted;



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
import com.tibco.tgdb.query.TGResultSet;

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
 * Query for members in the House of Bonaparte graph 
 * born between the start and end years
 * and display the member attributes.
 * 
 * Usage : java QueryGraph -startyear 1900 -endyear 2000
 * 
 */
public class QueryGraph {
	
	static String url = "tcp://127.0.0.1:8222";
	static String user = "napoleon";
	static String pwd = "bonaparte";
	
	static int startYear = 1700;
    static int endYear = 1800;
	
	public static void main(String[] args) throws Exception {
		
		
        try {
            for (int i=0; i<args.length; i++) {
                if (args[i].equals("-startyear")) {
                    startYear = Integer.parseInt(args[i+1]);
                } else if (args[i].equals("-endyear")) {
                    endYear = Integer.parseInt(args[i+1]);
                }
            }
        } catch (NumberFormatException ex) {
            System.out.printf("Invalid year value specified\n");
            return;
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
			
			
			System.out.printf("Querying for member born between %d and %d\n", startYear, endYear);
            String queryString = "@nodetype = 'houseMemberType' and yearBorn > " + startYear + " and yearBorn < " + endYear + ";";
            TGResultSet resultSet = conn.executeQuery(queryString, null);

      		if (resultSet != null) {
                while (resultSet.hasNext()) {
                    TGEntity houseMember = resultSet.next();
                    SimpleDateFormat simpleFormat = new SimpleDateFormat("dd MMM yyyy");
                    System.out.printf("House member '%s' found\n",houseMember.getAttribute("memberName").getAsString());
                    for (TGAttribute attr : houseMember.getAttributes()) {
                        if (attr.getValue() == null)
                            System.out.printf("\t%s: %s\n", attr.getAttributeDescriptor().getName(), "");
                        else
                            System.out.printf("\t%s: %s\n", attr.getAttributeDescriptor().getName(), (attr.getValue() instanceof Calendar)?(simpleFormat.format(((Calendar)attr.getValue()).getTime())):attr.getValue());
                    }
                }
      		} else {
			    System.out.printf("Querying for member born between %d and %d not found\n", startYear, endYear);
      		}
		}
		finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}
