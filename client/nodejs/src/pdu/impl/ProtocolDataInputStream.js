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

function ProtocolDataInputStream(buffer) {
    this._buffer = buffer;
    //Current position in the internal buffer
    this._currentPos = 0;
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
    var value = bool != 0;
	console.log("after readBoolean buffer (%s) position at : %d", value, this._currentPos);

    return value
};

/**
 *
 * @returns {*}
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
	console.log("after readByte buffer (%s) position at : %d", value, this._currentPos);
    return value;
};

/**
 *
 * @returns {*}
 */
ProtocolDataInputStream.prototype.readShort = function () {
    var buffer = this._buffer;
    var currentPos = this._currentPos;
    var short = (buffer[currentPos++] << 8) + (buffer[currentPos++] & 0xff);
    this._currentPos = currentPos;
    //Converting unsigned to signed
    //Take 2s complement = 1s complement + add 1
    return (short & 0x8000) ? short | 0xFFFF0000 : short;
};

ProtocolDataInputStream.prototype.readUnsignedShort = function () {
    var ushort = (((this._buffer[this._currentPos++] << 8)) + 
    		       (this._buffer[this._currentPos++] & 0x00FF));
    return (ushort & 0x0000FFFF);
};

/**
 *
 * @returns {*}
 */
ProtocolDataInputStream.prototype.readInt = function () {
    var buffer = this._buffer;
    var currentPos = this._currentPos;
    var int = (((buffer[currentPos++] << 24)) |
               ((buffer[currentPos++] << 16)) |
               ((buffer[currentPos++] << 8)) |
               ((buffer[currentPos++])));
    this._currentPos = currentPos;
	console.log("after readInt buffer (%s) position at : %d", int, this._currentPos);
    return int;
};

/**
 *
 * @returns {*}
 */
ProtocolDataInputStream.prototype.readLong = function () {
	/* Should use bitwise 'or' to construct multibyte number */
    var value = (this.readInt() << 32)&0xFFFF0000 | (this.readInt());
	console.log("after readLong buffer (%s) position at : %d", value, this._currentPos);
	return value;
};

ProtocolDataInputStream.prototype.readUTF = function () {
	var length = this.readUnsignedShort();
	var bytes = [];
	for(var i = 0; i<length; i++) {
		bytes.push(this._buffer[this._currentPos++]);
	}
        
	var utf = (new Buffer(bytes)).toString('utf8');
        
	return utf;
};

ProtocolDataInputStream.prototype.available = function () {
    return this._buffer.length - this._currentPos;
};

exports.ProtocolDataInputStream = ProtocolDataInputStream;
