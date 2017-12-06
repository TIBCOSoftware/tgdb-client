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
 *
 * File name : TGQueryOption.java
 * SVN Id: $Id$
 */


package com.tibco.tgdb.query;

import com.tibco.tgdb.query.impl.QueryOptionImpl;
import com.tibco.tgdb.utils.TGProperties;

/**
 * A Set of QueryOption that allows the user manipulate the results of the query.
 */
public interface TGQueryOption extends TGProperties<String,String> {

    public static TGQueryOption DEFAULT_QUERY_OPTION = new QueryOptionImpl(false);

    /**
     * Create a configurable queryOptions
     * @return
     */
    public static TGQueryOption createQueryOption() {
        return new QueryOptionImpl(true);
    }

    /**
     * Set a limit on the number of entities(nodes and edges) return in a query. Default is 1000
     *
     * @param size number of entities to be returned
     */
    public abstract void setPrefetchSize(int size);

    /**
     * Return the current value of the pre-fetch size
     */
    public abstract int getPrefetchSize();

    /**
     * Set the additional level of traversal from the query result set. Default is 3.
     *
     * @param depth starts with value 1 for one level from the result nodes.
     */
    public abstract void setTraversalDepth(int depth);

    /** 
     * Return the current value of traversal depth
     */
    public abstract int getTraversalDepth();

    /**
     * Set the number of edges per node to be returned in a query.  Default is 0 which means unlimited.
     * 
     * @param limit number of edges
     */
    public abstract void setEdgeLimit(int limit);

    /**
     * Return the current value of edge limit
     */
    public abstract int getEdgeLimit();

    public abstract void setSortAttrName(String name);
    public abstract String getSortAttrName();

    public abstract void setSortOrderDsc(boolean isDsc);
    public abstract boolean isSortOrderDsc();

    public abstract void setSortResultLimit(int limit);
    public abstract int getSortResultLimit();

}
