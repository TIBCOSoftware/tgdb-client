package types

import (
	"fmt"
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
 * File name: TGError.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= Various Error Types =======
type TGExceptionType int

const (
	TGSuccess = iota
	TGErrorBadAuthentication
	TGErrorBadMagic
	TGErrorBadVerb
	TGErrorChannelDisconnected
	TGErrorConnectionTimeout
	TGErrorGeneralException
	TGErrorInvalidMessageLength
	TGErrorIOException
	TGErrorProtocolNotSupported
	TGErrorRetryIOException
	TGErrorSecurityException
	TGErrorTransactionException
	TGErrorTypeCoercionNotSupported
	TGErrorTypeNotSupported
	TGErrorVersionMismatchException
	TGErrorInvalidErrorCode
)

type TGError interface {
	error
	// Get the detail error message
	GetErrorCode() string
	GetErrorType() int
	GetErrorMsg() string
	GetErrorDetails() string
}

type TGDBError struct {
	ErrorCode    string
	ErrorType    int
	ErrorMsg     string
	ErrorDetails string
	//ErrorTimestamp  int64
}

// Code message map containing all code messages with the code as key
var PreDefinedErrors = map[int]TGDBError{
	TGSuccess:                       {ErrorCode: "TGSuccess", ErrorType: TGSuccess, ErrorMsg: "", ErrorDetails: ""},
	TGErrorBadAuthentication:        {ErrorCode: "TGErrorBadAuthentication", ErrorType: TGErrorBadAuthentication, ErrorMsg: "", ErrorDetails: ""},
	TGErrorBadMagic:                 {ErrorCode: "TGErrorBadMagic", ErrorType: TGErrorBadMagic, ErrorMsg: "", ErrorDetails: ""},
	TGErrorBadVerb:                  {ErrorCode: "TGErrorBadVerb", ErrorType: TGErrorBadVerb, ErrorMsg: "", ErrorDetails: ""},
	TGErrorChannelDisconnected:      {ErrorCode: "TGErrorChannelDisconnected", ErrorType: TGErrorChannelDisconnected, ErrorMsg: "", ErrorDetails: ""},
	TGErrorConnectionTimeout:        {ErrorCode: "TGErrorConnectionTimeout", ErrorType: TGErrorConnectionTimeout, ErrorMsg: "", ErrorDetails: ""},
	TGErrorGeneralException:         {ErrorCode: "TGErrorGeneralException", ErrorType: TGErrorGeneralException, ErrorMsg: "", ErrorDetails: ""},
	TGErrorInvalidMessageLength:     {ErrorCode: "TGErrorInvalidMessageLength", ErrorType: TGErrorInvalidMessageLength, ErrorMsg: "", ErrorDetails: ""},
	TGErrorIOException:              {ErrorCode: "TGErrorIOException", ErrorType: TGErrorIOException, ErrorMsg: "", ErrorDetails: ""},
	TGErrorProtocolNotSupported:     {ErrorCode: "TGErrorProtocolNotSupported", ErrorType: TGErrorProtocolNotSupported, ErrorMsg: "", ErrorDetails: ""},
	TGErrorRetryIOException:         {ErrorCode: "TGErrorRetryIOException", ErrorType: TGErrorRetryIOException, ErrorMsg: "", ErrorDetails: ""},
	TGErrorSecurityException:        {ErrorCode: "TGErrorSecurityException", ErrorType: TGErrorSecurityException, ErrorMsg: "", ErrorDetails: ""},
	TGErrorTransactionException:     {ErrorCode: "TGErrorTransactionException", ErrorType: TGErrorTransactionException, ErrorMsg: "", ErrorDetails: ""},
	TGErrorTypeCoercionNotSupported: {ErrorCode: "TGErrorTypeCoercionNotSupported", ErrorType: TGErrorTypeCoercionNotSupported, ErrorMsg: "", ErrorDetails: ""},
	TGErrorTypeNotSupported:         {ErrorCode: "TGErrorTypeNotSupported", ErrorType: TGErrorTypeNotSupported, ErrorMsg: "", ErrorDetails: ""},
	TGErrorVersionMismatchException: {ErrorCode: "TGErrorVersionMismatchException", ErrorType: TGErrorVersionMismatchException, ErrorMsg: "", ErrorDetails: ""},
	TGErrorInvalidErrorCode:         {ErrorCode: "TGErrorInvalidErrorCode", ErrorType: TGErrorInvalidErrorCode, ErrorMsg: "", ErrorDetails: ""},
}

func DefaultTGDBError() *TGDBError {
	newTGDBError := TGDBError{ErrorCode: "", ErrorType: TGSuccess, ErrorMsg: "", ErrorDetails: ""}
	return &newTGDBError
}

func NewTGDBError(eCode string, eType int, eMsg, eDetails string) *TGDBError {
	newTGDBError := DefaultTGDBError()
	newTGDBError.ErrorCode = eCode
	newTGDBError.ErrorType = eType
	newTGDBError.ErrorMsg = eMsg
	newTGDBError.ErrorDetails = eDetails
	return newTGDBError
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGDBError
/////////////////////////////////////////////////////////////////

func GetPreDefinedErrors(code string) *TGDBError {
	for _, tgError := range PreDefinedErrors {
		if tgError.ErrorCode == code {
			return &tgError
		}
	}
	invalid := PreDefinedErrors[TGErrorInvalidErrorCode]
	return &invalid
}

/////////////////////////////////////////////////////////////////
// Helper functions for Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGDBError) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGDBError) GetErrorType() int {
	return e.ErrorType
}

func (e *TGDBError) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGDBError) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGDBError) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}
