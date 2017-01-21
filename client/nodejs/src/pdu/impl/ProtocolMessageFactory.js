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

var AbstractProtocolMessage   = require('../AbstractProtocolMessage').AbstractProtocolMessage,
    HandshakeRequest          = require('./HandshakeRequest').HandshakeRequest,
    HandshakeResponse         = require('./HandshakeResponse').HandshakeResponse,
    AuthenticateRequest       = require('./AuthenticateRequest').AuthenticateRequest,
    AuthenticateResponse      = require('./AuthenticateResponse').AuthenticateResponse,
    CommitTransactionResponse = require('./CommitTransactionResponse').CommitTransactionResponse,
    QueryRequest              = require('./QueryRequest').QueryRequest,
    QueryResponse             = require('./QueryResponse').QueryResponse,
    MetadataRequest           = require('./MetadataRequest').MetadataRequest,
    MetadataResponse          = require('./MetadataResponse').MetadataResponse,
    GetEntityRequest          = require('./GetEntityRequest').GetEntityRequest,
    GetEntityResponse         = require('./GetEntityResponse').GetEntityResponse,
    DisconnectChannelRequest  = require('./DisconnectChannelRequest').DisconnectChannelRequest,
    ProtocolDataInputStream   = require('./ProtocolDataInputStream').ProtocolDataInputStream,
    VerbId                    = require('./VerbId').VerbId,
    TGException               = require('../../exception/TGException').TGException,
    TGLogManager              = require('../../log/TGLogManager'),
    TGLogLevel                = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();
  
function ProtocolMessageFactory() {

}

ProtocolMessageFactory.createMessageFromVerbIdValue = function(verbIdValue, callback) {
	logger.logDebugWire( 
			'ProtocolMessageFactory.createMessageFromVerbIdValue : verbIdValue = %s',
			verbIdValue);
    var message;
    switch (verbIdValue) {
        case VerbId.HANDSHAKE_REQUEST.value :
            message = new HandshakeRequest();
            break;
        case VerbId.HANDSHAKE_RESPONSE.value :
            message = new HandshakeResponse();
            break;
        case VerbId.AUTHENTICATE_REQUEST.value :
            message = new AuthenticateRequest();
            break;
        case VerbId.AUTHENTICATE_RESPONSE.value :
            message = new AuthenticateResponse();
            break;
        //case VerbId.COMMIT_TRANS_REQUEST.value :
        //    message = new CommitTransactionRequest();
        //    break;
        case VerbId.COMMIT_TRANS_RESPONSE.value :
            message = new CommitTransactionResponse();
            break;
        case VerbId.QUERY_REQUEST.value :
            message = new QueryRequest(callback);
            break;
        case VerbId.QUERY_RESPONSE.value :
            message = new QueryResponse();
            break;
        case VerbId.METADATA_REQUEST.value :
            message = new MetadataRequest(callback);
            break;
        case VerbId.METADATA_RESPONSE.value :
            message = new MetadataResponse();
            break;
        case VerbId.GET_ENTITY_REQUEST.value :
            message = new GetEntityRequest(callback);
            break;
        case VerbId.GET_ENTITY_RESPONSE.value :
            message = new GetEntityResponse();
            break;
        case VerbId.DISCONNECT_CHANNEL_REQUEST.value :
            message = new DisconnectChannelRequest();
            break;
        default:
        	throw new TGException('Unknown message verbid : ' + verbIdValue);
    }
    return message;
};

ProtocolMessageFactory.createMessageFromVerbId = function(verbId, authToken, sessionId, callback) {
	var protocolMsg = this.createMessageFromVerbIdValue(verbId.value, callback);
	if(authToken) {
		protocolMsg.setAuthToken(authToken);
	}
	if(sessionId) {
		protocolMsg.setSessionId(sessionId);
	}
	return protocolMsg;
};

ProtocolMessageFactory.createMessage = function(buffer, offset, length) {
    var modifiedBuffer;
    if (buffer.length === length) {
        modifiedBuffer = buffer;
    } else {
        buffer.copy(modifiedBuffer, 0, offset, offset + length);
    }
    var verbIdValue = VerbId.verbIdFromBytes(buffer);
    var message = ProtocolMessageFactory.createMessageFromVerbIdValue(verbIdValue);
    message.fromBytes(modifiedBuffer);
    return message;
};

exports.ProtocolMessageFactory = ProtocolMessageFactory;