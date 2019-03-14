package com.tibco.tgdb.gremlin;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.Comparator;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.function.BiFunction;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

import org.apache.tinkerpop.gremlin.process.computer.VertexProgram;
import org.apache.tinkerpop.gremlin.process.traversal.Order;
import org.apache.tinkerpop.gremlin.process.traversal.P;
import org.apache.tinkerpop.gremlin.process.traversal.Path;
import org.apache.tinkerpop.gremlin.process.traversal.Pop;
import org.apache.tinkerpop.gremlin.process.traversal.Scope;
import org.apache.tinkerpop.gremlin.process.traversal.Step;
import org.apache.tinkerpop.gremlin.process.traversal.Traversal;
import org.apache.tinkerpop.gremlin.process.traversal.Traverser;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal.Admin;
import org.apache.tinkerpop.gremlin.process.traversal.step.util.Tree;
import org.apache.tinkerpop.gremlin.process.traversal.util.TraversalMetrics;
import org.apache.tinkerpop.gremlin.structure.Column;
import org.apache.tinkerpop.gremlin.structure.Direction;
import org.apache.tinkerpop.gremlin.structure.Edge;
import org.apache.tinkerpop.gremlin.structure.Graph;
import org.apache.tinkerpop.gremlin.structure.Property;
import org.apache.tinkerpop.gremlin.structure.T;
import org.apache.tinkerpop.gremlin.structure.Vertex;
import org.apache.tinkerpop.gremlin.structure.VertexProperty.Cardinality;

public class DefaultGraphTraversal<S, E> extends DefaultTraversal<S, E> implements Admin<S, E> {

	public DefaultGraphTraversal() {
        super();
    }

    public DefaultGraphTraversal(final GraphTraversalSource graphTraversalSource) {
        super(graphTraversalSource);
        setConnection(graphTraversalSource.getConnection());
    }

    public DefaultGraphTraversal(final Graph graph) {
        super(graph);
    }

    @Override
    public GraphTraversal.Admin<S, E> asAdmin() {
        return this;
    }

    @Override
    public GraphTraversal<S, E> iterate() {
    	//FIXME: Do we need iterate?
        //return GraphTraversal.Admin.super.iterate();
        return this;
    }

    @Override
    public DefaultGraphTraversal<S, E> clone() {
        return (DefaultGraphTraversal<S, E>) super.clone();
    }

	@Override
	public <E2> GraphTraversal<S, E2> map(Function<Traverser<E>, E2> function) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.map, function);
        return (GraphTraversal<S, E2>) this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> map(Traversal<?, E2> mapTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.map, mapTraversal);
        return (GraphTraversal<S, E2>) this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> flatMap(Function<Traverser<E>, Iterator<E2>> function) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.flatMap, function);
		return (GraphTraversal<S, E2>) this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> flatMap(Traversal<?, E2> flatMapTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.flatMap, flatMapTraversal);
		return (GraphTraversal<S, E2>) this;
	}

	@Override
	public GraphTraversal<S, Object> id() {
		// TODO Auto-generated method stub
		return (GraphTraversal<S, Object>) this;
	}

	@Override
	public GraphTraversal<S, String> label() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.label);
		return (GraphTraversal<S, String>) this;
	}

	@Override
	public GraphTraversal<S, E> identity() {
		// TODO Auto-generated method stub
		return (GraphTraversal<S,E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> constant(E2 e) {
		// TODO Auto-generated method stub
		return (GraphTraversal<S, E2>) this;
	}

	@Override
	public GraphTraversal<S, Vertex> V(Object... vertexIdsOrElements) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.V, vertexIdsOrElements);
		return (GraphTraversal<S, Vertex>) this;
	}

	@Override
	public GraphTraversal<S, Vertex> to(Direction direction, String... edgeLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.to, direction, edgeLabels);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> out(String... edgeLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.out, (Object[])edgeLabels);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> in(String... edgeLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.in, (Object[])edgeLabels);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> both(String... edgeLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.both, (Object[])edgeLabels);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Edge> toE(Direction direction, String... edgeLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.toE, direction, edgeLabels);
		return (GraphTraversal<S, Edge>)this;
	}

	@Override
	public GraphTraversal<S, Edge> outE(String... edgeLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.outE, (Object[])edgeLabels);
		return (GraphTraversal<S, Edge>)this;
	}

	@Override
	public GraphTraversal<S, Edge> inE(String... edgeLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.inE, (Object[])edgeLabels);
		return (GraphTraversal<S, Edge>)this;
	}

	@Override
	public GraphTraversal<S, Edge> bothE(String... edgeLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.bothE, (Object[])edgeLabels);
		return (GraphTraversal<S, Edge>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> toV(Direction direction) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.toV, direction);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> inV() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.inV);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> outV() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.outV);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> bothV() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.bothV);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> otherV() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.otherV);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, E> order() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.order);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> order(Scope scope) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.order, scope);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, ? extends Property<E2>> properties(String... propertyKeys) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.properties, (Object[])propertyKeys);
		return (GraphTraversal<S, ? extends Property<E2>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> values(String... propertyKeys) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.values, (Object[])propertyKeys);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, Map<String, E2>> propertyMap(String... propertyKeys) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.propertyMap, (Object[])propertyKeys);
		return (GraphTraversal<S, Map<String, E2>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, Map<String, E2>> valueMap(String... propertyKeys) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.valueMap, (Object[])propertyKeys);
		return (GraphTraversal<S, Map<String, E2>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, Map<Object, E2>> valueMap(boolean includeTokens, String... propertyKeys) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.valueMap, includeTokens, propertyKeys);
		return (GraphTraversal<S, Map<Object, E2>>)this;
	}

	@Override
	public GraphTraversal<S, String> key() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.key);
		return (GraphTraversal<S, String>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> value() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.value);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, Path> path() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.path);
		return (GraphTraversal<S, Path>)this;
	}

	@Override
	public <E2> GraphTraversal<S, Map<String, E2>> match(Traversal<?, ?>... matchTraversals) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.match, (Object[])matchTraversals);
		return (GraphTraversal<S, Map<String, E2>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> sack() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.sack);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, Integer> loops() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.loops);
		return (GraphTraversal<S, Integer>)this;
	}

	@Override
	public <E2> GraphTraversal<S, Map<String, E2>> project(String projectKey, String... otherProjectKeys) {
		// TODO Auto-generated method stub
        final String[] projectKeys = new String[otherProjectKeys.length + 1];
        projectKeys[0] = projectKey;
        System.arraycopy(otherProjectKeys, 0, projectKeys, 1, otherProjectKeys.length);
        this.asAdmin().getBytecode().addStep(Symbols.project, projectKey, otherProjectKeys);
		return (GraphTraversal<S, Map<String, E2>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, Map<String, E2>> select(Pop pop, String selectKey1, String selectKey2,
			String... otherSelectKeys) {
		// TODO Auto-generated method stub
        final String[] selectKeys = new String[otherSelectKeys.length + 2];
        selectKeys[0] = selectKey1;
        selectKeys[1] = selectKey2;
        System.arraycopy(otherSelectKeys, 0, selectKeys, 2, otherSelectKeys.length);
        this.asAdmin().getBytecode().addStep(Symbols.select, pop, selectKey1, selectKey2, otherSelectKeys);
		return (GraphTraversal<S, Map<String, E2>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, Map<String, E2>> select(String selectKey1, String selectKey2,
			String... otherSelectKeys) {
		// TODO Auto-generated method stub
        final String[] selectKeys = new String[otherSelectKeys.length + 2];
        selectKeys[0] = selectKey1;
        selectKeys[1] = selectKey2;
        System.arraycopy(otherSelectKeys, 0, selectKeys, 2, otherSelectKeys.length);
        this.asAdmin().getBytecode().addStep(Symbols.select, selectKey1, selectKey2, otherSelectKeys);
		return (GraphTraversal<S, Map<String, E2>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> select(Pop pop, String selectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.select, pop, selectKey);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> select(String selectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.select, selectKey);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> select(Pop pop, Traversal<S, E2> keyTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.select, pop, keyTraversal);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> select(Traversal<S, E2> keyTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.select, keyTraversal);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, Collection<E2>> select(Column column) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.select, column);
		return (GraphTraversal<S, Collection<E2>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> unfold() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.unfold);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, List<E>> fold() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.fold);
		return (GraphTraversal<S, List<E>>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> fold(E2 seed, BiFunction<E2, E, E2> foldFunction) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.fold, seed, foldFunction);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, Long> count() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.count);
		return (GraphTraversal<S, Long>)this;
	}

	@Override
	public GraphTraversal<S, Long> count(Scope scope) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.count, scope);
		return (GraphTraversal<S, Long>)this;
	}

	@Override
	public <E2 extends Number> GraphTraversal<S, E2> sum() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.sum);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2 extends Number> GraphTraversal<S, E2> sum(Scope scope) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.sum, scope);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2 extends Number> GraphTraversal<S, E2> max() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.max);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2 extends Number> GraphTraversal<S, E2> max(Scope scope) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.max, scope);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2 extends Number> GraphTraversal<S, E2> min() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.min);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2 extends Number> GraphTraversal<S, E2> min(Scope scope) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.min, scope);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2 extends Number> GraphTraversal<S, E2> mean() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.mean);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2 extends Number> GraphTraversal<S, E2> mean(Scope scope) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.mean, scope);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <K, V> GraphTraversal<S, Map<K, V>> group() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.group);
		return (GraphTraversal<S, Map<K, V>>)this;
	}

	@Override
	public <K> GraphTraversal<S, Map<K, Long>> groupCount() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.groupCount);
		return (GraphTraversal<S, Map<K, Long>>)this;
	}

	@Override
	public GraphTraversal<S, Tree> tree() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.tree);
		return (GraphTraversal<S, Tree>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> addV(String vertexLabel) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.addV, vertexLabel);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> addV(Traversal<?, String> vertexLabelTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.addV, vertexLabelTraversal);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Vertex> addV() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.addV);
		return (GraphTraversal<S, Vertex>)this;
	}

	@Override
	public GraphTraversal<S, Edge> addE(String edgeLabel) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.addE, edgeLabel);
		return (GraphTraversal<S, Edge>)this;
	}

	@Override
	public GraphTraversal<S, Edge> addE(Traversal<?, String> edgeLabelTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.addE, edgeLabelTraversal);
		return (GraphTraversal<S, Edge>)this;
	}

	@Override
	public GraphTraversal<S, E> to(String toStepLabel) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.to, toStepLabel);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> from(String fromStepLabel) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.from, fromStepLabel);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> to(Traversal<?, Vertex> toVertex) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.to, toVertex);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> from(Traversal<?, Vertex> fromVertex) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.from, fromVertex);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> to(Vertex toVertex) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.to, toVertex);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> from(Vertex fromVertex) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.from, fromVertex);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, Double> math(String expression) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.math, expression);
		return (GraphTraversal<S, Double>)this;
	}

	@Override
	public GraphTraversal<S, E> filter(Predicate<Traverser<E>> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.filter, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> filter(Traversal<?, ?> filterTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.filter, filterTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> or(Traversal<?, ?>... orTraversals) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.or, (Object[])orTraversals);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> and(Traversal<?, ?>... andTraversals) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.and, (Object[])andTraversals);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> inject(E... injections) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.inject, injections);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> dedup(Scope scope, String... dedupLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.dedup, scope, dedupLabels);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> dedup(String... dedupLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.dedup, (Object[])dedupLabels);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> where(String startKey, P<String> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.where, startKey, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> where(P<String> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.where, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> where(Traversal<?, ?> whereTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.where, whereTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> has(String propertyKey, P<?> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.has, propertyKey, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> has(T accessor, P<?> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.has, accessor, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> has(String propertyKey, Object value) {
		// TODO Auto-generated method stub
        if (value instanceof P)
            return this.has(propertyKey, (P) value);
        else if (value instanceof Traversal)
            return this.has(propertyKey, (Traversal) value);
        else {
            this.asAdmin().getBytecode().addStep(Symbols.has, propertyKey, value);
            return (GraphTraversal<S, E>)this;
        }
	}

	@Override
	public GraphTraversal<S, E> has(T accessor, Object value) {
		// TODO Auto-generated method stub
        if (value instanceof P)
            return this.has(accessor, (P) value);
        else if (value instanceof Traversal)
            return this.has(accessor, (Traversal) value);
        else {
            this.asAdmin().getBytecode().addStep(Symbols.has, accessor, value);
            return (GraphTraversal<S, E>)this;
        }
	}

	@Override
	public GraphTraversal<S, E> has(String label, String propertyKey, P<?> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.has, label, propertyKey, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> has(String label, String propertyKey, Object value) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.has, label, propertyKey, value);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> has(T accessor, Traversal<?, ?> propertyTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.has, accessor, propertyTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> has(String propertyKey, Traversal<?, ?> propertyTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.has, propertyKey, propertyTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> has(String propertyKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.has, propertyKey);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> hasNot(String propertyKey) {
		// TODO Auto-generated method stub
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> hasLabel(String label, String... otherLabels) {
		// TODO Auto-generated method stub
        final String[] labels = new String[otherLabels.length + 1];
        labels[0] = label;
        System.arraycopy(otherLabels, 0, labels, 1, otherLabels.length);
        this.asAdmin().getBytecode().addStep(Symbols.hasLabel, (Object[])labels);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> hasLabel(P<String> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.hasLabel, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> hasId(Object id, Object... otherIds) {
		// TODO Auto-generated method stub
        if (id instanceof P)
            return this.hasId((P) id);
        else {
            final List<Object> ids = new ArrayList<>();
            if (id instanceof Object[]) {
                for (final Object i : (Object[]) id) {
                    ids.add(i);
                }
            } else
                ids.add(id);
            for (final Object i : otherIds) {
                if (i.getClass().isArray()) {
                    for (final Object ii : (Object[]) i) {
                        ids.add(ii);
                    }
                } else
                    ids.add(i);
            }
            this.asAdmin().getBytecode().addStep(Symbols.hasId, ids.toArray());
		    return (GraphTraversal<S, E>)this;
        }
	}

	@Override
	public GraphTraversal<S, E> hasId(P<Object> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.hasId, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> hasKey(String label, String... otherLabels) {
		// TODO Auto-generated method stub
        final String[] labels = new String[otherLabels.length + 1];
        labels[0] = label;
        System.arraycopy(otherLabels, 0, labels, 1, otherLabels.length);
        this.asAdmin().getBytecode().addStep(Symbols.hasKey, (Object[])labels);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> hasKey(P<String> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.hasKey, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> hasValue(Object value, Object... otherValues) {
		// TODO Auto-generated method stub
        if (value instanceof P)
            return this.hasValue((P) value);
        else {
            final List<Object> values = new ArrayList<>();
            if (value instanceof Object[]) {
                for (final Object v : (Object[]) value) {
                    values.add(v);
                }
            } else
                values.add(value);
            for (final Object v : otherValues) {
                if (v instanceof Object[]) {
                    for (final Object vv : (Object[]) v) {
                        values.add(vv);
                    }
                } else
                    values.add(v);
            }
            this.asAdmin().getBytecode().addStep(Symbols.hasValue, values.toArray());
		    return (GraphTraversal<S, E>)this;
        }
	}

	@Override
	public GraphTraversal<S, E> hasValue(P<Object> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.hasValue, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> is(P<E> predicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.is, predicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> is(Object value) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.is, value);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> not(Traversal<?, ?> notTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.not, notTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> coin(double probability) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.coin, probability);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> range(long low, long high) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.range, low, high);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> range(Scope scope, long low, long high) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.range, scope, low, high);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, E> limit(long limit) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.limit, limit);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> limit(Scope scope, long limit) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.limit, scope, limit);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, E> tail() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.tail);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> tail(long limit) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.tail, limit);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> tail(Scope scope) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.tail, scope);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> tail(Scope scope, long limit) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.tail, scope, limit);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, E> skip(long skip) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.skip, skip);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> skip(Scope scope, long skip) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.skip, scope, skip);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, E> timeLimit(long timeLimit) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.timeLimit, timeLimit);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> simplePath() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.simplePath);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> cyclicPath() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.cyclicPath);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> sample(int amountToSample) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.sample, amountToSample);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> sample(Scope scope, int amountToSample) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.sample, scope, amountToSample);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> drop() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.drop);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> sideEffect(Consumer<Traverser<E>> consumer) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.sideEffect, consumer);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> sideEffect(Traversal<?, ?> sideEffectTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.sideEffect, sideEffectTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> cap(String sideEffectKey, String... sideEffectKeys) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.cap, sideEffectKey, sideEffectKeys);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, Edge> subgraph(String sideEffectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.subgraph, sideEffectKey);
		return (GraphTraversal<S, Edge>)this;
	}

	@Override
	public GraphTraversal<S, E> aggregate(String sideEffectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.aggregate, sideEffectKey);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> group(String sideEffectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.group, sideEffectKey);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> groupCount(String sideEffectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.groupCount, sideEffectKey);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> tree(String sideEffectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.tree, sideEffectKey);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <V, U> GraphTraversal<S, E> sack(BiFunction<V, U, V> sackOperator) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.sack, sackOperator);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> store(String sideEffectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.store, sideEffectKey);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> profile(String sideEffectKey) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Traversal.Symbols.profile, sideEffectKey);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, TraversalMetrics> profile() {
		// TODO Auto-generated method stub
        // FIXME: Will this work?
        return (GraphTraversal<S, TraversalMetrics>) this;
	}

	@Override
	public GraphTraversal<S, E> property(Cardinality cardinality, Object key, Object value, Object... keyValues) {
		// TODO Auto-generated method stub
        if (null == cardinality)
            this.asAdmin().getBytecode().addStep(Symbols.property, key, value, keyValues);
        else
            this.asAdmin().getBytecode().addStep(Symbols.property, cardinality, key, value, keyValues);
        return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> property(Object key, Object value, Object... keyValues) {
		// TODO Auto-generated method stub
        return key instanceof Cardinality ?
                this.property((Cardinality) key, value, keyValues[0],
                        keyValues.length > 1 ?
                                Arrays.copyOfRange(keyValues, 1, keyValues.length) :
                                new Object[]{}) :
                this.property(null, key, value, keyValues);
	}

	@Override
	public <M, E2> GraphTraversal<S, E2> branch(Traversal<?, M> branchTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.branch, branchTraversal);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <M, E2> GraphTraversal<S, E2> branch(Function<Traverser<E>, M> function) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.branch, function);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <M, E2> GraphTraversal<S, E2> choose(Traversal<?, M> choiceTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.choose, choiceTraversal);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> choose(Traversal<?, ?> traversalPredicate, Traversal<?, E2> trueChoice,
			Traversal<?, E2> falseChoice) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.choose, traversalPredicate, trueChoice, falseChoice);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> choose(Traversal<?, ?> traversalPredicate, Traversal<?, E2> trueChoice) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.choose, traversalPredicate, trueChoice);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <M, E2> GraphTraversal<S, E2> choose(Function<E, M> choiceFunction) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.choose, choiceFunction);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> choose(Predicate<E> choosePredicate, Traversal<?, E2> trueChoice,
			Traversal<?, E2> falseChoice) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.choose, choosePredicate, trueChoice, falseChoice);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> choose(Predicate<E> choosePredicate, Traversal<?, E2> trueChoice) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.choose, choosePredicate, trueChoice);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> optional(Traversal<?, E2> optionalTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.optional, optionalTraversal);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> union(Traversal<?, E2>... unionTraversals) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.union, (Object[])unionTraversals);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> coalesce(Traversal<?, E2>... coalesceTraversals) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.coalesce, (Object[])coalesceTraversals);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, E> repeat(Traversal<?, E> repeatTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.repeat, repeatTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> emit(Traversal<?, ?> emitTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.emit, emitTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> emit(Predicate<Traverser<E>> emitPredicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.emit, emitPredicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> emit() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.emit);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> until(Traversal<?, ?> untilTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.until, untilTraversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> until(Predicate<Traverser<E>> untilPredicate) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.until, untilPredicate);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> times(int maxLoops) {
		// TODO Auto-generated method stub
		//Time modulating or Loop travsersal
        this.asAdmin().getBytecode().addStep(Symbols.times, maxLoops);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E2> local(Traversal<?, E2> localTraversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.local, localTraversal);
		return (GraphTraversal<S, E2>)this;
	}

	@Override
	public GraphTraversal<S, E> pageRank() {
		return pageRank(0.20);  //Standard value provide nvGraph

	}

	@Override
	public GraphTraversal<S, E> pageRank(double alpha) {
		// TODO Auto-generated method stub
        return pageRank(alpha, 0.5f, 500);
	}

	@Override
	public GraphTraversal<S,E> pageRank(double alpha, float tolerance, int maxiter)
	{
		this.asAdmin().getBytecode().addStep(Symbols.pageRank, alpha, tolerance, maxiter);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> sssp(String costlabel) {
		this.asAdmin().getBytecode().addStep(Symbols.sssp, costlabel);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> spectralClustering(String weightlabel, int numberOfCluster, ClusteringAlgorithm algorithm) {
		this.asAdmin().getBytecode().addStep(Symbols.spectralClustering, weightlabel, numberOfCluster, algorithm);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> peerPressure() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.peerPressure);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> program(VertexProgram<?> vertexProgram) {
		// TODO Auto-generated method stub
		//FIXME: There is no byte code for Program
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> as(String stepLabel, String... stepLabels) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.as, stepLabel, stepLabels);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> barrier() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.barrier);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> barrier(int maxBarrierSize) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.barrier, maxBarrierSize);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> barrier(Consumer<Set<Object>> barrierConsumer) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.barrier, barrierConsumer);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> by() {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> by(Traversal<?, ?> traversal) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, traversal);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> by(T token) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, token);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> by(String key) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, key);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <V> GraphTraversal<S, E> by(Function<V, Object> function) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, function);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <V> GraphTraversal<S, E> by(Traversal<?, ?> traversal, Comparator<V> comparator) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, traversal, comparator);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> by(Comparator<E> comparator) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, comparator);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public GraphTraversal<S, E> by(Order order) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, order);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <V> GraphTraversal<S, E> by(String key, Comparator<V> comparator) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, key, comparator);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <U> GraphTraversal<S, E> by(Function<U, Object> function, Comparator comparator) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.by, function, comparator);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <M, E2> GraphTraversal<S, E> option(M pickToken, Traversal<?, E2> traversalOption) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.option, pickToken, traversalOption);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> GraphTraversal<S, E> option(Traversal<?, E2> traversalOption) {
		// TODO Auto-generated method stub
        this.asAdmin().getBytecode().addStep(Symbols.option, traversalOption);
		return (GraphTraversal<S, E>)this;
	}

	@Override
	public <E2> org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal.Admin<S, E2> addStep(
			Step<?, E2> step) {
		// TODO Auto-generated method stub
		return (GraphTraversal.Admin<S, E2>)this;
	}

	/*
	@Override
	public GraphTraversal<S, E> iterate() {
		// TODO Auto-generated method stub
		return this;
	}

	@Override
	public org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal.Admin<S, E> clone() {
		// TODO Auto-generated method stub
		return this;
	}
	*/

    public static final class Symbols {

        private Symbols() {
            // static fields only
        }

        public static final String map = "map";
        public static final String flatMap = "flatMap";
        public static final String id = "id";
        public static final String label = "label";
        public static final String identity = "identity";
        public static final String constant = "constant";
        public static final String V = "V";
        public static final String E = "E";
        public static final String to = "to";
        public static final String out = "out";
        public static final String in = "in";
        public static final String both = "both";
        public static final String toE = "toE";
        public static final String outE = "outE";
        public static final String inE = "inE";
        public static final String bothE = "bothE";
        public static final String toV = "toV";
        public static final String outV = "outV";
        public static final String inV = "inV";
        public static final String bothV = "bothV";
        public static final String otherV = "otherV";
        public static final String order = "order";
        public static final String properties = "properties";
        public static final String values = "values";
        public static final String propertyMap = "propertyMap";
        public static final String valueMap = "valueMap";
        public static final String select = "select";
        public static final String key = "key";
        public static final String value = "value";
        public static final String path = "path";
        public static final String match = "match";
        public static final String math = "math";
        public static final String sack = "sack";
        public static final String loops = "loops";
        public static final String project = "project";
        public static final String unfold = "unfold";
        public static final String fold = "fold";
        public static final String count = "count";
        public static final String sum = "sum";
        public static final String max = "max";
        public static final String min = "min";
        public static final String mean = "mean";
        public static final String group = "group";
        public static final String groupCount = "groupCount";
        public static final String tree = "tree";
        public static final String addV = "addV";
        public static final String addE = "addE";
        public static final String from = "from";
        public static final String filter = "filter";
        public static final String or = "or";
        public static final String and = "and";
        public static final String inject = "inject";
        public static final String dedup = "dedup";
        public static final String where = "where";
        public static final String has = "has";
        public static final String hasNot = "hasNot";
        public static final String hasLabel = "hasLabel";
        public static final String hasId = "hasId";
        public static final String hasKey = "hasKey";
        public static final String hasValue = "hasValue";
        public static final String is = "is";
        public static final String not = "not";
        public static final String range = "range";
        public static final String limit = "limit";
        public static final String skip = "skip";
        public static final String tail = "tail";
        public static final String coin = "coin";

        public static final String timeLimit = "timeLimit";
        public static final String simplePath = "simplePath";
        public static final String cyclicPath = "cyclicPath";
        public static final String sample = "sample";

        public static final String drop = "drop";

        public static final String sideEffect = "sideEffect";
        public static final String cap = "cap";
        public static final String property = "property";
        public static final String store = "store";
        public static final String aggregate = "aggregate";
        public static final String subgraph = "subgraph";
        public static final String barrier = "barrier";
        public static final String local = "local";
        public static final String emit = "emit";
        public static final String repeat = "repeat";
        public static final String until = "until";
        public static final String branch = "branch";
        public static final String union = "union";
        public static final String coalesce = "coalesce";
        public static final String choose = "choose";
        public static final String optional = "optional";


        public static final String pageRank = "pageRank";
        public static final String peerPressure = "peerPressure";
        public static final String program = "program";
		public static final String spectralClustering = "spectralClustering";
		public static final String sssp = "sssp";

        public static final String by = "by";
        public static final String times = "times";
        public static final String as = "as";
        public static final String option = "option";

    }
}
