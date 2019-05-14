/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"encoding/json"
	"fmt"
)

//-================================-//
//     Define Attributefinition
//-================================-//

type Attributefinition struct {
	name     string
	dataType DataType
}

func (this *Attributefinition) GetDataType() DataType {
	return this.dataType
}

func NewAttributeModel(name string, dataTypeStr string) *Attributefinition {
	dataType, _ := ToTypeEnum(dataTypeStr)
	return &Attributefinition{name, dataType}
}

//-================================-//
//     Define EntityDefinition
//-===============================-//

type EntityDefinition struct {
	_keyDefinition []string
	_type          string
	_skipCondition string
	_attributes    map[string]*Attributefinition
}

func (this *EntityDefinition) GetAttributeDefinitions() map[string]*Attributefinition {
	return this._attributes
}

func NewEntityModel(skipCondition string, entityInfo map[string]interface{}) *EntityDefinition {
	var key []string
	if nil != entityInfo["key"] {
		keyInfo := entityInfo["key"].([]interface{})
		key = make([]string, len(keyInfo))
		for i := 0; i < len(key); i++ {
			key[i] = keyInfo[i].(string)
		}
	}

	entityType := entityInfo["name"].(string)

	attributesModel := make(map[string]*Attributefinition)

	if nil != entityInfo["attributes"] {
		attributesInfo := entityInfo["attributes"].([]interface{})
		for _, attributeInfo := range attributesInfo {
			attribute := attributeInfo.(map[string]interface{})
			attrName := attribute["name"].(string)

			var attrType string
			if nil != attribute["type"] {
				attrType = attribute["type"].(string)
			} else {
				attrType = "String"
			}
			attributesModel[attrName] = NewAttributeModel(attrName, attrType)
		}
	}

	return &EntityDefinition{key, entityType, skipCondition, attributesModel}
}

//-============================-//
//     Define NodeDefinition
//-============================-//

type NodeDefinition struct {
	*EntityDefinition
}

func NewNodeModel(skipCondition string, nodeInfo map[string]interface{}) *NodeDefinition {
	var nodeModel NodeDefinition
	nodeModel.EntityDefinition = NewEntityModel(skipCondition, nodeInfo)
	return &nodeModel
}

//-============================-//
//     Define NodeDefinition
//-============================-//

type EdgeDefinition struct {
	_fromNodeType string
	_toNodeType   string
	*EntityDefinition
}

func NewEdgeModel(skipCondition string, edgeInfo map[string]interface{}) *EdgeDefinition {
	var edgeModel EdgeDefinition
	edgeModel._fromNodeType = edgeInfo["from"].(string)
	edgeModel._toNodeType = edgeInfo["to"].(string)
	edgeModel.EntityDefinition = NewEntityModel(skipCondition, edgeInfo)

	return &edgeModel
}

//-============================-//
//     Define GraphDefinition
//-============================-//

type GraphDefinition struct {
	_id              string
	_nodeDefinitions map[string]*NodeDefinition
	_edgeDefinitions map[string]*EdgeDefinition
}

func (gd *GraphDefinition) GetId() string {
	return gd._id
}

func (gd *GraphDefinition) GetNodeDefinition(nodeType string) *NodeDefinition {
	return gd._nodeDefinitions[nodeType]
}

func (gd *GraphDefinition) GetNodeDefinitions() map[string]*NodeDefinition {
	return gd._nodeDefinitions
}

func (gd *GraphDefinition) GetEdgeDefinition(edgeType string) *EdgeDefinition {
	return gd._edgeDefinitions[edgeType]
}

func (gd *GraphDefinition) GetEdgeDefinitions() map[string]*EdgeDefinition {
	return gd._edgeDefinitions
}

func (gd *GraphDefinition) Export() map[string]interface{} {
	nodeTypes := make([]string, 0)
	nodeKeyMap := make(map[string][]string)
	attrTypeMap := make(map[string](map[string]string))
	for nodeType, definition := range gd._nodeDefinitions {
		nodeTypes = append(nodeTypes, nodeType)
		nodeKeyMap[nodeType] = definition._keyDefinition
		nodeAttrTypeMap := make(map[string]string)
		attrTypeMap[nodeType] = nodeAttrTypeMap
		for attrName, attrDef := range definition._attributes {
			nodeAttrTypeMap[attrName] = attrDef.dataType.String()
		}
	}
	nodeModels := make(map[string]interface{})
	nodeModels["types"] = nodeTypes
	nodeModels["keyMap"] = nodeKeyMap
	nodeModels["attrTypeMap"] = attrTypeMap

	edgeTypes := make([]string, 0)
	edgeKeyMap := make(map[string][]string)
	attrTypeMap = make(map[string](map[string]string))
	for edgeType, definition := range gd._edgeDefinitions {
		edgeTypes = append(edgeTypes, edgeType)
		edgeKeyMap[edgeType] = definition._keyDefinition
		edgeAttrTypeMap := make(map[string]string)
		attrTypeMap[edgeType] = edgeAttrTypeMap
		for attrName, attrDef := range definition._attributes {
			edgeAttrTypeMap[attrName] = attrDef.dataType.String()
		}
	}
	edgeModels := make(map[string]interface{})
	edgeModels["types"] = edgeTypes
	edgeModels["keyMap"] = edgeKeyMap
	edgeModels["attrTypeMap"] = attrTypeMap

	graphModel := make(map[string]interface{})
	graphModel["nodes"] = nodeModels
	graphModel["edges"] = edgeModels

	return graphModel
}

func NewGraphModel(id string, graphmodel string) *GraphDefinition {
	var rootObject interface{}
	err := json.Unmarshal([]byte(graphmodel), &rootObject)
	if nil != err {
		return nil
	}

	return parseTGBModel(id, rootObject)
}

func parseTGBModel(id string, rootObject interface{}) *GraphDefinition {
	fmt.Println("id = ", id, ", root obj = ", rootObject)
	dataMap := rootObject.(map[string]interface{})

	nodeModels := make(map[string]*NodeDefinition)
	nodes := dataMap["nodes"].([]interface{})
	for _, node := range nodes {
		nodeInfo := node.(map[string]interface{})
		nodeType := nodeInfo["name"].(string)
		nodeModels[nodeType] = NewNodeModel("skipCondition", nodeInfo)
		fmt.Println("nodeType = ", nodeType, ", nodeInfo = ", nodeInfo)
	}

	edgeModels := make(map[string]*EdgeDefinition)
	edges := dataMap["edges"].([]interface{})
	for _, edge := range edges {
		edgeInfo := edge.(map[string]interface{})
		fmt.Println("edgeType = ", edgeInfo["name"], ", edgeInfo = ", edgeInfo)
		edgeModels[edgeInfo["name"].(string)] = NewEdgeModel("skipCondition", edgeInfo)
	}
	fmt.Println("nodeModels = ", nodeModels, ", edgeModels = ", edgeModels)

	graphModel := &GraphDefinition{id, nodeModels, edgeModels}
	fmt.Println("graphModel = ", graphModel)
	return graphModel
}
