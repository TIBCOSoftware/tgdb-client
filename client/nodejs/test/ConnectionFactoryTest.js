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

var conFactory = require('../lib/connection/TGConnectionFactory'),
    TGLogManager  = require('../lib/log/TGLogManager'),
    TGLogLevel    = require('../lib/log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

var counter = 0;
function connectionTest(conn, callback) {
	logger.logInfo('\n*** Start test %d !!! ***\n', ++counter);

    logger.logInfo('conn, my id is = ' + conn.getId());
    conn.on('connect', function(connectionStatus){
        if (connectionStatus) {
        	logger.logInfo( 'Connection to server successful.');
            conn.disconnect();
        }
        callback('From onConnect clause .. ');
    })
    .on('exception', function(exception){
        if (exception) {
        	logger.logError( 'Exception happens : %s', exception.message);
            conn.disconnect();
        }
        callback('From onException clause ..');
    })
    .connect();
}

function test() {
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.Debug);
	
	var factory01 = conFactory.getFactory();
	logger.logInfo('factory01, my id is = ' + factory01.getId());
	
	var factory02 = conFactory.getFactory();
	logger.logInfo('factory02, my id is = ' + factory02.getId());

	var linkURL = 'tcp://scott@192.168.1.6:8222';
	var properties = { 
//		ConnectionImpl : './tcp/TCPConnection',
//		ConnectionPoolSize : 5
	};
    var conn01 = factory01.createConnection(linkURL, 'scott', 'scott', properties);
  
    logger.logInfo('conn01, my id is = ' + conn01.getId());
    
    connectionTest(conn01, function(msg){
		logger.logInfo( '\n*** Callback from %s 1 ***\n', msg);
	    var linkURL = 'tcp://scott@255.255.255.255:8222';
	    properties.ConnectionPoolSize = 1;
	    var conn02 = conFactory.getFactory().createConnection(linkURL, 'scott', 'scott', properties);
	    connectionTest(conn02, function(msg){
			logger.logInfo( '\n*** Callback from %s 2 ***\n', msg);
		});
	});
}

test();