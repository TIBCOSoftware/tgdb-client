
/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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

 * File name : TGResultSet.java
 * Created on: 1/22/15
 * Created by: suresh

 * SVN Id: $Id: TGResultSet.java 3952 2020-05-04 20:18:13Z vchung $
 */

package com.tibco.tgdb.query;

import com.tibco.tgdb.exception.TGException;

import java.util.Collection;
import java.util.Iterator;
import java.util.List;


public interface TGResultSet<T> extends Iterator<T>, AutoCloseable {

    /**
     * Does the Resultset have any Exceptions
     * @return true has exceptions and false has no exception
     */
    boolean hasExceptions();

    /**
     * Get the Exceptions in the ResultSet
     * @return list of exceptions
     */
    List<TGException> getExceptions();

    /**
     * Return nos of entities returned by the query. The result set has a cursor which prefetches "n" rows as
     * per the query constraint. If the nos of entities returned by the query is less than prefetch count, then
     * all are returned.
     * @return number of entities in the result set
     */
    int count();

    /**
     * Return the first entity in the ResultSet
     * @return first entity in the result set
     */
    T first();

    /**
     * Return the last Entity in the ResultSet
     * @return last entity in the result set
     */
    T last();

    /**
     * Return the prev entity w.r.t to the current cursor position in the ResultSet
     * @return previous entity from the current result set position
     */
    T prev();

    /**
     * Get the Current cursor position. A resultset upon creation is set to the position 0.
     * @return current result set position
     */
    int getPosition();

    /**
     * Get the entity at the position.
     * @param position get entity at this position
     * @return entity at specific position
     */
    T getAt(int position);

    /**
     * Skip a number of position
     * @param position skip the number of position 
     * 
     */
    void skip(int position);
    
    /**
     * Convert the result set into a collection
     * @return collection of the result set data
     */
    Collection<T> toCollection();

    /**
     * Get result set meta data for elements in the result set
     * @return ResultSetMetaData
     */
    TGResultSetMetaData getMetaData();
}
