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

import java.util.HashSet;
import java.util.Iterator;
import java.util.NoSuchElementException;
import java.util.Set;

/**
 * An {@link Element} is the base class for both {@link Vertex} and {@link Edge}. An {@link Element} has an identifier
 * that must be unique to its inheriting classes ({@link Vertex} or {@link Edge}). An {@link Element} can maintain a
 * collection of {@link Property} objects.  Typically, objects are Java primitives (e.g. String, long, int, boolean,
 * etc.)
 *
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 * @author Stephen Mallette (http://stephen.genoprime.com)
 * @author suresh (suresh.subramani@tibco.com) 
 *
 */
public interface Element {

    /**
     * Gets the unique identifier for the graph {@code Element}.
     *
     * @return The id of the element
     */
    Object id();

    /**
     * Gets the label for the graph {@code Element} which helps categorize it.
     *
     * @return The label of the element
     */
    String label();

    /**
     * Get the graph that this element is within.
     *
     * @return the graph of this element
     */
    Graph graph();

    /**
     * Get the keys of the properties associated with this element.
     * The  implementation iterators the properties and stores the keys into a {@link HashSet}.
     *
     * @return The property key set
     */
    Set<String> keys();
    /**
     * Get a {@link Property} for the {@code Element} given its key.
     * The  implementation calls the raw {@link Element#properties}.
     */
    <V> Property<V> property(final String key);
    /**
     * Add or set a property value for the {@code Element} given its key.
     */
    <V> Property<V> property(final String key, final V value);

    /**
     * Get the value of a {@link Property} given it's key.
     * The  implementation calls {@link Element#property} and then returns the associated value.
     *
     * @throws NoSuchElementException if the property does not exist on the {@code Element}.
     */
    <V> V value(final String key) throws NoSuchElementException;

    /**
     * Removes the {@code Element} from the graph.
     */
    void remove();


    /**
     * Get the values of properties as an {@link Iterator}.
     */
    <V> Iterator<V> values(final String... propertyKeys);

    /**
     * Get an {@link Iterator} of properties.
     */
    <V> Iterator<? extends Property<V>> properties(final String... propertyKeys);


    /**
     * Common exceptions to use with an element.
     */
    public static class Exceptions {

        protected Exceptions() {
        }

        public static IllegalArgumentException providedKeyValuesMustBeAMultipleOfTwo() {
            return new IllegalArgumentException("The provided key/value array length must be a multiple of two");
        }

        public static IllegalArgumentException providedKeyValuesMustHaveALegalKeyOnEvenIndices() {
            return new IllegalArgumentException("The provided key/value array must have a String or T on even array indices");
        }

        public static IllegalStateException propertyAdditionNotSupported() {
            return new IllegalStateException("Property addition is not supported");
        }

        public static IllegalArgumentException labelCanNotBeNull() {
            return new IllegalArgumentException("Label can not be null");
        }

        public static IllegalArgumentException labelCanNotBeEmpty() {
            return new IllegalArgumentException("Label can not be empty");
        }

        public static IllegalArgumentException labelCanNotBeAHiddenKey(final String label) {
            return new IllegalArgumentException("Label can not be a hidden key: " + label);
        }

        /**
         * @deprecated As of release 3.1.0, not replaced - this exception is no longer enforced by the test suite.
         * @see <a href="https://issues.apache.org/jira/browse/TINKERPOP-297">TINKERPOP-297</a>
         */
        @Deprecated
        public static IllegalStateException elementAlreadyRemoved(final Class<? extends Element> clazz, final Object id) {
            return new IllegalStateException(String.format("%s with id %s was removed.", clazz.getSimpleName(), id));
        }
    }
}
