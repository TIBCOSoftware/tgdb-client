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
    VerbId                  = require('./VerbId').VerbId,
    Response                = require('./Response').Response,
	TGEdge                  = require('../../model/TGEdge'),
	TGNodeType              = require('../../model/TGNodeType'),
	TGEdgeType              = require('../../model/TGEdgeType'),
	TGSystemObject          = require('../../model/TGSystemObject').TGSystemObject,
	TGAttributeType         = require('../../model/TGAttributeType').TGAttributeType,
	TGAttributeDescriptor   = require('../../model/TGAttributeDescriptor').TGAttributeDescriptor,
    TGException             = require('../../exception/TGException').TGException,
    TGLogManager            = require('../../log/TGLogManager'),
    TGLogLevel              = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function MetadataResponse () {
	MetadataResponse.super_.call(this);
    this._attrDescList = [];
    this._nodeTypeList = []; 
    this._edgeTypeList = []; 
}

util.inherits(MetadataResponse, Response);

MetadataResponse.prototype.writePayload = function (outputStream) {
};

MetadataResponse.prototype.readPayload = function (inputStream) {
    	if (inputStream.available() === 0) {
    		logger.logDebugWire( 
    				"Entering metadata response has no data");
    		return;
    	}
    	var count = inputStream.readInt();
    	while (count > 0) {
    		var i = 0;
    		var sysType = inputStream.readByte();
    		var typeCount = inputStream.readInt();
    		if (sysType === TGSystemObject.TGSystemType.AttributeDescriptor.value) {
    			for (i=0; i<typeCount; i++) {
    				var attrDesc = new TGAttributeDescriptor("temp", TGAttributeType.String);
    				attrDesc.readExternal(inputStream);
    				this._attrDescList.push(attrDesc);
    			}
    		} else if (sysType === TGSystemObject.TGSystemType.NodeType.value) {
    			for (i=0; i<typeCount; i++) {
    				var nodeType = new TGNodeType("temp", null);
    				nodeType.readExternal(inputStream);
    				this._nodeTypeList.push(nodeType);
    			}
    		} else if (sysType === TGSystemObject.TGSystemType.EdgeType.value) {
    			for (i=0; i<typeCount; i++) {
    				var edgeType = new TGEdgeType("temp", TGEdge.DirectionType.BiDirectional, null);
    				edgeType.readExternal(inputStream);
    				this._edgeTypeList.push(edgeType);
    			}
    		} else {
    			throw new TGException("Invalid meta data type received %d\n", sysType);
    		}
    		count -= typeCount;
    	}
};

MetadataResponse.prototype.isUpdateable = function () {
    return false;
};

MetadataResponse.prototype.getVerbId = function () {
    return VerbId.METADATA_RESPONSE;
};

MetadataResponse.prototype.getAttrDescList = function () {
    return this._attrDescList;
};

MetadataResponse.prototype.getNodeTypeList = function () {
    return this._nodeTypeList;
};

MetadataResponse.prototype.getEdgeTypeList = function () {
    return this._edgeTypeList;
};

exports.MetadataResponse = MetadataResponse;