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
 *  File name :admin
 *  Created on: 10/24/19
 *  Created by: suresh
 *
 *		SVN Id: $Id:  $
 *
 """
import abc
import typing
import enum
import datetime

import tgdb.version as tgvers
import tgdb.model as tgmodel
import tgdb.connection as tgconn
import tgdb.log as tglog


class TGImportDescriptor(abc.ABC):
    """Keeps track of the types that were sent to the server."""

    @property
    @abc.abstractmethod
    def typename(self) -> str:
        """The name of the type."""

    @property
    @abc.abstractmethod
    def isnode(self) -> bool:
        """Whether this type is a node."""

    @property
    @abc.abstractmethod
    def numinstances(self) -> int:
        """The number of instances of this type."""


class TGCacheStat(abc.ABC):
    @property
    @abc.abstractmethod
    def cacheMaxEntries(self) -> int:
        """Gets the maximum number of entries for this cache."""

    @property
    @abc.abstractmethod
    def cacheNumEntries(self) -> int:
        """Gets the current number of entries in this cache."""

    @property
    @abc.abstractmethod
    def cacheHits(self) -> int:
        """Gets the number of hits for this cache."""

    @property
    @abc.abstractmethod
    def cacheMisses(self) -> int:
        """Gets the number of misses for this cache."""

    @property
    @abc.abstractmethod
    def cacheMaxMemory(self) -> int:
        """Gets the maximum memory that this cache can use."""


class TGEntityCache(TGCacheStat, abc.ABC):

    @property
    @abc.abstractmethod
    def maxEntitySize(self) -> int:
        """Gets the size of the largest entity possible in the entity cache."""

    @property
    @abc.abstractmethod
    def averageEntitySize(self) -> int:
        """Gets the average size of the entities currently in the entity cache."""


class TGCacheStatistics(abc.ABC):

    @property
    @abc.abstractmethod
    def dataCache(self) -> TGCacheStat:
        """Gets the cache statistics for the data cache."""

    @property
    @abc.abstractmethod
    def indexCache(self) -> TGCacheStat:
        """Gets the cache statistics for the index cache."""

    @property
    @abc.abstractmethod
    def sharedCache(self) -> TGCacheStat:
        """Gets the cache statistics for the shared cache."""

    @property
    @abc.abstractmethod
    def entityCache(self) -> TGEntityCache:
        """Gets the cache statistics for the entity cache."""


class TGConnectionInfo(abc.ABC):

    @property
    @abc.abstractmethod
    def listenerName(self) -> str:
        """Gets the name of the listener."""

    @property
    @abc.abstractmethod
    def clientID(self) -> str:
        """Gets the client's identifier."""

    @property
    @abc.abstractmethod
    def sessionID(self) -> int:
        """Gets the session identifier."""

    @property
    @abc.abstractmethod
    def userName(self) -> str:
        """Gets the logged-in user's name for the connection."""

    @property
    @abc.abstractmethod
    def remoteAddress(self) -> str:
        """Gets the client-side address of the connection."""

    @property
    @abc.abstractmethod
    def createdTimeInSeconds(self) -> float:
        """Number of seconds that have elapsed since the connection was created."""


class TGDatabaseStatistics(abc.ABC):

    @property
    @abc.abstractmethod
    def dbSize(self) -> int:
        """Total size (includes both index and data size) of the database (in bytes)"""

    @property
    @abc.abstractmethod
    def dbPath(self) -> str:
        """Total size (includes both index and data size) of the database (in bytes)"""

    @property
    @abc.abstractmethod
    def numDataSegments(self) -> int:
        """Total number of data segments"""

    @property
    @abc.abstractmethod
    def dataSize(self) -> int:
        """Total ize of the data chunks of the database when in storage (in bytes)"""

    @property
    @abc.abstractmethod
    def dataUsed(self) -> int:
        """Size of the data used by the database in storage (in bytes)"""

    @property
    @abc.abstractmethod
    def dataFree(self) -> int:
        """Size of the data free by the database in storage (in bytes)"""

    @property
    @abc.abstractmethod
    def dataBlockSize(self) -> int:
        """Size of data blocks by the database in storage (in bytes)"""

    @property
    @abc.abstractmethod
    def numIndexSegments(self) -> int:
        """Total number of data segments"""

    @property
    @abc.abstractmethod
    def indexSize(self) -> int:
        """Total size of the index chunks of the database when in storage (in bytes)"""

    @property
    @abc.abstractmethod
    def usedIndexSize(self) -> int:
        """Size of the index used by the database in storage (in bytes)"""

    @property
    @abc.abstractmethod
    def freeIndexSize(self) -> int:
        """Size of the index free by the database in storage (in bytes)"""

    @property
    @abc.abstractmethod
    def blockSize(self) -> int:
        """Size of index blocks by the database in storage (in bytes)"""


class TGMemoryInfo(abc.ABC):

    @property
    @abc.abstractmethod
    def usedMemorySize(self) -> int:
        """Gets the used memory size (in bytes)"""

    @property
    @abc.abstractmethod
    def freeMemorySize(self) -> int:
        """Gets the free memory size (in bytes)"""

    @property
    @abc.abstractmethod
    def maxMemorySize(self) -> int:
        """Gets the max memory size (in bytes)"""

    @property
    @abc.abstractmethod
    def sharedMemoryFileLocation(self) -> str:
        """Gets location of the shared memory file."""


class TGNetListenerInfo(abc.ABC):

    @property
    @abc.abstractmethod
    def listenerName(self) -> str:
        """Gets name of the net listener."""

    @property
    @abc.abstractmethod
    def currentConnections(self) -> int:
        """Gets the number of currently connected clients."""

    @property
    @abc.abstractmethod
    def maxConnections(self) -> int:
        """Gets the maximum allowable number of connections."""

    @property
    @abc.abstractmethod
    def portNumber(self) -> str:
        """Gets the port number on the server side for this net listener"""


class MemoryType(enum.Enum):
    """Represents one of the two memory types: process-specific or shared."""
    PROCESS = 0
    SHARED = 1


class TGServerMemoryInfo(abc.ABC):

    @abc.abstractmethod
    def getMemoryInfo(self, type: MemoryType) -> TGMemoryInfo:
        """Gets the memory type specified by the type parameter"""


class ServerState(enum.Enum):
    """Represents the current server's state."""
    INVALID = -1
    CREATED = 0
    INIT_PHASE_0 = 1
    INIT_PHASE_1 = 2
    INITIALIZED = 3
    STARTED = 4
    SUSPENDED = 5
    INTERRUPTED = 6
    REQUESTSTOP = 7
    STOPPED = 8
    SHUTDOWN = 9

    @classmethod
    def fromId(cls, id: int):
        for ac in cls:
            if ac.value == id:
                return ac
        return cls.INVALID


class TGServerStatus(abc.ABC):
    """Information about the server."""

    @property
    @abc.abstractmethod
    def name(self) -> str:
        """The name of the server."""

    @property
    @abc.abstractmethod
    def version(self) -> tgvers.TGVersion:
        """The corresponding version of the server."""

    @property
    @abc.abstractmethod
    def state(self) -> ServerState:
        """The server's current state"""

    @property
    @abc.abstractmethod
    def processID(self) -> str:
        """The process identifier for the server."""

    @property
    @abc.abstractmethod
    def uptime(self) -> datetime.timedelta:
        """The server's uptime."""

    @property
    @abc.abstractmethod
    def serverPath(self) -> str:
        """The server's current working directory."""


class TGProcessorSetStatistics(abc.ABC):
    """Processing statistics for queries and transactions."""

    @property
    @abc.abstractmethod
    def processorsCount(self) -> int:
        """The total number of transaction processors (these are each capable of handling a transaction)"""

    @property
    @abc.abstractmethod
    def processedCount(self) -> int:
        """The total number of processed transactions."""

    @property
    @abc.abstractmethod
    def successfulCount(self) -> int:
        """The total number of successfully processed transactions."""

    @property
    @abc.abstractmethod
    def averageProcessingTime(self) -> float:
        """The average time for processing a transaction (in seconds)."""

    @property
    @abc.abstractmethod
    def pendingCount(self) -> int:
        """The number of pending transactions (transactions not yet finalized, successful or not)."""


class TGTransactionStatistics(TGProcessorSetStatistics, abc.ABC):
    """Represents statistics for processing transactions."""

    @property
    @abc.abstractmethod
    def transactionLoggerQueueDepth(self) -> int:
        """The depth of the transaction logger (the number of transactions sent to the database that are not being\
        processed right now. """


class TGUserInfo(tgmodel.TGSystemObject):
    """The user information."""


class TGServerInfo(abc.ABC):
    """Gets all relevent server information and statistics."""

    @property
    @abc.abstractmethod
    def serverStatus(self) -> TGServerStatus:
        """Gets the server's status."""

    @property
    @abc.abstractmethod
    def netListenersInfo(self) -> typing.List[TGNetListenerInfo]:
        """Gets the information for all of the server's net listeners."""

    @abc.abstractmethod
    def memoryInfo(self, type: MemoryType) -> TGMemoryInfo:
        """Gets the server's memory info for the given type."""

    @property
    @abc.abstractmethod
    def transactionsInfo(self) -> TGTransactionStatistics:
        """Gets the server's transactional statistics."""

    @property
    @abc.abstractmethod
    def queryInfo(self) -> TGProcessorSetStatistics:
        """Gets the server's query statistics."""

    @property
    @abc.abstractmethod
    def cacheInfo(self) -> TGCacheStatistics:
        """Gets the server's cache information."""

    @property
    @abc.abstractmethod
    def databaseInfo(self) -> TGDatabaseStatistics:
        """Gets the server's database information."""


class TGAdminDescription(abc.ABC):
    """
    Used for describing a system type, like a user, role, or to get more information about some of the other types.
    """

    @property
    @abc.abstractmethod
    def primary(self):
        """Gets the primary object described.

        :returns: The primarily-described system object.
        :rtype: tgdb.model.TGSystemObject
        """

    @property
    @abc.abstractmethod
    def attributes(self):
        """Gets the attributes of the primary object described.

        :returns: The attribute descriptors of the object. None if this object cannot have indices. Empty list if it\
        could, but there are no such indices.
        :rtype: typing.Optional[typing.List[tgdb.model.TGAttributeDescriptor]]
        """

    @property
    @abc.abstractmethod
    def indices(self):
        """Gets the indices of the primary object described.

        :returns: The indices on this type of object. None if this object cannot have indices. Empty list if it could,\
        but there are no such indices.
        :rtype: typing.Optional[typing.List[tgdb.model.TGIndex]]
        """

    @property
    @abc.abstractmethod
    def roles(self):
        """Gets the roles associated with the primary object. Should only be set when the primary is a user.

        :returns: The roles for this object.
        :rtype: typing.Optional[typing.List[tgdb.model.TGRole]]
        """

    @property
    @abc.abstractmethod
    def fromType(self):
        """Gets the from type of the primary object described.

        :returns: The from node's type, or None if there is no such type.
        :rtype: typing.Optional[tgdb.model.TGNodeType]
        """

    @property
    @abc.abstractmethod
    def toType(self):
        """Gets the to type of the primary object described.

        :returns: The to node's type, or None if there is no such type.
        :rtype: typing.Optional[tgdb.model.TGNodeType]
        """


class TGAdminConnection(tgconn.TGConnection, abc.ABC):
    """
    This class represents an administrative connection to the graph database server.

    Should only be used when modifying metadata or accessing server state information programmatically.
    """

    @abc.abstractmethod
    def getInfo(self) -> TGServerInfo:
        """
        Gets the server state information.
        """
    @abc.abstractmethod
    def getUsers(self) -> typing.List[TGUserInfo]:
        """
        Gets all users that the database server knows about.

        :return: List of all known users.
        """

    @abc.abstractmethod
    def getRoles(self) -> typing.List[tgmodel.TGRole]:
        """
        Gets all roles that the database server knows about.

        :return: List of all known roles.
        """

    @abc.abstractmethod
    def getConnections(self) -> typing.List[TGConnectionInfo]:
        """
        Gets all current live connections to the server.

        :return: List of all connections that are alive between the server and any clients, including this one.
        """

    @abc.abstractmethod
    def getStoredProcedures(self) -> typing.List[tgmodel.TGStoredProcedure]:
        """Gets the stored procedures."""

    @abc.abstractmethod
    def setSPDirectory(self, dir_path: str):
        """Sets the stored procedure's directory."""

    @abc.abstractmethod
    def refreshStoredProcedures(self):
        """Refreshes the server's stored procedures."""

    @abc.abstractmethod
    def stopServer(self):
        """
        Stops the server now.
        """

    @abc.abstractmethod
    def dumpServerStacktrace(self):
        """
        Programmatically dumps the server's stacktrace to the server's console and log files.

        Could be useful if one thread on the server is not responding.

        :return: Nothing.
        """

    @abc.abstractmethod
    def checkpointServer(self):
        """
        Allows the programmatic controller to checkpoint the server.

        :return: Nothing.
        """

    @abc.abstractmethod
    def killConnection(self, sessionid: typing.Optional[int] = None) -> int:
        """
        Kills a particular connection with the sessionid passed in.

        :param sessionid: The connection's sessionid of which to terminate. Default: None, which kills all connections.
        :return: The number of connections killed.
        """

    @abc.abstractmethod
    def getAttributeDescriptors(self) -> typing.List[tgmodel.TGAttributeDescriptor]:
        """
        Gets the attribute descriptors on this server.

        :return: The list of attribute descriptors.
        """

    @abc.abstractmethod
    def getIndices(self) -> typing.List[tgmodel.TGIndex]:
        """
        Gets the list of indices from the server.

        :return: The list of index information objects.
        """

    @abc.abstractmethod
    def getTypes(self) -> typing.List[tgmodel.TGEntityType]:
        """
        Gets a list of types from the server. Will only include entity types, but will include the default, built-in
        types too.

        :return: The types returned by the server.
        """

    @abc.abstractmethod
    def showTypes(self) -> None:
        """
        Shows a list of types from the server. Will only include entity types, but will include the default, built-in
        types too.

        :return: Nothing.
        """


    @abc.abstractmethod
    def describe(self, name: str) -> TGAdminDescription:
        """
        Gets description information for a particular type, attribute descriptor, or index.

        :return: The type's information, as well as any necessary supporting types (like the attribute descriptors that\
        are mentioned in a nodetype's attribute list).
        """

    @abc.abstractmethod
    def setServerLogLevel(self, logcomponent: typing.Union[tglog.TGLogComponent, typing.List[tglog.TGLogComponent]],
                          level: tglog.TGLevel):
        """
        Sets the server's log level for a particular component.

        :param logcomponent: The log component to set.
        :param level: The logging level at which to set the component.
        :return: Nothing.
        """

    @abc.abstractmethod
    def createNodetype(self, name: str, attrDescs: typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]],
                    pageSize: int = None,
                    pkeys: typing.Optional[typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]]] = None):
        """
        Creates a nodetype on the server.

        :param name: The name of the nodetype. Should not conflict with any other type, attribute descriptor, index, or\
        user already in the system.
        :param attrDescs: The attributes for this nodetype.
        :param pageSize: Optional, the pagesize for this nodes of this type. Might want to make it smaller for higher\
        efficiency or larger for when each nodetype will have many nodes and/or attributes are numerous and/or large.\
        Default: None, which implies using the system default of 512 bytes.
        :param pkeys: Optional, the primary keys of this nodetype. Must be a subset of the attribute descriptors\
        specified by the attrDescs parameter. Default: None, which implies using the system default of the globally\
        unique identifier.
        :return: Nothing. Will raise an error if not successful.
        """

    @abc.abstractmethod
    def createEdgetype(self, name: str, attrDescs: typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]],
                    dirType: typing.Optional[typing.Union[str, tgmodel.DirectionType]] = None,
                    fromType: typing.Optional[typing.Union[str, tgmodel.TGNodeType]] = None,
                    toType: typing.Optional[typing.Union[str, tgmodel.TGNodeType]] = None,
                    pkeys: typing.Optional[typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]]] = None):
        """
        Creates an edgetype on the server.

        :param name: The name of the edgetype. Should not conflict with any other type, attribute descriptor, index, or\
        user already in the system.
        :param attrDescs: The attributes for this edgetype.
        :param dirType: The direction that this edgetype has. If it is a string, will be checked against approved\
        strings before sending to the server. Default is None which implies that directionality is not important, and\
        for the API to pick.
        :param fromType: The nodetype of all nodes that edge comes from. Default implies that the type is not specified.
        :param toType: The nodetype of all nodes that edge comes from. Default implies that the type is not specified.
        :param pkeys: The primary key for this edge, must be a subset of the attrDescs parameter.
        :return: Nothing. Will raise an error if not successful.
        """

    @abc.abstractmethod
    def createIndex(self, name: str, attrDescs: typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]],
                        isUnique: bool = False,
                        types: typing.Optional[typing.List[typing.Union[tgmodel.TGEntityType, str]]] = None):
        """
        Creates an index on the server.
        :param name: Name of the index. Should not conflict with any other index, attribute descriptor, type, or\
        user already in the system.
        :param attrDescs: The attributes for this edgetype.
        :param isUnique: Whether every element on the index must be unique. Will prevent additions to the database if\
        those additions would make an index from being unique if the isUnique paramter is True.
        :param types: The types that this index is on. Default is any combination of the attribute descriptors\
        specified on any type.
        :return: Nothing. Will raise an error if not successful.
        """

    @abc.abstractmethod
    def createUser(self, name: str, password: str, roles: typing.List[str] = None):
        """
        Creates a user on the server.

        :param name: The name of the new user. Should not conflict with any other user, attribute descriptor, type, or\
        index already in the system.
        :param password: The new user's password.
        :param roles: The role(s) of the new user. If no role is specified, user will have whatever the database's\
            default role when accessing the database. Otherwise, the permissions and privileges of the user are a union\
            of the roles the user has, always including the database's default role.
        :return: Nothing. Will raise an error if not successful.
        """

    @abc.abstractmethod
    def createAttrDesc(self, name: str, ad_type: tgmodel.TGAttributeType, isArray: bool = False,
                           isEncrypted: bool = False, prec: int = None, scale: int = None):
        """
        Create a new attribute descriptor on the server.

        :param name: Name of the attribute descriptor. Should not conflict with any other attribute descriptor, index,\
        type, or user already in the system.
        :param ad_type: The type of the attribute descriptor.
        :param isArray: Whether the new attribute descriptor is an array (i.e. can allow multiple values on one node or\
        edge of this attribute descriptor)
        :param isEncrypted: Whether the new attribute descriptor should be encrypted.
        :param prec: The precision of the attribute descriptor. Should only be set when the ad_type parameter is set to\
        Number.
        :param scale: The scale of the attribute descriptor. Should only be set when the ad_type parameter is set to\
        Number. Should not be greater than the prec paramter.
        :return:
        """

    @abc.abstractmethod
    def fullImport(self, dir_path: str = "import") -> typing.List[TGImportDescriptor]:
        """
        Imports all of the files found in the directory specified to the server.

        :param dir_path: The directory to get all of the files from.
        :return: A list of import descriptors that describe what the types are and how many were imported.
        """

    @abc.abstractmethod
    def createRole(self, rolename: str, permissions: typing.List[typing.Tuple[typing.Union[str, tgmodel.TGNodeType,
                  tgmodel.TGEdgeType, tgmodel.TGAttributeDescriptor], typing.Union[tgmodel.TGPermissionType, str]]],
                   privileges: typing.Union[typing.Iterable[tgmodel.TGPrivilege], tgmodel.TGPrivilege] = tuple()):
        """Creates a role with the permissions and privileges passed in."""

    @abc.abstractmethod
    def grantPermissions(self, rolename: str, permissions: typing.List[typing.Tuple[
                             typing.Union[str, tgmodel.TGNodeType, tgmodel.TGEdgeType, tgmodel.TGAttributeDescriptor],
                             typing.Union[tgmodel.TGPermissionType, str]]]):
        """Grants to a role the permissions passed in."""

    @abc.abstractmethod
    def grantPrivileges(self, rolename: str, privileges: typing.Union[typing.Iterable[tgmodel.TGPrivilege],
                                                                      tgmodel.TGPrivilege] = tuple()):
        """Grants to a role the privileges passed in."""

    @abc.abstractmethod
    def revokePermissions(self, rolename: str, permissions: typing.List[typing.Tuple[
                             typing.Union[str, tgmodel.TGNodeType, tgmodel.TGEdgeType, tgmodel.TGAttributeDescriptor],
                             typing.Union[tgmodel.TGPermissionType, str]]]):
        """Revokes from a role the permissions passed in."""

    @abc.abstractmethod
    def revokePrivileges(self, rolename: str, privileges: typing.Union[typing.Iterable[tgmodel.TGPrivilege],
                                                                       tgmodel.TGPrivilege] = tuple()):
        """Revokes from a role the privileges passed in."""

    @abc.abstractmethod
    def clearRole(self, rolename: str):
        """Clears the permissions from the role from the server."""

    @abc.abstractmethod
    def grantUser(self, username: str, roles: typing.List[str]):
        """Grants to a user the roles specified in the list of roles."""

    @abc.abstractmethod
    def revokeUser(self, username: str, roles: typing.List[str]):
        """Revokes from a user the roles specified in the list of roles."""
