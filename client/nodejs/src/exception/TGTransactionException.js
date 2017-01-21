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

function TGTransactionException(message) { 
  TGTransactionException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionException, Error);

function TGTransactionAlreadyInProgressException(message) {
	TGTransactionAlreadyInProgressException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionAlreadyInProgressException, TGTransactionException);

function TGTransactionMalFormedException(message) {
	TGTransactionMalFormedException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionMalFormedException, TGTransactionException);

function TGTransactionGeneralErrorException(message) {
	TGTransactionGeneralErrorException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionGeneralErrorException, TGTransactionException);

function TGTransactionInBadStateException(message) {
	TGTransactionInBadStateException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionInBadStateException, TGTransactionException);

function TGTransactionVerificationErrorException(message) {
	TGTransactionVerificationErrorException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionVerificationErrorException, TGTransactionException);

function TGTransactionUniqueConstraintViolationException(message) {
	TGTransactionUniqueConstraintViolationException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionUniqueConstraintViolationException, TGTransactionException);

function TGTransactionOptimisticLockFailedException(message) {
	TGTransactionOptimisticLockFailedException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionOptimisticLockFailedException, TGTransactionException);

function TGTransactionResourceExceededException(message) {
	TGTransactionResourceExceededException.super_.call(this);
  this.message = message;
}

util.inherits(TGTransactionResourceExceededException, TGTransactionException);

exports.TGTransactionException                          = TGTransactionException;
exports.TGTransactionAlreadyInProgressException         = TGTransactionAlreadyInProgressException;
exports.TGTransactionMalFormedException                 = TGTransactionMalFormedException;
exports.TGTransactionGeneralErrorException              = TGTransactionGeneralErrorException;
exports.TGTransactionInBadStateException                = TGTransactionInBadStateException;
exports.TGTransactionVerificationErrorException         = TGTransactionVerificationErrorException;
exports.TGTransactionUniqueConstraintViolationException = TGTransactionUniqueConstraintViolationException;
exports.TGTransactionOptimisticLockFailedException      = TGTransactionOptimisticLockFailedException;
exports.TGTransactionResourceExceededException          = TGTransactionResourceExceededException;
