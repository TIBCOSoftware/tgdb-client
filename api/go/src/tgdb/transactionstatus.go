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
 * File name: transactionstatus.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: transactionstatus.go 3513 2019-11-13 19:49:04Z nimish $
 */

package tgdb

import "bytes"

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

func FromStatus(status int) TGTransactionStatus {
	switch TGTransactionStatus(status) {
	case TGTransactionSuccess:
		return TGTransactionSuccess
	case TGTransactionAlreadyInProgress:
		return TGTransactionAlreadyInProgress
	case TGTransactionClientDisconnected:
		return TGTransactionClientDisconnected
	case TGTransactionMalFormed:
		return TGTransactionMalFormed
	case TGTransactionGeneralError:
		return TGTransactionGeneralError
	case TGTransactionVerificationError:
		return TGTransactionVerificationError
	case TGTransactionInBadState:
		return TGTransactionInBadState
	case TGTransactionUniqueConstraintViolation:
		return TGTransactionUniqueConstraintViolation
	case TGTransactionOptimisticLockFailed:
		return TGTransactionOptimisticLockFailed
	case TGTransactionResourceExceeded:
		return TGTransactionResourceExceeded
	case TGCurrentThreadNotInTransaction:
		return TGCurrentThreadNotInTransaction
	case TGTransactionUniqueIndexKeyAttributeNullError:
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

