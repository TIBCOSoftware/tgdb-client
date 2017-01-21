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

exports.TGAttributeType = {
	    INVALID     : {value : 0, type  : null},
	    BOOLEAN     : {value : 1, type  : 'Boolean'},  //A single bit representing the truth value
	    BYTE        : {value : 2, type  : 'Byte'},     //8bit octet
	    CHAR        : {value : 3, type  : 'Char'},     //Fixed 8-Bit octet of N length
	    SHORT       : {value : 4, type  : 'Short'},    //16bit
	    INT         : {value : 5, type  : 'Integer'},  //32bit signed integer
	    LONG        : {value : 6, type  : 'Long'},     //64bit
	    FLOAT       : {value : 7, type  : 'Float'},    //32bit float
	    DOUBLE      : {value : 8, type  : 'Double'},   //64bit float
	    NUMBER      : {value : 9, type  : 'Number'},   //Number with precision
	    STRING      : {value : 10, type : 'String'},   //Varying length String < 64K
	    DATE        : {value : 11, type : 'Date'},     //Only the Date part of the DateTime
	    TIME        : {value : 12, type : 'Time'},     //Time
	    TIMESTAMP   : {value : 13, type : 'Timestamp'},//64bit Timestamp - engine time - upto nanosecond precision, if the OS provides //SS:TODO
	    CLOB        : {value : 14, type : 'Clob'},     //Character -UTF-8 encoded string or large length > 64K
	    BLOB        : {value : 15, type : 'Blob'},  //Binary object - a stream of octets (unsigned 8bit char) with length. A variation of such Blobs could
	    fromTypeId : function(typeid) {
	    	switch(typeid) {
	    		case this.BOOLEAN.value:   return this.BOOLEAN;
	    		case this.BYTE.value:      return this.BYTE;
	    		case this.CHAR.value:      return this.CHAR;
	    		case this.SHORT.value:     return this.SHORT;
	    		case this.INT.value:       return this.INT;
	    		case this.LONG.value:      return this.LONG;
	    		case this.FLOAT.value:     return this.FLOAT;
	    		case this.DOUBLE.value:    return this.DOUBLE;
	    		case this.NUMBER.value:    return this.NUMBER;
	    		case this.STRING.value:    return this.STRING;
	    		case this.DATE.value:      return this.DATE;
	    		case this.TIME.value:      return this.TIME;
	    		case this.TIMESTAMP.value: return this.TIMESTAMP;
	    		case this.CLOB.value:      return this.CLOB;
	    		case this.BLOB.value:      return this.BLOB;
	    	}
	        return this.INVALID;
	     },
		 fromTypeName : function(typename) {
			 //console.log('TGAttributeType::fromTypeName Not implemented yet!');
		 }
	};
