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

var TCPConnection        = require('./tcp/TCPConnection'),
    TGLogManager         = require('../log/TGLogManager'),
    TGLogLevel           = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

var globleConnPoolIds = 1;

function TGConnectionPool(channels, poolSize, properties) {
	this._connPoolId = globleConnPoolIds++;
	
    this._connections = [];
    this._channels = channels;
    this._properties = properties;

    var connection = null;
    var connections = this._connections;
    var me = this;
    this._channels.forEach(function(channel) {
    	if(!properties.ConnectionImpl) {
        	connection = new TCPConnection(me, channel, properties);    		
    	} else {
    		var ConnectionModule = require(properties.ConnectionImpl);
    		connection = new ConnectionModule(me, channel, properties);
    	}
        
        channel.setLinkEventHandler(connection);
        connections.push(connection);
        //channel.setLinkEventHandler(this);
    	logger.logDebug('[TGConnectionPool] connection initialized, id = %s', connection.getId());
    });
    
    
}

/**
 * Exception listener.
 */
TGConnectionPool.prototype.setExceptionListener = function() {

};

/**
 *   LinkEventHandler interface
 */

TGConnectionPool.prototype.onException = function(ex, duringClose) {
	logger.logInfo('[TGConnectionPool.prototype.onException] Exception happens!!');
};

TGConnectionPool.prototype.onReconnect = function() {
	
};

TGConnectionPool.prototype.getTerminatedText = function () {
	
};

/**
 * All connections in the pool to Graph DB server.
 */
TGConnectionPool.prototype.connect = function() {
    var connections = this._connections;
    connections.map(function(con) {
        con.connect();
    });
};

/**
 * All connections in the pool disconnected from the Graph DB server.
 */
TGConnectionPool.prototype.disconnect = function() {
    var connections = this._connections;
    connections.map(function(con) {
        con.disconnect();
    });
};

/**
 * Get a free connection from the pool.
 */
TGConnectionPool.prototype.get = function() {
    //Check if any connection is available.
    var connections = this._connections;
    if (connections.length === 0) {
        //Return null for now
        return null;
    }
    return connections.pop();
    //return connections[0];
};

/**
 * Release a connection to the pool.
 */
TGConnectionPool.prototype.release = function(connection) {
    var connections = this._connections;
    connections.push(connection);
};

/**
 * Get number of free connections.
 */
TGConnectionPool.prototype.getNumFreeConnections = function() {
    var connections = this._connections;
    return connections.length;
};

TGConnectionPool.prototype.onException = function (ex, duringClose) {
    this.disconnect();

    if (!lsnr) {
        lsnr.onException(ex);
    }
};

TGConnectionPool.prototype.onReconnect = function () {
	return false;
};

TGConnectionPool.prototype.getTerminatedText = function () {
	
};


exports.TGConnectionPool = TGConnectionPool;
