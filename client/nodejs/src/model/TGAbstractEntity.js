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

var globalVirtualId = 0;

var TGAttribute = require('./TGAttribute').TGAttribute,
	TGAttributeDescriptor = require('./TGAttributeDescriptor').TGAttributeDescriptor,
	TGEntitySequencer = require('./TGEntitySequencer').TGEntitySequencer;

//Class definition
function TGAbstractEntity(graphMetadata) {
    var _graphMetadata = graphMetadata;
    var _attributes    = {};
    var _modifiedAttributeName = [];
    
    //1 represents true and 0 is false
    var _entityId      = -1;
    var _version       = 0;
    var _isNew         = true;
    var _virtualId     = --globalVirtualId;
    
    this.isNew = function () {
    	console.log("TGAbstractEntity.isNew = " + _isNew);
    	return _isNew;
    }

    /**
     * Set entityId.
     * @param name
     */
    this.setIsNew = function(isNew) {
    	_isNew = isNew;
    };
    
    /**
     * Get virtualId.
     * @param name
     */
    this.getVirtualId = function() {
        return _isNew ? _virtualId : _entityId;
    };

    /**
     * Set entityId.
     * @param name
     */
    this.setEntityId = function(entityId) {
    	_virtualId = 0;
    	_entityId = entityId;
    };
    
    /**
     * Get attribute value for attribute with given name.
     * @param name
     */
    this.getAttribute = function(name) {
        return _attributes[name];
    };

    /**
     * Return attributes for this entity.
     */
    this.getAttributes = function() {
        return Object.keys(_attributes).map(function(key){
            return _attributes[key];
        });
    };

    /**
     * Check if attribute with given name is set.
     * @param name
     */
    this.isAttributeSet = function(name) {
        return this.getAttribute(name) != null;
    };

    /**
     * Set attribute on the node.
     * @param name
     */
    this.getAttribute = function(name) {
        return _attributes[name];
    };

    /**
     * Set attribute on the node.
     * @param name
     * @param value
     */
    this.setAttribute = function(name, value, type) {
    	console.log("SetAttribute : " + name + ", " + value);
    	
    	var attr = _attributes[name];
    	if(!attr && (value===null || typeof(value) == "undefined") && !type) {
    		throw new Error("Unable to determine type of value!");
    	}
    	
        if (!attr) {
        	var attrDesc = _graphMetadata.getAttributeDescriptor(name);
            if (!attrDesc) {
                attrDesc = _graphMetadata.createAttributeDescriptor(name, type, false);
            }
            
            if(!attrDesc) {
            	throw new Error('Unable to create attribute descriptor for name : ' + name + ', value : ' + value);
            }
            
            attr = new TGAttribute(this, attrDesc, value);
            _attributes[name] = attr;
        }
    	
    	attr.setValue(value);
    	_modifiedAttributeName.push(name);
    };

    /**
     * Return a byte buffer.
     */
    this.writeAttributes = function(outputStream) {
    	console.log("****Enering TGAbstractEntity.writeAttributes at output buffer position at : %s %d", _virtualId, outputStream.getPosition());
        outputStream.writeBoolean(_isNew);
        outputStream.writeByte(this.getEntityKind());
        outputStream.writeLong(this.getVirtualId());
        outputStream.writeInt(_version);
        
        var modifiedCount = 0;
        
        for (var attName in _attributes) {
        	if(_attributes[attName].isModified()) {
        		modifiedCount++;
        	}
        }
        outputStream.writeInt(modifiedCount);
        
        for (attName in _attributes) {
        	if(_attributes[attName].isModified()) {
        		outputStream.writeInt(_attributes[attName].getAttributeDescriptor().getAttributeId());
        		outputStream.writeBoolean(_attributes[attName].isNull());
        		_attributes[attName].writeExternal(outputStream);
        	}
        }
       	console.log("****Leaving TGAbstractEntity.writeAttributes at output buffer position at : %s %d", _virtualId, outputStream.getPosition());
    };

    this.readExternal = function(inputStream) {
        var kind = inputStream.readByte();
        if (this.getEntityKind() != kind) {
        	throw new Error("Invalid object for deserialization. Expecting..."); 
        }

        _isNew = false;
        _entityId = inputStream.readLong();  //Overwrite
        _version = inputStream.readInt();

        var count = inputStream.readInt();
        for (var i = 0; i <count; i++) {
            var attr = new TGAttribute(this);
            attr.readExternal(inputStream);
            setAttribute(attr);
        }
    };

    this.resetModifiedAttributes = function () {
    	for (var key in _modifiedAttributeName) {
    		_attributes[_modifiedAttributeName[key]].resetIsModified();
    	}
    	_modifiedAttributeName.length = 0;
    };
}

TGAbstractEntity.EntityKind = {
		INVALIDKIND : 0,
	    ATTRIBUTE_DESCRIPTOR : 1,
	    NODETYPE : 2,
	    EDGETYPE : 3,
	    ATTRIBUTE : 4,
	    NODE : 5,
	    EDGE : 6,
	    GRAPH : 7,
	    HYPEREDGE : 8
};

exports.TGAbstractEntity = TGAbstractEntity;