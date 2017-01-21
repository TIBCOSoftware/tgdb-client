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

var	TGGraphObjectFactory = require('./TGGraphObjectFactory').TGGraphObjectFactory,
    TGAbstractEntity     = require('./TGAbstractEntity').TGAbstractEntity,
    TGEntityKind         = require('./TGEntityKind').TGEntityKind,
    TGException          = require('../exception/TGException').TGException,
    PrintUtility         = require('../utils/PrintUtility'),
    TGLogManager         = require('../log/TGLogManager'),
    TGLogLevel           = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

function checkEntity(entity) {
	if(entity.getEntityKind() === TGEntityKind.EDGE) {
		var nodes = entity.getVertices();
		nodes.forEach(function(node){
			if(!node) {
				throw new TGException('A undefined node found for an edge.');
			}
		});
	}
}

function TGEntityManager() {
	
    var _graphObjectFactory = new TGGraphObjectFactory(this);
    var _addedEntities      = {};
    var _updatedEntities    = {};
    var _removedEntities    = {};
    
    this.getGraphObjectFactory = function() { return _graphObjectFactory; };
    
    this.addedEntities = function() { return _addedEntities; };

    this.updatedEntities = function() { return _updatedEntities; };

    this.removedEntities = function() { return _removedEntities; };
    
    this.newAttrDecsSet = function() {
        return _graphObjectFactory.getGraphMetaData().getNewAttributeDescriptors(); 
    };
    
    this.getCachedGraphMetaData = function() {
		return _graphObjectFactory.getGraphMetaData();
    };
    
    this.updateGraphMetaData = function(attrDescList, nodeTypeList, edgeTypeList) {
		var graphMetadata = _graphObjectFactory.getGraphMetaData();
		graphMetadata.updateMetadata(attrDescList, nodeTypeList, edgeTypeList);

		return graphMetadata; 
    };
    
    /**
     * New Entity Created.
     * 
     * @param entity
     */
    this.entityCreated = function(tgEntity) {
    	var attributes = tgEntity.getAttributes();
    	var attKeys = Object.keys(attributes);
    	logger.logDebugWire(
    		'Entity is created - Id = %s - %s', tgEntity.getId(),  
    		(0!==attKeys.length ? attributes[attKeys[0]].getValue() : "no attribute found"));
        // Should be using the virtualId here because it's brand new
        _addedEntities[tgEntity.getId().getHexString()] = tgEntity;    	
    };

    /**
     * Called when a Node is added to Graph
     * 
     * @param graph
     * @param node
     */
    this.nodeAdded = function(tgGraph, tgNode) {
    	var attributes = tgNode.getAttributes();
    	var attKeys = Object.keys(attributes);
    	logger.logDebugWire(
    		'tgNode is added - Id = %s - %s', tgNode.getId(),  
    		(0!==attKeys.length ? attributes[attKeys[0]].getValue() : "no attribute found"));
        _addedEntities[tgNode.getId().getHexString()] = tgNode;
    };

    /**
     * Called when an attribute is Added to an entity.
     * 
     * @param attribute
     */
    this.attributeAdded = function(tgAttribute, tgEntityOwner) {
        //console.log("Attribute is created");
    };

    this.entityUpdated = function(tgEntity) {
    	var attributes = tgEntity.getAttributes();
    	var attKeys = Object.keys(attributes);
    	logger.logDebugWire(
    		'Entity is updated - Id = %s - %s', tgEntity.getId(),  
    		(0!==attKeys.length ? attributes[attKeys[0]].getValue() : "no attribute found"));
        // Should be using the virtualId here because it's brand new
        _updatedEntities[tgEntity.getId().getHexString()] = tgEntity;    	
    };
    
    this.updateChangedListForNewEdge = function () {
        //Include existing nodes to the changed list if it's part of a new edge
    	
        Object.keys(_addedEntities).forEach(function(key) {
			var entity = _addedEntities[key];
        	if (entity.getEntityKind() === TGEntityKind.EDGE) {
        		// Since it's an edge, we can call getVertices to get nodes which asscociated with this edge.
        		entity.getVertices().forEach(function(node){
        			if (!node.isNew()) {
        				_updatedEntities[node.getId().getHexString()] = node;
        				logger.logDebugWire(
        						"New edge added to an existing node : %s", 
        						node.getId().getHexString());
        			}
        		});
        	}
        });
    };
    
    this.updateChangedListForUpdatedEdge = function () {
        Object.keys(_updatedEntities).forEach(function(key) {
			var entity = _updatedEntities[key];
        	if (entity.getEntityKind() === TGEntityKind.EDGE) {
        		// Since it's an edge, we can call getVertices to get nodes which asscociated with this edge.
        		entity.getVertices().forEach(function(node) {
        			if (!node.isNew()) {
        				_updatedEntities[node.getId().getHexString()] = node;
        				logger.logDebugWire(
        						"Edge updated for an existing node : %s", 
        						node.getId().getHexString());
        			}
        		});
        	}
        });
    };
    
    this.updateRemovedListForRemovesEdge = function () {
        Object.keys(_removedEntities).forEach(function(key) {
			var entity = _removedEntities[key];
        	if (entity.getEntityKind() === TGEntityKind.EDGE) {
        		// Since it's an edge, we can call getVertices to get nodes which asscociated with this edge.
        		entity.getVertices().forEach(function(node) {
        			if (!node.isNew()) {
        				if(!_removedEntities[node.getId().getHexString()]) {
        					_updatedEntities[node.getId().getHexString()] = node;
        					logger.logDebugWire(
        							"Edge removed but node is not : %s", 
        							node.getId().getHexString());
        				}
        			}
        		});
        	}
        });
    };    
    
    /**
     * Called when an attribute is set.
     * 
     * @param attribute
     * @param oldValue
     * @param newValue
     */
    this.attributeChanged = function(tgAttribute, oldValue,
    		newValue) {
        //console.log("Attribute is changed");
    };

    /**
     * Called when an attribute is removed from the entity.
     * 
     * @param attribute
     */
    this.attributeRemoved = function(tgAttribute,
    		tgEntityOwner) {
        //console.log("Attribute is removed");
    };

    /**
     * Called when an node is from Graph
     * 
     * @param graph
     * @param node
     */
    this.nodeRemoved = function(tgGraph, tgNode) {
        //console.log("Node is removed");
        _removedEntities[tgNode.getVirtualId()] = tgNode;
    };

    /**
     * Called when the entity is deleted.
     * 
     * @param entity
     */
    this.entityDeleted = function(tgEntity) {
    	var attributes = tgEntity.getAttributes();
    	var attKeys = Object.keys(attributes);
        logger.logDebugWire(
        		'Entity is deleted - Id = %s - %s', tgEntity.getId(),  
        		(0!==attKeys.length ? attributes[attKeys[0]].getValue() : "no attribute found"));
        _removedEntities[tgEntity.getId().getHexString()] = tgEntity;
    };
    
    this.updateEntityIds = function(response) {
    	logger.logDebugWire(
    			'(((((((((((((((((((((((((( updateEntityIds )))))))))))))))))))))))))))');
    	fixUpAttrDescriptors(response, _graphObjectFactory.getGraphMetaData().getNewAttributeDescriptors());
        fixUpEntityIds(response, _addedEntities, _updatedEntities);
    	
        var changeList = [];
        Object.keys(_removedEntities).forEach(function(key) {
        	_removedEntities[key].markDeleted();
        	changeList.push(_removedEntities[key]);
        });
        
        //for (var key in _removedEntities) {
        //	_removedEntities[key].markDeleted();
        //	changeList.push(_removedEntities[key]);
        //}
    	
        Object.keys(_updatedEntities).forEach(function(key) {
        	_updatedEntities[key].resetModifiedAttributes();
        	changeList.push(_updatedEntities[key]);
        });
        
        //for(var key in _updatedEntities) {
        //	_updatedEntities[key].resetModifiedAttributes();
        //	changeList.push(_updatedEntities[key]);
        //}
        
        Object.keys(_addedEntities).forEach(function(key) {
        	_addedEntities[key].resetModifiedAttributes();
        	changeList.push(_addedEntities[key]);
        });
        
        //for(var key in _addedEntities) {
        //	_addedEntities[key].resetModifiedAttributes();
        //	changeList.push(_addedEntities[key]);
        //}
        
        this.clear();
        return changeList;
    };
    
    this.clear = function (){
    	
        if(logger.isDebug) {
        	logger.logDebug(
			    '[TGEntityManager.clear] Before clean up !!!!');
            PrintUtility.printEntityMap(_addedEntities, 'addedEntities');
            PrintUtility.printEntityMap(_updatedEntities, 'updatedEntities');
            PrintUtility.printEntityMap(_removedEntities, 'removedEntities');        	
        }
        
        _addedEntities = {};
        _updatedEntities = {};
        _removedEntities = {};
        
        if(logger.isDebug) {
        	logger.logDebug(
		        '[TGEntityManager.clear] After clean up !!!!');
            PrintUtility.printEntityMap(_addedEntities, 'addedEntities');
            PrintUtility.printEntityMap(_updatedEntities, 'updatedEntities');
            PrintUtility.printEntityMap(_removedEntities, 'removedEntities');        	
        }
    };
}

/*
 *  Private methods
 * */

function fixUpAttrDescriptors (response, attrDescSet) {
	logger.logDebugWire( 
			"[TGEntityManager::fixUpAttrDescriptors] Fixup attribute descriptor ids");
    var attrDescCount = response.getAttrDescCount();
    var attrDescIdList = response.getAttrDescIdList();
    for (var i=0; i<attrDescCount; i++) {
    	var tempId = attrDescIdList[i*2]; 
    	var realId = attrDescIdList[(i*2) + 1];
    	
		for(var index in attrDescSet) {
    		if (attrDescSet[index].getAttributeId() === tempId) {
    			logger.logDebugWire( 
    					"[TGEntityManager::fixUpAttrDescriptors] Replace descriptor id : %d by %d", 
    					attrDescSet[index].getAttributeId(), realId);
    			attrDescSet[index].setAttributeId(realId);
    			break;
    		}
		}
    }
}

function fixUpEntityIds (response, addedEntities, updatedEntities) {
	logger.logDebugWire( 
			"[TGEntityManager::fixUpEntityIds] Fixup entity ids");
    var addedIdCount = response.getAddedEntityCount();
    var addedIdList = response.getAddedIdList();
    var key = null;
    var version = null;
    for (var i=0; i<addedIdCount; i++) {
    	var tempId = addedIdList[i*3]; 
    	var realId = addedIdList[(i*3) + 1];
    	version = addedIdList[(i*3) + 2];
		logger.logDebugWire( 
				"[TGEntityManager::fixUpEntityIds::added] Try to replace entity id : %s by %s, version %s", 
				tempId, realId, version);
		for(key in addedEntities) {	
    		if (addedEntities[key].getId().equalsToBytes(tempId)) {
    			logger.logDebugWire( 
    				"[TGEntityManager::fixUpEntityIds::added] Replace entity id : %s by %s", 
    				tempId, realId);
    			addedEntities[key].getId().setBytes(realId);
    			addedEntities[key].setIsNew(false);
    			addedEntities[key].setVersion(version);
    		}
    	}
    }

    var updatedIdCount = response.getUpdatedEntityCount();
    var updatedIdList = response.getUpdatedIdList();
    for (i=0; i<updatedIdCount; i++) {
    	var id = updatedIdList[i*2]; 
    	version = updatedIdList[(i*2) + 1];
		logger.logDebugWire( 
				"[TGEntityManager::fixUpEntityIds::updated] Try to update for id : %s, version %s", 
				id, version);
   		for(key in updatedEntities) {
			logger.logDebugWire( 
    				"[TGEntityManager::fixUpEntityIds::updated] Replace entity id : %s, version %s", 
    				id, version);
   			if (updatedEntities[key].getId().equalsToBytes(id)) {
   				updatedEntities[key].setVersion(version);
   			    break;
   			}
   		}
    }
}

exports.TGEntityManager = TGEntityManager;