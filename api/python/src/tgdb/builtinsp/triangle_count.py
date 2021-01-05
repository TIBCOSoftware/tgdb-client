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
 *   File name   :triangle_count.py
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


def _logCurrentCent(cent_ret: typing.Union[typing.List[typing.Tuple[tgmodel.TGNode, float]],
                                           typing.Dict[tgmodel.TGNode, float]]):
    if isinstance(cent_ret, dict):
        for node in cent_ret:
            cent = cent_ret[node]
            generics._log("Node Name: {0}, Centrality: {1}\n".format(node.getAttribute("airportID").getValue(),
                                                                     str(cent)))
    else:
        for node, cent in cent_ret:
            generics._log("Node Name: {0}, Centrality: {1}\n".format(node.getAttribute("airportID").getValue(),
                                                                     str(cent)))


def _getDirectlyReachable(node: tgmodel.TGNode,
                          edges: typing.Union[str, typing.Set[tgmodel.TGEdge], tgsp.TGGraphContext],
                          nodes: typing.Set[tgmodel.TGNode]
                          ) -> typing.Set[tgmodel.TGNode]:
    ret = set()
    for edge in generics.getEdges(node, edges):
        otherNode = generics.otherVertex(edge, node)
        if otherNode != node and otherNode in nodes:
            ret.add(otherNode)
    return ret


@tgsp.tgstoredproc
def triangleCounts(g: tgsp.TGGraphContext,
                   source: str = "",
                   edges: str = "BOTH"
                   ) -> "[(V,l)]":
    """
    Computes the triangle count (the number of 3 distinct node paths that run through each node) on a per node basis.

    :param g: The graph context that all stored procedures get from the server.
    :param source: Optionally, the set of nodes to determine the triangle counts.
    :type source: typing.Optional[str] A string representing a gremlin query.
    :param edges: A set of edges to determine the triangle counts between the given vertices.
    :type edges: typing.Union[str, typing.Set[tgmodel.TGEdge], tgsp.TGGraphContext] A string representing a gremlin
    query.
    :return:
    """
    # TODO Remove all logging references and indicators...
    source: typing.Optional[str]
    edges: typing.Union[str, typing.Set[tgmodel.TGEdge], tgsp.TGGraphContext]
    generics._log("Before getNodes('airportType')\n", mode='w')
    nodes = list(g.getNodes()) if source is None or source == "" else list(g.executeQuery(source).toCollection())
    generics._log("After getNodes('airportType')\n")
    nodeSet = set(nodes)
    if isinstance(edges, tgsp.TGGraphContext):
        edges = set(edges.getEdges())
    elif edges not in generics.EDGE_TYPES:
        edges = set(g.executeQuery(edges).toCollection())

    nodeLookup: typing.Dict[tgmodel.TGNode, typing.Set[tgmodel.TGNode]] =\
        dict(((node, _getDirectlyReachable(node, edges, nodeSet)) for node in nodes))
    ret: typing.List[typing.Tuple[tgmodel.TGNode, int]] = []

    for node in nodes:
        count = 0
        alreadyFound = {node}
        for node_1 in nodeLookup[node]:
            alreadyFound.add(node_1)
            for node_2 in nodeLookup[node_1]:
                if node_2 in nodeLookup[node] and node_2 not in alreadyFound:
                    count += 1
        ret.append((node, count))
    _logCurrentCent(ret)
    return ret
