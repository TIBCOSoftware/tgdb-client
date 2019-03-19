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
 * File name: TGErrorBadAuthentication.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type BadAuthentication struct {
	*types.TGDBError
	realm string
}

// Create New BadAuthentication Instance
func DefaultTGBadAuthentication() *BadAuthentication {
	newException := BadAuthentication{
		TGDBError: types.DefaultTGDBError(),
	}
	newException.ErrorType = types.TGErrorBadAuthentication
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
