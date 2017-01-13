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

    public abstract void setPrefetchSize(int size);
    public abstract int getPrefetchSize();

    public abstract void setTraversalDepth(int depth);
    public abstract int getTraversalDepth();

    public abstract void setEdgeLimit(int depth);
    public abstract int getEdgeLimit();




}
