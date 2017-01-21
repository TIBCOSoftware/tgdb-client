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
    ProtocolDataOutputStream = require('./ProtocolDataOutputStream').ProtocolDataOutputStream,
    AbstractProtocolMessage  = require('../AbstractProtocolMessage').AbstractProtocolMessage,
    TGProtocolVersion        = require('../../TGProtocolVersion').TGProtocolVersion,
    TGLogManager             = require('../../log/TGLogManager'),
    TGLogLevel               = require('../../log/TGLogger').TGLogLevel,
    TGNumber                 = require('../../datatype/TGNumber');

var logger = TGLogManager.getLogger();

var counter = 0;
function Request() {
	Request.super_.call(this);
    this._requestId = TGNumber.getLong(counter++);
    this._bufLength = -1;
}

util.inherits(Request, AbstractProtocolMessage);

Request.prototype.getRequestId = function() {
    return this._requestId;
};

/**
 * Convert message to bytes.
 * TODO Override.
 */
Request.prototype.toBytes = function() {
	var outputStream = new ProtocolDataOutputStream();
	//Write to this array and finally convert to buffer
	this.writeHeader(this, outputStream);
	//Write payload
	this.writePayload(outputStream);
	//Write length at the beginning
	this._bufLength = outputStream.length();
	outputStream.writeIntAt(0, this._bufLength);
	return outputStream.toBuffer();
};

/**
 * Get the MessageByteBufLength, call this method after the toBytes is called.
 */
Request.prototype.getMessageByteBufLength = function() {
    return this._bufLength;
};

Request.prototype.writeHeader = function(message, outputStream) {
	logger.logDebugWire( 
			"**** Entering AbstractProtocolMessage.writeHeader at output buffer position at : %d", 
			outputStream.getPosition());
    outputStream.writeInt(0);
    outputStream.writeInt(TGProtocolVersion.getMagic());
    outputStream.writeShort(TGProtocolVersion.getProtocolVersion());
    outputStream.writeShort(message.getVerbId().value);
    outputStream.writeLong(message.getSequenceNo());
    outputStream.writeDate(message.getTimestamp());
    outputStream.writeTGLong(message.getRequestId());
    
    // Moved from AuthenticatedMessage
    outputStream.writeTGLong(message.getAuthToken());    
    outputStream.writeTGLong(message.getSessionId());
    
    outputStream.writeShort(outputStream.length() + 2);
    logger.logDebugWire( 
    		"**** Leaving AbstractProtocolMessage.writeHeader at output buffer position at : %d", 
    		outputStream.getPosition());
};

Request.prototype.writePayload = function(dynamicBuffer) {
	logger.logDebugWire( '[Request.prototype.writePayload] should override this method.');
};

exports.Request = Request;