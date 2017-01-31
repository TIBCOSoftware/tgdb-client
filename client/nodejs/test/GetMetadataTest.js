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

var conFactory   = require('../lib/connection/TGConnectionFactory'),
    TGException  = require('../lib/exception/TGException').TGException,
    TGLogManager = require('../lib/log/TGLogManager'),
    TGLogLevel   = require('../lib/log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function fetchMetadata(conn) {
    
    var gof = conn.getGraphObjectFactory();

    conn.getGraphMetadata(true, function(gm) {
    	
    	//The TestAttrString1 and TestNodeType should be loaded as test data on the server.
    	//They are there only for testing.  
    	
    	var descriptors = gm.getAttributeDescriptors();
    	descriptors.forEach(function(attributeDescriptor) {
        	if (!(!attributeDescriptor)) {
        		logger.logInfo( 	
        			"'%s' has id %d and data type : %d", attributeDescriptor.getName(), 
        			attributeDescriptor.getAttributeId(), attributeDescriptor.getType().value);
        	} else {
        		logger.logInfo( 
        			"'%s' is not found from meta data fetch", attributeDescriptor);
        	}
    	});

    	var nodeType = gm.getNodeType("testnode");
    	if (!nodeType) {
    		logger.logInfo( 
    				"'testnode' is not found from meta data fetch");
    	} else {
    		logger.logInfo( 
    				"'testnode' is found with %d attributes\n", nodeType.getAttributeDescriptors().length);
    	}
    	
        conn.disconnect();
        logger.logInfo( 
        		"Connection test connection disconnected.");

    });
}

function test() {
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.DebugWire);

	var connectionFactory = conFactory.getFactory();
    var linkURL = 'tcp://scott@192.168.1.6:8222';
    var conn = connectionFactory.createConnection(linkURL, 'scott', 'scott', null);
    
    conn.on('exceptionx', function(exception){
    	logger.logError( "Exception happens : " + exception.message);    	
    }).on('connect', function(connectionStatus) {
        if (connectionStatus) {
        	logger.logInfo( 'Connection to server successful');
            fetchMetadata(conn);
        }
    }).connect();
}

test();