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

var util            = require('util'),
    TGSystemObject  = require('./TGSystemObject').TGSystemObject,
    TGAttributeType = require('./TGAttributeType').TGAttributeType,
    TGLogManager    = require('../log/TGLogManager'),
    TGLogLevel      = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

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
    this._scale = null;
    this._precision = null;
    if (this._type === TGAttributeType.NUMBER) {
    	this._scale = 5;
    	this._precision = 20;
    }
}

util.inherits(TGAttributeDescriptor, TGSystemObject);

TGAttributeDescriptor.prototype.setAttributeId = function(attributeId) {
    this._attributeId = attributeId;
};

TGAttributeDescriptor.prototype.getAttributeId = function() {
    return this._attributeId;
};

/**
 * Get the system type enum
 * @return the system type of the object
 */
TGAttributeDescriptor.prototype.getSystemType = function () {
	return TGSystemObject.TGSystemType.AttributeDescriptor;
};

/**
 * Get the type name.
 * @return the name of the object
 */
TGAttributeDescriptor.prototype.getName = function() {
    return this._name;
};

TGAttributeDescriptor.prototype.getType = function() {
    return this._type;
};

TGAttributeDescriptor.prototype.isArray = function() {
    return this._isArray;
};

TGAttributeDescriptor.prototype.setPrecision = function(precision) {
    if (this._type === TGAttributeType.Number) {
    	this._precision = precision;
    }
};

TGAttributeDescriptor.prototype.setScale = function(scale) {
    if (this._type === TGAttributeType.Number) {
    	this._scale = scale;
    }
};

TGAttributeDescriptor.prototype.getPrecision = function() { 
	return this._precision; 
};

TGAttributeDescriptor.prototype.getScale = function() { 
	return this._scale; 
};

TGAttributeDescriptor.prototype.writeExternal = function (outputStream) {
    logger.logDebugWire(
    		"****Entering TGAttributeDescriptor.writeExternal at output buffer position at : %s %d", 
    		this._name, outputStream.getPosition());
	outputStream.writeByte(TGSystemObject.TGSystemType.AttributeDescriptor.value);  // sysobject type attribute descriptor
	outputStream.writeInt(this._attributeId);
	outputStream.writeUTF(this._name);
	outputStream.writeByte(this._type.value);   // Need to double check
	outputStream.writeBoolean(this._isArray);
    if (this._type === TGAttributeType.NUMBER) {
    	outputStream.writeShort(this._precision);
    	outputStream.writeShort(this._scale);
    }
    logger.logDebugWire(
    		"****Leaving TGAttributeDescriptor.writeExternal at output buffer position at : %s %d",
    		this._name, outputStream.getPosition());
};

TGAttributeDescriptor.prototype.readExternal = function (inputStream) {
    logger.logDebugWire(
    		"****Entering TGAttributeDescriptor.readExternal at output buffer position at : %s %d",
    		this._name, inputStream.getPosition());
	var sysObjectType = inputStream.readByte(); // read the sysobject type field which should be 0 for attribute descriptor
	var stype = TGSystemObject.TGSystemType.fromValue(sysObjectType);
	if (stype !== TGSystemObject.TGSystemType.AttributeDescriptor) {
		logger.logWarning("Attribute descriptor has invalid input stream value : %d", sysObjectType);
		//FIXME: Throw exception is needed
	}
	
    this._attributeId = inputStream.readInt();
    logger.logDebugWire(
    		'Attribute descriptor retrieved from server, aid = %d', this._attributeId);
    this._name = inputStream.readUTF();
    this._type = TGAttributeType.fromTypeId(inputStream.readByte());
    this._isArray = inputStream.readBoolean();
    if (this._type === TGAttributeType.NUMBER) {
    	this._precision = inputStream.readShort();
    	this._scale = inputStream.readShort();
    }
    logger.logDebugWire(
    		"****Leaving TGAttributeDescriptor.readExternal at output buffer position at : %s %d", 
    		this._name, inputStream.getPosition());
};

TGAttributeDescriptor.prototype.toString = function () {
	return 	'attributeId : ' + this._attributeId + 
			', name : ' + this._name + 
			', type : ' + this._type + 
			', isArray : ' + this._isArray;
};

exports.TGAttributeDescriptor = TGAttributeDescriptor;