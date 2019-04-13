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
 * File name: VerbGetEntityRequest.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type GetEntityRequestMessage struct {
	*AbstractProtocolMessage
	commandType    int16 //0 - get, 1 - getbyid, 2 - get multiples, 10 - continue, 20 - close
	fetchSize      int
	batchSize      int
	traversalDepth int
	edgeLimit      int
	resultId       int
	key            types.TGKey
}

func DefaultGetEntityRequestMessage() *GetEntityRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GetEntityRequestMessage{})

	newMsg := GetEntityRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.isUpdatable = true
	newMsg.commandType = 0
	newMsg.fetchSize = 1000
	newMsg.batchSize = 50
	newMsg.traversalDepth = 3
	newMsg.verbId = VerbGetEntityRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewGetEntityRequestMessage(authToken, sessionId int64) *GetEntityRequestMessage {
	newMsg := DefaultGetEntityRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for GetEntityRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *GetEntityRequestMessage) GetBatchSize() int {
	return msg.batchSize
}

func (msg *GetEntityRequestMessage) GetCommand() int16 {
	return msg.commandType
}

func (msg *GetEntityRequestMessage) GetEdgeLimit() int {
	return msg.edgeLimit
}

func (msg *GetEntityRequestMessage) GetFetchSize() int {
	return msg.fetchSize
}

func (msg *GetEntityRequestMessage) GetKey() types.TGKey {
	return msg.key
}

func (msg *GetEntityRequestMessage) GetResultId() int {
	return msg.resultId
}

func (msg *GetEntityRequestMessage) GetTraversalDepth() int {
	return msg.traversalDepth
}

func (msg *GetEntityRequestMessage) SetBatchSize(size int) {
	if size < 10 || size > 32767 {
		msg.batchSize = 50
	} else {
		msg.batchSize = size
	}
}

func (msg *GetEntityRequestMessage) SetCommand(cmd int16) {
	msg.commandType = cmd
}

func (msg *GetEntityRequestMessage) SetEdgeLimit(size int) {
	if size < 0 || size > 32767 {
		msg.edgeLimit = 1000
	} else {
		msg.edgeLimit = size
	}
}

func (msg *GetEntityRequestMessage) SetFetchSize(size int) {
	if size < 0 {
		msg.fetchSize = 1000
	} else {
		msg.fetchSize = size
	}
}

func (msg *GetEntityRequestMessage) SetKey(key types.TGKey) {
	msg.key = key
}

func (msg *GetEntityRequestMessage) SetResultId(resultId int) {
	msg.resultId = resultId
}

func (msg *GetEntityRequestMessage) SetTraversalDepth(depth int) {
	if depth < 1 || depth > 1000 {
		msg.traversalDepth = 3
	} else {
		msg.traversalDepth = depth
	}
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *GetEntityRequestMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering GetEntityRequestMessage:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside GetEntityRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside GetEntityRequestMessage:FromBytes - about to APMReadHeader"))
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside GetEntityRequestMessage:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("GetEntityRequestMessage::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *GetEntityRequestMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering GetEntityRequestMessage:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside GetEntityRequestMessage:ToBytes - about to APMWriteHeader"))
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside GetEntityRequestMessage:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("GetEntityRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *GetEntityRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *GetEntityRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *GetEntityRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *GetEntityRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *GetEntityRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *GetEntityRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *GetEntityRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *GetEntityRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *GetEntityRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *GetEntityRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *GetEntityRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *GetEntityRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *GetEntityRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *GetEntityRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *GetEntityRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *GetEntityRequestMessage) SetTimestamp(timestamp int64) types.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *GetEntityRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *GetEntityRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("GetEntityRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("CommandType: %d", msg.commandType))
	buffer.WriteString(fmt.Sprintf(", FetchSize: %d", msg.fetchSize))
	buffer.WriteString(fmt.Sprintf(", BatchSize: %d", msg.batchSize))
	buffer.WriteString(fmt.Sprintf(", TraversalDepth: %d", msg.traversalDepth))
	buffer.WriteString(fmt.Sprintf(", EdgeLimit: %d", msg.edgeLimit))
	buffer.WriteString(fmt.Sprintf(", ResultId: %d", msg.resultId))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *GetEntityRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *GetEntityRequestMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *GetEntityRequestMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *GetEntityRequestMessage) ReadPayload(is types.TGInputStream) types.TGError {
	// No-op for Now
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *GetEntityRequestMessage) WritePayload(os types.TGOutputStream) types.TGError {
	startPos := os.GetPosition()
	logger.Log(fmt.Sprintf("Entering GetEntityRequestMessage:WritePayload at output buffer position: '%d'", startPos))
	os.(*iostream.ProtocolDataOutputStream).WriteShort(int(msg.GetCommand()))
	os.(*iostream.ProtocolDataOutputStream).WriteInt(msg.GetResultId())
	if msg.GetCommand() == 0 || msg.GetCommand() == 1 || msg.GetCommand() == 2 {
		os.(*iostream.ProtocolDataOutputStream).WriteInt(msg.GetFetchSize())
		os.(*iostream.ProtocolDataOutputStream).WriteShort(msg.GetBatchSize())
		os.(*iostream.ProtocolDataOutputStream).WriteShort(msg.GetTraversalDepth())
		os.(*iostream.ProtocolDataOutputStream).WriteShort(msg.GetEdgeLimit())
		err := msg.GetKey().WriteExternal(os)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning GetEntityRequestMessage:WritePayload w/ Error in writing key to message buffer"))
			return err
		}
	}
	currPos := os.GetPosition()
	length := currPos - startPos
	logger.Log(fmt.Sprintf("Returning GetEntityRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *GetEntityRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.commandType, msg.fetchSize,
		msg.batchSize, msg.traversalDepth, msg.edgeLimit, msg.resultId, msg.key)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GetEntityRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *GetEntityRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.commandType, &msg.fetchSize, &msg.batchSize, &msg.traversalDepth, &msg.edgeLimit, &msg.resultId, &msg.key)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GetEntityRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
