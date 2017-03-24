

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import java.util.*;
import java.text.SimpleDateFormat;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGEntity;

import org.apache.commons.cli.*;
/*
    importCsv -file <> -command <>
    ex of command : create node (type {a:[colname1], b:[<colname2>]})
    create node houseMemberType {memberName:%memberName%, crownName:%crownName%, houseHead:%houseHead%, yearBorn:%yearBorn%, yearDied:%yearDied%, reignStart:%reignStart%, reignEnd:%reignEnd%, crownTitle:%crownTitle%}
	javac -cp ..\lib\commons-cli-1.3.1.jar;..\lib\tgdb-client.jar;. ImportCsv.java

	C:\tgdb\1.0\examples>
	java -cp ..\lib\commons-cli-1.3.1.jar;..\lib\tgdb-client.jar;. ImportCsv -file ./hierarchy/csv-nodes.csv -command "create node houseMemberType {memberName:%memberName%, crownName:%crownName%, houseHead:%houseHead%, yearBorn:%yearBorn%, yearDied:%yearDied%, reignStart:%reignStart%, reignEnd:%reignEnd%, crownTitle:%crownTitle%}"
    java -cp ..\lib\commons-cli-1.3.1.jar;..\lib\tgdb-client.jar;. ImportCsv -file ./hierarchy/marriages.csv -command "create node marriageType {marriageId:%couple%}"
    java -cp ..\lib\commons-cli-1.3.1.jar;..\lib\tgdb-client.jar;. ImportCsv -file ./hierarchy/marriages.csv -command "create edge (houseMemberType {memberName:%From%})-[{relType:\"spouse\"}]-(marriageType {marriageId:%couple%})"
    java -cp ..\lib\commons-cli-1.3.1.jar;..\lib\tgdb-client.jar;. ImportCsv -file ./hierarchy/marriages.csv -command "create edge (houseMemberType {memberName:%To%})-[{relType:\"spouse\"}]-(marriageType {marriageId:%couple%})"
    java -cp ..\lib\commons-cli-1.3.1.jar;..\lib\tgdb-client.jar;. ImportCsv -file ./hierarchy/parents.csv -command "create edge (houseMemberType {memberName:%id%})-[{relType:\"child\"}]-(marriageType {marriageId:%couple%})"





    java -cp ..\lib\commons-cli-1.3.1.jar;..\lib\tgdb-client.jar;. ImportCsv -file ./hierarchy/parents.csv


    create edge (nodetypefrom {attribute:value,...})-[{attribute:value, ...}]-(nodetypeto {attribute:value, ....})

    java -cp ..\lib\commons-cli-1.3.1.jar;..\lib\tgdb-client.jar;. ImportCsv -file ./hierarchy/csv-edges.csv -command "create edge (houseMemberType {memberName:%From%})-[{relType:%relation%}]-(houseMemberType {memberName:%To%})"
*/



 abstract class GraphCmd {
    static final String url = "tcp://127.0.0.1:8222";
	static final String user = "admin";
	static final String pwd = "admin";
	TGGraphObjectFactory gof;
	TGGraphMetadata gmd;
	TGConnection conn = null;
	abstract void run();
	abstract Map<String,String> getFieldsMap();
	abstract void setFields(Map<String,String> map);
	abstract void parseCommand(String cmd);
            public void close() {
            	System.out.println("Closing connection");
            	if (conn!=null)
            		conn.disconnect();
            }
	public void connect() throws Exception {
		conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
					conn.connect();

					gof = conn.getGraphObjectFactory();
					if (gof == null) {
						throw new Exception("Graph object not found");
					}

					gmd = conn.getGraphMetadata(true);
					
	}
}
 class CreateNodeCmd extends GraphCmd {

    String _cmd=null;
    String nodetype=null;

    TGNodeType graphNodeType;

          
            Map<String,String> fieldsDef;
            Map<String,String> replaced;
       
            CreateNodeCmd(String cmd){
            	fieldsDef = new LinkedHashMap<String,String>();
            	_cmd=cmd;
                parseCommand(cmd);
            }

            public Map<String,String> getFieldsMap() {
            	return fieldsDef;
            }
            public void setFields(Map<String,String> map) {
            	// replace script
                for (String key : fieldsDef.keySet()) {
                	String keyreplace=fieldsDef.get(key);
                	if (keyreplace.startsWith("%")) {
					   // System.out.println("Key : " + key +" = "+keyreplace+ " Value : " + map.get(keyreplace));
						replaced.put(key,map.get(keyreplace));
					}
					
				}
            }
            public void run() {
                    
					SimpleDateFormat df = new SimpleDateFormat("MM/dd/yyyy");
					SimpleDateFormat df2 = new SimpleDateFormat("dd MMM yyyy");
					Calendar cal=null;
					TGNode graphNode;
					System.out.println("Inserting node  "+nodetype+" "+ replaced.toString());
					//
					// Insert node data into database
					//
					try {
						graphNode = gof.createNode(graphNodeType);
						for (String key : replaced.keySet()) {
							String val = replaced.get(key);
							if (!val.equals("null")) {
								//System.out.println(key+" = ["+val+"]");
								//TGAttributeDescriptor attr= graphNodeType.getAttributeDescriptor(key);
								
									//System.out.println(key +" attr type "+attr.getType().toString());
									if (val.equals("true") || val.equals("false")) {
										//System.out.println(key +" attr type boolean");
										graphNode.setAttribute(key, Boolean.valueOf(val));
									} else {
										try { 
	        								int i = Integer.parseInt(val);
	        								//System.out.println(key +" attr type Integer");
	        								graphNode.setAttribute(key, i);
	        							} catch (NumberFormatException e) { 
	        							   try {
		        							   	Date d=df.parse(val);
		        							   	//System.out.println(key +" attr type Date");
		        							   	cal=Calendar.getInstance();
	  											cal.setTime(d);
		        							   	graphNode.setAttribute(key, cal);
	        							   } catch (Exception e2) {
		        							   	try {
			        							   	Date d=df2.parse(val);
			        							   	//System.out.println(key +" attr type Date");
			        							   	cal=Calendar.getInstance();
		  											cal.setTime(d);
			        							   	graphNode.setAttribute(key, cal);
		        							   } catch (Exception e3) {
		        							      //System.out.println(key +" attr type String "+e3.toString());
											      graphNode.setAttribute(key, replaced.get(key));
											   }
	        							      
										   }
									    }
									}

							}
						}
						
							conn.insertEntity(graphNode);
				            conn.commit(); // Write data to database
				            System.out.println("Node Inserted "+nodetype+" "+ replaced.toString());
				        } catch (Exception ex) {
				        	System.out.println("Insert exception "+ex.toString());
				        }
            }
            public  void parseCommand(String cmd)  {
		    	/* (?:[,]?(\w+):(?:\W+)) */
		      String pattern = "create node (\\w+)([\\w,\\W]*)";
		      String fieldpattern = 	"(?:[,]?\\W*(\\w+):([%,\"]\\w+[%,\"]))";
			  Matcher m2;
		      // Create a Pattern object
		      Pattern rp = Pattern.compile(pattern);
			  Pattern rf = Pattern.compile(fieldpattern);
		      // Now create matcher object.
		      Matcher m = rp.matcher(cmd);
		      if (m.find()){
		      	nodetype = m.group(1);
		      	  m2 = rf.matcher(m.group(2));
		      	  System.out.println("node type: " + nodetype );
		      	  System.out.println("field to parse: " + m.group(2) );
			      while (m2.find( )) {
			         System.out.println("Found value: " + m2.group(0) );
			         System.out.println("Found value: " + m2.group(1) );
			         System.out.println("Found value: " + m2.group(2) );
			         fieldsDef.put(m2.group(1).trim(),m2.group(2).trim());
			      }
		  	  }
		      replaced = new LinkedHashMap<String,String>(fieldsDef);
		      // prepare the graph 
		      
			  try {
					connect();
					graphNodeType = gmd.getNodeType(nodetype);
					if (graphNodeType == null)
						throw new Exception("Node type not found");
				    
				} 
				catch (Exception e) {
					   System.out.println("graph exception "+e.toString());
				}
				
            }
            // Valid in JDK 8 and later:

//            public void printOriginalNumbers() {
//                System.out.println("Original numbers are " + phoneNumber1 +
//                    " and " + phoneNumber2);
//            }
        }
class CreateEdgeCmd extends GraphCmd {

    String _cmd=null;
    String nodetypefrom=null;
    String nodetypeto=null;
    
 
    TGNodeType graphNodeTypeFrom;
    TGNodeType graphNodeTypeTo;

          
            Map<String,String> fieldsDefFrom, fieldsDefTo, fieldsDefEdge;
            Map<String,String> replacedFrom,replacedTo,replacedEdge;
       
            CreateEdgeCmd(String cmd){
            	fieldsDefFrom = new LinkedHashMap<String,String>();
            	fieldsDefTo = new LinkedHashMap<String,String>();
            	fieldsDefEdge = new LinkedHashMap<String,String>();
            	_cmd=cmd;
                parseCommand(cmd);
            }

            public Map<String,String> getFieldsMap() {
            	return fieldsDefEdge;
            }
            public void setFields(Map<String,String> map) {
            	// replace script
                for (String key : fieldsDefFrom.keySet()) {
                	String keyreplace=fieldsDefFrom.get(key);
                	if (keyreplace.startsWith("%")) {
					   // System.out.println("Key : " + key +" = "+keyreplace+ " Value : " + map.get(keyreplace));
						replacedFrom.put(key,map.get(keyreplace));
					}
					
				}
				for (String key : fieldsDefTo.keySet()) {
                	String keyreplace=fieldsDefTo.get(key);
                	if (keyreplace.startsWith("%")) {
					   // System.out.println("Key : " + key +" = "+keyreplace+ " Value : " + map.get(keyreplace));
						replacedTo.put(key,map.get(keyreplace));
					}
					
				}
				for (String key : fieldsDefEdge.keySet()) {
                	String keyreplace=fieldsDefEdge.get(key);
                	if (keyreplace.startsWith("%")) {
					   // System.out.println("Key : " + key +" = "+keyreplace+ " Value : " + map.get(keyreplace));
						replacedEdge.put(key,map.get(keyreplace));
					}
					
				}
            }
            public void run() {
            	/* perform edge creation
            	we need to search of from and to nodes first and create edge 
            	*/
            	System.out.print("creating edge from "+replacedFrom.toString());
							System.out.print(" to "+replacedTo.toString());
							System.out.println(" attributes "+replacedEdge.toString());
            	try {
            		
            		System.out.println("create key for ["+nodetypefrom+"]");
		            	TGKey keyfrom = gof.createCompositeKey(nodetypefrom);
		            	System.out.println("create key for ["+nodetypeto+"]");
						TGKey keyto = gof.createCompositeKey(nodetypeto);
		                for (String key : replacedFrom.keySet()) {
		                	String keyreplace=replacedFrom.get(key);
		                	System.out.println("from set key "+key+" "+keyreplace);
		                	keyfrom.setAttribute(key, keyreplace);
						}
					    for (String key : replacedTo.keySet()) {
		                	String keyreplace=replacedTo.get(key);
		                	System.out.println("to set key "+key+" "+keyreplace);
		                	keyto.setAttribute(key, keyreplace);
						}
					TGNode from = (TGNode)conn.getEntity(keyfrom, null);
					TGNode to = (TGNode)conn.getEntity(keyto, null);
					TGEdge edge = gof.createEdge(from, to, TGEdge.DirectionType.Directed);
					for (String key : replacedEdge.keySet()) {
							String val = replacedEdge.get(key);
							if (!val.equals("null")) {
								System.out.println("edge set attr "+key+" "+val);
								edge.setAttribute(key,val);
							}
						}
						conn.insertEntity(edge);
						conn.commit();

            	} catch (Exception ex) {
				        	System.out.println("create edge exception ");
				        	ex.printStackTrace();
				        }
            }
            public  void parseCommand(String cmd)  {
		    	/* (?:[,]?(\w+):(?:\W+)) */
		    	/* parsing line 
		    	create edge (nodetypefrom {attribute:value,...})-[{attribute:value, ...}]-(nodetypeto {attribute:value, ....})
		    	*/
		      String pattern = "create edge \\((\\w+) ([^\\)]*)\\)-\\[([^\\]]*)\\]-\\((\\w+) ([^\\)]*)\\)";
		      String fieldpattern = 	"(?:[,]?\\W*(\\w+):([%,\"](\\w+)[%,\"]))";
			  Matcher m2;
		      // Create a Pattern object
		      Pattern rp = Pattern.compile(pattern);
			  Pattern rf = Pattern.compile(fieldpattern);
		      // Now create matcher object.
		      Matcher m = rp.matcher(cmd);
		      if (m.find()){
		      	  nodetypefrom = m.group(1);
		      	  m2 = rf.matcher(m.group(2));
		      	  System.out.println("node from type: " + nodetypefrom );
		      	  System.out.println("field to parse: " + m.group(2) );
			      while (m2.find( )) {
			         fieldsDefFrom.put(m2.group(1).trim(),m2.group(2).trim());
			      }
			      nodetypeto = m.group(4);
			      m2 = rf.matcher(m.group(5));
		      	  System.out.println("node to type: " + nodetypeto );
		      	  System.out.println("field to parse: " + m.group(5) );
			      while (m2.find( )) {
			         fieldsDefTo.put(m2.group(1).trim(),m2.group(2).trim());
			      }
			      System.out.println("attributes for edge to parse: " + m.group(3) );
			      m2 = rf.matcher(m.group(3));
			      while (m2.find( )) {
			      	// use %xxx% or yyy in "yyy"
			      	 if (m2.group(2).trim().startsWith("%")) {
			         	fieldsDefEdge.put(m2.group(1).trim(),m2.group(2).trim());
			     		} else {
						fieldsDefEdge.put(m2.group(1).trim(),m2.group(3).trim());
			     	}
			      }
		  	  }
		      replacedFrom = new LinkedHashMap<String,String>(fieldsDefFrom);
		      replacedTo = new LinkedHashMap<String,String>(fieldsDefTo);
		      replacedEdge = new LinkedHashMap<String,String>(fieldsDefEdge);
		      // prepare the graph 
		      
			  try {
					connect();

					graphNodeTypeFrom = gmd.getNodeType(nodetypefrom);
					graphNodeTypeTo = gmd.getNodeType(nodetypeto);
					if ((graphNodeTypeFrom == null) || (graphNodeTypeTo == null))
						throw new Exception("Node type not found");
				    
				} 
				catch (Exception e) {
					   System.out.println("graph exception "+e.toString());
				}
				
            }

        }

public class ImportCsv 
{
	//Delimiters used in the CSV file
    private static final String COMMA_DELIMITER = ",";

    public static void main(String args[])
    {
        
        CommandLine commandLine;
 
       
    	Option csvfile   = OptionBuilder.withArgName( "file" )
                                .hasArg()
                                .withDescription(  "csv file to import" )
                                .create( "file" );
        Option command   = OptionBuilder.withArgName( "cmd" )
                                .hasArg()
                                .withDescription(  "command to run on each line" )
                                .create( "command" );
        Options options = new Options();
        CommandLineParser parser = new GnuParser();
        options.addOption(csvfile);
        options.addOption(command);
        
        try
        {
            commandLine = parser.parse(options, args);

            if (commandLine.hasOption("file"))
            {
                System.out.print("Option file is present.  The value is: ");
                System.out.println(commandLine.getOptionValue("file"));
                
            } else {
            	throw (new ParseException("missing file"));
            }
            if (commandLine.hasOption("command"))
            { 
            	String cmd=commandLine.getOptionValue("command");
                System.out.print("Command file is present.  The value is: ");
                System.out.println(commandLine.getOptionValue("command"));
                if (cmd.startsWith("create node")) {
                  GraphCmd script = new CreateNodeCmd(cmd);
                  runImport(commandLine.getOptionValue("file"), script);
                  script.close();
				}
				if (cmd.startsWith("create edge")) {
                  GraphCmd script = new CreateEdgeCmd(cmd);
                  runImport(commandLine.getOptionValue("file"), script);
                  script.close();
				}
                
            } else {
            	throw (new ParseException("missing ccommand"));
            }
        }
        catch (ParseException exception)
        {
            System.out.print("Parse error: ");
            System.out.println(exception.getMessage());
        }
    }

    private static void runImport(String filename, GraphCmd cmd) {

        BufferedReader br = null;
        Map<String,String> fields = cmd.getFieldsMap();

        
      try {

            //Reading the csv file
            br = new BufferedReader(new FileReader(filename));
            
            
            
            String line = "";
            //Read to skip the header
            String header = br.readLine();
            String[] colNames = header.split(COMMA_DELIMITER);
            for (String col : colNames) {
            	System.out.println(col);
            }
            //Reading from the second line
            while ((line = br.readLine()) != null) 
            {
            	
                String[] cols = line.split(COMMA_DELIMITER);
                Map<String,String> map = new LinkedHashMap<String,String>();
                for (int i=0; i<colNames.length; i++ ){
                	map.put("%"+colNames[i].trim()+"%",cols[i].trim());
                }
                //System.out.println("line : "+map.toString());
                // replace script
                cmd.setFields(map);
                cmd.run();

            }
            System.out.println("End of file");         
        }
        catch(Exception ee)
        {   
        	try
            {
                br.close();
            }
            catch(IOException ie)
            {
                System.out.println("Error occured while closing the BufferedReader");
                
            }
            ee.printStackTrace();
        }
        
    }
}