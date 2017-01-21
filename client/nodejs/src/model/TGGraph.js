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

var TGNode = require('./TGNode').TGNode,
    util   = require('util');

//Class definition
function TGGraph(graphMetadata, name) {
    TGGraph.caTGNode.super_.call(this, graphMetadata);
    this._graphMetadata = graphMetadata;
    this._name          = name;
}

util.inherits(TGGraph, TGNode);


/**
 * Get the Graphs to which this Node is a member of.
 */
TGGraph.prototype.getGraphs = function() {

};

/**
 * Get the List of Edges.
 */
TGGraph.prototype.getEdges = function() {

};

/**
 * Add a new node to the graph.
 */
TGGraph.prototype.addNode = function(node) {
    //TODO add new node to the graph
};

/**
 * Add a new node to the graph.
 */
TGGraph.prototype.removeNode = function(node) {
    //TODO remove node from graph
};

module.exports = TGGraph;