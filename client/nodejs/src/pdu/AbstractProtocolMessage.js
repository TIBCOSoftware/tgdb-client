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

var TGNumber                 = require('../datatype/TGNumber'),
    TGException              = require('../exception/TGException').TGException,
    TGLogManager             = require('../log/TGLogManager'),
    TGLogLevel               = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

var globalSequenceNo = 0;

//Class definition
function AbstractProtocolMessage() {
	var _sequenceNo = globalSequenceNo++;
	var _authToken =  TGNumber.getLong(0);
    var _sessionId =  TGNumber.getLong(0);
    var _timestamp = new Date();
    var _dataOffset = -1;
    
    var _exception = null;
    
    /**
     * Return exception.
     */
    this.getException = function() {
    	return _exception;
    };

    /**
     * Set exception.
     */
    this.setException = function(exception) {
    	_exception = exception;
    };   

    /**
     * Return current sequence number.
     */
    this.getSequenceNo = function() {
    	return _sequenceNo;
    };

    /**
     * Set current sequence number.
     */
    this.setSequenceNo = function(sequenceNo) {
    	_sequenceNo = sequenceNo;
    };
    
    this.incrementThenUpdateSequenceNo = function() {
    	this.getSequenceNo(globalSequenceNo++);
    };
    
    /**
     * Get current timestamp.
     */
    this.getTimestamp = function() {
    	if (!_timestamp) {
    		return new Date();
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

    /**
     * Get the RequestId for the Message. This will be used as the Co-relationId
     */
    //this.getRequestId = function() {
    //	return _requestId;
    //};

    /**
     * Set the RequestId for the Message.
     */
    //this.setRequestId = function(requestId) {
    //	_requestId = requestId;
    //	logger.logDebugWire( 'AbstractProtocolMessage::setRequestId %s, %s', requestId, _requestId);
    //};
    
    this.getAuthToken = function() {
    	return _authToken;
    };

    this.setAuthToken = function(authToken) {
        _authToken = authToken;
        logger.logDebugWire( 'AbstractProtocolMessage::setAuthToken %s', _authToken.getHexString());
    };
    
    this.getSessionId = function() {
        return _sessionId;
    };

    this.setSessionId = function(sessionId) {
        _sessionId = sessionId;
        logger.logDebugWire( 'AbstractProtocolMessage::setSessionId %s', _sessionId.getHexString());
    };

    this.setDataOffset = function(dataOffset) {
    	_dataOffset = dataOffset;
    };
    
}

exports.AbstractProtocolMessage = AbstractProtocolMessage;