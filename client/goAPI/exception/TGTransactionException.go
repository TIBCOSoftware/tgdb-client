package exception

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * File name: TGErrorTransactionException.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TransactionException struct {
	*types.TGDBError
}

// Create New TransactionException Instance
func DefaultTGTransactionException() *TransactionException {
	newException := TransactionException{
		TGDBError: types.DefaultTGDBError(),
	}
	newException.ErrorType = types.TGErrorTransactionException
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

func BuildException(ts types.TGTransactionStatus, msg string) *TransactionException {
	newException := NewTGTransactionExceptionWithMsg(msg)
	switch ts {
	case types.TGTransactionAlreadyInProgress:
		newException = NewTGTransactionAlreadyInProgressException(msg)
	case types.TGTransactionMalFormed:
		newException = NewTGTransactionMalFormed(msg)
	case types.TGTransactionGeneralError:
		newException = NewTGTransactionGeneralError(msg)
	case types.TGTransactionVerificationError:
		newException = NewTGTransactionVerificationError(msg)
	case types.TGTransactionInBadState:
		newException = NewTGTransactionInBadState(msg)
	case types.TGTransactionUniqueConstraintViolation:
		newException = NewTGTransactionUniqueConstraintViolation(msg)
	case types.TGTransactionOptimisticLockFailed:
		newException = NewTGTransactionOptimisticLockFailed(msg)
	case types.TGTransactionResourceExceeded:
		newException = NewTGTransactionResourceExceeded(msg)
	case types.TGTransactionUniqueIndexKeyAttributeNullError:
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
