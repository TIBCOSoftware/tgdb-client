/*
 * Copyright 2018 TIBCO Software Inc. All rights reserved.
 *
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

import org.apache.commons.configuration.Configuration;
import org.apache.tinkerpop.gremlin.process.computer.GraphComputer;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal;
import org.apache.tinkerpop.gremlin.process.traversal.step.TraversalParent;
import org.apache.tinkerpop.gremlin.process.traversal.traverser.TraverserRequirement;
import org.apache.tinkerpop.gremlin.process.traversal.util.TraversalExplanation;
import org.apache.tinkerpop.gremlin.process.traversal.util.TraversalMetrics;
import org.apache.tinkerpop.gremlin.structure.Graph;

import java.io.Serializable;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.stream.Stream;

/**
 * A {@link Traversal} represents a directed walk over a {@link Graph}.
 * This is the base interface for all traversal's, where each extending interface is seen as a domain specific language.
 * For example, {@link GraphTraversal} is a domain specific language for traversing a graph using "graph concepts" (e.g. vertices, edges).
 * Another example may represent the graph using "social concepts" (e.g. people, cities, artifacts).
 * A {@link Traversal} is evaluated in one of two ways: iterator-based OLTP or {@link GraphComputer}-based OLAP.
 * OLTP traversals leverage an iterator and are executed within a single JVM (with data access allowed to be remote).
 * OLAP traversals leverage {@link GraphComputer} and are executed between multiple JVMs (and/or cores).
 *
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 */
public interface Traversal<S, E> extends Iterator<E>, Serializable, Cloneable, AutoCloseable {

    public static class Symbols {
        private Symbols() {
            // static fields only
        }

        public static final String profile = "profile";
        public static final String none = "none";
    }

    /**
     * Get access to administrative methods of the traversal via its accompanying {@link Traversal.Admin}.
     *
     * @return the admin of this traversal
     */
    public Traversal.Admin<S, E> asAdmin();

    /**
     * Return an {@link Optional} of the next E object in the traversal.
     * If the traversal is empty, then an {@link Optional#empty()} is returned.
     *
     * @return an optional of the next object in the traversal
     */
    public Optional<E> tryNext();

    /**
     * Get the next n-number of results from the traversal.
     * If the traversal has less than n-results, then only that number of results are returned.
     *
     * @param amount the number of results to get
     * @return the n-results in a {@link List}
     */
    public List<E> next(final int amount);

    /**
     * Put all the results into an {@link ArrayList}.
     *
     * @return the results in a list
     */
    public List<E> toList();

    /**
     * Put all the results into a {@link HashSet}.
     *
     * @return the results in a set
     */
    public Set<E> toSet();

    /**
     * Put all the results into a {@link BulkSet}.
     * This can reduce both time and space when aggregating results by ensuring a weighted set.
     *
     * @return the results in a bulk set
     */
    public Set<E> toBulkSet();
    /**
     * Return the traversal as a {@link Stream}.
     *
     * @return the traversal as a stream.
     */
    public Stream<E> toStream();

    /**
     * Starts a promise to execute a function on the current {@code Traversal} that will be completed in the future.
     * Note that this method can only be used if the {@code Traversal} is constructed using
     * {@link TraversalSource#withRemote(Configuration)}. Calling this method otherwise will yield an
     * {@code IllegalStateException}.
     */
    public <T> CompletableFuture<T> promise(final Function<Traversal<S, E>, T> traversalFunction);

    /**
     * Add all the results of the traversal to the provided collection.
     *
     * @param collection the collection to fill
     * @return the collection now filled
     */
    public <C extends Collection<E>> C fill(final C collection);

    /**
     * Iterate all the {@link Traverser} instances in the traversal.
     * What is returned is the empty traversal.
     * It is assumed that what is desired from the computation is are the sideEffects yielded by the traversal.
     *
     * @return the fully drained traversal
     */
    public <A, B> Traversal<A, B> iterate();

    /**
     * Filter all traversers in the traversal.
     * @return the updated traversal with respective {@link NoneStep}.
     */
    public Traversal<S, E> none();

    /**
     * Profile the traversal.
     * @return the updated traversal with respective {@link ProfileSideEffectStep}.
     */
    public Traversal<S, TraversalMetrics> profile();

    /**
     * Return a {@link TraversalExplanation} that shows how this traversal will mutate with each applied {@link TraversalStrategy}.
     * @return a traversal explanation
     */
    public TraversalExplanation explain();
    /**
     * A traversal can be rewritten such that its defined end type E may yield objects of a different type.
     * This helper method allows for the casting of the output to the known the type.
     * @param endType  the true output type of the traversal
     * @param consumer a {@link Consumer} to process each output
     * @param <E2>     the known output type of the traversal
     */
    public <E2> void forEachRemaining(final Class<E2> endType, final Consumer<E2> consumer);


    public void forEachRemaining(final Consumer<? super E> action);

    /**
     * Releases resources opened in any steps that implement {@link AutoCloseable}.
     */
    @Override
    public void close() throws Exception;

    /**
     * A collection of {@link Exception} types associated with Traversal execution.
     */
    public static class Exceptions {

        public static IllegalStateException traversalIsLocked() {
            return new IllegalStateException("The traversal strategies are complete and the traversal can no longer be modulated");
        }

        public static IllegalStateException traversalIsNotReversible() {
            return new IllegalStateException("The traversal is not reversible as it contains steps that are not reversible");
        }
    }

    public interface Admin<S, E> extends Traversal<S, E> {

        /**
         * Get the {@link Bytecode} associated with the construction of this traversal.
         * @return the byte code representation of the traversal
         */
        public Bytecode getBytecode();

        /**
         * Add an iterator of {@link Traverser.Admin} objects to the head/start of the traversal.
         * Users should typically not need to call this method. For dynamic inject of data, they should use {@link org.apache.tinkerpop.gremlin.process.traversal.step.sideEffect.InjectStep}.
         * @param starts an iterators of traversers
         */
        public void addStarts(final Iterator<Traverser.Admin<S>> starts);

        /**
         * Add a single {@link Traverser.Admin} object to the head of the traversal.
         * Users should typically not need to call this method. For dynamic inject of data, they should use {@link org.apache.tinkerpop.gremlin.process.traversal.step.sideEffect.InjectStep}.
         * @param start a traverser to add to the traversal
         */
        public void addStart(final Traverser.Admin<S> start);

        /**
         * Get the {@link Step} instances associated with this traversal.
         * The steps are ordered according to their linked list structure as defined by {@link Step#getPreviousStep()} and {@link Step#getNextStep()}.
         *
         * @return the ordered steps of the traversal
         */
        public List<Step> getSteps();

        /**
         * Add a {@link Step} to the end of the traversal. This method should link the step to its next and previous step accordingly.
         *
         * @param step the step to add
         * @param <E2> the output of the step
         * @return the updated traversal
         * @throws IllegalStateException if the {@link TraversalStrategies} have already been applied
         */
        public <E2> Traversal.Admin<S, E2> addStep(final Step<?, E2> step) throws IllegalStateException ;

        /**
         * Add a {@link Step} to an arbitrary point in the traversal.
         *
         * @param index the location in the traversal to insert the step
         * @param step  the step to add
         * @param <S2>  the new start type of the traversal (if the added step was a start step)
         * @param <E2>  the new end type of the traversal (if the added step was an end step)
         * @return the newly modulated traversal
         * @throws IllegalStateException if the {@link TraversalStrategies} have already been applied
         */
        public <S2, E2> Traversal.Admin<S2, E2> addStep(final int index, final Step<?, ?> step) throws IllegalStateException;

        /**
         * Remove a {@link Step} from the traversal.
         *
         * @param step the step to remove
         * @param <S2> the new start type of the traversal (if the removed step was a start step)
         * @param <E2> the new end type of the traversal (if the removed step was an end step)
         * @return the newly modulated traversal
         * @throws IllegalStateException if the {@link TraversalStrategies} have already been applied
         */
        public <S2, E2> Traversal.Admin<S2, E2> removeStep(final Step<?, ?> step) throws IllegalStateException;

        /**
         * Remove a {@link Step} from the traversal.
         *
         * @param index the location in the traversal of the step to be evicted
         * @param <S2>  the new start type of the traversal (if the removed step was a start step)
         * @param <E2>  the new end type of the traversal (if the removed step was an end step)
         * @return the newly modulated traversal
         * @throws IllegalStateException if the {@link TraversalStrategies} have already been applied
         */
        public <S2, E2> Traversal.Admin<S2, E2> removeStep(final int index) throws IllegalStateException;

        /**
         * Get the start/head of the traversal. If the traversal is empty, then an {@link EmptyStep} instance is returned.
         *
         * @return the start step of the traversal
         */
        public Step<S, ?> getStartStep();

        /**
         * Get the end/tail of the traversal. If the traversal is empty, then an {@link EmptyStep} instance is returned.
         *
         * @return the end step of the traversal
         */
        public Step<?, E> getEndStep();

        /**
         * Apply the registered {@link TraversalStrategies} to the traversal.
         * Once the strategies are applied, the traversal is "locked" and can no longer have steps added to it.
         * The order of operations for strategy applications should be: globally id steps, apply strategies to root traversal, then to nested traversals.
         *
         * @throws IllegalStateException if the {@link TraversalStrategies} have already been applied
         */
        public void applyStrategies() throws IllegalStateException;

        /**
         * Get the {@link TraverserGenerator} associated with this traversal.
         * The traversal generator creates {@link Traverser} instances that are respective of the traversal's {@link org.apache.tinkerpop.gremlin.process.traversal.traverser.TraverserRequirement}.
         *
         * @return the generator of traversers
         */
        public TraverserGenerator getTraverserGenerator();

        /**
         * Get the set of all {@link TraverserRequirement}s for this traversal.
         *
         * @return the features of a traverser that are required to execute properly in this traversal
         */
        public Set<TraverserRequirement> getTraverserRequirements();

        /**
         * Call the {@link Step#reset} method on every step in the traversal.
         */
        public default void reset() {
            this.getSteps().forEach(Step::reset);
        }

        /**
         * Set the {@link TraversalSideEffects} of this traversal.
         *
         * @param sideEffects the sideEffects to set for this traversal.
         */
        public void setSideEffects(final TraversalSideEffects sideEffects);

        /**
         * Get the {@link TraversalSideEffects} associated with the traversal.
         *
         * @return The traversal sideEffects
         */
        public TraversalSideEffects getSideEffects();

        /**
         * Set the {@link TraversalStrategies} to be used by this traversal at evaluation time.
         *
         * @param strategies the strategies to use on this traversal
         */
        public void setStrategies(final TraversalStrategies strategies);

        /**
         * Get the {@link TraversalStrategies} associated with this traversal.
         *
         * @return the strategies associated with this traversal
         */
        public TraversalStrategies getStrategies();

        /**
         * Set the {@link org.apache.tinkerpop.gremlin.process.traversal.step.TraversalParent} {@link Step} that is the parent of this traversal.
         * Traversals can be nested and this is the means by which the traversal tree is connected.
         *
         * @param step the traversal holder parent step
         */
        public void setParent(final TraversalParent step);

        /**
         * Get the {@link org.apache.tinkerpop.gremlin.process.traversal.step.TraversalParent} {@link Step} that is the parent of this traversal.
         * Traversals can be nested and this is the means by which the traversal tree is walked.
         *
         * @return the traversal holder parent step
         */
        public TraversalParent getParent();

        /**
         * Cloning is used to duplicate the traversal typically in OLAP environments.
         *
         * @return The cloned traversal
         */
        @SuppressWarnings("CloneDoesntDeclareCloneNotSupportedException")
        public Traversal.Admin<S, E> clone();

        /**
         * When the traversal has had its {@link TraversalStrategies} applied to it, it is locked.
         *
         * @return whether the traversal is locked
         */
        public boolean isLocked();

        public Optional<Graph> getGraph();

        public void setGraph(final Graph graph);

        public boolean equals(final Traversal.Admin<S, E> other);

        public Traverser.Admin<E> nextTraverser();

    }

}
