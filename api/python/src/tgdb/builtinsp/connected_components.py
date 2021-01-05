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
 *   File name   :connected_components.py
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


@tgsp.tgstoredproc
def connectedComponents(g: tgsp.TGGraphContext,
                        source: str = "",
                        edges: str = "BOTH"
                        ) -> "[(V,l)]":
    """
    Determines connected components in the subgraph described by source and edges.

    :param g: The graph context that all stored procedures get from the server.
    :param source: Optionally, the set of nodes to determine the connected components on.
    :type source: typing.Union[str, None] A string representing a gremlin query.
    :param edges: A set of edges to determine the connected components between the given vertices.
    :type edges: typing.Union[str, typing.Set[tgmodel.TGEdge], tgsp.TGGraphContext] A string representing a gremlin
    query.
    :return:
    """
    # TODO Remove all logging references and indicators...
    source: typing.Optional[str]
    edges: typing.Union[str, typing.Set[tgmodel.TGEdge], tgsp.TGGraphContext]
    if source == "":
        source = None
    sourceNodes = list(g.getNodes()) if source is None else list(g.executeQuery(source).toCollection())

    nodeSet = set(sourceNodes)

    generics._log("", mode='w')
    for node in sourceNodes:
        generics._log("Node: Id {0} has {1} edges.\n".format(str(node.getAttribute('airportID').getValue()),
                                                             len(list(node.getEdges()))))
    alreadyFound: typing.Set[tgmodel.TGNode] = set()
    ret: typing.List[typing.Tuple[tgmodel.TGNode, int]] = []

    if isinstance(edges, tgsp.TGGraphContext):
        edges = set(edges.getEdges())
    elif edges not in generics.EDGE_TYPES:
        edges = set(g.executeQuery(edges).toCollection())

    comp = 0
    nodeIdx = 0
    if not isinstance(comp, int):
        generics._log("Unknown type for comp: {0}\n".format(str(type(comp))))
    else:
        generics._log("type(comp) is int\n")
    while nodeIdx < len(sourceNodes):
        curNode = sourceNodes[nodeIdx]
        nodeIdx += 1
        if curNode not in alreadyFound:
            comp += 1
            toFind = queue.SimpleQueue()
            toFind.put(curNode)
            while not toFind.empty():
                node = toFind.get()
                if node in alreadyFound:
                    continue
                alreadyFound.add(node)
                ret.append((node, comp))
                for edge in generics.getEdges(node, edges):
                    otherVertex = generics.otherVertex(edge, node)
                    if otherVertex not in alreadyFound and otherVertex in nodeSet:
                        toFind.put(otherVertex)

    generics._log("About to return value...\n")
    return ret
