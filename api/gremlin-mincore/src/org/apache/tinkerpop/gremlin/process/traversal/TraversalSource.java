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

import org.apache.commons.configuration.Configuration;

import org.apache.tinkerpop.gremlin.process.computer.Computer;
import org.apache.tinkerpop.gremlin.process.computer.GraphComputer;
import org.apache.tinkerpop.gremlin.process.remote.RemoteConnection;
import org.apache.tinkerpop.gremlin.structure.Graph;
import java.util.Optional;
import java.util.function.BinaryOperator;
import java.util.function.Supplier;
import java.util.function.UnaryOperator;

/**
 * A {@code TraversalSource} is used to create {@link Traversal} instances.
 * A traversal source can generate any number of {@link Traversal} instances.
 * A traversal source is primarily composed of a {@link Graph} and a {@link TraversalStrategies}.
 * Various {@code withXXX}-based methods are used to configure the traversal strategies (called "configurations").
 * Various other methods (dependent on the traversal source type) will then generate a traversal given the graph and configured strategies (called "spawns").
 * A traversal source is immutable in that fluent chaining of configurations create new traversal sources.
 * This is unlike {@link Traversal} and {@link GraphComputer}, where chained methods configure the same instance.
 * Every traversal source implementation must maintain two constructors to enable proper reflection-based construction.
 * <p/>
 * {@code TraversalSource(Graph)} and {@code TraversalSource(Graph,TraversalStrategies)}
 *
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 */
public interface TraversalSource extends Cloneable, AutoCloseable {

    public static final String GREMLIN_REMOTE = "gremlin.remote.";
    public static final String GREMLIN_REMOTE_CONNECTION_CLASS = GREMLIN_REMOTE + "remoteConnectionClass";

    /**
     * Get the {@link TraversalStrategies} associated with this traversal source.
     *
     * @return the traversal strategies of the traversal source
     */
    public TraversalStrategies getStrategies();

    /**
     * Get the {@link Graph} associated with this traversal source.
     *
     * @return the graph of the traversal source
     */
    public Graph getGraph();

    /**
     * Get the {@link Bytecode} associated with the current state of this traversal source.
     *
     * @return the traversal source byte code
     */
    public Bytecode getBytecode();

    /////////////////////////////

    public static class Symbols {

        private Symbols() {
            // static fields only
        }

        public static final String withSack = "withSack";
        public static final String withStrategies = "withStrategies";
        public static final String withoutStrategies = "withoutStrategies";
        public static final String withComputer = "withComputer";
        public static final String withSideEffect = "withSideEffect";
        public static final String withRemote = "withRemote";
    }

    /////////////////////////////

    /**
     * Add an arbitrary collection of {@link TraversalStrategy} instances to the traversal source.
     *
     * @param traversalStrategies a collection of traversal strategies to add
     * @return a new traversal source with updated strategies
     */
    TraversalSource withStrategies(final TraversalStrategy... traversalStrategies);

    /**
     * Remove an arbitrary collection of {@link TraversalStrategy} classes from the traversal source.
     *
     * @param traversalStrategyClasses a collection of traversal strategy classes to remove
     * @return a new traversal source with updated strategies
     */
    @SuppressWarnings({"unchecked", "varargs"})
    public TraversalSource withoutStrategies(final Class<? extends TraversalStrategy>... traversalStrategyClasses);

    /**
     * Add a {@link Computer} that will generate a {@link GraphComputer} from the {@link Graph} that will be used to execute the traversal.
     * This adds a {@link VertexProgramStrategy} to the strategies.
     *
     * @param computer a builder to generate a graph computer from the graph
     * @return a new traversal source with updated strategies
     */
    TraversalSource withComputer(final Computer computer);

    /**
     * Add a {@link GraphComputer} class used to execute the traversal.
     * This adds a {@link VertexProgramStrategy} to the strategies.
     *
     * @param graphComputerClass the graph computer class
     * @return a new traversal source with updated strategies
     */
    public TraversalSource withComputer(final Class<? extends GraphComputer> graphComputerClass);

    /**
     * Add the standard {@link GraphComputer} of the graph that will be used to execute the traversal.
     * This adds a {@link VertexProgramStrategy} to the strategies.
     *
     * @return a new traversal source with updated strategies
     */
    public TraversalSource withComputer();

    /**
     * Add a sideEffect to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SideEffectStrategy} to the strategies.
     *
     * @param key          the key of the sideEffect
     * @param initialValue a supplier that produces the initial value of the sideEffect
     * @param reducer      a reducer to merge sideEffect mutations into a single result
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSideEffect(final String key, final Supplier<A> initialValue, final BinaryOperator<A> reducer);


    /**
     * Add a sideEffect to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SideEffectStrategy} to the strategies.
     *
     * @param key          the key of the sideEffect
     * @param initialValue the initial value of the sideEffect
     * @param reducer      a reducer to merge sideEffect mutations into a single result
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSideEffect(final String key, final A initialValue, final BinaryOperator<A> reducer);

    /**
     * Add a sideEffect to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SideEffectStrategy} to the strategies.
     *
     * @param key          the key of the sideEffect
     * @param initialValue a supplier that produces the initial value of the sideEffect
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSideEffect(final String key, final Supplier<A> initialValue);

    /**
     * Add a sideEffect to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SideEffectStrategy} to the strategies.
     *
     * @param key          the key of the sideEffect
     * @param initialValue the initial value of the sideEffect
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSideEffect(final String key, final A initialValue);

    /**
     * Add a sack to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SackStrategy} to the strategies.
     *
     * @param initialValue  a supplier that produces the initial value of the sideEffect
     * @param splitOperator the sack split operator
     * @param mergeOperator the sack merge operator
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSack(final Supplier<A> initialValue, final UnaryOperator<A> splitOperator, final BinaryOperator<A> mergeOperator);

    /**
     * Add a sack to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SackStrategy} to the strategies.
     *
     * @param initialValue  the initial value of the sideEffect
     * @param splitOperator the sack split operator
     * @param mergeOperator the sack merge operator
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSack(final A initialValue, final UnaryOperator<A> splitOperator, final BinaryOperator<A> mergeOperator);

    /**
     * Add a sack to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SackStrategy} to the strategies.
     *
     * @param initialValue the initial value of the sideEffect
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSack(final A initialValue);

    /**
     * Add a sack to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SackStrategy} to the strategies.
     *
     * @param initialValue a supplier that produces the initial value of the sideEffect
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSack(final Supplier<A> initialValue);

    /**
     * Add a sack to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SackStrategy} to the strategies.
     *
     * @param initialValue  a supplier that produces the initial value of the sideEffect
     * @param splitOperator the sack split operator
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSack(final Supplier<A> initialValue, final UnaryOperator<A> splitOperator);

    /**
     * Add a sack to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SackStrategy} to the strategies.
     *
     * @param initialValue  the initial value of the sideEffect
     * @param splitOperator the sack split operator
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSack(final A initialValue, final UnaryOperator<A> splitOperator);

    /**
     * Add a sack to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SackStrategy} to the strategies.
     *
     * @param initialValue  a supplier that produces the initial value of the sideEffect
     * @param mergeOperator the sack merge operator
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSack(final Supplier<A> initialValue, final BinaryOperator<A> mergeOperator);

    /**
     * Add a sack to be used throughout the life of a spawned {@link Traversal}.
     * This adds a {@link org.apache.tinkerpop.gremlin.process.traversal.strategy.decoration.SackStrategy} to the strategies.
     *
     * @param initialValue  the initial value of the sideEffect
     * @param mergeOperator the sack merge operator
     * @return a new traversal source with updated strategies
     */
    public <A> TraversalSource withSack(final A initialValue, final BinaryOperator<A> mergeOperator);

    /**
     * Configures the {@code TraversalSource} as a "remote" to issue the {@link Traversal} for execution elsewhere.
     * Expects key for {@link #GREMLIN_REMOTE_CONNECTION_CLASS} as well as any configuration required by
     * the underlying {@link RemoteConnection} which will be instantiated. Note that the {@code Configuration} object
     * is passed down without change to the creation of the {@link RemoteConnection} instance.
     */
    public TraversalSource withRemote(final Configuration conf);

    /**
     * Configures the {@code TraversalSource} as a "remote" to issue the {@link Traversal} for execution elsewhere.
     * Calls {@link #withRemote(Configuration)} after reading the properties file specified.
     */
    public TraversalSource withRemote(final String configFile) throws Exception;

    /**
     * Configures the {@code TraversalSource} as a "remote" to issue the {@link Traversal} for execution elsewhere.
     * Implementations should track {@link RemoteConnection} instances that are created and call
     * {@link RemoteConnection#close()} on them when the {@code TraversalSource} itself is closed.
     *
     * @param connection the {@link RemoteConnection} instance to use to submit the {@link Traversal}.
     */
    public TraversalSource withRemote(final RemoteConnection connection);

    public Optional<Class> getAnonymousTraversalClass();

    /**
     * The clone-method should be used to create immutable traversal sources with each call to a configuration "withXXX"-method.
     * The clone-method should clone the {@link Bytecode}, {@link TraversalStrategies}, mutate the cloned strategies accordingly,
     * and then return the cloned traversal source leaving the original unaltered.
     *
     * @return the cloned traversal source
     */
    @SuppressWarnings("CloneDoesntDeclareCloneNotSupportedException")
    public TraversalSource clone();

    @Override
    public void close() throws Exception;

}
