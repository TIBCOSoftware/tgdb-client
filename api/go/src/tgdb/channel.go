/*
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 *
 * File name: channel.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: channel.go 4539 2020-10-22 17:53:09Z nimish $
 */

package tgdb

import (
	"bytes"
	"sync"
)

type TGChannel interface {
	// ChannelLock locks the communication channel between TGDB client and server
	ChannelLock()
	// ChannelUnlock unlocks the communication channel between TGDB client and server
	ChannelUnlock()
	// Connect connects the underlying channel using the URL end point
	Connect() TGError
	// Disconnect disconnects the channel from its URL end point
	Disconnect() TGError
	// DisablePing disables the pinging ability to the channel
	DisablePing()
	// ExceptionLock locks the communication channel between TGDB client and server in case of business exceptions
	ExceptionLock()
	// ExceptionUnlock unlocks the communication channel between TGDB client and server in case of business exceptions
	ExceptionUnlock()
	// GetAuthToken gets Authorization Token
	GetAuthToken() int64
	// GetClientId gets Client Name
	GetClientId() string
	// GetClientProtocolVersion gets Client Protocol Version
	//GetClientProtocolVersion() int
	// GetChannelURL gets the Channel URL
	GetChannelURL() TGChannelUrl
	// GetConnectionIndex gets the Connection Index
	GetConnectionIndex() int
	// GetDataCryptoGrapher gets the Data Crypto Grapher
	GetDataCryptoGrapher() TGDataCryptoGrapher
	// GetExceptionCondition gets the Exception Condition
	GetExceptionCondition() *sync.Cond
	// GetLinkState gets the Link/Channel State
	GetLinkState() LinkState
	// GetNoOfConnections gets number of connections this channel has
	GetNoOfConnections() int32
	// GetPrimaryURL gets the Primary URL
	GetPrimaryURL() TGChannelUrl
	// GetProperties gets the Channel Properties
	GetProperties() TGProperties
	// GetReader gets the Channel Reader
	GetReader() TGChannelReader
	// GetResponses gets the Channel Response Map
	GetResponses() map[int64]TGChannelResponse
	// GetServerProtocolVersion gets Server Protocol Version
	//GetServerProtocolVersion() int
	// GetSessionId gets Session id
	GetSessionId() int64
	// GetTracer gets the channel Tracer
	GetTracer() TGTracer
	// EnablePing enables the pinging ability to the channel
	EnablePing()
	// IsChannelPingable checks whether the channel is pingable or not
	IsChannelPingable() bool
	// IsClosed checks whether channel is open or closed
	IsClosed() bool
	// SendMessage sends a Message on this channel, and returns immediately - An Asynchronous or Non-Blocking operation
	SendMessage(msg TGMessage) TGError
	// SendRequest sends a Message, waits for a response in the message format, and blocks the thread till it gets the response
	SendRequest(msg TGMessage, response TGChannelResponse) (TGMessage, TGError)
	// SetChannelLinkState sets the Link/Channel State
	SetChannelLinkState(state LinkState)
	// SetChannelURL sets the channel URL
	SetChannelURL(channelUrl TGChannelUrl)
	// SetConnectionIndex sets the connection index
	SetConnectionIndex(index int)
	// SetNoOfConnections sets number of connections
	SetNoOfConnections(count int32)
	// SetResponse sets the ChannelResponse Map
	SetResponse(reqId int64, response TGChannelResponse)
	// SetAuthToken sets the Authorization Token
	SetAuthToken(token int64)
	// DoAuthenticate helps authenticate the external user's credentials
	DoAuthenticate() TGError
	// Start starts the channel so that it can send and receive messages
	Start() TGError
	// Stop stops the channel forcefully or gracefully
	Stop(bForcefully bool)

	// Additional in GO - Abstract declared in Java - to implement function inheritance
	// CreateSocket creates a network socket to transfer the messages in the byte format
	CreateSocket() TGError
	// CloseSocket closes the network socket
	CloseSocket() TGError
	// OnConnect executes functional logic after successfully establishing the connection to the server
	OnConnect() TGError
	// ReadWireMsg read the message from the network in the byte format
	ReadWireMsg() (TGMessage, TGError)
	// Send sends the message to the server, compress and or encrypt.
	// Hence it is abstraction, that the Channel knows about it.
	// @param msg       The message that needs to be sent to the server
	Send(msg TGMessage) TGError
	// Additional Method to help debugging
	String() string
}

// Channel URL is an encapsulation of all the attributes necessary to construct a valid and meaningful URL to connect to TGDB server
type TGChannelUrl interface {
	// GetFTUrls gets the Fault Tolerant URLs
	GetFTUrls() []TGChannelUrl
	// GetHost gets the host part of the URL
	GetHost() string
	// GetPort gets the port on which it is connected
	GetPort() int
	// GetProperties gets the URL Properties
	GetProperties() TGProperties
	// GetProtocol gets the protocol used as part of the URL
	GetProtocol() TGProtocol
	// GetUrlAsString gets the string form of the URL
	GetUrlAsString() string
	// GetUser gets the user associated with the URL
	GetUser() string
}

// ======= Protocol Types Supported =======
type TGProtocol int

const (
	ProtocolTCP TGProtocol = 1 << iota
	ProtocolSSL
	ProtocolHTTP
	ProtocolHTTPS
)

// ======= Link State Types =======
type LinkState int

/*
func (state LinkState) String() interface{} {
	
}
*/


const (
	LinkNotConnected LinkState = 1 << iota
	LinkConnected
	LinkClosing
	LinkClosed
	LinkFailedOnSend
	LinkFailedOnRecv
	LinkFailedOnProcessing
	LinkReconnecting
	LinkTerminated
)

// Channel Reader is an independent thread that starts and stops with the channel, and continuously monitors
// network communication socket to read and process any message that is sent by the TGDB server
type TGChannelReader interface {
	// Start starts the channel reader
	Start()
	// Stop stops the channel reader
	Stop()
	// Additional Method to help debugging
	String() string
}

// Channel Response is an independent thread that starts and stops with the channel, and continuously monitors
// whether server has replied with a message event or not
type TGChannelResponse interface {
	// Await waits (loops) till the channel response receives reply message from the server
	Await(tester StatusTester)
	// GetCallback gets a Callback object
	GetCallback() Callback
	// GetReply gets Reply object
	GetReply() TGMessage
	// GetRequestId gets Request id
	GetRequestId() int64
	// GetStatus gets Status
	GetStatus() ChannelResponseStatus
	// IsBlocking checks whether this channel response is blocking or not
	IsBlocking() bool
	// Reset resets the state of channel response and initializes everything
	Reset()
	// SetReply sets the reply message received from the server
	SetReply(msg TGMessage)
	// SetRequestId sets Request id
	SetRequestId(requestId int64)
	// Signal lets other listeners of channel response know the status of this channel response
	Signal(status ChannelResponseStatus)
}

type Callback interface {
	OnResponse(msg TGMessage)
}

type StatusTester interface {
	// Test checks whether the channel response is in WAIT mode or not
	Test(Status ChannelResponseStatus) bool
}

// ======= Channel Response Status =======
type ChannelResponseStatus int

const (
	Waiting ChannelResponseStatus = iota
	Ok
	Pushed
	Resend
	Disconnected
	Closed
)

type TGTracer interface {
	// Start starts the channel tracer
	Start()
	// Stop stops the channel tracer
	Stop()
	// Trace traces the path the message has taken
	Trace(msg TGMessage)
}

// DataCryptoGrapher is an event listener that gets triggered when that event occurs
type TGDataCryptoGrapher interface {
	// Decrypt decrypts the buffer
	//Decrypt(encBuffer []byte) ([]byte, TGError)
	Decrypt(is TGInputStream) ([]byte, TGError)
	// Encrypt encrypts the buffer
	Encrypt(decBuffer []byte) ([]byte, TGError)
}

type LinkEventHandler interface {
	OnException(exception TGError, duringClose bool)
	OnReconnect() bool
	GetTerminatedText() string
}


func (proType TGProtocol) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if proType&ProtocolTCP == ProtocolTCP {
		buffer.WriteString("tcp")
	} else if proType&ProtocolSSL == ProtocolSSL {
		buffer.WriteString("ssl")
	} else if proType&ProtocolHTTP == ProtocolHTTP {
		buffer.WriteString("http")
	} else if proType&ProtocolHTTPS == ProtocolHTTPS {
		buffer.WriteString("https")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}



func (linkState LinkState) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if linkState&LinkNotConnected == LinkNotConnected {
		buffer.WriteString("Link Not Connected")
	} else if linkState&LinkConnected == LinkConnected {
		buffer.WriteString("Link Connected")
	} else if linkState&LinkClosing == LinkClosing {
		buffer.WriteString("Link Closing")
	} else if linkState&LinkClosed == LinkClosed {
		buffer.WriteString("Link Closed")
	} else if linkState&LinkFailedOnSend == LinkFailedOnSend {
		buffer.WriteString("Link Failed On Send")
	} else if linkState&LinkFailedOnRecv == LinkFailedOnRecv {
		buffer.WriteString("Link Failed On Receiving")
	} else if linkState&LinkFailedOnProcessing == LinkFailedOnProcessing {
		buffer.WriteString("Link Failed On Processing")
	} else if linkState&LinkReconnecting == LinkReconnecting {
		buffer.WriteString("Link Reconnecting")
	} else if linkState&LinkTerminated == LinkTerminated {
		buffer.WriteString("Link Terminated")
	}
	return buffer.String()
}

type ResendMode int

// ======= Resend Mode Types =======
const (
	ModeDontReconnectAndIgnore ResendMode = 1 << iota
	ModeReconnectAndResend
	ModeReconnectAndRaiseException
	ModeReconnectAndIgnore
)

func (resendMode ResendMode) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if resendMode&ModeDontReconnectAndIgnore == ModeDontReconnectAndIgnore {
		buffer.WriteString("DontReconnectAndIgnore")
	} else if resendMode&ModeReconnectAndResend == ModeReconnectAndResend {
		buffer.WriteString("ReconnectAndResend")
	} else if resendMode&ModeReconnectAndRaiseException == ModeReconnectAndRaiseException {
		buffer.WriteString("ReconnectAndRaiseException")
	} else if resendMode&ModeReconnectAndIgnore == ModeReconnectAndIgnore {
		buffer.WriteString("ReconnectAndIgnore")
	}
	return buffer.String()
}

// ======= Reconnect State Types =======
type ReconnectState int

const (
	ReconnectStateChannelClosed ReconnectState = 1 << iota
	ReconnectStateSuccess
	ReconnectStateFailed
	ReconnectStateFailedAllAttempts
)

func (reconnectState ReconnectState) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if reconnectState&ReconnectStateChannelClosed == ReconnectStateChannelClosed {
		buffer.WriteString("Reconnect Channel Closed")
	} else if reconnectState&ReconnectStateSuccess == ReconnectStateSuccess {
		buffer.WriteString("Reconnect Success")
	} else if reconnectState&ReconnectStateFailed == ReconnectStateFailed {
		buffer.WriteString("Reconnect Failed")
	} else if reconnectState&ReconnectStateFailedAllAttempts == ReconnectStateFailedAllAttempts {
		buffer.WriteString("Reconnect Failed All Attempts")
	}
	return buffer.String()
}

// ======= Channel Status Types =======
type ChannelStatus int

const (
	ChannelStatusWaiting ChannelStatus = 1 << iota
	ChannelStatusOk
	ChannelStatusPushed
	ChannelStatusResend
	ChannelStatusDisconnected
	ChannelStatusClosed
)

func (channelStatus ChannelStatus) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if channelStatus&ChannelStatusWaiting == ChannelStatusWaiting {
		buffer.WriteString("Channel Status Waiting")
	} else if channelStatus&ChannelStatusOk == ChannelStatusOk {
		buffer.WriteString("Channel Status Ok")
	} else if channelStatus&ChannelStatusPushed == ChannelStatusPushed {
		buffer.WriteString("Channel Status Pushed")
	} else if channelStatus&ChannelStatusResend == ChannelStatusResend {
		buffer.WriteString("Channel Status Resend")
	} else if channelStatus&ChannelStatusDisconnected == ChannelStatusDisconnected {
		buffer.WriteString("Channel Status Disconnected")
	} else if channelStatus&ChannelStatusClosed == ChannelStatusClosed {
		buffer.WriteString("Channel Status Closed")
	}
	return buffer.String()
}

func (responseStatus ChannelResponseStatus) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if responseStatus&Waiting == Waiting {
		buffer.WriteString("Channel Response Status is Waiting")
	} else if responseStatus&Ok == Ok {
		buffer.WriteString("Channel Response Status is Ok")
	} else if responseStatus&Pushed == Pushed {
		buffer.WriteString("Channel Response Status is Pushed")
	} else if responseStatus&Resend == Resend {
		buffer.WriteString("Channel Response Status is Resend")
	} else if responseStatus&Disconnected == Disconnected {
		buffer.WriteString("Channel Response Status is Disconnected")
	} else if responseStatus&Closed == Closed {
		buffer.WriteString("Channel Response Status is Closed")
	}
	return buffer.String()
}

