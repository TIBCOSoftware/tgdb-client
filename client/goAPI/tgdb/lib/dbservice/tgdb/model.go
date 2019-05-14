package tgdb

import (
	"fmt"

	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

type TGEntityWrapper struct {
	id    string
	isNew bool
}

type TGNodeWrapper struct {
	TGEntityWrapper
	node *types.TGNode
}

func NewTGNodeWrapper(node types.TGNode, isNew bool) *TGNodeWrapper {
	nodeWrapper := TGNodeWrapper{}
	nodeWrapper.isNew = isNew
	nodeWrapper.SetNode(&node)
	return &nodeWrapper
}

func (this *TGNodeWrapper) SetNode(node *types.TGNode) {
	this.node = node
}

func (this *TGNodeWrapper) GetNode() *types.TGNode {
	return this.node
}

type TGEdgeWrapper struct {
	TGEntityWrapper
	edge *types.TGEdge
}

func NewTGEdgeWrapper(edge types.TGEdge, isNew bool) *TGEdgeWrapper {
	edgeWrapper := TGEdgeWrapper{}
	edgeWrapper.isNew = isNew
	edgeWrapper.SetEdge(&edge)
	return &edgeWrapper
}

func (this *TGEdgeWrapper) SetEdge(edge *types.TGEdge) {
	this.edge = edge
}

func (this *TGEdgeWrapper) GetEdge() *types.TGEdge {
	return this.edge
}

type ReadyEntityKeeper struct {
	nodes map[string](map[interface{}]*TGNodeWrapper)
	edges map[string](map[interface{}]*TGEdgeWrapper)
}

func NewReadyEntityKeeper() ReadyEntityKeeper {
	return ReadyEntityKeeper{
		nodes: make(map[string](map[interface{}]*TGNodeWrapper)),
		edges: make(map[string](map[interface{}]*TGEdgeWrapper))}
}

func (r *ReadyEntityKeeper) AddNode(nodeTypeStr string, key interface{}, tgnode *TGNodeWrapper) {

	cachedNodePerType := r.nodes[nodeTypeStr]
	if nil == cachedNodePerType {
		cachedNodePerType = make(map[interface{}]*TGNodeWrapper)
		r.nodes[nodeTypeStr] = cachedNodePerType
	}
	cachedNodePerType[key] = tgnode
}

func (r *ReadyEntityKeeper) AddEdge(edgeTypeStr string, key interface{}, tgedge *TGEdgeWrapper) {
	cachedEdgePerType := r.edges[edgeTypeStr]
	if nil == cachedEdgePerType {
		cachedEdgePerType = make(map[interface{}]*TGEdgeWrapper)
		r.edges[edgeTypeStr] = cachedEdgePerType
	}
	cachedEdgePerType[key] = tgedge
}

func (r *ReadyEntityKeeper) GetNode(nodeTypeStr string, key interface{}) *TGNodeWrapper {

	logger.Debug(fmt.Sprintf("[RedayEntityKeeper:getNode] Target node : %s, %s", nodeTypeStr, key))
	logger.Debug(fmt.Sprintf("[RedayEntityKeeper:getNode] Available ready node : %s", r.nodes))
	if val, ok := r.nodes[nodeTypeStr]; ok {
		logger.Debug(fmt.Sprintf("[RedayEntityKeeper:getNode] Available target type ready node : %s ", r.nodes[nodeTypeStr]))
		return val[key]
	}
	return nil
}

func (r *ReadyEntityKeeper) GetEdge(edgeTypeStr string, key interface{}) *TGEdgeWrapper {
	if val, ok := r.edges[edgeTypeStr]; ok {
		return val[key]
	}
	return nil
}

func (r *ReadyEntityKeeper) GetNodes() []*TGNodeWrapper {
	allNodes := make([]*TGNodeWrapper, 0)

	for _, nodesInType := range r.nodes {
		for _, node := range nodesInType {
			allNodes = append(allNodes, node)
		}
	}
	return allNodes
}

func (r *ReadyEntityKeeper) GetEdges() []*TGEdgeWrapper {
	allEdges := make([]*TGEdgeWrapper, 0)
	for _, edgesInType := range r.edges {
		for _, edge := range edgesInType {
			allEdges = append(allEdges, edge)
		}
	}
	return allEdges
}

func (r *ReadyEntityKeeper) Clear() {
	for nodesByType := range r.nodes {
		delete(r.nodes, nodesByType)
	}
	for edgesByType := range r.edges {
		delete(r.edges, edgesByType)
	}
}

func (r *ReadyEntityKeeper) Print() {
	logger.Debug("+++++++ ready nodes +++++++")
	logger.Debug(fmt.Sprintf("%s", r.nodes))
	logger.Debug("+++++++ reday edges +++++++")
	logger.Debug(fmt.Sprintf("%s", r.edges))
	logger.Debug("++++++++++++++++++++++++++++")
}
