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

var util         = require('util'),
    VerbId       = require('./VerbId').VerbId,
    Response     = require('./Response').Response,
    TGException  = require('../../exception/TGException').TGException,
    TGLogManager = require('../../log/TGLogManager'),
    TGLogLevel   = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function QueryResponse () {
	QueryResponse.super_.call(this);
	this._result = null;
	this._queryHashId = null;
	
	this._entityStream = null;
	this._hasResult = false;
	this._totalCount = 0;
	this._resultCount = 0;
}

util.inherits(QueryResponse, Response);

QueryResponse.prototype.getVerbId = function() {
    return VerbId.QUERY_RESPONSE;
};

QueryResponse.prototype.writePayload = function (outputStream) {
};

QueryResponse.prototype.readPayload = function (inputStream) {
	logger.logDebugWire(
			"****Entering query response readPayload at output buffer position at : %d", 
			inputStream.getPosition());
	if (inputStream.available() === 0) {
		return;
	}
	this._entityStream = inputStream;
	inputStream.readInt(); // buf length
	inputStream.readInt(); // checksum
	this._result = inputStream.readInt(); // query result
	this._queryHashId = inputStream.readLong(); // query hash id

    this._resultCount = inputStream.readInt();
    this._totalCount = inputStream.readInt();
	if (this._resultCount > 0) {
		this._hasResult = true;
	}
	logger.logDebugWire(
			"****Leaving query response readPayload at output buffer position at : %d", 
			inputStream.getPosition());
};

QueryResponse.prototype.isUpdateable = function () {
	return false;
};

QueryResponse.prototype.getVerbId = function () {
	return VerbId.QUERY_RESPONSE;
};
    
QueryResponse.prototype.getResult = function () {
	return this._result;
};

QueryResponse.prototype.getQueryHashId = function () {
	return this._queryHashId;
};

QueryResponse.prototype.getEntityStream = function () {
	return this._entityStream;
};

QueryResponse.prototype.hasResult = function () {
	return this._hasResult;
};

QueryResponse.prototype.getTotalCount = function () {
    return this._totalCount;
};

QueryResponse.prototype.getResultCount = function () {
    return this._resultCount;
};

exports.QueryResponse = QueryResponse;
