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
 *   File name   :generic.py
 *   Created on  :8/31/20.
 *   Created by  :dhudson
 *   SVN Id      :$Id:$
"""

import tgdb.storedproc as tgsp
import tgdb.model as tgmodel
import tgdb.exception as tgexcept
import typing
import re
import queue
import ctypes as ct


EDGE_TYPES = {'OUT', 'IN', 'BOTH'}


def _addOrKeep(nodeEdges: typing.Dict[tgmodel.TGNode, tgmodel.TGEdge], edge: tgmodel.TGEdge, node: tgmodel.TGNode,
               keepMax: bool, distance: typing.Union[str, float], defaultDist: float):
    if node in nodeEdges:
        other = nodeEdges[node]
        if isinstance(distance, str):       # if distance is a float, then it is a constant factor and can be ignored.
            otherValue = getEdgeDist(other, distance, defaultDist)
            thisValue = getEdgeDist(edge, distance, defaultDist)
            if keepMax and thisValue > otherValue:
                nodeEdges[node] = edge
            elif not keepMax and thisValue < otherValue:
                nodeEdges[node] = edge
    else:
        nodeEdges[node] = edge


def getEdges(node: tgmodel.TGNode, edges: typing.Union[typing.Set[tgmodel.TGEdge], str,
                                                       typing.Dict[tgmodel.TGNode,typing.Iterable[tgmodel.TGEdge]]],
             distance: typing.Union[str, float] = 1.0, defaultDist: float = float('inf'), dedup: bool = True,
             keepMax: bool = False) -> typing.List[tgmodel.TGEdge]:
    """
    Gets all of the edges for this node, possibly running deduplication on the potential return values.

    :param node: The node whose edges we want to acquire.
    :param edges: The edge acquire method. If it is a set of edges, then includes all edges that are incident to this
    node that are in the edges set. Otherwise, it gets the edges that are in the direction specified by this parameter.
    :param distance: The distance to use when deduplicating and to determine which edge to keep.
    :param defaultDist: The default distance to use when deduplicating and an edge does not have the distance attribute.
    :param dedup: Whether to deduplicate the edges (in this case deduplication removes edges incident on the same pair
    of nodes).
    :param keepMax: When deduplicating, only keeps the edges with the largest distance (if True) or smallest distance
    (if False)
    :return:
    """
    ret = []
    if isinstance(edges, set):
        for edge in node.getEdges():
            if edge in edges:
                ret.append(edge)
    if isinstance(edges, dict):
        if node in edges:
            ret.extend(edges[node])
    elif edges == "BOTH":
        ret.extend(node.getEdges())
    elif edges == "IN":
        filter = lambda edge: edge.getDirection() != tgmodel.DirectionType.Directed.value or\
                                              edge.getVertices()[1] == node
        ret.extend([edge for edge in node.getEdges() if filter(edge)])
    elif edges == "OUT":
        filter = lambda edge: edge.getDirection() != tgmodel.DirectionType.Directed.value or \
                                              edge.getVertices()[0] == node
        ret.extend([edge for edge in node.getEdges() if filter(edge)])
    else:
        raise tgexcept.TGException("Unknown edge direction: {0}".format(edges))
    if dedup:
        fromNodes: typing.Dict[tgmodel.TGNode, tgmodel.TGEdge] = {}
        toNodes: typing.Dict[tgmodel.TGNode, tgmodel.TGEdge] = {}
        for edge in ret:
            if edge.getVertices()[0] == node:
                _addOrKeep(toNodes, edge, edge.getVertices()[1], keepMax, distance, defaultDist)
            else:
                _addOrKeep(fromNodes, edge, edge.getVertices()[0], keepMax, distance, defaultDist)
        ret = []
        for node in fromNodes:
            ret.append(fromNodes[node])
        for node in toNodes:
            ret.append(toNodes[node])
    return ret


def otherVertex(edge: tgmodel.TGEdge, node: tgmodel.TGNode) -> tgmodel.TGNode:
    """
    Gets the edge's other vertex.

    :param edge: The edge.
    :param node: 'This' vertex.
    :return: The other vertex.
    """
    vertices = edge.getVertices()
    return vertices[0] if node == vertices[1] else vertices[1]


def getEdgeDist(edge: tgmodel.TGEdge, distance: typing.Union[str, float], default: float = float('inf')) -> float:
    """
    Gets the given edge's distance/cost/weight based on the distance parameter, with the default parameter as the result
    when the distance/cost/weight does not exist on the given edge.

    :param edge: The edge to determine the value
    :param distance: The edge's cost, weight, or distance value. If it is a string, we interpret this as an attribute
    name.
    :param default: The default value if the distance parameter is a string and the edge does not have that attribute
    set.
    :return:
    """
    if isinstance(distance, float):
        return distance
    else:
        try:
            return float(distance)
        except:
            attr = edge.getAttribute(distance)
            if attr is None:
                return default
            else:
                return float(attr.getValue())


def _log(log: str, filename: str = 'tb_file_print.txt', mode: str = 'a'):
    """
    Writes the to the log file the string passed in.

    :param log: The log message to write to the file.
    :param filename: The name of the file to write to.
    :param mode: The file mode to write.
    """
    with open(filename, mode) as tb_file:
        tb_file.write(log)

