package utils

import (
	"encoding/binary"
	"math"
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

func BigDecimalToByteArray(f float32) ([]byte, error) {
	var buf [8]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:], nil
}

func LongToCalendar(l int64) time.Time {
	return time.Unix(l, 0)
}

func DoubleToByteArray(f float64) ([]byte, error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:], nil
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
