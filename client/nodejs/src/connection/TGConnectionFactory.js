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

var  CONFIG_NAMES      = require('../utils/ConfigName').CONFIG_NAMES,
    TGChannelURL       = require('../channel/TGChannelURL').TGChannelURL,
    TGChannelFactory  = require('../channel/TGChannelFactory').TGChannelFactory,
    TGConnectionPool  = require('./TGConnectionPool').TGConnectionPool,
    TGLogManager  = require('../log/TGLogManager'),
    TGLogLevel    = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();
var globleCounter = 0;
//Class definition.Factory class
function DefaultConnectionFactory() {
	var _id = globleCounter++;
	
	this.getId = function() {
		return _id;
	};
}

/**
 * Get a free connection.
 * @param serverURL
 * @param username
 * @param password
 * @param properties
 */
DefaultConnectionFactory.prototype.createConnection = function(serverURL, username, password, properties) {
	logger.logDebug("[DefaultConnectionFactory.prototype.createConnection] url = %s, username = %s, hasPassword = %s", 
			serverURL, username, !(!password));
    if (!properties) {
        properties = {};
    }

    properties[CONFIG_NAMES.USE_DEDICATED_CHN_PER_CONN.name] = 'true';
    
    var connectionPool = this.createConnectionPool(serverURL, username, password, properties.ConnectionPoolSize, properties);
    var conn = connectionPool.get();
    logger.logDebug('[DefaultConnectionFactory.prototype.createConnection] Connection id = %s', conn.getId());
    
    return conn;
};

/**
 * Create a connection pool of size defined pool size.
 * @param serverURL
 * @param username
 * @param password
 * @param poolSize
 * @param properties
 */
DefaultConnectionFactory.prototype.createConnectionPool = function(serverURL, username, password, poolSize, properties) {
	if (properties === null) {
	    properties = {};
	}
	if (!poolSize || poolSize <= 0 ) {
	    poolSize = parseInt(CONFIG_NAMES.CONN_POOL_DEFAULT_SIZE.defaultValue);
	}
	//Convert to proper url format
	var channelURL = new TGChannelURL(serverURL);
	var channels = [];
	for (var loop = 0; loop < poolSize; loop++) {
	    var channel = TGChannelFactory.createChannel(channelURL, username, password, properties);
	    channels.push(channel);
	}
	this._pool =  new TGConnectionPool(channels, poolSize, properties);
	
	return this._pool;
};

var defaultFactory = null;
var TGConnectionFactory = {
	    getFactory : function() {
		if(!defaultFactory) {
			defaultFactory = new DefaultConnectionFactory();
		}
		return defaultFactory;
	}
};

module.exports = TGConnectionFactory;
