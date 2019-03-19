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
	VerbId     int
	SequenceNo int64 // Intentionally to be Kept Private? - Message SystemTypeSequence No from the client
	Timestamp  int64 // Timestamp of the message sent
	RequestId  int64 // Unique _request Identifier from the client, which is returned
	AuthToken  int64 // Only Authenticated messages (Post Successful User Login) will have proper value
	SessionId  int64 // Only Authenticated messages (Post Successful User Login) will have proper value
	DataOffset int16 // Offset from where the payload begins
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
		VerbId:     AbstractMessage,
		SequenceNo: seqNo,
		Timestamp:  time.Now().UnixNano(), // C|Java had -1 as initialization - replaced for testing
		RequestId:  -1,
		AuthToken:  0,
		SessionId:  0,
		DataOffset: -1,
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

// NOTE: Maintain the order of structure elements as shown above for streamlined server communication
type AbstractProtocolMessage struct {
	*MessageHeader
	IsUpdatable bool
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
		IsUpdatable:   false,
		bytesBuffer:   make([]byte, 0),
	}
	newMsg.BufLength = binary.Size(reflect.ValueOf(newMsg))
	return &newMsg
}

func NewAbstractProtocolMessage(authToken, sessionId int64) *AbstractProtocolMessage {
	newMsg := DefaultAbstractProtocolMessage()
	newMsg.AuthToken = authToken
	newMsg.SessionId = sessionId
	newMsg.BufLength = binary.Size(reflect.ValueOf(newMsg))
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Private functions for AbstractProtocolMessage / TGMessage - used in all derived messages
/////////////////////////////////////////////////////////////////

func (msg *AbstractProtocolMessage) getAuthToken() int64 {
	return msg.AuthToken
}

func (msg *AbstractProtocolMessage) getIsUpdatable() bool {
	return msg.IsUpdatable
}

func (msg *AbstractProtocolMessage) getMessageByteBufLength() int {
	return msg.BufLength
}

func (msg *AbstractProtocolMessage) getRequestId() int64 {
	return msg.RequestId
}

func (msg *AbstractProtocolMessage) getSequenceNo() int64 {
	return msg.SequenceNo
}

func (msg *AbstractProtocolMessage) getSessionId() int64 {
	return msg.SessionId
}

func (msg *AbstractProtocolMessage) getTimestamp() int64 {
	if msg.Timestamp == -1 {
		msg.Timestamp = time.Now().Unix()
	}
	return msg.Timestamp
}

func (msg *AbstractProtocolMessage) getVerbId() int {
	return msg.VerbId
}

func (msg *AbstractProtocolMessage) setAuthToken(authToken int64) {
	msg.AuthToken = authToken
}

func (msg *AbstractProtocolMessage) setIsUpdatable(updateFlag bool) {
	msg.IsUpdatable = updateFlag
}

func (msg *AbstractProtocolMessage) setMessageByteBufLength(bufLength int) {
	msg.BufLength = bufLength
}

func (msg *AbstractProtocolMessage) setRequestId(requestId int64) {
	msg.RequestId = requestId
}

func (msg *AbstractProtocolMessage) setSequenceNo(sequenceNo int64) {
	msg.SequenceNo = sequenceNo
}

func (msg *AbstractProtocolMessage) setSessionId(sessionId int64) {
	msg.SessionId = sessionId
}

func (msg *AbstractProtocolMessage) setTimestamp(timestamp int64) types.TGError {
	if !(msg.IsUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning readHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.VerbId).Name)
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Timestamp = timestamp
	return nil
}

func (msg *AbstractProtocolMessage) setVerbId(verbId int) {
	msg.VerbId = verbId
}

func (msg *AbstractProtocolMessage) updateSequenceAndTimeStamp(timestamp int64) types.TGError {
	err := msg.setTimestamp(timestamp)
	if err != nil {
		return err
	}
	if msg.IsUpdatable {
		msg.SequenceNo = atomic.AddInt64(&AtomicSequenceNumber, 1)
		msg.BufLength = -1
		msg.bytesBuffer = make([]byte, 0)
	}
	return nil
}

func (msg *AbstractProtocolMessage) readHeader(is types.TGInputStream) types.TGError {
	logger.Error(fmt.Sprint("Entering AbstractProtocolMessage:readHeader"))
	// First member attribute / element of message header is BufLength
	// It has already been read before reaching here - in FromBytes()

	magic, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading magicId from message buffer"))
		return err
	}
	if magic != utils.GetMagic() {
		errMsg := fmt.Sprint("Bad Magic id")
		return exception.GetErrorByType(types.TGErrorBadMagic, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read MagicId as '%d'", magic))

	protocolVersion, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading protocolVersion from message buffer"))
		return err
	}
	if protocolVersion != int16(utils.GetProtocolVersion()) {
		errMsg := fmt.Sprint("Unsupported protocol version")
		return exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read ProtocolVersion as '%d'", protocolVersion))

	verbId, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading verbId from message buffer"))
		return err
	}
	if verbId != int16(msg.GetVerbId()) {
		errMsg := fmt.Sprint("Incorrect Message Type")
		return exception.GetErrorByType(types.TGErrorBadVerb, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read VerbId as '%d'", verbId))

	sequenceNo, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading sequenceNo from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read SequenceNo as '%d'", sequenceNo))

	timestamp, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading timestamp from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read Timestamp as '%d'", timestamp))

	requestId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading requestId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read RequestId as '%d'", requestId))

	authToken, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading authToken from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read AuthToken as '%d'", authToken))

	sessionId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading sessionId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read SessionId as '%d'", sessionId))

	dataOffset, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in reading protocolVersion from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractProtocolMessage:readHeader read DataOffset as '%d'", dataOffset))

	// (Re)Set the message attributes with correct values from input stream
	//msg.MagicId = magic
	//msg.ProtocolVersion = uint16(protocolVersion)
	msg.SetVerbId(int(verbId))
	msg.SetSequenceNo(sequenceNo)
	err = msg.SetTimestamp(timestamp) // Ignore Error handling
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:readHeader w/ Error in setting timestamp as message attribute"))
		return err
	}
	msg.SetRequestId(requestId)
	msg.SetAuthToken(authToken)
	msg.SetSessionId(sessionId)
	msg.DataOffset = dataOffset
	//msg.SetMessageByteBufLength(binary.Size(reflect.ValueOf(msg)))
	msg.SetMessageByteBufLength(int(binary.Size(reflect.ValueOf(msg))))
	logger.Error(fmt.Sprint("Returning AbstractProtocolMessage:readHeader"))
	return nil
}

func (msg *AbstractProtocolMessage) writeHeader(os types.TGOutputStream) types.TGError {
	startPos := os.GetPosition()
	logger.Log(fmt.Sprintf("Entering AbstractProtocolMessage:writeHeader at output buffer position: '%d'", startPos))
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
	logger.Log(fmt.Sprintf("Returning AbstractProtocolMessage::writeHeader at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return nil
}

func (msg *AbstractProtocolMessage) messageToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("AbstractProtocolMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	//buffer.WriteString(fmt.Sprintf("MagicId: %d", msg.MagicId))
	//buffer.WriteString(fmt.Sprintf(", ProtocolVersion: %d", msg.ProtocolVersion))
	buffer.WriteString(fmt.Sprintf(", VerbId: %d", msg.VerbId))
	buffer.WriteString(fmt.Sprintf(", SequenceNo: %d", msg.SequenceNo))
	buffer.WriteString(fmt.Sprintf(", Timestamp: %d", msg.Timestamp))
	buffer.WriteString(fmt.Sprintf(", RequestId: %d", msg.RequestId))
	buffer.WriteString(fmt.Sprintf(", AuthToken: %d", msg.AuthToken))
	buffer.WriteString(fmt.Sprintf(", SessionId: %d", msg.SessionId))
	buffer.WriteString(fmt.Sprintf(", DataOffset: %d", msg.DataOffset))
	buffer.WriteString(fmt.Sprintf(", IsUpdatable: %+v", msg.IsUpdatable))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Helper functions for AbstractProtocolMessage
/////////////////////////////////////////////////////////////////

// SetIsUpdatable sets the updatable flag
func (msg *AbstractProtocolMessage) SetIsUpdatable(updateFlag bool) {
	msg.setIsUpdatable(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AbstractProtocolMessage) SetMessageByteBufLength(bufLength int) {
	msg.setMessageByteBufLength(bufLength)
}

// SetSequenceNo sets the sequenceNo
func (msg *AbstractProtocolMessage) SetSequenceNo(sequenceNo int64) {
	msg.setSequenceNo(sequenceNo)
}

// SetVerbId sets verbId of the message
func (msg *AbstractProtocolMessage) SetVerbId(verbId int) {
	msg.setVerbId(verbId)
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
	err = msg.readHeader(is)
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

		err := msg.writeHeader(os)
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
	return msg.getAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AbstractProtocolMessage) GetIsUpdatable() bool {
	return msg.getIsUpdatable()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AbstractProtocolMessage) GetMessageByteBufLength() int {
	return msg.getMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AbstractProtocolMessage) GetRequestId() int64 {
	return msg.getRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AbstractProtocolMessage) GetSequenceNo() int64 {
	return msg.getSequenceNo()
}

// GetSessionId gets the session id
func (msg *AbstractProtocolMessage) GetSessionId() int64 {
	return msg.getSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AbstractProtocolMessage) GetTimestamp() int64 {
	return msg.getTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AbstractProtocolMessage) GetVerbId() int {
	return msg.getVerbId()
}

// SetAuthToken sets the authToken
func (msg *AbstractProtocolMessage) SetAuthToken(authToken int64) {
	msg.setAuthToken(authToken)
}

// SetRequestId sets the request id
func (msg *AbstractProtocolMessage) SetRequestId(requestId int64) {
	msg.setRequestId(requestId)
}

// SetSessionId sets the session id
func (msg *AbstractProtocolMessage) SetSessionId(sessionId int64) {
	msg.setSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AbstractProtocolMessage) SetTimestamp(timestamp int64) types.TGError {
	return msg.setTimestamp(timestamp)
}

func (msg *AbstractProtocolMessage) String() string {
	return msg.messageToString()
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AbstractProtocolMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.updateSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AbstractProtocolMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return msg.readHeader(is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *AbstractProtocolMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return msg.writeHeader(os)
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
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.VerbId, msg.SequenceNo,
		msg.Timestamp, msg.RequestId, msg.AuthToken, msg.SessionId, msg.DataOffset, msg.IsUpdatable)
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
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.VerbId, &msg.SequenceNo,
		&msg.Timestamp, &msg.RequestId, &msg.AuthToken, &msg.SessionId, &msg.DataOffset, &msg.IsUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractProtocolMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
