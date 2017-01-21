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

var CONFIG_NAMES = require('../utils/ConfigName').CONFIG_NAMES,
    LINK_STATE   = require('./LinkState').LINK_STATE,
    TGProperties = require('../utils/TGProperties').TGProperties,
    TGException  = require('../exception/TGException').TGException,
    TGLogManager = require('../log/TGLogManager'),
    TGLogLevel   = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

/**
 *
 * @param serverURL
 * @param properties - config properties as dictionary
 * @constructor
 */
function TGChannel(tgChannelURL, properties) {
	logger.logDebug( '[TGChannel::TGChannel] serverURL = %s', tgChannelURL);
	
	this._linkEventHandler = null;
	
    var _tgChannelURL  = tgChannelURL;
    var _properties = new TGProperties(properties);
    var _linkState  = LINK_STATE.NOT_CONNECTED;
    //var _requestId  = -1;
    var _sessionId  = -1;
    var _authToken  = -1;
    
    var _requests   = {};
    
    var _numConnections = 0;

    this.setLinkEventHandler = function(handler) {
    	this._linkEventHandler = handler;
    };

    this.connect = function() {
        if (this.isConnected()) {
        	this.incrementNumConnectionAndGet();
        	logger.logDebug('[TGChannel::connect] number of connections : %d ', this.numConnections());
            return this.deliverConnectionEvent(this.isConnected());
        }
    	
    	if (this.isClosed() || _linkState === LINK_STATE.NOT_CONNECTED) {
   		    this.makeConnection();
    	} else {
            throw new TGException("Connect called on an invalid state :=" + _linkState);
        }
    };

    /**
     * Forcefully reconnect to FT urls.
     */
    this.reconnect = function() {
    };

    /**
     * Disconnect from the graph DB TCP server.
     */
    this.disconnect = function() {
        if (this.numConnections() === 0) {
        	logger.logDebug( 
        		"[TCPChannel.prototype.disconnect] Calling disconnect more than the number of connects.");
        } else {
    	    this.decrementNumConnectionAndGet();
        }
        
        if (this.numConnections() === 0) {
    	    this.stop();
        }
    };
    
    this.deliverException = function(ex) {
        if (this._linkEventHandler !== null) {
            this._linkEventHandler.onException(ex, true);
        }
        return 0;
    };
    
    this.deliverConnectionEvent = function(connectionStatus) {
        if (!(!this._linkEventHandler)) {
            this._linkEventHandler.onConnect(connectionStatus);
        }
        return 0;
    };

    this.deliverResponse = function(response) {
        if (!(!this._linkEventHandler)) {
            this._linkEventHandler.onResponse(response);
        }
        return 0;
    };
    
    this.numConnections = function() {
    	return _numConnections;
    };
    
    this.incrementNumConnectionAndGet = function() {
    	return ++_numConnections;
    };
    
    this.decrementNumConnectionAndGet = function() {
    	return --_numConnections;
    };

    this.updateLinkState = function(linkState) {
    	logger.logDebugWire( '[TGChannel::updateLinkState] %s', _linkState);
    	_linkState = linkState;
    	if(_linkState===LINK_STATE.CLOSED) {
            this._linkEventHandler.onDisconnected();    		
    	}
    };
    
    this.isConnected = function() {
        return _linkState === LINK_STATE.CONNECTED;
    };
    
    this.isClosing = function() { 
    	return _linkState === LINK_STATE.CLOSING; 
    };

    this.isClosed = function() {
        return _linkState === LINK_STATE.CLOSED;
    };

    this.putRequest = function(id, request) {
    	_requests[id] = request;
    };
    
    this.removeRequest = function(id) {
    	var request = _requests[id];
    	if(request) {
    	    delete _requests[id];
    	}
    	return request;
    };
    
    this.getSessionId = function() {
    	logger.logDebugWire( '[TGChannel::getSessionId] %s', _sessionId);
    	return _sessionId;
    };
    
    this.setSessionId  = function(sessionId) {
    	_sessionId = sessionId;
    	logger.logDebugWire( '[TGChannel::setSessionId] %s', _sessionId);
    };
    
    this.getAuthToken = function() {
    	logger.logDebugWire( '[TGChannel::getAuthToken] %s', _authToken);
    	return _authToken;
    };
    
    this.setAuthToken = function(authToken) {
    	_authToken = authToken;
    	logger.logDebugWire( '[TGChannel::setAuthToken] %s', _authToken);
    };
    
/*    
    this.getRequestId = function() {
    	logger.logDebugWire( '[TGChannel::getRequestId] %s', _requestId);
    	return _requestId;
    };
    
    this.setRequestId = function(requestId) {
    	_requestId = requestId;
    	logger.logDebugWire( '[TGChannel::setRequestId] %s', _requestId);
    };
*/    
    this.getClientId = function() {
    	return _properties.getProperty(CONFIG_NAMES.CHANNEL_CLIENTID, CONFIG_NAMES.CHANNEL_CLIENTID.defaultValue);
    };
    
    this.getUserName = function() {
    	return _properties.getProperty(CONFIG_NAMES.CHANNEL_USER_ID, CONFIG_NAMES.CHANNEL_USER_ID.defaultValue);
    };
    
    this.getPassword = function() {
    	return _properties.getProperty(CONFIG_NAMES.CHANNEL_PASSWD, CONFIG_NAMES.CHANNEL_PASSWD.defaultValue);
    };
}
exports.TGChannel = TGChannel;
