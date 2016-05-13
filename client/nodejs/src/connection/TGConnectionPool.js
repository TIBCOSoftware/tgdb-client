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

var channel = require('../channel/AbstractChannel'),
    con     = require('../connection/TGCConnection');

//Class definition
function TGConnectionPool(channels, poolSize, properties) {
    this._connections = [];
    this._channels = channels;
    this._properties = properties;

    for (var loop = 0; loop < poolSize; loop++) {
        var channel = this._channels[loop];
        var connection = new con.TGCConnection(this, channel, properties);
        this._connections.push(connection);
    }
}

/**
 * Exception listener.
 */
TGConnectionPool.prototype.setExceptionListener = function() {

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
    if (connections.length == 0) {
        //Return null for now
        return null;
    }
    return connections.pop();
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

exports.TGConnectionPool = TGConnectionPool;
