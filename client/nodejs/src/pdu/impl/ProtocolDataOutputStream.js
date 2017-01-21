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
var TGNumber = require('../../datatype/TGNumber');
var HexUtils = require('../../utils/HexUtils').HexUtils;
var TGLogManager  = require('../../log/TGLogManager');
var TGLogLevel    = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function ProtocolDataOutputStream() {
	// Internal buffer
	this._buffer = [];
	this._currentLength = 0;
}

ProtocolDataOutputStream.prototype.getPosition = function() {
	return this._currentLength;
};

ProtocolDataOutputStream.prototype.writeBoolean = function(value) {
	var start = this.getPosition();
	logger.logDebugWire("before writeBoolean buffer (%s) position at : %d", value, this.getPosition());
	this._buffer.push(value ? 1 : 0);
	this._currentLength++;
	var end = this.getPosition();
	logger.logDebugWire("after writeBoolean buffer (%s) position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

ProtocolDataOutputStream.prototype.writeByte = function(value, verbose) {
	var start = this.getPosition();
	if(!verbose)
		logger.logDebugWire("before writeByte buffer (%s) position at : %d", value, this.getPosition());
	if (value > 255 || value < -1) {
		throw new Error('Invalid byte : ' + value);
	}
	this._buffer.push(value);
	this._currentLength++;
	var end = this.getPosition();
	if(!verbose)
	    logger.logDebugWire("after writeByte buffer (%s) position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

ProtocolDataOutputStream.prototype.writeBytes = function(bytes) {
	var length = bytes.length;
	// Write its length as int
	this.writeInt(length);
	for (var i = 0; i < length; i++) {
		this.writeByte(bytes[i]);
	}
};

ProtocolDataOutputStream.prototype.writeTGLong = function(tgLong) {
	var bytes = tgLong.getBytes();
	logger.logDebugWire("before writeTGLong buffer (%s) position at : %d", tgLong.getHexString(), this.getPosition());
	var length = bytes.length;
	if(length!==8) {
		throw new Error('Invalid buffer length, expected : 8');
	}
	for (var i = 0; i < length; i++) {
		this.writeByte(bytes[i], 'yes');
	}
	logger.logDebugWire("after writeTGLong buffer (%s) position at : %d", tgLong.getHexString(), this.getPosition());
};

ProtocolDataOutputStream.prototype.writeLongAsBytes = function(bytes) {
	logger.logDebugWire("before writeLongAsBytes buffer (%s) position at : %d", bytes, this.getPosition());
	var length = bytes.length;
	if(length!==8) {
		throw new Error('Invalid buffer length, expected : 8');
	}
	for (var i = 0; i < length; i++) {
		this.writeByte(bytes[i], 'yes');
	}
	logger.logDebugWire("after writeLongAsBytes buffer (%s) position at : %d", bytes, this.getPosition());
};

ProtocolDataOutputStream.prototype.writeDate = function(date) {
	var start = this.getPosition();
	logger.logDebugWire("before writeDate buffer (%s) position at : %d", date, this.getPosition());
	this._currentLength += writeLong64(date.getTime(), this._buffer);
	var end = this.getPosition();
	logger.logDebugWire("after writeDate buffer (%s) position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

/**
 * Write java-like short (16 bits)
 * 
 * @param value
 */
ProtocolDataOutputStream.prototype.writeShort = function(value) {
	var start = this.getPosition();
	logger.logDebugWire("before writeShort buffer (%s) position at : %d", value, this.getPosition());
	this._currentLength += bitwiseOperation(value, this._buffer, 2);
	var end = this.getPosition();
	logger.logDebugWire("after writeShort buffer (%s) position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

ProtocolDataOutputStream.prototype.writeChar = function(value) {
	var start = this.getPosition();
	logger.logDebugWire("before writeChar (%s) buffer position at : %d\", value", this.getPosition());
	this._buffer[this._currentLength++] = (value >> 8) & (0xFF);
	this._buffer[this._currentLength++] = (value) & (0xFF);
	var end = this.getPosition();
	logger.logDebugWire("after writeChar (%s) buffer position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

/**
 * Write java - like int (32 bits)
 * 
 * @param value
 */
ProtocolDataOutputStream.prototype.writeInt = function(value, verbose) {
	var start = this.getPosition();
	logger.logDebugWire("before writeInt buffer (%s) position at : %d", value, this.getPosition());
	this._currentLength += bitwiseOperation(value, this._buffer, 4);
	var end = this.getPosition();
	logger.logDebugWire("after writeInt buffer (%s) position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

/**
 * Write java - like int (32 bits)
 * 
 * @param position ->
 *            Index in the internal array to write
 * @param value
 */
ProtocolDataOutputStream.prototype.writeIntAt = function(position, value) {
	var start = this.getPosition();
	logger.logDebugWire("before writeIntAt buffer (%s) position at : %d", value, position);
	for (var byte = 24; byte >= 0; byte -= 8) {
		this._buffer[position++] = (value >> byte) & (0xFF);
	}
	var end = this.getPosition();
	logger.logDebugWire("after writeIntAt buffer (%s) position at : %d", bufferDelta(start, end, this._buffer), position);
};

/**
 * Write java - like float (32 bits)
 * 
 * @param value
 */
/**
 * Returns the {@code float} value corresponding to a given bit representation.
 * The argument is considered to be a representation of a floating-point value
 * according to the IEEE 754 floating-point "single format" bit layout.
 * 
 * <p>
 * If the argument is {@code 0x7f800000}, the result is positive infinity.
 * 
 * <p>
 * If the argument is {@code 0xff800000}, the result is negative infinity.
 * 
 * <p>
 * If the argument is any value in the range {@code 0x7f800001} through
 * {@code 0x7fffffff} or in the range {@code 0xff800001} through
 * {@code 0xffffffff}, the result is a NaN. No IEEE 754 floating-point
 * operation provided by Java can distinguish between two NaN values of the same
 * type with different bit patterns. Distinct values of NaN are only
 * distinguishable by use of the {@code Float.floatToRawIntBits} method.
 * 
 * <p>
 * In all other cases, let <i>s</i>, <i>e</i>, and <i>m</i> be three values
 * that can be computed from the argument:
 * 
 * <blockquote>
 * 
 * <pre>
 * {@code
 * int s = ((bits &gt;&gt; 31) == 0) ? 1 : -1;
 * int e = ((bits &gt;&gt; 23) &amp; 0xff);
 * int m = (e == 0) ?
 *                 (bits &amp; 0x7fffff) &lt;&lt; 1 :
 *                 (bits &amp; 0x7fffff) | 0x800000;
 * }
 * </pre>
 * 
 * </blockquote>
 * 
 * Then the floating-point result equals the value of the mathematical
 * expression <i>s</i>&middot;<i>m</i>&middot;2<sup><i>e</i>-150</sup>.
 * 
 * <p>
 * Note that this method may not be able to return a {@code float} NaN with
 * exactly same bit pattern as the {@code int} argument. IEEE 754 distinguishes
 * between two kinds of NaNs, quiet NaNs and <i>signaling NaNs</i>. The
 * differences between the two kinds of NaN are generally not visible in Java.
 * Arithmetic operations on signaling NaNs turn them into quiet NaNs with a
 * different, but often similar, bit pattern. However, on some processors merely
 * copying a signaling NaN also performs that conversion. In particular, copying
 * a signaling NaN to return it to the calling method may perform this
 * conversion. So {@code intBitsToFloat} may not be able to return a
 * {@code float} with a signaling NaN bit pattern. Consequently, for some
 * {@code int} values, {@code floatToRawIntBits(intBitsToFloat(start))} may
 * <i>not</i> equal {@code start}. Moreover, which particular bit patterns
 * represent signaling NaNs is platform dependent; although all NaN bit
 * patterns, quiet or signaling, must be in the NaN range identified above.
 * 
 * @param bits
 *            an integer.
 * @return the {@code float} floating-point value with the same bit pattern.
 */

ProtocolDataOutputStream.prototype.writeFloat = function(value) {
	var start = this.getPosition();
	logger.logDebugWire("before writeFloat (%s) buffer position at : %d", value, this.getPosition());
	// floating[sign, exponent, fraction]
	var encodingPara = getIEEE754EncodingParameters(value, 8, 23);

	this._buffer.push((encodingPara[1] >> 1 & 0x7F)
			+ ((encodingPara[0] ? 1 : 0) * 0x80));

	this._buffer.push(encodingPara[2] >> 16 & 0xFF + (encodingPara[1] & 0x01) * 0x10);
	this._buffer.push(encodingPara[2] >> 8 & 0xFF);
	this._buffer.push(encodingPara[2] >> 0 & 0xFF);

	this._currentLength += 4;
	var end = this.getPosition();
	logger.logDebugWire("after writeFloat (%s) buffer position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

/**
 * Write java - like long (64 bits)
 * 
 * @param value
 */
ProtocolDataOutputStream.prototype.writeLong = function(value) {
	var start = this.getPosition();
	logger.logDebugWire("before writeLong buffer (%s) position at : %d", value, this.getPosition());
	this._currentLength += writeLong64(value, this._buffer);
	var end = this.getPosition();
	logger.logDebugWire("after writeLong buffer (%s) position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

/**
 * Write java - like double (64 bits)
 * 
 * @param value
 */
/**
 * Returns a representation of the specified floating-point value according to
 * the IEEE 754 floating-point "double format" bit layout, preserving
 * Not-a-Number (NaN) values.
 * 
 * <p>
 * Bit 63 (the bit that is selected by the mask {@code 0x8000000000000000L})
 * represents the sign of the floating-point number. Bits 62-52 (the bits that
 * are selected by the mask {@code 0x7ff0000000000000L}) represent the
 * exponent. Bits 51-0 (the bits that are selected by the mask
 * {@code 0x000fffffffffffffL}) represent the significand (sometimes called the
 * mantissa) of the floating-point number.
 * 
 * <p>
 * If the argument is positive infinity, the result is
 * {@code 0x7ff0000000000000L}.
 * 
 * <p>
 * If the argument is negative infinity, the result is
 * {@code 0xfff0000000000000L}.
 * 
 * <p>
 * If the argument is NaN, the result is the {@code long} integer representing
 * the actual NaN value. Unlike the {@code doubleToLongBits} method,
 * {@code doubleToRawLongBits} does not collapse all the bit patterns encoding a
 * NaN to a single "canonical" NaN value.
 * 
 * <p>
 * In all cases, the result is a {@code long} integer that, when given to the
 * {@link #longBitsToDouble(long)} method, will produce a floating-point value
 * the same as the argument to {@code doubleToRawLongBits}.
 * 
 * @param value
 *            a {@code double} precision floating-point number.
 * @return the bits that represent the floating-point number.
 * @since 1.3
 */

/**
 * Returns the {@code double} value corresponding to a given bit representation.
 * The argument is considered to be a representation of a floating-point value
 * according to the IEEE 754 floating-point "double format" bit layout.
 * 
 * <p>
 * If the argument is {@code 0x7ff0000000000000L}, the result is positive
 * infinity.
 * 
 * <p>
 * If the argument is {@code 0xfff0000000000000L}, the result is negative
 * infinity.
 * 
 * <p>
 * If the argument is any value in the range {@code 0x7ff0000000000001L} through
 * {@code 0x7fffffffffffffffL} or in the range {@code 0xfff0000000000001L}
 * through {@code 0xffffffffffffffffL}, the result is a NaN. No IEEE 754
 * floating-point operation provided by Java can distinguish between two NaN
 * values of the same type with different bit patterns. Distinct values of NaN
 * are only distinguishable by use of the {@code Double.doubleToRawLongBits}
 * method.
 * 
 * <p>
 * In all other cases, let <i>s</i>, <i>e</i>, and <i>m</i> be three values
 * that can be computed from the argument:
 * 
 * <blockquote>
 * 
 * <pre>
 * {@code
 * int s = ((bits &gt;&gt; 63) == 0) ? 1 : -1;
 * int e = (int)((bits &gt;&gt; 52) &amp; 0x7ffL);
 * long m = (e == 0) ?
 *                 (bits &amp; 0xfffffffffffffL) &lt;&lt; 1 :
 *                 (bits &amp; 0xfffffffffffffL) | 0x10000000000000L;
 * }
 * </pre>
 * 
 * </blockquote>
 * 
 * Then the floating-point result equals the value of the mathematical
 * expression <i>s</i>&middot;<i>m</i>&middot;2<sup><i>e</i>-1075</sup>.
 * 
 * <p>
 * Note that this method may not be able to return a {@code double} NaN with
 * exactly same bit pattern as the {@code long} argument. IEEE 754 distinguishes
 * between two kinds of NaNs, quiet NaNs and <i>signaling NaNs</i>. The
 * differences between the two kinds of NaN are generally not visible in Java.
 * Arithmetic operations on signaling NaNs turn them into quiet NaNs with a
 * different, but often similar, bit pattern. However, on some processors merely
 * copying a signaling NaN also performs that conversion. In particular, copying
 * a signaling NaN to return it to the calling method may perform this
 * conversion. So {@code longBitsToDouble} may not be able to return a
 * {@code double} with a signaling NaN bit pattern. Consequently, for some
 * {@code long} values, {@code doubleToRawLongBits(longBitsToDouble(start))} may
 * <i>not</i> equal {@code start}. Moreover, which particular bit patterns
 * represent signaling NaNs is platform dependent; although all NaN bit
 * patterns, quiet or signaling, must be in the NaN range identified above.
 * 
 * @param bits
 *            any {@code long} integer.
 * @return the {@code double} floating-point value with the same bit pattern.
 */
ProtocolDataOutputStream.prototype.writeDouble = function(value) {
	var start = this.getPosition();
	logger.logDebugWire("before writeDouble buffer (%s) position at : %d", value, this.getPosition());
	// [sign, exponent, fraction]
	var encodingPara = getIEEE754EncodingParameters(value, 11, 52);

	fh = Math.floor(encodingPara[2] / 0x100000000);

	this._buffer.push((encodingPara[1] >> 4 & 0x7F)
			+ ((encodingPara[0] ? 1 : 0) * 0x80));

	this._buffer.push((fh >> 16 & 0xFF) + (encodingPara[1] & 0x0F) * 0x10);

	this._buffer.push(fh >> 8 & 0xFF);
	this._buffer.push(fh >> 0 & 0xFF);

	this._buffer.push(encodingPara[2] >> 24 & 0xFF);
	this._buffer.push(encodingPara[2] >> 16 & 0xFF);
	this._buffer.push(encodingPara[2] >> 8 & 0xFF);
	this._buffer.push(encodingPara[2] >> 0 & 0xFF);

	this._currentLength += 8;
	var end = this.getPosition();
	logger.logDebugWire("after writeDouble buffer (%s) position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

function getIEEE754EncodingParameters(value, exponentBits, fractionBits) {

	var encodingPara = []; // encodingPara[sign, exponent, fraction];
	var biasedExponent = (1 << (exponentBits - 1)) - 1;

	if (isNaN(value)) {
		encodingPara.push(0); // sign
		encodingPara.push((1 << biasedExponent) - 1); // exponent
		encodingPara.push(1); // fraction
	} else if (value === Infinity || value === -Infinity) {
		encodingPara.push((value < 0) ? 1 : 0);
		encodingPara.push((1 << biasedExponent) - 1);
		encodingPara.push(0);
	} else if (value === 0) {
		encodingPara.push((1 / value === -Infinity) ? 1 : 0);
		encodingPara.push(0);
		encodingPara.push(0);
	} else {
		encodingPara.push((value < 0) ? 1 : 0);
		value = Math.abs(value);

		if (value >= Math.pow(2, 1 - biasedExponent)) {
			var ln = Math.min(Math.floor(Math.log(value) / Math.LN2),
					biasedExponent);
			encodingPara.push(ln + biasedExponent);
			encodingPara.push(value * Math.pow(2, fractionBits - ln)
					- Math.pow(2, fractionBits));
		} else {
			encodingPara.push(0);
			encodingPara.push(value
					/ Math.pow(2, 1 - biasedExponent - fractionBits));
		}
	}
	//logger.logDebugWire('sign = ' + encodingPara[0] + ', exponent = ' + encodingPara[1]
	//		+ ', fraction = ' + encodingPara[2]);

	return encodingPara;
}

/**
 * 
 * @param string
 */
ProtocolDataOutputStream.prototype.writeUTF = function(string) {
	var start = this.getPosition();
	logger.logDebugWire("before writeUTF (%s) buffer position at : %d", string, this.getPosition());
	var utfLength = 0;

	for (var i = 0; i < string.length; i++) {
		var char = string.charCodeAt(i);
		if ((char >= 0x0001) && (char <= 0x007F)) {
			// single byte chars
			utfLength++;
		} else if (char > 0x07FF) {
			// Triple byte chars
			utfLength += 3;
		} else {
			// Double byte chars
			utfLength += 2;
		}
	}
	// Write length
	this.writeShort(utfLength);
	writeUTFString(this, string, utfLength);
	var end = this.getPosition();
	logger.logDebugWire("after writeUTF (%s) buffer position at : %d", bufferDelta(start, end, this._buffer), this.getPosition());
};

/**
 * @param outputStream
 * @param string
 * @param utfLength
 */
function writeUTFString(outputStream, string, utfLength) {
	var internal = outputStream._buffer;
	var buffer = new Int8Array(utfLength);
	var currentCount = 0;

	for (var i = 0; i < string.length; i++) {
		var char = string.charCodeAt(i);

		if ((char >= 0x0001) && (char <= 0x007F)) {
			internal.push(char & (0xFF));
		} else if (char > 0x07FF) {
			internal.push(0xE0 | ((char >> 12) & 0x0F));
			internal.push(0x80 | ((char >> 6) & 0x3F));
			internal.push(0x80 | ((char >> 0) & 0x3F));
		} else {
			internal.push(0xC0 | ((char >> 6) & 0x1F));
			internal.push(0x80 | ((char >> 0) & 0x3F));
		}
	}
	outputStream._currentLength += utfLength;
}

/**
 * Return the length of the stream
 */
ProtocolDataOutputStream.prototype.length = function() {
	return this._currentLength;
};

/**
 * Convert the outputstream to cumulative buffer
 * 
 * @param value
 */
ProtocolDataOutputStream.prototype.toBuffer = function(value) {
	var buffer = new Buffer(this._currentLength);
	var counter = 0;

	// Copy contents into custom buffer
	// We cannot use in built buffer because it
	// performs bitwise & for every entry with 255
	for (var j = 0; j < this._currentLength; j++) {
		buffer[counter++] = this._buffer[j];
	}
	return buffer;
};

function writeLong64(value, buffer) {
	var length = 0;
	var valueHigh = Math.floor(value/0x100000000);
	var temp = null;
	for (var byte = 56; byte >= 32; byte -= 8) {
		temp = (valueHigh >> byte) & (0xFF);
		
		
		buffer.push(temp);
		length += 1;
	}
	
	for (var byte = 24; byte >= 0; byte -= 8) {
		buffer.push((value >> byte) & (0xFF));
		length += 1;
	}
	
	return length;
}


function bitwiseOperation(value, buffer, byteWidth) {
	var length = 0;

	if (byteWidth > 4) {
		length += bitwiseOperation32(
				(value - (value % 0x100000000)) / 0x100000000, buffer, 4, 7);
		byteWidth = 4;
	}
	length += bitwiseOperation32(value, buffer, 0, byteWidth - 1);

	return length;
}

function bitwiseOperation32(value, buffer, lowestByte, highestByte) {
	var length = 0;
	for (var byte = 8 * highestByte; byte >= 8 * lowestByte; byte -= 8) {
		buffer.push((value >> byte) & (0xFF));
		length += 1;
	}
	return length;
}

function bufferDelta(start, end, buffer) {
	var writer = '';
	var number = null;
	var hex = null;
    for (var i=start; i < end; i++)
    {        	
    	number = buffer[i];
    	
    	if(number!==null) {
    		if (number < 0) {
    			number = 0xFFFFFFFF + number + 1;
    		}
    		
    		hex = number.toString(16).toUpperCase();
    	
    		while (hex.length < 2) {
    			hex = "0" + hex;
    		}
    	} else {
    		hex = '@@'
    	}

        writer += hex;
        writer += (' ');
    }
    return writer;
}

exports.ProtocolDataOutputStream = ProtocolDataOutputStream;

function testUTF() {
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.DebugWire);

	var utfString = '是超人嗎';
	for(var i = 0 ; i<utfString.length; i++){
		logger.logDebugWire(utfString.charCodeAt(i));
	}
	
	var outputStream = new ProtocolDataOutputStream();
	outputStream.writeUTF(utfString);
	logger.logDebugWire(HexUtils.formatHex(outputStream._buffer));
}

//testUTF();

