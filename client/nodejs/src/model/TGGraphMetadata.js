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

var TGNodeType = require('./TGNodeType'),
    TGEdgeType = require('./TGEdgeType'),
	TGAttributeDescriptor = require('./TGAttributeDescriptor').TGAttributeDescriptor;

//Class definition
function TGGraphMetadata() {
    this._attrDescriptorMap = {};
    this._nodeTypeMap       = {};
    this._edgeTypeMap       = {};
}

/**
 * Return a set of Node Type defined in the System.
 */
TGGraphMetadata.prototype.getNodeTypes = function() {
    var nodeTypeMap = this._nodeTypeMap;
    return Object.keys(nodeTypeMap).map(function(key){
        return nodeTypeMap[key];
    });
};

/**
 * Get Node Type by name.
 * @param typeName
 */
TGGraphMetadata.prototype.getNodeType = function(typeName) {
    var nodeTypeMap = this._nodeTypeMap;
    return nodeTypeMap[typeName];
};

/**
 * Return a set of known edge Types.
 */
TGGraphMetadata.prototype.getEdgeTypes = function() {
    var edgeTypeMap = this._edgeTypeMap;
    return Object.keys(edgeTypeMap).map(function(key){
        return edgeTypeMap[key];
    });
};

/**
 * Get Edge Type by name.
 * @param typeName
 */
TGGraphMetadata.prototype.getEdgeType = function(typeName) {
    var edgeTypeMap = this._edgeTypeMap;
    return edgeTypeMap[typeName];
};

/**
 * Return a set of known attribute descriptors.
 */
TGGraphMetadata.prototype.getAttributeDescriptors = function() {
    return Object.keys(this._attrDescriptorMap).map(function(key){
        return attrDescriptorMap[key];
    });
};

TGGraphMetadata.prototype.getNewAttributeDescriptors = function () {
	var attrDescriptorMap = this._attrDescriptorMap;
	var newAttributeDescriptors = Object.keys(this._attrDescriptorMap).filter(function(key) {
	    return attrDescriptorMap[key].getAttributeId() < 0;
	}).map(function(key) {
	    return attrDescriptorMap[key];
	});
	
	for (index in newAttributeDescriptors) {
		console.log("Confirm ne attribute descriptor :  index = " + index + ", id = " + newAttributeDescriptors[index].getAttributeId());
	}
	
	return newAttributeDescriptors;
}

/**
 * Get the Attribute Descriptor.
 * @param attributeName
 */
TGGraphMetadata.prototype.getAttributeDescriptor = function(attributeName) {
	return this._attrDescriptorMap[attributeName];
};

/**
 * Create an Attribute Descriptor.
 * @param attributeName
 * @param attributeType
 * @param minOccurs
 * @param maxOccurs
 */
TGGraphMetadata.prototype.createAttributeDescriptor = function(attributeName, attributeType, isArray) {
    var attributeDescriptor = new TGAttributeDescriptor(attributeName, attributeType, isArray);
    this._attrDescriptorMap[attributeName] = attributeDescriptor;	
	
    return attributeDescriptor;
};

/**
 * Create Node type for a given name, and derive it from the following type.
 * @param nodeName
 * @param parentNodeType - #TGNodeType
 */
TGGraphMetadata.prototype.createNode = function(nodeName, parentNodeType) {
    var nodeType = new TGNodeType(this, nodeName, parentNodeType);
    this._nodeTypeMap[nodeName] = nodeType;
};

/**
 * Create Edge type for a given name, and derive it from the following type.
 * @param edgeName
 * @param parentEdgeType - #TGEdgeType
 */
TGGraphMetadata.prototype.createEdge = function(edgeName, parentEdgeType) {
    var edgeType = new TGEdgeType(this, edgeName, parentEdgeType);
    this._edgeTypeMap[edgeName] = parentEdgeType;
};

exports.TGGraphMetadata = TGGraphMetadata;
