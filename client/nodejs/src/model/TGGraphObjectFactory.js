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

var TGGraphMetadata = require('./TGGraphMetadata').TGGraphMetadata;
var TGAbstractEntity = require('./TGAbstractEntity').TGAbstractEntity;
var TGEntityKind        = require('./TGEntityKind').TGEntityKind;
var TGCompositeKey = require('./TGCompositeKey').TGCompositeKey;
var TGEdgeDirectionType = require('./TGEdgeDirectionType').TGEdgeDirectionType;
//var TGEntityId = require('./TGEntityId').TGEntityId;
var TGNode = require('./TGNode').TGNode;
var TGEdge = require('./TGEdge').TGEdge;
var TGGraph = require('./TGGraph');

//Class definition
/**
 *
 * @param connection - #ConnectionImpl
 * @constructor
 */
function TGGraphObjectFactory(entityMgr) {
    this._graphMetadata = new TGGraphMetadata();
    this._entityMgr    = entityMgr;
}

TGGraphObjectFactory.prototype.isNode = function (entity) {
    return (entity instanceof TGNode);
};

TGGraphObjectFactory.prototype.isEdge = function (entity) {
    return (entity instanceof TGEdge);
}

TGGraphObjectFactory.prototype.getGraphMetaData = function () {
    return this._graphMetadata;
};

TGGraphObjectFactory.prototype.getEntityManager = function () {
    return this._entityMgr;
};

/**
 * Create Entity based on kind. This is used for deserialization purpose only. Does not notify  the listener.
 * @param kind
 * @return
 */
TGGraphObjectFactory.prototype.createEntity = function(kind) {
    switch (kind) {
        case TGEntityKind.NODE:
            return new TGNode(this._graphMetadata);
        case TGEntityKind.EDGE:
            return new TGEdge(this._graphMetadata);
        case TGEntityKind.GRAPH:
            return new TGGraph(this._graphMetadata);
    }
    return null;
};

TGGraphObjectFactory.prototype.createNode = function(nodeType) {
    return new TGNode(this._graphMetadata, nodeType);
};
/**
 * Create an undirected edge using from and to nodes.
 * @param fromNode
 * @param toNode
 * @param edgeType - #TGEdgeType
 * @param directionType - #TGEdge#DirectionType
 */
TGGraphObjectFactory.prototype.createUndirectedEdge = function(fromNode, toNode, edgeType) {
    return this.createEdge(fromNode, toNode, edgeType, TGEdge.UNDIRECTED);
};

/**
 * Create an directed edge using from and to nodes.
 * @param fromNode
 * @param toNode
 * @param edgeType - #TGEdgeType
 * @param directionType - #TGEdge#DirectionType
 */
TGGraphObjectFactory.prototype.createDirectedEdge = function(fromNode, toNode, edgeType) {
    return this.createEdge(fromNode, toNode, edgeType, TGEdge.DIRECTED);
};

/**
 * Create an bidirectional edge using from and to nodes.
 * @param fromNode
 * @param toNode
 * @param edgeType - #TGEdgeType
 * @param directionType - #TGEdge#DirectionType
 */
TGGraphObjectFactory.prototype.createBidirectionalEdge = function(fromNode, toNode, edgeType) {
    return this.createEdge(fromNode, toNode, TGEdgeDirectionType.BIDIRECTIONAL, edgeType);
};

/**
 * Create an edge using from and to nodes.
 * @param fromNode
 * @param toNode
 * @param edgeType - #TGEdgeType
 * @param directionType - #TGEdge#DirectionType
 */
TGGraphObjectFactory.prototype.createEdge = function(fromNode, toNode, directionType, edgeType) {
    directionType = directionType || TGEdgeDirectionType.BIDIRECTIONAL;
    var edge = new TGEdge(this._graphMetadata, fromNode, toNode, directionType, edgeType);
    if(fromNode)
    	fromNode.addEdge(edge);
    if(toNode)
    	toNode.addEdge(edge);
 
    return edge;
};

/**
 * Create a new graph with given name.
 * @param name
 */
TGGraphObjectFactory.prototype.createGraph = function(name) {
    return new TGGraph(this._graphMetadata, name);
};

TGGraphObjectFactory.prototype.createEntityId = function(buf) {
    return new ByteArrayEntityId(buf);
};

TGGraphObjectFactory.prototype.createCompositeKey = function(nodeTypeName) {
    var tgKey =  new TGCompositeKey(this._graphMetadata, nodeTypeName);
    //console.log('----> new tgKey = ' + tgKey);
    return tgKey;
};

exports.TGGraphObjectFactory = TGGraphObjectFactory;

