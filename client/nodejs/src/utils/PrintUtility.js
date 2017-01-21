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

var TGEntityKind     = require('../model/TGEntityKind').TGEntityKind,
    TGAbstractEntity = require('../model/TGAbstractEntity').TGAbstractEntity,
    TGLogManager     = require('../log/TGLogManager'),
    TGLogLevel       = require('../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();

var PrintUtility = {
	printEntityArray : function(entities) {
		logger.logDebug( 
				'{{{{{{{{{{{{{{{{{{{{ printEntityArray start }}}}}}}}}}}}}}}}}}}}}} %s',
				Object.keys(entities).length);
		for ( var key in entities) {
			logger.logDebug( 
					'Entity : key = %s, vid = %s', 
					key, entities[key].getId().getHexString());
			var attributes = entities[key].getAttributes();
			for ( var index in attributes) {
				logger.logDebug( 
						'                  Attrinute : %d - %s', 
						index, attributes[index].getValue() 
						/** +* attributes[index].getAttributeDescriptor().getName()*/);
			}
		}
		logger.logDebug( 
				'{{{{{{{{{{{{{{{{{{{{ printEntityArray end  }}}}}}}}}}}}}}}}}}}}}}');
	},
	printEntityMap : function(entities, type) {
		logger.logDebug( 
				'{{{{{{{{{{{{{{{{{{{{ printEntityMap start }}}}}}}}}}}}}}}}}}}}}} %d',
				Object.keys(entities).length);
		for ( var key in entities) {
			logger.logDebug( 
					'%s - Entity : key = %s, vid = %s', 
					type, key, entities[key].getId().getHexString());
			var attributes = entities[key].getAttributes();
			for ( var index in attributes) {
				logger.logDebug( 
						'                  Attrinute : %s - %s', 
						index, attributes[index].getValue() /* * + * attributes[index].getAttributeDescriptor().getName()*/);
			}
		}
		logger.logDebug( 
				'{{{{{{{{{{{{{{{{{{{{ printEntityMap end  }}}}}}}}}}}}}}}}}}}}}}');
	},
    printEntitiesBreadth : function(node, maxDepth) {
    	logger.logDebug( 
    			'{{{{{{{{{{{{{{{{{{{{ printEntitiesBreadth start }}}}}}}}}}}}}}}}}}}}}} %s',
				node.getId().getHexString());
        var indent = '';
        var traverseMap = {};
        var nodeList = [];
        nodeList.push(node);

        var level = 1;

        while (nodeList.length > 0) {
            var listSize = nodeList.length;
            while (listSize > 0) {
                listSize--;
                node = nodeList.shift();
                
                if (!(!traverseMap[node.getId().getHexString()])) {
                	logger.logDebug( 
                			"%s%d:Node(%s) visited", indent, level, node.getId().getHexString());
                    continue;
                } 
                logger.logDebug( "%s%d:Node(%s)", indent, level, node.getId().getHexString());
                for (var index in node.getAttributes()) {
                	logger.logDebug( 
                			"%s %d:Attr : %s", indent, level, node.getAttributes()[index].getValue());
                }
                traverseMap[node.getId().getHexString()] = node;
                logger.logDebug( "%s %d:has %d edges", indent, level, node.getEdges().length);
                for (var index in node.getEdges()) {
                	var edge = node.getEdges()[index];
                    var nodes = edge.getVertices();
                    if (nodes.length > 0) {
                        if (nodes[0] === null) {
                        	logger.logDebug( "%s   %d:passed last edge\n", indent, level);
                        	continue;
                        }
                        for (var j=0; j<nodes.length; j++) {
                            if (nodes[j].getId().getHexString() === node.getId().getHexString()) {
                                continue;
                            } else {
                            	var attrs = edge.getAttributes();
                            	if (attrs.length>0) {
                            		var attrib = attrs[0];
                            		logger.logDebug( 
                            				"%s   %d:edge(%d)(%s) with node(%s)", 
                            				indent, level, edge.getId().getHexString(), 
                            				attrib.getValue(), nodes[j].getId().getHexString());
                            	} else {
                            		logger.logDebug( 
                            				"%s   %d:edge(%d) with node(%s)", 
                            				indent, level, edge.getId().getHexString(), 
                            				nodes[j].getId().getHexString());
                            	}
                                nodeList.push(nodes[j]);
                            }
                        }
                    }
                }
            }
            indent += "    ";
            level++;
        }
        logger.logDebug( 
        		'{{{{{{{{{{{{{{{{{{{{ printEntitiesBreadth end  }}}}}}}}}}}}}}}}}}}}}}');
    },
    printCommitResponseEntity : function(commitresponse) {
    	logger.logDebug( 
    			'{{{{{{{{{{{{{{{{{{{{ printCommitResponseEntity start }}}}}}}}}}}}}}}}}}}}}}');
    	var gof = this.getGraphObjectFactory();
    	var entityStream = commitresponse.getEntityStream();
    	var fetchedEntities = null;
        var count = entityStream.readInt();
        if (count > 0) {
         	fetchedEntities = {};
        }
        entityStream.setReferenceMap(fetchedEntities);
        for (var i=0; i<count; i++) {
        	var kind = TGAbstractEntity.TGEntityKind.fromValue(entityStream.readByte());
          	if (kind !== TGAbstractEntity.TGEntityKind.InvalidKind) {
           		var id = entityStream.readLong();
           		var entity = fetchedEntities[id];
           		if (kind === TGAbstractEntity.TGEntityKind.Node) {
           			//Need to put shell object into hashmap to be deserialized later
           			var node = entity;
           			if (node === null) {
           				node = gof.createNode();
           				entity = node;
           				fetchedEntities[node.getEntityId()] = node;
           			}
           			node.readExternal(entityStream);
           		} else if (kind === TGAbstractEntity.TGEntityKind.Edge) {
           			var edge = entity;
           			if (edge === null) {
           				edge = gof.createEdge(null, null, TGEdge.DirectionType.BiDirectional);
           				entity = edge;
           				fetchedEntities[edge.getEntityId()] = edge;
           			}
           			edge.readExternal(entityStream);
           		}
           		logger.logDebug( 
           				"Kind : %d, Id : %d", entity.getEntityKind().kind(), entity.getEntityId());
           		var attribs = entity.getAttributes();
    		    for (var key in attribs) {
    		    	logger.logDebug( "Attr : %s", attribs[key].getValue());
    		    }
    		    if (entity.getEntityKind() === TGAbstractEntity.TGEntityKind.Node) {
    		    	var edges = entity.getEdges();
    		    	for (var index in edges) {
    		    		logger.logDebug( "    Edge : %d", edges[index].getEntityId());
    		    	}
    		    } if (entity.getEntityKind() === TGAbstractEntity.TGEntityKind.Edge) {
    		    	var nodes = entity.getVertices();
    		    	for (var index in nodes) {
    		    		logger.logDebug( "    Node : %d", nodes[index].getEntityId());
    		    	}
    		    }
           	} else {
           		//console.log("Received invalid entity kind %d", kind);
           	}
        }
        logger.logDebug( 
        		'{{{{{{{{{{{{{{{{{{{{ printCommitResponseEntity end  }}}}}}}}}}}}}}}}}}}}}}');
    },
    printEntities : function(ent, maxDepth, currDepth, indent, traverseMap) {
    	logger.logDebug( 
    			'{{{{{{{{{{{{{{{{{{{{ printEntities start }}}}}}}}}}}}}}}}}}}}}}');
    	if(ent) {
    	    logger.logDebug('Size of enity list : ' + Object.keys(ent).length);
    	    printEntities(ent, maxDepth, currDepth, indent, traverseMap);
    	} else {
    	    logger.logDebug('Null entity list !!!!');
    	}
        logger.logDebug(
        		'{{{{{{{{{{{{{{{{{{{{ printEntities end  }}}}}}}}}}}}}}}}}}}}}}');
    }
};

function printEntities(ent, maxDepth, currDepth, indent, traverseMap) {
    if (currDepth === maxDepth) {
    	return;
    }
	if (!ent) {
		return;
	}
	
	var entityId = ent.getId().getHexString();
	
	if (!(!traverseMap[entityId])) {
		logger.logDebug( "%sKind : %s, entityId : %s visited", indent, 
			(ent.getEntityKind() === TGEntityKind.NODE ? "Node" : "Edge"), entityId);
		return;
	} 
    traverseMap[entityId] = ent;
    logger.logDebug( "%sKind : %s, entityId : %s", indent, 
   		(ent.getEntityKind() === TGEntityKind.NODE ? "Node" : "Edge"), entityId);
    for (var index in ent.getAttributes()) {
    	logger.logDebug( "%s Attr : %s", indent, ent.getAttributes()[index].getValue());
    }
    if (ent.getEntityKind() === TGEntityKind.NODE) {
    	var edgeCount = ent.getEdges().length;
    	logger.logDebug( "%s Has %d edges", indent, edgeCount);
    	for (var index in ent.getEdges()) {
    		printEntities(ent.getEdges()[index], maxDepth, currDepth, indent.concat(" "), traverseMap);
    	}
    } else if (ent.getEntityKind() === TGEntityKind.EDGE) {
    	var nodes = ent.getVertices();
    	if (nodes.length > 0) {
    		logger.logDebug( "%s Has end nodes", indent);
    	}
    	++currDepth;
    	for (var j=0; j<nodes.length; j++) {
    		printEntities(nodes[j], maxDepth, currDepth, indent.concat("  "), traverseMap);
    	}
    }
}

module.exports = PrintUtility;