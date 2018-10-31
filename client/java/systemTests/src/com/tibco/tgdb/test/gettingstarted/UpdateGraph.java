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

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Calendar;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;


/**
 * For a given member of the House, update the attributes
 * 
 * Usage : java UpdateGraph [options]
 * 
 *  where options are:
 *   -memberName <memberName> Required. Member name - "Napoleon Bonaparte"
 *   -crownName <crownName>   Optional. Name while reigning - "Napoleon XVIII"
 *   -crownTitle <crownTitle> Optional. Title while reigning - "King of USA"    
 *   -houseHead <houseHead>   Optional. Head of the house - true or false
 *   -yearBorn <yearBorn>     Optional. Year of birth - 2004
 *   -yearDied <yearDied>     Optional. Year of death - 2016 or null if still alive
 *   -reignStart <reignStart> Optional. Date reign starts (format dd MMM yyyy) - 20 Jan 2008 or null if never reigned
 *   -reignEnd <reignEnd>     Optional. Date reign ends (format dd MMM yyyy) - 08 Nov 2016 or null if never reigned or still reigning
 *   
 *  For instance to update the house member named "Napoleon Bonaparte" :
 *  java UpdateGraph -memberName "Napoleon Bonaparte" -crownName "Napoleon XVIII" -crownTitle "King of USA" -yearDied null -reignEnd "31 Jan 2016"
 *
 */
public class UpdateGraph {
	
	static String url = "tcp://127.0.0.1:8222";
	static String user = "napoleon";
	static String pwd = "bonaparte";
	
	static String memberName = null;
	static String crownName = null;
	static String crownTitle = null;
	static String yearBorn = null;
	static String yearDied = null;
	static String reignStart = null;
	static String reignEnd= null;
	static String houseHead = null;
	
	static void parseArgs(String[] args) throws Exception {
		for (int i=0; i<args.length; i++) {
			if (args[i].equals("-memberName"))
				memberName = args[i+1];
			else if (args[i].equals("-crownName"))
				crownName = args[i+1];
			else if (args[i].equals("-crownTitle"))
				crownTitle = args[i+1];
			else if (args[i].equals("-yearBorn"))
				yearBorn = args[i+1];
			else if (args[i].equals("-yearDied"))
				yearDied = args[i+1];
			else if (args[i].equals("-houseHead"))
				houseHead = args[i+1];
			else if (args[i].equals("-reignStart")) 
				reignStart = args[i+1];
			else if (args[i].equals("-reignEnd")) 
				reignEnd = args[i+1];
		}
	}
	
	public static void main(String[] args) throws Exception {
		
		parseArgs(args);
		
		if (memberName == null) {
			System.out.println("No house member to update.\nArguments example: -memberName \"Napoleon Bonaparte\" -crownName \"Grand Napoleon\" -crownTitle \"King of the world\" -reignStart \"8 Nov 2001\" -yearDied 2016");
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
			
			TGKey houseKey = gof.createCompositeKey("houseMemberType");
			
			houseKey.setAttribute("memberName", memberName);
			System.out.printf("Searching for member '%s'...\n",memberName);
      		TGEntity houseMember = conn.getEntity(houseKey, null);
      		if (houseMember != null) {
      			System.out.printf("House member '%s' found\n",houseMember.getAttribute("memberName").getAsString());
      			if (crownName != null)
      				houseMember.setAttribute("crownName", crownName);
      			if (crownTitle != null)
      				houseMember.setAttribute("crownTitle", crownTitle);
      			if (houseHead != null)
      				houseMember.setAttribute("houseHead", Boolean.parseBoolean(houseHead));
      			if (yearBorn != null)
      				houseMember.setAttribute("yearBorn", Integer.parseInt(yearBorn));
      			if (yearDied != null) { 
      				if (yearDied.equals("null"))
      					houseMember.setAttribute("yearDied", null);
      				else
      					houseMember.setAttribute("yearDied", Integer.parseInt(yearDied));
      			}
      			Calendar calReignStart = Calendar.getInstance();
      			if (reignStart != null) {
      				if (reignStart.equals("null"))
      					houseMember.setAttribute("reignStart", null);
      				else {
      					try {
      						calReignStart.setTime((new SimpleDateFormat("dd MMM yyyy").parse(reignStart)));
      						houseMember.setAttribute("reignStart", calReignStart); 
      					}
      					catch (ParseException e) {
      						throw new Exception("Member update failed - Wrong parameter: -reignStart format should be \"dd MMM yyyy\"");
      					}
      				}
      			}
      			Calendar calReignEnd = Calendar.getInstance();
      			if (reignEnd != null) {
      				if (reignEnd.equals("null"))
      					houseMember.setAttribute("reignEnd", null);
      				else {
      					try {
      						calReignEnd.setTime((new SimpleDateFormat("dd MMM yyyy").parse(reignEnd)));
      						houseMember.setAttribute("reignEnd", calReignEnd); 
      					}
      					catch (ParseException e) {
      						throw new Exception("Member update failed - Wrong parameter: -reignEnd format should be \"dd MMM yyyy\"");
      					}
      				}
      			}
      			
      			conn.updateEntity(houseMember);
      			conn.commit();
      			System.out.printf("House member '%s' updated successfully\n", memberName);	
      		} 
      		else {
      			System.out.printf("House member '%s' not found\n", memberName);
      		}
		}
		finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}
