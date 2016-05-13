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

var ProtocolDataInputStream = require('../pdu/impl/ProtocolDataInputStream').ProtocolDataInputStream,
    ProtocolDataOutputStream = require('../pdu/impl/ProtocolDataOutputStream').ProtocolDataOutputStream,
    TGProtocolVersion       = require('../TGProtocolVersion').TGProtocolVersion;

var globalSequenceNo = 0;

//Class definition
function AbstractProtocolMessage() {
	var _sequenceNo = globalSequenceNo;
	globalSequenceNo++;
    var _requestId = -1;
	var _authToken =  0;
    var _sessionId =  0;
    var _timestamp = -1;
    var _bufLength = -1;
    var_dataOffset = -1;

    /**
     * Return current sequence number.
     */
    this.getSequenceNo = function() {
    	return _sequenceNo;
    };

    /**
     * Get current timestamp.
     */
    this.getTimestamp = function() {
    	if (_timestamp == null || _timestamp == -1) {
    		return new Date().getTime();
    	}
    	return _timestamp;
    };

    /**
     * set current timestamp.
     */
    this.setTimestamp = function(timestamp) {
    	if (this.isUpdateable()) {
    		this._timestamp = timestamp;
    	}
    };


    this.updateSequenceAndTimestamp = function(timestamp) {
    	if (this.isUpdateable()) {
    		_sequenceNo = globalSequenceNo;
    		globalSequenceNo++;
    		console.log('after updateSequenceAndTimestamp -----------------> ' + _sequenceNo);
    		_timestamp  = timestamp;
    		_bufLength  = -1;
    	} else {
    		throw new Error('Mutating a readonly message');
    	}
    };

    /**
     * Get the RequestId for the Message. This will be used as the Co-relationId
     */
    this.getRequestId = function() {
    	return _requestId;
    };

    /**
     * Set the RequestId for the Message.
     */
    this.setRequestId = function(requestId) {
    	_requestId = requestId;
    	console.log('AbstractProtocolMessage::setRequestId ' + requestId + ', ' + _requestId);
    };
    
    this.getAuthToken = function() {
    	return _authToken;
    };

    this.setAuthToken = function(authToken) {
        _authToken = authToken;
    };
    
    this.getSessionId = function() {
        return _sessionId;
    };

    this.setSessionId = function(sessionId) {
        _sessionId = sessionId;
    };

    /**
     * Convert message to bytes.
     * TODO Override.
     */
    this.toBytes = function() {
    	console.log("Entering MessageUtils.messageToBytes ..............");
    	var outputStream = new ProtocolDataOutputStream();
    	//Write to this array and finally convert to buffer
    	writeHeader(this, outputStream);
    	//Write payload
    	this.writePayload(outputStream);
    	//Write length at the beginning
    	outputStream.writeIntAt(0, outputStream.length());
    	return outputStream.toBuffer();
    };

    /**
     * Get the MessageByteBufLength, call this method after the toBytes is called.
     */
    this.getMessageByteBufLength = function() {

    };

    this.readHeader = function(inputStream) {
    	var magic = inputStream.readInt();
    	if (magic != TGProtocolVersion.getMagic()) {
    		throw new Error('Bad Message Magic');
    	}

    	var protocolVersion = inputStream.readShort();
    	if (!TGProtocolVersion.isCompatible(protocolVersion)) {
    		throw new Error('Unsupported Protocol version');
    	}

    	var verbId = inputStream.readShort();

    	_sequenceNo = inputStream.readLong();
    	_timestamp  = inputStream.readLong();
    	_requestId  = inputStream.readLong();
    	
        // Moved from AuthenticatedMessage
        _authToken 	 = inputStream.readLong();
        _sessionId    = inputStream.readLong();

    	_dataOffset = inputStream.readShort();
    };
    
    /**
     * Reconstruct the message from the buffer.
     * @param byteBuffer
     */
    this.fromBytes = function(byteBuffer) {

    };
}

function writeHeader(message, outputStream) {
	console.log("**** Entering commit AbstractProtocolMessage.writeHeader at output buffer position at : %d", outputStream.getPosition());
    outputStream.writeInt(0);
    outputStream.writeInt(TGProtocolVersion.getMagic());
    outputStream.writeShort(TGProtocolVersion.getProtocolVersion());
    outputStream.writeShort(message.getVerbId().value);
    outputStream.writeLong(message.getSequenceNo());
    outputStream.writeLong(message.getTimestamp());
    outputStream.writeLong(message.getRequestId());
    
    // Moved from AuthenticatedMessage
    outputStream.writeLong(message.getAuthToken());
    outputStream.writeLong(message.getSessionId());
    
    outputStream.writeShort(outputStream.length() + 2);
	console.log("**** Leaving commit AbstractProtocolMessage.writeHeader at output buffer position at : %d", outputStream.getPosition());
};

exports.AbstractProtocolMessage = AbstractProtocolMessage;