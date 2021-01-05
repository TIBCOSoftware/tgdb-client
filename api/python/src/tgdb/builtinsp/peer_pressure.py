"""
 * Copyright (c) 2020 TIBCO Software Inc. All rights reserved.
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
 *   File name   :peer_pressure.py
 *   Created on  :8/31/20.
 *   Created by  :dhudson
 *   SVN Id      :$Id:$
"""

import tgdb.storedproc as tgsp
import tgdb.model as tgmodel
import tgdb.exception as tgexcept
import typing
import queue
import ctypes as ct
import tgdb.builtinsp.generic as generics


def _computeVertexWeights(nodes: typing.List[tgmodel.TGNode],
                          edges: typing.Union[typing.Dict[tgmodel.TGNode, typing.Iterable[tgmodel.TGEdge]], str],
                          weight: typing.Union[str, float], defaultWeight: float) -> typing.Dict[tgmodel.TGNode, float]:
    ret: typing.Dict[tgmodel.TGNode, float] = {}
    for node in nodes:
        curWeight = 0.0
        for edge in generics.getEdges(node, edges, weight, defaultWeight, keepMax=True):
            if edge.getVertices()[0] == node:
                curWeight += generics.getEdgeDist(edge, weight, defaultWeight)
        if curWeight > 0.0:
            ret[node] = 1.0 / curWeight
        else:
            ret[node] = 0.0
    return ret


def _computeVertexCluster(nodes: typing.Dict[tgmodel.TGNode, int],
                          edges: typing.Union[typing.Dict[tgmodel.TGNode, typing.Iterable[tgmodel.TGEdge]], str],
                          nodeWeights: typing.Dict[tgmodel.TGNode, float], weight: typing.Union[str, float],
                          defaultWeight: float) -> typing.Dict[tgmodel.TGNode, int]:
    ret: typing.Dict[tgmodel.TGNode, int] = {}
    for node in nodes:
        curNodeWeight: typing.Dict[int, float] = {}
        for edge in generics.getEdges(node, edges, weight, defaultWeight, keepMax=True):
            otherNode = generics.otherVertex(edge, node)
            otherWeight = nodeWeights[otherNode] * generics.getEdgeDist(edge, weight, defaultWeight)
            if otherNode in curNodeWeight:
                curNodeWeight[nodes[otherNode]] += otherWeight
            else:
                curNodeWeight[nodes[otherNode]] = otherWeight
        nextNodeWeights = [(cluster, curNodeWeight[cluster]) for cluster in curNodeWeight]
        nextNodeWeights = list(sorted(nextNodeWeights, key=lambda nextNodeWeight: nextNodeWeights[0]))
        if len(nextNodeWeights) > 0:
            ret[node] = list(sorted(nextNodeWeights, key=lambda nextNodeWeight: nextNodeWeight[1]))[-1][0]
        else:
            ret[node] = nodes[node]     # Don't change the vertex cluster if it has no influencers.

    return ret


@tgsp.tgstoredproc
def peerPressure(g: tgsp.TGGraphContext,
                 source: str = "",
                 weight: str = "1.0",
                 defaultWeight: float = 0.0,
                 maxIterations: int = 50,
                 edges: str = "IN"
                 ) -> "[(V,l)]":
    """
    Run the peer pressure graph algorithm on the sources to determine clustering based on the parameter arguments
    :param g: The graph context that the server passes into all stored procedures.
    :param source: Optionally, a query representing what nodes should form the basis of the peer pressure.
    :type source: typing.Optional[str] A string representing a query.
    :param weight: The weight to use for the edges. Only makes sense to set for when it is a string that determines the
    attribute name to use for the basis.
    :type weight: typing.Union[str, float] A string or a float.
    :param defaultWeight: What to use when the edge does not have the weight attribute specified. Only used when weight
    is set to an attribute name.
    :type defaultWeight: A floating point.
    :param maxIterations: The maximum iterations to run through the graph.
    :type maxIterations: An integer.
    :param edges: The edges to determine the subgraph that the peer pressure runs on. "IN" indicates that only incoming
    edges affect a given vertex's cluster assignment, "OUT" indicates that only outgoing edges affect a given vertex's
    cluster assignment, "BOTH" indicates that any edge incident on a given vertex affects that vertex's assignment.
    :type edges: typing.Union[str, typing.Set[tgmodel.TGEdge], tgsp.TGGraphContext] A string representing either a
    direction type or a gremlin query.
    :return:
    """
    # TODO Remove all logging references and indicators...
    source: typing.Optional[str]
    weight: typing.Union[str, float]
    sourceNodes: typing.List[tgmodel.TGNode]
    edges: typing.Union[str, typing.Set[tgmodel.TGEdge], tgsp.TGGraphContext,
                        typing.Dict[tgmodel.TGNode, typing.Iterable[tgmodel.TGEdge]]]
    generics._log("Starting Peer Pressure...\n", mode='w')
    # source = None
    # weight = 1.0
    # defaultWeight = 0.0
    # maxIterations = 50
    # edges = "OUT"
    try:
        if source is None or source == '':
            sourceNodes = list(g.getNodes())
        else:
            sourceNodes = list(g.executeQuery(source).toCollection())
        sourceSet = set(g.getNodes())
        if isinstance(edges, tgsp.TGGraphContext):
            edges = set(edges.getEdges())
        elif edges not in generics.EDGE_TYPES:
            edges = set(g.executeQuery(edges).toCollection())
        else:
            tmpEdges = {}
            for node in sourceNodes:
                tmpE = []
                for edge in generics.getEdges(node, edges, weight, defaultWeight, keepMax=True):
                    if generics.otherVertex(edge, node) in sourceSet:
                        tmpE.append(edge)
                generics._log("Node: %s has %d edges...\n" % (node.getAttribute("airportID").getValue(), len(tmpE)))
                tmpEdges[node] = tmpE
            edges = tmpEdges

        vertWeights = _computeVertexWeights(sourceNodes, edges, weight, defaultWeight)
        edges = "IN" if edges == "OUT" else "OUT" if edges == "IN" else edges       # Invert the edge direction.

        clusters: typing.Dict[tgmodel.TGNode, int] = {}
        for i in range(len(sourceNodes)):
            clusters[sourceNodes[i]] = i

        for iteration in range(maxIterations):
            newClusters = _computeVertexCluster(clusters, edges, vertWeights, weight, defaultWeight)

            different: bool = False
            for node in newClusters:
                different = newClusters[node] != clusters[node]
                if different:
                    break

            clusters = newClusters
            if not different:
                break

        ret = [(node, clusters[node]) for node in clusters]
        for node, cluster in ret:
            generics._log("Node %s has cluster %d\n" % (node.getAttribute("airportID").getValue(), cluster))
        return ret
    except BaseException as e:
        import traceback as tb
        with open('tb_file_print.txt', 'a') as tb_file:
            tb.print_tb(e.__traceback__, file=tb_file)

        raise e
