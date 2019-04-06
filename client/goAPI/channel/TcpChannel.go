package channel

import (
	"bytes"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/pdu"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"net"
	"strings"
	"sync"
	"time"
)

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
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
 * File name: TcpChannel.go
 * Created on: Dec 01, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	dataBufferSize = 32 * 1024 // 32 KB
)

type TCPChannel struct {
	*AbstractChannel
	shutdownLock   sync.RWMutex // rw-lock for synchronizing read-n-update of env configuration
	isSocketClosed bool         // indicate if the connection is already closed
	msgCh          chan types.TGMessage
	socket         *net.TCPConn
	input          *iostream.ProtocolDataInputStream
	output         *iostream.ProtocolDataOutputStream
}

func DefaultTCPChannel() *TCPChannel {
	newChannel := TCPChannel{
		AbstractChannel: DefaultAbstractChannel(),
		msgCh:           make(chan types.TGMessage),
		isSocketClosed:  false,
	}
	buff := make([]byte, 0)
	newChannel.input = iostream.NewProtocolDataInputStream(buff)
	newChannel.output = iostream.NewProtocolDataOutputStream(0)
	//newChannel.reconnecting = false
	newChannel.exceptionCond = sync.NewCond(&newChannel.exceptionLock) // Condition for lock
	newChannel.reader = NewChannelReader(&newChannel)
	return &newChannel
}

func NewTCPChannel(linkUrl *LinkUrl, props *utils.SortedProperties) *TCPChannel {
	//logger.Log(fmt.Sprintf("======> Entering TCPChannel:NewTCPChannel w/ linkUrl: '%s'", linkUrl.String()))
	newChannel := DefaultTCPChannel()
	newChannel.channelUrl = linkUrl
	newChannel.primaryUrl = linkUrl
	newChannel.channelProperties = props
	//logger.Log(fmt.Sprintf("======> Entering TCPChannel:NewTCPChannel w/ newChannel: '%+v'", newChannel.String()))
	return newChannel
}

/////////////////////////////////////////////////////////////////
// Private functions for TCPChannel
/////////////////////////////////////////////////////////////////

func (obj *TCPChannel) doAuthenticate() types.TGError {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:doAuthenticate"))
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := pdu.CreateMessageForVerb(pdu.VerbAuthenticateRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: TCPChannel::doAuthenticate pdu.CreateMessageForVerb(pdu.VerbAuthenticateRequest) failed w/ '%+v'", err.Error()))
		return err
	}

	msgRequest.(*pdu.AuthenticateRequestMessage).SetClientId(obj.clientId)
	msgRequest.(*pdu.AuthenticateRequestMessage).SetInboxAddr(obj.inboxAddress)
	msgRequest.(*pdu.AuthenticateRequestMessage).SetUserName(obj.getChannelUserName())
	msgRequest.(*pdu.AuthenticateRequestMessage).SetPassword(obj.getChannelPassword())

	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:doAuthenticate about to request reply for request '%+v'", msgRequest.String()))
	msgResponse, err := channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::doAuthenticate channelRequestReply failed w/ '%+v'", err.Error()))
		return err
	}
	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:doAuthenticate received reply as '%+v'", msgResponse.String()))
	if ! msgResponse.(*pdu.AuthenticateResponseMessage).IsSuccess() {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::doAuthenticate msgResponse.(*pdu.AuthenticateResponseMessage).IsSuccess() failed"))
		return exception.NewTGBadAuthenticationWithRealm(types.INTERNAL_SERVER_ERROR, types.TGErrorBadAuthentication, "Bad username/password combination", "", "tgdb")
	}

	obj.setChannelAuthToken(msgResponse.GetAuthToken())
	obj.setChannelSessionId(msgResponse.GetSessionId())

	cryptoDataGrapher, err := NewDataCryptoGrapher(msgResponse.GetSessionId(), msgResponse.(*pdu.AuthenticateResponseMessage).GetServerCertBuffer())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::doAuthenticate NewDataCryptoGrapher failed w/ '%+v'", err.Error()))
		return err
	}
	obj.setDataCryptoGrapher(cryptoDataGrapher)
	logger.Log(fmt.Sprintf("======> Returning TCPChannel:doAuthenticate Successfully authenticated for user: '%s'", obj.getChannelUserName()))
	return nil
}

func (obj *TCPChannel) performHandshake(sslMode bool) types.TGError {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:performHandshake"))
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := pdu.CreateMessageForVerb(pdu.VerbHandShakeRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::performHandshake pdu.CreateMessageForVerb(pdu.VerbAuthenticateRequest) failed w/ '%+v'", err.Error()))
		return err
	}

	msgRequest.(*pdu.HandShakeRequestMessage).SetRequestType(pdu.InitiateRequest)

	logger.Log(fmt.Sprintf("======> Inside TCPChannel:performHandshake about to request reply for InitiateRequest '%+v'", msgRequest.String()))
	msgResponse, err := channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::doAuthenticate channelRequestReply failed w/ '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("======> Inside TCPChannel:performHandshake received reply as '%+v'", msgResponse.String()))
	if msgResponse.GetVerbId() != pdu.VerbHandShakeResponse {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::performHandshake HandshakeResponse message response NOT received"))
		if msgResponse.GetVerbId() == pdu.VerbSessionForcefullyTerminated {
			errMsg := msgResponse.(*pdu.SessionForcefullyTerminatedMessage).GetKillString()
			return exception.NewTGChannelDisconnectedWithMsg(errMsg)
		}
		errMsg := fmt.Sprintf("Expecting a HandshakeResponse message, and received: '%d'. Cannot connect to the server at: '%s'", msgResponse.GetVerbId(), obj.channelUrl.GetUrlAsString())
		return exception.NewTGGeneralException(types.TGDB_HNDSHKRESP_ERROR, types.TGErrorGeneralException, errMsg, "")
	}

	response := msgResponse.(*pdu.HandShakeResponseMessage)
	if response.GetResponseStatus() != pdu.ResponseAcceptChallenge {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::performHandshake response.GetResponseStatus() is NOT pdu.ResponseAcceptChallenge"))
		errMsg := fmt.Sprintf("'%s': Handshake Failed. Cannot connect to the server at: '%s'", types.TGDB_HNDSHKRESP_ERROR, obj.channelUrl.GetUrlAsString())
		return exception.NewTGGeneralException(types.TGDB_HNDSHKRESP_ERROR, types.TGErrorGeneralException, errMsg, "")
	}

	// Validate the version specific information on the response object
	serverVersion := response.GetChallenge()
	clientVersion := utils.GetClientVersion()
	err = obj.validateHandshakeResponseVersion(serverVersion, clientVersion)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::performHandshake validateHandshakeResponseVersion failed w/ '%+v'", err.Error()))
		return err
	}

	challenge := clientVersion.GetVersionAsLong()

	// Ignore Error Handling
	_ = msgRequest.(*pdu.HandShakeRequestMessage).UpdateSequenceAndTimeStamp(-1)
	msgRequest.(*pdu.HandShakeRequestMessage).SetRequestType(pdu.ChallengeAccepted)
	msgRequest.(*pdu.HandShakeRequestMessage).SetSslMode(sslMode)
	msgRequest.(*pdu.HandShakeRequestMessage).SetChallenge(challenge)

	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:performHandshake about to request reply for ChallengeAccepted '%+v'", msgRequest.String()))
	msgResponse, err = channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::performHandshake channelRequestReply failed w/ '%+v'", err.Error()))
		return err
	}
	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:performHandshake received reply (2) as '%+v'", msgResponse.String()))
	if msgResponse.GetVerbId() != pdu.VerbHandShakeResponse {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::performHandshake HandshakeResponse message response NOT received"))
		if msgResponse.GetVerbId() == pdu.VerbSessionForcefullyTerminated {
			errMsg := msgResponse.(*pdu.SessionForcefullyTerminatedMessage).GetKillString()
			return exception.NewTGChannelDisconnectedWithMsg(errMsg)
		}
		errMsg := fmt.Sprintf("Expecting a HandshakeResponse message, and received: '%d'. Cannot connect to the server at: '%s'", msgResponse.GetVerbId(), obj.channelUrl.GetUrlAsString())
		return exception.NewTGGeneralException(types.TGDB_HNDSHKRESP_ERROR, types.TGErrorGeneralException, errMsg, "")
	}

	response = msgResponse.(*pdu.HandShakeResponseMessage)
	if response.GetResponseStatus() != pdu.ResponseProceedWithAuthentication {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::performHandshake response.GetResponseStatus() is NOT pdu.ResponseAcceptChallenge"))
		errMsg := fmt.Sprintf("'%s': Handshake Failed. Cannot connect to the server at: '%s'", types.TGDB_HNDSHKRESP_ERROR, obj.channelUrl.GetUrlAsString())
		return exception.NewTGGeneralException(types.TGDB_HNDSHKRESP_ERROR, types.TGErrorGeneralException, errMsg, "")
	}
	logger.Log(fmt.Sprintf("======> Returning TCPChannel::performHandshake Handshake w/ Remote Server is successful."))
	return nil
}

func (obj *TCPChannel) setSocket(newSocket *net.TCPConn) types.TGError {
	obj.socket = newSocket
	err := obj.socket.SetNoDelay(true)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel:setSocket Failed to set NoDelay flag to true w/ '%+v'", err.Error()))
		failureMessage := fmt.Sprint("Failed to set NoDelay flag to true")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, err.Error())
	}

	err = obj.socket.SetLinger(0) // <= 0 means Do not linger
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel:setSocket Failed to set NoLinger flag to true w/ '%+v'", err.Error()))
		failureMessage := fmt.Sprint("Failed to set NoLinger flag to true")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, err.Error())
	}

	buff := make([]byte, dataBufferSize)
	obj.input = iostream.NewProtocolDataInputStream(buff)
	obj.input.BufLen = 0
	obj.output = iostream.NewProtocolDataOutputStream(dataBufferSize)
	//clientId = properties.get(ConfigName.ChannelClientId.getName());
	//if (clientId == null) {
	//	clientId = properties.get(ConfigName.ChannelClientId.getAlias());
	//	if (clientId == null) {
	//		clientId = TGEnvironment.getInstance().getChannelClientId();
	//	}
	//}
	clientId := obj.GetProperties().GetProperty(utils.GetConfigFromKey(utils.ChannelClientId), "")
	obj.setChannelClientId(clientId)
	obj.setChannelInboxAddr(obj.socket.RemoteAddr().String()) //SS:TODO: Is this correct
	return nil
}

func (obj *TCPChannel) setBuffers(newSocket *net.TCPConn) types.TGError {
	sendSize := obj.channelProperties.GetPropertyAsInt(utils.GetConfigFromKey(utils.ChannelSendSize))
	if sendSize > 0 {
		err := newSocket.SetWriteBuffer(sendSize*1024)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::setBuffers newSocket.SetWriteBuffer failed w/ '%+v'", err.Error()))
			errMsg := fmt.Sprintf("TCPChannel:setBuffers unable to set write buffer limit to '%d'", sendSize*1024)
			return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
		}
	}
	receiveSize := obj.channelProperties.GetPropertyAsInt(utils.GetConfigFromKey(utils.ChannelRecvSize))
	if receiveSize > 0 {
		err := newSocket.SetReadBuffer(receiveSize*1024)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::setBuffers SetReadBuffer failed w/ '%+v'", err.Error()))
			errMsg := fmt.Sprintf("TCPChannel:setBuffers unable to set read buffer limit to '%d'", receiveSize*1024)
			return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
		}
	}
	return nil
}

func (obj *TCPChannel) tryRead() (types.TGMessage, types.TGError) {
	//logger.Log(fmt.Sprintf("======> Entering TCPChannel:tryRead"))
	n, err := obj.input.Available()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::tryRead obj.input.Available() failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::tryRead there is no data available to be read"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	if n <= 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning TCPChannel::tryRead as there are no bytes to read from the wire"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("======> Inside TCPChannel:tryRead about to request message '%d' bytes from the wire", n))
	return obj.ReadWireMsg()
}

func (obj *TCPChannel) validateHandshakeResponseVersion(sVersion int64, cVersion *utils.TGClientVersion) types.TGError {
	serverVersion := utils.NewTGServerVersion(sVersion)
	sStrVer := serverVersion.GetVersionString()

	cStrVer := cVersion.GetVersionString()

	if 	serverVersion.GetMajor() == cVersion.GetMajor() &&
		serverVersion.GetMinor() == cVersion.GetMinor() &&
		serverVersion.GetUpdate() == cVersion.GetUpdate() {
		return nil
	}

	errMsg := fmt.Sprintf("======> Inside SSLChannel:validateHandshakeResponseVersion - Version mismatch between client(%s) & server(%s)", cStrVer, sStrVer)
	logger.Log(errMsg)
	return exception.GetErrorByType(types.TGErrorVersionMismatchException, "", errMsg, "")
}

func (obj *TCPChannel) writeLoop(done chan bool) {
	logger.Log(fmt.Sprint("======> Entering TCPChannel:writeLoop"))
	for {
		logger.Log(fmt.Sprintf("======> Inside TCPChannel:writeLoop entering infinite loop"))
		select { // Non-blocking channel operation
		case msg, ok := <-obj.msgCh: // Retrieve the message from the channel
			if !ok {
				//if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
				//	logMessage("TCPChannel::writeLoop unable to retrieve msg from the channel);
				//}
				logger.Error(fmt.Sprint("ERROR: Returning TCPChannel:writeLoop unable to retrieve message from obj.msgCh"))
				return
			}
			logger.Log(fmt.Sprintf("======> Inside TCPChannel:writeLoop retrieved message from obj.msgCh as '%+v'", msg.String()))

			err := obj.writeToWire(msg)
			if err != nil {
				// TODO: Revisit later - Do something
				logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::writeLoop unable to obj.writeToWire w/ '%+v'", err.Error()))
				return
			}

			logger.Log(fmt.Sprintf("======> Inside TCPChannel:writeLoop successfully wrote message '%+v' on the socket", msg.String()))
			break
		default:
			// TODO: Revisit later - Do something
		}
	}	// End of Infinite Loop
	// Send an acknowledgement of completion to the parent thread
	done <- true
	logger.Log(fmt.Sprint("======> Returning TCPChannel:writeLoop"))
}

func (obj *TCPChannel) writeToWire(msg types.TGMessage) types.TGError {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:writeToWire w/ Msg: '%+v'", msg.String()))
	obj.DisablePing()
	msgBytes, bufLen, err := msg.ToBytes()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::writeToWire unable to convert message into byte format w/ '%+v'", err.Error()))
		errMsg := fmt.Sprintf("TCPChannel::writeToWire unable to convert message into byte format")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.GetErrorMsg())
	}

	// Clear timeout deadlines set at the time of creation of the socket
	sErr := obj.socket.SetDeadline(time.Time{})
	if sErr != nil {
		logger.Error(fmt.Sprintf("Returning TCPChannel::writeToWire unable to clear the deadline over TCP socket w/ '%+v'", sErr.Error()))
		errMsg := fmt.Sprintf("TCPChannel::writeToWire unable to clear the deadline over TCP socket")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}

	// Reset timeout deadlines starting from NOW!!!
	timeout := utils.NewTGEnvironment().GetChannelConnectTimeout()
	sErr = obj.socket.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if sErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::writeToWire unable to reset the deadline over TCP socket w/ '%+v'", sErr.Error()))
		errMsg := fmt.Sprintf("TCPChannel::writeToWire unable to reset the deadline over TCP socket")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}

	logger.Log(fmt.Sprintf("======> Inside TCPChannel:writeToWire about to write message bytes on the socket as '%+v'", msgBytes[0:bufLen]))
	// Put the data packet on the socket for network transmission
	_, sErr = obj.socket.Write(msgBytes[0:bufLen])
	if sErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::writeToWire unable to send message bytes over TCP socket w/ '%+v'", sErr.Error()))
		errMsg := fmt.Sprintf("TCPChannel::writeToWire unable to send message bytes over TCP socket")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}
	logger.Log(fmt.Sprintf("======> Returning TCPChannel:writeToWire successfully wrote message bytes on the socket as '%+v'", msgBytes[0:bufLen]))
	return nil
}

func intToBytes(value int, bytes []byte, offset int) {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:intToBytes w/ value as '%d', byteArray as '%+v' and offset '%d'", value, bytes, offset))
	for i := 0; i < 4; i++ {
		bytes[offset+i] = byte((value >> uint(8*(3-i))) & 0xff)
	}
	logger.Log(fmt.Sprintf("======> Returning TCPChannel:intToBytes w/ byteArray as '%+v'", bytes))
}

/////////////////////////////////////////////////////////////////
// Helper functions for TCPChannel
/////////////////////////////////////////////////////////////////

func (obj *TCPChannel) GetIsClosed() bool {
	return obj.isSocketClosed
}

func (obj *TCPChannel) SetIsClosed(flag bool) {
	obj.isSocketClosed = flag
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChannel
/////////////////////////////////////////////////////////////////

// ChannelLock locks the communication channel between TGDB client and server
func (obj *TCPChannel) ChannelLock() {
	obj.sendLock.Lock()
}

// ChannelUnlock unlocks the communication channel between TGDB client and server
func (obj *TCPChannel) ChannelUnlock() {
	obj.sendLock.Unlock()
}

// Connect connects the underlying channel using the URL end point
func (obj *TCPChannel) Connect() types.TGError {
	return channelConnect(obj)
}

// DisablePing disables the pinging ability to the channel
func (obj *TCPChannel) DisablePing() {
	obj.needsPing = false
}

// Disconnect disconnects the channel from its URL end point
func (obj *TCPChannel) Disconnect() types.TGError {
	return channelDisConnect(obj)
}

// EnablePing enables the pinging ability to the channel
func (obj *TCPChannel) EnablePing() {
	obj.needsPing = true
}

// ExceptionLock locks the communication channel between TGDB client and server in case of business exceptions
func (obj *TCPChannel) ExceptionLock() {
	obj.exceptionLock.Lock()
}

// ExceptionUnlock unlocks the communication channel between TGDB client and server in case of business exceptions
func (obj *TCPChannel) ExceptionUnlock() {
	obj.exceptionLock.Unlock()
}

// GetAuthToken gets Authorization Token
func (obj *TCPChannel) GetAuthToken() int64 {
	return obj.authToken
}

// GetClientId gets Client Name
func (obj *TCPChannel) GetClientId() string {
	return obj.clientId
}

// GetChannelURL gets the Channel URL
func (obj *TCPChannel) GetChannelURL() types.TGChannelUrl {
	return obj.channelUrl
}

// GetConnectionIndex gets the Connection Index
func (obj *TCPChannel) GetConnectionIndex() int {
	return obj.connectionIndex
}

// GetExceptionCondition gets the Exception Condition
func (obj *TCPChannel) GetExceptionCondition() *sync.Cond {
	return obj.exceptionCond
}

// GetLinkState gets the Link/Channel State
func (obj *TCPChannel) GetLinkState() types.LinkState {
	return obj.channelLinkState
}

// GetNoOfConnections gets number of connections this channel has
func (obj *TCPChannel) GetNoOfConnections() int32 {
	return obj.numOfConnections
}

// GetPrimaryURL gets the Primary URL
func (obj *TCPChannel) GetPrimaryURL() types.TGChannelUrl {
	return obj.primaryUrl
}

// GetProperties gets the Channel Properties
func (obj *TCPChannel) GetProperties() types.TGProperties {
	return obj.channelProperties
}

// GetReader gets the Channel Reader
func (obj *TCPChannel) GetReader() types.TGChannelReader {
	return obj.reader
}

// GetResponses gets the Channel Response Map
func (obj *TCPChannel) GetResponses() map[int64]types.TGChannelResponse {
	return obj.responses
}

// GetSessionId gets Session id
func (obj *TCPChannel) GetSessionId() int64 {
	return obj.sessionId
}

// IsChannelPingable checks whether the channel is pingable or not
func (obj *TCPChannel) IsChannelPingable() bool {
	return obj.needsPing
}

// IsClosed checks whether channel is open or closed
func (obj *TCPChannel) IsClosed() bool {
	return isChannelClosed(obj)
}

// SendMessage sends a Message on this channel, and returns immediately - An Asynchronous or Non-Blocking operation
func (obj *TCPChannel) SendMessage(msg types.TGMessage) types.TGError {
	return channelSendMessage(obj, msg, true)
}

// SendRequest sends a Message, waits for a response in the message format, and blocks the thread till it gets the response
func (obj *TCPChannel) SendRequest(msg types.TGMessage, response types.TGChannelResponse) (types.TGMessage, types.TGError) {
	return channelSendRequest(obj, msg, response, true)
}

// SetChannelLinkState sets the Link/Channel State
func (obj *TCPChannel) SetChannelLinkState(state types.LinkState) {
	obj.channelLinkState = state
}

// SetChannelURL sets the channel URL
func (obj *TCPChannel) SetChannelURL(url types.TGChannelUrl) {
	obj.channelUrl = url.(*LinkUrl)
}

// SetConnectionIndex sets the connection index
func (obj *TCPChannel) SetConnectionIndex(index int) {
	obj.connectionIndex = index
}

// SetNoOfConnections sets number of connections
func (obj *TCPChannel) SetNoOfConnections(count int32) {
	obj.numOfConnections = count
}

// SetResponse sets the ChannelResponse Map
func (obj *TCPChannel) SetResponse(reqId int64, response types.TGChannelResponse) {
	obj.responses[reqId] = response
}

// Start starts the channel so that it can send and receive messages
func (obj *TCPChannel) Start() types.TGError {
	return channelStart(obj)
}

// Stop stops the channel forcefully or gracefully
func (obj *TCPChannel) Stop(bForcefully bool) {
	channelStop(obj, bForcefully)
}

// CreateSocket creates the physical link socket
func (obj *TCPChannel) CreateSocket() types.TGError {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:CreateSocket"))
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	obj.SetChannelLinkState(types.LinkNotConnected)
	host := obj.channelUrl.urlHost
	port := obj.channelUrl.urlPort
	serverAddr := fmt.Sprint(host, ":", port)
	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:CreateSocket attempting to resolve address for '%s'", serverAddr))

	tcpAddr, tErr := net.ResolveTCPAddr(types.ProtocolTCP.String(), serverAddr)
	if tErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket net.ResolveTCPAddr failed w/ '%+v'", tErr.Error()))
		errMsg := fmt.Sprintf("TCPChannel:CreateSocket unable to resolve channel address '%s'", serverAddr)
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, tErr.Error())
	}
	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:CreateSocket resolved TCP address for '%s' as '%+v'", serverAddr, tcpAddr))

	tcpConn, cErr := net.DialTCP(types.ProtocolTCP.String(), nil, tcpAddr)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to connect to the server at '%s' w/ '%+v'", serverAddr, cErr.Error()))
		failureMessage := fmt.Sprintf("Failed to connect to the server at '%s'" + serverAddr)
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
	}
	logger.Log(fmt.Sprintf("======> Inside TCPChannel:CreateSocket created TCP connection for '%s' as '%+v'", serverAddr, tcpConn))

	timeout := utils.NewTGEnvironment().GetChannelConnectTimeout()
	dErr := tcpConn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if dErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to set deadline of '%+v' seconds on the connection to the server w/ '%+v'", time.Duration(timeout) * time.Second, dErr.Error()))
		failureMessage := fmt.Sprintf("Failed to set the timeout '%d' on socket", timeout)
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, dErr.Error())
	}

	err := tcpConn.SetKeepAlive(true)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to set keep alive flag to true w/ '%+v'", err.Error()))
		failureMessage := fmt.Sprint("Failed to set keep alive flag to true")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
	}

	// Set Read / Write Buffer Size on the socket
	tcErr := obj.setBuffers(tcpConn)
	if tcErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to set buffers w/ '%+v'", tcErr.Error()))
		return tcErr
	}
	tcErr = obj.setSocket(tcpConn)
	if tcErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to set socket value to the object w/ '%+v'", tcErr.Error()))
		return tcErr
	}
	obj.SetIsClosed(false)
	logger.Log(fmt.Sprintf("======> Returning TCPChannel:CreateSocket w/ TCP Connection as '%+v'", *obj.socket))
	return nil
}

// CloseSocket closes the socket
func (obj *TCPChannel) CloseSocket() types.TGError {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:CloseSocket w/ socket: '%+v'", obj.socket))
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	if obj.socket != nil {
		cErr := obj.socket.Close()
		if cErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CloseSocket obj.socket.Close() failed w/ '%+v'", cErr.Error()))
			failureMessage := "Failed to close the socket to the server"
			return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
		}
	}
	obj.SetIsClosed(true)
	obj.socket = nil
	logger.Log(fmt.Sprintf("======> Returning TCPChannel:CloseSocket for socket: '%+v'", obj.socket))
	return nil
}

// OnConnect executes all the channel specific activities
func (obj *TCPChannel) OnConnect() types.TGError {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:OnConnect about to tryRead w/ socket: '%+v'", obj.socket))
	msg, err := obj.tryRead()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::OnConnect obj.tryRead() failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::OnConnect there is no data available to be read"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}
	if msg != nil {
		logger.Log(fmt.Sprintf("======> Inside TCPChannel:OnConnect tryRead() read Message as '%+v'", msg.String()))
	}

	if msg != nil && msg.GetVerbId() == pdu.VerbSessionForcefullyTerminated {
		logger.Warning(fmt.Sprint("WARNING: Returning TCPChannel:OnConnect since Message is of Forceful Termination Type"))
		return exception.NewTGChannelDisconnectedWithMsg(msg.(*pdu.SessionForcefullyTerminatedMessage).GetKillString())
	}

	logger.Log(fmt.Sprintf("======> Inside TCPChannel:OnConnect about to performHandshake"))
	err = obj.performHandshake(false)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::OnConnect obj.performHandshake() failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::OnConnect error in performing handshake with server"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	logger.Log(fmt.Sprintf("======> Inside TCPChannel:OnConnect about to doAuthenticate"))
	err = obj.doAuthenticate()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::OnConnect obj.doAuthenticate() failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::OnConnect error in authentication with server"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	logger.Log(fmt.Sprintf("======> Returning TCPChannel:OnConnect w/ socket: '%+v'", obj.socket))
	return nil
}

// ReadWireMsg reads the message from the wire in the form of byte stream
func (obj *TCPChannel) ReadWireMsg() (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:ReadWireMsg  w/ socket: '%+v'", obj.socket))
	obj.input.BufLen = dataBufferSize
	in := obj.input
	if in == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning TCPChannel:ReadWireMsg since obj.input is NIL"))
		// TODO: Revisit later - Should we not return an error?
		return nil, nil
	}

	obj.DisablePing()
	if obj.GetIsClosed() {
		logger.Warning(fmt.Sprint("WARNING: Returning TCPChannel:ReadWireMsg since TCP Channel is Closed"))
		// TODO: Revisit later - Should we not return an error?
		return nil, nil
	}
	obj.lastActiveTime = time.Now()

	// Read the data on the socket
	buff := make([]byte, dataBufferSize)
	n, sErr := obj.socket.Read(buff)
	if sErr != nil || n <= 0 {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::ReadWireMsg obj.socket.Read failed w/ '%+v'", sErr.Error()))
		errMsg := "TCPChannel::ReadWireMsg obj.socket.Read failed"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, sErr.Error())
	}
	logger.Log(fmt.Sprintf("======> Inside TCPChannel:ReadWireMsg Read '%d' bytes from the wire in buff '%+v'", n, buff[:(2*n)]))
	copy(in.Buf, buff[:n])
	in.BufLen = n
	logger.Log(fmt.Sprintf("======> Inside TCPChannel:ReadWireMsg Input Stream Buffer('%d') is '%+v'", in.BufLen, in.Buf[:(2*n)]))

	// Needed to avoid dirty data in the buffer when we handle the message
	buffer := make([]byte, n)
	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:ReadWireMsg in.ReadFullyAtPos read msgBytes as '%+v'", msgBytes))
	copy(buffer, buff[:n])
	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:ReadWireMsg copied into buffer as '%+v'", buffer))

	//intToBytes(size, msgBytes, 0)
	//bytesRead, _ := utils.FormatHex(msgBytes)
	//logger.Log(fmt.Sprintf("======> Inside TCPChannel:ReadWireMsg bytes read: '%s'", bytesRead))

	msg, err := pdu.CreateMessageFromBuffer(buffer, 0, n)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::ReadWireMsg pdu.CreateMessageFromBuffer failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::ReadWireMsg unable to create a message from the input stream bytes"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.GetErrorMsg())
	}
	logger.Log(fmt.Sprintf("======> Inside TCPChannel:ReadWireMsg Created Message from buffer as '%+v'", msg.String()))

	if msg.GetVerbId() == pdu.VerbExceptionMessage {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::ReadWireMsg msg.GetVerbId() == pdu.VerbExceptionMessage"))
		errMsg := msg.(*pdu.ExceptionMessage).GetExceptionMsg()
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	if msg.GetVerbId() == pdu.VerbHandShakeResponse {
		if msg.(*pdu.HandShakeResponseMessage).GetResponseStatus() == pdu.ResponseChallengeFailed {
			errMsg := msg.(*pdu.HandShakeResponseMessage).GetErrorMessage()
			logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::ReadWireMsg msg.GetVerbId() == pdu.VerbHandShakeResponse w/ '%+v'", errMsg))
			return nil, exception.GetErrorByType(types.TGErrorVersionMismatchException, "", errMsg, "")
		}
	}
	logger.Log(fmt.Sprintf("======> Returning TCPChannel:ReadWireMsg w/ Socket '%+v' and Message as '%+v'", obj.socket, msg.String()))
	return msg, nil
}

// Send sends the message to the server, compress and or encrypt.
// Hence it is abstraction, that the Channel knows about it.
// @param msg       The message that needs to be sent to the server
func (obj *TCPChannel) Send(msg types.TGMessage) types.TGError {
	logger.Log(fmt.Sprintf("======> Entering TCPChannel:Send w/ Socket '%+v' and Message as '%+v'", obj.socket, msg.String()))
	if obj.output == nil || obj.GetIsClosed() {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::Send as the channel is closed"))
		errMsg := fmt.Sprintf("TCPChannel:Send - unable to send message to server as the channel is closed")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, "")
	}

	err := obj.writeToWire(msg)
	logger.Log(fmt.Sprintf("======> Returning TCPChannel:Send w/ error '%+v'", err))
	return err

	//// Wait for success notification from the GO routine
	//done := make(chan bool, 1)
	//// Push the message on to the channel (FIFO pipe)
	//obj.msgCh <- msg
	//
	//// TODO: Revisit later - for performance and optimization
	//// Execute sending each message content in another thread/go-routine
	////go obj.writeLoop(done)
	//// This is a common function called by both SendMessage (non-blocking) and SendRequest (blocking)
	//// Hence this should be handled called in another thread/go-routing in SendMessage ONLY
	//obj.writeLoop(done)
	//<-done
	//logger.Log(fmt.Sprintf("======> Exiting TCPChannel:Send w/ Message as '%+v'", msg.String()))
	//return nil
}

func (obj *TCPChannel) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TCPChannel:{")
	buffer.WriteString(fmt.Sprintf("IsSocketClosed: %+v", obj.isSocketClosed))
	buffer.WriteString(fmt.Sprintf(", MsgCh: %+v", obj.msgCh))
	buffer.WriteString(fmt.Sprintf(", Socket: %+v", obj.socket))
	//buffer.WriteString(fmt.Sprintf(", Input: %+v", obj.input.String()))
	//buffer.WriteString(fmt.Sprintf(", Output: %+v", obj.output.String()))
	strArray := []string{buffer.String(), obj.channelToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}
