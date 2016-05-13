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
    TGEdge           = require('./TGEdge'),
    util             = require('util');

//Class definition
function TGNode(graphMetadata) {
	TGNode.super_.call(this, graphMetadata);
    this._edges         = [];
}

util.inherits(TGNode, TGAbstractEntity);

TGNode.prototype.getEntityKind = function() {
    return TGAbstractEntity.EntityKind.NODE;
};

/**
 * Get list of graphs to which node belongs.
 */
TGNode.prototype.getGraphs = function() {
    return null;
};

/**
 * Get all edges incoming/outgoing.
 */
TGNode.prototype.getEdges = function() {
    return this._edges;
};

/**
 * Get all edges based on direction.
 * @param edgeDirection
 */
TGNode.prototype.getEdges = function(edgeDirection) {
    var allEdges = this._edges;
    var edges    = [];

    for (var loop = 0; loop < allEdges.length; loop++) {
        var edge = allEdges[loop];
        if (edge.getDirection() == edgeDirection) {
            edges.push(edge);
        }
    }
    return edges;
};

/**
 * Add edge from this node to other node.
 * @param toNode - To node #TGNode
 * @param edgeType - Type of edge #TGEdgeType
 * @param direction
 */
TGNode.prototype.addEdge = function(edge) {
    this._edges.push(edge);
    return edge;
};


TGNode.prototype.writeExternal = function (outputStream) {
	console.log("****Entering NodeImpl.writeExternal at output buffer position at : %d", outputStream.getPosition());
	var startPos = outputStream.getPosition();
	outputStream.writeInt(0);
	//write attributes from the based class
    this.writeAttributes(outputStream);
    //write the edges ids
    var newCount = 0;
    
    console.log('?????????????????????' + this._edges.length);
    
    for(var i = 0; i < this._edges.length; i++) {
    	if(this._edges[i].isNew()) {
    		newCount++;
    	}
    }
    outputStream.writeInt(newCount); // only include the new edges
    
    for(i = 0; i < this._edges.length; i++) {
    	if(this._edges[i].isNew()) {
    		outputStream.writeLong(this._edges[i].getVirtualId());
    	}
    }
    var currPos = outputStream.getPosition();
    var length = currPos - startPos;
    outputStream.writeIntAt(startPos, length);
	console.log("****Leaving NodeImpl.writeExternal at output buffer position at : %d", outputStream.getPosition());
};

exports.TGNode = TGNode;
