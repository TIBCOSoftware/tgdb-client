package tgdb

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/logging"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/tgdb/lib/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

const (
	Query                    = "query"
	Query_QueryString        = "queryString"
	Query_TraversalCondition = "traversalCondition"
	Query_EdgeFilter         = "edgeFilter"
	Query_EndCondition       = "endCondition"
	Query_OPT_PrefetchSize   = "prefetchSize"
	Query_OPT_EdgeLimit      = "edgeLimit"
	Query_OPT_TraversalDepth = "traversalDepth"
)

var logger = logging.DefaultTGLogManager().GetLogger()

type TGDBServiceFactory struct {
}

func (this *TGDBServiceFactory) GetService(properties map[string]interface{}) *TGDBService {
	logger.SetLogLevel(types.InfoLog)
	tgdbService := TGDBService{}
	if nil != properties["url"] {
		tgdbService._url = properties["url"].(string)
	}
	if nil != properties["user"] {
		tgdbService._user = properties["user"].(string)
	}
	if nil != properties["password"] {
		tgdbService._password = properties["password"].(string)
	}
	tgdbService._keyMap = make(map[string][]string)
	tgdbService._readyEntity = NewReadyEntityKeeper()

	return &tgdbService
}

func NewTGDBServiceFactory() *TGDBServiceFactory {
	return &TGDBServiceFactory{}
}

type TGDBService struct {
	_url             string
	_user            string
	_password        string
	_connectionProps map[string]string
	_tgConnection    types.TGConnection
	_gof             types.TGGraphObjectFactory
	_gmd             types.TGGraphMetadata
	_keyMap          map[string][]string
	_readyEntity     ReadyEntityKeeper
	_mux             sync.Mutex
}

func (this *TGDBService) ensureConnection() types.TGError {
	var err types.TGError
	if nil == this._tgConnection {
		fmt.Println("\n\nXXXXXXXXXXXXXXXXXXXXXX Start connecting XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n\n")
		fmt.Println("[TGDBService::ensureConnection] Will try to connect ..........")
		fmt.Println("[TGDBService::ensureConnection] url = " + this._url)
		fmt.Println("[TGDBService::ensureConnection] user = " + this._user)
		fmt.Println("[TGDBService::ensureConnection] password = " + this._password)
		fmt.Println("[TGDBService::ensureConnection] connectionProps = ", this._connectionProps)
		this._tgConnection, err = connection.NewTGConnectionFactory().CreateConnection(this._url, this._user, this._password, this._connectionProps)
		if nil != err {
			this._tgConnection = nil
			return err
		}
		err = this._tgConnection.Connect()
		if nil != err {
			this._tgConnection = nil
			this._gof = nil
			return err
		}

		err := this.fetchMetadata()
		if nil != err {
			return err
		}

		this._gof, err = this._tgConnection.GetGraphObjectFactory()
		if nil != err {
			this._tgConnection = nil
			return err
		}
		this._tgConnection.SetExceptionListener(this)
		fmt.Println("\n\nXXXXXXXXXXXXXXXXXXXXXX End connecting normally XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n\n")
	}
	return nil
}

//-====================-//
//    Metadata API
//-====================-//

func (this *TGDBService) GetFactory() (types.TGGraphObjectFactory, error) {
	err := this.ensureConnection()
	return this._gof, err
}

func (this *TGDBService) GetMetadata() (types.TGGraphMetadata, error) {
	err := this.ensureConnection()
	return this._gmd, err
}

func (this *TGDBService) GetEdgeType(edgeTypeStr string) (types.TGEdgeType, types.TGError) {
	this.ensureConnection()
	return this._gmd.GetEdgeType(edgeTypeStr)
}

func (this *TGDBService) GetNodeKeyfields(nodeTypeStr string) []string {
	this.ensureConnection()
	nodeType, err := this._gmd.GetNodeType(nodeTypeStr)
	if nil != err {
		return nil
	}

	attributeDescriptors := nodeType.GetAttributeDescriptors()
	attrLenth := len(attributeDescriptors)

	key := make([]string, attrLenth)
	for i := 0; i < attrLenth; i++ {
		key[i] = attributeDescriptors[i].GetName()
	}
	return key
}

func (this *TGDBService) commit() types.TGError {
	fmt.Println("\n\n\n\n[TGDBDirectInsert::commit] Entering ........ ")
	this._readyEntity.Print()
	nodeWrappers := this._readyEntity.GetNodes()
	for _, nodeWrapper := range nodeWrappers {
		node := *nodeWrapper.node
		if nodeWrapper.isNew {
			fmt.Println("Insert node -> ", node.GetEntityType().GetName(), ", ", node)
			this._tgConnection.InsertEntity(node)
		} else {
			fmt.Println("Update node -> ", node.GetEntityType().GetName(), ", ", node)
			this._tgConnection.UpdateEntity(node)
		}
	}

	edgeWrappers := this._readyEntity.GetEdges()
	for _, edgeWrapper := range edgeWrappers {
		edge := *edgeWrapper.edge
		if edgeWrapper.isNew {
			fmt.Println("Insert edge -> "+edge.GetEntityType().GetName()+", ", edge)
			this._tgConnection.InsertEntity(edge)
		} else {
			if !strings.HasSuffix(edge.GetEntityType().GetName(), "_event") {
				fmt.Println("Update edge -> ", edge.GetEntityType().GetName(), ", ", edge)
				this._tgConnection.UpdateEntity(edge)
			} else {
				fmt.Println("Blocking update edge -> ", edge.GetEntityType().GetName(), ", ", edge)
			}
		}
	}

	_, err := this._tgConnection.Commit()

	// Shall handle exception not just rollback
	if nil != err {
		fmt.Println("** Error comitting : Message = " + err.GetErrorMsg())
		this._tgConnection.Rollback()
		if "Channel is Closed" == err.GetErrorMsg() {

		}
		return err
	}

	this._readyEntity.Clear()

	fmt.Println("[TGDBDirectInsert::commit] Exit ........ \n\n\n")
	return nil
}

//-====================-//
//    Upsert Graph
//-====================-//

func (this *TGDBService) UpsertGraph(graph map[string]interface{}) error {
	nodes := graph["nodes"].(map[string]interface{})
	edges := graph["edges"].(map[string]interface{})

	fmt.Println("\n\n********************************1*********************************\n\n")

	for edgeId, edge := range edges {
		edgeDetail := edge.(map[string]interface{})
		fmt.Println("Process -> edge : ", edge)
		edgeType := edgeDetail["type"]
		fromNodeId := edgeDetail["from"].(string)
		toNodeId := edgeDetail["to"].(string)
		key := edgeDetail["key"]
		fmt.Println("EdgeType: ", edgeType, ", Key: ", key, ", FromNode: ", fromNodeId, ", ToNode: ", toNodeId)
		ok, err := this.UpsertEdge(
			edgeId,
			edgeDetail,
			nodes[fromNodeId].(map[string]interface{}),
			nodes[toNodeId].(map[string]interface{}),
			true)

		if nil != err {
			fmt.Errorf("\n\n >>>>>>>>>>>>>>>>> Upsert edge failed >>>>>>>>>>>>>>>>>>>> : %s \n\n\n", err)
		}

		if !ok {
			fmt.Println("\n\n >>>>>>>>>>>>>> Edge was not upserted! >>>>>>>>>>>>>>>>>>\n\n\n")
		}

	}

	fmt.Println("\n\n********************************2*********************************\n\n")

	for _, node := range nodes {
		nodeDetail := node.(map[string]interface{})
		fmt.Println("Process -> node : ", node)
		nodeType := nodeDetail["type"]
		key := nodeDetail["key"]
		fmt.Println(nodeType, key)
		err := this.UpsertNode(nodeDetail)
		if nil != err {
			fmt.Errorf(" >>>>>>>>>>>>>>>>> Upsert node failed >>>>>>>>>>>>>>>>>>>> : \n", err)
		}

	}

	fmt.Println("\n\n********************************3*********************************\n\n")

	err := this.commit()

	fmt.Println("\n\n********************************4*********************************\n\n")

	return err
}

//-====================-//
//    Node API
//-====================-//

func (this *TGDBService) UpsertNode(nodeData map[string]interface{}) error {
	fmt.Println("\n\nXXXXXXXXXXXXXXXXXXXXXX Start UpsertNode XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n\n")

	nodeTypeStr := nodeData["type"].(string)

	this._mux.Lock()
	defer this._mux.Unlock()

	nodeKey, err := this.BuildNodeKey(nodeTypeStr, nodeData)
	fmt.Println("[TGDBService::upsertNode] node -> type = "+nodeTypeStr+", nodeKey = ", nodeKey)

	if nil != err {
		return err
	}

	var node types.TGNode
	//	var err types.TGError
	var isNew bool
	nodeWrapper := this._readyEntity.GetNode(nodeTypeStr, nodeKey)
	if nil == nodeWrapper {
		node, err = this.fetchNode(nodeTypeStr, nodeData)
		if nil != err {
			fmt.Println("/* Error to find node */")
			return err
		}

		if nil != node {
			fmt.Println("/* Node found in TGDB (%s) */", node)
			isNew = false
		} else {
			isNew = true
			node, err = this.buildNode(nodeTypeStr, nodeData)
			if nil != err {
				fmt.Println("/* Error to builde node */")
				return err
			}
			//			fmt.Println("/* So build a new node (%s) */", node)
		}
		nodeWrapper = NewTGNodeWrapper(node, isNew)
		this._readyEntity.AddNode(nodeTypeStr, nodeKey, nodeWrapper)
	} else {
		isNew = false
		node = *nodeWrapper.node
		fmt.Printf("/* Node exist in ready list (warpper = %s) */\n", nodeWrapper)
	}

	var attributesData map[string]interface{}
	if !isNew && nil != nodeData["attributes"] {
		attributesData = nodeData["attributes"].(map[string]interface{})
		nodeType, err := this._gmd.GetNodeType(nodeTypeStr)
		if nil == nodeType {
			return err
		}
		attrDescs := nodeType.GetAttributeDescriptors()
		for index := range attrDescs {
			name := attrDescs[index].GetName()
			if nil != attributesData[name] {
				attributeData := attributesData[name].(map[string]interface{})
				if nil != attributeData["value"] {
					node.SetOrCreateAttribute(name, attributeData["value"])
				}
			}
		}
	}

	if nodeWrapper.isNew {
		fmt.Println("[TGDBService::upsertNode] insert Node : type = %s, id = %s, attribuesData = %s", node.GetEntityType().GetName(), nodeKey, attributesData)
	} else {
		fmt.Println("[TGDBService::upsertNode] update Node : type = %s, id = %s, attribuesData = %s", node.GetEntityType().GetName(), nodeKey, attributesData)
	}

	return nil
}

func (this *TGDBService) GetNode(nodeType string, parameter map[string]interface{}) (types.TGEntity, types.TGError) {
	return this.fetchNode(nodeType, parameter)
}

//---------- Edge API -----------

func (this *TGDBService) UpsertEdge(
	id interface{},
	edgeData map[string]interface{},
	firstNodeData map[string]interface{},
	secondNodeData map[string]interface{},
	autoCommit bool) (bool, error) {
	fmt.Println("\n\nXXXXXXXXXXXXXXXXXXXXXX Start UpsertEdge XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n\n")

	edgeTypeStr := edgeData["type"].(string)
	attributesData := edgeData["attributes"].(map[string]interface{})
	//	allowDuplicate := false
	//	if nil != edgeData["allowDuplicate"] {
	//		allowDuplicate == edgeData["allowDuplicate"].(bool)
	//	}

	iDirection := edgeData["direction"]

	firstNodeType := firstNodeData["type"].(string)
	secondNodeType := secondNodeData["type"].(string)

	this._mux.Lock()
	defer this._mux.Unlock()

	this.ensureConnection()
	edgeType, err := this._gmd.GetEdgeType(edgeTypeStr)

	fmt.Println("[TGDBService::upsertEdge] TGDB edge type = ", edgeType, ", model edge type = ", edgeTypeStr, ", edgeID = ", id)

	if nil != err {
		return false, err
	}
	var targetEdge types.TGEdge
	edgeWrapper := this._readyEntity.GetEdge(edgeTypeStr, id)
	if nil == edgeWrapper {
		firstNodeKey, err := this.BuildNodeKey(firstNodeType, firstNodeData)
		if nil != err {
			return false, err
		}
		secondNodeKey, err := this.BuildNodeKey(secondNodeType, secondNodeData)
		if nil != err {
			return false, err
		}
		var fromNodeWrapper *TGNodeWrapper
		var toNodeWrapper *TGNodeWrapper
		// Check existence in tgdb
		var iEdgeType interface{}
		if nil != edgeType {
			iEdgeType = edgeTypeStr
		}

		var keyAttributes []string
		if nil != edgeData["keyAttributeName"] {
			keyAttributes = edgeData["keyAttributeName"].([]string)
		}

		fromNode, err := this.search(firstNodeType, firstNodeData, secondNodeType, secondNodeData, iEdgeType, keyAttributes, attributesData)
		if nil != err {
			return false, err
		}

		var toNode types.TGNode
		if nil != fromNode {
			fmt.Println("/* Edge found - will do update (%s) */", fromNode)
			edges := fromNode.GetEdges()
			fmt.Println("\nEdges : ", edges)
			for _, edge := range edges {
				fmt.Println("\nAn Edge : ", edge)
				toNode = edge.GetVertices()[1]
				fmt.Println("\nAn Edge Vertices : ", edge.GetVertices())
				if nil != toNode {
					targetEdge = edge
					fmt.Println("\nTarget Edge : " + targetEdge.GetEntityType().GetName())
					break
				}
			}

			fmt.Println("\nfrom : ", fromNode, ", to : ", toNode, ", edge : ", targetEdge)

			if nil == toNode || nil == targetEdge {
				return false, fmt.Errorf("Incomplete query (Edge type node defined in TGDB) : toNode or edge node coming back!!")
			}

			fromNodeWrapper = NewTGNodeWrapper(fromNode, false)
			toNodeWrapper = NewTGNodeWrapper(toNode, false)
			edgeWrapper = NewTGEdgeWrapper(targetEdge, false)
		} else {
			fmt.Println("/* fromNode, toNode or edge not exist. Will find fromNode, toNode saperatly */)")
			fromNodeWrapper = this._readyEntity.GetNode(firstNodeType, firstNodeKey)
			if nil == fromNodeWrapper {
				fromNode, err = this.fetchNode(firstNodeType, firstNodeData)
				if nil != err {
					return false, err
				}
				if nil != fromNode {
					fmt.Println("/* fromNode found (%s) */", fromNode)
					fromNodeWrapper = NewTGNodeWrapper(fromNode, false)
				} else {
					fromNode, err = this.buildNode(firstNodeType, firstNodeData)
					if nil != err {
						fmt.Errorf("\nError to build fromNode !!!\n")
						return false, err
					}
					if nil == fromNode {
						fmt.Errorf("\nUnable to build fromNode !!!\n")
						return false, nil
					}
					fmt.Println("/* so create fromNode (%s) */", fromNode)
					fromNodeWrapper = NewTGNodeWrapper(fromNode, true)
				}
			} else {
				fromNode = *fromNodeWrapper.node
				fmt.Println("/* fromNode exist - might be created from previous edge (%s) */", fromNode)
			}

			toNodeWrapper = this._readyEntity.GetNode(secondNodeType, secondNodeKey)
			if nil == toNodeWrapper {
				toNode, err = this.fetchNode(secondNodeType, secondNodeData)
				if nil != err {
					return false, err
				}
				if nil != toNode {
					fmt.Println("/* toNode found (%s) */", toNode)
					toNodeWrapper = NewTGNodeWrapper(toNode, false)
				} else {
					toNode, err = this.buildNode(secondNodeType, secondNodeData)
					if nil != err {
						fmt.Errorf("\nError to build toNode !!!\n")
						return false, err
					}
					if nil == toNode {
						fmt.Errorf("\nUnable to build toNode !!!\n")
						return false, nil
					}
					fmt.Println("/* so create toNode (%s) */", toNode)
					toNodeWrapper = NewTGNodeWrapper(toNode, true)
				}
			} else {
				toNode = *toNodeWrapper.node
				fmt.Println("/* toNode exist - might be created from previous edge (%s) */", toNode)
			}

			if nil != edgeType {
				targetEdge, err = this._gof.CreateEdgeWithEdgeType(fromNode, toNode, edgeType)
			} else {
				direction := types.DirectionTypeUnDirected
				if nil != iDirection {
					switch iDirection.(int) {
					case 0:
						direction = types.DirectionTypeUnDirected
						break
					case 1:
						direction = types.DirectionTypeDirected
						break
					case 2:
						direction = types.DirectionTypeBiDirectional
						break
					}
				}
				targetEdge, err = this._gof.CreateEdgeWithDirection(fromNode, toNode, direction)
			}

			if nil != err {
				return false, err
			}
			edgeWrapper = NewTGEdgeWrapper(targetEdge, true)
		}
		this._readyEntity.AddNode(firstNodeType, firstNodeKey, fromNodeWrapper)
		this._readyEntity.AddNode(secondNodeType, secondNodeKey, toNodeWrapper)
		this._readyEntity.AddEdge(edgeTypeStr, id, edgeWrapper)
	} else {
		targetEdge = *edgeWrapper.edge
	}

	if nil != attributesData {
		if nil != edgeType {
			attrDescs := edgeType.GetAttributeDescriptors()
			for index := range attrDescs {
				name := attrDescs[index].GetName()
				if nil != attributesData[name] {
					attributeData := attributesData[name].(map[string]interface{})
					if nil != attributeData["value"] {
						targetEdge.SetOrCreateAttribute(name, attributeData["value"])
					}
				}
			}
		} else {
			for name := range attributesData {
				attributeData := attributesData[name].(map[string]interface{})
				if nil != attributeData["value"] {
					targetEdge.SetOrCreateAttribute(name, attributeData["value"])
				}
			}
		}
	}

	fmt.Println("\n\nXXXXXXXXXXXXXXXXXXXXXX End UpsertEdge XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n\n")

	return true, nil
}

func (this *TGDBService) Query(para map[string]interface{}) (types.TGResultSet, types.TGError) {
	var queryString string
	var edgeFilter string
	var traversalCondition string
	var endCondition string

	if nil == para["query"] {
		return nil, &types.TGDBError{ErrorMsg: "No query parameter defined!"}
	}
	queryObj := para["query"].(map[string]interface{})
	if nil != queryObj {
		if nil != queryObj[Query_QueryString] {
			queryString = queryObj[Query_QueryString].(string)
		}

		if nil != queryObj[Query_EdgeFilter] {
			edgeFilter = queryObj[Query_EdgeFilter].(string)
		}

		if nil != queryObj[Query_TraversalCondition] {
			traversalCondition = queryObj[Query_TraversalCondition].(string)
		}

		if nil != queryObj[Query_EndCondition] {
			endCondition = queryObj[Query_EndCondition].(string)
		}
	}

	option := this.buildQueryOption(para)

	fmt.Println("* queryString           = ", queryString)
	fmt.Println("* edgeFilter            = ", edgeFilter)
	fmt.Println("* traversalCondition    = ", traversalCondition)
	fmt.Println("* endCondition          = ", endCondition)
	fmt.Println("* Option.prefetchSize   = ", option.GetPreFetchSize())
	fmt.Println("* Option.edgeLimit      = ", option.GetEdgeLimit())
	fmt.Println("* Option.traversalDepth = ", option.GetTraversalDepth())

	this.ensureConnection()

	resultSet, err := this._tgConnection.ExecuteQueryWithFilter(
		queryString, edgeFilter, traversalCondition, endCondition, option)

	return resultSet, err
}

func (this *TGDBService) Destroy() {
	if nil != this._tgConnection {
		this._tgConnection.Disconnect()
	}
}

func (this *TGDBService) BuildNodeKey(nodeTypeStr string, nodeData map[string]interface{}) (string, types.TGError) {

	attribuesData := nodeData["attributes"].(map[string]interface{})
	//fmt.Println("Type = ", nodeTypeStr, ", nodeData = ", nodeData, ", attribuesData = ", attribuesData)
	fmt.Println("Type = ", strings.Trim(nodeTypeStr, " "), ", nodeData = ", nodeData, ", attribuesData = ", attribuesData)
	this.ensureConnection()
	nodeType, err := this._gmd.GetNodeType("houseMemberType")
	if nil != err {
		return "", err
	}
	fmt.Printf("======> NodeType: '%+v'\n", nodeType)

	pKeyAttributeDescriptors := nodeType.GetPKeyAttributeDescriptors()
	attrLength := len(pKeyAttributeDescriptors)
	key := make([]interface{}, attrLength)
	for i := 0; i < attrLength; i++ {
		name := pKeyAttributeDescriptors[i].GetName()
		if nil != attribuesData[name] {
			attributeData := attribuesData[name].(map[string]interface{})
			key[i] = attributeData["value"]
		} else {
			key[i] = nil
		}
	}

	return model.Hash(key), nil
}

//-====================-//
//       private
//-====================-//

func (this *TGDBService) fetchNode(fromType string, fromNode map[string]interface{}) (types.TGNode, types.TGError) {
	return this.search(fromType, fromNode, "", nil, nil, nil, nil)
}

func extractAttribute(entity map[string]interface{}) map[string]interface{} {
	attributes := entity["attributes"]
	if nil == attributes {
		attributes = make(map[string]interface{})
	}
	return attributes.(map[string]interface{})
}

func (this *TGDBService) search(
	fromType string,
	fromNode map[string]interface{},
	toType string,
	toNode map[string]interface{},
	edgeType interface{},
	edgeKeyAttributes []string,
	edgeKey map[string]interface{}) (types.TGNode, types.TGError) {

	fromAttributes := extractAttribute(fromNode)
	toAttributes := extractAttribute(toNode)

	fmt.Println("fromType = ", fromType, ", fromAttributes = ", fromAttributes)
	fmt.Println("toType = ", toType, ", toAttributes = ", toAttributes)
	fmt.Println("edgeType = ", edgeType, ", edgeKeyAttributes = ", edgeKeyAttributes, ", edgeKey = ", edgeKey)

	this.ensureConnection()
	nodeType, err := this._gmd.GetNodeType(fromType)
	if nil != err {
		return nil, err
	}

	fromCondition := ""
	//var fromNodeCondition strings.Builder
	attrDescs := nodeType.GetPKeyAttributeDescriptors()
	for index := range attrDescs {
		name := attrDescs[index].GetName()
		//fromNodeCondition.WriteString("and ")
		fromCondition += "and "
		//fromNodeCondition.WriteString(name)
		fromCondition += name
		//fromNodeCondition.WriteString(" = '")
		fromCondition += " = '"
		fromKeyAttributeData := fromAttributes[name].(map[string]interface{})
		//fromNodeCondition.WriteString(fromKeyAttributeData["value"].(string))
		fromCondition += fromKeyAttributeData["value"].(string)
		//fromNodeCondition.WriteString("'")
		fromCondition += "'"
	}

	edgeFilterCondition := ""
	//var edgeCondition strings.Builder
	if nil != edgeKeyAttributes {
		for index := range edgeKeyAttributes {
			//edgeCondition.WriteString("and @edge.")
			edgeFilterCondition += "and @edge."
			//edgeCondition.WriteString(edgeKeyAttributes[index])
			edgeFilterCondition += edgeKeyAttributes[index]
			//edgeCondition.WriteString(" = '")
			edgeFilterCondition += " = '"
			edgeKeyAttributeData := edgeKey[edgeKeyAttributes[index]].(map[string]interface{})
			//edgeCondition.WriteString(edgeKeyAttributeData["value"].(string))
			edgeFilterCondition += edgeKeyAttributeData["value"].(string)
			//edgeCondition.WriteString("'")
			edgeFilterCondition += "'"
		}
	}

	toCondition := ""
	//var toNodeCondition strings.Builder
	if "" != toType {
		toNodeType, err := this._gmd.GetNodeType(toType)
		if nil != err {
			return nil, err
		}
		attrDescs := toNodeType.GetPKeyAttributeDescriptors()
		for index := range attrDescs {
			name := attrDescs[index].GetName()
			//toNodeCondition.WriteString("and @tonode.")
			toCondition += "and @tonode."
			//toNodeCondition.WriteString(name)
			toCondition += name
			//toNodeCondition.WriteString(" = '")
			toCondition += " = '"
			toKeyAttributeData := toAttributes[name].(map[string]interface{})
			//toNodeCondition.WriteString(toKeyAttributeData["value"].(string))
			toCondition += toKeyAttributeData["value"].(string)
			//toNodeCondition.WriteString("'")
			toCondition += "'"
		}
	}

	para := make(map[string]interface{})
	query := make(map[string]interface{})
	//query["queryString"] = fmt.Sprintf("@nodetype = '%s' %s;", fromType, fromNodeCondition.String())
	query["queryString"] = fmt.Sprintf("@nodetype = '%s' %s;", fromType, fromCondition)
	if nil != edgeType {
		//query["traversalCondition"] = fmt.Sprintf("@edgetype = '%s' %s and @tonodetype = '%s' %s and @isfromedge = 1 and @degree = 1;", edgeType, edgeCondition.String(), toType, toNodeCondition.String())
		query["traversalCondition"] = fmt.Sprintf("@edgetype = '%s' %s and @tonodetype = '%s' %s and @isfromedge = 1 and @degree = 1;", edgeType, edgeFilterCondition, toType, toCondition)
	} else {
		/* need to search by edge index */
		//query["traversalCondition"] = fmt.Sprintf("@tonodetype = '%s' %s and @isfromedge = 1 and @degree = 1;", toType, toNodeCondition.String())
		query["traversalCondition"] = fmt.Sprintf("@tonodetype = '%s' %s and @isfromedge = 1 and @degree = 1;", toType, toCondition)
	}

	fmt.Printf("[TGDBService::searchNode] query string = %s \n", query["queryString"])
	fmt.Printf("[TGDBService::searchNode] traversalCondition = %s \n", query["traversalCondition"])

	para[Query] = query
	para[Query_OPT_PrefetchSize] = 500
	para[Query_OPT_TraversalDepth] = 1
	para[Query_OPT_EdgeLimit] = 500

	var startingNode types.TGNode
	resultSet, err := this.Query(para)
	if nil == err {
		if nil != resultSet {
			result := resultSet.GetAt(0)
			//			fmt.Println("result : ", result)
			if nil == result {
				return nil, nil
			}
			startingNode = result.(types.TGNode)
		}
	} else {
		fmt.Println("Message : " + err.GetErrorMsg())
		if "Channel is Closed" == err.GetErrorMsg() {
			return nil, err
		}
	}

	return startingNode, err
}

func (this *TGDBService) populateAttributes(node types.TGNode, attributesData map[string]interface{}) types.TGNode {
	for _, attrDesc := range node.GetEntityType().GetAttributeDescriptors() {
		name := attrDesc.GetName()
		if nil != attributesData[name] {
			attributeData := attributesData[name].(map[string]interface{})
			node.SetOrCreateAttribute(name, attributeData["value"])
		}
	}

	return node
}

func (this *TGDBService) buildNode(nodeTypeStr string, nodeData map[string]interface{}) (types.TGNode, types.TGError) {

	nodeType, err := this._gmd.GetNodeType(nodeTypeStr)
	if nil == nodeType {
		return nil, err
	}

	node, _ := this._gof.CreateNodeInGraph(nodeType)

	if nil != nodeData["attributes"] {
		this.populateAttributes(node, nodeData["attributes"].(map[string]interface{}))
	}

	return node, nil
}

func (this *TGDBService) fetchMetadata() types.TGError {
	gmd, err := this._tgConnection.GetGraphMetadata(true)

	if nil != err {
		return err
	}

	this._gmd = gmd

	for value := range this._keyMap {
		delete(this._keyMap, value)
	}

	nodeTypes, err2 := this._gmd.GetNodeTypes()

	if nil != err2 {
		return err2
	}

	for _, nodeType := range nodeTypes {
		keys := make([]string, 0)
		attrDescs := nodeType.GetPKeyAttributeDescriptors()
		for index := range attrDescs {
			keys = append(keys, attrDescs[index].GetName())
		}
		this._keyMap[nodeType.GetName()] = keys
	}
	return nil
}

func (this *TGDBService) buildQueryOption(para map[string]interface{}) types.TGQueryOption {
	option := query.NewQueryOption()

	fmt.Println("Parameter : ", para)

	if nil == para[Query_OPT_PrefetchSize] || 0 == para[Query_OPT_PrefetchSize].(int) {
		option.SetPreFetchSize(500)
	} else {
		option.SetPreFetchSize(para[Query_OPT_PrefetchSize].(int))
	}

	if nil == para[Query_OPT_TraversalDepth] || 0 == para[Query_OPT_TraversalDepth].(int) {
		option.SetTraversalDepth(5)
	} else {
		option.SetTraversalDepth(para[Query_OPT_TraversalDepth].(int))
	}

	if nil == para[Query_OPT_EdgeLimit] || 0 == para[Query_OPT_EdgeLimit].(int) {
		option.SetEdgeLimit(100)
	} else {
		option.SetEdgeLimit(para[Query_OPT_EdgeLimit].(int))
	}

	return option
}

func (this *TGDBService) OnException(ex types.TGError) {
	fmt.Println("[TGDBService::onException] Exception Happens : " + ex.GetErrorMsg())
	if nil != this._tgConnection {
		this._tgConnection.Disconnect()
	}
	this._tgConnection = nil
}
