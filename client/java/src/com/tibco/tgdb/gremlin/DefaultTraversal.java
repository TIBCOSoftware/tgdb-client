package com.tibco.tgdb.gremlin;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashSet;
import java.util.Iterator;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.stream.Stream;

import org.apache.tinkerpop.gremlin.process.traversal.Bytecode;
import org.apache.tinkerpop.gremlin.process.traversal.Step;
import org.apache.tinkerpop.gremlin.process.traversal.Traversal;
import org.apache.tinkerpop.gremlin.process.traversal.TraversalSideEffects;
import org.apache.tinkerpop.gremlin.process.traversal.TraversalSource;
import org.apache.tinkerpop.gremlin.process.traversal.TraversalStrategies;
import org.apache.tinkerpop.gremlin.process.traversal.TraverserGenerator;
import org.apache.tinkerpop.gremlin.process.traversal.step.TraversalParent;
import org.apache.tinkerpop.gremlin.process.traversal.traverser.TraverserRequirement;
import org.apache.tinkerpop.gremlin.process.traversal.util.TraversalExplanation;
import org.apache.tinkerpop.gremlin.structure.Graph;
import org.apache.tinkerpop.gremlin.structure.util.empty.EmptyGraph;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.query.TGResultSet;

public class DefaultTraversal<S, E> implements Traversal.Admin<S, E> {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    protected transient Graph graph;
	protected transient TGConnection connection;
    protected TraversalStrategies strategies;
    protected Bytecode bytecode; // TODO: perhaps make transient until 3.3.0?
    
    private DefaultTraversal(final Graph graph, final TraversalStrategies traversalStrategies, final Bytecode bytecode) {
        this.graph = graph;
        this.strategies = traversalStrategies;
        this.bytecode = bytecode;
    }

    public DefaultTraversal(final Graph graph) {
        this(graph, null, new Bytecode());
    }

    public DefaultTraversal(final TraversalSource traversalSource) {
        this(traversalSource.getGraph(), traversalSource.getStrategies(), traversalSource.getBytecode());
    }

    public DefaultTraversal(final TraversalSource traversalSource, final DefaultTraversal.Admin<S,E> traversal) {
        this(traversalSource.getGraph(), traversalSource.getStrategies(), traversal.getBytecode());
//        steps.addAll(traversal.getSteps());
    }

    // TODO: clean up unused or redundant constructors

    public DefaultTraversal() {
        this(EmptyGraph.instance(), null, new Bytecode());
    }

    public DefaultTraversal(final Bytecode bytecode) {
        this(EmptyGraph.instance(), null, bytecode);
    }

	@Override
	public boolean hasNext() {
		// TODO Auto-generated method stub
		return false;
	}

	@Override
	public E next() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Admin asAdmin() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Optional tryNext() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public List next(int amount) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public List toList() {
		// TODO Auto-generated method stub
        return this.fill(new ArrayList<>());
	}

	@Override
	public Set toSet() {
		// TODO Auto-generated method stub
        return this.fill(new HashSet<>());
	}

	@Override
	public Set toBulkSet() {
		// TODO Auto-generated method stub
        return null;
	}

	@Override
	public Stream toStream() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public CompletableFuture promise(Function traversalFunction) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <C extends Collection<E>> C fill(final C collection) {
		// TODO Auto-generated method stub
		//connection.executeQuery(, option)
		try {
			((ConnectionImpl)connection).executeGremlinQuery(bytecode.toString(), collection, TGQueryOption.DEFAULT_QUERY_OPTION);
			return collection;
		} catch (TGException e) {
			gLogger.logException("Gremlin query failed", e);
		}
		return collection;
	}

	@Override
	public Traversal iterate() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Traversal none() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Traversal profile() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public TraversalExplanation explain() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public void forEachRemaining(Class endType, Consumer consumer) {
		// TODO Auto-generated method stub

	}

	@Override
	public void forEachRemaining(Consumer action) {
		// TODO Auto-generated method stub

	}

	@Override
	public void close() throws Exception {
		// TODO Auto-generated method stub

	}

	@Override
	public Bytecode getBytecode() {
		// TODO Auto-generated method stub
		return bytecode;
	}

	public void setConnection(TGConnection conn) {
		this.connection = conn;
	}

	public TGConnection getConnection() {
		return this.connection;
	}

	@Override
	public void addStarts(Iterator<org.apache.tinkerpop.gremlin.process.traversal.Traverser.Admin<S>> starts) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void addStart(org.apache.tinkerpop.gremlin.process.traversal.Traverser.Admin<S> start) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public List<Step> getSteps() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <E2> Admin<S, E2> addStep(Step<?, E2> step) throws IllegalStateException {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <S2, E2> Admin<S2, E2> addStep(int index, Step<?, ?> step) throws IllegalStateException {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <S2, E2> Admin<S2, E2> removeStep(Step<?, ?> step) throws IllegalStateException {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public <S2, E2> Admin<S2, E2> removeStep(int index) throws IllegalStateException {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Step<S, ?> getStartStep() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Step<?, E> getEndStep() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public void applyStrategies() throws IllegalStateException {
		// TODO Auto-generated method stub
		
	}

	@Override
	public TraverserGenerator getTraverserGenerator() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Set<TraverserRequirement> getTraverserRequirements() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public void setSideEffects(TraversalSideEffects sideEffects) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public TraversalSideEffects getSideEffects() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public void setStrategies(TraversalStrategies strategies) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public TraversalStrategies getStrategies() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public void setParent(TraversalParent step) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public TraversalParent getParent() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public Admin<S, E> clone() {
		// TODO Auto-generated method stub
        try {
            final DefaultTraversal<S, E> clone = (DefaultTraversal<S, E>) super.clone();
//            clone.lastTraverser = EmptyTraverser.instance();
//            clone.steps = new ArrayList<>();
//            clone.unmodifiableSteps = Collections.unmodifiableList(clone.steps);
//            clone.sideEffects = this.sideEffects.clone();
            clone.strategies = this.strategies;
            clone.bytecode = this.bytecode.clone();
            /*
            for (final Step<?, ?> step : this.steps) {
                final Step<?, ?> clonedStep = step.clone();
                clonedStep.setTraversal(clone);
                final Step previousStep = clone.steps.isEmpty() ? EmptyStep.instance() : clone.steps.get(clone.steps.size() - 1);
                clonedStep.setPreviousStep(previousStep);
                previousStep.setNextStep(clonedStep);
                clone.steps.add(clonedStep);
            }
            */
//            clone.finalEndStep = clone.getEndStep();
            return clone;
        } catch (final CloneNotSupportedException e) {
            throw new IllegalStateException(e.getMessage(), e);
        }
		//return null;
	}

	@Override
	public boolean isLocked() {
		// TODO Auto-generated method stub
		return false;
	}

	@Override
	public Optional<Graph> getGraph() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public void setGraph(Graph graph) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public boolean equals(Admin<S, E> other) {
		// TODO Auto-generated method stub
		return false;
	}

	@Override
	public org.apache.tinkerpop.gremlin.process.traversal.Traverser.Admin<E> nextTraverser() {
		// TODO Auto-generated method stub
		return null;
	}

}
