package types

import "bytes"

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
 * File name: TGTransactionStatus.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= Various Transaction Status Returned from TGDB server =======
type TGTransactionStatus int

const (
	lastStatus                                                        = 8000
	TGTransactionInvalid                          TGTransactionStatus = -1
	TGTransactionSuccess                          TGTransactionStatus = 0
	TGTransactionAlreadyInProgress                TGTransactionStatus = lastStatus + 1
	TGTransactionClientDisconnected               TGTransactionStatus = lastStatus + 2
	TGTransactionMalFormed                        TGTransactionStatus = lastStatus + 3
	TGTransactionGeneralError                     TGTransactionStatus = lastStatus + 4
	TGTransactionVerificationError                TGTransactionStatus = lastStatus + 5
	TGTransactionInBadState                       TGTransactionStatus = lastStatus + 6
	TGTransactionUniqueConstraintViolation        TGTransactionStatus = lastStatus + 7
	TGTransactionOptimisticLockFailed             TGTransactionStatus = lastStatus + 8
	TGTransactionResourceExceeded                 TGTransactionStatus = lastStatus + 9
	TGCurrentThreadNotInTransaction               TGTransactionStatus = lastStatus + 10
	TGTransactionUniqueIndexKeyAttributeNullError TGTransactionStatus = lastStatus + 11
)

func (txnStatus TGTransactionStatus) FromStatus(status int) TGTransactionStatus {
	if txnStatus&TGTransactionSuccess == TGTransactionSuccess {
		return TGTransactionSuccess
	} else if txnStatus&TGTransactionAlreadyInProgress == TGTransactionAlreadyInProgress {
		return TGTransactionAlreadyInProgress
	} else if txnStatus&TGTransactionClientDisconnected == TGTransactionClientDisconnected {
		return TGTransactionClientDisconnected
	} else if txnStatus&TGTransactionMalFormed == TGTransactionMalFormed {
		return TGTransactionMalFormed
	} else if txnStatus&TGTransactionGeneralError == TGTransactionGeneralError {
		return TGTransactionGeneralError
	} else if txnStatus&TGTransactionVerificationError == TGTransactionVerificationError {
		return TGTransactionVerificationError
	} else if txnStatus&TGTransactionInBadState == TGTransactionInBadState {
		return TGTransactionInBadState
	} else if txnStatus&TGTransactionUniqueConstraintViolation == TGTransactionUniqueConstraintViolation {
		return TGTransactionUniqueConstraintViolation
	} else if txnStatus&TGTransactionOptimisticLockFailed == TGTransactionOptimisticLockFailed {
		return TGTransactionOptimisticLockFailed
	} else if txnStatus&TGTransactionResourceExceeded == TGTransactionResourceExceeded {
		return TGTransactionResourceExceeded
	} else if txnStatus&TGCurrentThreadNotInTransaction == TGCurrentThreadNotInTransaction {
		return TGCurrentThreadNotInTransaction
	} else if txnStatus&TGTransactionUniqueIndexKeyAttributeNullError == TGTransactionUniqueIndexKeyAttributeNullError {
		return TGTransactionUniqueIndexKeyAttributeNullError
	}
	return TGTransactionInvalid
}

func (txnStatus TGTransactionStatus) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if txnStatus&TGTransactionInvalid == TGTransactionInvalid {
		buffer.WriteString("TransactionInvalid")
	} else if txnStatus&TGTransactionSuccess == TGTransactionSuccess {
		buffer.WriteString("TransactionSuccess")
	} else if txnStatus&TGTransactionAlreadyInProgress == TGTransactionAlreadyInProgress {
		buffer.WriteString("TransactionAlreadyInProgress")
	} else if txnStatus&TGTransactionClientDisconnected == TGTransactionClientDisconnected {
		buffer.WriteString("TransactionClientDisconnected")
	} else if txnStatus&TGTransactionMalFormed == TGTransactionMalFormed {
		buffer.WriteString("TransactionMalFormed")
	} else if txnStatus&TGTransactionGeneralError == TGTransactionGeneralError {
		buffer.WriteString("TransactionGeneralError")
	} else if txnStatus&TGTransactionVerificationError == TGTransactionVerificationError {
		buffer.WriteString("TransactionVerificationError")
	} else if txnStatus&TGTransactionInBadState == TGTransactionInBadState {
		buffer.WriteString("TransactionInBadState")
	} else if txnStatus&TGTransactionUniqueConstraintViolation == TGTransactionUniqueConstraintViolation {
		buffer.WriteString("TransactionUniqueConstraintViolation")
	} else if txnStatus&TGTransactionOptimisticLockFailed == TGTransactionOptimisticLockFailed {
		buffer.WriteString("TransactionOptimisticLockFailed")
	} else if txnStatus&TGTransactionResourceExceeded == TGTransactionResourceExceeded {
		buffer.WriteString("TransactionResourceExceeded")
	} else if txnStatus&TGCurrentThreadNotInTransaction == TGCurrentThreadNotInTransaction {
		buffer.WriteString("CurrentThreadNotInTransaction")
	} else if txnStatus&TGTransactionUniqueIndexKeyAttributeNullError == TGTransactionUniqueIndexKeyAttributeNullError {
		buffer.WriteString("TransactionUniqueIndexKeyAttributeNullError")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}
