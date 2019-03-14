package com.tibco.tgdb.test;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.HashMap;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.utils.EntityUtils;

/**
 * Modified from example QueryGraph to demonstrate traversal filter
 * Uses example inithousedb and housetgdb config
 * Needs to run BuildGraph first
 *
 * Traversal starts from the nodes in the query results set.  
 * Starting from each edge of those nodes the traversal filter is applied.
 * If evaluation is true, it will apply the end filter to see the ending
 * condition is fulfilled.  If not, it will continue the traversal to the
 * next level. The final result set only contains the nodes and edges 
 * fulfilled all the conditions.
 *
 * The following reserved keyword are introduced for traversal filtering
 * @fromnodetype - string - the node desc name of the node we are getting the edge from
 * @tonodetype - string - the node desc name of the other end of the edge
 * @isfromedge - is the node where the edge is retrieve from on the from side of the edge. 1 - true, 0 - false
 * @fromnode.<attr name> - retrieve the from node attribute value
 * @tonode.<attr name> - retrieve the to node attribute value
 * @edge.<attr name> - retrieve the edge attribute value
 * @edgetype - string - edge desc name
 * @degree/depth - int - degree of separation or what we call the depth
 *                       both degree and depth are valid keywords
 *
 * e.g.  If we starts from Napoleon Bonaparte and get the offspring edge to Carlo Bonaparte,
 *       the isfromedge will be 0 because the edge is created from Carlo to Napoleon.
 *       But the fromnode is Napoleon because we are traversing from Napoleon to Carlo
 *       and therefore tonode is Carlo.  It may sound confusing.  May be a different 
 *       naming for fromnode and tonode can help.
 *
 * The isfromedge is used to control which direction of the edge you want to traverse.  
 * Using offspring edge as an example, if you only interested in traversing from
 * parent to child, you should specify isfromedge to 1. If you start from 
 * Napoleon Bonaparte and traverse offspring edge with isfromedge = 1, you will 
 * only get to Francois Bonaparte but not to Carlo or Letitia.
 *
 * The second argument of the new executeQuery method is not used right now.  The end
 * condition is required but the traversal condition is optional.
 *
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
	
	static int startYear = 1760;
    static int endYear = 1770;

    static boolean runTest1 = false;
    static boolean runTest2 = true;
    static boolean runTest3 = false;
    static boolean runTest4 = false;
    static boolean runTest5 = false;
	
	public static void main(String[] args) throws Exception {
		
		
        try {
            for (int i=0; i<args.length; i++) {
                if (args[i].equals("-startyear")) {
                    startYear = Integer.parseInt(args[i+1]);
                } else if (args[i].equals("-endyear")) {
                    endYear = Integer.parseInt(args[i+1]);
                } else if (args[i].equals("test1")) {
                    runTest1 = true;
                } else if (args[i].equals("test2")) {
                    runTest2 = true;
                } else if (args[i].equals("test3")) {
                    runTest3 = true;
                } else if (args[i].equals("test4")) {
                    runTest4 = true;
                } else if (args[i].equals("test5")) {
                    runTest5 = true;
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
			
			
            String queryString = null;
            String traverseString = null;
            String endString = null;
            TGResultSet resultSet = null;
            int dumpDepth = 5;
            int currDepth = 0;
            boolean dumpBreadth = false;
            boolean showAllPath = true;
            TGQueryOption option = TGQueryOption.createQueryOption();

            if (runTest1 == true) {
                //Simple query
                dumpBreadth = true;
			    System.out.printf("Querying for member born between %d and %d\n", startYear, endYear);
                queryString = "@nodetype = 'houseMemberType' and yearBorn > " + startYear + " and yearBorn < " + endYear + ";";
                resultSet = conn.executeQuery(queryString, null);
            } else if (runTest2 == true) {
            	//Identify a single path from Napoleon Bonaparte  to Francois Bonaparte
                queryString = "@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';";
                traverseString = "@edgetype = 'offspringEdge' and @isfromedge = 1 and @edge.birthOrder = 1 and @degree < 3" + ";";
                endString = "@tonodetype = 'houseMemberType' and @tonode.memberName = 'Francois Bonaparte'" + ";"; 
                resultSet = conn.executeQuery(queryString, null, traverseString, endString, null);
            } else if (runTest3 == true) {
            	//Identify all paths from Napoleon Bonaparte  to Napoleon IV Eugene with no traversal filter except 10 level deep restriction
                dumpDepth = 10;
                queryString = "@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';";
                endString = "@tonodetype = 'houseMemberType' and @tonode.memberName = 'Napoleon IV Eugene'" + ";"; 
                option.setTraversalDepth(10);
                resultSet = conn.executeQuery(queryString, null, null, endString, option);
            } else if (runTest4 == true) {
            	//Identify all paths from Napoleon Bonaparte  to Napoleon IV Eugene using only offspringEdge desc and within 10 level deep
                dumpDepth = 10;
                queryString = "@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';";
                traverseString = "@edgetype = 'offspringEdge' and @degree <= 10" + ";";
                endString = "@tonodetype = 'houseMemberType' and @tonode.memberName = 'Napoleon IV Eugene'" + ";"; 
                option.setTraversalDepth(10);
                resultSet = conn.executeQuery(queryString, null, traverseString, endString, option);
            } else if (runTest5 == true) {
            	//Identify specific path from Napoleon Bonaparte -> his parents -> Louis Bonaparte -> Louis Napoleon -> Napoleon IV Eugene
                dumpDepth = 10;
                queryString = "@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';";
                traverseString = "(@edgetype = 'offspringEdge' and @isfromedge = 0 and @degree = 1)" +
                                           "or (@edgetype = 'offspringEdge' and @isfromedge = 1 and @degree = 2)" +
                                           "or (@edgetype = 'offspringEdge' and @isfromedge = 1 and @degree = 3)" +
                                           "or (@edgetype = 'offspringEdge' and @isfromedge = 1 and @degree = 4)" +
                                            ";";
                endString = "@tonodetype = 'houseMemberType' and @tonode.memberName = 'Napoleon IV Eugene'" + ";"; 
                option.setTraversalDepth(10);
                resultSet = conn.executeQuery(queryString, null, traverseString, endString, option);
            }

      		if (resultSet != null) {
                while (resultSet.hasNext()) {
                    TGEntity houseMember = resultSet.next();
                    SimpleDateFormat simpleFormat = new SimpleDateFormat("dd MMM yyyy");
                    System.out.printf("House member '%s' found\n",houseMember.getAttribute("memberName").getAsString());
                    if (dumpBreadth) {
                        EntityUtils.printEntitiesBreadth((TGNode) houseMember, dumpDepth);
                    } else {
                        EntityUtils.printEntities(houseMember, dumpDepth, currDepth, "", showAllPath, new HashMap<Integer, TGEntity>());
                    }

                    /*
                    for (TGAttribute attr : houseMember.getAttributes()) {
                        if (attr.getValue() == null)
                            System.out.printf("\t%s: %s\n", attr.getAttributeDescriptor().getName(), "");
                        else
                            System.out.printf("\t%s: %s\n", attr.getAttributeDescriptor().getName(), (attr.getValue() instanceof Calendar)?(simpleFormat.format(((Calendar)attr.getValue()).getTime())):attr.getValue());
                    }
                    int i=0;
                    for (TGEdge edge: ((TGNode)houseMember).getEdges()) {
                            for (TGAttribute attr : edge.getAttributes()) {
                                if (attr.getValue() == null)
                                    System.out.printf("\t\t%d:%s: %s\n", i, attr.getAttributeDescriptor().getName(), "");
                                else
                                    System.out.printf("\t\t%d:%s: %s\n", i, attr.getAttributeDescriptor().getName(), (attr.getValue() instanceof Calendar)?(simpleFormat.format(((Calendar)attr.getValue()).getTime())):attr.getValue());
                            }
                            i++;
                    }
                    */
                }
      		} else {
			    System.out.printf("Query return no result set\n");
      		}
		}
		finally {
			if (conn != null)
				conn.disconnect();
		}
	}
}
