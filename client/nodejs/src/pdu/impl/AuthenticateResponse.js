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

var util                    = require('util'),
    AbstractProtocolMessage = require('../AbstractProtocolMessage').AbstractProtocolMessage,
    VerbId                  = require('./VerbId').VerbId;

//Class Definition
function AuthenticateResponse() {
	AuthenticateResponse.super_.call(this);
    this._responseStatus   = false;
    this._authToken        = -1;
    this._sessionId        = -1;
}

util.inherits(AuthenticateResponse, AbstractProtocolMessage);

AuthenticateResponse.prototype.getVerbId = function() {
    return VerbId.AUTHENTICATE_RESPONSE;
};

/**
 * Is this message updateable.
 */
AuthenticateResponse.prototype.isUpdateable = function() {
    return false;
};

AuthenticateResponse.prototype.readPayload = function(inputStream) {
    this._responseStatus   = inputStream.readBoolean();
    this._authToken        = inputStream.readLong();
    this._sessionId        = inputStream.readLong();
};

AuthenticateResponse.prototype.writePayload = function(outputStream) {
//	console.log('AuthenticateResponse.prototype.writePayload ... in');
    outputStream.writeBoolean(this._responseStatus);
    outputStream.writeLong(this._authToken);
    outputStream.writeLong(this._sessionId);
//	console.log('AuthenticateResponse.prototype.writePayload ... out');
};

AuthenticateResponse.prototype.getResponseStatus = function() {
    return this._responseStatus;
};

AuthenticateResponse.prototype.setResponseStatus = function(responseStatus) {
    this._responseStatus = responseStatus;
};

AuthenticateResponse.prototype.getAuthToken = function() {
    return this._authToken;
};

AuthenticateResponse.prototype.setAuthToken = function(authToken) {
    this._authToken = authToken;
};

AuthenticateResponse.prototype.getSessionId = function() {
    return this._sessionId;
};

AuthenticateResponse.prototype.setSessionId = function(sessionId) {
    this._sessionId = sessionId;
};

exports.AuthenticateResponse = AuthenticateResponse;



