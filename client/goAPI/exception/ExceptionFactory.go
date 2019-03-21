package exception

import (
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
 * File name: TGExceptionFactory.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// Create new exception instance based on the input type
func CreateExceptionByType(excpTypeId int) types.TGError {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputExcpTypeId := excpTypeId

	// Use a switch case to switch between exception types, if a type exist then error is nil (null)
	// Whenever new exception type gets into the mix, just add a case below
	switch inputExcpTypeId {
	case types.TGSuccess:
		return DefaultTGSuccess()

	case types.TGErrorBadAuthentication:
		return DefaultTGBadAuthentication()
	case types.TGErrorBadMagic:
		return DefaultTGBadMagic()
	case types.TGErrorBadVerb:
		return DefaultTGBadVerb()
	case types.TGErrorChannelDisconnected:
		return DefaultTGChannelDisconnected()
	case types.TGErrorConnectionTimeout:
		return DefaultTGConnectionTimeout()
	case types.TGErrorGeneralException:
		return DefaultTGGeneralException()
	case types.TGErrorInvalidMessageLength:
		return DefaultTGInvalidMessageLength()
	case types.TGErrorIOException:
		return DefaultTGIOException()
	case types.TGErrorProtocolNotSupported:
		return DefaultTGProtocolNotSupported()
	case types.TGErrorRetryIOException:
		return DefaultTGRetryIOException()
	case types.TGErrorSecurityException:
		return DefaultTGSecurityException()
	case types.TGErrorTransactionException:
		return DefaultTGTransactionException()
	case types.TGErrorTypeCoercionNotSupported:
		return DefaultTGTypeCoercionNotSupported()
	case types.TGErrorTypeNotSupported:
		return DefaultTGTypeNotSupported()
	case types.TGErrorVersionMismatchException:
		return DefaultTGVersionMismatchException()

	case types.TGErrorInvalidErrorCode:
		fallthrough
	default:
		return types.GetPreDefinedErrors("TGErrorInvalidErrorCode")
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

// Get the error struct with timestamp
func GetErrorByType(excpTypeId int, errorCode, errorMsg, errorDetails string) types.TGError {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputExcpTypeId := excpTypeId

	// Use a switch case to switch between exception types, if a type exist then error is nil (null)
	// Whenever new exception type gets into the mix, just add a case below
	switch inputExcpTypeId {
	case types.TGSuccess:
		return NewTGSuccess(errorCode, excpTypeId, errorMsg, errorDetails)

	case types.TGErrorBadAuthentication:
		return NewTGBadAuthentication(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorBadMagic:
		return NewTGBadMagic(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorBadVerb:
		return NewTGBadVerb(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorChannelDisconnected:
		return NewTGChannelDisconnected(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorConnectionTimeout:
		return NewTGConnectionTimeout(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorGeneralException:
		return NewTGGeneralException(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorInvalidMessageLength:
		return NewTGInvalidMessageLength(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorIOException:
		return NewTGIOException(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorProtocolNotSupported:
		return NewTGProtocolNotSupported(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorRetryIOException:
		return NewTGRetryIOException(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorSecurityException:
		return NewTGSecurityException(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorTransactionException:
		return NewTGTransactionException(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorTypeCoercionNotSupported:
		return NewTGTypeCoercionNotSupported(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorTypeNotSupported:
		return NewTGTypeNotSupported(errorCode, excpTypeId, errorMsg, errorDetails)
	case types.TGErrorVersionMismatchException:
		return NewTGVersionMismatchException(errorCode, excpTypeId, errorMsg, errorDetails)

	case types.TGErrorInvalidErrorCode:
		fallthrough
	default:
		return types.GetPreDefinedErrors("TGErrorInvalidErrorCode")
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
