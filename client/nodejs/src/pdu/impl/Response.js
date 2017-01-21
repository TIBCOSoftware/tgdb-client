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

var util                     = require('util'),
    VerbId                   = require('./VerbId').VerbId,
    ProtocolDataInputStream  = require('./ProtocolDataInputStream').ProtocolDataInputStream,
    AbstractProtocolMessage  = require('../AbstractProtocolMessage').AbstractProtocolMessage,
    TGProtocolVersion        = require('../../TGProtocolVersion').TGProtocolVersion,
    TGException              = require('../../exception/TGException').TGException,
    TGLogManager             = require('../../log/TGLogManager'),
    TGLogLevel               = require('../../log/TGLogger').TGLogLevel,
    TGNumber                 = require('../../datatype/TGNumber');

var logger = TGLogManager.getLogger();

function Response() {
	Response.super_.call(this);
	this._request = null;
    this._requestId = TGNumber.getLong(-1);
}

util.inherits(Response, AbstractProtocolMessage);

Response.prototype.setRequest = function(request){
	this._request = request;
};

Response.prototype.getRequest = function(){
	return this._request;
};

Response.prototype.setRequestId = function(requestId) {
	this._requestId = requestId;
	logger.logDebugWire( 
			'[Response.prototype.setRequestId] %s', 
			requestId.getHexString());
};

Response.prototype.getRequestId = function() {
    return this._requestId;
};

/**
 * Reconstruct the message from the buffer.
 * @param byteBuffer
 */
Response.prototype.fromBytes = function(byteBuffer) {
    var protocolInputStream = new ProtocolDataInputStream(byteBuffer);
    var length = protocolInputStream.readInt();
    if (length !== byteBuffer.length) {
        throw new TGException('Buffer length mismatch');
    }
    //Read header
    this.readHeader(this, protocolInputStream);
    //Read payload
    this.readPayload(protocolInputStream);
};

Response.prototype.readHeader = function(message, inputStream) {
	logger.logDebugWire( 
			"**** Entering AbstractProtocolMessage.readHeader at input buffer position at : %d", 
			inputStream.getPosition());
	var magic = inputStream.readInt();
	if (magic !== TGProtocolVersion.getMagic()) {
		throw new Error('Bad Message Magic');
	}

	var protocolVersion = inputStream.readShort();
	if (!TGProtocolVersion.isCompatible(protocolVersion)) {
		throw new Error('Unsupported Protocol version');
	}

	var verbId = inputStream.readShort();

	message.setSequenceNo(inputStream.readLong());
	message.setTimestamp(inputStream.readDate());
	message.setRequestId(inputStream.readTGLong());
	
    // Moved from AuthenticatedMessage
    message.setAuthToken(inputStream.readTGLong());
    message.setSessionId(inputStream.readTGLong());

	message.setDataOffset(inputStream.readShort());
	logger.logDebugWire( 
			"**** Leaving AbstractProtocolMessage.readHeader at input buffer position at : %d", 
			inputStream.getPosition());
};

exports.Response = Response;