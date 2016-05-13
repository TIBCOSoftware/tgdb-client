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

var TGAbstractEntity = require('./TGAbstractEntity').TGAbstractEntity,
    util             = require('util');

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
function TGEdge(graphMetadata, fromNode, toNode, edgeType, directionType) {
	TGEdge.super_.call(this, graphMetadata);
    this._fromNode      = fromNode;
    this._toNode        = toNode;
    this._edgeType      = edgeType;
    this._directionType = directionType;
}

TGEdge.DirectionType = {
    UNDIRECTED : 0,
    DIRECTED : 1,
    BIDIRECTIONAL : 2
};

util.inherits(TGEdge, TGAbstractEntity);

TGEdge.prototype.getEntityKind = function() {
    return TGAbstractEntity.EntityKind.EDGE;
};

/**
 * Return the nodes connected by the edge.
 */
TGEdge.prototype.getVertices = function() {
    var nodes = [];
    nodes.push(this._fromNode, this._toNode);
};

/**
 * Get direction for this edge
 */
TGEdge.prototype.getDirection = function() {
    return (this._edgeType != null) ? this._edgeType.getDirectionType() : this._directionType;
};

TGEdge.prototype.writeExternal = function (outputStream) {
	console.log("**** Enering TGEdge.writeExternal at output buffer position at : %d", outputStream.getPosition());
	var startPos = outputStream.getPosition();
	outputStream.writeInt(0);
	//write attributes from the based class
    this.writeAttributes(outputStream);
    //write the edges ids
    if (this.isNew()) {
    	outputStream.writeLong(this._fromNode.getVirtualId());
    	outputStream.writeLong(this._toNode.getVirtualId());
    }
    var currPos = outputStream.getPosition();
    var length = currPos - startPos;
    outputStream.writeIntAt(startPos, length);
	console.log("**** Leaving TGEdge.writeExternal at output buffer position at : %d", outputStream.getPosition());
}

//TGEdge.prototype.readExternal = function (inputStream) { 
//};

module.exports = TGEdge;

