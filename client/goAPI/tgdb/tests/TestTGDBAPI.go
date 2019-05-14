// TestTGDBAPI
package main

import (
	"encoding/json"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/tgdb/lib/dbservice/tgdb"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

type EdgeStruct struct {
	Name         string       `json:"edges,omitempty"`
	Attributes   []AttrStruct `json:"attributes,omitempty"`
	FromNodeName string       `json:"from,omitempty"`
	Keys         []string     `json:"key,omitempty"`
	KeyAttrName  string       `json:"keyAttributeName,omitempty"`
	NodeType     string       `json:"type,omitempty"`
	ToNodeName   string       `json:"to,omitempty"`
}

type AttrStruct struct {
	AttrName  string `json:"name,omitempty"`
	AttrType  string `json:"type,omitempty"`
	AttrValue string `json:"value,omitempty"`
}

type NodeStruct struct {
	Name        string       `json:"edges,omitempty"`
	Attributes  []AttrStruct `json:"attributes,omitempty"`
	Keys        []string     `json:"key,omitempty"`
	KeyAttrName string       `json:"keyAttributeName,omitempty"`
	NodeType    string       `json:"type,omitempty"`
}

type IntermediateKeyMap struct {
	KeyMapName   string   `json:"keyMapName,omitempty"`
	KeyMapValues []string `json:"keyMapValues,omitempty"`
}

type IntermediateNode struct {
	KeyMap IntermediateKeyMap `json:"keyMap,omitempty"`
	Types  []string           `json:"types,omitempty"`
}

type IntermediateEdge struct {
	KeyMap IntermediateKeyMap `json:"keyMap,omitempty"`
	Types  []string           `json:"types,omitempty"`
}

type ModelStruct struct {
	ModelId string             `json:"modelId,omitempty"`
	Edges   []IntermediateEdge `json:"attributes,omitempty"`
	Nodes   []IntermediateNode `json:"key,omitempty"`
}

type IntGraphStruct struct {
	Edges   []EdgeStruct  `json:"edges,omitempty"`
	Id      string        `json:"id"`
	Model   []ModelStruct `json:"model,omitempty"`
	ModelId string        `json:"modelId"`
	Nodes   []NodeStruct  `json:"nodes,omitempty"`
}

type GraphStruct struct {
	Tag string         `json:"graph"`
	Obj IntGraphStruct `json:"obj"`
}

func main() {
	properties := make(map[string]interface{})
	properties["url"] = "tcp://localhost:8222"
	properties["user"] = "scott"
	properties["password"] = "scott"
	properties["password"] = "scott"
	fmt.Println(properties)

	tgdbSrv := tgdb.NewTGDBServiceFactory().GetService(properties)

	upsertGraph(tgdbSrv)

	query(tgdbSrv)
}

func upsertGraph(tgdbSrv *tgdb.TGDBService) {
	//a := loadJSONData()
	//fmt.Printf("JSON Data Loaded: '%+v'\n", a)
	graph := getGraph()
	if nil != graph {
		tgdbSrv.UpsertGraph(graph.(map[string]interface{}))
	}
}

func query(tgdbSrv *tgdb.TGDBService) {
	pamarms := make(map[string]interface{})
	query := make(map[string]interface{})
	query["queryString"] = "@nodetype = 'houseMemberType' and memberName = 'Carlo Bonaparte';"
	query["traversalCondition"] = "@edgetype = 'relation'  and @tonodetype = 'houseMemberType' and @tonode.memberName = 'Letizia Ramolino' and @isfromedge = 1 and @degree = 1;"

	pamarms[tgdb.Query] = query
	pamarms[tgdb.Query_OPT_PrefetchSize] = 500
	pamarms[tgdb.Query_OPT_TraversalDepth] = 1
	pamarms[tgdb.Query_OPT_EdgeLimit] = 500

	var startingNode types.TGNode
	var targetEdge types.TGEdge
	var toNode types.TGNode
	resultSet, err := tgdbSrv.Query(pamarms)
	if nil == err {
		if nil != resultSet {
			result := resultSet.GetAt(0)
			fmt.Println("result : ", result)
			startingNode = result.(types.TGNode)

			fmt.Println("\n1. Starting Node found - ", startingNode)
			edges := startingNode.GetEdges()
			fmt.Println("\n2. Edges associated with node - ", edges)
			for _, edge := range edges {
				fmt.Println("\n3. An Edge : ", edge)
				toNode = edge.GetVertices()[1]
				fmt.Println("\n4. Got a Vertice from Edge : ", edge.GetVertices())
				if nil != toNode {
					targetEdge = edge
					fmt.Println("\n5 .Target Edge : " + targetEdge.GetEntityType().GetName())
					break
				}
			}

			fmt.Println("\n*************************** Query Result *********************************")
			fmt.Println("- startingNode = ", startingNode)
			fmt.Println("- toNode = ", toNode)
			fmt.Println("- edge = ", targetEdge)

			if nil == toNode || nil == targetEdge {
				fmt.Println("\nIncomplete query result (Edge type node defined in TGDB) : toNode or edge node coming back!!")
			}
			fmt.Println("\n**************************************************************************")
		} else {
			fmt.Println("\nNo Result !!!!!!!!!!!!!!!!!")
		}
	} else {
		fmt.Println("\nMessage : " + err.GetErrorMsg())
	}

}

func loadJSONData() interface{} {
	graphString := "{\"graph\":{\"edges\":{\"relation_d41d8cd98f00b204e9800998ecf8427e_5b0533f7c787ece47dd36f7008b04975_dc4c735265e66942e6478ef1d067072a\":{\"attributes\":{\"relType\":{\"name\":\"relType\",\"type\":\"String\",\"value\":\"spouse\"}},\"from\":\"houseMemberType_5b0533f7c787ece47dd36f7008b04975\",\"key\":[],\"keyAttributeName\":null,\"to\":\"houseMemberType_dc4c735265e66942e6478ef1d067072a\",\"type\":\"relation\"}},\"id\":\"hierarchy\",\"model\":{\"edges\":{\"keyMap\":{\"relation\":null},\"types\":[\"relation\"]},\"nodes\":{\"keyMap\":{\"houseMemberType\":[\"memberName\"]},\"types\":[\"houseMemberType\"]}},\"modelId\":\"hierarchy\",\"nodes\":{\"houseMemberType_5b0533f7c787ece47dd36f7008b04975\":{\"attributes\":{\"memberName\":{\"name\":\"memberName\",\"type\":\"String\",\"value\":\"Carlo Bonaparte\"}},\"key\":[\"Carlo Bonaparte\"],\"keyAttributeName\":[\"memberName\"],\"type\":\"houseMemberType\"},\"houseMemberType_dc4c735265e66942e6478ef1d067072a\":{\"attributes\":{\"memberName\":{\"name\":\"memberName\",\"type\":\"String\",\"value\":\"Letizia Ramolino\"}},\"key\":[\"Letizia Ramolino\"],\"keyAttributeName\":[\"memberName\"],\"type\":\"houseMemberType\"}}}}"
	var graph = GraphStruct{}
	err := json.Unmarshal([]byte(graphString), &graph)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return graph.Obj
}

func getGraph() interface{} {
	graphString := "{\"graph\":{\"edges\":{\"relation_d41d8cd98f00b204e9800998ecf8427e_5b0533f7c787ece47dd36f7008b04975_dc4c735265e66942e6478ef1d067072a\":{\"attributes\":{\"relType\":{\"name\":\"relType\",\"type\":\"String\",\"value\":\"spouse\"}},\"from\":\"houseMemberType_5b0533f7c787ece47dd36f7008b04975\",\"key\":[],\"keyAttributeName\":null,\"to\":\"houseMemberType_dc4c735265e66942e6478ef1d067072a\",\"type\":\"relation\"}},\"id\":\"hierarchy\",\"model\":{\"edges\":{\"keyMap\":{\"relation\":null},\"types\":[\"relation\"]},\"nodes\":{\"keyMap\":{\"houseMemberType\":[\"memberName\"]},\"types\":[\"houseMemberType\"]}},\"modelId\":\"hierarchy\",\"nodes\":{\"houseMemberType_5b0533f7c787ece47dd36f7008b04975\":{\"attributes\":{\"memberName\":{\"name\":\"memberName\",\"type\":\"String\",\"value\":\"Carlo Bonaparte\"}},\"key\":[\"Carlo Bonaparte\"],\"keyAttributeName\":[\"memberName\"],\"type\":\"houseMemberType\"},\"houseMemberType_dc4c735265e66942e6478ef1d067072a\":{\"attributes\":{\"memberName\":{\"name\":\"memberName\",\"type\":\"String\",\"value\":\"Letizia Ramolino\"}},\"key\":[\"Letizia Ramolino\"],\"keyAttributeName\":[\"memberName\"],\"type\":\"houseMemberType\"}}}}"
	var rootObject interface{}
	err := json.Unmarshal([]byte(graphString), &rootObject)
	if nil != err {
		fmt.Println(err)
		return nil
	}
	return rootObject.(map[string]interface{})["graph"]
}
