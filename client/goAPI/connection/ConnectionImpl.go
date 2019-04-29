/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: ConnectionImpl.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

package connection

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/channel"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/pdu"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"strings"
	"sync/atomic"
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
	channel         types.TGChannel
	connId          int64
	connPoolImpl    types.TGConnectionPool    // Connection belongs to a connection pool
	graphObjFactory *model.GraphObjectFactory // Intentionally kept private to ensure execution of InitMetaData() before accessing graph objects
	connProperties  types.TGProperties
	addedList       map[int64]types.TGEntity
	changedList     map[int64]types.TGEntity
	removedList     map[int64]types.TGEntity
	attrByTypeList  map[int][]types.TGAttribute
}

func DefaultTGDBConnection() *TGDBConnection {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGDBConnection{})

	newSGDBConnection := &TGDBConnection{
		addedList:      make(map[int64]types.TGEntity, 0),
		changedList:    make(map[int64]types.TGEntity, 0),
		removedList:    make(map[int64]types.TGEntity, 0),
		attrByTypeList: make(map[int][]types.TGAttribute, 0),
	}
	//newSGDBConnection.channel = DefaultAbstractChannel()
	newSGDBConnection.connId = atomic.AddInt64(&connectionIds, 1)
	// We cannot get meta data before we connect to the server
	newSGDBConnection.graphObjFactory = model.NewGraphObjectFactory(newSGDBConnection)
	newSGDBConnection.connProperties = utils.NewSortedProperties()
	return newSGDBConnection
}

func NewTGDBConnection(conPool *ConnectionPoolImpl, channel types.TGChannel, props types.TGProperties) *TGDBConnection {
	newSGDBConnection := DefaultTGDBConnection()
	newSGDBConnection.connPoolImpl = conPool
	newSGDBConnection.channel = channel
	newSGDBConnection.connProperties = props.(*utils.SortedProperties)
	return newSGDBConnection
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGConnection
/////////////////////////////////////////////////////////////////

func (obj *TGDBConnection) GetConnectionPool() types.TGConnectionPool {
	return obj.connPoolImpl
}

func (obj *TGDBConnection) GetConnectionGraphObjectFactory() *model.GraphObjectFactory {
	return obj.graphObjFactory
}

func (obj *TGDBConnection) InitMetadata() types.TGError {
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
func (obj *TGDBConnection) SetConnectionPool(connPool types.TGConnectionPool) {
	obj.connPoolImpl = connPool
}

// SetConnectionProperties sets connection properties
func (obj *TGDBConnection) SetConnectionProperties(connProps types.TGProperties) {
	obj.connProperties = connProps.(*utils.SortedProperties)
}

/////////////////////////////////////////////////////////////////
// Private functions for types.TGConnection
/////////////////////////////////////////////////////////////////

func fixUpAttrDescriptors(response *pdu.CommitTransactionResponse, attrDescSet []types.TGAttributeDescriptor) {
	logger.Log(fmt.Sprint("Entering TGDBConnection:fixUpAttrDescriptors"))
	attrDescCount := response.GetAttrDescCount()
	attrDescIdList := response.GetAttrDescIdList()
	for i := 0; i < attrDescCount; i++ {
		tempId := attrDescIdList[(i * 2)]
		realId := attrDescIdList[((i * 2) + 1)]

		for _, attrDesc := range attrDescSet {
			desc := attrDesc.(*model.AttributeDescriptor)
			if attrDesc.GetAttributeId() == tempId {
				logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:fixUpAttrDescriptors - Replace descriptor: '%d' by '%d'", attrDesc.GetAttributeId(), realId))
				desc.SetAttributeId(realId)
				break
			}
		}
	}
	logger.Log(fmt.Sprint("Returning TGDBConnection:fixUpAttrDescriptors"))
}

func fixUpEntities(obj types.TGConnection, response *pdu.CommitTransactionResponse) {
	logger.Log(fmt.Sprint("Entering TGDBConnection:fixUpEntities"))
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
	logger.Log(fmt.Sprint("Returning TGDBConnection:fixUpEntities"))
}

func createChannelRequest(obj types.TGConnection, verb int) (types.TGMessage, types.TGChannelResponse, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:createChannelRequest for Verb: '%s'", pdu.GetVerb(verb).GetName()))
	cn := utils.GetConfigFromKey(utils.ConnectionOperationTimeoutSeconds)
	//logger.Log(fmt.Sprintf("Inside AbstractChannel::channelTryRepeatConnect config for ConnectionOperationTimeoutSeconds is '%+v", cn))
	timeout := obj.GetConnectionProperties().GetPropertyAsInt(cn)
	requestId := atomic.AddInt64(&requestIds, 1)

	//logger.Log(fmt.Sprint("Inside TGDBConnection::createChannelRequest about to create channel.NewBlockingChannelResponse()"))
	// Create a non-blocking channel response
	channelResponse := channel.NewBlockingChannelResponse(requestId, int64(timeout))

	// Use Message Factory method to create appropriate message structure (class) based on input type
	//msgRequest, err := pdu.CreateMessageForVerb(verb)
	msgRequest, err := pdu.CreateMessageWithToken(verb, obj.GetChannel().GetAuthToken(), obj.GetChannel().GetSessionId())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:createChannelRequest - CreateMessageForVerb failed for Verb: '%s' w/ '%+v'", pdu.GetVerb(verb).GetName(), err.Error()))
		return nil, nil, err
	}
	logger.Log(fmt.Sprintf("Returning TGDBConnection:createChannelRequest for Verb: '%s' w/ MessageRequest: '%+v' ChannelResponse: '%+v'", pdu.GetVerb(verb).GetName(), msgRequest, channelResponse))
	return msgRequest, channelResponse, nil
}

func configureGetRequest(getReq *pdu.GetEntityRequestMessage, reqProps types.TGProperties) {
	if getReq == nil || reqProps == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::configureGetRequest as getReq == nil || reqProps == nil"))
		return
	}
	logger.Log(fmt.Sprintf("Entering TGDBConnection:configureGetRequest w/ EntityRequest: '%+v'", getReq.String()))
	props := reqProps.(*query.TGQueryOptionImpl)
	getReq.SetFetchSize(props.GetPreFetchSize())
	getReq.SetTraversalDepth(props.GetTraversalDepth())
	getReq.SetEdgeLimit(props.GetEdgeLimit())
	getReq.SetBatchSize(props.GetBatchSize())
	logger.Log(fmt.Sprintf("Returning TGDBConnection:configureGetRequest w/ EntityRequest: '%+v'", getReq.String()))
	return
}

func configureQueryRequest(qryReq *pdu.QueryRequestMessage, reqProps types.TGProperties) {
	if qryReq == nil || reqProps == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection:configureQueryRequest as getReq == nil || reqProps == nil"))
		return
	}
	logger.Log(fmt.Sprintf("Entering TGDBConnection:configureQueryRequest w/ QueryRequest: '%+v'", qryReq.String()))
	props := reqProps.(*query.TGQueryOptionImpl)
	qryReq.SetFetchSize(props.GetPreFetchSize())
	qryReq.SetTraversalDepth(props.GetTraversalDepth())
	qryReq.SetEdgeLimit(props.GetEdgeLimit())
	qryReq.SetBatchSize(props.GetBatchSize())
	qryReq.SetSortAttrName(props.GetSortAttrName())
	qryReq.SetSortOrderDsc(props.IsSortOrderDsc())
	qryReq.SetSortResultLimit(props.GetSortResultLimit())
	logger.Log(fmt.Sprintf("Returning TGDBConnection:configureQueryRequest w/ QueryRequest: '%+v'", qryReq.String()))
	return
}

func (obj *TGDBConnection) populateResultSetFromQueryResponse(resultId int, msgResponse *pdu.QueryResponseMessage) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:populateResultSetFromQueryResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse as msgResponse does not have any results"))
		errMsg := "TGDBConnection::populateResultSetFromQueryResponse does not have any results in QueryResponseMessage"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]types.TGEntity, 0)
	var rSet *query.ResultSet

	currResultCount := 0
	resultCount := msgResponse.GetResultCount()
	logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse read resultCount: '%d' FetchedEntityCount: '%d'", resultCount, len(fetchedEntities)))
	if resultCount > 0 {
		respStream.SetReferenceMap(fetchedEntities)
		rSet = query.NewResultSet(obj, resultId)
	}

	totalCount := msgResponse.GetTotalCount()
	logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse read totalCount: '%d'", totalCount))
	for i := 0; i < totalCount; i++ {
		entityType, err := respStream.(*iostream.ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse - unable to read entityType in the response stream"))
			errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to read entity type in the response stream"
			return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
		}
		kindId := types.TGEntityKind(entityType)
		logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse read #'%d'-entityType: '%+v', kindId: '%s'", i, entityType, kindId.String()))
		if kindId != types.EntityKindInvalid {
			entityId, err := respStream.(*iostream.ProtocolDataInputStream).ReadLong()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse - unable to read entityId in the response stream"))
				errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to read entity type in the response stream"
				return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
			}
			entity := fetchedEntities[entityId]
			logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse read entityId: '%d', kindId: '%s', entity: '%+v'", entityId, kindId.String(), entity))
			switch kindId {
			case types.EntityKindNode:
				var node model.Node
				if entity == nil {
					node, nErr := obj.graphObjFactory.CreateNode()
					if nErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to create a new node from the response stream"
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.Error())
					}
					entity = node
					fetchedEntities[entityId] = node
					if currResultCount < resultCount {
						rSet.AddEntityToResultSet(node)
					}
					logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse created new node: '%+v' FetchedEntityCount: '%d'", node, len(fetchedEntities)))
				}
				node = *(entity.(*model.Node))
				err := node.ReadExternal(respStream)
				if err != nil {
					errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to node.ReadExternal() from the response stream"
					logger.Error(errMsg)
					return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.GetErrorDetails())
				}
				//logger.Log(fmt.Sprintf("======> ======> After node.ReadExternal() FetchedEntityCount: '%d'", len(fetchedEntities)))
				logger.Log(fmt.Sprintf("======> ======> Node w/ Edges: '%+v'\n", node.GetEdges()))
			case types.EntityKindEdge:
				var edge model.Edge
				if entity == nil {
					//edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, types.DirectionTypeBiDirectional)
					edge, eErr := obj.graphObjFactory.CreateEntity(types.EntityKindEdge)
					if eErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromQueryResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "TGDBConnection::populateResultSetFromQueryResponse unable to create a new bi-directional edge from the response stream"
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, eErr.Error())
					}
					entity = edge
					fetchedEntities[entityId] = edge
					if currResultCount < resultCount {
						rSet.AddEntityToResultSet(edge)
					}
					logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromQueryResponse created new edge: '%+v' FetchedEntityCount: '%d'", edge, len(fetchedEntities)))
				}
				edge = *(entity.(*model.Edge))
				err := edge.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromQueryResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				//logger.Log(fmt.Sprintf("======> ======> After edge.ReadExternal() FetchedEntityCount: '%d'", len(fetchedEntities)))
				logger.Log(fmt.Sprintf("======> ======> Edge w/ Vertices: '%+v'\n", edge.GetVertices()))
			case types.EntityKindGraph:
				// TODO: Revisit later - Should we break after throwing/logging an error
				continue
			}
			//if entity != nil {
			//	logger.Log(fmt.Sprintf("======> TGDBConnection::populateResultSetFromQueryResponse entityId: '%d', kindId: '%d', entityType: '%+v'\n", entityId, kindId, kindId.String()))
			//	attrList, _ := entity.GetAttributes()
			//	for _, attrib := range attrList {
			//		logger.Log(fmt.Sprintf("======> Attribute Value: '%+v'\n", attrib.GetValue()))
			//	}
			//	if kindId == types.EntityKindNode {
			//		edges := (entity.(*model.Node)).GetEdges()
			//		logger.Log(fmt.Sprintf("======> Node w/ Edges: '%+v'\n", edges))
			//	} else if kindId == types.EntityKindEdge {
			//		vertices := (entity.(*model.Edge)).GetVertices()
			//		logger.Log(fmt.Sprintf("======> Edge w/ Vertices: '%+v'\n", vertices))
			//	}
			//}
		} else {
			logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:populateResultSetFromQueryResponse - Received invalid entity kind %d", kindId))
		} // Valid entity types
	} // End of for loop
	logger.Log(fmt.Sprintf("Returning TGDBConnection:populateResultSetFromQueryResponse w/ ResultSet: '%+v'", rSet))
	return rSet, nil
}

func (obj *TGDBConnection) populateResultSetFromGetEntitiesResponse(msgResponse *pdu.GetEntityResponseMessage) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:populateResultSetFromGetEntitiesResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse as msgResponse does not have any results"))
		errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse does not have any results"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]types.TGEntity, 0)

	totalCount, err := respStream.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read totalCount in the response stream w/ error: '%s'", err.Error()))
		errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read total count of entities in the response stream"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse extracted totalCount: '%d'", totalCount))
	if totalCount > 0 {
		respStream.SetReferenceMap(fetchedEntities)
	}

	rSet := query.NewResultSet(obj, msgResponse.GetResultId())
	// Number of entities matches the search.  Exclude the related entities
	currResultCount := 0
	_, err = respStream.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read count of result entities in the response stream w/ error: '%s'", err.Error()))
		errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read count of result entities in the response stream"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}

	for i := 0; i < totalCount; i++ {
		isResult, err := respStream.(*iostream.ProtocolDataInputStream).ReadBoolean()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read isResult in the response stream w/ error: '%s'", err.Error()))
			errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read count of result entities in the response stream"
			return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
		}
		logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse read isResult: '%+v'", isResult))
		entityType, err := respStream.(*iostream.ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
			errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read entity type in the response stream"
			return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
		}
		kindId := types.TGEntityKind(entityType)
		logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse extracted entityType: '%+v', kindId: '%d'", entityType, kindId))
		if kindId != types.EntityKindInvalid {
			entityId, err := respStream.(*iostream.ProtocolDataInputStream).ReadLong()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to read entityId in the response stream w/ error: '%s'", err.Error()))
				errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to read entity type in the response stream"
				return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
			}
			entity := fetchedEntities[entityId]
			logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse extracted entityId: '%d', entity: '%+v'", entityId, entity))
			switch kindId {
			case types.EntityKindNode:
				var node model.Node
				if entity == nil {
					node, nErr := obj.graphObjFactory.CreateNode()
					if nErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to create a new node from the response stream"
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.GetErrorDetails())
					}
					entity = node
					fetchedEntities[entityId] = node
					logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse created new node: '%+v'", node))
				}
				node = *(entity.(*model.Node))
				err := node.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromGetEntitiesResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				if isResult {
					rSet.AddEntityToResultSet(entity)
					currResultCount++
				}
			case types.EntityKindEdge:
				var edge model.Edge
				if entity == nil {
					edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, types.DirectionTypeBiDirectional)
					if eErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntitiesResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
						// TODO: Revisit later - Should we break after throwing/logging an error
						//continue
						errMsg := "TGDBConnection::populateResultSetFromGetEntitiesResponse unable to create a new bi-directional edge from the response stream"
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, eErr.Error())
					}
					entity = edge
					fetchedEntities[entityId] = edge
					logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntitiesResponse created new edge: '%+v'", edge))
				}
				edge = *(entity.(*model.Edge))
				err := edge.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromGetEntitiesResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
					logger.Error(errMsg)
					return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				if isResult {
					rSet.AddEntityToResultSet(entity)
					currResultCount++
				}
			case types.EntityKindGraph:
				// TODO: Revisit later - Should we break after throwing/logging an error
				continue
			}
		} else {
			logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:populateResultSetFromGetEntitiesResponse - Received invalid entity kind %d", kindId))
		} // Valid entity types
	} // End of for loop
	logger.Log(fmt.Sprintf("Returning TGDBConnection:populateResultSetFromGetEntitiesResponse w/ ResultSe: '%+v'", rSet))
	return rSet, nil
}

func (obj *TGDBConnection) populateResultSetFromGetEntityResponse(msgResponse *pdu.GetEntityResponseMessage) (types.TGEntity, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:populateResultSetFromGetEntityResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse as msgResponse does not have any results"))
		errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse does not have any results"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]types.TGEntity, 0)

	var entityFound types.TGEntity
	respStream.SetReferenceMap(fetchedEntities)

	count, err := respStream.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to read count in the response stream w/ error: '%s'", err.Error()))
		errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to read count of entities in the response stream"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse extracted Count: '%d'", count))
	if count > 0 {
		respStream.SetReferenceMap(fetchedEntities)
		for i := 0; i < count; i++ {
			entityType, err := respStream.(*iostream.ProtocolDataInputStream).ReadByte()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
				errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to read entity type in the response stream"
				return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
			}
			kindId := types.TGEntityKind(entityType)
			logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse extracted entityType: '%+v', kindId: '%d'", entityType, kindId))
			if kindId != types.EntityKindInvalid {
				entityId, err := respStream.(*iostream.ProtocolDataInputStream).ReadLong()
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to read entityId in the response stream w/ error: '%s'", err.Error()))
					errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to read entity type in the response stream"
					return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				entity := fetchedEntities[entityId]
				logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse extracted entityId: '%d', entity: '%+v'", entityId, entity))
				switch kindId {
				case types.EntityKindNode:
					// Need to put shell object into map to be deserialized later
					var node model.Node
					if entity == nil {
						node, nErr := obj.graphObjFactory.CreateNode()
						if nErr != nil {
							logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
							// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
							//continue
							errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to create a new node from the response stream"
							return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.Error())
						}
						entity = node
						fetchedEntities[entityId] = node
						if entityFound == nil {
							entityFound = node
						}
						logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse created new node: '%+v'", entity))
					}
					node = *(entity.(*model.Node))
					err := node.ReadExternal(respStream)
					if err != nil {
						errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromGetEntityResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
						logger.Error(errMsg)
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
					}
				case types.EntityKindEdge:
					var edge model.Edge
					if entity == nil {
						edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, types.DirectionTypeBiDirectional)
						if eErr != nil {
							logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:populateResultSetFromGetEntityResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
							// TODO: Revisit later - Should we break after throwing/logging an error
							//continue
							errMsg := "TGDBConnection::populateResultSetFromGetEntityResponse unable to create a new bi-directional edge from the response stream"
							return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, eErr.Error())
						}
						entity = edge
						fetchedEntities[entityId] = edge
						if entityFound == nil {
							entityFound = edge
						}
						logger.Log(fmt.Sprintf("Inside TGDBConnection::populateResultSetFromGetEntityResponse created new edge: '%+v'", edge))
					}
					edge = *(entity.(*model.Edge))
					err := edge.ReadExternal(respStream)
					if err != nil {
						errMsg := fmt.Sprintf("TGDBConnection::populateResultSetFromGetEntityResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
						logger.Error(errMsg)
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
					}
				case types.EntityKindGraph:
					// TODO: Revisit later - Should we break after throwing/logging an error
					continue
				}
			} else {
				logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:populateResultSetFromGetEntityResponse - Received invalid entity kind %d", kindId))
			} // Valid entity types
		} // End of for loop
	}
	logger.Log(fmt.Sprintf("Returning TGDBConnection:populateResultSetFromGetEntityResponse w/ Entity: '%+v'", entityFound))
	return entityFound, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConnection
/////////////////////////////////////////////////////////////////

// Commit commits the current transaction on this connection
func (obj *TGDBConnection) Commit() (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprint("Entering TGDBConnection:Commit"))
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::Commit - about to loop through addedList to include existing nodes to the changed list if it's part of a new edge"))
	// Include existing nodes to the changed list if it's part of a new edge
	for _, addEntity := range obj.GetAddedList() {
		if addEntity.GetEntityKind() == types.EntityKindEdge {
			nodes := addEntity.(*model.Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*model.Node)
					if !node.GetIsNew() {
						obj.changedList[node.GetVirtualId()] = node
						logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:Commit - Existing node '%d' added to change list for a new edge", node.GetVirtualId()))
					}
				}
			}
		}
	}

	logger.Log(fmt.Sprint("Inside TGDBConnection::Commit - about to loop through changedList to include existing nodes to the changed list even for edge update"))
	// Need to include existing node to the changed list even for edge update
	for _, modEntity := range obj.GetChangedList() {
		if modEntity.GetEntityKind() == types.EntityKindEdge {
			nodes := modEntity.(*model.Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*model.Node)
					if !node.GetIsNew() {
						obj.changedList[node.GetVirtualId()] = node
						logger.Warning(fmt.Sprintf("WARNING: TGDBConnection:Commit - Existing node '%d' added to change list for an existing edge '%d'", node.GetVirtualId(), modEntity.GetVirtualId()))
					}
				}
			}
		}
	}

	logger.Log(fmt.Sprint("Inside TGDBConnection::Commit - about to loop through removedList to include existing nodes to the changed list even for edge update"))
	// Need to include existing node to the changed list even for edge update
	for _, delEntity := range obj.GetRemovedList() {
		if delEntity.GetEntityKind() == types.EntityKindEdge {
			nodes := delEntity.(*model.Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*model.Node)
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

	logger.Log(fmt.Sprint("Inside TGDBConnection::Commit about to createChannelRequest() for: pdu.VerbCommitTransactionRequest"))
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, pdu.VerbCommitTransactionRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:Commit - unable to createChannelRequest(pdu.VerbCommitTransactionRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*pdu.CommitTransactionRequest)

	logger.Log(fmt.Sprint("Inside TGDBConnection::Commit about to GetAttributeDescriptors() for: pdu.VerbCommitTransactionRequest"))
	gof := obj.graphObjFactory
	attrDescSet, aErr := gof.GetGraphMetaData().GetNewAttributeDescriptors()
	if aErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:Commit - unable to gmd.GetAttributeDescriptors() w/ error: '%s'", aErr.Error()))
		return nil, aErr
	}
	queryRequest.AddCommitLists(obj.addedList, obj.changedList, obj.removedList, attrDescSet)

	logger.Log(fmt.Sprint("Inside TGDBConnection::Commit about to channelSendRequest() for: pdu.VerbCommitTransactionRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:Commit - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::Commit received response for: pdu.VerbCommitTransactionRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.CommitTransactionResponse)

	if response.HasException() {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:Commit as response has exceptions"))
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil, nil
	}

	logger.Log(fmt.Sprint("Inside TGDBConnection::Commit about to fixUpAttrDescriptors()"))
	fixUpAttrDescriptors(response, attrDescSet)
	logger.Log(fmt.Sprint("Inside TGDBConnection::Commit about to obj.fixUpEntities()"))
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

	obj.addedList = make(map[int64]types.TGEntity, 0)
	obj.changedList = make(map[int64]types.TGEntity, 0)
	obj.removedList = make(map[int64]types.TGEntity, 0)
	obj.attrByTypeList = make(map[int][]types.TGAttribute, 0)

	logger.Log(fmt.Sprint("Returning TGDBConnection:Commit"))
	return nil, nil
}

// Connect establishes a network connection to the TGDB server
func (obj *TGDBConnection) Connect() types.TGError {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:Connect for connection: '%+v'", obj))
	err := obj.GetChannel().Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection::Connect - error in obj.GetChannel().Connect() as '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprint("Inside TGDBConnection::Connect about to obj.GetChannel().Start()"))
	err = obj.GetChannel().Start()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection::Connect - error in obj.GetChannel().Start() as '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprint("Returning TGDBConnection:Connect"))
	return nil
}

// CloseQuery closes a specific query and associated objects
func (obj *TGDBConnection) CloseQuery(queryHashId int64) (types.TGQuery, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:CloseQuery for QueryHashId: '%+v'", queryHashId))
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::CloseQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, pdu.VerbQueryRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:CloseQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(CLOSE)
	queryRequest.SetQueryHashId(queryHashId)

	logger.Log(fmt.Sprint("Inside TGDBConnection::CloseQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	_, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:CloseQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	logger.Log(fmt.Sprintf("Returning TGDBConnection:CloseQuery"))
	return nil, nil
}

// CreateQuery creates a reusable query object that can be used to execute one or more statement
func (obj *TGDBConnection) CreateQuery(expr string) (types.TGQuery, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:CreateQuery for Query: '%+v'", expr))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::CreateQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, pdu.VerbQueryRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:CreateQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(CREATE)
	queryRequest.SetQuery(expr)

	logger.Log(fmt.Sprint("Inside TGDBConnection::CreateQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:CreateQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::CreateQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))

	response := msgResponse.(*pdu.QueryResponseMessage)
	queryHashId := response.GetQueryHashId()

	logger.Log(fmt.Sprint("Returning TGDBConnection:CreateQuery"))
	if response.GetResult() == 0 && queryHashId > 0 {
		return query.NewQuery(obj, queryHashId), nil
	}
	return nil, nil
}

// DecryptBuffer decrypts the encrypted buffer by sending a DecryptBufferRequest to the server
func (obj *TGDBConnection) DecryptBuffer(is types.TGInputStream) ([]byte, types.TGError) {
	logger.Log(fmt.Sprint("Entering TGDBConnection:DecryptBuffer w/ EncryptedBuffer"))
	//err := obj.InitMetadata()
	//if err != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:DecryptBuffer - unable to initialize metadata w/ error: '%s'", err.Error()))
	//	return nil, err
	//}
	//obj.connPoolImpl.AdminLock()
	//defer obj.connPoolImpl.AdminUnlock()
	//
	//logger.Log(fmt.Sprint("Inside TGDBConnection::DecryptBuffer about to createChannelRequest() for: pdu.VerbDecryptBufferRequest"))
	//// Create a channel request
	//msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbDecryptBufferRequest)
	//if cErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:DecryptBuffer - unable to createChannelRequest(pdu.VerbDecryptBufferRequest w/ error: '%s'", cErr.Error()))
	//	return nil, cErr
	//}
	//decryptRequest := msgRequest.(*pdu.DecryptBufferRequestMessage)
	//decryptRequest.SetEncryptedBuffer(encryptedBuf)
	//
	//logger.Log(fmt.Sprint("Inside TGDBConnection::DecryptBuffer about to obj.GetChannel().SendRequest() for: pdu.VerbDecryptBufferRequest"))
	//// Execute request on channel and get the response
	//msgResponse, channelErr := obj.GetChannel().SendRequest(decryptRequest, channelResponse.(*channel.BlockingChannelResponse))
	//if channelErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:DecryptBuffer - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
	//	return nil, channelErr
	//}
	//logger.Log(fmt.Sprintf("Inside TGDBConnection::DecryptBuffer received response for: pdu.VerbGetLargeObjectRequest as '%+v'", msgResponse))
	//response := msgResponse.(*pdu.DecryptBufferResponseMessage)
	//
	//if response == nil {
	//	errMsg := "TGDBConnection::DecryptBuffer does not have any results in GetLargeObjectResponseMessage"
	//	logger.Error(errMsg)
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	//}
	//
	//logger.Log(fmt.Sprintf("Returning TGDBConnection:DecryptBuffer w/ '%+v'", response.GetDecryptedBuffer()))
	//return response.GetDecryptedBuffer(), nil
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Decrypt(is)
}

// DecryptEntity decrypts the encrypted entity using channel's data cryptographer
func (obj *TGDBConnection) DecryptEntity(entityId int64) ([]byte, types.TGError) {
	buf, err := obj.GetLargeObjectAsBytes(entityId, true)
	if err != nil {
		return nil, err
	}
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Decrypt(iostream.NewProtocolDataInputStream(buf))
}

// DeleteEntity marks an ENTITY for delete operation. Upon commit, the entity will be deleted from the database
func (obj *TGDBConnection) DeleteEntity(entity types.TGEntity) types.TGError {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:DeleteEntity for Entity: '%+v'", entity))
	obj.removedList[entity.GetVirtualId()] = entity
	logger.Log(fmt.Sprint("Returning TGDBConnection:DeleteEntity"))
	return nil
}

// Disconnect breaks the connection from the TGDB server
func (obj *TGDBConnection) Disconnect() types.TGError {
	logger.Log(fmt.Sprint("Entering TGDBConnection:Disconnect"))
	err := obj.GetChannel().Disconnect()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:Disconnect - unable to channel.Disconnect() w/ error: '%s'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprint("Inside TGDBConnection::Disconnect about to obj.GetChannel().Stop()"))
	obj.GetChannel().Stop(false)
	logger.Log(fmt.Sprint("Returning TGDBConnection:Disconnect"))
	return nil
}

// EncryptEntity encrypts the encrypted entity using channel's data cryptographer
func (obj *TGDBConnection) EncryptEntity(rawBuffer []byte) ([]byte, types.TGError) {
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Encrypt(rawBuffer)
}

// ExecuteGremlinQuery executes a Gremlin Grammer-Based query with  query options
func (obj *TGDBConnection) ExecuteGremlinQuery(expr string, collection []interface{}, options types.TGQueryOption) ([]interface{}, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:ExecuteGremlinQuery for Query: '%+v'", expr))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTEGREMLIN)
	queryRequest.SetQuery(expr)
	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinQuery about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::ExecuteGremlinQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::ExecuteGremlinQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	respStream := response.GetEntityStream()
	logger.Log(fmt.Sprintf("Returning TGDBConnection:ExecuteGremlinQuery w/ '%+v'", response))
	err = query.FillCollection(respStream, obj.graphObjFactory, collection)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinQuery - unable to query.FillCollection w/ error: '%s'", err.Error()))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Returning TGDBConnection:ExecuteGremlinQuery w/ '%+v'", collection))
	return collection, nil
}

// ExecuteGremlinStrQuery executes a Gremlin Grammer-Based string query with  query options
func (obj *TGDBConnection) ExecuteGremlinStrQuery(strQuery string, options types.TGQueryOption) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:ExecuteGremlinStrQuery for Query: '%+v'", strQuery))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinStrQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinStrQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTEGREMLINSTR)
	queryRequest.SetQuery(strQuery)
	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinStrQuery about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteGremlinStrQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinStrQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::ExecuteGremlinStrQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::ExecuteGremlinStrQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	// TODO: Revisit later once Gremlin Package and GremlinQueryResult are implemented
	//This is just a dummy value where > 0 means it has results
	//resultCount := response.GetResultCount()
	resultSet := query.NewResultSet(obj, 0)
	respStream := response.GetEntityStream()
	err = query.FillCollection(respStream, obj.graphObjFactory, resultSet.GetResults())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteGremlinQuery - unable to query.FillCollection w/ error: '%s'", err.Error()))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Returning TGDBConnection:ExecuteGremlinStrQuery w/ '%+v'", resultSet))
	return resultSet, nil
}

// ExecuteQuery executes an immediate query with associated query options
func (obj *TGDBConnection) ExecuteQuery(expr string, options types.TGQueryOption) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:ExecuteQuery for Query: '%+v'", expr))

	// TODO: Revisit later once Gremlin Package and GremlinQueryResult are implemented
	cn := utils.GetConfigFromKey(utils.ConnectionDefaultQueryLanguage)
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
			return query.NewResultSet(obj, 0), nil
		}
	} else {
		switch queryLang {
		case "tgql":
			return obj.ExecuteTGDBQuery(expr, options)
		case "gremlin":
			return obj.ExecuteGremlinStrQuery(expr, options)
		default:
			return query.NewResultSet(obj, 0), nil
		}
	}
	return nil, nil
}

// ExecuteTGDBQuery executes an immediate query with associated query options
func (obj *TGDBConnection) ExecuteTGDBQuery(expr string, options types.TGQueryOption) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:ExecuteTGDBQuery for Query: '%+v'", expr))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteTGDBQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteTGDBQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery(expr)
	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteTGDBQuery about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteTGDBQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteTGDBQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::ExecuteTGDBQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::ExecuteTGDBQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("Returning TGDBConnection:ExecuteTGDBQuery w/ '%+v'", response))
	return obj.populateResultSetFromQueryResponse(0, response)
}

// ExecuteQueryWithFilter executes an immediate query with specified filter & query options
// The query option is place holder at this time
// @param expr A subset of SQL-92 where clause
// @param edgeFilter filter used for selecting edges to be returned
// @param traversalCondition condition used for selecting edges to be traversed and returned
// @param endCondition condition used to stop the traversal
// @param option Query options for executing. Can be null, then it will use the default option
func (obj *TGDBConnection) ExecuteQueryWithFilter(expr, edgeFilter, traversalCondition, endCondition string, options types.TGQueryOption) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:ExecuteQueryWithFilter for Query: '%+v', EdgeFilter: '%+v', Traversal: '%+v', EndCondition: '%+v'", expr, edgeFilter, traversalCondition, endCondition))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithFilter about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteQueryWithFilter - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery(expr)
	queryRequest.SetEdgeFilter(edgeFilter)
	queryRequest.SetTraversalCondition(traversalCondition)
	queryRequest.SetEndCondition(endCondition)
	logger.Log(fmt.Sprintf("Inside TGDBConnection::ExecuteQueryWithFilter about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithFilter about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteQueryWithFilter - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::ExecuteQueryWithFilter received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::ExecuteQueryWithFilter - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("Returning TGDBConnection:ExecuteQueryWithFilter w/ '%+v'", response))
	return obj.populateResultSetFromQueryResponse(0, response)
}

// ExecuteQueryWithId executes an immediate query for specified id & query options
func (obj *TGDBConnection) ExecuteQueryWithId(queryHashId int64, options types.TGQueryOption) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:ExecuteQueryWithId for QueryHashId: '%+v'", queryHashId))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithId about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteQueryWithId - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTED)
	queryRequest.SetQueryHashId(queryHashId)
	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithId about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside TGDBConnection::ExecuteQueryWithId about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	_, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:ExecuteQueryWithId - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	logger.Log(fmt.Sprint("Returning TGDBConnection:ExecuteQueryWithId"))
	return nil, nil
}

// GetAddedList gets a list of added entities
func (obj *TGDBConnection) GetAddedList() map[int64]types.TGEntity {
	return obj.addedList
}

// GetChangedList gets a list of changed entities
func (obj *TGDBConnection) GetChangedList() map[int64]types.TGEntity {
	return obj.changedList
}

// GetChangedList gets the communication channel associated with this connection
func (obj *TGDBConnection) GetChannel() types.TGChannel {
	return obj.channel
}

// GetConnectionId gets connection identifier
func (obj *TGDBConnection) GetConnectionId() int64 {
	return obj.connId
}

// GetConnectionProperties gets a list of connection properties
func (obj *TGDBConnection) GetConnectionProperties() types.TGProperties {
	return obj.connProperties
}

// GetEntities gets a result set of entities given an non-uniqueKey
func (obj *TGDBConnection) GetEntities(qryKey types.TGKey, props types.TGProperties) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:GetEntities for QueryKey: '%+v'", qryKey))
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:GetEntities - unable to InitMetadata"))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if props == nil {
		props = query.NewQueryOption()
	}
	logger.Log(fmt.Sprint("Inside TGDBConnection::GetEntities about to createChannelRequest() for: pdu.VerbGetEntityRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbGetEntityRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetEntities - unable to createChannelRequest(pdu.VerbGetEntityRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.GetEntityRequestMessage)
	queryRequest.SetCommand(2)
	queryRequest.SetKey(qryKey)
	configureGetRequest(queryRequest, props)

	logger.Log(fmt.Sprint("Inside TGDBConnection::GetEntities about to obj.GetChannel().SendRequest() for: pdu.VerbGetEntityRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetEntities - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::GetEntities received response for: pdu.VerbGetEntityRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.GetEntityResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::GetEntities - The request does not have any results in GetEntityResponseMessage"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("Returning TGDBConnection:GetEntities w/ '%+v'", response))
	return obj.populateResultSetFromGetEntitiesResponse(response)
}

// GetEntity gets an Entity given an UniqueKey for the Object
func (obj *TGDBConnection) GetEntity(qryKey types.TGKey, options types.TGQueryOption) (types.TGEntity, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:GetEntity for QueryKey: '%+v'", qryKey))
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TGDBConnection:GetEntity - unable to InitMetadata"))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if options == nil {
		options = query.NewQueryOption()
	}
	logger.Log(fmt.Sprint("Inside TGDBConnection::GetEntity about to createChannelRequest() for: pdu.VerbGetEntityRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbGetEntityRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetEntity - unable to createChannelRequest(pdu.VerbGetEntityRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.GetEntityRequestMessage)
	queryRequest.SetCommand(0)
	queryRequest.SetKey(qryKey)
	configureGetRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside TGDBConnection::GetEntity about to obj.GetChannel().SendRequest() for: pdu.VerbGetEntityRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetEntity - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::GetEntity received response for: pdu.VerbGetEntityRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.GetEntityResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning TGDBConnection::GetEntity - The request does not have any results in GetEntityResponseMessage"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("Returning TGDBConnection:GetEntity w/ '%+v'", response))
	return obj.populateResultSetFromGetEntityResponse(response)
}

// GetGraphMetadata gets the Graph Metadata
func (obj *TGDBConnection) GetGraphMetadata(refresh bool) (types.TGGraphMetadata, types.TGError) {
	logger.Log(fmt.Sprint("Entering TGDBConnection:GetGraphMetadata"))
	if refresh {
		obj.connPoolImpl.AdminLock()
		defer obj.connPoolImpl.AdminUnlock()

		logger.Log(fmt.Sprint("Inside TGDBConnection::GetGraphMetadata about to createChannelRequest() for: pdu.VerbMetadataRequest"))
		// Create a channel request
		msgRequest, channelResponse, err := createChannelRequest(obj, pdu.VerbMetadataRequest)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetGraphMetadata - unable to createChannelRequest(pdu.VerbMetadataRequest w/ error: '%s'", err.Error()))
			return nil, err
		}
		metaRequest := msgRequest.(*pdu.MetadataRequest)
		logger.Log(fmt.Sprintf("Inside TGDBConnection::GetGraphMetadata createChannelRequest() returned MsgRequest: '%+v' ChannelResponse: '%+v'", msgRequest, channelResponse.(*channel.BlockingChannelResponse)))

		logger.Log(fmt.Sprint("Inside TGDBConnection::GetGraphMetadata about to obj.GetChannel().SendRequest() for: pdu.VerbMetadataRequest"))
		// Execute request on channel and get the response
		msgResponse, channelErr := obj.GetChannel().SendRequest(metaRequest, channelResponse.(*channel.BlockingChannelResponse))
		if channelErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetGraphMetadata - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
			return nil, channelErr
		}
		//logger.Log(fmt.Sprintf("Inside TGDBConnection::GetGraphMetadata received response for: pdu.VerbMetadataRequest as '%+v'", msgResponse))

		response := msgResponse.(*pdu.MetadataResponse)
		attrDescList := response.GetAttrDescList()
		edgeTypeList := response.GetEdgeTypeList()
		nodeTypeList := response.GetNodeTypeList()

		gmd := obj.graphObjFactory.GetGraphMetaData()
		logger.Log(fmt.Sprint("Inside TGDBConnection::GetGraphMetadata about to update GraphMetadata"))
		uErr := gmd.UpdateMetadata(attrDescList, nodeTypeList, edgeTypeList)
		if uErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetGraphMetadata - unable to gmd.UpdateMetadata() w/ error: '%s'", uErr.Error()))
			return nil, uErr
		}
		logger.Log(fmt.Sprint("Inside TGDBConnection::GetGraphMetadata successfully updated GraphMetadata"))
	}
	logger.Log(fmt.Sprint("Returning TGDBConnection:GetGraphMetadata"))
	return obj.graphObjFactory.GetGraphMetaData(), nil
}

// GetGraphObjectFactory gets the Graph Object Factory for Object creation
func (obj *TGDBConnection) GetGraphObjectFactory() (types.TGGraphObjectFactory, types.TGError) {
	logger.Log(fmt.Sprint("Entering TGDBConnection:GetGraphObjectFactory"))
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetGraphObjectFactory - unable to initialize metadata w/ error: '%s'", err.Error()))
		return nil, err
	}
	logger.Log(fmt.Sprint("Returning TGDBConnection:GetGraphObjectFactory"))
	return obj.graphObjFactory, nil
}

// GetLargeObjectAsBytes gets an Binary Large Object Entity given an UniqueKey for the Object
func (obj *TGDBConnection) GetLargeObjectAsBytes(entityId int64, decryptFlag bool) ([]byte, types.TGError) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:GetLargeObjectAsBytes for EntityId: '%+v'", entityId))
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetLargeObjectAsBytes - unable to initialize metadata w/ error: '%s'", err.Error()))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside TGDBConnection::GetLargeObjectAsBytes about to createChannelRequest() for: pdu.VerbGetLargeObjectRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbGetLargeObjectRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetLargeObjectAsBytes - unable to createChannelRequest(pdu.VerbGetLargeObjectRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.GetLargeObjectRequestMessage)
	queryRequest.SetEntityId(entityId)
	queryRequest.SetDecryption(decryptFlag)

	logger.Log(fmt.Sprint("Inside TGDBConnection::GetLargeObjectAsBytes about to obj.GetChannel().SendRequest() for: pdu.VerbGetLargeObjectRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGDBConnection:GetLargeObjectAsBytes - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside TGDBConnection::GetLargeObjectAsBytes received response for: pdu.VerbGetLargeObjectRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.GetLargeObjectResponseMessage)

	if response == nil {
		errMsg := "TGDBConnection::GetLargeObjectAsBytes does not have any results in GetLargeObjectResponseMessage"
		logger.Error(errMsg)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	logger.Log(fmt.Sprintf("Returning TGDBConnection:GetLargeObjectAsBytes w/ '%+v'", response.GetBuffer()))
	return response.GetBuffer(), nil
}

// GetRemovedList gets a list of removed entities
func (obj *TGDBConnection) GetRemovedList() map[int64]types.TGEntity {
	return obj.removedList
}

// InsertEntity marks an ENTITY for insert operation. Upon commit, the entity will be inserted in the database
func (obj *TGDBConnection) InsertEntity(entity types.TGEntity) types.TGError {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:InsertEntity to insert Entity: '%+v'", entity.GetEntityType()))
	obj.addedList[entity.GetVirtualId()] = entity
	logger.Log(fmt.Sprint("Returning TGDBConnection:InsertEntity"))
	return nil
}

// Rollback rolls back the current transaction on this connection
func (obj *TGDBConnection) Rollback() types.TGError {
	logger.Log(fmt.Sprint("Entering TGDBConnection:Rollback"))
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	// Reset all the lists to empty contents
	obj.addedList = make(map[int64]types.TGEntity, 0)
	obj.changedList = make(map[int64]types.TGEntity, 0)
	obj.removedList = make(map[int64]types.TGEntity, 0)
	obj.attrByTypeList = make(map[int][]types.TGAttribute, 0)

	logger.Log(fmt.Sprint("Returning TGDBConnection:Rollback"))
	return nil
}

// SetExceptionListener sets exception listener
func (obj *TGDBConnection) SetExceptionListener(listener types.TGConnectionExceptionListener) {
	obj.connPoolImpl.SetExceptionListener(listener) //delegate it to the Pool.
}

// UpdateEntity marks an ENTITY for update operation. Upon commit, the entity will be updated in the database
// When commit is called, the object is resolved to check if it is dirty. Entity.setAttribute calls make the entity
// dirty. If it is dirty, then the object is send to the server for update, otherwise it is ignored.
// Calling multiple times, does not change the behavior.
// The same entity cannot be updated on multiple connections. It will result an TGException of already associated to a connection.
func (obj *TGDBConnection) UpdateEntity(entity types.TGEntity) types.TGError {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:UpdateEntity to update Entity: '%+v'", entity))
	obj.changedList[entity.GetVirtualId()] = entity
	logger.Log(fmt.Sprint("Returning TGDBConnection:UpdateEntity"))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChangeListener
/////////////////////////////////////////////////////////////////

// AttributeAdded gets called when an attribute is Added to an entity.
func (obj *TGDBConnection) AttributeAdded(attr types.TGAttribute, owner types.TGEntity) {
	logger.Log(fmt.Sprint("Returning TGDBConnection:AttributeAdded"))
}

// AttributeChanged gets called when an attribute is set.
func (obj *TGDBConnection) AttributeChanged(attr types.TGAttribute, oldValue, newValue interface{}) {
	logger.Log(fmt.Sprint("Returning TGDBConnection:AttributeChanged"))
}

// AttributeRemoved gets called when an attribute is removed from the entity.
func (obj *TGDBConnection) AttributeRemoved(attr types.TGAttribute, owner types.TGEntity) {
	logger.Log(fmt.Sprint("Returning TGDBConnection:AttributeRemoved"))
}

// EntityCreated gets called when an entity is Added
func (obj *TGDBConnection) EntityCreated(entity types.TGEntity) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:EntityCreated to add Entity: '%+v'", entity))
	entityId := entity.(*model.AbstractEntity).GetVirtualId()
	obj.addedList[entityId] = entity
	logger.Log(fmt.Sprint("Returning TGDBConnection:EntityCreated"))
}

// EntityDeleted gets called when the entity is deleted
func (obj *TGDBConnection) EntityDeleted(entity types.TGEntity) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:EntityDeleted to delete Entity: '%+v'", entity))
	entityId := entity.(*model.AbstractEntity).GetVirtualId()
	obj.removedList[entityId] = entity
	logger.Log(fmt.Sprint("Returning TGDBConnection:EntityDeleted"))
}

// NodeAdded gets called when a node is Added
func (obj *TGDBConnection) NodeAdded(graph types.TGGraph, node types.TGNode) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:NodeAdded to add Node: '%+v' to Graph: '%+v'", node, graph))
	entityId := graph.(*model.Graph).GetVirtualId()
	obj.addedList[entityId] = graph
	logger.Log(fmt.Sprint("Returning TGDBConnection:NodeAdded"))
}

// NodeRemoved gets called when a node is removed
func (obj *TGDBConnection) NodeRemoved(graph types.TGGraph, node types.TGNode) {
	logger.Log(fmt.Sprintf("Entering TGDBConnection:NodeRemoved to remove Node: '%+v' to Graph: '%+v'", node, graph))
	entityId := graph.(*model.Graph).GetVirtualId()
	obj.removedList[entityId] = graph
	logger.Log(fmt.Sprint("Returning TGDBConnection:NodeRemoved"))
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
