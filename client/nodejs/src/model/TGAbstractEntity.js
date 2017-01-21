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

//var globalVirtualId = 0;

var TGAttribute           = require('./TGAttribute').TGAttribute,
	TGAttributeType       = require('./TGAttributeType').TGAttributeType,
	TGAttributeDescriptor = require('./TGAttributeDescriptor').TGAttributeDescriptor,
	TGEntitySequencer     = require('./TGEntitySequencer').TGEntitySequencer,
	TGEntityId            = require('./TGEntityId'),
    TGException           = require('../exception/TGException').TGException,
    TGLogManager          = require('../log/TGLogManager'),
    TGLogLevel            = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function TGAbstractEntity(graphMetadata, entityType) {
    var _graphMetadata = graphMetadata;
    var _attributes    = {};
    var _modifiedAttribute = [];
    
    //1 represents true and 0 is false
    var _entityType     = entityType;
    var _version       = 0;
    var _isNew         = true;
    var _isDeleted     = false; // Need to set this to true once the entity is confirmed deleted by the server
    var _isInitialized  = true;
    
    var _id = TGEntityId.createId();
    
    this.getId = function() {
    	return _id;
    };
    
    this.TGEntityKind = {
    		INVALIDKIND : {value :  0},
    	    ENTITY      : {value : 1},
    	    NODE        : {value : 2},
    	    EDGE        : {value : 3},
    	    GRAPH       : {value : 4},
    	    HYPEREDGE   : {value : 5},
    	    fromValue   : function (kind) {
            	switch(kind) {
            		case 0 : return this.INVALIDKIND;
            		case 1 : return this.ENTITY;
            		case 2 : return this.NODE;
            		case 3 : return this.EDGE;
            		case 4 : return this.GRAPH;
            		case 5 : return this.HYPEREDGE;
    	    }
    	}
    };
    
    this.getGraphMetadata = function() {
    	return _graphMetadata;
    };
    
    this.isNew = function () {
    	//console.log("TGAbstractEntity.isNew = " + _isNew);
    	return _isNew;
    };

    /**
     * Set entityId.
     * @param name
     */
    this.setIsNew = function(isNew) {
    	_isNew = isNew;
    };
    
    this.getVersion = function() {
        return _version;
    };

    this.setVersion = function(version) {
    	_version = version;
    };

    this.markDeleted = function() {
    	_isDeleted = true;
    };

    this.isDeleted = function() {
    	return _isDeleted;
    };
    
    this.setEntityType = function(entityType) {
    	_entityType = entityType;
    };
    
    this.getEntityType = function() {
    	return _entityType;
    };
    
    /**
     * Get attribute value for attribute with given name.
     * @param name
     */
    this.getAttribute = function(name) {
        return _attributes[name];
    };
    
    this.getAttributeKeys = function() {
    	return Object.keys(_attributes);
    };
    
    /**
     * Get attribute value for attribute with given name.
     * @param name
     */
    this.getAttributeValue = function(name) {
    	if (_attributes[name]) {
            return _attributes[name].getValue();
    	}
    };
    
    /**
     * Return attributes for this entity.
     */
    this.getAttributes = function() {
    	// convert a map to an array.
        return Object.keys(_attributes).map(function(key){
            return _attributes[key];
        });
    };

    /**
     * Check if attribute with given name is set.
     * @param name
     */
    this.isAttributeSet = function(name) {
        return this.getAttribute(name) !== null;
    };

    /**
     * Set attribute on the node.
     * @param name
     */
    this.getAttribute = function(name) {
        return _attributes[name];
    };

    this.setBooleanAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.BOOLEAN);
    };
    
    this.setByteAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.BYTE);
    };
    
    this.setCharAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.CHAR);
    };
    
    this.setShortAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.SHORT);
    };
    
    this.setIntegerAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.INT);
    };

    this.setLongAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.LONG);
    };

    this.setFloatAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.FLOAT);
    };

    this.setDoubleAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.DOUBLE);
    };

    this.setNumberAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.NUMBER);
    };

    this.setStringAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.STRING);
    };
    
    this.setDateAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.DATE);
    };

    this.setDateTimeAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.DATETIME);
    };

    this.setTimestampAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.TIMESTAMP);
    };
    
    this.setClobAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.CLOB);
    };
    
    this.setBlobAttribute = function(name, value) {
    	this.setAttribute(name, value, TGAttributeType.BLOB);
    };
    
    /**
     * Set attribute on the node.
     * @param name
     * @param value
     */
    this.setAttribute = function(name, value, type) {
    	logger.logDebugWire("SetAttribute : name = %s, value = %s", name, value);
    	
    	if(!name) {
    		throw new TGException("Name of the attribute cannot be null");
    	}
    	
    	var attr = _attributes[name];    	
        if (!attr) {
        	var attrDesc = _graphMetadata.getAttributeDescriptor(name);
            if(!attrDesc) {
            	if ((value===null || typeof(value)==="undefined") && !type) {
            		//If the attribute has not been set and has no descriptor, it cannot have a null value
            		//because we cannot figure out the data type
                	throw new TGException('Null value specified for an undefined attribuyte : ' + name + ', value : ' + value);
            	} else {
            		console.log(name + " is a new attribute type");
                    attrDesc = _graphMetadata.createAttributeDescriptor(name, type, false);
            	}
            }
            
            attr = new TGAttribute(this, attrDesc, value);
        }
        
        if (!attr.isModified()) {
        	_modifiedAttribute.push(attr);
        }
        attr.setValue(value);
        _attributes[name] = attr;
    };

    /**
     * Return a byte buffer.
     */
    this.writeAttributes = function(outputStream) {
    	logger.logDebugWire(
    			"**** Enering TGAbstractEntity.writeAttributes at output buffer position at : %s %d", 
    			this.getId().getHexString(), outputStream.getPosition());
        outputStream.writeBoolean(_isNew);
        outputStream.writeByte(this.getEntityKind().value);
        outputStream.writeLongAsBytes(this.getId().getBytes());   //outputStream.writeLong(this.getVirtualId());
        outputStream.writeInt(_version);
        outputStream.writeInt(!_entityType ? 0 : _entityType.getId());

        var modifiedCount = 0;
        
        for (var attName in _attributes) {
        	if(_attributes[attName].isModified()) {
        		modifiedCount++;
        	}
        }
        outputStream.writeInt(modifiedCount);
        
        for (attName in _attributes) {
        	if(_attributes[attName].isModified()) {
        		_attributes[attName].writeExternal(outputStream);
        	}
        }
        logger.logDebugWire(
        		"****Leaving TGAbstractEntity.writeAttributes at output buffer position at : %s %d", 
        		this.getId().getHexString(), outputStream.getPosition());
    };

    this.readAttributes = function(inputStream) {
    	logger.logDebugWire(
    			"**** Enering TGAbstractEntity.readAttributes at output buffer position at : %s %d", 
    			this.getId().getHexString(), inputStream.getPosition());
    	_isNew = inputStream.readBoolean();  //Should always be False.
        if (_isNew === true) {
        	//console.log("Deserializing a new entity is not expected");
        	_isNew = false;
        }
        
        var kind = inputStream.readByte();
        if (this.getEntityKind().value !== kind) {
        	throw new TGException("Invalid object for deserialization. Expecting..."); 
        }

        this.getId().setBytes(inputStream.readLongAsBytes());  //Overwrite
        _version = inputStream.readInt();
        var entityTypeId = inputStream.readInt();
        if (entityTypeId !== 0) {
        	_entityType = _graphMetadata.getNodeType(entityTypeId);
        	if (!_entityType) {
        		//FIXME: retrieve entity type together with the entity?
        		//console.log("Cannot lookup entity type %d from graph meta data cache", entityTypeId);
        	}
        }

        var count = inputStream.readInt();
        for (var i = 0; i <count; i++) {
            var attr = new TGAttribute(this);
            attr.readExternal(inputStream);
            if (attr && attr.getAttributeDescriptor()) {
                _attributes[attr.getAttributeDescriptor().getName()] = attr;
            }
        }
        logger.logDebugWire(
        		"****Leaving TGAbstractEntity.readAttributes at output buffer position at : %s %d", 
        		this.getId().getHexString(), inputStream.getPosition());
    };

    this.resetModifiedAttributes = function () {
    	for (var index in _modifiedAttribute) {
    		_modifiedAttribute[index].resetIsModified();
    	}
    	_modifiedAttribute.length = 0;
    };
    
    this.setInitialized = function (isInit) {
    	_isInitialized = isInit;
    };

    this.isInitialized = function () {
    	return _isInitialized;
    };
}

exports.TGAbstractEntity = TGAbstractEntity;