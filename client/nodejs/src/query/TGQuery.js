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
var TGQueryOption = require('./TGQueryOption');

function TGQuery(connection, queryHashId) {
	this._connection = connection;
    this._queryHashId = queryHashId;
    this._option = TGQueryOption.DEFAULT_QUERY_OPTION;
    this._parameters = {};

    this.NULLPARAMETER = null;
}
    
TGQuery.prototype.setBoolean = function (name, value) {
	this._parameters[name] = value;
};

TGQuery.prototype.setChar = function (name, value) {    
	this._parameters[name] = value;
};

TGQuery.prototype.setShort = function (name, value) {
	this._parameters[name] = value;
};

TGQuery.prototype.setInt = function (name, value) {
	this._parameters[name] = value;
};

TGQuery.prototype.setLong = function (name, value) {
	this._parameters[name] = value;
};

TGQuery.prototype.setFloat = function (name, value) {
	this._parameters[name] = value;
};

TGQuery.prototype.setDouble = function (name, value) {
	this._parameters[name] = value;
};

TGQuery.prototype.setString = function (name,  value) {
	this._parameters[name] = value;
};
    
TGQuery.prototype.setDate = function (name, value) {
	this._parameters[name] = value;
};
    
TGQuery.prototype.setBytes = function (name, value) {
	this._parameters[name] = value;
};

TGQuery.prototype.setNull = function (name) {
	this._parameters[name] = this.NULLPARAMETER;
};

TGQuery.prototype.setOption = function(option) {
    this._option = option;
};

TGQuery.prototype.execute = function (callback) {
	this._connection.executeQueryWithId(this._queryHashId, callback);
};

TGQuery.prototype.close = function (callback) {
	this._connection.closeQuery(this._queryHashId, callback);
};

exports.TGQuery = TGQuery;