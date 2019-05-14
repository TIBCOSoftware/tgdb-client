/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"strings"
)

//-====================-//
//   Define Attribute
//-====================-//

type Attribute struct {
	_name  string
	_value interface{}
	_type  DataType
}

func (a *Attribute) GetName() string {
	return a._name
}

func (a *Attribute) SetName(myName string) {
	a._name = myName
}

func (a *Attribute) GetValue() interface{} {
	return a._value
}

func (a *Attribute) SetValue(myValue interface{}) {
	a._value = myValue
}

func (a *Attribute) GetType() DataType {
	return a._type
}

func (a *Attribute) SetType(myType DataType) {
	a._type = myType
}

func NewAttribute(attributeModel *Attributefinition, value interface{}) *Attribute {
	return &Attribute{attributeModel.name, value, attributeModel.dataType}
}

//-========================-//
//     Define EntityId
//-========================-//

type EntityId struct {
	_keyHash string
	_type    string
}

func (this *EntityId) GetType() string {
	return this._type
}

func (this *EntityId) GetKeyHash() string {
	return this._keyHash
}

func (this *EntityId) ConvertToString() string {
	var sb string
	sb = this._type
	sb += "_"
	sb += this._keyHash
	return sb
}

//func (this *EntityId) ToStringBuffer() *strings.Builder {
//	var sb strings.Builder
//	sb.WriteString(this._type)
//	sb.WriteString("_")
//	sb.WriteString(this._keyHash)
//	return &sb
//}

//-====================-//
//     Define Entity
//-====================-//

type Entity struct {
	_key        []interface{}
	_attributes map[string]*Attribute
}

func (e *Entity) SetKey(key []interface{}) {
	e._key = key
}

func (e *Entity) GetKey() []interface{} {
	return e._key
}

func (e *Entity) GetAttribute(name string) *Attribute {
	return e._attributes[name]
}

func (e *Entity) SetAttribute(name string, attr *Attribute) {
	e._attributes[name] = attr
}

func (e *Entity) GetAttributes() map[string]*Attribute {
	return e._attributes
}

//-====================-//
//     Define NodeId
//-====================-//

type NodeId struct {
	EntityId
}

func (this *NodeId) ToString() string {
	//return this.EntityId.ToStringBuffer().String()
	return this.EntityId.ConvertToString()
}

func (this *NodeId) FromString(id string) *NodeId {
	idComp := strings.Split(id, "_")
	this._type = idComp[0]
	this._keyHash = idComp[1]
	return this
}

//-====================-//
//     Define Node
//-====================-//

type Node struct {
	NodeId
	Entity
}

func NewNode(myType string, myKey []interface{}) *Node {
	var n Node
	n._key = myKey
	n._keyHash = Hash(myKey)
	n._type = myType
	n._attributes = make(map[string]*Attribute)
	return &n
}

//-====================-//
//     Define EdgeId
//-====================-//

type EdgeId struct {
	EntityId
	_fromNodeKeyHash string
	_toNodeKeyHash   string
}

func (eid *EdgeId) ToString() string {
	//sb := eid.EntityId.ToStringBuffer()
	sb := eid.EntityId.ConvertToString()
	sb += "_"
	sb += eid._fromNodeKeyHash
	sb += "_"
	sb += eid._toNodeKeyHash
	//sb.WriteString("_")
	//sb.WriteString(eid._fromNodeKeyHash)
	//sb.WriteString("_")
	//sb.WriteString(eid._toNodeKeyHash)
	//return sb.String()
	return sb
}

//-====================-//
//     Define Edge
//-====================-//

type Edge struct {
	EdgeId
	Entity
	_fromNodeId  string
	_fromNodeKey []interface{}
	_toNodeId    string
	_toNodeKey   []interface{}
}

func (this *Edge) GetFromId() string {
	return this._fromNodeId
}

func (this *Edge) GetToId() string {
	return this._toNodeId
}

func NewEdge(myType string, myKey []interface{}, fromNode *Node, toNode *Node) *Edge {
	var e Edge

	e._fromNodeId = fromNode.NodeId.ToString()
	e._fromNodeKey = fromNode.GetKey()
	e._fromNodeKeyHash = fromNode.GetKeyHash()
	e._toNodeId = toNode.NodeId.ToString()
	e._toNodeKey = toNode.GetKey()
	e._toNodeKeyHash = toNode.GetKeyHash()
	e._type = myType
	e._key = myKey
	e._keyHash = Hash(myKey)
	e._attributes = make(map[string]*Attribute)

	return &e
}

//-====================-//
//     Define Graph
//-====================-//

type Graph interface {
	GetId() string
	GetModelId() string
	GetModel() map[string]interface{}
	GetNodes() map[NodeId]*Node
	GetEdges() map[EdgeId]*Edge
	UpsertGraph(graph map[string]interface{})
	UpsertNode(nodeType string, nodeKey []interface{}) *Node
	UpsertEdge(edgeType string, edgeKey []interface{}, fromNode *Node, toNode *Node) *Edge
	GetNodeByTypeByKey(nodeType string, nodeId NodeId) *Node
	GetNodesByType(nodeType string) map[NodeId]*Node
	GetEntityKeyNamesForNode(entityName string) []string
	GetEntityKeyNamesForEdge(entityName string) []string
}

//-=========================-//
//     Define GraphFragment
//-=========================-//

type GraphFragment struct {
	id    string
	edges map[string]*Edge
	nodes map[string]*Node
}

func (g *GraphFragment) GetId() string {
	return g.id
}

func NewGraphFragment(id string) *GraphFragment {
	var g GraphFragment
	g.id = id
	g.edges = make(map[string]*Edge)
	g.nodes = make(map[string]*Node)
	return &g
}

//-=========================-//
//     Define GraphImpl
//-=========================-//

type GraphImpl struct {
	id          string
	modelId     string
	model       map[string]interface{}
	edges       map[EdgeId]*Edge
	nodes       map[NodeId]*Node
	edgesByType map[string](map[EdgeId]*Edge)
	nodesByType map[string](map[NodeId]*Node)
}

func (g *GraphImpl) GetId() string {
	return g.id
}

func (g *GraphImpl) GetModelId() string {
	return g.modelId
}

func (g *GraphImpl) SetModel(model map[string]interface{}) {
	g.model = model
}

func (g *GraphImpl) GetModel() map[string]interface{} {
	return g.model
}

func (g *GraphImpl) UpsertGraph(graph map[string]interface{}) {

}

func (g *GraphImpl) UpsertNode(nodeType string, nodeKey []interface{}) *Node {
	node := NewNode(nodeType, nodeKey)
	if nil != g.nodes[node.NodeId] {
		return g.nodes[node.NodeId]
	}
	g.SetNode(node.NodeId, node)
	return node
}

func (g *GraphImpl) GetNode(id NodeId) *Node {
	return g.nodes[id]
}

func (g *GraphImpl) GetNodeByTypeByKey(nodeType string, nodeId NodeId) *Node {
	typedNodeById := g.nodesByType[nodeType]
	return typedNodeById[nodeId]
}

func (g *GraphImpl) GetNodesByType(nodeType string) map[NodeId]*Node {
	return g.nodesByType[nodeType]
}

func (g *GraphImpl) GetNodes() map[NodeId]*Node {
	return g.nodes
}

func (g *GraphImpl) SetNode(id NodeId, node *Node) {
	g.nodes[id] = node
	nodeMap := g.nodesByType[node._type]
	if nil == nodeMap {
		nodeMap = make(map[NodeId]*Node)
		g.nodesByType[node._type] = nodeMap
	}
	nodeMap[id] = node
}

func (g *GraphImpl) UpsertEdge(edgeType string, edgeKey []interface{}, fromNode *Node, toNode *Node) *Edge {
	edge := NewEdge(edgeType, edgeKey, fromNode, toNode)
	if nil != g.edges[edge.EdgeId] {
		return g.edges[edge.EdgeId]
	}
	g.SetEdge(edge.EdgeId, edge)
	return edge
}

func (g *GraphImpl) GetEdge(id EdgeId) *Edge {
	return g.edges[id]
}

func (g *GraphImpl) SetEdge(id EdgeId, edge *Edge) {
	g.edges[id] = edge
	edgeMap := g.edgesByType[edge._type]
	if nil == edgeMap {
		edgeMap = make(map[EdgeId]*Edge)
		g.edgesByType[edge._type] = edgeMap
	}
	edgeMap[id] = edge
}

func (g *GraphImpl) GetEdges() map[EdgeId]*Edge {
	return g.edges
}

func (g *GraphImpl) GetEntityKeyNamesForNode(nodeType string) []string {
	nodeKeyMap := g.model["nodes"].(map[string]interface{})["keyMap"].(map[string][]string)
	return nodeKeyMap[nodeType]
}

func (g *GraphImpl) GetEntityKeyNamesForEdge(edgeType string) []string {
	edgeKeyMap := g.model["edges"].(map[string]interface{})["keyMap"].(map[string][]string)
	return edgeKeyMap[edgeType]
}

func NewGraphImpl(modelId string, id string) *GraphImpl {
	var g GraphImpl
	g.id = id
	g.modelId = modelId
	g.edges = make(map[EdgeId]*Edge)
	g.nodes = make(map[NodeId]*Node)
	g.edgesByType = make(map[string](map[EdgeId]*Edge))
	g.nodesByType = make(map[string](map[NodeId]*Node))
	return &g
}
