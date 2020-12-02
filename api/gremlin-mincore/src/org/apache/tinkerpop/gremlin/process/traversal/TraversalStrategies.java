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
package org.apache.tinkerpop.gremlin.process.traversal;

import org.apache.tinkerpop.gremlin.process.computer.GraphComputer;

import java.io.Serializable;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Optional;

/**
 * A {@link Traversal} maintains a set of {@link TraversalStrategy} instances within a TraversalStrategies object.
 * TraversalStrategies are responsible for compiling a traversal prior to its execution.
 *
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 * @author Matthias Broecheler (me@matthiasb.com)
 */
public interface TraversalStrategies extends Serializable, Cloneable {

    static List<Class<? extends TraversalStrategy>> STRATEGY_CATEGORIES = Collections.unmodifiableList(
            Arrays.asList(
                    TraversalStrategy.DecorationStrategy.class,
                    TraversalStrategy.OptimizationStrategy.class,
                    TraversalStrategy.ProviderOptimizationStrategy.class,
                    TraversalStrategy.FinalizationStrategy.class,
                    TraversalStrategy.VerificationStrategy.class));

    /**
     * Return all the {@link TraversalStrategy} singleton instances associated with this {@link TraversalStrategies}.
     */
    public List<TraversalStrategy<?>> toList();

    /**
     * Return the {@link TraversalStrategy} instance associated with the provided class.
     *
     * @param traversalStrategyClass the class of the strategy to get
     * @param <T>                    the strategy class type
     * @return an optional containing the strategy instance or not
     */
    public default <T extends TraversalStrategy> Optional<T> getStrategy(final Class<T> traversalStrategyClass) {
        return (Optional) toList().stream().filter(s -> traversalStrategyClass.isAssignableFrom(s.getClass())).findAny();
    }

    /**
     * Apply all the {@link TraversalStrategy} optimizers to the {@link Traversal} for the stated {@link TraversalEngine}.
     * This method must ensure that the strategies are sorted prior to application.
     *
     * @param traversal the traversal to apply the strategies to
     */
    public void applyStrategies(final Traversal.Admin<?, ?> traversal);

    /**
     * Add all the provided {@link TraversalStrategy} instances to the current collection.
     * When all the provided strategies have been added, the collection is resorted.
     *
     * @param strategies the traversal strategies to add
     * @return the newly updated/sorted traversal strategies collection
     */
    public TraversalStrategies addStrategies(final TraversalStrategy<?>... strategies);

    /**
     * Remove all the provided {@link TraversalStrategy} classes from the current collection.
     * When all the provided strategies have been removed, the collection is resorted.
     *
     * @param strategyClasses the traversal strategies to remove by their class
     * @return the newly updated/sorted traversal strategies collection
     */
    @SuppressWarnings({"unchecked", "varargs"})
    public TraversalStrategies removeStrategies(final Class<? extends TraversalStrategy>... strategyClasses);

    @SuppressWarnings("CloneDoesntDeclareCloneNotSupportedException")
    public TraversalStrategies clone();

    }
