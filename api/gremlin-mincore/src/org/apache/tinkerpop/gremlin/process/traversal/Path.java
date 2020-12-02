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


import org.apache.tinkerpop.gremlin.process.computer.KeyValue;
import org.apache.tinkerpop.gremlin.structure.Graph;


import java.util.ArrayList;
import java.util.Collections;
import java.util.Iterator;
import java.util.List;
import java.util.Set;
import java.util.function.BiConsumer;
import java.util.stream.IntStream;
import java.util.stream.Stream;

/**
 * A Path denotes a particular walk through a {@link Graph} as defined by a {@link Traversal}.
 * In abstraction, any Path implementation maintains two lists: a list of sets of labels and a list of objects.
 * The list of labels are the labels of the steps traversed. The list of objects are the objects traversed.
 *
 * @author Marko A. Rodriguez (http://markorodriguez.com)
 */
public interface Path extends Cloneable, Iterable<Object> {

    /**
     * Get the number of step in the path.
     *
     * @return the size of the path
     */
      int size();

    /**
     * Determine if the path is empty or not.
     *
     * @return whether the path is empty or not.
     */
      boolean isEmpty();

    /**
     * Get the head of the path.
     *
     * @param <A> the type of the head of the path
     * @return the head of the path
     */
      <A> A head();

    /**
     * Add a new step to the path with an object and any number of associated labels.
     *
     * @param object the new head of the path
     * @param labels the labels at the head of the path
     * @return the extended path
     */
     Path extend(final Object object, final Set<String> labels);

    /**
     * Add labels to the head of the path.
     *
     * @param labels the labels at the head of the path
     * @return the path with added labels
     */
     Path extend(final Set<String> labels);

    /**
     * Remove labels from path.
     *
     * @param labels the labels to remove
     * @return the path with removed labels
     */
     Path retract(final Set<String> labels);

    /**
     * Get the object associated with the particular label of the path.
     * If the path as multiple labels of the type, then return a {@link List} of those objects.
     *
     * @param label the label of the path
     * @param <A>   the type of the object associated with the label
     * @return the object associated with the label of the path
     * @throws IllegalArgumentException if the path does not contain the label
     */
      <A> A get(final String label) throws IllegalArgumentException;


    /**
     * Pop the object(s) associated with the label of the path.
     *
     * @param pop   first for least recent, last for most recent, and all for all in a list
     * @param label the label of the path
     * @param <A>   the type of the object associated with the label
     * @return the object associated with the label of the path
     * @throws IllegalArgumentException if the path does not contain the label
     */
      <A> A get(final Pop pop, final String label) throws IllegalArgumentException;


    /**
     * Get the object associated with the specified index into the path.
     *
     * @param index the index of the path
     * @param <A>   the type of the object associated with the index
     * @return the object associated with the index of the path
     */
      <A> A get(final int index);

    /**
     * Return true if the path has the specified label, else return false.
     *
     * @param label the label to search for
     * @return true if the label exists in the path
     */
      boolean hasLabel(final String label);

    /**
     * An ordered list of the objects in the path.
     *
     * @return the objects of the path
     */
     List<Object> objects();

    /**
     * An ordered list of the labels associated with the path
     * The set of labels for a particular step are ordered by the order in which {@link Path#extend(Object, Set)} was called.
     *
     * @return the labels of the path
     */
     List<Set<String>> labels();

    @SuppressWarnings("CloneDoesntDeclareCloneNotSupportedException")
     Path clone();

    /**
     * Determines whether the path is a simple or not.
     * A simple path has no cycles and thus, no repeated objects.
     *
     * @return Whether the path is simple or not
     */
      boolean isSimple();

      Iterator<Object> iterator();

      void forEach(final BiConsumer<Object, Set<String>> consumer);

      Stream<KeyValue<Object, Set<String>>> stream();

      boolean popEquals(final Pop pop, final Object other);


    /**
     * Isolate a sub-path from the path object. The isolation is based solely on the path labels.
     * The to-label is inclusive. Thus, from "b" to "c" would isolate the example path as follows {@code a,[b,c],d}.
     * Note that if there are multiple path segments with the same label, then its the last occurrence that is isolated.
     * For instance, from "b" to "c" would be {@code a,b,[b,c,d,c]}.
     *
     * @param fromLabel The label to start recording the sub-path from.
     * @param toLabel   The label to end recording the sub-path to.
     * @return the isolated sub-path.
     */
      Path subPath(final String fromLabel, final String toLabel);

    public static class Exceptions {

        public static IllegalArgumentException stepWithProvidedLabelDoesNotExist(final String label) {
            return new IllegalArgumentException("The step with label " + label + " does not exist");
        }

        public static IllegalArgumentException couldNotLocatePathFromLabel(final String fromLabel) {
            return new IllegalArgumentException("Could not locate path from-label: " + fromLabel);
        }

        public static IllegalArgumentException couldNotLocatePathToLabel(final String toLabel) {
            return new IllegalArgumentException("Could not locate path to-label: " + toLabel);
        }

        public static IllegalArgumentException couldNotIsolatedSubPath(final String fromLabel, final String toLabel) {
            return new IllegalArgumentException("Could not isolate path because from comes after to: " + fromLabel + "->" + toLabel);
        }
    }
}
