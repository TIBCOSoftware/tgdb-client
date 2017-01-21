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

exports.HexUtils = {
	_NullString : "0000",
	_Space : ' ',
	_NewLine : "\r\n",

	formatHex : function(buf) {

		if (!buf) {
			return NullString;
		}

		return this.formatHexToWriter(buf, 48, 0);
	},

	formatHexToWriter : function(buf, lineLength, actualLength) {
		var blen = buf.length;
		var bNewLine = false;
		var lineNo = 1;
		var writer = ("Formatted Byte Array:");

		writer += (this._NewLine);
		writer += (decimalToHexString(0, 8));
		writer += (this._Space);

		if (actualLength > 0) {
			blen = actualLength;
		}

		for (var i = 0; i < blen; i++) {
			if (bNewLine) {
				bNewLine = false;
				writer += (this._NewLine);
				writer += (decimalToHexString(lineNo * lineLength, 8));
				writer += (this._Space);
			}

			writer += (decimalToHexString(buf[i]));

			if ((i + 1) % 2 == 0)
				writer += (this._Space);

			if ((i + 1) % lineLength == 0) {
				bNewLine = true;
				++lineNo;
			}
		}

		return writer;
	}
};

function decimalToHexString(number, padding) {
	if (!padding) {
		padding = 2;
	}
	if (number < 0) {
		number = 0xFFFFFFFF + number + 1;
	}

	var hex = number.toString(16).toUpperCase();
	while (hex.length < padding) {
		hex = "0" + hex;
	}
	return hex;
}
