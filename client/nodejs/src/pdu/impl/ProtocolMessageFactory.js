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
    PingMessage               = require('./PingMessage').PingMessage,
    HandshakeRequest          = require('./HandshakeRequest').HandshakeRequest,
    HandshakeResponse         = require('./HandshakeResponse').HandshakeResponse,
    AuthenticateRequest       = require('./AuthenticateRequest').AuthenticateRequest,
    AuthenticateResponse      = require('./AuthenticateResponse').AuthenticateResponse,
    CommitTransactionResponse = require('./CommitTransactionResponse').CommitTransactionResponse,
    QueryRequest              = require('./QueryRequest').QueryRequest,
    QueryResponse             = require('./QueryResponse').QueryResponse,
    ProtocolDataInputStream   = require('./ProtocolDataInputStream').ProtocolDataInputStream,
    VerbId                    = require('./VerbId').VerbId;
  
function ProtocolMessageFactory() {

}

ProtocolMessageFactory.createMessageFromVerbIdValue = function(verbIdValue) {
	console.log('ProtocolMessageFactory.createMessageFromVerbIdValue : verbIdValue = ' + verbIdValue);
    var message;
    switch (verbIdValue) {
        case VerbId.PING_MESSAGE.value :
            message = new PingMessage();
            break;
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
            message = new QueryRequest();
            break;
        case VerbId.QUERY_RESPONSE.value :
            message = new QueryResponse();
            break;
        default:
        	throw new Error('Unknown message verbid : ' + verbIdValue);
    }
    return message;
};

ProtocolMessageFactory.createMessageFromVerbId = function(verbId) {
	return this.createMessageFromVerbIdValue(verbId.value);
};

ProtocolMessageFactory.createMessage = function(buffer, offset, length) {
    var modifiedBuffer;
    if (buffer.length == length) {
        modifiedBuffer = buffer;
    } else {
        buffer.copy(modifiedBuffer, 0, offset, offset + length);
    }
    var verbIdValue = VerbId.verbIdFromBytes(buffer);
    var message = ProtocolMessageFactory.createMessageFromVerbIdValue(verbIdValue);
    messageFromBytes(message, modifiedBuffer);
    return message;
};

function messageFromBytes(message, byteBuffer) {
    var protocolInputStream = new ProtocolDataInputStream(byteBuffer);
    var length = protocolInputStream.readInt();
    if (length != byteBuffer.length) {
        throw new Error('Buffer length mismatch');
    }
    //Read header
    message.readHeader(protocolInputStream);
    //Read payload
    message.readPayload(protocolInputStream);
}

exports.ProtocolMessageFactory = ProtocolMessageFactory;