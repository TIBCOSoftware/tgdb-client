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

var util                      = require('util'),
    TGEntityManager           = require('../model/TGEntityManager').TGEntityManager,
    TGException               = require('../exception/TGException').TGException,
    TGLogManager              = require('../log/TGLogManager'),
    TGLogLevel                = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

var globleConnIds = 1;

function TGConnection(connectionPool, channel, properties) {
	var _connId             = globleConnIds++;
    var _connectionPool     = connectionPool;
    var _channel            = channel;
    var _properties         = properties;
    var _eventListenerMap   = {
    	exception : [],
    	connect : [],
    	disconnect : []
    };
    
    _channel.setLinkEventHandler(this);
    
    this._entities = new TGEntityManager(); 
    
    this.getGraphObjectFactory = function() {
    	return this._entities.getGraphObjectFactory();
    };

    this.insertEntity = function(tgEntity) {
    	this._entities.entityCreated(tgEntity);
    };

    this.updateEntity = function(tgEntity) {
    	this._entities.entityUpdated(tgEntity);
    };

    this.deleteEntity = function(tgEntity) {
    	this._entities.entityDeleted(tgEntity);
    };    
    
    this.getId = function() {
    	return _connectionPool._connPoolId + '-' +_connId;
    };
    
    this.getChannel = function() {
    	return _channel;
    };
    
    this.on = function(event, listener) {    	
        if (!listener || (typeof listener !== 'function')) {
            throw new TGException('listener should be a function');
        }
        if(!(!_eventListenerMap[event])) {
        	_eventListenerMap[event].push(listener);
        }
    	return this;
    };    
    
    this.dispatchEvent = function(eventType, event) {
    	var listeners = _eventListenerMap[eventType];
    	
    	if(listeners.length===0) {
    		return false;
    	}
    	
    	listeners.forEach(function(listener) {
    		listener(event);
    	});
    	return true;
    };
    
    this.throwException = function(exception) {
    	if(!this.dispatchEvent('exception', exception)){
    		throw exception;
    	}
    };

    this.handleException = function(response) {
    	this.throwException(response.getException());
    };

    /**
     *   LinkEventHandler interface
     */

    this.onException = function (exception, duringClose) {
    	this.throwException(exception);
    };

    this.onConnect = function (connectionStatus) {
    	var callback = this._callback;
    	this._callback = null;
    	if(!this.dispatchEvent('connect', connectionStatus)) {
    	    if (!callback || (typeof callback !== 'function')) {
    	        throw new TGException('listener should be a function');
    	    }
    	    callback(connectionStatus);
    	}
    };

    this.onDisconnected = function() {
    	_connectionPool.release(this);
    };
    
    this.onResponse = function(response) {
    	this.handleResponse(response);
    };

    this.getTerminatedText = function () {
    };
}

TGConnection.command = {
	CREATE : 1,
	EXECUTE : 2,
	EXECUTEID : 3,
	CLOSE : 4
};

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
TGConnection.prototype.connect = function(callback) {
    if (!(!callback) && (typeof callback !== 'function')) {
    	var exception = new TGException('Callback should be a function');
    	if(!this.throwException(exception)) {
            throw exception;
    	}
    }
    
    this._callback = callback;

    try {
    	this.getChannel().connect();
    } catch (exception) {
    	if(!this.throwException(exception)) {
            throw exception;
    	}    	
    }
};

/**
 * Disconnect from the graph DB server.
 */
TGConnection.prototype.disconnect = function() {
	logger.logDebugWire( 
			'[TGCConnection.prototype.disconnect] disconnect is called');
    this.getChannel().disconnect();
};

module.exports = TGConnection;
