/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/tgdb/lib/util"
)

//-============================-//
//     Define GraphBuilder
//-============================-//

type GraphBuilder struct {
	mux sync.Mutex
}

func NewGraphBuilder() *GraphBuilder {
	builder := &GraphBuilder{}
	return builder
}

func (builder *GraphBuilder) CreateGraph(graphId string, model *GraphDefinition) Graph {
	graph := NewGraphImpl(model.GetId(), graphId)
	graph.SetModel(model.Export())
	return graph
}

func (builder *GraphBuilder) CreateUndefinedGraph(modelId string, graphId string) Graph {
	graph := NewGraphImpl(modelId, graphId)
	return graph
}

func (builder *GraphBuilder) BuildGraph(graph Graph, model *GraphDefinition, nodes interface{}, edges interface{}) {
	nodesWrapper, nodesWrapperValid := nodes.([]interface{})
	if nodesWrapperValid {
		if nil != nodesWrapper[0] {
			nodes, nodesValid := nodesWrapper[0].(map[string]interface{})
			if nodesValid {
				for nodekey, node := range nodes {

					fmt.Println("nodekey before = ", nodekey)
					checkPos := strings.LastIndex(nodekey, "_")
					if 0 <= checkPos && util.IsInteger(nodekey[checkPos+1:len(nodekey)]) {
						nodekey = string(nodekey[0:checkPos])
					}
					fmt.Println("nodekey after = " + nodekey)

					var attrs map[string]interface{}
					attrWrapper, attrWrapperValid := node.([]interface{})
					if attrWrapperValid {
						if nil != attrWrapper[0] {
							attrs = attrWrapper[0].(map[string]interface{})
						}
					}
					if nil == attrs {
						attrs = make(map[string]interface{})
					}
					builder.BuildNode(graph, model, nodekey, attrs)
				}
			}
		}
	}

	edgesWrapper, edgesWrapperValid := edges.([]interface{})
	fmt.Println("edgesWrapper = ", edgesWrapper, ", edgesWrapperValid = ", edgesWrapperValid)
	if edgesWrapperValid {
		if nil != edgesWrapper[0] {
			edges, edgesValid := edgesWrapper[0].(map[string]interface{})
			fmt.Println("edges = ", edges, ", edgesValid = ", edgesValid)
			if edgesValid {
				for edgekey, edge := range edges {
					var attrs map[string]interface{}
					attrWrapper, attrWrapperValid := edge.([]interface{})
					fmt.Println("attrWrapper = ", attrWrapper, ", attrWrapperValid = ", attrWrapperValid)
					if attrWrapperValid {
						if nil != attrWrapper[0] {
							attrs = attrWrapper[0].(map[string]interface{})
						}
					}
					if nil == attrs {
						attrs = make(map[string]interface{})
					}
					fmt.Println("attrs = ", attrs)
					builder.BuildEdge(graph, model, edgekey, attrs)
				}
			}
		}
	}

}

func (builder *GraphBuilder) BuildNode(graph Graph, model *GraphDefinition, nodeType string, attributesInfo map[string]interface{}) {
	builder.mux.Lock()
	nodeModel := model.GetNodeDefinition(nodeType)
	keyDefinition := nodeModel._keyDefinition
	key := make([]interface{}, len(keyDefinition))
	for i := 0; i < len(keyDefinition); i++ {
		key[i] = attributesInfo[keyDefinition[i]]
	}

	node := graph.UpsertNode(nodeType, key)
	for attrKey, attrVal := range attributesInfo {
		attribute := NewAttribute(nodeModel._attributes[attrKey], attrVal)
		node._attributes[attribute.GetName()] = attribute
	}
	builder.mux.Unlock()
}

func (builder *GraphBuilder) BuildEdge(graph Graph, model *GraphDefinition, edgeType string, attributesInfo map[string]interface{}) {
	builder.mux.Lock()
	defer builder.mux.Unlock()

	edgeModel := model.GetEdgeDefinition(edgeType)
	keyDefinition := edgeModel._keyDefinition

	fromNodes, toNodes := builder.buildVerexes(graph, model, edgeType, attributesInfo)

	fmt.Println("from = ", fromNodes, ", to = ", toNodes)

	var fromNode *Node
	var toNode *Node

	for _, fromNode = range fromNodes {
		for _, toNode = range toNodes {
			/* Allow duplicate? */
			key := make([]interface{}, len(keyDefinition))
			for i := 0; i < len(keyDefinition); i++ {
				key[i] = attributesInfo[keyDefinition[i]]
			}
			edge := graph.UpsertEdge(edgeType, key, fromNode, toNode)
			for attrKey, attrVal := range attributesInfo {
				attribute := NewAttribute(edgeModel._attributes[attrKey], attrVal)
				edge._attributes[attribute.GetName()] = attribute
			}
		}
	}
}

func (builder *GraphBuilder) buildVerexes(
	graph Graph,
	model *GraphDefinition,
	edgeType string,
	attributesInfo map[string]interface{}) (map[NodeId]*Node, map[NodeId]*Node) {

	edgeModel := model.GetEdgeDefinition(edgeType)

	var fromNodeKey []interface{}
	var fromNodeKeyAttrValues map[string]interface{}
	if nil != attributesInfo["from"] {
		rawKeyArray := attributesInfo["from"].([]interface{})
		if 0 < len(rawKeyArray) {
			fromNodeKeyAttrValues = rawKeyArray[0].(map[string]interface{})
			fromNodeKey = make([]interface{}, len(fromNodeKeyAttrValues))
			for index, keyElement := range model.GetNodeDefinition(edgeModel._fromNodeType)._keyDefinition {
				fromNodeKey[index] = fromNodeKeyAttrValues[keyElement]
			}
		}

		delete(attributesInfo, "from")
	}
	var toNodeKey []interface{}
	var toNodeKeyAttrValues map[string]interface{}
	if nil != attributesInfo["to"] {
		rawKeyArray := attributesInfo["to"].([]interface{})
		if 0 < len(rawKeyArray) {
			toNodeKeyAttrValues = rawKeyArray[0].(map[string]interface{})
			toNodeKey = make([]interface{}, len(toNodeKeyAttrValues))
			for index, keyElement := range model.GetNodeDefinition(edgeModel._toNodeType)._keyDefinition {
				toNodeKey[index] = toNodeKeyAttrValues[keyElement]
			}
		}

		delete(attributesInfo, "to")
	}

	var fromNodes map[NodeId]*Node
	var toNodes map[NodeId]*Node
	if nil != fromNodeKey && nil != toNodeKey {
		fromNodeId := NewNodeId(edgeModel._fromNodeType, fromNodeKey)
		fromNode := graph.GetNodeByTypeByKey(edgeModel._fromNodeType, fromNodeId)
		if nil == fromNode {
			fromNode = graph.UpsertNode(edgeModel._fromNodeType, fromNodeKey)
			for attrKey, attrVal := range fromNodeKeyAttrValues {
				attribute := NewAttribute(model.GetNodeDefinition(edgeModel._fromNodeType)._attributes[attrKey], attrVal)
				fromNode._attributes[attribute.GetName()] = attribute
			}

		}
		fromNodes = make(map[NodeId]*Node)
		fromNodes[fromNodeId] = fromNode

		toNodeId := NewNodeId(edgeModel._toNodeType, toNodeKey)
		toNode := graph.GetNodeByTypeByKey(edgeModel._toNodeType, toNodeId)
		if nil == toNode {
			toNode = graph.UpsertNode(edgeModel._toNodeType, toNodeKey)
			for attrKey, attrVal := range toNodeKeyAttrValues {
				attribute := NewAttribute(model.GetNodeDefinition(edgeModel._toNodeType)._attributes[attrKey], attrVal)
				toNode._attributes[attribute.GetName()] = attribute
			}
		}
		toNodes = make(map[NodeId]*Node)
		toNodes[toNodeId] = toNode
	} else {
		fromNodes = graph.GetNodesByType(edgeModel._fromNodeType)
		toNodes = graph.GetNodesByType(edgeModel._toNodeType)
	}
	return fromNodes, toNodes
}

func (builder *GraphBuilder) buildVerexesBak(
	graph Graph,
	model *GraphDefinition,
	edgeType string,
	attributesInfo map[string]interface{}) (map[NodeId]*Node, map[NodeId]*Node) {

	edgeModel := model.GetEdgeDefinition(edgeType)

	var fromNodeKey []interface{}
	if nil != attributesInfo["from"] {
		rawKeyArray := attributesInfo["from"].([]interface{})
		if 0 < len(rawKeyArray) {
			fromNodeKeyAttrValues := rawKeyArray[0].(map[string]interface{})
			fromNodeKey = make([]interface{}, len(fromNodeKeyAttrValues))
			for index, keyElement := range model.GetNodeDefinition(edgeModel._fromNodeType)._keyDefinition {
				fromNodeKey[index] = fromNodeKeyAttrValues[keyElement]
			}
		}

		delete(attributesInfo, "from")
	}
	var toNodeKey []interface{}
	if nil != attributesInfo["to"] {
		rawKeyArray := attributesInfo["to"].([]interface{})
		if 0 < len(rawKeyArray) {
			toNodeKeyAttrValues := rawKeyArray[0].(map[string]interface{})
			toNodeKey = make([]interface{}, len(toNodeKeyAttrValues))
			for index, keyElement := range model.GetNodeDefinition(edgeModel._toNodeType)._keyDefinition {
				toNodeKey[index] = toNodeKeyAttrValues[keyElement]
			}
		}

		delete(attributesInfo, "to")
	}

	var fromNodes map[NodeId]*Node
	var toNodes map[NodeId]*Node
	if nil != fromNodeKey && nil != toNodeKey {
		fromNodeId := NewNodeId(edgeModel._fromNodeType, fromNodeKey)
		fromNode := graph.GetNodeByTypeByKey(edgeModel._fromNodeType, fromNodeId)
		if nil == fromNode {
			//fromNode = NewNode(edgeModel._fromNodeType, fromNodeKey)
			fromNode = graph.UpsertNode(edgeModel._fromNodeType, fromNodeKey)
		}
		fromNodes = make(map[NodeId]*Node)
		fromNodes[fromNodeId] = fromNode

		toNodeId := NewNodeId(edgeModel._toNodeType, toNodeKey)
		toNode := graph.GetNodeByTypeByKey(edgeModel._toNodeType, toNodeId)
		if nil == toNode {
			//toNode = NewNode(edgeModel._toNodeType, toNodeKey)
			toNode = graph.UpsertNode(edgeModel._toNodeType, toNodeKey)
		}
		toNodes = make(map[NodeId]*Node)
		toNodes[toNodeId] = toNode
	} else {
		fromNodes = graph.GetNodesByType(edgeModel._fromNodeType)
		toNodes = graph.GetNodesByType(edgeModel._toNodeType)
	}
	return fromNodes, toNodes
}

func (builder *GraphBuilder) Export(g Graph, graphModel *GraphDefinition) map[string]interface{} {

	nodeDefinitions := graphModel._nodeDefinitions
	edgeDefinitions := graphModel._edgeDefinitions

	data := make(map[string]interface{})
	data["id"] = g.GetId()
	data["modelId"] = g.GetModelId()
	data["model"] = graphModel.Export()

	nodesData := make(map[string]interface{})
	for nodeId, node := range g.GetNodes() {
		nodeData := make(map[string]interface{})
		attrsData := make(map[string]interface{})
		for attrName, attribute := range node._attributes {
			attrData := make(map[string]interface{})
			attrData["name"] = attribute._name
			attrData["value"] = attribute._value
			attrData["type"] = attribute._type.String()
			attrsData[attrName] = attrData
		}
		nodeData["type"] = node._type
		nodeData["keyAttributeName"] = nodeDefinitions[node._type]._keyDefinition
		nodeData["key"] = node._key
		nodeData["attributes"] = attrsData
		nodesData[nodeId.ToString()] = nodeData
	}
	data["nodes"] = nodesData

	edgesData := make(map[string]interface{})
	for edgeId, edge := range g.GetEdges() {
		edgeData := make(map[string]interface{})
		attrsData := make(map[string]interface{})
		for attrName, attribute := range edge._attributes {
			attrData := make(map[string]interface{})
			attrData["name"] = attribute._name
			attrData["value"] = attribute._value
			attrData["type"] = attribute._type.String()
			attrsData[attrName] = attrData
		}
		edgeData["type"] = edge._type
		edgeData["from"] = edge._fromNodeId
		edgeData["to"] = edge._toNodeId
		edgeData["keyAttributeName"] = edgeDefinitions[edge._type]._keyDefinition
		edgeData["key"] = edge._key
		edgeData["attributes"] = attrsData
		edgesData[edgeId.ToString()] = edgeData
	}
	data["edges"] = edgesData

	//log.Debug("[GraphBuilder::Export] graph : ", data)

	return data
}

func ReconstructGraph(graphData map[string]interface{}) Graph {

	graph := NewGraphImpl(graphData["modelId"].(string), graphData["id"].(string))
	graph.SetModel(graphData["model"].(map[string]interface{}))

	nodes := util.CastGenMap(graphData["nodes"])
	for _, value := range nodes {
		nodeData := util.CastGenMap(value)
		node := NewNode(
			util.CastString(nodeData["type"]),
			util.CastGenArray(nodeData["key"]),
		)

		attributes := util.CastGenMap(nodeData["attributes"])
		for attrName, value := range attributes {
			attrData := util.CastGenMap(value)
			dataType, ok := ToTypeEnum(util.CastString(attrData["type"]))
			if !ok {
				dataType = TypeString
			}
			attribute := Attribute{
				_name:  util.CastString(attrData["name"]),
				_value: attrData["value"],
				_type:  dataType,
			}
			node.SetAttribute(attrName, &attribute)
		}
		graph.SetNode(node.NodeId, node)
	}

	edges := util.CastGenMap(graphData["edges"])
	for _, value := range edges {
		edgeData := util.CastGenMap(value)
		//fromId := *
		edge := NewEdge(
			util.CastString(edgeData["type"]),
			util.CastGenArray(edgeData["key"]),
			graph.GetNode(*(&NodeId{}).FromString(util.CastString(edgeData["from"]))),
			graph.GetNode(*(&NodeId{}).FromString(util.CastString(edgeData["to"]))),
		)

		attributes := util.CastGenMap(edgeData["attributes"])
		for attrName, value := range attributes {
			attrData := util.CastGenMap(value)
			dataType, ok := ToTypeEnum(util.CastString(attrData["type"]))
			if !ok {
				dataType = TypeString
			}
			attribute := Attribute{
				_name:  util.CastString(attrData["name"]),
				_value: attrData["value"],
				_type:  dataType,
			}
			edge.SetAttribute(attrName, &attribute)
		}
		graph.SetEdge(edge.EdgeId, edge)
	}

	return graph
}
