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

var	TGGraphObjectFactory = require('../model/TGGraphObjectFactory').TGGraphObjectFactory,
PrintUtility             = require('../utils/PrintUtility').PrintUtility;

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
    }
    
    /**
     * New Entity Created.
     * 
     * @param entity
     */
    this.entityCreated = function(tgEntity) {
    	var attributes = tgEntity.getAttributes();
    	var attKeys = Object.keys(attributes);
        console.log('Entity is created - virtualId = %s - %s', tgEntity.getVirtualId(),  (0!=attKeys.length ? attributes[attKeys[0]].getValue() : "no attribute found"));
        // Should be using the virtualId here because it's brand new
        _addedEntities[tgEntity.getVirtualId()] = tgEntity;    	
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
        console.log('tgNode is added - virtualId = %s - %s', tgNode.getVirtualId(),  (0!=attKeys.length ? attributes[attKeys[0]].getValue() : "no attribute found"));
        _addedEntities[tgNode.getVirtualId()] = tgNode;
    };

    /**
     * Called when an attribute is Added to an entity.
     * 
     * @param attribute
     */
    this.attributeAdded = function(tgAttribute, tgEntityOwner) {
        console.log("Attribute is created");
    };

    this.entityUpdated = function(tgEntity) {
    	var attributes = tgEntity.getAttributes();
    	var attKeys = Object.keys(attributes);
        console.log('Entity is created - virtualId = %s - %s', tgEntity.getVirtualId(),  (0!=attKeys.length ? attributes[attKeys[0]].getValue() : "no attribute found"));
        // Should be using the virtualId here because it's brand new
        _updatedEntities[tgEntity.getVirtualId()] = tgEntity;    	
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
        console.log("Attribute is changed");
    };

    /**
     * Called when an attribute is removed from the entity.
     * 
     * @param attribute
     */
    this.attributeRemoved = function(tgAttribute,
    		tgEntityOwner) {
        console.log("Attribute is removed");
    };

    /**
     * Called when an node is from Graph
     * 
     * @param graph
     * @param node
     */
    this.nodeRemoved = function(tgGraph, tgNode) {
        console.log("Node is removed");
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
        console.log('Entity is deleted - virtualId = %s - %s', tgEntity.getVirtualId(),  (0!=attKeys.length ? attributes[attKeys[0]].getValue() : "no attribute found"));
        _removedEntities[tgEntity.getVirtualId()] = tgEntity;
    };
    
    this.updateEntityIds = function(response) {
    	console.log('(((((((((((((((((((((((((( updateEntityIds )))))))))))))))))))))))))))');
        fixUpAttrDescIds(response, _graphObjectFactory.getGraphMetaData().getNewAttributeDescriptors());
        fixUpEntityIds(response, _addedEntities);
    }
    
    this.clear = function (){
        for(var key in _updatedEntities) {
        	_updatedEntities[key].resetModifiedAttributes();
        }
        
        for(var key in _addedEntities) {
        	_addedEntities[key].resetModifiedAttributes();
        }
        
        _addedEntities = {};
        _updatedEntities = {};
        _removedEntities = {};
        
        PrintUtility.printEntityMap(_addedEntities, 'addedEntities');
        PrintUtility.printEntityMap(_updatedEntities, 'updatedEntities');
        PrintUtility.printEntityMap(_removedEntities, 'removedEntities');
    };
}

function fixUpAttrDescIds (response, attrDescSet) {
    console.log("Fixup attribute descriptor ids");
    var attrDescCount = response.getAttrDescCount();
    var attrDescIdList = response.getAttrDescIdList();
    for (var i=0; i<attrDescCount; i++) {
    	var tempId = attrDescIdList[i*2]; 
    	var realId = attrDescIdList[(i*2) + 1];
    	
		for(var index in attrDescSet) {
    		if (attrDescSet[index].getAttributeId() == tempId) {
    			console.log("Replace descriptor id : %d by %d", attrDescSet[index].getAttributeId(), realId);
    			attrDescSet[index].setAttributeId(realId);
    		}
		}
    }
}

function fixUpEntityIds (response, addedList) {
    console.log("Fixup entity ids");
    var addedIdCount = response.getAddedEntityCount();
    var addedIdList = response.getAddedIdList();
    for (var i=0; i<addedIdCount; i++) {
    	var tempId = addedIdList[i*2]; 
    	var realId = addedIdList[(i*2) + 1];

		for(var key in addedList) {
			
//			console.log('---------------- key = ' + key + ', vid = ' + addedList[key].getVirtualId() + ', tempId = ' + tempId);
			
    		if (addedList[key].getVirtualId() == tempId) {
    			console.log("Replace entity id : %d by %d", tempId, realId);
    			addedList[key].setEntityId(realId);
    			addedList[key].setIsNew(false);
    		}
    	}
    }
//    PrintUtility.printEntityMap(addedList, 'FixedAddedEntities');

}

exports.TGEntityManager = TGEntityManager;