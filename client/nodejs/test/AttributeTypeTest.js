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
	PrintUtility  = require('../lib/utils/PrintUtility').PrintUtility,
    TGNumber      = require('../lib/datatype/TGNumber'),
    TGLogManager  = require('../lib/log/TGLogManager'),
    TGLogLevel    = require('../lib/log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

var datasForInsert = [
	{ name:'testBOOLEAN'  , value:true },          //A single bit representing the truth value
	{ name:'testBYTE'     , value:21   },          //8bit octet
	{ name:'testCHAR'     , value:22   },          //Fixed 8-Bit octet of N length
	{ name:'testSHORT'    , value:23   },          //16bit
	{ name:'testINT'      , value:24   },          //32bit signed integer
	{ name:'testLONG'     , value:25   },          //64bit float
	{ name:'testFLOAT'    , value:26   },          //32bit float
	{ name:'testDOUBLE'   , value:27   },          //64bit float
//	{ name:'testNUMBER'   , value:'28' },          //Number with precision
	{ name:'testSTRING'   , value:'29' },          //Varying length String < 64K
	{ name:'testDATE'     , value:new Date()},     //Only the Date part of the DateTime
	{ name:'testTIME'     , value:new Date()},     //Only the Time part of the DateTime
	{ name:'testTIMESTAMP', value:new Date()},     //64bit Timestamp - engine time - upto nanosecond precision, if the OS provides //SS:TODO
];

var datasForUpdate = [
    { name:'testBOOLEAN'  , value:false },         //A single bit representing the truth value
    { name:'testBYTE'     , value:37   },          //8bit octet
    { name:'testCHAR'     , value:36   },          //Fixed 8-Bit octet of N length
    { name:'testSHORT'    , value:35   },          //16bit
    { name:'testINT'      , value:34   },          //32bit signed integer
    { name:'testLONG'     , value:33   },          //64bit float
    { name:'testFLOAT'    , value:32   },          //32bit float
    { name:'testDOUBLE'   , value:31   },          //64bit float
//    { name:'testNUMBER'   , value:'30' },        //Number with precision
    { name:'testSTRING'   , value:'29' },          //Varying length String < 64K
    { name:'testDATE'     , value:new Date()},     //Only the Date part of the DateTime
    { name:'testTIME'     , value:new Date()},     //Only the Time part of the DateTime
    { name:'testTIMESTAMP', value:new Date()},     //64bit Timestamp - engine time - upto nanosecond precision, if the OS provides //SS:TODO
];

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
		testNodeType = gmd.getNodeType("typetestnode");
		if (testNodeType !== null) {
			logger.logInfo("'typetestnode' is found with %d attributes", 
					testNodeType.getAttributeDescriptors().length);
		} else {
			logger.logInfo("'typetestnode' is not found from meta data fetch");
		}
		callback();
	});
}

function createThenUpdate(conn, callback) {
    var gof = conn.getGraphObjectFactory();
    
    if (gof === null) {
    	console.log("Graph Object Factory is null...exiting");
    }

    logger.logInfo("Start transaction 1");
    logger.logInfo("Create node");
    var node = createNode(gof, testNodeType);
    datasForInsert.forEach(function(data){
        node.setAttribute(data.name, data.value);
    });

    conn.insertEntity(node);

  	console.log("Commit transaction 1");
    conn.commit(function(){
    	logger.logInfo("Commit transaction 1 completed");

      	var key = gof.createCompositeKey("typetestnode");
      	key.setAttribute('testSTRING', '29');
      	conn.getEntity(key, null, function(ent01){
      		if (!!ent01) {
      			logger.logInfo("getEntity for '29' - found - good");
    	  	  	PrintUtility.printEntities(ent01, 5, 0, "", {});
      		} else {
      			logger.logInfo("getEntity for '29' returns nothing - wrong");
      		}
      		
        	logger.logInfo("Start transaction 2");
        	logger.logInfo("Update node");
      	    datasForUpdate.forEach(function(data){
      	    	ent01.setAttribute(data.name, data.value);
      	    });
            conn.updateEntity(ent01);
            logger.logInfo("Commit transaction 2");
            conn.commit(function(){
            	logger.logInfo("Commit transaction 2 completed");
              	conn.getEntity(key, null, function(ent02){
              		if (!!ent02) {
              			logger.logInfo("getEntity for '29' - found - good");
            	  	  	PrintUtility.printEntities(ent02, 5, 0, "", {});
              		} else {
              			logger.logInfo("getEntity for '29' returns nothing - wrong");
              		}
                	callback();
              	});
            }); // update John doe's value --------.

      	});
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
            	createThenUpdate(conn, function(){
	  	  	        conn.disconnect();
	  	  	        logger.logInfo("Connection test connection disconnected.");
            	}); 
            });
        }
    });	
}

test();