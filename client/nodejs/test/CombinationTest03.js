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
 */

var conFactory    = require('../lib/connection/TGConnectionFactory'),
	TGQueryOption = require('../lib/query/TGQueryOption').TGQueryOption,
	PrintUtility  = require('../lib/utils/PrintUtility'),
    TGLogManager  = require('../lib/log/TGLogManager'),
    TGLogLevel    = require('../lib/log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function createNode(gof, nodeType) {
	if (!nodeType) {
		return gof.createNode();
	} else {
		return gof.createNode(nodeType);
	}
}

var testNodeType = null;
function fetchMetadata(conn, callback) {
	    conn.getGraphMetadata(true, function(gmd) {
		var attrDesc = gmd.getAttributeDescriptor("factor");
		if (attrDesc !== null) {
			logger.logInfo( 
					"'factor' has id %d and data type : %d", 
					attrDesc.getAttributeId(), attrDesc.getType().value);
		} else {
			logger.logInfo( 
					"'factor' is not found from meta data fetch");
		}
		attrDesc = gmd.getAttributeDescriptor("level");
		if (attrDesc !== null) {
			logger.logInfo( 
					"'level' has id %d and data type : %d", 
					attrDesc.getAttributeId(), attrDesc.getType().value);
		} else {
			logger.logInfo( 
					"'level' is not found from meta data fetch");
		}
		testNodeType = gmd.getNodeType("testnode");
		if (testNodeType !== null) {
			logger.logInfo( 
					"'testnode' is found with %d attributes", 
					testNodeType.getAttributeDescriptors().length);
		} else {
			logger.logInfo( 
					"'testnode' is not found from meta data fetch");
		}
		callback();
	});
}

var node13 = null;
var node32 = null;
function createCommit(conn, callback) {
    var gof = conn.getGraphObjectFactory();
	logger.logInfo(
			"Test Transaction : Insert Simple Node(John) of testnode type with a few properties");
	logger.logInfo(
			"Create node");
    var node1 = createNode(gof, testNodeType);
    node1.setAttribute("name", "Bruce Wayne");
    node1.setAttribute("multiple", 7);
    node1.setAttribute("rate", 5.5);
    node1.setAttribute("nickname", "è¶…ç´šè‹±é›„");
    node1.setAttribute("level", 4.0);
    node1.setAttribute("age", 38);
    conn.insertEntity(node1);
    var node11 = createNode(gof, testNodeType);
    node11.setAttribute("name", "Peter Parker");
    node11.setAttribute("multiple", 7);
    node11.setAttribute("rate", 3.3);
    node11.setAttribute("nickname", "è¶…ç´šè‹±é›„");
    node11.setAttribute("level", 4.0);
    node11.setAttribute("age", 24);
    conn.insertEntity(node11);
    var node12 = createNode(gof, testNodeType);
    node12.setAttribute("name", "Clark Kent");
    node12.setAttribute("multiple", 7);
    node12.setAttribute("rate", 6.6);
    node12.setAttribute("nickname", "è¶…ç´šè‹±é›„");
    node12.setAttribute("level", 4.0);
    node12.setAttribute("age", 32);
    conn.insertEntity(node12);
    node13 = createNode(gof, testNodeType);
    node13.setAttribute("name", "James Logan Howlett");
    node13.setAttribute("multiple", 7);
    node13.setAttribute("rate", 4.4);
    node13.setAttribute("nickname", "è¶…ç´šè‹±é›„");
    node13.setAttribute("level", 4.0);
    node13.setAttribute("age", 40);
    conn.insertEntity(node13);

    var node2 = createNode(gof, testNodeType);
    node2.setAttribute("name", "Mary Jane Watson");
    node2.setAttribute("multiple", 14);
    node2.setAttribute("rate", 6.3);
    node2.setAttribute("nickname", "ç¾Žéº—");
    node2.setAttribute("level", 3.0);
    node2.setAttribute("age", 22);
    conn.insertEntity(node2);
    var node21 = createNode(gof, testNodeType);
    node21.setAttribute("name", "Lois Lane");
    node21.setAttribute("multiple", 14);
    node21.setAttribute("rate", 6.4);
    node21.setAttribute("nickname", "ç¾Žéº—");
    node21.setAttribute("level", 3.0);
    node21.setAttribute("age", 30);
    conn.insertEntity(node21);
    var node22 = createNode(gof, testNodeType);
    node22.setAttribute("name", "Jean Grey");
    node22.setAttribute("multiple", 14);
    node22.setAttribute("rate", 6.2);
    node22.setAttribute("nickname", "ç¾Žéº—");
    node22.setAttribute("level", 3.0);
    node22.setAttribute("age", 30);
    conn.insertEntity(node22);
    var node23 = createNode(gof, testNodeType);
    node23.setAttribute("name", "Selina Kyle");
    node23.setAttribute("multiple", 14);
    node23.setAttribute("rate", 6.5);
    node23.setAttribute("nickname", "Criminal");
    node23.setAttribute("level", 3.0);
    node23.setAttribute("age", 30);
    conn.insertEntity(node23);
    var node24 = createNode(gof, testNodeType);
    node24.setAttribute("name", "Harley Quinn");
    node24.setAttribute("multiple", 14);
    node24.setAttribute("rate", 5.8);
    node24.setAttribute("nickname", "Criminal");
    node24.setAttribute("level", 3.0);
    node24.setAttribute("age", 30);
    conn.insertEntity(node24);

    var node3 = createNode(gof, testNodeType);
    node3.setAttribute("name", "Lex Luthor");
    node3.setAttribute("multiple", 52);
    node3.setAttribute("rate", 7.1);
    node3.setAttribute("nickname", "å£žäºº");
    node3.setAttribute("level", 11.0);
    node3.setAttribute("age", 26);
    conn.insertEntity(node3);
    var node31 = createNode(gof, testNodeType);
    node31.setAttribute("name", "Harvey Dent");
    node31.setAttribute("multiple", 52);
    node31.setAttribute("rate", 4.6);
    node31.setAttribute("nickname", "å£žäºº");
    node31.setAttribute("level", 11.0);
    node31.setAttribute("age", 40);
    conn.insertEntity(node31);
    node32 = createNode(gof, testNodeType);
    node32.setAttribute("name", "Victor Creed");
    node32.setAttribute("multiple", 52);
    node32.setAttribute("rate", 6.4);
    node32.setAttribute("nickname", "å£žäºº");
    node32.setAttribute("level", 11.0);
    node32.setAttribute("age", 40);
    conn.insertEntity(node32);
    var node33 = createNode(gof, testNodeType);
    node33.setAttribute("name", "Norman Osborn");
    node33.setAttribute("multiple", 52);
    node33.setAttribute("rate", 6.3);
    node33.setAttribute("nickname", "å£žäºº");
    node33.setAttribute("level", 11.0);
    node33.setAttribute("age", 50);
    conn.insertEntity(node33);

    var edge1 = gof.createBidirectionalEdge(node1, node31);
    edge1.setAttribute("name", "Nemesis");
    conn.insertEntity(edge1);
    var edge2 = gof.createBidirectionalEdge(node1, node23);
    edge2.setAttribute("name", "Frenemy");
    conn.insertEntity(edge2);
    var edge3 = gof.createBidirectionalEdge(node1, node24);
    edge3.setAttribute("name", "Nemesis");
    conn.insertEntity(edge3);
    var edge4 = gof.createBidirectionalEdge(node1, node12);
    edge4.setAttribute("name", "Teammate");
    conn.insertEntity(edge4);

    logger.logInfo( "Commit transaction 1");
    conn.commit(function(){    	
        callback();
    });
}

var depth = 5;
var printDepth = 5;
var resultCount = 100;
var edgeLimit = 0;
function getEntity(conn, callback) {
    var gof = conn.getGraphObjectFactory();
    logger.logInfo("Start unique get operation");

  	var key = gof.createCompositeKey("testnode");
  	key.setAttribute("name", "John Doe");
    var option = TGQueryOption.createQueryOption();
    option.setPrefetchSize(resultCount);
    option.setTraversalDepth(depth);
    option.setEdgeLimit(edgeLimit);

  	conn.getEntity(key, option, function(ent){
  		if (!!ent) {
  			logger.logInfo("getEntity for 'John Doe' - found - wrong");
	  	  	PrintUtility.printEntities(ent, printDepth, 0, "", {});
  		} else {
  			logger.logInfo("getEntity for 'John Doe' returns nothing - good");
  		}

  		key.setAttribute("name", "Bruce Wayne");
  	  	conn.getEntity(key, null, function(ent){
      		if (!!ent) {
      			logger.logInfo("getEntity for 'Bruce Wayne' - found - good");
  	  	  	    PrintUtility.printEntities(ent, printDepth, 0, "", {});
      		} else {
      			logger.logInfo("getEntity for 'Bruce Wayne' returns nothing - wrong");
      		}
      		key.setAttribute("name", "Peter Parker");
  	  	  	conn.getEntity(key, null, function(ent){
  	      		if (!!ent) {
  	      			logger.logInfo("getEntity for 'Peter Parker' - found - good");
	  	  	  	    PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	      		} else {
  	      			logger.logInfo("getEntity for 'Peter Parker' returns nothing - wrong");
  	      		}
  	      		key.setAttribute("name", "Mary Jane Watson");
  	  	  	  	conn.getEntity(key, null, function(ent){
  	        		if (!!ent) {
  	        			logger.logInfo("getEntity for 'Mary Jane Watson' - found - good");
	  	  	  	  	    PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	        		} else {
  	        			logger.logInfo("getEntity for 'Mary Jane' returns nothing - wrong");
  	        		}
  	        		key.setAttribute("name", "Super Jane");
  	  	  	  	  	conn.getEntity(key, null, function(ent){
  	  	  	  	  		if (!!ent) {
  	  	  	  	  			logger.logInfo("getEntity for 'Super Jane' - found - wrong");
  	  	  	  	  	        PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  		} else {
  	  	  	  	  			logger.logInfo("getEntity for 'Super Jane' returns nothing - good");
  	  	  	  	  		}

  	  	  	  	  		var key1 = gof.createCompositeKey("testnode");
  	  	  	  	  		key1.setAttribute("nickname", "Stupid");
  	  	  	  	  		conn.getEntity(key1, null, function(ent){
  	  	  	  	  			if (!!ent) {
  	  	  	  	  				logger.logInfo("getEntity for 'Stupid' - found - wrong");
  	  	  	  	  	  	        PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  			} else {
  	  	  	  	  				logger.logInfo("getEntity for nickname 'Stupid' returns nothing - good");
  	  	  	  	  			}
  	  	  	  	  			key1.setAttribute("nickname", "ç¾Žéº—");
  	  	  	  	  			conn.getEntity(key1, null, function(ent){
  	  	  	  	  				if (!!ent) {
  	  	  	  	  					logger.logInfo("getEntity for 'ç¾Žéº—' - found - good");
  	    	  	  	  	  	        PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  				} else {
  	  	  	  	  					logger.logInfo("getEntity for nickname 'ç¾Žéº—' returns nothing - wrong");
  	  	  	  	  				}
  	  	  	  	  				key1.setAttribute("nickname", "Criminal");
  	  	  	  	  				conn.getEntity(key1, null, function(ent){
  	  	  	  	  					if (!!ent) {
  	  	  	  	  						logger.logInfo("getEntity for 'Criminal' - found - good");
  	  	  	  	  	  	  	            PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  					} else {
  	  	  	  	  						logger.logInfo("getEntity for nickname 'Criminal' returns nothing - wrong");
  	  	  	  	  					}
  	  	  	  	  					key1.setAttribute("nickname", "è¶…ç´šè‹±é›„");
  	  	  	  	  					conn.getEntity(key1, null, function(ent){
  	  	  	  	  						if (!!ent) {
  	  	  	  	  							logger.logInfo("getEntity for 'è¶…ç´šè‹±é›„' - found - good");
  	  	  	  	  	  	  	  	            PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  						} else {
  	  	  	  	  							logger.logInfo("getEntity for nickname 'è¶…ç´šè‹±é›„' returns nothing - wrong");
  	  	  	  	  						}

  	  	  	  	  						var key2 = gof.createCompositeKey("testnode");
  	  	  	  	  						key2.setAttribute("nickname", "å£žäºº");
  	  	  	  	  						key2.setAttribute("level", 11.0);
  	  	  	  	  						key2.setAttribute("rate", 7.1);
  	  	  	  	  						conn.getEntity(key2, null, function(ent){
  	  	  	  	  							if (!!ent) {
  	  	  	  	  								logger.logInfo("getEntity for 'å£žäºº'(11.0)(7.1) found - good");
  	  	  	    	  	  	  	  	            PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  							} else {
  	  	  	  	  								logger.logInfo("getEntity for 'å£žäºº'(11.0)(7.1) returns nothing - wrong");
  	  	  	  	  							}

  	  	  	  	  							var key3 = gof.createCompositeKey("testnode");
  	  	  	  	  							key3.setAttribute("multiple", 52);
  	  	  	  	  							conn.getEntity(key3, null, function(ent){
  	  	  	  	  								if (!!ent) {
  	  	  	  	  									logger.logInfo("getEntity for multiple 52 found - good");
  	  	  	  	  	  	  	  	  	                PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  								} else {
  	  	  	  	  									logger.logInfo("getEntity for multiple 52 returns nothing - wrong");
  	  	  	  	  								}
  	  	  	  	  								logger.logInfo("End unique get operation");

  	  	  	  	  								logger.logInfo("Start update test");
  	  	  	  	  								node13.setAttribute("age", 43);
  	  	  	  	  								//node13.setAttribute("nickname", "Wolverine");
  	  	  	  	  								node13.setAttribute("nickname", "è¶…ç´š");
  	  	  	  	  								conn.updateEntity(node13);

  	  	  	  	  								//node32.setAttribute("nickname", "Sabertooth");
  	  	  	  	  								node32.setAttribute("nickname", "å®³æ€•");
  	  	  	  	  								node32.setAttribute("level", 11.0);
  	  	  	  	  								node32.setAttribute("age", 38);
  	  	  	  	  								conn.updateEntity(node32);
  	  	  	  	  								conn.commit(function(){
  	  	  	  	  									logger.logInfo("End update test");

  	  	  	  	  									logger.logInfo("Get update entity");
  	  	  	  	  									key2.setAttribute("nickname", "å®³æ€•");
  	  	  	  	  									key2.setAttribute("level", 11.0);
  	  	  	  	  									key2.setAttribute("rate", 6.4);
  	  	  	  	  									conn.getEntity(key2, null, function(ent){
  	  	  	  	  										if (!!ent) {
  	  	  	  	  											logger.logInfo("getEntity for 'Sabertooth'(11.0)(6.4) found - good");
  	  	  	  	  		  	  	  	  	  	                PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  										} else {
  	  	  	  	  											logger.logInfo("getEntity for 'Sabertooth'(11.0)(6.4) returns nothing - wrong");
  	  	  	  	  										}
  	  	  	  	  										logger.logInfo("Get update entity end");
  	  	  	  	  										callback();
  	  	  	  	  									});
  	  	  	  	  								});
  	  	  	  	  							});
  	  	  	  	  						});
  	  	  	  	  					});
  	  	  	  	  				});
  	  	  	  	  			});
  	  	  	  	  		});
   	  	  	  	  	});
  	  	  	  	});
  	  	  	});
  	  	});
  	});
}

function query(conn, callback) {
	logger.logInfo("Start test query");
    //TGQuery Query1 = conn.createQuery("testquery < X '5ef';");
    //Query1.execute();
    //Query1.close();
    conn.executeQuery("@nodetype = 'testnode' and ((nickname = 'å£žäºº' and level = 11.0) or (level = 3.0));",
    				  TGQueryOption.DEFAULT_QUERY_OPTION,
    				  function(resultSet){
	    var i=1;
	    while (resultSet.hasNext()) {
	        var node = resultSet.next();
	        logger.logInfo( "[CombinationTest03::testPKey] Node : %d", i++);
	        PrintUtility.printEntitiesBreadth(node, 5);
	    }
    	callback();
    });
}

function test() {
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.DebugWire);

    var connectionFactory = conFactory.getFactory();
    var linkURL = 'tcp://scott@192.168.1.6:8222';
    var conn = connectionFactory.createConnection(linkURL, 'scott', 'scott', null);
    
    conn.on('exception', function(exception){
    	logger.logError( "Exception Happens, message : " + exception.message);
    	conn.disconnect();
    	//throw exception;
    }).connect(function(connectionStatus) {
        if (connectionStatus) {
        	logger.logInfo('Connection to server successful');
            fetchMetadata(conn, function(){
            	createCommit(conn, function(){
            		getEntity(conn, function(){
            			query(conn, function(){
    	  	  	            conn.disconnect();
    	  	  	            logger.logInfo("Connection test connection disconnected.");            				
            			});
            		});
            	}); 
            });
        }
    });	
}

test();