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
 * File Name: pduimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: pduimpl.go 4575 2020-10-27 00:21:18Z nimish $
 */

package impl

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"reflect"
	"tgdb"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)


const (
	AbstractMessage int = -100
)

const (
	// Ping Message - Heart beats
	VerbPingMessage int = 0
	// HandShake Request/Response protocol
	VerbHandShakeRequest  int = 1
	VerbHandShakeResponse int = 2
	// Authenticate Request/Response protocol
	VerbAuthenticateRequest  int = 3
	VerbAuthenticateResponse int = 4
	// Transaction - begin/commit/rollback protocol verbs
	VerbBeginTransactionRequest     int = 5
	VerbBeginTransactionResponse    int = 6
	VerbCommitTransactionRequest    int = 7
	VerbCommitTransactionResponse   int = 8
	VerbRollbackTransactionRequest  int = 9
	VerbRollbackTransactionResponse int = 10
	// Query Request/Response verbs
	VerbQueryRequest  int = 11
	VerbQueryResponse int = 12
	// Graph Traversal verbs
	VerbTraverseRequest  int = 13
	VerbTraverseResponse int = 14
	// Admin Request/Response verbs
	VerbAdminRequest  int = 15
	VerbAdminResponse int = 16
	// Retrieve meta data
	VerbMetadataRequest  int = 19
	VerbMetadataResponse int = 20
	// Get entities
	VerbGetEntityRequest  int = 21
	VerbGetEntityResponse int = 22
	// Get LargeObject
	VerbGetLargeObjectRequest  int = 23
	VerbGetLargeObjectResponse int = 24
	// Import/Export verbs - They are admin request, and not supported by Java
	VerbBeginExportRequest    int = 25
	VerbBeginExportResponse   int = 26
	VerbPartialExportRequest  int = 27
	VerbPartialExportResponse int = 28
	VerbCancelExportRequest   int = 29
	VerbBeginImportRequest    int = 31
	BeginImportResponse       int = 32
	VerbPartialImportRequest  int = 33
	VerbPartialImportResponse int = 34
	// Dump Stacktrace request verb
	VerbDumpStacktraceRequest int = 39
	// Disconnect Request verbs
	VerbDisconnectChannelRequest    int = 40
	VerbSessionForcefullyTerminated int = 41
	// Decryption Request verbs
	VerbDecryptBufferRequest  int = 44
	VerbDecryptBufferResponse int = 45
	// Unknown Exception Message on the server.
	VerbExceptionMessage int = 100
	VerbInvalidMessage   int = -1
)

type CommandVerbs struct {
	id          int
	name        string
	implementor string
}

var PreDefinedVerbs = map[int]CommandVerbs{
	VerbPingMessage:                 {id: VerbPingMessage, name: "VerbPingMessage", implementor: "pdu.VerbPingMessage"},
	VerbHandShakeRequest:            {id: VerbHandShakeRequest, name: "VerbHandShakeRequest", implementor: "pdu.HandshakeRequest"},
	VerbHandShakeResponse:           {id: VerbHandShakeResponse, name: "VerbHandShakeResponse", implementor: "pdu.HandshakeResponse"},
	VerbAuthenticateRequest:         {id: VerbAuthenticateRequest, name: "VerbAuthenticateRequest", implementor: "pdu.VerbAuthenticateRequest"},
	VerbAuthenticateResponse:        {id: VerbAuthenticateResponse, name: "VerbAuthenticateResponse", implementor: "pdu.VerbAuthenticateResponse"},
	VerbBeginTransactionRequest:     {id: VerbBeginTransactionRequest, name: "VerbBeginTransactionRequest", implementor: "pdu.VerbBeginTransactionRequest"},
	VerbBeginTransactionResponse:    {id: VerbBeginTransactionResponse, name: "VerbBeginTransactionResponse", implementor: "pdu.VerbBeginTransactionResponse"},
	VerbCommitTransactionRequest:    {id: VerbCommitTransactionRequest, name: "VerbCommitTransactionRequest", implementor: "pdu.VerbCommitTransactionRequest"},
	VerbCommitTransactionResponse:   {id: VerbCommitTransactionResponse, name: "VerbCommitTransactionResponse", implementor: "pdu.VerbCommitTransactionResponse"},
	VerbRollbackTransactionRequest:  {id: VerbRollbackTransactionRequest, name: "VerbRollbackTransactionRequest", implementor: "pdu.VerbRollbackTransactionRequest"},
	VerbRollbackTransactionResponse: {id: VerbRollbackTransactionResponse, name: "VerbRollbackTransactionResponse", implementor: "pdu.VerbRollbackTransactionResponse"},
	VerbQueryRequest:                {id: VerbQueryRequest, name: "VerbQueryRequest", implementor: "pdu.VerbQueryRequest"},
	VerbQueryResponse:               {id: VerbQueryResponse, name: "VerbQueryResponse", implementor: "pdu.VerbQueryResponse"},
	VerbTraverseRequest:             {id: VerbTraverseRequest, name: "VerbTraverseRequest", implementor: "pdu.VerbTraverseRequest"},
	VerbTraverseResponse:            {id: VerbTraverseResponse, name: "VerbTraverseResponse", implementor: "pdu.VerbTraverseResponse"},
	VerbAdminRequest:                {id: VerbAdminRequest, name: "VerbAdminRequest", implementor: "admin.VerbAdminRequest"},
	VerbAdminResponse:               {id: VerbAdminResponse, name: "VerbAdminResponse", implementor: "admin.VerbAdminResponse"},
	VerbMetadataRequest:             {id: VerbMetadataRequest, name: "VerbMetadataRequest", implementor: "pdu.VerbMetadataRequest"},
	VerbMetadataResponse:            {id: VerbMetadataResponse, name: "VerbMetadataResponse", implementor: "pdu.VerbMetadataResponse"},
	VerbGetEntityRequest:            {id: VerbGetEntityRequest, name: "VerbGetEntityRequest", implementor: "pdu.VerbGetEntityRequest"},    //0 = mean immediate, AttributeTypeInteger Max for indefinite
	VerbGetEntityResponse:           {id: VerbGetEntityResponse, name: "VerbGetEntityResponse", implementor: "pdu.VerbGetEntityResponse"}, //Represented in ms. Default Value is 10sec
	VerbGetLargeObjectRequest:       {id: VerbGetLargeObjectRequest, name: "VerbGetLargeObjectRequest", implementor: "pdu.VerbGetLargeObjectRequest"},
	VerbGetLargeObjectResponse:      {id: VerbGetLargeObjectResponse, name: "VerbGetLargeObjectResponse", implementor: "pdu.VerbGetLargeObjectResponse"},
	VerbDumpStacktraceRequest:       {id: VerbDumpStacktraceRequest, name: "VerbDumpStacktraceRequest", implementor: "pdu.VerbDumpStacktraceRequest"},
	VerbDisconnectChannelRequest:    {id: VerbDisconnectChannelRequest, name: "VerbDisconnectChannelRequest", implementor: "pdu.VerbDisconnectChannelRequest"},
	VerbSessionForcefullyTerminated: {id: VerbSessionForcefullyTerminated, name: "VerbSessionForcefullyTerminated", implementor: "pdu.VerbSessionForcefullyTerminated"},
	VerbDecryptBufferRequest:        {id: VerbDecryptBufferRequest, name: "VerbDecryptBufferRequest", implementor: "pdu.VerbDecryptBufferRequest"},
	VerbDecryptBufferResponse:       {id: VerbDecryptBufferResponse, name: "VerbDecryptBufferResponse", implementor: "pdu.VerbDecryptBufferResponse"},
	VerbExceptionMessage:            {id: VerbExceptionMessage, name: "VerbExceptionMessage", implementor: "pdu.VerbExceptionMessage"},
	VerbInvalidMessage:              {id: VerbInvalidMessage, name: "VerbInvalidMessage", implementor: "pdu.VerbInvalidMessage"},
}

func NewVerbId(id int, name, impl string) *CommandVerbs {
	newConfig := &CommandVerbs{id: id, name: name, implementor: impl}
	return newConfig
}

// Return the commandVerbs given its id
// @param id VerbId
// @return commandVerbs associated to the id
func GetVerb(id int) *CommandVerbs {
	verb, ok := PreDefinedVerbs[id]
	if ok {
		return &verb
	} else {
		invalid := PreDefinedVerbs[VerbInvalidMessage]
		return &invalid
	}
}

/////////////////////////////////////////////
// Helper Public functions for CommonVerbs //
/////////////////////////////////////////////

func (obj *CommandVerbs) GetID() int {
	return obj.id
}

func (obj *CommandVerbs) GetName() string {
	return obj.name
}

func (obj *CommandVerbs) GetImplementor() string {
	return obj.implementor
}


type ProtocolDataInputStream struct {
	mark                int	// Private Internal Marker - Not to be serialized / de-serialized
	Buf                 []byte
	BufLen              int
	iStreamCurPos       int
	iStreamEncoding     string
	iStreamReferenceMap map[int64]tgdb.TGEntity
}

/////////////////////////////////////////////////////////////////
// Helper functions for ProtocolDataInputStream
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataInputStream) ReadBoolean() (bool, tgdb.TGError) {
	//logger.Log(fmt.Sprintf(Entering ProtocolDataInputStream::ReadBoolean()"))
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadBoolean('%+v') as msg.CurPos >= msg.BufLen", false))
		errMsg := fmt.Sprint("End of data stream")
		return false, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := msg.Buf[msg.iStreamCurPos]
	msg.iStreamCurPos += 1
	if a == 0x00 {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadBoolean('%+v') as a == 0x00", false))
		}
		return false, nil
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadBoolean from the contents w/ ('%+v')", true))
	return true, nil
}

func (msg *ProtocolDataInputStream) ReadByte() (byte, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadByte()"))
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadByte('%+v') as msg.CurPos >= msg.BufLen", 0))
		errMsg := fmt.Sprint("End of data stream")
		return 0, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := msg.Buf[msg.iStreamCurPos]
	msg.iStreamCurPos += 1
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadByte from the contents w/ ('%+v')", a))
	return a, nil
}

func (msg *ProtocolDataInputStream) ReadChar() (string, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadChar()"))
	if (msg.iStreamCurPos + 2) > msg.BufLen {
		logger.Error(fmt.Sprintf("Returning ProtocolDataInputStream::ReadChar('%+v') as msg.CurPos+2 > msg.BufLen", ""))
		errMsg := fmt.Sprint("End of data stream")
		return "", GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int(msg.Buf[msg.iStreamCurPos]) << 8
	b := int(msg.Buf[msg.iStreamCurPos+1]) & 0xff
	msg.iStreamCurPos += 2
	c := fmt.Sprintf("%c", a+b)
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadChar from the contents w/ ('%+v')", c))
	return c, nil
}

func (msg *ProtocolDataInputStream) ReadDouble() (float64, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadDouble()"))
	a, err := msg.ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadDouble('%+v') as msg.ReadInt resulted in error", 0))
		return 0, err
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadDouble from the contents w/ ('%+v')", float64(a)))
	return math.Float64frombits(uint64(a)), nil
}

func (msg *ProtocolDataInputStream) ReadFloat() (float32, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadFloat()"))
	a, err := msg.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadFloat('%+v') as msg.ReadInt resulted in error", 0))
		return 0, err
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadFloat from the contents w/ ('%+v')", float32(a)))
	return math.Float32frombits(uint32(a)), nil
}

func (msg *ProtocolDataInputStream) ReadFully(b []byte) ([]byte, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadFully() from contents '%+v'", msg.Buf))
	if b == nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadFully as input byte array is EMPTY"))
		errMsg := fmt.Sprint("Input argument to ReadFully is null")
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return msg.ReadFullyAtPos(b, msg.iStreamCurPos, len(b))
}

func (msg *ProtocolDataInputStream) ReadFullyAtPos(b []byte, readCurPos, readLen int) ([]byte, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadFullyAtPos('%d') for '%d' from '%+v'", readCurPos, readLen, msg.Buf[msg.CurPos:(msg.CurPos+readLen)]))
	if readLen <= 0 {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadFullyAtPos('%+v') as readLen <= 0", readCurPos))
		errMsg := fmt.Sprint("Input argument representing how many bytes to read in ReadFullyAtPos is null")
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if (msg.iStreamCurPos + readLen) > msg.BufLen {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadFullyAtPos as (msg.CurPos + readLen) > msg.BufLen"))
		errMsg := fmt.Sprint("Input readLen to ReadFully is invalid")
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// Copy in input buffer b at offset for length
	copy(b, msg.Buf[msg.iStreamCurPos:(msg.iStreamCurPos +readLen)])
	//logger.Debug(fmt.Sprintf("Inside ProtocolDataInputStream::ReadFullyAtPos() temp / b is '%+v'", b))
	msg.iStreamCurPos = msg.iStreamCurPos + readLen
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadFullyAtPos('%d') for length '%d' from '%+v' in contents", readCurPos, readLen, b))
	return b, nil
}

func (msg *ProtocolDataInputStream) ReadInt() (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadInt() from Stream '%+v'", msg.String()))
	if (msg.iStreamCurPos + 4) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadInt as (msg.CurPos + 4) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int(msg.Buf[msg.iStreamCurPos]) << 24
	b := int(msg.Buf[msg.iStreamCurPos+1]) << 16 & 0x00ff0000
	c := int(msg.Buf[msg.iStreamCurPos+2]) << 8 & 0x0000ff00
	d := int(msg.Buf[msg.iStreamCurPos+3]) & 0xff
	msg.iStreamCurPos += 4
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadInt from the contents w/ ('%d')", a+b+c+d))
	return a + b + c + d, nil
}

func (msg *ProtocolDataInputStream) ReadLine() (string, tgdb.TGError) {
	// TODO: Revisit later - No-op For Now
	return "", nil
}

func (msg *ProtocolDataInputStream) ReadLong() (int64, tgdb.TGError) {
	//logger.Log(fmt.Sprint("Entering ProtocolDataInputStream::ReadLong()"))
	if (msg.iStreamCurPos + 8) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadLong as (msg.CurPos + 8) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a, err := msg.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadLong('%d') as msg.ReadInt resulted in error", 0))
		return 0, err
	}
	b, err := msg.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadLong('%d') as msg.ReadInt resulted in error", 0))
		return 0, err
	}
	c := int64((a << 32) + (b & 0xffffffff))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadLong from the contents w/ ('%d')", c))
	return c, nil
}

func (msg *ProtocolDataInputStream) ReadShort() (int16, tgdb.TGError) {
	//logger.Log(fmt.Sprint(Entering ProtocolDataInputStream::ReadShort()"))
	if (msg.iStreamCurPos + 2) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadShort as (msg.CurPos + 2) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int16(msg.Buf[msg.iStreamCurPos]) << 8
	b := int16(msg.Buf[msg.iStreamCurPos+1]) & 0xff
	msg.iStreamCurPos += 2
	c := a + b
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadShort from the contents w/ ('%d')", c))
	return c, nil
}

func (msg *ProtocolDataInputStream) ReadUnsignedByte() (int, tgdb.TGError) {
	//logger.Log(fmt.Sprint("Entering ProtocolDataInputStream::ReadUnsignedByte()"))
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUnsignedByte as msg.CurPos > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int(msg.Buf[msg.iStreamCurPos]) & 0xff
	msg.iStreamCurPos += 1
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadUnsignedByte from the contents w/ ('%d')", a))
	return a, nil
}

func (msg *ProtocolDataInputStream) ReadUnsignedShort() (uint16, tgdb.TGError) {
	//logger.Log(fmt.Sprint("Entering ProtocolDataInputStream::ReadUnsignedShort()"))
	if (msg.iStreamCurPos + 2) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUnsignedShort as (msg.CurPos + 2) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return 0, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int(msg.Buf[msg.iStreamCurPos]) << 8
	b := int(msg.Buf[msg.iStreamCurPos+1]) & 0xff
	msg.iStreamCurPos += 2
	c := uint16((a + b) & 0x0000ffff)
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadUnsignedShort from the contents w/ ('%d')", c))
	return c, nil
}

func (msg *ProtocolDataInputStream) ReadUTF() (string, tgdb.TGError) {
	//logger.Log(fmt.Sprint("Entering ProtocolDataInputStream::ReadUTF()"))
	start := msg.iStreamCurPos

	utfLen, err := msg.ReadUnsignedShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTF as msg.ReadUnsignedShort resulted in error"))
		// Reset the current position
		msg.iStreamCurPos = start
		return "", err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside ProtocolDataInputStream::ReadUTF() utfLength to read: '%d'", utfLen))
	}
	str, err := msg.ReadUTFString(int(utfLen))
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadUTF from the contents w/ ('%s') & error ('%+v')", str, err))
	}
	return str, err
}

func (msg *ProtocolDataInputStream) ReadUTFString(utfLen int) (string, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadUTFString('%d') CurPos: %d BufLen: %d Buf: %+v", utfLen, msg.CurPos, msg.BufLen, msg.Buf))
	if (msg.iStreamCurPos + utfLen) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as (msg.CurPos + utfLen) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return "", GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	charBuf := ""
	start := msg.iStreamCurPos
	lastPos := msg.iStreamCurPos + utfLen
	strLen := 0
	var c, cVal, char2, char3 int

	// Loop through message buffer
	for {
		if !(msg.iStreamCurPos < lastPos) {
			//logger.Log(fmt.Sprintf("ProtocolDataInputStream::ReadUTFString - Breaking Loop @msg.CurPos '%d' lastPos '%d' c '%d' cVal '%d' char2 '%d' char3 '%d'", msg.CurPos, lastPos, c, cVal, char2, char3))
			break
		}
		c = int(msg.Buf[msg.iStreamCurPos]) & 0xff
		//logger.Debug(fmt.Sprintf("ProtocolDataInputStream::ReadUTFString - Inside Loop msg.CurPos '%d' lastPos '%d' c '%d' cVal '%d' char2 '%d' char3 '%d'", msg.CurPos, lastPos, c, cVal, char2, char3))
		msg.iStreamCurPos += 1
		cVal = c >> 4
		if cVal <= 7 { // 0xxxxxxx
			charBuf = charBuf + string(c)
			strLen++
		} else if cVal == 12 || cVal == 13 { // 110x xxxx   10xx xxxx
			if msg.iStreamCurPos+1 > lastPos {
				msg.iStreamCurPos = start
				logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as msg.CurPos+1 > lastPos"))
				errMsg := fmt.Sprint("Data Format Issue")
				return "", GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
			}
			char2 = int(msg.Buf[msg.iStreamCurPos])
			msg.iStreamCurPos += 1
			if (char2 & 0xC0) != 0x80 {
				msg.iStreamCurPos = start
				logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as (char2 & 0xC0) != 0x80"))
				errMsg := fmt.Sprint("Data Format Issue")
				return "", GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
			}
			c1 := ((c & 0x1F) << 6) | (char2 & 0x3F)
			charBuf = charBuf + string(c1)
			strLen++
		} else if cVal == 14 { // 1110 xxxx  10xx xxxx  10xx xxxx
			if msg.iStreamCurPos+2 > lastPos {
				msg.iStreamCurPos = start
				logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as msg.CurPos+2 > lastPos"))
				errMsg := fmt.Sprint("Data Format Issue")
				return "", GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
			}
			char2 = int(msg.Buf[msg.iStreamCurPos])
			msg.iStreamCurPos += 1
			char3 = int(msg.Buf[msg.iStreamCurPos])
			msg.iStreamCurPos += 1
			if ((char2 & 0xC0) != 0x80) || ((char3 & 0xC0) != 0x80) {
				msg.iStreamCurPos = start
				logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as ((char2 & 0xC0) != 0x80) || ((char3 & 0xC0) != 0x80)"))
				errMsg := fmt.Sprint("Data Format Issue")
				return "", GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
			}
			c1 := ((c & 0x0F) << 12) | ((char2 & 0x3F) << 6) | ((char3 & 0x3F) << 0)
			charBuf = charBuf + string(c1)
			strLen++
		} else { // 10xx xxxx,  1111 xxxx
			msg.iStreamCurPos = start
			logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString for 10xxxxxx, 1111xxxx"))
			errMsg := fmt.Sprint("Data Format Issue")
			return "", GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
		}
	}

	if msg.iStreamCurPos != lastPos {
		msg.iStreamCurPos = start
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as msg.CurPos != lastPos"))
		errMsg := fmt.Sprint("Data Format Issue")
		return "", GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadUTFString from the contents w/ ('%+v')", charBuf))
	return charBuf, nil
}

func (msg *ProtocolDataInputStream) SkipBytes(n int) (int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering ProtocolDataInputStream::SkipBytes('%d'), contents are '%+v'", n, msg.Buf))
	}
	if n < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::SkipBytes as n < 0"))
		errMsg := fmt.Sprint("Invalid bytes to skip")
		return 0, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	avail := msg.BufLen - msg.iStreamCurPos
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside ProtocolDataInputStream::SkipBytes - Available bytes ('%d') compared w/ To-Be-Skipped '%d'", avail, n))
	}
	if avail <= n {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::SkipBytes as avail <= n from the contents w/ ('%d')", avail))
		}
		msg.iStreamCurPos = msg.BufLen
		return avail, nil
	}
	msg.iStreamCurPos = msg.iStreamCurPos + n
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::SkipBytes from the contents w/ ('%d')", n))
	}
	return n, nil
}

func (msg *ProtocolDataInputStream) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ProtocolDataInputStream:{")
	buffer.WriteString(fmt.Sprintf("Buffer Length: %d", msg.BufLen))
	buffer.WriteString(fmt.Sprintf(", Current Position: %d", msg.iStreamCurPos))
	buffer.WriteString(fmt.Sprintf(", Encoding: %s", msg.iStreamEncoding))
	buffer.WriteString(fmt.Sprintf(", Buffer: %s", bytes.NewBuffer(msg.Buf).String()))
	buffer.WriteString(fmt.Sprintf(", Reference Map: %+v", msg.iStreamReferenceMap))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGInputStream
/////////////////////////////////////////////////////////////////

// Available checks whether there is any data available on the stream to read
func (msg *ProtocolDataInputStream) Available() (int, tgdb.TGError) {
	//if msg.BufLen == 0 || (msg.BufLen-msg.CurPos) < 0 {
	if (msg.BufLen-msg.iStreamCurPos) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::Available as (msg.BufLen-msg.CurPos) < 0"))
		errMsg := fmt.Sprint("Invalid data length of Protocol Data Input Stream")
		return 0, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	return msg.BufLen - msg.iStreamCurPos, nil
}

// GetPosition gets the current position of internal cursor
func (msg *ProtocolDataInputStream) GetPosition() int64 {
	return int64(msg.iStreamCurPos)
}

// GetReferenceMap returns a user maintained reference map
func (msg *ProtocolDataInputStream) GetReferenceMap() map[int64]tgdb.TGEntity {
	return msg.iStreamReferenceMap
}

// Mark marks the current position
func (msg *ProtocolDataInputStream) Mark(readlimit int) {
	msg.mark = readlimit
}

// MarkSupported checks whether the marking is supported or not
func (msg *ProtocolDataInputStream) MarkSupported() bool {
	return true
}

// Read reads the current byte
func (msg *ProtocolDataInputStream) Read() (int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering ProtocolDataInputStream::Read(), contents are '%+v'", msg.Buf))
	}
	if msg.iStreamCurPos > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::Read as msg.CurPos > msg.BufLen"))
		return -1, nil
	}
	if msg.iStreamCurPos < msg.BufLen {
		b := int(msg.Buf[msg.iStreamCurPos] & 0xff)
		msg.iStreamCurPos += 1
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::Read from the contents w/ ('%d')", b))
		}
		return b, nil
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::Read from the contents w/ ('%d')", 0))
	}
	return 0, nil
}

// ReadIntoBuffer copies bytes in specified buffer
// The buffer cannot be NIL
func (msg *ProtocolDataInputStream) ReadIntoBuffer(b []byte) (int, tgdb.TGError) {
	return msg.ReadAtOffset(b, 0, len(b))
}

// ReadAtOffset is similar to readFully.
func (msg *ProtocolDataInputStream) ReadAtOffset(b []byte, off int, length int) (int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering ProtocolDataInputStream::ReadAtOffset('%d') for length '%d' from '%+v', contents are '%+v'", off, length, b, msg.Buf))
	}
	if b == nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadAtOffset as b == nil"))
		return -1, CreateExceptionByType(TGErrorIOException)
	}
	if off < 0 || off > len(b) || length < 0 || (off+length) > len(b) || (off+length) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadAtOffset as invalid input data / offset / length OR data is corrupt"))
		errMsg := fmt.Sprint("Invalid input data / offset / length OR data is corrupt")
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if msg.iStreamCurPos >= msg.BufLen {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadAtOffset('%d') for length '%d' from the contents w/ ('%d')", off, -1, b))
		}
		return -1, nil
	}
	if (msg.iStreamCurPos + length) > msg.BufLen {
		length = msg.BufLen - msg.iStreamCurPos
	}
	if length <= 0 {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadAtOffset('%d') for length '%d' from the contents w/ ('%d')", off, 0, b))
		}
		return 0, nil
	}
	// Copy from input buffer b at offset for length
	// Create a slice of input buffer
	temp := b[off:] // Slice starting at offset 'off' of input buffer 'b'
	tempLen := len(temp) - length
	temp = temp[:tempLen]
	copy(temp, msg.Buf)
	msg.iStreamCurPos = msg.iStreamCurPos + length
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadAtOffset('%d') for length '%d' from the contents w/ ('%d')", off, length, b))
	}
	return length, nil
}

// ReadBytes reads an encoded byte array. writeBytes encodes the length, and the byte[].
// This is equivalent to do a readInt, and read(byte[])
func (msg *ProtocolDataInputStream) ReadBytes() ([]byte, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering ProtocolDataInputStream::ReadBytes(), contents are '%+v'", msg.Buf))
	}
	emptyByteArray := make([]byte, 0)
	length, err := msg.ReadInt()
	if length == 0 {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning ProtocolDataInputStream::ReadBytes as len == 0"))
		}
		return emptyByteArray, nil
	}
	if length < -1 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadBytes as len == -1"))
		errMsg := fmt.Sprint("Read data corrupt")
		return emptyByteArray, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	buf := make([]byte, length)
	_, err = msg.ReadIntoBuffer(buf)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadBytes w/ error in reading into buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadBytes from the contents w/ ('%+v')", buf))
	}
	return buf, nil
}

// ReadVarLong reads a Variable long field
func (msg *ProtocolDataInputStream) ReadVarLong() (int64, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering ProtocolDataInputStream::ReadVarLong(), contents are '%+v'", msg.Buf))
	}
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadVarLong as msg.CurPos >= msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	curByte := msg.Buf[msg.iStreamCurPos]
	if curByte == U64PACKED_NULL {
		msg.iStreamCurPos++
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", U64_NULL))
		}
		return U64_NULL, nil
	}

	if curByte&0x80 == 0 {
		msg.iStreamCurPos++
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", int64(curByte)))
		}
		return int64(curByte), nil
	}

	if curByte&0x40 == 0 {
		if (msg.iStreamCurPos + 2) > msg.BufLen {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadVarLong as (msg.CurPos + 2) > msg.BufLen"))
			errMsg := fmt.Sprint("End of data stream")
			return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
		}
		a := msg.Buf[msg.iStreamCurPos] << 8
		msg.iStreamCurPos++
		b := msg.Buf[msg.iStreamCurPos] & 0xff
		msg.iStreamCurPos++
		c := int16(a + b)
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", int64(c&0x3fff)))
		}
		return int64(c & 0x3fff), nil
	}

	if curByte&0x20 == 0 {
		if (msg.iStreamCurPos + 4) > msg.BufLen {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadVarLong as (msg.CurPos + 4) > msg.BufLen"))
			errMsg := fmt.Sprint("End of data stream")
			return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
		}
		a := int(msg.Buf[msg.iStreamCurPos] << 24)
		msg.iStreamCurPos++
		b := int(msg.Buf[msg.iStreamCurPos]<<16) & 0x00ff0000
		msg.iStreamCurPos++
		c := int(msg.Buf[msg.iStreamCurPos]<<8) & 0x0000ff00
		msg.iStreamCurPos++
		d := int((msg.Buf[msg.iStreamCurPos]) & 0xff)
		msg.iStreamCurPos++
		lValue := (a + b + c + d) & 0x1fffffff
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", int64(lValue)))
		}
		return int64(lValue), nil
	}

	count := int(curByte & 0x0f)
	msg.iStreamCurPos++
	if (msg.iStreamCurPos + count) > msg.BufLen {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadVarLong as (msg.CurPos + count) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	lValue := 0
	for i := 0; i < count; i++ {
		lValue <<= 8
		lValue |= int(msg.Buf[msg.iStreamCurPos] & 0x00ff)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", int64(lValue)))
	}
	return int64(lValue), nil
}

// Reset brings internal moving cursor back to the old position
func (msg *ProtocolDataInputStream) Reset() {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering ProtocolDataInputStream::Reset( Marker = '%d') in the contents", msg.mark))
	}
	msg.iStreamCurPos = msg.mark
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::Reset( Marker = '%d'), contents are '%+v'", msg.mark, msg.Buf))
	}
}

// SetPosition sets the position of reading.
func (msg *ProtocolDataInputStream) SetPosition(position int64) int64 {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering ProtocolDataInputStream::SetPosition('%d') in the contents", position))
	}
	oldPos := msg.iStreamCurPos
	msg.iStreamCurPos = int(position)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::SetPosition('%d'), contents are '%+v'", int64(oldPos), msg.Buf))
	}
	return int64(oldPos)
}

// SetReferenceMap sets a user maintained map as reference data
func (msg *ProtocolDataInputStream) SetReferenceMap(rMap map[int64]tgdb.TGEntity) {
	msg.iStreamReferenceMap = rMap
}

// Skip skips n bytes
func (msg *ProtocolDataInputStream) Skip(n int64) (int64, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering ProtocolDataInputStream::Skip('%d') in the contents", n))
	}
	a, err := msg.SkipBytes(int(n))
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataInputStream::Skip('%d'), contents are '%+v'", int64(a), msg.Buf))
	}
	return int64(a), err
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataInputStream) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.Buf, msg.BufLen, msg.iStreamCurPos, msg.iStreamEncoding, msg.iStreamReferenceMap)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}


var writeBuf = ""

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataInputStream) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.Buf, &msg.BufLen, &msg.iStreamCurPos, &msg.iStreamEncoding, &msg.iStreamReferenceMap)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type ProtocolDataOutputStream struct {
	Buf              []byte
	oStreamBufLen    int
	oStreamByteCount int
}


/////////////////////////////////////////////////////////////////
// Helper functions for ProtocolDataInputStream
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataOutputStream) Ensure(len int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::Ensure('%d') in contents", len))
	if (msg.oStreamByteCount + len) <= msg.oStreamBufLen {
		//logger.Log(fmt.Sprint("Returning ProtocolDataOutputStream::Ensure as (msg.Count + len) <= msg.BufLen"))
		return
	}
	newLen := 0
	if len > 100000 {
		newLen = msg.oStreamByteCount + len + 2048
	} else {
		newLen = (msg.oStreamByteCount + len) * 2
	}
	b := make([]byte, newLen)
	copy(b, msg.Buf[0:msg.oStreamByteCount])
	msg.Buf = b
	msg.oStreamBufLen = newLen
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::Ensure('%d') in contents '%+v'", len, msg.Buf))
}

func (msg *ProtocolDataOutputStream) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ProtocolDataOutputStream:{")
	buffer.WriteString(fmt.Sprintf("Buffer Length: %d", msg.oStreamBufLen))
	buffer.WriteString(fmt.Sprintf(", Count: %d", msg.oStreamByteCount))
	buffer.WriteString(fmt.Sprintf(", Buffer: "))
	strArray := []string{buffer.String(), bytes.NewBuffer(msg.Buf).String()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

func (msg *ProtocolDataOutputStream) WriteBoolean(value bool) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBoolean('%+v') in contents", value))
	if msg.oStreamByteCount >= msg.oStreamBufLen {
		msg.Ensure(1)
	}
	if value {
		msg.Buf[msg.oStreamByteCount] = byte(1)
	} else {
		msg.Buf[msg.oStreamByteCount] = byte(0)
	}
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBoolean('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteByte(value int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteByte('%+v') in contents", value))
	if msg.oStreamByteCount >= msg.oStreamBufLen {
		msg.Ensure(1)
	}
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteByte('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteBytesFromString(value string) tgdb.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBytesFromString('%s') in contents", value))
	buf := []byte(value)
	return msg.WriteBytes(buf)
}

func (msg *ProtocolDataOutputStream) WriteChar(value int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteChar('%+v') in contents", value))
	if (msg.oStreamByteCount + 2) >= msg.oStreamBufLen {
		msg.Ensure(2)
	}
	msg.Buf[msg.oStreamByteCount] = byte(value >> 8)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteChar('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteChars(value string) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteChars('%+v') in contents", value))
	sLen := len(value)
	msg.Ensure(sLen)
	for i := 0 ; i < sLen ; i++ {
		msg.WriteChar(int(value[i]))
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteChars('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteDouble(value float64) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteDouble('%+v') in contents", value))
	msg.WriteLong(int64(math.Float64bits(value)))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteDouble('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteFloat(value float32) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteFloat('%+v') in contents", value))
	msg.WriteInt(int(math.Float32bits(value)))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteFloat('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteInt(value int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteInt('%+v') in contents", value))
	if (msg.oStreamByteCount + 4) >= msg.oStreamBufLen {
		msg.Ensure(4)
	}
	//intToBytes(value, msg.Buf, msg.Count)
	msg.Buf[msg.oStreamByteCount] = byte(value >> 24)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value >> 16)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value >> 8)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteInt('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteLong(value int64) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteLong('%+v') in contents", value))
	if (msg.oStreamByteCount + 8) >= msg.oStreamBufLen {
		msg.Ensure(8)
	}
	// Splitting in into two ints, is much faster than shifting bits for a long i.e (byte)val >> 56, (byte)val >> 48 ..
	a := int(value >> 32)
	b := int(value)
	//intToBytes(a, msg.Buf, msg.Count)
	msg.Buf[msg.oStreamByteCount] = byte(a >> 24)
	msg.Buf[msg.oStreamByteCount+1] = byte(a >> 16)
	msg.Buf[msg.oStreamByteCount+2] = byte(a >> 8)
	msg.Buf[msg.oStreamByteCount+3] = byte(a)
	//intToBytes(b, msg.Buf, msg.Count+4)
	msg.Buf[msg.oStreamByteCount+4] = byte(b >> 24)
	msg.Buf[msg.oStreamByteCount+5] = byte(b >> 16)
	msg.Buf[msg.oStreamByteCount+6] = byte(b >> 8)
	msg.Buf[msg.oStreamByteCount+7] = byte(b)
	msg.oStreamByteCount += 8
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteLong('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteShort(value int)  {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteShort('%+v') in contents", value))
	if (msg.oStreamByteCount + 2) >= msg.oStreamBufLen {
		msg.Ensure(2)
	}
	msg.Buf[msg.oStreamByteCount] = byte(value >> 8)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteShort('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteBytesFromPos(value []byte, writePos, writeLen int) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBytesFromPos('%d') for len '%d' from '%+v' in contents", writePos, writeLen, value))
	}
	if value == nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBytesFromPos - invalid byte value specified to be written at pos %d", writePos))
		errMsg := fmt.Sprintf("Invalid byte value specified to be written at pos %d", writePos)
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if writePos < 0 || writePos > len(value) || writeLen < 0 || (writePos + writeLen) > len(value) || (writePos + writeLen) < 0 {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBytesFromPos - either invalid pos %d or write length %d specified", writeLen, writePos))
		errMsg := fmt.Sprintf("Either invalid pos %d or write length %d specified", writeLen, writePos)
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if writeLen == 0 {
		return nil
	}
	msg.Ensure(writeLen)
	tempBuf := value[writePos:]
	copy(msg.Buf[msg.oStreamByteCount:], tempBuf[:writeLen])
	msg.oStreamByteCount += writeLen
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBytesFromPos(''%d') for len '%d' from '%+v' in contents '%+v'", writePos, writeLen, value, msg.Buf))
	}
	return nil
}

func (msg *ProtocolDataOutputStream) WriteUTF(str string) tgdb.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteUTF('%+v') in contents", str))
	start := msg.oStreamByteCount
	sLen := 0

	if msg.oStreamByteCount+ 2 > msg.oStreamBufLen {
		msg.Ensure(2 + len(str) * 3)
	}
	msg.oStreamByteCount += 2

	sLen, err := msg.WriteUTFString(str)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteUTF - UTF Data Format Issue"))
		msg.oStreamByteCount = start
		errMsg := fmt.Sprint("UTF Data Format Issue")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
	}

	// Now write length
	msg.Buf[start]   = byte((sLen >> 8) & 0xFF)
	msg.Buf[start+1] = byte(sLen & 0xFF)
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteUTF('%+v') in contents '%+v'", str, msg.Buf))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGInputStream
/////////////////////////////////////////////////////////////////

// GetBuffer gets the underlying Buffer
func (msg *ProtocolDataOutputStream) GetBuffer() []byte {
	return msg.Buf
}

// GetLength gets the total write length
func (msg *ProtocolDataOutputStream) GetLength() int {
	return msg.oStreamByteCount
}

// GetPosition gets the current write position
func (msg *ProtocolDataOutputStream) GetPosition() int {
	return msg.oStreamByteCount
}

// SkipNBytes skips n bytes. Allocate if necessary
func (msg *ProtocolDataOutputStream) SkipNBytes(n int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::SkipNBytes('%d') in contents", n))
	if msg.oStreamByteCount+ n > msg.oStreamBufLen {
		msg.Ensure(n)
	}
	msg.oStreamByteCount += n
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataOutputStream::SkipNBytes('%d') in contents are '%+v'", n, msg))
	}
}

// ToByteArray returns a new constructed byte array of the data that is being streamed.
func (msg *ProtocolDataOutputStream) ToByteArray() ([]byte, tgdb.TGError) {
	buf := make([]byte, msg.oStreamByteCount)
	_, err := io.Copy(bytes.NewBuffer(buf), bytes.NewBuffer(msg.Buf))
	if err != nil {
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, err.Error(), "")
	}
	return buf, nil
}

// WriteBooleanAt writes boolean at a given position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteBooleanAt(pos int, value bool) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBooleanAt(Pos '%d' - '%+v') in contents", pos, value))
	if pos >= msg.oStreamByteCount {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBooleanAt - Invalid position '%d' specified", pos))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if value {
		msg.Buf[pos] = byte(1)
	} else {
		msg.Buf[pos] = byte(0)
	}
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBooleanAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteByteAt writes a byte at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteByteAt(pos int, value int) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteByteAt(Pos '%d' - '%+v') in contents", pos, value))
	if pos >= msg.oStreamByteCount {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteByteAt - Invalid position '%d' specified", pos))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Buf[pos] = byte(value)
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteByteAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteBytes writes the len, and the byte array into the buffer
func (msg *ProtocolDataOutputStream) WriteBytes(buf []byte) tgdb.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBytes('%+v') in contents", buf))
	bLen := len(buf)
	msg.Ensure(bLen+4)
	msg.WriteInt(bLen)

	for i := 0; i < bLen ; i++ {
		msg.WriteByte(int(buf[i]))
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBytes('%+v') in contents are '%+v'", buf, msg.Buf))
	return nil
}

// WriteBytesAt writes string at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteBytesAt(pos int, s string) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBytesAt(Pos '%d' - '%+v') in contents", pos, s))
	sLen := len([]rune(s))
	if (pos + sLen) >= msg.oStreamByteCount {
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	for i := 0; i < sLen ; i++ {
		v, err := fmt.Scanf("%d", s[i])
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBytesAt - Unable to get character at position %d from string %s", i, s))
			errMsg := fmt.Sprintf("Unable to get character at position %d from string %s", i, s)
			return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		pos, err := msg.WriteByteAt(pos, v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBytesAt - Invalid position '%d' specified", pos))
			errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
			return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBytesAt(Pos '%d' - '%+v') in contents are '%+v'", pos, s, msg.Buf))
	return pos, nil
}

// WriteCharAt writes a Java Char at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteCharAt(pos int, value int) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteCharAt(Pos '%d' - '%+v') in contents", pos, value))
	if (pos + 2) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteCharAt as (pos + 2) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Buf[pos] = byte(value >> 8)
	pos++
	msg.Buf[pos] = byte(value)
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteCharAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteCharsAt writes Chars at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteCharsAt(pos int, s string) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteCharsAt(Pos '%d' - '%+v') in contents", pos, s))
	sLen := len([]rune(s))
	if (pos + sLen) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteCharsAt as (pos + sLen) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	for i := 0; i < sLen ; i++ {
		v, err := fmt.Scanf("%d", s[i])
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteCharsAt - Unable to get character at position %d from string %s", i, s))
			errMsg := fmt.Sprintf("Unable to get character at position %d from string %s", i, s)
			return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
		}
		pos, err := msg.WriteCharAt(pos, v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteCharsAt - Invalid position '%d' specified", pos))
			errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
			return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
		}
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteCharsAt(Pos '%d' - '%+v') in contents are '%+v'", pos, s, msg.Buf))
	return pos, nil
}

// WriteDoubleAt writes Double at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteDoubleAt(pos int, value float64) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteDoubleAt(Pos '%d' - '%+v') in contents", pos, value))
	v, err := msg.WriteLongAt(pos, int64(value))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteDoubleAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return v, err
}

// WriteFloatAt writes Float at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteFloatAt(pos int, value float32) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteFloatAt(Pos '%d' - '%+v') in contents", pos, value))
	v, err := msg.WriteIntAt(pos, int(value))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteFloatAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return v, err
}

// WriteIntAt writes Integer at the position.Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteIntAt(pos int, value int) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteIntAt(Pos '%d' - '%+v') in contents", pos, value))
	if (pos + 4) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteIntAt as (pos + 4) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Buf[pos] = byte(value >> 24)
	pos++
	msg.Buf[pos] = byte(value >> 16)
	pos++
	msg.Buf[pos] = byte(value >> 8)
	pos++
	msg.Buf[pos] = byte(value)
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteIntAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteLongAsBytes writes Long in byte format
func (msg *ProtocolDataOutputStream) WriteLongAsBytes(value int64) tgdb.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteLongAsBytes(Pos '%d' - '%+v') in contents", pos, value))
	if (msg.oStreamByteCount + 8) > len(msg.Buf) {
		msg.Ensure(8)
	}
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.Buf[msg.oStreamByteCount+1] = byte(value >> 8)
	msg.Buf[msg.oStreamByteCount+2] = byte(value >> 16)
	msg.Buf[msg.oStreamByteCount+3] = byte(value >> 24)
	msg.Buf[msg.oStreamByteCount+4] = byte(value >> 32)
	msg.Buf[msg.oStreamByteCount+5] = byte(value >> 40)
	msg.Buf[msg.oStreamByteCount+6] = byte(value >> 48)
	msg.Buf[msg.oStreamByteCount+7] = byte(value >> 56)
	msg.oStreamByteCount += 8

	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteLongAsBytes(Pos '%d' - '%+v') in contents are '%+v'", pos+8, value, msg.Buf)
	return nil
}

// WriteLongAt writes Long at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteLongAt(pos int, value int64) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteLongAt(Pos '%d' - '%+v') in contents", pos, value))
	if (pos + 8) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteLongAt as (pos + 8) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// Splitting in into two ints, is much faster than shifting bits for a long i.e (byte)val >> 56, (byte)val >> 48 ..
	a := int(value >> 32)
	b := int(value)

	msg.Buf[pos] = byte(a >> 24)
	msg.Buf[pos+1] = byte(a >> 16)
	msg.Buf[pos+2] = byte(a >> 8)
	msg.Buf[pos+3] = byte(a)
	msg.Buf[pos+4] = byte(b >> 24)
	msg.Buf[pos+5] = byte(b >> 16)
	msg.Buf[pos+6] = byte(b >> 8)
	msg.Buf[pos+7] = byte(b)

	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteLongAt(Pos '%d' - '%+v') in contents are '%+v'", pos+8, value, msg.Buf)
	return pos+8, nil
}

// WriteShortAt writes a Java Char at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteShortAt(pos int, value int) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteShortAt(Pos '%d' - '%+v') in contents", pos, value))
	if (pos + 2) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteShortAt as (pos + 2) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Buf[pos] = byte(value >> 8)
	pos++
	msg.Buf[pos] = byte(value)
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteShortAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteUTFString writes UTFString
func (msg *ProtocolDataOutputStream) WriteUTFString(str string) (int, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteUTFString('%+v') in contents", str))
	start := msg.oStreamByteCount
	sLen := len(str)
	i := 0
	c := 0

	if (msg.oStreamByteCount + 3 * sLen) > msg.oStreamBufLen {
		msg.Ensure(len(str) * 3)
	}
	writeBuf = str[0:sLen]

	for {
		if ! (i < sLen) {
			//logger.Debug(fmt.Sprintf("ProtocolDataOutputStream::WriteUTFString - Breaking loop i='%d', c='%d', msg.Count='%d', msg.Buf[msg.Count]='%d'", i, c, msg.Count, msg.Buf[msg.Count]))
			break
		}
		c = int(writeBuf[i])
		i++
		if (c >= 0x0001) && (c <= 0x007F) {
			msg.Buf[msg.oStreamByteCount] = byte(c)
			msg.oStreamByteCount++
		} else if c > 0x07FF {
			msg.Buf[msg.oStreamByteCount] = byte(0xE0 | ((c >> 12) & 0x0F))
			msg.oStreamByteCount++
			msg.Buf[msg.oStreamByteCount] = byte(0x80 | ((c >> 6) & 0x3F))
			msg.oStreamByteCount++
			msg.Buf[msg.oStreamByteCount] = byte(0x80 | ((c >> 0) & 0x3F))
			msg.oStreamByteCount++
		} else {
			msg.Buf[msg.oStreamByteCount] = byte(0xC0 | ((c >> 6) & 0x1F))
			msg.oStreamByteCount++
			msg.Buf[msg.oStreamByteCount] = byte(0x80 | ((c >> 0) & 0x3F))
			msg.oStreamByteCount++
		}
		//logger.Debug(fmt.Sprintf("ProtocolDataOutputStream::WriteUTFString - Inside loop i='%d', c='%d', msg.Count='%d', msg.Buf[msg.Count]='%d'", i, c, msg.Count, msg.Buf[msg.Count]))
	}

	writtenLen := msg.oStreamByteCount - start
	if writtenLen > 65535 {
		msg.oStreamByteCount = start
		errMsg := fmt.Sprint("Input String is too long")
		return -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteUTFString('%+v') in contents are '%+v'", str, msg.Buf))
	return writtenLen, nil
}

// WriteVarLong writes a long value as varying length into the buffer.
func (msg *ProtocolDataOutputStream) WriteVarLong(value int64) tgdb.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteVarLong('%+v') in contents", value))
	if value < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteVarLong as value < 0"))
		errMsg := fmt.Sprint("Can not pack negative long value")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if value == U64_NULL {
		if msg.oStreamByteCount >= msg.oStreamBufLen {
			msg.Ensure(1)
		}
		msg.Buf[msg.oStreamByteCount] = U64PACKED_NULL
		msg.oStreamByteCount++
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
		}
		return nil
	}

	if value <= 0x7f {
		if msg.oStreamByteCount >= msg.oStreamBufLen {
			msg.Ensure(1)
		}
		msg.Buf[msg.oStreamByteCount] = byte(value)
		msg.oStreamByteCount++
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
		}
		return nil
	}

	if value <= 0x3fff {
		if (msg.oStreamByteCount + 2) >= msg.oStreamBufLen {
			msg.Ensure(2)
		}
		value |= 0x00008000
		msg.Buf[msg.oStreamByteCount] = byte(value >> 8)
		msg.oStreamByteCount++
		msg.Buf[msg.oStreamByteCount] = byte(value)
		msg.oStreamByteCount++
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
		}
		return nil
	}

	if value <= 0x1fffffff {
		if (msg.oStreamByteCount + 4) >= msg.oStreamBufLen {
			msg.Ensure(4)
		}
		value |= 0xC0000000
		msg.oStreamByteCount += 4
		for i := 1; i < 4; i++ {
			msg.Buf[msg.oStreamByteCount- 1] = byte(value)
			value = value >> 8
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
		}
		return nil
	}

	// We may need up to 9 bytes
	if (msg.oStreamByteCount + 9) > msg.oStreamBufLen {
		msg.Ensure(9)
	}

	// Calculate the number of non-zero bytes
	mask := 0xff0000000000000
	cnt := 8
	for i := 0; i < 8; i++ {
		if (value & int64(mask)) != 0 {
			break
		}
		cnt--
		mask = mask >> 8
	}

	b := byte(cnt | 0xE0)
	msg.Buf[msg.oStreamByteCount] = b
	msg.oStreamByteCount++

	cnt += cnt
	for i:=1; i<=cnt; i++ {
		msg.Buf[cnt - i] = byte(value)
		value = value >> 8
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataOutputStream) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.Buf, msg.oStreamBufLen, msg.oStreamByteCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataOutputStream) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.Buf, &msg.oStreamBufLen, &msg.oStreamByteCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


func DefaultProtocolDataInputStream() *ProtocolDataInputStream {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ProtocolDataInputStream{})

	newStream := ProtocolDataInputStream{
		mark:                0,
		Buf:                 make([]byte, 0),
		BufLen:              0,
		iStreamCurPos:       0,
		iStreamReferenceMap: make(map[int64]tgdb.TGEntity, 0),
	}
	return &newStream
}

// Create New Input Stream Instance
func NewProtocolDataInputStream(buf []byte) *ProtocolDataInputStream {
	newStream := DefaultProtocolDataInputStream()
	newStream.Buf = buf
	newStream.BufLen = len(buf)
	return newStream
}



func DefaultProtocolDataOutputStream() *ProtocolDataOutputStream {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ProtocolDataOutputStream{})

	newStream := ProtocolDataOutputStream{
		Buf:              make([]byte, 256),
		oStreamBufLen:    256,
		oStreamByteCount: 0,
	}
	return &newStream
}

// Create New Input Stream Instance
func NewProtocolDataOutputStream(len int) *ProtocolDataOutputStream {
	newStream := DefaultProtocolDataOutputStream()
	newStream.Buf = make([]byte, len)
	newStream.oStreamBufLen = len
	return newStream
}

/////////////////////////////////////////////////////////////////
// Private functions for ProtocolDataInputStream
/////////////////////////////////////////////////////////////////

/*
// Already defined in another file (rf_channelimpl.go)
func intToBytes(value int, bytes []byte, offset int) {
	for i:=0; i<4; i++ {
		bytes[offset+i] = byte(value >> uint(8*(3-i)))
	}
}
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
	tenantId int
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
		tenantId: 0,
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

func (msg *AbstractProtocolMessage) GetTenantId() int{
	return msg.tenantId
}

func (msg *AbstractProtocolMessage) SetTenantId(id int) {
	msg.tenantId = id
}


func (msg *AbstractProtocolMessage) SetSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	if !(msg.GetIsUpdatable() || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
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

func APMReadHeader(msg tgdb.TGMessage, is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractProtocolMessage:APMReadHeader"))
	}
	// First member attribute / element of message header is BufLength
	// It has already been read before reaching here - in FromBytes()

	magic, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading magicId from message buffer"))
		return err
	}
	if magic != GetMagic() {
		errMsg := fmt.Sprint("Bad Magic id")
		return GetErrorByType(TGErrorBadMagic, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read MagicId as '%d'", magic))
	}

	protocolVersion, err := is.(*ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading protocolVersion from message buffer"))
		return err
	}
	if protocolVersion != int16(GetProtocolVersion()) {
		errMsg := fmt.Sprint("Unsupported protocol version")
		return GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read ProtocolVersion as '%d'", protocolVersion))
	}

	verbId, err := is.(*ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading verbId from message buffer"))
		return err
	}
	if verbId != int16(msg.GetVerbId()) {
		errMsg := fmt.Sprint("Incorrect Message Type")
		return GetErrorByType(TGErrorBadVerb, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read VerbId as '%d'", verbId))
	}

	sequenceNo, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading sequenceNo from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read SequenceNo as '%d'", sequenceNo))
	}

	timestamp, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading timestamp from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read Timestamp as '%d'", timestamp))
	}

	requestId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading requestId from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read RequestId as '%d'", requestId))
	}

	authToken, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading authToken from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read AuthToken as '%d'", authToken))
	}

	sessionId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading sessionId from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read SessionId as '%d'", sessionId))
	}

	tenantId, err := is.(*ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading tenantId from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read tenantId as '%d'", tenantId))
	}

	dataOffset, err := is.(*ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:APMReadHeader w/ Error in reading protocolVersion from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:APMReadHeader read DataOffset as '%d'", dataOffset))
	}

	// (Re)Set the message Attributes with correct values from input stream
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
	msg.SetTenantId(int(tenantId))
	msg.SetDataOffset(dataOffset)
	//msg.SetMessageByteBufLength(binary.Size(reflect.ValueOf(msg)))
	msg.SetMessageByteBufLength(int(binary.Size(reflect.ValueOf(msg))))
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AbstractProtocolMessage:APMReadHeader"))
	}
	return nil
}

func APMWriteHeader(msg tgdb.TGMessage, os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AbstractProtocolMessage:APMWriteHeader at output buffer position: '%d'", startPos))
	}
	os.(*ProtocolDataOutputStream).WriteInt(0) //The length is written later.
	os.(*ProtocolDataOutputStream).WriteInt(GetMagic())
	os.(*ProtocolDataOutputStream).WriteShort(int(GetProtocolVersion()))
	os.(*ProtocolDataOutputStream).WriteShort(msg.GetVerbId())

	os.(*ProtocolDataOutputStream).WriteLong(msg.GetSequenceNo())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetTimestamp())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetRequestId())

	os.(*ProtocolDataOutputStream).WriteLong(msg.GetAuthToken())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetSessionId())
	os.(*ProtocolDataOutputStream).WriteShort(msg.GetTenantId())
	os.(*ProtocolDataOutputStream).WriteShort(os.GetPosition() + 2) //DataOffset.
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractProtocolMessage::APMWriteHeader at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

// VerbIdFromBytes extracts message type from the input buffer in the byte format
func VerbIdFromBytes(buffer []byte) (*CommandVerbs, tgdb.TGError) {
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AbstractProtocolMessage:VerbIdFromBytes - received input buffer as '%+v'", buffer))
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:VerbIdFromBytes - extracted bufLen: '%d'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Second member attribute / element of message header is MagicId
	magic, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error in reading magicId from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:VerbIdFromBytes - extracted magic id: '%d'", magic))
	}
	if magic != GetMagic() {
		errMsg := fmt.Sprint("Bad Magic id")
		return nil, GetErrorByType(TGErrorBadMagic, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Third member attribute / element of message header is ProtocolVersion
	protocolVersion, err := is.ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:VerbIdFromBytes w/ Error in reading protocolVersion from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:VerbIdFromBytes - extracted protocolVersion: '%d'", protocolVersion))
	}
	if protocolVersion != int16(GetProtocolVersion()) {
		errMsg := fmt.Sprint("Unsupported protocol version")
		return nil, GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
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
func (msg *AbstractProtocolMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AbstractProtocolMessage:FromBytes - received input buffer as '%+v'", buffer))
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractProtocolMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractProtocolMessage:FromBytes read bufLen as '%d'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractProtocolMessage:FromBytes about to read Header data elements"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractProtocolMessage:FromBytes about to read Payload data elements"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AbstractProtocolMessage:FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AbstractProtocolMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AbstractProtocolMessage:ToBytes"))
	}

	msg.contentLock.Lock()
	defer msg.contentLock.Unlock()

	var bufLength int
	if msg.bytesBuffer == nil {
		os := DefaultProtocolDataOutputStream()

		err := APMWriteHeader(msg, os)
		if err != nil {
			errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
			return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
		}

		err = msg.WritePayload(os)
		if err != nil {
			errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
			return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
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
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AbstractProtocolMessage:ToBytes resulted bytes-on-the-wire in '%+v'", msg.bytesBuffer))
	}
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
func (msg *AbstractProtocolMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
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
func (msg *AbstractProtocolMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AbstractProtocolMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *AbstractProtocolMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *AbstractProtocolMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *AbstractProtocolMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
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
func (msg *AuthenticatedMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AuthenticatedMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AuthenticatedMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AuthenticatedMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AuthenticatedMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AuthenticatedMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AuthenticatedMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AuthenticatedMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticatedMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticatedMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AuthenticatedMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *AuthenticatedMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AuthenticatedMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AuthenticatedMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AuthenticatedMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AuthenticatedMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *AuthenticatedMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AuthenticatedMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AuthenticatedMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *AuthenticatedMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *AuthenticatedMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *AuthenticatedMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AuthenticatedMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *AuthenticatedMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *AuthenticatedMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *AuthenticatedMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AuthenticatedMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *AuthenticatedMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *AuthenticatedMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AuthenticatedMessage:{")
	buffer.WriteString(fmt.Sprintf("ConnectionId: %d", msg.connectionId))
	buffer.WriteString(fmt.Sprintf(", ClientId: %s", msg.clientId))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AuthenticatedMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AuthenticatedMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *AuthenticatedMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *AuthenticatedMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AuthenticatedMessage:ReadPayload"))
	}

	authToken, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:ReadPayload w/ Error in reading authToken from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AuthenticatedMessage:ReadPayload read authToken as '%+v'", authToken))
	}

	sessionId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticatedMessage:ReadPayload w/ Error in reading sessionId from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read sessionId as '%+v'", sessionId))
	}

	msg.SetAuthToken(authToken)
	msg.SetSessionId(sessionId)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AuthenticatedMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *AuthenticatedMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AuthenticatedMessage:WritePayload at output buffer position: '%d'", startPos))
	}
	if msg.GetAuthToken() == -1 || msg.GetSessionId() == -1 {
		errMsg := fmt.Sprint("Message not authenticated")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetAuthToken())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetSessionId())
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AuthenticatedMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
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



type DisconnectChannelRequestMessage struct {
	*AuthenticatedMessage
}

func DefaultDisconnectChannelRequestMessage() *DisconnectChannelRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(DisconnectChannelRequestMessage{})

	newMsg := DisconnectChannelRequestMessage{
		AuthenticatedMessage: DefaultAuthenticatedMessage(),
	}
	newMsg.verbId = VerbDisconnectChannelRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewDisconnectChannelRequestMessage(authToken, sessionId int64) *DisconnectChannelRequestMessage {
	newMsg := DefaultDisconnectChannelRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *DisconnectChannelRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering DisconnectChannelRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning DisconnectChannelRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DisconnectChannelRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside DisconnectChannelRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DisconnectChannelRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DisconnectChannelRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning DisconnectChannelRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *DisconnectChannelRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering DisconnectChannelRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DisconnectChannelRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DisconnectChannelRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DisconnectChannelRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning DisconnectChannelRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *DisconnectChannelRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *DisconnectChannelRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *DisconnectChannelRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *DisconnectChannelRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *DisconnectChannelRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *DisconnectChannelRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *DisconnectChannelRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *DisconnectChannelRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *DisconnectChannelRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *DisconnectChannelRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *DisconnectChannelRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *DisconnectChannelRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *DisconnectChannelRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *DisconnectChannelRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *DisconnectChannelRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *DisconnectChannelRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *DisconnectChannelRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *DisconnectChannelRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DisconnectChannelRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *DisconnectChannelRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *DisconnectChannelRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *DisconnectChannelRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *DisconnectChannelRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering DisconnectChannelRequestMessage:ReadPayload"))
	}
	authToken, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DisconnectChannelRequestMessage:ReadPayload w/ Error in reading authToken from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("DisconnectChannelRequestMessage:ReadPayload read authToken as '%+v'", authToken))
	}

	sessionId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DisconnectChannelRequestMessage:ReadPayload w/ Error in reading sessionId from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("DisconnectChannelRequestMessage:ReadPayload read sessionId as '%+v'", sessionId))
	}
	msg.SetAuthToken(authToken)
	msg.SetSessionId(sessionId)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning DisconnectChannelRequestMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *DisconnectChannelRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering DisconnectChannelRequestMessage:WritePayload at output buffer position: '%d'", startPos))
	}
	if msg.GetAuthToken() == -1 || msg.GetSessionId() == -1 {
		errMsg := fmt.Sprint("Message not authenticated")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetAuthToken())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetSessionId())
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning DisconnectChannelRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *DisconnectChannelRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DisconnectChannelRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *DisconnectChannelRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DisconnectChannelRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type AuthenticateRequestMessage struct {
	*AbstractProtocolMessage
	clientId  string
	inboxAddr string
	userName  string
	password  []byte
	dbName	string
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
func (msg *AuthenticateRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AuthenticateRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
	if err != nil {
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AuthenticateRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticateRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticateRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AuthenticateRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AuthenticateRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AuthenticateRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticateRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticateRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AuthenticateRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
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
func (msg *AuthenticateRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
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
func (msg *AuthenticateRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AuthenticateRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *AuthenticateRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *AuthenticateRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AuthenticateRequestMessage:ReadPayload"))
	}
	// For Testing purpose only.
	bIsClientId, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading clientId flag from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AuthenticateRequestMessage:ReadPayload read bIsClientId as '%+v'", bIsClientId))
	}
	strClient := ""
	if ! bIsClientId {
		strClient, err = is.(*ProtocolDataInputStream).ReadUTF()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading clientId from message buffer"))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AuthenticateRequestMessage:ReadPayload read strClient as '%+v'", strClient))
		}
	}

	inboxAddr, err := is.(*ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading inboxAddr flag from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AuthenticateRequestMessage:ReadPayload read inboxAddr as '%+v'", inboxAddr))
	}

	userName, err := is.(*ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading username from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AuthenticateRequestMessage:ReadPayload read userName as '%+v'", userName))
	}

	password, err := is.(*ProtocolDataInputStream).ReadBytes()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:ReadPayload w/ Error in reading password from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AuthenticateRequestMessage:ReadPayload read password as '%+v'", password))
	}

	msg.SetClientId(strClient)
	msg.SetInboxAddr(inboxAddr)
	msg.SetUserName(userName)
	msg.SetPassword(password)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AuthenticateRequestMessage:ReadPayload"))
	}
	return nil
}


//  SetDatabaseName sets the database-name on request message
func (msg *AuthenticateRequestMessage) SetDatabaseName (dbname string) {
	msg.dbName = dbname
}

// GetDatabaseName gets the database-name from request message
func (msg *AuthenticateRequestMessage) GetDatabaseName () string {
	return msg.dbName
}


// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *AuthenticateRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AuthenticateRequestMessage:WritePayload at output buffer position: '%d'", startPos))
	}

	if msg.GetDatabaseName() == "" {
		os.(*ProtocolDataOutputStream).WriteBoolean(true)
	} else {
		os.(*ProtocolDataOutputStream).WriteBoolean(false)
		err := os.(*ProtocolDataOutputStream).WriteUTF(msg.dbName)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:WritePayload w/ Error in writing dbName to message buffer"))
			return err
		}
	}

	if msg.GetClientId() == "" {
		os.(*ProtocolDataOutputStream).WriteBoolean(true)
	} else {
		os.(*ProtocolDataOutputStream).WriteBoolean(false) // No client id
		err := os.(*ProtocolDataOutputStream).WriteUTF(msg.GetClientId())
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:WritePayload w/ Error in writing clientId to message buffer"))
			return err
		}
	}
	if msg.GetInboxAddr() == "" {
		os.(*ProtocolDataOutputStream).WriteBoolean(true)
	} else {
		os.(*ProtocolDataOutputStream).WriteBoolean(false)
		err := os.(*ProtocolDataOutputStream).WriteUTF(msg.GetInboxAddr())
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:WritePayload w/ Error in writing inboxAddr to message buffer"))
			return err
		}
	}

	//TODO in Java client as well
	//FIXME: Need to add full support for specifying roles
	//Currently set to use all roles
	os.(*ProtocolDataOutputStream).WriteInt(-1)


	if msg.GetUserName() == "" {
		os.(*ProtocolDataOutputStream).WriteBoolean(true)
	} else {
		os.(*ProtocolDataOutputStream).WriteBoolean(false)
		err := os.(*ProtocolDataOutputStream).WriteUTF(msg.GetUserName())	// Can't be null.
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AuthenticateRequestMessage:WritePayload w/ Error in writing userName to message buffer"))
			return err
		}
	}
	err := os.(*ProtocolDataOutputStream).WriteBytes(msg.GetPassword())
	if err != nil {
		return err
	}
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AuthenticateRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return os.(*ProtocolDataOutputStream).WriteBytes(msg.GetPassword())
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



type AuthenticateResponseMessage struct {
	*AbstractProtocolMessage
	successFlag      bool
	errorStatus      int
	serverCertBuffer []byte
}

func DefaultAuthenticateResponseMessage() *AuthenticateResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AuthenticateResponseMessage{})

	newMsg := AuthenticateResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
		successFlag:             false,
	}
	newMsg.authToken = -1
	newMsg.sessionId = -1
	newMsg.verbId = VerbAuthenticateResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewAuthenticateResponseMessage(authToken, sessionId int64) *AuthenticateResponseMessage {
	newMsg := DefaultAuthenticateResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for AuthenticateResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *AuthenticateResponseMessage) GetServerCertBuffer() []byte {
	return msg.serverCertBuffer
}

func (msg *AuthenticateResponseMessage) IsSuccess() bool {
	return msg.successFlag
}

func (msg *AuthenticateResponseMessage) SetErrorStatus(status int) {
	msg.errorStatus = status
}

func (msg *AuthenticateResponseMessage) SetServerCertBuffer(buffer []byte) {
	msg.serverCertBuffer = buffer
}

func (msg *AuthenticateResponseMessage) SetSuccess(flag bool) {
	msg.successFlag = flag
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *AuthenticateResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AuthenticateResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AuthenticateResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AuthenticateResponseMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticateResponseMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AuthenticateResponseMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AuthenticateResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AuthenticateResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticateResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AuthenticateResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AuthenticateResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AuthenticateResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *AuthenticateResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AuthenticateResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AuthenticateResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AuthenticateResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AuthenticateResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *AuthenticateResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AuthenticateResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AuthenticateResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *AuthenticateResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *AuthenticateResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *AuthenticateResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AuthenticateResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *AuthenticateResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *AuthenticateResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *AuthenticateResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AuthenticateResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *AuthenticateResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *AuthenticateResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AuthenticateResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("SuccessFlag: %+v", msg.successFlag))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	buffer.WriteString(fmt.Sprintf(", ErrorStatus: %+v", msg.errorStatus))
	//buffer.WriteString(fmt.Sprintf(", ServerCertBuffer: %+v", msg.serverCertBuffer))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AuthenticateResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AuthenticateResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *AuthenticateResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *AuthenticateResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AuthenticateResponseMessage:ReadPayload"))
	}
	bSuccess, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading success flag from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("AuthenticateResponseMessage:ReadPayload read bSuccess as '%+v'", bSuccess))
	}

	if !bSuccess {
		errorStatus, err := is.(*ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading errorStatus from message buffer"))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("AuthenticateResponseMessage:ReadPayload read errorStatus as '%+v'", errorStatus))
		}
		msg.SetErrorStatus(errorStatus)
	}

	authToken, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading authToken from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("AuthenticateResponseMessage:ReadPayload read authToken as '%+v'", authToken))
	}

	sessionId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading sessionId from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("AuthenticateResponseMessage:ReadPayload read sessionId as '%+v'", sessionId))
	}

	certBuffer, err := is.(*ProtocolDataInputStream).ReadBytes()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading certBuffer from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("AuthenticateResponseMessage:ReadPayload read certBuffer as '%+v'", certBuffer))
	}

	msg.SetSuccess(bSuccess)
	msg.SetAuthToken(authToken)
	msg.SetSessionId(sessionId)
	msg.SetServerCertBuffer(certBuffer)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AuthenticateResponseMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *AuthenticateResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AuthenticateResponseMessage:WritePayload at output buffer position: '%d'", startPos))
	}
	os.(*ProtocolDataOutputStream).WriteBoolean(msg.IsSuccess())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetAuthToken())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetSessionId())
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AuthenticateResponseMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *AuthenticateResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable, msg.successFlag, msg.errorStatus, msg.serverCertBuffer)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AuthenticateResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *AuthenticateResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable,
		&msg.successFlag, &msg.errorStatus, &msg.serverCertBuffer)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AuthenticateResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type BeginTransactionRequestMessage struct {
	*AbstractProtocolMessage
}

func DefaultBeginTransactionRequestMessage() *BeginTransactionRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(BeginTransactionRequestMessage{})

	newMsg := BeginTransactionRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbBeginTransactionRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewBeginTransactionRequestMessage(authToken, sessionId int64) *BeginTransactionRequestMessage {
	newMsg := DefaultBeginTransactionRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *BeginTransactionRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering BeginTransactionRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning BeginTransactionRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning BeginTransactionRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside BeginTransactionRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside BeginTransactionRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside BeginTransactionRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning BeginTransactionRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *BeginTransactionRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering BeginTransactionRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside BeginTransactionRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside BeginTransactionRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning BeginTransactionRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning BeginTransactionRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *BeginTransactionRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *BeginTransactionRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *BeginTransactionRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *BeginTransactionRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *BeginTransactionRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *BeginTransactionRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *BeginTransactionRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *BeginTransactionRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *BeginTransactionRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *BeginTransactionRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *BeginTransactionRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *BeginTransactionRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *BeginTransactionRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *BeginTransactionRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *BeginTransactionRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *BeginTransactionRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *BeginTransactionRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *BeginTransactionRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("BeginTransactionRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *BeginTransactionRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *BeginTransactionRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *BeginTransactionRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *BeginTransactionRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *BeginTransactionRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *BeginTransactionRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BeginTransactionRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *BeginTransactionRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BeginTransactionRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type BeginTransactionResponseMessage struct {
	*AbstractProtocolMessage
	transactionId int64
}

func DefaultBeginTransactionResponseMessage() *BeginTransactionResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(BeginTransactionResponseMessage{})

	newMsg := BeginTransactionResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.transactionId = -1
	newMsg.verbId = VerbBeginTransactionResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewBeginTransactionResponseMessage(authToken, sessionId int64) *BeginTransactionResponseMessage {
	newMsg := DefaultBeginTransactionResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for BeginTransactionResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *BeginTransactionResponseMessage) GetTransactionId() int64 {
	return msg.transactionId
}

func (msg *BeginTransactionResponseMessage) SetTransactionId(txnId int64) {
	msg.transactionId = txnId
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *BeginTransactionResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering BeginTransactionResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning BeginTransactionResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning BeginTransactionResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside BeginTransactionResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside BeginTransactionResponseMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside BeginTransactionResponseMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning BeginTransactionResponseMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *BeginTransactionResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering BeginTransactionResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside BeginTransactionResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside BeginTransactionResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning BeginTransactionResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning BeginTransactionResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *BeginTransactionResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *BeginTransactionResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *BeginTransactionResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *BeginTransactionResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *BeginTransactionResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *BeginTransactionResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *BeginTransactionResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *BeginTransactionResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *BeginTransactionResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *BeginTransactionResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *BeginTransactionResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *BeginTransactionResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *BeginTransactionResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *BeginTransactionResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *BeginTransactionResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *BeginTransactionResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *BeginTransactionResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *BeginTransactionResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("BeginTransactionResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("ClientId: %d", msg.transactionId))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *BeginTransactionResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *BeginTransactionResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *BeginTransactionResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *BeginTransactionResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering BeginTransactionResponseMessage:ReadPayload"))
	}
	txnId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning BeginTransactionResponseMessage:ReadPayload w/ Error in reading txnId from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AuthenticateRequestMessage:ReadPayload read txnId as '%+v'", txnId))
	}
	msg.SetTransactionId(txnId)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning BeginTransactionResponseMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *BeginTransactionResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering BeginTransactionResponseMessage:WritePayload at output buffer position: '%d'", startPos))
	}
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetTransactionId())
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning BeginTransactionResponseMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *BeginTransactionResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable, msg.transactionId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BeginTransactionResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *BeginTransactionResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable, &msg.transactionId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BeginTransactionResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type CommitTransactionRequest struct {
	*AbstractProtocolMessage
	addedList   map[int64]tgdb.TGEntity
	updatedList map[int64]tgdb.TGEntity
	removedList map[int64]tgdb.TGEntity
	attrDescSet []tgdb.TGAttributeDescriptor
}

func DefaultCommitTransactionRequestMessage() *CommitTransactionRequest {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(CommitTransactionRequest{})

	newMsg := CommitTransactionRequest{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbCommitTransactionRequest
	newMsg.addedList = make(map[int64]tgdb.TGEntity, 0)
	newMsg.updatedList = make(map[int64]tgdb.TGEntity, 0)
	newMsg.removedList = make(map[int64]tgdb.TGEntity, 0)
	newMsg.attrDescSet = make([]tgdb.TGAttributeDescriptor, 0)
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewCommitTransactionRequestMessage(authToken, sessionId int64) *CommitTransactionRequest {
	newMsg := DefaultCommitTransactionRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for VerbCommitTransactionRequest message
/////////////////////////////////////////////////////////////////

func (msg *CommitTransactionRequest) AddCommitLists(addedList, updatedList, removedList map[int64]tgdb.TGEntity, attrDescriptors []tgdb.TGAttributeDescriptor) *CommitTransactionRequest {
	if len(addedList) > 0 {
		msg.addedList = addedList
	}
	if len(updatedList) > 0 {
		msg.updatedList = updatedList
	}
	if len(removedList) > 0 {
		msg.removedList = removedList
	}
	if len(attrDescriptors) > 0 {
		msg.attrDescSet = attrDescriptors
	}
	return msg
}

func (msg *CommitTransactionRequest) GetAddedList() map[int64]tgdb.TGEntity {
	return msg.addedList
}

func (msg *CommitTransactionRequest) GetUpdatedList() map[int64]tgdb.TGEntity {
	return msg.updatedList
}

func (msg *CommitTransactionRequest) GetRemovedList() map[int64]tgdb.TGEntity {
	return msg.removedList
}

func (msg *CommitTransactionRequest) GetAttrDescSet() []tgdb.TGAttributeDescriptor {
	return msg.attrDescSet
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *CommitTransactionRequest) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering CommitTransactionRequest:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside CommitTransactionRequest:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside CommitTransactionRequest:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside CommitTransactionRequest:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning CommitTransactionRequest::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *CommitTransactionRequest) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering CommitTransactionRequest::ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside CommitTransactionRequest:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside CommitTransactionRequest:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning CommitTransactionRequest::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()[:os.GetLength()]))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *CommitTransactionRequest) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *CommitTransactionRequest) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *CommitTransactionRequest) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *CommitTransactionRequest) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *CommitTransactionRequest) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *CommitTransactionRequest) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *CommitTransactionRequest) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *CommitTransactionRequest) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *CommitTransactionRequest) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *CommitTransactionRequest) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *CommitTransactionRequest) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *CommitTransactionRequest) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *CommitTransactionRequest) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *CommitTransactionRequest) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *CommitTransactionRequest) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *CommitTransactionRequest) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *CommitTransactionRequest) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *CommitTransactionRequest) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("CommitTransactionRequest:{")
	buffer.WriteString(fmt.Sprint("AttrDescSet: "))
	buffer.WriteString("{")
	//for _, d := range msg.AttrDescSet {
	//	buffer.WriteString(fmt.Sprintf("Attribute Descriptor: %+v ", d))
	//}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", AddedList:{"))
	for k, v := range msg.addedList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d=Entity: %+v ", k, v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", UpdatedList:{"))
	for k, v := range msg.updatedList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d=Entity: %+v ", k, v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", RemovedList:{"))
	for k, v := range msg.removedList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d=Entity: %+v ", k, v))
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
func (msg *CommitTransactionRequest) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *CommitTransactionRequest) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *CommitTransactionRequest) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *CommitTransactionRequest) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	//Commit response need to send back real id for all entities and descriptors.
	// No-op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *CommitTransactionRequest) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering CommitTransactionRequest:ReadPayload at output buffer position: '%d'", startPos))
	}
	os.(*ProtocolDataOutputStream).WriteInt(0) // This is for the commit buffer length
	os.(*ProtocolDataOutputStream).WriteInt(0) // This is for the checksum for the commit buffer to be added later.  Currently not used
	////<A> for attribute descriptor, <N> for node desc definitions, <E> for edge desc definitions
	////meta should be sent before the instance data
	if len(msg.attrDescSet) > 0 {
		os.(*ProtocolDataOutputStream).WriteShort(0x1010) // For attribute descriptor
		// There should be nothing after the marker due to no new attribute descriptor
		// Need to check for new descriptor only with attribute id as negative number
		// Check for size overrun
		newAttrCount := 0
		for _, attrDesc := range msg.attrDescSet {
			if attrDesc.GetAttributeId() < 0 {
				newAttrCount++
			}
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CommitTransactionRequest:ReadPayload - There are '%d' new attribute descriptors", newAttrCount))
		}
		os.(*ProtocolDataOutputStream).WriteInt(newAttrCount)
		for _, attrDesc := range msg.attrDescSet {
			err := attrDesc.WriteExternal(os)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning CommitTransactionRequest:WritePayload w/ Error in writing AttrDesc '%+v' to message buffer", attrDesc))
				return err
			}
			//}
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CommitTransactionRequest:ReadPayload - '%d' attribute descriptors are written in byte format", len(msg.attrDescSet)))
		}
	}
	if len(msg.addedList) > 0 {
		os.(*ProtocolDataOutputStream).WriteShort(0x1011) // For entity creation
		os.(*ProtocolDataOutputStream).WriteInt(len(msg.addedList))
		for _, entity := range msg.addedList {
			err := entity.WriteExternal(os)
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:WritePayload w/ Error in writing addedEntity to message buffer"))
				return err
			}
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CommitTransactionRequest:WritePayload - '%d' new entities are written in byte format", len(msg.addedList)))
		}
	}
	//TODO: Ask TGDB Engineering Team - Need to write only the modified Attributes
	if len(msg.updatedList) > 0 {
		os.(*ProtocolDataOutputStream).WriteShort(0x1012) // For entity update
		os.(*ProtocolDataOutputStream).WriteInt(len(msg.updatedList))
		for _, entity := range msg.updatedList {
			err := entity.WriteExternal(os)
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:WritePayload w/ Error in writing updatedEntity to message buffer"))
				return err
			}
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CommitTransactionRequest:WritePayload - '%d' updateable entities are written in byte format", len(msg.updatedList)))
		}
	}
	//TODO: Ask TGDB Engineering Team - Need to write the id only
	if len(msg.removedList) > 0 {
		os.(*ProtocolDataOutputStream).WriteShort(0x1013) // For deleted entities
		os.(*ProtocolDataOutputStream).WriteInt(len(msg.removedList))
		for _, entity := range msg.removedList {
			err := entity.WriteExternal(os)
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:WritePayload w/ Error in writing removedEntity to message buffer"))
				return err
			}
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CommitTransactionRequest:WritePayload - '%d' removable entities are written in byte format", len(msg.removedList)))
		}
	}
	currPos := os.GetPosition()
	length := currPos - startPos
	_, err := os.(*ProtocolDataOutputStream).WriteIntAt(startPos, length)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning CommitTransactionRequest::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *CommitTransactionRequest) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable, msg.addedList,
		msg.updatedList, msg.removedList, msg.attrDescSet)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CommitTransactionRequest:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *CommitTransactionRequest) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable,
		&msg.addedList, &msg.updatedList, &msg.removedList, &msg.attrDescSet)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CommitTransactionRequest:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



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
	graphObjFact   tgdb.TGGraphObjectFactory
	exception      *TransactionException
	entityStream   tgdb.TGInputStream
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

func ProcessTransactionStatus(is tgdb.TGInputStream, status int) *TransactionException {
	txnStatus := tgdb.FromStatus(status)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering CommitTransactionResponse:ProcessTransactionStatus read txnStatus as '%d'", txnStatus))
	}
	if txnStatus == tgdb.TGTransactionSuccess {
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Returning CommitTransactionResponse:ProcessTransactionStatus NO EXCEPTION for txnStatus:'%+v'", txnStatus))
		}
		return nil
	}
	switch txnStatus {
	case tgdb.TGTransactionSuccess:
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Returning CommitTransactionResponse:ProcessTransactionStatus NO EXCEPTION for txnStatus:'%+v'", txnStatus))
		}
		return nil
	case tgdb.TGTransactionAlreadyInProgress:
		fallthrough
	case tgdb.TGTransactionClientDisconnected:
		fallthrough
	case tgdb.TGTransactionMalFormed:
		fallthrough
	case tgdb.TGTransactionGeneralError:
		fallthrough
	case tgdb.TGTransactionVerificationError:
		fallthrough
	case tgdb.TGTransactionInBadState:
		fallthrough
	case tgdb.TGTransactionUniqueConstraintViolation:
		fallthrough
	case tgdb.TGTransactionOptimisticLockFailed:
		fallthrough
	case tgdb.TGTransactionResourceExceeded:
		fallthrough
	case tgdb.TGCurrentThreadNotInTransaction:
		fallthrough
	case tgdb.TGTransactionUniqueIndexKeyAttributeNullError:
		fallthrough
	default:
		errMsg, err := is.(*ProtocolDataInputStream).ReadUTF()
		if err != nil || errMsg == "" {
			errMsg = "Error not available"
		}
		logger.Error(fmt.Sprintf("Returning CommitTransactionResponse:ProcessTransactionStatus for txnStatus:'%+v' w/ error: '%+v'", txnStatus, errMsg))
		return BuildException(txnStatus, errMsg)
	}
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

func (msg *CommitTransactionResponse) GetException() *TransactionException {
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
func (msg *CommitTransactionResponse) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering CommitTransactionResponse:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside CommitTransactionResponse:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		//errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		//return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		return nil, err
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside CommitTransactionResponse:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		//errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		//return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		return nil, err
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning CommitTransactionResponse::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *CommitTransactionResponse) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering CommitTransactionResponse:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside CommitTransactionResponse:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside CommitTransactionResponse:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning CommitTransactionResponse::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
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
func (msg *CommitTransactionResponse) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
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
func (msg *CommitTransactionResponse) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *CommitTransactionResponse) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *CommitTransactionResponse) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *CommitTransactionResponse) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering CommitTransactionResponse:ReadPayload"))
	}
	bLen, err := is.(*ProtocolDataInputStream).ReadInt() // buf length
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading bLen from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read buf length as '%+v'", bLen))
	}

	checkSum, err := is.(*ProtocolDataInputStream).ReadInt() // checksum
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading checkSum from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read checkSum as '%+v'", checkSum))
	}

	status, err := is.(*ProtocolDataInputStream).ReadInt() // status code - currently zero
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading status from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read status as '%+v'", status))
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload - about to ProcessTransactionStatus for status '%+v'", status))
	}
	txnException := ProcessTransactionStatus(is, status)
	if txnException != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ txnException"))
		return txnException
	}

	for {
		avail, err := is.(*ProtocolDataInputStream).Available()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading availability bytes from message buffer"))
			return err
		}
		if avail <= 0 {
			break
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read avail as '%+v'", avail))
		}

		opCode, err := is.(*ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading opCode from message buffer"))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read opCode as '%+v'", opCode))
		}

		switch opCode {
		case 0x1010:
			attrDescCount, err := is.(*ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading attrDescCount from message buffer"))
				return err
			}
			msg.attrDescCount = attrDescCount
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read attrDescCount for opCode=0x1010 as '%+v'", attrDescCount))
			}
			for i := 0; i < attrDescCount; i++ {
				tempId, _ := is.(*ProtocolDataInputStream).ReadInt()
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read tempId for opCode=0x1010 as '%+v'", tempId))
				}
				msg.attrDescIdList = append(msg.attrDescIdList, int64(tempId))
				realId, _ := is.(*ProtocolDataInputStream).ReadInt()
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read realId for opCode=0x1010 as '%+v'", realId))
				}
				msg.attrDescIdList = append(msg.attrDescIdList, int64(realId))
			}
			break
		case 0x1011:
			addedCount, err := is.(*ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading addedCount from message buffer"))
				return err
			}
			msg.addedCount = addedCount
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read addedCount for opCode=0x1011 as '%+v'", addedCount))
			}
			for i := 0; i < addedCount; i++ {
				longTempId, _ := is.(*ProtocolDataInputStream).ReadLong()
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read longTempId for opCode=0x1011 as '%+v'", longTempId))
				}
				msg.addedIdList = append(msg.addedIdList, longTempId)
				longRealId, _ := is.(*ProtocolDataInputStream).ReadLong()
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read longRealId for opCode=0x1011 as '%+v'", longRealId))
				}
				msg.addedIdList = append(msg.addedIdList, longRealId)
				longVersion, _ := is.(*ProtocolDataInputStream).ReadLong()
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read longVersion for opCode=0x1011 as '%+v'", longVersion))
				}
				msg.addedIdList = append(msg.addedIdList, longVersion)
			}
			break
		case 0x1012:
			updatedCount, err := is.(*ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading updatedCount from message buffer"))
				return err
			}
			msg.updatedCount = updatedCount
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read updatedCount for opCode=0x1012 as '%+v'", updatedCount))
			}
			for i := 0; i < updatedCount; i++ {
				id, _ := is.(*ProtocolDataInputStream).ReadLong()
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read id for opCode=0x1012 as '%+v'", id))
				}
				msg.updatedIdList = append(msg.updatedIdList, id)
				version, _ := is.(*ProtocolDataInputStream).ReadLong()
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read version for opCode=0x1012 as '%+v'", version))
				}
				msg.updatedIdList = append(msg.updatedIdList, version)
			}
			break
		case 0x1013:
			removedCount, err := is.(*ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading removedCount from message buffer"))
				return err
			}
			msg.removedCount = removedCount
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read removedCount for opCode=0x1013 as '%+v'", removedCount))
			}
			for i := 0; i < removedCount; i++ {
				id, _ := is.(*ProtocolDataInputStream).ReadLong()
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read id for opCode=0x1013 as '%+v'", id))
				}
				msg.removedIdList = append(msg.removedIdList, id)
			}
			break
		case 0x6789:
			msg.entityStream = is
			pos := is.(*ProtocolDataInputStream).GetPosition()
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read pos for opCode=0x6789 as '%+v'", pos))
			}
			count, err := is.(*ProtocolDataInputStream).ReadInt()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionResponse:ReadPayload w/ Error in reading count from message buffer"))
				return err
			}
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside CommitTransactionResponse:ReadPayload read count for opCode=0x6789 as '%+v'", count))
			}
			is.(*ProtocolDataInputStream).SetPosition(pos)
			break
		default:
			break
		}
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning CommitTransactionResponse:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *CommitTransactionResponse) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
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



type DecryptBufferRequestMessage struct {
	*AbstractProtocolMessage
	encryptedBuffer []byte
}

func DefaultDecryptBufferRequestMessage() *DecryptBufferRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(DecryptBufferRequestMessage{})

	newMsg := DecryptBufferRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbDecryptBufferResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewDecryptBufferRequestMessage(authToken, sessionId int64) *DecryptBufferRequestMessage {
	newMsg := DefaultDecryptBufferRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for DecryptBufferRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *DecryptBufferRequestMessage) GetEncryptedBuffer() []byte {
	return msg.encryptedBuffer
}

func (msg *DecryptBufferRequestMessage) SetEncryptedBuffer(buf []byte) {
	msg.encryptedBuffer = buf
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *DecryptBufferRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering DecryptBufferRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning DecryptBufferRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DecryptBufferRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside DecryptBufferRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside DecryptBufferRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside DecryptBufferRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning DecryptBufferRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *DecryptBufferRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering DecryptBufferRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DecryptBufferRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DecryptBufferRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DecryptBufferRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning DecryptBufferRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *DecryptBufferRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *DecryptBufferRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *DecryptBufferRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *DecryptBufferRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *DecryptBufferRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *DecryptBufferRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *DecryptBufferRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *DecryptBufferRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *DecryptBufferRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *DecryptBufferRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *DecryptBufferRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *DecryptBufferRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *DecryptBufferRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *DecryptBufferRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *DecryptBufferRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *DecryptBufferRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *DecryptBufferRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *DecryptBufferRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DecryptBufferRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *DecryptBufferRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *DecryptBufferRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *DecryptBufferRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *DecryptBufferRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *DecryptBufferRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	if msg.GetEncryptedBuffer() == nil {
		errMsg := fmt.Sprint("Encrypted Buffer is EMPTY")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return os.(*ProtocolDataOutputStream).WriteBytes(msg.GetEncryptedBuffer())
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *DecryptBufferRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable, msg.encryptedBuffer)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DecryptBufferRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *DecryptBufferRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable, &msg.encryptedBuffer)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DecryptBufferRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type DecryptBufferResponseMessage struct {
	*AbstractProtocolMessage
	decryptedBuffer []byte
}

func DefaultDecryptBufferResponseMessage() *DecryptBufferResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(DecryptBufferResponseMessage{})

	newMsg := DecryptBufferResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbDecryptBufferResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewDecryptBufferResponseMessage(authToken, sessionId int64) *DecryptBufferResponseMessage {
	newMsg := DefaultDecryptBufferResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for DecryptBufferResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *DecryptBufferResponseMessage) GetDecryptedBuffer() []byte {
	return msg.decryptedBuffer
}

func (msg *DecryptBufferResponseMessage) SetDecryptedBuffer(buf []byte) {
	msg.decryptedBuffer = buf
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *DecryptBufferResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering DecryptBufferResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning DecryptBufferResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DecryptBufferResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside DecryptBufferResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside DecryptBufferResponseMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside DecryptBufferResponseMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning DecryptBufferResponseMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *DecryptBufferResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering DecryptBufferResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DecryptBufferResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DecryptBufferResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DecryptBufferResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning DecryptBufferResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *DecryptBufferResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *DecryptBufferResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *DecryptBufferResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *DecryptBufferResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *DecryptBufferResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *DecryptBufferResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *DecryptBufferResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *DecryptBufferResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *DecryptBufferResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *DecryptBufferResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *DecryptBufferResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *DecryptBufferResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *DecryptBufferResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *DecryptBufferResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *DecryptBufferResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *DecryptBufferResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *DecryptBufferResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *DecryptBufferResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DecryptBufferResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *DecryptBufferResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *DecryptBufferResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *DecryptBufferResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *DecryptBufferResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering DecryptBufferResponseMessage:ReadPayload"))
	}
	decryptBuffer, err := is.(*ProtocolDataInputStream).ReadBytes()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DecryptBufferResponseMessage:ReadPayload w/ Error in reading decryptBuffer from message buffer"))
		return err
	}
	//logger.Debug(fmt.Sprintf("DecryptBufferResponseMessage:ReadPayload read decryptBuffer as '%+v'", decryptBuffer))
	msg.SetDecryptedBuffer(decryptBuffer)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning DecryptBufferResponseMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *DecryptBufferResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *DecryptBufferResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.authToken, msg.sessionId, msg.dataOffset, msg.isUpdatable, msg.decryptedBuffer)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DecryptBufferResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *DecryptBufferResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.authToken, &msg.sessionId, &msg.dataOffset, &msg.isUpdatable, &msg.decryptedBuffer)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DecryptBufferResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type DumpStacktraceRequestMessage struct {
	*AbstractProtocolMessage
}

func DefaultDumpStacktraceRequestMessage() *DumpStacktraceRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(DumpStacktraceRequestMessage{})

	newMsg := DumpStacktraceRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbDumpStacktraceRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewDumpStacktraceRequestMessage(authToken, sessionId int64) *DumpStacktraceRequestMessage {
	newMsg := DefaultDumpStacktraceRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *DumpStacktraceRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering DumpStacktraceRequest:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning DumpStacktraceRequest:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DumpStacktraceRequest:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside DumpStacktraceRequest:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DumpStacktraceRequest:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside DumpStacktraceRequest:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("DumpStacktraceRequest::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *DumpStacktraceRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering DumpStacktraceRequest:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DumpStacktraceRequest:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside DumpStacktraceRequest:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DumpStacktraceRequest:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("DumpStacktraceRequest::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *DumpStacktraceRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *DumpStacktraceRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *DumpStacktraceRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *DumpStacktraceRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *DumpStacktraceRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *DumpStacktraceRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *DumpStacktraceRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *DumpStacktraceRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *DumpStacktraceRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *DumpStacktraceRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *DumpStacktraceRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *DumpStacktraceRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *DumpStacktraceRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *DumpStacktraceRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *DumpStacktraceRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *DumpStacktraceRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *DumpStacktraceRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *DumpStacktraceRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DumpStacktraceRequest:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param Timestamp
// @return TGMessage on success, error on failure
func (msg *DumpStacktraceRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *DumpStacktraceRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *DumpStacktraceRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *DumpStacktraceRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *DumpStacktraceRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *DumpStacktraceRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DumpStacktraceRequest:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *DumpStacktraceRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DumpStacktraceRequest:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type ExceptionMessage struct {
	*AbstractProtocolMessage
	servercode int
	exceptionMsg  string
	exceptionType int
}

func DefaultExceptionMessage() *ExceptionMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ExceptionMessage{})

	newMsg := ExceptionMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbExceptionMessage
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	newMsg.servercode = -1
	return &newMsg
}

// Create New Message Instance
func NewExceptionMessage(authToken, sessionId int64) *ExceptionMessage {
	newMsg := DefaultExceptionMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return newMsg
}

func NewExceptionMessageWithType(exType int, msg string) *ExceptionMessage {
	newMsg := DefaultExceptionMessage()
	newMsg.exceptionMsg = msg
	newMsg.exceptionType = exType
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

func NewExceptionMessageWithTypeWithServerErrorCode(exType int, msg string, sCode int) *ExceptionMessage {
	newMsg := NewExceptionMessageWithType(exType, msg)
	newMsg.servercode = sCode
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper Public functions for VerbExceptionMessage
/////////////////////////////////////////////////////////////////

func BuildFromException(ex tgdb.TGError) *ExceptionMessage {
	exceptionMsg := NewExceptionMessageWithType(ex.GetErrorType(), ex.Error())
	return exceptionMsg
}

func (msg *ExceptionMessage) GetExceptionMsg() string {
	return msg.exceptionMsg
}

func (msg *ExceptionMessage) GetExceptionType() int {
	return msg.exceptionType
}

func (msg *ExceptionMessage) GetServerErrorCode() int {
	return msg.servercode
}

func (msg *ExceptionMessage) SetExceptionMsg(exMsg string) {
	msg.exceptionMsg = exMsg
}

func (msg *ExceptionMessage) SetExceptionType(exType int) {
	msg.exceptionType = exType
}

func (msg *ExceptionMessage) SetServerErrorCode(sCode int) {
	msg.servercode = sCode
}


/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *ExceptionMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering ExceptionMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ExceptionMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ExceptionMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside ExceptionMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside ExceptionMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside ExceptionMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("ExceptionMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *ExceptionMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering ExceptionMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside ExceptionMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside ExceptionMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ExceptionMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("ExceptionMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *ExceptionMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *ExceptionMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *ExceptionMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *ExceptionMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *ExceptionMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *ExceptionMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *ExceptionMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *ExceptionMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *ExceptionMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *ExceptionMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *ExceptionMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *ExceptionMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *ExceptionMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *ExceptionMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *ExceptionMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *ExceptionMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *ExceptionMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *ExceptionMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ExceptionMessage:{")
	buffer.WriteString(fmt.Sprintf("ExceptionMsg: %s", msg.exceptionMsg))
	buffer.WriteString(fmt.Sprintf(", ExceptionType: %d", msg.exceptionType))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *ExceptionMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *ExceptionMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *ExceptionMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

func (msg *ExceptionMessage) mapServerCodeToExceptionType () {
	// TODO
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *ExceptionMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering ExceptionMessage:ReadPayload"))
	}
	sCode, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ExceptionMessage:ReadPayload w/ Error in reading error code from message buffer"))
		return err
	}
	msg.servercode = sCode
	if msg.servercode == 0 {
		return nil
	}

	isNull, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ExceptionMessage:ReadPayload w/ Error in reading isNull from message buffer"))
		return err
	}

	if !isNull {
		msg.exceptionMsg, err = is.(*ProtocolDataInputStream).ReadUTF()
	}

	msg.mapServerCodeToExceptionType()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning ExceptionMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *ExceptionMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	os.(*ProtocolDataOutputStream).WriteByte(msg.GetExceptionType())
	return os.(*ProtocolDataOutputStream).WriteUTF(msg.GetExceptionMsg())
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *ExceptionMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.exceptionMsg, msg.exceptionType)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ExceptionMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *ExceptionMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.exceptionMsg, &msg.exceptionType)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ExceptionMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type GetEntityRequestMessage struct {
	*AbstractProtocolMessage
	commandType    int16 //0 - get, 1 - getbyid, 2 - get multiples, 10 - continue, 20 - close
	fetchSize      int
	batchSize      int
	traversalDepth int
	edgeLimit      int
	resultId       int
	key            tgdb.TGKey
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

func (msg *GetEntityRequestMessage) GetKey() tgdb.TGKey {
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

func (msg *GetEntityRequestMessage) SetKey(key tgdb.TGKey) {
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
func (msg *GetEntityRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering GetEntityRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetEntityRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside GetEntityRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside GetEntityRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning GetEntityRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *GetEntityRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering GetEntityRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetEntityRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetEntityRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning GetEntityRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
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
func (msg *GetEntityRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
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
func (msg *GetEntityRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *GetEntityRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *GetEntityRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *GetEntityRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *GetEntityRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	os.(*ProtocolDataOutputStream).WriteShort(int(msg.GetCommand()))
	os.(*ProtocolDataOutputStream).WriteInt(msg.GetResultId())
	if msg.GetCommand() == 0 || msg.GetCommand() == 1 || msg.GetCommand() == 2 {
		os.(*ProtocolDataOutputStream).WriteInt(msg.GetFetchSize())
		os.(*ProtocolDataOutputStream).WriteShort(msg.GetBatchSize())
		os.(*ProtocolDataOutputStream).WriteShort(msg.GetTraversalDepth())
		os.(*ProtocolDataOutputStream).WriteShort(msg.GetEdgeLimit())

		startPos := os.(*ProtocolDataOutputStream).GetPosition()

		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering GetEntityRequestMessage:WritePayload at output buffer position: '%d'", startPos))
		}

		os.(*ProtocolDataOutputStream).WriteInt(0)
		bufPos := os.(*ProtocolDataOutputStream).GetPosition()

		err := msg.GetKey().WriteExternal(os)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning GetEntityRequestMessage:WritePayload w/ Error in writing key to message buffer"))
			return err
		}

		currPos := os.GetPosition()
		length := currPos - bufPos
		os.(*ProtocolDataOutputStream).WriteIntAt(startPos, length)

		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning GetEntityRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
		}

	}
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


type GetEntityResponseMessage struct {
	*AbstractProtocolMessage
	entityStream tgdb.TGInputStream
	hasResult    bool
	resultId     int
	totalCount   int
	resultCount  int
}

func DefaultGetEntityResponseMessage() *GetEntityResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GetEntityResponseMessage{})

	newMsg := GetEntityResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.hasResult = false
	newMsg.resultId = 0
	newMsg.totalCount = 0
	newMsg.resultCount = 0
	newMsg.verbId = VerbGetEntityResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create a new message instance
func NewGetEntityResponseMessage(authToken, sessionId int64) *GetEntityResponseMessage {
	newMsg := DefaultGetEntityResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for GetEntityResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *GetEntityResponseMessage) GetHasResult() bool {
	return msg.hasResult
}

func (msg *GetEntityResponseMessage) GetResultId() int {
	return msg.resultId
}

func (msg *GetEntityResponseMessage) GetTotalCount() int {
	return msg.totalCount
}

func (msg *GetEntityResponseMessage) GetResultCount() int {
	return msg.resultCount
}

func (msg *GetEntityResponseMessage) GetEntityStream() tgdb.TGInputStream {
	return msg.entityStream
}

func (msg *GetEntityResponseMessage) SetHasResult(rFlag bool) {
	msg.hasResult = rFlag
}

func (msg *GetEntityResponseMessage) SetResultId(resultId int) {
	msg.resultId = resultId
}

func (msg *GetEntityResponseMessage) SetTotalCount(tCount int) {
	msg.totalCount = tCount
}

func (msg *GetEntityResponseMessage) SetResultCount(rCount int) {
	msg.resultCount = rCount
}

func (msg *GetEntityResponseMessage) SetEntityStream(eStream tgdb.TGInputStream) {
	msg.entityStream = eStream
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *GetEntityResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering GetEntityResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside GetEntityResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetEntityResponseMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetEntityResponseMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning GetEntityResponseMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *GetEntityResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering GetEntityResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetEntityResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetEntityResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning GetEntityResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *GetEntityResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *GetEntityResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *GetEntityResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *GetEntityResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *GetEntityResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *GetEntityResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *GetEntityResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *GetEntityResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *GetEntityResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *GetEntityResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *GetEntityResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *GetEntityResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *GetEntityResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *GetEntityResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *GetEntityResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *GetEntityResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *GetEntityResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *GetEntityResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("GetEntityResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("HasResult: %+v", msg.hasResult))
	buffer.WriteString(fmt.Sprintf(", ResultId: %d", msg.resultId))
	buffer.WriteString(fmt.Sprintf(", TotalCount: %d", msg.totalCount))
	buffer.WriteString(fmt.Sprintf(", ResultCount: %d", msg.resultCount))
	//buffer.WriteString(fmt.Sprintf(", EntityStream: %+v", msg.entityStream))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *GetEntityResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *GetEntityResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *GetEntityResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *GetEntityResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering GetEntityResponseMessage:ReadPayload"))
	}
	avail, err := is.(*ProtocolDataInputStream).Available()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityResponseMessage:ReadPayload w/ Error in reading available bytes from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetEntityResponseMessage:ReadPayload read avail as '%+v'", avail))
	}
	if avail == 0 {
		errMsg := fmt.Sprint("Get entity response has no data")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	msg.SetEntityStream(is)

	resultId, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityResponseMessage:ReadPayload w/ Error in reading resultId from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetEntityResponseMessage:ReadPayload read resultId as '%+v'", resultId))
	}

	pos := is.(*ProtocolDataInputStream).GetPosition()

	totCount, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetEntityResponseMessage:ReadPayload w/ Error in reading totCount from message buffer"))
		return err
	}
	if totCount > 0 {
		msg.SetHasResult(true)
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetEntityResponseMessage:ReadPayload read totCount as '%+v'", totCount))
	}

	is.(*ProtocolDataInputStream).SetPosition(pos)

	msg.SetResultId(resultId)
	msg.SetTotalCount(totCount)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning GetEntityResponseMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *GetEntityResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *GetEntityResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.hasResult, msg.resultId, msg.totalCount, msg.resultCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GetEntityResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *GetEntityResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.hasResult, &msg.resultId, &msg.totalCount, &msg.resultCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GetEntityResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type GetLargeObjectRequestMessage struct {
	*AbstractProtocolMessage
	entityId    int64
	decryptFlag bool
}

func DefaultGetLargeObjectRequestMessage() *GetLargeObjectRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GetEntityRequestMessage{})

	newMsg := GetLargeObjectRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.isUpdatable = true
	newMsg.entityId = 0
	newMsg.verbId = VerbGetLargeObjectRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewGetLargeObjectRequestMessage(authToken, sessionId int64) *GetLargeObjectRequestMessage {
	newMsg := DefaultGetLargeObjectRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for GetEntityRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *GetLargeObjectRequestMessage) GetDecryptFlag() bool {
	return msg.decryptFlag
}

func (msg *GetLargeObjectRequestMessage) GetEntityId() int64 {
	return msg.entityId
}

func (msg *GetLargeObjectRequestMessage) SetDecryption(flag bool) {
	msg.decryptFlag = flag
}

func (msg *GetLargeObjectRequestMessage) SetEntityId(id int64) {
	msg.entityId = id
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *GetLargeObjectRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering GetLargeObjectRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetLargeObjectRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetLargeObjectRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetLargeObjectRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning GetLargeObjectRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *GetLargeObjectRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering GetLargeObjectRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetLargeObjectRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetLargeObjectRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("GetLargeObjectRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *GetLargeObjectRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *GetLargeObjectRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *GetLargeObjectRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *GetLargeObjectRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *GetLargeObjectRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *GetLargeObjectRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *GetLargeObjectRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *GetLargeObjectRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *GetLargeObjectRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *GetLargeObjectRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *GetLargeObjectRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *GetLargeObjectRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *GetLargeObjectRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *GetLargeObjectRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *GetLargeObjectRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *GetLargeObjectRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *GetLargeObjectRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *GetLargeObjectRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("GetLargeObjectRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("EntityId: %d", msg.entityId))
	buffer.WriteString(fmt.Sprintf("DecryptFlag: %+v", msg.decryptFlag))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString() + "}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *GetLargeObjectRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *GetLargeObjectRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *GetLargeObjectRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *GetLargeObjectRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering GetLargeObjectRequestMessage:ReadPayload"))
	}
	entityId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectRequestMessage:ReadPayload w/ Error in reading entityId from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetLargeObjectRequestMessage:ReadPayload read entityId as '%+v'", entityId))
	}
	msg.SetEntityId(entityId)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning GetLargeObjectRequestMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *GetLargeObjectRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering GetLargeObjectRequestMessage:WritePayload at output buffer position: '%d'", startPos))
	}
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetEntityId())
	os.(*ProtocolDataOutputStream).WriteBoolean(msg.GetDecryptFlag())
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning GetLargeObjectRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *GetLargeObjectRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.entityId, msg.decryptFlag)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GetLargeObjectRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *GetLargeObjectRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.entityId, &msg.decryptFlag)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GetLargeObjectRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type GetLargeObjectResponseMessage struct {
	*AbstractProtocolMessage
	entityId int64
	boStream bytes.Buffer
}

func DefaultGetLargeObjectResponseMessage() *GetLargeObjectResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(GetEntityRequestMessage{})

	newMsg := GetLargeObjectResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.isUpdatable = true
	newMsg.entityId = 0
	//newMsg.boStream = new(bytes.Buffer)
	newMsg.verbId = VerbGetLargeObjectResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewGetLargeObjectResponseMessage(authToken, sessionId int64) *GetLargeObjectResponseMessage {
	newMsg := DefaultGetLargeObjectResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for GetEntityRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *GetLargeObjectResponseMessage) GetEntityId() int64 {
	return msg.entityId
}

func (msg *GetLargeObjectResponseMessage) GetBuffer() []byte {
	return msg.boStream.Bytes()
}

func (msg *GetLargeObjectResponseMessage) SetEntityId(id int64) {
	msg.entityId = id
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *GetLargeObjectResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering GetLargeObjectResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetLargeObjectResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetLargeObjectResponseMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside GetLargeObjectResponseMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning GetLargeObjectResponseMessage::FromBytes results in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *GetLargeObjectResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering GetLargeObjectResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetLargeObjectResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside GetLargeObjectResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning GetLargeObjectResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *GetLargeObjectResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *GetLargeObjectResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *GetLargeObjectResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *GetLargeObjectResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *GetLargeObjectResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *GetLargeObjectResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *GetLargeObjectResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *GetLargeObjectResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *GetLargeObjectResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *GetLargeObjectResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *GetLargeObjectResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *GetLargeObjectResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *GetLargeObjectResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *GetLargeObjectResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *GetLargeObjectResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *GetLargeObjectResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *GetLargeObjectResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *GetLargeObjectResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("GetLargeObjectResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("EntityId: %d", msg.entityId))
	buffer.WriteString(fmt.Sprintf(", BoStream: %s", msg.boStream.String()))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *GetLargeObjectResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *GetLargeObjectResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *GetLargeObjectResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *GetLargeObjectResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering GetLargeObjectResponseMessage:ReadPayload"))
	}
	status, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectResponseMessage:ReadPayload w/ Error in reading status from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetLargeObjectResponseMessage:ReadPayload read status as '%+v'", status))
	}
	if status > 0 {
		errMsg := fmt.Sprintf("Read Large Object failed with status: %d", status)
		return GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Read the chunks
	entityId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectResponseMessage:ReadPayload w/ Error in reading entityId from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetLargeObjectResponseMessage:ReadPayload read entityId as '%+v'", entityId))
	}

	bHasData, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectResponseMessage:ReadPayload w/ Error in reading data availability flag from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside GetLargeObjectResponseMessage:ReadPayload read bHasData as '%+v'", bHasData))
	}
	if bHasData {
		numChunks, err := is.(*ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectResponseMessage:ReadPayload w/ Error in reading chunk count from message buffer"))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside GetLargeObjectResponseMessage:ReadPayload read numChunks as '%+v'", numChunks))
		}
		for i := 0; i < numChunks; i++ {
			chunk, err := is.(*ProtocolDataInputStream).ReadBytes()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning GetLargeObjectResponseMessage:ReadPayload w/ Error in reading chunk bytes from message buffer"))
				return err
			}
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("GetLargeObjectResponseMessage:ReadPayload read chunk as '%+v'", chunk))
			}
			msg.boStream.Write(chunk)
		}
	}

	msg.SetEntityId(entityId)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning GetLargeObjectResponseMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *GetLargeObjectResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	//os.(*ProtocolDataOutputStream).WriteLong(msg.GetEntityId())
	// No-op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *GetLargeObjectResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.entityId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GetLargeObjectResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *GetLargeObjectResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.entityId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GetLargeObjectResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


const (
	InvalidRequest = iota
	InitiateRequest
	ChallengeAccepted
)

type HandShakeRequestMessage struct {
	*AbstractProtocolMessage
	sslMode       bool
	challenge     int64
	handshakeType int
	version       int64
}

func DefaultHandShakeRequestMessage() *HandShakeRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(HandShakeRequestMessage{})

	newMsg := HandShakeRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.isUpdatable = true
	newMsg.verbId = VerbHandShakeRequest
	newMsg.sslMode = false
	newMsg.challenge = 0
	newMsg.version = 0
	newMsg.handshakeType = InvalidRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewHandShakeRequestMessage(authToken, sessionId int64) *HandShakeRequestMessage {
	newMsg := DefaultHandShakeRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for HandShakeRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *HandShakeRequestMessage) GetSslMode() bool {
	return msg.sslMode
}

func (msg *HandShakeRequestMessage) GetChallenge() int64 {
	return msg.challenge
}

func (msg *HandShakeRequestMessage) GetRequestType() int {
	return msg.handshakeType
}

func (msg *HandShakeRequestMessage) GetVersion() int64 {
	return msg.version
}

func (msg *HandShakeRequestMessage) SetSslMode(mode bool) {
	msg.sslMode = mode
}

func (msg *HandShakeRequestMessage) SetChallenge(challenge int64) {
	msg.challenge = challenge
}

func (msg *HandShakeRequestMessage) SetRequestType(rType int) {
	msg.handshakeType = rType
}

func (msg *HandShakeRequestMessage) SetVersion(version int64) {
	msg.version = version
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *HandShakeRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering HandShakeRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside HandShakeRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside HandShakeRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside HandShakeRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning HandShakeRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *HandShakeRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering HandShakeRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside HandShakeRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside HandShakeRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning HandShakeRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *HandShakeRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *HandShakeRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *HandShakeRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *HandShakeRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *HandShakeRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *HandShakeRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *HandShakeRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *HandShakeRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *HandShakeRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *HandShakeRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *HandShakeRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *HandShakeRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *HandShakeRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *HandShakeRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *HandShakeRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *HandShakeRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *HandShakeRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *HandShakeRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("HandShakeRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("SslMode: %+v", msg.sslMode))
	buffer.WriteString(fmt.Sprintf(", Challenge: %d", msg.challenge))
	buffer.WriteString(fmt.Sprintf(", HandshakeType: %d", msg.handshakeType))
	buffer.WriteString(fmt.Sprintf(", Version: %d", msg.version))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *HandShakeRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *HandShakeRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *HandShakeRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *HandShakeRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering HandShakeRequestMessage:ReadPayload"))
	}
	//For Testing purpose only.
	rType, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:ReadPayload w/ Error in reading rType from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside HandShakeRequestMessage:ReadPayload read rType as '%+v'", rType))
	}

	mode, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:ReadPayload w/ Error in reading mode from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside HandShakeRequestMessage:ReadPayload read mode as '%+v'", mode))
	}

	challenge, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:ReadPayload w/ Error in reading challenge from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside HandShakeRequestMessage:ReadPayload read challenge as '%+v'", challenge))
	}

	msg.SetRequestType(int(rType))
	msg.SetSslMode(mode)
	msg.SetChallenge(challenge)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning HandShakeRequestMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *HandShakeRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering HandShakeRequestMessage:WritePayload at output buffer position: '%d'", startPos))
	}
	os.(*ProtocolDataOutputStream).WriteByte(msg.GetRequestType())
	os.(*ProtocolDataOutputStream).WriteBoolean(msg.GetSslMode())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetChallenge())
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning HandShakeRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaler
/////////////////////////////////////////////////////////////////

func (msg *HandShakeRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.sslMode, msg.challenge, msg.handshakeType, msg.version)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning HandShakeRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaler
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *HandShakeRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.sslMode, &msg.challenge, &msg.handshakeType, &msg.version)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning HandShakeRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



const (
	ResponseInvalid = iota
	ResponseAcceptChallenge
	ResponseProceedWithAuthentication
	ResponseChallengeFailed
)

type HandShakeResponseMessage struct {
	*AbstractProtocolMessage
	challenge      int64
	responseStatus int
	version        int64
	errorMessage   string
}

func DefaultHandShakeResponseMessage() *HandShakeResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(HandShakeResponseMessage{})

	newMsg := HandShakeResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.challenge = 0
	newMsg.version = 0
	newMsg.responseStatus = ResponseInvalid
	newMsg.verbId = VerbHandShakeResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewHandShakeResponseMessage(authToken, sessionId int64) *HandShakeResponseMessage {
	newMsg := DefaultHandShakeResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for HandShakeResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *HandShakeResponseMessage) GetChallenge() int64 {
	return msg.challenge
}

func (msg *HandShakeResponseMessage) GetErrorMessage() string {
	return msg.errorMessage
}

func (msg *HandShakeResponseMessage) GetResponseStatus() int {
	return msg.responseStatus
}

func (msg *HandShakeResponseMessage) GetVersion() int64 {
	return msg.version
}

func (msg *HandShakeResponseMessage) SetChallenge(challenge int64) {
	msg.challenge = challenge
}

func (msg *HandShakeResponseMessage) SetErrorMessage(errMsg string) {
	msg.errorMessage = errMsg
}

func (msg *HandShakeResponseMessage) SetResponseStatus(rStatus int) {
	msg.responseStatus = rStatus
}

func (msg *HandShakeResponseMessage) SetVersion(version int64) {
	msg.version = version
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *HandShakeResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering HandShakeResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering HandShakeResponseMessage:FromBytes - received input buffer as '%+v'", buffer))
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside HandShakeResponseMessage:FromBytes read bufLen as '%d'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside HandShakeResponseMessage:FromBytes about to read Header data elements"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside HandShakeResponseMessage:FromBytes about to read Payload data elements"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning HandShakeResponseMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *HandShakeResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering HandShakeResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside HandShakeResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside HandShakeResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning HandShakeResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *HandShakeResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *HandShakeResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *HandShakeResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *HandShakeResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *HandShakeResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *HandShakeResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *HandShakeResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *HandShakeResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *HandShakeResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *HandShakeResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *HandShakeResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *HandShakeResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *HandShakeResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *HandShakeResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *HandShakeResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *HandShakeResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *HandShakeResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *HandShakeResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("HandShakeResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("Challenge: %d ", msg.challenge))
	buffer.WriteString(fmt.Sprintf(", ResponseStatus: %d ", msg.responseStatus))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *HandShakeResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *HandShakeResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *HandShakeResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *HandShakeResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering HandShakeResponseMessage:ReadPayload"))
	}
	rStatus, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeResponseMessage:ReadPayload w/ Error in reading rStatus from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside HandShakeResponseMessage:ReadPayload read rStatus as '%+v'", rStatus))
	}

	challenge, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeResponseMessage:ReadPayload w/ Error in reading challenge from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside HandShakeResponseMessage:ReadPayload read challenge as '%+v'", challenge))
	}

	if int(rStatus) == ResponseChallengeFailed {
		errMsgBytes, err := is.(*ProtocolDataInputStream).ReadBytes()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning HandShakeResponseMessage:ReadPayload w/ Error in reading errMsgBytes from message buffer"))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside HandShakeResponseMessage:ReadPayload read rStatus as '%+v'", errMsgBytes))
		}
		msg.SetErrorMessage(string(errMsgBytes))
	}
	msg.SetResponseStatus(int(rStatus))
	msg.SetChallenge(challenge)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning HandShakeResponseMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *HandShakeResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering HandShakeResponseMessage:WritePayload at output buffer position: '%d'", startPos))
	}
	//This is purely for testing. Client never writes out the response.
	os.(*ProtocolDataOutputStream).WriteByte(msg.GetResponseStatus())
	os.(*ProtocolDataOutputStream).WriteLong(msg.GetChallenge())
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning HandShakeResponseMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *HandShakeResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.challenge, msg.responseStatus,
		msg.version, msg.errorMessage)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning HandShakeResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *HandShakeResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.challenge, &msg.responseStatus, &msg.version, &msg.errorMessage)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning HandShakeResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type InvalidMessage struct {
	*AbstractProtocolMessage
}

func DefaultInvalidMessage() *InvalidMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(InvalidMessage{})

	newMsg := InvalidMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbInvalidMessage
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewInvalidMessage(authToken, sessionId int64) *InvalidMessage {
	newMsg := DefaultInvalidMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *InvalidMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering InvalidMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning InvalidMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning InvalidMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside InvalidMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside InvalidMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside InvalidMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning InvalidMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *InvalidMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering InvalidMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside InvalidMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside InvalidMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning InvalidMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning InvalidMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *InvalidMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *InvalidMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *InvalidMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *InvalidMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *InvalidMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *InvalidMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *InvalidMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *InvalidMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *InvalidMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *InvalidMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *InvalidMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *InvalidMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *InvalidMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *InvalidMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *InvalidMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *InvalidMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *InvalidMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *InvalidMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("InvalidMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *InvalidMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *InvalidMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *InvalidMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *InvalidMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *InvalidMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *InvalidMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning InvalidMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *InvalidMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning InvalidMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type MetadataRequest struct {
	*AbstractProtocolMessage
}

func DefaultMetadataRequestMessage() *MetadataRequest {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(MetadataRequest{})

	newMsg := MetadataRequest{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbMetadataRequest
	newMsg.isUpdatable = true
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewMetadataRequestMessage(authToken, sessionId int64) *MetadataRequest {
	newMsg := DefaultMetadataRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *MetadataRequest) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering MetadataRequest:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataRequest:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataRequest:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside MetadataRequest:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside MetadataRequest:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside MetadataRequest:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning MetadataRequest::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *MetadataRequest) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering MetadataRequest:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside MetadataRequest:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside MetadataRequest:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataRequest:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning MetadataRequest::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *MetadataRequest) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *MetadataRequest) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *MetadataRequest) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *MetadataRequest) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *MetadataRequest) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *MetadataRequest) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *MetadataRequest) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *MetadataRequest) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *MetadataRequest) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *MetadataRequest) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *MetadataRequest) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *MetadataRequest) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *MetadataRequest) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *MetadataRequest) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *MetadataRequest) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *MetadataRequest) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *MetadataRequest) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *MetadataRequest) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("MetadataRequest:{")
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *MetadataRequest) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *MetadataRequest) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *MetadataRequest) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *MetadataRequest) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *MetadataRequest) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *MetadataRequest) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MetadataRequest:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *MetadataRequest) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MetadataRequest:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type MetadataResponse struct {
	*AbstractProtocolMessage
	attrDescList []tgdb.TGAttributeDescriptor
	nodeTypeList []tgdb.TGNodeType
	edgeTypeList []tgdb.TGEdgeType
}

func DefaultMetadataResponseMessage() *MetadataResponse {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(MetadataResponse{})

	newMsg := MetadataResponse{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbMetadataResponse
	newMsg.attrDescList = make([]tgdb.TGAttributeDescriptor, 0)
	newMsg.nodeTypeList = make([]tgdb.TGNodeType, 0)
	newMsg.edgeTypeList = make([]tgdb.TGEdgeType, 0)
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewMetadataResponseMessage(authToken, sessionId int64) *MetadataResponse {
	newMsg := DefaultMetadataResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for VerbMetadataResponse
/////////////////////////////////////////////////////////////////

func (msg *MetadataResponse) GetAttrDescList() []tgdb.TGAttributeDescriptor {
	return msg.attrDescList
}

func (msg *MetadataResponse) GetNodeTypeList() []tgdb.TGNodeType {
	return msg.nodeTypeList
}

func (msg *MetadataResponse) GetEdgeTypeList() []tgdb.TGEdgeType {
	return msg.edgeTypeList
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *MetadataResponse) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering MetadataResponse:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside MetadataResponse:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside MetadataResponse:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside MetadataResponse:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning MetadataResponse::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *MetadataResponse) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering MetadataResponse:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside MetadataResponse:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside MetadataResponse:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning MetadataResponse::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *MetadataResponse) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *MetadataResponse) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *MetadataResponse) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *MetadataResponse) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *MetadataResponse) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *MetadataResponse) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *MetadataResponse) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *MetadataResponse) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *MetadataResponse) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *MetadataResponse) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *MetadataResponse) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *MetadataResponse) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *MetadataResponse) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *MetadataResponse) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *MetadataResponse) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *MetadataResponse) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *MetadataResponse) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *MetadataResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("MetadataResponse:{")
	buffer.WriteString(fmt.Sprint("AttrDescSet:{"))
	for _, d := range msg.attrDescList {
		buffer.WriteString(fmt.Sprintf("Attribute Descriptor: %+v ", d))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", NodeTypeList:{"))
	for _, v := range msg.nodeTypeList {
		buffer.WriteString(fmt.Sprintf("NodeType: %+v ", v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", EdgeTypeList:{"))
	for _, v := range msg.edgeTypeList {
		buffer.WriteString(fmt.Sprintf("EdgeType: %+v ", v))
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
func (msg *MetadataResponse) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *MetadataResponse) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *MetadataResponse) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *MetadataResponse) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering MetadataResponse:ReadPayload"))
	}
	avail, err := is.(*ProtocolDataInputStream).Available()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading available bytes from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read avail as '%+v'", avail))
	}
	if avail == 0 {
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Metadata response has no data"))
		}
		errMsg := fmt.Sprint("Metadata Response has no data")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	count, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading count from message buffer"))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read count as '%+v'", count))
	}
	for {
		if count <= 0 {
			break
		}

		sysType, err := is.(*ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading SysType from message buffer"))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read SysType as '%+v'", sysType))
		}

		typeCount, err := is.(*ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading typeCount from message buffer"))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read typeCount as '%+v'", typeCount))
		}

		if tgdb.TGSystemType(sysType) == tgdb.SystemTypeAttributeDescriptor {
			for i := 0; i < typeCount; i++ {
				attrDesc := NewAttributeDescriptorWithType("temp", AttributeTypeString)
				err := attrDesc.ReadExternal(is)
				if err != nil {
					logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading AttrDesc from message buffer"))
					return err
				}
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read AttrDesc as '%+v'", attrDesc))
				}
				msg.attrDescList = append(msg.attrDescList, attrDesc)
			}
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read '%d' AttrDesc and assigned as '%+v'", typeCount, msg.attrDescList))
			}
		} else if tgdb.TGSystemType(sysType) == tgdb.SystemTypeNode {
			for i := 0; i < typeCount; i++ {
				nodeType := NewNodeType("temp", nil)
				err := nodeType.ReadExternal(is)
				if err != nil {
					logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading node from message buffer"))
					return err
				}
				name := nodeType.GetName()
				if strings.HasPrefix(name, "@") || strings.HasPrefix(name, "$") {
					continue
				}
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read nodeType as '%+v'", nodeType))
				}
				msg.nodeTypeList = append(msg.nodeTypeList, nodeType)
			}
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read '%d' nodes and assigned as '%+v'", typeCount, msg.nodeTypeList))
			}
		} else if tgdb.TGSystemType(sysType) == tgdb.SystemTypeEdge {
			for i := 0; i < typeCount; i++ {
				edgeType := NewEdgeType("temp", tgdb.DirectionTypeBiDirectional, nil)
				err := edgeType.ReadExternal(is)
				if err != nil {
					logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading edge from message buffer"))
					return err
				}
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read edgeType as '%+v'", edgeType))
				}
				msg.edgeTypeList = append(msg.edgeTypeList, edgeType)
			}
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside MetadataResponse:ReadPayload read '%d' edges and assigned as '%+v'", typeCount, msg.edgeTypeList))
			}
		} else {
			logger.Warning(fmt.Sprintf("WARNING: MetadataResponse:ReadPayload - Invalid metadata desc '%d' received", sysType))
			//TODO: Revisit later - Do we need to throw exception?
		}
		count -= typeCount
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning MetadataResponse:ReadPayload w/ MedataResponse as '%+v'", msg))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *MetadataResponse) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *MetadataResponse) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.attrDescList,
		msg.nodeTypeList, msg.edgeTypeList)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MetadataResponse:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *MetadataResponse) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.attrDescList, &msg.nodeTypeList, &msg.edgeTypeList)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MetadataResponse:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



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
func (msg *PingMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering PingMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning PingMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning PingMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside PingMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside PingMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside PingMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning PingMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *PingMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering PingMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside PingMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside PingMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning PingMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning PingMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *PingMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *PingMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *PingMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *PingMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *PingMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *PingMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *PingMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *PingMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *PingMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *PingMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *PingMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *PingMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *PingMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *PingMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *PingMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *PingMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *PingMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *PingMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("PingMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *PingMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *PingMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *PingMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *PingMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *PingMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
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

//var logger = logging.DefaultTGLogManager().GetLogger()

/////////////////////////////////////////////////////////////////
// Helper Public functions for TGMessageFactory
/////////////////////////////////////////////////////////////////

// Create new message instance based on the input type
func CreateMessageForVerb(verbId int) (tgdb.TGMessage, tgdb.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputVerbId := verbId

	// Use a switch case to switch between message types, if a type exist then error is nil (null)
	// Whenever new message type gets into the mix, just add a case below
	switch inputVerbId {
	case VerbPingMessage:
		return DefaultPingMessage(), nil
	case VerbHandShakeRequest:
		return DefaultHandShakeRequestMessage(), nil
	case VerbHandShakeResponse:
		return DefaultHandShakeResponseMessage(), nil
	case VerbAuthenticateRequest:
		return DefaultAuthenticateRequestMessage(), nil
	case VerbAuthenticateResponse:
		return DefaultAuthenticateResponseMessage(), nil
	case VerbBeginTransactionRequest:
		return DefaultBeginTransactionRequestMessage(), nil
	case VerbBeginTransactionResponse:
		return DefaultBeginTransactionResponseMessage(), nil
	case VerbCommitTransactionRequest:
		return DefaultCommitTransactionRequestMessage(), nil
	case VerbCommitTransactionResponse:
		return DefaultCommitTransactionResponseMessage(), nil
	case VerbRollbackTransactionRequest:
		return DefaultRollbackTransactionRequestMessage(), nil
	case VerbRollbackTransactionResponse:
		return DefaultRollbackTransactionResponseMessage(), nil
	case VerbQueryRequest:
		return DefaultQueryRequestMessage(), nil
	case VerbQueryResponse:
		return DefaultQueryResponseMessage(), nil
	case VerbTraverseRequest:
		return DefaultTraverseRequestMessage(), nil
	case VerbTraverseResponse:
		return DefaultTraverseResponseMessage(), nil
	case VerbAdminRequest:
		return DefaultAdminRequestMessage(), nil
	case VerbAdminResponse:
		return DefaultAdminResponseMessage(), nil
	case VerbMetadataRequest:
		return DefaultMetadataRequestMessage(), nil
	case VerbMetadataResponse:
		return DefaultMetadataResponseMessage(), nil
	case VerbGetEntityRequest:
		return DefaultGetEntityRequestMessage(), nil
	case VerbGetEntityResponse:
		return DefaultGetEntityResponseMessage(), nil
	case VerbGetLargeObjectRequest:
		return DefaultGetLargeObjectRequestMessage(), nil
	case VerbGetLargeObjectResponse:
		return DefaultGetLargeObjectResponseMessage(), nil
	case VerbDumpStacktraceRequest:
		fallthrough
		//return DefaultDumpStacktraceRequestMessage(), nil
	case VerbDisconnectChannelRequest:
		return DefaultDisconnectChannelRequestMessage(), nil
	case VerbSessionForcefullyTerminated:
		return DefaultSessionForcefullyTerminatedMessage(), nil
	case VerbDecryptBufferRequest:
		return DefaultDecryptBufferRequestMessage(), nil
	case VerbDecryptBufferResponse:
		return DefaultDecryptBufferResponseMessage(), nil
	case VerbExceptionMessage:
		return DefaultExceptionMessage(), nil
	case VerbInvalidMessage:
		return DefaultInvalidMessage(), nil
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Message Type '%s'", GetVerb(inputVerbId).name)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
}

func CreateMessageWithToken(verbId int, authToken, sessionId int64) (tgdb.TGMessage, tgdb.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputVerbId := verbId

	// Use a switch case to switch between message types, if a type exist then error is nil (null)
	// Whenever new message type gets into the mix, just add a case below
	switch inputVerbId {
	case VerbPingMessage:
		return NewPingMessage(authToken, sessionId), nil
	case VerbHandShakeRequest:
		return NewHandShakeRequestMessage(authToken, sessionId), nil
	case VerbHandShakeResponse:
		return NewHandShakeResponseMessage(authToken, sessionId), nil
	case VerbAuthenticateRequest:
		return NewAuthenticateRequestMessage(authToken, sessionId), nil
	case VerbAuthenticateResponse:
		return NewAuthenticateResponseMessage(authToken, sessionId), nil
	case VerbBeginTransactionRequest:
		return NewBeginTransactionRequestMessage(authToken, sessionId), nil
	case VerbBeginTransactionResponse:
		return NewBeginTransactionResponseMessage(authToken, sessionId), nil
	case VerbCommitTransactionRequest:
		return NewCommitTransactionRequestMessage(authToken, sessionId), nil
	case VerbCommitTransactionResponse:
		return NewCommitTransactionResponseMessage(authToken, sessionId), nil
	case VerbRollbackTransactionRequest:
		return NewRollbackTransactionRequestMessage(authToken, sessionId), nil
	case VerbRollbackTransactionResponse:
		return NewRollbackTransactionResponseMessage(authToken, sessionId), nil
	case VerbQueryRequest:
		return NewQueryRequestMessage(authToken, sessionId), nil
	case VerbQueryResponse:
		return NewQueryResponseMessage(authToken, sessionId), nil
	case VerbTraverseRequest:
		return NewTraverseRequestMessage(authToken, sessionId), nil
	case VerbTraverseResponse:
		return NewTraverseResponseMessage(authToken, sessionId), nil
	case VerbAdminRequest:
		return NewAdminRequestMessage(authToken, sessionId), nil
	case VerbAdminResponse:
		return NewAdminResponseMessage(authToken, sessionId), nil
	case VerbMetadataRequest:
		return NewMetadataRequestMessage(authToken, sessionId), nil
	case VerbMetadataResponse:
		return NewMetadataResponseMessage(authToken, sessionId), nil
	case VerbGetEntityRequest:
		return NewGetEntityRequestMessage(authToken, sessionId), nil
	case VerbGetEntityResponse:
		return NewGetEntityResponseMessage(authToken, sessionId), nil
	case VerbGetLargeObjectRequest:
		return NewGetLargeObjectRequestMessage(authToken, sessionId), nil
	case VerbGetLargeObjectResponse:
		return NewGetLargeObjectResponseMessage(authToken, sessionId), nil
	case VerbDumpStacktraceRequest:
		fallthrough
		//return NewDumpStacktraceRequestMessage(authToken, sessionId), nil
	case VerbDisconnectChannelRequest:
		return NewDisconnectChannelRequestMessage(authToken, sessionId), nil
	case VerbSessionForcefullyTerminated:
		return NewSessionForcefullyTerminatedMessage(authToken, sessionId), nil
	case VerbDecryptBufferRequest:
		return NewDecryptBufferRequestMessage(authToken, sessionId), nil
	case VerbDecryptBufferResponse:
		return NewDecryptBufferResponseMessage(authToken, sessionId), nil
	case VerbExceptionMessage:
		return NewExceptionMessage(authToken, sessionId), nil
	case VerbInvalidMessage:
		return NewInvalidMessage(authToken, sessionId), nil
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Message Type '%s'", GetVerb(inputVerbId).name)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
}

func CreateMessageFromBuffer(buffer []byte, offset int, length int) (tgdb.TGMessage, tgdb.TGError) {
	buf := make([]byte, 0)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering MessageFactory::CreateMessageFromBuffer received buffer(BufLen: %d, Offset: %d, Len: %d)", len(buffer), offset, length))
	}

	if len(buffer) == length {
		buf = buffer
	} else {
		buf = append(buffer[offset:length])
	}
	//logger.Debug(fmt.Sprintf("Inside MessageFactory::CreateMessageFromBuffer buf is '%+v'", buf))

	commandVerb, err := VerbIdFromBytes(buf)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MessageFactory::CreateMessageFromBuffer w/ error in extracting verbId from message buffer: %s", err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside MessageFactory::CreateMessageFromBuffer retrieved VerbId '%s'", commandVerb.name))
	}
	msg, err := CreateMessageForVerb(commandVerb.id)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MessageFactory::CreateMessageFromBuffer w/ error in creating message for verb('%s'): %s", commandVerb.name, err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside MessageFactory::CreateMessageFromBuffer created '%+v'", msg))
	}
	msg1, err := msg.FromBytes(buffer)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MessageFactory::CreateMessageFromBuffer w/ error in updating message contents from buffer: %s", err.Error()))
		return nil, err
	}
	return msg1, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

func FromBytes(verbId int, buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return nil, err
	}
	// Execute Individual Message's method
	return msg.FromBytes(buffer)
}

func ToBytes(verbId int) ([]byte, int, tgdb.TGError) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return nil, -1, err
	}
	// Execute Individual Message's method
	return msg.ToBytes()
}

// GetAuthToken gets Authorization Token for specified message type
func GetAuthToken(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetAuthToken()
}

func GetIsUpdatable(verbId int) bool {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetIsUpdatable()
}

func GetMessageByteBufLength(verbId int) int {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetMessageByteBufLength()
}

func GetRequestId(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetRequestId()
}

func GetSequenceNo(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetSequenceNo()
}

// GetSessionId gets Session id for specified message type
func GetSessionId(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetSessionId()
}

func GetTimestamp(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetTimestamp()
}

func GetVerbId(verbId int) int {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetVerbId()
}

func SetAuthToken(verbId int, authToken int64) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return
	}
	// Execute Derived / Dependent Message's method
	msg.SetAuthToken(authToken)
}

//func SetIsUpdatable(verbId int, updateFlag bool) {
//	msg, err := CreateMessageForVerb(verbId)
//	if err != nil {
//		return
//	}
//	// Execute Derived / Dependent Message's method
//	return msg.SetIsUpdatable(updateFlag)
//}

//func SetMessageByteBufLength(verbId int, bufLength int) {
//	msg, err := CreateMessageForVerb(verbId)
//	if err != nil {
//		return
//	}
//	// Execute Derived / Dependent Message's method
//	return msg.SetMessageByteBufLength(bufLength)
//}

func SetRequestId(verbId int, requestId int64) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return
	}
	// Execute Derived / Dependent Message's method
	msg.SetRequestId(requestId)
}

//func SetSequenceNo(verbId int, sequenceNo int64) {
//	msg, err := CreateMessageForVerb(verbId)
//	if err != nil {
//		return
//	}
//	// Execute Derived / Dependent Message's method
//	return msg.SetSequenceNo(sequenceNo)
//}

func SetSessionId(verbId int, sessionId int64) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return
	}
	// Execute Derived / Dependent Message's method
	msg.SetSessionId(sessionId)
}

func SetTimestamp(verbId int, timestamp int64) tgdb.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Derived / Dependent Message's method
	return msg.SetTimestamp(timestamp)
}

//func SetVerbId(verbId int) {
//	msg, err := CreateMessageForVerb(verbId)
//	if err != nil {
//		return
//	}
//	// Execute Derived / Dependent Message's method
//	return msg.SetVerbId(verbId)
//}

func UpdateSequenceAndTimeStamp(verbId int, timestamp int64) tgdb.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Individual Message's method
	return msg.UpdateSequenceAndTimeStamp(timestamp)
}

func ReadHeader(verbId int, is tgdb.TGInputStream) tgdb.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	//
	// Execute Individual Message's method
	return msg.ReadHeader(is)
}

func WriteHeader(verbId int, os tgdb.TGOutputStream) tgdb.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Individual Message's method
	return msg.WriteHeader(os)
}

func ReadPayload(verbId int, is tgdb.TGInputStream) tgdb.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Individual Message's method
	return msg.ReadPayload(is)
}

func WritePayload(verbId int, os tgdb.TGOutputStream) tgdb.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Individual Message's method
	return msg.WritePayload(os)
}


type QueryRequestMessage struct {
	*AbstractProtocolMessage
	queryExpr       string
	edgeExpr        string
	traverseExpr    string
	endExpr         string
	queryHashId     int64
	command         int
	queryObject     tgdb.TGQuery
	fetchSize       int
	batchSize       int
	traversalDepth  int
	edgeLimit       int
	sortAttrName    string
	sortOrderDsc    bool
	sortResultLimit int
}

func DefaultQueryRequestMessage() *QueryRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(QueryRequestMessage{})

	newMsg := QueryRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.fetchSize = 1000
	newMsg.batchSize = 50
	newMsg.traversalDepth = 3
	newMsg.edgeLimit = 0
	newMsg.sortOrderDsc = false
	newMsg.sortResultLimit = 50
	newMsg.verbId = VerbQueryRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewQueryRequestMessage(authToken, sessionId int64) *QueryRequestMessage {
	newMsg := DefaultQueryRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for QueryRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *QueryRequestMessage) GetBatchSize() int {
	return msg.batchSize
}

func (msg *QueryRequestMessage) GetEdgeLimit() int {
	return msg.edgeLimit
}

func (msg *QueryRequestMessage) GetFetchSize() int {
	return msg.fetchSize
}

func (msg *QueryRequestMessage) GetSortAttrName() string {
	return msg.sortAttrName
}

func (msg *QueryRequestMessage) GetSortOrderDsc() bool {
	return msg.sortOrderDsc
}

func (msg *QueryRequestMessage) GetSortResultLimit() int {
	return msg.sortResultLimit
}

func (msg *QueryRequestMessage) GetTraversalDepth() int {
	return msg.traversalDepth
}

func (msg *QueryRequestMessage) GetQuery() string {
	return msg.queryExpr
}

func (msg *QueryRequestMessage) GetEdgeFilter() string {
	return msg.edgeExpr
}

func (msg *QueryRequestMessage) GetTraversalCondition() string {
	return msg.traverseExpr
}

func (msg *QueryRequestMessage) GetEndCondition() string {
	return msg.endExpr
}

func (msg *QueryRequestMessage) GetQueryHashId() int64 {
	return msg.queryHashId
}

func (msg *QueryRequestMessage) GetCommand() int {
	return msg.command
}

func (msg *QueryRequestMessage) GetQueryObject() tgdb.TGQuery {
	return msg.queryObject
}

func (msg *QueryRequestMessage) SetBatchSize(size int) {
	if size < 10 || size > 32767 {
		msg.batchSize = 50
	} else {
		msg.batchSize = size
	}
}

func (msg *QueryRequestMessage) SetEdgeLimit(size int) {
	if size < 0 || size > 32767 {
		msg.edgeLimit = 1000
	} else {
		msg.edgeLimit = size
	}
}

func (msg *QueryRequestMessage) SetFetchSize(size int) {
	if size < 0 {
		msg.fetchSize = 1000
	} else {
		msg.fetchSize = size
	}
}

func (msg *QueryRequestMessage) SetSortAttrName(name string) {
	if len(name) != 0 {
		msg.sortAttrName = name
	}
}

func (msg *QueryRequestMessage) SetSortOrderDsc(order bool) {
	msg.sortOrderDsc = order
}

func (msg *QueryRequestMessage) SetSortResultLimit(limit int) {
	if limit < 0 {
		msg.sortResultLimit = 50
	} else {
		msg.sortResultLimit = limit
	}
}

func (msg *QueryRequestMessage) SetTraversalDepth(depth int) {
	if depth < 1 || depth > 1000 {
		msg.traversalDepth = 3
	} else {
		msg.traversalDepth = depth
	}
}

func (msg *QueryRequestMessage) SetQuery(expr string) {
	msg.queryExpr = expr
}

func (msg *QueryRequestMessage) SetEdgeFilter(expr string) {
	msg.edgeExpr = expr
}

func (msg *QueryRequestMessage) SetTraversalCondition(expr string) {
	msg.traverseExpr = expr
}

func (msg *QueryRequestMessage) SetEndCondition(expr string) {
	msg.endExpr = expr
}

func (msg *QueryRequestMessage) SetQueryHashId(hash int64) {
	msg.queryHashId = hash
}

func (msg *QueryRequestMessage) SetCommand(cmd int) {
	msg.command = cmd
}

func (msg *QueryRequestMessage) SetQueryObject(query tgdb.TGQuery) {
	msg.queryObject = query
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *QueryRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering QueryRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning QueryRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside QueryRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside QueryRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside QueryRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning QueryRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *QueryRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering QueryRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside QueryRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside QueryRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning QueryRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *QueryRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *QueryRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *QueryRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *QueryRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *QueryRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *QueryRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *QueryRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *QueryRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *QueryRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *QueryRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *QueryRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *QueryRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *QueryRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *QueryRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *QueryRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *QueryRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *QueryRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *QueryRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("QueryRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("QueryExpr: %s", msg.queryExpr))
	buffer.WriteString(fmt.Sprintf(", EdgeExpr: %s", msg.edgeExpr))
	buffer.WriteString(fmt.Sprintf(", TraverseExpr: %s", msg.traverseExpr))
	buffer.WriteString(fmt.Sprintf(", EndExpr: %s", msg.endExpr))
	buffer.WriteString(fmt.Sprintf(", QueryHashId: %d", msg.queryHashId))
	buffer.WriteString(fmt.Sprintf(", Command: %d", msg.command))
	//buffer.WriteString(fmt.Sprintf(", QueryObject: %+v", msg.queryObject))
	buffer.WriteString(fmt.Sprintf(", FetchSize: %d", msg.fetchSize))
	buffer.WriteString(fmt.Sprintf(", BatchSize: %d", msg.batchSize))
	buffer.WriteString(fmt.Sprintf(", TraversalDepth: %d", msg.traversalDepth))
	buffer.WriteString(fmt.Sprintf(", EdgeLimit: %d", msg.edgeLimit))
	buffer.WriteString(fmt.Sprintf(", SortAttrName: %s", msg.sortAttrName))
	buffer.WriteString(fmt.Sprintf(", SortOrderDsc: %+v", msg.sortOrderDsc))
	buffer.WriteString(fmt.Sprintf(", SortResultLimit: %d", msg.sortResultLimit))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *QueryRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *QueryRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *QueryRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *QueryRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *QueryRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	startPos := os.(*ProtocolDataOutputStream).GetPosition()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering QueryRequestMessage::WritePayload at output buffer position at: %d", startPos))
	}
	os.(*ProtocolDataOutputStream).WriteInt(1) // datalength
	os.(*ProtocolDataOutputStream).WriteInt(1) //checksum
	os.(*ProtocolDataOutputStream).WriteInt(msg.GetCommand())
	os.(*ProtocolDataOutputStream).WriteInt(msg.GetFetchSize())
	os.(*ProtocolDataOutputStream).WriteShort(msg.GetBatchSize())
	os.(*ProtocolDataOutputStream).WriteShort(msg.GetTraversalDepth())
	os.(*ProtocolDataOutputStream).WriteShort(msg.GetEdgeLimit())
	// Has sort attr
	if msg.GetSortAttrName() != "" {
		os.(*ProtocolDataOutputStream).WriteBoolean(true)
		err := os.(*ProtocolDataOutputStream).WriteUTF(msg.GetSortAttrName())
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning QueryRequestMessage:WritePayload w/ Error in writing sortAttrName to message buffer"))
			return err
		}
		os.(*ProtocolDataOutputStream).WriteBoolean(msg.GetSortOrderDsc())
		os.(*ProtocolDataOutputStream).WriteInt(msg.GetSortResultLimit())
	} else {
		os.(*ProtocolDataOutputStream).WriteBoolean(false)
	}

	// CREATE, EXECUTE.
	if msg.GetCommand() == 1 || msg.GetCommand() == 2 || msg.GetCommand() == 3 || msg.GetCommand() == 4 {
		if msg.GetQuery() == "" {
			// isNull is true
			os.(*ProtocolDataOutputStream).WriteBoolean(true)
		} else {
			os.(*ProtocolDataOutputStream).WriteBoolean(false)
			err := os.(*ProtocolDataOutputStream).WriteUTF(msg.GetQuery())
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning QueryRequestMessage:WritePayload w/ Error in writing queryStr to message buffer"))
				return err
			}
		}
		if msg.GetEdgeFilter() == "" {
			os.(*ProtocolDataOutputStream).WriteBoolean(true)
		} else {
			os.(*ProtocolDataOutputStream).WriteBoolean(false)
			err := os.(*ProtocolDataOutputStream).WriteUTF(msg.GetEdgeFilter())
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning QueryRequestMessage:WritePayload w/ Error in writing edgeFilter to message buffer"))
				return err
			}
		}
		if msg.GetTraversalCondition() == "" {
			os.(*ProtocolDataOutputStream).WriteBoolean(true)
		} else {
			os.(*ProtocolDataOutputStream).WriteBoolean(false)
			err := os.(*ProtocolDataOutputStream).WriteUTF(msg.GetTraversalCondition())
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning QueryRequestMessage:WritePayload w/ Error in writing traversalFilter to message buffer"))
				return err
			}
		}
		if msg.GetEndCondition() == "" {
			os.(*ProtocolDataOutputStream).WriteBoolean(true)
		} else {
			os.(*ProtocolDataOutputStream).WriteBoolean(false)
			err := os.(*ProtocolDataOutputStream).WriteUTF(msg.GetEndCondition())
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning QueryRequestMessage:WritePayload w/ Error in writing endFilter to message buffer"))
				return err
			}
		}
	} else if msg.GetCommand() == 5 || msg.GetCommand() == 6 {
		// EXECUTED, CLOSE
		os.(*ProtocolDataOutputStream).WriteLong(msg.GetQueryHashId())
	}
	currPos := os.GetPosition()
	length := currPos - startPos
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning QueryRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaler
/////////////////////////////////////////////////////////////////

func (msg *QueryRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.queryExpr, msg.edgeExpr,
		msg.traverseExpr, msg.endExpr, msg.queryHashId, msg.command, msg.queryObject, msg.fetchSize, msg.batchSize,
		msg.traversalDepth, msg.edgeLimit, msg.sortAttrName, msg.sortOrderDsc, msg.sortResultLimit)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning QueryRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaler
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *QueryRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.queryExpr, &msg.edgeExpr, &msg.traverseExpr, &msg.endExpr, &msg.queryHashId, &msg.command, &msg.queryObject,
		&msg.fetchSize, &msg.batchSize, &msg.traversalDepth, &msg.edgeLimit, &msg.sortAttrName, &msg.sortOrderDsc,
		&msg.sortResultLimit)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning QueryRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type QueryResponseMessage struct {
	*AbstractProtocolMessage
	entityStream tgdb.TGInputStream
	hasResult    bool
	totalCount   int
	resultCount  int
	result       int
	queryHashId  int64
	exception 	tgdb.TGError
	//exception TGQueryException
	resultTypeAnnot string
}

func DefaultQueryResponseMessage() *QueryResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(QueryResponseMessage{})

	newMsg := QueryResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.hasResult = false
	newMsg.totalCount = 0
	newMsg.resultCount = 0
	newMsg.result = 0
	newMsg.queryHashId = 0
	newMsg.verbId = VerbQueryResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewQueryResponseMessage(authToken, sessionId int64) *QueryResponseMessage {
	newMsg := DefaultQueryResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for QueryResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *QueryResponseMessage) GetResultCount() int {
	return msg.resultCount
}

func (msg *QueryResponseMessage) GetResult() int {
	return msg.result
}

func (msg *QueryResponseMessage) GetQueryHashId() int64 {
	return msg.queryHashId
}

func (msg *QueryResponseMessage) GetEntityStream() tgdb.TGInputStream {
	return msg.entityStream
}

func (msg *QueryResponseMessage) GetHasResult() bool {
	return msg.hasResult
}

func (msg *QueryResponseMessage) GetTotalCount() int {
	return msg.totalCount
}

func (msg *QueryResponseMessage) GetResultTypeAnnot() string {
	return msg.resultTypeAnnot
}

func (msg *QueryResponseMessage) SetResultCount(size int) {
	msg.resultCount = size
}

func (msg *QueryResponseMessage) SetResult(size int) {
	msg.result = size
}

func (msg *QueryResponseMessage) SetQueryHashId(id int64) {
	msg.queryHashId = id
}

func (msg *QueryResponseMessage) SetEntityStream(eStream tgdb.TGInputStream) {
	msg.entityStream = eStream
}

func (msg *QueryResponseMessage) SetHasResult(resultFlag bool) {
	msg.hasResult = resultFlag
}

func (msg *QueryResponseMessage) SetTotalCount(count int) {
	msg.totalCount = count
}

func (msg *QueryResponseMessage) SetResultTypeAnnot(anon string) {
	msg.resultTypeAnnot = anon
}
/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *QueryResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering QueryResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside QueryResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside QueryResponseMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside QueryResponseMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning QueryResponseMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *QueryResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering QueryResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside QueryResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside QueryResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning QueryResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *QueryResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *QueryResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *QueryResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *QueryResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *QueryResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *QueryResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *QueryResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *QueryResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *QueryResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *QueryResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *QueryResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *QueryResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *QueryResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *QueryResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *QueryResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *QueryResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *QueryResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *QueryResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("QueryResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("HasResult: %+v", msg.hasResult))
	buffer.WriteString(fmt.Sprintf(", TotalCount: %d", msg.totalCount))
	buffer.WriteString(fmt.Sprintf(", ResultCount: %d", msg.resultCount))
	buffer.WriteString(fmt.Sprintf(", Result: %d", msg.result))
	buffer.WriteString(fmt.Sprintf(", QueryHashId: %d", msg.queryHashId))
	//buffer.WriteString(fmt.Sprintf(", EntityStream: %+v", msg.entityStream))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *QueryResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *QueryResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *QueryResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

func (msg *QueryResponseMessage) processQueryStatus (is tgdb.TGInputStream, status int) (tgdb.TGError) {
	if (status == 0) {
		return nil
	}

	decodedStatus := fromStatus(status)
	message, err := is.(*ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload-processQueryStatus w/ Error in reading error message buffer"))
		return err
	}

	err = buildQueryException(decodedStatus, message, status)
	return err
}

func fromStatus(status int) int {
	if status > TGQueryInvalid && status < TGQueryErrorCodeEndMarker {
		return status
	}
	return TGQueryInvalid;
}


// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *QueryResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering QueryResponseMessage:ReadPayload"))
	}
	avail, err := is.(*ProtocolDataInputStream).Available()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading available bytes from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside QueryResponseMessage:ReadPayload read avail as '%+v'", avail))
	}
	if avail == 0 {
		logger.Warning(fmt.Sprint("WARNING: Query response has no data"))
		errMsg := fmt.Sprint("Query response has no data")
		return GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.entityStream = is

	_, err = is.(*ProtocolDataInputStream).ReadInt() // buf length
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading buffer length from message buffer"))
		return err
	}

	_, err = is.(*ProtocolDataInputStream).ReadInt() // checksum
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading checksum from message buffer"))
		return err
	}

	result, err := is.(*ProtocolDataInputStream).ReadInt() // query result
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading query result from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside QueryResponseMessage:ReadPayload read result as '%+v'", result))
	}

	hashId, err := is.(*ProtocolDataInputStream).ReadLong() // query hash id
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading query hashId from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside QueryResponseMessage:ReadPayload read hashId as '%+v'", hashId))
	}

	syntax, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading syntax from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside QueryResponseMessage:ReadPayload read syntax as '%+v'", syntax))
	}


	//TODO: exception needs to be handled here
	err = msg.processQueryStatus(is, result)
	if err != nil {
		msg.exception = err
	}

	annon, err := is.(*ProtocolDataInputStream).ReadUTF() // buf length
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading annotation from message buffer"))
		return err
	}

	msg.SetResultTypeAnnot(annon)

	resultCount, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading resultCount from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside QueryResponseMessage:ReadPayload read resultCount as '%+v'", resultCount))
	}
	if resultCount > 0 {
		msg.SetHasResult(true)
	}

	if syntax == 1 {
		totalCount, err := is.(*ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading totalCount from message buffer"))
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside QueryResponseMessage:ReadPayload read totalCount as '%+v'", totalCount))
		}
		msg.SetTotalCount(totalCount)
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Query has '%d' result entities and %d total entities", resultCount, totalCount))
		}
	} else {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Query has '%d' result count", resultCount))
		}
	}

	msg.SetResult(result)
	msg.SetQueryHashId(hashId)
	msg.SetResultCount(resultCount)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning QueryResponseMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *QueryResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *QueryResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.hasResult, msg.totalCount,
		msg.resultCount, msg.result, msg.queryHashId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning QueryResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *QueryResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.hasResult, &msg.totalCount, &msg.resultCount, &msg.result, &msg.queryHashId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning QueryResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type RollbackTransactionRequestMessage struct {
	*AbstractProtocolMessage
}

func DefaultRollbackTransactionRequestMessage() *RollbackTransactionRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(RollbackTransactionRequestMessage{})

	newMsg := RollbackTransactionRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbRollbackTransactionRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewRollbackTransactionRequestMessage(authToken, sessionId int64) *RollbackTransactionRequestMessage {
	newMsg := DefaultRollbackTransactionRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *RollbackTransactionRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering RollbackTransactionRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning RollbackTransactionRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning RollbackTransactionRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside RollbackTransactionRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside RollbackTransactionRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside RollbackTransactionRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning RollbackTransactionRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *RollbackTransactionRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering RollbackTransactionRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside RollbackTransactionRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside RollbackTransactionRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning RollbackTransactionRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning RollbackTransactionRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *RollbackTransactionRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *RollbackTransactionRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *RollbackTransactionRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *RollbackTransactionRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *RollbackTransactionRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *RollbackTransactionRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *RollbackTransactionRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *RollbackTransactionRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *RollbackTransactionRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *RollbackTransactionRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *RollbackTransactionRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *RollbackTransactionRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *RollbackTransactionRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *RollbackTransactionRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *RollbackTransactionRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *RollbackTransactionRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *RollbackTransactionRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *RollbackTransactionRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("RollbackTransactionRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *RollbackTransactionRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *RollbackTransactionRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *RollbackTransactionRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *RollbackTransactionRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *RollbackTransactionRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *RollbackTransactionRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning RollbackTransactionRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *RollbackTransactionRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning RollbackTransactionRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type RollbackTransactionResponseMessage struct {
	*AbstractProtocolMessage
}

func DefaultRollbackTransactionResponseMessage() *RollbackTransactionResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(RollbackTransactionResponseMessage{})

	newMsg := RollbackTransactionResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbRollbackTransactionResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewRollbackTransactionResponseMessage(authToken, sessionId int64) *RollbackTransactionResponseMessage {
	newMsg := DefaultRollbackTransactionResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *RollbackTransactionResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering RollbackTransactionResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning RollbackTransactionResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning RollbackTransactionResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside RollbackTransactionResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside RollbackTransactionResponseMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside RollbackTransactionResponseMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning RollbackTransactionResponseMessage::FromBytes results in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *RollbackTransactionResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering RollbackTransactionResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside RollbackTransactionResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside RollbackTransactionResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning RollbackTransactionResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning RollbackTransactionResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *RollbackTransactionResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *RollbackTransactionResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *RollbackTransactionResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *RollbackTransactionResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *RollbackTransactionResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *RollbackTransactionResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *RollbackTransactionResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *RollbackTransactionResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *RollbackTransactionResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *RollbackTransactionResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *RollbackTransactionResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *RollbackTransactionResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *RollbackTransactionResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *RollbackTransactionResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *RollbackTransactionResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *RollbackTransactionResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *RollbackTransactionResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *RollbackTransactionResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("RollbackTransactionResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *RollbackTransactionResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *RollbackTransactionResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *RollbackTransactionResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *RollbackTransactionResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *RollbackTransactionResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *RollbackTransactionResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning RollbackTransactionResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *RollbackTransactionResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning RollbackTransactionResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type SessionForcefullyTerminatedMessage struct {
	*ExceptionMessage
}

func DefaultSessionForcefullyTerminatedMessage() *SessionForcefullyTerminatedMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(SessionForcefullyTerminatedMessage{})

	newMsg := SessionForcefullyTerminatedMessage{
		ExceptionMessage: DefaultExceptionMessage(),
	}
	newMsg.verbId = VerbSessionForcefullyTerminated
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewSessionForcefullyTerminatedMessage(authToken, sessionId int64) *SessionForcefullyTerminatedMessage {
	newMsg := DefaultSessionForcefullyTerminatedMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for SessionForcefullyTerminatedMessage
/////////////////////////////////////////////////////////////////

func (msg *SessionForcefullyTerminatedMessage) GetKillString() string {
	return msg.exceptionMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *SessionForcefullyTerminatedMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering SessionForcefullyTerminatedMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning SessionForcefullyTerminatedMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SessionForcefullyTerminatedMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside SessionForcefullyTerminatedMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside SessionForcefullyTerminatedMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside SessionForcefullyTerminatedMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning SessionForcefullyTerminatedMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *SessionForcefullyTerminatedMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering SessionForcefullyTerminatedMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside SessionForcefullyTerminatedMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside SessionForcefullyTerminatedMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SessionForcefullyTerminatedMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning SessionForcefullyTerminatedMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *SessionForcefullyTerminatedMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *SessionForcefullyTerminatedMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *SessionForcefullyTerminatedMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *SessionForcefullyTerminatedMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *SessionForcefullyTerminatedMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *SessionForcefullyTerminatedMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *SessionForcefullyTerminatedMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *SessionForcefullyTerminatedMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *SessionForcefullyTerminatedMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *SessionForcefullyTerminatedMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *SessionForcefullyTerminatedMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *SessionForcefullyTerminatedMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *SessionForcefullyTerminatedMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *SessionForcefullyTerminatedMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *SessionForcefullyTerminatedMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *SessionForcefullyTerminatedMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *SessionForcefullyTerminatedMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *SessionForcefullyTerminatedMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("SessionForcefullyTerminatedMessage:{")
	buffer.WriteString(fmt.Sprintf("ExceptionMsg: %s", msg.exceptionMsg))
	buffer.WriteString(fmt.Sprintf(", ExceptionType: %d", msg.exceptionType))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *SessionForcefullyTerminatedMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *SessionForcefullyTerminatedMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *SessionForcefullyTerminatedMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *SessionForcefullyTerminatedMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering SessionForcefullyTerminatedMessage:ReadPayload"))
	}
	bType, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside SessionForcefullyTerminatedMessage:ReadPayload read bType as '%+v'", bType))
	}

	exMsg, err := is.(*ProtocolDataInputStream).ReadUTF()
	if err != nil {
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside SessionForcefullyTerminatedMessage:ReadPayload read exMsg as '%+v'", exMsg))
	}

	msg.SetExceptionMsg(exMsg)
	msg.SetExceptionType(int(bType))
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning SessionForcefullyTerminatedMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *SessionForcefullyTerminatedMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	os.(*ProtocolDataOutputStream).WriteByte(msg.GetExceptionType())
	return os.(*ProtocolDataOutputStream).WriteUTF(msg.GetExceptionMsg())
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *SessionForcefullyTerminatedMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.exceptionMsg, msg.exceptionType)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning SessionForcefullyTerminatedMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *SessionForcefullyTerminatedMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.exceptionMsg, &msg.exceptionType)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning SessionForcefullyTerminatedMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type TraverseRequestMessage struct {
	*AbstractProtocolMessage
}

func DefaultTraverseRequestMessage() *TraverseRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TraverseRequestMessage{})

	newMsg := TraverseRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbTraverseRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewTraverseRequestMessage(authToken, sessionId int64) *TraverseRequestMessage {
	newMsg := DefaultTraverseRequestMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *TraverseRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TraverseRequestMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning TraverseRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TraverseRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TraverseRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TraverseRequestMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TraverseRequestMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TraverseRequestMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// GetSessionId gets Session id
func (msg *TraverseRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TraverseRequestMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TraverseRequestMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TraverseRequestMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TraverseRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TraverseRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *TraverseRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *TraverseRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *TraverseRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *TraverseRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *TraverseRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *TraverseRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *TraverseRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *TraverseRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *TraverseRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *TraverseRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *TraverseRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *TraverseRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *TraverseRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *TraverseRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *TraverseRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *TraverseRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *TraverseRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *TraverseRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TraverseRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *TraverseRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *TraverseRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *TraverseRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *TraverseRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *TraverseRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *TraverseRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TraverseRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *TraverseRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TraverseRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type TraverseResponseMessage struct {
	*AbstractProtocolMessage
}

func DefaultTraverseResponseMessage() *TraverseResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TraverseResponseMessage{})

	newMsg := TraverseResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbTraverseResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewTraverseResponseMessage(authToken, sessionId int64) *TraverseResponseMessage {
	newMsg := DefaultTraverseResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *TraverseResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TraverseResponseMessage:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning TraverseResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TraverseResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TraverseResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TraverseResponseMessage:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TraverseResponseMessage:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TraverseResponseMessage::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *TraverseResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TraverseResponseMessage:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TraverseResponseMessage:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TraverseResponseMessage:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TraverseResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TraverseResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *TraverseResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *TraverseResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *TraverseResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *TraverseResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *TraverseResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *TraverseResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *TraverseResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *TraverseResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *TraverseResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *TraverseResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *TraverseResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *TraverseResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *TraverseResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *TraverseResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *TraverseResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *TraverseResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *TraverseResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *TraverseResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TraverseResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *TraverseResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *TraverseResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *TraverseResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *TraverseResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *TraverseResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *TraverseResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TraverseResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *TraverseResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TraverseResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
