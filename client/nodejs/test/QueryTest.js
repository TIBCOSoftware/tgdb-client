/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not
 * use this file except in compliance with the License. A copy of the License is
 * included in the distribution package with this file. You also may obtain a
 * copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

var conFactory    = require('../lib/connection/TGConnectionFactory'),
	TGQueryOption = require('../lib/query/TGQueryOption').TGQueryOption,
	PrintUtility  = require('../lib/utils/PrintUtility'),
	TGLogManager  = require('../lib/log/TGLogManager'),
	TGLogLevel    = require('../lib/log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();
	
var testNodeType = null;

// name is primary key
function testPKey(conn, callback) {
	logger.logInfo( "[QueryTest::testPKey] Start unique query");
    // use primary key
    	conn.executeQuery("@nodetype = 'testnode' and name = 'Lex Luthor';", 
    		              TGQueryOption.DEFAULT_QUERY_OPTION,
    		              function(resultSet) {
    	if (!resultSet) {
    		logger.logInfo( "[QueryTest::testPKey] Empty result set !!!!!!");
    	} else {
    		logger.logInfo( "[QueryTest::testPKey] Result set size : " + resultSet.count());
    	    var i=1;
    	    while (resultSet.hasNext()) {
    	        var node = resultSet.next();
    	        logger.logInfo( "[QueryTest::testPKey] Node : %d", i++);
    	        PrintUtility.printEntitiesBreadth(node, 5);
    	    }
    	}
    	logger.logInfo( "[QueryTest::testPKey] End unique query");
    	callback();
   	});
}

function testPartialUnique(conn, callback) {
    logger.logInfo( 
    		"[QueryTest::testPartialUnique] Start partial unique query");
    // use unique key testidx1
    	conn.executeQuery("@nodetype = 'testnode' and nickname = 'Bad guy' and level = 11.0;", 
    	                  TGQueryOption.DEFAULT_QUERY_OPTION,
    	                  function(resultSet){
        if (!resultSet) {
        	logger.logInfo( 
        			"[QueryTest::testPartialUnique] Empty result set !!!!!!");
        } else {
        	var i=1;
        	while (resultSet.hasNext()) {
        	    var node = resultSet.next();
        	    logger.logInfo( 
        	    		"[QueryTest::testPartialUnique] Node : %d", i++);
        	    PrintUtility.printEntitiesBreadth(node, 5);
        	}
        }
        logger.logInfo( 
        		"[QueryTest::testPartialUnique] End partial unique query");
       	callback();
    });
}

function testNonunique(conn, callback) {
	logger.logInfo( 
			"[QueryTest::testNonunique] Start non-unique query");
    // use nonunique key testidx2
	conn.executeQuery("@nodetype = 'testnode' and nickname = 'Superhero';", 
                      TGQueryOption.DEFAULT_QUERY_OPTION,
                      function(resultSet){
        if (!resultSet) {
            logger.logInfo( 
            		"[QueryTest::testNonunique] Empty result set !!!!!!");
        } else {
            var i=1;
            while (resultSet.hasNext()) {
                var node = resultSet.next();
                logger.logInfo( 
                		"[QueryTest::testNonunique] Node : %d", i++);
                PrintUtility.printEntitiesBreadth(node, 5);
            }
        }
        logger.logInfo( 
        		"[QueryTest::testNonunique] End non-unique query");
        callback();
    });
}

function testGreaterThan(conn, callback) {
	logger.logInfo( 
			"[QueryTest::testGreaterThan] Start GT query");
    // use nonunique key testidx2
	conn.executeQuery("@nodetype = 'testnode' and age > 32;", 
                      TGQueryOption.DEFAULT_QUERY_OPTION,
                      function(resultSet){
        if (!resultSet) {
            logger.logInfo( 
            		"[QueryTest::testGreaterThan] Empty result set !!!!!!");
        } else {
            var i=1;
            while (resultSet.hasNext()) {
                var node = resultSet.next();
                logger.logInfo( 
                		"[QueryTest::testGreaterThan] Node : %d", i++);
                PrintUtility.printEntitiesBreadth(node, 5);
            }
        }
        logger.logInfo( 
        		"[QueryTest::testGreaterThan] End GT query");
        callback();
    });
}

function testLessThan(conn, callback) {
	logger.logInfo( 
			"[QueryTest::testLessThan] Start LT query");
    // use nonunique key testidx2
	conn.executeQuery("@nodetype = 'testnode' and age < 28;", 
                      TGQueryOption.DEFAULT_QUERY_OPTION,
                      function(resultSet){
        if (!resultSet) {
            logger.logInfo( 
            		"[QueryTest::testLessThan] Empty result set !!!!!!");
        } else {
            var i=1;
            while (resultSet.hasNext()) {
                var node = resultSet.next();
                logger.logInfo( 
                		"[QueryTest::testLessThan] Node : %d", i++);
                PrintUtility.printEntitiesBreadth(node, 5);
            }
        }
        logger.logInfo( 
        		"[QueryTest::testLessThan] End LT query");
        callback();
    });
}

function testRange(conn, callback) {
	logger.logInfo(
			"[QueryTest::testRange] Start range query");
    // use nonunique key testidx2
	conn.executeQuery("@nodetype = 'testnode' and age > 28 and age < 33 and level > 2.9 and level < 4.1;", 
                      TGQueryOption.DEFAULT_QUERY_OPTION,
                      function(resultSet){
        if (!resultSet) {
            logger.logInfo( 
            		"[QueryTest::testRange] Empty result set !!!!!!");
        } else {
            var i=1;
            while (resultSet.hasNext()) {
                var node = resultSet.next();
                logger.logInfo( 
                		"[QueryTest::testRange] Node : %d", i++);
                PrintUtility.printEntitiesBreadth(node, 5);
            }
        }
        logger.logInfo( 
        		"[QueryTest::testRange] End range query");
        callback();
    });
}

function createNode(gof, nodeType) {
	if (!nodeType) {
		return gof.createNode();
	} else {
		return gof.createNode(nodeType);
	}
}

function fetchMetadata(conn, callback) {
	logger.logInfo( 
			"[QueryTest::fetchMetadata] Start fetch metadata");
	conn.getGraphMetadata(true, function(graphMetadata) {
	    testNodeType = graphMetadata.getNodeType("testnode");
	    if (!testNodeType) {
 		    logger.logInfo( 
 		    	"[QueryTest::fetchMetadata] 'testnode' is not found from meta data fetch");
	    } else {
   		    logger.logInfo( 
   		    	"[QueryTest::fetchMetadata] 'testnode' is found with %d attributes", 
   		    	testNodeType.getAttributeDescriptors().length);
	    }
		logger.logInfo( 
				"[QueryTest::fetchMetadata] End fetch metadata");
		callback();
	});
}

function createData(conn, callback) {

	logger.logInfo( 
			'[QueryTest::createData] Im in createData .........');

    var gof = conn.getGraphObjectFactory();

	logger.logInfo( "Create test nodes");
    var node1 = createNode(gof, testNodeType);
    node1.setAttribute("name", "Bruce Wayne");
    node1.setAttribute("multiple", 7);
    node1.setAttribute("nickname", "Superhero");
    node1.setAttribute("level", 4.0);
    node1.setAttribute("rate", 5.5);
    node1.setAttribute("age", 40);
    conn.insertEntity(node1);
    var node11 = createNode(gof, testNodeType);
    node11.setAttribute("name", "Peter Parker");
    node11.setAttribute("multiple", 7);
    node11.setAttribute("nickname", "Superhero");
    node11.setAttribute("level", 4.0);
    node11.setAttribute("rate", 3.3);
    node11.setAttribute("age", 24);
    conn.insertEntity(node11);
    var node12 = createNode(gof, testNodeType);
    node12.setAttribute("name", "Clark Kent");
    node12.setAttribute("multiple", 7);
    node12.setAttribute("nickname", "Superhero");
    node12.setAttribute("level", 4.0);
    node12.setAttribute("rate", 6.6);
    node12.setAttribute("age", 32);
    conn.insertEntity(node12);
    var node13 = createNode(gof, testNodeType);
    node13.setAttribute("name", "James Logan Howlett");
    node13.setAttribute("multiple", 7);
    node13.setAttribute("nickname", "Superhero");
    node13.setAttribute("level", 4.0);
    node13.setAttribute("rate", 4.4);
    node13.setAttribute("age", 40);
    conn.insertEntity(node13);
    var node14 = createNode(gof, testNodeType);
    node14.setAttribute("name", "Diana Prince");
    node14.setAttribute("multiple", 7);
    node14.setAttribute("nickname", "Superheroine");
    node14.setAttribute("level", 4.0);
    node14.setAttribute("rate", 5.9);
    node14.setAttribute("age", 35);
    conn.insertEntity(node14);
    var node15 = createNode(gof, testNodeType);
    node15.setAttribute("name", "Jean Grey");
    node15.setAttribute("multiple", 7);
    node15.setAttribute("nickname", "Superheroine");
    node15.setAttribute("level", 4.0);
    node15.setAttribute("rate", 6.2);
    node15.setAttribute("age", 35);
    conn.insertEntity(node15);

    var node2 = createNode(gof, testNodeType);
    node2.setAttribute("name", "Mary Jane Watson");
    node2.setAttribute("multiple", 14);
    node2.setAttribute("nickname", "Girlfriend");
    node2.setAttribute("level", 3.0);
    node2.setAttribute("rate", 6.3);
    node2.setAttribute("age", 22);
    conn.insertEntity(node2);
    var node21 = createNode(gof, testNodeType);
    node21.setAttribute("name", "Lois Lane");
    node21.setAttribute("multiple", 14);
    node21.setAttribute("nickname", "Girlfriend");
    node21.setAttribute("level", 3.0);
    node21.setAttribute("rate", 6.4);
    node21.setAttribute("age", 30);
    conn.insertEntity(node21);
    var node22 = createNode(gof, testNodeType);
    node22.setAttribute("name", "Raven DarkhÃ¶lme");
    node22.setAttribute("multiple", 14);
    node22.setAttribute("nickname", "Criminal");
    node22.setAttribute("level", 3.0);
    node22.setAttribute("rate", 7.2);
    node22.setAttribute("age", 36);
    conn.insertEntity(node22);
    var node23 = createNode(gof, testNodeType);
    node23.setAttribute("name", "Selina Kyle");
    node23.setAttribute("multiple", 14);
    node23.setAttribute("nickname", "Criminal");
    node23.setAttribute("level", 3.0);
    node23.setAttribute("rate", 6.5);
    node23.setAttribute("age", 30);
    conn.insertEntity(node23);
    var node24 = createNode(gof, testNodeType);
    node24.setAttribute("name", "Harley Quinn");
    node24.setAttribute("multiple", 14);
    node24.setAttribute("nickname", "Criminal");
    node24.setAttribute("level", 3.0);
    node24.setAttribute("rate", 5.8);
    node24.setAttribute("age", 30);
    conn.insertEntity(node24);

    var node3 = createNode(gof, testNodeType);
    node3.setAttribute("name", "Lex Luthor");
    node3.setAttribute("multiple", 52);
    node3.setAttribute("nickname", "Bad guy");
    node3.setAttribute("level", 11.0);
    node3.setAttribute("rate", 7.1);
    node3.setAttribute("age", 26);
    conn.insertEntity(node3);
    var node31 = createNode(gof, testNodeType);
    node31.setAttribute("name", "Harvey Dent");
    node31.setAttribute("multiple", 52);
    node31.setAttribute("nickname", "Bad guy");
    node31.setAttribute("level", 11.0);
    node31.setAttribute("rate", 4.6);
    node31.setAttribute("age", 40);
    conn.insertEntity(node31);
    var node32 = createNode(gof, testNodeType);
    node32.setAttribute("name", "Victor Creed");
    node32.setAttribute("multiple", 52);
    node32.setAttribute("nickname", "Bad guy");
    node32.setAttribute("level", 11.0);
    node32.setAttribute("rate", 6.4);
    node32.setAttribute("age", 40);
    conn.insertEntity(node32);
    var node33 = createNode(gof, testNodeType);
    node33.setAttribute("name", "Norman Osborn");
    node33.setAttribute("multiple", 52);
    node33.setAttribute("nickname", "Bad guy");
    node33.setAttribute("level", 11.0);
    node33.setAttribute("rate", 6.3);
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

    var edge5 = gof.createBidirectionalEdge(node13, node15);
    edge5.setAttribute("name", "Friend");
    conn.insertEntity(edge5);
    var edge6 = gof.createBidirectionalEdge(node13, node22);
    edge6.setAttribute("name", "Enemy");
    conn.insertEntity(edge6);
    var edge7 = gof.createBidirectionalEdge(node13, node32);
    edge7.setAttribute("name", "Enemy");
    conn.insertEntity(edge7);

	conn.commit ( 
		function(status) {
		logger.logInfo( 
				"[QueryTest::createData] Commit transaction completed");
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
    	logger.logError( "Exception happens : " + exception.message); 
    	throw exception;
    }).connect(function(connectionStatus) {
		if (connectionStatus) {
			fetchMetadata(conn, function() {
				createData(conn, function() {
				    testPKey(conn, function(){
				    	testPartialUnique(conn, function() {
				    		testNonunique(conn, function() {
				    			testGreaterThan(conn, function() {
				    				testLessThan(conn, function() {
				    					testRange(conn, function() {
					    					
									    	conn.disconnect();
									    	logger.logInfo(
									    		"[QueryTest::test] Disconnected.");

								        });  //testRange
							        });  //testLessThan
						        });  //testGreaterThan
					        });  //testNonunique
				        });  //testPartialUnique
				    });	 //testPKey			
				});  //createData				
			});  //fetchMetadata
		}  //if
	});  //connect
}

test();