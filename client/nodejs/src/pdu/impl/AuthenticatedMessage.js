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
    TGException             = require('../../exception/TGException').TGException,
    TGLogManager            = require('../../log/TGLogManager'),
    TGLogLevel              = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function AuthenticatedMessage() {
	AuthenticatedMessage.super_.call(this);
    this._connectionId = null;
    this._clientId = null;
}

util.inherits(AuthenticatedMessage, AbstractProtocolMessage);

AuthenticatedMessage.prototype.getConnectionId = function() {
    return this._connectionId;
};

AuthenticatedMessage.prototype.setConnectionId = function(connectionId) {
    this._connectionId = connectionId;
};

AuthenticatedMessage.prototype.getClientId = function () {
    return this._clientId;
};

AuthenticatedMessage.prototype.setClientId = function (clientId) {
    this._clientId = clientId;
};

AuthenticatedMessage.prototype.writePayload = function (outputStream) {
	logger.logDebugWire( 
			"**** Entering commit AuthenticatedMessage.writePayload at output buffer position at : %d", 
			outputStream.getPosition());
    if ((this._authToken === -1) || (this._sessionId === -1)) {
        throw new Error("Message not authenticated");
    }

    outputStream.writeLong(this._authToken);
    outputStream.writeLong(this._sessionId);
    logger.logDebugWire( 
    		"**** Leaving commit AuthenticatedMessage.writePayload at output buffer position at : %d", 
    		outputStream.getPosition());
};

AuthenticatedMessage.prototype.readPayload = function (inputStream) {
    this.authToken = inputStream.readLong();
    this.sessionId = inputStream.readLong();
};

exports.AuthenticatedMessage = AuthenticatedMessage;