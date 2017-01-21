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

var util                = require('util'),
    TGAbstractEntity    = require('./TGAbstractEntity').TGAbstractEntity,
	TGEntityKind        = require('./TGEntityKind').TGEntityKind,
	TGEdgeDirectionType = require('./TGEdgeDirectionType').TGEdgeDirectionType,
    TGEntityId          = require('./TGEntityId'),
	TGNode              = require('./TGNode').TGNode,
    TGException         = require('../exception/TGException').TGException,
    TGLogManager        = require('../log/TGLogManager'),
    TGLogLevel          = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

/**
 *
 * @param graphMetadata - #TGGraphMetadata
 * @param fromNode - #TGNode
 * @param toNode - #TGNode
 * @param edgeType - #TGEdgeType
 * @param directionType - #DirectionType
 * @constructor
 */
//Class definition
function TGEdge(graphMetadata, fromNode, toNode, directionType, edgeType) {
	TGEdge.super_.call(this, graphMetadata, edgeType);
    this._fromNode      = fromNode;
    this._toNode        = toNode;
    this._directionType = edgeType?edgeType.getDirectionType():directionType;
}

util.inherits(TGEdge, TGAbstractEntity);

TGEdge.prototype.getEntityKind = function() {
    return TGEntityKind.EDGE;
};

/**
 * Return the nodes connected by the edge.
 */
TGEdge.prototype.getVertices = function() {
    var nodes = [];
    nodes.push(this._fromNode, this._toNode);
    
    return nodes;
};

/**
 * Get direction for this edge
 */
TGEdge.prototype.getDirection = function() {
	var edgeType = this.getEntityType();
    return (edgeType !== null) ? edgeType.getDirectionType() : this._directionType;
};

TGEdge.prototype.writeExternal = function (outputStream) {
	logger.logDebugWire(
			"**** Enering TGEdge.writeExternal at output buffer position at : %d", 
			outputStream.getPosition());
	var startPos = outputStream.getPosition();
	outputStream.writeInt(0);
	//write attributes from the based class
    this.writeAttributes(outputStream);
    //write the edges ids
    if (this.isNew()) {
    	outputStream.writeByte(this._directionType.ordinal);
    	outputStream.writeLongAsBytes(this._fromNode.getId().getBytes());
    	outputStream.writeLongAsBytes(this._toNode.getId().getBytes());
    } else {
    	//FIXME:  Not sending the direction, is it ok?
    	outputStream.writeLongAsBytes(this._fromNode.getId().getBytes());
    	outputStream.writeLongAsBytes(this._toNode.getId().getBytes());
    }
    var currPos = outputStream.getPosition();
    var length = currPos - startPos;
    outputStream.writeIntAt(startPos, length);
    logger.logDebugWire(
    		"**** Leaving TGEdge.writeExternal at output buffer position at : %d", 
    		outputStream.getPosition());
};

TGEdge.prototype.readExternal = function (inputStream, gof) { 
	//FIXME: Need to validate length
	var buflen = inputStream.readInt();
    this.readAttributes(inputStream);
    var dir = inputStream.readByte();
    switch (dir) {
    	case TGEdgeDirectionType.UNDIRECTED.ordinal :
    		this._directionType = TGEdgeDirectionType.UNDIRECTED;
    		break;
    	case TGEdgeDirectionType.DIRECTED.ordinal :
    		this._directionType = TGEdgeDirectionType.DIRECTED;
    		break;
    	case TGEdgeDirectionType.BIDIRECTIONAL.ordinal :
    		this._directionType = TGEdgeDirectionType.BIDIRECTIONAL;
    		break;
    }
    var fromNode = null;
    var toNode = null;
    var idBytes = inputStream.readLongAsBytes();
    var id = TGEntityId.bytesToString(idBytes);
    
    var refMap = inputStream.getReferenceMap();
    if (refMap) {
    	fromNode = refMap[id];
    }

    if (!fromNode) {
    	fromNode = gof.createNode();//new TGNode(this._graphMetadata);
    	fromNode.getId().setBytes(idBytes);
    	fromNode.setInitialized(false);
    	refMap[id] = fromNode;
    }
    this._fromNode = fromNode;

    idBytes = inputStream.readLongAsBytes();
    id = TGEntityId.bytesToString(idBytes);
    
    if (refMap) {
    	toNode = refMap[id];
    }
    
    if (!toNode) {
    	toNode = gof.createNode();//new TGNode(this._graphMetadata);
    	toNode.getId().setBytes(idBytes);
    	toNode.setInitialized(false);
    	refMap[id] = toNode;
    }
    this._toNode = toNode;
	this._isInitialized = true;
};

exports.TGEdge = TGEdge;

