/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"sync"
)

type GraphManager struct {
	graphs map[string]Graph
}

var (
	instance *GraphManager
	once     sync.Once
	mux      sync.Mutex
)

func GetGraphManager() *GraphManager {
	once.Do(func() {
		instance = &GraphManager{}
	})
	return instance
}

func (this *GraphManager) GetGraph(
	modelId string,
	graphId string) Graph {

	graph := this.graphs[graphId]
	if nil == graph {
		mux.Lock()
		defer mux.Unlock()
		graph = this.graphs[graphId]
		if nil == graph {
			graph = CreateUndefinedGraph(modelId, graphId)
			this.graphs[graphId] = graph
		}
	}

	return graph
}

func CreateUndefinedGraph(modelId string, graphId string) Graph {
	graph := NewGraphImpl(modelId, graphId)
	return graph
}
