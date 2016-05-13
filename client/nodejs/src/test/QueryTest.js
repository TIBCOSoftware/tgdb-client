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

var conFactory              = require('../connection/TGConnectionFactory'),
	ProtocolDataInputStream = require('../pdu/impl/ProtocolDataInputStream').ProtocolDataInputStream, 
	StringUtils             = require('../utils/StringUtils').StringUtils, 
    TGEdge                  = require('../model/TGEdge'), 
	TGAttributeType         = require('../model/TGAttributeType').TGAttributeType;

function test() {
	var connectionFactory = new conFactory.DefaultConnectionFactory();
	var linkURL = 'tcp://scott@192.168.1.18:8222';
	var conn = connectionFactory.createConnection(linkURL, 'scott', 'scott',
			null);

	var callback = function(connectionStatus) {
		if (connectionStatus) {
			console.log('Connection to server successful');
			// executeQuery(conn);
			createThenExecuteQuery(conn);
		}
	};

	conn.connect(callback);
}

function createThenExecuteQuery(conn) {

	console.log('Im in createThenExecuteQuery .........');

	console.log('create query');

	conn.createQuery("testquery < X '5ef';", function(query) {

		console.log('execute query');

		query.execute(function(resultSet) {

			console.log('Got result set from TGQuery.execute().');
			console.log('Process result set.');
			console.log('close query');

			query.close(function() {
				console.log('query closed');
			});
		});
	});
}

function executeQuery(conn) {
	console.log('Im in executeQuery .........');

	conn.executeQuery("testquery < X '5ef';", function(resultSet) {
		console.log('Got result set from TGConnection.executeQuery().');
	});

	console.log('Im after executeQuery .........');

}

test();