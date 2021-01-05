"""
.. Necessary to reject for documentation building
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
 *  File name :query.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: query.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates channel interfaces
 """


import enum
import abc
import typing
import tgdb.utils as tgutils
import tgdb.pdu as tgpdu
import tgdb.exception as tgexception
import tgdb.model as tgmodel


class TGQueryOption(tgutils.TGProperties, tgpdu.TGSerializable):
    """Stores information about the query.

    .. attribute:: prefetchsize

        :type: int
    .. attribute:: traversalDepth

        :type: int
    .. attribute:: edgeLimit

        :type: int
    .. attribute:: queryExpr

        :type: str
    .. attribute:: edgeExpr

        :type: str
    .. attribute:: sortAttrName
        The attribute to sort the results (only used when usign TGQL, not Gremlin).
        :type: str
    .. attribute:: sortOrderDesc
        Whether to sort on the attribute descending or ascending, with True for sort descending and False for \
        ascending (only used when usign TGQL, not Gremlin).
        :type: bool
    .. attribute:: sortResultLimit
        The maximum number of objects to return in the result (only used when usign TGQL, not Gremlin).
        :type: int
    """

    QUERY_OPTION_FETCHSIZE = 'fetchsize'
    QUERY_OPTION_TRAVERSALDEPTH = 'traversaldepth'
    QUERY_OPTION_EDGELIMIT = 'edgelimit'
    QUERY_OPTION_BATCHSIZE = 'batchSize'
    QUERY_OPTION_QUERY_EXPR = 'queryExpr'
    QUERY_OPTION_EDGE_EXPR = 'edgeExpr'
    QUERY_OPTION_TRAVERSE_EXPR = 'traverseExpr'
    QUERY_OPTION_END_EXPR = 'endExpr'
    QUERY_OPTION_SORT_ATTR_NAME = 'sortAttrName'
    QUERY_OPTION_SORT_DESC = 'sortOrderDesc'
    QUERY_OPTION_SORT_RESULT_LIMIT = 'sortResultLimit'

    def __init__(self):
        super().__init__()
        self[TGQueryOption.QUERY_OPTION_FETCHSIZE] = -1
        self[TGQueryOption.QUERY_OPTION_TRAVERSALDEPTH] = -1
        self[TGQueryOption.QUERY_OPTION_EDGELIMIT] = -1
        self[TGQueryOption.QUERY_OPTION_BATCHSIZE] = -1
        self[TGQueryOption.QUERY_OPTION_QUERY_EXPR] = None
        self[TGQueryOption.QUERY_OPTION_EDGE_EXPR] = None
        self[TGQueryOption.QUERY_OPTION_TRAVERSE_EXPR] = None
        self[TGQueryOption.QUERY_OPTION_END_EXPR] = None
        self[TGQueryOption.QUERY_OPTION_SORT_ATTR_NAME] = None
        self[TGQueryOption.QUERY_OPTION_SORT_DESC] = False
        self[TGQueryOption.QUERY_OPTION_SORT_RESULT_LIMIT] = 0

    @property
    def prefetchSize(self):
        """
        Represents the maximum number of objects to fetch during a single request/response to the server.

        :return: The number of bytes to prefetch.
        :rtype: int
        """
        return self[TGQueryOption.QUERY_OPTION_FETCHSIZE]

    @prefetchSize.setter
    def prefetchSize(self, v: int):
        """
        Represents the maximum number of objects to fetch during a single request/response to the server.

        :param v: The number of bytes to prefetch.
        :type v: int
        :returns: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_FETCHSIZE] = v

    @property
    def traversalDepth(self):
        """
        Represents the maximum depth of recursion for the query (only used when using TGQL, not Gremlin).

        :return: Maximum traversal depth
        :rtype: int
        """
        return int(self[TGQueryOption.QUERY_OPTION_TRAVERSALDEPTH])

    @traversalDepth.setter
    def traversalDepth(self, v: int):
        """
        Represents the maximum depth of recursion for the query (only used when using TGQL, not Gremlin).

        :param v: Maximum traversal depth
        :type v: int
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_TRAVERSALDEPTH] = v

    @property
    def edgeLimit(self) -> int:
        """
        Represents the maximum number of edges for the server to respond with (only used when using TGQL, not Gremlin).

        :return: The maximum number of edges from the server.
        :rtype: int
        """
        return int(self[TGQueryOption.QUERY_OPTION_EDGELIMIT])

    @edgeLimit.setter
    def edgeLimit(self, v: int):
        """
        Represents the maximum number of edges for the server to respond with (only used when using TGQL, not Gremlin).

        :param v: The maximum number of edges from the server.
        :type v: int
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_EDGELIMIT] = v

    @property
    def batchSize(self):
        """

        :return: Whether the sort should be descending.
        :rtype: int
        """
        return self[TGQueryOption.QUERY_OPTION_BATCHSIZE]

    @batchSize.setter
    def batchSize(self, v: int):
        """

        :param v:
        :type v:
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_BATCHSIZE] = v

    @property
    def queryExpr(self) -> str:
        """
        Gets the query string.

        :return: The query string.
        :rtype: str
        """
        return self[TGQueryOption.QUERY_OPTION_QUERY_EXPR]

    @queryExpr.setter
    def queryExpr(self, v: str):
        """
        Sets the query string.

        :param v: The query string
        :type v: str
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_QUERY_EXPR] = v

    @property
    def edgeExpr(self) -> str:
        """
        Gets the expression to traverse an edge to another node (only used when usign TGQL, not Gremlin).

        :return: The expression to traverse an edge to another node (only used when usign TGQL, not Gremlin).
        :rtype: str
        """
        return self[TGQueryOption.QUERY_OPTION_EDGE_EXPR]

    @edgeExpr.setter
    def edgeExpr(self, v: str):
        """

        :param v:
        :type v:
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_EDGE_EXPR] = v

    @property
    def traverseExpr(self) -> str:
        """

        :return: Whether the sort should be descending.
        :rtype: int
        """
        return self[TGQueryOption.QUERY_OPTION_TRAVERSE_EXPR]

    @traverseExpr.setter
    def traverseExpr(self, v: str):
        """

        :param v:
        :type v:
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_TRAVERSE_EXPR] = v

    @property
    def endExpr(self) -> str:
        """

        :return: Whether the sort should be descending.
        :rtype: int
        """
        return self[TGQueryOption.QUERY_OPTION_END_EXPR]

    @endExpr.setter
    def endExpr(self, v: str):
        """

        :param v:
        :type v:
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_END_EXPR] = v

    @property
    def sortAttrName(self) -> str:
        """

        :return: Whether the sort should be descending.
        :rtype: int
        """
        return self[TGQueryOption.QUERY_OPTION_SORT_ATTR_NAME]

    @sortAttrName.setter
    def sortAttrName(self, v: str):
        """

        :param v:
        :type v:
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_SORT_ATTR_NAME] = v

    @property
    def sortOrderDesc(self) -> bool:
        """

        :return: Whether the sort should be descending.
        :rtype: int
        """
        return self[TGQueryOption.QUERY_OPTION_SORT_DESC]

    @sortOrderDesc.setter
    def sortOrderDesc(self, v: bool):
        """

        :param v:
        :type v:
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_SORT_DESC] = v

    @property
    def sortResultLimit(self) -> int:
        """

        :return: The sort result limit.
        :rtype: int
        """
        return self[TGQueryOption.QUERY_OPTION_SORT_RESULT_LIMIT]

    @sortResultLimit.setter
    def sortResultLimit(self, v: int):
        """

        :param v:
        :type v:
        :return: Nothing
        :rtype: None
        """
        self[TGQueryOption.QUERY_OPTION_SORT_RESULT_LIMIT] = v

    def writeExternal(self, os: tgpdu.TGOutputStream, writeSort: bool = False, writeQueryStr: bool = False):
        """Handles communicating the query options with the server.

        :param os: The stream to write this to.
        :param writeSort: Whether to write the sort information.
        :param writeQueryStr: Whether to write the query string options.
        """
        os.writeInt(self.prefetchSize)
        os.writeShort(self.batchSize)
        os.writeShort(self.traversalDepth)
        os.writeShort(self.edgeLimit)

        if writeSort:
            if self.sortAttrName is None:
                os.writeBoolean(False)
            else:
                os.writeBoolean(True)
                os.writeUTF(self.sortAttrName)
                os.writeBoolean(self.sortOrderDesc)
                os.writeInt(self.sortResultLimit)
        if writeQueryStr:
            def writeStrIfThere(string: str):
                if string is None:
                    os.writeBoolean(True)
                else:
                    os.writeBoolean(False)
                    os.writeUTF(string)

            writeStrIfThere(self.queryExpr)
            writeStrIfThere(self.edgeExpr)
            writeStrIfThere(self.traverseExpr)
            writeStrIfThere(self.endExpr)

    def readExternal(self, istream: tgpdu.TGInputStream):
        pass


DefaultQueryOption = TGQueryOption()


class TGQueryErrorStatus(enum.Enum):
    """
    This class converts a query result from an integer to the meaning behind it.
    """

    Invalid = (8100, "Unknown query error")
    ProviderNotInitialized = (8101, "Provider not initialized")
    ParsingError = (8102, "Parsing error")
    StepNotSupported = (8103, "Gremlin step not supported")
    StepNotAllowed = (8104, "Gremlin step not allowed")
    StepArgMissing = (8105, "Gremlin step argument missing")
    StepArgNotSupported = (8106, "Gremlin step argument not supported")
    StepMissing = (8107, "Gremlin step missing")
    NotDefined = (8108, "Query not defined")
    AttrDescNotFound = (8109, "Attribute descriptor not found")
    EdgeTypeNotFound = (8110, "Edge type not found")
    NodeTypeNotFound = (8111, "Not type not found")
    InternalDataMismatchError = (8112, "Data type mismatch")
    StepSignatureNotSupported = (8113, "Gremlin step signature not supported")
    InvalidDataType = (8114, "Invalid data type")
    TGQueryExecSPFailure = (8115, "Executing stored procedure failed")
    TGQuerySPNotFound = (8116, "Stored procedure not found")
    TGQuerySPArgMissing = (8117, "Stored procedure argument missing")
    TGQueryStepArgInvalid = (8118, "Step argument invalid")
    TGQueryStepModulationInvalid = (8119, "Step modulation invalid")
    AccessDenied = (8120, "Query access denied")

    def __init__(self, result: int, errorMsg: str):
        self.__result = result
        self.__errMsg = errorMsg

    @property
    def result(self) -> int:
        """
        The integer representation of the query result. Must correspond to a result on the server.

        :return: The result.
        :rtype: int
        """
        return self.__result

    @property
    def errorMsg(self) -> str:
        """
        The message for the user.

        :return: The error message.
        :rtype: str
        """
        return self.__errMsg

    @classmethod
    def fromId(cls, id: int):
        """
        Converts an integer identifier to a :class:`TGQueryErrorStatus<tgdb.query.TGQueryErrorStatus>`.

        :param id: The integer identifier that the server has for a given result.
        :type id: int
        :return: The error status in a more readable and supportable manner.
        :rtype: tgdb.query.TGQueryErrorStatus
        """
        for qc in TGQueryErrorStatus:
            if qc.result == id:
                return qc
        return TGQueryErrorStatus.Invalid


class TGQueryException(tgexception.TGException):
    """Represents an exception occurring in the execution of a query."""

    def __init__(self, result: int, msg: typing.Optional[str]):
        super().__init__(TGQueryErrorStatus.fromId(result).errorMsg + ("" if msg is None else ": " + msg) + ".",
                         errorcode=result)
        self.__errStatus = TGQueryErrorStatus.fromId(result)

    @property
    def errStatus(self) -> TGQueryErrorStatus:
        """
        Gets the error status.

        :return: The error status
        :rtype: tgdb.query.TGQueryErrorStatus
        """
        return self.__errStatus


class TGQueryResultDataType(enum.Enum):
    """Data type used to describe the metadata of the resulting type of a gremlin query."""
    Unknown = (0, '')
    Object = (1, '')
    Entity = (2, '')
    Attr = (3, '')
    Node = (4, 'V')
    Edge = (5, 'E')
    List = (6, '[')
    Map = (7, '{')
    Tuple = (8, '(')
    Scalar = (9, 'S ' + ' '.join([attr.char for attr in tgmodel.TGAttributeType if attr.char is not None]))
    Path = (10, 'P')

    def __init__(self, ident, knownChars):
        self.__id = ident
        self.__knownChars = frozenset(knownChars.split())

    @classmethod
    def getTypeForChar(cls, char: str):
        """
        Converts a character into the appropriate query result data type.

        :param char: A single character.
        :type char: str
        :returns: The query result data type.
        :rtype: tgdb.query.TGQueryResultDataType
        """
        for qrdt in TGQueryResultDataType:
            if char in qrdt.__knownChars:
                return qrdt
        return TGQueryResultDataType.Unknown


class TGResultDataDescriptor(abc.ABC):
    """The resulting set's data descriptor."""

    @property
    @abc.abstractmethod
    def dataType(self) -> TGQueryResultDataType:
        """
        Gets the data type for this data descriptor.

        :returns: The result data type.
        :rtype: tgdb.query.TGQueryResultDataType
        """

    @property
    @abc.abstractmethod
    def containedDataSize(self) -> int:
        """
        Gets the number of descriptors in this descriptor (should only not be None when this is of type list or \
        path).

        :returns: The number of data types that are sub-types.
        :rtype: int
        """

    @property
    @abc.abstractmethod
    def isMap(self) -> bool:
        """
        Whether this descriptor represents a map.

        :returns: Whether this descriptor represents a map.
        :rtype: bool
        """

    @property
    @abc.abstractmethod
    def isArray(self) -> bool:
        """
        Whether this descriptor represents an array.

        :returns: Whether this descriptor represents an array.
        :rtype: bool
        """

    @property
    @abc.abstractmethod
    def hasConcreteType(self) -> bool:
        """
        Whether this attribute has a concrete type (should be False if this is a composite type, like a list or a \
        map).

        :returns: Whether this descriptor represents a concrete type.
        :rtype: bool
        """

    @property
    @abc.abstractmethod
    def attributeType(self) -> tgmodel.TGAttributeType:
        """
        Gets the attribute type (if this descriptor has type attribute).

        :returns: The attribute type.
        :rtype: tgdb.model.TGAttributeType
        """

    @property
    @abc.abstractmethod
    def sysobj(self) -> tgmodel.TGSystemObject:
        """
        The corresponding system object (if this is of type object).

        :returns: The corresponding system object.
        :rtype: tgdb.model.TGSystemObject
        """

    @property
    @abc.abstractmethod
    def keyDescriptor(self):
        """
        Gets the descriptor used as the key for this descriptor. Should only be defined when this corresponds to a map.

        :returns: The descriptor for the key.
        :rtype: tgdb.query.TGResultDataDescriptor
        """

    @property
    @abc.abstractmethod
    def valueDescriptor(self):
        """
        Gets the descriptor used as the value for this descriptor. Should only be defined when this corresponds to a \
        map.

        :returns: The descriptor for the value.
        :rtype: tgdb.query.TGResultDataDescriptor
        """

    @property
    @abc.abstractmethod
    def containedDescriptors(self):
        """
        Gets the descriptors contained by this descriptor. Should only be defined when this is of type list or path.

        :returns: The contained descriptors.
        :rtype: typing.List[tgdb.query.TGResultDataDescriptor]
        """

    @abc.abstractmethod
    def __getitem__(self, position: int):
        """
        Gets the descriptor at the position described

        :param position: The position of the described object
        :type position: int
        :returns: The descriptor at that position. Should raise an error if this data descriptor does not have \
            sub-descriptors.
        :rtype: tgdb.query.TGResultDataDescriptor
        """


class TGResultSetMetaData(abc.ABC):
    """The resulting set's metadata descriptor."""

    @property
    @abc.abstractmethod
    def resultDataDescriptor(self) -> TGResultDataDescriptor:
        """
        Gets the data descriptor for this result set's metadata.

        :returns: The result set's data descriptor.
        :rtype: tgdb.query.TGResultDataDescriptor
        """

    @property
    @abc.abstractmethod
    def resultType(self) -> TGQueryResultDataType:
        """
        Gets the result type of this result set's metadata.

        :returns: The top level result data type for this result.
        :rtype: tgdb.query.TGQueryResultDataType
        """

    @property
    @abc.abstractmethod
    def annotation(self) -> str:
        """
        Gets the result type of this result set's metadata.

        :returns: The annotation for this result. Can be used to determine the language.
        :rtype: str
        """


class TGQueryCommand(enum.Enum):
    """The query commands that the server recognizes."""
    Invalid = 0
    Create = 1
    Execute = 2
    ExecuteGremlin = 3
    ExecuteGremlinStr = 4
    ExecuteID = 5
    Close = 6

    @classmethod
    def fromId(cls, id: int):
        """
        Converts a query command identifier into a query execution type.

        :param id: An integer identifier to get the query command.
        :type id: int
        :return: The query command.
        :rtype: tgdb.query.TGQueryCommand
        """
        for qc in TGQueryCommand:
            if qc.value is id:
                return qc
        return TGQueryCommand.Invalid


class TGResultSet(abc.ABC):
    """
    This represents a single query response from the server

    It will in the future handle acting like a cursor.
    """

    @property
    @abc.abstractmethod
    def hasExceptions(self) -> bool:
        """
        Returns whether the server responded with any exceptions.

        :returns: Whether the query caused any exceptions.
        :rtype: bool
        """

    @property
    @abc.abstractmethod
    def exceptions(self) -> typing.List[Exception]:
        """
        Returns what the server's exceptions are.

        To be either Null or an empty List if there are no exceptions.

        :returns: A list of the exceptions that occurred during recovery of the data requested by the query.
        :rtype: typing.List[Exception]
        """
        pass

    @property
    @abc.abstractmethod
    def metadata(self) -> TGResultSetMetaData:
        """
        Returns the metadata for this query result set.

        :returns: This queries metadata.
        :rtype: tgdb.query.TGResultSetMetaData
        """

    @property
    @abc.abstractmethod
    def position(self) -> int:
        """
        The current position of this cursor.

        :returns: The current position.
        :rtype: int
        """

    @abc.abstractmethod
    def __len__(self):
        """
        Gets the number of elements for this cursor.

        :returns: Number of elements in this query.
        :rtype: int
        """

    @abc.abstractmethod
    def __iter__(self):
        """
        Iterates through the elements of this result.

        :return: An iterator for the elements of this query. Should match up with the result metadata.
        :rtype: typing.Iterator[typing.Any]
        """

    @abc.abstractmethod
    def __getitem__(self, item):
        """
        Gets a particular item at the given index of this query.

        :param item: The index for the item to get.
        :type item: int
        :return: The query element requested.
        """

    @abc.abstractmethod
    def skip(self, count):
        """Skips the objects

        :param count: The number of objects to skip past.
        :type count: int
        :returns: Nothing
        :rtype: None
        """

    @abc.abstractmethod
    def toCollection(self) -> typing.List:
        """
        Gets all of the objects as a list

        Note: Do not use when the number of objects returned could be large. Would recommend iterating through because \
        of the penalty that Python has with deserializing objects.

        :returns: A list of all of the query results. The types should match up with those described by the metadata.
        :rtype: typing.List
        """
        pass

    @property
    @abc.abstractmethod
    def size(self) -> int:
        """
        Gets the number of objects returned as an integer.

        :returns: The number of objects in this query.
        :rtype: int
        """


class TGQuery(abc.ABC):
    """Used for parameterized queries."""

    @abc.abstractmethod
    def __setitem__(self, key, value):
        """
        Sets a particular parameter to the value.

        :param key: The key parameter value.
        :param value: The value of that parameter.
        :returns: Nothing
        :rtype: None
        """

    @abc.abstractmethod
    def __getitem__(self, item):
        """
        Gets a particular parameter's value.

        :param item: The key parameter value.
        :returns: The value of that parameter.
        """

    @property
    @abc.abstractmethod
    def option(self) -> TGQueryOption:
        """
        Gets the query options for this query.

        :returns: The options for this query.
        :rtype: tgdb.query.TGQueryOption
        """

    @option.setter
    def option(self, v: TGQueryOption):
        """
        Sets the query options for this query.
        :param v: The query option object for this query.
        :type v: tgdb.query.TGQueryOption
        :returns: Nothing
        :rtype: None
        """
        self._option(v)

    def _option(self, v: TGQueryOption):
        pass

    @abc.abstractmethod
    def execute(self) -> TGResultSet:
        """
        Execute this query with the current parameter list.

        :returns: The results from executing this query with the given parameters.
        :rtype: tgdb.query.TGResultSet
        """

    @abc.abstractmethod
    def close(self):
        """
        Close the query and clean up the compiled query on the server side.

        :returns: Nothing
        :rtype: None
        """

