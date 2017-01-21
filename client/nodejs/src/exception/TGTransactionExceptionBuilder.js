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
 * limitations under the License.
 */

var assert = require('assert');  
var util = require('util');
var TransactionStatus = require('../pdu/impl/TransactionStatus');
var TGTransactionExceptions = require('./TGTransactionException');

var TGTransactionExceptionBuilder = {
	build : function(msg, ts) {
	    switch (ts) {
            case TransactionStatus.TGTransactionAlreadyInProgress:
                return new TGTransactionExceptions.TGTransactionAlreadyInProgressException(msg);

            case TransactionStatus.TGTransactionMalFormed:
                return new TGTransactionExceptions.TGTransactionMalFormedException(msg);

            case TransactionStatus.TGTransactionGeneralError:
                return new TGTransactionExceptions.TGTransactionGeneralErrorException(msg);

            case TransactionStatus.TGTransactionInBadState:
                return new TGTransactionExceptions.TGTransactionInBadStateException(msg);

            case TransactionStatus.TGTransactionVerificationError:
                return new TGTransactionExceptions.TGTransactionVerificationErrorException(msg);

            case TransactionStatus.TGTransactionUniqueConstraintViolation:
                return new TGTransactionExceptions.TGTransactionUniqueConstraintViolationException(msg);

            case TransactionStatus.TGTransactionOptimisticLockFailed:
                return new TGTransactionExceptions.TGTransactionOptimisticLockFailedException(msg);

            case TransactionStatus.TGTransactionResourceExceeded:
                return new TGTransactionExceptions.TGTransactionResourceExceededException(msg);
                
            default :
                return new TGTransactionExceptions.TGTransactionException(msg);
        }
	}
};

module.exports = TGTransactionExceptionBuilder;
