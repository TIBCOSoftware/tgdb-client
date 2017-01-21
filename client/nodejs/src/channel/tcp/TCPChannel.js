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

var net                       = require('net'),
    http                      = require('http'),
    inherits                  = require('util').inherits,
    TCPChannelURL             = require('./TCPChannelURL').TCPChannelURL,
    TCPConnectionEventEmitter = require('./TCPConnectionEventEmitter').TCPConnectionEventEmitter,
    TGChannel                 = require('../TGChannel').TGChannel,
    LINK_STATE                = require('../LinkState').LINK_STATE,
    VerbId                    = require('../../pdu/impl/VerbId').VerbId,
    HandshakeRequest          = require('../../pdu/impl/HandshakeRequest').HandshakeRequest,
    HandshakeResponse         = require('../../pdu/impl/HandshakeResponse').HandshakeResponse,
    AuthenticateResponse      = require('../../pdu/impl/AuthenticateResponse').AuthenticateResponse,
    RequestType               = require('../../pdu/impl/RequestType').RequestType,
    ProtocolMessageFactory    = require('../../pdu/impl/ProtocolMessageFactory').ProtocolMessageFactory,
    ACCEPT_CHALLENGE          = require('../../pdu/impl/HandshakeResponse').ACCEPT_CHALLENGE,
    PROCEED_WITH_AUTH         = require('../../pdu/impl/HandshakeResponse').PROCEED_WITH_AUTH,
    TGException               = require('../../exception/TGException').TGException,
    CONFIG_NAMES              = require('../../utils/ConfigName').CONFIG_NAMES,
    HexUtils                  = require('../../utils/HexUtils').HexUtils,
    StringUtils               = require('../../utils/StringUtils').StringUtils,
    TGLogManager              = require('../../log/TGLogManager'),
    TGLogLevel                = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

/**
 *
 * @param tgChannelURL - #TGChannelURL
 * @param properties - Dictionary for configuration
 * @constructor
 */
function TCPChannel(tgChannelURL, properties) {
	TCPChannel.super_.call(this, tgChannelURL, properties);
	
	this._channelURL     = new TCPChannelURL(tgChannelURL);
	this._clientSocket   = null;
	this._inboxAddr      = null;
	this._isConnected    = false;
	this._currentRequest = null;
	this._buffer         = null;
	this._contentLength  = 0;
	this._currentPosition = 0;
	
    this.getHost = function() {
    	return this._channelURL.getHost();
    };

    this.getPort = function() {
    	return this._channelURL.getPort();
    };
}

inherits(TCPChannel, TGChannel);

/**
 * Connect using TCP.
 * This operation is asynchronous, upon completion
 * will invoke the callback specified in the function.
 */
TCPChannel.prototype.makeConnection = function() {
	
	logger.logDebug( 
		'[TCPChannel.prototype.makeConnection] connect to ..... %s:%d',
		this.getHost(), this.getPort());
	this._accumulatedLength = 0;
	
    var options = {
        host: this.getHost(),
        port: this.getPort()
    };
    
    var eventEmitter = new TCPConnectionEventEmitter(this);
    eventEmitter.on('connected', initiateHandshake);
    var channel = this;

    this._clientSocket = net.connect(options, function() {
    	channel.updateLinkState(LINK_STATE.CONNECTED);
    	channel.incrementNumConnectionAndGet();
    	logger.logDebug( 
    		'[TCPChannel.prototype.makeConnection] number of connections : %d', 
    		channel.numConnections());
        eventEmitter.emit('connected');
    });

    this._inboxAddr = this._clientSocket.localAddress;
    if(!this._inboxAddr) {
    	this._inboxAddr = ''; // read from environment file?
    }
   
    this._clientSocket.on('data', function(data) {
    	try {
        	if(!channel._buffer) {
        		channel._contentLength = data.readInt32BE();
        		channel._buffer = new Buffer(channel._contentLength);
        	}
        	
        	data.copy(channel._buffer, channel._currentPosition, 0, data.length);
        	channel._currentPosition += data.length;
    		
        	if(channel._currentPosition<channel._contentLength) {
        		// keep accumulate data
        		return;
        	}
        	
            //Pass socket reference
            onData(channel, channel._buffer);
            
            channel._buffer = null;
            channel._contentLength = 0;
            channel._currentPosition = 0;
    	} catch (exception) {
    		channel.deliverException(exception);
    	}
    });
   
    this._clientSocket.on('end', function(data) {
    	logger.logDebug( 
    		'[TCPChannel.prototype.makeConnection] Socket end called, will close shortly.');
    });
    
    this._clientSocket.on('close', function(data) {
        channel.updateLinkState(LINK_STATE.CLOSED);
    	logger.logDebug( '[TCPChannel.prototype.makeConnection] Socket closed.');
    });
    
    this._clientSocket.on('error', function(error) {
    	if(!(!channel._clientSocket)) {
            channel.updateLinkState(LINK_STATE.CLOSED);
            channel._clientSocket.end();
    	}
        channel.deliverException(new TGException(error.name + ' : ' + error.message));        	
    });
};

/**
 * Start channel for send and receive messages.
 */
TCPChannel.prototype.start = function() {
    if (!this.isConnected) {
        throw new TGException("Channel not connected");
    }
};

/**
 * Stop TCP channel.
 */
TCPChannel.prototype.stop = function() {
	this.stop(false);
};

TCPChannel.prototype.stop = function(bForcefully) {
	logger.logDebug( 
		'[TCPChannel.prototype.stop] number of connections : %s', this.numConnections());

    try {
        if (this.isClosed() || this.isClosing()) {
        	logger.logDebug( 
        			'[TCPChannel.prototype.stop] Stop a channel while it\'s closing or closed.');
            return;
        } else if (!this.isConnected()) {
        	logger.logDebug( 
			'[TCPChannel.prototype.stop] Stop a channel while it\'s not connected.');
            return;
        }

        if ((bForcefully) || (this.numConnections() === 0))
        {
            // Send the disconnect request.
            var request = ProtocolMessageFactory.createMessageFromVerbId(VerbId.DISCONNECT_CHANNEL_REQUEST);
            // sendRequest will not receive a channel response since the channel will be disconnected.
            this.send(request);

            this.updateLinkState(LINK_STATE.CLOSING);
            this._clientSocket.end();
        }
    } catch (error) {
    	logger.logInfo( 
    		'[TCPChannel.prototype.stop] Error when stop channel, message : %s',
    		error.message);
    	this._clientSocket.destroy();
    } finally {
        this._clientSocket = null;
        this.updateLinkState(LINK_STATE.CLOSED);
    }
};

/**
* Send message to Graph DB server.
*/
TCPChannel.prototype.send = function(request) {
	this.putRequest(request.getRequestId().getHexString(), request);
	send (this._clientSocket, request);
};

/**
 *  Private methods
 */

function send (socket, request) {
	var buffer = request.toBytes();
	logger.logDebugWire( "----------------- outgoing message --------------------");
	logger.logDebugWire( "name : %s", request.getVerbId().name);
	logger.logDebugWire( HexUtils.formatHex(buffer));
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
function onData(channel, data) {
    //Read wire message, all about responses from server
    var response = readWireMsg(data);
    logger.logDebugWire( '[TCPChannel::onData] %s', response.getVerbId().name);
    // for further request
    //channel.setRequestId(response.getRequestId()); 

    if (response instanceof AuthenticateResponse) {
        var responseStatus = response.getResponseStatus();
        if (responseStatus) {
            channel.setAuthToken(response.getAuthToken());
            channel.setSessionId(response.getSessionId());
        }
        //Set state appropriately
        channel._linkState = LINK_STATE.CONNECTED;
        //Call back the function
        channel.deliverConnectionEvent(responseStatus);
    } else if (response instanceof HandshakeResponse) {
        if (response.getResponseStatus() === ACCEPT_CHALLENGE) {
            //Perform handshake part 2.
            completeHandshake(channel, response);
        } else if (response.getResponseStatus() === PROCEED_WITH_AUTH) {
            //Perform authentication
            performAuthentication(channel);
        } else {
            throw new TGException('Handshake failed with DB Server');
        }
    } else {
    	response.setRequest(channel.removeRequest(response.getRequestId().getHexString()));
    	channel.deliverResponse(response);
    }
}

/**
 * Complete second part of handshake
 * @param response
 */
function completeHandshake(channel, response) {
    var challenge = response.getChallenge();
    var handshakeRequest = ProtocolMessageFactory.createMessageFromVerbId(VerbId.HANDSHAKE_REQUEST);
    handshakeRequest.updateSequenceAndTimestamp(new Date());
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
    logger.logDebugWire("----------------- incoming message --------------------");
    logger.logDebugWire("name : " + message.getVerbId().name);
    logger.logDebugWire(HexUtils.formatHex(inputData));
    if (message.getVerbId().value === VerbId.EXCEPTION_MESSAGE.value) {
        throw new TGException('Exception during handshake');
    }
    return message;
}

exports.TCPChannel = TCPChannel;