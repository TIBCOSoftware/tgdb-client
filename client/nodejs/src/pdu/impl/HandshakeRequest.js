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

var util        = require('util'),
    VerbId      = require('./VerbId').VerbId,
    Request     = require('./Request').Request,
    TGException = require('../../exception/TGException').TGException;

function HandshakeRequest() {
	HandshakeRequest.super_.call(this);
    this._sslMode = false;
    this._challenge = 0;
}

util.inherits(HandshakeRequest, Request);

HandshakeRequest.prototype.getVerbId = function() {
    return VerbId.HANDSHAKE_REQUEST;
};

HandshakeRequest.prototype.isUpdateable = function() {
    return true;
};

HandshakeRequest.prototype.updateSequenceAndTimestamp = function(timestamp) {
	if (this.isUpdateable()) {
		this.incrementThenUpdateSequenceNo();
		if(!(!timestamp)) {
			this.setTimestamp(timestamp);
		}
		//_bufLength  = -1;
	} else {
		throw new TGException('Mutating a readonly message');
	}
};

HandshakeRequest.prototype.writePayload = function(outputStream) {
    outputStream.writeByte(this._requestType.valueOf());
    //SSL enabled or not
    outputStream.writeBoolean(this._sslMode);
    //Challenge parameter
    outputStream.writeInt(this._challenge);
};

HandshakeRequest.prototype.readPayload = function() {

};

HandshakeRequest.prototype.getSSLMode = function() {
    return this._sslMode;
};

HandshakeRequest.prototype.setSSLMode = function(sslMode) {
    this._sslMode = sslMode;
};

HandshakeRequest.prototype.getChallenge = function() {
    return this._challenge;
};

HandshakeRequest.prototype.setChallenge = function(challenge) {
    this._challenge = challenge;
};

HandshakeRequest.prototype.getRequestType = function() {
    return this.getRequestType();
};

HandshakeRequest.prototype.setRequestType = function(requestType) {
    this._requestType = requestType;
};

exports.HandshakeRequest = HandshakeRequest;