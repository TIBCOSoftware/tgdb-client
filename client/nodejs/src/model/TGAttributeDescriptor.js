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

var TGAttributeType = require('./TGAttributeType').TGAttributeType;

/**
 *
 * @param name
 * @param type
 * @param isArray
 * @constructor
 */
var gLocalAttributeId = 0;

function TGAttributeDescriptor(name, type, isArray) {
	this._attributeId = --gLocalAttributeId;
    this._name   = name;
    this._type   = type || TGAttributeType.STRING;
    this._isArray = isArray;
}

TGAttributeDescriptor.prototype.setAttributeId = function(attributeId) {
    this._attributeId = attributeId;
};

TGAttributeDescriptor.prototype.getAttributeId = function() {
    return this._attributeId;
};

TGAttributeDescriptor.prototype.getName = function() {
    return this._name;
};


TGAttributeDescriptor.prototype.getType = function() {
    return this._type;
};

TGAttributeDescriptor.prototype.isArray = function() {
    return this._isArray;
};

TGAttributeDescriptor.prototype.writeExternal = function (outputStream) {
	console.log("****Entering TGAttributeDescriptor.writeExternal at output buffer position at : %s %d", this._name, outputStream.getPosition());
	outputStream.writeByte(0);  // sysobject type attribute descriptor
	outputStream.writeInt(this._attributeId);
	outputStream.writeUTF(this._name);
	outputStream.writeByte(this._type.value);   // Need to double check
	outputStream.writeBoolean(this._isArray);
	console.log("****Leaving TGAttributeDescriptor.writeExternal at output buffer position at : %s %d", this._name, outputStream.getPosition());
};

TGAttributeDescriptor.prototype.readExternal = function (inputStream) {
	inputStream.readByte(); // read the sysobject type field which should be 0 for attribute descriptor
    this._attributeId = inputStream.readInt();
    this._name = inputStream.readUTF();
    this._type = TGAttributeType.fromTypeId(inputStream.readByte());
    this._isArray = inputStream.readBoolean();
};

TGAttributeDescriptor.prototype.toString = function () {
	return 	'attributeId : ' + this._attributeId + 
			', name : ' + this._name + 
			', type : ' + this._type + 
			', isArray : ' + this._isArray;
}

exports.TGAttributeDescriptor = TGAttributeDescriptor;