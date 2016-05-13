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

var globolRequestId = 0;

//Class Definition
function QueryRequest() {
	QueryRequest.super_.call(this);
    this._queryExpr   = null;
    this._command     = null;
    this._queryObject = null;
}

util.inherits(QueryRequest, AbstractProtocolMessage);

QueryRequest.prototype.getVerbId = function() {
    return VerbId.QUERY_REQUEST;
};

QueryRequest.prototype.setQuery = function(queryExpr) {
    this._queryExpr = queryExpr;
};

QueryRequest.prototype.setQueryObject = function(queryObject) {
    this._queryObject = queryObject;
};

QueryRequest.prototype.setCommand = function(command) {
    this._command = command;
};

QueryRequest.prototype.getCommand = function() {
    return this._command;
};

QueryRequest.prototype.getVerbId = function() {
    return VerbId.QUERY_REQUEST;
};

QueryRequest.prototype.writePayload = function(outputStream) {
    var startPos = outputStream.getPosition();
    outputStream.writeInt(1);
    outputStream.writeInt(1);
    console.log("****Entering query request writePayload at output buffer position at : %d", startPos);

    outputStream.writeInt(this._command);
    
    if(null==this._queryExpr|| (typeof this._queryExpr == 'undefined')) {
    	throw new Error('Not a valid query expression : ' + this._queryExpr);
    }
    
    // createQuery
    var strQuertExpression = String(this._queryExpr);
    if (strQuertExpression) {
    	outputStream.writeUTF(this._queryExpr);
    }
};

QueryRequest.prototype.readPayload = function() {

};

QueryRequest.prototype.isUpdateable = function() {
    return false;
};

exports.QueryRequest = QueryRequest;
