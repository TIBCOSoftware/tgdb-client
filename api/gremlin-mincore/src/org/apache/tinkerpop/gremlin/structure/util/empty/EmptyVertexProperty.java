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
package org.apache.tinkerpop.gremlin.structure.util.empty;

import org.apache.tinkerpop.gremlin.structure.Graph;
import org.apache.tinkerpop.gremlin.structure.Property;
import org.apache.tinkerpop.gremlin.structure.Vertex;
import org.apache.tinkerpop.gremlin.structure.VertexProperty;
import org.apache.tinkerpop.gremlin.structure.util.StringFactory;

import java.util.Collections;
import java.util.Iterator;
import java.util.NoSuchElementException;
import java.util.Set;
import java.util.function.Consumer;
import java.util.function.Supplier;

/**
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 */
public final class EmptyVertexProperty<V> implements VertexProperty<V> {

    private static final EmptyVertexProperty INSTANCE = new EmptyVertexProperty();

    public static <U> VertexProperty<U> instance() {
        return INSTANCE;
    }

    @Override
    public Vertex element() {
        throw Property.Exceptions.propertyDoesNotExist();
    }

    @Override
    public Object id() {
        throw Property.Exceptions.propertyDoesNotExist();
    }

    @Override
    public Graph graph() {
        throw Property.Exceptions.propertyDoesNotExist();
    }

    @Override
    public <U> Property<U> property(String key) {
        return Property.<U>empty();
    }

    @Override
    public <U> Property<U> property(String key, U value) {
        return Property.<U>empty();
    }

    @Override
    public String key() {
        throw Property.Exceptions.propertyDoesNotExist();
    }

    @Override
    public V value() throws NoSuchElementException {
        throw Property.Exceptions.propertyDoesNotExist();
    }

    @Override
    public boolean isPresent() {
        return false;
    }

    @Override
    public void remove() {

    }

    @Override
    public String label() {
        return null;
    }

    @Override
    public Set<String> keys() {
        return null;
    }

    @Override
    public <V> V value(String key) throws NoSuchElementException {
        return null;
    }

    @Override
    public <V> Iterator<V> values(String... propertyKeys) {
        return null;
    }

    @Override
    public void ifPresent(Consumer<? super V> consumer) {

    }

    @Override
    public V orElse(V otherValue) {
        return null;
    }

    @Override
    public V orElseGet(Supplier<? extends V> valueSupplier) {
        return null;
    }

    @Override
    public <E extends Throwable> V orElseThrow(Supplier<? extends E> exceptionSupplier) throws E {
        return null;
    }

    @Override
    public String toString() {
        return StringFactory.propertyString(this);
    }

    @Override
    public <U> Iterator<Property<U>> properties(String... propertyKeys) {
        return Collections.emptyIterator();
    }
}
