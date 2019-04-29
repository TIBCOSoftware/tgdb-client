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
 * File name: AdminConnectionImpl.go
 * Created on: Apr 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

package connection

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/admin"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/channel"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/pdu"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
)

//static TGLogger gLogger        = TGLogManager.getInstance().getLogger();
//var connectionIds int64
//var requestIds int64

//type TGConnectionCommand int

//const (
//	CREATE = 1 + iota
//	EXECUTE
//	EXECUTED
//	CLOSE
//)

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

func NewAdminConnection(conPool *ConnectionPoolImpl, channel types.TGChannel, props types.TGProperties) *AdminConnectionImpl {
	newSGDBConnection := DefaultAdminConnection()
	newSGDBConnection.connPoolImpl = conPool
	newSGDBConnection.channel = channel
	newSGDBConnection.connProperties = props.(*utils.SortedProperties)
	return newSGDBConnection
}

/////////////////////////////////////////////////////////////////
// Helper functions for AdminConnectionImpl
/////////////////////////////////////////////////////////////////

func (obj *AdminConnectionImpl) GetConnectionPool() types.TGConnectionPool {
	return obj.connPoolImpl
}

func (obj *AdminConnectionImpl) GetConnectionGraphObjectFactory() *model.GraphObjectFactory {
	return obj.graphObjFactory
}

func (obj *AdminConnectionImpl) InitMetadata() types.TGError {
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
func (obj *AdminConnectionImpl) SetConnectionPool(connPool types.TGConnectionPool) {
	obj.connPoolImpl = connPool
}

// SetConnectionProperties sets connection properties
func (obj *AdminConnectionImpl) SetConnectionProperties(connProps types.TGProperties) {
	obj.connProperties = connProps.(*utils.SortedProperties)
}

/////////////////////////////////////////////////////////////////
// Private functions for AdminConnectionImpl
/////////////////////////////////////////////////////////////////

//func fixUpAttrDescriptors(response *pdu.CommitTransactionResponse, attrDescSet []types.TGAttributeDescriptor) {
//	logger.Log(fmt.Sprint("Entering AdminConnectionImpl:fixUpAttrDescriptors"))
//	attrDescCount := response.GetAttrDescCount()
//	attrDescIdList := response.GetAttrDescIdList()
//	for i := 0; i < attrDescCount; i++ {
//		tempId := attrDescIdList[(i * 2)]
//		realId := attrDescIdList[((i * 2) + 1)]
//
//		for _, attrDesc := range attrDescSet {
//			desc := attrDesc.(*model.AttributeDescriptor)
//			if attrDesc.GetAttributeId() == tempId {
//				logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:fixUpAttrDescriptors - Replace descriptor: '%d' by '%d'", attrDesc.GetAttributeId(), realId))
//				desc.SetAttributeId(realId)
//				break
//			}
//		}
//	}
//	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:fixUpAttrDescriptors"))
//}
//
//func fixUpEntities(obj types.TGConnection, response *pdu.CommitTransactionResponse) {
//	logger.Log(fmt.Sprint("Entering AdminConnectionImpl:fixUpEntities"))
//	addedIdCount := response.GetAddedEntityCount()
//	addedIdList := response.GetAddedIdList()
//	for i := 0; i < addedIdCount; i++ {
//		tempId := addedIdList[(i * 3)]
//		realId := addedIdList[((i * 3) + 1)]
//		version := addedIdList[((i * 3) + 2)]
//
//		for _, addEntity := range obj.GetAddedList() {
//			if addEntity.GetVirtualId() == tempId {
//				logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:fixUpEntities - Replace entity id: '%d' by '%d'", tempId, realId))
//				addEntity.SetEntityId(realId)
//				addEntity.SetIsNew(false)
//				addEntity.SetVersion(int(version))
//				break
//			}
//		}
//	}
//
//	updatedIdCount := response.GetUpdatedEntityCount()
//	updatedIdList := response.GetUpdatedIdList()
//	for i := 0; i < updatedIdCount; i++ {
//		id := updatedIdList[(i * 2)]
//		version := updatedIdList[((i * 2) + 1)]
//
//		for _, modEntity := range obj.GetChangedList() {
//			if modEntity.GetVirtualId() == id {
//				logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:fixUpEntities - Replace entity version: '%d' to '%d'", id, version))
//				modEntity.SetVersion(int(version))
//				break
//			}
//		}
//	}
//	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:fixUpEntities"))
//}
//
//func createChannelRequest(obj types.TGConnection, verb int) (types.TGMessage, types.TGChannelResponse, types.TGError) {
//	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:createChannelRequest for Verb: '%s'", pdu.GetVerb(verb).GetName()))
//	cn := utils.GetConfigFromKey(utils.ConnectionOperationTimeoutSeconds)
//	//logger.Log(fmt.Sprintf("Inside AbstractChannel::channelTryRepeatConnect config for ConnectionOperationTimeoutSeconds is '%+v", cn))
//	timeout := obj.GetConnectionProperties().GetPropertyAsInt(cn)
//	requestId := atomic.AddInt64(&connectionIds, 1)
//
//	//logger.Log(fmt.Sprint("Inside AdminConnectionImpl::createChannelRequest about to create channel.NewBlockingChannelResponse()"))
//	// Create a non-blocking channel response
//	channelResponse := channel.NewBlockingChannelResponse(requestId, int64(timeout))
//
//	// Use Message Factory method to create appropriate message structure (class) based on input type
//	//msgRequest, err := pdu.CreateMessageForVerb(verb)
//	msgRequest, err := pdu.CreateMessageWithToken(verb, obj.GetChannel().GetAuthToken(), obj.GetChannel().GetSessionId())
//	if err != nil {
//		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:createChannelRequest - CreateMessageForVerb failed for Verb: '%s' w/ '%+v'", pdu.GetVerb(verb).GetName(), err.Error()))
//		return nil, nil, err
//	}
//	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:createChannelRequest for Verb: '%s' w/ MessageRequest: '%+v' ChannelResponse: '%+v'", pdu.GetVerb(verb).GetName(), msgRequest, channelResponse))
//	return msgRequest, channelResponse, nil
//}
//
//func configureGetRequest(getReq *pdu.GetEntityRequestMessage, reqProps types.TGProperties) {
//	if getReq == nil || reqProps == nil {
//		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::configureGetRequest as getReq == nil || reqProps == nil"))
//		return
//	}
//	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:configureGetRequest w/ EntityRequest: '%+v'", getReq.String()))
//	props := reqProps.(*query.TGQueryOptionImpl)
//	getReq.SetFetchSize(props.GetPreFetchSize())
//	getReq.SetTraversalDepth(props.GetTraversalDepth())
//	getReq.SetEdgeLimit(props.GetEdgeLimit())
//	getReq.SetBatchSize(props.GetBatchSize())
//	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:configureGetRequest w/ EntityRequest: '%+v'", getReq.String()))
//	return
//}
//
//func configureQueryRequest(qryReq *pdu.QueryRequestMessage, reqProps types.TGProperties) {
//	if qryReq == nil || reqProps == nil {
//		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl:configureQueryRequest as getReq == nil || reqProps == nil"))
//		return
//	}
//	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:configureQueryRequest w/ QueryRequest: '%+v'", qryReq.String()))
//	props := reqProps.(*query.TGQueryOptionImpl)
//	qryReq.SetFetchSize(props.GetPreFetchSize())
//	qryReq.SetTraversalDepth(props.GetTraversalDepth())
//	qryReq.SetEdgeLimit(props.GetEdgeLimit())
//	qryReq.SetBatchSize(props.GetBatchSize())
//	qryReq.SetSortAttrName(props.GetSortAttrName())
//	qryReq.SetSortOrderDsc(props.IsSortOrderDsc())
//	qryReq.SetSortResultLimit(props.GetSortResultLimit())
//	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:configureQueryRequest w/ QueryRequest: '%+v'", qryReq.String()))
//	return
//}

/**
	private Object getCommandWithoutParameters (TGAdminCommand command) throws TGException {
		connPool.adminlock();

		TGChannelResponse channelResponse;
		long timeout = Long.parseLong(properties.getProperty(ConfigName.ConnectionOperationTimeoutSeconds, "-1"));

		long requestId  = requestIds.getAndIncrement();
		channelResponse = new BlockingChannelResponse(requestId, timeout);

		AdminResponse response = null;
		try {
			AdminRequest request = (AdminRequest) TGMessageFactory.getInstance().createMessage(VerbId.AdminRequest);
			request.setCommand(command);
			response = (AdminResponse) this.channel.sendRequest(request, channelResponse);
		}
		finally {
			connPool.adminUnlock();
		}

		Object result = null;
		switch (command)
		{
			case ShowInfo: {
				result = response.getAdminCommandInfoResult();
				break;
			}
			case ShowUsers: {
				result = response.getAdminCommandUsersResult();
				break;
			}
			case ShowConnections: {
				result = response.getAdminCommandConnectionsResult();
				break;
			}
			case ShowAttrDescs: {
				result = response.getAdminCommandAttrDescsResult();
				break;
			}
			case ShowIndices: {
				result = response.getIndices();
				break;
			}
			case StopServer: {
			break;
			}
			case SetLogLevel: {
			break;
			}
		}
		return result;
	}

	private Object getCommandWithParameters (TGAdminCommand command, Object parameters) throws TGException {
		connPool.adminlock();

		TGChannelResponse channelResponse;
		long timeout = Long.parseLong(properties.getProperty(ConfigName.ConnectionOperationTimeoutSeconds, "-1"));

		long requestId  = requestIds.getAndIncrement();
		channelResponse = new BlockingChannelResponse(requestId, timeout);

		AdminResponse response = null;
		try {
			AdminRequest request = (AdminRequest) TGMessageFactory.getInstance().createMessage(VerbId.AdminRequest);
			request.setCommand(command);

			switch (command)
			{
				case KillConnection: {
					//request.setKillConnectionInfo((TGAdminKillConnectionInfo)parameters);
					request.setSessionId((Long)parameters);
					//request.setKillConnectionInfo(parameters);
					break;
				}
				case SetLogLevel: {
					request.setLogLevel((TGServerLogDetails) parameters);
				}
			}

			response = (AdminResponse) this.channel.sendRequest(request, channelResponse);
		}
		finally {
			connPool.adminUnlock();
		}

		Object result = null;

		//
		// The switch case may be needed for other commands later
		//

		switch (command)
		{
			case SHOW_INFO: {
				result = response.getAdminCommandInfoResult();
				break;
			}
			case SHOW_USERS: {
				result = response.getAdminCommandUsersResult();
				break;
			}
			case SHOW_CONNECTIONS: {
				result = response.getAdminCommandConnectionsResult();
				break;
			}
			case STOP_SERVER: {
				break;
			}
		}
		return result;
	}
*/

func (obj *AdminConnectionImpl) populateResultSetFromQueryResponse(resultId int, msgResponse *pdu.QueryResponseMessage) (types.TGResultSet, types.TGError) {
																																									   logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:populateResultSetFromQueryResponse w/ MsgResponse: '%+v'", msgResponse.String()))
																																									   if !msgResponse.GetHasResult() {
																																									   logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse as msgResponse does not have any results"))
																																									   errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse does not have any results in QueryResponseMessage"
																																									   return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
																																									   }

																																									   respStream := msgResponse.GetEntityStream()
																																									   fetchedEntities := make(map[int64]types.TGEntity, 0)
																																									   var rSet *query.ResultSet

																																									   currResultCount := 0
																																									   resultCount := msgResponse.GetResultCount()
																																									   logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse read resultCount: '%d' FetchedEntityCount: '%d'", resultCount, len(fetchedEntities)))
																																									   if resultCount > 0 {
																																									   respStream.SetReferenceMap(fetchedEntities)
																																									   rSet = query.NewResultSet(obj, resultId)
																																									   }

																																									   totalCount := msgResponse.GetTotalCount()
																																									   logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse read totalCount: '%d'", totalCount))
																																									   for i := 0; i < totalCount; i++ {
																																									   entityType, err := respStream.(*iostream.ProtocolDataInputStream).ReadByte()
																																									   if err != nil {
																																									   logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse - unable to read entityType in the response stream"))
																																									   errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to read entity type in the response stream"
																																									   return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
																																									   }
																																									   kindId := types.TGEntityKind(entityType)
																																									   logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse read #'%d'-entityType: '%+v', kindId: '%s'", i, entityType, kindId.String()))
																																									   if kindId != types.EntityKindInvalid {
																																									   entityId, err := respStream.(*iostream.ProtocolDataInputStream).ReadLong()
																																									   if err != nil {
																																									   logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse - unable to read entityId in the response stream"))
																																									   errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to read entity type in the response stream"
																																									   return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
																																									   }
																																									   entity := fetchedEntities[entityId]
																																									   logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse read entityId: '%d', kindId: '%s', entity: '%+v'", entityId, kindId.String(), entity))
																																									   switch kindId {
																																									   case types.EntityKindNode:
																																									   var node model.Node
																																									   if entity == nil {
																																									   node, nErr := obj.graphObjFactory.CreateNode()
																																									   if nErr != nil {
																																									   logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
																																									   // TODO: Revisit later - Should we continue OR break after throwing/logging an error?
																																									   //continue
																																									   errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to create a new node from the response stream"
																																									   return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.Error())
																																									   }
																																									   entity = node
																																									   fetchedEntities[entityId] = node
																																									   if currResultCount < resultCount {
																																									   rSet.AddEntityToResultSet(node)
																																									   }
																																									   logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse created new node: '%+v' FetchedEntityCount: '%d'", node, len(fetchedEntities)))
																																									   }
																																									   node = *(entity.(*model.Node))
																																									   err := node.ReadExternal(respStream)
																																									   if err != nil {
																																									   errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to node.ReadExternal() from the response stream"
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
																																									   logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromQueryResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
																																									   // TODO: Revisit later - Should we continue OR break after throwing/logging an error?
																																									   //continue
																																									   errMsg := "AdminConnectionImpl::populateResultSetFromQueryResponse unable to create a new bi-directional edge from the response stream"
																																									   return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, eErr.Error())
																																									   }
																																									   entity = edge
																																									   fetchedEntities[entityId] = edge
																																									   if currResultCount < resultCount {
																																									   rSet.AddEntityToResultSet(edge)
																																									   }
																																									   logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromQueryResponse created new edge: '%+v' FetchedEntityCount: '%d'", edge, len(fetchedEntities)))
																																									   }
																																									   edge = *(entity.(*model.Edge))
																																									   err := edge.ReadExternal(respStream)
																																									   if err != nil {
																																									   errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromQueryResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
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
																																									   //	logger.Log(fmt.Sprintf("======> AdminConnectionImpl::populateResultSetFromQueryResponse entityId: '%d', kindId: '%d', entityType: '%+v'\n", entityId, kindId, kindId.String()))
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
																																									   logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:populateResultSetFromQueryResponse - Received invalid entity kind %d", kindId))
																																									   } // Valid entity types
																																									   } // End of for loop
																																									   logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:populateResultSetFromQueryResponse w/ ResultSet: '%+v'", rSet))
																																									   return rSet, nil
																																									   }

func (obj *AdminConnectionImpl) populateResultSetFromGetEntitiesResponse(msgResponse *pdu.GetEntityResponseMessage) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:populateResultSetFromGetEntitiesResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse as msgResponse does not have any results"))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse does not have any results"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]types.TGEntity, 0)

	totalCount, err := respStream.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read totalCount in the response stream w/ error: '%s'", err.Error()))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read total count of entities in the response stream"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse extracted totalCount: '%d'", totalCount))
	if totalCount > 0 {
		respStream.SetReferenceMap(fetchedEntities)
	}

	rSet := query.NewResultSet(obj, msgResponse.GetResultId())
	// Number of entities matches the search.  Exclude the related entities
	currResultCount := 0
	_, err = respStream.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read count of result entities in the response stream w/ error: '%s'", err.Error()))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read count of result entities in the response stream"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}

	for i := 0; i < totalCount; i++ {
		isResult, err := respStream.(*iostream.ProtocolDataInputStream).ReadBoolean()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read isResult in the response stream w/ error: '%s'", err.Error()))
			errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read count of result entities in the response stream"
			return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
		}
		logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse read isResult: '%+v'", isResult))
		entityType, err := respStream.(*iostream.ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
			errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read entity type in the response stream"
			return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
		}
		kindId := types.TGEntityKind(entityType)
		logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse extracted entityType: '%+v', kindId: '%d'", entityType, kindId))
		if kindId != types.EntityKindInvalid {
			entityId, err := respStream.(*iostream.ProtocolDataInputStream).ReadLong()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to read entityId in the response stream w/ error: '%s'", err.Error()))
				errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to read entity type in the response stream"
				return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
			}
			entity := fetchedEntities[entityId]
			logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse extracted entityId: '%d', entity: '%+v'", entityId, entity))
			switch kindId {
			case types.EntityKindNode:
				var node model.Node
				if entity == nil {
					node, nErr := obj.graphObjFactory.CreateNode()
					if nErr != nil {
						logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
						// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
						//continue
						errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to create a new node from the response stream"
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.GetErrorDetails())
					}
					entity = node
					fetchedEntities[entityId] = node
					logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse created new node: '%+v'", node))
				}
				node = *(entity.(*model.Node))
				err := node.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
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
						logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
						// TODO: Revisit later - Should we break after throwing/logging an error
						//continue
						errMsg := "AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to create a new bi-directional edge from the response stream"
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, eErr.Error())
					}
					entity = edge
					fetchedEntities[entityId] = edge
					logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntitiesResponse created new edge: '%+v'", edge))
				}
				edge = *(entity.(*model.Edge))
				err := edge.ReadExternal(respStream)
				if err != nil {
					errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromGetEntitiesResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
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
			logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:populateResultSetFromGetEntitiesResponse - Received invalid entity kind %d", kindId))
		} // Valid entity types
	} // End of for loop
	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:populateResultSetFromGetEntitiesResponse w/ ResultSe: '%+v'", rSet))
	return rSet, nil
}

func (obj *AdminConnectionImpl) populateResultSetFromGetEntityResponse(msgResponse *pdu.GetEntityResponseMessage) (types.TGEntity, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:populateResultSetFromGetEntityResponse w/ MsgResponse: '%+v'", msgResponse.String()))
	if !msgResponse.GetHasResult() {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse as msgResponse does not have any results"))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse does not have any results"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	respStream := msgResponse.GetEntityStream()
	fetchedEntities := make(map[int64]types.TGEntity, 0)

	var entityFound types.TGEntity
	respStream.SetReferenceMap(fetchedEntities)

	count, err := respStream.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to read count in the response stream w/ error: '%s'", err.Error()))
		errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to read count of entities in the response stream"
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse extracted Count: '%d'", count))
	if count > 0 {
		respStream.SetReferenceMap(fetchedEntities)
		for i := 0; i < count; i++ {
			entityType, err := respStream.(*iostream.ProtocolDataInputStream).ReadByte()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
				errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to read entity type in the response stream"
				return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
			}
			kindId := types.TGEntityKind(entityType)
			logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse extracted entityType: '%+v', kindId: '%d'", entityType, kindId))
			if kindId != types.EntityKindInvalid {
				entityId, err := respStream.(*iostream.ProtocolDataInputStream).ReadLong()
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to read entityId in the response stream w/ error: '%s'", err.Error()))
					errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to read entity type in the response stream"
					return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				entity := fetchedEntities[entityId]
				logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse extracted entityId: '%d', entity: '%+v'", entityId, entity))
				switch kindId {
				case types.EntityKindNode:
					// Need to put shell object into map to be deserialized later
					var node model.Node
					if entity == nil {
						node, nErr := obj.graphObjFactory.CreateNode()
						if nErr != nil {
							logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to CreateNode() w/ error: '%s'", nErr.Error()))
							// TODO: Revisit later - Should we continue OR break after throwing/logging an error?
							//continue
							errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to create a new node from the response stream"
							return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.Error())
						}
						entity = node
						fetchedEntities[entityId] = node
						if entityFound == nil {
							entityFound = node
						}
						logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse created new node: '%+v'", entity))
					}
					node = *(entity.(*model.Node))
					err := node.ReadExternal(respStream)
					if err != nil {
						errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to node.ReadExternal() from the response stream w/ error: '%s'", err.Error())
						logger.Error(errMsg)
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
					}
				case types.EntityKindEdge:
					var edge model.Edge
					if entity == nil {
						edge, eErr := obj.graphObjFactory.CreateEdgeWithDirection(nil, nil, types.DirectionTypeBiDirectional)
						if eErr != nil {
							logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse - unable to CreateEdgeWithDirection() w/ error: '%s'", eErr.Error()))
							// TODO: Revisit later - Should we break after throwing/logging an error
							//continue
							errMsg := "AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to create a new bi-directional edge from the response stream"
							return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, eErr.Error())
						}
						entity = edge
						fetchedEntities[entityId] = edge
						if entityFound == nil {
							entityFound = edge
						}
						logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::populateResultSetFromGetEntityResponse created new edge: '%+v'", edge))
					}
					edge = *(entity.(*model.Edge))
					err := edge.ReadExternal(respStream)
					if err != nil {
						errMsg := fmt.Sprintf("AdminConnectionImpl::populateResultSetFromGetEntityResponse unable to edge.ReadExternal() from the response stream w/ error: '%s'", err.Error())
						logger.Error(errMsg)
						return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
					}
				case types.EntityKindGraph:
					// TODO: Revisit later - Should we break after throwing/logging an error
					continue
				}
			} else {
				logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:populateResultSetFromGetEntityResponse - Received invalid entity kind %d", kindId))
			} // Valid entity types
		} // End of for loop
	}
	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:populateResultSetFromGetEntityResponse w/ Entity: '%+v'", entityFound))
	return entityFound, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAdminConnection
/////////////////////////////////////////////////////////////////

// CheckpointServer allows the programmatic control to do a checkpoint on server
func (obj *AdminConnectionImpl) CheckpointServer() types.TGError {
	/**
	Object obj = getCommandWithoutParameters (TGAdminCommand.CheckpointServer);
	 */
	return nil
}

// DumpServerStackTrace prints the stack trace
func (obj *AdminConnectionImpl) DumpServerStackTrace() types.TGError {
	/**
	connPool.adminlock();

	try {
		DumpStacktraceRequest request = (DumpStacktraceRequest) TGMessageFactory.getInstance().createMessage(VerbId.DumpStacktraceRequest);
		this.channel.sendMessage(request);
	}
	finally {
		connPool.adminUnlock();
	}
	 */
	return nil
}

// GetAttributeDescriptors gets the list of attribute descriptors
func (obj *AdminConnectionImpl) GetAttributeDescriptors() ([]types.TGAttributeDescriptor, types.TGError) {
	/**
	return ((Collection<TGAttributeDescriptor>) getCommandWithoutParameters (TGAdminCommand.ShowAttrDescs));
	 */
	return nil, nil
}

// GetConnections gets the list of all socket connections using this connection type
func (obj *AdminConnectionImpl) GetConnections() ([]admin.TGConnectionInfo, types.TGError) {
	/**
	return ((Collection<TGConnectionInfo>) getCommandWithoutParameters (TGAdminCommand.ShowConnections));
	 */
	return nil, nil
}

// GetIndices gets the list of all indices
func (obj *AdminConnectionImpl) GetIndices() ([]admin.TGIndexInfo, types.TGError) {
	/**
	return ((Collection<TGIndexInfo>) getCommandWithoutParameters (TGAdminCommand.ShowIndices));
	 */
	return nil, nil
}

// GetInfo gets the information about this connection type
func (obj *AdminConnectionImpl) GetInfo() (admin.TGServerInfo, types.TGError) {
	/**
	return ((TGServerInfo) getCommandWithoutParameters (TGAdminCommand.ShowInfo));
	 */
	return nil, nil
}

// GetUsers gets the list of users
func (obj *AdminConnectionImpl) GetUsers() ([]admin.TGUserInfo, types.TGError) {
	/**
	return ((Collection<TGUserInfo>) getCommandWithoutParameters (TGAdminCommand.ShowUsers));
	 */
	return nil, nil
}

// KillConnection terminates the connection forcefully
func (obj *AdminConnectionImpl) KillConnection(sessionId int64) types.TGError {
	/**
	getCommandWithParameters(TGAdminCommand.KillConnection, new Long (sessionId));
	 */
	return nil
}

// SetServerLogLevel set the log level
func (obj *AdminConnectionImpl) SetServerLogLevel(logLevel int, logComponent int64) types.TGError {
	/**
	TGServerLogDetails logDetails = new TGServerLogDetails(logComponent, logLevel);
	getCommandWithParameters (TGAdminCommand.SetLogLevel, logDetails);
	 */
	return nil
}

// StopServer stops the admin connection
func (obj *AdminConnectionImpl) StopServer() types.TGError {
	/**
	getCommandWithoutParameters (TGAdminCommand.StopServer);
	 */
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConnection
/////////////////////////////////////////////////////////////////

// Commit commits the current transaction on this connection
func (obj *AdminConnectionImpl) Commit() (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprint("Entering AdminConnectionImpl:Commit"))
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Commit - about to loop through addedList to include existing nodes to the changed list if it's part of a new edge"))
	// Include existing nodes to the changed list if it's part of a new edge
	for _, addEntity := range obj.GetAddedList() {
		if addEntity.GetEntityKind() == types.EntityKindEdge {
			nodes := addEntity.(*model.Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*model.Node)
					if !node.GetIsNew() {
						obj.changedList[node.GetVirtualId()] = node
						logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:Commit - Existing node '%d' added to change list for a new edge", node.GetVirtualId()))
					}
				}
			}
		}
	}

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Commit - about to loop through changedList to include existing nodes to the changed list even for edge update"))
	// Need to include existing node to the changed list even for edge update
	for _, modEntity := range obj.GetChangedList() {
		if modEntity.GetEntityKind() == types.EntityKindEdge {
			nodes := modEntity.(*model.Edge).GetVertices()
			if len(nodes) > 0 {
				for _, vNode := range nodes {
					node := vNode.(*model.Node)
					if !node.GetIsNew() {
						obj.changedList[node.GetVirtualId()] = node
						logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:Commit - Existing node '%d' added to change list for an existing edge '%d'", node.GetVirtualId(), modEntity.GetVirtualId()))
					}
				}
			}
		}
	}

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Commit - about to loop through removedList to include existing nodes to the changed list even for edge update"))
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
							logger.Warning(fmt.Sprintf("WARNING: AdminConnectionImpl:Commit - Existing node '%d' added to change list for an edge %d to be deleted", node.GetVirtualId(), delEntity.GetVirtualId()))
						}
					}
				}
			}
		}
	}
	//For deleted edge and node, we don't immediately change the effected nodes or edges.

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Commit about to createChannelRequest() for: pdu.VerbCommitTransactionRequest"))
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, pdu.VerbCommitTransactionRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:Commit - unable to createChannelRequest(pdu.VerbCommitTransactionRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*pdu.CommitTransactionRequest)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Commit about to GetAttributeDescriptors() for: pdu.VerbCommitTransactionRequest"))
	gof := obj.graphObjFactory
	attrDescSet, aErr := gof.GetGraphMetaData().GetNewAttributeDescriptors()
	if aErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:Commit - unable to gmd.GetAttributeDescriptors() w/ error: '%s'", aErr.Error()))
		return nil, aErr
	}
	queryRequest.AddCommitLists(obj.addedList, obj.changedList, obj.removedList, attrDescSet)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Commit about to channelSendRequest() for: pdu.VerbCommitTransactionRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:Commit - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::Commit received response for: pdu.VerbCommitTransactionRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.CommitTransactionResponse)

	if response.HasException() {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:Commit as response has exceptions"))
		// TODO: Revisit later - Should we not throw an appropriate exception?
		return nil, nil
	}

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Commit about to fixUpAttrDescriptors()"))
	fixUpAttrDescriptors(response, attrDescSet)
	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Commit about to obj.fixUpEntities()"))
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

	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:Commit"))
	return nil, nil
}

// Connect establishes a network connection to the TGDB server
func (obj *AdminConnectionImpl) Connect() types.TGError {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:Connect for connection: '%+v'", obj))
	err := obj.GetChannel().Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl::Connect - error in obj.GetChannel().Connect() as '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Connect about to obj.GetChannel().Start()"))
	err = obj.GetChannel().Start()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl::Connect - error in obj.GetChannel().Start() as '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:Connect"))
	return nil
}

// CloseQuery closes a specific query and associated objects
func (obj *AdminConnectionImpl) CloseQuery(queryHashId int64) (types.TGQuery, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:CloseQuery for QueryHashId: '%+v'", queryHashId))
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::CloseQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, pdu.VerbQueryRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:CloseQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(CLOSE)
	queryRequest.SetQueryHashId(queryHashId)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::CloseQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	_, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:CloseQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:CloseQuery"))
	return nil, nil
}

// CreateQuery creates a reusable query object that can be used to execute one or more statement
func (obj *AdminConnectionImpl) CreateQuery(expr string) (types.TGQuery, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:CreateQuery for Query: '%+v'", expr))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::CreateQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, err := createChannelRequest(obj, pdu.VerbQueryRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:CreateQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", err.Error()))
		return nil, err
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(CREATE)
	queryRequest.SetQuery(expr)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::CreateQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:CreateQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::CreateQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))

	response := msgResponse.(*pdu.QueryResponseMessage)
	queryHashId := response.GetQueryHashId()

	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:CreateQuery"))
	if response.GetResult() == 0 && queryHashId > 0 {
		return query.NewQuery(obj, queryHashId), nil
	}
	return nil, nil
}

// DecryptBuffer decrypts the encrypted buffer by sending a DecryptBufferRequest to the server
func (obj *AdminConnectionImpl) DecryptBuffer(is types.TGInputStream) ([]byte, types.TGError) {
	logger.Log(fmt.Sprint("Entering AdminConnectionImpl:DecryptBuffer w/ EncryptedBuffer"))
	//err := obj.InitMetadata()
	//if err != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:DecryptBuffer - unable to initialize metadata w/ error: '%s'", err.Error()))
	//	return nil, err
	//}
	//obj.connPoolImpl.AdminLock()
	//defer obj.connPoolImpl.AdminUnlock()
	//
	//logger.Log(fmt.Sprint("Inside AdminConnectionImpl::DecryptBuffer about to createChannelRequest() for: pdu.VerbDecryptBufferRequest"))
	//// Create a channel request
	//msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbDecryptBufferRequest)
	//if cErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:DecryptBuffer - unable to createChannelRequest(pdu.VerbDecryptBufferRequest w/ error: '%s'", cErr.Error()))
	//	return nil, cErr
	//}
	//decryptRequest := msgRequest.(*pdu.DecryptBufferRequestMessage)
	//decryptRequest.SetEncryptedBuffer(encryptedBuf)
	//
	//logger.Log(fmt.Sprint("Inside AdminConnectionImpl::DecryptBuffer about to obj.GetChannel().SendRequest() for: pdu.VerbDecryptBufferRequest"))
	//// Execute request on channel and get the response
	//msgResponse, channelErr := obj.GetChannel().SendRequest(decryptRequest, channelResponse.(*channel.BlockingChannelResponse))
	//if channelErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:DecryptBuffer - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
	//	return nil, channelErr
	//}
	//logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::DecryptBuffer received response for: pdu.VerbGetLargeObjectRequest as '%+v'", msgResponse))
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
func (obj *AdminConnectionImpl) DecryptEntity(entityId int64) ([]byte, types.TGError) {
	buf, err := obj.GetLargeObjectAsBytes(entityId, true)
	if err != nil {
		return nil, err
	}
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Decrypt(iostream.NewProtocolDataInputStream(buf))
}

// DeleteEntity marks an ENTITY for delete operation. Upon commit, the entity will be deleted from the database
func (obj *AdminConnectionImpl) DeleteEntity(entity types.TGEntity) types.TGError {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:DeleteEntity for Entity: '%+v'", entity))
	obj.removedList[entity.GetVirtualId()] = entity
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:DeleteEntity"))
	return nil
}

// Disconnect breaks the connection from the TGDB server
func (obj *AdminConnectionImpl) Disconnect() types.TGError {
	logger.Log(fmt.Sprint("Entering AdminConnectionImpl:Disconnect"))
	err := obj.GetChannel().Disconnect()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:Disconnect - unable to channel.Disconnect() w/ error: '%s'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::Disconnect about to obj.GetChannel().Stop()"))
	obj.GetChannel().Stop(false)
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:Disconnect"))
	return nil
}

// EncryptEntity encrypts the encrypted entity using channel's data cryptographer
func (obj *AdminConnectionImpl) EncryptEntity(rawBuffer []byte) ([]byte, types.TGError) {
	cryptoGrapher := obj.GetChannel().GetDataCryptoGrapher()
	return cryptoGrapher.Encrypt(rawBuffer)
}

// ExecuteGremlinQuery executes a Gremlin Grammer-Based query with  query options
func (obj *AdminConnectionImpl) ExecuteGremlinQuery(expr string, collection []interface{}, options types.TGQueryOption) ([]interface{}, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteGremlinQuery for Query: '%+v'", expr))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteGremlinQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteGremlinQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery("gbc : " + expr)
	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteGremlinQuery about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteGremlinQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteGremlinQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteGremlinQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::ExecuteGremlinQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	// TODO: Revisit later once Gremlin Package and GremlinQueryResult are implemented
	//respStream := response.GetEntityStream()
	//GremlinResult.fillCollection(entityStream, gof, collection);
	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:ExecuteGremlinQuery w/ '%+v'", response))
	return collection, nil
}

// ExecuteQuery executes an immediate query with associated query options
func (obj *AdminConnectionImpl) ExecuteQuery(expr string, options types.TGQueryOption) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteQuery for Query: '%+v'", expr))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQuery about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQuery - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery(expr)
	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQuery about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQuery about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQuery - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteQuery received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::ExecuteQuery - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:ExecuteQuery w/ '%+v'", response))
	return obj.populateResultSetFromQueryResponse(0, response)
}

// ExecuteQueryWithFilter executes an immediate query with specified filter & query options
// The query option is place holder at this time
// @param expr A subset of SQL-92 where clause
// @param edgeFilter filter used for selecting edges to be returned
// @param traversalCondition condition used for selecting edges to be traversed and returned
// @param endCondition condition used to stop the traversal
// @param option Query options for executing. Can be null, then it will use the default option
func (obj *AdminConnectionImpl) ExecuteQueryWithFilter(expr, edgeFilter, traversalCondition, endCondition string, options types.TGQueryOption) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteQueryWithFilter for Query: '%+v', EdgeFilter: '%+v', Traversal: '%+v', EndCondition: '%+v'", expr, edgeFilter, traversalCondition, endCondition))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithFilter about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQueryWithFilter - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTE)
	queryRequest.SetQuery(expr)
	queryRequest.SetEdgeFilter(edgeFilter)
	queryRequest.SetTraversalCondition(traversalCondition)
	queryRequest.SetEndCondition(endCondition)
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteQueryWithFilter about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithFilter about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQueryWithFilter - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::ExecuteQueryWithFilter received response for: pdu.VerbQueryRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.QueryResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::ExecuteQueryWithFilter - The query does not have any results in QueryResponseMessage"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:ExecuteQueryWithFilter w/ '%+v'", response))
	return obj.populateResultSetFromQueryResponse(0, response)
}

// ExecuteQueryWithId executes an immediate query for specified id & query options
func (obj *AdminConnectionImpl) ExecuteQueryWithId(queryHashId int64, options types.TGQueryOption) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:ExecuteQueryWithId for QueryHashId: '%+v'", queryHashId))
	err := obj.InitMetadata()
	if err != nil {
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithId about to createChannelRequest() for: pdu.VerbQueryRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbQueryRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQueryWithId - unable to createChannelRequest(pdu.VerbQueryRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.QueryRequestMessage)
	queryRequest.SetCommand(EXECUTED)
	queryRequest.SetQueryHashId(queryHashId)
	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithId about to obj.configureQueryRequest() for: pdu.VerbQueryRequest"))
	configureQueryRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::ExecuteQueryWithId about to obj.GetChannel().SendRequest() for: pdu.VerbQueryRequest"))
	// Execute request on channel and get the response
	_, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:ExecuteQueryWithId - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}

	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:ExecuteQueryWithId"))
	return nil, nil
}

// GetAddedList gets a list of added entities
func (obj *AdminConnectionImpl) GetAddedList() map[int64]types.TGEntity {
	return obj.addedList
}

// GetChangedList gets a list of changed entities
func (obj *AdminConnectionImpl) GetChangedList() map[int64]types.TGEntity {
	return obj.changedList
}

// GetChangedList gets the communication channel associated with this connection
func (obj *AdminConnectionImpl) GetChannel() types.TGChannel {
	return obj.channel
}

// GetConnectionId gets connection identifier
func (obj *AdminConnectionImpl) GetConnectionId() int64 {
	return obj.connId
}

// GetConnectionProperties gets a list of connection properties
func (obj *AdminConnectionImpl) GetConnectionProperties() types.TGProperties {
	return obj.connProperties
}

// GetEntities gets a result set of entities given an non-uniqueKey
func (obj *AdminConnectionImpl) GetEntities(qryKey types.TGKey, props types.TGProperties) (types.TGResultSet, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:GetEntities for QueryKey: '%+v'", qryKey))
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:GetEntities - unable to InitMetadata"))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if props == nil {
		props = query.NewQueryOption()
	}
	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetEntities about to createChannelRequest() for: pdu.VerbGetEntityRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbGetEntityRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetEntities - unable to createChannelRequest(pdu.VerbGetEntityRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.GetEntityRequestMessage)
	queryRequest.SetCommand(2)
	queryRequest.SetKey(qryKey)
	configureGetRequest(queryRequest, props)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetEntities about to obj.GetChannel().SendRequest() for: pdu.VerbGetEntityRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetEntities - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::GetEntities received response for: pdu.VerbGetEntityRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.GetEntityResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::GetEntities - The request does not have any results in GetEntityResponseMessage"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:GetEntities w/ '%+v'", response))
	return obj.populateResultSetFromGetEntitiesResponse(response)
}

// GetEntity gets an Entity given an UniqueKey for the Object
func (obj *AdminConnectionImpl) GetEntity(qryKey types.TGKey, options types.TGQueryOption) (types.TGEntity, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:GetEntity for QueryKey: '%+v'", qryKey))
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminConnectionImpl:GetEntity - unable to InitMetadata"))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	if options == nil {
		options = query.NewQueryOption()
	}
	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetEntity about to createChannelRequest() for: pdu.VerbGetEntityRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbGetEntityRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetEntity - unable to createChannelRequest(pdu.VerbGetEntityRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.GetEntityRequestMessage)
	queryRequest.SetCommand(0)
	queryRequest.SetKey(qryKey)
	configureGetRequest(queryRequest, options)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetEntity about to obj.GetChannel().SendRequest() for: pdu.VerbGetEntityRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetEntity - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::GetEntity received response for: pdu.VerbGetEntityRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.GetEntityResponseMessage)

	if !response.GetHasResult() {
		logger.Warning(fmt.Sprint("WARNING: Returning AdminConnectionImpl::GetEntity - The request does not have any results in GetEntityResponseMessage"))
		return nil, nil
	}

	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:GetEntity w/ '%+v'", response))
	return obj.populateResultSetFromGetEntityResponse(response)
}

// GetGraphMetadata gets the Graph Metadata
func (obj *AdminConnectionImpl) GetGraphMetadata(refresh bool) (types.TGGraphMetadata, types.TGError) {
	logger.Log(fmt.Sprint("Entering AdminConnectionImpl:GetGraphMetadata"))
	if refresh {
		obj.connPoolImpl.AdminLock()
		defer obj.connPoolImpl.AdminUnlock()

		logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetGraphMetadata about to createChannelRequest() for: pdu.VerbMetadataRequest"))
		// Create a channel request
		msgRequest, channelResponse, err := createChannelRequest(obj, pdu.VerbMetadataRequest)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetGraphMetadata - unable to createChannelRequest(pdu.VerbMetadataRequest w/ error: '%s'", err.Error()))
			return nil, err
		}
		metaRequest := msgRequest.(*pdu.MetadataRequest)
		logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::GetGraphMetadata createChannelRequest() returned MsgRequest: '%+v' ChannelResponse: '%+v'", msgRequest, channelResponse.(*channel.BlockingChannelResponse)))

		logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetGraphMetadata about to obj.GetChannel().SendRequest() for: pdu.VerbMetadataRequest"))
		// Execute request on channel and get the response
		msgResponse, channelErr := obj.GetChannel().SendRequest(metaRequest, channelResponse.(*channel.BlockingChannelResponse))
		if channelErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetGraphMetadata - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
			return nil, channelErr
		}
		//logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::GetGraphMetadata received response for: pdu.VerbMetadataRequest as '%+v'", msgResponse))

		response := msgResponse.(*pdu.MetadataResponse)
		attrDescList := response.GetAttrDescList()
		edgeTypeList := response.GetEdgeTypeList()
		nodeTypeList := response.GetNodeTypeList()

		gmd := obj.graphObjFactory.GetGraphMetaData()
		logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetGraphMetadata about to update GraphMetadata"))
		uErr := gmd.UpdateMetadata(attrDescList, nodeTypeList, edgeTypeList)
		if uErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetGraphMetadata - unable to gmd.UpdateMetadata() w/ error: '%s'", uErr.Error()))
			return nil, uErr
		}
		logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetGraphMetadata successfully updated GraphMetadata"))
	}
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:GetGraphMetadata"))
	return obj.graphObjFactory.GetGraphMetaData(), nil
}

// GetGraphObjectFactory gets the Graph Object Factory for Object creation
func (obj *AdminConnectionImpl) GetGraphObjectFactory() (types.TGGraphObjectFactory, types.TGError) {
	logger.Log(fmt.Sprint("Entering AdminConnectionImpl:GetGraphObjectFactory"))
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetGraphObjectFactory - unable to initialize metadata w/ error: '%s'", err.Error()))
		return nil, err
	}
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:GetGraphObjectFactory"))
	return obj.graphObjFactory, nil
}

// GetLargeObjectAsBytes gets an Binary Large Object Entity given an UniqueKey for the Object
func (obj *AdminConnectionImpl) GetLargeObjectAsBytes(entityId int64, decryptFlag bool) ([]byte, types.TGError) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:GetLargeObjectAsBytes for EntityId: '%+v'", entityId))
	err := obj.InitMetadata()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetLargeObjectAsBytes - unable to initialize metadata w/ error: '%s'", err.Error()))
		return nil, err
	}
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetLargeObjectAsBytes about to createChannelRequest() for: pdu.VerbGetLargeObjectRequest"))
	// Create a channel request
	msgRequest, channelResponse, cErr := createChannelRequest(obj, pdu.VerbGetLargeObjectRequest)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetLargeObjectAsBytes - unable to createChannelRequest(pdu.VerbGetLargeObjectRequest w/ error: '%s'", cErr.Error()))
		return nil, cErr
	}
	queryRequest := msgRequest.(*pdu.GetLargeObjectRequestMessage)
	queryRequest.SetEntityId(entityId)
	queryRequest.SetDecryption(decryptFlag)

	logger.Log(fmt.Sprint("Inside AdminConnectionImpl::GetLargeObjectAsBytes about to obj.GetChannel().SendRequest() for: pdu.VerbGetLargeObjectRequest"))
	// Execute request on channel and get the response
	msgResponse, channelErr := obj.GetChannel().SendRequest(queryRequest, channelResponse.(*channel.BlockingChannelResponse))
	if channelErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminConnectionImpl:GetLargeObjectAsBytes - unable to channel.SendRequest() w/ error: '%s'", channelErr.Error()))
		return nil, channelErr
	}
	logger.Log(fmt.Sprintf("Inside AdminConnectionImpl::GetLargeObjectAsBytes received response for: pdu.VerbGetLargeObjectRequest as '%+v'", msgResponse))
	response := msgResponse.(*pdu.GetLargeObjectResponseMessage)

	if response == nil {
		errMsg := "AdminConnectionImpl::GetLargeObjectAsBytes does not have any results in GetLargeObjectResponseMessage"
		logger.Error(errMsg)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	logger.Log(fmt.Sprintf("Returning AdminConnectionImpl:GetLargeObjectAsBytes w/ '%+v'", response.GetBuffer()))
	return response.GetBuffer(), nil
}

// GetRemovedList gets a list of removed entities
func (obj *AdminConnectionImpl) GetRemovedList() map[int64]types.TGEntity {
	return obj.removedList
}

// InsertEntity marks an ENTITY for insert operation. Upon commit, the entity will be inserted in the database
func (obj *AdminConnectionImpl) InsertEntity(entity types.TGEntity) types.TGError {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:InsertEntity to insert Entity: '%+v'", entity.GetEntityType()))
	obj.addedList[entity.GetVirtualId()] = entity
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:InsertEntity"))
	return nil
}

// Rollback rolls back the current transaction on this connection
func (obj *AdminConnectionImpl) Rollback() types.TGError {
	logger.Log(fmt.Sprint("Entering AdminConnectionImpl:Rollback"))
	obj.connPoolImpl.AdminLock()
	defer obj.connPoolImpl.AdminUnlock()

	// Reset all the lists to empty contents
	obj.addedList = make(map[int64]types.TGEntity, 0)
	obj.changedList = make(map[int64]types.TGEntity, 0)
	obj.removedList = make(map[int64]types.TGEntity, 0)
	obj.attrByTypeList = make(map[int][]types.TGAttribute, 0)

	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:Rollback"))
	return nil
}

// SetExceptionListener sets exception listener
func (obj *AdminConnectionImpl) SetExceptionListener(listener types.TGConnectionExceptionListener) {
	obj.connPoolImpl.SetExceptionListener(listener) //delegate it to the Pool.
}

// UpdateEntity marks an ENTITY for update operation. Upon commit, the entity will be updated in the database
// When commit is called, the object is resolved to check if it is dirty. Entity.setAttribute calls make the entity
// dirty. If it is dirty, then the object is send to the server for update, otherwise it is ignored.
// Calling multiple times, does not change the behavior.
// The same entity cannot be updated on multiple connections. It will result an TGException of already associated to a connection.
func (obj *AdminConnectionImpl) UpdateEntity(entity types.TGEntity) types.TGError {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:UpdateEntity to update Entity: '%+v'", entity))
	obj.changedList[entity.GetVirtualId()] = entity
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:UpdateEntity"))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChangeListener
/////////////////////////////////////////////////////////////////

// AttributeAdded gets called when an attribute is Added to an entity.
func (obj *AdminConnectionImpl) AttributeAdded(attr types.TGAttribute, owner types.TGEntity) {
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:AttributeAdded"))
}

// AttributeChanged gets called when an attribute is set.
func (obj *AdminConnectionImpl) AttributeChanged(attr types.TGAttribute, oldValue, newValue interface{}) {
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:AttributeChanged"))
}

// AttributeRemoved gets called when an attribute is removed from the entity.
func (obj *AdminConnectionImpl) AttributeRemoved(attr types.TGAttribute, owner types.TGEntity) {
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:AttributeRemoved"))
}

// EntityCreated gets called when an entity is Added
func (obj *AdminConnectionImpl) EntityCreated(entity types.TGEntity) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:EntityCreated to add Entity: '%+v'", entity))
	entityId := entity.(*model.AbstractEntity).GetVirtualId()
	obj.addedList[entityId] = entity
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:EntityCreated"))
}

// EntityDeleted gets called when the entity is deleted
func (obj *AdminConnectionImpl) EntityDeleted(entity types.TGEntity) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:EntityDeleted to delete Entity: '%+v'", entity))
	entityId := entity.(*model.AbstractEntity).GetVirtualId()
	obj.removedList[entityId] = entity
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:EntityDeleted"))
}

// NodeAdded gets called when a node is Added
func (obj *AdminConnectionImpl) NodeAdded(graph types.TGGraph, node types.TGNode) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:NodeAdded to add Node: '%+v' to Graph: '%+v'", node, graph))
	entityId := graph.(*model.Graph).GetVirtualId()
	obj.addedList[entityId] = graph
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:NodeAdded"))
}

// NodeRemoved gets called when a node is removed
func (obj *AdminConnectionImpl) NodeRemoved(graph types.TGGraph, node types.TGNode) {
	logger.Log(fmt.Sprintf("Entering AdminConnectionImpl:NodeRemoved to remove Node: '%+v' to Graph: '%+v'", node, graph))
	entityId := graph.(*model.Graph).GetVirtualId()
	obj.removedList[entityId] = graph
	logger.Log(fmt.Sprint("Returning AdminConnectionImpl:NodeRemoved"))
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
