/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not
 * use this file except in compliance with the License. A copy of the License is
 * included in the distribution package with this file. You also may obtain a
 * copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */
var util            = require('util'),
    TGAttributeType  = require('./TGAttributeType'),
    TGAttributeDescriptor = require('./TGAttributeDescriptor').TGAttributeDescriptor,
    TGSystemObject  = require('./TGSystemObject').TGSystemObject;

function TGEntityType() {

	this._attributes = {};

	this._id     = null; // issued only for creation and not valid later
	this._name   = null;
	this._parent = null;

}

util.inherits(TGEntityType, TGSystemObject);

TGEntityType.prototype.getAttributeDescriptors = function() {
	var values = [];
	for ( var key in this._attributes) {
		values.push(this._attributes[key]);
	}

	return values;
};

TGEntityType.prototype.getAttributeDescriptor = function(attrName) {
	return this._attributes[attrName];
};

TGEntityType.prototype.getId = function() {
	return this._id;
};

TGEntityType.prototype.getName = function() {
	return this._name;
};

TGEntityType.prototype.derivedFrom = function() {
	return this._parent;
};

TGEntityType.prototype.writeExternal = function(outputStream) {
	//console.log("writeExternal for entity type is not implemented");
};

TGEntityType.prototype.readAttributeDescriptors = function(inputStream) {
	// FIXME: Do we save the type value??
	var typeValue = inputStream.readByte();
	var type = TGSystemObject.TGSystemType.fromValue(typeValue);
	if (type === TGSystemObject.TGSystemType.InvalidType) {
		//console.log("Entity type input stream has invalid type value : %d", typeValue);
		// FIXME: Need to throw Exception
	}

	this._id = inputStream.readInt();
	this._name = inputStream.readUTF();

	var pageSize = inputStream.readInt(); // pagesize
	var attrCount = inputStream.readShort();
	for (var i = 0; i < attrCount; i++) {
		var name = inputStream.readUTF();
		// FIXME: The stream only contains name of the descriptor.
		var attrDesc = new TGAttributeDescriptor(name, TGAttributeType.STRING);
		this._attributes[name] = attrDesc;
	}
};

module.exports = TGEntityType;
