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
 * File name: VerbPingMessage.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type PingMessage struct {
	*AbstractProtocolMessage
}

func DefaultPingMessage() *PingMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(PingMessage{})

	newMsg := PingMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbPingMessage
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewPingMessage(authToken, sessionId int64) *PingMessage {
	newMsg := DefaultPingMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *PingMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering PingMessage:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning PingMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning PingMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside PingMessage:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside PingMessage:FromBytes - about to readHeader"))
	err = msg.readHeader(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside PingMessage:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("PingMessage::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *PingMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering PingMessage:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside PingMessage:ToBytes - about to writeHeader"))
	err := msg.writeHeader(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside PingMessage:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning PingMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("PingMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *PingMessage) GetAuthToken() int64 {
	return msg.getAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *PingMessage) GetIsUpdatable() bool {
	return msg.getIsUpdatable()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *PingMessage) GetMessageByteBufLength() int {
	return msg.getMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *PingMessage) GetRequestId() int64 {
	return msg.getRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *PingMessage) GetSequenceNo() int64 {
	return msg.getSequenceNo()
}

// GetSessionId gets the session id
func (msg *PingMessage) GetSessionId() int64 {
	return msg.getSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *PingMessage) GetTimestamp() int64 {
	return msg.getTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *PingMessage) GetVerbId() int {
	return msg.getVerbId()
}

// SetAuthToken sets the authToken
func (msg *PingMessage) SetAuthToken(authToken int64) {
	msg.setAuthToken(authToken)
}

// SetRequestId sets the request id
func (msg *PingMessage) SetRequestId(requestId int64) {
	msg.setRequestId(requestId)
}

// SetSessionId sets the session id
func (msg *PingMessage) SetSessionId(sessionId int64) {
	msg.setSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *PingMessage) SetTimestamp(timestamp int64) types.TGError {
	return msg.setTimestamp(timestamp)
}

func (msg *PingMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("PingMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.messageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *PingMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.updateSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *PingMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return msg.readHeader(is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *PingMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return msg.writeHeader(os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *PingMessage) ReadPayload(is types.TGInputStream) types.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *PingMessage) WritePayload(os types.TGOutputStream) types.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *PingMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning PingMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *PingMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning PingMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
