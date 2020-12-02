/*
 * Copyright 2020 TIBCO Software Inc. All rights reserved.
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
 * File name: cytoscapeconverter.go
 * Created on: 06/27/2020
 * Created by: nimish
 *
 * SVN Id: $Id$
 */

package spotfire

import (
	"encoding/json"
	"strconv"
	"tgdb"
	"tgdb/impl"
)

var logger = impl.DefaultTGLogManager().GetLogger()

type CytoscapeDataGroupNode struct {
	Group string					`json:"group"`
	Data map[string]interface{}		`json:"data"`
}

type CytoscapeDataNode struct {
	Id string		`json:"id"`
	Name string		`json:"name"`
}

type CytoscapeDataEdge struct {
	Id string		`json:"id"`
	Source string	`json:"source"`
	Target string	`json:"target"`
}


func TestFormCytoscapeResult () []interface{} {
	groupNode0 := CytoscapeDataGroupNode{
		Group: "nodes",
	}
	groupNode0.Data = make(map[string]interface{})
	groupNode0.Data["id"] = "A0"
	groupNode0.Data["name"] = "A0 Name"
	groupNode0.Data["udt0"] = 20

	groupNode1 := CytoscapeDataGroupNode{
		Group: "nodes",
	}
	groupNode1.Data = make(map[string]interface{})

	groupNode1.Data["id"] = "B0"
	groupNode1.Data["name"] = "B0 Name"
	groupNode1.Data["udt1"] = 34.5


	resultsetArray := make([]interface{}, 2)
	resultsetArray[0] = groupNode0
	resultsetArray[1] = groupNode1


	result, er := json.MarshalIndent(resultsetArray, "", "\t")
	if er != nil {
		//fmt.Println("error:", er)
		logger.Error("Error: " + er.Error())
	}

	if logger.IsDebug() {
		logger.Debug("Query Result:" + string(result))
	}
	//fmt.Println("Query Result:" + string(result))

	return resultsetArray
}

func FormCytoscapeResult (resultSet tgdb.TGResultSet) []CytoscapeDataGroupNode {
	collection := resultSet.ToCollection()
	resultsetArray := make([]CytoscapeDataGroupNode, 0)
	for i := 0; i < len(collection); i++ {
		resultsetArray = dealWithEntityRecursively(collection[i], resultsetArray, true)
		countOfNodes, bNodesPresent := findNodeCount(collection[i])
		if bNodesPresent {
			updateCytoscapteDataGroupNode (countOfNodes, collection[i], resultsetArray)
		}
	}
	return resultsetArray
}

func updateCytoscapteDataGroupNode(countOfNodes int, tgdbCollection interface{}, cytoNodes []CytoscapeDataGroupNode) {
	arrayList, ok := tgdbCollection.([]interface{})
	if ok {
		if len(arrayList) > 0 {
			for j := 0; j < len(arrayList); j++ {
				tgdbNode, ok := arrayList[j].(*impl.Node)
				if ok {
					for i := 0; i < len(cytoNodes); i++ {
						cytoId, ok := cytoNodes[i].Data["id"]
						if ok {
							if tgdbNode.GetVirtualId() == cytoId {
								minDistance, ok := cytoNodes[i].Data["$minDistance"]
								if ok {
									if countOfNodes < minDistance.(int) {
										cytoNodes[i].Data["$minDistance"] = countOfNodes
										//countOfNodes--
									}
								} else {
									cytoNodes[i].Data["$minDistance"] = countOfNodes
									//countOfNodes--
								}

								maxDistance, ok := cytoNodes[i].Data["$maxDistance"]
								if ok {
									if countOfNodes > maxDistance.(int) {
										cytoNodes[i].Data["$maxDistance"] = countOfNodes
										//countOfNodes--
									}
								} else {
									cytoNodes[i].Data["$maxDistance"] = countOfNodes
									//countOfNodes--
								}
								countOfNodes--
								break;
							}
						}
					}
				}
			}
		}
	}
}


func findNodeCount(entity interface{}) (int, bool) {
	count := 0
	arrayList, ok := entity.([]interface{})
	if ok {
		for i := 0; i < len(arrayList); i++ {
			_, ok := arrayList[i].(*impl.Node)
			if ok {
				count++
			}
		}
	} else {
		return -1, false
	}
	return count, true
}

func dealWithEntityRecursively (entity interface{}, resultsetArray []CytoscapeDataGroupNode, isEdgeNodeToBePut bool) ([]CytoscapeDataGroupNode){
	node, ok := entity.(*impl.Node)
	if ok {
		groupNode0 := dealWithNode(resultsetArray, node, true)
		if groupNode0 != nil {
			resultsetArray = append(resultsetArray, *groupNode0)
			return resultsetArray
		}
		return resultsetArray
	} else {
		edge, ok := entity.(*impl.Edge)
		if ok {
			edge, fromNode, toNode := dealWithEdge(resultsetArray, edge)
			if edge != nil {
				resultsetArray = append(resultsetArray, *edge)
				if isEdgeNodeToBePut {
					if fromNode != nil {
						resultsetArray = append(resultsetArray, *fromNode)
					}
					if toNode != nil {
						resultsetArray = append(resultsetArray, *toNode)
					}
				}
				return resultsetArray
			}
			return resultsetArray
		} else {
			arrayList, ok := entity.([]interface{})
			if ok {
				for j := 0; j < len(arrayList); j++ {
					resultsetArray = dealWithEntityRecursively(arrayList[j], resultsetArray, false)
					//if tempResultSet != nil {
					//	resultsetArray = tempResultSet
					//}
				}
				return resultsetArray
			}
			return resultsetArray
		}
	}
}


func dealWithEdge(resultsetArray []CytoscapeDataGroupNode, edge *impl.Edge) (*CytoscapeDataGroupNode, *CytoscapeDataGroupNode, *CytoscapeDataGroupNode) /* edge, fromNode, toNode*/ {
	index, isPresent := IsEdgePresent(resultsetArray, *edge)
	if !isPresent {
		groupNode1 := CytoscapeDataGroupNode{
			Group: "edges",
		}

		groupNode1.Data = make(map[string]interface{})
		groupNode1.Data["id"] = edge.EntityId
		groupNode1.Data["source"] = edge.GetFromNode().GetVirtualId()
		groupNode1.Data["target"] = edge.GetToNode().GetVirtualId()

		entityName := edge.EntityType.GetName()
		if len(entityName) <= 0 {
			entityName = "TGDB_UNKNOWN"
		}
		groupNode1.Data["$type"] = entityName

		fillCytoscapeEdge (groupNode1, edge)
		//resultsetArray = append(resultsetArray, groupNode1)
		// return

		var groupNodeFrom *CytoscapeDataGroupNode
		fromNode, ok := edge.GetFromNode().(*impl.Node)
		if ok {
			groupNodeFrom = dealWithNode(resultsetArray, fromNode, false)
			//if groupNodeFrom != nil {
				//resultsetArray = append(resultsetArray, *groupNodeFrom)
			//}
		}

		var groupNodeTo *CytoscapeDataGroupNode
		toNode, ok := edge.GetToNode().(*impl.Node)
		if ok {
			groupNodeTo = dealWithNode(resultsetArray, toNode, false)
			//if groupNodeTo != nil {
				//resultsetArray = append(resultsetArray, *groupNodeTo)
			//}
		}
		return &groupNode1, groupNodeFrom, groupNodeTo
	} else {
		currentCount, isOK := resultsetArray[index].Data["$count"]
		if isOK {
			currentCount = currentCount.(int) + 1
			resultsetArray[index].Data["$count"] = currentCount
		} else {
			resultsetArray[index].Data["$count"] = 2
		}
		return nil, nil, nil
	}
}

func dealWithNode(collection []CytoscapeDataGroupNode, node *impl.Node, isPresentInResultDirectly bool) *CytoscapeDataGroupNode {
	index, isPresent := IsNodePresent(collection, *node)
	if !isPresent {
		groupNode0 := CytoscapeDataGroupNode{
			Group: "nodes",
		}

		groupNode0.Data = make(map[string]interface{})
		groupNode0.Data["id"] = node.EntityId//strconv.FormatInt(node.EntityId, 10)
		groupNode0.Data["name"] = strconv.FormatInt(node.EntityId, 10)
		if !isPresentInResultDirectly {
			groupNode0.Data["$internal"] = -1
		}
		entityName := node.EntityType.GetName()
		if len(entityName) <= 0 {
			entityName = "TGDB_UNKNOWN"
		}
		groupNode0.Data["$type"] = entityName

		fillCytoscapeNode (groupNode0, node)
		return &groupNode0

		//returnedArray := fillEdgesFromNode (node, resultsetArray)
		//resultsetArray = append(resultsetArray, returnedArray)
	} else {
		currentCount, isOK := collection[index].Data["$count"]
		if isOK {
			currentCount = currentCount.(int) + 1
			collection[index].Data["$count"] = currentCount
		} else {
			collection[index].Data["$count"] = 2
		}
	}

	return nil
}

func fillCytoscapeEdge(cytoNode CytoscapeDataGroupNode, tgdbEdge *impl.Edge) {
	for k, v := range tgdbEdge.Attributes {
		cytoNode.Data[k] = v.GetValue()
	}
}

func fillCytoscapeNode(cytoNode CytoscapeDataGroupNode, tgdbNode *impl.Node) {
	for k, v := range tgdbNode.Attributes {
		cytoNode.Data[k] = v.GetValue()
	}
}

//func fillEdgesFromNode(node *impl.Node, collection []interface{}) []interface{} {
//	resultsetArray := make([]interface{}, 0)
//	edges := node.GetEdges()
//	for i := 0; i < len(edges); i++ {
//		if !IsEdgePresent_(collection, edges[i]) {
//			groupNode1 := CytoscapeDataGroupNode{
//				Group: "edges",
//				Data: CytoscapeDataEdge{
//					Id:     strconv.FormatInt(edges[i].GetVirtualId(), 10),
//					Source: strconv.FormatInt(edges[i].GetVertices()[0].GetVirtualId(), 10),
//					Target: strconv.FormatInt(edges[i].GetVertices()[1].GetVirtualId(), 10),
//				},
//			}
//			resultsetArray = append(resultsetArray, groupNode1)
//		}
//	}
//	return resultsetArray
//}

func IsNodePresent (collection []CytoscapeDataGroupNode, nodeToTest impl.Node) (int, bool) {
	for i := 0; i < len(collection); i++ {
		currentNode := collection[i]
		if currentNode.Data["id"] == nodeToTest.EntityId {
			return i, true
		}
	}
	return -1, false
}

/*
func IsNodePresent (collection []interface{}, nodeToTest tgdb.TGNode) bool {
	for i := 0; i < len(collection); i++ {
		currentNode, ok := collection[i].(impl.Node)
		if ok {
			if currentNode.EntityId == nodeToTest.GetVirtualId() {
				return true
			}
		}
	}
	return false
}
*/



func IsEdgePresent (collection []CytoscapeDataGroupNode, edgeToTest impl.Edge) (int, bool) {
	for i := 0; i < len(collection); i++ {
		currentEdge := collection[i]

		if currentEdge.Data["id"] == edgeToTest.EntityId {
			return i, true
		}
	}
	return -1, false
}

////
//// TODO: Revisit this, and possibly merge with IsEdgePresent
////
//func IsEdgePresent_ (collection []interface{}, edgeToTest tgdb.TGEdge) bool {
//	for i := 0; i < len(collection); i++ {
//		currentEdge, ok := collection[i].(impl.Edge)
//		if ok {
//			if currentEdge.EntityId == edgeToTest.GetVirtualId() {
//				return true
//			}
//		}
//	}
//	return false
//}

