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
    TGCompositeKey = require('./TGCompositeKey'),
	TGAttributeDescriptor = require('./TGAttributeDescriptor').TGAttributeDescriptor;

//Class definition
function TGGraphMetadata() {
    this._attrDescriptorMap = {};
    this._nodeTypeMap       = {};
    this._edgeTypeMap       = {};
	this._descriptorMapById = {};
	this._nodeTypeMapById   = {};
	this._edgeTypeMapById   = {};
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
	var attrDescriptorMap = this._attrDescriptorMap;
	return Object.keys(attrDescriptorMap).filter(function(key) {
	    return true;
	}).map(function(key) {
	    return attrDescriptorMap[key];
	});
};

TGGraphMetadata.prototype.getNewAttributeDescriptors = function () {
	var attrDescriptorMap = this._attrDescriptorMap;
	var newAttributeDescriptors = Object.keys(attrDescriptorMap).filter(function(key) {
	    return attrDescriptorMap[key].getAttributeId() < 0;
	}).map(function(key) {
	    return attrDescriptorMap[key];
	});
	
	for (var index in newAttributeDescriptors) {
		console.log("Confirm ne attribute descriptor :  index = " + index + ", id = " + newAttributeDescriptors[index].getAttributeId()
				                                                          + ", name = " + newAttributeDescriptors[index].getName());
	}
	
	return newAttributeDescriptors;
};

/**
 * Get the Attribute Descriptor.
 * @param attributeName
 */
TGGraphMetadata.prototype.getAttributeDescriptor = function(attributeName) {
	//console.log('getAttributeDescriptor : attributeName = ' + attributeName);
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
	//console.log('createAttributeDescriptor : attributeName = ' + attributeName);
    this._attrDescriptorMap[attributeName] = attributeDescriptor;	
	
    return attributeDescriptor;
};

/**
 * Create Node type for a given name, and derive it from the following type.
 * @param nodeName
 * @param parentNodeType - #TGNodeType
 */
TGGraphMetadata.prototype.createNodeType = function(nodeTypeName, parentNodeType) {
    var nodeType = new TGNodeType(this, nodeTypeName, parentNodeType);
    //this._nodeTypeMap[nodeName] = nodeType;
    return nodeType;
};

/**
 * Create Edge type for a given name, and derive it from the following type.
 * @param edgeName
 * @param parentEdgeType - #TGEdgeType
 */
TGGraphMetadata.prototype.createEdgeType = function(edgeTypeName, parentEdgeType) {
    var edgeType = new TGEdgeType(this, edgeTypeName, parentEdgeType.getDirectionType(), parentEdgeType);
    //this._edgeTypeMap[edgeName] = parentEdgeType;
    return edgeType;
};

TGGraphMetadata.prototype.createCompositeKey = function (typeName) {
    return new TGCompositeKey(this, typeName);
};

TGGraphMetadata.prototype.updateMetadata = function(attrDescList, nodeTypeList, edgeTypeList) {
	var index = null;
	if (attrDescList) {
		for (index in attrDescList) {
			var desc = attrDescList[index];
			//console.log('Update attribute descriptor : id ' + desc.getAttributeId() + ', name ' + desc.getName());
			this._descriptorMapById[desc.getAttributeId()] = desc;
			this._attrDescriptorMap[desc.getName()] = desc;
		}
	}
	if (nodeTypeList) {
		for (index in nodeTypeList) {
			var nt = nodeTypeList[index];
			this._nodeTypeMap[nt.getName()] = nt;
			this._nodeTypeMapById[nt.getId()] = nt;
		}
	}
	if (edgeTypeList) {
		for (index in edgeTypeList) {
			var et = edgeTypeList[index];
			this._edgeTypeMap[et.getName()] = et;
			this._edgeTypeMapById[et.getId()] = et;
		}
	}
};

TGGraphMetadata.prototype.getAttributeDescriptorById = function(id) {
	return this._descriptorMapById[id];
};

TGGraphMetadata.prototype.getNodeTypeById = function(id) {
	return this._nodeTypeMapById[id];
};

TGGraphMetadata.prototype.getEdgeTypeById = function(id) {
	return this._edgeTypeMapById[id];
};

TGGraphMetadata.prototype.writeExternal = function(outputStream) {

};

TGGraphMetadata.prototype.readExternal = function(inputStream) {

};

exports.TGGraphMetadata = TGGraphMetadata;
