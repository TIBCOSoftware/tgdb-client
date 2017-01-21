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


function TGResultSet (conn, resultId) {
    this._conn = conn;
    this._resultId = resultId;
    this._resultList = [];
    this._isOpen = true;
    this._currPos = -1;
}

TGResultSet.prototype.isOpen = function () {
    return this._isOpen;
};

TGResultSet.prototype.addEntityToResultSet = function (entity) {
    this._resultList.push(entity);
};

TGResultSet.prototype.getResultId = function () {
    return this._resultId;
};

TGResultSet.prototype.hasNext = function () {
	if (this._isOpen === false) {
		return false;
	}
	if (this._currPos < (this._resultList.length - 1)) {
		return true;
	}
	return false;
};

TGResultSet.prototype.close = function() {
	this._isOpen = false;
};

TGResultSet.prototype.hasExceptions = function() {
	// TODO Auto-generated method stub
	return false;
};

TGResultSet.prototype.getExceptions = function() {
	// TODO Auto-generated method stub
	return null;
};

TGResultSet.prototype.count = function() {
	if (this._isOpen === false) {
		return 0;
	}
	return this._resultList.length;
};

TGResultSet.prototype.first = function() {
	if (this._isOpen === false) {
		return null;
	}
	this._currPos = 0;
	if (this._resultList.length === 0) {
		return null;
	}
	return this._resultList[this._currPos];
};

TGResultSet.prototype.last = function() {
	if (this._isOpen === false) {
		return null;
	}
	this._currPos = this._resultList.length - 1;
	if (this._resultList.length === 0) {
		this._currPos = 0;
		return null;
	}
	return this._resultList[this._currPos];
};

TGResultSet.prototype.prev = function() {
	if (this._isOpen === false) {
		return null;
	}
	if (this._currPos > 0) {
		this._currPos--;
		return this._resultList[this._currPos];
	}
	return null;
};

TGResultSet.prototype.next = function() {
	if (this._isOpen === false) {
		return null;
	}
	if (this._currPos < (this._resultList.length - 1)) {
		this._currPos++;
		return this._resultList[this._currPos];
	}
	return null;
};

TGResultSet.prototype.getPosition = function() {
	if (this._isOpen === false) {
		return 0;
	}
	return this._currPos;
};

TGResultSet.prototype.getAt = function(position) {
	if (this._isOpen === false) {
		return null;
	}
	if (position >= 0 && position < this._resultList.length) {
		return this._resultList[position];
	}
	return null;
};

TGResultSet.prototype.skip = function(position) {
	if (this._isOpen === false) {
		return;
	}
	var newPos = this._currPos + position;
	if (newPos >=0 && newPos < this._resultList.length) {
		this._currPos = newPos;
	}
};

module.exports = TGResultSet;
