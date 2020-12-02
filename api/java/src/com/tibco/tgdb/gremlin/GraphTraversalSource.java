/**
 * Copyright (c) 2019 TIBCO Software Inc.
 * All rights reserved.
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
 * <p/>
 * File name: GraphTraversalSource.java
 * Created on: 2018-10-24
 * Created by: vincent
 * <p/>
 * SVN Id: $Id: GraphTraversalSource.java 3158 2019-04-26 20:49:24Z kattaylo $
 */

package com.tibco.tgdb.gremlin;

import java.util.Optional;
import java.util.function.BinaryOperator;
import java.util.function.Supplier;
import java.util.function.UnaryOperator;

import org.apache.commons.configuration.Configuration;
import org.apache.tinkerpop.gremlin.process.computer.Computer;
import org.apache.tinkerpop.gremlin.process.computer.GraphComputer;
import org.apache.tinkerpop.gremlin.process.remote.RemoteConnection;
import org.apache.tinkerpop.gremlin.process.traversal.Bytecode;
import org.apache.tinkerpop.gremlin.process.traversal.TraversalSource;
import org.apache.tinkerpop.gremlin.process.traversal.TraversalStrategies;
import org.apache.tinkerpop.gremlin.process.traversal.TraversalStrategy;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal;
import org.apache.tinkerpop.gremlin.structure.Edge;
import org.apache.tinkerpop.gremlin.structure.Graph;
import org.apache.tinkerpop.gremlin.structure.Vertex;
import org.apache.tinkerpop.gremlin.structure.util.StringFactory;

import com.tibco.tgdb.connection.TGConnection;

public class GraphTraversalSource implements TraversalSource {
	//protected transient RemoteConnection connection;
	//Allow sharing the connection with DefaultTraversal
	transient TGConnection connection;
    protected final Graph graph;
    protected TraversalStrategies strategies;
    protected Bytecode bytecode = new Bytecode();

	public GraphTraversalSource(final Graph graph) {
		this.graph = graph;
	}
		   
	@Override
	public TraversalStrategies getStrategies() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Graph getGraph() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Bytecode getBytecode() {
		// TODO Auto-generated method stub
		return this.bytecode;
	}

	@Override
	public TraversalSource withStrategies(TraversalStrategy... traversalStrategies) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public TraversalSource withoutStrategies(Class<? extends TraversalStrategy>... traversalStrategyClasses) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public TraversalSource withComputer(Computer computer) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public TraversalSource withComputer(Class<? extends GraphComputer> graphComputerClass) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public TraversalSource withComputer() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSideEffect(String key, Supplier<A> initialValue, BinaryOperator<A> reducer) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSideEffect(String key, A initialValue, BinaryOperator<A> reducer) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSideEffect(String key, Supplier<A> initialValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSideEffect(String key, A initialValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSack(Supplier<A> initialValue, UnaryOperator<A> splitOperator,
			BinaryOperator<A> mergeOperator) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSack(A initialValue, UnaryOperator<A> splitOperator,
			BinaryOperator<A> mergeOperator) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSack(A initialValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSack(Supplier<A> initialValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSack(Supplier<A> initialValue, UnaryOperator<A> splitOperator) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSack(A initialValue, UnaryOperator<A> splitOperator) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSack(Supplier<A> initialValue, BinaryOperator<A> mergeOperator) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <A> TraversalSource withSack(A initialValue, BinaryOperator<A> mergeOperator) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public TraversalSource withRemote(Configuration conf) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public TraversalSource withRemote(String configFile) throws Exception {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public GraphTraversalSource withRemote(RemoteConnection connection) {
		// TODO Auto-generated method stub
		return null;
	}

	//This is created specific for handling TGConnection
	public GraphTraversalSource withRemote(TGConnection connection) {
		// TODO Auto-generated method stub
        this.connection = connection;
        final TraversalSource clone = this.clone();
        return (GraphTraversalSource) clone;
	}

	@Override
	public Optional<Class> getAnonymousTraversalClass() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public GraphTraversalSource clone() {
		// TODO Auto-generated method stub
        try {
            final GraphTraversalSource clone = (GraphTraversalSource) super.clone();
            //clone.strategies = this.strategies.clone();
            clone.bytecode = this.bytecode.clone();
            clone.connection = this.connection;
            return clone;
        } catch (final CloneNotSupportedException e) {
            throw new IllegalStateException(e.getMessage(), e);
        }
	}

	@Override
	public void close() throws Exception {
		// TODO Auto-generated method stub

	}

	public GraphTraversalSource withBulk(final boolean useBulk) {
		//Copy from gremlin
		return null;
	}
	
	public GraphTraversalSource withPath() {
		return null;
	}
	
	public <S> GraphTraversal<S, S> inject(S... starts) {
        final GraphTraversalSource clone = this.clone();
        clone.bytecode.addStep(GraphTraversal.Symbols.inject, starts);
        final GraphTraversal.Admin<S, S> traversal = new DefaultGraphTraversal<>(clone);
        return traversal;
	}
	
	public GraphTraversal<Vertex, Vertex> V(final Object... vertexIds) {
		final GraphTraversalSource clone = this.clone();
        clone.bytecode.addStep(GraphTraversal.Symbols.V, vertexIds);
        final GraphTraversal.Admin<Vertex, Vertex> traversal = new DefaultGraphTraversal<>(clone);
//        return traversal.addStep(new GraphStep<>(traversal, Vertex.class, true, vertexIds));
        return traversal;
	}
	
	public GraphTraversal<Edge, Edge> E(final Object... edgesIds) {
        GraphTraversalSource clone = this.clone();
        clone.bytecode.addStep(GraphTraversal.Symbols.E, edgesIds);
        final GraphTraversal.Admin<Edge, Edge> traversal = new DefaultGraphTraversal<>(clone);
		return traversal;
	}
	
	public String toString() {
        return StringFactory.traversalSourceString(this);
    }
	
	public TGConnection getConnection() {
		return this.connection;
	}
}
