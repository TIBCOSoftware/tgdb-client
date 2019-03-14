package com.tibco.tgdb.test;

import java.util.Collection;
import java.util.List;

import org.apache.tinkerpop.gremlin.process.remote.RemoteConnection;
import org.apache.tinkerpop.gremlin.process.traversal.P;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal;
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

    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();
        
        //Following two lines function the same.  traversal() returns GraphTraversalSource
	    //EmptyGraph.instance().traversal().withRemote(conn);
        //GraphTraversalSource g = EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
//        GraphTraversalSource g = (GraphTraversalSource) EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
        //Pass in TGConnection instead of RemoteConnection from Gremlin
        //We may look into supporting RemoteConnection
        GraphTraversalSource g = EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
        GraphTraversal t = g.V();

		System.out.println("Test values");
        //simple query
        //This should return a list of primitive values
        List valueList = g.V().has("cdi", "cdiid", "172CDIXEAY44").values().toList();
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

		System.out.println("Test between conditions valueMap");
        valueList = g.V().has("cdi", "groupid", P.between(700, 701)).valueMap("itemid", "itemname", "groupid").toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test between conditions valueMap ends");
		System.out.println("");

		System.out.println("Test and/or steps");
        valueList = g.V().hasLabel("cdi").and(
        		__.or(__.has("cdiid","172CBAFEVZ08"),__.has("cdiid", "172CBAFFPU57")),
        		__.has("groupid",P.gte(1700))).valueMap().toList();
		for (Object value : valueList) {
			System.out.println(value);
		}
		System.out.println("Test and/or steps ends");

		/*
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
