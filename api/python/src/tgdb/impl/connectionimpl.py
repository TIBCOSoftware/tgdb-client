"""
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
 *  File name :connectionimpl.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: connectionimpl.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates connection interfaces
 """

from tgdb.log import *
import tgdb.log as tglog
from tgdb.utils import *
from tgdb.impl.atomics import *
import typing

import tgdb.channel as tgchannel
import tgdb.impl.channelimpl as tgchannelimpl
import tgdb.pdu as tgpdu
import tgdb.impl.pduimpl as tgpduimpl
import tgdb.connection as tgconn
import tgdb.model as tgmodel
import tgdb.impl.entityimpl as tgentimpl
import tgdb.impl.gmdimpl as tggmdimpl
import tgdb.query as tgquery
import tgdb.impl.queryimpl as tgqueryimpl
import tgdb.exception as tgexception
import tgdb.bulkio as tgbulk
import tgdb.admin as tgadm


def findCommandForLang(lang: str) -> tgquery.TGQueryCommand:
    retCommand: tgquery.TGQueryCommand
    if lang == "tgql":
        retCommand = tgquery.TGQueryCommand.Execute
    elif lang == "gremlin":
        retCommand = tgquery.TGQueryCommand.ExecuteGremlinStr
    elif lang == "gbc":
        retCommand = tgquery.TGQueryCommand.ExecuteGremlin
    else:
        raise tgexception.TGException("Unknown property for ConnectionDefaultQueryLanguage: %s", lang)
    return retCommand


def findCommandAndQueryString(query: str, props: tgchannel.TGProperties) -> typing.Tuple[tgquery.TGQueryCommand, str]:
    lang: str = props.get(ConfigName.ConnectionDefaultQueryLanguage,
                          ConfigName.ConnectionDefaultQueryLanguage.defaultvalue)
    retCommand: tgquery.TGQueryCommand
    retStr = query
    try:
        idx: int = query.index("://")
        prefix = query[:idx].lower()
        retCommand = findCommandForLang(prefix)
        retStr = query[idx + 3:]
    except ValueError:
        lang = lang.lower()
        retCommand = findCommandForLang(lang)
    return retCommand, retStr


class ConnectionImpl(tgconn.TGConnection):

    def __init__(self, url, username, password, dbName: typing.Optional[str], env):
        self.__url__ = url
        self.__username__ = username
        self.__password__ = password
        self.__props__: TGProperties = TGProperties(env)
        self._dbName = dbName
        self.__channel__: tgchannel.TGChannel = tgchannel.TGChannel.createChannel(url, username, password, dbName,
                                                                                  self.__props__)
        self.__props__.update(tgchannelimpl.LinkUrl.parse(url).properties)
        self.__gof__: tggmdimpl.GraphObjectFactoryImpl = tggmdimpl.GraphObjectFactoryImpl(self)
        self.__addEntities__: typing.Dict[int, tgentimpl.AbstractEntity] = {}
        self.__updateEntities__: typing.Dict[int, tgentimpl.AbstractEntity] = {}
        self.__removeEntities__: typing.Dict[int, tgentimpl.AbstractEntity] = {}
        self.__requestIds__ = AtomicReference('i', 0)

    def _genBCRWaiter(self) -> tgchannelimpl.BlockingChannelResponseWaiter:
        timeout = self.__props__.get(ConfigName.ConnectionOperationTimeoutSeconds, None)
        if timeout is not None and isinstance(timeout, str):
            timeout = float(timeout)
        requestId = self.__requestIds__.increment()
        return tgchannelimpl.BlockingChannelResponseWaiter(requestId, timeout)

    def connect(self):
        tglog.gLogger.log(tglog.TGLevel.Debug, "Attempting to connect")
        self.__channel__.connect()
        tglog.gLogger.log(tglog.TGLevel.Debug, "Connected, now logging in.")
        self.__channel__.start()
        tglog.gLogger.log(tglog.TGLevel.Debug, "Logged in, now acquiring metadata.")
        self.__initMetadata__()
        tglog.gLogger.log(tglog.TGLevel.Debug, "Acquired metadata, now sending connection properties.")
        self.__sendConnectionProperties()
        tglog.gLogger.log(tglog.TGLevel.Debug, 'Connected successfully')

    def __initMetadata__(self):
        waiter = self._genBCRWaiter()
        request = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.MetadataRequest,
                                           authtoken=self.__channel__.authtoken, sessionid=self.__channel__.sessionid)
        response = self.__channel__.send(request, waiter)
        if response.verbid != tgpdu.VerbId.MetadataResponse:
            raise tgexception.TGException('Invalid response object received')
        self.__gof__.graphmetadata.registry = response.typeregistry

    def disconnect(self):
        self.__channel__.disconnect()
        self.__channel__.stop()

    def commit(self):
        channelResponse = self._genBCRWaiter()

        try:
            if gLogger.level is TGLevel.Debug:
                def echoAttributes(ent: tgmodel.TGEntity):
                    gLogger.log(TGLevel, "Entity ID: %d", ent.virtualId)
                    attr: tgmodel.TGAttribute
                    for attr in ent.attributes:
                        gLogger.log(TGLevel, " Attribute: %s", attr._value)

                [echoAttributes(ent) for ent in self.__addEntities__.values()]
                [echoAttributes(ent) for ent in self.__updateEntities__.values()]
                [echoAttributes(ent) for ent in self.__removeEntities__.values()]

            request: tgpduimpl.CommitTransactionRequestMessage = tgpduimpl.TGMessageFactory.createMessage(
                                                tgpdu.VerbId.CommitTransactionRequest, authtoken=self.__channel__.authtoken,
                                                sessionid=self.__channel__.sessionid)
            attrDescSet = self.graphObjectFactory.graphmetadata.attritubeDescriptors

            request.addCommitList(self.__addEntities__, self.__updateEntities__, self.__removeEntities__, attrDescSet)

            response: tgpduimpl.CommitTransactionResponseMessage = self.__channel__.send(request, channelResponse)

            if response.exception is not None:
                 raise response.exception

            response.finishReadWith(self.__addEntities__, self.__updateEntities__, self.__removeEntities__,
                                    self.__gof__.graphmetadata.registry)

            for id in self.__removeEntities__:
                self.__removeEntities__[id].markDeleted()

            if gLogger.isEnabled(TGLevel.Debug):
                gLogger.log(TGLevel.Debug, "Transaction commit succeeded")

        except IOError as e:
            raise tgexception.TGException.buildException("IO Error", cause=e)
        finally:
            for id in self.__addEntities__:
                self.__addEntities__[id].resetModifiedAttributes()

            for id in self.__updateEntities__:
                self.__updateEntities__[id].resetModifiedAttributes()

            self.__addEntities__.clear()
            self.__updateEntities__.clear()
            self.__removeEntities__.clear()

    def refreshMetadata(self):
        self.__initMetadata__()

    def rollback(self):
        self.__addEntities__.clear()
        self.__updateEntities__.clear()
        self.__removeEntities__.clear()

    def __sendConnectionProperties(self):
        request: tgpduimpl.ConnectionPropertiesMessage = tgpduimpl.TGMessageFactory.createMessage(
                                                tgpdu.VerbId.ConnectionPropertiesMessage, authtoken=self.__channel__.authtoken,
                                                sessionid=self.__channel__.sessionid)
        request.props = self.__channel__.properties
        self.__channel__.send(request)

    """
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//          Begin Bulk Import Stuff                                                                                   //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    """

    def startImport(self, loadopt: typing.Union[str, tgbulk.TGLoadOptions] = tgbulk.TGLoadOptions.Insert,
                    erroropt: typing.Union[str, tgbulk.TGErrorOptions] = tgbulk.TGErrorOptions.Stop,
                    dateformat: typing.Union[str, tgbulk.TGDateFormat] = tgbulk.TGDateFormat.YMD,
                    props: typing.Optional[TGProperties] = None):
        import tgdb.impl.bulkioimpl as tgbulkimpl
        ret: tgbulkimpl.BulkImportImpl
        channelResponseWaiter = self._genBCRWaiter()
        request: tgpduimpl.BeginImportSessionRequest
        request = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.BeginImportRequest,
                                                           authtoken=self.__channel__.authtoken, sessionid=self.__channel__.sessionid)
        if isinstance(loadopt, str):
            loadopt = tgbulk.TGErrorOptions.findVal(loadopt)
        if loadopt == tgbulk.TGLoadOptions.Invalid:
            raise tgexception.TGException("Bad argument: cannot have an invalid load option!")
        if isinstance(erroropt, str):
            erroropt = tgbulk.TGErrorOptions.findVal(erroropt)
        if erroropt == tgbulk.TGErrorOptions.Invalid:
            raise tgexception.TGException("Bad argument: cannot have an invalid error option!")
        if isinstance(dateformat, str):
            dateformat = tgbulk.TGDateFormat.findVal(dateformat)
        if dateformat == tgbulk.TGDateFormat.Invalid:
            raise tgexception.TGException("Bad argument: cannot have an invalid Date-Time Format!")

        request.loadopt = loadopt
        request.erroropt = erroropt
        request.dtformat = dateformat

        response: tgpduimpl.BeginImportSessionResponse = self.__channel__.send(request, channelResponseWaiter)

        if response.error is not None:
            raise response.error

        ret = tgbulkimpl.BulkImportImpl(self, props)

        return ret

    def partialImportEntity(self, entType: tgmodel.TGEntityType, reqIdx: int, totReqs: int, data: str,
                            attrList: typing.List[str]) -> typing.List[tgadm.TGImportDescriptor]:
        channelResponseWaiter = self._genBCRWaiter()
        request: tgpduimpl.PartialImportRequest = tgpduimpl.TGMessageFactory.createMessage(
                                            tgpdu.VerbId.PartialImportRequest,
                                            authtoken=self.__channel__.authtoken, sessionid=self.__channel__.sessionid)
        request.type = entType
        request.reqIdx = reqIdx
        request.totalRequestsForType = totReqs
        request.data = data
        request.attrList = attrList
        response: tgpduimpl.PartialImportResponse = self.__channel__.send(request, channelResponseWaiter)

        if response.error is not None:
            raise response.error

        return response.resultList

    def endBulkImport(self):
        channelResponseWaiter = self._genBCRWaiter()
        request: tgpduimpl.EndBulkImportSessionRequest = tgpduimpl.TGMessageFactory.createMessage(
            tgpdu.VerbId.EndImportRequest,
            authtoken=self.__channel__.authtoken, sessionid=self.__channel__.sessionid)
        response: tgpduimpl.PartialImportResponse = self.__channel__.send(request, channelResponseWaiter)
        return response.resultList

    """
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//          End Bulk Import Stuff                                                                                     //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    """

    """
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//          Begin Bulk Export Stuff                                                                                   //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    """

    def startExport(self, props: typing.Optional[TGProperties] = None, zip: typing.Optional[str] = None,
                    isBatch: bool = True):
        import tgdb.impl.bulkioimpl as tgbulkimpl
        channelResponseWaiter = self._genBCRWaiter()

        request: tgpduimpl.BeginExportRequest = tgpduimpl.TGMessageFactory.createMessage(
            tgpdu.VerbId.BeginExportRequest, authtoken=self.__channel__.authtoken, sessionid=self.__channel__.sessionid)

        request.zipName = zip
        request.isBatch = isBatch
        request.maxBatchEntities = int(ConfigName.BulkIOEntityBatchSize.defaultvalue)\
                    if props is None or props[ConfigName.BulkIOEntityBatchSize] is None else\
                    int(props[ConfigName.BulkIOEntityBatchSize])

        response: tgpduimpl.BeginExportResponse = self.__channel__.send(request, channelResponseWaiter)

        if response.error is not None:
            raise response.error

        return tgbulkimpl.BulkExportImpl(self, props, response.typeList, response.numRequests)

    def partialExport(self, reqNum: int) -> typing.Tuple[str, bytes, bool, int,
                                                         typing.Optional[typing.Tuple[str, typing.List[str]]]]:
        channelResponseWaiter = self._genBCRWaiter()

        request: tgpduimpl.PartialExportRequest = tgpduimpl.TGMessageFactory.createMessage(
            tgpdu.VerbId.PartialExportRequest, authtoken=self.__channel__.authtoken,
            sessionid=self.__channel__.sessionid)

        request.requestNum = reqNum

        response: tgpduimpl.PartialExportResponse = self.__channel__.send(request, channelResponseWaiter)

        return response.fileName, response.data, response.hasMore, response.numEntities,\
               (response.typeName, response.attrList) if response.newType else None

    """
    def startExport(self, props: Optional[TGProperties] = None) -> tgbulk.TGBulkExport:
        channelResponseWaiter = self.__genBCRWaiter()
        request = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.BeginBulkExportSessionRequest,
                                                           authtoken=self.__channel__.authtoken,
                                                           sessionid=self.__channel__.sessionid)
        _ = self.__channel__.send(request, channelResponseWaiter)
        return tgbulkimpl.BulkExportImpl(self, props)

    def beginBatchExportEntity(self, entkind: tgmodel.TGEntityKind, enttype: tgmodel.TGEntityType, batchSize: int) \
            -> Tuple[int, List[str]]:
        channelResponseWaiter = self.__genBCRWaiter()
        request: tgpduimpl.BeginBatchExportEntityRequest = tgpduimpl.TGMessageFactory.createMessage(
            tgpdu.VerbId.BeginBatchExportEntityRequest, authtoken=self.__channel__.authtoken,
            sessionid=self.__channel__.sessionid)

        request.entKind = entkind
        request.entType = enttype
        request.batchSize = batchSize

        response: tgpduimpl.BeginBatchExportEntityResponse = self.__channel__.send(request, channelResponseWaiter)

        return response.descriptor, response.columnLabels

    def singleBatchExportEntity(self, desc: int) -> Tuple[int, str, bool]:
        channelResponseWaiter = self.__genBCRWaiter()
        request: tgpduimpl.SingleBatchExportEntityRequest = tgpduimpl.TGMessageFactory.createMessage(
            tgpdu.VerbId.SingleBatchExportEntityRequest, authtoken=self.__channel__.authtoken,
            sessionid=self.__channel__.sessionid)

        request.descriptor = desc

        response: tgpduimpl.SingleBatchExportEntityResponse = self.__channel__.send(request, channelResponseWaiter)

        return response.numEnts, response.data, response.hasMore

    def endBulkExportSession(self):
        channelResponseWaiter = self.__genBCRWaiter()
        request: tgpduimpl.EndBulkExportSessionRequest = tgpduimpl.TGMessageFactory.createMessage(
            tgpdu.VerbId.EndBulkExportSessionRequest, authtoken=self.__channel__.authtoken,
            sessionid=self.__channel__.sessionid)

        _ = self.__channel__.send(request, channelResponseWaiter)
    """

    """
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//          End Bulk Export Stuff                                                                                     //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    """

    def getEntity(self, key: tgmodel.TGKey, option: tgquery.TGQueryOption = tgquery.DefaultQueryOption) ->\
            tgmodel.TGEntity:
        channelResponseWaiter = self._genBCRWaiter()

        requestMessage: tgpduimpl.GetEntityRequestMessage
        retV: tgmodel.TGEntity = None
        try:
            requestMessage = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.GetEntityRequest,
                                                                      authtoken=self.__channel__.authtoken,
                                                                      sessionid=self.__channel__.sessionid)
            requestMessage.command = tgpduimpl.GetEntityCommand.GetEntity
            requestMessage.key = key
            response: tgpduimpl.GetEntityResponseMessage = self.__channel__.send(requestMessage, channelResponseWaiter)
            if response.hasResult:
                response.finishReadWith(self.graphObjectFactory)
                fetchedEntities = response.fetchedEntities
                for id in fetchedEntities:
                    fetchedEnt: tgmodel.TGEntity = fetchedEntities[id]
                    if key.matches(fetchedEnt):
                        retV = fetchedEnt
                        break
        finally:
            pass
        return retV

    def insertEntity(self, entity: tgmodel.TGEntity):
        if not entity.isNew:
            raise tgexception.TGException("Should only be calling insertEntity on a new entity!")
        if entity.virtualId not in self.__removeEntities__:
            self.__addEntities__[entity.virtualId] = entity
            self.__updateEdge__(entity)
            if gLogger.isEnabled(TGLevel.Debug):
                gLogger.log(TGLevel.Debug, 'Insert entity called')

    def updateEntity(self, entity: tgmodel.TGEntity):
        if entity.isNew:
            raise tgexception.TGException('Should not be calling update on a new entity!')
        if entity.isDeleted:
            raise tgexception.TGException('Should not be calling update on an already deleted entity!')
        if entity.virtualId not in self.__removeEntities__:
            self.__updateEntities__[entity.virtualId] = entity
            self.__updateEdge__(entity)

    def __updateEdge__(self, entity: tgmodel.TGEntity):
        if isinstance(entity, tgentimpl.EdgeImpl):
            edge: tgmodel.TGEdge = entity
            fr, to = edge.vertices
            if not fr.isNew and fr.virtualId not in self.__removeEntities__:
                self.__updateEntities__[fr.virtualId] = fr
            if not to.isNew and to.virtualId not in self.__removeEntities__:
                self.__updateEntities__[to.virtualId] = to

    def deleteEntity(self, entity: tgentimpl.AbstractEntity):
        if entity.isDeleted:
            raise tgexception.TGException('Should not be calling delete on an already deleted entity!')
        # Remove any entities added to the add changelist
        if entity.virtualId in self.__addEntities__:
            del self.__addEntities__[entity.virtualId]
        # Remove any entities added to the update changelist
        if entity.virtualId in self.__updateEntities__:
            del self.__updateEntities__[entity.virtualId]
        if entity.isNew:
            entity.markDeleted()
        else:
            self.__removeEntities__[entity.virtualId] = entity
            self.__updateEdge__(entity)

    def createQuery(self, query: str) -> tgquery.TGQuery:
        channelResponseWaiter: tgchannel.TGChannelResponseWaiter
        result: int
        ret: tgquery.TGQuery = None
        channelResponseWaiter = self._genBCRWaiter()
        try:
            request: tgpduimpl.QueryRequestMessage = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.QueryRequest,
                                                                                      authtoken=self.__channel__.authtoken,
                                                                                      sessionid=self.__channel__.sessionid)
            request.command = tgquery.TGQueryCommand.Create
            request.query = query
            response: tgpduimpl.QueryResponseMessage = self.__channel__.send(request, channelResponseWaiter)
            gLogger.log(TGLevel.Debug, "Send query completed")
            result: int = response.result
            queryHashId: int = response.queryHashId
            if result == 0 and queryHashId > 0:         #TODO Create error reporting for query result.
                ret = tgqueryimpl.QueryImpl(self, queryHashId)

        finally:
            pass
        return ret

    def executeQuery(self, query: typing.Optional[str] = None,
                     option: tgquery.TGQueryOption = tgquery.DefaultQueryOption) -> tgquery.TGResultSet:
        if query is None:
            try:
                query = option.queryExpr
            except KeyError as e:
                raise tgexception.TGException("Need to specify a query string!", cause=e)
        channelResponseWaiter: tgchannel.TGChannelResponseWaiter = self._genBCRWaiter()
        result: int
        try:
            request: tgpduimpl.QueryRequestMessage = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.QueryRequest,
                                                                                              authtoken=self.__channel__.authtoken,
                                                                                              sessionid=self.__channel__.sessionid)
            request.option = option
            request.command, request.query = findCommandAndQueryString(query, self.__props__)

            response: tgpduimpl.QueryResponseMessage = self.__channel__.send(request, channelResponseWaiter)

            if response.error is not None:
                raise response.error

            return response.finishReadWith(request.command, self.__gof__)
        except (Exception, tgexception.TGException):
            raise

    # TODO implement some form of compiled queries
    def executeQueryWithId(self, queryId: int, option: tgquery.TGQueryOption = tgquery.DefaultQueryOption) -> \
            tgquery.TGResultSet:
        result: int
        channelResponseWaiter: tgchannel.TGChannelResponseWaiter = self._genBCRWaiter()
        try:
            request: tgpduimpl.QueryRequestMessage = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.QueryRequest,
                                                                                              authtoken=self.__channel__.authtoken,
                                                                                              sessionid=self.__channel__.sessionid)
            request.command = tgquery.TGQueryCommand.ExecuteID
            request.queryHashId = queryId
            request.option = option
            response: tgpduimpl.QueryResponseMessage = self.__channel__.send(request, channelResponseWaiter)

            return response.finishReadWith(tgquery.TGQueryCommand.ExecuteID, self.__gof__)
        except Exception as e:
            raise tgexception.TGException("Exception in executeQueryWithId", cause=e)

    def closeQuery(self, queryId: int):
        channelResponseWaiter: tgchannel.TGChannelResponseWaiter = self._genBCRWaiter()
        try:
            request: tgpduimpl.QueryRequestMessage = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.QueryRequest,
                                                                                              authtoken=self.__channel__.authtoken,
                                                                                              sessionid=self.__channel__.sessionid)
            request.command = tgquery.TGQueryCommand.Close
            request.queryHashId = queryId
            _: tgpduimpl.QueryResponseMessage = self.__channel__.send(request, channelResponseWaiter)
            # TODO check response state
            gLogger.log(TGLevel.Debug, "Send close query completed")
        except Exception as e:
            raise tgexception.TGException("Exception in closeQuery", cause=e)

    def getLargeObjectAsBytes(self, entityId: int, encrypted: bool = False) -> bytes:
        channelResponseWaiter = self._genBCRWaiter()
        if encrypted:                                   # TODO Decrypt encrypted entities
            raise tgexception.TGProtocolNotSupported("Blob/Clob encryption/decryption not implemented.")
        request: tgpduimpl.GetLargeObjectRequestMessage = tgpduimpl.TGMessageFactory.createMessage(
                            tgpdu.VerbId.GetLargeObjectRequest, authtoken=self.__channel__.authtoken, sessionid=self.__channel__.sessionid)

        request.entityId = entityId
        request.decrypt = encrypted
        response: tgpduimpl.GetLargeObjectResponseMessage = self.__channel__.send(request, channelResponseWaiter)
        if entityId != response.entityId:
            raise tgexception.TGException("Server responded with different entityId than expected!")
        data = bytes() if response.data is None else response.data
        return data

    @property
    def linkState(self) -> tgchannel.LinkState:
        return self.__channel__.linkstate

    @property
    def outboxaddr(self) -> str:
        return self.__channel__.outboxaddr

    @property
    def connectedUsername(self) -> str:
        return self.__username__

    @property
    def graphMetadata(self) -> tgmodel.TGGraphMetadata:
        return self.__gof__.graphmetadata

    @property
    def graphObjectFactory(self) -> tgmodel.TGGraphObjectFactory:
        return self.__gof__
