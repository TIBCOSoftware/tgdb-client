package tgdb

import (
	"reflect"

	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

func BuildMetadata(metadata types.TGGraphMetadata) map[string]interface{} {
	data := make(map[string]interface{})
	edgeTypes, _ := metadata.GetEdgeTypes()
	edgeTypeInfos := make([]map[string]interface{}, len(edgeTypes))
	data["edgeTypes"] = edgeTypeInfos
	for index, edgeType := range edgeTypes {
		edgeTypeInfo := make(map[string]interface{})
		edgeTypeInfo["id"] = edgeType.GetEntityTypeId()
		edgeTypeInfo["name"] = edgeType.GetName()
		edgeTypeInfo["systemType"] = edgeType.GetSystemType()

		fromNode := edgeType.GetFromNodeType()
		if nil != fromNode {
			fromNodeInfo := make(map[string]interface{})
			edgeTypeInfo["fromNodeType"] = fromNodeInfo
			fromNodeInfo["id"] = fromNode.GetEntityTypeId()
			fromNodeInfo["name"] = fromNode.GetName()
		}

		toNode := edgeType.GetToNodeType()
		if nil != toNode {
			toNodeInfo := make(map[string]interface{})
			edgeTypeInfo["toNodeType"] = toNodeInfo
			toNodeInfo["id"] = toNode.GetEntityTypeId()
			toNodeInfo["name"] = toNode.GetName()
		}

		attributeDescriptors := edgeType.GetAttributeDescriptors()
		attributeDescriptorInfos := make([]map[string]interface{}, len(attributeDescriptors))
		edgeTypeInfo["attributeDescriptors"] = attributeDescriptorInfos
		for index, attributeDescriptor := range attributeDescriptors {
			attributeDescriptorInfo := make(map[string]interface{})
			attributeDescriptorInfo["name"] = attributeDescriptor.GetName()
			attributeDescriptorInfo["type"] = attributeDescriptor.GetAttrType()
			attributeDescriptorInfo["attributeId"] = attributeDescriptor.GetAttributeId()
			attributeDescriptorInfo["scale"] = attributeDescriptor.GetScale()
			attributeDescriptorInfo["precision"] = attributeDescriptor.GetPrecision()
			attributeDescriptorInfo["systemType"] = attributeDescriptor.GetSystemType()
			attributeDescriptorInfos[index] = attributeDescriptorInfo
		}

		edgeTypeInfos[index] = edgeTypeInfo
	}

	nodeTypes, _ := metadata.GetNodeTypes()
	nodeTypeInfos := make([]map[string]interface{}, len(nodeTypes))
	data["nodeTypes"] = nodeTypeInfos
	for index, nodeType := range nodeTypes {
		nodeTypeInfo := make(map[string]interface{})
		nodeTypeInfo["id"] = nodeType.GetEntityTypeId()
		nodeTypeInfo["name"] = nodeType.GetName()
		nodeTypeInfo["systemType"] = nodeType.GetSystemType()

		pkeyAttributeDescriptors := nodeType.GetPKeyAttributeDescriptors()
		pkeyAttributeDescriptorInfos := make([]map[string]interface{}, len(pkeyAttributeDescriptors))
		nodeTypeInfo["pkeyAttributeDescriptors"] = pkeyAttributeDescriptorInfos
		for index, pkeyAttributeDescriptor := range pkeyAttributeDescriptors {
			pkeyAttributeDescriptorInfo := make(map[string]interface{})
			pkeyAttributeDescriptorInfo["name"] = pkeyAttributeDescriptor.GetName()
			pkeyAttributeDescriptorInfo["type"] = pkeyAttributeDescriptor.GetAttrType()
			pkeyAttributeDescriptorInfo["attributeId"] = pkeyAttributeDescriptor.GetAttributeId()
			pkeyAttributeDescriptorInfo["scale"] = pkeyAttributeDescriptor.GetScale()
			pkeyAttributeDescriptorInfo["precision"] = pkeyAttributeDescriptor.GetPrecision()
			pkeyAttributeDescriptorInfo["systemType"] = pkeyAttributeDescriptor.GetSystemType()
			pkeyAttributeDescriptorInfos[index] = pkeyAttributeDescriptorInfo
		}

		attributeDescriptors := nodeType.GetAttributeDescriptors()
		attributeDescriptorInfos := make([]map[string]interface{}, len(attributeDescriptors))
		nodeTypeInfo["attributeDescriptors"] = attributeDescriptorInfos
		for index, attributeDescriptor := range attributeDescriptors {
			attributeDescriptorInfo := make(map[string]interface{})
			attributeDescriptorInfo["name"] = attributeDescriptor.GetName()
			attributeDescriptorInfo["type"] = attributeDescriptor.GetAttrType()
			attributeDescriptorInfo["attributeId"] = attributeDescriptor.GetAttributeId()
			attributeDescriptorInfo["scale"] = attributeDescriptor.GetScale()
			attributeDescriptorInfo["precision"] = attributeDescriptor.GetPrecision()
			attributeDescriptorInfo["systemType"] = attributeDescriptor.GetSystemType()
			attributeDescriptorInfos[index] = attributeDescriptorInfo
		}

		nodeTypeInfos[index] = nodeTypeInfo
	}
	/*
		attributeDescriptors, _ := metadata.GetAttributeDescriptors()
		attributeDescriptorInfos := make([]map[string]interface{}, len(attributeDescriptors))
		data["attributeDescriptors"] = attributeDescriptorInfos
		for index, attributeDescriptor := range attributeDescriptors {
			attributeDescriptorInfo := make(map[string]interface{})
			attributeDescriptorInfo["name"] = attributeDescriptor.GetName()
			attributeDescriptorInfo["type"] = attributeDescriptor.GetAttrType()
			attributeDescriptorInfo["attributeId"] = attributeDescriptor.GetAttributeId()
			attributeDescriptorInfo["scale"] = attributeDescriptor.GetScale()
			attributeDescriptorInfo["precision"] = attributeDescriptor.GetPrecision()
			attributeDescriptorInfo["systemType"] = attributeDescriptor.GetSystemType()
			attributeDescriptorInfos[index] = attributeDescriptorInfo
		}
	*/
	return data
}

func BuildNode(tgdb *TGDBService, node types.TGNode) map[string]interface{} {
	tgResult := make(map[string][]types.TGEntity)
	tgResult["nodes"] = make([]types.TGEntity, 0)
	tgResult["edges"] = make([]types.TGEntity, 0)
	traverse(tgResult, node, 0)
	return buildResult(tgdb, tgResult)
}

func traverse(
	tgResult map[string]([]types.TGEntity),
	entity types.TGEntity,
	currDepth int) {
	if "types.TGNode" == reflect.TypeOf(entity).String() {
		node := entity.(types.TGNode)
		tgResult["nodes"] = append(tgResult["nodes"], node)
		for _, edge := range node.GetEdges() {
			if !contains(tgResult["edges"], edge) {
				currDepth += 1
				traverse(tgResult, edge, currDepth)
			}
		}
	} else if "types.TGEdge" == reflect.TypeOf(entity).String() {
		edge := entity.(types.TGEdge)
		tgResult["edges"] = append(tgResult["edges"], edge)
		for _, node := range edge.GetVertices() {
			if !contains(tgResult["nodes"], node) {
				traverse(tgResult, node, currDepth)
			}
		}
	}
}

func buildResult(tgdb *TGDBService, tgResult map[string]([]types.TGEntity)) map[string]interface{} {
	result := make(map[string]interface{})
	result["nodes"] = make([]map[string]interface{}, 0)
	result["edges"] = make([]map[string]interface{}, 0)

	for _, entity := range tgResult["nodes"] {
		node := entity.(types.TGNode)
		if nil != node && nil != node.GetEntityType() {
			aNode := make(map[string]interface{})
			nodeInfos := result["nodes"].([]map[string]interface{})
			nodeInfos = append(nodeInfos, aNode)
			aNode["_type"] = node.GetEntityType().GetName()
			id := ExtractNodeKeyAttrValue(tgdb, node)
			if nil != id {
				aNode["id"] = id
			}
			attributes, _ := node.GetAttributes()
			for _, attr := range attributes {
				aNode[attr.GetAttributeDescriptor().GetName()] = attr.GetValue()
			}
		}
	}

	for _, entity := range tgResult["edges"] {
		edge := entity.(types.TGEdge)
		if nil != edge.GetEntityType() {
			fromNodeId := ExtractNodeKeyAttrValue(tgdb, edge.GetVertices()[0])
			toNodeId := ExtractNodeKeyAttrValue(tgdb, edge.GetVertices()[1])
			if nil == fromNodeId || nil == toNodeId {
				continue
			}

			anEdge := make(map[string]interface{})
			anEdge["_type"] = edge.(types.TGEdge).GetEntityType().GetName()
			anEdge["fromNode"] = fromNodeId
			anEdge["toNode"] = toNodeId
			attributes, _ := edge.GetAttributes()
			for _, attr := range attributes {
				anEdge[attr.GetAttributeDescriptor().GetName()] = attr.GetValue()
			}
			edgeInfos := result["edges"].([]map[string]interface{})
			edgeInfos = append(edgeInfos, anEdge)
		}
	}

	return result
}

func contains(arrays []types.TGEntity, target types.TGEntity) bool {
	for _, element := range arrays {
		if element == target {
			return true
		}
	}
	return false
}

func ExtractNodeKeyAttrValue(tgdb *TGDBService, node types.TGNode) []interface{} {
	keyFields := tgdb.GetNodeKeyfields(node.GetEntityType().String())
	if nil == keyFields {
		return nil
	}
	key := make([]interface{}, 0)
	for _, keyField := range keyFields {
		key = append(key, node.GetAttribute(keyField).GetValue())
	}

	return key
}
