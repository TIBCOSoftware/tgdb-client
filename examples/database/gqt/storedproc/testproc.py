"""
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 *  File name :testproc.py
 *  Created on: 6/14/2019
 *  Created by: katie 
 *
 *  This is an example stored procedure file
 *
 *  Usage:
 *  1.  In the database configuration file:
 *          * define a nodetype named "basicnode"
 *          * point to the directory containing this file in the [storedproc] section
 *  2.  Initialize and start the server as normal
 *  3.  In the admin console, execute the stored procedures defined below via gremlin query. e.g.:
 *          > g.V().execsp('retD');
 *          > g.V().execsp('retN');
 *
 """

from tgdb.storedproc import *

"""
    A stored procedure which operates on the routes data set, returning a list
    of airports based on the provided city and country
"""
@tgstoredproc(alias="getAirports")
def getAirports(g: TGGraphContext, country : str, city : str) -> "[V]":
    result = []
    nodes = g.getNodes("airportType")
    for node in nodes:
        co = node.getAttribute("country")
        ci = node.getAttribute("city")
        if co.getValue() == country and ci.getValue() == city:
            result.append(node)
    return result

# return single double value
@tgstoredproc(alias="retD")
def retD(g: TGGraphContext) -> "d":
    return 8.24

# return single Node value
@tgstoredproc(alias="retN")
def retN(g: TGGraphContext) -> "V":
    nodes = g.getNodes("airportType")
    for node in nodes:
        return node

# return single tuple of double
@tgstoredproc(alias="retT")
def retT(g: TGGraphContext) -> "(d)":
    tup = (8.0,)
    return tup

@tgstoredproc(alias="retT2D")
def retT2D(g: TGGraphContext) -> "(d, d)":
    tup = (8.0, 9.0)
    return tup

@tgstoredproc(alias="retTL")
def retTL(g: TGGraphContext) -> "([d])":
    l = [8.0, 9.0]
    return (l,)

@tgstoredproc(alias="retNestedT")
def retNestedT(g:TGGraphContext) -> "((i), s)":
    return ((3,),"hello")

# return empty list of node
@tgstoredproc(alias="retEmptyL")
def retEmptyL(g: TGGraphContext) -> "[V]":
    return []

# return list of node
@tgstoredproc(alias="retLN")
def retLN(g: TGGraphContext) -> "[V]":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    l1 = [v1,v1,v1]
    return l1

# return list of list of Node
@tgstoredproc(alias="retLLN")
def retLLN(g: TGGraphContext) -> "[[V]]":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    l1 = [v1,v1]
    l2 = [v1]
    l3 = [v1,v1]
    result = [l1,l2,l3]
    return result

# return list of tuple of Node,long
@tgstoredproc(alias="retLTNl")
def retLTNl(g: TGGraphContext) -> "[(V,l)]":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    t = (v1, 8)
    l1 = [t, t]
    result = l1
    return result

# return List of Path of fixed length
@tgstoredproc(alias="retLP")
def retLP(g: TGGraphContext) -> "[P(d, s, V) ]":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    s = "sample str"
    s1 = "another sample str"
    p = []
    p.append([8.3, s, v1])
    p.append([3.3, s1, v1])
    return p

# return List of Path with variable # of string 
@tgstoredproc(alias="retLPLs")
def retLPLs(g: TGGraphContext) -> "[P( [s] ) ]":
    s = "sample str"
    p = []
    p.append([s, s, s])
    return p

# return List of Path with variable length string, node value pair
@tgstoredproc(alias="retLPLsV")
def retLPLsV(g: TGGraphContext) -> "[P( [(s, V)] ) ]":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    s = "sample str"
    p = []
    p.append([s, v1, s, v1])
    return p

# return List of Path of double with variable length string, node value pair
@tgstoredproc(alias="retLPdL")
def retLPdL(g: TGGraphContext) -> "[P(d, [(s, V)] ) ]":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    s = "sample str"
    p = []
    p.append([8.3, s, v1, s, v1])
    return p

# return map with K:string, V: List of tuple of Node,long
@tgstoredproc(alias="retML")
def retML(g: TGGraphContext) -> "{s,[(V,l)]}":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    t = (v1, 8)
    l1 = [t, t]
    return {"test result": l1}

# return map with K:string, V: List of Path of generic type
@tgstoredproc(alias="retMLP")
def retMLP(g: TGGraphContext) -> "{s,[P]}":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    p = []
    p.append([8, "test123", v1])
    return {"test result": p}

# return map with K:string, V: List of Path of double with variable length string, node value pair
@tgstoredproc(alias="retMLPL")
def retMLPL(g: TGGraphContext) -> "{s, [P(d, [(s, V)] ) ] }":
    v1 = None
    nodes = g.getNodes("airportType")
    for node in nodes:
        v1 = node
        break
    s = "sample str"
    p = []
    p.append([8.3, s, v1, s, v1]) 
    return {"test result": p}
