/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 * File name : GremlinQueryTest.${EXT}
 * Created on: 11/07/2018
 * Created by: vincent
 * SVN Id: $Id: GremlinQueryTest.java 3159 2019-04-26 21:16:15Z vchung $
 */

package com.tibco.tgdb.test;

import java.util.Collection;
import java.util.List;

import org.apache.tinkerpop.gremlin.process.remote.RemoteConnection;
import org.apache.tinkerpop.gremlin.process.traversal.P;
import org.apache.tinkerpop.gremlin.process.traversal.Path;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal;
import org.apache.tinkerpop.gremlin.structure.Column;
import org.apache.tinkerpop.gremlin.structure.T;
import org.apache.tinkerpop.gremlin.structure.util.empty.EmptyGraph;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.gremlin.GraphTraversalSource;
import com.tibco.tgdb.gremlin.__;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.query.TGResultSet;

//This test uses data from OpenFlight, Brembo and Bonaparte
public class GremlinQueryTest {
	public String url = "tcp://scott@localhost:8222";
    public String passwd = "scott";
    public TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    public int depth = 5;
    public int printDepth = 5;
    public int resultCount = 100;
    public int edgeLimit = 0;

	public static void main(String[] args) throws Exception {
		String url = "tcp://scott@localhost:8222";
		String passwd = "scott";
		TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
		List valueList = null;
		int i = 1;

    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();
        
        //Following two lines function the same.  traversal() returns GraphTraversalSource
	    //EmptyGraph.instance().traversal().withRemote(conn);
        //GraphTraversalSource g = EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
        //GraphTraversalSource g = (GraphTraversalSource) EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
        //Pass in TGConnection instead of RemoteConnection from Gremlin
        //We may look into supporting RemoteConnection
        GraphTraversalSource g = EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
        GraphTraversal t = g.V();

        //Condition tests
		System.out.println("Test values");
        //simple query
        //This should return a list of primitive values
        valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").values().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test values ends");
		System.out.println("");

		System.out.println("Test values count");
        valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").values().count().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test values count ends");
		System.out.println("");

		System.out.println("Test values fold");
		//returns a list with a single element of another list
        valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").values().fold().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test values fold ends");
		System.out.println("");

		System.out.println("Test valueMap all");
		//returns a list with a single element of a map of key/value
        valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test valueMap all ends");
		System.out.println("");

		System.out.println("Test valueMap select");
		//returns a list with a single element of a map of key/value
        valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").valueMap("itemname", "oops", "itemid").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test valueMap  select ends");
		System.out.println("");

		System.out.println("Test valueMap count");
		//returns a list with a single element of a map of key/value
        valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").valueMap().count().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test valueMap count ends");
		System.out.println("");

		System.out.println("Test valueMap fold");
		//returns a list with a single element of a map of key/value
        valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").valueMap().fold().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test valueMap fold ends");
		System.out.println("");

		System.out.println("Test and-and condition");
		//return TGNode directly
        valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").has("groupid", 2200).toList();
		for (Object value : valueList) {
			TGNode node = (TGNode)value;
			Collection<TGAttribute> attrs = node.getAttributes();
			for (TGAttribute attr : attrs) {
				System.out.printf("Attr name : %s, value : %s\n", attr.getAttributeDescriptor().getName(),
						attr.getValue().toString());
			}
		}
		System.out.println("Test and-and condition ends");
		System.out.println("");

		System.out.println("Test cdi batch");
        valueList = g.V().has("cdibatch", "batchid", "17012LXFP342").values("batchid", "oops").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test cdi batch ends");
		System.out.println("");

		System.out.println("Test pagerank");
        valueList = g.V().pageRank().valueMap("@pagerank", "itemname").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test pagerank ends");
		System.out.println("");

		System.out.println("Test conditions");
        valueList = g.V().has("cdi", "cdiid", P.eq("172CBAFEVZ08").or(P.eq("172CBAFFPU57")).or(P.eq("172CBAFGLK39"))).toList();
		for (Object value : valueList) {
			TGNode node = (TGNode)value;
			Collection<TGAttribute> attrs = node.getAttributes();
			for (TGAttribute attr : attrs) {
				System.out.printf("Attr name : %s, value : %s\n", attr.getAttributeDescriptor().getName(),
						attr.getValue().toString());
			}
			System.out.println("");
		}
		System.out.println("Test conditions ends");
		System.out.println("");

		System.out.println("Test conditions valueMap");
        valueList = g.V().has("cdi", "cdiid", P.eq("172CBAFEVZ08").or(P.eq("172CBAFFPU57")).or(P.eq("172CBAFGLK39"))).valueMap("itemid", "itemname").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test conditions valueMap ends");
		System.out.println("");

        /*FIXME: Disable this test 4/16/2019.  Takes up too much memory. Revisit later. 
		System.out.println("Test between conditions valueMap");
         * Can add a non-unique index for groupid to test out the index and
         * at the same time reduce memory usage due to scanning all the entities in the db.
        valueList = g.V().has("cdi", "groupid", P.between(700, 701)).valueMap("itemid", "itemname", "groupid").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test between conditions valueMap ends");
		System.out.println("");
        */

        //FIXME:  Need to rework index check logic on the server side.
        //This query should do a unique get instead of the 
        //range up prefetch on the server side.
		/* disable for now to speed up the test 4/17/2019
		System.out.println("Test and/or steps");
//        List valueList = g.V().hasLabel("cdi").and(
        valueList = g.V().hasLabel("cdi").and(
        		__.or(__.has("cdiid","172CBAFEVZ08"),__.has("cdiid", "172CBAFFPU57")),
        		__.has("groupid",P.gte(1700))).valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test and/or steps ends");
		System.out.println("");
		*/

		System.out.println("Test and/or with empty or steps");
        valueList = g.V().hasLabel("cdi").and(
        		__.or(__.has("cdiid","172CBAFEVZ08"),__.has("cdiid", "172CBAFFPU57"))).or().
        		hasLabel("cdibatch").and(__.has("batchid","17012LXFP342")).valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test and/or with empty or steps ends");
		System.out.println("");

		System.out.println("Test or after V");
        valueList = g.V().or(__.hasLabel("cdi").and(
        		__.or(__.has("cdiid","172CBAFEVZ08"),__.has("cdiid", "172CBAFFPU57"))),
        		__.hasLabel("cdibatch").and(__.has("batchid","17012LXFP342"))).valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test or after V ends");
		System.out.println("");

        //Traversal tests
		System.out.println("Traversal 1");
		valueList = g.V().has("cdi", "cdiid", "172CBAFEVZ08").outE("produces").has("quantity", P.gt(5)).
				inV().has("groupid", P.gt(1000)).out("produces").valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 1 ends");
		System.out.println("");

		System.out.println("Traversal 2");
		valueList = g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXEI727")).outE("contains").simplePath().
				inV().has("groupid", P.neq(1000)).limit(10).out("produces").valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 2 ends");
		System.out.println("");

		System.out.println("Traversal 3");
		valueList = g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXEI727")).outE("contains").inV().valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 3 ends");
		System.out.println("");
		
		System.out.println("Traversal 4");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").valueMap("quantity").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 4 ends");
		System.out.println("");

		System.out.println("Traversal 5");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
				outE("produces").inV().path().by("cdiid").by("quantity").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 5 ends");
		System.out.println("");

		System.out.println("Traversal 6 multiple folds");
		valueList = g.V().has("cdi", "cdiid", "172CBAFEVZ08").outE("produces").has("quantity", P.gt(5)).
				inV().has("groupid", P.gt(1000)).out("produces").values().fold().fold().fold().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 6 multiple folds ends");
		System.out.println("");

		System.out.println("Traversal 7.0");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
				or(P.eq("172CBAFFPU57")))).
				outE("produces").inV().outE("produces").inV().valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 7.0 ends");
		System.out.println("");

		System.out.println("Traversal 7");
		/* Not working. -- Need to investigate -- similar behavior using gremlin console also
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18").or(__.has("cdiid", "172CBAFEVZ08")).
				or(__.has("cdiid", "172CBAFFPU57"))).
				outE("produces").inV().outE("produces").inV().path().by("cdiid").by("quantity").toList();
				*/
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
				or(P.eq("172CBAFFPU57")))).
				outE("produces").inV().outE("produces").inV().path().by("cdiid").by("quantity").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 7 ends");
		System.out.println("");

		/*
		System.out.println("Traversal 7.1");
		//Not supporting by just entity in the path
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
				or(P.eq("172CBAFFPU57")))).
				outE("produces").inV().outE("produces").inV().path().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 7.1 ends");
		System.out.println("");
		*/

		System.out.println("Traversal 7.2");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
				or(P.eq("172CBAFFPU57")))).
				outE("produces").inV().outE("produces").inV().path().count().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 7.2 ends");
		System.out.println("");

		System.out.println("Traversal 7.3");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
				or(P.eq("172CBAFFPU57")))).
				outE("produces").inV().outE("produces").inV().simplePath().path().by("cdiid").by("quantity").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 7.3 ends");
		System.out.println("");
		
		System.out.println("Traversal 7.4");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
				or(P.eq("172CBAFFPU57")))).
				outE("produces").inV().outE("produces").inV().simplePath().path().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 7.4 ends");
		System.out.println("");

		System.out.println("Traversal 7.5");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
				or(P.eq("172CBAFFPU57")))).
				outE("produces").inV().outE("produces").inV().simplePath().path().count().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal 7.5 ends");
		System.out.println("");
		
		System.out.println("Unsupported step 1");
		valueList = g.V().hasLabel("cdi").has("cdiid", "172CDIXDQC18").outE("produces").inV().outE("pduce").out().path().by("cdiid").by("quantity").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Unsupported step 1 ends");
		System.out.println("");
		
		System.out.println("Traversal Flight data 1");
		//Not supporting by just entity in the path
		valueList = g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
				outE("routeType").inV().outE("routeType").inV().path().by("iataCode").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal Flight data 1 ends");
		System.out.println("");

		/*
		System.out.println("Traversal Flight data 2");
		valueList = g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
				outE("routeType").inV().outE("routeType").inV().outE("routeType").inV().
				outE("routeType").inV().has("iataCode", "CDG").path().count().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal Flight data 2 ends");
		System.out.println("");
		*/

		System.out.println("Traversal Flight data 2.0");
		valueList = g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
				outE("routeType").inV().
				outE("routeType").inV().
				outE("routeType").inV().
				simplePath().
				has("iataCode", "CDG").path().by("iataCode").toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i) + " " + value);
			i++;
		}
		System.out.println("Traversal Flight data 2.0 ends");
		System.out.println("");

		System.out.println("Traversal Flight data 2.0 - count ");
		valueList = g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
				outE("routeType").has("iataCode", "UA").inV().
				outE("routeType").has("iataCode", "UA").inV().
				outE("routeType").has("iataCode", "UA").inV().
				outE("routeType").has("iataCode", "UA").inV().
//				outE("routeType").inV().
//				outE("routeType").inV().
//				outE("routeType").inV().
//				outE("routeType").inV().
				simplePath().
				has("iataCode", "CDG").path().count().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Traversal Flight data 2.0  - count ends");
		System.out.println("");

		System.out.println("Traversal Flight data 2.1");
		valueList = g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
				outE("routeType").inV().
				outE("routeType").inV().
				outE("routeType").inV().
				has("iataCode", "CDG").valueMap().toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i) + " " + value);
			i++;
		}
		System.out.println("Traversal Flight data 2.1 ends");
		System.out.println("");

		System.out.println("Traversal Flight data 2.2");
		valueList = g.V().has("airportType", "iataCode", "SFO").
				out("routeType").
				out("routeType").
				//has("iataCode", "CDG").
				path().by("iataCode").toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i) + " " + value);
			i++;
		}
		System.out.println("Traversal Flight data 2.2 ends");
		System.out.println("");

        //Gremlin string tests
		System.out.println("Traversal string 1");
		TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().hasLabel('cdi').has('cdiid', '172CDIXDQC18').outE('produces').inV().outE('produces').inV().path();", null);
		for (Object value : resultSet.toCollection()) {
			System.out.println(value);
		}
		System.out.println("Traversal string 1 ends");
		System.out.println("");
		
		System.out.println("Traversal bad string 1");
		resultSet = conn.executeQuery("gremlin://g.V().hasLabel('cdi').limitt(10).has('cdiid', '172CDIXDQC18').outE('produces').inV().outE('produces').inV().path();", null);
		for (Object value : resultSet.toCollection()) {
			System.out.println(value);
		}
		System.out.println("Traversal bad string 1 ends");
		System.out.println("");
		
        //Aggregation tests
		System.out.println("Aggregation raw data");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
				outE("produces").inV().outE().values("quantity").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Aggregation raw data ends");
		System.out.println("");

		System.out.println("Sum");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
				outE("produces").inV().outE().values("quantity").sum().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Sum ends");
		System.out.println("");

		System.out.println("Max");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
				outE("produces").inV().outE().values("quantity").max().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Max ends");
		System.out.println("");

		System.out.println("Min");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
				outE("produces").inV().outE().values("quantity").min().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Min ends");
		System.out.println("");

		System.out.println("Mean");
		valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
				outE("produces").inV().outE().values("quantity").mean().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Mean ends");
		System.out.println("");

        //Edge tests
		/*
		System.out.println("Get all edges");
		valueList = g.E().toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i++) + " " + value);
		}
		System.out.println("Get all edges ends");
		System.out.println("");
		*/

		/*
		System.out.println("Get all 'produces' edges");
		valueList = g.E().hasLabel("produces").has("quantity", P.gt(30)).toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i++) + " " + value);
		}
		System.out.println("Get all 'produces' edges ends");
		System.out.println("");
		*/

		System.out.println("Get all UA edges");
		valueList = g.E().hasLabel("routeType").has("iataCode", "UA").valueMap().toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i++) + " " + value);
		}
		System.out.println("Get all UA edges ends");
		System.out.println("");
		
		System.out.println("Traverse from UA edge to node to SW edge");
		valueList = g.E().hasLabel("routeType").has("iataCode", "UA").
				inV().outE("routeType").has("iataCode", "SW").
				path().by("iataCode").toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i++) + " " + value);
		}
		System.out.println("Traverse from UA edge to node to SW edge");
		System.out.println("");
		
		System.out.println("Traverse from UA edge to node");
		valueList = g.E().hasLabel("routeType").has("iataCode", "UA").
				inV().
				path().by("iataCode").toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i++) + " " + value);
		}
		System.out.println("Traverse from UA edge to node");
		System.out.println("");
		
		System.out.println("Out edges from 'SFO'");
		valueList = g.V().hasLabel("airportType").has("iataCode", "SFO").
				outE("routeType").values("iataCode").toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i++) + " " + value);
		}
		System.out.println("Out edges from 'SFO' ends");
		System.out.println("");
		
		System.out.println("In edges from 'SFO'");
		valueList = g.V().hasLabel("airportType").has("iataCode", "SFO").
				inE("routeType").values("iataCode").toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i++) + " " + value);
		}
		System.out.println("In edges from 'SFO' ends");
		System.out.println("");
		
		System.out.println("Flight from SFO to RNO");
		valueList = g.V().hasLabel("airportType").has("iataCode", "SFO").
				outE("routeType").inV().has("airportType", "iataCode", "RNO").path().by("iataCode").toList();
		i = 1;
		for (Object value : valueList) {
			System.out.println(String.valueOf(i++) + " " + value);
		}
		System.out.println("Flight from SFO to RNO");
		System.out.println("");

        //Test for client side validation by gremlib mincore. This is no validation in mincore.
        //valueList = g.V().has("cdi", "groupid", 1000).group().by(T.label).by(__.count()).by("foo").toList();
        //valueList = g.V().group().by(T.label).by(__.count()).by("foo");
        //Column.keys;
		/*
		 * Followings are not active cases
        valueList = g.V().has("houseMemberType", "yearBorn", P.between(1906, 1936)).toList();

        valueList = g.V().has("airport","code",P.eq("SFO").or(P.eq("JFK"))).toList();
        valueList = g.V().hasLabel("airport").and(__.or(__.has("code","JFK"),__.has("code", "SFO")),__.has("runways",P.gt(1))).valueMap().toList();

        int startYear = 1846;
        int endYear = 1905;
        valueList = g.V().has("houseMemberType", "yearBorn", P.gt(startYear).and(P.lt(endYear))).toList();

		System.out.println("Predicate");
        valueList = g.V().has("airport","code",P.eq("SFO").or(P.eq("JFK"))).toList();
		System.out.println("Predicate ends");


		//Path traversal query
        valueList = g.V().
        		has("airport", "code", "AUS").
        		repeat(__.out().simplePath()).
        		emit().
        		times(5).
        		has("airport", "code", "AGR").
        		path().by("code").
        		limit(10).toList();
		for (Object value : valueList) {
			System.out.println(value);
		}

        valueList = g.V().
        		has("airport", "code", "AUS").
        		repeat(__.out()).
        		emit().
        		times(5).
        		has("airport", "code", "AGR").
        		path().by("code").
        		limit(10).toList();
		for (Object value : valueList) {
			System.out.println(value);
		}

		valueList = g.V().has("code","AUS").
				match(__.as("aus").values("runways").as("ausr"), 
						__.as("aus").out("route").as("outa").values("runways").as("outr")
				.where("ausr",P.eq("outr"))).
				select("outa").valueMap().select("code","runways").toList();

		//Referencing system id
		g.V().hasId(6).out().has(T.id,P.lt(46)).path().by("code").toList();
		
		//another way to specify the vertex type/label
		valueList = g.V().where(__.label().is(P.eq("airport"))).count().toList();
		

		valueList = g.V().hasLabel("airport").has("city", P.between("Dal","Dat")).
				values("city").order().dedup().toList();

		valueList = g.V().has("airport","code","AUS").out().
				not(__.where(__.in("contains").has("code","US"))).
				valueMap("code","city").toList();
		
		valueList = g.V().has("airport","code","AUS").out().
				where(__.values("runways").is(P.gt(6)).or().values("runways").is(4)).
				valueMap("code","runways").toList();

		valueList = g.V().has("code","AUS").as("a").out().as("b").
				filter(__.select("a","b").by("runways").where("a",P.eq("b"))).
				valueMap("code","runways").toList();

		valueList = g.V().has("airport","city","London").as("a","r").in("contains").as("b").
				where("a",P.eq("b")).by("country").by("code").
				select("a","r","b").by("code").by("region").toList();


		//Bonaparte queries
        valueList = g.V().has("houseMemberType", "yearBorn", P.between(startYear, endYear)).toList();

        valueList = g.V().has("houseMemberType", "yearBorn", P.gt(startYear).and(P.lt(endYear))).toList();

        //In progress, model conv UI query
        valueList = g.V().has("project", "cost", P.gte(500)).in("project_school").in("school_state").has("getstate", "uuid",
        	    		"USA-California").toList();

        //select and airport if outgoing routes > 100 and outgoing route distance > 1000 and the destination airport has > 7 runways
        valueList = g.V().hasLabel("airport").where(
        		__.and(__.out("route").count().is(P.gt(100)), 
        			__.outE("route").has("dist", P.gt(1000)).inV().has("runways",P.gt(7)))
        		).values("code").toList();

        //One way of using 'and' step
        valueList = g.V().where(__.outE("created").and().outE("knows")).values("name").toList();
        
        //select an edge to return
        valueList = g.V().has("code","MIA").outE().as("e").inV().has("code","DFW").select("e").toList();

        //select 10 airports and return result with a count of 10 and props of each of the 10 airports
        valueList = g.V().hasLabel("airport").limit(10).fold().
        		project("count","fields").
        		by(__.unfold().count()).
        		by(__.unfold().project("props").by(__.values("code", "desc")).fold()).toList();

        valueList = g.V().hasLabel("airport").and(__.has("region","US-TX"),__.has("longest",P.gte(12000))).values("code").toList();
        
        valueList = g.V().has("region","US-TX").has("longest",P.gte(12000)).toList();
        */

		conn.disconnect();
        System.out.println("Disconnected.");
    }
}
