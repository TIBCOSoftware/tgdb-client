/*
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
package org.apache.tinkerpop.gremlin.structure.util;

import org.apache.tinkerpop.gremlin.structure.Direction;
import org.apache.tinkerpop.gremlin.structure.Edge;
import org.apache.tinkerpop.gremlin.structure.Element;
import org.apache.tinkerpop.gremlin.structure.Graph;
import org.apache.tinkerpop.gremlin.structure.Property;
import org.apache.tinkerpop.gremlin.structure.T;
import org.apache.tinkerpop.gremlin.structure.Vertex;
import org.apache.tinkerpop.gremlin.structure.VertexProperty;

import java.util.Iterator;
import java.util.Optional;
import java.util.function.Function;

/**
 * An interface that provides methods for detached properties and elements to be re-attached to the {@link Graph}.
 * There are two general ways in which they can be attached: {@link Method#get} or {@link Method#create}.
 * A {@link Method#get} will find the property/element at the host location and return it.
 * A {@link Method#create} will create the property/element at the host location and return it.
 *
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 * @author Stephen Mallette (http://stephen.genoprime.com)
 */
public interface Attachable<V> {

    /**
     * Get the raw object trying to be attached.
     *
     * @return the raw object to attach
     */
    public V get();

    /**
     * Provide a way to attach an {@link Attachable} implementation to a host.  Note that the context of the host
     * is not defined by way of the attachment method itself that is supplied as an argument.  It is up to the
     * implementer to supply that context.
     *
     * @param method a {@link Function} that takes an {@link Attachable} and returns the "re-attached" object
     * @return the return value of the {@code method}
     * @throws IllegalStateException if the {@link Attachable} is not a "graph" object (i.e. host or
     *                               attachable don't work together)
     */
    public default V attach(final Function<Attachable<V>, V> method) throws IllegalStateException {
        return method.apply(this);
    }

    public static class Exceptions {

        private Exceptions() {
        }

        public static IllegalStateException canNotGetAttachableFromHostVertex(final Attachable<?> attachable, final Vertex hostVertex) {
            return new IllegalStateException("Can not get the attachable from the host vertex: " + attachable + "-/->" + hostVertex);
        }

        public static IllegalStateException canNotGetAttachableFromHostGraph(final Attachable<?> attachable, final Graph hostGraph) {
            return new IllegalStateException("Can not get the attachable from the host vertex: " + attachable + "-/->" + hostGraph);
        }

        public static IllegalArgumentException providedAttachableMustContainAGraphObject(final Attachable<?> attachable) {
            return new IllegalArgumentException("The provided attachable must contain a graph object: " + attachable);
        }
    }

}