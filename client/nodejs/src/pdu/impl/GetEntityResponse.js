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

function GetEntityResponse () {
	GetEntityResponse.super_.call(this);
    
    this._entityStream = null;
    this._hasResult = false;
    this._resultId = 0;
    this._totalCount = 0;
    this._resultCount = 0;
}

util.inherits(GetEntityResponse, Response);

GetEntityResponse.prototype.getVerbId = function () {
    return VerbId.GET_ENTITY_RESPONSE;
};

GetEntityResponse.prototype.isUpdateable = function () {
    return false;
};

GetEntityResponse.prototype.getEntityStream = function () {
    return this._entityStream;
};

GetEntityResponse.prototype.hasResult = function () {
    return this._hasResult;
};

GetEntityResponse.prototype.getResultId = function () {
    return this._resultId;
};

GetEntityResponse.prototype.writePayload = function (outputStream) {
};

GetEntityResponse.prototype.readPayload = function (inputStream) {
	logger.logDebugWire( 
			"Entering get entity response readPayload");
    if (inputStream.available() === 0) {
    	logger.logDebugWire( 
    			"Entering metadata response has no data");
    	return;
    }
    this._entityStream = inputStream;
    this._resultId = inputStream.readInt();
    var pos = inputStream.getPosition();
    this._totalCount = inputStream.readInt();
    if (this._totalCount > 0) {
    	this._hasResult = true;
    }
    inputStream.setPosition(pos);
    logger.logDebugWire( 
    		"Received %d debug entities", this._totalCount);
};


GetEntityResponse.prototype.getTotalCount = function () {
    return this._totalCount;
};

GetEntityResponse.prototype.getResultCount = function () {
    return this._resultCount;
};

exports.GetEntityResponse = GetEntityResponse;    

