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

var util                          = require('util'),
    VerbId                        = require('./VerbId').VerbId,
    TransactionStatus             = require('./TransactionStatus').TransactionStatus,
    Response                      = require('./Response').Response,
    TGTransactionExceptionBuilder = require('../../exception/TGTransactionExceptionBuilder'),
    TGException                   = require('../../exception/TGException').TGException,
    TGLogManager                  = require('../../log/TGLogManager'),
    TGLogLevel                    = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function CommitTransactionResponse () {
	CommitTransactionResponse.super_.call(this);
    this._addedIdList = []; 
    this._updatedIdList = []; 
    this._removedIdList = []; 
    this._attrDescIdList = []; 
    this._attrDescCount = 0;
    this._addedCount = 0;
    this._updatedCount = 0;
    this._removedCount = 0;
    
    this._entityStream = null;
}

util.inherits(CommitTransactionResponse, Response);

CommitTransactionResponse.prototype.writePayload = function (outputStream) {
};

CommitTransactionResponse.prototype.readPayload = function (inputStream) {
	inputStream.readInt(); // buf length
	inputStream.readInt(); // checksum
	var status = inputStream.readInt();// status code - currently zero
    var exception = this.processTransactionStatus(inputStream, status);
    if (exception !== null) {
    	this.setException(exception);
    	return;
    }
	
    while (inputStream.available() > 0) {
    	var opCode = inputStream.readShort();
    	switch (opCode) {
    		case 0x1010:
                this._attrDescCount = inputStream.readInt();
                for (var i = 0; i < this._attrDescCount; i++) {
                    this._attrDescIdList.push(inputStream.readInt());  // temp id ?
                    this._attrDescIdList.push(inputStream.readInt());  // real id ?
                }
                logger.logDebugWire( 
                		"[CommitTransactionResponse.prototype.readPayload] Received %d attr desc", 
                		this._attrDescCount);
    			break;
    		case 0x1011:
                this._addedCount = inputStream.readInt();
                for (i = 0; i < this._addedCount; i++) {
                    this._addedIdList.push(inputStream.readLongAsBytes());  // temp id
                    this._addedIdList.push(inputStream.readLongAsBytes());  // real id
                    this._addedIdList.push(inputStream.readLong());  // version
                }
                logger.logDebugWire( 
                		"[CommitTransactionResponse.prototype.readPayload] Received %d added entities", 
                		this._addedCount);
    			break;
            case 0x1012:
                this._updatedCount = inputStream.readInt();
                for (i = 0; i < this._updatedCount; i++) {
                    this._updatedIdList.push(inputStream.readLongAsBytes());  //id
                    this._updatedIdList.push(inputStream.readLong());  //version
                }
                logger.logDebugWire(
                		"[CommitTransactionResponse.prototype.readPayload] Received %d updated entities", 
                		this._updatedCount);
                break;
            case 0x1013:
                this._removedCount = inputStream.readInt();
                for (i = 0; i < this._removedCount; i++) {
                    this._removedIdList.push(inputStream.readLongAsBytes());  //id
                }
                logger.logDebugWire(
                		"[CommitTransactionResponse.prototype.readPayload] Received %d delete results", 
                		this._removedCount);
                break;
            case 0x6789:
            	this._entityStream = inputStream;
            	var pos = inputStream.getPosition();
            	this._entityCount = inputStream.readInt();
            	inputStream.setPosition(pos);
            	logger.logDebugWire(
            			"[CommitTransactionResponse.prototype.readPayload] Received %d debug entities", 
            			this._entityCount);
                return;
            default:
                break;
    		}
        }
};

CommitTransactionResponse.prototype.processTransactionStatus = function (inputStream, status) {
    var ts = TransactionStatus.fromStatus(status);
    var msg;
    switch (ts) {
        case TransactionStatus.TGTransactionSuccess:
            return null;

        case TransactionStatus.TGTransactionAlreadyInProgress:
        case TransactionStatus.TGTransactionClientDisconnected:
        case TransactionStatus.TGTransactionMalFormed:
        case TransactionStatus.TGTransactionGeneralError:
        case TransactionStatus.TGTransactionInBadState:
        case TransactionStatus.TGTransactionUniqueConstraintViolation:
        case TransactionStatus.TGTransactionOptimisticLockFailed:
        case TransactionStatus.TGTransactionResourceExceeded:
        default:
            try {
                msg = inputStream.readUTF();
            }
            catch (Error) {
                msg = "Error not available";
            }
            return TGTransactionExceptionBuilder.build(msg, ts);
    }

};

CommitTransactionResponse.prototype.isUpdateable = function () {
    return false;
};

CommitTransactionResponse.prototype.getVerbId = function () {
    return VerbId.COMMIT_TRANS_RESPONSE;
};

CommitTransactionResponse.prototype.getAttrDescCount = function () {
    return this._attrDescCount;
};

CommitTransactionResponse.prototype.getAddedEntityCount = function () {
    return this._addedCount;
}; 

CommitTransactionResponse.prototype.getUpdatedEntityCount = function () {
    return this._updatedCount;
};

CommitTransactionResponse.prototype.getAttrDescIdList = function () {
    return this._attrDescIdList;
};

CommitTransactionResponse.prototype.getAddedIdList = function () {
    return this._addedIdList;
};

CommitTransactionResponse.prototype.getUpdatedIdList = function () {
    return this._updatedIdList;
};

CommitTransactionResponse.prototype.getEntityStream = function () {
	return this._entityStream;
};

CommitTransactionResponse.prototype.hasException = function () { 
	//return exception !== null;
};

CommitTransactionResponse.prototype.getException = function () {
	//return exception;
};

exports.CommitTransactionResponse = CommitTransactionResponse;
