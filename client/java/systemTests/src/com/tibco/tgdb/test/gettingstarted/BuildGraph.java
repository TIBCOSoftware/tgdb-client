package com.tibco.tgdb.test.gettingstarted;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Hashtable;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

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
 * Build the House of Bonaparte graph.
 * The House of Bonaparte is an imperial and royal European dynasty founded in 1804 by Napoleon I.
 * 
 * Each member has the following characteristics : 
 * -memberName : Name of member (primary key)
 * -crownName : Name while reigning
 * -crownTitle : Title while reigning
 * -houseHead : There is always a head of the house at a given time.
 * -yearBorn : Year of birth
 * -yearDied : Year of death
 * -reignStart : Date reign started
 * -reignEnd : Date reign ended
 * 
 * Usage : java BuildGraph
 * 
 */
public class BuildGraph {
	
	// Members of the House to be inserted in the database
	final static Object houseMemberData[][] = {
			// memberName, crownName, houseHead, yearBorn, yearDied, reignStart, reignEnd, crownTitle
			{ "Carlo Bonaparte", null, false, 1746, 1785, null, null, null },
			{ "Letizia Ramolino", null, false, 1750, 1836, null, null, null },
			{ "Joseph Bonaparte", "Joseph I", false, 1768, 1844, "6 Jun 1808", "11 Dec 1813", "King of Spain" },
			{ "Napoleon Bonaparte", "Napoleon I", false, 1769, 1821, "18 May 1804", "22 Jun 1815", "Emperor of the French" },
			{ "Lucien Bonaparte", null, false, 1775, 1840, null, null, null },
			{ "Elisa Bonaparte", "Elisa Bonaparte", false, 1777, 1820, "3 Mar 1809", "1 Feb 1814", "Grand Duchess of Tuscany" },
			{ "Louis Bonaparte", "Louis I", false, 1778, 1846, "5 Jun 1806", "1 Jul 1810", "King of Holland" },
			{ "Pauline Bonaparte", null, false, 1780, 1825, null, null, null },
			{ "Caroline Bonaparte", null, false, 1782, 1839, null, null, null },
			{ "Jerome Bonaparte", "Jerome I", false, 1784, 1860, "8 Jul 1807", "26 Oct 1813", "King of Westphalia" },
			{ "Marie Louise of Austria", null, false, 1791, 1847, null, null, "Empress Consort of the French" },
			{ "Josephine of Beauharnais", null, false, 1763, 1814, null, null, "Empress Consort of the French" },
			{ "Alexandre of Beauharnais", null, false, 1760, 1794, null, null, null },
			{ "Betsy Patterson", null, false, 1785, 1879, null, null, null },
			{ "Catharina of Wurttemberg", null, false, 1783, 1835, null, null, "Queen Consort of Westphalia" },
			{ "Francois Bonaparte", "Napoleon II", false, 1811, 1832, "22 Jun 1815", "7 Jul 1815", "Emperor of the French" },
			{ "Hortense of Beauharnais", null, false, 1783, 1837, null, null, "Queen Consort of Holland" },
			{ "Jerome Napoleon", null, false, 1805, 1870, null, null, null },
			{ "Prince Napoleon", null, false, 1822, 1891, null, null, null },
			{ "Louis Napoleon", "Napoleon III", true, 1808, 1873, "2 Dec 1852", "4 Sep 1870", "Emperors of the French" },
			{ "Napoleon-Louis Bonaparte", "Louis II", false, 1804, 1831, "1 Jul 1810", "13 Jul 1810", "King of Holland" },
			{ "Napoleon IV Eugene", null, true, 1856, 1879, null, null, null },
			{ "Napoleon V Victor", null, true, 1862, 1926, null, null, null },
			{ "Marie Clotilde Bonaparte", null, false, 1912, 1996, null, null, null },
			{ "Napoleon VI Louis", null, true, 1914, 1997, null, null, null },
			{ "Napoleon VII Charles", null, true, 1950, null, null, null, null },
			{ "Napoleon VIII Jean-Christophe", null, true, 1986, null, null, null, null },
			{ "Sophie Catherine Bonaparte", null, false, 1992, null, null, null, null } };

	// Relation among the members of the House to be inserted in the database
	final static Object houseRelationData[][] = {
			// From memberName, To memberName, relation type
			{ "Carlo Bonaparte", "Letizia Ramolino", "spouse" }, 
			{ "Carlo Bonaparte", "Joseph Bonaparte", "child" },
			{ "Letizia Ramolino", "Joseph Bonaparte", "child" }, 
			{ "Carlo Bonaparte", "Napoleon Bonaparte", "child" },
			{ "Letizia Ramolino", "Napoleon Bonaparte", "child" }, 
			{ "Carlo Bonaparte", "Lucien Bonaparte", "child" },
			{ "Letizia Ramolino", "Lucien Bonaparte", "child" }, 
			{ "Carlo Bonaparte", "Elisa Bonaparte", "child" },
			{ "Letizia Ramolino", "Elisa Bonaparte", "child" }, 
			{ "Carlo Bonaparte", "Louis Bonaparte", "child" },
			{ "Letizia Ramolino", "Louis Bonaparte", "child" }, 
			{ "Carlo Bonaparte", "Pauline Bonaparte", "child" },
			{ "Letizia Ramolino", "Pauline Bonaparte", "child" }, 
			{ "Carlo Bonaparte", "Caroline Bonaparte", "child" },
			{ "Letizia Ramolino", "Caroline Bonaparte", "child" }, 
			{ "Carlo Bonaparte", "Jerome Bonaparte", "child" },
			{ "Letizia Ramolino", "Jerome Bonaparte", "child" },

			{ "Napoleon Bonaparte", "Marie Louise of Austria", "spouse" },
			{ "Napoleon Bonaparte", "Francois Bonaparte", "child" },
			{ "Marie Louise of Austria", "Francois Bonaparte", "child" },

			{ "Napoleon Bonaparte", "Josephine of Beauharnais", "spouse" },

			{ "Alexandre of Beauharnais", "Josephine of Beauharnais", "spouse" },
			{ "Alexandre of Beauharnais", "Hortense of Beauharnais", "child" },
			{ "Josephine of Beauharnais", "Hortense of Beauharnais", "child" },

			{ "Louis Bonaparte", "Hortense of Beauharnais", "spouse" },
			{ "Louis Bonaparte", "Louis Napoleon", "child" },
			{ "Hortense of Beauharnais", "Louis Napoleon", "child" },
			{ "Louis Bonaparte", "Napoleon-Louis Bonaparte", "child" },
			{ "Hortense of Beauharnais", "Napoleon-Louis Bonaparte", "child" },

			{ "Jerome Bonaparte", "Betsy Patterson", "spouse" },
			{ "Jerome Bonaparte", "Jerome Napoleon", "child" },
			{ "Betsy Patterson", "Jerome Napoleon", "child" },

			{ "Jerome Bonaparte", "Catharina of Wurttemberg", "spouse" },
			{ "Jerome Bonaparte", "Prince Napoleon", "child" },
			{ "Catharina of Wurttemberg", "Prince Napoleon", "child" },

			{ "Louis Napoleon", "Napoleon IV Eugene", "child" },

			{ "Prince Napoleon", "Napoleon V Victor", "child" },

			{ "Napoleon V Victor", "Napoleon VI Louis", "child" },

			{ "Napoleon VI Louis", "Napoleon VII Charles", "child" },

			{ "Napoleon VII Charles", "Napoleon VIII Jean-Christophe", "child" },
			{ "Napoleon VII Charles", "Sophie Catherine Bonaparte", "child" } };

	
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
			if (houseMemberType == null)
				throw new Exception("Node type not found");
			TGNode houseMember;
			Hashtable<String, TGNode> houseMemberTable = new Hashtable<String, TGNode>();
			
			//
			// Insert node data into database
			//
			for (int i = 0; i < houseMemberData.length; i++) {
				houseMember = gof.createNode(houseMemberType);
				houseMember.setAttribute("memberName", houseMemberData[i][0]);
				houseMember.setAttribute("crownName", houseMemberData[i][1]);
				houseMember.setAttribute("houseHead", houseMemberData[i][2]);
				houseMember.setAttribute("yearBorn", houseMemberData[i][3]);
				houseMember.setAttribute("yearDied", houseMemberData[i][4]);
				houseMember.setAttribute("crownTitle", houseMemberData[i][7]);
				
				Calendar reignStart = Calendar.getInstance();
				if (houseMemberData[i][5] != null) {
					reignStart.setTime((new SimpleDateFormat("dd MMM yyyy").parse((String) houseMemberData[i][5])));
					houseMember.setAttribute("reignStart", reignStart);
				}
				else 
					houseMember.setAttribute("reignStart", null);

				Calendar reignEnd = Calendar.getInstance();
				if (houseMemberData[i][6] != null) {
					reignEnd.setTime((new SimpleDateFormat("dd MMM yyyy").parse((String) houseMemberData[i][6])));
					houseMember.setAttribute("reignEnd", reignEnd);
				}
				else 
					houseMember.setAttribute("reignEnd", null);
				
				conn.insertEntity(houseMember);
				conn.commit(); // Write data to database
				System.out.println(
						"Transaction completed for Node : " + houseMember.getAttribute("memberName").getAsString());
				houseMemberTable.put(houseMember.getAttribute("memberName").getAsString(), houseMember);
			}

			System.out.println("-------------------------------------------------");
			
			// 
			// Insert edge data into database
			// 
			TGNode houseMemberFrom;
			TGNode houseMemberTo;
			TGEdge houseRelation;
			TGEdge.DirectionType houseRelationDirection;
			for (int i = 0; i < houseRelationData.length; i++) {
				houseMemberFrom = houseMemberTable.get(houseRelationData[i][0]);
				houseMemberTo = houseMemberTable.get(houseRelationData[i][1]);
				houseRelationDirection = houseRelationData[i][2].equals("spouse") ? TGEdge.DirectionType.UnDirected
						: TGEdge.DirectionType.Directed;
				houseRelation = gof.createEdge(houseMemberFrom, houseMemberTo, houseRelationDirection);
				houseRelation.setAttribute("relType", houseRelationData[i][2]);
				conn.insertEntity(houseRelation);
				conn.commit();
				System.out.println(
						"Transaction completed for Edge : " + houseMemberFrom.getAttribute("memberName").getAsString()
								+ " to " + houseMemberTo.getAttribute("memberName").getAsString());
			}
			System.out.println("\nHouse of Bonaparte graph completed successfully");
		} 
		finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}