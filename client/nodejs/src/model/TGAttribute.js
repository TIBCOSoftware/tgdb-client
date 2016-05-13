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

var TGAttributeDescriptor = require('./TGAttributeDescriptor').TGAttributeDescriptor,
    TGAttributeType = require('./TGAttributeType').TGAttributeType;


function TGAttribute (owner, descriptor, value) {
	this._owner = owner;
	this._descriptor = descriptor;
	this._value = value;
	this._isModified = false;
}

TGAttribute.prototype = function (owner) {
	this._owner = owner;
};

TGAttribute.prototype.getOwner = function () {
	return this._owner;
};

TGAttribute.prototype.getAttributeDescriptor = function () {
	return this._descriptor;
};

TGAttribute.prototype.getValue = function() {
    return this._value;
};

TGAttribute.prototype.setValue = function (value) {
	//FIXME: Need to match the type of the attribute descriptor
    this._value = value;
    this._isModified = true;
	console.log('IsModified flag is set +++++++++++++++++++++++++++ ' + this._isModified);
};

TGAttribute.prototype.getAsBoolean = function () {
    return null;
};

TGAttribute.prototype.getAsByte = function () {
    return null;
};

TGAttribute.prototype.getAsChar = function () {
    return null;
};

TGAttribute.prototype.getAsShort = function () {
    return null;
};

TGAttribute.prototype.getAsInt = function () {
    return null;
};

TGAttribute.prototype.getAsLong = function () {
    return null;
};

TGAttribute.prototype.getAsFloat = function () {
    return null;
};

TGAttribute.prototype.getAsDouble = function () {
    return null;
};

TGAttribute.prototype.getAsString = function () {
    return null;
};

TGAttribute.prototype.writeExternal = function (outputStream) {
	console.log("**** Entering TGAttribute.writeExternal at output buffer position at : %s %d", this._descriptor.getName(), outputStream.getPosition());
	switch(this._descriptor.getType()) {
		case TGAttributeType.BOOLEAN:
			outputStream.writeBoolean(this._value);
			break;
		case TGAttributeType.BYTE:
			outputStream.writeByte(this._value);
			break;
		case TGAttributeType.CHAR:
			outputStream.writeChar(this._value);
			break;
		case TGAttributeType.SHORT:
			outputStream.writeShort(this._value);
			break;
		case TGAttributeType.INT:
			outputStream.writeInt(this._value);
			break;
		case TGAttributeType.LONG:
			outputStream.writeLong(this._value);
			break;
		case TGAttributeType.FLOAT:
			outputStream.writeFloat(this._value);
			break;
		case TGAttributeType.DOUBLE:
			outputStream.writeDouble(this._value);
			break;
		case TGAttributeType.NUMBER:
			outputStream.writeUTF(this._value);
			break;
		case TGAttributeType.STRING:
			outputStream.writeUTF(this._value);
   		 	break;
		default:
			throw new Error('Write external method does not support : ' + this._descriptor.getType().type);
	}
	console.log("**** Leavinging TGAttribute.writeExternal at output buffer position at : %s %d", this._descriptor.getName(), outputStream.getPosition());
};

TGAttribute.prototype.readExternal = function (inputStream) {

};

//FIXME: Do we need it?
TGAttribute.prototype.resetIsModified = function () {
	this._isModified = false;
};

TGAttribute.prototype.isModified = function () {
	return this._isModified;
};

TGAttribute.prototype.isNull = function () {
	return (this._value == null);
};

exports.TGAttribute = TGAttribute;
