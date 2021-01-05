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
 *  File name :queryimpl.py
 *  Created on: 8/22/2019
 *  Created by: derek
 *
 *		SVN Id: $Id$
 *
 *  This file encapsulates all the Graph metadata and its related operations
 """

import tgdb.connection as tgconn
import tgdb.query as tgquery
import tgdb.exception as tgexception
import tgdb.pdu as tgpdu
import tgdb.model as tgmodel
import enum
import re
import typing
from tgdb.impl.entityimpl import AbstractEntity
from tgdb.impl.attrimpl import AbstractAttribute
from tgdb.log import *
from typing import *

from tgdb.query import TGQueryResultDataType, TGResultDataDescriptor

T = TypeVar('T')


class QueryImpl(tgquery.TGQuery):

    def __init__(self, conn: tgconn.TGConnection, queryHashId: int):
        self.__conn__: tgconn.TGConnection = conn
        self.__queryHashId__: int = queryHashId
        self.__option__: tgquery.TGQueryOption = tgquery.TGQueryOption()

    def __setitem__(self, key, value):
        pass

    def __getitem__(self, item):
        return None

    def execute(self) -> tgquery.TGResultSet:
        return self.__conn__.executeQueryWithId(self.__queryHashId__, self.option)

    def close(self):
        self.__conn__.closeQuery(self.__queryHashId__)

    @property
    def option(self) -> tgquery.TGQueryOption:
        return self.__option__

    @tgquery.TGQuery.option.setter
    def option(self, val: tgquery.TGQueryOption):
        self.__option__ = val


def _matchArg(grem_str: str) -> typing.Tuple[str, str]:
    ret: str
    result: str
    grem_str = grem_str.strip()
    quote_match = re.match('(\')(?P<ret>([^\'\\\\]|(\\\\.))*)(\')(?P<result>.+)', grem_str)
    num_match = re.match('(?P<ret>(\\+-)?[0-9]+(\\.[0-9]+([eE][+\\-]?[0-9]))?)(?P<result>.+)', grem_str)
    pred_match = re.match('(\\w+\\()(.+)', grem_str)
    list_match = re.match('(\\[)(.+)', grem_str)
    keyword_match = re.match('(?P<ret>\\w+)(?P<result>([^(]|$).*)', grem_str)
    if quote_match:
        ret = quote_match.group('ret')
        result = quote_match.group('result')
    elif num_match:
        ret = num_match.group('ret')
        result = num_match.group('result')
    elif pred_match:
        ret, result = _matchMethod(grem_str, sub_pred=True)
    elif list_match:
        ret = '['
        lis, result = _matchArgList(grem_str[1:].strip())
        ret += lis + ']'
        result = result.strip()
        if len(result) <= 0 or result[0] != ']':
            raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')
        result = result[1:]
    elif keyword_match:
        ret = keyword_match.group('ret')
        result = keyword_match.group('result')
        if ret not in {'local'}:        # TODO Find out more keywords.
            raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')
    else:
        raise tgexception.TGException('Unable to conform the gremlin query to bytecode:')
    result = result.strip()
    return ret, result


def _matchArgList(grem_str: str) -> typing.Tuple[str, str]:
    another_arg = (grem_str[0] != ')' and grem_str[0] != ']')
    if len(grem_str) < 0:
        raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')

    found_prior = False

    ret = ''

    while another_arg:

        arg_res, grem_str = _matchArg(grem_str)

        if found_prior:
            ret += ', '
        found_prior = True
        ret += arg_res

        grem_str = grem_str.strip()

        if len(grem_str) < 1 or (grem_str[0] != ')' and grem_str[0] != ',' and grem_str[0] != ']'):
            raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')

        another_arg = grem_str[0] == ','
        if another_arg:
            grem_str = grem_str[1:]
    return ret, grem_str


def _matchMethod(grem_str: str, sub_pred: bool = False) -> typing.Tuple[str, str]:
    match = re.match('(^\\w+\\()(.+)', grem_str)
    if not match:
        raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')

    ret = match.group(1)
    if ret == 'between(' and sub_pred:
        ret = 'and(gte('
        first_arg, grem_str = _matchArg(match.group(2).strip())
        grem_str = grem_str.strip()
        if grem_str[0] != ',':
            raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')
        ret += first_arg + '), lt('
        second_arg, grem_str = _matchArg(grem_str[1:].strip())
        ret += second_arg + ')'
    elif ret == 'outside(' and sub_pred:
        ret = 'or(lt('
        first_arg, grem_str = _matchArg(match.group(2).strip())
        grem_str = grem_str.strip()
        if grem_str[0] != ',':
            raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')
        ret += first_arg + '), gt('
        second_arg, grem_str = _matchArg(grem_str[1:].strip())
        ret += second_arg + ')'
    else:
        grem_str = match.group(2).strip()

        arg_list, grem_str = _matchArgList(grem_str)

        ret += arg_list

    if grem_str[0] != ')':
        raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')

    ret += ')'

    grem_str = grem_str[1:].strip()

    return ret, grem_str


def convertToBytecode(gremlin: str) -> str:
    if len(gremlin) < 2 or gremlin[:2] != 'g.':
        raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')

    ret = '[[], ['

    another_step = True

    gremlin = gremlin[2:]

    has_other = False

    while another_step:
        gremlin = gremlin.strip()

        full_step, gremlin = _matchMethod(gremlin)

        if has_other:
            ret += ', '
        has_other = True

        ret += full_step

        gremlin = gremlin.strip()

        if len(gremlin) < 1 or (gremlin[0] != '.' and gremlin[0] != ';'):
            raise tgexception.TGException('Unable to conform the gremlin query to bytecode.')

        another_step = gremlin[0] == '.'

        gremlin = gremlin[1:]

    ret += ']]'

    return ret


class ResultDataDescriptor(tgquery.TGResultDataDescriptor):

    @property
    def dataType(self) -> tgquery.TGQueryResultDataType:
        return self.__type

    @property
    def containedDataSize(self) -> int:
        return None if self._container is None else len(self._container)

    @property
    def isMap(self) -> bool:
        return self.__type == tgquery.TGQueryResultDataType.Map

    @property
    def isArray(self) -> bool:
        return self.__type == tgquery.TGQueryResultDataType.List

    @property
    def hasConcreteType(self) -> bool:
        return self._attrType is not None and self._attrType != tgmodel.TGAttributeType.Invalid

    @property
    def attributeType(self) -> tgmodel.TGAttributeType:
        return self._attrType

    @property
    def sysobj(self) -> tgmodel.TGSystemObject:
        return self._sysobj

    @property
    def keyDescriptor(self):
        return self._keyDesc

    @property
    def valueDescriptor(self):
        return self._valueDesc

    @property
    def containedDescriptors(self):
        return self._container

    def __getitem__(self, position: int):
        return None if self._container is None else self._container[position]

    def __init__(self, char: str):
        self.__type = tgquery.TGQueryResultDataType.getTypeForChar(char)
        self._container: typing.List[tgquery.TGResultDataDescriptor] = None
        self._attrType: tgmodel.TGAttributeType = None
        if self.__type == tgquery.TGQueryResultDataType.Scalar:
            self._attrType = tgmodel.TGAttributeType.fromTypeChar(char)
        self._keyDesc: tgquery.TGResultDataDescriptor = None
        self._valueDesc: tgquery.TGResultDataDescriptor = None
        self._sysobj: tgmodel.TGSystemObject = None


class ResultMetaDataDesc(tgquery.TGResultSetMetaData):

    def __init__(self, annot: str = None):
        self.__annot = annot
        self.__dataDesc: TGResultDataDescriptor = None

    @property
    def resultDataDescriptor(self) -> TGResultDataDescriptor:
        return self.__dataDesc

    @property
    def resultType(self) -> TGQueryResultDataType:
        return tgquery.TGQueryResultDataType.Unknown if self.__dataDesc is None else self.__dataDesc.dataType

    @property
    def annotation(self) -> str:
        return self.__annot

    def initialize(self, gmd: tgmodel.TGGraphMetadata):
        if self.__annot is not None:
            self.__dataDesc = self.__generateDataDescriptors(gmd, self.__annot)

    def __split(self, insideAnnot: str) -> typing.List[str]:
        ret = []
        insideCount = 0
        startIdx = 0
        for curIdx in range(len(insideAnnot)):
            cur = insideAnnot[curIdx]
            if cur == ',' and insideCount == 0:
                ret.append(insideAnnot[startIdx:curIdx])
                startIdx = curIdx + 1
            if cur == '(' or cur == '{' or cur == '[':
                insideCount += 1
            if cur == ')' or cur == '}' or cur == ']':
                insideCount -= 1
        if startIdx < len(insideAnnot):
            ret.append(insideAnnot[startIdx:])
        return ret

    def __generateDataDescriptors(self, gmd: tgmodel.TGGraphMetadata, annot: str) -> ResultDataDescriptor:
        if len(annot) == 0:
            return None
        desc = ResultDataDescriptor(annot[0])
        if desc.isMap:
            desc._keyDesc = ResultDataDescriptor('s')           # Corresponds to string scalar type
            desc._valueDesc = ResultDataDescriptor('S')         # Corresponds to generic scalar type
        elif desc.isArray:
            subDesc = self.__generateDataDescriptors(gmd, annot[1:-1])
            desc._container = [subDesc]
        elif desc.dataType == tgquery.TGQueryResultDataType.Path or\
                desc.dataType == tgquery.TGQueryResultDataType.Tuple:
            start = 1 if desc.dataType == tgquery.TGQueryResultDataType.Tuple else 2
            data = self.__split(annot[start:-1])
            cont = []
            for datum in data:
                cont.append(self.__generateDataDescriptors(gmd, datum))
            desc._container = cont
        elif desc.dataType == tgquery.TGQueryResultDataType.Unknown:
            gLogger.log(TGLevel.Info, "Failed to initialize data descriptor with type annotation: %s" % annot)
            desc = None
        return desc


class GremlinElementType(enum.Enum):
    Invalid = 0
    List = 1
    Attr = 2
    AttrValue = 3
    AttrValueTransient = 4
    Entity = 5
    Map = 6
    MixedList = 7
    Void = 8

    @classmethod
    def fromId(cls, id: int):
        for gEntityType in GremlinElementType:
            if gEntityType.value == id:
                return gEntityType
        return GremlinElementType.Invalid


class ResultSetImpl(tgquery.TGResultSet, Generic[T]):

    def __init__(self, conn: tgconn.TGConnection, resultId: int):
        self.__conn__: tgconn.TGConnection = conn
        self.__resultId__: int = resultId
        self.__messageBytes: bytes = None
        self.__resultList__: List[T] = None
        self.__isOpen__: bool = True
        self.__currPos__: int = -1
        self.__resultTypeAnnot: str = None
        self.__metadata: ResultMetaDataDesc = None
        self.__length: int = 0
        self.__exceptions: typing.List[BaseException] = []

    @property
    def hasExceptions(self) -> bool:
        return False

    @property
    def exceptions(self) -> List[Exception]:
        return self.__exceptions

    @property
    def position(self) -> int:
        ret: int = -1
        if self.isOpen:
            ret = self.__currPos__
        return ret

    @property
    def messageBytes(self) -> bytes:
        return self.__messageBytes

    @messageBytes.setter
    def messageBytes(self, val: bytes):
        self.__messageBytes = val
        if val is not None:
            import tgdb.impl.pduimpl as pduimpl
            instream = pduimpl.ProtocolDataInputStream(bytearray(val))
            self.__length = instream.readInt()
        else:
            self.__length = 0

    def __len__(self):
        return self.__length

    def __iter__(self):
        import tgdb.impl.pduimpl as pduimpl
        if self.__messageBytes is not None:
            instream = pduimpl.ProtocolDataInputStream(bytearray(self.__messageBytes))
            for toYield in ResultSetImpl.parseGremlinList(instream, self.__conn__.graphObjectFactory):
                yield toYield

    def __getitem__(self, key):
        if isinstance(key, int):
            if self.__resultList__ is not None:
                return self.__resultList__[key]
            useKey = key
            if key < 0:
                useKey = len(self) + key
            if useKey < 0 or useKey >= len(self):
                raise KeyError("Key requested is too large! Key: {0}".format(key))
            idx = 0
            for result in self:
                if idx == useKey:
                    return result
        else:
            import tgdb.log as tglog
            tglog.gLogger.log(tglog.TGLevel.Warning, "Key of type %s is not supported for indexing.", str(type(key)))
            raise NotImplementedError("Key type given ({0}) is not supported.".format(str(type(key))))

    def close(self):
        self.__isOpen__ = False

    def skip(self, count):
        if count < 0:
            raise tgexception.TGException("Cannot skip a negative number!")
        self.__currPos__ += count

    def toCollection(self) -> [T]:
        if self.__resultList__ is None:
            self.__resultList__ = list(self)
        return self.__resultList__

    @property
    def isOpen(self) -> bool:
        return self.__isOpen__

    def addEntityToResultSet(self, ent: T):
        self.__resultList__.append(ent)

    @property
    def resultId(self) -> int:
        return self.__resultId__

    @property
    def hasNext(self) -> bool:
        ret: bool = False
        if self.isOpen and len(self.__resultList__) > 0:
            ret = self.__currPos__ < len(self.__resultList__) - 1
        return ret

    @property
    def size(self) -> int:
        return self.__length

    def setResultTypeAnnotation(self, annot: str):
        if annot is not None and len(annot) > 0:
            self.__resultTypeAnnot = annot
            self.__metadata = ResultMetaDataDesc(annot)
            try:
                self.__metadata.initialize(self.__conn__.graphMetadata)
            except BaseException as e:
                gLogger.log(TGLevel.Warning, "Failed to initialize result set metadata with error %s." % str(e))
                self.__metadata = None

    @property
    def metadata(self) -> tgquery.TGResultSetMetaData:
        return self.__metadata

    @classmethod
    def parseGremlinObject(cls, instream: tgpdu.TGInputStream, gof: tgmodel.TGGraphObjectFactory,
                           gremlinObjectType: GremlinElementType) -> Any:
        ret = None
        dummyNode: tgmodel.TGNode = gof.createNode(None)
        if gremlinObjectType is GremlinElementType.Entity:
            ret = AbstractEntity.readNewEntityFromStream(instream, gof)
        elif gremlinObjectType is GremlinElementType.Attr:
            ret = AbstractAttribute.createFromStream(dummyNode, instream)
        elif gremlinObjectType is GremlinElementType.AttrValue or gremlinObjectType is\
                GremlinElementType.AttrValueTransient:
            temp: tgmodel.TGAttribute = AbstractAttribute.createFromStream(dummyNode, instream)
            ret = temp.value
        elif gremlinObjectType is GremlinElementType.List:
            ret = list(ResultSetImpl.parseGremlinList(instream, gof))
        elif gremlinObjectType is GremlinElementType.Map:
            ret = {}
            ResultSetImpl.parseGremlinDict(instream, gof, ret)
        elif gremlinObjectType is GremlinElementType.MixedList:
            ret = list(ResultSetImpl.parseGremlinMixedList(instream, gof))
        else:
            gLogger.log(TGLevel.Warning, "Unknown element type: %s", str(gremlinObjectType))
        return ret

    @classmethod
    def parseGremlinList(cls, instream: tgpdu.TGInputStream, gof: tgmodel.TGGraphObjectFactory) -> typing.Iterable[Any]:
        size: int = instream.readInt()
        elemType: GremlinElementType = GremlinElementType.fromId(instream.readByte())
        for i in range(size):
            res = ResultSetImpl.parseGremlinObject(instream, gof, elemType)
            if res is not None:
                yield res

    @classmethod
    def parseGremlinMixedList(cls, instream: tgpdu.TGInputStream, gof: tgmodel.TGGraphObjectFactory
                              ) -> typing.Iterable[Any]:
        size: int = instream.readInt()
        for i in range(size):
            elemType: GremlinElementType = GremlinElementType.fromId(instream.readByte())
            res = ResultSetImpl.parseGremlinObject(instream, gof, elemType)
            if res is not None:
                yield res

    @classmethod
    def parseGremlinDict(cls, instream: tgpdu.TGInputStream, gof: tgmodel.TGGraphObjectFactory, addTo: Dict[str, Any]):
        size: int = instream.readInt()
        for i in range(size):
            key: str = instream.readUTF()
            elemType: GremlinElementType = GremlinElementType.fromId(instream.readByte())
            res = ResultSetImpl.parseGremlinObject(instream, gof, elemType)
            if res is not None:
                addTo[key] = res

    @classmethod
    def parseEntityStream(cls, instream: tgpdu.TGInputStream, gof: tgmodel.TGGraphObjectFactory, requestId: int, count: int,
                          resultCount: int = -1) -> tgquery.TGResultSet:
        incResCount: bool = True
        if resultCount is None:
            resultCount = 0
        else:
            incResCount = False
        refMap: Dict[int, tgmodel.TGEntity]
        refMap = instream.referencemap
        resultSet = EntityResultSet(gof.connection, requestId)
        for i in range(count):
            entKind: tgmodel.TGEntityKind = tgmodel.TGEntityKind.fromKindId(instream.readByte())
            if entKind is tgmodel.TGEntityKind.InvalidKind:
                # Log that we have an invalid kind entity, as this should not happen with any use case known
                gLogger.log(TGLevel.Warning, "Received invalid entity kind: %s", str(entKind))
                continue
            id: int = instream.readLong()
            entity: tgmodel.TGEntity
            if id in refMap:
                entity = refMap[id]
                entity.readExternal(instream)
            else:
                entity = AbstractEntity.readNewEntityFromStream(instream, gof, entKind)
                if incResCount:
                    resultCount += 1
            if resultSet is not None and resultSet.size < resultCount:
                resultSet.addEntityToResultSet(entity)
            # Log the information
            gLogger.logEntity(TGLevel.Debug, entity)

        return resultSet


class EntityResultSet(ResultSetImpl[tgmodel.TGEntity]):

    def __init__(self, conn: tgconn.TGConnection, resultId: int):
        super().__init__(conn, resultId)


class ObjectResultSet(ResultSetImpl[Any]):

    def __init__(self, conn: tgconn.TGConnection, resultId: int):
        super().__init__(conn, resultId)



