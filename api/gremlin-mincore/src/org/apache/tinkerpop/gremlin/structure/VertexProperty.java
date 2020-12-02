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
package org.apache.tinkerpop.gremlin.structure;



import org.apache.tinkerpop.gremlin.structure.util.empty.EmptyVertexProperty;

import java.util.Iterator;

/**
 * A {@code VertexProperty} is similar to a {@link Property} in that it denotes a key/value pair associated with an
 * {@link Vertex}, however it is different in the sense that it also represents an entity that it is an {@link Element}
 * that can have properties of its own.
 * <p/>
 * A property is much like a Java8 {@link java.util.Optional} in that a property can be not present (i.e. empty).
 * The key of a property is always a String and the value of a property is an arbitrary Java object.
 * Each underlying graph engine will typically have constraints on what Java objects are allowed to be used as values.
 *
 * @author Matthias Broecheler (me@matthiasb.com)
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 * @author Stephen Mallette (http://stephen.genoprime.com)
 */
public interface VertexProperty<V> extends Property<V>, Element {

     static final String _LABEL = "vertexProperty";

     enum Cardinality {
        single, list, set
    }

    /**
     * Constructs an empty {@code VertexProperty}.
     */
     static <V> VertexProperty<V> empty() {
        return EmptyVertexProperty.instance();
    }

    /**
     * {@inheritDoc}
     */
    @Override
     <U> Iterator<Property<U>> properties(final String... propertyKeys);


    /**
     * Common exceptions to use with a property.
     */
    public static class Exceptions {

        private Exceptions() {
        }

        public static UnsupportedOperationException userSuppliedIdsNotSupported() {
            return new UnsupportedOperationException("VertexProperty does not support user supplied identifiers");
        }

        public static UnsupportedOperationException userSuppliedIdsOfThisTypeNotSupported() {
            return new UnsupportedOperationException("VertexProperty does not support user supplied identifiers of this type");
        }

        public static UnsupportedOperationException multiPropertiesNotSupported() {
            return new UnsupportedOperationException("Multiple properties on a vertex is not supported");
        }

        public static UnsupportedOperationException identicalMultiPropertiesNotSupported() {
            return new UnsupportedOperationException("Multiple properties on a vertex is supported, but a single key may not hold the same value more than once");
        }

        public static UnsupportedOperationException metaPropertiesNotSupported() {
            return new UnsupportedOperationException("Properties on a vertex property is not supported");
        }
    }
}
