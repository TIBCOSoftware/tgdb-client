package pdu

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"reflect"
	"sync"
	"sync/atomic"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: AbstractProtocolMessage.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/**
 * The Server describes the pdu header as below.
 * struct _tg_pduheader_t_ {
	 tg_int32    length;         //length of the message including the header
	 tg_int32    magic;          //Magic to recognize this is our message
	 tg_int16    protVersion;    //protocol version
	 tg_pduverb  verbId;         //we write the verb as a short value
	 tg_uint64   sequenceNo;     //message SystemTypeSequence No from the client
	 tg_uint64   timestamp;      //Timestamp of the message sent.
	 tg_uint64   requestId;      //Unique _request Identifier from the client, which is returned
	 tg_int32    dataOffset;     //Offset from where the payload begins
 }
*/
// NOTE: DO NOT CHANGE THIS ORDER OF STRUCTURE ELEMENTS
type MessageHeader struct {
	BufLength int // Length of the message including the header
	//MagicId         int    // Intentionally to be Kept Private? - Magic to recognize this is our message
	//ProtocolVersion uint16 // protocol Version
	verbId     int
	sequenceNo int64 // Intentionally to be Kept Private? - Message SystemTypeSequence No from the client
	timestamp  int64 // Timestamp of the message sent
	requestId  int64 // Unique _request Identifier from the client, which is returned
	authToken  int64 // Only Authenticated messages (Post Successful User Login) will have proper value
	sessionId  int64 // Only Authenticated messages (Post Successful User Login) will have proper value
	dataOffset int16 // Offset from where the payload begins
}

var AtomicSequenceNumber int64

func defaultMessageHeader() *MessageHeader {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(MessageHeader{})

	seqNo := atomic.AddInt64(&AtomicSequenceNumber, 1)
	newMsg := MessageHeader{
		BufLength: -1,
		//MagicId:         utils.GetMagic(),
		//ProtocolVersion: utils.GetProtocolVersion(),
		verbId:     AbstractMessage,
		sequenceNo: seqNo,
		timestamp:  time.Now().UnixNano(), // C|Java had -1 as initialization - replaced for testing
		requestId:  -1,
		authToken:  0,
		sessionId:  0,
		dataOffset: -1,
	}
	return &newMsg
}

/** Various ways to get the size of a structure
func main() {
	fmt.Println("Hello, playground")
	newMsg := MessageHeader{
		BufLength:       -1,
		MagicId:         229822948,
		ProtocolVersion: 256,
		VerbId:          100,
		SequenceNo:      1,
		Timestamp:       time.Now().Unix(),
		RequestId:       -1,
		AuthToken:       0,
		SessionId:       0,
		DataOffset:      -1,
	}
	fmt.Printf("NewMsg: '%+v' w/ Message BufLength: '%+v'\n", newMsg, newMsg.BufLength)
	a := reflect.TypeOf(newMsg.SequenceNo).Size()
	b := reflect.TypeOf(newMsg).Size()
	c := unsafe.Sizeof(reflect.TypeOf(newMsg))
	fmt.Printf("A: '%d' B: '%d' C: '%d'\n", a, b, c)
	newMsg.BufLength = int(b)
	fmt.Printf("NewMsg: '%+v' w/ Message BufLength: '%+v'\n", newMsg, newMsg.BufLength)
}
*/

/////////////////////////////////////////////////////////////////
// Helper functions for MessageHeader
/////////////////////////////////////////////////////////////////

func (hdr *MessageHeader) HeaderGetMessageByteBufLength() int {
	return hdr.BufLength
}

func (hdr *MessageHeader) HeaderGetVerbId() int {
	return hdr.verbId
}

func (hdr *MessageHeader) HeaderGetSequenceNo() int64 {
	return hdr.sequenceNo
}

func (hdr *MessageHeader) HeaderGetTimestamp() int64 {
	if hdr.timestamp == -1 {
		hdr.timestamp = time.Now().Unix()
	}
	return hdr.timestamp
}

func (hdr *MessageHeader) HeaderGetRequestId() int64 {
	return hdr.requestId
}

func (hdr *MessageHeader) HeaderGetAuthToken() int64 {
	return hdr.authToken
}

func (hdr *MessageHeader) HeaderGetSessionId() int64 {
	return hdr.sessionId
}

func (hdr *MessageHeader) HeaderGetDataOffset() int16 {
	return hdr.dataOffset
}

func (hdr *MessageHeader) HeaderSetMessageByteBufLength(bufLength int) {
	hdr.BufLength = bufLength
}

func (hdr *MessageHeader) HeaderSetVerbId(verbId int) {
	hdr.verbId = verbId
}

func (hdr *MessageHeader) HeaderSetSequenceNo(sequenceNo int64) {
	hdr.sequenceNo = sequenceNo
}

func (hdr *MessageHeader) HeaderSetRequestId(requestId int64) {
	hdr.requestId = requestId
}

func (hdr *MessageHeader) HeaderSetAuthToken(authToken int64) {
	hdr.authToken = authToken
}

func (hdr *MessageHeader) HeaderSetSessionId(sessionId int64) {
	hdr.sessionId = sessionId
}

func (hdr *MessageHeader) HeaderSetDataOffset(dataOffset int16) {
	hdr.dataOffset = dataOffset
}

func (hdr *MessageHeader) HeaderSetTimestamp(timestamp int64) {
	hdr.timestamp = timestamp
}

// NOTE: Maintain the order of structure elements as shown above for streamlined server communication
type AbstractProtocolMessage struct {
	*MessageHeader
	isUpdatable bool
	bytesBuffer []byte
	contentLock sync.Mutex // reentrant-lock for synchronizing sending/receiving messages over the wire
}

func DefaultAbstractProtocolMessage() *AbstractProtocolMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AbstractProtocolMessage{})

	newMsg := AbstractProtocolMessage{
		MessageHeader: defaultMessageHeader(),
		isUpdatable:   false,
		bytesBuffer:   make([]byte, 0),
	}
	newMsg.BufLength = binary.Size(reflect.ValueOf(newMsg))
	return &newMsg
}

func NewAbstractProtocolMessage(authToken, sessionId int64) *AbstractProtocolMessage {
	newMsg := DefaultAbstractProtocolMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = binary.Size(reflect.ValueOf(newMsg))
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for AbstractProtocolMessage
/////////////////////////////////////////////////////////////////

func (msg *AbstractProtocolMessage) GetUpdatableFlag() bool {
	return msg.isUpdatable
}

func (msg *AbstractProtocolMessage) SetUpdatableFlag(updateFlag bool) {
	msg.isUpdatable = updateFlag
}

func (msg *AbstractProtocolMessage) SetSequenceAndTimeStamp(timestamp int64) types.TGError {
	if !(msg.GetIsUpdatable() || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	//err := msg.setTimestamp(timestamp)
	//if err != nil {
	//	return err
	//}
	if msg.GetIsUpdatable() {
		msg.sequenceNo = atomic.AddInt64(&AtomicSequenceNumber, 1)
		msg.BufLength = -1
		msg.bytesBuffer = make([]byte, 0)
	}
	return nil
}

func (msg *AbstractProtocolMessage) APMMessageToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("AbstractProtocolMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.HeaderGetMessageByteBufLength()))
	//buffer.WriteString(fmt.Sprintf("MagicId: %d", msg.MagicId))
	//buffer.WriteString(fmt.Sprintf(", ProtocolVersion: %d", msg.ProtocolVersion))
	buffer.WriteString(fmt.Sprintf(", VerbId: %d", msg.HeaderGetVerbId()))
	buffer.WriteString(fmt.Sprintf(", SequenceNo: %d", msg.HeaderGetSequenceNo()))
	buffer.WriteString(fmt.Sprintf(", Timestamp: %d", msg.HeaderGetTimestamp()))
	buffer.WriteString(fmt.Sprintf(", RequestId: %d", msg.HeaderGetRequestId()))
	buffer.WriteString(fmt.Sprintf(", AuthToken: %d", msg.HeaderGetAuthToken()))
	buffer.WriteString(fmt.Sprintf(", SessionId: %d", msg.HeaderGetSessionId()))
	buffer.WriteString(fmt.Sprintf(", DataOffset: %d", msg.HeaderGetDataOffset()))
	buffer.WriteString(fmt.Sprintf(", IsUpdatable: %+v", msg.GetIsUpdatable()))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGMessage
/////////////////////////////////////////////////////////////////

func APMReadHeader(msg types.TGMessage, is types.TGInputStream) types.TGError {
	logger.Error(fmt.Sprint("Entering AbstractProtocolMessage:APMReadHeader"))
	// First member attribute / element of message header is BufLength
	// It has already been read before reaching here - in FromBytes()

	magic, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading magicId from message buffer"))
		return err
	}
	if magic != utils.GetMagic() {
		errMsg := fmt.Sprint("Bad Magic id")
		return exception.GetErrorByType(types.TGErrorBadMagic, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read MagicId as '%d'", magic))

	protocolVersion, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading protocolVersion from message buffer"))
		return err
	}
	if protocolVersion != int16(utils.GetProtocolVersion()) {
		errMsg := fmt.Sprint("Unsupported protocol version")
		return exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read ProtocolVersion as '%d'", protocolVersion))

	verbId, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading verbId from message buffer"))
		return err
	}
	if verbId != int16(msg.GetVerbId()) {
		errMsg := fmt.Sprint("Incorrect Message Type")
		return exception.GetErrorByType(types.TGErrorBadVerb, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read VerbId as '%d'", verbId))

	sequenceNo, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading sequenceNo from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read SequenceNo as '%d'", sequenceNo))

	timestamp, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading timestamp from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read Timestamp as '%d'", timestamp))

	requestId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading requestId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read RequestId as '%d'", requestId))

	authToken, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading authToken from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read AuthToken as '%d'", authToken))

	sessionId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading sessionId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read SessionId as '%d'", sessionId))

	dataOffset, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading protocolVersion from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read DataOffset as '%d'", dataOffset))

	// (Re)Set the message attributes with correct values from input stream
	//msg.MagicId = magic
	//msg.ProtocolVersion = uint16(protocolVersion)
	msg.SetVerbId(int(verbId))
	msg.SetSequenceNo(sequenceNo)
	err = msg.SetTimestamp(timestamp) // Ignore Error handling
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in setting timestamp as message attribute"))
		return err
	}
	msg.SetRequestId(requestId)
	msg.SetAuthToken(authToken)
	msg.SetSessionId(sessionId)
	msg.SetDataOffset(dataOffset)
	//msg.SetMessageByteBufLength(binary.Size(reflect.ValueOf(msg)))
	msg.SetMessageByteBufLength(int(binary.Size(reflect.ValueOf(msg))))
	logger.Error(fmt.Sprint("Returning AbstractProtocolMessage:APMReadHeader"))
	return nil
}

func APMWriteHeader(msg types.TGMessage, os types.TGOutputStream) types.TGError {
	startPos := os.GetPosition()
	logger.Log(fmt.Sprintf("Entering AbstractProtocolMessage:APMWriteHeader at output buffer position: '%d'", startPos))
	os.(*iostream.ProtocolDataOutputStream).WriteInt(0) //The length is written later.
	os.(*iostream.ProtocolDataOutputStream).WriteInt(utils.GetMagic())
	os.(*iostream.ProtocolDataOutputStream).WriteShort(int(utils.GetProtocolVersion()))
	os.(*iostream.ProtocolDataOutputStream).WriteShort(msg.GetVerbId())

	os.(*iostream.ProtocolDataOutputStream).WriteLong(msg.GetSequenceNo())
	os.(*iostream.ProtocolDataOutputStream).WriteLong(msg.GetTimestamp())
	os.(*iostream.ProtocolDataOutputStream).WriteLong(msg.GetRequestId())

	os.(*iostream.ProtocolDataOutputStream).WriteLong(msg.GetAuthToken())
	os.(*iostream.ProtocolDataOutputStream).WriteLong(msg.GetSessionId())
	os.(*iostream.ProtocolDataOutputStream).WriteShort(os.GetPosition() + 2) //DataOffset.
	currPos := os.GetPosition()
	length := currPos - startPos
	logger.Log(fmt.Sprintf("Returning AbstractProtocolMessage::APMWriteHeader at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return nil
}

// VerbIdFromBytes extracts message type from the input buffer in the byte format
func VerbIdFromBytes(buffer []byte) (*CommandVerbs, types.TGError) {
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}
	logger.Log(fmt.Sprintf("Entering AbstractProtocolMessage:VerbIdFromBytes - received input buffer as '%+v'", buffer))

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:VerbIdFromBytes - extracted bufLen: '%d'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Second member attribute / element of message header is MagicId
	magic, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error in reading magicId from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:VerbIdFromBytes - extracted magic id: '%d'", magic))
	if magic != utils.GetMagic() {
		errMsg := fmt.Sprint("Bad Magic id")
		return nil, exception.GetErrorByType(types.TGErrorBadMagic, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Third member attribute / element of message header is ProtocolVersion
	protocolVersion, err := is.ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error in reading protocolVersion from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:VerbIdFromBytes - extracted protocolVersion: '%d'", protocolVersion))
	if protocolVersion != int16(utils.GetProtocolVersion()) {
		errMsg := fmt.Sprint("Unsupported protocol version")
		return nil, exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Fourth member attribute / element of message header is VerbId
	verbId, err := is.ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error in reading verbId from message buffer"))
		return nil, err
	}

	return GetVerb(int(verbId)), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *AbstractProtocolMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}
	logger.Log(fmt.Sprintf("Entering AbstractProtocolMessage:FromBytes - received input buffer as '%+v'", buffer))

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:FromBytes read bufLen as '%d'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AbstractProtocolMessage:FromBytes about to read Header data elements"))
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
	}

	logger.Log(fmt.Sprint("Inside AbstractProtocolMessage:FromBytes about to read Payload data elements"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
	}

	logger.Log(fmt.Sprintf("Returning AbstractProtocolMessage:FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AbstractProtocolMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering AbstractProtocolMessage:ToBytes"))

	msg.contentLock.Lock()
	defer msg.contentLock.Unlock()

	var bufLength int
	if msg.bytesBuffer == nil {
		os := iostream.DefaultProtocolDataOutputStream()

		err := APMWriteHeader(msg, os)
		if err != nil {
			errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
			return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
		}

		err = msg.WritePayload(os)
		if err != nil {
			errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
			return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
		}

		bufLength = os.GetLength()
		_, err = os.WriteIntAt(0, bufLength)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:ToBytes w/ Error in writing buffer length"))
			return nil, -1, err
		}
		msg.bytesBuffer = os.GetBuffer()
	} else {
		bufLength = len(msg.bytesBuffer)
	}
	logger.Log(fmt.Sprintf("Returning AbstractProtocolMessage:ToBytes resulted bytes-on-the-wire in '%+v'", msg.bytesBuffer))
	return msg.bytesBuffer, bufLength, nil
}

// GetAuthToken gets the authToken
func (msg *AbstractProtocolMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AbstractProtocolMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AbstractProtocolMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AbstractProtocolMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AbstractProtocolMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *AbstractProtocolMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AbstractProtocolMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AbstractProtocolMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *AbstractProtocolMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *AbstractProtocolMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *AbstractProtocolMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AbstractProtocolMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *AbstractProtocolMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *AbstractProtocolMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *AbstractProtocolMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AbstractProtocolMessage) SetTimestamp(timestamp int64) types.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *AbstractProtocolMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *AbstractProtocolMessage) String() string {
	return msg.APMMessageToString()
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AbstractProtocolMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AbstractProtocolMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *AbstractProtocolMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *AbstractProtocolMessage) ReadPayload(is types.TGInputStream) types.TGError {
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *AbstractProtocolMessage) WritePayload(os types.TGOutputStream) types.TGError {
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *AbstractProtocolMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo,
		msg.timestamp, msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractProtocolMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *AbstractProtocolMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractProtocolMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
