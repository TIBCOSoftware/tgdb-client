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

var util             = require('util'),
    TGAbstractEntity = require('./TGAbstractEntity').TGAbstractEntity,
	TGEntityKind     = require('./TGEntityKind').TGEntityKind,
    TGNodeType       = require('./TGNodeType'),
    TGEntityId       = require('./TGEntityId'),
    TGEdge           = require('./TGEdge').TGEdge,
    TGException      = require('../exception/TGException').TGException,
    TGLogManager     = require('../log/TGLogManager'),
    TGLogLevel       = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

//Class definition
function TGNode(graphMetadata, nodeType) {
	TGNode.super_.call(this, graphMetadata, nodeType);
    this._edges         = [];
    //if(nodeType) {
    //	this.setEntityType(nodeType);
    //}
}

util.inherits(TGNode, TGAbstractEntity);

TGNode.prototype.getEntityKind = function() {
    return TGEntityKind.NODE;
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
        if (!edgeDirection || (edge.getDirection() === edgeDirection)) {
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
	//console.log('Add edge : ' + edge);
    this._edges.push(edge);
    return edge;
};


TGNode.prototype.writeExternal = function (outputStream) {
	logger.logDebugWire(
			"****Entering NodeImpl.writeExternal at output buffer position at : %d", 
			outputStream.getPosition());
	var startPos = outputStream.getPosition();
	outputStream.writeInt(0);
	//write attributes from the based class
    this.writeAttributes(outputStream);
    //write the edges ids
    var newCount = 0;
    
    for(var i = 0; i < this._edges.length; i++) {
    	if(this._edges[i].isNew()) {
    		newCount++;
    	}
    }
    outputStream.writeInt(newCount); // only include the new edges
    
    for(i = 0; i < this._edges.length; i++) {
    	if(this._edges[i].isNew()) {
    		outputStream.writeLongAsBytes(this._edges[i].getId().getBytes());
    	}
    }
    var currPos = outputStream.getPosition();
    var length = currPos - startPos;
    outputStream.writeIntAt(startPos, length);
    logger.logDebugWire(
    		"****Leaving NodeImpl.writeExternal at output buffer position at : %d", 
    		outputStream.getPosition());
};

TGNode.prototype.readExternal = function (inputStream, gof) {
	logger.logDebugWire(
			"****Entering NodeImpl.readExternal at output buffer position at : %d", 
			inputStream.getPosition());
	//FIXME: Need to validate length
	var buflen = inputStream.readInt();
	this.readAttributes(inputStream);

    var edgeCount = inputStream.readInt();
    for (var i=0; i<edgeCount; i++){
    	var edge = null;
        var idBytes = inputStream.readLongAsBytes();
        var id = TGEntityId.bytesToString(idBytes);
    	var refMap = inputStream.getReferenceMap();
    	if (refMap) {
    		edge = refMap[id];
    	}
    	
    	if (!edge) {
    		edge = gof.createEdge();//new TGEdge(this._graphMetadata);
    		edge.setInitialized(false);
    		edge.getId().setBytes(idBytes);
    		if(refMap) {
       			refMap[id] = edge;   			
    		}
     	}
    	// Add this edge to my collection
    	this._edges.push(edge);
    }
	this._isInitialized = true;
	logger.logDebugWire(
			"****Leaving NodeImpl.readExternal at output buffer position at : %d", 
			inputStream.getPosition());
};

exports.TGNode = TGNode;
