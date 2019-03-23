package iostream

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/logging"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"math"
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
 * File name: ProtocolDataInputStream.go
 * Created on: Nov 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/** Test to ensure float/double gets converted to/from int/int64 correctly
package main

import (
    "fmt"
    "math"
)

func main() {
	fmt.Println("Hello, playground")
	input := 1.23456789
	tmp := math.Float64bits(input)	// For float, use Float32bits()
	output := int64(tmp)
	fmt.Printf("Input: '%+v' Output: '%+v'\n", input, output)
	input1 := tmp
	output1 := math.Float64frombits(input1)	// For float, use Float32frombits()
	fmt.Printf("Input1: '%+v' Output1: '%+v'\n", input1, output1)
}
*/

//var mark = 0
var logger = logging.DefaultTGLogManager().GetLogger()

type ProtocolDataInputStream struct {
	mark                int	// Private Internal Marker - Not to be serialized / de-serialized
	Buf                 []byte
	BufLen              int
	iStreamCurPos       int
	iStreamEncoding     string
	iStreamReferenceMap map[int64]types.TGEntity
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
		iStreamReferenceMap: make(map[int64]types.TGEntity, 0),
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

/////////////////////////////////////////////////////////////////
// Helper functions for ProtocolDataInputStream
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataInputStream) ReadBoolean() (bool, types.TGError) {
	//logger.Log(fmt.Sprintf(Entering ProtocolDataInputStream::ReadBoolean()"))
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadBoolean('%+v') as msg.CurPos >= msg.BufLen", false))
		errMsg := fmt.Sprint("End of data stream")
		return false, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := msg.Buf[msg.iStreamCurPos]
	msg.iStreamCurPos += 1
	if a == 0x00 {
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadBoolean('%+v') as a == 0x00", false))
		return false, nil
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadBoolean from the contents w/ ('%+v')", true))
	return true, nil
}

func (msg *ProtocolDataInputStream) ReadByte() (byte, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadByte()"))
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadByte('%+v') as msg.CurPos >= msg.BufLen", 0))
		errMsg := fmt.Sprint("End of data stream")
		return 0, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := msg.Buf[msg.iStreamCurPos]
	msg.iStreamCurPos += 1
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadByte from the contents w/ ('%+v')", a))
	return a, nil
}

func (msg *ProtocolDataInputStream) ReadChar() (string, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadChar()"))
	if (msg.iStreamCurPos + 2) > msg.BufLen {
		logger.Error(fmt.Sprintf("Returning ProtocolDataInputStream::ReadChar('%+v') as msg.CurPos+2 > msg.BufLen", ""))
		errMsg := fmt.Sprint("End of data stream")
		return "", exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int(msg.Buf[msg.iStreamCurPos]) << 8
	b := int(msg.Buf[msg.iStreamCurPos+1]) & 0xff
	msg.iStreamCurPos += 2
	c := fmt.Sprintf("%c", a+b)
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadChar from the contents w/ ('%+v')", c))
	return c, nil
}

func (msg *ProtocolDataInputStream) ReadDouble() (float64, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadDouble()"))
	a, err := msg.ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadDouble('%+v') as msg.ReadInt resulted in error", 0))
		return 0, err
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadDouble from the contents w/ ('%+v')", float64(a)))
	return math.Float64frombits(uint64(a)), nil
}

func (msg *ProtocolDataInputStream) ReadFloat() (float32, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadFloat()"))
	a, err := msg.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadFloat('%+v') as msg.ReadInt resulted in error", 0))
		return 0, err
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadFloat from the contents w/ ('%+v')", float32(a)))
	return math.Float32frombits(uint32(a)), nil
}

func (msg *ProtocolDataInputStream) ReadFully(b []byte) ([]byte, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadFully() from contents '%+v'", msg.Buf))
	if b == nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadFully as input byte array is EMPTY"))
		errMsg := fmt.Sprint("Input argument to ReadFully is null")
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return msg.ReadFullyAtPos(b, msg.iStreamCurPos, len(b))
}

func (msg *ProtocolDataInputStream) ReadFullyAtPos(b []byte, readCurPos, readLen int) ([]byte, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadFullyAtPos('%d') for '%d' from '%+v'", readCurPos, readLen, msg.Buf[msg.CurPos:(msg.CurPos+readLen)]))
	if readLen <= 0 {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadFullyAtPos('%+v') as readLen <= 0", readCurPos))
		errMsg := fmt.Sprint("Input argument representing how many bytes to read in ReadFullyAtPos is null")
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if (msg.iStreamCurPos + readLen) > msg.BufLen {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadFullyAtPos as (msg.CurPos + readLen) > msg.BufLen"))
		errMsg := fmt.Sprint("Input readLen to ReadFully is invalid")
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// Copy in input buffer b at offset for length
	copy(b, msg.Buf[msg.iStreamCurPos:(msg.iStreamCurPos +readLen)])
	//logger.Log(fmt.Sprintf("Inside ProtocolDataInputStream::ReadFullyAtPos() temp / b is '%+v'", b))
	msg.iStreamCurPos = msg.iStreamCurPos + readLen
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadFullyAtPos('%d') for length '%d' from '%+v' in contents", readCurPos, readLen, b))
	return b, nil
}

func (msg *ProtocolDataInputStream) ReadInt() (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadInt() from Stream '%+v'", msg.String()))
	if (msg.iStreamCurPos + 4) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadInt as (msg.CurPos + 4) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int(msg.Buf[msg.iStreamCurPos]) << 24
	b := int(msg.Buf[msg.iStreamCurPos+1]) << 16 & 0x00ff0000
	c := int(msg.Buf[msg.iStreamCurPos+2]) << 8 & 0x0000ff00
	d := int(msg.Buf[msg.iStreamCurPos+3]) & 0xff
	msg.iStreamCurPos += 4
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadInt from the contents w/ ('%d')", a+b+c+d))
	return a + b + c + d, nil
}

func (msg *ProtocolDataInputStream) ReadLine() (string, types.TGError) {
	// TODO: Revisit later - No-op For Now
	return "", nil
}

func (msg *ProtocolDataInputStream) ReadLong() (int64, types.TGError) {
	//logger.Log(fmt.Sprint("Entering ProtocolDataInputStream::ReadLong()"))
	if (msg.iStreamCurPos + 8) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadLong as (msg.CurPos + 8) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
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

func (msg *ProtocolDataInputStream) ReadShort() (int16, types.TGError) {
	//logger.Log(fmt.Sprint(Entering ProtocolDataInputStream::ReadShort()"))
	if (msg.iStreamCurPos + 2) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadShort as (msg.CurPos + 2) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int16(msg.Buf[msg.iStreamCurPos]) << 8
	b := int16(msg.Buf[msg.iStreamCurPos+1]) & 0xff
	msg.iStreamCurPos += 2
	c := a + b
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadShort from the contents w/ ('%d')", c))
	return c, nil
}

func (msg *ProtocolDataInputStream) ReadUnsignedByte() (int, types.TGError) {
	//logger.Log(fmt.Sprint("Entering ProtocolDataInputStream::ReadUnsignedByte()"))
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUnsignedByte as msg.CurPos > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int(msg.Buf[msg.iStreamCurPos]) & 0xff
	msg.iStreamCurPos += 1
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadUnsignedByte from the contents w/ ('%d')", a))
	return a, nil
}

func (msg *ProtocolDataInputStream) ReadUnsignedShort() (uint16, types.TGError) {
	//logger.Log(fmt.Sprint("Entering ProtocolDataInputStream::ReadUnsignedShort()"))
	if (msg.iStreamCurPos + 2) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUnsignedShort as (msg.CurPos + 2) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return 0, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	a := int(msg.Buf[msg.iStreamCurPos]) << 8
	b := int(msg.Buf[msg.iStreamCurPos+1]) & 0xff
	msg.iStreamCurPos += 2
	c := uint16((a + b) & 0x0000ffff)
	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadUnsignedShort from the contents w/ ('%d')", c))
	return c, nil
}

func (msg *ProtocolDataInputStream) ReadUTF() (string, types.TGError) {
	//logger.Log(fmt.Sprint("Entering ProtocolDataInputStream::ReadUTF()"))
	start := msg.iStreamCurPos

	utfLen, err := msg.ReadUnsignedShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTF as msg.ReadUnsignedShort resulted in error"))
		// Reset the current position
		msg.iStreamCurPos = start
		return "", err
	}
	logger.Log(fmt.Sprintf("Inside ProtocolDataInputStream::ReadUTF() utfLength to read: '%d'", utfLen))
	str, err := msg.ReadUTFString(int(utfLen))
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadUTF from the contents w/ ('%s') & error ('%+v')", str, err))
	return str, err
}

func (msg *ProtocolDataInputStream) ReadUTFString(utfLen int) (string, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadUTFString('%d') CurPos: %d BufLen: %d Buf: %+v", utfLen, msg.CurPos, msg.BufLen, msg.Buf))
	if (msg.iStreamCurPos + utfLen) > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as (msg.CurPos + utfLen) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return "", exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
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
		//logger.Log(fmt.Sprintf("ProtocolDataInputStream::ReadUTFString - Inside Loop msg.CurPos '%d' lastPos '%d' c '%d' cVal '%d' char2 '%d' char3 '%d'", msg.CurPos, lastPos, c, cVal, char2, char3))
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
				return "", exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
			}
			char2 = int(msg.Buf[msg.iStreamCurPos])
			msg.iStreamCurPos += 1
			if (char2 & 0xC0) != 0x80 {
				msg.iStreamCurPos = start
				logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as (char2 & 0xC0) != 0x80"))
				errMsg := fmt.Sprint("Data Format Issue")
				return "", exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
			}
			c1 := ((c & 0x1F) << 6) | (char2 & 0x3F)
			charBuf = charBuf + string(c1)
			strLen++
		} else if cVal == 14 { // 1110 xxxx  10xx xxxx  10xx xxxx
			if msg.iStreamCurPos+2 > lastPos {
				msg.iStreamCurPos = start
				logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as msg.CurPos+2 > lastPos"))
				errMsg := fmt.Sprint("Data Format Issue")
				return "", exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
			}
			char2 = int(msg.Buf[msg.iStreamCurPos])
			msg.iStreamCurPos += 1
			char3 = int(msg.Buf[msg.iStreamCurPos])
			msg.iStreamCurPos += 1
			if ((char2 & 0xC0) != 0x80) || ((char3 & 0xC0) != 0x80) {
				msg.iStreamCurPos = start
				logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as ((char2 & 0xC0) != 0x80) || ((char3 & 0xC0) != 0x80)"))
				errMsg := fmt.Sprint("Data Format Issue")
				return "", exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
			}
			c1 := ((c & 0x0F) << 12) | ((char2 & 0x3F) << 6) | ((char3 & 0x3F) << 0)
			charBuf = charBuf + string(c1)
			strLen++
		} else { // 10xx xxxx,  1111 xxxx
			msg.iStreamCurPos = start
			logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString for 10xxxxxx, 1111xxxx"))
			errMsg := fmt.Sprint("Data Format Issue")
			return "", exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		}
	}

	if msg.iStreamCurPos != lastPos {
		msg.iStreamCurPos = start
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadUTFString as msg.CurPos != lastPos"))
		errMsg := fmt.Sprint("Data Format Issue")
		return "", exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	//logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadUTFString from the contents w/ ('%+v')", charBuf))
	return charBuf, nil
}

func (msg *ProtocolDataInputStream) SkipBytes(n int) (int, types.TGError) {
	logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::SkipBytes('%d'), contents are '%+v'", n, msg.Buf))
	if n < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::SkipBytes as n < 0"))
		errMsg := fmt.Sprint("Invalid bytes to skip")
		return 0, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	avail := msg.BufLen - msg.iStreamCurPos
	logger.Log(fmt.Sprintf("ProtocolDataInputStream::SkipBytes - Available bytes ('%d') compared w/ To-Be-Skipped '%d'", avail, n))
	if avail <= n {
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::SkipBytes as avail <= n from the contents w/ ('%d')", avail))
		msg.iStreamCurPos = msg.BufLen
		return avail, nil
	}
	msg.iStreamCurPos = msg.iStreamCurPos + n
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::SkipBytes from the contents w/ ('%d')", n))
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
func (msg *ProtocolDataInputStream) Available() (int, types.TGError) {
	//if msg.BufLen == 0 || (msg.BufLen-msg.CurPos) < 0 {
	if (msg.BufLen-msg.iStreamCurPos) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::Available as (msg.BufLen-msg.CurPos) < 0"))
		errMsg := fmt.Sprint("Invalid data length of Protocol Data Input Stream")
		return 0, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	return msg.BufLen - msg.iStreamCurPos, nil
}

// GetPosition gets the current position of internal cursor
func (msg *ProtocolDataInputStream) GetPosition() int64 {
	return int64(msg.iStreamCurPos)
}

// GetReferenceMap returns a user maintained reference map
func (msg *ProtocolDataInputStream) GetReferenceMap() map[int64]types.TGEntity {
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
func (msg *ProtocolDataInputStream) Read() (int, types.TGError) {
	logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::Read(), contents are '%+v'", msg.Buf))
	if msg.iStreamCurPos > msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::Read as msg.CurPos > msg.BufLen"))
		return -1, nil
	}
	if msg.iStreamCurPos < msg.BufLen {
		b := int(msg.Buf[msg.iStreamCurPos] & 0xff)
		msg.iStreamCurPos += 1
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::Read from the contents w/ ('%d')", b))
		return b, nil
	}
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::Read from the contents w/ ('%d')", 0))
	return 0, nil
}

// ReadIntoBuffer copies bytes in specified buffer
// The buffer cannot be NIL
func (msg *ProtocolDataInputStream) ReadIntoBuffer(b []byte) (int, types.TGError) {
	return msg.ReadAtOffset(b, 0, len(b))
}

// ReadAtOffset is similar to readFully.
func (msg *ProtocolDataInputStream) ReadAtOffset(b []byte, off int, length int) (int, types.TGError) {
	logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadAtOffset('%d') for length '%d' from '%+v', contents are '%+v'", off, length, b, msg.Buf))
	if b == nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadAtOffset as b == nil"))
		return -1, exception.CreateExceptionByType(types.TGErrorIOException)
	}
	if off < 0 || off > len(b) || length < 0 || (off+length) > len(b) || (off+length) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadAtOffset as invalid input data / offset / length OR data is corrupt"))
		errMsg := fmt.Sprint("Invalid input data / offset / length OR data is corrupt")
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadAtOffset('%d') for length '%d' from the contents w/ ('%d')", off, -1, b))
		return -1, nil
	}
	if (msg.iStreamCurPos + length) > msg.BufLen {
		length = msg.BufLen - msg.iStreamCurPos
	}
	if length <= 0 {
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadAtOffset('%d') for length '%d' from the contents w/ ('%d')", off, 0, b))
		return 0, nil
	}
	// Copy from input buffer b at offset for length
	// Create a slice of input buffer
	temp := b[off:] // Slice starting at offset 'off' of input buffer 'b'
	tempLen := len(temp) - length
	temp = temp[:tempLen]
	copy(temp, msg.Buf)
	msg.iStreamCurPos = msg.iStreamCurPos + length
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadAtOffset('%d') for length '%d' from the contents w/ ('%d')", off, length, b))
	return length, nil
}

// ReadBytes reads an encoded byte array. writeBytes encodes the length, and the byte[].
// This is equivalent to do a readInt, and read(byte[])
func (msg *ProtocolDataInputStream) ReadBytes() ([]byte, types.TGError) {
	logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadBytes(), contents are '%+v'", msg.Buf))
	emptyByteArray := make([]byte, 0)
	length, err := msg.ReadInt()
	if length == 0 {
		logger.Log(fmt.Sprint("Returning ProtocolDataInputStream::ReadBytes as len == 0"))
		return emptyByteArray, nil
	}
	if length == -1 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadBytes as len == -1"))
		errMsg := fmt.Sprint("Read data corrupt")
		return emptyByteArray, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	buf := make([]byte, length)
	_, err = msg.ReadIntoBuffer(buf)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadBytes w/ error in reading into buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadBytes from the contents w/ ('%+v')", buf))
	return buf, nil
}

// ReadVarLong reads a Variable long field
func (msg *ProtocolDataInputStream) ReadVarLong() (int64, types.TGError) {
	logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::ReadVarLong(), contents are '%+v'", msg.Buf))
	if msg.iStreamCurPos >= msg.BufLen {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataInputStream::ReadVarLong as msg.CurPos >= msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	curByte := msg.Buf[msg.iStreamCurPos]
	if curByte == types.U64PACKED_NULL {
		msg.iStreamCurPos++
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", types.U64_NULL))
		return types.U64_NULL, nil
	}

	if curByte&0x80 == 0 {
		msg.iStreamCurPos++
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", int64(curByte)))
		return int64(curByte), nil
	}

	if curByte&0x40 == 0 {
		if (msg.iStreamCurPos + 2) > msg.BufLen {
			logger.Log(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadVarLong as (msg.CurPos + 2) > msg.BufLen"))
			errMsg := fmt.Sprint("End of data stream")
			return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		}
		a := msg.Buf[msg.iStreamCurPos] << 8
		msg.iStreamCurPos++
		b := msg.Buf[msg.iStreamCurPos] & 0xff
		msg.iStreamCurPos++
		c := int16(a + b)
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", int64(c&0x3fff)))
		return int64(c & 0x3fff), nil
	}

	if curByte&0x20 == 0 {
		if (msg.iStreamCurPos + 4) > msg.BufLen {
			logger.Log(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadVarLong as (msg.CurPos + 4) > msg.BufLen"))
			errMsg := fmt.Sprint("End of data stream")
			return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
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
		logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", int64(lValue)))
		return int64(lValue), nil
	}

	count := int(curByte & 0x0f)
	msg.iStreamCurPos++
	if (msg.iStreamCurPos + count) > msg.BufLen {
		logger.Log(fmt.Sprintf("ERROR: Returning ProtocolDataInputStream::ReadVarLong as (msg.CurPos + count) > msg.BufLen"))
		errMsg := fmt.Sprint("End of data stream")
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	lValue := 0
	for i := 0; i < count; i++ {
		lValue <<= 8
		lValue |= int(msg.Buf[msg.iStreamCurPos] & 0x00ff)
	}
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::ReadVarLong from the contents w/ ('%d')", int64(lValue)))
	return int64(lValue), nil
}

// Reset brings internal moving cursor back to the old position
func (msg *ProtocolDataInputStream) Reset() {
	logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::Reset( Marker = '%d') in the contents", msg.mark))
	msg.iStreamCurPos = msg.mark
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::Reset( Marker = '%d'), contents are '%+v'", msg.mark, msg.Buf))
}

// SetPosition sets the position of reading.
func (msg *ProtocolDataInputStream) SetPosition(position int64) int64 {
	logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::SetPosition('%d') in the contents", position))
	oldPos := msg.iStreamCurPos
	msg.iStreamCurPos = int(position)
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::SetPosition('%d'), contents are '%+v'", int64(oldPos), msg.Buf))
	return int64(oldPos)
}

// SetReferenceMap sets a user maintained map as reference data
func (msg *ProtocolDataInputStream) SetReferenceMap(rMap map[int64]types.TGEntity) {
	msg.iStreamReferenceMap = rMap
}

// Skip skips n bytes
func (msg *ProtocolDataInputStream) Skip(n int64) (int64, types.TGError) {
	logger.Log(fmt.Sprintf("Entering ProtocolDataInputStream::Skip('%d') in the contents", n))
	a, err := msg.SkipBytes(int(n))
	logger.Log(fmt.Sprintf("Returning ProtocolDataInputStream::Skip('%d'), contents are '%+v'", int64(a), msg.Buf))
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
