package org.apache.tinkerpop.gremlin.process.traversal.util;

/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
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

import org.apache.tinkerpop.gremlin.process.computer.KeyValue;
import org.apache.tinkerpop.gremlin.process.traversal.Traversal;
import org.apache.tinkerpop.gremlin.process.traversal.TraversalStrategy;

import java.util.List;

/**
 * A TraversalExplanation takes a {@link Traversal} and, for each registered {@link TraversalStrategy}, it creates a
 * mapping reflecting how each strategy alters the traversal. This is useful for understanding how each traversal
 * strategy mutates the traversal. This is useful in debugging and analysis of traversal compilation. The
 * {@link TraversalExplanation#toString()} has a pretty-print representation that is useful in the Gremlin Console.
 *
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 * @author suresh (suresh.subramani@tibco.com)
 * Changed this to a Interface instead of concrete implementation.
 */


public interface TraversalExplanation {

    /**
     * Get the list of {@link TraversalStrategy} applications. For strategy, the resultant mutated {@link Traversal} is provided.
     *
     * @return the list of strategy/traversal pairs
     */
    List<KeyValue<TraversalStrategy, Traversal.Admin<?,?>>> getStrategyTraversals();
    /**
     * A pretty-print representation of the traversal explanation.
     *
     * @return a {@link String} representation of the traversal explanation
     */
    String prettyPrint(int maxLineLength);

}
