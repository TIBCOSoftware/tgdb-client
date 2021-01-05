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
 *   File name   :page_rank.py
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


def _calculatePRNodes(sourceNodes: typing.List[tgmodel.TGNode], edges: typing.Union[str, typing.Set[tgmodel.TGEdge]],
                      weight: typing.Union[str, float], defaultWeight: float,
                      matrix: typing.Dict[tgmodel.TGNode, typing.Dict[tgmodel.TGNode, float]])\
        -> typing.Set[tgmodel.TGNode]:
    """Returns the nodes which have no outbound edges and should therefore be stricken out of the matrix, for now."""
    strikeOut: typing.Set[tgmodel.TGNode] = set()
    sourceSet = set(sourceNodes)
    num = 0
    for node in sourceNodes:
        generics._log("Computing node {0} to prevent it's out vertices from being striken out...\n".format(str(num)))
        num += 1
        curOut: typing.Dict[tgmodel.TGNode, float] = {}
        curTotal = 0.0
        for edge in generics.getEdges(node, edges, weight, dedup=False, defaultDist=defaultWeight, keepMax=True):
            weight = generics.getEdgeDist(edge, weight, defaultWeight)
            if weight <= 0.0:
                continue
            otherVertex = generics.otherVertex(edge, node)
            if otherVertex not in sourceSet:
                continue
            if otherVertex not in curOut:
                curOut[otherVertex] = weight
            else:
                curOut[otherVertex] += weight
            curTotal += weight
        if curTotal == 0.0:
            strikeOut.add(node)
        else:
            for otherVertex in curOut:
                curOut[otherVertex] = curOut[otherVertex] / curTotal
        matrix[node] = curOut

    return strikeOut


def _multiplyMatrix(matrix: typing.Dict[tgmodel.TGNode, typing.Dict[tgmodel.TGNode, float]],
                    input: typing.Dict[tgmodel.TGNode, float], damping: float) -> typing.Dict[tgmodel.TGNode, float]:
    ret: typing.Dict[tgmodel.TGNode, float] = dict(((node, 0.0) for node in input))
    for node in input:
        curWeight = input[node]
        for otherNode in matrix[node]:
            if otherNode in ret:
                ret[otherNode] += curWeight * matrix[node][otherNode]

    if len(input) <= 0:
        adjustment = 1
    else:
        adjustment = (1.0 - damping)/len(input)
    for node in ret:
        ret[node] = adjustment + ret[node] * damping
    return ret


def _logCurrentRanks(vector: typing.Dict[tgmodel.TGNode, float]):
    for node in vector:
        generics._log("Node Name: {0}, Page Rank: {1}\n".format(node.getAttribute("airportID").getValue(),
                                                                str(vector[node])))


@tgsp.tgstoredproc
def pageRank(g: tgsp.TGGraphContext,
             source: str = "",
             weight: str = "1.0",
             defaultWeight: float = 0.0,
             dampingFactor: float = 0.85,
             epsilon: float = 0.00001,
             maxIterations: int = 20,
             edges: str = "OUT"
             ) -> "[(V,d)]":
    """
    Computes the PageRank for the source nodes and edges with the damping factor, weight, epsilon, and max iterations
    given.

    :param g: The graph context that all stored procedures get from the server.
    :param source: Optionally the source nodes for the subgraph that is used to determine the PageRank. Will acquire
    the source nodes from the current graph context if sources is not set.
    :type source: typing.Optional[str] A string representing a gremlin query.
    :param weight: What to use for the weight. Only makes sense to set to an attribute name to determine each edge's
    outgoing weight.
    :type weight: typing.Union[str, float] A string representing an attribute name or floating point.
    :param defaultWeight: What to use when the edge does not have the weight attribute specified. Only used when weight
    is set to an attribute name.
    :type defaultWeight: A floating point.
    :param dampingFactor: The damping factor that determines how likely it is for a given user to select the next node
    based on the outgoing relations for the current node versus a random jump to another node.
    :type dampingFactor: A positive floating point less than 1.0
    :param epsilon: If the maximum difference between the new PageRank scores and the previous iterations PageRank score
    drops below epsilon, then the PageRank algorithm will halt, even if this occurs before the maxIterations.
    :type epsilon: A positive floating point (likely close to 0.0)
    :param maxIterations: The maximum number of iterations before stopping the PageRank algorithm regardless of the
    maximum difference between the new PageRank scores and the previous iterations PageRank score.
    :type maxIterations: An Intger
    :param edges: Optionally the source edges for the subgraph that is used to determine the PageRank.
    :type edges: typing.Union[tgsp.TGGraphContext, str, typing.Set[tgmodel.TGEdge]] A string representing either a
    gremlin query or a direction from each node.
    :return:
    """
    # TODO Remove all logging references and indicators...
    source: typing.Optional[str]
    weight: typing.Union[str, float]
    edges: typing.Union[tgsp.TGGraphContext, str, typing.Set[tgmodel.TGEdge]]
    # source = None
    # weight = 1.0
    # defaultWeight = 0.0
    # dampingFactor = 0.85
    # epsilon = 0.00001
    # maxIterations = 20
    # edges = "OUT"
    generics._log("In pageRank...\n", mode='w')
    sourceNodes: typing.List[tgmodel.TGNode]
    if source is None or source == '':
        sourceNodes = list(g.getNodes())
    else:
        sourceNodes = list(g.executeQuery(source).toCollection())

    if isinstance(edges, tgsp.TGGraphContext):
        edges = set(edges.getEdges())
    elif edges not in generics.EDGE_TYPES:
        edges = set(g.executeQuery(edges).toCollection())

    generics._log("About to compute PageRank nodes (number of nodes: {0})...\n".format(len(sourceNodes)))
    finalMatrix: typing.Dict[tgmodel.TGNode, typing.Dict[tgmodel.TGNode, float]] = {}
    strikeOut = _calculatePRNodes(sourceNodes, edges, weight, defaultWeight, finalMatrix)
    useNodes: typing.List[tgmodel.TGNode] = list(sourceNodes)
    useMatrix: typing.Dict[tgmodel.TGNode, typing.Dict[tgmodel.TGNode, float]] = dict(finalMatrix)
    striken: typing.List[typing.Set[tgmodel.TGNode]] = []

    while len(strikeOut) > 0:
        useNodes = [node for node in useNodes if node not in strikeOut]
        striken.append(strikeOut)
        useMatrix = {}
        strikeOut = _calculatePRNodes(useNodes, edges, weight, defaultWeight, useMatrix)

    generics._log("Finsihed Computing PageRanke nodes\nAbout to start matrix multiplication...\n")
    if len(useNodes) <= 0:
        initialStart = 0
    else:
       initialStart = 1.0 / len(useNodes)
    inputVector = dict(((node, initialStart) for node in useNodes))

    for number in range(maxIterations):
        outVec = _multiplyMatrix(useMatrix, inputVector, dampingFactor)

        maxDist = 0.0
        for node in outVec:
            maxDist = max(maxDist, abs(outVec[node] - inputVector[node]))
        if maxDist < epsilon:
            break
        _logCurrentRanks(outVec)
        generics._log("Finished revision {0}\n".format(number))
        inputVector = outVec

    for strikedOut in reversed(striken):
        for node in strikedOut:
            inputVector[node] = 0.0
        inputVector = _multiplyMatrix(finalMatrix, inputVector, dampingFactor)

    for node in sourceNodes:
        if node not in inputVector:
            inputVector[node] = 0.0

    ret: typing.List[typing.Tuple[tgmodel.TGNode, float]] = []
    for node in inputVector:
        ret.append((node, inputVector[node]))
    _logCurrentRanks(inputVector)
    return ret
