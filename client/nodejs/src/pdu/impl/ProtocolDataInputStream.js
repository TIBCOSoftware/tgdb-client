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

var TGNumber     = require('../../datatype/TGNumber');
var TGLogManager = require('../../log/TGLogManager');
var TGLogLevel   = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function ProtocolDataInputStream(buffer) {
    this._buffer = buffer;
    //Current position in the internal buffer
    this._currentPos = 0;
    var referenceMap = null;
}

/**
 *
 * @returns {boolean}
 */
ProtocolDataInputStream.prototype.readBoolean = function () {
    var buffer = this._buffer;
    var currentPos = this._currentPos;
    var bool = buffer[currentPos++];
    this._currentPos = currentPos;
    var value = bool !== 0;
    logger.logDebugWire("after readBoolean buffer (%s) position at : %d", value, this._currentPos);

    return value;
};

/**
 *
 * @returns {byte}
 */
ProtocolDataInputStream.prototype.readByte = function () {
    var buffer = this._buffer;
    var currentPos = this._currentPos;
    var byte = buffer[currentPos++];
    this._currentPos = currentPos;
    //Converting unsigned to signed
    //Take 2s complement = 1s complement + add 1
    //This translates to subtracting from 256 & multiple by -1
    var value = !(byte & 0x80) ? byte : (0xff - byte + 1) * -1;
	logger.logDebugWire("after readByte buffer (%s) position at : %d", value, this._currentPos);
    return value;
};

/**
 *
 * @returns {short}
 */
ProtocolDataInputStream.prototype.readShort = function () {
    var buffer = this._buffer;
    var currentPos = this._currentPos;
    var short = (buffer[currentPos++] << 8) + (buffer[currentPos++] & 0xff);
    this._currentPos = currentPos;
    //Converting unsigned to signed
    //Take 2s complement = 1s complement + add 1
    var value =  (short & 0x8000) ? short | 0xFFFF0000 : short;
	logger.logDebugWire("after readByte buffer (%s) position at : %d", value, this._currentPos);
    return value;
};

/**
*
* @returns {unsigned short}
*/
ProtocolDataInputStream.prototype.readUnsignedShort = function () {
    var ushort = (((this._buffer[this._currentPos++] << 8)) + 
    		       (this._buffer[this._currentPos++] & 0x00FF));
    return (ushort & 0x0000FFFF);
};

/**
*
* @returns {char as a 16 bits number}
*/
ProtocolDataInputStream.prototype.readChar = function() {
    var buffer = this._buffer;
    var char = (((buffer[this._currentPos++] << 8)) |
    		    (buffer[this._currentPos++]));
	logger.logDebugWire("after readChar buffer (%d) position at : %d", char, this._currentPos);
};

/**
 *
 * @returns {int}
 */
ProtocolDataInputStream.prototype.readInt = function () {
    var buffer = this._buffer;
    var currentPos = this._currentPos;
    var int = (((buffer[currentPos++] << 24)) |
               ((buffer[currentPos++] << 16)) |
               ((buffer[currentPos++] << 8)) |
               ((buffer[currentPos++])));
    this._currentPos = currentPos;
	logger.logDebugWire("after readInt buffer (%s) position at : %d", int, this._currentPos);
    return int;
};

/**
 *
 * @returns {long}
 */
ProtocolDataInputStream.prototype.readLong = function () {
	/* Should use bitwise 'or' to construct multibyte number */
    var value = (this.readInt() << 32)&0xFFFF0000 | (this.readInt());
	logger.logDebugWire("after readLong buffer (%s) position at : %d", value, this._currentPos);
	return value;
};

ProtocolDataInputStream.prototype.readDate = function () {
    var date = new Date(this.readInt() << 32)&0xFFFF0000 | (this.readInt());
	logger.logDebugWire("after readDate buffer (%s) position at : %d", date, this._currentPos);
	return date;
};

ProtocolDataInputStream.prototype.readTGLong = function () {
    var bytes = [];
    for(var i=0; i<8; i++) {
    	bytes.push(this._buffer[this._currentPos++]);	   
    }
    var value = TGNumber.getLongFromBytes(bytes);
	logger.logDebugWire("after readLongBytes buffer (%s) position at : %d", value.getHexString(), this._currentPos);
	return value;
};

ProtocolDataInputStream.prototype.readLongAsBytes = function () {
	/* Should use bitwise 'or' to construct multibyte number */
    var value = [];
    for(var i=0; i<8; i++) {
        value.push(this._buffer[this._currentPos++]);	   
    }
	logger.logDebugWire("after readBytes buffer (%s) position at : %d", value, this._currentPos);
	return value;
};

/**
*
* @returns {utf}
*/
ProtocolDataInputStream.prototype.readUTF = function () {
	var length = this.readUnsignedShort();
	var bytes = [];
	for(var i = 0; i<length; i++) {
		bytes.push(this._buffer[this._currentPos++]);
	}
        
	var utf = (new Buffer(bytes)).toString('utf8');
	logger.logDebugWire("after readUTF buffer (%s) position at : %d, length : %d", utf, this._currentPos, length);
        
	return utf;
};

/**
*
* @returns {
* 	IEEE754 float
* 	bits 0  ~ 22 fraction
* 	bits 23 ~ 30 exponent
* 	bits 32      sign
* }
*/
ProtocolDataInputStream.prototype.readFloat = function () {
	
	var bytes = (this._buffer[this._currentPos++]*Math.pow(2,24)) + 
	            (this._buffer[this._currentPos++]*Math.pow(2,16)) +
	            (this._buffer[this._currentPos++]*Math.pow(2,8)) +
	             this._buffer[this._currentPos++] ;
	
    var sign = (bytes & 0x80000000) ? -1 : 1;    
    var exponent = ((bytes & 0x7f800000) >> 23) - 0x7f;
    var fraction = (bytes & 0x007fffff)

    if (exponent == 128) {
        return sign * ((fraction) ? Number.NaN : Number.POSITIVE_INFINITY);
    }
    
    if (exponent == -127) {
        if (fraction == 0) return sign * 0.0;
        exponent = -126;
        fraction = fraction/0x400000;
    } else {
    	fraction = (fraction + 0x800000)/0x800000;
    }

	var float = sign * fraction * Math.pow(2, exponent);
	
	logger.logDebugWire("after readFloat buffer (%f) position at : %d", float, this._currentPos);
    return float;
};

/**
*
* @returns {
* 	IEEE754 double
* 	bits 0  ~ 51 fraction
* 	bits 52 ~ 62 exponent
* 	bits 63      sign
* }
*/
ProtocolDataInputStream.prototype.readDouble = function () {
	
	var high32 = (this._buffer[this._currentPos++]*Math.pow(2,24)) +
	             (this._buffer[this._currentPos++]*Math.pow(2,16)) +
	             (this._buffer[this._currentPos++]*Math.pow(2,8)) +
	              this._buffer[this._currentPos++] ;
	
	var low32 = (this._buffer[this._currentPos++]*Math.pow(2,24)) +
	            (this._buffer[this._currentPos++]*Math.pow(2,16)) +
	            (this._buffer[this._currentPos++]*Math.pow(2,8)) +
	             this._buffer[this._currentPos++] ;

	var sign = ( high32 & 0x80000000 ) ? -1 : 1;
	var exponent = ((high32 & 0x7ff00000) >> 20) - 1023;
	var fractionLong = (high32&0x000fffff) * Math.pow(2, 32) + low32;
	var fraction = 
		(exponent == 0) ? (fractionLong << 1) : fractionLong + 0x10000000000000;
	
	var double = sign * fraction * Math.pow(2, exponent-52);
	
	logger.logDebugWire("after readDouble buffer (%d) position at : %d", double, this._currentPos);
	return double;
};

ProtocolDataInputStream.prototype.available = function () {
    return this._buffer.length - this._currentPos;
};

ProtocolDataInputStream.prototype.getPosition = function () {
    return this._currentPos;
};

ProtocolDataInputStream.prototype.setPosition = function (pos) {
    var oldPos = this._currentPos;
    this._currentPos = pos;

    return oldPos;
};

ProtocolDataInputStream.prototype.setReferenceMap = function (map) {
	this._referenceMap = map;
};

ProtocolDataInputStream.prototype.getReferenceMap = function () {
	return this._referenceMap;
};

exports.ProtocolDataInputStream = ProtocolDataInputStream;
