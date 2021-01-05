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
 *  File name :channel.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: log.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates the log package of java
 """

import enum
import abc
import re
import logging
import typing
import inspect
import traceback
import tgdb.model as tgmodel


__all__ = ['TGLevel', 'TGLogComponent', 'TGPyLogComponent', 'TGPyLogCompManager', 'TGLogManager', 'TGLogger', 'gLogger']


class TGLevel(enum.Enum):
    """ The level of debugging, including at what level a log should be recorded.

    Renamed to more closely resemble the TGDB scheme. From 12/01/2019, transition all references to TGL* from current.
    """
    TGLogLevelInvalid = -1
    TGLConsole = -2
    TGLFatal = 0
    TGLError = 1
    TGLWarn = 2
    TGLInfo = 3
    TGLUser = 4
    TGLDebug = 5
    TGLDebugFine = 6
    TGLDebugFiner = 7
    TGLMaxLogLevel = 8
    NotSet = -1
    DebugWire = 7
    Debug = 5
    Info = 3
    Warning = 2
    Error = 1
    Fatal = 0


all_tglc = {}


class TGLogComponent(enum.Enum):
    """
    This enumeration is for any components of the server that a user may want to keep track of.

    They are broken down by larger components, and individual portions can be invoked for fine-grained control, or an
    aggregate one may be preferential. This class allows both of those possibilities.

    Only intended when setting the loglevel on the server or the administrative client. The bits property and the
    aggregate_bits method are not intended for use in client code.
    """

    TGLC_INVALID = (-1, set(), "__UNUSED__")

    # Common components.
    TGLC_COMMON_COREMEMORY = (0, set(), "common.core.memory")
    TGLC_COMMON_CORECOLLECTIONS = (1, set(), "common.core.collections")
    TGLC_COMMON_COREPLATFORM = (2, set(), "common.core.platform")
    TGLC_COMMON_CORESTRING = (3, set(), "common.core.string")
    TGLC_COMMON_UTILS = (4, set(), "common.utils")
    TGLC_COMMON_GRAPH = (5, set(), "common.graph")
    TGLC_COMMON_MODEL = (6, set(), "common.model")
    TGLC_COMMON_NET = (7, set(), "common.net")
    TGLC_COMMON_PDU = (8, set(), "common.pdu")
    TGLC_COMMON_SEC = (9, set(), "common.sec")
    TGLC_COMMON_FILES = (10, set(), "common.files")
    TGLC_COMMON_RESV2 = (11, set(), "__UNUSED__")

    # Server only components
    TGLC_SERVER_CDMP = (12, set(), "server.cdmp")
    TGLC_SERVER_DB = (13, set(), "server.db")
    TGLC_SERVER_EXPIMP = (14, set(), "server.expimp")
    TGLC_SERVER_INDEX = (15, set(), "server.index")
    TGLC_SERVER_INDEXBTREE = (16, set(), "server.index.btree")
    TGLC_SERVER_INDEXISAM = (17, set(), "server.index.isam")
    TGLC_SERVER_QUERY = (18, set(), "server.query")
    TGLC_SERVER_QUERY_RESV1 = (19, set(), "__UNUSED__")
    TGLC_SERVER_QUERY_RESV2 = (20, set(), "__UNUSED__")
    TGLC_SERVER_TXN = (21, set(), "server.txn")
    TGLC_SERVER_TXNLOG = (22, set(), "server.txn.log")
    TGLC_SERVER_TXNWRITER = (23, set(), "server.txn.writer")
    TGLC_SERVER_STORAGE = (24, set(), "server.storage")
    TGLC_SERVER_STORAGEPAGEMANAGER = (25, set(), "server.storage.pagemanager")
    TGLC_SERVER_GRAPH = (26, set(), "server.graph")
    TGLC_SERVER_MAIN = (27, set(), "server.main")
    TGLC_SERVER_RESV2 = (28, set(), "__UNUSED__")
    TGLC_SERVER_RESV3 = (29, set(), "__UNUSED__")
    TGLC_SERVER_RESV4 = (30, set(), "__UNUSED__")

    # Security components
    TGLC_SECURITY_DATA = (31, set(), "security.data")
    TGLC_SECURITY_NET = (32, set(), "security.net")
    TGLC_SECURITY_RESV1 = (33, set(), "__UNUSED__")
    TGLC_SECURITY_RESV2 = (34, set(), "__UNUSED__")

    # Administration components
    TGLC_ADMIN_LANG = (35, set(), "admin.lang")
    TGLC_ADMIN_CMD = (36, set(), "admin.cmd")
    TGLC_ADMIN_MAIN = (37, set(), "admin.main")
    TGLC_ADMIN_AST = (38, set(), "admin.ast")
    TGLC_ADMIN_GREMLIN = (39, set(), "admin.gremlin")

    # CUDA components
    TGLC_CUDA_GRAPHMGR = (40, set(), "cuda.graphmgr")
    TGLC_CUDA_KERNELEXECUTIVE = (41, set(), "cuda.kernelexecutive")
    TGLC_CUDA_RESV1 = (42, set(), "__UNUSED__")

    # Log all components
    TGLC_LOG_GLOBAL = (-1, set(), "*")

    # Log core components
    TGLC_LOG_COREALL = (-1, {'TGLC_COMMON_COREMEMORY', 'TGLC_COMMON_CORECOLLECTIONS', 'TGLC_COMMON_COREPLATFORM',
                             'TGLC_COMMON_CORESTRING'}, "core.all")

    # Log a particular set of components or features
    TGLC_LOG_GRAPHALL = (-1, {'TGLC_COMMON_GRAPH', 'TGLC_SERVER_GRAPH'}, "graph.all")
    TGLC_LOG_NET = (-1, {'TGLC_COMMON_NET'}, "__UNUSED__")
    TGLC_LOG_PDUALL = (-1, {'TGLC_COMMON_PDU', 'TGLC_SERVER_CDMP'}, "pdu.all")
    TGLC_LOG_SECALL = (-1, {'TGLC_COMMON_SEC', 'TGLC_SECURITY_DATA', 'TGLC_SECURITY_NET'}, "sec.all")
    TGLC_LOG_CUDAALL = (-1, {'TGLC_LOG_GRAPHALL', 'TGLC_CUDA_GRAPHMGR', 'TGLC_CUDA_KERNELEXECUTIVE'}, "cuda.all")
    TGLC_LOG_TXNALL = (-1, {'TGLC_SERVER_TXN', 'TGLC_SERVER_TXNLOG', 'TGLC_SERVER_TXNWRITER'}, "txn.all")
    TGLC_LOG_STORAGEALL = (-1, {'TGLC_SERVER_STORAGE', 'TGLC_SERVER_STORAGEPAGEMANAGER'}, "storage.all")
    TGLC_LOG_PAGEMANAGER = (-1, {'TGLC_SERVER_STORAGEPAGEMANAGER'}, "__UNUSED__")
    TGLC_LOG_ADMINALL = (-1, {'TGLC_ADMIN_LANG', 'TGLC_ADMIN_CMD', 'TGLC_ADMIN_MAIN', 'TGLC_ADMIN_AST',
                              'TGLC_ADMIN_GREMLIN'}, "admin.all")
    TGLC_LOG_MAIN = (-1, {'TGLC_ADMIN_MAIN', 'TGLC_SERVER_MAIN'}, "main.all")

    def __init__(self, num_shift: int, sub_tglc_comps: typing.Union[typing.List, typing.Set], longName: str):
        self.__longname = longName
        self.__num_shift = num_shift
        self.__sub_tglc_comps = sub_tglc_comps

        all_tglc[self.name] = self

    def __hash__(self):
        return hash((self.__num_shift, tuple(self.__sub_tglc_comps)))

    def __eq__(self, other):
        if isinstance(other, TGLogComponent):
            return self.__num_shift == other.__num_shift and self.__sub_tglc_comps == other.__sub_tglc_comps
        else:
            return NotImplemented

    @property
    def longname(self):
        return self.__longname

    def _to_bits(self, visited: typing.Set) -> int:
        if self in visited:
            raise Exception('Unexpected loop in TGLogComponent')
        ret = 0
        if self.__num_shift == -1:
            if len(self.__sub_tglc_comps) > 0:
                visited.add(self)
                for tglc in self.__sub_tglc_comps:
                    ret |= all_tglc[tglc]._to_bits(visited)
                visited.remove(self)
            else:
                ret = 0xFFFFFFFFFFFFFFFF
        else:
            ret = 1 << self.__num_shift
        return ret

    @property
    def bits(self) -> int:
        """
        Gets the bit representation of this TGLogComponent so that it may be sent to the server.

        Not intended for client use.

        :return: A bit string (packed as a Python int) that contains all of the 'on' components.
        """
        return self._to_bits(set())

    @classmethod
    def fromName(cls, longname: str):
        for tgl in TGLogComponent:
            if tgl.longname == longname:
                return tgl
        return cls.TGLC_INVALID

    @classmethod
    def aggregate_bits(cls, arr: typing.Iterable) -> int:
        """
        Gets the bit representation of the TGLogComponents in the iterable so that it may be sent to the server.

        Not intended for client use.

        :return: A bit string (packed as a Python int) that contains all of the 'on' components.
        """
        ret = 0
        for elem in arr:
            if not isinstance(elem, cls):
                raise NotImplementedError("Not allowed to call aggregate_bits with anything but TGLogComponent instance"
                                          "s.")
            ret |= elem.bits
        return ret


all_tgpylc = {}


class TGPyLogComponent(enum.Enum):
    """This class represents logging inside of the Python API and the TGDB Admin Console. Only intended to set the
    loglevel on specific components when running into an error."""

    ADMIN = ("admin", {"tgdb.admin", "tgdb.impl.adminimpl"}, set(), 0)
    BULKIO = ("bulkio", {"tgdb.bulkio", "tgdb.impl.bulkioimpl"}, set(), 1)
    CHANNEL = ("channel", {"tgdb.channel", "tgdb.impl.channelimpl"}, set(), 2)
    CONNECTION = ("connection", {"tgdb.connection", "tgdb.impl.connectionimpl"}, set(), 3)
    PDU = ("pdu", {"tgdb.pdu", "tgdb.impl.pduimpl"}, set(), 4)
    QUERY = ("query", {"tgdb.query", "tgdb.impl.queryimpl"}, set(), 5)
    MODEL = ("model", {"tgdb.model", "tgdb.impl.attrimpl", "tgdb.impl.entityimpl", "tgdb.impl.gmdimpl"}, set(), 6)
    EXCEPTION = ("exception", {"tgdb.exception", }, set(), 7)
    LOG = ("log", {"tgdb.log", }, set(), 8)
    STOREDPROC = ("storedproc", {"tgdb.storedproc", }, set(), 9)
    UTILS = ("utils", {"tgdb.utils", }, set(), 10)
    VERSION = ("version", {"tgdb.version", }, set(), 11)
    LICENSE = ("license", {"tgdb.license.*", }, set(), 16)
    ALL = ("*", {"tgdb.*", }, {'admin', 'bulkio', 'channel', 'connection', 'pdu', 'query', 'model', 'exception', 'log',
           'storedproc', 'utils', 'version', 'license'}, -1)

    def __init__(self, displayName: str, matchingModules: typing.Iterable[str], subversions: typing.Iterable[str],
                 bitNum: int):
        self.__matchingModules = set()
        for matchMod in matchingModules:
            self.__matchingModules.add(self.__convertMatchingModule(matchMod))
        self.__subversions = subversions
        if bitNum >= 0:
            self.__bitString = 1 << bitNum
            self.__getBitString = 1 << bitNum
        else:
            self.__getBitString = 0xFFFFFFFF
            self.__bitString = 0
        for name in subversions:
            self.__bitString |= all_tgpylc[name].setBitString
        self.__displayName = displayName
        all_tgpylc[self.__displayName] = self

    @property
    def setBitString(self):
        return self.__bitString

    @property
    def getBitString(self):
        return self.__getBitString

    @property
    def displayName(self) -> str:
        """The fancy display name used in the admin console or when lazily setting the level of a particular component
        using the API."""
        return self.__displayName

    def __convertMatchingModule(self, inp) -> str:
        """Not for use outside of this class."""
        import os, re
        inp = inp.replace('.', '(?:/|%s)' % re.escape(os.sep))
        inp = inp.replace('*', '.*')
        return '^.*' + inp + '[.]pyc?$'

    @classmethod
    def getBestMatch(cls, filename: str):
        """Gets the best match for the corresponding filename, not for use by clients."""
        import re
        lc_match = cls.ALL
        for comp in cls:
            for mm in comp.__matchingModules:
                if re.match(mm, filename) and len(comp.displayName) > len(lc_match.displayName):
                    lc_match = comp
                    break
        return lc_match

    @classmethod
    def getLogComponentByName(cls, compName: str):
        if compName in all_tgpylc:
            return all_tgpylc[compName]
        return None


class TGPyLogCompManager:
    """Only for internal use by the Python API/Admin Console"""

    def __init__(self, level: TGLevel):
        self.__lvl = level
        self.__loglevelMap = {}
        for lvl in TGLevel:
            if lvl.value <= level.value:
                self.__loglevelMap[lvl.value] = 0xFFFFFFFF
            else:
                self.__loglevelMap[lvl.value] = 0

    def setLC(self, lc: TGPyLogComponent, level: TGLevel):
        for lvl in TGLevel:
            if lvl.value <= level.value:
                self.__loglevelMap[lvl.value] |= lc.setBitString
            else:
                self.__loglevelMap[lvl.value] &= lc.setBitString

    def enabledForComponent(self, lc: TGPyLogComponent, level: TGLevel) -> bool:
        return self.__loglevelMap[level.value] & lc.getBitString > 0


class TGLogger(abc.ABC):
    """Represents a single logger for the client to use to record information."""

    @abc.abstractmethod
    def log(self, level: TGLevel, fmt: str, *args, **kwargs):
        """Logs some piece of information.

        :param level: the level at which this logger should activate (for example, if the logging level is set to info,
        then all info, warning, error, and fatal logs will be recorded, but debug logs will not be.
        :param fmt: the format string, similar to the format string for C's printf.
        :param args: list of objects to pass into the formatted string.
        :param kwargs: dictionary of keyword arguments to pass into the formatted string."""
        pass

    @abc.abstractmethod
    def logEntity(self, level: TGLevel, ent: tgmodel.TGEntity):
        """Logs an entity at the particular log level."""
        pass

    @abc.abstractmethod
    def logException(self, level, msg, e: BaseException):
        """Logs an exception with a message.
        :param level:
        """
        pass

    @abc.abstractmethod
    def isEnabled(self, level: TGLevel) -> bool:
        """Checks if a particular log level is enabled on this logger."""
        pass

    @property
    @abc.abstractmethod
    def level(self) -> TGLevel:
        """Gets this logger's general log level."""
        pass

    @level.setter
    @abc.abstractmethod
    def level(self, lvl: TGLevel):
        """Sets this logger's general log level."""
        pass

    @abc.abstractmethod
    def setLevel(self, lvl: TGLevel, lc: TGPyLogComponent):
        """Sets this logger's level for the specific log component specified"""

    @classmethod
    def getBestLogCompFromStackTrace(cls, trace: typing.Union[traceback.StackSummary, typing.List[inspect.FrameInfo]]):
        """Only for internal use by the Python API/Admin Console"""
        import os
        if isinstance(trace, traceback.StackSummary):
            tmp_trace = []                          # Need to reverse the stack trace
            for i in range(len(trace)-1, -1, -1):
                tmp_trace.append(trace[i])
            trace = tmp_trace
        for frame in trace:
            if re.search("(?:/|{0})tgdb(?:/|{0})".format(re.escape(os.sep)), frame.filename):
                return TGPyLogComponent.getBestMatch(frame.filename)
        return TGPyLogComponent.ALL


class SimplePythonLogger(TGLogger):
    """Default logging implementation.

    Not for client use.
    """

    def __init__(self):
        logging.addLevelName(self._getLevel(TGLevel.DebugWire), 'DebugWire')
        self.__logger__ = logging.getLogger('tgdb')
        self.__logger__.setLevel(self._getLevel(TGLevel.TGLDebugFiner)-1)
        self.__pylclvl = TGPyLogCompManager(TGLevel.Info)

    def _getLevel(self, level: TGLevel) -> int:
        return 10 - level.value                 # Need to invert the levels show that the logs show up properly.

    def log(self, level: TGLevel, fmt: str, *args, **kwargs):
        if self.__pylclvl.enabledForComponent(TGLogger.getBestLogCompFromStackTrace(inspect.stack()), level):
            self.__logger__.log(self._getLevel(level), fmt, *args, **kwargs)

    def logException(self, level: TGLevel, msg, e: BaseException):
        if self.__pylclvl.enabledForComponent(TGLogger.getBestLogCompFromStackTrace(
                                                        traceback.extract_tb(e.__traceback__)), level):
            self.__logger__.exception(msg)

    def logEntity(self, level: TGLevel, ent: tgmodel.TGEntity):
        if not self.__pylclvl.enabledForComponent(TGLogger.getBestLogCompFromStackTrace(inspect.stack()), level):
            return
        self.log(level, "Kind: %s Id: %d", str(ent.entityKind), ent.virtualId)
        attr: tgmodel.TGAttribute
        node: tgmodel.TGNode
        edge: tgmodel.TGEdge
        for attr in ent.attributes:
            self.log(level, " Attr[%s]: %s", attr.descriptor.name, attr.value)
        if ent.entityKind is tgmodel.TGEntityKind.Node:
            node = ent
            for edge in node.edges:
                self.log(level, "   Edge: %d", edge.virtualId)
        elif ent.entityKind is tgmodel.TGEntityKind.Edge:
            edge = ent
            for node in edge.vertices:
                self.log(level, "   Node: %d", node.virtualId)

    def isEnabled(self, level: TGLevel) -> bool:
        return self.__pylclvl.enabledForComponent(TGLogger.getBestLogCompFromStackTrace(inspect.stack()), level)

    @property
    def level(self) -> TGLevel:
        minLvl = TGLevel.TGLConsole
        for lvl in TGLevel:
            if self.__pylclvl.enabledForComponent(TGPyLogComponent.ALL, lvl) and lvl.value > minLvl.value:
                minLvl = lvl
        return minLvl

    @level.setter
    def level(self, lvl: TGLevel):
        self.__pylclvl.setLC(TGPyLogComponent.ALL, lvl)
        self.__logger__.setLevel(self._getLevel(lvl))
        return

    def setLevel(self, lvl: TGLevel, lc: TGPyLogComponent):
        self.__pylclvl.setLC(lc, lvl)


class TGLogManager:
    """Class intended to initialize and start the logger."""

    __logger__: TGLogger = None

    @classmethod
    def getLogger(cls) -> TGLogger:
        """Gets the logger (creating one if one does not already exist)."""
        if cls.__logger__ is None:
            cls.__logger__ = SimplePythonLogger()
        return cls.__logger__

    @classmethod
    def setLogger(cls, logger: TGLogger, setGlobal: bool = False) -> TGLogger:
        """Sets the logger, returning the old one."""
        global gLogger
        oldlogger = cls.__logger__
        cls.__logger__ = logger
        if setGlobal:
            gLogger = logger
        return oldlogger


"""Default 'global' logger."""
gLogger = TGLogManager.getLogger()
