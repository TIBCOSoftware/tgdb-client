package exception

import (
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"testing"
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
 * File name: ExceptionFactory_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func TestCreateExceptionByType(t *testing.T) {
	for inputExcpTypeId, predefExcp := range types.PreDefinedErrors {
		//if inputExcpTypeId != types.TGSuccess {
		//	continue
		//}
		errMsg := CreateExceptionByType(inputExcpTypeId)
		t.Logf("ExceptionFactory returned error message code %s for verbId: '%+v' as '%+v'", predefExcp.ErrorCode, inputExcpTypeId, errMsg)
	}
}

func TestGetErrorByType(t *testing.T) {
	for inputExcpTypeId, predefExcp := range types.PreDefinedErrors {
		//if inputExcpTypeId != types.TGSuccess {
		//	continue
		//}
		errMsg := GetErrorByType(inputExcpTypeId, predefExcp.ErrorCode, predefExcp.ErrorMsg, predefExcp.ErrorDetails)
		t.Logf("ExceptionFactory returned error message code %s for verbId: '%+v' as '%+v'", predefExcp.ErrorCode, inputExcpTypeId, errMsg)
	}
}

func TestGetDetailsWithContext(t *testing.T) {
	for inputExcpTypeId, predefExcp := range types.PreDefinedErrors {
		if inputExcpTypeId == types.TGErrorTypeCoercionNotSupported || inputExcpTypeId == types.TGErrorTransactionException ||
			inputExcpTypeId == types.TGErrorInvalidErrorCode  {
			continue
		}
		errMsg := GetErrorDetails(inputExcpTypeId)
		t.Logf("ExceptionFactory returned error message code %s for verbId: '%+v' as '%+v'", predefExcp.ErrorCode, inputExcpTypeId, errMsg)
	}
}
