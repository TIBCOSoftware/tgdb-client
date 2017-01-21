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
var TGNumber = require('../datatype/TGNumber');

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
		str += buf[i].toString(16);
	}
	return str;
}

var TGEntityId = {
	globalSequenceNo : 0,
    createId : function(buf) {
    	if(!buf) {
    		return TGNumber.getLong(--this.globalSequenceNo);
    	} else {
        	var id = TGNumber.getLong();
        	id.setBytes(buf);
        	return id;
    	}
    },
    bytesToString : function(bytes) {
    	return bytesToHexString(bytes);
    }
};

module.exports = TGEntityId;

function test() {
	var buf = longToBytes(13);
	var buf1 = longToBytes(13);
	var id01 = TGEntityId.createId(buf);
	var id02 = TGEntityId.createId(buf1);

	console.log('id01 = ' + id01.getBytes());
	console.log('id02 = ' + id02.getHexString());
	console.log('Equals ? ' + id01.equals(id02));
}

//test();
