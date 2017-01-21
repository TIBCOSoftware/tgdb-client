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
    CallableRequest  = require('./CallableRequest').CallableRequest,
    VerbId                  = require('./VerbId').VerbId;

function GetEntityRequest (callback) {
	GetEntityRequest.super_.call(this, callback);
	
	this._key = null;

    //0 - get, 1 - getbyid, 2 - get multiples, 10 - continue, 20 - close
	this._getCommand = 0;
	this._fetchSize = 1000;
	this._batchSize = 50;
	this._traversalDepth = 3;
	this._edgeFetchSize = 0; // zero means no limitation
	this._resultSetId = 0;
}

util.inherits(GetEntityRequest, CallableRequest);

GetEntityRequest.prototype.getVerbId = function () {
    return VerbId.GET_ENTITY_REQUEST;
};

GetEntityRequest.prototype.isUpdateable = function () {
    return true;
};

GetEntityRequest.prototype.setKey = function (key) {
    this._key = key;
};

GetEntityRequest.prototype.setBatchSize = function (size) {
    if (size <= 10 || size >= 32767) {
    	this._batchSize = 50;
    } else {
    	this._batchSize = size;
    }
};
    
GetEntityRequest.prototype.getBatchSize = function () {
    return this._batchSize;
};
    
GetEntityRequest.prototype.setFetchSize = function (size) {
    if (size < 0) {
    	this._fetchSize = 1000;
    } else {
    	this._fetchSize = size;
    }
};
    
GetEntityRequest.prototype.getFetchSize = function () {
    return this._fetchSize;
};

GetEntityRequest.prototype.setEdgeFetchSize = function (size) {
	if (size < 0 || size > 32767) {
		this._edgeFetchSize = 1000;
	} else {
		this._edgeFetchSize = size;
	}
};

GetEntityRequest.prototype.getEdgeFetchSize = function () {
	return this._edgeFetchSize;
};

GetEntityRequest.prototype.setTraversalDepth = function (depth) {
    if (depth < 1 || depth > 1000) {
    	this._traversalDepth = 3;
    } else {
    	this._traversalDepth = depth;
    }
};
    
GetEntityRequest.prototype.getTraversalDepth = function () {
    return this._traversalDepth;
};

GetEntityRequest.prototype.setResultSetId = function (id) {
	this._resultSetId = id;
};

GetEntityRequest.prototype.getResultSetId = function () {
    return this._resultSetId;
};

GetEntityRequest.prototype.setCommand = function (cmd) {
    this._getCommand = cmd;
};

GetEntityRequest.prototype.getCommand = function () {
    return this._getCommand;
};

GetEntityRequest.prototype.writePayload = function (outputStream) {
	outputStream.writeShort(this._getCommand);
	outputStream.writeInt(this._resultId);
    if (this._getCommand === 0 || 
    	this._getCommand === 1 || 
    	this._getCommand === 2) {
    	outputStream.writeInt(this._fetchSize);
    	outputStream.writeShort(this._batchSize);
    	outputStream.writeShort(this._traversalDepth);
    	outputStream.writeShort(this._edgeFetchSize);
	    this._key.writeExternal(outputStream);
    }
};

GetEntityRequest.prototype.readPayload = function (inputStream) {
		// TODO Auto-generated method stub
		
};

exports.GetEntityRequest = GetEntityRequest;