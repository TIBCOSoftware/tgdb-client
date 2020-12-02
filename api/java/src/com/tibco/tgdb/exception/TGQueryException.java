/**
 * Copyright (c) 2019 TIBCO Software Inc.
 * All rights reserved.
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
 * <p/>
 * File name: TGQueryException.java
 * Created on: 1/10/20
 * Created by: chung
 * <p/>
 * SVN Id: $Id: TGQueryException.java 3158 2019-04-26 20:49:24Z kattaylo $
 */
package com.tibco.tgdb.exception;

import com.tibco.tgdb.pdu.impl.QueryResponse;

import static com.tibco.tgdb.exception.TGException.TGExceptionType.*;
import static com.tibco.tgdb.pdu.impl.QueryResponse.QueryErrorStatus;
import static com.tibco.tgdb.pdu.impl.QueryResponse.QueryErrorStatus.*;

public class TGQueryException extends TGException {

    private TGQueryException(String ex, TGExceptionType et, int serverCode) {
        super(ex, et, serverCode);
    }

    public static TGQueryException buildException(QueryErrorStatus ts, String msg, int serverCode) {
        switch (ts) {
            case TGQueryProviderNotInitialized:
                return new TGQueryException(msg, QryProviderNotInitialized, serverCode);
            case TGQueryParsingError:
                return new TGQueryException(msg, QryParsingError, serverCode);
            case TGQueryStepNotSupported:
                return new TGQueryException(msg, QryStepNotSupported, serverCode);
            case TGQueryStepNotAllowed:
                return new TGQueryException(msg, QryStepNotAllowed, serverCode);
            case TGQueryStepArgMissing:
                return new TGQueryException(msg, QryStepArgMissing, serverCode);
            case TGQueryStepArgNotSupported:
                return new TGQueryException(msg, QryStepArgNotSupported, serverCode);
            case TGQueryStepMissing:
                return new TGQueryException(msg, QryStepMissing, serverCode);
            case TGQueryNotDefined:
                return new TGQueryException(msg, QryNotDefined, serverCode);
            case TGQueryAttrDescNotFound:
                return new TGQueryException(msg, QryAttrDescNotFound, serverCode);
            case TGQueryEdgeTypeNotFound:
                return new TGQueryException(msg, QryEdgeTypeNotFound, serverCode);
            case TGQueryNodeTypeNotFound:
                return new TGQueryException(msg, QryNodeTypeNotFound, serverCode);
            case TGQueryInternalDataMismatchError:
                return new TGQueryException(msg, QryInternalDataMismatchError, serverCode);
            case TGQueryStepSignatureNotSupported:
                return new TGQueryException(msg, QryStepSignatureNotSupported, serverCode);
            case TGQueryInvalidDataType:
                return new TGQueryException(msg, QryInvalidDataType, serverCode);
            case TGQueryExecSPFailure:
                return new TGQueryException(msg, QryExecSPFailure, serverCode);
            case TGQuerySPNotFound:
                return new TGQueryException(msg, QrySPNotFound, serverCode);
            case TGQuerySPArgMissing:
                return new TGQueryException(msg, QrySPArgMissing, serverCode);
            case TGQueryStepArgInvalid:
                return new TGQueryException(msg, QryStepArgInvalid, serverCode);
            case TGQueryStepModulationInvalid:
                return new TGQueryException(msg, QryStepModulationInvalid, serverCode);
            case TGQueryAccessDenied:
                return new TGQueryException(msg, QryAccessDenied, serverCode);
            default:
                return new TGQueryException(msg, QryError, serverCode);
        }
    }
}

