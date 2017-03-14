/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * limitations under the License.*
 *
 * <p/>
 * File name : TGTransactionException.java
 * Created on: 10/4/16
 * Created by: suresh
 * <p/>
 *
 * SVN Id: $Id$
 */


package com.tibco.tgdb.exception;

import com.tibco.tgdb.pdu.impl.CommitTransactionResponse;

public class TGTransactionException extends TGException {

    private TGTransactionException(String ex) {
        super(ex);
    }

    public static TGTransactionException buildException(CommitTransactionResponse.TransactionStatus ts, String msg) {
        switch (ts) {
            case TGTransactionAlreadyInProgress:
                return new TGTransactionAlreadyInProgressException(msg);

            case TGTransactionMalFormed:
                return new TGTransactionMalFormedException(msg);

            case TGTransactionGeneralError:
                return new TGTransactionGeneralErrorException(msg);

            case TGTransactionInBadState:
                return new TGTransactionInBadStateException(msg);

            case TGTransactionVerificationError:
                return new TGTransactionVerificationErrorException(msg);

            case TGTransactionUniqueConstraintViolation:
                return new TGTransactionUniqueConstraintViolationException(msg);

            case TGTransactionOptimisticLockFailed:
                return new TGTransactionOptimisticLockFailedException(msg);

            case TGTransactionResourceExceeded:
                return new TGTransactionResourceExceededException(msg);

            case TGTransactionUniqueIndexKeyAttributeNullError:
                return new TGTransactionUniqueIndexKeyAttributeNullError(msg);
        }
        return new TGTransactionException(msg);
    }


    private static class TGTransactionAlreadyInProgressException extends TGTransactionException {
        TGTransactionAlreadyInProgressException(String ex) {
            super(ex);
        }
    }

    private static class TGTransactionMalFormedException extends TGTransactionException {
        TGTransactionMalFormedException(String ex) {
            super(ex);
        }
    }

    private static class TGTransactionGeneralErrorException extends TGTransactionException {
        TGTransactionGeneralErrorException(String ex) {
            super(ex);
        }
    }

    private static class TGTransactionInBadStateException extends TGTransactionException {
        TGTransactionInBadStateException(String ex) {
            super(ex);
        }
    }

    private static class TGTransactionVerificationErrorException extends TGTransactionException {
        TGTransactionVerificationErrorException(String ex) {
            super(ex);
        }
    }
    private static class TGTransactionUniqueConstraintViolationException extends TGTransactionException {
        TGTransactionUniqueConstraintViolationException(String ex) {
            super(ex);
        }
    }

    private static class TGTransactionOptimisticLockFailedException extends TGTransactionException {
        TGTransactionOptimisticLockFailedException(String ex) {
            super(ex);
        }
    }

    private static class TGTransactionResourceExceededException extends TGTransactionException {
        TGTransactionResourceExceededException(String ex) {
            super(ex);
        }
    }

    private static class TGTransactionUniqueIndexKeyAttributeNullError extends TGTransactionException {
        TGTransactionUniqueIndexKeyAttributeNullError(String ex) { super(ex); }
    }
}

