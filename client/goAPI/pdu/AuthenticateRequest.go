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
 * File name: VerbAuthenticateRequest.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type AuthenticateRequestMessage struct {
	*AbstractProtocolMessage
	clientId  string
	inboxAddr string
	userName  string
	password  []byte
}

func DefaultAuthenticateRequestMessage() *AuthenticateRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AuthenticateRequestMessage{})

	newMsg := AuthenticateRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
		clientId:                "",
		inboxAddr:               "",
		userName:                "",
		password:                make([]byte, 0), // TGConstants.EmptyByteArray
	}
	newMsg.verbId = VerbAuthenticateRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewAuthenticateRequestMessage(authToken, sessionId int64) *AuthenticateRequestMessage {
	newMsg := DefaultAuthenticateRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Public functions for AuthenticateRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *AuthenticateRequestMessage) GetClientId() string {
	return msg.clientId
}

func (msg *AuthenticateRequestMessage) GetInboxAddr() string {
	return msg.inboxAddr
}

func (msg *AuthenticateRequestMessage) GetUserName() string {
	return msg.userName
}

func (msg *AuthenticateRequestMessage) GetPassword() []byte {
	return msg.password
}

func (msg *AuthenticateRequestMessage) SetClientId(client string) {
	msg.clientId = client
}

func (msg *AuthenticateRequestMessage) SetInboxAddr(inbox string) {
	msg.inboxAddr = inbox
}

func (msg *AuthenticateRequestMessage) SetUserName(user string) {
	msg.userName = user
}

func (msg *AuthenticateRequestMessage) SetPassword(pwd []byte) {
	msg.password = pwd
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *AuthenticateRequestMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering AuthenticateRequestMessage:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
	if err != nil {
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AuthenticateRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AuthenticateRequestMessage:FromBytes - about to APMReadHeader"))
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AuthenticateRequestMessage:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("AuthenticateRequestMessage::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AuthenticateRequestMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering AuthenticateRequestMessage:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside AuthenticateRequestMessage:ToBytes - about to APMWriteHeader"))
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AuthenticateRequestMessage:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("AuthenticateRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *AuthenticateRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AuthenticateRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AuthenticateRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AuthenticateRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AuthenticateRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *AuthenticateRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AuthenticateRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AuthenticateRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *AuthenticateRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *AuthenticateRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *AuthenticateRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AuthenticateRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *AuthenticateRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *AuthenticateRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *AuthenticateRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AuthenticateRequestMessage) SetTimestamp(timestamp int64) types.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *AuthenticateRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *AuthenticateRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AuthenticateRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("ClientId: %s", msg.clientId))
	buffer.WriteString(fmt.Sprintf(", InboxAddr: %s", msg.inboxAddr))
	buffer.WriteString(fmt.Sprintf(", UserName: %s", msg.userName))
	buffer.WriteString(fmt.Sprintf(", Password: %s", string(msg.password)))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AuthenticateRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AuthenticateRequestMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *AuthenticateRequestMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *AuthenticateRequestMessage) ReadPayload(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering AuthenticateRequestMessage:ReadPayload"))
	// For Testing purpose only.
	bIsClientId, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading clientId flag from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("AuthenticateRequestMessage:ReadPayload read bIsClientId as '%+v'", bIsClientId))
	strClient := ""
	if ! bIsClientId {
		strClient, err = is.(*iostream.ProtocolDataInputStream).ReadUTF()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading clientId from message buffer"))
			return err
		}
		logger.Log(fmt.Sprintf("AuthenticateRequestMessage:ReadPayload read strClient as '%+v'", strClient))
	}

	inboxAddr, err := is.(*iostream.ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading inboxAddr flag from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("AuthenticateRequestMessage:ReadPayload read inboxAddr as '%+v'", inboxAddr))

	userName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading username from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("AuthenticateRequestMessage:ReadPayload read userName as '%+v'", userName))

	password, err := is.(*iostream.ProtocolDataInputStream).ReadBytes()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading password from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("AuthenticateRequestMessage:ReadPayload read password as '%+v'", password))

	msg.SetClientId(strClient)
	msg.SetInboxAddr(inboxAddr)
	msg.SetUserName(userName)
	msg.SetPassword(password)
	logger.Log(fmt.Sprint("Returning AuthenticateRequestMessage:ReadPayload"))
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *AuthenticateRequestMessage) WritePayload(os types.TGOutputStream) types.TGError {
	startPos := os.GetPosition()
	logger.Log(fmt.Sprintf("Entering AuthenticateRequestMessage:WritePayload at output buffer position: '%d'", startPos))
	if msg.GetClientId() == "" {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(true)
	} else {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(false) // No client id
		err := os.(*iostream.ProtocolDataOutputStream).WriteUTF(msg.GetClientId())
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:WritePayload w/ Error in writing clientId to message buffer"))
			return err
		}
	}
	if msg.GetInboxAddr() == "" {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(true)
	} else {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(false)
		err := os.(*iostream.ProtocolDataOutputStream).WriteUTF(msg.GetInboxAddr())
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:WritePayload w/ Error in writing inboxAddr to message buffer"))
			return err
		}
	}
	if msg.GetUserName() == "" {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(true)
	} else {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(false)
		err := os.(*iostream.ProtocolDataOutputStream).WriteUTF(msg.GetUserName())	// Can't be null.
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:WritePayload w/ Error in writing userName to message buffer"))
			return err
		}
	}
	err := os.(*iostream.ProtocolDataOutputStream).WriteBytes(msg.GetPassword())
	if err != nil {
		return err
	}
	currPos := os.GetPosition()
	length := currPos - startPos
	logger.Log(fmt.Sprintf("Returning AuthenticateRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return os.(*iostream.ProtocolDataOutputStream).WriteBytes(msg.GetPassword())
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *AuthenticateRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable, msg.clientId, msg.inboxAddr,
		msg.userName, msg.password)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AuthenticateRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *AuthenticateRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable,
		&msg.clientId, &msg.inboxAddr, &msg.userName, &msg.password)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AuthenticateRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
