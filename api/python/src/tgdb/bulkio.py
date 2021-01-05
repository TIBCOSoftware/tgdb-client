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
 *  File name :bulkio.py
 *  Created on: 9/3/2019
 *  Created by: derek
 *
 *		SVN Id:
 *
 *  This file encapsulates all the unit test case for bulk import test
 """

import typing
import enum
import abc

import tgdb.model as tgmodel


class TGLoadOptions(enum.Enum):
    Invalid = "INVALID"
    Insert = "insert"
    Upsert = "upsert"

    @classmethod
    def findVal(cls, loadopt: str):
        for val in cls:
            if val.value == loadopt:
                return val
        return cls.Invalid


class TGErrorOptions(enum.Enum):
    Invalid = "INVALID"
    Stop = "stop"
    Ignore = "ignore"

    @classmethod
    def findVal(cls, erroropt: str):
        for val in cls:
            if val.value == erroropt:
                return val
        return cls.Invalid


class TGDateFormat(enum.Enum):
    Invalid = "INVALID"
    MDY = "MDY"
    DMY = "DMY"
    YMD = "YMD"

    @classmethod
    def findVal(cls, dtformat: str):
        for val in cls:
            if val.value == dtformat:
                return val
        return cls.Invalid


class TGBulkExportDescriptor(abc.ABC):
    """Keeps track of various information about all of the entities of a particular type that are exported."""

    @abc.abstractmethod
    def __len__(self) -> int:
        """Gets the number of entities for this type that have been exported."""

    @property
    @abc.abstractmethod
    def entitytype(self) -> tgmodel.TGEntityType:
        """Gets the type of this entity"""

    @property
    @abc.abstractmethod
    def data(self) -> typing.Optional[typing.Any]:
        """Gets the data (if any) for this type.

        :returns: None if no data exists or a DataFrame that is correctly formatted by the attributes
        :rtype: pandas.DataFrame
        """

    @abc.abstractmethod
    def addHandler(self, handler: typing.Callable[[typing.Any], None]):
        """Add a handler to the list.

        :param handler: Handles entities of this type when they arrive from the server during an export.
        :type handler: typing.Callable[[TGBulkImportDescriptor], None]
        """

    @abc.abstractmethod
    def clearHandlers(self):
        """Removes all handlers for this export of this type."""

    @property
    @abc.abstractmethod
    def hasHandler(self) -> bool:
        """Is there more than one handler for this type."""


class TGBulkImport(abc.ABC):
    """Represents a single import session with the server.

    Properties
    --------------------
    nodeBatchSize: int
        Maximum number of nodes to import in a single request/response with the server.
    edgeBatchSize: int
        Maximum number of nodes to import in a single in request/response with the server.
    idColumnName: str
        Name of the column to get the ids from in all of the DataFrames imported to the server. Default is 'id'.
    fromColumnName: str
        Name of the column to get the from vertex from in all of the DataFrames imported to the server. Only used when
        importing edges. Default is 'from'.
    toColumnName: str
        Name of the column to get the to vertex from in all of the DataFrames imported to the server. Only used when
        importing edges. Default is 'to'.
    """

    @abc.abstractmethod
    def _set_nodeBatchSize(self, value: int):
        pass

    @abc.abstractmethod
    def _get_nodeBatchSize(self) -> int:
        pass

    nodeBatchSize = property(_get_nodeBatchSize, _set_nodeBatchSize)

    @abc.abstractmethod
    def _set_edgeBatchSize(self, value: int):
        pass

    @abc.abstractmethod
    def _get_edgeBatchSize(self) -> int:
        pass

    edgeBatchSize = property(_get_edgeBatchSize, _set_edgeBatchSize)

    @abc.abstractmethod
    def _set_idColumnName(self, value: str):
        pass

    @abc.abstractmethod
    def _get_idColumnName(self) -> str:
        pass

    idColumnName = property(_get_idColumnName, _set_idColumnName)

    @abc.abstractmethod
    def _set_fromColumnName(self, value: str):
        pass

    @abc.abstractmethod
    def _get_fromColumnName(self) -> str:
        pass

    fromColumnName = property(_get_fromColumnName, _set_fromColumnName)

    @abc.abstractmethod
    def _set_toColumnName(self, value: str):
        pass

    @abc.abstractmethod
    def _get_toColumnName(self) -> str:
        pass

    toColumnName = property(_get_toColumnName, _set_toColumnName)

    @abc.abstractmethod
    def importNodes(self, data: typing.Any, nodetype: typing.Union[str, tgmodel.TGNodeType]):
        """Imports the nodes by the specified type.

        Note: can only import nodes of one type one time during an import session. Trying to do more than one will cause
        an error.

        :param data: The data to send over to the server.
        :type data: pandas.DataFrame
        :param nodetype: the name of the type (or the type itself) of the nodes imported under this command.

        :returns: the bulk import descriptor for the nodes imported
        :rtype: tgdb.admin.TGImportDescriptor
        """

    @abc.abstractmethod
    def importEdges(self, data: typing.Any, edgetype: typing.Union[str, tgmodel.TGEdgeType] = None,
                    edgedir: tgmodel.DirectionType = tgmodel.DirectionType.Directed):
        """Imports the edges of the specified type or direction (if no type was specified).

        Note: can only import edges of one type one time during an import session. Trying to do more than one will cause
        an error.

        :param data: The data to send over to the server.
        :type data: pandas.DataFrame
        :param defaultType: the name of the type (or the type itself) of the edges imported under this command, or
            (when None) to use a default edgetype corresponding to the defaultDir. Default is None.
        :param defaultDir: the direction for these imported edges. Will be ignored if defaultType is not None. Default
            is tgdb.DirectionType.Directed

        :returns: the bulk import descriptor for the edges imported
        :rtype: tgdb.admin.TGImportDescriptor
        """

    @abc.abstractmethod
    def finish(self):
        """Tells the server to clean up the import process."""


class TGBulkExport(abc.ABC):
    """Represents a single export session with the server.

    Properties
    --------------------
    idColumnName: str
        Name of the column to set the ids to in all of the DataFrames exported from the server. Default is 'id'.
    fromColumnName: str
        Name of the column to set the from vertex to in all of the DataFrames exported from the server. Only used when
        exporting edges. Default is 'from'.
    toColumnName: str
        Name of the column to set the to vertex to in all of the DataFrames exported from the server. Only used when
        exporting edges. Default is 'to'.
    """

    @abc.abstractmethod
    def _set_idColumnName(self, value: str):
        pass

    @abc.abstractmethod
    def _get_idColumnName(self) -> str:
        pass

    idColumnName = property(_get_idColumnName, _set_idColumnName)

    @abc.abstractmethod
    def _set_fromColumnName(self, value: str):
        pass

    @abc.abstractmethod
    def _get_fromColumnName(self) -> str:
        pass

    fromColumnName = property(_get_fromColumnName, _set_fromColumnName)

    @abc.abstractmethod
    def _set_toColumnName(self, value: str):
        pass

    @abc.abstractmethod
    def _get_toColumnName(self) -> str:
        pass

    toColumnName = property(_get_toColumnName, _set_toColumnName)

    @abc.abstractmethod
    def __getitem__(self, typename: str) -> TGBulkExportDescriptor:
        """Gets the TGBulkExportDescriptor corresponding to an exported type.

        Intended to set the handlers for a particular entity type.

        :param typename: The name of the type of the wanted export descriptor.

        :returns: The desired exported descriptor.
        """

    @abc.abstractmethod
    def exportEntities(self, handler: typing.Optional[typing.Callable[[TGBulkExportDescriptor], None]] = None,
                       callPassedHandlerOnAll: bool = False):
        """Exports all of the entities from the server.

        :param handler: Handles all of the entities exported. Entities will be exported in batches. The data exported
            will be set in the data property of the BulkExportDescriptor that is passed to the handler.
        :param callPassedHandlerOnAll: If True will call the passed handler on all batches called, regardless of whether
            that export descriptor already has any handlers. Will call all handlers for that export descriptor
            regardless of this flag.
        """
