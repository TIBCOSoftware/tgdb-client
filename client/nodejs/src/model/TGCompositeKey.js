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

var TGAttribute = require('./TGAttribute').TGAttribute,
    TGKey       = require('./TGKey').TGKey,
    util        = require('util');

function TGCompositeKey (graphMetadata, typeName) {

    this._graphMetadata = graphMetadata;
    this._typeName = typeName;
    this._attributes = {};
/*  Not require to have a type name
    if (this._graphMetadata.getNodeType(this._typeName) === null) {
    	throw Error("Invalid NodeType specified :" + typeName);
    }
*/
}

util.inherits(TGCompositeKey, TGKey);

TGCompositeKey.prototype.setAttribute = function(name, value, attributeType) {
	//console.log("Entering TGCompositeKey.prototype.setAttribute : type = " + this._typeName + ', ' + name + ' = ' + value);

	if (value === null) {
		throw new Error("Value is null");
	}

	var attr = this._attributes[name];
	if (!attr) {
		var attrDesc = this._graphMetadata.getAttributeDescriptor(name);
		if (attrDesc === null) {
	    	if(attributeType===null || typeof(attributeType)==="undefined") {
	    		/* Force type or throw exception? */
	    		//attributeType = TGAttributeType.STRING;
	    		throw new Error("Attribute type can not be determine");
	    	}
            attrDesc = this._graphMetadata.createAttributeDescriptor(name, attributeType);
		}

		attr = new TGAttribute(null, attrDesc, value); //There is no owner
	}
	//try {
		attr.setValue(value);
		this._attributes[name] = attr;
	//}
	//catch (Error) {
	//	throw {message:"Can't set value to attribute"};
	//}
};

TGCompositeKey.prototype.writeExternal = function (outputStream) {
    if (this._typeName !== null) {
    	outputStream.writeBoolean(true); //TypeName exists
    	outputStream.writeUTF(this._typeName);
    }
    else {
    	outputStream.writeBoolean(false);
    }
    outputStream.writeShort(Object.keys(this._attributes).length);
	for (var key in this._attributes) {
		//Null value is not allowed and therefore no need to include isNull flag
		this._attributes[key].writeExternal(outputStream);
	}
};

TGCompositeKey.prototype.readExternal = function (inputStream) {
	throw new Error("Not Supported operation");
};

exports.TGCompositeKey = TGCompositeKey;