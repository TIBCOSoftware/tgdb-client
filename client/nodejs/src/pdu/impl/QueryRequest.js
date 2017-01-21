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

var util            = require('util'),
    VerbId          = require('./VerbId').VerbId,
    CallableRequest = require('./CallableRequest').CallableRequest,
    TGException     = require('../../exception/TGException').TGException,
    TGLogManager    = require('../../log/TGLogManager'),
    TGLogLevel      = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

var globolRequestId = 0;

//Class Definition
function QueryRequest(callback) {
	QueryRequest.super_.call(this, callback);
    this._queryExpr   = null;
    this._queryHashId = null;
    this._command     = null;
    this._queryObject = null;
    this._fetchSize = 1000;
    this._batchSize = 50;
    this._traversalDepth = 3;
    this._edgeFetchSize = 0; // zero means no limitation}
}

util.inherits(QueryRequest, CallableRequest);

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

QueryRequest.prototype.setBatchSize = function(size) {
	if (size < 10 || size > 32767) {
		this._batchSize = 50;
	} else {
		this._batchSize = size;
	}
};

QueryRequest.prototype.getBatchSize = function() {
	return this._batchSize;
};

QueryRequest.prototype.setFetchSize = function(size) {
	if (size < 0) {
		this._fetchSize = 1000;
	} else {
		this._fetchSize = size;
	}
};

QueryRequest.prototype.getFetchSize = function() {
	return this._fetchSize;
};

QueryRequest.prototype.setEdgeFetchSize = function(size) {
	if (size < 0 || size > 32767) {
		this._edgeFetchSize = 1000;
	} else {
		this._edgeFetchSize = size;
	}
};

QueryRequest.prototype.getEdgeFetchSize = function() {
	return this._edgeFetchSize;
};

QueryRequest.prototype.setTraversalDepth = function(depth) {
	if (depth < 1 || depth > 1000) {
		this._traversalDepth = 3;
	} else {
		this._traversalDepth = depth;
	}
};

QueryRequest.prototype.getTraversalDepth = function() {
	return this._traversalDepth;
};

QueryRequest.prototype.writePayload = function(outputStream) {
	logger.logDebugWire(
			"****Entering query request writePayload at output buffer position at : %d", 
			outputStream.getPosition());
    var startPos = outputStream.getPosition();
    outputStream.writeInt(1); // datalength
    outputStream.writeInt(1); //checksum

    outputStream.writeInt(this._command);
    outputStream.writeInt(this._fetchSize);
    outputStream.writeShort(this._batchSize);
    outputStream.writeShort(this._traversalDepth);
    outputStream.writeShort(this._edgeFetchSize);
    
    if(null===this._queryExpr|| (typeof this._queryExpr === 'undefined')) {
    	throw new TGException('Not a valid query expression : ' + this._queryExpr);
    }
        
    // CREATE, EXECUTE.
    if (this._command === 1 || this._command === 2) {
        var strQuertExpression = String(this._queryExpr);
        if (strQuertExpression) {
        	outputStream.writeUTF(this._queryExpr);
        }
    }
    // EXECUTEID, CLOSE
    else if (this._command === 3 || this._command === 4){
    	outputStream.writeLong(this._queryHashId);
    }
    logger.logDebugWire(
    		"****Leaving query request writePayload at output buffer position at : %d", 
    		outputStream.getPosition());
};

QueryRequest.prototype.readPayload = function() {

};

QueryRequest.prototype.isUpdateable = function() {
    return false;
};

exports.QueryRequest = QueryRequest;