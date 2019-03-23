package pdu

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"reflect"
	"strings"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: AuthenticatedMessage.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type AuthenticatedMessage struct {
	*AbstractProtocolMessage
	connectionId int
	clientId     string
}

func DefaultAuthenticatedMessage() *AuthenticatedMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AuthenticatedMessage{})

	newMsg := AuthenticatedMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
		connectionId:            -1,
		clientId:                "",
	}
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewAuthenticatedMessage(authToken, sessionId int64) *AuthenticatedMessage {
	newMsg := DefaultAuthenticatedMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for AuthenticatedMessage
/////////////////////////////////////////////////////////////////

func (msg *AuthenticatedMessage) GetClientId() string {
	return msg.clientId
}

func (msg *AuthenticatedMessage) GetConnectionId() int {
	return msg.connectionId
}

func (msg *AuthenticatedMessage) SetClientId(client string) {
	msg.clientId = client
}

func (msg *AuthenticatedMessage) SetConnectionId(connId int) {
	msg.connectionId = connId
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *AuthenticatedMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering AuthenticatedMessage:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AuthenticatedMessage:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AuthenticatedMessage:FromBytes - about to readHeader"))
	err = msg.readHeader(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AuthenticatedMessage:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("AuthenticatedMessage::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AuthenticatedMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering AuthenticatedMessage:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside AuthenticatedMessage:ToBytes - about to writeHeader"))
	err := msg.writeHeader(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AuthenticatedMessage:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("AuthenticatedMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *AuthenticatedMessage) GetAuthToken() int64 {
	return msg.getAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AuthenticatedMessage) GetIsUpdatable() bool {
	return msg.getIsUpdatable()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AuthenticatedMessage) GetMessageByteBufLength() int {
	return msg.getMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AuthenticatedMessage) GetRequestId() int64 {
	return msg.getRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AuthenticatedMessage) GetSequenceNo() int64 {
	return msg.getSequenceNo()
}

// GetSessionId gets the session id
func (msg *AuthenticatedMessage) GetSessionId() int64 {
	return msg.getSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AuthenticatedMessage) GetTimestamp() int64 {
	return msg.getTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AuthenticatedMessage) GetVerbId() int {
	return msg.getVerbId()
}

// SetAuthToken sets the authToken
func (msg *AuthenticatedMessage) SetAuthToken(authToken int64) {
	msg.setAuthToken(authToken)
}

// SetRequestId sets the request id
func (msg *AuthenticatedMessage) SetRequestId(requestId int64) {
	msg.setRequestId(requestId)
}

// SetSessionId sets the session id
func (msg *AuthenticatedMessage) SetSessionId(sessionId int64) {
	msg.setSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AuthenticatedMessage) SetTimestamp(timestamp int64) types.TGError {
	return msg.setTimestamp(timestamp)
}

func (msg *AuthenticatedMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AuthenticatedMessage:{")
	buffer.WriteString(fmt.Sprintf("ConnectionId: %d", msg.connectionId))
	buffer.WriteString(fmt.Sprintf(", ClientId: %s", msg.clientId))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.messageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AuthenticatedMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.updateSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AuthenticatedMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return msg.readHeader(is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *AuthenticatedMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return msg.writeHeader(os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *AuthenticatedMessage) ReadPayload(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering AuthenticatedMessage:ReadPayload"))

	authToken, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:ReadPayload w/ Error in reading authToken from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AuthenticatedMessage:ReadPayload read authToken as '%+v'", authToken))

	sessionId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:ReadPayload w/ Error in reading sessionId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read sessionId as '%+v'", sessionId))

	msg.SetAuthToken(authToken)
	msg.SetSessionId(sessionId)
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *AuthenticatedMessage) WritePayload(os types.TGOutputStream) types.TGError {
	startPos := os.GetPosition()
	logger.Log(fmt.Sprintf("Entering AuthenticatedMessage:WritePayload at output buffer position: '%d'", startPos))
	if msg.GetAuthToken() == -1 || msg.GetSessionId() == -1 {
		errMsg := fmt.Sprint("Message not authenticated")
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	os.(*iostream.ProtocolDataOutputStream).WriteLong(msg.GetAuthToken())
	os.(*iostream.ProtocolDataOutputStream).WriteLong(msg.GetSessionId())
	currPos := os.GetPosition()
	length := currPos - startPos
	logger.Log(fmt.Sprintf("Returning AuthenticatedMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *AuthenticatedMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.connectionId, msg.clientId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AuthenticatedMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *AuthenticatedMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.connectionId, &msg.clientId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AuthenticatedMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
