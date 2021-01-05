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
 *  File name :bulkioimpl.py
 *  Created on: 9/3/2019
 *  Created by: derek
 *
 *		SVN Id:
 *
 *  This file encapsulates all the unit test case for bulk import test
 """


import typing
import io
import re

import tgdb.log as tglog
import tgdb.utils as tgutils
import tgdb.model as tgmodel
import tgdb.connection as tgconn
import tgdb.impl.connectionimpl as tgconnimpl
import tgdb.exception as tgexception
import tgdb.bulkio as tgbulk
import tgdb.impl.adminimpl as tgadmimpl
import csv


class BulkExportDescImpl(tgbulk.TGBulkExportDescriptor):

    def __init__(self, type: tgmodel.TGEntityType, data: typing.Optional[typing.Any] = None, numEnts: int = 0):
        self.__type__ = type

        self.__data__ = data
        self.__numEnts__ = numEnts
        self.__handlers__: typing.List[typing.Callable[[tgbulk.TGBulkExportDescriptor], None]] = []
        self.__raw_data: str = ""
        self.leftover: str = ""

    def __len__(self) -> int:
        return self.__numEnts__

    @property
    def raw_data(self) -> str:
        return self.__raw_data

    @raw_data.setter
    def raw_data(self, new_data):
        self.__raw_data = new_data

    @property
    def entitytype(self) -> tgmodel.TGEntityType:
        return self.__type__

    @property
    def data(self) -> typing.Optional[typing.Any]:
        return self.__data__

    @data.setter
    def data(self, val: typing.Optional[typing.Any]):
        oldData = self.__data__
        self.__data__ = val
        if self.hasHandler and val is not None:
            for handler in self.__handlers__:
                handler(self)
        elif oldData is not None and self.__data__ is not None:
            self.__data__ = oldData.append(self.__data__)

    def addHandler(self, handler: typing.Callable[[tgbulk.TGBulkExportDescriptor], None]):
        self.__handlers__.append(handler)

    def clearHandlers(self):
        self.__handlers__ = []

    def _runHandlers(self):
        for handler in self.__handlers__:
            handler(self)

    @property
    def hasHandler(self) -> bool:
        return len(self.__handlers__) > 0


def findNumNewlineOutsideQuotes(data: str) -> int:
    index = 0
    ret = 0
    in_quotes = False
    next_match = re.search('["\n]', data[index:])
    while next_match:
        index += next_match.start()
        if data[index] == '"':
            in_quotes = not in_quotes
        elif not in_quotes:
            ret += 1
        index += 1
        next_match = re.search('["\n]', data[index:])
    return ret


class BulkImportImpl(tgbulk.TGBulkImport):
    idColName = "@TGDB_SPECIAL@_@id@"
    fromColName = "@TGDB_SPECIAL@_@from@"
    toColName = "@TGDB_SPECIAL@_@to@"
    edgeImportStart: int = 1027
    MAX_STRING_LEN: int = 32000

    def __init__(self, conn: tgconn.TGConnection, props: tgutils.TGProperties = None):
        self.__props__: tgutils.TGProperties
        if props is None:
            self.__props__ = tgutils.TGProperties()
        else:
            self.__props__ = props
        self.__conn__: tgconnimpl.ConnectionImpl = conn
        self.__descs: typing.Dict[str, tgadmimpl.ImportDescImpl] = {}
        self.__edgeId__: int = 0

    @property
    def _c(self) -> tgconn.TGConnection:
        return self.__conn__

    def _set_nodeBatchSize(self, value: int):
        if not isinstance(value, int):
            raise TypeError("Bad type! Check the API!")
        self.__props__[tgutils.ConfigName.BulkIONodeBatchSize] = str(value)

    def _get_nodeBatchSize(self) -> int:
        return int(self.__props__.get(tgutils.ConfigName.BulkIONodeBatchSize.propertyname,
                                      tgutils.ConfigName.BulkIONodeBatchSize.defaultvalue))

    nodeBatchSize = property(_get_nodeBatchSize, _set_nodeBatchSize)

    def _set_edgeBatchSize(self, value: int):
        if not isinstance(value, int):
            raise TypeError("Bad type! Check the API!")
        self.__props__[tgutils.ConfigName.BulkIOEdgeBatchSize] = str(value)

    def _get_edgeBatchSize(self) -> int:
        return int(self.__props__.get(tgutils.ConfigName.BulkIOEdgeBatchSize.propertyname,
                                      tgutils.ConfigName.BulkIOEdgeBatchSize.defaultvalue))

    edgeBatchSize = property(_get_edgeBatchSize, _set_edgeBatchSize)

    def _set_idColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[tgutils.ConfigName.BulkIOIdColName] = value

    def _get_idColumnName(self) -> str:
        return self.__props__.get(tgutils.ConfigName.BulkIOIdColName.propertyname,
                                  tgutils.ConfigName.BulkIOIdColName.defaultvalue)

    idColumnName = property(_get_idColumnName, _set_idColumnName)

    def _set_fromColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[tgutils.ConfigName.BulkIOFromColName] = value

    def _get_fromColumnName(self) -> str:
        return self.__props__.get(tgutils.ConfigName.BulkIOFromColName.propertyname,
                                  tgutils.ConfigName.BulkIOFromColName.defaultvalue)

    fromColumnName = property(_get_fromColumnName, _set_fromColumnName)

    def _set_toColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[tgutils.ConfigName.BulkIOToColName] = value

    def _get_toColumnName(self) -> str:
        return self.__props__.get(tgutils.ConfigName.BulkIOToColName.propertyname,
                                  tgutils.ConfigName.BulkIOToColName.defaultvalue)

    toColumnName = property(_get_toColumnName, _set_toColumnName)

    def __isApprovedNodeColumn(self, col: str) -> bool:
        return col in self._c.graphMetadata.attritubeDescriptors or col == self.idColumnName

    def __pdToCSV(self, data: typing.Any) -> str:
        return data.to_csv(header=None, index=False, quoting=csv.QUOTE_MINIMAL,
                                                                doublequote=False, escapechar='\\', na_rep="null")

    def __findMaxNumberPerBatch(self, data: typing.Any, maxBatchSize: int) -> int:
        data_buf: str = self.__pdToCSV(data[:maxBatchSize])
        encoded_data = data_buf.encode('utf-8')
        str_data = encoded_data.decode('ascii', errors='replace')
        str_data = str_data[:BulkImportImpl.MAX_STRING_LEN]
        num_matches = findNumNewlineOutsideQuotes(str_data)
        return num_matches if num_matches < maxBatchSize else maxBatchSize

    def __findNumBatches(self, data: typing.Any, maxBatchSize: int) -> int:
        startIdx = 0
        numIter = 0
        while startIdx < len(data):
            numIter += 1
            startIdx += self.__findMaxNumberPerBatch(data[startIdx:], maxBatchSize)
        return numIter

    def __stripColumns(self, columns: typing.List[str]) -> typing.List[str]:
        types: typing.List[str] = []
        for colName in columns:
            if colName not in (BulkImportImpl.idColName, BulkImportImpl.fromColName, BulkImportImpl.toColName):
                types.append(colName)
        return types

    def __pushBatch(self, data: typing.Any, type: tgmodel.TGEntityType, batchIdx: int, numBatches: int,
                    maxBatchSize: int) -> int:
        sendNum = self.__findMaxNumberPerBatch(data, maxBatchSize)
        resultList = self.__conn__.partialImportEntity(type, batchIdx, numBatches,
                                          self.__pdToCSV(data[:sendNum]), self.__stripColumns(data.columns))
        for result in resultList:
            if result.typename in self.__descs:
                self.__descs[result.typename]._inc(result.numinstances)
            else:
                self.__descs[result.typename] = result
            if result.typename == type.name and result.numinstances != sendNum:
                tglog.gLogger.log(tglog.TGLevel.Warning, "Should have sent %d number of entities, server loaded %d.",
                                  sendNum, result.numinstances)
            elif result.typename != type.name:
                tglog.gLogger.log(tglog.TGLevel.Warning, "Server received type %s, when all client sent was %s.",
                                  result.typename, type.name)
        return sendNum

    def __cleanPdDataFrame(self, data):
        # Do some funky stuff with \ so that the server can easily read the csv. ;)
        data = data.replace({"": "\"\""})
        data = data.replace("\\\\", "\\\\\\\\", regex=True)       # 8 \ in a row - peek regex performance?
        return data.replace("'", "\\'", regex=True)

    def importNodes(self, data: typing.Any, nodetype: typing.Union[str, tgmodel.TGNodeType]) \
            -> "tgadmimpl.ImportDescImpl":

        import pandas as pd
        data: pd.DataFrame = data

        hasIds: bool = self.idColumnName in data.columns

        if isinstance(nodetype, str):
            nodetype = self._c.graphMetadata.nodetypes[nodetype]

        if not hasIds:
            raise tgexception.TGException("Must have a unique ID when importing nodes.")

        if hasIds:
            self.__findFirstBadIfAny(data.isna()[self.idColumnName],
                                     "Should not have any NA values in the 'id' column. Index is {0}")

        numNodes = data.shape[0]

        attrDesc: tgmodel.TGAttributeDescriptor

        # Ensure that all primary key attributes are in this DataFrame

        for attrDesc in nodetype.pkeyAttributeDescriptors:
            # If the attribute descriptor is in the primary key attributes, but is not a column, unless it is a special
            # attribute descriptor (denoted by the first symbol of '@')
            if attrDesc.name not in data.columns and attrDesc.name[:1] != '@':
                raise tgexception.TGException("Need to specify all primary key attributes in the columns field to legal"
                                              "ly import this data set. Could not find the primary key attribute {0}."
                                              .format(attrDesc.name))

        toDrop = []

        # Remove all unneeded column names i.e. ones that are not part of the nodetype
        [toDrop.append(column) for column in data.columns if not self.__isApprovedNodeColumn(column)]

        data = data.drop(columns=toDrop)

        # Add all columns specified by the nodetype that are not in this dataframe, while maintaining the approved
        # columns. Set because we don't want to add the same row twice.

        reIndexCols = set()

        data = self.__cleanPdDataFrame(data)

        [reIndexCols.add(col) for col in data.columns if col != self.idColumnName]

        toReindex: typing.List[str] = [BulkImportImpl.idColName]

        toReindex.extend(reIndexCols)

        data = data.rename(columns={self.idColumnName: BulkImportImpl.idColName}).reindex(columns=toReindex)

        numBatches = self.__findNumBatches(data, self.nodeBatchSize)

        curBatchStart: int = 0
        sent: int

        batchIdx: int = 0

        while curBatchStart < numNodes:
            sent = self.__pushBatch(data[curBatchStart:], nodetype, batchIdx, numBatches, self.nodeBatchSize)

            batchIdx += 1
            curBatchStart += sent

        result = self.__descs[nodetype.name] if nodetype.name in self.__descs else None
        return result

    def __findFirstBadIfAny(self, isna: typing.Any, message: str):
        import pandas as pd
        isna: pd.DataFrame = isna
        if isna.values.any():
            for i in range(isna.shape[0]):
                if isna.iat[0, i]:
                    raise tgexception.TGException(message.format(i))

    def __isApprovedEdgeColumn(self, col: str) -> bool:
        return col in self._c.graphMetadata.attritubeDescriptors or col == self.fromColumnName or \
               col == self.toColumnName

    def importEdges(self, data: typing.Any, edgetype: typing.Union[str, tgmodel.TGEdgeType] = None,
                    edgedir: tgmodel.DirectionType = tgmodel.DirectionType.Directed) -> "tgadmimpl.ImportDescImpl":
        import pandas as pd
        data: pd.DataFrame = data

        if data.ndim != 2:
            raise tgexception.TGException("bulkImportEdges only supports Pandas DataFrames with dimension 2, not {0}".
                                          format(data.ndim))

        if self.fromColumnName not in data.columns:
            raise tgexception.TGException("Must have one column labeled '{0}' that represents the 'from' nodes"
                                          .format(self.fromColumnName))

        if self.toColumnName not in data.columns:
            raise tgexception.TGException("Must have one column labeled '{0}' that represents the 'to' nodes"
                                          .format(self.toColumnName))

        isna = data.isna()

        self.__findFirstBadIfAny(isna[self.fromColumnName], "Should not have any NA values in the 'from' column. "
                                                            "Index is {0}")

        self.__findFirstBadIfAny(isna[self.toColumnName], "Should not have any NA values in the 'to' column. Ind"
                                                          "ex is {0}")

        if isinstance(edgetype, str) and edgetype is not None:
            edgetype = self._c.graphMetadata.edgetypes[edgetype]

        if edgetype is None:
            edgetype = self._c.graphMetadata.defaultEdgetype(edgedir)

        data = self.__cleanPdDataFrame(data)

        numEdges = data.shape[0]

        toDrop = []
        # Drop all columns not corresponding to an approved function: either linkage information, or attribute
        [toDrop.append(column) for column in data.columns if not self.__isApprovedEdgeColumn(column)]

        for column in toDrop:
            tglog.gLogger.log(tglog.TGLevel.Warning, "Stripping column %s because no corresponding attribute descriptor"
                                                     " exists.", column)

        data = data.drop(columns=toDrop)

        reIndexCols = set()

        [reIndexCols.add(col) for col in data.columns if col != self.fromColumnName and
         col != self.toColumnName]

        toReindex: typing.List[str] = [BulkImportImpl.fromColName, BulkImportImpl.toColName]

        toReindex.extend(reIndexCols)

        # Rename the columns to show to reduce any work the server needs to do when parsing the column labels.
        data = data.rename(columns={self.fromColumnName: BulkImportImpl.fromColName,
                                    self.toColumnName: BulkImportImpl.toColName})
        data = data.reindex(columns=toReindex)

        numBatches = self.__findNumBatches(data, self.edgeBatchSize)

        curBatchStart: int = 0
        sent: int

        batchIdx: int = 0

        while curBatchStart < numEdges:
            sent = self.__pushBatch(data[curBatchStart:], edgetype, batchIdx, numBatches, self.edgeBatchSize)

            curBatchStart += sent
            batchIdx += 1

        result = self.__descs[edgetype.name] if edgetype.name in self.__descs else None
        return result

    def _close(self, nodetype: tgmodel.TGNode):
        # TODO Send close signal for this nodetype to release the instances from the server, if we want to implement
        #  that
        pass

    def finish(self):
        resultList = self.__conn__.endBulkImport()
        return resultList


class BulkExportImpl(tgbulk.TGBulkExport):

    def _set_idColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[tgutils.ConfigName.BulkIOIdColName] = value

    def _get_idColumnName(self) -> str:
        return self.__props__.get(tgutils.ConfigName.BulkIOIdColName.propertyname,
                                  tgutils.ConfigName.BulkIOIdColName.defaultvalue)

    idColumnName = property(_get_idColumnName, _set_idColumnName)

    def _set_fromColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[tgutils.ConfigName.BulkIOFromColName] = value

    def _get_fromColumnName(self) -> str:
        return self.__props__.get(tgutils.ConfigName.BulkIOFromColName.propertyname,
                                  tgutils.ConfigName.BulkIOFromColName.defaultvalue)

    fromColumnName = property(_get_fromColumnName, _set_fromColumnName)

    def _set_toColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[tgutils.ConfigName.BulkIOToColName] = value

    def _get_toColumnName(self) -> str:
        return self.__props__.get(tgutils.ConfigName.BulkIOToColName.propertyname,
                                  tgutils.ConfigName.BulkIOToColName.defaultvalue)

    toColumnName = property(_get_toColumnName, _set_toColumnName)

    def __processTypeList(self, typeList: typing.List[typing.Tuple[str, bool, int]]):
        for tup in typeList:
            typename = tup[0]
            type: tgmodel.TGEntityType = None
            gmd = self._c.graphMetadata
            if typename in gmd.nodetypes:
                type = gmd.nodetypes[typename]
            elif typename in gmd.edgetypes:
                type = gmd.edgetypes[typename]
            elif typename == 'Default Bidirected Edgetype':
                type = gmd.defaultEdgetype(tgmodel.DirectionType.BiDirectional)
            elif typename == 'Default Undirected Edgetype':
                type = gmd.defaultEdgetype(tgmodel.DirectionType.UnDirected)
            elif typename == 'Default Directed Edgetype':
                type = gmd.defaultEdgetype(tgmodel.DirectionType.Directed)
            else:
                tglog.gLogger.log(tglog.TGLevel.Error, "Unknown type returned from server: %s", typename)

            self.__descMap__[typename] = BulkExportDescImpl(type, numEnts=tup[2])

    def __getitem__(self, typename: str) -> tgbulk.TGBulkExportDescriptor:
        return self.__descMap__[typename]

    def __init__(self, conn: tgconn.TGConnection, props: typing.Optional[tgutils.TGProperties],
                 typeList: typing.List[typing.Tuple[str, bool, int]], numRequests: int):
        if props is None:
            props = tgutils.TGProperties()

        self.__props__ = props
        self.__conn__ = conn
        self.__numRequests__ = numRequests
        self.__descMap__: typing.Dict[str, BulkExportDescImpl] = {}
        self.__typeList = typeList
        self.__processTypeList(typeList)
        self.__conf = ""

    @property
    def _c(self) -> tgconn.TGConnection:
        return self.__conn__

    def __renameCols(self, colLabels: typing.List[str]):
        for i in range(len(colLabels)):
            if colLabels[i] == "id":
                colLabels[i] = self.idColumnName
            elif colLabels[i] == "from":
                colLabels[i] = self.fromColumnName
            elif colLabels[i] == "to":
                colLabels[i] = self.toColumnName

    def exportEntities(self, handler: typing.Optional[typing.Callable[[tgbulk.TGBulkExportDescriptor], None]] = None,
                       callPassedHandlerOnAll: bool = False):
        import pandas as pd
        typename = None
        hasNext = len(self.__typeList) > 0
        attributeList: typing.List[str] = None
        while hasNext:
            #           fileName  data  hasNext  numEntities  newTypeInformation
            resp: typing.Tuple[str, bytes, bool, int, typing.Optional[typing.Tuple[str, int, typing.List[str]]]] =\
                self.__conn__.partialExport(-1)

            if resp[4] is not None:
                new_type_info = resp[4]
                typename = new_type_info[0]
                desc = self.__descMap__[typename]
                attributeList = ['id']
                if desc.entitytype.systemType == tgmodel.TGSystemType.EdgeType:
                    attributeList = ['from', 'to']
                attributeList.extend(new_type_info[1])
                self.__renameCols(attributeList)

            desc = self.__descMap__[typename]

            data = resp[1].decode('utf-8', errors='replace')

            desc.raw_data = data

            temp_data = desc.leftover

            temp_data += data

            lastNew = len(temp_data)

            index = temp_data.rfind('\n')

            if index != -1:
                lastNew = index

            desc.leftover = temp_data[lastNew:]

            file = io.StringIO(initial_value=temp_data[:lastNew])

            desc.data = pd.read_csv(file, header=None, doublequote=False, escapechar='\\', na_values=["null"],
                                    keep_default_na=False, comment='#', names=attributeList)

            file.close()

            if (not desc.hasHandler or callPassedHandlerOnAll) and handler is not None:
                handler(desc)

            hasNext = resp[2]


"""
class BulkExportImpl(tgbulk.TGBulkExport):

    def _set_nodeBatchSize(self, value: int):
        if not isinstance(value, int):
            raise TypeError("Bad type! Check the API!")
        self.__props__[ConfigName.BulkIONodeBatchSize] = str(value)

    def _get_nodeBatchSize(self) -> int:
        return int(self.__props__.get(ConfigName.BulkIONodeBatchSize.propertyname,
                                      ConfigName.BulkIONodeBatchSize.defaultvalue))

    nodeBatchSize = property(_get_nodeBatchSize, _set_nodeBatchSize)

    def _set_edgeBatchSize(self, value: int):
        if not isinstance(value, int):
            raise TypeError("Bad type! Check the API!")
        self.__props__[ConfigName.BulkIOEdgeBatchSize] = str(value)

    def _get_edgeBatchSize(self) -> int:
        return int(self.__props__.get(ConfigName.BulkIOEdgeBatchSize.propertyname,
                                      ConfigName.BulkIOEdgeBatchSize.defaultvalue))

    edgeBatchSize = property(_get_edgeBatchSize, _set_edgeBatchSize)

    def _set_idColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[ConfigName.BulkIOIdColName] = value

    def _get_idColumnName(self) -> str:
        return self.__props__.get(ConfigName.BulkIOIdColName.propertyname,
                                  ConfigName.BulkIOIdColName.defaultvalue)

    idColumnName = property(_get_idColumnName, _set_idColumnName)

    def _set_fromColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[ConfigName.BulkIOFromColName] = value

    def _get_fromColumnName(self) -> str:
        return self.__props__.get(ConfigName.BulkIOFromColName.propertyname,
                                  ConfigName.BulkIOFromColName.defaultvalue)

    fromColumnName = property(_get_fromColumnName, _set_fromColumnName)

    def _set_toColumnName(self, value: str):
        if not isinstance(value, str) and value is not None:
            raise TypeError("Bad type! Check the API!")
        self.__props__[ConfigName.BulkIOToColName] = value

    def _get_toColumnName(self) -> str:
        return self.__props__.get(ConfigName.BulkIOToColName.propertyname,
                                  ConfigName.BulkIOToColName.defaultvalue)

    toColumnName = property(_get_toColumnName, _set_toColumnName)

    def __init__(self, conn: tgconn.TGConnection, props: TGProperties = None):
        if props is None:
            props = TGProperties()

        self.__props__ = props
        self.__conn__: tgconnimpl.ConnectionImpl = conn

    @property
    def _c(self) -> tgconn.TGConnection:
        return self.__conn__

    def __renameCols(self, colLabels: List[str]):
        for i in range(len(colLabels)):
            if colLabels[i] == "id":
                colLabels[i] = self.idColumnName
            elif colLabels[i] == "from":
                colLabels[i] = self.fromColumnName
            elif colLabels[i] == "to":
                colLabels[i] = self.toColumnName

    def __exportEntBatch(self, exportDesc: int, colLabels: List[str], dataPtr: List[pd.DataFrame], maxEntities: int,
                         processEntities: Optional[Callable[[pd.DataFrame], Any]], batchSize: int) -> Tuple[bool, int]:
        csvStr: str
        haveMore: bool
        numEnts, csvStr, haveMore = self.__conn__.singleBatchExportEntity(exportDesc)

        tmpData = pd.read_csv(csvStr, header=None, columns=colLabels)

        dataPtr[0] = dataPtr[0].append(tmpData)

        if haveMore and numEnts != batchSize:
            gLogger.log(TGLevel.Warning, "Database has more entities to export, but the number of nodes exported is not"
                                         " the requested node batch size.")

        if numEnts != tmpData.shape[0]:
            gLogger.log(TGLevel.Warning, "Number of entities read is different than what the database stated was the nu"
                                         "mber of nodes sent. Possibly a reading/corruption error?")

        if processEntities is not None:
            data = dataPtr[0]
            while data.shape[0] > maxEntities:
                processEntities(data[:maxEntities])
                data = data[maxEntities:]
            dataPtr[0] = data

        return haveMore, tmpData.shape[0]

    def exportNodes(self, nodetype: Union[str, tgmodel.TGNodeType],
                    processNodes: Optional[Callable[[pd.DataFrame], Any]] = None,
                    maxNodes: int = -1) -> tgbulk.TGBulkExportDescriptor:
        if nodetype is None:
            raise tgexception.TGException("Cannot export a None nodetype!")

        if isinstance(nodetype, str):
            nodetype = self._c.graphMetadata.nodetypes[nodetype]

        if maxNodes != -1 and maxNodes < 1:
            raise tgexception.TGException("Must specify either -1 to use whatever the batch size is, or something great"
                                          "er than 0!")

        useProcess: bool = True
        if processNodes is None:
            useProcess = False
            maxNodes = -1
        elif maxNodes == -1:
            maxNodes = self.nodeBatchSize

        desc: int
        colLabels: List[str]
        desc, colLabels = self.__conn__.beginBatchExportEntity(tgmodel.TGEntityKind.Node, nodetype, self.nodeBatchSize)

        self.__renameCols(colLabels)

        data = pd.DataFrame()
        dataPtr: List[pd.DataFrame] = [data]
        total: int = 0

        haveMore = True

        while haveMore:
            haveMore, num = self.__exportEntBatch(desc, colLabels, dataPtr, maxNodes, processNodes, self.nodeBatchSize)
            total += num

        data = dataPtr[0]
        if data.shape[0] > 0 and useProcess:
            processNodes(data)

        return BulkExportDescImpl(nodetype, None if useProcess else data, total)

    def exportNodesAsGenerator(self, nodetype: Union[str, tgmodel.TGNodeType]) \
            -> Generator[pd.DataFrame, None, BulkExportDescImpl]:
        if nodetype is None:
            raise tgexception.TGException("Cannot export a None nodetype!")

        if isinstance(nodetype, str):
            nodetype = self._c.graphMetadata.nodetypes[nodetype]

        desc: int
        colLabels: List[str]
        desc, colLabels = self.__conn__.beginBatchExportEntity(tgmodel.TGEntityKind.Node, nodetype, self.nodeBatchSize)
        total: int = 0

        haveMore = True
        while haveMore:
            data = pd.DataFrame()
            dataPtr: List[pd.DataFrame] = [data]

            haveMore, num = self.__exportEntBatch(desc, colLabels, dataPtr, -1, None, self.nodeBatchSize)

            total += num

            yield dataPtr[0]

        return BulkExportDescImpl(nodetype, None, total)

    def exportEdges(self, edgetype: Union[str, tgmodel.TGEdgeType] = None,
                    edgedir: tgmodel.DirectionType = tgmodel.DirectionType.Directed,
                    processEdges: Optional[Callable[[pd.DataFrame], Any]] = None,
                    maxEdges: int = -1) -> tgbulk.TGBulkExportDescriptor:
        if isinstance(edgetype, str) and edgetype is not None:
            edgetype = self._c.graphMetadata.edgetypes[edgetype]

        if edgetype is None:
            edgetype = self._c.graphMetadata.defaultEdgetype(edgedir)

        useProcess: bool = True
        if processEdges is None:
            useProcess = False
            maxEdges = -1
        elif maxEdges == -1:
            maxEdges = self.nodeBatchSize

        desc: int
        colLabels: List[str]
        desc, colLabels = self.__conn__.beginBatchExportEntity(tgmodel.TGEntityKind.Node, edgetype, self.edgeBatchSize)
        total: int = 0

        data = pd.DataFrame()
        dataPtr: List[pd.DataFrame] = [data]

        haveMore = True

        while haveMore:
            haveMore, num = self.__exportEntBatch(desc, colLabels, dataPtr, maxEdges, processEdges, self.edgeBatchSize)
            total += num

        data = dataPtr[0]
        if data.shape[0] > 0 and processEdges is not None:
            processEdges(data)

        return BulkExportDescImpl(edgetype, None if useProcess else data, total)

    def exportEdgeAsGenerator(self, edgetype: Union[str, tgmodel.TGEdgeType] = None,
                              edgedir: tgmodel.DirectionType = tgmodel.DirectionType.Directed) \
            -> Generator[pd.DataFrame, None, BulkExportDescImpl]:
        if isinstance(edgetype, str):
            edgetype = self._c.graphMetadata.edgetypes[edgetype]

        if edgetype is None:
            edgetype = self._c.graphMetadata.defaultEdgetype(edgedir)

        desc: int
        colLabels: List[str]
        desc, colLabels = self.__conn__.beginBatchExportEntity(tgmodel.TGEntityKind.Edge, edgetype, self.edgeBatchSize)
        total = 0

        haveMore = True
        while haveMore:
            data = pd.DataFrame()
            dataPtr: List[pd.DataFrame] = [data]

            haveMore, num = self.__exportEntBatch(desc, colLabels, dataPtr, -1, None, self.edgeBatchSize)

            total += num

            yield dataPtr[0]

        return BulkExportDescImpl(edgetype, None, total)

    def exportAll(self) -> Dict[str, tgbulk.TGBulkExportDescriptor]:
        ret = {}  # TODO find a better version that takes more advantage of callbacks??

        for nodetype in self._c.graphMetadata.nodetypes:
            ret[nodetype] = self.exportNodes(nodetype)

        for edgetype in self._c.graphMetadata.edgetypes:
            ret[edgetype] = self.exportEdges(edgetype)

        return ret

    def exportAllAsGenerator(self) -> \
            Generator[Tuple[str, Generator[pd.DataFrame, None, BulkExportDescImpl]], None, None]:
        for nodetype in self._c.graphMetadata.nodetypes:
            yield nodetype, self.exportNodes(nodetype)

        for edgetype in self._c.graphMetadata.edgetypes:
            yield edgetype, self.exportEdges(edgetype)

    def finish(self):
        self.__conn__.endBulkExportSession()
"""