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

var util         = require('util'),
    VerbId       = require('./VerbId').VerbId,
    Response     = require('./Response').Response,
    TGException  = require('../../exception/TGException').TGException,
    TGLogManager = require('../../log/TGLogManager'),
    TGLogLevel   = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function AuthenticateResponse() {
	AuthenticateResponse.super_.call(this);
    this._responseStatus   = false;
}

util.inherits(AuthenticateResponse, Response);

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
	logger.logDebugWire(
			'AuthenticateResponse.prototype.readPayload ... in');
    this._responseStatus   = inputStream.readBoolean();
    this.setAuthToken(inputStream.readTGLong());
    this.setSessionId(inputStream.readTGLong());
    logger.logDebugWire(
    		'AuthenticateResponse.prototype.readPayload ... out');
};

AuthenticateResponse.prototype.writePayload = function(outputStream) {
	logger.logDebugWire(
			'AuthenticateResponse.prototype.writePayload ... in');
    outputStream.writeBoolean(this._responseStatus);
    outputStream.writeTGLong(this._authToken);
    outputStream.writeTGLong(this._sessionId);
    logger.logDebugWire(
    		'AuthenticateResponse.prototype.writePayload ... out');
};

AuthenticateResponse.prototype.getResponseStatus = function() {
    return this._responseStatus;
};

AuthenticateResponse.prototype.setResponseStatus = function(responseStatus) {
    this._responseStatus = responseStatus;
};

exports.AuthenticateResponse = AuthenticateResponse;