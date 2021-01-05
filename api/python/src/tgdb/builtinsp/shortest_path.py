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
 *   File name   :shortest_path.py
 *   Created on  :8/21/20.
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


def _singleSourceShortestPath(source: tgmodel.TGNode, targets: typing.Optional[typing.Set[tgmodel.TGNode]],
                              edges: typing.Union[str, typing.Set[tgmodel.TGEdge]], dist: typing.Union[str, float],
                              defaultDist: float, maxDist: float, includeEdges: bool)\
        -> typing.Dict[tgmodel.TGNode, typing.List[tgmodel.TGEntity]]:
    pque = queue.PriorityQueue()
    pque.put((0.0, source, None))
    visited: typing.Dict[tgmodel.TGNode, typing.List[tgmodel.TGEntity]] = {}
    ret: typing.Dict[tgmodel.TGNode, typing.List[tgmodel.TGEntity]] = {}

    while pque.empty() is False:
        cur_dist, cur_node, cur_path = pque.get()

        if targets is not None and targets.issubset(visited):
            return ret

        if cur_node in visited:
            continue

        comp_path: typing.List = []
        if cur_path is None:
            comp_path.append(cur_node)
        else:
            comp_path.append(visited[cur_path[0]])
            comp_path.append(cur_path[1])

        visited[cur_node] = comp_path

        if targets is None or cur_node in targets:
            # This code "inflates" the visited's reference based map
            tmp_node = comp_path[0]
            ret_path = [tmp_node]
            while len(visited[tmp_node]) > 1:
                if includeEdges:
                    ret_path.append(visited[tmp_node][1])
                tmp_node = visited[tmp_node][0]
                ret_path.append(tmp_node)
            if cur_dist < maxDist:
                ret[cur_node] = list(reversed(ret_path))

        for edge in generics.getEdges(cur_node, edges, dist, defaultDist):
            next_dist = generics.getEdgeDist(edge, dist, defaultDist)
            next_node = generics.otherVertex(edge, cur_node)
            pque.put((next_dist, next_node, (cur_node, edge)))

    return ret


def _allPairsShortestPaths(nodes: typing.Iterable[tgmodel.TGNode], edges: typing.Union[typing.Set[tgmodel.TGEdge], str],
                           distance: typing.Union[str, float], defaultDistance: float, maxDistance: float,
                           includeEdges: bool)\
        -> typing.Dict[typing.Tuple[tgmodel.TGNode, tgmodel.TGNode], typing.List[tgmodel.TGEntity]]:
    ret: typing.Dict[typing.Tuple[tgmodel.TGNode, tgmodel.TGNode], typing.List[tgmodel.TGEntity]] = {}
    matrix: typing.Dict[tgmodel.TGNode, typing.Dict[tgmodel.TGNode, float]] = {}
    nodes = list(nodes)
    nodeSet = set(nodes)

    for node in nodes:
        curEdgeWeight = {}
        matrix[node] = curEdgeWeight
        for edge in generics.getEdges(node, edges, distance, defaultDistance, True):
            otherNode = generics.otherVertex(edge, node)
            if otherNode not in nodeSet:
                continue
            curEdgeWeight[otherNode] = generics.getEdgeDist(edge, distance, defaultDistance)
            ret[(node, otherNode)] = [node, edge, otherNode] if includeEdges else [node, otherNode]
        curEdgeWeight[node] = 0.0
        ret[(node, node)] = list((node,))

    for intermediateNode in nodes:
        for startNode in nodes:
            for endNode in nodes:
                potDistance = float('inf')
                curDistance = float('inf')
                if endNode in matrix[startNode]:
                    curDistance = matrix[startNode][endNode]
                if intermediateNode in matrix[startNode] and endNode in matrix[intermediateNode]:
                    potDistance = matrix[startNode][intermediateNode] + matrix[intermediateNode][endNode]
                if potDistance < curDistance:
                    path = list(ret[(startNode, intermediateNode)])
                    path.extend(ret[(intermediateNode, endNode)][1:])
                    # if startNode.getAttribute("airportID").getValue() == 'AIRPORT340':
                    #     generics._log("AIRPORT340 (temporary) path to {0} is {1}.\n".format(
                    #         endNode.getAttribute("airportID").getValue(),
                    #         str([node.getAttribute("airportID").getValue() for node in path])
                    #     ))
                    ret[(startNode, endNode)] = path
                    matrix[startNode][endNode] = potDistance

    toDelete: typing.List[typing.Tuple[tgmodel.TGNode, tgmodel.TGNode]] = []
    for start, end in ret:
        if matrix[start][end] > maxDistance:
            toDelete.append((start, end))

    for toDel in toDelete:
        del ret[toDel]

    return ret


@tgsp.tgstoredproc
def shortestPath(g: tgsp.TGGraphContext,
                 source: str = "",
                 target: str = "",
                 distance: str = "1.0",
                 defaultDistance: float = -1.0,
                 edges: str = "OUT",
                 maxDistance: float = -1.0,
                 includeEdges: bool = False,
                 setResult: bool = True
                 ) -> "[[V]]":        # TODO Change signature to List of Paths (i.e. "[P([T])]"
    """
    Determines the shortest path between source and target nodes, using the distance metric, the edges specified, the
    furthest distance between two nodes for the path to be included, and whether to include edges in the path.

    If only one source node (or one target node) is specified, does a single-source shortest path.

    If multiple source nodes are set and multiple targets, then does a single source shortest path to all the targets
    from each source node.

    If multiple source nodes are set and targets are not set, then does an all-pairs shortest path between the sources.

    :param g: The graph context. Should be passed by the stored procedure manager.
    :param source: The sources for the shortest path algorithm. The start for all returned paths will be in the nodes in
    `source`.
    :type source: typing.Optional[str]  Optional, a graph traversal that returns a list of nodes or a string
    representation of a query.
    :param target: The targets for the shortest path algorithm. When set, the destination for all returned paths will be
    in the nodes in `target`.
    :type target: typing.Union[None, tgsp.TGGraphContext, str] Optional, a graph traversal that returns a list of nodes.
    :param distance: The distance parameter to use for determining the shortest path. If set to a string, will use that
    attribute to determine what the distance is. If that attribute is not set, will assume infinity (i.e. to not use
    it).
    :type distance: typing.Union[str, float] Attribute name, or a floating point.
    :param defaultDistance: What to use when the edge does not have the weight attribute specified. Only used when
    weight is set to an attribute name.
    :type defaultDistance: A floating point.
    :param edges: A graph traversal that indicates the edges to use or the direction for the those edge's direction.
    :type edges: typing.Union[tgsp.TGGraphContext, str, typing.Set[tgmodel.TGEdge]] A graph traversal that returns a
    list of edges or a string that is in the set {"BOTH", "IN", "OUT"}
    :param maxDistance: The maximum distance for a returned path.
    :type maxDistance: A floating point number.
    :param includeEdges: Whether to include the edges in the output path.
    :type includeEdges: Boolean, defaults to False.
    :return: A list of paths (which each are a list of entities)
    """
    # TODO Remove all logging references and indicators...
    try:
        source: typing.Union[str, None, typing.Set[tgmodel.TGNode]]
        target: typing.Union[None, tgsp.TGGraphContext, str]
        distance: typing.Union[str, float]
        defaultDistance: float
        edges: typing.Union[tgsp.TGGraphContext, str, typing.Set[tgmodel.TGEdge]]
        maxDistance: float
        includeEdges: bool
        setResult: bool
        if target == '':
            target = None
        if source == '':
            source = None
        if maxDistance < 0:
            maxDistance = float('inf')
        if defaultDistance < 0:
            defaultDistance = float('inf')
        sourceNodes: typing.List[tgmodel.TGNode]
        if source is None or source == '':
            sourceNodes = list(g.getNodes())
        elif isinstance(source, set):
            sourceNodes = list(source)
        else:
            sourceNodes = list(g.executeQuery(source).toCollection())
        targetNodes = sourceNodes
        if isinstance(target, str) and target != '':
            targetNodes = list(g.executeQuery(target).toCollection())
        elif isinstance(target, set):
            targetNodes = list(target)
        elif isinstance(target, tgsp.TGGraphContext):
            targetNodes = list(target.getNodes())

        if isinstance(edges, tgsp.TGGraphContext):
            edges = set(edges.getEdges())
        elif edges not in generics.EDGE_TYPES:
            edges = set(g.executeQuery(edges).toCollection())

        ret: typing.List[typing.Tuple[typing.List[tgmodel.TGEntity]]] = []
        if len(sourceNodes) == 1 and target is None:
            retVal = _singleSourceShortestPath(sourceNodes[0], target, edges, distance, defaultDistance, maxDistance,
                                               includeEdges)
            ret = [(retVal[node],) for node in retVal]
        elif len(sourceNodes) == 1:
            retVal = _singleSourceShortestPath(sourceNodes[0], set(targetNodes), edges, distance, defaultDistance,
                                               maxDistance, includeEdges)
            ret = [(retVal[node],) for node in retVal]
        elif target is None:
            retVal = _allPairsShortestPaths(sourceNodes, edges, distance, defaultDistance, maxDistance, includeEdges)
            ret = [(retVal[nodePair],) for nodePair in retVal]
        elif len(targetNodes) == 1:
            if edges == "OUT":
                edges = "IN"
            elif edges == "IN":
                edges = "OUT"
            retVal = _singleSourceShortestPath(targetNodes[0], set(sourceNodes), edges, distance, defaultDistance,
                                               maxDistance, includeEdges)
            ret = [(retVal[node],) for node in retVal]
        else:
            pass    # Not sure what to do here. Shortest Path from every source node to every target node, using
                    # intermediate edges (i.e. edges with neither vertex in either sourceNodes or targetNodes).

        return ret
    except BaseException as e:
        import traceback as tb
        with open('tb_file_print.txt', 'w') as tb_file:
            tb.print_tb(e.__traceback__, file=tb_file)
        raise e


# TODO Max Flow, Max Cut, Min Cut, Label Propagation, Louvain Modularity
