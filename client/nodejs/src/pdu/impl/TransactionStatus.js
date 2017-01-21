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

exports.TransactionStatus = {
            TGTransactionInvalid                   : {value :  -1},
            TGTransactionSuccess                   : {value :  0},
            TGTransactionAlreadyInProgress         : {value :  1},
            TGTransactionClientDisconnected        : {value :  2},
            TGTransactionMalFormed                 : {value :  3},
            TGTransactionGeneralError              : {value :  4},
            TGTransactionVerificationError         : {value :  5},
            TGTransactionInBadState                : {value :  6},
            TGTransactionUniqueConstraintViolation : {value :  7},
            TGTransactionOptimisticLockFailed      : {value :  8},
            TGTransactionResourceExceeded          : {value :  9},
            TGCurrentThreadNotinTransaction        : {value :  10},
    	    fromStatus : function (status) {
            	switch(status) {
            		case -1 : return this.TGTransactionInvalid;
            		case 0 : return this.TGTransactionSuccess;
            		case 1 : return this.TGTransactionAlreadyInProgress;
            		case 2 : return this.TGTransactionClientDisconnected;
            		case 3 : return this.TGTransactionMalFormed;
            		case 4 : return this.TGTransactionGeneralError;
            		case 5 : return this.TGTransactionVerificationError;
            		case 6 : return this.TGTransactionInBadState;
            		case 7 : return this.TGTransactionUniqueConstraintViolation;
            		case 8 : return this.TGTransactionOptimisticLockFailed;
            		case 9 : return this.TGTransactionResourceExceeded;
            		case 10 : return this.TGCurrentThreadNotinTransaction;
    	    }
    	}
    };