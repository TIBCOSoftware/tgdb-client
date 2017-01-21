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
    Request      = require('./Request').Request,
    TGException  = require('../../exception/TGException').TGException,
    TGLogManager = require('../../log/TGLogManager'),
    TGLogLevel   = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function AuthenticateRequest() {
	AuthenticateRequest.super_.call(this);
    this._clientId  = null;
    this._inboxAddr = null;
    this._userName  = null;
    this._password  = null;
}

util.inherits(AuthenticateRequest, Request);

AuthenticateRequest.prototype.getVerbId = function() {
    return VerbId.AUTHENTICATE_REQUEST;
};

/**
 * Is this message updateable.
 */
AuthenticateRequest.prototype.isUpdateable = function() {
    return false;
};

AuthenticateRequest.prototype.writePayload = function(outputStream) {
	logger.logDebugWire(
			'AuthenticateRequest.prototype.writePayload  in');
    if (this._clientId === null || this._clientId.length === 0) {
        outputStream.writeBoolean(true);
    } else {
        outputStream.writeBoolean(false);
        outputStream.writeUTF(this._clientId);
    }
    if (this._inboxAddr === null || this._inboxAddr.length === 0) {
        outputStream.writeBoolean(true);
    }
    else {
        outputStream.writeBoolean(false);
        outputStream.writeUTF(this._inboxAddr);
    }
    if (this._username === null || this._username.length === 0) {
        outputStream.writeBoolean(true);
    }
    else {
        outputStream.writeBoolean(false);
        outputStream.writeUTF(this._username);
    }
    outputStream.writeBytes(this._password);
    logger.logDebugWire(
    		'AuthenticateRequest.prototype.writePayload  out');
};

AuthenticateRequest.prototype.readPayload = function(inputstream) {

};

AuthenticateRequest.prototype.getClientId = function() {
    return this._clientId;
};

AuthenticateRequest.prototype.setClientId = function(clientId) {
    this._clientId = clientId;
};

AuthenticateRequest.prototype.getUsername = function() {
    return this._username;
};

AuthenticateRequest.prototype.setUsername = function(username) {
    this._username = username;
};

AuthenticateRequest.prototype.getInboxAddr = function() {
    return this._inboxAddr;
};

AuthenticateRequest.prototype.setInboxAddr = function(inboxAddr) {
    this._inboxAddr = inboxAddr;
};

AuthenticateRequest.prototype.getPassword = function() {
    return this._password;
};

AuthenticateRequest.prototype.setPassword = function(password) {
    this._password = password;
};

exports.AuthenticateRequest = AuthenticateRequest;