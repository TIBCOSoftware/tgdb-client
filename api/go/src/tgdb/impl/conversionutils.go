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
 * File Name: conversionutils.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: conversionutils.go 3626 2019-12-09 19:35:03Z nimish $
 */

package impl

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"tgdb"
	"strconv"
	"time"
)

func BigDecimalToByteArray(bd float64) ([]byte, error) {
	strVal := NewTGDecimalFromFloat(bd).String()
	buf := []byte(strVal)
	return buf, nil
}

func ByteArrayToBigDecimal(buf []byte) (*TGDecimal, error) {
	iStream := NewProtocolDataInputStream(buf)
	scale, err := iStream.ReadInt()
	if err != nil {
		errMsg := fmt.Sprint("Unable to iStream.ReadInt()")
		return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	prcLen, err := iStream.ReadInt()
	if err != nil {
		errMsg := fmt.Sprint("Unable to iStream.ReadInt() for prcLen")
		return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	tBuf := make([]byte, prcLen)
	buf1, err := iStream.ReadFully(tBuf)
	if err != nil {
		errMsg := fmt.Sprint("Unable to iStream.ReadFully(buf)")
		return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	bd := big.NewInt(1).SetBytes(buf1)
	dec := NewTGDecimalFromBigInt(bd, int32(scale))
	return &dec, nil
}

func CalendarToString(t time.Time) (string, error) {
	env := NewTGEnvironment()
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

func InputStreamToByteArray(is tgdb.TGInputStream) ([]byte, error) {
	oStream := DefaultProtocolDataOutputStream()
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
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		err = oStream.WriteBytesFromPos(buf, 0, readCnt)
		if err != nil {
			errMsg := fmt.Sprint("Unable to oStream.WriteBytesFromPos(buf, 0, readCnt)")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
	}
	return oStream.GetBuffer(), nil
}

func LongToCalendar(l int64) time.Time {
	return time.Unix(l, 0)
}

func StringToCalendar(s string) (time.Time, error) {
	env := NewTGEnvironment()
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

func ObjectFromByteArray(value []byte, attrType int) (interface{}, tgdb.TGError) {
	if value == nil || len(value) == 0 {
		errMsg := fmt.Sprint("Invalid/Null attribute value received")
		return nil, GetErrorByType(TGErrorIOException, TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	iStream := NewProtocolDataInputStream(value)
	switch attrType {
	case AttributeTypeBoolean:
		return iStream.ReadBoolean()
	case AttributeTypeByte:
		return iStream.ReadByte()
	case AttributeTypeChar:
		return iStream.ReadChar()
	case AttributeTypeShort:
		return iStream.ReadShort()
	case AttributeTypeInteger:
		return iStream.ReadInt()
	case AttributeTypeLong:
		return iStream.ReadLong()
	case AttributeTypeFloat:
		return iStream.ReadFloat()
	case AttributeTypeDouble:
		return iStream.ReadDouble()
	case AttributeTypeNumber:
		bufLen, err := iStream.ReadInt()
		if err != nil {
			errMsg := fmt.Sprint("Unable to read bufLen from input stream")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		buf := make([]byte, bufLen)
		buf1, err := iStream.ReadFully(buf)
		if err != nil {
			errMsg := fmt.Sprint("Unable to iStream.ReadFully(buf)")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		bd, err1 := ByteArrayToBigDecimal(buf1)
		if err1 != nil {
			errMsg := fmt.Sprint("Unable to ByteArrayToBigDecimal(buf1)")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		return bd, nil
	case AttributeTypeString:
		return iStream.ReadUTF()
	case AttributeTypeDate:
		fallthrough
	case AttributeTypeTime:
		fallthrough
	case AttributeTypeTimeStamp:
		strVal, err := iStream.ReadUTF()
		if err != nil {
			errMsg := fmt.Sprint("Unable to iStream.ReadUTF()")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		t, err1 := StringToCalendar(strVal)
		if err1 != nil {
			errMsg := fmt.Sprint("Unable to StringToCalendar(strVal)")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		return t, nil
	case AttributeTypeBlob:
		bufLen, err := iStream.ReadInt()
		if err != nil {
			errMsg := fmt.Sprint("Unable to read bufLen from input stream")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		buf := make([]byte, bufLen)
		_, err = iStream.ReadAtOffset(buf, 0, bufLen)
		if err != nil {
			errMsg := fmt.Sprint("Unable to iStream.ReadAtOffset(buf, 0, bufLen)")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		return buf, nil
	case AttributeTypeClob:
		return iStream.ReadUTF()
	default:
		errMsg := fmt.Sprint("Unable to convert object into byte array")
		return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
}

func ObjectToByteArray(value interface{}, attrType int) ([]byte, tgdb.TGError) {
	if value == nil {
		errMsg := fmt.Sprint("Invalid/Null attribute value received")
		return nil, GetErrorByType(TGErrorIOException, TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	oStream := DefaultProtocolDataOutputStream()
	switch attrType {
	case AttributeTypeBoolean:
		oStream.WriteBoolean(value.(bool))
	case AttributeTypeByte:
		oStream.WriteByte(value.(int))
	case AttributeTypeChar:
		oStream.WriteChar(value.(int))
	case AttributeTypeShort:
		oStream.WriteShort(value.(int))
	case AttributeTypeInteger:
		oStream.WriteInt(value.(int))
	case AttributeTypeLong:
		oStream.WriteLong(value.(int64))
	case AttributeTypeFloat:
		oStream.WriteFloat(value.(float32))
	case AttributeTypeDouble:
		oStream.WriteDouble(value.(float64))
	case AttributeTypeNumber:
		buf, err := BigDecimalToByteArray(value.(float64))
		if err != nil {
			errMsg := fmt.Sprint("Unable to convert object of type Number into byte array")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		oStream.WriteInt(len(buf))
		_ = oStream.WriteBytes(buf)
	case AttributeTypeString:
		_ = oStream.WriteUTF(value.(string))
	case AttributeTypeDate:
		fallthrough
	case AttributeTypeTime:
		fallthrough
	case AttributeTypeTimeStamp:
		strVal, err := CalendarToString(value.(time.Time))
		if err != nil {
			errMsg := fmt.Sprint("Unable to convert object of type TimeStamp into byte array")
			return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
		}
		_ = oStream.WriteUTF(strVal)
	case AttributeTypeBlob:
		buf := value.([]byte)
		oStream.WriteInt(len(buf))
		_ = oStream.WriteBytes(buf)
	case AttributeTypeClob:
		_ = oStream.WriteUTF(value.(string))
	default:
		errMsg := fmt.Sprint("Unable to convert object into byte array")
		return nil, GetErrorByType(TGErrorTypeCoercionNotSupported, TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	return oStream.GetBuffer(), nil
}
