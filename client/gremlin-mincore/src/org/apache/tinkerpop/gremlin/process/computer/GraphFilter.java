/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.tinkerpop.gremlin.process.computer;

import org.apache.tinkerpop.gremlin.process.traversal.Traversal;
import org.apache.tinkerpop.gremlin.structure.Direction;
import org.apache.tinkerpop.gremlin.structure.Edge;
import org.apache.tinkerpop.gremlin.structure.Vertex;

import java.io.Serializable;
import java.util.Iterator;
import java.util.Set;



public interface GraphFilter extends Cloneable, Serializable {

    /**
     * A enum denoting whether a particular result will be allowed or not.
     * {@link Legal#YES} means that the specified element set will definitely not be removed by {@link GraphFilter}.
     * {@link Legal#MAYBE} means that the element set may or may not be filtered out by the {@link GraphFilter}.
     * {@link Legal#NO} means that the specified element set will definitely be removed by {@link GraphFilter}.
     */
    public enum Legal {
        YES, MAYBE, NO;

        /**
         * The enum is either {@link Legal#YES} or {@link Legal#MAYBE}.
         *
         * @return true if potentially legal.
         */
        public boolean positive() {
            return this != NO;
        }

        /**
         * The enum is {@link Legal#NO}.
         *
         * @return true if definitely not legal.
         */
        public boolean negative() {
            return this == NO;
        }
    }

    public void setVertexFilter(final Traversal<Vertex, Vertex> vertexFilter);

    public void setEdgeFilter(final Traversal<Vertex, Edge> edgeFilter);

    public Traversal<Vertex, Vertex> getVertexFilter();

    public Traversal<Vertex, Edge> getEdgeFilter();

    public boolean legalVertex(final Vertex vertex);

    public Iterator<Edge> legalEdges(final Vertex vertex);

    public boolean hasFilter();

    public boolean hasEdgeFilter();

    public boolean hasVertexFilter();

    public Set<String> getLegallyPositiveEdgeLabels(final Direction direction);

    public Legal checkEdgeLegality(final Direction direction, final String label);

    public Legal checkEdgeLegality(final Direction direction);


}
