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
 * File name : IndexTest.${EXT}
 * Created on: 04/24/2019
 * Created by: suresh
 * SVN Id: $Id: IndexTest.java 3148 2019-04-26 00:35:38Z sbangar $
 */

package com.tibco.tgdb.test;

import java.io.BufferedReader;
import java.io.FileReader;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.Iterator;
import java.util.List;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.query.TGQueryOption;

/*
 * Run the test with the following arguments:
 * -afile <path>/actors.csv -all
 * this will run all the test cases.
 * 
 * When run the test the first time, db needs to be initialized.
 * Run the test with the following arguments:
 * -afile <path>/actors.csv -i -all
 * 
 */
public class IndexTest {
	private TGConnection conn = null;
	private TGGraphObjectFactory gof = null;
	private TGNodeType testNodeType = null;
	private String url = "tcp://scott@localhost:8222";
    private String passwd = "scott";
    private TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    private int depth = 5;
    private int printDepth = 5;
    private int resultCount = 100;
    private int edgeLimit = 0;
    private boolean initdb = false;
    private boolean testGet = false;
    private boolean testUnique = false;
    private boolean testPartialUnique = false;
    private boolean testNonUnique = false;
    private boolean testPartialNonUnique = false;
    private boolean testNonUniqueRangeUp = false;
    private boolean testNonUniqueRangeDown = false;
    private boolean testNonUniqueRange = false;
    private boolean testFullScan = false;
    private boolean testAll = false;
    private String actorsFile = "./actors.csv";
    private String[] cities = {"Palo Alto", "San Francisco", "San Mateo", 
    		"Belmont", "Daly City", "Foster City", "San Carlos",
    		"Burlingame", "Millbrae", "Menlo Park"
    };
    private String[] citycodes = {"PA", "SF", "SM", "BL", "DC", "FC", "SC", "BG", "MB", "MP"};
    private List<String> actorProfiles = new ArrayList<String>();
    private int numActors = 0;
    private int commitSize = 100;

    private String getStringValue(Iterator<String> argIter) {
    	while (argIter.hasNext()) {
    		String s = argIter.next();
    		return s;
    	}
    	return null;
    }
    
    private String getStringValue(Iterator<String> argIter, String defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
    		return s;
    	}
    }

    private int getIntValue(Iterator<String> argIter, int defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
    		try {
    			int i = Integer.valueOf(s);
    			return i;
    		} catch (NumberFormatException e) {
    			System.out.printf("Invalid number : %s\n", s);
    		}
    		return defaultValue;
    	}
    }

    private boolean getBoolValue(Iterator<String> argIter, boolean defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
			boolean b = Boolean.valueOf(s);
			return b;
    	}
    }

    private void getArgs(String[] args) {
    	List<String> argList = Arrays.asList(args);
    	Iterator<String> argIter = argList.iterator();
    	while (argIter.hasNext()) {
    		String s = argIter.next();
    		System.out.printf("Arg : \"%s\"\n", s);
    		if (s.equalsIgnoreCase("-url")) {
    			url = getStringValue(argIter, "tcp://scott@localhost:8222");
    		} else if (s.equalsIgnoreCase("-password") || s.equalsIgnoreCase("-pw")) {
    			passwd = getStringValue(argIter, "scott");
    		} else if (s.equalsIgnoreCase("-loglevel") || s.equalsIgnoreCase("-ll")) {
    			String ll = getStringValue(argIter, "Debug");
    			try {
    				logLevel = TGLogger.TGLevel.valueOf(ll);
    			} catch(IllegalArgumentException e) {
    				System.out.printf("Invalid log level value '%s'...ignored\n", ll);
    			}
    		} else if (s.equalsIgnoreCase("-init") || s.equalsIgnoreCase("-i")) {
    			initdb = true;
    		} else if (s.equalsIgnoreCase("-testall") || s.equalsIgnoreCase("-all")) {
    			testAll = true;
    		} else if (s.equalsIgnoreCase("-get") || s.equalsIgnoreCase("-g")) {
    			testGet = true;
    		} else if (s.equalsIgnoreCase("-uquery") || s.equalsIgnoreCase("-u")) {
    			testUnique = true;
    		} else if (s.equalsIgnoreCase("-upquery") || s.equalsIgnoreCase("-up")) {
    			testPartialUnique = true;
    		} else if (s.equalsIgnoreCase("-nuquery") || s.equalsIgnoreCase("-nu")) {
    			testNonUnique = true;
    		} else if (s.equalsIgnoreCase("-nupquery") || s.equalsIgnoreCase("-nup")) {
    			testPartialNonUnique = true;
    		} else if (s.equalsIgnoreCase("-nurangeup") || s.equalsIgnoreCase("-nurup")) {
    			testNonUniqueRangeUp = true;
    		} else if (s.equalsIgnoreCase("-nurangedown") || s.equalsIgnoreCase("-nurdn")) {
    			testNonUniqueRangeDown = true;
    		} else if (s.equalsIgnoreCase("-nurange") || s.equalsIgnoreCase("-nur")) {
    			testNonUniqueRange = true;
    		} else if (s.equalsIgnoreCase("-fullscan") || s.equalsIgnoreCase("-fs")) {
    			testFullScan = true;
    		} else if (s.equalsIgnoreCase("-actorsfile") || s.equalsIgnoreCase("-afile")) {
    			actorsFile = getStringValue(argIter, actorsFile);
    		} else if (s.equalsIgnoreCase("-commitsize") || s.equalsIgnoreCase("-cc")) {
    			commitSize = getIntValue(argIter, commitSize);
    		} else {
    			System.out.printf("Skip argument %s\n", s);
    		}
    	}
    }

    private TGNode createNode(TGGraphObjectFactory gof, TGNodeType nodeType) {
    	if (nodeType != null) {
    		return gof.createNode(nodeType);
    	} else {
    		return gof.createNode();
    	}
    }
    
    void startup() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);
        conn.connect();

        gof = conn.getGraphObjectFactory();
       	testNodeType = null;
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        	throw new Exception("Cannot create graph object factory");
        }

       	TGGraphMetadata gmd = conn.getGraphMetadata(false);
        testNodeType = gmd.getNodeType("testnode");
        if (testNodeType != null) {
        	System.out.printf("'testnode' type is found with %d attributes\n\n", testNodeType.getAttributeDescriptors().size());
        } else {
        	System.out.println("'testnode' type is not found from meta data fetch");
        }
    }

    void shutdown() throws Exception {
    	if (conn != null) {
    		conn.disconnect();
    	}
    }

    private void createActorsList() throws Exception {
    	String file = actorsFile;
    	String line;
    	int actorCount = 0;
    	
    	BufferedReader br = new BufferedReader(new FileReader(file)); 
        while ((line = br.readLine()) != null) {
            String[] tokens = line.split("~");
            if (tokens == null) {
            	continue;
            }
            if (tokens.length == 0) {
            	continue;
            }
            int recordSize = tokens.length;
            String tstr = new String();
            for (int i=0; i<4; i++) {
            	if (i >= recordSize) {
            		actorProfiles.add(null);
            		tstr = tstr.concat(" ");
            	} else {
            		String tok = tokens[i];
            		if (tok.length() == 0) {
            			actorProfiles.add(null);
            			tstr = tstr.concat(" ");
            		} else {
            			if(i == 2) {
            				String dobstr = tok;
            				int datepos = dobstr.lastIndexOf(", ");
            				tok = dobstr.substring(datepos + 2);
            			}
            			actorProfiles.add(tok);
            			tstr = tstr.concat(tok);
            		}
            	}
            	tstr = tstr.concat(",");
            }
            System.out.println(">> " + tstr);
            actorCount++;
        }
        br.close();
        System.out.printf("Actors processed : %d\n", actorCount); 
    }

    /*
     * The data has total of 10K nodes.  It's partitioned into 10 1K sections.
     * In each section, all nodes have the same city name with a 1-1000 sequence
     * number.  Together form a unique key for the entire 10K data set.
     * Within each 1K section, the times10 attribute store a 10 multiple sequence
     * number.  The citycode and times10 attibutes form a non-unique key.
     * Finally, a list of 1000 actor information is repeated for each 1K section
     * of the data set.  It's used to test query filtering beyond the initial
     * index based prefetch.
     */
    void initDB() throws Exception {
    	if (!initdb) {
    		return;
    	}
    	System.out.println("Initializing DB with data");
    	createActorsList();
    	for (int i=0; i<cities.length; i++) {
    		String city = cities[i];
    		String cc = citycodes[i];
    		for (int j=0; j<1000; j++) {
    			TGNode node = createNode(gof, testNodeType);
				node.setAttribute("cityseqno1k", city+(j+1));
				node.setAttribute("city", city);
				node.setAttribute("seqno1k", j+1);
				int times10 = 10 * ((j/10) + 1);
				node.setAttribute("citycode10", cc + times10);
				node.setAttribute("citycode", cc);
				node.setAttribute("times10", times10);
				node.setAttribute("actor", actorProfiles.get(4 * j));
				node.setAttribute("movie", actorProfiles.get(4 * j + 1));
				String str = actorProfiles.get(4 * j + 2);
				if (str != null) {
					node.setAttribute("birthyear", Integer.valueOf(str));
				}
				str = actorProfiles.get(4 * j + 3);
				if (str != null) {
					node.setAttribute("birthcity", str);
				}
				conn.insertEntity(node);
				if (((j + 1) % commitSize) == 0) {
					conn.commit();
					System.out.printf("Committed %d nodes\n", i * 1000 + (j + 1));
				}
    		}
    	}
    	System.out.println("Finished initializing DB with data\n");
    }

    void testGetByKey() throws Exception {
    	if (!testGet && !testAll) {
    		System.out.println("Skipped getEntity test\n");
    		return;
    	}
    	System.out.println("getEntity test start");
    	System.out.println("getEntity with key 'Palo Alto500'");
	    TGKey key = gof.createCompositeKey("testnode");
	    key.setAttribute("cityseqno1k", "Palo Alto500");
	    TGEntity entity = conn.getEntity(key, TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (entity == null) {
	    	System.out.println("getEntity test failed - cannot lookup 'Palo Alto1000'");
	    } else {
	    	printEntity(entity);
	    	System.out.println("getEntity with key 'Palo Alto500' successful");
	    }

    	System.out.println("getEntity with compound key ('San Mateo', 123)");
	    key = gof.createCompositeKey("testnode");
	    key.setAttribute("city", "San Mateo");
	    key.setAttribute("seqno1k", "123");
	    entity = conn.getEntity(key, TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (entity == null) {
	    	System.out.println("getEntity test failed - cannot lookup ('San Mateo', 123)");
	    } else {
	    	printEntity(entity);
	    	System.out.println("getEntity with compound key ('San Mateo', 123) successful");
	    }
    	System.out.println("getEntity test end\n");
    }

    public void testUniqueQuery() throws Exception {
    	if (!testUnique && !testAll) {
    		System.out.println("Skipped unique query test\n");
    		return;
    	}
	    int count = 0;
    	System.out.println("Unique query test start");
    	System.out.println("Query using a single key - Belmont312");
	    TGResultSet<TGEntity> resultSet = conn.executeQuery("@nodetype='testnode' and cityseqno1k = 'Belmont312';", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Unique query returned %d nodes - %s\n", count, count == 1 ? "correct" : "wrong");
    	System.out.println("Query using a compound key - ('San Carlos', 973)");
	    count = 0;
	    resultSet = conn.executeQuery("@nodetype='testnode' and city = 'San Carlos' and seqno1k = 973;", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Unique query returned %d nodes - %s\n", count, count == 1 ? "correct" : "wrong");
	    System.out.println("Unique query test end\n");
    }

    public void testPartialUniqueQuery() throws Exception {
    	if (!testPartialUnique && !testAll) {
    		System.out.println("Skipped partial unique query test\n");
    		return;
    	}
    	/* This functionality is not working yet.
	    System.out.println("Partial unique query test start");
	    TGResultSet resultSet = conn.executeQuery("@nodetype='testnode' and city = 'Belmont';", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet == null) {
	    	System.out.println("Result set is null - failed(known issue)");
	    	System.out.println("Partial unique query test end\n");
	    	return;
	    }
	    if (resultSet.hasNext() == false) {
	    	System.out.println("Partial unique query returned no result - failed");
	    	System.out.println("Partial unique query test end\n");
	    	return;
	    }
	    while (resultSet.hasNext()) {
	    	TGEntity entity = resultSet.next();
	    	printEntity(entity);
	    }
	    System.out.println("Partial unique query test end\n");
	    */
    }

    public void testNonUniqueQuery() throws Exception {
    	if (!testNonUnique && !testAll) {
    		System.out.println("Skipped non-unique query test\n");
    		return;
    	}
	    System.out.println("Non-unique query test start");
    	System.out.println("Query using a single key - BL320");
	    int count = 0;
	    TGResultSet<TGEntity> resultSet = conn.executeQuery("@nodetype='testnode' and citycode10 = 'BL320';", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Result set has %d nodes - %s\n", count, count == 10 ? "correct" : "wrong");
    	System.out.println("Query using a compound key - ('BL', 320)");
	    count = 0;
	    resultSet = conn.executeQuery("@nodetype='testnode' and citycode = 'BL' and times10 = 320;", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Result set has %d nodes - %s\n", count, count == 10 ? "correct" : "wrong");
	    System.out.println("Non-unique query test end\n");
    }

    public void testPartialNonUniqueQuery() throws Exception {
    	if (!testPartialNonUnique && !testAll) {
    		System.out.println("Skipped partial non-unique query test\n");
    		return;
    	}
    }

    public void testNonUniqueRangeUpQuery() throws Exception {
    	if (!testNonUniqueRangeUp && !testAll) {
    		System.out.println("Skipped non-unique greater than query test\n");
    		return;
    	}
	    System.out.println("Non-unique greater than query test start");
	    int count = 0;
	    TGResultSet<TGEntity> resultSet = conn.executeQuery("@nodetype='testnode' and citycode = 'SC' and times10 > 940;", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Result set has %d nodes - %s\n", count, count == 60 ? "correct" : "wrong");
	    System.out.println("Non-unique greater than query test end\n");
    }

    public void testNonUniqueRangeDownQuery() throws Exception {
    	if (!testNonUniqueRangeDown && !testAll) {
    		System.out.println("Skipped non-unique less than query test\n");
    		return;
    	}
	    System.out.println("Non-unique less than query test start");
	    int count = 0;
	    TGResultSet<TGEntity> resultSet = conn.executeQuery("@nodetype='testnode' and citycode = 'BG' and times10 < 230;", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Result set has %d nodes - %s\n", count, count == 220 ? "correct" : "wrong");
	    System.out.println("Non-unique less than query test end\n");
    }

    public void testNonUniqueRangeQuery() throws Exception {
    	if (!testNonUniqueRange && !testAll) {
    		System.out.println("Skipped non-unique range query test\n");
    		return;
    	}
    	System.out.println("Non-unique range query test start");
	    int count = 0;
	    TGResultSet<TGEntity> resultSet;
		resultSet = conn.executeQuery("@nodetype='testnode' and citycode = 'SC' and times10 > 940 and times10 < 1000;", TGQueryOption.DEFAULT_QUERY_OPTION);
		if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Result set has %d nodes - %s\n", count, count == 50 ? "correct" : "wrong");
    	System.out.println("Non-unique range query test end\n");
    }

    public void testFullScanQuery() throws Exception {
    	if (!testFullScan && !testAll) {
    		System.out.println("Skipped full scan test\n");
    		return;
    	}
    	System.out.println("Full scan test start");
    	System.out.println("Query actress Michelle Pfeiffer");
	    int count = 0;
	    TGResultSet<TGEntity> resultSet = conn.executeQuery("@nodetype='testnode' and actor = 'Michelle Pfeiffer';", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Query 'Michelle Pfeiffer' returns : %d nodes - %s\n", count, count == 10 ? "correct" : "wrong");
    	System.out.println("Query actor with birth year 1958");
	    count = 0;
	    resultSet = conn.executeQuery("@nodetype='testnode' and birthyear = 1958;", TGQueryOption.DEFAULT_QUERY_OPTION);
	    if (resultSet != null) {
	    	while (resultSet.hasNext()) {
	    		TGEntity entity = resultSet.next();
	    		printEntity(entity);
	    		count++;
	    	}
	    }
	    System.out.printf("Query birth year 1958 returns : %d nodes - %s\n", count, count == 60 ? "correct" : "wrong");
    	System.out.println("Full scan test end\n");
    }

    private void printEntity(TGEntity entity) {
    	if (entity == null) {
    		System.out.println("Entity is null");
    		return;
    	}
    	Collection<TGAttribute> attrs = entity.getAttributes();
    	int i = 0;
    	for (TGAttribute attr : attrs) {
    		System.out.printf("%s%s : %s", i > 0 ? ", " : "", attr.getAttributeDescriptor().getName(), attr.getValue());
    		i++;
    	}
    	System.out.println("");
    }

    public static void main(String[] args) throws Exception {
    	IndexTest idxt = new IndexTest();
    	idxt.getArgs(args);
    	idxt.startup();
    	idxt.initDB();
    	idxt.testGetByKey();
    	idxt.testUniqueQuery();
    	idxt.testPartialUniqueQuery();
    	idxt.testNonUniqueQuery();
    	idxt.testPartialNonUniqueQuery();
    	idxt.testNonUniqueRangeUpQuery();
    	idxt.testNonUniqueRangeDownQuery();
    	idxt.testNonUniqueRangeQuery();
    	idxt.testFullScanQuery();
    	idxt.shutdown();
    }
}

