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
 * File name: TGTransactionException.java
 * Created on: 10/4/16
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGTransactionException.java 4018 2020-05-22 20:53:23Z ssubrama $
 */


package com.tibco.tgdb.exception;

import com.tibco.tgdb.pdu.impl.CommitTransactionResponse;
import static com.tibco.tgdb.exception.TGException.TGExceptionType.*;

public class TGTransactionException extends TGException {

    private TGTransactionException(String ex, TGExceptionType exceptionType, int serverCode) {
        super(ex, exceptionType, serverCode);
    }

    public static TGTransactionException buildException(CommitTransactionResponse.TransactionStatus ts, String msg,
                                                        int serverCode) {
        switch (ts) {
            case TGTransactionAlreadyInProgress:
                return new TGTransactionException(msg, TxnAlreadyInProgress, serverCode);
            case TGTransactionMalFormed:
                return new TGTransactionException(msg, TxnMalFormed, serverCode);
            case TGTransactionGeneralError:
                return new TGTransactionException(msg, TxnGeneralError, serverCode);
            case TGTransactionInBadState:
                return new TGTransactionException(msg, TxnInBadState, serverCode);
            case TGTransactionVerificationError:
                return new TGTransactionException(msg, TxnVerificationError, serverCode);
            case TGTransactionUniqueConstraintViolation:
                return new TGTransactionException(msg, TxnUniqueConstraintViolation, serverCode);
            case TGTransactionOptimisticLockFailed:
                return new TGTransactionException(msg, TxnOptimisticLockFailed, serverCode);
            case TGTransactionResourceExceeded:
                return new TGTransactionException(msg, TxnResourceExceeded, serverCode);
            case TGTransactionUniqueIndexKeyAttributeNullError:
                return new TGTransactionException(msg, TxnUniqueIndexKeyAttributeNullError, serverCode);
            case TGTxnInvalidPhase:
                return new TGTransactionException(msg, TxnInvalidPhase, serverCode);
            case TGTransactionWriteAheadLogBusy:
                return new TGTransactionException(msg, TxnWriteAheadLogBusy, serverCode);
            case TGTransactionWriteAheadLogFailed2Commit:
                return new TGTransactionException(msg, TxnWriteAheadLogFailed2Commit, serverCode);
            default:
                return new TGTransactionException(msg, TxnGeneralError, serverCode);
        }
    }
}

