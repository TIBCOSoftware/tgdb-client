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
    TGAttributeType       = require('./TGAttributeType').TGAttributeType,
    TGNumber              = require('../datatype/TGNumber'),
    TGException           = require('../exception/TGException').TGException,
    TGLogManager          = require('../log/TGLogManager'),
    TGLogLevel            = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

var DATE_ONLY    = 0;
var TIME_ONLY    = 1;
var TIMESTAMP    = 2;
var TGNoZone     = -1;
var TGZoneOffset = 0;
var TGZoneId     = 1;
var TGZoneName   = 2;

function checkType(descriptor, value) {
	if(value===null || typeof value==='undefined') {
		return;
	}
	var compatible = false;
	switch(descriptor.getType()) {
	case TGAttributeType.BOOLEAN:
		compatible = (typeof value === 'boolean');
		break;
	case TGAttributeType.BYTE:
		compatible = (typeof value === 'number');
		break;
	case TGAttributeType.CHAR:
		compatible = (typeof value === 'number');
		break;
	case TGAttributeType.SHORT:
		compatible = (typeof value === 'number');
		break;
	case TGAttributeType.INT:
		compatible = (typeof value === 'number');
		break;
	case TGAttributeType.LONG:
		compatible = (typeof value === 'number');
		break;
	case TGAttributeType.FLOAT:
		compatible = (typeof value === 'number');
		break;
	case TGAttributeType.DOUBLE:
		compatible = (typeof value === 'number');
		break;
	case TGAttributeType.NUMBER:
		compatible = (typeof value === 'string');
		break;
	case TGAttributeType.STRING:
		compatible = (typeof value === 'string');
		break;
    case TGAttributeType.DATE:
    	compatible = (value instanceof Date);
        break;
    case TGAttributeType.TIME:
    	compatible = (value instanceof Date);
        break;
    case TGAttributeType.TIMESTAMP:
    	compatible = (value instanceof Date);
        break;
	}
	if(!compatible) {
		throw new TGException('value set for '+ descriptor.getName() +
				' is incompatible with type ' + descriptor.getType().type);		
	}
}

function TGAttribute (owner, descriptor, value) {
	this._owner = owner;
	this._descriptor = descriptor;
	this._value = value;
	this._isModified = false;
	checkType(this._descriptor, this._value);
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

TGAttribute.prototype.getName = function() {
    return this._descriptor.getName();
};

TGAttribute.prototype.getValue = function() {
    return this._value;
};

TGAttribute.prototype.setValue = function (value) {
	//FIXME: Need to match the type of the attribute descriptor
	checkType(this._descriptor, value);
	if(this._descriptor.getType()===TGAttributeType.Long) {
	    this._value = value;
	} else {
	    this._value = value;		
	}
    this._isModified = true;
    if (!value) {
    	return;
    }
    
    if (this._descriptor.getType()===TGAttributeType.Number) {
        var precision = this._descriptor.getPrecision();
        var scale = this._descriptor.getScale();
        this.setPrecisionAndScale(precision, scale);
    }
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
	logger.logDebugWire(
			"**** Entering TGAttribute.writeExternal at output buffer position at : %s %d", 
			this._descriptor.getName(), outputStream.getPosition());
   	var aid = this._descriptor.getAttributeId();
    //null attribute is not allowed during entity creation
   	outputStream.writeInt(aid);
   	outputStream.writeBoolean(this.isNull());
   	if (this.isNull()) {
   		return;
   	}
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
			if(typeof this._value === 'number') {
				outputStream.writeLong(this._value);
			} else if(TGNumber.isTGLong(this._value)) {
				outputStream.writeTGLong(this._value);
			} else {
				throw new TGException('Invalid value for Long attribute, %s' + this._value);
			}
			break;
		case TGAttributeType.FLOAT:
			outputStream.writeFloat(this._value);
			break;
		case TGAttributeType.DOUBLE:
			outputStream.writeDouble(this._value);
			break;
        case TGAttributeType.NUMBER:
            this.writeNumber(outputStream);
            break;
		case TGAttributeType.STRING:
			outputStream.writeUTF(this._value);
   		 	break;
        case TGAttributeType.DATE:
            this.writeTimestamp(outputStream, DATE_ONLY);
            break;
        case TGAttributeType.TIME:
            this.writeTimestamp(outputStream, TIME_ONLY);
            break;
        case TGAttributeType.TIMESTAMP:
            this.writeTimestamp(outputStream, TIMESTAMP);
            break;
        default:
			throw new TGException('Write external method does not support : ' + this._descriptor.getType().type);
	}
	logger.logDebugWire(
			"**** Leaving TGAttribute.writeExternal at output buffer position at : %s %d", 
			this._descriptor.getName(), outputStream.getPosition());
};

TGAttribute.prototype.readExternal = function (inputStream) {
	logger.logDebugWire(
			"**** Entering TGAttribute.readExternal at output buffer position at : %d", 
			inputStream.getPosition());
	var aid = inputStream.readInt();  // attribute descriptor id
	this._descriptor = this._owner.getGraphMetadata().getAttributeDescriptorById(aid);
	
	if (!this._descriptor) {
		//FIXME: retrieve entity type together with the entity?
		logger.logWarning(
				"**** Leaving [TGAttribute.prototype.readExternal] unable to find attribute descriptor %d from metadata cache", aid);
	} else {
		logger.logDebugWire(
				'Attribute name = %s, type = %s', 
				this._descriptor.getName(), this._descriptor.getType().type); 
	
		if (inputStream.readByte() === 1) {  // Null value ??????
			this._value = null;
			return;
		}
	
		if (this._descriptor) {
			switch(this._descriptor.getType()) {
				case TGAttributeType.BOOLEAN:
					this._value = inputStream.readBoolean();
					break;
				case TGAttributeType.BYTE:
					this._value = inputStream.readByte();
					break;
				case TGAttributeType.CHAR:
					this._value = inputStream.readChar();
					break;
				case TGAttributeType.SHORT:
					this._value = inputStream.readShort();
					break;
				case TGAttributeType.INT:
					this._value = inputStream.readInt();
					break;
				case TGAttributeType.LONG:
					this._value = inputStream.readTGLong().getBytes();
					break;
				case TGAttributeType.FLOAT:
					this._value = inputStream.readFloat();
					break;
				case TGAttributeType.DOUBLE:
					this._value = inputStream.readDouble();
					break;
				case TGAttributeType.STRING:
					this._value = inputStream.readUTF();
					break;
				case TGAttributeType.NUMBER:
					this.readNumber(inputStream);
					break;
				case TGAttributeType.DATE:
					this.readTimestamp(inputStream, DATE_ONLY);
					break;
				case TGAttributeType.TIME:
					this.readTimestamp(inputStream, TIME_ONLY);
					break;
				case TGAttributeType.TIMESTAMP:
					this.readTimestamp(inputStream, TIMESTAMP);
					break;
				default:
					break;
			}
		}
		logger.logDebugWire(
				"**** Leavinging TGAttribute.readExternal at output buffer position at : %s %d", 
				this._descriptor.getName(), inputStream.getPosition());
	}
};

//FIXME: Do we need it?
TGAttribute.prototype.resetIsModified = function () {
	this._isModified = false;
};

TGAttribute.prototype.isModified = function () {
	return this._isModified;
};

TGAttribute.prototype.isNull = function () {
	return (!this._value);
};


TGAttribute.prototype.writeTimestamp = function (outputStream, component2Write){
    if (!this._value || !(this._value instanceof Date)) {
    	throw new TGException("value is null or not a Date type");
    }
    
    switch (component2Write) {
        case DATE_ONLY: {
            outputStream.writeBoolean(this._value>=new Date(0,0,0,0,0,0,0));
            outputStream.writeShort(this._value.getFullYear());
            outputStream.writeByte(this._value.getMonth());
            outputStream.writeByte(this._value.getDate);
            outputStream.writeByte(0); //HR
            outputStream.writeByte(0); //Min
            outputStream.writeByte(0); //Sec
            outputStream.writeShort(0); //msec
            outputStream.writeByte(TGNoZone); //First to indicate we have no zone support
            //outputStream.writeShort(TGNoZone); //This is for the zone ID
            break;
        }

        case TIME_ONLY: {
        	outputStream.writeBoolean(true);
        	outputStream.writeShort(0);
        	outputStream.writeByte(0);
        	outputStream.writeByte(0);
        	outputStream.writeByte(this._value.getHours()); //HR
        	outputStream.writeByte(this._value.getMinutes()); //Min
        	outputStream.writeByte(this._value.getSeconds()); //Sec
        	outputStream.writeShort(this._value.getMilliseconds()); //msec
        	outputStream.writeByte(TGNoZone); //First to indicate we have no zone support
            //outputStream.writeShort(TGNoZone); //This is for the zone ID
            break;
        }

        case TIMESTAMP:
        {
            outputStream.writeBoolean(this._value>=new Date(0,0,0,0,0,0,0));
            outputStream.writeShort(this._value.getFullYear());
            outputStream.writeByte(this._value.getMonth());
            outputStream.writeByte(this._value.getDate);
        	outputStream.writeByte(this._value.getHours()); //HR
        	outputStream.writeByte(this._value.getMinutes()); //Min
        	outputStream.writeByte(this._value.getSeconds()); //Sec
        	outputStream.writeShort(this._value.getMilliseconds()); //msec
            outputStream.writeByte(TGNoZone); //First to indicate we have no zone support
            //outputStream.writeShort(TGNoZone); //This is for the zone ID
            break;
        }

        default:
            throw new TGException("Invalid spec provided to write the Calendar");
    }
};

//SS:TODO Support for Timezone - Post v1.0
//SS:TODO Only support Gregorian Calendar - There is Japanese, Thai, ...
TGAttribute.prototype.readTimestamp = function (inputStream, component2read) {
    var year, mon, dom, hr, min, sec, ms, tztype, tzid, era;
    era     = inputStream.readBoolean();
    year    = inputStream.readShort();
    mon     = inputStream.readByte();
    dom     = inputStream.readByte();
    hr      = inputStream.readByte();
    min     = inputStream.readByte();
    sec     = inputStream.readByte();
    ms      = inputStream.readUnsignedShort();
    tztype  = inputStream.readByte();
    //tzid    = in.readShort();

    var value = null;
    switch (component2read) {
        case DATE_ONLY:
        	this._value = new Date(year,mon,dom,0,0,0,0);

            break;
        case TIME_ONLY:
        	this._value = new Date(0,0,0,hr,min,sec,ms);
        	
            break;
        case TIMESTAMP:
        	this._value = new Date(year,mon,dom,hr,min,sec,ms);
        	
            break;
        default:
            throw new TGException("Invalid spec provided to read the Calendar");

    }
};

TGAttribute.prototype.writeNumber = function (outputStream) {
	outputStream.writeShort(this._descriptor.getPrecision());
	outputStream.writeShort(this._descriptor.getScale());
	outputStream.writeUTF(this._value);
};

TGAttribute.prototype.readNumber = function (inputStream) {
    var precision = inputStream.readShort();
    var scale = inputStream.readShort();
    var bdstr = inputStream.readUTF();
    this._value = TGNumber.getBigDecimal(bdstr, precision, scale);
    //setPrecisionAndScale(precision, scale);
};

TGAttribute.prototype.setPrecisionAndScale = function (precision, scale) {

};

exports.TGAttribute = TGAttribute;
