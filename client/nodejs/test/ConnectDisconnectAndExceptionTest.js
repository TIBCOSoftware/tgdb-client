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

function connectThenDisconnect01(conn, callback) {
	logger.logInfo('\n*** Start connectThenDisconnect01 !!! ***\n');

    conn.connect(function(connectionStatus){
        if (connectionStatus) {
        	logger.logInfo( 'Connection to server successful 01');
            conn.disconnect();
            callback();
        }
    });
}

function connectThenDisconnect02(conn, callback) {
	logger.logInfo('\n*** Start connectThenDisconnect02 !!! ***\n');

    conn.on('connect', function(connectionStatus){
        if (connectionStatus) {
        	logger.logInfo( 'Connection to server successful 02');
            conn.disconnect();
            callback();
        }
    })
    .connect();
}

function exceptionThentryToDisconnect(callback) {
	logger.logInfo('\n*** Start exceptionThentryToDisconnect !!! ***\n');

    var linkURL = 'tcp://scott@255.255.255.255:8222';
    var conn = conFactory.getFactory().createConnection(linkURL, 'scott', 'scott', null);
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
	
    var linkURL = 'tcp://scott@192.168.1.6:8222';
    var conn = conFactory.getFactory().createConnection(linkURL, 'scott', 'scott', null);
	connectThenDisconnect01(conn, function(){
		connectThenDisconnect02(conn, function(){
			exceptionThentryToDisconnect(function(msg){
				logger.logInfo( '\n*** Callback from %s ***\n', msg);
			}); 
		});
	});
}

test();