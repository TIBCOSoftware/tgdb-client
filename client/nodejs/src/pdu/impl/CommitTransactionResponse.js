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

var util                    = require('util'),
    AbstractProtocolMessage = require('../AbstractProtocolMessage').AbstractProtocolMessage,
    VerbId                  = require('./VerbId').VerbId;

function CommitTransactionResponse () {
	CommitTransactionResponse.super_.call(this);
    this._addedIdList = []; 
    this._updatedIdList = []; 
    this._removedIdList = []; 
    this._attrDescIdList = []; 
    this._attrDescCount = 0;
    this._entityCount = 0;
}
    
util.inherits(CommitTransactionResponse, AbstractProtocolMessage);
    
CommitTransactionResponse.prototype.writePayload = function (outputStream) {
};

CommitTransactionResponse.prototype.readPayload = function (inputStream) {
	inputStream.readInt(); // buf length
	inputStream.readInt(); // checksum
    inputStream.readInt();// status code - currently zero
    while (inputStream.available() > 0) {
    	var opCode = inputStream.readShort();
    	switch (opCode) {
    		case 0x1010:
                this._attrDescCount = inputStream.readInt();
                for (var i = 0; i < this._attrDescCount; i++) {
                    this._attrDescIdList.push(inputStream.readInt());  // temp id ?
                    this._attrDescIdList.push(inputStream.readInt());  // real id ?
                }
                console.log("Received %d attr desc", this._attrDescCount);
    			break;
    		case 0x1011:
                this._entityCount = inputStream.readInt();
                for (i = 0; i < this._entityCount; i++) {
                    this._addedIdList.push(inputStream.readLong());  // temp id
                    this._addedIdList.push(inputStream.readLong());  // real id
                }
                console.log("Received %d entity", this._entityCount);
    			break;
            case 0x1012:
                console.log("Received update results");
                break;
            case 0x1013:
                this._entityCount = inputStream.readInt();
                for (i = 0; i < this._entityCount; i++) {
                    this._removedIdList.push(inputStream.readLong());  //id
                }
                console.log("Received %d delete results", this._entityCount);
                break;
            default:
                break;
    		}
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
    return this._entityCount;
}; 

CommitTransactionResponse.prototype.getAttrDescIdList = function () {
    return this._attrDescIdList;
};

CommitTransactionResponse.prototype.getAddedIdList = function () {
    return this._addedIdList;
};

exports.CommitTransactionResponse = CommitTransactionResponse;
