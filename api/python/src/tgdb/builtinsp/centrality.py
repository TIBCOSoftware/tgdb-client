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
 *   File name   :centrality.py
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
import tgdb.builtinsp.shortest_path as shortest_path


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


@tgsp.tgstoredproc
def closenessCentrality(g: tgsp.TGGraphContext,
                        source: str = '',
                        target: str = '',
                        distance: str = "1.0",
                        defaultDistance: float = -1.0,
                        edges: str = 'OUT'
                        ) -> "[(V,d)]":
    """
    Determines the closeness centrality of the source nodes. If sources is not specified, takes them from the current
    graph context. If targets is specified, then the only vertices and their closeness returned are those in the targets
    query. If distance is specified and not a floating point value, then the value for that attribute on the edge is
    used.

    :param g: The graph context that the server passes into this stored procedure automatically.
    :param source: Specifies a query to acquire the nodes that determines the vertices for the graph.
    :type source: typing.Optional[str] A string representing a query.
    :param target: Specifies which vertices' closeness centrality score are desired. When not set, returns every
    vertices' closeness centrality score.
    :type target: typing.Optional[str] A string representing a query.
    :param distance: Specifies the distance parameter to use for closeness centrality.
    :type distance: A string representing an attribute descriptor name or a string of a float for the constant distance
    parameter.
    :param defaultDistance: What to use when the edge does not have the weight attribute specified. Only used when
    weight is set to an attribute name.
    :type defaultDistance: A floating point.
    :param edges: Specifies which edges to use to make the graph used to find the closeness centrality.
    :type edges: typing.Union[str, typing.Set[tgmodel.TGEdge]] A string representing a query.
    :return:
    """
    # TODO Remove all logging references and indicators...
    try:
        generics._log("Starting centrality...\n", mode='w')
        source: typing.Optional[str]
        target: typing.Optional[str]
        defaultDistance: float
        edges: typing.Union[str, typing.Set[tgmodel.TGEdge]]
        if source == '':
            source = None
        if target == '':
            target = None
        # source = None
        # target = None
        # distance = 1.0
        # defaultDistance = float('inf')
        # edges = 'OUT'
        if defaultDistance < 0:
            defaultDistance = float('inf')
        generics._log("distance attribute: %s, defaultDistance: %f\n" % (distance, defaultDistance))
        sourceNodes: typing.List[tgmodel.TGNode] = list(g.getNodes()) if source is None else\
                                                   list(g.executeQuery(source).toCollection())
        vertexPaths: typing.List[typing.Tuple[typing.List[tgmodel.TGNode]]] =\
            shortest_path.shortestPath(g, source=source, distance=distance, defaultDistance=defaultDistance,
                                       edges=edges, setResult=False)
        sourceSet: typing.Set[tgmodel.TGNode] = set(sourceNodes)
        targetSet: typing.Set[tgmodel.TGNode] = set(sourceNodes) if target is None else\
                                                set(g.executeQuery(target).toCollection())
        numReachable: typing.Dict[tgmodel.TGNode, int] = dict(((node, 0) for node in sourceNodes))
        distMap: typing.Dict[tgmodel.TGNode, float] = dict(((node, float('inf')) for node in sourceNodes))
        reached: typing.Dict[tgmodel.TGNode, typing.Set[tgmodel.TGNode]] =\
            dict(((node, {node}) for node in sourceNodes))
        if edges not in generics.EDGE_TYPES:
            edges = set(g.executeQuery(edges).toCollection())

        edgeMap: typing.Dict[tgmodel.TGNode, typing.Dict[tgmodel.TGNode, float]] = {}
        if not isinstance(distance, float):
            for node in sourceNodes:
                eMap = {}
                for edge in generics.getEdges(node, edges, distance, defaultDistance):
                    other = generics.otherVertex(edge, node)
                    eMap[other] = generics.getEdgeDist(edge, distance, defaultDistance)
                edgeMap[node] = eMap

        for pathTuple in vertexPaths:
            path = pathTuple[0]
            if len(path) <= 1:
                continue
            start = path[0]
            end = path[-1]
            # generics._log("%s pathDist: %s\n" % (start.getAttribute('airportID').getValue(),
            #                                    end.getAttribute('airportID').getValue()))
            if end in reached[start]:
                continue
            reached[start].add(end)
            pathDistance: float
            if isinstance(distance, float):
                pathDistance = (len(path) - 1) * distance
            else:
                pathDistance = 0.0
                for i in range(1, len(path)):
                    prev = path[i-1]
                    cur = path[i]
                    dist = defaultDistance
                    if cur in edgeMap[prev]:
                        dist = edgeMap[prev][cur]
                    pathDistance += dist
                # generics._log("%s pathDist: %d\n" % (start.getAttribute('airportID').getValue(), pathDistance))
            if distMap[start] < float('inf'):
                distMap[start] += pathDistance
            else:
                distMap[start] = pathDistance
            numReachable[start] += 1

        ret: typing.List[typing.Tuple[tgmodel.TGNode, float]] = []
        totalNodes: int = len(sourceNodes) - 1
        generics._log("Size of totalNodes {0}.\n".format(str(totalNodes)))
        #_logCurrentCent(distMap)
        #_logCurrentCent(numReachable)
        for node in distMap:
            reachable = numReachable[node] - 1
            sumDist = distMap[node]
            wassermanFaustCloseness = 0.0 if reachable == 0 else reachable * reachable / sumDist / totalNodes
            if node in targetSet:
                ret.append((node, wassermanFaustCloseness))

        _logCurrentCent(ret)
        return ret
    except BaseException as e:
        import traceback as tb
        with open('tb_file_print.txt', 'a') as tb_file:
            tb.print_tb(e.__traceback__, file=tb_file)
        raise e



@tgsp.tgstoredproc
def betweennessCentrality(g: tgsp.TGGraphContext,
                          source: str = "",
                          target: str = "",
                          distance: str = "1.0",
                          defaultDistance: float = -1.0,
                          edges: str = 'OUT'
                          ) -> "[(V,d)]":
    """
    Determines the betweenness centrality of the source nodes. If sources is not specified, takes them from the current
    graph context. If targets is specified, then the only vertices and their closeness returned are those in the targets
    query. If distance is specified and not a floating point value,  then the value for that attribute on the edge is
    used.

    :param g: The graph context that the server passes into this stored procedure automatically.
    :param source: Specifies a query to acquire the nodes that determines the vertices for the graph.
    :type source: A string representing a query.
    :param target: Specifies which vertices' betweenness centrality score are desired. When not set, returns every
    vertices' betweenness centrality score.
    :type target: A string representing a query.
    :param distance: Specifies the distance parameter to use for closeness centrality.
    :type distance: A string representing an attribute descriptor name or a string of a float for the constant distance
    parameter.
    :param defaultDistance: What to use when the edge does not have the weight attribute specified. Only used when
    weight is set to an attribute name.
    :type defaultDistance: A floating point.
    :param edges: Specifies which edges to use to make the graph used to find the closeness centrality.
    :type edges: A string representing a query.
    :return:
    """
    # TODO Remove all logging references and indicators...
    source: typing.Optional[str]
    target: typing.Optional[str]
    edges: typing.Union[str, typing.Set[tgmodel.TGEdge]]
    if source == '':
        source = None
    if target == '':
        target = None
    if defaultDistance < 0:
        defaultDistance = float('inf')
    generics._log("", mode='w')
    sourceNodes: typing.List[tgmodel.TGNode] = list(g.getNodes()) if source is None else\
                                               list(g.executeQuery(source).toCollection())
    vertexPaths =\
        shortest_path.shortestPath(g, source=source, distance=distance, defaultDistance=defaultDistance, edges=edges,
                                   setResult=False)
    targetSet: typing.Set[tgmodel.TGNode] = set(sourceNodes) if target is None else\
                                            set(g.executeQuery(target).toCollection())
    numBetween: typing.Dict[tgmodel.TGNode, int] = dict(((node, 0) for node in sourceNodes))

    for pathTuple in vertexPaths:
        path = pathTuple[0]
        if len(path) < 2:
            continue
        for intermediate in path[1:-1]:
            numBetween[intermediate] += 1

    ret: typing.List[typing.Tuple[tgmodel.TGNode, float]] = []
    for node in numBetween:
        if node in targetSet:
            ret.append((node, numBetween[node] / len(vertexPaths)))

    _logCurrentCent(ret)
    return ret
