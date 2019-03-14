package com.tibco.tgdb.model;

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
 *  File name :TGEdge.java
 *  Created by: suresh
 *
 *		SVN Id: $Id: TGEdge.java 2344 2018-06-11 23:21:45Z ssubrama $
 *
 */

/**
 * An Egde in a graph that connects 2 vertices.
 */
public interface TGEdge extends TGEntity {

    enum DirectionType {
        UnDirected,
        Directed,
        BiDirectional
    }

    enum Direction {
    	Inbound,
    	Outbound,
    	Any
    }

    /**
     * Get the pair of vertices for this Edge
     * @return An array of Nodes (typically 2) for this edge
     */
    TGNode[] getVertices();

    /**
     * Get the direction desc of this edge.
     * @return the DirectionTyoe of this edge
     */
    DirectionType getDirectionType();

}
