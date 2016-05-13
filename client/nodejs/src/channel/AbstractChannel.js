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
    TGProperties = require('../utils/TGProperties').TGProperties;

/**
 *
 * @param serverURL
 * @param properties - config properties as dictionary
 * @constructor
 */
function AbstractChannel(serverURL, properties) {
	console.log("[AbstractChannel::AbstractChannel] serverURL = " + serverURL);
    var _serverURL  = serverURL;
    var _properties = new TGProperties(properties);
    var _linkState  = LINK_STATE.NOT_CONNECTED;
    var _requestId  = -1;
    var _sessionId  = -1;
    var _authToken  = -1;

    this.connect = function(callback) {
    	if (_linkState == undefined || _linkState == LINK_STATE.NOT_CONNECTED) {
    		this.makeConnection(callback);
    	}
    };

    this.getSessionId = function() {
    	console.log('AbstractChannel::getSessionId ' + _sessionId);
    	return _sessionId;
    };
    
    this.setSessionId  = function(sessionId) {
    	_sessionId = sessionId;
    	console.log('AbstractChannel::setSessionId  ' + _sessionId);
    }
    
    
    this.getAuthToken = function() {
    	console.log('AbstractChannel::getAuthToken ' + _authToken);
    	return _authToken;
    };
    
    this.setAuthToken = function(authToken) {
    	_authToken = authToken;
    	console.log('AbstractChannel::setAuthToken ' + _authToken);
    }
    
    
    this.getRequestId = function() {
    	console.log('AbstractChannel::getRequestId ' + _requestId);
    	return _requestId;
    };
    
    this.setRequestId = function(requestId) {
    	_requestId = requestId;
    	console.log('AbstractChannel::setRequestId ' + _requestId);
    }
    
    this.getHost = function() {
    	return _serverURL.getHost();
    };

    this.getPort = function() {
    	return _serverURL.getPort();
    };

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
exports.AbstractChannel = AbstractChannel;
