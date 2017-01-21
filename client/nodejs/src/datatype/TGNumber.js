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

// Number.MAX_SAFE_INTEGER

function longToBytes(value) {
	var buf = [];
	var valueHigh = Math.floor(value/0x100000000);
	for (var byte = 56; byte >= 32; byte -= 8) {
		buf.push((valueHigh >> byte) & (0xFF));
	}
	
	for (byte = 24; byte >= 0; byte -= 8) {
		buf.push((value >> byte) & (0xFF));
	}
	return buf;
}

function bytesToHexString(buf) {
	var str = '';
	for(var i=0; i<8; i++) {
		if(i!==0) {
			str += ' ';
		}
		str += decimalToHexString(buf[i]);
	}
	return str;
}

function decimalToHexString(number) {
	if (number < 0) {
		number = 0xFFFFFFFF + number + 1;
	}

	var hex = number.toString(16).toUpperCase();
	while (hex.length < 2) {
		hex = "0" + hex;
	}
	return hex;
}

function TGLong (value) {
	if(value===null&&value===undefined) {
		return;
	}
	this._bytes = longToBytes(value);
	this._hexStr = bytesToHexString(this._bytes);
}

TGLong.prototype.equals = function (longObj) {
	if (! (longObj instanceof TGLong)) {
		return false;
	}
	
	return this.equalsToBytes(longObj._bytes);
};

TGLong.prototype.equalsToBytes = function (bytes) {
	for(var index in this._bytes) {
		if(this._bytes[index]!==bytes[index]) {
			return false;
		}
	}

	return true;
};

TGLong.prototype.setBytes = function (bytes) {
	if (bytes.length !== 8) {
		throw new Error("Invalid bytes length, expects 8");
	}
    this._bytes = bytes.slice(0, bytes.length);
    this._hexStr = bytesToHexString(this._bytes);
};

TGLong.prototype.getBytes = function () {
	return this._bytes;
};

TGLong.prototype.getHexString = function () {
	return this._hexStr;
};



function BigDecimal (strValue, precision, scale) {
	if(!strValue) {
		return;
	}
	this._strValue = strValue;
	this._precision = precision;
	this._scale =scale;
}

TGLong.prototype.getStringValue = function () {
	return this._strValue;
};

TGLong.prototype.getPrecision = function () {
	return this._precision;
};

TGLong.prototype.getScale = function () {
	return this._scale;
};

var TGNumber = {
	getLongFromBytes : function(bytes) {
	    if(!bytes) {
	    	return new Error('Invalid byte array for create Long.');
	    } else {
	        var long = new TGLong();
	        long.setBytes(bytes);
	        return long;
	    }
	},
	getLong : function(value) {
	    return new TGLong(value);
	},
	getBigDecimal : function (strValue, precision, scale) {
		return new BigDecimal (strValue, precision, scale);
	},
	isTGLong : function (data) {
		return (data instanceof TGLong);
	}
};

module.exports = TGNumber;

function test() {
	var buf = longToBytes(13);
	var buf1 = longToBytes(13);
	var id01 = TGNumber.getLongFromBytes(buf);
	var id02 = TGNumber.getLongFromBytes(buf);
	var id03 = TGNumber.getLong(0);

	console.log('id01 = ' + id01.getBytes());
	console.log('id02 = ' + id02.getHexString());
	console.log('id03 = ' + id03.getBytes());
	console.log('Equals ? ' + id01.equals(id02));
}

//test();
