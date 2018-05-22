/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * <p/>
 * File name : ConnectionTest1.${EXT}
 * Created on: 1/13/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: ConnectionTest1.java 748 2016-04-25 17:10:38Z vchung $
 */


package com.tibco.tgdb.test;

import java.io.BufferedReader;
import java.io.FileReader;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

public class ImportData {
	public String url = "tcp://scott@localhost:8222";
    public String passwd = "scott";
    public TGLogger.TGLevel logLevel = TGLogger.TGLevel.Info;
    public int edgeFetchCount = -1;
    public int nodeFetchCount = -1;
    public int nodeCommitCount = 1000;
    public int edgeCommitCount = 1000;
    public String hEscFile = "./hESC_mESC.csv";
    public String hEscNetworkFile = "./hESC_comp_network_025.dat";
	public String mEscNetworkFile = "./mESC_comp_network_025.dat";


    boolean treatDoubleAsString = false;
    
    String getStringValue(Iterator<String> argIter) {
    	while (argIter.hasNext()) {
    		String s = argIter.next();
    		return s;
    	}
    	return null;
    }
    
    String getStringValue(Iterator<String> argIter, String defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
    		return s;
    	}
    }

    int getIntValue(Iterator<String> argIter, int defaultValue) {
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

    void getArgs(String[] args) {
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
    			String ll = getStringValue(argIter, "Info");
    			try {
    				logLevel = TGLogger.TGLevel.valueOf(ll);
    			} catch(IllegalArgumentException e) {
    				System.out.printf("Invalid log level value '%s'...ignored\n", ll);
    			}
    		} else if (s.equalsIgnoreCase("-edgecount") || s.equalsIgnoreCase("-ec")) {
    			edgeFetchCount = getIntValue(argIter, edgeFetchCount);
    		} else if (s.equalsIgnoreCase("-nodecount") || s.equalsIgnoreCase("-nc")) {
    			nodeFetchCount = getIntValue(argIter, nodeFetchCount);
    		} else if (s.equalsIgnoreCase("-nodecommitcount") || s.equalsIgnoreCase("-ncc")) {
    			nodeCommitCount = getIntValue(argIter, nodeCommitCount);
    		} else if (s.equalsIgnoreCase("-edgecommitcount") || s.equalsIgnoreCase("-ecc")) {
    			edgeCommitCount = getIntValue(argIter, edgeCommitCount);
    		} else if (s.equalsIgnoreCase("-treatdoubleasstring") || s.equalsIgnoreCase("-dtos")) {
				treatDoubleAsString = true;
			} else if (s.equalsIgnoreCase("-hESC")) {
				hEscFile = getStringValue(argIter, hEscFile);
			} else if (s.equalsIgnoreCase("-hESCNET")) {
				hEscNetworkFile = getStringValue(argIter, hEscNetworkFile);
			} else if (s.equalsIgnoreCase("-mESCNET")) {
				mEscNetworkFile = getStringValue(argIter, mEscNetworkFile);
    		} else {
    			System.out.printf("Skip argument %s\n", s);
    		}
    	}
    }

    void run() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	System.out.printf(" max node count : %d, max edge count : %d, node commit count : %d, edge commit count : %d\n",
            nodeFetchCount, edgeFetchCount, nodeCommitCount, edgeCommitCount);
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        TGGraphObjectFactory gof = conn.getGraphObjectFactory();
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        	conn.disconnect();
        	return;
        }

        // Nodes
        String line;
    	String file = hEscFile;
    	TGGraphMetadata gmd = conn.getGraphMetadata(true);
        HashMap<String, TGNode> mESCmap = new HashMap<String, TGNode>();
        HashMap<String, TGNode> hESCmap = new HashMap<String, TGNode>();
        
    	try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String fromId, toId;
            int count = 0;
            float correlation;
            TGNodeType hESCNodetype = gmd.getNodeType("hESC");
            TGNodeType mESCNodetype = gmd.getNodeType("mESC");
            
            while ((line = br.readLine()) != null) {
                // process the line.
                String[] arr = line.split(";");
            	if (!arr[0].matches("^\\d+$")) {
            		continue;
            	}
                TGNode hESCNode = gof.createNode(hESCNodetype);
				//System.out.printf("Human - Symbol:%s name:%s\n", arr[1], arr[4]);
				//System.out.printf("Mouse - Symbol:%s name:%s\n", arr[3], arr[4]);
				hESCNode.setAttribute("symbol", arr[1]);
                hESCNode.setAttribute("name", arr[4]);
                hESCmap.put(arr[0], hESCNode);
                
                TGNode mESCNode = gof.createNode(mESCNodetype);
                mESCNode.setAttribute("symbol", arr[3]);
                mESCNode.setAttribute("name", arr[4]);
                mESCmap.put(arr[2], mESCNode);

                conn.insertEntity(hESCNode);
                conn.insertEntity(mESCNode);

                count+=2;
                if (count%nodeCommitCount == 0) {
					conn.commit();
                }
                if (count == nodeFetchCount) {
                	break;
                }
            }
            // Last commit for hESC, mESC nodes
            if (count%nodeCommitCount != 0) {
            	conn.commit();
            }
            System.out.printf("Finished processing %d nodes\n", count);
        }
    	
        // hESC Edges
    	TGNode fromNode;
    	TGNode toNode;
    	file = hEscNetworkFile;
    	try (BufferedReader hESCEdgeReader = new BufferedReader(new FileReader(file))) {
            int count = 0;

            while ((line = hESCEdgeReader.readLine()) != null) {
                // process the line.
                String[] arr = line.split("\t");
                fromNode = hESCmap.get(arr[0]);
                toNode = hESCmap.get(arr[1]);
                double infscore = Double.valueOf(arr[2]);
                if (fromNode != null && toNode != null) {
	                TGEdge edge = gof.createEdge(fromNode, toNode, TGEdge.DirectionType.BiDirectional);
                    if (treatDoubleAsString == true) {
	                    edge.setAttribute("infscore", String.valueOf(infscore));
                    } else {
	                    edge.setAttribute("infscore", infscore);
                    }
	                conn.insertEntity(edge);	
	                count++;
	                if (count%edgeCommitCount == 0) {
						//System.out.println("Waiting to commit edges");
						//System.in.read();
	                	conn.commit();
	                	System.out.printf("Count:%d\n", count);
	                	//System.exit(0);
	                }
                }
            }
            // Last commit for hESC
            if (count%edgeCommitCount != 0) {
            	conn.commit();
            }
            System.out.printf("Finished processing %d hESC edges\n", count);
        }
    	
        // mESC Edges
    	file = mEscNetworkFile;
    	try (BufferedReader mESCEdgeReader = new BufferedReader(new FileReader(file))) {
            int count = 0;

            while ((line = mESCEdgeReader.readLine()) != null) {
                // process the line.
                String[] arr = line.split("\t");
                fromNode = mESCmap.get(arr[0]);
                toNode = mESCmap.get(arr[1]);
                double infscore = Double.valueOf(arr[2]);
                if (fromNode != null && toNode != null) {
                	TGEdge edge = gof.createEdge(fromNode, toNode, TGEdge.DirectionType.BiDirectional);
                    if (treatDoubleAsString == true) {
	                    edge.setAttribute("infscore", String.valueOf(infscore));
                    } else {
	                    edge.setAttribute("infscore", infscore);
                    }
                	conn.insertEntity(edge);
                	count++;
                	if (count%edgeCommitCount == 0) {
                		conn.commit();
                		System.out.printf("Count:%d\n", count);
                	}
                }
            }
            // Last commit for mESC
            if (count%edgeCommitCount != 0) {
            	conn.commit();
            }
            System.out.printf("Finished processing %d mESC edges\n", count);
        }

        conn.disconnect();
        System.out.println("Connection test connection disconnected.");
    }

    public static void main(String[] args) throws Exception {
    	ImportData importData = new ImportData();
        importData.getArgs(args);
    	importData.run();
    }
}
