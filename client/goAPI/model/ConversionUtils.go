package model

import (
	"encoding/binary"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"math"
	"math/big"
	"strconv"
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
 * File name: ConversionUtils.go
 * Created on: Dec 14, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func BigDecimalToByteArray(bd float64) ([]byte, error) {
	strVal := utils.NewTGDecimalFromFloat(bd).String()
	buf := []byte(strVal)
	return buf, nil
}

func ByteArrayToBigDecimal(buf []byte) (*utils.TGDecimal, error) {
	iStream := iostream.NewProtocolDataInputStream(buf)
	scale, err := iStream.ReadInt()
	if err != nil {
		errMsg := fmt.Sprint("Unable to iStream.ReadInt()")
		return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	prcLen, err := iStream.ReadInt()
	if err != nil {
		errMsg := fmt.Sprint("Unable to iStream.ReadInt() for prcLen")
		return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	tBuf := make([]byte, prcLen)
	buf1, err := iStream.ReadFully(tBuf)
	if err != nil {
		errMsg := fmt.Sprint("Unable to iStream.ReadFully(buf)")
		return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	bd := big.NewInt(1).SetBytes(buf1)
	dec := utils.NewTGDecimalFromBigInt(bd, int32(scale))
	return &dec, nil
}

func CalendarToString(t time.Time) (string, error) {
	env := utils.NewTGEnvironment()
	return t.Format(env.GetDefaultDateTimeFormat()), nil
}

func DoubleToByteArray(f float64) ([]byte, error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:], nil
}

func FloatToByteArray(f float32) ([]byte, error) {
	var buf [8]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:], nil
}

func InputStreamToByteArray(is types.TGInputStream) ([]byte, error) {
	oStream := iostream.DefaultProtocolDataOutputStream()
	avail := 0
	for {
		avail, _ = is.Available()
		if avail <= 0 {
			break
		}
		buf := make([]byte, avail)
		readCnt, err := is.ReadAtOffset(buf, 0, avail)
		if err != nil {
			errMsg := fmt.Sprint("Unable to is.ReadAtOffset(buf, 0 , avail")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		err = oStream.WriteBytesFromPos(buf, 0, readCnt)
		if err != nil {
			errMsg := fmt.Sprint("Unable to oStream.WriteBytesFromPos(buf, 0, readCnt)")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
	}
	return oStream.GetBuffer(), nil
}

func LongToCalendar(l int64) time.Time {
	return time.Unix(l, 0)
}

func StringToCalendar(s string) (time.Time, error) {
	env := utils.NewTGEnvironment()
	return time.ParseInLocation(env.GetDefaultDateTimeFormat(), s, time.Local)
}

func StringToDouble(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func StringToFloat(s string) (float32, error) {
	f, err := StringToDouble(s)
	return float32(f), err
}

func StringToInteger(s string) (int, error) {
	return strconv.Atoi(s)
}

func StringToLong(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func StringToShort(s string) (int16, error) {
	v, err := strconv.ParseInt(s, 10, 16)
	return int16(v), err
}

func StringToCharacter(s string) (string, error) {
	if len(s) == 1 {
		return string(s[0]), nil
	}
	v, err := strconv.ParseInt(s, 10, 16)
	return string(v), err
}

func ObjectFromByteArray(value []byte, attrType int) (interface{}, types.TGError) {
	if value == nil {
		errMsg := fmt.Sprint("Invalid/Null attribute value received")
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	iStream := iostream.NewProtocolDataInputStream(value)
	switch attrType {
	case types.AttributeTypeBoolean:
		return iStream.ReadBoolean()
	case types.AttributeTypeByte:
		return iStream.ReadByte()
	case types.AttributeTypeChar:
		return iStream.ReadChar()
	case types.AttributeTypeShort:
		return iStream.ReadShort()
	case types.AttributeTypeInteger:
		return iStream.ReadInt()
	case types.AttributeTypeLong:
		return iStream.ReadLong()
	case types.AttributeTypeFloat:
		return iStream.ReadFloat()
	case types.AttributeTypeDouble:
		return iStream.ReadDouble()
	case types.AttributeTypeNumber:
		bufLen, err := iStream.ReadInt()
		if err != nil {
			errMsg := fmt.Sprint("Unable to read bufLen from input stream")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		buf := make([]byte, bufLen)
		buf1, err := iStream.ReadFully(buf)
		if err != nil {
			errMsg := fmt.Sprint("Unable to iStream.ReadFully(buf)")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		bd, err1 := ByteArrayToBigDecimal(buf1)
		if err1 != nil {
			errMsg := fmt.Sprint("Unable to ByteArrayToBigDecimal(buf1)")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		return bd, nil
	case types.AttributeTypeString:
		return iStream.ReadUTF()
	case types.AttributeTypeDate:
		fallthrough
	case types.AttributeTypeTime:
		fallthrough
	case types.AttributeTypeTimeStamp:
		strVal, err := iStream.ReadUTF()
		if err != nil {
			errMsg := fmt.Sprint("Unable to iStream.ReadUTF()")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		t, err1 := StringToCalendar(strVal)
		if err1 != nil {
			errMsg := fmt.Sprint("Unable to StringToCalendar(strVal)")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		return t, nil
	case types.AttributeTypeBlob:
		bufLen, err := iStream.ReadInt()
		if err != nil {
			errMsg := fmt.Sprint("Unable to read bufLen from input stream")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		buf := make([]byte, bufLen)
		_, err = iStream.ReadAtOffset(buf, 0, bufLen)
		if err != nil {
			errMsg := fmt.Sprint("Unable to iStream.ReadAtOffset(buf, 0, bufLen)")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		return buf, nil
	case types.AttributeTypeClob:
		return iStream.ReadUTF()
	default:
		errMsg := fmt.Sprint("Unable to convert object into byte array")
		return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
}

func ObjectToByteArray(value interface{}, attrType int) ([]byte, types.TGError) {
	if value == nil {
		errMsg := fmt.Sprint("Invalid/Null attribute value received")
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	oStream := iostream.DefaultProtocolDataOutputStream()
	switch attrType {
	case types.AttributeTypeBoolean:
		oStream.WriteBoolean(value.(bool))
	case types.AttributeTypeByte:
		oStream.WriteByte(value.(int))
	case types.AttributeTypeChar:
		oStream.WriteChar(value.(int))
	case types.AttributeTypeShort:
		oStream.WriteShort(value.(int))
	case types.AttributeTypeInteger:
		oStream.WriteInt(value.(int))
	case types.AttributeTypeLong:
		oStream.WriteLong(value.(int64))
	case types.AttributeTypeFloat:
		oStream.WriteFloat(value.(float32))
	case types.AttributeTypeDouble:
		oStream.WriteDouble(value.(float64))
	case types.AttributeTypeNumber:
		buf, err := BigDecimalToByteArray(value.(float64))
		if err != nil {
			errMsg := fmt.Sprint("Unable to convert object of type Number into byte array")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		oStream.WriteInt(len(buf))
		_ = oStream.WriteBytes(buf)
	case types.AttributeTypeString:
		_ = oStream.WriteUTF(value.(string))
	case types.AttributeTypeDate:
		fallthrough
	case types.AttributeTypeTime:
		fallthrough
	case types.AttributeTypeTimeStamp:
		strVal, err := CalendarToString(value.(time.Time))
		if err != nil {
			errMsg := fmt.Sprint("Unable to convert object of type TimeStamp into byte array")
			return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		_ = oStream.WriteUTF(strVal)
	case types.AttributeTypeBlob:
		buf := value.([]byte)
		oStream.WriteInt(len(buf))
		_ = oStream.WriteBytes(buf)
	case types.AttributeTypeClob:
		_ = oStream.WriteUTF(value.(string))
	default:
		errMsg := fmt.Sprint("Unable to convert object into byte array")
		return nil, exception.GetErrorByType(types.TGErrorTypeCoercionNotSupported, types.TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	return oStream.GetBuffer(), nil
}
