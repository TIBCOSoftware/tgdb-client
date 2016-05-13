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

var HexUtils                  = require('../../utils/HexUtils').HexUtils,
    net                       = require('net'),
    http                      = require('http'),
    VerbId                    = require('../../pdu/impl/VerbId').VerbId,
    StringUtils               = require('../../utils/StringUtils').StringUtils,
    inherits                  = require('util').inherits,
    AbstractChannel           = require('../AbstractChannel').AbstractChannel,
    TCPConnectionEventEmitter = require('./TCPConnectionEventEmitter').TCPConnectionEventEmitter,
    HandshakeRequest          = require('../../pdu/impl/HandshakeRequest').HandshakeRequest,
    HandshakeResponse         = require('../../pdu/impl/HandshakeResponse').HandshakeResponse,
    AuthenticateResponse      = require('../../pdu/impl/AuthenticateResponse').AuthenticateResponse,
    CommitTransactionResponse = require('../../pdu/impl/CommitTransactionResponse').CommitTransactionResponse,
    QueryResponse             = require('../../pdu/impl/QueryResponse').QueryResponse,
    RequestType               = require('../../pdu/impl/RequestType').RequestType,
    ProtocolMessageFactory    = require('../../pdu/impl/ProtocolMessageFactory').ProtocolMessageFactory,
    CONFIG_NAMES              = require('../../utils/ConfigName').CONFIG_NAMES,
    LINK_STATE                = require('../LinkState').LINK_STATE,
    ACCEPT_CHALLENGE          = require('../../pdu/impl/HandshakeResponse').ACCEPT_CHALLENGE,
    PROCEED_WITH_AUTH         = require('../../pdu/impl/HandshakeResponse').PROCEED_WITH_AUTH;



//Class definition
/**
 *
 * @param serverURL - #TGChannelURL
 * @param properties - Dictionary for configuration
 * @constructor
 */
function TCPChannel(serverURL, properties) {
	TCPChannel.super_.call(this, serverURL, properties);
	this._clientSocket   = null;
	this._inboxAddr      = null;
	this._caller         = null;
	this._isConnected    = false;
	this._currentRequest = null;
}

inherits(TCPChannel, AbstractChannel);

/**
 * Connect using TCP.
 * This operation is asynchronous, upon completion
 * will invoke the callback specified in the function.
 */
TCPChannel.prototype.makeConnection = function(callback) {
	
	console.log("[TCPChannel::makeConnection] connect to ..... " + this.getHost() + ":" + this.getPort());
	
    var options = {
        host: this.getHost(),
        port: this.getPort()
    };
    var eventEmitter = new TCPConnectionEventEmitter(this);
    var channel = this;
    
    this._clientSocket = net.connect(options, function() {
        eventEmitter.emit('connected');
    });

    eventEmitter.on('connected', initiateHandshake);

    this._inboxAddr = this._clientSocket.localAddress;
    if(!this._inboxAddr) {
    	this._inboxAddr = '/192.168.1.18';
    }

    this._clientSocket.on('data', function(data) {
        //Pass socket reference
        onData(channel, data, callback);
    });
};

/**
 * Forcefully reconnect to FT urls.
 */
TCPChannel.prototype.reconnect = function() {
};

/**
 * Disconnect from the graph DB TCP server.
 */
TCPChannel.prototype.disconnect = function() {
    this._clientSocket.end();
};

/**
 * Start channel for send and receive messages.
 */
TCPChannel.prototype.start = function() {
};

/**
 * Stop TCP channel.
 */
TCPChannel.prototype.stop = function() {
};

/**
* Send message to Graph DB server.
*/
TCPChannel.prototype.send = function(request, caller) {
	if(null!==this._caller) {
		throw new Error('[TCPChannel.send] Channel is busy for previous request!');
	}
	this._caller = caller;
	this._request = request;
	
	send (this._clientSocket, request);
};

function send (socket, request) {
	var buffer = request.toBytes();
	console.log("----------------- outgoing message --------------------");
	console.log("name : " + request.getVerbId().name);
	console.log("Send buf : Formatted Byte Array:");
	console.log(HexUtils.formatHex(buffer));
	socket.write(buffer);
}
	
/**
 * Internal method
 */
function initiateHandshake() {
    var handshakeRequest = new HandshakeRequest();
    handshakeRequest.setRequestType(RequestType.Type.INITIATE);
    send(this.getChannel()._clientSocket, handshakeRequest);
}

/**
 * Upon receiving data.
 * @param channel
 * @param data
 * @param callback - Mandatory function
 */
function onData(channel, data, callback) {
    //Read wire message, all about responses from server
    var response = readWireMsg(data);

    // for further request
    channel.setRequestId(response.getRequestId()); 

    if (response instanceof AuthenticateResponse) {
        var responseStatus = response.getResponseStatus();
        if (responseStatus) {
            channel.setAuthToken(response.getAuthToken());
            channel.setSessionId(response.getSessionId());
        }
        //Set state appropriately
        channel._linkState = LINK_STATE.CONNECTED;
        //Call back the function
        callback(responseStatus);
    } else if (response instanceof HandshakeResponse) {
        if (response.getResponseStatus() == ACCEPT_CHALLENGE) {
            //Perform handshake part 2.
            completeHandshake(channel, response);
        } else if (response.getResponseStatus() == PROCEED_WITH_AUTH) {
            //Perform authentication
            performAuthentication(channel);
        } else {
            throw new Error('Handshake failed with DB Server');
        }
    } else if (response instanceof CommitTransactionResponse) {
        console.log('onData for CommitTransactionResponse ......... set response : ' + response.getRequestId() + ', ' + response);
        if(null!==this._caller) {
        	var caller = channel._caller;
        	channel._caller = null;
        	caller.handleCommitResponse(response);
        }
    } else if (response instanceof QueryResponse) {
        console.log('onData for QueryResponse ......... set response : ' + response.getRequestId() + ', ' + response);
        if(null!==this._caller) {
        	var caller = channel._caller;
        	channel._caller = null;
        	var request = channel._request;
        	caller.handleQueryResponse(request, response);
        }
    } else {
    	throw new Error('Unknow response from server');
    }
}

/**
 * Complete second part of handshake
 * @param response
 */
function completeHandshake(channel, response) {
    var challenge = response.getChallenge();
    var handshakeRequest = ProtocolMessageFactory.createMessageFromVerbId(VerbId.HANDSHAKE_REQUEST);
    handshakeRequest.updateSequenceAndTimestamp(-1);
    handshakeRequest.setRequestType(RequestType.Type.CHALLENGE_ACCEPTED);
    //TODO Handle SSL later
    handshakeRequest.setSSLMode(false);
    handshakeRequest.setChallenge(challenge * 2 / 3);
    return send(channel._clientSocket, handshakeRequest);
}

/**
 * Perform authn with server
 * @param channel
 */
function performAuthentication(channel) {
    var authenticateRequest = ProtocolMessageFactory.createMessageFromVerbId(VerbId.AUTHENTICATE_REQUEST);
    authenticateRequest.setClientId(channel.getClientId);
    authenticateRequest.setInboxAddr(channel._inboxAddr);
    authenticateRequest.setUsername(channel.getUserName());
    authenticateRequest.setPassword(StringUtils.toBytes(channel.getPassword()));
    return send(channel._clientSocket, authenticateRequest);
}

function readWireMsg(inputData) {
    var message = ProtocolMessageFactory.createMessage(inputData, 0, inputData.length);
//  console.log("----------------- incoming message --------------------");
    console.log("name : " + message.getVerbId().name);
    console.log("ReadMsg : Formatted Byte Array:");
    console.log(HexUtils.formatHex(inputData));
    if (message.getVerbId().value == VerbId.EXCEPTION_MESSAGE.value) {
        throw new Error('Exception during handshake');
    }
    return message;
}

exports.TCPChannel = TCPChannel;


