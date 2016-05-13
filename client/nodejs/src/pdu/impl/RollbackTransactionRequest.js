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

var util                      = require('util'),
    AuthenticatedMessage      = require('./AuthenticatedMessage').AuthenticatedMessage,
    AbstractProtocolMessage   = require('../AbstractProtocolMessage').AbstractProtocolMessage,
    VerbId                    = require('./VerbId').VerbId;

//Class Definition
function RollbackTransactionRequest() {
    AuthenticatedMessage.call();
}

util.inherits(RollbackTransactionRequest, AuthenticatedMessage);

RollbackTransactionRequest.prototype.getVerbId = function() {
    return VerbId.ROLLBCK_TRANS_REQUEST;
};

RollbackTransactionRequest.prototype.writePayload = function(dynamicBuffer) {

};

RollbackTransactionRequest.prototype.readPayload = function() {

};

RollbackTransactionRequest.prototype.isUpdateable = function() {
    return false;
};

exports.RollbackTransactionRequest = RollbackTransactionRequest;