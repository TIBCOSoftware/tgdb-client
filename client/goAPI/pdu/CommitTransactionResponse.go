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
 * File name: VerbCommitTransactionResponse.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type CommitTransactionResponse struct {
	*AbstractProtocolMessage
	addedIdList    []int64
	addedCount     int
	updatedIdList  []int64
	updatedCount   int
	removedIdList  []int64
	removedCount   int
	attrDescIdList []int64
	attrDescCount  int
	graphObjFact   types.TGGraphObjectFactory
	exception      *exception.TransactionException
	entityStream   types.TGInputStream
}

func DefaultCommitTransactionResponseMessage() *CommitTransactionResponse {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(CommitTransactionResponse{})

	newMsg := CommitTransactionResponse{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbCommitTransactionResponse
	newMsg.addedIdList = make([]int64, 0)
	newMsg.addedCount = 0
	newMsg.updatedIdList = make([]int64, 0)
	newMsg.updatedCount = 0
	newMsg.removedIdList = make([]int64, 0)
	newMsg.removedCount = 0
	newMsg.attrDescIdList = make([]int64, 0)
	newMsg.attrDescCount = 0
	newMsg.graphObjFact = nil
	//newMsg.entityStream = nil
	newMsg.exception = nil
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewCommitTransactionResponseMessage(authToken, sessionId int64) *CommitTransactionResponse {
	newMsg := DefaultCommitTransactionResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for VerbCommitTransactionResponse message
/////////////////////////////////////////////////////////////////

func ProcessTransactionStatus(is types.TGInputStream, status int) *exception.TransactionException {
	txnStatus := types.TGTransactionStatus(status)
	errMsg := "Error not available"
	logger.Log(fmt.Sprintf("Entering CommitTransactionResponse:ProcessTransactionStatus read txnStatus as '%d'", txnStatus))
	if txnStatus == types.TGTransactionSuccess {
		logger.Log(fmt.Sprintf("Returning CommitTransactionResponse:ProcessTransactionStatus NO EXCEPTION for txnStatus:'%+v'", txnStatus))
		return nil
	}
	//switch txnStatus {
	//case types.TGTransactionSuccess:
	//	logger.Log(fmt.Sprintf("Returning CommitTransactionResponse:ProcessTransactionStatus NO EXCEPTION for txnStatus:'%+v'", txnStatus))
	//	return nil
	//case types.TGTransactionAlreadyInProgress:
	//	//fallthrough
	//case types.TGTransactionClientDisconnected:
	//	fallthrough
	//case types.TGTransactionMalFormed:
	//	fallthrough
	//case types.TGTransactionGeneralError:
	//	fallthrough
	//case types.TGTransactionVerificationError:
	//	fallthrough
	//case types.TGTransactionInBadState:
	//	fallthrough
	//case types.TGTransactionUniqueConstraintViolation:
	//	fallthrough
	//case types.TGTransactionOptimisticLockFailed:
	//	fallthrough
	//case types.TGTransactionResourceExceeded:
	//	fallthrough
	//case types.TGCurrentThreadNotInTransaction:
	//	fallthrough
	//case types.TGTransactionUniqueIndexKeyAttributeNullError:
	//	fallthrough
	//default:
		errMsg, _ = is.(*iostream.ProtocolDataInputStream).ReadUTF()
		//if err == nil {
		//	errMsg = "Error not available"
		//}
		logger.Error(fmt.Sprintf("Returning CommitTransactionResponse:ProcessTransactionStatus for txnStatus:'%+v' w/ error: '%+v'", txnStatus, errMsg))
		return exception.BuildException(txnStatus, errMsg)
	//}
}

func (msg *CommitTransactionResponse) GetAddedEntityCount() int {
	return msg.addedCount
}

func (msg *CommitTransactionResponse) GetAddedIdList() []int64 {
	return msg.addedIdList
}

func (msg *CommitTransactionResponse) GetUpdatedEntityCount() int {
	return msg.updatedCount
}

func (msg *CommitTransactionResponse) GetUpdatedIdList() []int64 {
	return msg.updatedIdList
}

func (msg *CommitTransactionResponse) GetRemovedEntityCount() int {
	return msg.removedCount
}

func (msg *CommitTransactionResponse) GetRemovedIdList() []int64 {
	return msg.removedIdList
}

func (msg *CommitTransactionResponse) GetAttrDescCount() int {
	return msg.attrDescCount
}

func (msg *CommitTransactionResponse) GetAttrDescIdList() []int64 {
	return msg.attrDescIdList
}

func (msg *CommitTransactionResponse) HasException() bool {
	return msg.exception != nil
}

func (msg *CommitTransactionResponse) GetException() *exception.TransactionException {
	return msg.exception
}

func (msg *CommitTransactionResponse) SetAttrDescCount(count int) {
	msg.attrDescCount = count
}

func (msg *CommitTransactionResponse) SetAttrDescId(list []int64) {
	msg.attrDescIdList = list
}

func (msg *CommitTransactionResponse) SetAddEntityCount(count int) {
	msg.addedCount = count
}

func (msg *CommitTransactionResponse) SetAddedIdList(list []int64) {
	msg.addedIdList = list
}

func (msg *CommitTransactionResponse) SetUpdatedEntityCount(count int) {
	msg.updatedCount = count
}

func (msg *CommitTransactionResponse) SetUpdatedIdList(list []int64) {
	msg.updatedIdList = list
}

func (msg *CommitTransactionResponse) SetRemovedEntityCount(count int) {
	msg.removedCount = count
}

func (msg *CommitTransactionResponse) SetRemovedIdList(list []int64) {
	msg.removedIdList = list
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *CommitTransactionResponse) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering CommitTransactionResponse:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside CommitTransactionResponse:FromBytes - about to APMReadHeader"))
	err = APMReadHeader(msg, is)
	if err != nil {
		//errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		//return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		return nil, err
	}

	logger.Log(fmt.Sprint("Inside CommitTransactionResponse:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		//errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		//return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		return nil, err
	}

	logger.Log(fmt.Sprintf("CommitTransactionResponse::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *CommitTransactionResponse) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering CommitTransactionResponse:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside CommitTransactionResponse:ToBytes - about to APMWriteHeader"))
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside CommitTransactionResponse:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("CommitTransactionResponse::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *CommitTransactionResponse) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *CommitTransactionResponse) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *CommitTransactionResponse) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *CommitTransactionResponse) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *CommitTransactionResponse) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *CommitTransactionResponse) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *CommitTransactionResponse) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *CommitTransactionResponse) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *CommitTransactionResponse) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *CommitTransactionResponse) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *CommitTransactionResponse) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *CommitTransactionResponse) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *CommitTransactionResponse) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *CommitTransactionResponse) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *CommitTransactionResponse) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *CommitTransactionResponse) SetTimestamp(timestamp int64) types.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *CommitTransactionResponse) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *CommitTransactionResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("CommitTransactionResponse:{")
	buffer.WriteString(fmt.Sprintf("AttrDescCount: %d", msg.attrDescCount))
	buffer.WriteString(fmt.Sprint(", AttrDescSet:{"))
	//for _, d := range msg.attrDescIdList {
	//	buffer.WriteString(fmt.Sprintf("Attribute Descriptor: %d ", d))
	//}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprintf(", AddedListCount: %d", msg.addedCount))
	buffer.WriteString(fmt.Sprint(", AddedIdList:{"))
	for _, v := range msg.addedIdList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d ", v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprintf(", UpdatedListCount: %d", msg.updatedCount))
	buffer.WriteString(fmt.Sprint(", UpdatedIdList:{"))
	for _, v := range msg.updatedIdList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d ", v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprintf(", RemovedListCount: %d", msg.removedCount))
	buffer.WriteString(fmt.Sprint(", RemovedIdList:{"))
	for _, v := range msg.removedIdList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d ", v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *CommitTransactionResponse) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *CommitTransactionResponse) ReadHeader(is types.TGInputStream) types.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *CommitTransactionResponse) WriteHeader(os types.TGOutputStream) types.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *CommitTransactionResponse) ReadPayload(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering CommitTransactionResponse:ReadPayload"))
	bLen, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // buf length
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading bLen from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read buf length as '%+v'", bLen))

	checkSum, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // checksum
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading checkSum from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read checkSum as '%+v'", checkSum))

	status, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // status code - currently zero
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading status from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read status as '%+v'", status))

	logger.Log(fmt.Sprint("Inside CommitTransactionResponse:ReadPayload - about to ProcessTransactionStatus"))
	txnException := ProcessTransactionStatus(is, status)
	if txnException != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ txnException"))
		return txnException
	}

	for {
		avail, err := is.(*iostream.ProtocolDataInputStream).Available()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading availability bytes from message buffer"))
			return err
		}
		if avail <= 0 {
			break
		}
		logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read avail as '%+v'", avail))

		opCode, err := is.(*iostream.ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading opCode from message buffer"))
			return err
		}
		logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read opCode as '%+v'", opCode))

		switch opCode {
		case 0x1010:
			attrDescCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading attrDescCount from message buffer"))
				return err
			}
			msg.attrDescCount = attrDescCount
			logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read attrDescCount for opCode=0x1010 as '%+v'", attrDescCount))
			for i := 0; i < attrDescCount; i++ {
				tempId, _ := is.(*iostream.ProtocolDataInputStream).ReadInt()
				logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read tempId for opCode=0x1010 as '%+v'", tempId))
				msg.attrDescIdList = append(msg.attrDescIdList, int64(tempId))
				realId, _ := is.(*iostream.ProtocolDataInputStream).ReadInt()
				logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read realId for opCode=0x1010 as '%+v'", realId))
				msg.attrDescIdList = append(msg.attrDescIdList, int64(realId))
			}
			break
		case 0x1011:
			addedCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading addedCount from message buffer"))
				return err
			}
			msg.addedCount = addedCount
			logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read addedCount for opCode=0x1011 as '%+v'", addedCount))
			for i := 0; i < addedCount; i++ {
				longTempId, _ := is.(*iostream.ProtocolDataInputStream).ReadLong()
				logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read longTempId for opCode=0x1011 as '%+v'", longTempId))
				msg.addedIdList = append(msg.addedIdList, longTempId)
				longRealId, _ := is.(*iostream.ProtocolDataInputStream).ReadLong()
				logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read longRealId for opCode=0x1011 as '%+v'", longRealId))
				msg.addedIdList = append(msg.addedIdList, longRealId)
				longVersion, _ := is.(*iostream.ProtocolDataInputStream).ReadLong()
				logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read longVersion for opCode=0x1011 as '%+v'", longVersion))
				msg.addedIdList = append(msg.addedIdList, longVersion)
			}
			break
		case 0x1012:
			updatedCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading updatedCount from message buffer"))
				return err
			}
			msg.updatedCount = updatedCount
			logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read updatedCount for opCode=0x1012 as '%+v'", updatedCount))
			for i := 0; i < updatedCount; i++ {
				id, _ := is.(*iostream.ProtocolDataInputStream).ReadLong()
				logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read id for opCode=0x1012 as '%+v'", id))
				msg.updatedIdList = append(msg.updatedIdList, id)
				version, _ := is.(*iostream.ProtocolDataInputStream).ReadLong()
				logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read version for opCode=0x1012 as '%+v'", version))
				msg.updatedIdList = append(msg.updatedIdList, version)
			}
			break
		case 0x1013:
			removedCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading removedCount from message buffer"))
				return err
			}
			msg.removedCount = removedCount
			logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read removedCount for opCode=0x1013 as '%+v'", removedCount))
			for i := 0; i < removedCount; i++ {
				id, _ := is.(*iostream.ProtocolDataInputStream).ReadLong()
				logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read id for opCode=0x1013 as '%+v'", id))
				msg.removedIdList = append(msg.removedIdList, id)
			}
			break
		case 0x6789:
			msg.entityStream = is
			pos := is.(*iostream.ProtocolDataInputStream).GetPosition()
			logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read pos for opCode=0x6789 as '%+v'", pos))
			count, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading count from message buffer"))
				return err
			}
			logger.Log(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read count for opCode=0x6789 as '%+v'", count))
			is.(*iostream.ProtocolDataInputStream).SetPosition(pos)
			break
		default:
			break
		}
	}
	logger.Log(fmt.Sprint("Returning CommitTransactionResponse:ReadPayload"))
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *CommitTransactionResponse) WritePayload(os types.TGOutputStream) types.TGError {
	// No-op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *CommitTransactionResponse) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable, msg.addedIdList, msg.addedCount,
		msg.updatedIdList, msg.updatedCount, msg.removedIdList, msg.removedCount, msg.attrDescIdList, msg.attrDescCount,
		msg.graphObjFact, msg.exception)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CommitTransactionResponse:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *CommitTransactionResponse) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable, &msg.addedIdList,
		&msg.addedCount, &msg.updatedIdList, &msg.updatedCount, &msg.removedIdList, &msg.removedCount, &msg.attrDescIdList,
		&msg.attrDescCount, &msg.graphObjFact, &msg.exception)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CommitTransactionResponse:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
