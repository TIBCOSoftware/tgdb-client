package com.tibco.tgdb.query;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEntity;

import java.util.Iterator;
import java.util.List;

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

 * File name : TGResultSet.java
 * Created on: 1/22/15
 * Created by: suresh

 * SVN Id: $Id: TGResultSet.java 622 2016-03-19 20:51:12Z ssubrama $
 */

public interface TGResultSet extends Iterator<TGEntity>, AutoCloseable {

    /**
     * Does the Resultset have any Exceptions
     * @return
     */
    boolean hasExceptions();

    /**
     * Get the Exceptions in the ResultSet
     * @return
     */
    List<TGException> getExceptions();

    /**
     * Return nos of entities returned by the query. The result set has a cursor which prefetches "n" rows as
     * per the query constraint. If the nos of entities returned by the query is less than prefetch count, then
     * all are returned.
     * @return
     */
    int count();

    /**
     * Return the first entity in the ResultSet
     * @return
     */
    TGEntity first();

    /**
     * Return the last Entity in the ResultSet
     * @return
     */
    TGEntity last();

    /**
     * Return the prev entity w.r.t to the current cursor position in the ResultSet
     * @return
     */
    TGEntity prev();

    /**
     * Return the next entity w.r.t to the current cursor position in the ResultSet
     * Purely from a completeness point.
     * @return
     */
    TGEntity next();

    /**
     * Get the Current cursor position. A resultset upon creation is set to the position 0.
     * @return
     */
    int getPosition();

    /**
     * Get the entity at the position.
     * @param position
     * @return
     */
    TGEntity getAt(int position);

    void skip(int position);


}
