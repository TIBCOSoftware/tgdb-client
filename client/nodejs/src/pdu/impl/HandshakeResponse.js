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

var util     = require('util'),
    VerbId   = require('./VerbId').VerbId,
    Response = require('./Response').Response;

function HandshakeResponse() {
	HandshakeResponse.super_.call(this);
    this._challenge = 0;
}

util.inherits(HandshakeResponse, Response);

HandshakeResponse.prototype.getVerbId = function() {
    return VerbId.HANDSHAKE_RESPONSE;
};

HandshakeResponse.prototype.getChallenge = function() {
    return this._challenge;
};

HandshakeResponse.prototype.getResponseStatus = function() {
    return this._responseStatus;
};

HandshakeResponse.prototype.isUpdateable = function() {
    return false;
};

HandshakeResponse.prototype.writePayload = function(outputStream) {

};

HandshakeResponse.prototype.readPayload = function(inputStream) {
    this._responseStatus = inputStream.readByte();
    this._challenge      = inputStream.readInt();
};

exports.HandshakeResponse      = HandshakeResponse;
exports.INVALID                = 0;
exports.ACCEPT_CHALLENGE       = 1;
exports.PROCEED_WITH_AUTH      = 2;
