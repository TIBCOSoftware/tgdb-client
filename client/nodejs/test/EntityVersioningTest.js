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
    TGNumber      = require('../lib/datatype/TGNumber'),
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
		testNodeType = gmd.getNodeType("testnode");
		if (testNodeType !== null) {
			logger.logInfo("'testnode' is found with %d attributes\n", 
					testNodeType.getAttributeDescriptors().length);
		} else {
			logger.logInfo("'testnode' is not found from meta data fetch");
		}
		callback();
	});
}

var nodeMap = {};
function createCommit(conn, callback) {
    var gof = conn.getGraphObjectFactory();
    
    if (gof === null) {
    	console.log("Graph Object Factory is null...exiting");
    }

    logger.logInfo("Start transaction 1");
    logger.logInfo("Create node1");
    var node1 = createNode(gof, testNodeType);
    node1.setAttribute("name", "john doe4");
    node1.setAttribute("multiple", 7);
    conn.insertEntity(node1);
    nodeMap[node1.getAttribute("name").getValue()] = node1;

  	console.log("Commit transaction 1");
    conn.commit(function(){
    	logger.logInfo("Commit transaction 1 completed");

    	logger.logInfo("Start transaction 2");
        // updates
    	logger.logInfo("Update node1");
        var node = nodeMap["john doe4"];
        node.setAttribute("multiple", 8);
        conn.updateEntity(node);
         
        logger.logInfo("Commit transaction 2");
        conn.commit(function(){
        	logger.logInfo("Commit transaction 2 completed");
        	
        	logger.logInfo("Start transaction 3");
            // updates
        	logger.logInfo("Update node1");
            var node = nodeMap["john doe4"];
            node.setAttribute("multiple", 9);
            conn.updateEntity(node);
             
            logger.logInfo("Commit transaction 3");
            conn.commit(function(){
            	logger.logInfo("Commit transaction 3 completed");
            	callback();
            });
        }); // update John doe's value --------.
    }); //----- write data to database ----------. Everything is create
}

function test() {
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.DebugWire);
	
    var connectionFactory = conFactory.getFactory();
    var linkURL = 'tcp://scott@192.168.1.6:8222';
    var conn = connectionFactory.createConnection(linkURL, 'scott', 'scott', null);
    
    conn.connect(function(connectionStatus) {
        if (connectionStatus) {
        	logger.logInfo('Connection to server successful');
            fetchMetadata(conn, function(){
            	createCommit(conn, function(){            			
	  	  	            conn.disconnect();
	  	  	            logger.logInfo("Connection test connection disconnected.");
            	}); 
            });
        }
    });	
}

test();