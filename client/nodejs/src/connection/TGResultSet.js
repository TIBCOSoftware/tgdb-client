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

function TGResultSet () {
}

/**
 * Does the Resultset have any Exceptions
 * @return
 */
TGResultSet.prototype.hasExceptions = function(){};

    /**
     * Get the Exceptions in the ResultSet
     * @return
     */
TGResultSet.prototype.getExceptions = function(){};

    /**
     * Return nos of entities returned by the query. The result set has a cursor which prefetches "n" rows as
     * per the query constraint. If the nos of entities returned by the query is less than prefetch count, then
     * all are returned.
     * @return
     */
TGResultSet.prototype.count = function(){};

    /**
     * Return the first entity in the ResultSet
     * @return
     */
TGResultSet.prototype.first = function(){};

    /**
     * Return the last Entity in the ResultSet
     * @return
     */
TGResultSet.prototype.last = function(){};

    /**
     * Return the prev entity w.r.t to the current cursor position in the ResultSet
     * @return
     */
TGResultSet.prototype.prev = function(){};

    /**
     * Return the next entity w.r.t to the current cursor position in the ResultSet
     * Purely from a completeness point.
     * @return
     */
TGResultSet.prototype.next = function(){};

    /**
     * Get the Current cursor position. A resultset upon creation is set to the position 0.
     * @return
     */
TGResultSet.prototype.getPosition = function(){};

    /**
     * Get the entity at the position.
     * @param position
     * @return
     */
TGResultSet.prototype.getAt = function(position){};

TGResultSet.prototype.skip = function(position){};

exports.TGResultSet = TGResultSet;
