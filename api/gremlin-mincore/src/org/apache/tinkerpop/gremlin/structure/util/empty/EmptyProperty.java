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

import org.apache.tinkerpop.gremlin.structure.Element;
import org.apache.tinkerpop.gremlin.structure.Property;
import org.apache.tinkerpop.gremlin.structure.util.StringFactory;

import java.util.NoSuchElementException;
import java.util.function.Consumer;
import java.util.function.Supplier;

/**
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 */
public final class EmptyProperty<V> implements Property<V> {

    private static final EmptyProperty INSTANCE = new EmptyProperty();

    private EmptyProperty() {

    }

    @Override
    public String key() {
        throw Exceptions.propertyDoesNotExist();
    }

    @Override
    public V value() throws NoSuchElementException {
        throw Exceptions.propertyDoesNotExist();
    }

    @Override
    public boolean isPresent() {
        return false;
    }

    @Override
    public Element element() {
        throw Exceptions.propertyDoesNotExist();
    }

    @Override
    public void remove() { throw Exceptions.propertyDoesNotExist(); }

    @Override
    public void ifPresent(Consumer<? super V> consumer) {   }

    @Override
    public V orElse(V otherValue) { return null;  }

    @Override
    public V orElseGet(Supplier<? extends V> valueSupplier) { return null; }

    @Override
    public <E extends Throwable> V orElseThrow(Supplier<? extends E> exceptionSupplier) throws E { return null; }

    @Override
    public String toString() {
        return StringFactory.propertyString(this);
    }

    @SuppressWarnings("EqualsWhichDoesntCheckParameterClass")
    @Override
    public boolean equals(final Object object) {
        return object instanceof EmptyProperty;
    }

    public int hashCode() {
        return 1281483122;
    }

    public static <V> Property<V> instance() {
        return INSTANCE;
    }
}
