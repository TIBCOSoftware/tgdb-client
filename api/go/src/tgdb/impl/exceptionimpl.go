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
 * File Name: exceptionimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: exceptionimpl.go 4048 2020-06-01 18:53:50Z nimish $
 */

package impl

import (
	"fmt"
	"tgdb"
)

/////////////////////////////////////////////////////////////////
// Helper functions for Interface ==> TGError
/////////////////////////////////////////////////////////////////
const (
	TGErrorBadVerb = iota
	TGErrorInvalidMessageLength
	TGErrorBadMagic
	TGErrorProtocolNotSupported
	TGErrorBadAuthentication
	TGErrorIOException
	TGErrorConnectionTimeout
	TGErrorGeneralException
	TGErrorRetryIOException
	TGErrorChannelDisconnected
	TGErrorSecurityException
	TGErrorTransactionException
	TGErrorTypeCoercionNotSupported
	TGErrorTypeNotSupported
	TGErrorVersionMismatchException
	TGErrorInvalidErrorCode
	TGSuccess
	TGQryError
    TGQryProviderNotInitialized
    TGQryParsingError
    TGQryStepNotSupported
    TGQryStepNotAllowed
    TGQryStepArgMissing
    TGQryStepArgNotSupported
    TGQryStepMissing
    TGQryNotDefined
    TGQryAttrDescNotFound
    TGQryEdgeTypeNotFound
    TGQryNodeTypeNotFound
    TGQryInternalDataMismatchError
    TGQryStepSignatureNotSupported
    TGQryInvalidDataType
)

const (
	TGQueryInvalid = iota + 8100
	TGQueryProviderNotInitialized
	TGQueryParsingError
	TGQueryStepNotSupported
	TGQueryStepNotAllowed
	TGQueryStepArgMissing
	TGQueryStepArgNotSupported
	TGQueryStepMissing
	TGQueryNotDefined
	TGQueryAttrDescNotFound
	TGQueryEdgeTypeNotFound
	TGQueryNodeTypeNotFound
	TGQueryInternalDataMismatchError
	TGQueryStepSignatureNotSupported
	TGQueryInvalidDataType
	TGQueryErrorCodeEndMarker
)

type TGDBError struct {
	ErrorCode    string
	ErrorType    int
	ErrorMsg     string
	ErrorDetails string
	ErrorServerErrorCode int
	//ErrorTimestamp  int64
}


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

func (e *TGDBError) GetServerErrorCode() int {
	return e.ErrorServerErrorCode
}


// Code message map containing all code messages with the code as key
var PreDefinedErrors = map[int]TGDBError{
	TGErrorBadVerb:                  {ErrorCode: "TGErrorBadVerb", ErrorType: TGErrorBadVerb, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorInvalidMessageLength:     {ErrorCode: "TGErrorInvalidMessageLength", ErrorType: TGErrorInvalidMessageLength, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorBadMagic:                 {ErrorCode: "TGErrorBadMagic", ErrorType: TGErrorBadMagic, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorProtocolNotSupported:     {ErrorCode: "TGErrorProtocolNotSupported", ErrorType: TGErrorProtocolNotSupported, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorBadAuthentication:        {ErrorCode: "TGErrorBadAuthentication", ErrorType: TGErrorBadAuthentication, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorIOException:              {ErrorCode: "TGErrorIOException", ErrorType: TGErrorIOException, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorConnectionTimeout:        {ErrorCode: "TGErrorConnectionTimeout", ErrorType: TGErrorConnectionTimeout, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorGeneralException:         {ErrorCode: "TGErrorGeneralException", ErrorType: TGErrorGeneralException, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorRetryIOException:         {ErrorCode: "TGErrorRetryIOException", ErrorType: TGErrorRetryIOException, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorChannelDisconnected:      {ErrorCode: "TGErrorChannelDisconnected", ErrorType: TGErrorChannelDisconnected, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorSecurityException:        {ErrorCode: "TGErrorSecurityException", ErrorType: TGErrorSecurityException, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorTransactionException:     {ErrorCode: "TGErrorTransactionException", ErrorType: TGErrorTransactionException, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorTypeCoercionNotSupported: {ErrorCode: "TGErrorTypeCoercionNotSupported", ErrorType: TGErrorTypeCoercionNotSupported, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorTypeNotSupported:         {ErrorCode: "TGErrorTypeNotSupported", ErrorType: TGErrorTypeNotSupported, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorVersionMismatchException: {ErrorCode: "TGErrorVersionMismatchException", ErrorType: TGErrorVersionMismatchException, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGErrorInvalidErrorCode:         {ErrorCode: "TGErrorInvalidErrorCode", ErrorType: TGErrorInvalidErrorCode, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
	TGSuccess:                       {ErrorCode: "TGSuccess", ErrorType: TGSuccess, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1},
}

func DefaultTGDBError() *TGDBError {
	newTGDBError := TGDBError{ErrorCode: "", ErrorType: TGSuccess, ErrorMsg: "", ErrorDetails: "", ErrorServerErrorCode:-1}
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

func NewTGDBErrorWithServerErrorCode(eCode string, eType int, eMsg, eDetails string, eServerErrorCode int) *TGDBError {
	newTGDBError := NewTGDBError(eCode, eType, eMsg, eDetails)
	newTGDBError.ErrorServerErrorCode = eServerErrorCode
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
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGDBError) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


// Create new exception instance based on the input type
func CreateExceptionByType(excpTypeId int) tgdb.TGError {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputExcpTypeId := excpTypeId

	// Use a switch case to switch between exception types, if a type exist then error is nil (null)
	// Whenever new exception type gets into the mix, just add a case below
	switch inputExcpTypeId {
	case TGSuccess:
		return DefaultTGSuccess()

	case TGErrorBadAuthentication:
		return DefaultTGBadAuthentication()
	case TGErrorBadMagic:
		return DefaultTGBadMagic()
	case TGErrorBadVerb:
		return DefaultTGBadVerb()
	case TGErrorChannelDisconnected:
		return DefaultTGChannelDisconnected()
	case TGErrorConnectionTimeout:
		return DefaultTGConnectionTimeout()
	case TGErrorGeneralException:
		return DefaultTGGeneralException()
	case TGErrorInvalidMessageLength:
		return DefaultTGInvalidMessageLength()
	case TGErrorIOException:
		return DefaultTGIOException()
	case TGErrorProtocolNotSupported:
		return DefaultTGProtocolNotSupported()
	case TGErrorRetryIOException:
		return DefaultTGRetryIOException()
	case TGErrorSecurityException:
		return DefaultTGSecurityException()
	case TGErrorTransactionException:
		return DefaultTGTransactionException()
	case TGErrorTypeCoercionNotSupported:
		return DefaultTGTypeCoercionNotSupported()
	case TGErrorTypeNotSupported:
		return DefaultTGTypeNotSupported()
	case TGErrorVersionMismatchException:
		return DefaultTGVersionMismatchException()

	case TGErrorInvalidErrorCode:
		fallthrough
	default:
		return GetPreDefinedErrors("TGErrorInvalidErrorCode")
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

// Get the error struct with timestamp
func GetErrorByType(excpTypeId int, errorCode, errorMsg, errorDetails string) tgdb.TGError {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputExcpTypeId := excpTypeId

	// Use a switch case to switch between exception types, if a type exist then error is nil (null)
	// Whenever new exception type gets into the mix, just add a case below
	switch inputExcpTypeId {
	case TGSuccess:
		return NewTGSuccess(errorCode, excpTypeId, errorMsg, errorDetails)

	case TGErrorBadAuthentication:
		return NewTGBadAuthentication(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorBadMagic:
		return NewTGBadMagic(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorBadVerb:
		return NewTGBadVerb(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorChannelDisconnected:
		return NewTGChannelDisconnected(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorConnectionTimeout:
		return NewTGConnectionTimeout(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorGeneralException:
		return NewTGGeneralException(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorInvalidMessageLength:
		return NewTGInvalidMessageLength(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorIOException:
		return NewTGIOException(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorProtocolNotSupported:
		return NewTGProtocolNotSupported(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorRetryIOException:
		return NewTGRetryIOException(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorSecurityException:
		return NewTGSecurityException(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorTransactionException:
		return NewTGTransactionException(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorTypeCoercionNotSupported:
		return NewTGTypeCoercionNotSupported(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorTypeNotSupported:
		return NewTGTypeNotSupported(errorCode, excpTypeId, errorMsg, errorDetails)
	case TGErrorVersionMismatchException:
		return NewTGVersionMismatchException(errorCode, excpTypeId, errorMsg, errorDetails)

	case TGErrorInvalidErrorCode:
		fallthrough
	default:
		return GetPreDefinedErrors("TGErrorInvalidErrorCode")
	}
	return nil
}

func GetErrorDetails(excpTypeId int) string {
	tgDbError := CreateExceptionByType(excpTypeId)
	if tgDbError != nil {
		return tgDbError.GetErrorDetails()
	}
	return tgDbError.GetErrorDetails()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func Error(excpTypeId int) string {
	tgDbError := CreateExceptionByType(excpTypeId)
	if tgDbError != nil {
		return tgDbError.Error()
	}
	return tgDbError.Error()
}


type Success struct {
	*TGDBError
}

// Create New Success Instance
func DefaultTGSuccess() *Success {
	newException := Success{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGSuccess
	return &newException
}

func NewTGSuccess(eCode string, eType int, eMsg, eDetails string) *Success {
	newException := DefaultTGSuccess()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGSuccessWithMsg(msg string) *Success {
	newException := DefaultTGSuccess()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *Success) GetErrorCode() string {
	return e.ErrorCode
}

func (e *Success) GetErrorType() int {
	return e.ErrorType
}

func (e *Success) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *Success) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *Success) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type BadAuthentication struct {
	*TGDBError
	realm string
}

// Create New BadAuthentication Instance
func DefaultTGBadAuthentication() *BadAuthentication {
	newException := BadAuthentication{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorBadAuthentication
	newException.realm = ""
	return &newException
}

func NewTGBadAuthentication(eCode string, eType int, eMsg, eDetails string) *BadAuthentication {
	newException := DefaultTGBadAuthentication()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGBadAuthenticationWithMsg(msg string) *BadAuthentication {
	newException := DefaultTGBadAuthentication()
	newException.ErrorMsg = msg
	return newException
}

func NewTGBadAuthenticationWithRealm(eCode string, eType int, eMsg, eDetails string, realm string) *BadAuthentication {
	newException := NewTGBadAuthentication(eCode, eType, eMsg, eDetails)
	newException.realm = realm
	return newException
}

/////////////////////////////////////////////////////////////////
// Helper functions for BadAuthentication
/////////////////////////////////////////////////////////////////

func (e *BadAuthentication) GetRealm() string {
	return e.realm
}

func (e *BadAuthentication) SetRealm(realmStr string) {
	e.realm = realmStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *BadAuthentication) GetErrorCode() string {
	return e.ErrorCode
}

func (e *BadAuthentication) GetErrorType() int {
	return e.ErrorType
}

func (e *BadAuthentication) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *BadAuthentication) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *BadAuthentication) Error() string {
	errMsg := fmt.Sprintf("realm: %s, ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.realm, e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type BadMagic struct {
	*TGDBError
}

// Create New BadMagic Instance
func DefaultTGBadMagic() *BadMagic {
	newException := BadMagic{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorBadMagic
	return &newException
}

func NewTGBadMagic(eCode string, eType int, eMsg, eDetails string) *BadMagic {
	newException := DefaultTGBadMagic()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGBadMagicWithMsg(msg string) *BadMagic {
	newException := DefaultTGBadMagic()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *BadMagic) GetErrorCode() string {
	return e.ErrorCode
}

func (e *BadMagic) GetErrorType() int {
	return e.ErrorType
}

func (e *BadMagic) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *BadMagic) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *BadMagic) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}



type BadVerb struct {
	*TGDBError
}

// Create New BadVerb Instance
func DefaultTGBadVerb() *BadVerb {
	newException := BadVerb{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorBadVerb
	return &newException
}

func NewTGBadVerb(eCode string, eType int, eMsg, eDetails string) *BadVerb {
	newException := DefaultTGBadVerb()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGBadVerbWithMsg(msg string) *BadVerb {
	newException := DefaultTGBadVerb()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *BadVerb) GetErrorCode() string {
	return e.ErrorCode
}

func (e *BadVerb) GetErrorType() int {
	return e.ErrorType
}

func (e *BadVerb) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *BadVerb) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *BadVerb) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type ChannelDisconnected struct {
	*TGDBError
}

// Create New ChannelDisconnected Instance
func DefaultTGChannelDisconnected() *ChannelDisconnected {
	newException := ChannelDisconnected{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorChannelDisconnected
	return &newException
}

func NewTGChannelDisconnected(eCode string, eType int, eMsg, eDetails string) *ChannelDisconnected {
	newException := DefaultTGChannelDisconnected()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGChannelDisconnectedWithMsg(msg string) *ChannelDisconnected {
	newException := DefaultTGChannelDisconnected()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *ChannelDisconnected) GetErrorCode() string {
	return e.ErrorCode
}

func (e *ChannelDisconnected) GetErrorType() int {
	return e.ErrorType
}

func (e *ChannelDisconnected) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *ChannelDisconnected) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *ChannelDisconnected) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}



type ConnectionTimeout struct {
	*TGDBError
}

// Create New ConnectionTimeout Instance
func DefaultTGConnectionTimeout() *ConnectionTimeout {
	newException := ConnectionTimeout{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorConnectionTimeout
	return &newException
}

func NewTGConnectionTimeout(eCode string, eType int, eMsg, eDetails string) *ConnectionTimeout {
	newException := DefaultTGConnectionTimeout()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGConnectionTimeoutWithMsg(msg string) *ConnectionTimeout {
	newException := DefaultTGConnectionTimeout()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *ConnectionTimeout) GetErrorCode() string {
	return e.ErrorCode
}

func (e *ConnectionTimeout) GetErrorType() int {
	return e.ErrorType
}

func (e *ConnectionTimeout) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *ConnectionTimeout) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *ConnectionTimeout) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type GeneralException struct {
	*TGDBError
}

// Create New GeneralException Instance
func DefaultTGGeneralException() *GeneralException {
	newException := GeneralException{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorGeneralException
	return &newException
}

func NewTGGeneralException(eCode string, eType int, eMsg, eDetails string) *GeneralException {
	newException := DefaultTGGeneralException()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGGeneralExceptionWithMsg(msg string) *GeneralException {
	newException := DefaultTGGeneralException()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////


func (e *GeneralException) GetErrorCode() string {
	return e.ErrorCode
}

func (e *GeneralException) GetErrorType() int {
	return e.ErrorType
}

func (e *GeneralException) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *GeneralException) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *GeneralException) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}



type InvalidMessageLength struct {
	*TGDBError
}

// Create New InvalidMessageLength Instance
func DefaultTGInvalidMessageLength() *InvalidMessageLength {
	newException := InvalidMessageLength{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorInvalidMessageLength
	return &newException
}

func NewTGInvalidMessageLength(eCode string, eType int, eMsg, eDetails string) *InvalidMessageLength {
	newException := DefaultTGInvalidMessageLength()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGInvalidMessageLengthWithMsg(msg string) *InvalidMessageLength {
	newException := DefaultTGInvalidMessageLength()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *InvalidMessageLength) GetErrorCode() string {
	return e.ErrorCode
}

func (e *InvalidMessageLength) GetErrorType() int {
	return e.ErrorType
}

func (e *InvalidMessageLength) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *InvalidMessageLength) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *InvalidMessageLength) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type IOException struct {
	*TGDBError
}

// Create New IOException Instance
func DefaultTGIOException() *IOException {
	newException := IOException{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorIOException
	return &newException
}

func NewTGIOException(eCode string, eType int, eMsg, eDetails string) *IOException {
	newException := DefaultTGIOException()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGIOExceptionWithMsg(msg string) *IOException {
	newException := DefaultTGIOException()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *IOException) GetErrorCode() string {
	return e.ErrorCode
}

func (e *IOException) GetErrorType() int {
	return e.ErrorType
}

func (e *IOException) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *IOException) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *IOException) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type ProtocolNotSupported struct {
	*TGDBError
}

// Create New ProtocolNotSupported Instance
func DefaultTGProtocolNotSupported() *ProtocolNotSupported {
	newException := ProtocolNotSupported{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorProtocolNotSupported
	return &newException
}

func NewTGProtocolNotSupported(eCode string, eType int, eMsg, eDetails string) *ProtocolNotSupported {
	newException := DefaultTGProtocolNotSupported()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGProtocolNotSupportedWithMsg(msg string) *ProtocolNotSupported {
	newException := DefaultTGProtocolNotSupported()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *ProtocolNotSupported) GetErrorCode() string {
	return e.ErrorCode
}

func (e *ProtocolNotSupported) GetErrorType() int {
	return e.ErrorType
}

func (e *ProtocolNotSupported) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *ProtocolNotSupported) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *ProtocolNotSupported) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type RetryIOException struct {
	*TGDBError
}

// Create New RetryIOException Instance
func DefaultTGRetryIOException() *RetryIOException {
	newException := RetryIOException{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorRetryIOException
	return &newException
}

func NewTGRetryIOException(eCode string, eType int, eMsg, eDetails string) *RetryIOException {
	newException := DefaultTGRetryIOException()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGRetryIOExceptionWithMsg(msg string) *RetryIOException {
	newException := DefaultTGRetryIOException()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *RetryIOException) GetErrorCode() string {
	return e.ErrorCode
}

func (e *RetryIOException) GetErrorType() int {
	return e.ErrorType
}

func (e *RetryIOException) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *RetryIOException) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *RetryIOException) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type SecurityException struct {
	*TGDBError
}

// Create New SecurityException Instance
func DefaultTGSecurityException() *SecurityException {
	newException := SecurityException{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorSecurityException
	return &newException
}

func NewTGSecurityException(eCode string, eType int, eMsg, eDetails string) *SecurityException {
	newException := DefaultTGSecurityException()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGSecurityExceptionWithMsg(msg string) *SecurityException {
	newException := DefaultTGSecurityException()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *SecurityException) GetErrorCode() string {
	return e.ErrorCode
}

func (e *SecurityException) GetErrorType() int {
	return e.ErrorType
}

func (e *SecurityException) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *SecurityException) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *SecurityException) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}



type TransactionException struct {
	*TGDBError
}

// Create New TransactionException Instance
func DefaultTGTransactionException() *TransactionException {
	newException := TransactionException{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorTransactionException
	return &newException
}

func NewTGTransactionExceptionWithMsg(msg string) *TransactionException {
	newException := DefaultTGTransactionException()
	newException.ErrorMsg = msg
	return newException
}

func NewTGTransactionException(eCode string, eType int, eMsg, eDetails string) *TransactionException {
	newException := DefaultTGTransactionException()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TransactionException) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TransactionException) GetErrorType() int {
	return e.ErrorType
}

func (e *TransactionException) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TransactionException) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TransactionException) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

func BuildException(ts tgdb.TGTransactionStatus, msg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(msg)
	switch ts {
	case tgdb.TGTransactionAlreadyInProgress:
		newException = NewTGTransactionAlreadyInProgressException(msg)
	case tgdb.TGTransactionMalFormed:
		newException = NewTGTransactionMalFormed(msg)
	case tgdb.TGTransactionGeneralError:
		newException = NewTGTransactionGeneralError(msg)
	case tgdb.TGTransactionVerificationError:
		newException = NewTGTransactionVerificationError(msg)
	case tgdb.TGTransactionInBadState:
		newException = NewTGTransactionInBadState(msg)
	case tgdb.TGTransactionUniqueConstraintViolation:
		newException = NewTGTransactionUniqueConstraintViolation(msg)
	case tgdb.TGTransactionOptimisticLockFailed:
		newException = NewTGTransactionOptimisticLockFailed(msg)
	case tgdb.TGTransactionResourceExceeded:
		newException = NewTGTransactionResourceExceeded(msg)
	case tgdb.TGTransactionUniqueIndexKeyAttributeNullError:
		newException = NewTGTransactionUniqueIndexKeyAttributeNullError(msg)
	default:
		newException = NewTGTransactionExceptionWithMsg(msg)
	}
	return newException
}

////////// TGTransactionAlreadyInProgressException //////////
type TGTransactionAlreadyInProgressException struct {
	*TransactionException
}

func NewTGTransactionAlreadyInProgressException(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionAlreadyInProgressException) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionAlreadyInProgressException) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionAlreadyInProgressException) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionAlreadyInProgressException) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionAlreadyInProgressException) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

////////// TGTransactionMalFormed //////////
type TGTransactionMalFormed struct {
	*TransactionException
}

func NewTGTransactionMalFormed(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionMalFormed) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionMalFormed) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionMalFormed) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionMalFormed) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionMalFormed) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

////////// TGTransactionGeneralError //////////
type TGTransactionGeneralError struct {
	*TransactionException
}

func NewTGTransactionGeneralError(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionGeneralError) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionGeneralError) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionGeneralError) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionGeneralError) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionGeneralError) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

////////// TGTransactionVerificationError //////////
type TGTransactionVerificationError struct {
	*TransactionException
}

func NewTGTransactionVerificationError(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionVerificationError) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionVerificationError) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionVerificationError) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionVerificationError) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionVerificationError) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

////////// TGTransactionInBadState //////////
type TGTransactionInBadState struct {
	*TransactionException
}

func NewTGTransactionInBadState(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionInBadState) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionInBadState) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionInBadState) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionInBadState) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionInBadState) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

////////// TGTransactionUniqueConstraintViolation //////////
type TGTransactionUniqueConstraintViolation struct {
	*TransactionException
}

func NewTGTransactionUniqueConstraintViolation(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionUniqueConstraintViolation) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionUniqueConstraintViolation) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionUniqueConstraintViolation) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionUniqueConstraintViolation) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionUniqueConstraintViolation) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

////////// TGTransactionOptimisticLockFailed //////////
type TGTransactionOptimisticLockFailed struct {
	*TransactionException
}

func NewTGTransactionOptimisticLockFailed(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionOptimisticLockFailed) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionOptimisticLockFailed) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionOptimisticLockFailed) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionOptimisticLockFailed) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionOptimisticLockFailed) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

////////// TGTransactionResourceExceeded //////////
type TGTransactionResourceExceeded struct {
	*TransactionException
}

func NewTGTransactionResourceExceeded(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionResourceExceeded) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionResourceExceeded) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionResourceExceeded) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionResourceExceeded) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionResourceExceeded) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}

////////// TGTransactionUniqueIndexKeyAttributeNullError //////////
type TGTransactionUniqueIndexKeyAttributeNullError struct {
	*TransactionException
}

func NewTGTransactionUniqueIndexKeyAttributeNullError(eMsg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(eMsg)
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TGTransactionUniqueIndexKeyAttributeNullError) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TGTransactionUniqueIndexKeyAttributeNullError) GetErrorType() int {
	return e.ErrorType
}

func (e *TGTransactionUniqueIndexKeyAttributeNullError) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TGTransactionUniqueIndexKeyAttributeNullError) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TGTransactionUniqueIndexKeyAttributeNullError) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type TypeCoercionNotSupported struct {
	*TGDBError
}

// Create New TypeCoercionNotSupported Instance
func DefaultTGTypeCoercionNotSupported() *TypeCoercionNotSupported {
	newException := TypeCoercionNotSupported{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorTypeCoercionNotSupported
	return &newException
}

func NewTGTypeCoercionNotSupported(eCode string, eType int, eMsg, eDetails string) *TypeCoercionNotSupported {
	newException := DefaultTGTypeCoercionNotSupported()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGTypeCoercionNotSupportedAttr(fromAttrTypeName, toAttrTypeName string) *TypeCoercionNotSupported {
	newException := DefaultTGTypeCoercionNotSupported()
	newException.ErrorMsg = fmt.Sprintf("Cannot coerce value of desc: '%s' to desc: '%s'", fromAttrTypeName, toAttrTypeName)
	return newException
}

func NewTGTypeCoercionNotSupportedWithMsg(msg string) *TypeCoercionNotSupported {
	newException := DefaultTGTypeCoercionNotSupported()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TypeCoercionNotSupported) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TypeCoercionNotSupported) GetErrorType() int {
	return e.ErrorType
}

func (e *TypeCoercionNotSupported) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TypeCoercionNotSupported) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TypeCoercionNotSupported) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type TypeNotSupported struct {
	*TGDBError
}

// Create New TypeNotSupported Instance
func DefaultTGTypeNotSupported() *TypeNotSupported {
	newException := TypeNotSupported{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorTypeNotSupported
	return &newException
}

func NewTGTypeNotSupported(eCode string, eType int, eMsg, eDetails string) *TypeNotSupported {
	newException := DefaultTGTypeNotSupported()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGTypeNotSupportedAttr(attrTypeName string) *TypeNotSupported {
	newException := DefaultTGTypeNotSupported()
	newException.ErrorMsg = fmt.Sprintf("Attribute descriptor: '%s' not supported", attrTypeName)
	return newException
}

func NewTGTypeNotSupportedWithMsg(msg string) *TypeNotSupported {
	newException := DefaultTGTypeNotSupported()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *TypeNotSupported) GetErrorCode() string {
	return e.ErrorCode
}

func (e *TypeNotSupported) GetErrorType() int {
	return e.ErrorType
}

func (e *TypeNotSupported) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *TypeNotSupported) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *TypeNotSupported) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}


type VersionMismatchException struct {
	*TGDBError
}

// Create New VersionMismatchException Instance
func DefaultTGVersionMismatchException() *VersionMismatchException {
	newException := VersionMismatchException{
		TGDBError: DefaultTGDBError(),
	}
	newException.ErrorType = TGErrorVersionMismatchException
	return &newException
}

func NewTGVersionMismatchException(eCode string, eType int, eMsg, eDetails string) *VersionMismatchException {
	newException := DefaultTGVersionMismatchException()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGVersionMismatchExceptionAttr(attrTypeName string) *VersionMismatchException {
	newException := DefaultTGVersionMismatchException()
	newException.ErrorMsg = fmt.Sprintf("Attribute descriptor: '%s' not supported", attrTypeName)
	return newException
}

func NewTGVersionMismatchExceptionWithMsg(msg string) *VersionMismatchException {
	newException := DefaultTGVersionMismatchException()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *VersionMismatchException) GetErrorCode() string {
	return e.ErrorCode
}

func (e *VersionMismatchException) GetErrorType() int {
	return e.ErrorType
}

func (e *VersionMismatchException) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *VersionMismatchException) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *VersionMismatchException) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}




///////////////// Exception Handling for TGDBQuery related error
type QueryError struct {
	*TGDBError
}

func DefaultQueryError() *QueryError {
	newException := QueryError {
		TGDBError: DefaultTGDBError(),
	}
	return &newException
}

func NewQueryError(eServerErrorCode int, eType int, eMsg string) *QueryError {
	newException := DefaultQueryError()
	newException.ErrorServerErrorCode = eServerErrorCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	return newException
}


func buildQueryException(ts int, msg string, serverCode int ) *QueryError {
	switch(ts) {
	case TGQueryProviderNotInitialized:
		return NewQueryError(serverCode, TGQryProviderNotInitialized, msg)
	case TGQueryParsingError:
		return NewQueryError(serverCode, TGQryParsingError, msg);
	case TGQueryStepNotSupported:
		return NewQueryError(serverCode, TGQryStepNotSupported, msg);
	case TGQueryStepNotAllowed:
		return NewQueryError(serverCode, TGQryStepNotAllowed, msg);
	case TGQueryStepArgMissing:
		return NewQueryError(serverCode, TGQryStepArgMissing, msg);
	case TGQueryStepArgNotSupported:
		return NewQueryError(serverCode, TGQryStepArgNotSupported, msg);
	case TGQueryStepMissing:
		return NewQueryError(serverCode, TGQryStepMissing, msg);
	case TGQueryNotDefined:
		return NewQueryError(serverCode, TGQryNotDefined, msg);
	case TGQueryAttrDescNotFound:
		return NewQueryError(serverCode, TGQryAttrDescNotFound, msg);
	case TGQueryEdgeTypeNotFound:
		return NewQueryError(serverCode, TGQryEdgeTypeNotFound, msg);
	case TGQueryNodeTypeNotFound:
		return NewQueryError(serverCode, TGQryNodeTypeNotFound, msg);
	case TGQueryInternalDataMismatchError:
		return NewQueryError(serverCode, TGQryInternalDataMismatchError, msg);
	case TGQueryStepSignatureNotSupported:
		return NewQueryError(serverCode, TGQryStepSignatureNotSupported, msg);
	case TGQueryInvalidDataType:
		return NewQueryError(serverCode, TGQryInvalidDataType, msg);
	default:
		return NewQueryError(serverCode, TGQryError, msg);
	}
}