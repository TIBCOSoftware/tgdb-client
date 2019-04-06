package channel

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/pdu"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"io/ioutil"
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

type SSLChannel struct {
	*AbstractChannel
	shutdownLock   sync.RWMutex // rw-lock for synchronizing read-n-update of env configuration
	isSocketClosed bool         // indicate if the connection is already closed
	msgCh          chan types.TGMessage
	socket         *tls.Conn
	tlsConfig      *tls.Config
	input          *iostream.ProtocolDataInputStream
	output         *iostream.ProtocolDataOutputStream
}

func DefaultSSLChannel() *SSLChannel {
	newChannel := SSLChannel{
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

func NewSSLChannel(linkUrl *LinkUrl, props *utils.SortedProperties) (*SSLChannel, types.TGError) {
	//logger.Log(fmt.Sprintf("======> Entering SSLChannel:NewSSLChannel w/ linkUrl: '%s'", linkUrl.String()))
	newChannel := DefaultSSLChannel()
	newChannel.channelUrl = linkUrl
	newChannel.primaryUrl = linkUrl
	newChannel.channelProperties = props
	config, err := initTLSConfig(props)
	if err != nil {
		return nil, err
	}
	newChannel.tlsConfig = config
	logger.Log(fmt.Sprintf("======> Returning SSLChannel:NewSSLChannel w/ TLSConfig: '%+v'", config))
	return newChannel, nil
}

/////////////////////////////////////////////////////////////////
// Private functions for SSLChannel
/////////////////////////////////////////////////////////////////

func initTLSConfig(props *utils.SortedProperties) (*tls.Config, types.TGError) {
	logger.Log(fmt.Sprint("======> Entering SSLChannel:initTLSConfig"))

	// Load System certificate
	rootCertPool, err := x509.SystemCertPool()
	if err != nil {
		errMsg := fmt.Sprint("ERROR: Returning SSLChannel::initTLSConfig Failed to read system certificate pool")
		logger.Error(errMsg)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
	}

	//sysTrustFile := fmt.Sprintf("%s%slib%ssecurity%scacerts", os.Getenv("JRE_HOME"), string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator))
	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:initTLSConfig about to ReadFile '%+v'", sysTrustFile))
	//pem, err := ioutil.ReadFile(sysTrustFile)
	//if err != nil {
	//	errMsg := fmt.Sprintf("ERROR: Returning SSLChannel::initTLSConfig Failed to read client certificate authority: %s", sysTrustFile)
	//	logger.Error(errMsg)
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
	//}

	certPool := x509.NewCertPool()
	clientCertificates := make([]tls.Certificate, 0)

	//logger.Log(fmt.Sprint("======> Inside SSLChannel:initTLSConfig about to add system certificate to certificate pool"))
	//if !certPool.AppendCertsFromPEM(pem) {
	//	errMsg := fmt.Sprintf("ERROR: Returning SSLChannel::initTLSConfig Can't parse client certificate data from '%s", sysTrustFile)
	//	logger.Error(errMsg)
	//	//return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, "")
	//}

	// Load the user defined certificates.
	trustedCerts := props.GetProperty(utils.GetConfigFromKey(utils.TlsTrustedCertificates), "")
	if trustedCerts == "" {
		errMsg := fmt.Sprint("WARNING: Returning SSLChannel::initTLSConfig There are no user defined certificates")
		logger.Warning(errMsg)
		//return nil, nil
	} else {
		userCertificateFilePaths := strings.Split(trustedCerts, ",")
		for _, userCertFile := range userCertificateFilePaths {
			userCertData, err := ioutil.ReadFile(userCertFile)
			if err != nil {
				errMsg := fmt.Sprintf("ERROR: Returning SSLChannel::initTLSConfig Failed to read user certificate file: %s", userCertFile)
				logger.Error(errMsg)
				return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
			}
			if !certPool.AppendCertsFromPEM(userCertData) {
				errMsg := fmt.Sprintf("ERROR: Returning SSLChannel::initTLSConfig Can't parse client certificate data from '%s", userCertFile)
				logger.Error(errMsg)
				return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, "")
			}
			clientCertificates = append(clientCertificates)
		}
	}

	tlsConfig := &tls.Config{
		Certificates:       clientCertificates,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          certPool,
		InsecureSkipVerify: false,
		Rand:               rand.Reader,
		RootCAs:            rootCertPool,
	}
	// TODO: This list may change if GO Language supports more or different - Ref. https://golang.org/pkg/crypto/tls/
	// TODO: Revisit later to find out if there is any API in GO to get this as a list.
	suites := []uint16{
		// TLS 1.0 - 1.2 cipher suites.
		tls.TLS_RSA_WITH_RC4_128_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		// TLS 1.3 cipher suites.
		//tls.TLS_AES_128_GCM_SHA256,
		//tls.TLS_AES_256_GCM_SHA384,
		//tls.TLS_CHACHA20_POLY1305_SHA256,
		// TLS_FALLBACK_SCSV isn't a standard cipher suite but an indicator
		// that the client is doing version fallback. See RFC 7507.
		tls.TLS_FALLBACK_SCSV,
	}
	supportedSuites := FilterSuitesById(suites)
	tlsConfig.CipherSuites = supportedSuites

	logger.Log(fmt.Sprint("======> Returning SSLChannel:initTLSConfig"))
	return tlsConfig, nil
}

//func (obj *SSLChannel) channelConnect() types.TGError {
//	//logger.Log(fmt.Sprint("======> Entering SSLChannel:channelConnect"))
//	if isChannelConnected(obj) {
//		logger.Log(fmt.Sprintf("======> SSLChannel::channelConnect channel is already connected"))
//		obj.setNoOfConnections(atomic.AddInt32(&ConnectionsToChannel, 1))
//		return nil
//	}
//	if isChannelClosed(obj) || obj.GetLinkState() == types.LinkNotConnected {
//		logger.Log(fmt.Sprintf("======> Inside SSLChannel:channelConnect about to channelTryRepeatConnect for object '%+v'", obj.String()))
//		err := channelTryRepeatConnect(obj, false)
//		if err != nil {
//			logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::channelConnect channelTryRepeatConnect failed"))
//			return err
//		}
//		obj.SetChannelLinkState(types.LinkConnected)
//		obj.setNoOfConnections(atomic.AddInt32(&ConnectionsToChannel, 1))
//		logger.Log(fmt.Sprintf("======> Returning SSLChannel:channelConnect successfully established socket connection and now has '%d' number of connections", obj.NumOfConnections))
//	} else {
//		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::channelConnect channelTryRepeatConnect - connect called on an invalid state := '%s'", obj.GetLinkState().String()))
//		errMsg := fmt.Sprintf("======> Connect called on an invalid state := '%s'", obj.GetLinkState().String())
//		return exception.NewTGGeneralExceptionWithMsg(errMsg)
//	}
//	logger.Log(fmt.Sprintf("======> Returning SSLChannel:channelConnect having '%d' number of connections", obj.GetNoOfConnections()))
//	return nil
//}

func (obj *SSLChannel) doAuthenticate() types.TGError {
	logger.Log(fmt.Sprint("======> Entering SSLChannel:doAuthenticate"))
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := pdu.CreateMessageForVerb(pdu.VerbAuthenticateRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: SSLChannel::doAuthenticate pdu.CreateMessageForVerb(pdu.VerbAuthenticateRequest) failed"))
		return err
	}

	msgRequest.(*pdu.AuthenticateRequestMessage).SetClientId(obj.clientId)
	msgRequest.(*pdu.AuthenticateRequestMessage).SetInboxAddr(obj.inboxAddress)
	msgRequest.(*pdu.AuthenticateRequestMessage).SetUserName(obj.getChannelUserName())
	msgRequest.(*pdu.AuthenticateRequestMessage).SetPassword(obj.getChannelPassword())

	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:doAuthenticate about to request reply for request '%+v'", msgRequest.String()))
	msgResponse, err := channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::doAuthenticate channelRequestReply failed"))
		return err
	}
	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:doAuthenticate received reply as '%+v'", msgResponse.String()))
	if !msgResponse.(*pdu.AuthenticateResponseMessage).IsSuccess() {
		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::doAuthenticate msgResponse.(*pdu.AuthenticateResponseMessage).IsSuccess() failed"))
		return exception.NewTGBadAuthenticationWithRealm(types.INTERNAL_SERVER_ERROR, types.TGErrorBadAuthentication, "Bad username/password combination", "", "tgdb")
	}

	obj.setChannelAuthToken(msgResponse.GetAuthToken())
	obj.setChannelSessionId(msgResponse.GetSessionId())

	cryptoDataGrapher, err := NewDataCryptoGrapher(msgResponse.GetSessionId(), msgResponse.(*pdu.AuthenticateResponseMessage).GetServerCertBuffer())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::doAuthenticate NewDataCryptoGrapher failed w/ '%+v'", err.Error()))
		return err
	}
	obj.setDataCryptoGrapher(cryptoDataGrapher)
	logger.Log(fmt.Sprintf("======> Returning SSLChannel:doAuthenticate Successfully authenticated for user: '%s'", obj.getChannelUserName()))
	return nil
}

func (obj *SSLChannel) performHandshake(sslMode bool) types.TGError {
	logger.Log(fmt.Sprint("======> Entering SSLChannel:performHandshake"))
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := pdu.CreateMessageForVerb(pdu.VerbHandShakeRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake pdu.CreateMessageForVerb(pdu.VerbAuthenticateRequest) failed"))
		return err
	}

	msgRequest.(*pdu.HandShakeRequestMessage).SetRequestType(pdu.InitiateRequest)

	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:performHandshake about to request reply for InitiateRequest '%+v'", msgRequest.String()))
	msgResponse, err := channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::doAuthenticate channelRequestReply failed"))
		return err
	}
	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:performHandshake received reply as '%+v'", msgResponse.String()))
	if msgResponse.GetVerbId() != pdu.VerbHandShakeResponse {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake HandshakeResponse message response NOT received"))
		if msgResponse.GetVerbId() == pdu.VerbSessionForcefullyTerminated {
			errMsg := msgResponse.(*pdu.SessionForcefullyTerminatedMessage).GetKillString()
			return exception.NewTGChannelDisconnectedWithMsg(errMsg)
		}
		errMsg := fmt.Sprintf("Expecting a HandshakeResponse message, and received: '%d'. Cannot connect to the server at: '%s'", msgResponse.GetVerbId(), obj.channelUrl.GetUrlAsString())
		return exception.NewTGGeneralException(types.TGDB_HNDSHKRESP_ERROR, types.TGErrorGeneralException, errMsg, "")
	}

	response := msgResponse.(*pdu.HandShakeResponseMessage)
	if response.GetResponseStatus() != pdu.ResponseAcceptChallenge {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake response.GetResponseStatus() is NOT pdu.ResponseAcceptChallenge"))
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

	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:performHandshake about to request reply for ChallengeAccepted '%+v'", msgRequest.String()))
	msgResponse, err = channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake channelRequestReply failed"))
		return err
	}
	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:performHandshake received reply (2) as '%+v'", msgResponse.String()))
	if msgResponse.GetVerbId() != pdu.VerbHandShakeResponse {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake HandshakeResponse message response NOT received"))
		if msgResponse.GetVerbId() == pdu.VerbSessionForcefullyTerminated {
			errMsg := msgResponse.(*pdu.SessionForcefullyTerminatedMessage).GetKillString()
			return exception.NewTGChannelDisconnectedWithMsg(errMsg)
		}
		errMsg := fmt.Sprintf("Expecting a HandshakeResponse message, and received: '%d'. Cannot connect to the server at: '%s'", msgResponse.GetVerbId(), obj.channelUrl.GetUrlAsString())
		return exception.NewTGGeneralException(types.TGDB_HNDSHKRESP_ERROR, types.TGErrorGeneralException, errMsg, "")
	}

	response = msgResponse.(*pdu.HandShakeResponseMessage)
	if response.GetResponseStatus() != pdu.ResponseProceedWithAuthentication {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake response.GetResponseStatus() is NOT pdu.ResponseAcceptChallenge"))
		errMsg := fmt.Sprintf("'%s': Handshake Failed. Cannot connect to the server at: '%s'", types.TGDB_HNDSHKRESP_ERROR, obj.channelUrl.GetUrlAsString())
		return exception.NewTGGeneralException(types.TGDB_HNDSHKRESP_ERROR, types.TGErrorGeneralException, errMsg, "")
	}
	logger.Log(fmt.Sprintf("======> Returning SSLChannel::performHandshake Handshake w/ Remote Server is successful."))
	return nil
}

func (obj *SSLChannel) setSocket(newSocket *tls.Conn) types.TGError {
	obj.socket = newSocket
	//err := obj.socket.SetNoDelay(true)
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning SSLChannel:setSocket Failed to set NoDelay flag to true"))
	//	failureMessage := fmt.Sprint("Failed to set NoDelay flag to true")
	//	return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, err.Error())
	//}
	//
	//err = obj.socket.SetLinger(0) // <= 0 means Do not linger
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning SSLChannel:setSocket Failed to set NoLinger flag to true"))
	//	failureMessage := fmt.Sprint("Failed to set NoLinger flag to true")
	//	return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, err.Error())
	//}

	buff := make([]byte, dataBufferSize)
	obj.input = iostream.NewProtocolDataInputStream(buff)
	obj.input.BufLen = 0
	obj.output = iostream.NewProtocolDataOutputStream(dataBufferSize)
	clientId := obj.GetProperties().GetProperty(utils.GetConfigFromKey(utils.ChannelClientId), "")
	obj.setChannelClientId(clientId)
	obj.setChannelInboxAddr(obj.socket.RemoteAddr().String()) //SS:TODO: Is this correct
	return nil
}

func (obj *SSLChannel) setBuffers(newSocket *tls.Conn) types.TGError {
	//sendSize := obj.ChannelProperties.GetPropertyAsInt(utils.GetConfigFromKey(utils.ChannelSendSize))
	//if sendSize > 0 {
	//	err := newSocket.SetWriteBuffer(sendSize*1024)
	//	if err != nil {
	//		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::setBuffers newSocket.SetWriteBuffer failed"))
	//		errMsg := fmt.Sprintf("SSLChannel:setBuffers unable to set write buffer limit to '%d'", sendSize*1024)
	//		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
	//	}
	//}
	//receiveSize := obj.ChannelProperties.GetPropertyAsInt(utils.GetConfigFromKey(utils.ChannelRecvSize))
	//if receiveSize > 0 {
	//	err := newSocket.SetReadBuffer(receiveSize*1024)
	//	if err != nil {
	//		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::setBuffers SetReadBuffer failed"))
	//		errMsg := fmt.Sprintf("SSLChannel:setBuffers unable to set read buffer limit to '%d'", receiveSize*1024)
	//		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
	//	}
	//}
	return nil
}

func (obj *SSLChannel) tryRead() (types.TGMessage, types.TGError) {
	//logger.Log(fmt.Sprint("======> Entering SSLChannel:tryRead"))
	n, err := obj.input.Available()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::tryRead obj.input.Available() failed"))
		errMsg := "SSLChannel::tryRead there is no data available to be read"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	if n <= 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning SSLChannel::tryRead as there are no bytes to read from the wire"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("======> Inside SSLChannel:tryRead about to request message '%d' bytes from the wire", n))
	return obj.ReadWireMsg()
}

func (obj *SSLChannel) validateHandshakeResponseVersion(sVersion int64, cVersion *utils.TGClientVersion) types.TGError {
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

func (obj *SSLChannel) writeLoop(done chan bool) {
	logger.Log(fmt.Sprint("======> Entering SSLChannel:writeLoop"))
	for {
		logger.Log(fmt.Sprintf("======> Inside SSLChannel:writeLoop entering infinite loop"))
		select { // Non-blocking channel operation
		case msg, ok := <-obj.msgCh: // Retrieve the message from the channel
			if !ok {
				//if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
				//	logMessage("SSLChannel::writeLoop unable to retrieve msg from the channel);
				//}
				logger.Error(fmt.Sprint("ERROR: Returning SSLChannel:writeLoop unable to retrieve message from obj.msgCh"))
				return
			}
			logger.Log(fmt.Sprintf("======> Inside SSLChannel:writeLoop retrieved message from obj.msgCh as '%+v'", msg.String()))

			err := obj.writeToWire(msg)
			if err != nil {
				// TODO: Revisit later - Do something
				logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::writeLoop unable to obj.writeToWire"))
				return
			}

			logger.Log(fmt.Sprintf("======> Inside SSLChannel:writeLoop successfully wrote message '%+v' on the socket", msg.String()))
			break
		default:
			// TODO: Revisit later - Do something
		}
	} // End of Infinite Loop
	// Send an acknowledgement of completion to the parent thread
	done <- true
	logger.Log(fmt.Sprint("======> Returning SSLChannel:writeLoop"))
}

func (obj *SSLChannel) writeToWire(msg types.TGMessage) types.TGError {
	logger.Log(fmt.Sprintf("======> Entering SSLChannel:writeToWire w/ Msg: '%+v'", msg.String()))
	obj.DisablePing()
	//bufLength := msg.GetMessageByteBufLength()
	//if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
	//	gLogger.log(TGLogger.TGLevel.DebugWire, "Send buf : %s", HexUtils.formatHex(buf, bufLen));
	//}
	msgBytes, bufLen, err := msg.ToBytes()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::writeToWire unable to convert message into byte format"))
		errMsg := fmt.Sprintf("SSLChannel::writeToWire unable to convert message into byte format")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.GetErrorMsg())
	}

	// Clear timeout deadlines set at the time of creation of the socket
	sErr := obj.socket.SetDeadline(time.Time{})
	if sErr != nil {
		logger.Error(fmt.Sprint("Returning SSLChannel::writeToWire unable to clear the deadline over SSL socket"))
		errMsg := fmt.Sprintf("SSLChannel::writeToWire unable to clear the deadline over SSL socket")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}

	// Reset timeout deadlines starting from NOW!!!
	timeout := utils.NewTGEnvironment().GetChannelConnectTimeout()
	sErr = obj.socket.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if sErr != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::writeToWire unable to reset the deadline over SSL socket"))
		errMsg := fmt.Sprintf("SSLChannel::writeToWire unable to reset the deadline over SSL socket")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}

	logger.Log(fmt.Sprintf("======> Inside SSLChannel:writeToWire about to write message bytes on the socket as '%+v'", msgBytes[0:bufLen]))
	// Put the data packet on the socket for network transmission
	_, sErr = obj.socket.Write(msgBytes[0:bufLen])
	if sErr != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::writeToWire unable to send message bytes over SSL socket"))
		errMsg := fmt.Sprintf("SSLChannel::writeToWire unable to send message bytes over SSL socket")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}
	logger.Log(fmt.Sprintf("======> Returning SSLChannel:writeToWire successfully wrote message bytes on the socket as '%+v'", msgBytes[0:bufLen]))
	return nil
}

/////////////////////////////////////////////////////////////////
// Helper functions for SSLChannel
/////////////////////////////////////////////////////////////////

func (obj *SSLChannel) GetIsClosed() bool {
	return obj.isSocketClosed
}

func (obj *SSLChannel) SetIsClosed(flag bool) {
	obj.isSocketClosed = flag
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChannel
/////////////////////////////////////////////////////////////////

// ChannelLock locks the communication channel between TGDB client and server
func (obj *SSLChannel) ChannelLock() {
	obj.sendLock.Lock()
}

// ChannelUnlock unlocks the communication channel between TGDB client and server
func (obj *SSLChannel) ChannelUnlock() {
	obj.sendLock.Unlock()
}

// Connect connects the underlying channel using the URL end point
func (obj *SSLChannel) Connect() types.TGError {
	return channelConnect(obj)
}

// DisablePing disables the pinging ability to the channel
func (obj *SSLChannel) DisablePing() {
	obj.needsPing = false
}

// Disconnect disconnects the channel from its URL end point
func (obj *SSLChannel) Disconnect() types.TGError {
	return channelDisConnect(obj)
}

// EnablePing enables the pinging ability to the channel
func (obj *SSLChannel) EnablePing() {
	obj.needsPing = true
}

// ExceptionLock locks the communication channel between TGDB client and server in case of business exceptions
func (obj *SSLChannel) ExceptionLock() {
	obj.exceptionLock.Lock()
}

// ExceptionUnlock unlocks the communication channel between TGDB client and server in case of business exceptions
func (obj *SSLChannel) ExceptionUnlock() {
	obj.exceptionLock.Unlock()
}

// GetAuthToken gets Authorization Token
func (obj *SSLChannel) GetAuthToken() int64 {
	return obj.authToken
}

// GetClientId gets Client Name
func (obj *SSLChannel) GetClientId() string {
	return obj.clientId
}

// GetChannelURL gets the Channel URL
func (obj *SSLChannel) GetChannelURL() types.TGChannelUrl {
	return obj.channelUrl
}

// GetConnectionIndex gets the Connection Index
func (obj *SSLChannel) GetConnectionIndex() int {
	return obj.connectionIndex
}

// GetExceptionCondition gets the Exception Condition
func (obj *SSLChannel) GetExceptionCondition() *sync.Cond {
	return obj.exceptionCond
}

// GetLinkState gets the Link/Channel State
func (obj *SSLChannel) GetLinkState() types.LinkState {
	return obj.channelLinkState
}

// GetNoOfConnections gets number of connections this channel has
func (obj *SSLChannel) GetNoOfConnections() int32 {
	return obj.numOfConnections
}

// GetPrimaryURL gets the Primary URL
func (obj *SSLChannel) GetPrimaryURL() types.TGChannelUrl {
	return obj.primaryUrl
}

// GetProperties gets the Channel Properties
func (obj *SSLChannel) GetProperties() types.TGProperties {
	return obj.channelProperties
}

// GetReader gets the Channel Reader
func (obj *SSLChannel) GetReader() types.TGChannelReader {
	return obj.reader
}

// GetResponses gets the Channel Response Map
func (obj *SSLChannel) GetResponses() map[int64]types.TGChannelResponse {
	return obj.responses
}

// GetSessionId gets Session id
func (obj *SSLChannel) GetSessionId() int64 {
	return obj.sessionId
}

// IsChannelPingable checks whether the channel is pingable or not
func (obj *SSLChannel) IsChannelPingable() bool {
	return obj.needsPing
}

// IsClosed checks whether channel is open or closed
func (obj *SSLChannel) IsClosed() bool {
	return isChannelClosed(obj)
}

// SendMessage sends a Message on this channel, and returns immediately - An Asynchronous or Non-Blocking operation
func (obj *SSLChannel) SendMessage(msg types.TGMessage) types.TGError {
	return channelSendMessage(obj, msg, true)
}

func (obj *SSLChannel) SendRequest(msg types.TGMessage, response types.TGChannelResponse) (types.TGMessage, types.TGError) {
	return channelSendRequest(obj, msg, response, true)
}

// SetChannelLinkState sets the Link/Channel State
func (obj *SSLChannel) SetChannelLinkState(state types.LinkState) {
	obj.channelLinkState = state
}

// SetChannelURL sets the channel URL
func (obj *SSLChannel) SetChannelURL(url types.TGChannelUrl) {
	obj.channelUrl = url.(*LinkUrl)
}

// SetConnectionIndex sets the connection index
func (obj *SSLChannel) SetConnectionIndex(index int) {
	obj.connectionIndex = index
}

// SetNoOfConnections sets number of connections
func (obj *SSLChannel) SetNoOfConnections(count int32) {
	obj.numOfConnections = count
}

// SetResponse sets the ChannelResponse Map
func (obj *SSLChannel) SetResponse(reqId int64, response types.TGChannelResponse) {
	obj.responses[reqId] = response
}

// Start starts the channel so that it can send and receive messages
func (obj *SSLChannel) Start() types.TGError {
	return channelStart(obj)
}

// Stop stops the channel forcefully or gracefully
func (obj *SSLChannel) Stop(bForcefully bool) {
	channelStop(obj, bForcefully)
}

// CreateSocket creates the physical link socket
func (obj *SSLChannel) CreateSocket() types.TGError {
	/**
			super.createSocket();
	        try {
	            sslSocket = (SSLSocket)sslSocketFactory.createSocket(this.socket, getHost(), getPort(), true);
	            String[] suites = sslSocket.getEnabledCipherSuites();
	            supportedSuites = TGCipherSuite.filterSuites(suites);
	            sslSocket.setEnabledCipherSuites(supportedSuites);
	            sslSocket.setEnabledProtocols(TLSProtocols);
	            sslSocket.setUseClientMode(true);
	        }
	        catch (IOException ioe) {
	            throw new TGChannelDisconnectedException(ioe);
	        }
	*/
	logger.Log(fmt.Sprintf("======> Entering SSLChannel:CreateSocket"))
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	obj.SetChannelLinkState(types.LinkNotConnected)
	host := obj.channelUrl.urlHost
	port := obj.channelUrl.urlPort
	serverAddr := fmt.Sprintf("%s:%d", host, port)
	logger.Log(fmt.Sprintf("======> Inside SSLChannel:CreateSocket attempting to resolve address for '%s'", serverAddr))

	//tcpAddr, tErr := net.ResolveSSLAddr(types.ProtocolSSL.String(), serverAddr)
	//if tErr != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CreateSocket net.ResolveSSLAddr failed"))
	//	errMsg := fmt.Sprintf("SSLChannel:CreateSocket unable to resolve channel address '%s'", serverAddr)
	//	return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, tErr.Error())
	//}
	////logger.Log(fmt.Sprintf("======> Inside SSLChannel:CreateSocket resolved SSL address for '%s' as '%+v'", serverAddr, tcpAddr))
	//
	sslConn, cErr := tls.Dial(types.ProtocolSSL.String(), serverAddr, obj.tlsConfig)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::CreateSocket Failed to connect to the server at '%s' w/ '%+v'", serverAddr, cErr.Error()))
		failureMessage := fmt.Sprintf("Failed to connect to the server at '%s'", serverAddr)
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
	}
	logger.Log(fmt.Sprintf("======> Inside SSLChannel:CreateSocket created SSL connection for '%s' as '%+v'", serverAddr, sslConn))

	timeout := utils.NewTGEnvironment().GetChannelConnectTimeout()
	dErr := sslConn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if dErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::CreateSocket Failed to set deadline of '%+v' seconds on the connection to the server", time.Duration(timeout)*time.Second))
		failureMessage := fmt.Sprintf("Failed to set the timeout '%d' on socket", timeout)
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, dErr.Error())
	}

	//err := sslConn.SetKeepAlive(true)
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CreateSocket Failed to set keep alive flag to true"))
	//	failureMessage := fmt.Sprint("Failed to set keep alive flag to true")
	//	return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
	//}

	// Set Read / Write Buffer Size on the socket
	tcErr := obj.setBuffers(sslConn)
	if tcErr != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CreateSocket Failed to set buffers"))
		return tcErr
	}
	tcErr = obj.setSocket(sslConn)
	if tcErr != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CreateSocket Failed to set socket value to the object"))
		return tcErr
	}
	obj.SetIsClosed(false)
	logger.Log(fmt.Sprintf("======> Returning SSLChannel:CreateSocket w/ SSL Connection as '%+v'", *obj.socket))
	return nil
}

// CloseSocket closes the socket
func (obj *SSLChannel) CloseSocket() types.TGError {
	logger.Log(fmt.Sprintf("======> Entering SSLChannel:CloseSocket w/ socket: '%+v'", obj.socket))
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	// TODO: Revisit later - Should the error be ignored?
	if obj.socket != nil {
		cErr := obj.socket.Close()
		if cErr != nil {
			logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CloseSocket obj.socket.Close() failed"))
			failureMessage := "Failed to close the socket to the server"
			return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
		}
	}
	obj.SetIsClosed(true)
	obj.socket = nil
	logger.Log(fmt.Sprintf("======> Returning SSLChannel:CloseSocket"))
	return nil
}

// OnConnect executes all the channel specific activities
func (obj *SSLChannel) OnConnect() types.TGError {
	logger.Log(fmt.Sprintf("======> Entering SSLChannel:OnConnect about to tryRead"))
	msg, err := obj.tryRead()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::OnConnect obj.tryRead() failed"))
		errMsg := "SSLChannel::OnConnect there is no data available to be read"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	if msg != nil {
		logger.Log(fmt.Sprintf("======> Inside SSLChannel:OnConnect tryRead() read Message as '%+v'", msg.String()))
	}

	if msg != nil && msg.GetVerbId() == pdu.VerbSessionForcefullyTerminated {
		logger.Warning(fmt.Sprint("WARNING: Returning SSLChannel:OnConnect since Message is of Forceful Termination Type"))
		return exception.NewTGChannelDisconnectedWithMsg(msg.(*pdu.SessionForcefullyTerminatedMessage).GetKillString())
	}

	// Check Host Verifier
	err1 := obj.socket.Handshake()
	if err1 != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::OnConnect obj.socket.Handshake() failed"))
		errMsg := "SSLChannel::OnConnect obj.socket.Handshake() failed"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err1.Error())
	}

	logger.Log(fmt.Sprintf("======> Inside SSLChannel:OnConnect about to performHandshake"))
	err = obj.performHandshake(true)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::OnConnect obj.performHandshake() failed"))
		errMsg := "SSLChannel::OnConnect error in performing handshake with server"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	logger.Log(fmt.Sprintf("======> Inside SSLChannel:OnConnect about to doAuthenticate"))
	err = obj.doAuthenticate()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::OnConnect obj.doAuthenticate() failed"))
		errMsg := "SSLChannel::OnConnect error in authentication with server"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	logger.Log(fmt.Sprintf("======> Returning SSLChannel:OnConnect"))
	return nil
}

// ReadWireMsg reads the message from the wire in the form of byte stream
func (obj *SSLChannel) ReadWireMsg() (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprintf("======> Entering SSLChannel:ReadWireMsg w/ SSLChannel as '%+v'", obj.String()))

	obj.input.BufLen = dataBufferSize
	in := obj.input
	if in == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning SSLChannel:ReadWireMsg since obj.input is NIL"))
		// TODO: Revisit later - Should we not return an error?
		return nil, nil
	}

	obj.DisablePing()
	if obj.GetIsClosed() {
		logger.Warning(fmt.Sprint("WARNING: Returning SSLChannel:ReadWireMsg since SSL Channel is Closed"))
		// TODO: Revisit later - Should we not return an error?
		return nil, nil
	}
	obj.lastActiveTime = time.Now()

	// Read the data on the socket
	buff := make([]byte, dataBufferSize)
	n, sErr := obj.socket.Read(buff)
	if sErr != nil || n <= 0 {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::ReadWireMsg obj.socket.Read failed"))
		errMsg := "SSLChannel::ReadWireMsg obj.socket.Read failed"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, sErr.Error())
	}
	logger.Log(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg Read '%d' bytes from the wire in buff '%+v'", n, buff[:(2*n)]))
	copy(in.Buf, buff[:n])
	in.BufLen = n
	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg Input Stream Buffer('%d') is '%+v'", in.BufLen, in.Buf[:(2*n)]))

	// Needed to avoid dirty data in the buffer when we handle the message
	buffer := make([]byte, n)
	//logger.Log(fmt.Sprint("======> Inside SSLChannel:ReadWireMsg in.ReadFullyAtPos read msgBytes as '%+v'", msgBytes))
	copy(buffer, buff[:n])
	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg copied into buffer as '%+v'", buffer))

	//intToBytes(size, msgBytes, 0)
	//bytesRead, _ := utils.FormatHex(msgBytes)
	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg bytes read: '%s'", bytesRead))

	msg, err := pdu.CreateMessageFromBuffer(buffer, 0, n)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::ReadWireMsg pdu.CreateMessageFromBuffer failed"))
		errMsg := "SSLChannel::ReadWireMsg unable to create a message from the input stream bytes"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}
	//logger.Log(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg Created Message from buffer as '%+v'", msg.String()))

	if msg.GetVerbId() == pdu.VerbExceptionMessage {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::ReadWireMsg msg.GetVerbId() == pdu.VerbExceptionMessage"))
		errMsg := msg.(*pdu.ExceptionMessage).GetExceptionMsg()
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	//if msg.GetVerbId() == pdu.VerbHandShakeResponse {
	//	if msg.GetResponseStatus() == pdu.ResponseChallengeFailed {
	//		errMsg := msg.GetErrorMessage()
	//		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::ReadWireMsg msg.GetVerbId() == pdu.VerbHandShakeResponse w/ '%+v'", errMsg))
	//		return nil, exception.GetErrorByType(types.TGErrorVersionMismatchException, "", errMsg, "")
	//	}
	//}
	logger.Log(fmt.Sprintf("======> Returning SSLChannel:ReadWireMsg w/ Message as '%+v'", msg.String()))
	return msg, nil
}

// Send Message to the server, compress and/or encrypt.
// Hence it is abstraction, that the Channel knows about it.
// @param msg       The message that needs to be sent to the server
func (obj *SSLChannel) Send(msg types.TGMessage) types.TGError {
	logger.Log(fmt.Sprintf("======> Entering SSLChannel:Send w/ Message as '%+v'", msg.String()))
	if obj.output == nil || obj.GetIsClosed() {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::Send as the channel is closed"))
		errMsg := fmt.Sprintf("SSLChannel:Send - unable to send message to server as the channel is closed")
		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, "")
	}

	return obj.writeToWire(msg)
}

func (obj *SSLChannel) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("SSLChannel:{")
	buffer.WriteString(fmt.Sprintf("IsSocketClosed: %+v", obj.isSocketClosed))
	buffer.WriteString(fmt.Sprintf(", MsgCh: %+v", obj.msgCh))
	buffer.WriteString(fmt.Sprintf(", Socket: %+v", obj.socket))
	//buffer.WriteString(fmt.Sprintf("Input: %+v ", obj.input.String()))
	//buffer.WriteString(fmt.Sprintf("Output: %+v ", obj.output.String()))
	//buffer.WriteString("\n")
	strArray := []string{buffer.String(), obj.channelToString() + "}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}
