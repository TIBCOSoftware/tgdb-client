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

var PrintUtility              = require('../../utils/PrintUtility').PrintUtility,
    TGResultSet               = require('../../query/TGResultSet'),
    TGQuery                   = require('../../query/TGQuery').TGQuery,
    TGEdge                    = require('../../model/TGEdge').TGEdge,
    TGGraphObjectFactory      = require('../../model/TGGraphObjectFactory').TGGraphObjectFactory,
    TGAbstractEntity          = require('../../model/TGAbstractEntity').TGAbstractEntity,
	TGEntityManager           = require('../../model/TGEntityManager').TGEntityManager;


var globleConnIds = 0;
var entities = {};

// Class definition
function TGConnectionTestImpl(connectionPool, channel, properties) {
	
	this._graphObjectFactory = new TGGraphObjectFactory(this);
	this._addedEntities      = {};
	this._updatedEntities    = {};
	this._removedEntities    = {};

	this._connId             = globleConnIds++;
    this._connectionPool     = connectionPool;
    this._properties         = properties;
}

TGConnectionTestImpl.command = {
	CREATE : 1,
	EXECUTE : 2,
	EXECUTEID : 3,
	CLOSE : 4
};

TGConnectionTestImpl.prototype.getGraphObjectFactory = function() { 
	return this._graphObjectFactory; 
};

/**
 * Connect to Graph DB server asynchronously. Upon successful handshake followed
 * by authentication, the callback function will be called to indicate
 * successful establishment of the connection.
 * <p>
 * Any operations relying on success of connection should be called only upon
 * responsestatus param in the callback function's value = true.
 * </p>
 * 
 * @param callback -
 *            Function. - signature function mycallback(responsestatus)
 */
TGConnectionTestImpl.prototype.connect = function(callback) {
    
    callback();
};

/**
 * Disconnect from the graph DB server.
 */
TGConnectionTestImpl.prototype.disconnect = function() {
};

TGConnectionTestImpl.prototype.getEntity = function (tgKey, properties, callback) {
};
    	

TGConnectionTestImpl.prototype.getEntities = function (tgKey, properties, callback) {	
	callback(entities);
};

// What is tgKey for ??????
TGConnectionTestImpl.prototype.insertEntity = function (tgEntity) {
    // Should be using the virtualId here because it's brand new
    this._addedEntities[tgEntity.getAttributeValue('key')] = tgEntity;
};

TGConnectionTestImpl.prototype.updateEntity = function (tgEntity) {
	this._updatedEntities[tgEntity.getAttributeValue('key')] = tgEntity;
};

TGConnectionTestImpl.prototype.deleteEntity = function (tgEntity) {
	////console.log("TGConnectionTestImpl.prototype.deleteEntity : " + tgEntity.getAttributeValue('key'));
	this._removedEntities[tgEntity.getAttributeValue('key')] = tgEntity;
};

/**
 * request metadata from Graph DB server.
 * 
 * @param callback -
 *            Consuming metadata from server.
 */
TGConnectionTestImpl.prototype.getGraphMetadata = function (callback) {
	callback();
};


/**
 * Commit transaction to Graph DB server.
 * 
 * @param callback -
 *            Convey status of commit operation.
 */
TGConnectionTestImpl.prototype.commit = function(callback) {
    var changeList = [];
    for (var key in this._removedEntities) {
    	delete entities[key];
    	this._removedEntities[key].markDeleted();
    	changeList.push(this._removedEntities[key]);
    }
	
    for(var key in this._updatedEntities) {
    	entities[key] = this._updatedEntities[key];
    	this._updatedEntities[key].resetModifiedAttributes();
    	changeList.push(this._updatedEntities[key]);
    }
    
    for(var key in this._addedEntities) {
    	entities[key] = this._addedEntities[key];
    	this._addedEntities[key].resetModifiedAttributes();
    	changeList.push(this._addedEntities[key]);
    }
    
    this._addedEntities = {};
    this._updatedEntities = {};
    this._removedEntities = {};
    
    ////console.log('Size of change list : ' + changeList.length);
    
	callback(changeList);
};

/**
 * Rollback transaction to Graph DB server.
 */
TGConnectionTestImpl.prototype.rollback = function() {
	
};

/****************************
 *                          *
 *        Query API         *
 *                          *
 ****************************/

TGConnectionTestImpl.prototype.createQuery = function (expr, callback) {

};

TGConnectionTestImpl.prototype.executeQuery = function (queryExpr, queryOption, callback) {	

};

TGConnectionTestImpl.prototype.executeQueryWithId = function (queryHashId, queryOption, callback) {

};

TGConnectionTestImpl.prototype.closeQuery = function (queryHashId, callback) {

};

exports.TGConnectionTestImpl = TGConnectionTestImpl;
