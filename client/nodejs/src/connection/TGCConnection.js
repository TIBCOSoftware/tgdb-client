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

var PrintUtility              = require('../utils/PrintUtility').PrintUtility,
	util                      = require('util'),
    VerbId                    = require('../pdu/impl/VerbId').VerbId,
    TGResultSet               = require('./TGResultSet').TGResultSet,
    TGQuery                   = require('./TGQuery').TGQuery,
	TGEntityManager           = require('../model/TGEntityManager').TGEntityManager,
    CommitTransactionRequest  = require('../pdu/impl/CommitTransactionRequest').CommitTransactionRequest,
    ProtocolMessageFactory    = require('../pdu/impl/ProtocolMessageFactory').ProtocolMessageFactory;


var globleConnIds = 0;

// Class definition
function TGCConnection(connectionPool, channel, properties) {
	TGCConnection.super_.call(this);
	this._connId             = globleConnIds++;
    this._connectionPool     = connectionPool;
    this._channel            = channel;
    this._properties         = properties;
    this._callback           = null;
}

util.inherits(TGCConnection, TGEntityManager);

/**
 * Connect to Graph DB server asynchronously. Upon successful handshake followed
 * by authentication, the callback function will be called to indicate
 * successful establishment of the connection.
 * <p>
 * Any operations relying on success of connection should be called only upon
 * responsestatus param in the callback function's value = true.
 * </p>
 * 
 * @param callback -
 *            Function. - signature function mycallback(responsestatus)
 */
TGCConnection.prototype.connect = function(callback) {
    if (callback == undefined || (typeof callback !== 'function')) {
        throw new Error('Callback should be a function');
    }
    var channel = this._channel;
    channel.connect(callback);
};

/**
 * Disconnect from the graph DB server.
 */
TGCConnection.prototype.disconnect = function() {
    var channel = this._channel;
    channel.disconnect();
};

TGCConnection.prototype.deleteEntity = function (tgEntity) {
    this.entityDeleted(tgEntity);
};

TGCConnection.prototype.updateEntity = function (tgEntity) {
    this.entityUpdated(tgEntity);
};

/**
 * Commit transaction to Graph DB server.
 * 
 * @param callback -
 *            Convey status of commit operation.
 */
TGCConnection.prototype.commit = function(callback) {
	if(null!==this._callback) {
		throw new Error('[TGCConnection.commit] Channel is busy for previous request!');
	}
	this._callback = callback;
	
	var addedEntities   = this.addedEntities();
    PrintUtility.printEntityMap(addedEntities, 'addedEntities');
    
	var updatedEntities = this.updatedEntities();
    PrintUtility.printEntityMap(updatedEntities, 'updatedEntities');
    
	var removedEntities = this.removedEntities();
    PrintUtility.printEntityMap(removedEntities, 'removedEntities');
	
    var attrDescSet = this.newAttrDecsSet(); 
    var commitTransRequest = new CommitTransactionRequest(addedEntities, updatedEntities, removedEntities, attrDescSet);
    commitTransRequest.setRequestId(this._channel.getRequestId());
    commitTransRequest.setAuthToken(this._channel.getAuthToken());
    commitTransRequest.setSessionId(this._channel.getSessionId());

    this._channel.send(commitTransRequest, this);
   
};

TGCConnection.prototype.handleCommitResponse = function(response) {
	// Fix all temp Ids
	this.updateEntityIds(response);

    // ToDo need to handle commit fail
	if(null!==this._callback) {
		var callback = this._callback;
		this._callback = null;
		callback(true);
	}
}

/**
 * Rollback transaction to Graph DB server.
 */
TGCConnection.prototype.rollback = function() {
	this.clear();
};

/****************************
 *                          *
 *        Query API         *
 *                          *
 ****************************/

TGCConnection.prototype.createQuery = function (expr, callback) {
	if(null!==this._callback) {
		throw new Error('[TGCConnection.commit] Channel is busy for previous request!');
	}
	this._callback = callback;

    //var timeout = Long.parseLong(properties.getProperty(CONFIG_NAMES.CONNECTION_OPERATION_TIMEOUT, "-1"));
    //var requestId  = globleRequestIds++;
	
	var request = ProtocolMessageFactory.createMessageFromVerbId(VerbId.QUERY_REQUEST);
	request.setCommand(TGCConnection.command.CREATE);
//	request.setConnectionId(connId);
	request.setQuery(expr);
	this._channel.send(request, this);
	console.log("Send create query completed");
};

TGCConnection.prototype.executeQuery = function (queryExpr, callback) {
	console.log("Entering TGCConnection.prototype.executeQuery .......... ");
	if(null!==this._callback) {
		throw new Error('[TGCConnection.commit] Channel is busy for previous request!');
	}
	this._callback = callback;
	
    //var timeout = Long.parseLong(properties.getProperty(CONFIG_NAMES.CONNECTION_OPERATION_TIMEOUT, "-1"));
	//var requestId  = QueryRequest.getThenIncrementQueryId();

	var request = ProtocolMessageFactory.createMessageFromVerbId(VerbId.QUERY_REQUEST);
	request.setCommand(TGCConnection.command.EXECUTE);
//	request.setConnectionId(this._connId);
	request.setQuery(queryExpr);
	this._channel.send(request, this);
	console.log("Send execute query completed");
};

TGCConnection.prototype.closeQuery = function (query, callback) {
	if(null!==this._callback) {
		throw new Error('[TGCConnection.commit] Channel is busy for previous request!');
	}
	this._callback = callback;
	
    //var timeout = Long.parseLong(properties.getProperty(CONFIG_NAMES.CONNECTION_OPERATION_TIMEOUT, "-1"));
    //var requestId  = requestIds.getAndIncrement();

	var request = ProtocolMessageFactory.createMessageFromVerbId(VerbId.QUERY_REQUEST);
	request.setCommand(TGCConnection.command.CLOSE);
//	request.setConnectionId(connId);
	request.setQuery(query);
	this._channel.send(request, this);
	console.log("Send close query completed");
};

TGCConnection.prototype.handleQueryResponse = function (request, response) {
	console.log('Entering TGCConnection.prototype.handleQueryResponse : ' + request.getCommand());
	var callback = this._callback;
	this._callback = null;
	switch (request.getCommand()) {
		case 1 :
		    var queryExpr = response.getQuery();
		    var queryObj = null;
			if(queryExpr != null) {
				queryObj = new TGQuery(this, queryExpr);
				console.log('Query back from server : ' + queryExpr);
			}
			callback(queryObj);
			break;
		case 2 :
			var resultSet = new TGResultSet();
			callback(resultSet);
			break;
		case 3 :
			callback();
			break;
		default :
			throw new Error('Unknow resopnse for query.');
	}
}

TGCConnection.command = {
	CREATE : 1,
	EXECUTE : 2,
	CLOSE : 3
};

exports.TGCConnection = TGCConnection;
