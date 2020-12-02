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
 * File Name: connectionimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: connectionimpl.go 4621 2020-11-01 00:22:49Z nimish $
 */

package impl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"tgdb"
	"time"
)

//
// TODO: Is this the correct place?
//

const (
	//public static final List<?> EmptyList = new ArrayList<>();
	//EmptyByteArray []byte = []byte{}
	EmptyString string = ""

	U64_NULL       int64 = 0xfffffffffffffff
	U64PACKED_NULL byte  = 0xf0

	INTERNAL_SERVER_ERROR    string = "TGDB-00001"
	TGDB_HNDSHKRESP_ERROR    string = "TGDB-HNDSHKRESP-ERR"
	TGDB_CHANNEL_ERROR       string = "TGDB-CHANNEL-ERR"
	TGDB_SEND_ERROR          string = "TGDB-SENDL-ERR"
	TGDB_CLIENT_READEXTERNAL string = "TGDB-CLIENT-READEXTERNAL"

	DebugEnabled bool = false
)

//static TGLogger gLogger        = TGLogManager.getInstance().getLogger();
var connectionIds int64
var requestIds int64

//type TGConnectionCommand int

const (
	CREATE = 1 + iota
	EXECUTE
	EXECUTEGREMLIN
	EXECUTEGREMLINSTR
	EXECUTED
	CLOSE
)

type TGDBConnection struct {
	channel         tgdb.TGChannel
	connId          int64
	connPoolImpl    tgdb.TGConnectionPool // Connection belongs to a connection pool
	graphObjFactory *GraphObjectFactory   // Intentionally kept private to ensure execution of InitMetaData() before accessing graph objects
	connProperties  tgdb.TGProperties
	addedList       map[int64]tgdb.TGEntity
	changedList     map[int64]tgdb.TGEntity
	removedList     map[int64]tgdb.TGEntity
	attrByTypeList  map[int][]tgdb.TGAttribute
}

func DefaultTGDBConnection() *TGDBConnection {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGDBConnection{})

	newSGDBConnection := &TGDBConnection{
		addedList:      make(map[int64]tgdb.TGEntity, 0),
		changedList:    make(map[int64]tgdb.TGEntity, 0),
		removedList:    make(map[int64]tgdb.TGEntity, 0),
		attrByTypeList: make(map[int][]tgdb.TGAttribute, 0),
	}
	//newSGDBConnection.channel = DefaultAbstractChannel()
	newSGDBConnection.connId = atomic.AddInt64(&connectionIds, 1)
	// We cannot get meta data before we connect to the server
	newSGDBConnection.graphObjFactory = NewGraphObjectFactory(newSGDBConnection)
	newSGDBConnection.connProperties = NewSortedProperties()
	return newSGDBConnection
}

func NewTGDBConnection(conPool *ConnectionPoolImpl, channel tgdb.TGChannel, props tgdb.TGProperties) *TGDBConnection {
	newSGDBConnection := DefaultTGDBConnection()
	newSGDBConnection.connPoolImpl = conPool
	newSGDBConnection.channel = channel
	newSGDBConnection.connProperties = props.(*SortedProperties)
	return newSGDBConnection
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGConnection
/////////////////////////////////////////////////////////////////

func (obj *TGDBConnection) GetConnectionPool() tgdb.TGConnectionPool {
	return obj.connPoolImpl
}

func (obj *TGDBConnection) GetConnectionGraphObjectFactory() *GraphObjectFactory {
	return obj.graphObjFactory
}

func (obj *TGDBConnection) InitMetadata() tgdb.TGError {
	if obj.graphObjFactory == nil {
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil
	}
	gmd := obj.graphObjFactory.GetGraphMetaData()
	if gmd.IsInitialized() {
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil
	}

	// Update the metadata and retrieve it fresh
	_, err := obj.GetGraphMetadata(true)
	if err != nil {
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil
	}
	return nil
}

// SetConnectionPool sets connection pool
func (obj *TGDBConnection) SetConnectionPool(connPool tgdb.TGConnectionPool) {
	obj.connPoolImpl = connPool
}

// SetConnectionProperties sets connection properties
func (obj *TGDBConnection) SetConnectionProperties(connProps tgdb.TGProperties) {
	obj.connProperties = connProps.(*SortedProperties)
}

/////////////////////////////////////////////////////////////////
// Private functions for types.TGConnection
/////////////////////////////////////////////////////////////////

func fixUpAttrDescriptors(response *CommitTransactionResponse, attrDescSet []tgdb.TGAttributeDescriptor) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering TGDBConnection:fixUpAttrDescriptors"))
	}
	attrDescCount := response.GetAttrDescCount()
	attrDescIdList := response.GetAttrDescIdList()
	for i := 0; i < attrDescCount; i++ {
		tempId := attrDescIdList[(i * 2)]
		realId := attrDescIdList[((i * 2) + 1)]

		for _, attrDesc := range attrDescSet {
			desc := attrDesc.(*AttributeDescriptor)
			if attrDesc.GetAttributeId() == tempId {
				logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:fixUpAttrDescriptors - Replace descriptor: '%d' by '%d'", attrDesc.GetAttributeId(), realId))
				desc.SetAttributeId(realId)
				break
			}
		}
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:fixUpAttrDescriptors"))
	}
}

func fixUpEntities(obj tgdb.TGConnection, response *CommitTransactionResponse) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering TGDBConnection:fixUpEntities"))
	}
	addedIdCount := response.GetAddedEntityCount()
	addedIdList := response.GetAddedIdList()
	for i := 0; i < addedIdCount; i++ {
		tempId := addedIdList[(i * 3)]
		realId := addedIdList[((i * 3) + 1)]
		version := addedIdList[((i * 3) + 2)]

		for _, addEntity := range obj.GetAddedList() {
			if addEntity.GetVirtualId() == tempId {
				logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:fixUpEntities - Replace entity id: '%d' by '%d'", tempId, realId))
				addEntity.SetEntityId(realId)
				addEntity.SetIsNew(false)
				addEntity.SetVersion(int(version))
				break
			}
		}
	}

	updatedIdCount := response.GetUpdatedEntityCount()
	updatedIdList := response.GetUpdatedIdList()
	for i := 0; i < updatedIdCount; i++ {
		id := updatedIdList[(i * 2)]
		version := updatedIdList[((i * 2) + 1)]

		for _, modEntity := range obj.GetChangedList() {
			if modEntity.GetVirtualId() == id {
				logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:fixUpEntities - Replace entity version: '%d' to '%d'", id, version))
				modEntity.SetVersion(int(version))
				break
			}
		}
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:fixUpEntities"))
	}
}

func createChannelRequest(obj tgdb.TGConnection, verb int) (tgdb.TGMessage, tgdb.TGChannelResponse, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:createChannelRequest for Verb: '%s'", GetVerb(verb).GetName()))
	}
	cn := GetConfigFromKey(ConnectionOperationTimeoutSeconds)
	//logger.Debug(fmt.Sprintf("Inside AbstractChannel::channelTryRepeatConnect config for ConnectionOperationTimeoutSeconds is '%+v", cn))
	timeout := obj.GetConnectionProperties().GetPropertyAsInt(cn)
	requestId := atomic.AddInt64(&requestIds, 1)

	//logger.Debug(fmt.Sprint("Inside TGDBConnection::createChannelRequest about to create channel.NewBlockingChannelResponse()"))
	// Create a non-blocking channel response
	channelResponse := NewBlockingChannelResponse(requestId, int64(timeout))

	// Use Message Factory method to create appropriate message structure (class) based on input type
	//msgRequest, err := pdu.CreateMessageForVerb(verb)
	msgRequest, err := CreateMessageWithToken(verb, obj.GetChannel().GetAuthToken(), obj.GetChannel().GetSessionId())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:createChannelRequest - CreateMessageForVerb failed for Verb: '%s' w/ '%+v'", GetVerb(verb).GetName(), err.Error()))
		return nil, nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGDBConnection:createChannelRequest for Verb: '%s' w/ MessageRequest: '%+v' ChannelResponse: '%+v'", GetVerb(verb).GetName(), msgRequest, channelResponse))
	}
	return msgRequest, channelResponse, nil
}

func configureGetRequest(getReq *GetEntityRequestMessage, reqProps tgdb.TGProperties) {
	if getReq == nil || reqProps == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::configureGetRequest as getReq == nil || reqProps == nil"))
		return
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:configureGetRequest w/ EntityRequest: '%+v'", getReq.String()))
	}
	props := reqProps.(*TGQueryOptionImpl)
	getReq.SetFetchSize(props.GetPreFetchSize())
	getReq.SetTraversalDepth(props.GetTraversalDepth())
	getReq.SetEdgeLimit(props.GetEdgeLimit())
	getReq.SetBatchSize(props.GetBatchSize())
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGDBConnection:configureGetRequest w/ EntityRequest: '%+v'", getReq.String()))
	}
	return
}

func configureQueryRequest(qryReq *QueryRequestMessage, reqProps tgdb.TGProperties) {
	if qryReq == nil || reqProps == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection:configureQueryRequest as getReq == nil || reqProps == nil"))
		return
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:configureQueryRequest w/ QueryRequest: '%+v'", qryReq.String()))
	}
	props := reqProps.(*TGQueryOptionImpl)
	qryReq.SetFetchSize(props.GetPreFetchSize())
	qryReq.SetTraversalDepth(props.GetTraversalDepth())
	qryReq.SetEdgeLimit(props.GetEdgeLimit())
	qryReq.SetBatchSize(props.GetBatchSize())
	qryReq.SetSortAttrName(props.GetSortAttrName())
	qryReq.SetSortOrderDsc(props.IsSortOrderDsc())
	qryReq.SetSortResultLimit(props.GetSortResultLimit())
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGDBConnection:configureQueryRequest w/ QueryRequest: '%+v'", qryReq.String()))
	}
	return
}

func (obj *TGDBConnection) populateResultSetFromQueryResponse(resultId int, msgResponse *QueryResponseMessage) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:populateResultSetFromQueryResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	}
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse as msgResponse does not have any results"))
		errMsg := "TGDBConnection::populateResultSetFromQueryResponse does not have any results in QueryResponseMessage"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]tgdb.TGEntity, 0)
	var rSet *ResultSet

	currResultCount := 0
	resultCount := msgResponse.GetResultCount()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse read resultCount: '%d' FetchedEntityCount: '%d'", resultCount, len(fetchedEntities)))
	}
	if resultCount > 0 {
		respStream.SetReferenceMap(fetchedEntities)
		rSet = NewResultSet(obj, resultId)
	}

	totalCount := msgResponse.GetTotalCount()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse read totalCount: '%d'", totalCount))
	}
	for i := 0; i < totalCount; i++ {
		entityType, err := respStream.(*ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse - unable to read entityType in the response stream"))
			errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to read entity type in the response stream"
			return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		kindId := tgdb.TGEntityKind(entityType)
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse read #'%d'-entityType: '%+v', kindId: '%s'", i, entityType, kindId.String()))
		}
		if kindId != tgdb.EntityKindInvalid {
			entityId, err := respStream.(*ProtocolDataInputStream).ReadLong()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse - unable to read entityId in the response stream"))
				errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to read entity type in the response stream"
				return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			entity := fetchedEntities[entityId]
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse read entityId: '%d', kindId: '%s', entity: '%+v'", entityId, kindId.String(), entity))
			}
			switch kindId {
			case tgdb.EntityKindNode:
				var node *Node
				if entity == nil {
					node, nErr := obj.graphObjFactory.CreateNode()
					if nErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to create a new node from the response stream"
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.GetErrorDetails())
					}
					entity = node
					fetchedEntities[entityId] = node
					if logger.IsDebug() {
						logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse created new node: '%+v' FetchedEntityCount: '%d'", node, len(fetchedEntities)))
					}
				}
				node = entity.(*Node)
				err := node.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromQueryResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.GetErrorDetails())
				}
				if currResultCount < resultCount {
					rSet.AddEntityToResultSet(node)
				}
				//logger.Debug(fmt.Sprintf("======> ======> After node.ReadExternal() FetchedEntityCount: '%d'", len(fetchedEntities)))
				if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("======> ======> Node w/ Edges: '%+v'\n", node.GetEdges()))
				}
			case tgdb.EntityKindEdge:
				var edge *Edge
				if entity == nil {
					//edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, types.DirectionTypeBiDirectional)
					edge, eErr := obj.graphObjFactory.CreateEntity(tgdb.EntityKindEdge)
					if eErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to create a new bi-directional edge from the response stream"
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, eErr.GetErrorDetails())
					}
					entity = edge
					fetchedEntities[entityId] = edge
					if logger.IsDebug() {
						logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse created new edge: '%+v' FetchedEntityCount: '%d'", edge, len(fetchedEntities)))
					}
				}
				edge = entity.(*Edge)
				err := edge.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromQueryResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.GetErrorDetails())
				}
				if currResultCount < resultCount {
					rSet.AddEntityToResultSet(edge)
				}
				//logger.Debug(fmt.Sprintf("======> ======> After edge.ReadExternal() FetchedEntityCount: '%d'", len(fetchedEntities)))
				if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("======> ======> Edge w/ Vertices: '%+v'\n", edge.GetVertices()))
				}
			case tgdb.EntityKindGraph:
				// TODO: Revisit later - Should we break after throwing/logging an error
				continue
			}
			//if entity != nil {
			//	logger.Debug(fmt.Sprintf("======> TGDBConnection::populateResultSetFromQueryResponse entityId: '%d', kindId: '%d', entityType: '%+v'\n", entityId, kindId, kindId.String()))
			//	attrList, _ := entity.GetAttributes()
			//	for _, attrib := range attrList {
			//		logger.Debug(fmt.Sprintf("======> Attribute Value: '%+v'\n", attrib.GetValue()))
			//	}
			//	if kindId == types.EntityKindNode {
			//		edges := (entity.(*model.Node)).GetEdges()
			//		logger.Debug(fmt.Sprintf("======> Node w/ Edges: '%+v'\n", edges))
			//	} else if kindId == types.EntityKindEdge {
			//		vertices := (entity.(*model.Edge)).GetVertices()
			//		logger.Debug(fmt.Sprintf("======> Edge w/ Vertices: '%+v'\n", vertices))
			//	}
			//}
		} else {
			logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:populateResultSetFromQueryResponse - Received invalid entity kind %d", kindId))
		} // Valid entity types
	} // End of for loop
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning TGDBConnection:populateResultSetFromQueryResponse w/ ResultSet: '%+v'", rSet))
	}
	return rSet, nil
}

func (obj *TGDBConnection) populateResultSetFromGetEntitiesResponse(msgResponse *GetEntityResponseMessage) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:populateResultSetFromGetEntitiesResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	}
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse as msgResponse does not have any results"))
		errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse does not have any results"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]tgdb.TGEntity, 0)

	totalCount, err := respStream.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read totalCount in the response stream w/ error: '%s'", err.Error()))
		errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read total count of entities in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse extracted totalCount: '%d'", totalCount))
	}
	if totalCount > 0 {
		respStream.SetReferenceMap(fetchedEntities)
	}

	rSet := NewResultSet(obj, msgResponse.GetResultId())
	// Number of entities matches the search.  Exclude the related entities
	currResultCount := 0
	_, err = respStream.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read count of result entities in the response stream w/ error: '%s'", err.Error()))
		errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read count of result entities in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}

	for i := 0; i < totalCount; i++ {
		isResult, err := respStream.(*ProtocolDataInputStream).ReadBoolean()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read isResult in the response stream w/ error: '%s'", err.Error()))
			errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read count of result entities in the response stream"
			return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse read isResult: '%+v'", isResult))
		}
		entityType, err := respStream.(*ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
			errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read entity type in the response stream"
			return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		kindId := tgdb.TGEntityKind(entityType)
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse extracted entityType: '%+v', kindId: '%d'", entityType, kindId))
		}
		if kindId != tgdb.EntityKindInvalid {
			entityId, err := respStream.(*ProtocolDataInputStream).ReadLong()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read entityId in the response stream w/ error: '%s'", err.Error()))
				errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read entity type in the response stream"
				return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			entity := fetchedEntities[entityId]
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse extracted entityId: '%d', entity: '%+v'", entityId, entity))
			}
			switch kindId {
			case tgdb.EntityKindNode:
				var node *Node
				if entity == nil {
					node, nErr := obj.graphObjFactory.CreateNode()
					if nErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to create a new node from the response stream"
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.GetErrorDetails())
					}
					entity = node
					fetchedEntities[entityId] = node
					if logger.IsDebug() {
											logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse created new node: '%+v'", node))
					}
				}
				node = entity.(*Node)
				err := node.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromGetEntitiesResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				if isResult {
					rSet.AddEntityToResultSet(node)
					currResultCount++
				}
			case tgdb.EntityKindEdge:
				var edge *Edge
				if entity == nil {
					edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, tgdb.DirectionTypeBiDirectional)
					if eErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
						// TODO: Revisit later - Should we break after throwing/logging an error
						//continue
						errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to create a new bi-directional edge from the response stream"
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, eErr.Error())
					}
					entity = edge
					fetchedEntities[entityId] = edge
					if logger.IsDebug() {
											logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse created new edge: '%+v'", edge))
					}
				}
				edge = entity.(*Edge)
				err := edge.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromGetEntitiesResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				if isResult {
					rSet.AddEntityToResultSet(edge)
					currResultCount++
				}
			case tgdb.EntityKindGraph:
				// TODO: Revisit later - Should we break after throwing/logging an error
				continue
			}
		} else {
			logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:populateResultSetFromGetEntitiesResponse - Received invalid entity kind %d", kindId))
		} // Valid entity types
	} // End of for loop
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning TGDBConnection:populateResultSetFromGetEntitiesResponse w/ ResultSe: '%+v'", rSet))
	}
	return rSet, nil
}

func (obj *TGDBConnection) populateResultSetFromGetEntityResponse(msgResponse *GetEntityResponseMessage) (tgdb.TGEntity, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:populateResultSetFromGetEntityResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	}
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse as msgResponse does not have any results"))
		errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse does not have any results"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]tgdb.TGEntity, 0)

	var entityFound tgdb.TGEntity
	respStream.SetReferenceMap(fetchedEntities)

	count, err := respStream.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to read count in the response stream w/ error: '%s'", err.Error()))
		errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to read count of entities in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse extracted Count: '%d'", count))
	}
	if count > 0 {
		respStream.SetReferenceMap(fetchedEntities)
		for i := 0; i < count; i++ {
			entityType, err := respStream.(*ProtocolDataInputStream).ReadByte()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
				errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to read entity type in the response stream"
				return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			kindId := tgdb.TGEntityKind(entityType)
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse extracted entityType: '%+v', kindId: '%d'", entityType, kindId))
			}
			if kindId != tgdb.EntityKindInvalid {
				entityId, err := respStream.(*ProtocolDataInputStream).ReadLong()
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to read entityId in the response stream w/ error: '%s'", err.Error()))
					errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to read entity type in the response stream"
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				entity := fetchedEntities[entityId]
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse extracted entityId: '%d', entity: '%+v'", entityId, entity))
				}
				switch kindId {
				case tgdb.EntityKindNode:
					// Need to put shell object into map to be deserialized later
					var node *Node
					if entity == nil {
						node, nErr := obj.graphObjFactory.CreateNode()
						if nErr != nil {
							logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
							// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
							//continue
							errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to create a new node from the response stream"
							return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.Error())
						}
						entity = node
						fetchedEntities[entityId] = node
						if entityFound == nil {
							entityFound = node
						}
						if logger.IsDebug() {
													logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse created new node: '%+v'", entity))
						}
					}
					node = entity.(*Node)
					err := node.ReadExternal(respStream)
					if err != nil {
						errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromGetEntityResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
						logger.Error(errMsg)
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
					}
				case tgdb.EntityKindEdge:
					var edge *Edge
					if entity == nil {
						edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, tgdb.DirectionTypeBiDirectional)
						if eErr != nil {
							logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
							// TODO: Revisit later - Should we break after throwing/logging an error
							//continue
							errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to create a new bi-directional edge from the response stream"
							return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, eErr.Error())
						}
						entity = edge
						fetchedEntities[entityId] = edge
						if entityFound == nil {
							entityFound = edge
						}
						if logger.IsDebug() {
													logger.Debug(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse created new edge: '%+v'", edge))
						}
					}
					edge = entity.(*Edge)
					err := edge.ReadExternal(respStream)
					if err != nil {
						errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromGetEntityResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
						logger.Error(errMsg)
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
					}
				case tgdb.EntityKindGraph:
					// TODO: Revisit later - Should we break after throwing/logging an error
					continue
				}
			} else {
				logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:populateResultSetFromGetEntityResponse - Received invalid entity kind %d", kindId))
			} // Valid entity types
		} // End of for loop
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning TGDBConnection:populateResultSetFromGetEntityResponse w/ Entity: '%+v'", entityFound))
	}
	return entityFound, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConnection
/////////////////////////////////////////////////////////////////

// Commit commits the current transaction on this connection
func (obj *TGDBConnection) Commit() (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TGDBConnection:Commit"))
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Commit - about to loop through addedList to include existing nodes to the changed list if it's part of a new edge"))
	}
	// Include existing nodes to the changed list if it's part of a new edge
	for _, addEntity := range obj.GetAddedList() {
		if addEntity.GetEntityKind() == tgdb.EntityKindEdge {
			nodes := addEntity.(*Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*Node)
					if !node.GetIsNew() {
						obj.changedList[node.GetVirtualId()] = node
						logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:Commit - Existing node '%d' added to change list for a new edge", node.GetVirtualId()))
					}
				}
			}
		}
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Commit - about to loop through changedList to include existing nodes to the changed list even for edge update"))
	}
	// Need to include existing node to the changed list even for edge update
	for _, modEntity := range obj.GetChangedList() {
		if modEntity.GetEntityKind() == tgdb.EntityKindEdge {
			nodes := modEntity.(*Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*Node)
					if !node.GetIsNew() {
						obj.changedList[node.GetVirtualId()] = node
						logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:Commit - Existing node '%d' added to change list for an existing edge '%d'", node.GetVirtualId(), modEntity.GetVirtualId()))
					}
				}
			}
		}
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Commit - about to loop through removedList to include existing nodes to the changed list even for edge update"))
	}
	// Need to include existing node to the changed list even for edge update
	for _, delEntity := range obj.GetRemovedList() {
		if delEntity.GetEntityKind() == tgdb.EntityKindEdge {
			nodes := delEntity.(*Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*Node)
					if !node.GetIsNew() {
						if obj.removedList[node.GetVirtualId()] == nil {
							obj.changedList[node.GetVirtualId()] = node
							logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:Commit - Existing node '%d' added to change list for an edge %d to be deleted", node.GetVirtualId(), delEntity.GetVirtualId()))
						}
					}
				}
			}
		}
	}
	//For deleted edge and node, we don't immediately change the effected nodes or edges.

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Commit about to createChannelRequest() for: VerbCommitTransactionRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, VerbCommitTransactionRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:Commit - unable to createChannelRequest(VerbCommitTransactionRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*CommitTransactionRequest)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Commit about to GetAttributeDescriptors() for: VerbCommitTransactionRequest"))
	}
	gof := obj.graphObjFactory
	attrDescSet, aErr := gof.GetGraphMetaData().GetNewAttributeDescriptors()
	if aErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:Commit - unable to gmd.GetAttributeDescriptors() w/ error: '%s'", aErr.Error()))
		return nil, aErr
	}
	queryRequest.AddCommitLists(obj.addedList, obj.changedList, obj.removedList, attrDescSet)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Commit about to channelSendRequest() for: VerbCommitTransactionRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:Commit - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGDBConnection::Commit received response for: VerbCommitTransactionRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*CommitTransactionResponse)

	if response.HasException() {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:Commit as response has exceptions"))
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil, nil
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Commit about to fixUpAttrDescriptors()"))
	}
	fixUpAttrDescriptors(response, attrDescSet)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Commit about to obj.fixUpEntities()"))
	}
	fixUpEntities(obj, response)

	for _, delEntity := range obj.GetRemovedList() {
		delEntity.SetIsDeleted(true)
	}

	// Reset and clear all the lists
	for _, modEntity := range obj.GetChangedList() {
		modEntity.ResetModifiedAttributes()
	}
	for _, newEntity := range obj.GetAddedList() {
		newEntity.ResetModifiedAttributes()
	}

	obj.addedList = make(map[int64]tgdb.TGEntity, 0)
	obj.changedList = make(map[int64]tgdb.TGEntity, 0)
	obj.removedList = make(map[int64]tgdb.TGEntity, 0)
	obj.attrByTypeList = make(map[int][]tgdb.TGAttribute, 0)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:Commit"))
	}
	return nil, nil
}

// Connect establishes a network connection to the TGDB server
func (obj *TGDBConnection) Connect() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:Connect for connection: '%+v'", obj))
	}
	err := obj.GetChannel().Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection::Connect - error in obj.GetChannel().Connect() as '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Connect about to obj.GetChannel().Start()"))
	}
	err = obj.GetChannel().Start()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection::Connect - error in obj.GetChannel().Start() as '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:Connect"))
	}
	return nil
}

// CloseQuery closes a specific query and associated objects
func (obj *TGDBConnection) CloseQuery(queryHashId int64) (tgdb.TGQuery, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:CloseQuery for QueryHashId: '%+v'", queryHashId))
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::CloseQuery about to createChannelRequest() for: VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, VerbQueryRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:CloseQuery - unable to createChannelRequest(VerbQueryRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(CLOSE)
	queryRequest.SetQueryHashId(queryHashId)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::CloseQuery about to obj.GetChannel().SendRequest() for: VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	_, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:CloseQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGDBConnection:CloseQuery"))
	}
	return nil, nil
}

// CreateQuery creates a reusable query object that can be used to execute one or more statement
func (obj *TGDBConnection) CreateQuery(expr string) (tgdb.TGQuery, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:CreateQuery for Query: '%+v'", expr))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::CreateQuery about to createChannelRequest() for: VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, VerbQueryRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:CreateQuery - unable to createChannelRequest(VerbQueryRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(CREATE)
	queryRequest.SetQuery(expr)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::CreateQuery about to obj.GetChannel().SendRequest() for: VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:CreateQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGDBConnection::CreateQuery received response for: VerbQueryRequest as '%+v'", msgResponse))
	}

	response := msgResponse.(*QueryResponseMessage)
	queryHashId := response.GetQueryHashId()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:CreateQuery"))
	}
	if response.GetResult() == 0 && queryHashId > 0 {
		return NewQuery(obj, queryHashId), nil
	}
	return nil, nil
}

// DecryptBuffer decrypts the encrypted buffer by sending a DecryptBufferRequest to the server
func (obj *TGDBConnection) DecryptBuffer(is tgdb.TGInputStream) ([]byte, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering TGDBConnection:DecryptBuffer w/ EncryptedBuffer"))
	}
	//err := obj.InitMetadata()
	//if err != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:DecryptBuffer - unable to initialize metadata w/ error: '%s'", err.Error()))
	//	return nil, err
	//}
	//obj.connPoolImpl.AdminLock()
	//defer obj.connPoolImpl.AdminUnlock()
	//
	//logger.Debug(fmt.Sprint("Inside TGDBConnection::DecryptBuffer about to createChannelRequest() for: VerbDecryptBufferRequest"))
	//// Create a channel request
	//msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbDecryptBufferRequest)
	//if cErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:DecryptBuffer - unable to createChannelRequest(VerbDecryptBufferRequest w/ error: '%s'", cErr.Error()))
	//	return nil, cErr
	//}
	//decryptRequest := msgRequest.(*DecryptBufferRequestMessage)
	//decryptRequest.SetEncryptedBuffer(encryptedBuf)
	//
	//logger.Debug(fmt.Sprint("Inside TGDBConnection::DecryptBuffer about to obj.GetChannel().SendRequest() for: VerbDecryptBufferRequest"))
	//// Execute request on channel and get the response
	//msgResponse, channelErr := obj.GetChannel().SendRequest(decryptRequest, channelResponse.(*channel.BlockingChannelResponse))
	//if channelErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:DecryptBuffer - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
	//	return nil, channelErr
	//}
	//logger.Debug(fmt.Sprintf("Inside TGDBConnection::DecryptBuffer received response for: VerbGetLargeObjectRequest as '%+v'", msgResponse))
	//response := msgResponse.(*DecryptBufferResponseMessage)
	//
	//if response == nil {
	//	errMsg := "TGDBConnection::DecryptBuffer does not have any results in GetLargeObjectResponseMessage"
	//	logger.Error(errMsg)
	//	return nil, GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	//}
	//
	//logger.Log(fmt.Sprintf("Returning TGDBConnection:DecryptBuffer w/ '%+v'", response.GetDecryptedBuffer()))
	//return response.GetDecryptedBuffer(), nil
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Decrypt(is)
}

// DecryptEntity decrypts the encrypted entity using channel's data cryptographer
func (obj *TGDBConnection) DecryptEntity(entityId int64) ([]byte, tgdb.TGError) {
	buf, err := obj.GetLargeObjectAsBytes(entityId, true)
	if err != nil {
		return nil, err
	}
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Decrypt(NewProtocolDataInputStream(buf))
}

// DeleteEntity marks an ENTITY for delete operation. Upon commit, the entity will be deleted from the database
func (obj *TGDBConnection) DeleteEntity(entity tgdb.TGEntity) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:DeleteEntity for Entity: '%+v'", entity))
	}
	obj.removedList[entity.GetVirtualId()] = entity
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:DeleteEntity"))
	}
	return nil
}

// Disconnect breaks the connection from the TGDB server
func (obj *TGDBConnection) Disconnect() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TGDBConnection:Disconnect"))
	}
	err := obj.GetChannel().Disconnect()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:Disconnect - unable to channel.Disconnect() w/ error: '%s'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::Disconnect about to obj.GetChannel().Stop()"))
	}
	obj.GetChannel().Stop(false)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:Disconnect"))
	}
	return nil
}

// EncryptEntity encrypts the encrypted entity using channel's data cryptographer
func (obj *TGDBConnection) EncryptEntity(rawBuffer []byte) ([]byte, tgdb.TGError) {
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Encrypt(rawBuffer)
}

// ExecuteGremlinQuery executes a Gremlin Grammer-Based query with  query options
func (obj *TGDBConnection) ExecuteGremlinQuery(expr string, collection []interface{}, options tgdb.TGQueryOption) ([]interface{}, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:ExecuteGremlinQuery for Query: '%+v'", expr))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinQuery about to createChannelRequest() for: VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinQuery - unable to createChannelRequest(VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(EXECUTEGREMLIN)
	queryRequest.SetQuery(expr)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinQuery about to obj.configureQueryRequest() for: VerbQueryRequest"))
	}
	configureQueryRequest(queryRequest, options)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinQuery about to obj.GetChannel().SendRequest() for: VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGDBConnection::ExecuteGremlinQuery received response for: VerbQueryRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::ExecuteGremlinQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	respStream := response.GetEntityStream()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGDBConnection:ExecuteGremlinQuery w/ '%+v'", response))
	}
	collection, err = FillCollection(respStream, obj.graphObjFactory, collection)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinQuery - unable to query.FillCollection w/ error: '%s'", err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGDBConnection:ExecuteGremlinQuery w/ '%+v'", collection))
	}
	return collection, nil
}

// ExecuteGremlinStrQuery executes a Gremlin Grammer-Based string query with  query options
func (obj *TGDBConnection) ExecuteGremlinStrQuery(strQuery string, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:ExecuteGremlinStrQuery for Query: '%+v'", strQuery))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinStrQuery about to createChannelRequest() for: VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinStrQuery - unable to createChannelRequest(VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(EXECUTEGREMLINSTR)
	queryRequest.SetQuery(strQuery)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinStrQuery about to obj.configureQueryRequest() for: VerbQueryRequest"))
	}
	configureQueryRequest(queryRequest, options)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinStrQuery about to obj.GetChannel().SendRequest() for: VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinStrQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGDBConnection::ExecuteGremlinStrQuery received response for: VerbQueryRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*QueryResponseMessage)

	if response.exception != nil {
		return nil, response.exception
	}

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::ExecuteGremlinStrQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	// TODO: Revisit later once Gremlin Package and GremlinQueryResult are implemented
	//This is just a dummy value where > 0 means it has results
	//resultCount := response.GetResultCount()
	resultSet := NewResultSet(obj, 0)

	// wiring of annotation to resultsetmetadata
	resultSet.SetResultTypeAnnotation((response.GetResultTypeAnnot()))

	respStream := response.GetEntityStream()
	resultList, err := FillCollection(respStream, obj.graphObjFactory, resultSet.GetResults())

	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinQuery - unable to query.FillCollection w/ error: '%s'", err.Error()))
		return nil, err
	}

	resultSet.ResultList = resultList
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGDBConnection:ExecuteGremlinStrQuery w/ '%+v'", resultSet))
	}
	return resultSet, nil
}

// ExecuteQuery executes an immediate query with associated query options
func (obj *TGDBConnection) ExecuteQuery(expr string, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:ExecuteQuery for Query: '%+v'", expr))
	}

	// TODO: Revisit later once Gremlin Package and GremlinQueryResult are implemented
	cn := GetConfigFromKey(ConnectionDefaultQueryLanguage)
	queryLang := obj.GetConnectionProperties().GetProperty(cn, "tgql")

	idx := strings.Index(expr, "://")
	if idx != -1 {
		tokens := strings.Split(expr, "://")
		switch tokens[0] {
		case "tgql":
			return obj.ExecuteTGDBQuery(tokens[1], options)
		case "gremlin":
			return obj.ExecuteGremlinStrQuery(tokens[1], options)
		default:
			return NewResultSet(obj, 0), nil
		}
	} else {
		switch queryLang {
		case "tgql":
			return obj.ExecuteTGDBQuery(expr, options)
		case "gremlin":
			return obj.ExecuteGremlinStrQuery(expr, options)
		default:
			return NewResultSet(obj, 0), nil
		}
	}
	return nil, nil
}

// ExecuteTGDBQuery executes an immediate query with associated query options
func (obj *TGDBConnection) ExecuteTGDBQuery(expr string, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:ExecuteTGDBQuery for Query: '%+v'", expr))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteTGDBQuery about to createChannelRequest() for: VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteTGDBQuery - unable to createChannelRequest(VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery(expr)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteTGDBQuery about to obj.configureQueryRequest() for: VerbQueryRequest"))
	}
	configureQueryRequest(queryRequest, options)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteTGDBQuery about to obj.GetChannel().SendRequest() for: VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteTGDBQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TGDBConnection::ExecuteTGDBQuery received response for: VerbQueryRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::ExecuteTGDBQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning TGDBConnection:ExecuteTGDBQuery w/ '%+v'", response))
	}
	return obj.populateResultSetFromQueryResponse(0, response)
}

// ExecuteQueryWithFilter executes an immediate query with specified filter & query options
// The query option is place holder at this time
// @param expr A subset of SQL-92 where clause
// @param edgeFilter filter used for selecting edges to be returned
// @param traversalCondition condition used for selecting edges to be traversed and returned
// @param endCondition condition used to stop the traversal
// @param option Query options for executing. Can be null, then it will use the default option
func (obj *TGDBConnection) ExecuteQueryWithFilter(expr, edgeFilter, traversalCondition, endCondition string, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:ExecuteQueryWithFilter for Query: '%+v', EdgeFilter: '%+v', Traversal: '%+v', EndCondition: '%+v'", expr, edgeFilter, traversalCondition, endCondition))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithFilter about to createChannelRequest() for: VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteQueryWithFilter - unable to createChannelRequest(VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery(expr)
	queryRequest.SetEdgeFilter(edgeFilter)
	queryRequest.SetTraversalCondition(traversalCondition)
	queryRequest.SetEndCondition(endCondition)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGDBConnection::ExecuteQueryWithFilter about to obj.configureQueryRequest() for: VerbQueryRequest"))
	}
	configureQueryRequest(queryRequest, options)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithFilter about to obj.GetChannel().SendRequest() for: VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteQueryWithFilter - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGDBConnection::ExecuteQueryWithFilter received response for: VerbQueryRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::ExecuteQueryWithFilter - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGDBConnection:ExecuteQueryWithFilter w/ '%+v'", response))
	}
	return obj.populateResultSetFromQueryResponse(0, response)
}

// ExecuteQueryWithId executes an immediate query for specified id & query options
func (obj *TGDBConnection) ExecuteQueryWithId(queryHashId int64, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:ExecuteQueryWithId for QueryHashId: '%+v'", queryHashId))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithId about to createChannelRequest() for: VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteQueryWithId - unable to createChannelRequest(VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(EXECUTED)
	queryRequest.SetQueryHashId(queryHashId)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithId about to obj.configureQueryRequest() for: VerbQueryRequest"))
	}
	configureQueryRequest(queryRequest, options)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithId about to obj.GetChannel().SendRequest() for: VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	_, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteQueryWithId - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:ExecuteQueryWithId"))
	}
	return nil, nil
}

// GetAddedList gets a list of added entities
func (obj *TGDBConnection) GetAddedList() map[int64]tgdb.TGEntity {
	return obj.addedList
}

// GetChangedList gets a list of changed entities
func (obj *TGDBConnection) GetChangedList() map[int64]tgdb.TGEntity {
	return obj.changedList
}

// GetChangedList gets the communication channel associated with this connection
func (obj *TGDBConnection) GetChannel() tgdb.TGChannel {
	return obj.channel
}

// GetConnectionId gets connection identifier
func (obj *TGDBConnection) GetConnectionId() int64 {
	return obj.connId
}

// GetConnectionProperties gets a list of connection properties
func (obj *TGDBConnection) GetConnectionProperties() tgdb.TGProperties {
	return obj.connProperties
}

// GetEntities gets a result set of entities given an non-uniqueKey
func (obj *TGDBConnection) GetEntities(qryKey tgdb.TGKey, props tgdb.TGProperties) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:GetEntities for QueryKey: '%+v'", qryKey))
	}
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:GetEntities - unable to InitMetadata"))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if props == nil {
		props = NewQueryOption()
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::GetEntities about to createChannelRequest() for: VerbGetEntityRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbGetEntityRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetEntities - unable to createChannelRequest(VerbGetEntityRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*GetEntityRequestMessage)
	queryRequest.SetCommand(2)
	queryRequest.SetKey(qryKey)
	configureGetRequest(queryRequest, props)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::GetEntities about to obj.GetChannel().SendRequest() for: VerbGetEntityRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetEntities - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TGDBConnection::GetEntities received response for: VerbGetEntityRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*GetEntityResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::GetEntities - The request does not have any results in GetEntityResponseMessage"))
		return nil, nil
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning TGDBConnection:GetEntities w/ '%+v'", response))
	}
	return obj.populateResultSetFromGetEntitiesResponse(response)
}

// GetEntity gets an Entity given an UniqueKey for the Object
func (obj *TGDBConnection) GetEntity(qryKey tgdb.TGKey, options tgdb.TGQueryOption) (tgdb.TGEntity, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:GetEntity for QueryKey: '%+v'", qryKey))
	}
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:GetEntity - unable to InitMetadata"))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if options == nil {
		options = NewQueryOption()
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::GetEntity about to createChannelRequest() for: VerbGetEntityRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbGetEntityRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetEntity - unable to createChannelRequest(VerbGetEntityRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*GetEntityRequestMessage)
	queryRequest.SetCommand(0)
	queryRequest.SetKey(qryKey)
	configureGetRequest(queryRequest, options)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::GetEntity about to obj.GetChannel().SendRequest() for: VerbGetEntityRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetEntity - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TGDBConnection::GetEntity received response for: VerbGetEntityRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*GetEntityResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::GetEntity - The request does not have any results in GetEntityResponseMessage"))
		return nil, nil
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning TGDBConnection:GetEntity w/ '%+v'", response))
	}
	return obj.populateResultSetFromGetEntityResponse(response)
}

// GetGraphMetadata gets the Graph Metadata
func (obj *TGDBConnection) GetGraphMetadata(refresh bool) (tgdb.TGGraphMetadata, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering TGDBConnection:GetGraphMetadata"))
	}
	if refresh {
		obj.connPoolImpl.AdminLock()
		defer obj.connPoolImpl.AdminUnlock()

		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside TGDBConnection::GetGraphMetadata about to createChannelRequest() for: VerbMetadataRequest"))
		}
		// Create a channel request
		msgRequest, channelResponse, err := createChannelRequest(obj, VerbMetadataRequest)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetGraphMetadata - unable to createChannelRequest(VerbMetadataRequest w/ error: '%s'", err.Error()))
			return nil, err
		}
		metaRequest := msgRequest.(*MetadataRequest)
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside TGDBConnection::GetGraphMetadata createChannelRequest() returned MsgRequest: '%+v' ChannelResponse: '%+v'", msgRequest, channelResponse.(*BlockingChannelResponse)))
		}

		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside TGDBConnection::GetGraphMetadata about to obj.GetChannel().SendRequest() for: VerbMetadataRequest"))
		}
		// Execute request on channel and get the response
		msgResponse, channelErr := obj.GetChannel().SendRequest(metaRequest, channelResponse.(*BlockingChannelResponse))
		if channelErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetGraphMetadata - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
			return nil, channelErr
		}
		//logger.Debug(fmt.Sprintf("Inside TGDBConnection::GetGraphMetadata received response for: VerbMetadataRequest as '%+v'", msgResponse))

		response := msgResponse.(*MetadataResponse)
		attrDescList := response.GetAttrDescList()
		edgeTypeList := response.GetEdgeTypeList()
		nodeTypeList := response.GetNodeTypeList()

		gmd := obj.graphObjFactory.GetGraphMetaData()
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside TGDBConnection::GetGraphMetadata about to update GraphMetadata"))
		}
		uErr := gmd.UpdateMetadata(attrDescList, nodeTypeList, edgeTypeList)
		if uErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetGraphMetadata - unable to gmd.UpdateMetadata() w/ error: '%s'", uErr.Error()))
			return nil, uErr
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside TGDBConnection::GetGraphMetadata successfully updated GraphMetadata"))
		}
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:GetGraphMetadata"))
	}
	return obj.graphObjFactory.GetGraphMetaData(), nil
}

// GetGraphObjectFactory gets the Graph Object Factory for Object creation
func (obj *TGDBConnection) GetGraphObjectFactory() (tgdb.TGGraphObjectFactory, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering TGDBConnection:GetGraphObjectFactory"))
	}
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetGraphObjectFactory - unable to initialize metadata w/ error: '%s'", err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:GetGraphObjectFactory"))
	}
	return obj.graphObjFactory, nil
}

// GetLargeObjectAsBytes gets an Binary Large Object Entity given an UniqueKey for the Object
func (obj *TGDBConnection) GetLargeObjectAsBytes(entityId int64, decryptFlag bool) ([]byte, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:GetLargeObjectAsBytes for EntityId: '%+v'", entityId))
	}
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetLargeObjectAsBytes - unable to initialize metadata w/ error: '%s'", err.Error()))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::GetLargeObjectAsBytes about to createChannelRequest() for: VerbGetLargeObjectRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbGetLargeObjectRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetLargeObjectAsBytes - unable to createChannelRequest(VerbGetLargeObjectRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*GetLargeObjectRequestMessage)
	queryRequest.SetEntityId(entityId)
	queryRequest.SetDecryption(decryptFlag)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside TGDBConnection::GetLargeObjectAsBytes about to obj.GetChannel().SendRequest() for: VerbGetLargeObjectRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetLargeObjectAsBytes - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TGDBConnection::GetLargeObjectAsBytes received response for: VerbGetLargeObjectRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*GetLargeObjectResponseMessage)

	if response == nil {
		errMsg := "TGDBConnection::GetLargeObjectAsBytes does not have any results in GetLargeObjectResponseMessage"
		logger.Error(errMsg)
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning TGDBConnection:GetLargeObjectAsBytes w/ '%+v'", response.GetBuffer()))
	}
	return response.GetBuffer(), nil
}

// GetRemovedList gets a list of removed entities
func (obj *TGDBConnection) GetRemovedList() map[int64]tgdb.TGEntity {
	return obj.removedList
}

// InsertEntity marks an ENTITY for insert operation. Upon commit, the entity will be inserted in the database
func (obj *TGDBConnection) InsertEntity(entity tgdb.TGEntity) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:InsertEntity to insert Entity: '%+v'", entity.GetEntityType()))
	}
	obj.addedList[entity.GetVirtualId()] = entity
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:InsertEntity"))
	}
	return nil
}

// Rollback rolls back the current transaction on this connection
func (obj *TGDBConnection) Rollback() tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering TGDBConnection:Rollback"))
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	// Reset all the lists to empty contents
	obj.addedList = make(map[int64]tgdb.TGEntity, 0)
	obj.changedList = make(map[int64]tgdb.TGEntity, 0)
	obj.removedList = make(map[int64]tgdb.TGEntity, 0)
	obj.attrByTypeList = make(map[int][]tgdb.TGAttribute, 0)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:Rollback"))
	}
	return nil
}

// SetExceptionListener sets exception listener
func (obj *TGDBConnection) SetExceptionListener(listener tgdb.TGConnectionExceptionListener) {
	obj.connPoolImpl.SetExceptionListener(listener) //delegate it to the Pool.
}

// UpdateEntity marks an ENTITY for update operation. Upon commit, the entity will be updated in the database
// When commit is called, the object is resolved to check if it is dirty. Entity.setAttribute calls make the entity
// dirty. If it is dirty, then the object is send to the server for update, otherwise it is ignored.
// Calling multiple times, does not change the behavior.
// The same entity cannot be updated on multiple connections. It will result an TGException of already associated to a connection.
func (obj *TGDBConnection) UpdateEntity(entity tgdb.TGEntity) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:UpdateEntity to update Entity: '%+v'", entity))
	}
	obj.changedList[entity.GetVirtualId()] = entity
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:UpdateEntity"))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChangeListener
/////////////////////////////////////////////////////////////////

// AttributeAdded gets called when an attribute is Added to an entity.
func (obj *TGDBConnection) AttributeAdded(attr tgdb.TGAttribute, owner tgdb.TGEntity) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:AttributeAdded"))
	}
}

// AttributeChanged gets called when an attribute is set.
func (obj *TGDBConnection) AttributeChanged(attr tgdb.TGAttribute, oldValue, newValue interface{}) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:AttributeChanged"))
	}
}

// AttributeRemoved gets called when an attribute is removed from the entity.
func (obj *TGDBConnection) AttributeRemoved(attr tgdb.TGAttribute, owner tgdb.TGEntity) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:AttributeRemoved"))
	}
}

// EntityCreated gets called when an entity is Added
func (obj *TGDBConnection) EntityCreated(entity tgdb.TGEntity) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:EntityCreated to add Entity: '%+v'", entity))
	}
	entityId := entity.(*AbstractEntity).GetVirtualId()
	obj.addedList[entityId] = entity
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:EntityCreated"))
	}
}

// EntityDeleted gets called when the entity is deleted
func (obj *TGDBConnection) EntityDeleted(entity tgdb.TGEntity) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGDBConnection:EntityDeleted to delete Entity: '%+v'", entity))
	}
	entityId := entity.(*AbstractEntity).GetVirtualId()
	obj.removedList[entityId] = entity
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning TGDBConnection:EntityDeleted"))
	}
}

// NodeAdded gets called when a node is Added
func (obj *TGDBConnection) NodeAdded(graph tgdb.TGGraph, node tgdb.TGNode) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:NodeAdded to add Node: '%+v' to Graph: '%+v'", node, graph))
	}
	entityId := graph.(*Graph).GetVirtualId()
	obj.addedList[entityId] = graph
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:NodeAdded"))
	}
}

// NodeRemoved gets called when a node is removed
func (obj *TGDBConnection) NodeRemoved(graph tgdb.TGGraph, node tgdb.TGNode) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering TGDBConnection:NodeRemoved to remove Node: '%+v' to Graph: '%+v'", node, graph))
	}
	entityId := graph.(*Graph).GetVirtualId()
	obj.removedList[entityId] = graph
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning TGDBConnection:NodeRemoved"))
	}
}

func (obj *TGDBConnection) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TGDBConnection:{")
	buffer.WriteString(fmt.Sprintf("channel: %+v", obj.channel))
	buffer.WriteString(fmt.Sprintf(", connId: %d", obj.connId))
	//buffer.WriteString(fmt.Sprintf(", connPoolImpl: %+v", obj.connPoolImpl))
	buffer.WriteString(fmt.Sprintf(", GraphObjFactory: %+v", obj.graphObjFactory))
	buffer.WriteString(fmt.Sprintf(", connProperties: %+v", obj.connProperties))
	buffer.WriteString(fmt.Sprintf(", addedList: %d", obj.addedList))
	buffer.WriteString(fmt.Sprintf(", changedList: %d", obj.changedList))
	buffer.WriteString(fmt.Sprintf(", removedList: %+v", obj.removedList))
	buffer.WriteString(fmt.Sprintf(", attrByTypeList: %+v", obj.attrByTypeList))
	buffer.WriteString("}")
	return buffer.String()
}


// ======= Various Connection Pool States =======
const (
	ConnectionPoolInitialized = iota
	ConnectionPoolConnecting
	ConnectionPoolConnected
	ConnectionPoolInUse
	ConnectionPoolDisconnecting
	ConnectionPoolDisconnected
	ConnectionPoolStopped
)

// 0 :      Indefinite
// -1 :     Immediate
// &gt; :   That many seconds
const (
	Immediate  = "-1"
	Indefinite = "0"
	IMMEDIATE  = 0
	INFINITE   = math.MaxInt32 // math.MaxInt64
)

type ConnectionPoolImpl struct {
	adminLock             sync.RWMutex // rw-lock for synchronizing read-n-update of connection pool properties
	connectReserveTimeOut time.Duration
	connList              []tgdb.TGConnection // Total Available Connections (Active + Dead/ToBeReused)
	connType              tgdb.TypeConnection
	chanPool              chan tgdb.TGConnection
	poolProperties        tgdb.TGProperties
	consumers             map[int64]tgdb.TGConnection        // Active/In-Use Connections
	exceptionListener     tgdb.TGConnectionExceptionListener // Function Pointer
	poolSize              int
	poolState             int
	useDedicateChannel    bool
}

var gInstance *ConnectionPoolImpl
var once sync.Once

func defaultTGConnectionPool() *ConnectionPoolImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ConnectionPoolImpl{})

	//once.Do(func() {
	gInstance := &ConnectionPoolImpl{
		connList:  make([]tgdb.TGConnection, 0),
		connType:  tgdb.TypeConventional,
		consumers: make(map[int64]tgdb.TGConnection, 0),
	}
	gInstance.poolSize, _ = strconv.Atoi(GetConfigFromKey(ConnectionPoolDefaultPoolSize).GetDefaultValue())
	gInstance.useDedicateChannel, _ = strconv.ParseBool(GetConfigFromKey(ConnectionPoolUseDedicatedChannelPerConnection).GetDefaultValue())
	//})
	return gInstance
}

func NewTGConnectionPool(url tgdb.TGChannelUrl, poolSize int, props *SortedProperties, connType tgdb.TypeConnection) *ConnectionPoolImpl {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering ConnectionPoolImpl:NewTGConnectionPool w/ ChannelURL: '%+v', Poolsize: '%d'", url.GetUrlAsString(), poolSize))
	}
	cp := defaultTGConnectionPool()
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside ConnectionPoolImpl:NewTGConnectionPool w/ Default Connection Pool: '%s'", cp.String()))
	}
	cp.connType = connType
	cp.poolProperties = props
	cp.chanPool = make(chan tgdb.TGConnection, poolSize+2)
	cp.poolSize = poolSize
	timeoutStr := GetConfigFromKey(ConnectionReserveTimeoutSeconds).GetDefaultValue()
	if timeoutStr == Immediate {
		cp.connectReserveTimeOut = time.Second * IMMEDIATE
	} else if timeoutStr == Indefinite {
		cp.connectReserveTimeOut = time.Second * INFINITE
	} else {
		timeout, _ := strconv.Atoi(timeoutStr)
		cp.connectReserveTimeOut = time.Second * time.Duration(timeout)
	}
	var ch tgdb.TGChannel
	channelFactory := GetChannelFactoryInstance()
	for i := 0; i < cp.poolSize; i++ {
		if ch == nil || cp.useDedicateChannel {
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside ConnectionPoolImpl:NewTGConnectionPool - about to channelFactory.CreateChannelWithUrlProperties() for URL: '%s'", url.GetUrlAsString()))
			}
			// Create a channel from channel factory
			channel1, err := channelFactory.CreateChannelWithUrlProperties(url, props)
			if err != nil {
				errMsg := fmt.Sprintf("ERROR: Returning ConnectionPoolImpl:NewTGConnectionPool Unable to create a channel for URL: '%s' via channel factory - '%+v'", url, err.Error())
				logger.Error(errMsg)
				continue
			}
			ch = channel1
		}
		//logger.Debug(fmt.Sprintf("Inside ConnectionPoolImpl:NewTGConnectionPool - about to NewTGDBConnection() for Channel: '%+v' and Properties: '%+v'", ch, props))
		// Create a connection
		var conn tgdb.TGConnection
		switch connType {
		case tgdb.TypeConventional:
			conn = NewTGDBConnection(cp, ch, props)
		case tgdb.TypeAdmin:
			conn = NewAdminConnection(cp, ch, props)
		default:
			conn = NewTGDBConnection(cp, ch, props)
		}
		conn.SetConnectionPool(cp)
		conn.SetConnectionProperties(props)
		// Add it in the pool to initialize the pool with a set number of initialized connections
		cp.connList = append(cp.connList, conn)
		//logger.Debug(fmt.Sprintf("Inside ConnectionPoolImpl:NewTGConnectionPool - about to add conn: '%+v' to the pool", conn.String()))
		cp.chanPool <- conn
	} // End of For loop for Pool Size
	cp.poolState = ConnectionPoolInitialized
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ConnectionPoolImpl:NewTGConnectionPool w/ Connection Pool: '%s'", cp.String()))
	}
	return cp
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGConnectionPool
/////////////////////////////////////////////////////////////////

// GetConnectionList returns all the connections = Active/In-Use + Un-used/Initialized
func (obj *ConnectionPoolImpl) GetConnectionList() []tgdb.TGConnection {
	return obj.connList
}

// GetConnectionProperties returns all the connection properties
func (obj *ConnectionPoolImpl) GetConnectionProperties() tgdb.TGProperties {
	return obj.poolProperties
}

// GetActiveConnections returns all the Active/In-Use connections
func (obj *ConnectionPoolImpl) GetActiveConnections() map[int64]tgdb.TGConnection {
	return obj.consumers
}

// GetNoOfActiveConnections returns the count of Active/In-Use connections
func (obj *ConnectionPoolImpl) GetNoOfActiveConnections() int {
	return len(obj.consumers)
}

// GetPoolState returns current state of the connection pool
func (obj *ConnectionPoolImpl) GetPoolState() int {
	return obj.poolState
}

// GetConnection returns an available connection from the pool that is NOT being used or nil if timeout elapses
func (obj *ConnectionPoolImpl) GetConnection() (tgdb.TGConnection, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering ConnectionPoolImpl:GetConnection for Pool Type: '%+v'", obj.connType))
	}
	obj.adminLock.Lock()
	defer obj.adminLock.Unlock()

	if obj.GetNoOfActiveConnections() == obj.GetPoolSize() {
		logger.Error(fmt.Sprint("ERROR: Returning ConnectionPoolImpl:GetConnection - as all the connections in the pool are in use."))
		errMsg := "ConnectionPoolImpl has already exhausted its limit. All the connections in the pool are in use. Please wait and retry."
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside ConnectionPoolImpl:GetConnection - about to pull connection from obj.ConnPool"))
	}
	var conn tgdb.TGConnection
	// Search through all available connections within the pooled set to find out which one is free to service next
	select {
	case conn = <-obj.chanPool:
		//logger.Debug(fmt.Sprintf("Inside ConnectionPoolImpl::GetConnection about to verify the state of the connection: '%+v' pulled from obj.ConnPool:", conn))
		// Proceed with the connection that is NOT being used already a.k.a. part of consumer connection map
		if _, ok := obj.consumers[conn.GetConnectionId()]; !ok {
			// Add it in the pool to initialize the pool with a set number of initialized connections
			obj.consumers[conn.GetConnectionId()] = conn
			break
		}
	case <-time.After(obj.connectReserveTimeOut):
		logger.Warning(fmt.Sprintf("WARNING: Returning ConnectionPoolImpl:GetConnection - trying to get a connection after wating for '%+v'", obj.connectReserveTimeOut))
		//log.Warning("Timed out trying to get connection from the pool")
		errMsg := fmt.Sprintf("Timed out trying to get a connection after wating for %v", obj.connectReserveTimeOut)
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ConnectionPoolImpl:GetConnection w/ Connection: '%+v'", conn))
	}
	return conn, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConnectionPool
/////////////////////////////////////////////////////////////////

// AdminLock locks the connection pool so that the list of connections can be updated
func (obj *ConnectionPoolImpl) AdminLock() {
	obj.adminLock.RLock()
}

// AdminUnlock unlocks the connection pool so that the list of connections can be updated
func (obj *ConnectionPoolImpl) AdminUnlock() {
	obj.adminLock.RUnlock()
}

// Connect establishes connection from this pool of available/configured connections to the TGDB server
// Exception could be BadAuthentication or BadUrl
func (obj *ConnectionPoolImpl) Connect() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering ConnectionPoolImpl:Connect"))
	}
	obj.adminLock.RLock()
	defer obj.adminLock.RUnlock()

	if obj.poolState == ConnectionPoolConnected {
		logger.Error(fmt.Sprint("ERROR: Returning ConnectionPoolImpl:Connect - ConnectionPoolImpl is already connected. Disconnect and then reconnect."))
		errMsg := "ConnectionPoolImpl is already connected. Disconnect and then reconnect"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}
	// Set the state to connecting
	obj.poolState = ConnectionPoolConnecting

	// Attempt to connect using each of the available connections in the pool
	for i := 0; i < len(obj.connList); i++ {
		conn := obj.connList[i]


		// Proceed when the connection is either in Initialized OR Disconnected (OR Stopped???) state
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside ConnectionPoolImpl::Connect Consumer Loop about to conn.Connect() using connection: '%+v'", conn))
		}
		err := conn.Connect()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ConnectionPoolImpl:Connect - unable to conn.Connect() w/ '%s'", err.Error()))
			// TODO: Revisit later - Decide what error to throw and whether to continue / break
			return err
		}
	}

	// Set the state to connecting
	obj.poolState = ConnectionPoolConnected
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning ConnectionPoolImpl:Connect"))
	}
	return nil
}

// Disconnect breaks the connection from the TGDB server and returns the connection back to this connection pool for reuse
func (obj *ConnectionPoolImpl) Disconnect() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering ConnectionPoolImpl:Disconnect"))
	}
	obj.adminLock.RLock()
	defer obj.adminLock.RUnlock()

	if obj.poolState != ConnectionPoolConnected {
		logger.Error(fmt.Sprint("ERROR: Returning ConnectionPoolImpl:Disconnect - ConnectionPoolImpl is NOT connected."))
		errMsg := fmt.Sprintf("ConnectionPoolImpl is not connected. State is: %d", obj.poolState)
		return GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}
	// Set the state to connecting
	obj.poolState = ConnectionPoolDisconnecting

	// Attempt to connect using each of the active connections in the pool
	for _, conn := range obj.connList {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside ConnectionPoolImpl::Disconnect active connection Loop about to conn.Disconnect() using connection: '%+v'", conn))
		}
		err := conn.Disconnect()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ConnectionPoolImpl:Disconnect - unable to conn.Disconnect() w/ '%s'", err.Error()))
			// TODO: Revisit later - Decide what error to throw and whether to continue / break
			return err
		}
	}

	// Set the state to connecting
	obj.poolState = ConnectionPoolDisconnected
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning ConnectionPoolImpl:Disconnect"))
	}
	return nil
}

// Get a free connection
// The property ConnectionReserveTimeoutSeconds or tgdb.connectionpool.ConnectionReserveTimeoutSeconds specifies the time
// to wait in seconds. It has the following meaning
// 0 :      Indefinite
// -1 :     Immediate
// &gt; :   That many seconds
func (obj *ConnectionPoolImpl) Get() (tgdb.TGConnection, tgdb.TGError) {
	return obj.GetConnection()
}

// GetPoolSize gets pool size
func (obj *ConnectionPoolImpl) GetPoolSize() int {
	return obj.poolSize
}

// ReleaseConnection frees the connection and sends back to the pool
func (obj *ConnectionPoolImpl) ReleaseConnection(conn tgdb.TGConnection) (tgdb.TGConnectionPool, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering ConnectionPoolImpl:ReleaseConnection"))
	}
	obj.adminLock.RLock()
	defer obj.adminLock.RUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside ConnectionPoolImpl::ReleaseConnection Consumer Loop about to remove connection from consumer list"))
	}

	tgConn, ok := conn.(*TGDBConnection)
	if ok {
		delete(obj.consumers, tgConn.GetConnectionId())
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside ConnectionPoolImpl::ReleaseConnection Consumer Loop about to return connection to obj.ConnPool"))
		}
		obj.chanPool <- tgConn
	} else {
		tgAdminConn, ok := conn.(*AdminConnectionImpl)
		if ok {
			delete(obj.consumers, tgAdminConn.TGDBConnection.GetConnectionId())
			if logger.IsDebug() {
							logger.Debug(fmt.Sprint("Inside ConnectionPoolImpl::ReleaseConnection Consumer Loop about to return connection to obj.ConnPool"))
			}
			obj.chanPool <- tgAdminConn
		}
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning ConnectionPoolImpl:ReleaseConnection"))
	}
	return obj, nil
}

// SetExceptionListener sets exception listener
func (obj *ConnectionPoolImpl) SetExceptionListener(listener tgdb.TGConnectionExceptionListener) {
	obj.exceptionListener = listener
}

func (obj *ConnectionPoolImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ConnectionPoolImpl:{")
	buffer.WriteString(fmt.Sprintf("ConnectReserveTimeOut: %+v", obj.connectReserveTimeOut))
	buffer.WriteString(fmt.Sprintf(", ConnType: %+v", obj.connType))
	buffer.WriteString(fmt.Sprintf(", ConnList: %+v", obj.connList))
	buffer.WriteString(fmt.Sprintf(", ConnPool: %+v", obj.chanPool))
	//buffer.WriteString(fmt.Sprintf(", PoolProperties: %+v", obj.poolProperties))
	buffer.WriteString(fmt.Sprintf(", Consumers: %+v", obj.consumers))
	buffer.WriteString(fmt.Sprintf(", PoolSize: %d", obj.poolSize))
	buffer.WriteString(fmt.Sprintf(", PoolState: %d", obj.poolState))
	buffer.WriteString(fmt.Sprintf(", UseDedicateChannel: %+v", obj.useDedicateChannel))
	buffer.WriteString("}")
	return buffer.String()
}





type AdminConnectionImpl struct {
	*TGDBConnection
	//channel         types.TGChannel
	//connId          int64
	//connPoolImpl    types.TGConnectionPool    // Connection belongs to a connection pool
	//graphObjFactory *model.GraphObjectFactory // Intentionally kept private to ensure execution of InitMetaData() before accessing graph objects
	//connProperties  types.TGProperties
	//addedList       map[int64]types.TGEntity
	//changedList     map[int64]types.TGEntity
	//removedList     map[int64]types.TGEntity
	//attrByTypeList  map[int][]types.TGAttribute
}

func DefaultAdminConnection() *AdminConnectionImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AdminConnectionImpl{})

	newSGDBConnection := &AdminConnectionImpl{
		TGDBConnection: DefaultTGDBConnection(),
		//addedList:      make(map[int64]types.TGEntity, 0),
		//changedList:    make(map[int64]types.TGEntity, 0),
		//removedList:    make(map[int64]types.TGEntity, 0),
		//attrByTypeList: make(map[int][]types.TGAttribute, 0),
	}
	////newSGDBConnection.channel = DefaultAbstractChannel()
	//newSGDBConnection.connId = atomic.AddInt64(&connectionIds, 1)
	//// We cannot get meta data before we connect to the server
	//newSGDBConnection.graphObjFactory = model.NewGraphObjectFactory(newSGDBConnection)
	//newSGDBConnection.connProperties = utils.NewSortedProperties()
	return newSGDBConnection
}

func NewAdminConnection(conPool *ConnectionPoolImpl, channel tgdb.TGChannel, props tgdb.TGProperties) *AdminConnectionImpl {
	newSGDBConnection := DefaultAdminConnection()
	newSGDBConnection.connPoolImpl = conPool
	newSGDBConnection.channel = channel
	newSGDBConnection.connProperties = props.(*SortedProperties)
	return newSGDBConnection
}

/////////////////////////////////////////////////////////////////
// Helper functions for AdminConnectionImpl
/////////////////////////////////////////////////////////////////

func (obj *AdminConnectionImpl) GetConnectionPool() tgdb.TGConnectionPool {
	return obj.connPoolImpl
}

func (obj *AdminConnectionImpl) GetConnectionGraphObjectFactory() *GraphObjectFactory {
	return obj.graphObjFactory
}

func (obj *AdminConnectionImpl) InitMetadata() tgdb.TGError {
	if obj.graphObjFactory == nil {
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil
	}
	gmd := obj.graphObjFactory.GetGraphMetaData()
	if gmd.IsInitialized() {
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil
	}

	// Update the metadata and retrieve it fresh
	_, err := obj.GetGraphMetadata(true)
	if err != nil {
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return err
	}
	return nil
}

// SetConnectionPool sets connection pool
func (obj *AdminConnectionImpl) SetConnectionPool(connPool tgdb.TGConnectionPool) {
	obj.connPoolImpl = connPool
}

// SetConnectionProperties sets connection properties
func (obj *AdminConnectionImpl) SetConnectionProperties(connProps tgdb.TGProperties) {
	obj.connProperties = connProps.(*SortedProperties)
}

/////////////////////////////////////////////////////////////////
// Private functions for AdminConnectionImpl
/////////////////////////////////////////////////////////////////

func configureAdminRequest(adminReq *AdminRequestMessage, option interface{}) {
	if adminReq == nil || option == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl:configureAdminRequest as getReq == nil || reqProps == nil"))
		return
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:configureAdminRequest w/ AdminRequest: '%+v'", adminReq.String()))
	}
	switch adminReq.GetCommand() {
	case AdminCommandKillConnection:
		adminReq.SetSessionId(option.(int64))
	case AdminCommandSetLogLevel:
		adminReq.SetLogLevel(option.(*ServerLogDetails))
	case AdminCommandShowAttrDescs:
	default:
		break
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:configureAdminRequest w/ AdminRequest: '%+v'", adminReq.String()))
	}
	return
}

// executeAdminRequest executes an immediate command with associated parameters
func (obj *AdminConnectionImpl) executeAdminRequest(command AdminCommand, parameters interface{}) (interface{}, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteAdminRequest for Admin Command: '%+v'", command))
	}
	//err := obj.InitMetadata()
	//if err != nil {
	//	return nil, err
	//}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteAdminRequest about to createChannelRequest() for: pdu.VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbAdminRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteAdminRequest - unable to createChannelRequest(pdu.VerbAdminRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	adminRequest := msgRequest.(*AdminRequestMessage)
	adminRequest.SetCommand(command)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteAdminRequest about to obj.configureQueryRequest() for: AdminRequestMessage"))
	}
	configureAdminRequest(adminRequest, parameters)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteAdminRequest about to obj.GetChannel().SendRequest() for: pdu.VerbAdminRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(adminRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteAdminRequest - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteAdminRequest received response for: pdu.VerbAdminRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*AdminResponseMessage)

	//if !response.GetHasResult() {
	//	logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::ExecuteAdminRequest - The query does not have any results in AdminResponseMessage"))
	//	return nil, nil
	//}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:ExecuteAdminRequest w/ '%+v'", response))
	}
	return obj.populateResultSetFromAdminResponse(command, response)
}

func (obj *AdminConnectionImpl) populateResultSetFromAdminResponse(command AdminCommand, msgResponse *AdminResponseMessage) (interface{}, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:populateResultSetFromAdminResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	}
	//if !msgResponse.GetHasResult() {
	//	logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromAdminResponse as msgResponse does not have any results"))
	//	errMsg := "AdminConnectionImpl::populateResultSetFromAdminResponse does not have any results"
	//	return nil, exception.GetErrorByType(rf.TGErrorGeneralException, "", errMsg, "")
	//}

	var results interface{}

	switch command {
	case AdminCommandShowAttrDescs:
		results = msgResponse.GetDescriptorList()
	case AdminCommandShowConnections:
		results = msgResponse.GetConnectionList()
	case AdminCommandShowIndices:
		results = msgResponse.GetIndexList()
	case AdminCommandShowInfo:
		results = msgResponse.GetServerInfo()
	case AdminCommandShowUsers:
		results = msgResponse.GetUserList()
	case AdminCommandSetLogLevel:
		fallthrough
	case AdminCommandStopServer:
		fallthrough
	default:
		break
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:populateResultSetFromAdminResponse w/ Results: '%+v'", results))
	}
	return results, nil
}

func (obj *AdminConnectionImpl) populateResultSetFromQueryResponse(resultId int, msgResponse *QueryResponseMessage) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:populateResultSetFromQueryResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	}
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse as msgResponse does not have any results"))
		errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse does not have any results in QueryResponseMessage"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]tgdb.TGEntity, 0)
	var rSet *ResultSet

	currResultCount := 0
	resultCount := msgResponse.GetResultCount()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse read resultCount: '%d' FetchedEntityCount: '%d'", resultCount, len(fetchedEntities)))
	}
	if resultCount > 0 {
		respStream.SetReferenceMap(fetchedEntities)
		rSet = NewResultSet(obj, resultId)
	}

	totalCount := msgResponse.GetTotalCount()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse read totalCount: '%d'", totalCount))
	}
	for i := 0; i < totalCount; i++ {
		entityType, err := respStream.(*ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse - unable to read entityType in the response stream"))
			errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to read entity type in the response stream"
			return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		kindId := tgdb.TGEntityKind(entityType)
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse read #'%d'-entityType: '%+v', kindId: '%s'", i, entityType, kindId.String()))
		}
		if kindId != tgdb.EntityKindInvalid {
			entityId, err := respStream.(*ProtocolDataInputStream).ReadLong()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse - unable to read entityId in the response stream"))
				errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to read entity type in the response stream"
				return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			entity := fetchedEntities[entityId]
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse read entityId: '%d', kindId: '%s', entity: '%+v'", entityId, kindId.String(), entity))
			}
			switch kindId {
			case tgdb.EntityKindNode:
				var node Node
				if entity == nil {
					node, nErr := obj.graphObjFactory.CreateNode()
					if nErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to create a new node from the response stream"
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.Error())
					}
					entity = node
					fetchedEntities[entityId] = node
					if currResultCount < resultCount {
						rSet.AddEntityToResultSet(node)
					}
					if logger.IsDebug() {
						logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse created new node: '%+v' FetchedEntityCount: '%d'", node, len(fetchedEntities)))
					}
				}
				node = *(entity.(*Node))
				err := node.ReadExternal(respStream)
				if err != nil {
					errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to node.ReadExternal() from the response stream"
					logger.Error(errMsg)
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.GetErrorDetails())
				}
				//logger.Debug(fmt.Sprintf("======> ======> After node.ReadExternal() FetchedEntityCount: '%d'", len(fetchedEntities)))
				if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("======> ======> Node w/ Edges: '%+v'\n", node.GetEdges()))
				}
			case tgdb.EntityKindEdge:
				var edge Edge
				if entity == nil {
					//edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, DirectionTypeBiDirectional)
					edge, eErr := obj.graphObjFactory.CreateEntity(tgdb.EntityKindEdge)
					if eErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to create a new bi-directional edge from the response stream"
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, eErr.Error())
					}
					entity = edge
					fetchedEntities[entityId] = edge
					if currResultCount < resultCount {
						rSet.AddEntityToResultSet(edge)
					}
					if logger.IsDebug() {
						logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse created new edge: '%+v' FetchedEntityCount: '%d'", edge, len(fetchedEntities)))
					}
				}
				edge = *(entity.(*Edge))
				err := edge.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromQueryResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				//logger.Debug(fmt.Sprintf("======> ======> After edge.ReadExternal() FetchedEntityCount: '%d'", len(fetchedEntities)))
				if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("======> ======> Edge w/ Vertices: '%+v'\n", edge.GetVertices()))
				}
			case tgdb.EntityKindGraph:
				// TODO: Revisit later - Should we break after throwing/logging an error
				continue
			}
			//if entity != nil {
			//	logger.Debug(fmt.Sprintf("======> AdminConnectionImpl::populateResultSetFromQueryResponse entityId: '%d', kindId: '%d', entityType: '%+v'\n", entityId, kindId, kindId.String()))
			//	attrList, _ := entity.GetAttributes()
			//	for _, attrib := range attrList {
			//		logger.Debug(fmt.Sprintf("======> Attribute Value: '%+v'\n", attrib.GetValue()))
			//	}
			//	if kindId == types.EntityKindNode {
			//		edges := (entity.(*model.Node)).GetEdges()
			//		logger.Debug(fmt.Sprintf("======> Node w/ Edges: '%+v'\n", edges))
			//	} else if kindId == types.EntityKindEdge {
			//		vertices := (entity.(*model.Edge)).GetVertices()
			//		logger.Debug(fmt.Sprintf("======> Edge w/ Vertices: '%+v'\n", vertices))
			//	}
			//}
		} else {
			logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:populateResultSetFromQueryResponse - Received invalid entity kind %d", kindId))
		} // Valid entity types
	} // End of for loop
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:populateResultSetFromQueryResponse w/ ResultSet: '%+v'", rSet))
	}
	return rSet, nil
}

func (obj *AdminConnectionImpl) populateResultSetFromGetEntitiesResponse(msgResponse *GetEntityResponseMessage) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:populateResultSetFromGetEntitiesResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	}
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse as msgResponse does not have any results"))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse does not have any results"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]tgdb.TGEntity, 0)

	totalCount, err := respStream.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read totalCount in the response stream w/ error: '%s'", err.Error()))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read total count of entities in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse extracted totalCount: '%d'", totalCount))
	}
	if totalCount > 0 {
		respStream.SetReferenceMap(fetchedEntities)
	}

	rSet := NewResultSet(obj, msgResponse.GetResultId())
	// Number of entities matches the search.  Exclude the related entities
	currResultCount := 0
	_, err = respStream.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read count of result entities in the response stream w/ error: '%s'", err.Error()))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read count of result entities in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}

	for i := 0; i < totalCount; i++ {
		isResult, err := respStream.(*ProtocolDataInputStream).ReadBoolean()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read isResult in the response stream w/ error: '%s'", err.Error()))
			errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read count of result entities in the response stream"
			return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse read isResult: '%+v'", isResult))
		}
		entityType, err := respStream.(*ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
			errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read entity type in the response stream"
			return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		kindId := tgdb.TGEntityKind(entityType)
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse extracted entityType: '%+v', kindId: '%d'", entityType, kindId))
		}
		if kindId != tgdb.EntityKindInvalid {
			entityId, err := respStream.(*ProtocolDataInputStream).ReadLong()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read entityId in the response stream w/ error: '%s'", err.Error()))
				errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read entity type in the response stream"
				return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			entity := fetchedEntities[entityId]
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse extracted entityId: '%d', entity: '%+v'", entityId, entity))
			}
			switch kindId {
			case tgdb.EntityKindNode:
				var node Node
				if entity == nil {
					node, nErr := obj.graphObjFactory.CreateNode()
					if nErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to create a new node from the response stream"
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.GetErrorDetails())
					}
					entity = node
					fetchedEntities[entityId] = node
					if logger.IsDebug() {
											logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse created new node: '%+v'", node))
					}
				}
				node = *(entity.(*Node))
				err := node.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				if isResult {
					rSet.AddEntityToResultSet(entity)
					currResultCount++
				}
			case tgdb.EntityKindEdge:
				var edge Edge
				if entity == nil {
					edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, tgdb.DirectionTypeBiDirectional)
					if eErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
						// TODO: Revisit later - Should we break after throwing/logging an error
						//continue
						errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to create a new bi-directional edge from the response stream"
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, eErr.Error())
					}
					entity = edge
					fetchedEntities[entityId] = edge
					if logger.IsDebug() {
											logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse created new edge: '%+v'", edge))
					}
				}
				edge = *(entity.(*Edge))
				err := edge.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				if isResult {
					rSet.AddEntityToResultSet(entity)
					currResultCount++
				}
			case tgdb.EntityKindGraph:
				// TODO: Revisit later - Should we break after throwing/logging an error
				continue
			}
		} else {
			logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - Received invalid entity kind %d", kindId))
		} // Valid entity types
	} // End of for loop
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse w/ ResultSe: '%+v'", rSet))
	}
	return rSet, nil
}

func (obj *AdminConnectionImpl) populateResultSetFromGetEntityResponse(msgResponse *GetEntityResponseMessage) (tgdb.TGEntity, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:populateResultSetFromGetEntityResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	}
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse as msgResponse does not have any results"))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse does not have any results"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]tgdb.TGEntity, 0)

	var entityFound tgdb.TGEntity
	respStream.SetReferenceMap(fetchedEntities)

	count, err := respStream.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to read count in the response stream w/ error: '%s'", err.Error()))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to read count of entities in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse extracted Count: '%d'", count))
	}
	if count > 0 {
		respStream.SetReferenceMap(fetchedEntities)
		for i := 0; i < count; i++ {
			entityType, err := respStream.(*ProtocolDataInputStream).ReadByte()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
				errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to read entity type in the response stream"
				return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			kindId := tgdb.TGEntityKind(entityType)
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse extracted entityType: '%+v', kindId: '%d'", entityType, kindId))
			}
			if kindId != tgdb.EntityKindInvalid {
				entityId, err := respStream.(*ProtocolDataInputStream).ReadLong()
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to read entityId in the response stream w/ error: '%s'", err.Error()))
					errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to read entity type in the response stream"
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				entity := fetchedEntities[entityId]
				if logger.IsDebug() {
									logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse extracted entityId: '%d', entity: '%+v'", entityId, entity))
				}
				switch kindId {
				case tgdb.EntityKindNode:
					// Need to put shell object into map to be deserialized later
					var node Node
					if entity == nil {
						node, nErr := obj.graphObjFactory.CreateNode()
						if nErr != nil {
							logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
							// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
							//continue
							errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to create a new node from the response stream"
							return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.Error())
						}
						entity = node
						fetchedEntities[entityId] = node
						if entityFound == nil {
							entityFound = node
						}
						if logger.IsDebug() {
													logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse created new node: '%+v'", entity))
						}
					}
					node = *(entity.(*Node))
					err := node.ReadExternal(respStream)
					if err != nil {
						errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
						logger.Error(errMsg)
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
					}
				case tgdb.EntityKindEdge:
					var edge Edge
					if entity == nil {
						edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, tgdb.DirectionTypeBiDirectional)
						if eErr != nil {
							logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
							// TODO: Revisit later - Should we break after throwing/logging an error
							//continue
							errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to create a new bi-directional edge from the response stream"
							return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, eErr.Error())
						}
						entity = edge
						fetchedEntities[entityId] = edge
						if entityFound == nil {
							entityFound = edge
						}
						if logger.IsDebug() {
													logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse created new edge: '%+v'", edge))
						}
					}
					edge = *(entity.(*Edge))
					err := edge.ReadExternal(respStream)
					if err != nil {
						errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
						logger.Error(errMsg)
						return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
					}
				case tgdb.EntityKindGraph:
					// TODO: Revisit later - Should we break after throwing/logging an error
					continue
				}
			} else {
				logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:populateResultSetFromGetEntityResponse - Received invalid entity kind %d", kindId))
			} // Valid entity types
		} // End of for loop
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse w/ Entity: '%+v'", entityFound))
	}
	return entityFound, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAdminConnection
/////////////////////////////////////////////////////////////////

// CheckpointServer allows the programmatic control to do a checkpoint on server
func (obj *AdminConnectionImpl) CheckpointServer() tgdb.TGError {
	_, err := obj.executeAdminRequest(AdminCommandCheckpointServer, nil)
	if err != nil {
		return err
	}
	return nil
}

// DumpServerStackTrace prints the stack trace
func (obj *AdminConnectionImpl) DumpServerStackTrace() tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AdminConnectionImpl:DumpServerStackTrace for Admin Command: DumpServerStackTrace"))
	}
	//err := obj.InitMetadata()
	//if err != nil {
	//	return nil, err
	//}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::DumpServerStackTrace about to createChannelRequest() for: pdu.DumpServerStackTrace"))
	}
	// Create a channel request
	msgRequest, _, cErr := createChannelRequest(obj, VerbDumpStacktraceRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:DumpServerStackTrace - unable to createChannelRequest(pdu.VerbDumpStacktraceRequest w/ error: '%s'", cErr.Error()))
		return cErr
	}
	dumpStackRequest := msgRequest.(*DumpStacktraceRequestMessage)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::DumpServerStackTrace about to obj.GetChannel().SendRequest() for: pdu.VerbDumpStacktraceRequest"))
	}
	// Execute request on channel and get the response
	channelErr := obj.GetChannel().SendMessage(dumpStackRequest)
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:DumpServerStackTrace - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return channelErr
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:DumpServerStackTrace for Admin Command: DumpServerStackTrace"))
	}
	return nil
}

// GetAttributeDescriptors gets the list of attribute descriptors
func (obj *AdminConnectionImpl) GetAttributeDescriptors() ([]tgdb.TGAttributeDescriptor, tgdb.TGError) {
	results, err := obj.executeAdminRequest(AdminCommandShowAttrDescs, nil)
	if err != nil {
		return nil, err
	}
	return results.([]tgdb.TGAttributeDescriptor), nil
}

// GetConnections gets the list of all socket connections using this connection type
func (obj *AdminConnectionImpl) GetConnections() ([]tgdb.TGConnectionInfo, tgdb.TGError) {
	results, err := obj.executeAdminRequest(AdminCommandShowConnections, nil)
	if err != nil {
		return nil, err
	}
	return results.([]tgdb.TGConnectionInfo), nil
}

// GetIndices gets the list of all indices
func (obj *AdminConnectionImpl) GetIndices() ([]tgdb.TGIndexInfo, tgdb.TGError) {
	results, err := obj.executeAdminRequest(AdminCommandShowIndices, nil)
	if err != nil {
		return nil, err
	}
	return results.([]tgdb.TGIndexInfo), nil
}

// GetInfo gets the information about this connection type
func (obj *AdminConnectionImpl) GetInfo() (tgdb.TGServerInfo, tgdb.TGError) {
	results, err := obj.executeAdminRequest(AdminCommandShowInfo, nil)
	if err != nil {
		return nil, err
	}
	return results.(tgdb.TGServerInfo), nil
}

// GetUsers gets the list of users
func (obj *AdminConnectionImpl) GetUsers() ([]tgdb.TGUserInfo, tgdb.TGError) {
	results, err := obj.executeAdminRequest(AdminCommandShowUsers, nil)
	if err != nil {
		return nil, err
	}
	return results.([]tgdb.TGUserInfo), nil
}

// KillConnection terminates the connection forcefully
func (obj *AdminConnectionImpl) KillConnection(sessionId int64) tgdb.TGError {
	_, err := obj.executeAdminRequest(AdminCommandKillConnection, sessionId)
	if err != nil {
		return err
	}
	return nil
}

// SetServerLogLevel set the log level
func (obj *AdminConnectionImpl) SetServerLogLevel(logLevel int, logComponent int64) tgdb.TGError {
	logDetails := NewServerLogDetails(TGLogLevel(logLevel), TGLogComponent(logComponent))
	_, err := obj.executeAdminRequest(AdminCommandSetLogLevel, logDetails)
	if err != nil {
		return err
	}
	return nil
}

// StopServer stops the admin connection
func (obj *AdminConnectionImpl) StopServer() tgdb.TGError {
	_, err := obj.executeAdminRequest(AdminCommandStopServer, nil)
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConnection
/////////////////////////////////////////////////////////////////

// Commit commits the current transaction on this connection
func (obj *AdminConnectionImpl) Commit() (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AdminConnectionImpl:Commit"))
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Commit - about to loop through addedList to include existing nodes to the changed list if it's part of a new edge"))
	}
	// Include existing nodes to the changed list if it's part of a new edge
	for _, addEntity := range obj.GetAddedList() {
		if addEntity.GetEntityKind() == tgdb.EntityKindEdge {
			nodes := addEntity.(*Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*Node)
					if !node.GetIsNew() {
						obj.changedList[node.GetVirtualId()] = node
						logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:Commit - Existing node '%d' added to change list for a new edge", node.GetVirtualId()))
					}
				}
			}
		}
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Commit - about to loop through changedList to include existing nodes to the changed list even for edge update"))
	}
	// Need to include existing node to the changed list even for edge update
	for _, modEntity := range obj.GetChangedList() {
		if modEntity.GetEntityKind() == tgdb.EntityKindEdge {
			nodes := modEntity.(*Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*Node)
					if !node.GetIsNew() {
						obj.changedList[node.GetVirtualId()] = node
						logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:Commit - Existing node '%d' added to change list for an existing edge '%d'", node.GetVirtualId(), modEntity.GetVirtualId()))
					}
				}
			}
		}
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Commit - about to loop through removedList to include existing nodes to the changed list even for edge update"))
	}
	// Need to include existing node to the changed list even for edge update
	for _, delEntity := range obj.GetRemovedList() {
		if delEntity.GetEntityKind() == tgdb.EntityKindEdge {
			nodes := delEntity.(*Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*Node)
					if !node.GetIsNew() {
						if obj.removedList[node.GetVirtualId()] == nil {
							obj.changedList[node.GetVirtualId()] = node
							logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:Commit - Existing node '%d' added to change list for an edge %d to be deleted", node.GetVirtualId(), delEntity.GetVirtualId()))
						}
					}
				}
			}
		}
	}
	//For deleted edge and node, we don't immediately change the effected nodes or edges.

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Commit about to createChannelRequest() for: pdu.VerbCommitTransactionRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, VerbCommitTransactionRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:Commit - unable to createChannelRequest(pdu.VerbCommitTransactionRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*CommitTransactionRequest)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Commit about to GetAttributeDescriptors() for: pdu.VerbCommitTransactionRequest"))
	}
	gof := obj.graphObjFactory
	attrDescSet, aErr := gof.GetGraphMetaData().GetNewAttributeDescriptors()
	if aErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:Commit - unable to gmd.GetAttributeDescriptors() w/ error: '%s'", aErr.Error()))
		return nil, aErr
	}
	queryRequest.AddCommitLists(obj.addedList, obj.changedList, obj.removedList, attrDescSet)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Commit about to channelSendRequest() for: pdu.VerbCommitTransactionRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:Commit - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::Commit received response for: pdu.VerbCommitTransactionRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*CommitTransactionResponse)

	if response.HasException() {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:Commit as response has exceptions"))
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil, nil
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Commit about to fixUpAttrDescriptors()"))
	}
	fixUpAttrDescriptors(response, attrDescSet)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Commit about to obj.fixUpEntities()"))
	}
	fixUpEntities(obj, response)

	for _, delEntity := range obj.GetRemovedList() {
		delEntity.SetIsDeleted(true)
	}

	// Reset and clear all the lists
	for _, modEntity := range obj.GetChangedList() {
		modEntity.ResetModifiedAttributes()
	}
	for _, newEntity := range obj.GetAddedList() {
		newEntity.ResetModifiedAttributes()
	}

	obj.addedList = make(map[int64]tgdb.TGEntity, 0)
	obj.changedList = make(map[int64]tgdb.TGEntity, 0)
	obj.removedList = make(map[int64]tgdb.TGEntity, 0)
	obj.attrByTypeList = make(map[int][]tgdb.TGAttribute, 0)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:Commit"))
	}
	return nil, nil
}

// Connect establishes a network connection to the TGDB server
func (obj *AdminConnectionImpl) Connect() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:Connect for connection: '%+v'", obj))
	}
	err := obj.GetChannel().Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl::Connect - error in obj.GetChannel().Connect() as '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Connect about to obj.GetChannel().Start()"))
	}
	err = obj.GetChannel().Start()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl::Connect - error in obj.GetChannel().Start() as '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:Connect"))
	}
	return nil
}

// CloseQuery closes a specific query and associated objects
func (obj *AdminConnectionImpl) CloseQuery(queryHashId int64) (tgdb.TGQuery, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:CloseQuery for QueryHashId: '%+v'", queryHashId))
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::CloseQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, VerbQueryRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:CloseQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(CLOSE)
	queryRequest.SetQueryHashId(queryHashId)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::CloseQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	_, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:CloseQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:CloseQuery"))
	}
	return nil, nil
}

// CreateQuery creates a reusable query object that can be used to execute one or more statement
func (obj *AdminConnectionImpl) CreateQuery(expr string) (tgdb.TGQuery, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:CreateQuery for Query: '%+v'", expr))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::CreateQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, VerbQueryRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:CreateQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(CREATE)
	queryRequest.SetQuery(expr)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::CreateQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:CreateQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::CreateQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	}

	response := msgResponse.(*QueryResponseMessage)
	queryHashId := response.GetQueryHashId()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:CreateQuery"))
	}
	if response.GetResult() == 0 && queryHashId > 0 {
		return NewQuery(obj, queryHashId), nil
	}
	return nil, nil
}

// DecryptBuffer decrypts the encrypted buffer by sending a DecryptBufferRequest to the server
func (obj *AdminConnectionImpl) DecryptBuffer(is tgdb.TGInputStream) ([]byte, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AdminConnectionImpl:DecryptBuffer w/ EncryptedBuffer"))
	}
	//err := obj.InitMetadata()
	//if err != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:DecryptBuffer - unable to initialize metadata w/ error: '%s'", err.Error()))
	//	return nil, err
	//}
	//obj.connPoolImpl.AdminLock()
	//defer obj.connPoolImpl.AdminUnlock()
	//
	//logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::DecryptBuffer about to createChannelRequest() for: pdu.VerbDecryptBufferRequest"))
	//// Create a channel request
	//msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbDecryptBufferRequest)
	//if cErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:DecryptBuffer - unable to createChannelRequest(pdu.VerbDecryptBufferRequest w/ error: '%s'", cErr.Error()))
	//	return nil, cErr
	//}
	//decryptRequest := msgRequest.(*pdu.DecryptBufferRequestMessage)
	//decryptRequest.SetEncryptedBuffer(encryptedBuf)
	//
	//logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::DecryptBuffer about to obj.GetChannel().SendRequest() for: pdu.VerbDecryptBufferRequest"))
	//// Execute request on channel and get the response
	//msgResponse, channelErr := obj.GetChannel().SendRequest(decryptRequest, channelResponse.(*channel.BlockingChannelResponse))
	//if channelErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:DecryptBuffer - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
	//	return nil, channelErr
	//}
	//logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::DecryptBuffer received response for: pdu.VerbGetLargeObjectRequest as '%+v'", msgResponse))
	//response := msgResponse.(*pdu.DecryptBufferResponseMessage)
	//
	//if response == nil {
	//	errMsg := "AdminConnectionImpl::DecryptBuffer does not have any results in GetLargeObjectResponseMessage"
	//	logger.Error(errMsg)
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	//}
	//
	//logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:DecryptBuffer w/ '%+v'", response.GetDecryptedBuffer()))
	//return response.GetDecryptedBuffer(), nil
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Decrypt(is)
}

// DecryptEntity decrypts the encrypted entity using channel's data cryptographer
func (obj *AdminConnectionImpl) DecryptEntity(entityId int64) ([]byte, tgdb.TGError) {
	buf, err := obj.GetLargeObjectAsBytes(entityId, true)
	if err != nil {
		return nil, err
	}
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Decrypt(NewProtocolDataInputStream(buf))
}

// DeleteEntity marks an ENTITY for delete operation. Upon commit, the entity will be deleted from the database
func (obj *AdminConnectionImpl) DeleteEntity(entity tgdb.TGEntity) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:DeleteEntity for Entity: '%+v'", entity))
	}
	obj.removedList[entity.GetVirtualId()] = entity
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:DeleteEntity"))
	}
	return nil
}

// Disconnect breaks the connection from the TGDB server
func (obj *AdminConnectionImpl) Disconnect() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AdminConnectionImpl:Disconnect"))
	}
	err := obj.GetChannel().Disconnect()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:Disconnect - unable to channel.Disconnect() w/ error: '%s'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::Disconnect about to obj.GetChannel().Stop()"))
	}
	obj.GetChannel().Stop(false)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:Disconnect"))
	}
	return nil
}

// EncryptEntity encrypts the encrypted entity using channel's data cryptographer
func (obj *AdminConnectionImpl) EncryptEntity(rawBuffer []byte) ([]byte, tgdb.TGError) {
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Encrypt(rawBuffer)
}

// ExecuteGremlinQuery executes a Gremlin Grammer-Based query with  query options
func (obj *AdminConnectionImpl) ExecuteGremlinQuery(expr string, collection []interface{}, options tgdb.TGQueryOption) ([]interface{}, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteGremlinQuery for Query: '%+v'", expr))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteGremlinQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteGremlinQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery("gbc : " + expr)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteGremlinQuery about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	}
	configureQueryRequest(queryRequest, options)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteGremlinQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteGremlinQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteGremlinQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::ExecuteGremlinQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	// TODO: Revisit later once Gremlin Package and GremlinQueryResult are implemented
	//respStream := response.GetEntityStream()
	//GremlinResult.fillCollection(entityStream, gof, collection);
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:ExecuteGremlinQuery w/ '%+v'", response))
	}
	return collection, nil
}

//
// The query needs to be on Connection object instead of AdminConnection
//
// ExecuteQuery executes an immediate query with associated query options

func (obj *AdminConnectionImpl) ExecuteQuery(expr string, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
	return obj.TGDBConnection.ExecuteQuery(expr, options)
}

//func (obj *AdminConnectionImpl) ExecuteQuery(expr string, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
//	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteQuery for Query: '%+v'", expr))
//	err := obj.InitMetadata()
//	if err != nil {
//		return nil, err
//	}
//	obj.connPoolImpl.AdminLock()
//	defer obj.connPoolImpl.AdminUnlock()
//
//	logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
//	// Create a channel request
//	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
//	if cErr != nil {
//		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
//		return nil, cErr
//	}
//	queryRequest := msgRequest.(*QueryRequestMessage)
//	queryRequest.SetCommand(EXECUTE)
//	queryRequest.SetQuery(expr)
//	logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQuery about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
//	configureQueryRequest(queryRequest, options)
//
//	logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
//	// Execute request on channel and get the response
//	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
//	if channelErr != nil {
//		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
//		return nil, channelErr
//	}
//	logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
//	response := msgResponse.(*QueryResponseMessage)
//
//	if !response.GetHasResult() {
//		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::ExecuteQuery - The query does not have any results in QueryResponseMessage"))
//		return nil, nil
//	}
//
//	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:ExecuteQuery w/ '%+v'", response))
//	return obj.populateResultSetFromQueryResponse(0, response)
//}

// ExecuteQueryWithFilter executes an immediate query with specified filter & query options
// The query option is place holder at this time
// @param expr A subset of SQL-92 where clause
// @param edgeFilter filter used for selecting edges to be returned
// @param traversalCondition condition used for selecting edges to be traversed and returned
// @param endCondition condition used to stop the traversal
// @param option Query options for executing. Can be null, then it will use the default option
func (obj *AdminConnectionImpl) ExecuteQueryWithFilter(expr, edgeFilter, traversalCondition, endCondition string, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteQueryWithFilter for Query: '%+v', EdgeFilter: '%+v', Traversal: '%+v', EndCondition: '%+v'", expr, edgeFilter, traversalCondition, endCondition))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithFilter about to createChannelRequest() for: pdu.VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQueryWithFilter - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery(expr)
	queryRequest.SetEdgeFilter(edgeFilter)
	queryRequest.SetTraversalCondition(traversalCondition)
	queryRequest.SetEndCondition(endCondition)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteQueryWithFilter about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	}
	configureQueryRequest(queryRequest, options)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithFilter about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQueryWithFilter - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteQueryWithFilter received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::ExecuteQueryWithFilter - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:ExecuteQueryWithFilter w/ '%+v'", response))
	}
	return obj.populateResultSetFromQueryResponse(0, response)
}

// ExecuteQueryWithId executes an immediate query for specified id & query options
func (obj *AdminConnectionImpl) ExecuteQueryWithId(queryHashId int64, options tgdb.TGQueryOption) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteQueryWithId for QueryHashId: '%+v'", queryHashId))
	}
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithId about to createChannelRequest() for: pdu.VerbQueryRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQueryWithId - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*QueryRequestMessage)
	queryRequest.SetCommand(EXECUTED)
	queryRequest.SetQueryHashId(queryHashId)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithId about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	}
	configureQueryRequest(queryRequest, options)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithId about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	}
	// Execute request on channel and get the response
	_, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQueryWithId - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:ExecuteQueryWithId"))
	}
	return nil, nil
}

// GetAddedList gets a list of added entities
func (obj *AdminConnectionImpl) GetAddedList() map[int64]tgdb.TGEntity {
	return obj.addedList
}

// GetChangedList gets a list of changed entities
func (obj *AdminConnectionImpl) GetChangedList() map[int64]tgdb.TGEntity {
	return obj.changedList
}

// GetChangedList gets the communication channel associated with this connection
func (obj *AdminConnectionImpl) GetChannel() tgdb.TGChannel {
	return obj.channel
}

// GetConnectionId gets connection identifier
func (obj *AdminConnectionImpl) GetConnectionId() int64 {
	return obj.connId
}

// GetConnectionProperties gets a list of connection properties
func (obj *AdminConnectionImpl) GetConnectionProperties() tgdb.TGProperties {
	return obj.connProperties
}

// GetEntities gets a result set of entities given an non-uniqueKey
func (obj *AdminConnectionImpl) GetEntities(qryKey tgdb.TGKey, props tgdb.TGProperties) (tgdb.TGResultSet, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:GetEntities for QueryKey: '%+v'", qryKey))
	}
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:GetEntities - unable to InitMetadata"))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if props == nil {
		props = NewQueryOption()
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetEntities about to createChannelRequest() for: pdu.VerbGetEntityRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbGetEntityRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetEntities - unable to createChannelRequest(pdu.VerbGetEntityRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*GetEntityRequestMessage)
	queryRequest.SetCommand(2)
	queryRequest.SetKey(qryKey)
	configureGetRequest(queryRequest, props)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetEntities about to obj.GetChannel().SendRequest() for: pdu.VerbGetEntityRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetEntities - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::GetEntities received response for: pdu.VerbGetEntityRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*GetEntityResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::GetEntities - The request does not have any results in GetEntityResponseMessage"))
		return nil, nil
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:GetEntities w/ '%+v'", response))
	}
	return obj.populateResultSetFromGetEntitiesResponse(response)
}

// GetEntity gets an Entity given an UniqueKey for the Object
func (obj *AdminConnectionImpl) GetEntity(qryKey tgdb.TGKey, options tgdb.TGQueryOption) (tgdb.TGEntity, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:GetEntity for QueryKey: '%+v'", qryKey))
	}
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:GetEntity - unable to InitMetadata"))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if options == nil {
		options = NewQueryOption()
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetEntity about to createChannelRequest() for: pdu.VerbGetEntityRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbGetEntityRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetEntity - unable to createChannelRequest(pdu.VerbGetEntityRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*GetEntityRequestMessage)
	queryRequest.SetCommand(0)
	queryRequest.SetKey(qryKey)
	configureGetRequest(queryRequest, options)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetEntity about to obj.GetChannel().SendRequest() for: pdu.VerbGetEntityRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetEntity - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::GetEntity received response for: pdu.VerbGetEntityRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*GetEntityResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::GetEntity - The request does not have any results in GetEntityResponseMessage"))
		return nil, nil
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:GetEntity w/ '%+v'", response))
	}
	return obj.populateResultSetFromGetEntityResponse(response)
}

// GetGraphMetadata gets the Graph Metadata
func (obj *AdminConnectionImpl) GetGraphMetadata(refresh bool) (tgdb.TGGraphMetadata, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AdminConnectionImpl:GetGraphMetadata"))
	}
	if refresh {
		obj.connPoolImpl.AdminLock()
		defer obj.connPoolImpl.AdminUnlock()

		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetGraphMetadata about to createChannelRequest() for: pdu.VerbMetadataRequest"))
		}
		// Create a channel request
		msgRequest, channelResponse, err := createChannelRequest(obj, VerbMetadataRequest)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetGraphMetadata - unable to createChannelRequest(pdu.VerbMetadataRequest w/ error: '%s'", err.Error()))
			return nil, err
		}
		metaRequest := msgRequest.(*MetadataRequest)
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::GetGraphMetadata createChannelRequest() returned MsgRequest: '%+v' ChannelResponse: '%+v'", msgRequest, channelResponse.(*BlockingChannelResponse)))
		}

		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetGraphMetadata about to obj.GetChannel().SendRequest() for: pdu.VerbMetadataRequest"))
		}
		// Execute request on channel and get the response
		msgResponse, channelErr := obj.GetChannel().SendRequest(metaRequest, channelResponse.(*BlockingChannelResponse))
		if channelErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetGraphMetadata - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
			return nil, channelErr
		}
		//logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::GetGraphMetadata received response for: pdu.VerbMetadataRequest as '%+v'", msgResponse))

		response := msgResponse.(*MetadataResponse)
		attrDescList := response.GetAttrDescList()
		edgeTypeList := response.GetEdgeTypeList()
		nodeTypeList := response.GetNodeTypeList()

		gmd := obj.graphObjFactory.GetGraphMetaData()
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetGraphMetadata about to update GraphMetadata"))
		}
		uErr := gmd.UpdateMetadata(attrDescList, nodeTypeList, edgeTypeList)
		if uErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetGraphMetadata - unable to gmd.UpdateMetadata() w/ error: '%s'", uErr.Error()))
			return nil, uErr
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetGraphMetadata successfully updated GraphMetadata"))
		}
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:GetGraphMetadata"))
	}
	return obj.graphObjFactory.GetGraphMetaData(), nil
}

// GetGraphObjectFactory gets the Graph Object Factory for Object creation
func (obj *AdminConnectionImpl) GetGraphObjectFactory() (tgdb.TGGraphObjectFactory, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AdminConnectionImpl:GetGraphObjectFactory"))
	}
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetGraphObjectFactory - unable to initialize metadata w/ error: '%s'", err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:GetGraphObjectFactory"))
	}
	return obj.graphObjFactory, nil
}

// GetLargeObjectAsBytes gets an Binary Large Object Entity given an UniqueKey for the Object
func (obj *AdminConnectionImpl) GetLargeObjectAsBytes(entityId int64, decryptFlag bool) ([]byte, tgdb.TGError) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:GetLargeObjectAsBytes for EntityId: '%+v'", entityId))
	}
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetLargeObjectAsBytes - unable to initialize metadata w/ error: '%s'", err.Error()))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetLargeObjectAsBytes about to createChannelRequest() for: pdu.VerbGetLargeObjectRequest"))
	}
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, VerbGetLargeObjectRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetLargeObjectAsBytes - unable to createChannelRequest(pdu.VerbGetLargeObjectRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*GetLargeObjectRequestMessage)
	queryRequest.SetEntityId(entityId)
	queryRequest.SetDecryption(decryptFlag)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AdminConnectionImpl::GetLargeObjectAsBytes about to obj.GetChannel().SendRequest() for: pdu.VerbGetLargeObjectRequest"))
	}
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetLargeObjectAsBytes - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminConnectionImpl::GetLargeObjectAsBytes received response for: pdu.VerbGetLargeObjectRequest as '%+v'", msgResponse))
	}
	response := msgResponse.(*GetLargeObjectResponseMessage)

	if response == nil {
		errMsg := "AdminConnectionImpl::GetLargeObjectAsBytes does not have any results in GetLargeObjectResponseMessage"
		logger.Error(errMsg)
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AdminConnectionImpl:GetLargeObjectAsBytes w/ '%+v'", response.GetBuffer()))
	}
	return response.GetBuffer(), nil
}

// GetRemovedList gets a list of removed entities
func (obj *AdminConnectionImpl) GetRemovedList() map[int64]tgdb.TGEntity {
	return obj.removedList
}

// InsertEntity marks an ENTITY for insert operation. Upon commit, the entity will be inserted in the database
func (obj *AdminConnectionImpl) InsertEntity(entity tgdb.TGEntity) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:InsertEntity to insert Entity: '%+v'", entity.GetEntityType()))
	}
	obj.addedList[entity.GetVirtualId()] = entity
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:InsertEntity"))
	}
	return nil
}

// Rollback rolls back the current transaction on this connection
func (obj *AdminConnectionImpl) Rollback() tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AdminConnectionImpl:Rollback"))
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	// Reset all the lists to empty contents
	obj.addedList = make(map[int64]tgdb.TGEntity, 0)
	obj.changedList = make(map[int64]tgdb.TGEntity, 0)
	obj.removedList = make(map[int64]tgdb.TGEntity, 0)
	obj.attrByTypeList = make(map[int][]tgdb.TGAttribute, 0)

	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:Rollback"))
	}
	return nil
}

// SetExceptionListener sets exception listener
func (obj *AdminConnectionImpl) SetExceptionListener(listener tgdb.TGConnectionExceptionListener) {
	obj.connPoolImpl.SetExceptionListener(listener) //delegate it to the Pool.
}

// UpdateEntity marks an ENTITY for update operation. Upon commit, the entity will be updated in the database
// When commit is called, the object is resolved to check if it is dirty. Entity.setAttribute calls make the entity
// dirty. If it is dirty, then the object is send to the server for update, otherwise it is ignored.
// Calling multiple times, does not change the behavior.
// The same entity cannot be updated on multiple connections. It will result an TGException of already associated to a connection.
func (obj *AdminConnectionImpl) UpdateEntity(entity tgdb.TGEntity) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:UpdateEntity to update Entity: '%+v'", entity))
	}
	obj.changedList[entity.GetVirtualId()] = entity
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:UpdateEntity"))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChangeListener
/////////////////////////////////////////////////////////////////

// AttributeAdded gets called when an attribute is Added to an entity.
func (obj *AdminConnectionImpl) AttributeAdded(attr tgdb.TGAttribute, owner tgdb.TGEntity) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:AttributeAdded"))
	}
}

// AttributeChanged gets called when an attribute is set.
func (obj *AdminConnectionImpl) AttributeChanged(attr tgdb.TGAttribute, oldValue, newValue interface{}) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:AttributeChanged"))
	}
}

// AttributeRemoved gets called when an attribute is removed from the entity.
func (obj *AdminConnectionImpl) AttributeRemoved(attr tgdb.TGAttribute, owner tgdb.TGEntity) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:AttributeRemoved"))
	}
}

// EntityCreated gets called when an entity is Added
func (obj *AdminConnectionImpl) EntityCreated(entity tgdb.TGEntity) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:EntityCreated to add Entity: '%+v'", entity))
	}
	entityId := entity.(*AbstractEntity).GetVirtualId()
	obj.addedList[entityId] = entity
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:EntityCreated"))
	}
}

// EntityDeleted gets called when the entity is deleted
func (obj *AdminConnectionImpl) EntityDeleted(entity tgdb.TGEntity) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:EntityDeleted to delete Entity: '%+v'", entity))
	}
	entityId := entity.(*AbstractEntity).GetVirtualId()
	obj.removedList[entityId] = entity
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:EntityDeleted"))
	}
}

// NodeAdded gets called when a node is Added
func (obj *AdminConnectionImpl) NodeAdded(graph tgdb.TGGraph, node tgdb.TGNode) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:NodeAdded to add Node: '%+v' to Graph: '%+v'", node, graph))
	}
	entityId := graph.(*Graph).GetVirtualId()
	obj.addedList[entityId] = graph
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:NodeAdded"))
	}
}

// NodeRemoved gets called when a node is removed
func (obj *AdminConnectionImpl) NodeRemoved(graph tgdb.TGGraph, node tgdb.TGNode) {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AdminConnectionImpl:NodeRemoved to remove Node: '%+v' to Graph: '%+v'", node, graph))
	}
	entityId := graph.(*Graph).GetVirtualId()
	obj.removedList[entityId] = graph
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AdminConnectionImpl:NodeRemoved"))
	}
}

func (obj *AdminConnectionImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AdminConnectionImpl:{")
	buffer.WriteString(fmt.Sprintf("channel: %+v", obj.channel))
	buffer.WriteString(fmt.Sprintf(", connId: %d", obj.connId))
	//buffer.WriteString(fmt.Sprintf(", connPoolImpl: %+v", obj.connPoolImpl))
	buffer.WriteString(fmt.Sprintf(", GraphObjFactory: %+v", obj.graphObjFactory))
	buffer.WriteString(fmt.Sprintf(", connProperties: %+v", obj.connProperties))
	buffer.WriteString(fmt.Sprintf(", addedList: %d", obj.addedList))
	buffer.WriteString(fmt.Sprintf(", changedList: %d", obj.changedList))
	buffer.WriteString(fmt.Sprintf(", removedList: %+v", obj.removedList))
	buffer.WriteString(fmt.Sprintf(", attrByTypeList: %+v", obj.attrByTypeList))
	buffer.WriteString("}")
	return buffer.String()
}




// Create a connection on the url using the Name and password
// Each connection will create a dedicated Channel for connection.
// @param url The url for connection.  A URL is represented as a string of the form <BR>
//            &lt;protocol&gt;://[user@]['['ipv6']'] | ipv4 [:][port][/]'{' Name:value;... '}' <BR>
//            protocol can be tcp or ssl.<BR>
//            For ssl protocol, the following properties must also be present <BR>
//            dbName:the database Name <BR>
//            verifyDBName:true<BR>
//
// @param userName The user Name for connection. The userId provided overrides all other userIds that can be infered.
//                 The rules for overriding are in this order<BR>
//                 a. The argument 'userId' is the highest priority. If Null then <BR>
//                 b. The user@url is considered. If that is Null <BR>
//                 c. the "userID=value" from the URL string is considered.<BR>
//                 d. If all of them is Null, then the default User associated to the installation will be taken.<BR>
//
// @param password The managled or unmanagled password
//
// @param env optional environment. This environment will override every other environment values infered, and is
// <table>
//  <caption>Connection Properties</caption>
// 	<thead>
// 		<tr>
// 			<th style="width:auto;text-align:left">Full Name</th>
// 			<th style="width:auto;text-align:left">Alias</th>
// 			<th style="width:auto;text-align:left">Default Value</th>
// 			<th style="width:auto;text-align:left">Description</th>
// 		</tr>
// </thead>
// 	<tbody>
// 		<tr>
// 			<td>tgdb.channel.defaultHost</td>
// 			<td>-</td>
// 			<td>localhost</td>
// 			<td>The default host specifier</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.defaultPort</td>
// 			<td>-</td>
// 			<td>8700</td>
// 			<td>The default port specifier</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.defaultProtocol</td>
// 			<td>-</td>
// 			<td>tcp</td>
// 			<td>The default protocol</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.sendSize</td>
// 			<td>sendSize</td>
// 			<td>122</td>
// 			<td>TCP send packet size in KBs</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.recvSize</td>
// 			<td>recvSize</td>
// 			<td>128</td>
// 			<td>TCP recv packet size in KB</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.pingInterval</td>
// 			<td>pingInterval</td>
// 			<td>30</td>
// 			<td>Keep alive ping intervals</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.connectTimeout</td>
// 			<td>connectTimeout</td>
// 			<td>1000</td>
// 			<td>Timeout for connection to establish, before it gives up and tries the ftUrls if specified</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.ftHosts</td>
// 			<td>ftHosts</td>
// 			<td>-</td>
// 			<td>Alternate fault tolerant list of &lt;host:port&gt; pair seperated by comma. </td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.ftRetryIntervalSeconds</td>
// 			<td>ftRetryIntervalSeconds</td>
// 			<td>10</td>
// 			<td>The connect retry interval to ftHosts</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.ftRetryCount</td>
// 			<td>ftRetryCount</td>
// 			<td>3</td>
// 			<td>The number of times ro retry </td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.defaultUserID</td>
// 			<td>-</td>
// 			<td>-</td>
// 			<td>The default user Id for the connection</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.userID</td>
// 			<td>userID</td>
// 			<td>-</td>
// 			<td>The user id for the connection if it is not specified in the API. See the rules for picking the user Name</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.password</td>
// 			<td>password</td>
// 			<td>-</td>
// 			<td>The password for the username</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.channel.clientId</td>
// 			<td>clientId</td>
// 			<td>tgdb.java-api.client</td>
// 			<td>The client id to be used for the connection</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.connection.dbName</td>
// 			<td>dbName</td>
// 			<td>-</td>
// 			<td>The database Name the client is connecting to. It is used as part of verification for ssl channels</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.connectionpool.useDedicatedChannelPerConnection</td>
// 			<td>useDedicatedChannelPerConnection</td>
// 			<td>false</td>
// 			<td>A boolean value indicating either to multiplex mulitple connections on a single tcp socket or use dedicate socket per connection. A true value consumes resource but provides good performance. Also check the max number of connections</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.connectionpool.defaultPoolSize</td>
// 			<td>defaultPoolSize</td>
// 			<td>10</td>
// 			<td>The default connection pool size to use when creating a ConnectionPoolImpl</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.connectionpool.connectionReserveTimeoutSeconds</td>
// 			<td>connectionReserveTimeoutSeconds</td>
// 			<td>10</td>
// 			<td>A timeout parameter indicating how long to wait before getting a connection from the pool</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.connection.operationTimeoutSeconds</td>
// 			<td>connectionOperationTimeoutSeconds</td>
// 			<td>10</td>
// 			<td>A timeout parameter indicating how long to wait for a operation before giving up. Some queries are long running, and may override this behavior.</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.tls.provider.Name</td>
// 			<td>tlsProviderName</td>
// 			<td>SunJSSE</td>
// 			<td>Transport level Security provider. Work with your InfoSec team to change this value</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.tls.provider.className</td>
// 			<td>tlsProviderClassName</td>
// 			<td>com.sun.net.ssl.internal.ssl.Provider</td>
// 			<td>The underlying Provider implementation. Work with your InfoSec team to change this value.</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.tls.provider.configFile</td>
// 			<td>tlsProviderConfigFile</td>
// 			<td>-</td>
// 			<td>Some providers require extra configuration paramters, and it can be passed as a file</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.tls.protocol</td>
// 			<td>tlsProtocol</td>
// 			<td>TLSv1.2</td>
// 			<td>tlsProtocol version. The system only supports 1.2+</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.tls.cipherSuites</td>
// 			<td>cipherSuites</td>
// 			<td>-</td>
// 			<td>A list cipher suites that the InfoSec team has cleared. The default list is a common list of JSSE's cipher list and Openssl list that supports 1.2 protocol </td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.tls.verifyDBName</td>
// 			<td>verifyDBName</td>
// 			<td>false</td>
// 			<td>Verify the Database Name in the certificate. TGDB provides self signed certificate for easy-to-use SSL.</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.tls.expectedHostName</td>
// 			<td>expectedHostName</td>
// 			<td>-</td>
// 			<td>The expected hostName for the certificate. This is for future use</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.tls.trustedCertificates</td>
// 			<td>trustedCertificates</td>
// 			<td>-</td>
//  			<td>The list of trusted Certificates</td>
// 		</tr>
// 		<tr>
// 			<td>tgdb.security.keyStorePassword</td>
// 			<td>keyStorePassword</td>
// 			<td>-</td>
// 			<td>The Keystore for the password</td>
// 		</tr>
// 	</tbody>
// </table>
//



type TransactionImpl struct {
	transactionId int64
}

// Make sure that the transactionImpl implements the TGTransaction interface
var _ tgdb.TGTransaction = (*TransactionImpl)(nil)

func DefaultTransaction() *TransactionImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TransactionImpl{})

	return &TransactionImpl{transactionId: -1}
}

func NewTransaction(txnId int64) *TransactionImpl {
	newTxn := DefaultTransaction()
	newTxn.transactionId = txnId
	return newTxn
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGTransaction
/////////////////////////////////////////////////////////////////

func (obj *TransactionImpl) GetTransactionId() int64 {
	return obj.transactionId
}

func (obj *TransactionImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TransactionImpl:{")
	buffer.WriteString(fmt.Sprintf("TransactionId: '%d'", obj.transactionId))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *TransactionImpl) ReadExternal(iStream tgdb.TGInputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *TransactionImpl) WriteExternal(oStream tgdb.TGOutputStream) tgdb.TGError {
	// No-op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TransactionImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.transactionId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TransactionImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TransactionImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.transactionId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TransactionImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type TGConnectionFactoryImpl struct {
}

// @return TGConnection - an instance of connection to the server with a dedicated channel
// @throws com.tibco.tgdb.exception.TGException - If it cannot create a connection to the server successfully
func (obj *TGConnectionFactoryImpl) CreateConnection(url, user, pwd string, env map[string]string) (tgdb.TGConnection, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGConnectionFactory:CreateConnection w/ URL: '%s' User: '%s', and Pwd: '%s'", url, user, pwd))
	}
	connPool, err := obj.CreateConnectionPool(url, user, pwd, 1, env)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGConnectionFactory:CreateConnection - unable to create connection pool - '%+v", err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGConnectionFactory:CreateConnection - about to connPool.Get()"))
	}
	conn, err := connPool.Get()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGConnectionFactory:CreateConnection - unable to connPool.Get() - '%+v", err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGConnectionFactory:CreateConnection w/ connection '%+v'",  conn))
	}
	return conn, nil
}

// Create an admin connection on the url using the Name and password
func (obj *TGConnectionFactoryImpl) CreateAdminConnection(url, user, pwd string, env map[string]string) (tgdb.TGConnection, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGConnectionFactory:CreateAdminConnection w/ URL: '%s' User: '%s', and Pwd: '%s'", url, user, pwd))
	}
	connPool, err := obj.CreateConnectionPoolWithType(url, user, pwd, 1, env, tgdb.TypeAdmin)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGConnectionFactory:CreateAdminConnection - unable to create connection pool - '%+v", err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TGConnectionFactory:CreateAdminConnection - about to connPool.Get()"))
	}
	conn, err := connPool.Get()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGConnectionFactory:CreateAdminConnection - unable to connPool.Get() - '%+v", err.Error()))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGConnectionFactory:CreateAdminConnection w/ connection '%+v'",  conn))
	}
	return conn, nil
}

// Create a connection Pool of pool size on the the url using the Name and password for a specific type of connections.
func (obj *TGConnectionFactoryImpl) CreateConnectionPool(url, user, pwd string, poolSize int, env map[string]string) (tgdb.TGConnectionPool, tgdb.TGError) {
	return obj.CreateConnectionPoolWithType(url, user, pwd, poolSize, env, tgdb.TypeConventional)
}

// Create a connection Pool of pool size on the the url using the Name and password. Each connection in the pool will default
// use a shared channel, but this can be overridden by setting the value property tgdb.connectionpool.useDedicatedChannel=true
// @param url The url for the channel used in the connection pool.
// @param userName  The user Name for connection.
// @param password  The password encrypted or un-encrypted
// @param poolSize the size of the pool
// @param env optional environment. This environment will override every other environment values infered, and is specific for this pool only
// @return A Connection Pool
// @throws com.tibco.tgdb.exception.TGException - If it cannot create a connection pool to the server successfully
func (obj *TGConnectionFactoryImpl) CreateConnectionPoolWithType(url, user, pwd string, poolSize int, env map[string]string, connType tgdb.TypeConnection) (tgdb.TGConnectionPool, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TGConnectionFactory:CreateConnectionPool w/ URL: '%s' User: '%s', Pwd: '%s' and Environment: '%+v'", url, user, pwd, env))
	}
	if poolSize <= 0 {
		poolSize = NewTGEnvironment().GetConnectionPoolDefaultPoolSize()
	}
	props := NewSortedProperties()
	// Add system default environment properties
	defEnvProperties := NewTGEnvironment().GetAsSortedProperties()
	for _, kvp := range defEnvProperties.(*SortedProperties).GetAllProperties() {
		props.AddProperty(kvp.KeyName, kvp.KeyValue)
	}
	// Add supplied environment values
	if env != nil {
		for n, v := range env {
			props.AddProperty(n, v)
		}
	}

	channelUrl := ParseChannelUrl(url)
	channelProps := channelUrl.GetProperties()
	// Add channel properties to the consolidated property set
	if channelProps != nil {
		for _, kvp := range channelProps.(*SortedProperties).GetAllProperties() {
			props.AddProperty(kvp.KeyName, kvp.KeyValue)
		}
	}
	_ = SetUserAndPassword(props, user, pwd) // Ignore Error Handling
	// At this point, the consolidated property set is already sorted by key a.k.a. property Name
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TGConnectionFactory:CreateConnectionPool about to initiate NewTGConnectionPool() for URL: '%+v'",  channelUrl.String()))
	}
	return NewTGConnectionPool(channelUrl, poolSize, props, connType), nil
}





//var logger = impl.DefaultTGLogManager().GetLogger()

//var globalConnectionFactory *TGConnectionFactoryImpl
var globalConnectionFactory tgdb.TGConnectionFactory
var sOnce sync.Once

//NewTGConnectionFactory

// This works
func NewTGConnectionFactory() tgdb.TGConnectionFactory {

	//func NewTGConnectionFactory() *TGConnectionFactory {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.

	gob.Register(TGConnectionFactoryImpl{})

	sOnce.Do(func() {
		globalConnectionFactory = &TGConnectionFactoryImpl{}
	})
	return globalConnectionFactory
}
