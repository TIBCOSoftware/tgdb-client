package types

import (
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
 * File name: TGQuery.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGQuery interface {
	// Close closes the Query
	Close()
	// Execute executes the Query
	Execute() TGResultSet
	// SetBoolean sets Boolean parameter
	SetBoolean(name string, value bool)
	// SetBytes sets Byte Parameter
	SetBytes(name string, bos []byte)
	// SetChar sets Character Parameter
	SetChar(name string, value string)
	// SetDate sets Date Parameter
	SetDate(name string, value time.Time)
	// SetDouble sets Double Parameter
	SetDouble(name string, value float64)
	// SetFloat sets Float Parameter
	SetFloat(name string, value float32)
	// SetInt sets Integer Parameter
	SetInt(name string, value int)
	// SetLong sets Long Parameter
	SetLong(name string, value int64)
	// SetNull sets the parameter to null
	SetNull(name string)
	// SetOption sets the Query Option
	SetOption(options TGQueryOption)
	// SetShort sets Short Parameter
	SetShort(name string, value int16)
	// SetString sets String Parameter
	SetString(name string, value string)
	// Additional Method to help debugging
	String() string
}
