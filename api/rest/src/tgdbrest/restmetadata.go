/*
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
 * File name: restmetadata.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: restmetadata.go 4610 2020-10-30 17:24:33Z nimish $
 */

package tgdbrest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"tgdb"
	"tgdb/impl"
	"tgdbrest/spotfire"
)

var logger = impl.DefaultTGLogManager().GetLogger()

const (
	url = "tcp://scott@localhost:8222"
	password         = "scott"
)

func RESTTransaction (conn tgdb.TGConnection, w http.ResponseWriter, r *http.Request, headers map[string] string, body map[string] interface{}) {

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}
	gmd, err := conn.GetGraphMetadata(true);
	nodeDetail, ok := body["CreateNode"]
	if ok {
		nodeDetailConcrete, ok := nodeDetail.(map[string]interface{})
		if ok {
			typeNameAny := nodeDetailConcrete["Name"]
			var typeName string
			if typeNameAny == nil {
				typeName = ""
			} else {
				typeName = typeNameAny.(string)
			}

			nodeType, err := gmd.GetNodeType(typeName)
			if err != nil {
				handleRESTError(err.Error(), w)
				return
			}

			node, err := gof.CreateNodeInGraph(nodeType)
			if err != nil {
				handleRESTError(err.Error(), w)
				return
			}

			attrib := nodeDetailConcrete["Attributes"]
			if attrib != nil {
				attribConcrete := attrib.([]interface{})
				//fmt.Printf("TypeName %T attributes %T", typeName, attrib)
				for i := 0; i < len(attribConcrete); i++ {
					currentAttribute := attribConcrete[i].(map[string]interface{})
					name := currentAttribute["Name"]
					value := currentAttribute["Value"]
					err = node.SetOrCreateAttribute(name.(string), value)
					if err != nil {
						handleRESTError(err.Error(), w)
						return
					}
				}
			}
			err = conn.InsertEntity(node)
			if err != nil {
				handleRESTError(err.Error(), w)
				return
			}

			_, err = conn.Commit()
			if err != nil {
				handleRESTError(err.Error(), w)
				return
			}

			var transactionCreateNodeBody TGDBRestTransactionCreateNodeBody
			transactionCreateNodeBody.Id = node.GetVirtualId()

			b, err1 := json.MarshalIndent(transactionCreateNodeBody, "", "\t")
			if err1 != nil {
				logger.Error("error: " + err1.Error())
				handleRESTError(err1.Error(), w)
				return
			}

			w.WriteHeader(200)
			w.Write([]byte(string(b)))
		}
		return
	}
	///*
	deleteNodeDetail, ok := body["DeleteNode"]
	if ok {
		deleteNodeDetailConcrete, ok := deleteNodeDetail.(map[string]interface{})
		if ok {
			typeNameAny := deleteNodeDetailConcrete["Name"]

			var typeName string
			if typeNameAny == nil {
				handleRESTError("NodeType Name is not specified in the request.", w)
				return
			} else {
				typeName = typeNameAny.(string)
			}

			metadata := gmd.(*impl.GraphMetadata)
			compositeKey := impl.NewCompositeKey(metadata, typeName)

			attrib := deleteNodeDetailConcrete["Attributes"]
			if attrib == nil {
				handleRESTError("Attributes are not specified in the request message.", w)
				return
			}
			attribConcrete := attrib.([]interface{})

			for i := 0; i < len(attribConcrete); i++ {
				currentAttribute := attribConcrete[i].(map[string]interface{})
				name := currentAttribute["Name"]
				value := currentAttribute["Value"]
				compositeKey.SetKeyName(typeName)
				err = compositeKey.SetOrCreateAttribute(name.(string), value)
				if err != nil {
					handleRESTError(err.Error(), w)
					return
				}
			}

			tgEntity, tgError := conn.GetEntity(compositeKey, nil)
			if tgError != nil {
				handleRESTError(tgError.Error(), w)
				return
			}
			err = conn.DeleteEntity(tgEntity)
			if err != nil {
				handleRESTError(err.Error(), w)
				return
			}
			_, err = conn.Commit()
			if err != nil {
				handleRESTError(err.Error(), w)
				return
			}

			var transactionCreateNodeBody TGDBRestTransactionCreateNodeBody
			transactionCreateNodeBody.Id = tgEntity.GetVirtualId()

			b, err1 := json.MarshalIndent(transactionCreateNodeBody, "", "\t")
			if err1 != nil {
				logger.Error("error:" + err1.Error())
				handleRESTError(err1.Error(), w)
				return
			}

			w.WriteHeader(200)
			w.Write([]byte(string(b)))

		}
		return
	}

	edgeDetail, ok := body["CreateEdge"]
	if ok {
		edgeDetailConcrete, ok := edgeDetail.(map[string]interface{})
		if ok {
			edgeTypeName, ok := edgeDetailConcrete["Name"]
			if !ok {
				handleRESTError("EdgeType Name is not specified in the request.", w)
				return
			}
			if ok {
				var tgEntityFromNode tgdb.TGEntity
				var tgError tgdb.TGError
				fromNodeDetail, ok := edgeDetailConcrete["FromNode"]
				if !ok {
					handleRESTError("FromNode is not specified in the request.", w)
					return
				}
				if ok {
					fromNodeDetailConcrete, ok := fromNodeDetail.(map[string]interface{})
					if ok {
						fromNodeTypeName, ok := fromNodeDetailConcrete["Name"]
						if !ok {
							handleRESTError("FromNode: NodeType Name is not specified in the request.", w)
							return
						}

						metadata := gmd.(*impl.GraphMetadata)
						compositeKeyForFromNode := impl.NewCompositeKey(metadata, fromNodeTypeName.(string))

						attrib := fromNodeDetailConcrete["Attributes"]
						if attrib == nil {
							handleRESTError("FromNode: Attributes are not specified in the request message.", w)
							return
						}
						attribConcrete := attrib.([]interface{})

						for i := 0; i < len(attribConcrete); i++ {
							currentAttribute := attribConcrete[i].(map[string]interface{})
							name := currentAttribute["Name"]
							value := currentAttribute["Value"]
							compositeKeyForFromNode.SetKeyName(fromNodeTypeName.(string))
							err = compositeKeyForFromNode.SetOrCreateAttribute(name.(string), value)
							if err != nil {
								handleRESTError(err.Error(), w)
								return
							}
						}

						tgEntityFromNode, tgError = conn.GetEntity(compositeKeyForFromNode, nil)
						if tgError != nil {
							handleRESTError(tgError.Error(), w)
							return
						}
					}

				}

				var tgEntityToNode tgdb.TGEntity
				toNodeDetail, ok := edgeDetailConcrete["ToNode"]
				if !ok {
					handleRESTError("ToNode is not specified in the request.", w)
					return
				}
				if ok {
					toNodeDetailConcrete, ok := toNodeDetail.(map[string]interface{})
					if ok {

						toNodeTypeName, ok := toNodeDetailConcrete["Name"]
						if !ok {
							handleRESTError("ToNode: NodeType Name is not specified in the request.", w)
							return
						}
						metadata := gmd.(*impl.GraphMetadata)
						compositeKeyForToNode := impl.NewCompositeKey(metadata, toNodeTypeName.(string))

						attrib := toNodeDetailConcrete["Attributes"]
						if attrib == nil {
							handleRESTError("ToNode: Attributes are not specified in the request message.", w)
							return
						}
						attribConcrete := attrib.([]interface{})

						for i := 0; i < len(attribConcrete); i++ {
							currentAttribute := attribConcrete[i].(map[string]interface{})
							name := currentAttribute["Name"]
							value := currentAttribute["Value"]
							compositeKeyForToNode.SetKeyName(toNodeTypeName.(string))
							err = compositeKeyForToNode.SetOrCreateAttribute(name.(string), value)
							if err != nil {
								handleRESTError(err.Error(), w)
								return
							}
						}

						tgEntityToNode, tgError = conn.GetEntity(compositeKeyForToNode, nil)
						if tgError != nil {
							handleRESTError(tgError.Error(), w)
							return
						}
					}
				}

				metadata := gmd.(*impl.GraphMetadata)
				edgeType, tgError := metadata.GetEdgeType(edgeTypeName.(string))
				if tgError != nil {
					handleRESTError(tgError.Error(), w)
					return
				}

				edge, tgError := gof.CreateEdgeWithEdgeType(tgEntityFromNode.(tgdb.TGNode), tgEntityToNode.(tgdb.TGNode), edgeType)
				if tgError != nil {
					handleRESTError(tgError.Error(), w)
					return
				}
				edgeAttributes, ok := edgeDetailConcrete["Attributes"]
				if !ok {
					handleRESTError("Edge Attributes are not specified in the request.", w)
					return
				}
				if ok {
					attribConcrete := edgeAttributes.([]interface{})
					for i := 0; i < len(attribConcrete); i++ {
						currentAttribute := attribConcrete[i].(map[string]interface{})
						name := currentAttribute["Name"]
						value := currentAttribute["Value"]
						err = edge.SetOrCreateAttribute(name.(string), value)
						if err != nil {
							handleRESTError(err.Error(), w)
							return
						}
					}
				}

				err = conn.InsertEntity(edge)
				if err != nil {
					handleRESTError(err.Error(), w)
					return
				}
				_, tgError = conn.Commit()
				if tgError != nil {
					handleRESTError(tgError.Error(), w)
					return
				}

				var transactionCreateNodeBody TGDBRestTransactionCreateNodeBody
				transactionCreateNodeBody.Id = edge.GetVirtualId()

				b, err1 := json.MarshalIndent(transactionCreateNodeBody, "", "\t")
				if err1 != nil {
					logger.Error("error:" + err1.Error())
					handleRESTError(err1.Error(), w)
					return
				}

				w.WriteHeader(200)
				w.Write([]byte(string(b)))
				return
			}
		}
	}
}

func handleRESTError (errorMsg string,  w http.ResponseWriter) {
	restError := TGDBRESTError{errorMsg}
	b, err1 := json.MarshalIndent(restError, "", "\t")
	if err1 != nil {
		logger.Error("error: " + err1.Error())
		w.WriteHeader(200)
		w.Write([]byte(err1.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func MetadataQueryForNodeTypes (conn tgdb.TGConnection, w http.ResponseWriter, r *http.Request, headers map[string] string, body map[string] string) {
	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}
	if gof == nil {
		return
	}

	//if prefetchMetaData {
	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}

	var nodeTypeName string

	if body != nil && len(body) > 0 {
		nodeTypeName = body["Name"]
	}

	if (len(nodeTypeName) < 1) {
		nodeTypes, err := gmd.GetNodeTypes()
		if err != nil {
			logger.Error("error:" + err.Error())
			handleRESTError(err.Error(), w)
			return
		}
		b, error := json.MarshalIndent(nodeTypes, "", "\t")
		if error != nil {
			logger.Error("error: " + error.Error())
			handleRESTError(error.Error(), w)
			return
		}

		fmt.Fprintf(w, string(b))
	} else {
		var nodeType tgdb.TGNodeType
		var tgError tgdb.TGError
		nodeType, tgError = gmd.GetNodeType(nodeTypeName)
		if tgError != nil {
			logger.Error("error: " + tgError.Error())
			handleRESTError(tgError.Error(), w)
			return
		}
		var result[] byte
		var er error
		result, er = json.MarshalIndent(nodeType, "", "\t")
		if er != nil {
			logger.Error("error: " + er.Error())
			handleRESTError(er.Error(), w)
			return
		}
		fmt.Fprintf(w, string(result))
	}
}

func MetadataQueryForEdgeTypes (conn tgdb.TGConnection, w http.ResponseWriter, r *http.Request, headers map[string]string, body map[string]string) {
	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}
	if gof == nil {
		return
	}


	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}

	var edgeTypeName string
	if body != nil && len(body) > 0 {
		edgeTypeName = body["Name"]
	}

	if (len(edgeTypeName) < 1) {
		edgeTypes, err := gmd.GetEdgeTypes()

		b, error := json.MarshalIndent(edgeTypes, "", "\t")
		if error != nil {
			logger.Error("error: " + err.Error())
			handleRESTError(error.Error(), w)
			return
		}

		fmt.Fprintf(w, string(b))
	} else {
		var edgeType tgdb.TGEdgeType
		var tgError tgdb.TGError
		edgeType, tgError = gmd.GetEdgeType(edgeTypeName)
		if tgError != nil {
			logger.Error("error: " + tgError.Error())
			handleRESTError(tgError.Error(), w)
			return
		}
		var result[] byte
		var er error
		result, er = json.MarshalIndent(edgeType, "", "\t")
		if er != nil {
			logger.Error("error: " + er.Error())
			handleRESTError(er.Error(), w)
			return
		}
		fmt.Fprintf(w, string(result))
	}

}

func MetadataQueryForAttributeDescriptors (conn tgdb.TGConnection, w http.ResponseWriter, r *http.Request, headers map[string]string, body map[string] string) {
	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}
	if gof == nil {
		return
	}

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}

	var attributeName string
	if body != nil && len(body) > 0 {
		attributeName = body["Name"]
	}

	if (len(attributeName) < 1) {
		attributeDescriptors, err := gmd.GetAttributeDescriptors()

		b, error := json.MarshalIndent(attributeDescriptors, "", "\t")
		if error != nil {
			logger.Error("error: " + err.Error())
			handleRESTError(error.Error(), w)
			return
		}
		fmt.Fprintf(w, string(b))
	} else {
		var attributeDescriptor tgdb.TGAttributeDescriptor
		var tgError tgdb.TGError
		attributeDescriptor, tgError = gmd.GetAttributeDescriptor(attributeName)
		if tgError != nil {
			logger.Error("error: " + tgError.Error())
			handleRESTError(tgError.Error(), w)
			return
		}
		var result[] byte
		var er error
		result, er = json.MarshalIndent(attributeDescriptor, "", "\t")
		if er != nil {
			logger.Error("error:" + er.Error())
			handleRESTError(er.Error(), w)
			return
		}

		fmt.Fprintf(w, string(result))
	}
}

func MetadataCreateForAttributeDescriptors (conn tgdb.TGConnection, w http.ResponseWriter, r *http.Request, headers map[string]string, body map[string] string) {
	attrDescriptor := impl.DefaultAttributeDescriptor()

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		return
	}
	if gof == nil {
		return
	}

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		return
	}

	gmd.CreateAttributeDescriptor(attrDescriptor.Name, attrDescriptor.AttrType, attrDescriptor.IsArray)

	fmt.Fprintf(w, "SUCCESS")

}

func MetadataQueryForUsers (conn tgdb.TGConnection, w http.ResponseWriter, r *http.Request, headers map[string]string, body map[string] string) {

	adminConnection := conn.(tgdb.TGAdminConnection)
	usersInfo, err := adminConnection.GetUsers()
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}

	var result[] byte
	var er error
	result, er = json.MarshalIndent(usersInfo, "", "\t")
	if er != nil {
		logger.Error("error: " + er.Error())
		handleRESTError(er.Error(), w)
		return
	}
	fmt.Fprintf(w, string(result))
}

func MetadataQueryForConnections (conn tgdb.TGConnection, w http.ResponseWriter, r *http.Request, headers map[string]string, body map[string] string) {

	adminConnection := conn.(tgdb.TGAdminConnection)
	connectionsInfo, err := adminConnection.GetConnections()
	if err != nil {
		handleRESTError(err.Error(), w)
		return
	}

	var result[] byte
	var er error
	result, er = json.MarshalIndent(connectionsInfo, "", "\t")
	if er != nil {
		logger.Error("error: " + er.Error())
		handleRESTError(er.Error(), w)
		return
	}
	fmt.Fprintf(w, string(result))
}


func Query (conn tgdb.TGConnection, w http.ResponseWriter, r *http.Request, headers map[string]string, body map[string] string) {
	queryOptions := initializeQueryOptions(headers)
	gremlinQuery := body["GremlinQuery"]
	//cytoscapeForm := true
	if len(gremlinQuery) > 0 {
		//resultSet, err := conn.ExecuteQuery(gremlinQuery, queryOptions)
		//return obj.TGDBConnection.ExecuteQuery(expr, options)
		//resultSet, err := conn.ExecuteQuery(gremlinQuery, nil)
		resultSet, err := conn.(*impl.AdminConnectionImpl).TGDBConnection.ExecuteQuery(gremlinQuery, queryOptions)
		if err != nil {
			serverErrorCode := err.GetServerErrorCode()
			serverErrorMsg := err.GetErrorMsg()
			logger.Error("ErrorCode: " + strconv.Itoa(serverErrorCode))
			logger.Error("ErrorMessage: " + serverErrorMsg)

			fmt.Fprintf(w, serverErrorMsg)
			return
		} else {
			cytoscapeForm := false
			/*
			responseType, ok := headers["ResponseType"]
			if ok {
				if strings.Compare(strings.ToLower(responseType), "tgdb") == 0 {
					cytoscapeForm = false
				}
			}
			*/

			if cytoscapeForm {
				var result [] byte
				var er error
				resultSetForCytoscape := spotfire.FormCytoscapeResult (resultSet)

				result, er = json.MarshalIndent(resultSetForCytoscape, "", "\t")
				if er != nil {
					//TODO: Handle error
					logger.Error("error: " + er.Error())
				}
				if logger.IsDebug() {
					logger.Debug("Query Result:" + string(result))
				}
				w.Header().Set("Access-Control-Allow-Origin", "*")
				fmt.Fprintf(w, string(result))
			} else {
				var result [] byte
				var er error
				result, er = json.MarshalIndent(resultSet, "", "\t")
				if er != nil {
					logger.Error("error: " + er.Error())
					w.Header().Set("Access-Control-Allow-Origin", "*")
					handleRESTError(er.Error(), w)
					return
				}
				w.Header().Set("Access-Control-Allow-Origin", "*")
				fmt.Fprintf(w, string(result))
			}
		}
	}
}

func initializeQueryOptions (headers map[string]string) (*impl.TGQueryOptionImpl) {
	queryOptions := impl.NewQueryOption()
	batchSize := headers["BatchSize"]
	if len(batchSize) > 0 {
		nBatchSize, err := strconv.Atoi(batchSize)
		if err == nil {
			queryOptions.SetBatchSize(nBatchSize)
		}
	}

	fetchSize := headers["FetchSize"]
	if len(fetchSize) > 0 {
		nFetchSize, err := strconv.Atoi(fetchSize)
		if err == nil {
			queryOptions.SetPreFetchSize(nFetchSize)
		}
	}

	traversalDepth := headers["TraversalDepth"]
	if len(traversalDepth) > 0 {
		nTraversalDepth, err := strconv.Atoi(traversalDepth)
		if err == nil {
			queryOptions.SetTraversalDepth(nTraversalDepth)
		}
	}

	edgeLimit := headers["EdgeLimit"]
	if len(edgeLimit) > 0 {
		nEdgeLimit, err := strconv.Atoi(edgeLimit)
		if err == nil {
			queryOptions.SetEdgeLimit(nEdgeLimit)
		}
	}
	sortAttrName := headers["SortAttrName"]
	if len(sortAttrName) > 0 {
		queryOptions.SetSortAttrName(sortAttrName)
	}

	sortOrder := headers["SortOrder"]
	if len (sortOrder) > 0 {
		if strings.Compare(sortOrder, "asc") == 0 {
			queryOptions.SetSortOrderDsc(false)
		} else if strings.Compare(sortOrder, "dsc") == 0 {
			queryOptions.SetSortOrderDsc(true)
		}
	}
	sortResultLimit := headers["SortResultLimit"]
	if len(sortResultLimit) > 0 {
		nSortResultLimit, err := strconv.Atoi(sortResultLimit)
		if err == nil {
			queryOptions.SetSortResultLimit(nSortResultLimit)
		}
	}
	return queryOptions
}