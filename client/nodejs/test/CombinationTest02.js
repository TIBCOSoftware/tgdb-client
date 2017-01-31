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

function createCommit(conn, callback) {
    var gof = conn.getGraphObjectFactory();
    
    if (gof === null) {
    	logger.logInfo( 
    			"Graph Object Factory is null...exiting");
    }

    logger.logInfo( "Start transaction 1");
    logger.logInfo( "Create node1");
  	//createNode/Edge by itself does not include the it in the transaction.
    //Explicit call to insertEntity to add it to the transaction.
    var node1 = createNode(gof, testNodeType);
    node1.setAttribute("name", "john doe");
    node1.setAttribute("multiple", 7);
    node1.setAttribute("rate", 3.3);
    node1.setAttribute("nickname", "蠢豬");
    conn.insertEntity(node1);

    logger.logInfo( "Create node2");
    var node2 = createNode(gof, testNodeType);
    node2.setAttribute("name", "julie");
    node2.setAttribute("factor", 3.3);
    conn.insertEntity(node2);

    logger.logInfo( "Create node3");
    var node3 = createNode(gof, testNodeType);
    node3.setAttribute("name", "margo");
    node3.setAttribute("factor", 2.3);
    node3.setAttribute("是超人嗎", false);
    conn.insertEntity(node3);

    logger.logInfo( "Create edge1");
    var edge1 = gof.createBidirectionalEdge(node1, node2);
    edge1.setAttribute("name", "spouse");
    conn.insertEntity(edge1);

    logger.logInfo( "Create edge2");
    var edge2 = gof.createDirectedEdge(node1, node3);
    edge2.setAttribute("name", "daughter");
    conn.insertEntity(edge2);

    logger.logInfo( "Commit transaction 1");
    conn.commit(function(){
    	logger.logInfo( "Commit transaction 1 completed");

    	logger.logInfo( "Start transaction 2");
        // updates
    	logger.logInfo( "Update node1");
        node1.setAttribute("age", 40);
        node1.setAttribute("multiple", 8);
        conn.updateEntity(node1);
        
        //logger.logInfo( "Delete edge1");
        //conn.deleteEntity(edge1);
        
        logger.logInfo( "Update edge2");
        //add a new attribute
        edge2.setAttribute("extra", true);
        //update existing one
        edge2.setAttribute("name", "kid");
        conn.updateEntity(edge2);
        
        logger.logInfo( "Create node4");
        var node4 = createNode(gof, testNodeType);
        node4.setAttribute("name", "Smith");
        node4.setAttribute("level", 3.0);
        conn.insertEntity(node4);

        logger.logInfo( "Create edge3");
        var edge3 = gof.createBidirectionalEdge(node1, node4);
        edge3.setAttribute("name", "Tennis Partner");
        conn.insertEntity(edge3);

        logger.logInfo( "Commit transaction 2");
        conn.commit(function(){
          	console.log("Commit transaction 2 completed");
        	callback();
        }); // update John doe's value --------.
    }); //----- write data to database ----------. Everything is create
}

var depth = 5;
var printDepth = 5;
var resultCount = 100;
var edgeLimit = 0;
function getEntity(conn, callback) {
    var gof = conn.getGraphObjectFactory();
    logger.logInfo( "Start unique get operation");
  	var key = gof.createCompositeKey("testnode");
  	key.setAttribute("name", "foo");
    var option = TGQueryOption.createQueryOption();
    option.setPrefetchSize(resultCount);
    option.setTraversalDepth(depth);
    option.setEdgeLimit(edgeLimit);

  	conn.getEntity(key, option, function(ent){
  	  	if (ent !== null) {
  	  	    PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	} else {
  	  	logger.logInfo( "getEntity for 'foo' returns nothing");
  	  	}

  	  	key.setAttribute("name", "margo");
  	  	conn.getEntity(key, null, function(ent){
  	  	  	if (ent !== null) {
  	  	  	    PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	} else {
  	  	  	logger.logInfo( "getEntity for 'margo' returns nothing");
  	  	  	}

  	  	  	key = gof.createCompositeKey("testnode");
  	  	  	key.setAttribute("nickname", "margo");
  	  	  	conn.getEntity(key, null, function(ent){
  	  	  	  	if (ent !== null) {
  	  	  	        PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	} else {
  	  	  	        logger.logInfo( "getEntity for nickname 'margo' returns nothing");
  	  	  	  	}
  	  	  	  		
  	  	  	  	key.setAttribute("nickname", "margo");
  	  	  	  	key.setAttribute("rate",  12.34);
  	  	  	  	conn.getEntity(key, null, function(ent){
  	  	  	  	  	if (ent !== null) {
  	  	  	  	        PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  	} else {
  	  	  	  	        logger.logInfo("getEntity for nickname 'margo' and rate - 12.34 returns nothing");
  	  	  	  	  	}
  	  	  	  	  		
  	  	  	  	  	key = gof.createCompositeKey("testnode");
  	  	  	  	  	key.setAttribute("multiple", 3);
  	  	  	  	  	conn.getEntity(key, null, function(ent){
  	  	  	  	  	  	if (ent !== null) {
  	  	  	  	  	        PrintUtility.printEntities(ent, printDepth, 0, "", {});
  	  	  	  	  	  	} else {
  	  	  	  	  	        logger.logInfo("getEntity for multiple - 3");
  	  	  	  	  	  	}
  	  	  	  	        logger.logInfo("End unique get operation"); 
  	  	  	  	        callback();
  	  	  	  	  	});
  	  	  	  	});
  	  	  	});
  	  	});
  	});
}

function test() {
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.DebugWire);

    var connectionFactory = conFactory.getFactory();
    var linkURL = 'tcp://scott@192.168.1.6:8222';
    var conn = connectionFactory.createConnection(linkURL, 'scott', 'scott', null);
    
    conn.on('exceptionx', function(exception){
    	logger.logError( "Exception Happens, message : " + exception.message);
    	conn.disconnect();
    }).connect(function(connectionStatus) {
        if (connectionStatus) {
        	logger.logInfo('Connection to server successful');
            fetchMetadata(conn, function(){
            	createCommit(conn, function(){
            		getEntity(conn, function(){
            			
	  	  	            conn.disconnect();
	  	  	            logger.logInfo("Connection test connection disconnected.");

            		});
            	}); 
            });
        }
    });	
}

test();