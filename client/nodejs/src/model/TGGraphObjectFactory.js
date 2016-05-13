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
var TGNode = require('./TGNode').TGNode;
var TGEdge = require('./TGEdge');

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

TGGraphObjectFactory.prototype.getGraphMetaData = function () {
    return this._graphMetadata;

}

TGGraphObjectFactory.prototype.getEntityManager = function () {
    return this._entityMgr;
}

/**
 * Create a root node in the Graph DB server.
 * @param nodeType - #TGNodeType
 */
TGGraphObjectFactory.prototype.createNode = function() {
    var node = new TGNode(this._graphMetadata);
    this._entityMgr.nodeAdded(null, node);
    //this._entityMgr.entityCreated(node);
    return node;
};

/**
 * Create an edge using from and to nodes.
 * @param fromNode
 * @param toNode
 * @param edgeType - #TGEdgeType
 * @param directionType - #TGEdge#DirectionType
 */
TGGraphObjectFactory.prototype.createEdge = function(fromNode, toNode, edgeType, directionType) {
    directionType = directionType || TGEdge.DirectionType.BIDIRECTIONAL;
    var edge = new TGEdge(this._graphMetadata, fromNode, toNode, edgeType);
    fromNode.addEdge(edge);
    toNode.addEdge(edge);
    this._entityMgr.entityCreated(edge);
    return edge;
};

/**
 * Create a new graph with given name.
 * @param name
 */
TGGraphObjectFactory.prototype.createGraph = function(name) {

};

exports.TGGraphObjectFactory = TGGraphObjectFactory;

