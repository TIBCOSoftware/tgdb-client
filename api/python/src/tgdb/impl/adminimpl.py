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
 *  File name : adminimpl.py
 *  Created on: 11/19/2019
 *  Created by: derek
 *  SVN Id: $Id$
 *
"""

import datetime
import time
import enum
import typing

import tgdb.impl.connectionimpl as connimpl
import tgdb.connection as tgconn
import tgdb.admin as tgadmin
from tgdb import log as tglog
from tgdb.admin import MemoryType, TGMemoryInfo, ServerState, TGDatabaseStatistics, TGCacheStatistics, \
    TGProcessorSetStatistics, TGNetListenerInfo, TGServerStatus, TGCacheStat
import tgdb.version as tgvers
import tgdb.pdu as tgpdu
import tgdb.impl.pduimpl as pduimpl
import tgdb.model as tgmodel
import tgdb.exception as tgexception
import tgdb.impl.attrimpl as attrimpl
import tgdb.impl.gmdimpl as gmdimpl
from tgdb.pdu import TGInputStream, TGOutputStream


class ImportDescImpl(tgadmin.TGImportDescriptor):
    def __init__(self, typename: str, isnode: bool, numinstances: int):
        self.__typename = typename
        self.__isnode = isnode
        self.__numinstances = numinstances

    @property
    def typename(self) -> str:
        return self.__typename

    @property
    def isnode(self) -> bool:
        return self.__isnode

    @property
    def numinstances(self) -> int:
        return self.__numinstances

    def _inc(self, val: int):
        self.__numinstances += val


class TGAdminCommand(enum.Enum):
    Invalid = -1
    CreateUser = 1
    CreateRole = 2
    CreateAttrDesc = 3
    CreateIndex = 4
    CreateNodeType = 5
    CreateEdgeType = 6
    ShowUsers = 7
    ShowRoles = 8
    ShowAttrDescs = 9
    ShowIndices = 10
    ShowTypes = 11
    ShowInfo = 12
    ShowConnections = 13
    ShowSP = 14
    Describe = 15
    SetLogLevel = 16
    SetSPDir = 17
    UpdateUser = 18
    UpdateRole = 19
    RefreshSPDir = 20
    StopServer = 21
    CheckpointServer = 22
    DisconnectClient = 23
    KillConnection = 24

    def __init__(self, val: int):
        self.__value = val

    @property
    def val(self):
        return self.__value

    @classmethod
    def fromId(cls, id: int):
        for ac in TGAdminCommand:
            if ac.__value == id:
                return ac
        return TGAdminCommand.Invalid


class CacheStatImpl(tgadmin.TGCacheStat):
    def __init__(self):
        self._cacheMaxEntries: int = -1
        self._cacheEntries: int = -1
        self._cacheHits: int  = -1
        self._cacheMisses: int = -1
        self._cacheMaxMemory: int = -1

    @property
    def cacheMaxEntries(self) -> int:
        return self._cacheMaxEntries

    @property
    def cacheNumEntries(self) -> int:
        return self._cacheEntries

    @property
    def cacheHits(self) -> int:
        return self._cacheHits

    @property
    def cacheMisses(self) -> int:
        return self._cacheMisses

    @property
    def cacheMaxMemory(self) -> int:
        return self._cacheMaxMemory
    
    def __str__(self):
        return "CacheStatImpl [cacheMaxEntries=" + str(self._cacheMaxEntries) + ", cacheEntries=" +\
                str(self._cacheEntries) + ", cacheHits=" + str(self._cacheHits) + ", cacheMisses="\
                + str(self._cacheMisses) + ", cacheMaxMemory=" + str(self._cacheMaxMemory) + "]"

    @classmethod
    def readExternal(cls, instream: tgpdu.TGInputStream):
        ret = CacheStatImpl()
        ret._cacheMaxEntries = instream.readInt()
        ret._cacheEntries = instream.readInt()
        ret._cacheHits = instream.readLong()
        ret._cacheMisses = instream.readLong()
        ret._cacheMaxMemory = instream.readLong()
        return ret


class EntityCacheImpl(tgadmin.TGEntityCache, CacheStatImpl):

    def __init__(self):
        super()
        self._averageEntitySize: int = -1
        self._maxEntitySize: int = -1

    @property
    def averageEntitySize(self) -> int:
        return self._averageEntitySize

    @property
    def maxEntitySize(self) -> int:
        return self._maxEntitySize

    def __str__(self):
        return "EntityCacheImpl [cacheMaxEntries=" + str(self._cacheMaxEntries) + ", cacheEntries=" +\
                str(self._cacheEntries) + ", cacheHits=" + str(self._cacheHits) + ", cacheMisses="\
                + str(self._cacheMisses) + ", cacheMaxMemory=" + str(self._cacheMaxMemory) + ", averageEntitySize=" +\
                str(self._averageEntitySize) + ", maxEntitySize=" + str(self._maxEntitySize) + "]"

    @classmethod
    def readExternal(cls, instream: tgpdu.TGInputStream):
        ret = EntityCacheImpl()
        ret._cacheMaxEntries = instream.readInt()
        ret._cacheEntries = instream.readInt()
        ret._cacheHits = instream.readLong()
        ret._cacheMisses = instream.readLong()
        ret._cacheMaxMemory = instream.readLong()
        ret._averageEntitySize = instream.readInt()
        ret._maxEntitySize = instream.readInt()
        return ret

class CacheStatisticsImpl(tgadmin.TGCacheStatistics):

    @property
    def dataCache(self) -> tgadmin.TGCacheStat:
        return self._dataCacheStat

    @property
    def indexCache(self) -> TGCacheStat:
        return self._indexCacheStat

    @property
    def sharedCache(self) -> TGCacheStat:
        return self._sharedCacheStat

    @property
    def entityCache(self) -> tgadmin.TGEntityCache:
        return self._entityCacheStat

    def __init__(self, dataCacheStat: CacheStatImpl, indexCacheStat: CacheStatImpl,
                 sharedCacheStat: CacheStatImpl, entityCacheStat: EntityCacheImpl):
        self._dataCacheStat = dataCacheStat
        self._indexCacheStat = indexCacheStat
        self._sharedCacheStat = sharedCacheStat
        self._entityCacheStat = entityCacheStat

    def __str__(self):
        return "CacheStatisticsImpl [dataCache=" + str(self._dataCacheStat) + ", indexCache=" +\
               str(self._indexCacheStat) + ", sharedCache=" + str(self._sharedCacheStat) + ", entityCache=" +\
               str(self._entityCacheStat) + "]"

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGCacheStatistics:
        data = CacheStatImpl.readExternal(instream)
        index = CacheStatImpl.readExternal(instream)
        shared = CacheStatImpl.readExternal(instream)
        entity = EntityCacheImpl.readExternal(instream)
        return cls(data, index, shared, entity)


class ConnectionInfoImpl(tgadmin.TGConnectionInfo):

    @property
    def listenerName(self) -> str:
        return self._listenerName

    @property
    def clientID(self) -> str:
        return self._clientID

    @property
    def sessionID(self) -> int:
        return self._sessionID

    @property
    def userName(self) -> str:
        return self._username

    @property
    def remoteAddress(self) -> str:
        return self._remoteAddress

    @property
    def createdTimeInSeconds(self) -> float:
        return self._createdTimeInSeconds

    def __init__(self, lname: str, cid: str, sid: int, uname: str, raddr: str, createdTime: float):
        self._listenerName: str = lname
        self._clientID: str = cid
        self._sessionID: int = sid
        self._username: str = uname
        self._remoteAddress: str = raddr
        self._createdTimeInSeconds: float = createdTime

    def __str__(self):
        return "ConnectionInfoImpl [listnerName=" + str(self._listenerName) + ", clientID=" + str(self._clientID) +\
               ", sessionID=" + str(self._sessionID) + ", userName=" + str(self._username) + ", remoteAddress=" +\
               str(self._remoteAddress) + ", createdTimeInSeconds=" + str(self._createdTimeInSeconds) + "]"

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGConnectionInfo:
        ln = instream.readUTF()
        cid = instream.readUTF()
        sid = instream.readUnsignedLong()
        un = instream.readUTF()
        raddr = instream.readUTF()
        ctis = (datetime.datetime.now().timestamp()) - (instream.readLong()/1000000000)
        return cls(ln, cid, sid, un, raddr, ctis)


class DatabaseStatisticsImpl(tgadmin.TGDatabaseStatistics):

    @property
    def dbSize(self) -> int:
        return self._dbSize

    @property
    def dbPath(self) -> str:
        return self._path

    @property
    def numDataSegments(self) -> int:
        return self._numDataSegments

    @property
    def dataSize(self) -> int:
        return self._dataSize

    @property
    def dataUsed(self) -> int:
        return self._dataUsed

    @property
    def dataFree(self) -> int:
        return self._dataFree

    @property
    def dataBlockSize(self) -> int:
        return self._dataBlockSize

    @property
    def numIndexSegments(self) -> int:
        return self._numIndexSegments

    @property
    def indexSize(self) -> int:
        return self._indexSize

    @property
    def usedIndexSize(self) -> int:
        return self._indexUsed

    @property
    def freeIndexSize(self) -> int:
        return self._indexFree

    @property
    def blockSize(self) -> int:
        return self._blockSize

    def __init__(self, dbSize: int, dbPath: str, nds: int, dataSize: int, dataUsed: int, dataFree: int,
                 dataBlockSize: int, numIdxSegs: int, idxSize: int, idxUsed: int, idxFree: int, idxBlockSize: int):
        self._dbSize: int = dbSize
        self._numDataSegments: int = nds
        self._dataSize: int = dataSize
        self._dataUsed: int = dataUsed
        self._dataFree: int = dataFree
        self._dataBlockSize: int = dataBlockSize

        self._numIndexSegments: int = numIdxSegs
        self._indexSize: int = idxSize
        self._indexUsed: int = idxUsed
        self._indexFree: int = idxFree
        self._blockSize: int = idxBlockSize
        self._path: str = dbPath

    def __str__(self):
        return "DatabaseStatisticsImpl [dbSize=" + str(self._dbSize) + ", dbPath=" + str(self._path) +\
               ", numDataSegments=" + str(self._numDataSegments) + ", dataSize=" + str(self._dataSize) + ", dataUsed="\
               + str(self._dataUsed) + ", dataFree=" + str(self._dataFree) + ", dataBlockSize=" +\
               str(self._dataBlockSize) + ", numIndexSegments=" + str(self._numIndexSegments) + ", indexSize=" +\
               str(self._indexSize) + ", indexUsed=" + str(self._indexUsed) + ", indexFree=" + str(self._indexFree) +\
               ", blockSize=" + str(self._blockSize) + "]"

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGDatabaseStatistics:
        dbs = instream.readLong()
        dbPath = instream.readUTF()
        nds = instream.readInt()
        ds = instream.readLong()
        du = instream.readLong()
        df = instream.readLong()
        dbls = instream.readInt()

        idxsg = instream.readInt()
        idxs = instream.readLong()
        idxu = instream.readLong()
        idxf = instream.readLong()
        idxbs = instream.readInt()
        return cls(dbs, dbPath, nds, ds, du, df, dbls, idxsg, idxs, idxu, idxf, idxbs)


class MemoryInfoImpl(tgadmin.TGMemoryInfo):

    @property
    def usedMemorySize(self) -> int:
        return self._usedMemory

    @property
    def freeMemorySize(self) -> int:
        return self._freeMemory

    @property
    def maxMemorySize(self) -> int:
        return self._maxMemory

    @property
    def sharedMemoryFileLocation(self) -> str:
        return self._sharedMemoryFileLocation

    def __init__(self, usedMem: int, freeMem: int, maxMem: int, sharedMemFileLoc: str):
        self._usedMemory: int = usedMem
        self._freeMemory: int = freeMem
        self._maxMemory: int = maxMem
        self._sharedMemoryFileLocation: str = sharedMemFileLoc

    def __str__(self):
        return "MemoryInfoImpl [usedMemory=" + str(self._usedMemory) + ", freeMemory=" + str(self._freeMemory) +\
                ", maxMemory=" + str(self._maxMemory) + ", sharedMemoryFileLocation=" + \
                str(self._sharedMemoryFileLocation) + "]"
    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream, hasFileLoc: bool = True) -> tgadmin.TGMemoryInfo:
        if hasFileLoc:
            sfl = instream.readUTF()
        else:
            sfl = None
        fm = instream.readLong()
        _ = instream.readInt()      # Memory usage percentage. Not used
        mm = instream.readLong()
        um = mm - fm
        return cls(um, fm, mm, sfl)


class NetListenerInfoImpl(tgadmin.TGNetListenerInfo):

    @property
    def listenerName(self) -> str:
        return self._listenerName

    @property
    def currentConnections(self) -> int:
        return self._currentConnections

    @property
    def maxConnections(self) -> int:
        return self._maxConnections

    @property
    def portNumber(self) -> str:
        return self._portNumber

    def __init__(self, lname: str, numConns: int, maxConns: int, portNum: str):
        self._listenerName: str = lname
        self._currentConnections: int = numConns
        self._maxConnections: int = maxConns
        self._portNumber: str = portNum

    def __str__(self):
        return "NetListenerInfoImpl [listenerName=" + self._listenerName + ", currentConnections=" +\
                str(self._currentConnections) + ", maxConnections=" + str(self._maxConnections) + ", portNumber=" +\
                self._portNumber + "]"

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGNetListenerInfo:
        lname = instream.readBytes().decode('utf-8')
        nconns = instream.readInt()
        mconns = instream.readInt()
        pnum = instream.readBytes().decode('utf-8')
        return cls(lname, nconns, mconns, pnum)


class ServerMemoryInfoImpl(tgadmin.TGServerMemoryInfo):

    def getMemoryInfo(self, type: MemoryType) -> TGMemoryInfo:
        if type == tgadmin.MemoryType.PROCESS:
            return self._processMemory
        elif type == tgadmin.MemoryType.SHARED:
            return self._sharedMemory
        else:
            return None

    def __init__(self, procMem: tgadmin.TGMemoryInfo, sharedMem: tgadmin.TGMemoryInfo):
        self._processMemory: MemoryInfoImpl = procMem
        self._sharedMemory: MemoryInfoImpl = sharedMem

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGServerMemoryInfo:
        pmem = MemoryInfoImpl.readFromStream(instream, False)
        smem = MemoryInfoImpl.readFromStream(instream, True)

        return cls(pmem, smem)


class ServerStatusImpl(tgadmin.TGServerStatus):

    @property
    def name(self) -> str:
        return self._name

    @property
    def version(self) -> tgvers.TGVersion:
        return self._vers

    @property
    def state(self) -> ServerState:
        return self._status

    @property
    def processID(self) -> str:
        return self._pid

    @property
    def uptime(self) -> datetime.timedelta:
        return self._uptime

    @property
    def serverPath(self) -> str:
        return self._path

    def __str__(self):
        return "ServerStatusImpl [name=" + self._name + ", version=" + str(self._vers) + ", status=" +\
               str(self._status) + ", processId=" + self._pid + ", uptime=" + str(self._uptime) + "]"

    def __init__(self, name: str, vers: tgvers.TGVersion, status: tgadmin.ServerState, pid: str,
                 uptime: datetime.timedelta, serverPath: str):

        self._name = name
        self._vers = vers
        self._status = status
        self._pid = pid
        self._uptime = uptime
        self._path = serverPath

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGServerStatus:
        name = instream.readUTF()
        statusByte = instream.readByte()
        status = tgadmin.ServerState.fromId(statusByte)
        pid = str(instream.readInt())
        durnum = instream.readLong() / 1000000000.0
        vers = tgvers.TGVersion.readExternal(instream, reverse=True)
        serverPath = instream.readUTF()

        dur = datetime.datetime.now() - datetime.datetime.fromtimestamp(durnum)

        return cls(name, vers, status, pid, dur, serverPath)


class ProcessorSetStatisticsImpl(tgadmin.TGProcessorSetStatistics):

    @property
    def processorsCount(self) -> int:
        return self._txnProcCount

    @property
    def processedCount(self) -> int:
        return self._txnProcdCount

    @property
    def successfulCount(self) -> int:
        return self._txnSuccessCount

    @property
    def averageProcessingTime(self) -> float:
        return self._avgProcTime

    @property
    def pendingCount(self) -> int:
        return self._pendTxnCount

    def __init__(self, txnProcCount: int, txnProcdCount: int, txnSuccessCount: int, txnAvgProcTime: float,
                 pendTxnCount: int):
        self._txnProcCount = txnProcCount
        self._txnProcdCount = txnProcdCount
        self._txnSuccessCount = txnSuccessCount
        self._avgProcTime = txnAvgProcTime
        self._pendTxnCount = pendTxnCount

    def __str__(self):
        return "ProcessorSetStatisticsImpl [processorsCount=" + str(self._txnProcCount) +\
                ", processedCount=" + str(self._txnProcdCount) + ", successfulCount="\
                + str(self._txnSuccessCount) + ", averageProcessingTime=" + str(self._avgProcTime)\
                + ", pendingCount=" + str(self._pendTxnCount) + "]"

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGProcessorSetStatistics:
        pcount = instream.readShort()
        pdcount = instream.readLong()
        scount = instream.readLong()
        avgptime = instream.readDouble()
        ptxncount = instream.readLong()

        return cls(pcount, pdcount, scount, avgptime, ptxncount)


class TransactionStatisticsImpl(ProcessorSetStatisticsImpl, tgadmin.TGTransactionStatistics):

    @property
    def transactionLoggerQueueDepth(self) -> int:
        return self._txnLogDepth

    def __init__(self, txnProcCount: int, txnProcdCount: int, txnSuccessCount: int, txnAvgProcTime: float,
                 pendTxnCount: int, txnLogDepth: int):
        super().__init__(txnProcCount, txnProcdCount, txnSuccessCount, txnAvgProcTime, pendTxnCount)
        self._txnLogDepth = txnLogDepth

    def __str__(self):
        return "TransactionStatisticsImpl [transactionProcessorsCount=" + str(self._txnProcCount) +\
                ", transactionProcessedCount=" + str(self._txnProcdCount) + ", transactionSuccessfulCount="\
                + str(self._txnSuccessCount) + ", averageProcessingTime=" + str(self._avgProcTime)\
                + ", pendingTransactionsCount=" + str(self._pendTxnCount) +  ", transactionLoggerQueueDepth="\
                + str(self._txnLogDepth) + "]"

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGTransactionStatistics:
        pcount = instream.readShort()
        pdcount = instream.readLong()
        scount = instream.readLong()
        avgptime = instream.readDouble()
        ptxncount = instream.readLong()
        txnlqd = instream.readInt()

        return cls(pcount, pdcount, scount, avgptime, ptxncount, txnlqd)


class UserInfoImpl(tgadmin.TGUserInfo):

    @property
    def systemType(self) -> tgmodel.TGSystemType:
        return tgmodel.TGSystemType.Principal

    def writeExternal(self, ostream: TGOutputStream):
        pass        # Should never need to write the user

    def readExternal(self, instream: TGInputStream):
        self._type = instream.readByte()
        self._id = instream.readInt()
        self._name = instream.readUTF()
        if self._type == tgmodel.TGSystemType.Principal.value:
            _ = instream.readBytes()  # Password Hash
            numRoles = instream.readInt()  # Number roles
            self._roles = []
            for i in range(numRoles):
                self._roles.append(instream.readInt())  # Role ID

    @property
    def type(self) -> int:
        return self._type

    @property
    def id(self) -> int:
        return self._id

    @property
    def name(self) -> str:
        return self._name

    @property
    def roles(self) -> typing.FrozenSet[int]:
        return frozenset(self._roles)

    def __init__(self, type: int, id: int, name: str):
        self._type = type
        self._id = id
        self._name = name
        self._roles = []

    def __str__(self):
        return "UserInfoImpl [type=" + str(self._type) + ", id=" + str(self._id) + ", name=" + str(self._name) + "]"

    @classmethod
    def createFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGUserInfo:
        ret = cls(tgmodel.TGSystemType.Principal.value, -1, '')
        ret.readExternal(instream)
        return ret


class ServerInfoImpl(tgadmin.TGServerInfo):

    @property
    def serverStatus(self) -> TGServerStatus:
        return self._status

    @property
    def netListenersInfo(self) -> typing.List[TGNetListenerInfo]:
        return self._netListenerInfos

    def memoryInfo(self, type: MemoryType) -> TGMemoryInfo:
        return self._memInfo.getMemoryInfo(type)

    @property
    def transactionsInfo(self) -> TGProcessorSetStatistics:
        return self._txnInfo

    @property
    def queryInfo(self) -> TGProcessorSetStatistics:
        return self._qryInfo

    @property
    def cacheInfo(self) -> TGCacheStatistics:
        return self._cacheInfo

    @property
    def databaseInfo(self) -> TGDatabaseStatistics:
        return self._dbInfo

    def __init__(self, status: tgadmin.TGServerStatus, netListInfos: typing.List[tgadmin.TGNetListenerInfo],
                 memInfo: tgadmin.TGServerMemoryInfo, txnInfo: tgadmin.TGTransactionStatistics,
                 cacheInfo: tgadmin.TGCacheStatistics, dbInfo: tgadmin.TGDatabaseStatistics,
                 qryInfo: tgadmin.TGProcessorSetStatistics):
        self._status = status
        self._netListenerInfos = netListInfos
        self._memInfo = memInfo
        self._txnInfo = txnInfo
        self._cacheInfo = cacheInfo
        self._dbInfo = dbInfo
        self._qryInfo = qryInfo

    def __str__(self):
        return "ServerInfoImpl [serverInfo=" + str(self._status) + ", netListenersInfo=" + str(self._netListenerInfos)\
                + ", memoryInfo=" + str(self._memInfo) + ", transactionsInfo=" + str(self._txnInfo) + ", cacheInfo=" +\
                str(self._cacheInfo) + ", databaseInfo=" + str(self._dbInfo) + "]"

    @classmethod
    def readFromStream(cls, instream: tgpdu.TGInputStream) -> tgadmin.TGServerInfo:

        serverStatus = ServerStatusImpl.readFromStream(instream)

        memInfo = ServerMemoryInfoImpl.readFromStream(instream)

        listenerInfos = []
        lCount = instream.readLong()
        for i in range(lCount):
            listenerInfos.append(NetListenerInfoImpl.readFromStream(instream))

        txnInfo = TransactionStatisticsImpl.readFromStream(instream)

        qryInfo = ProcessorSetStatisticsImpl.readFromStream(instream)

        cacheInfo = CacheStatisticsImpl.readFromStream(instream)

        dbInfo = DatabaseStatisticsImpl.readFromStream(instream)

        return cls(serverStatus, listenerInfos, memInfo, txnInfo, cacheInfo, dbInfo, qryInfo)


class AdminRequest(pduimpl.AbstractProtocolMessage):

    def writePayload(self, os: pduimpl.ProtocolDataOutputStream):
        pos = os.position
        os.writeInt(0)      # Data length
        os.writeInt(0)      # Data checksum

        os.writeInt(self.command.val)

        if self.additionalData is not None:
            self.additionalData(os)

        os.writeIntAt(pos, os.position - pos)
        #os.writeHash64At(pos + 4, pos + 8)

    def readPayload(self, instream: tgpdu.TGInputStream):
        raise tgexception.TGException('Read called on write only message!')

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.command: TGAdminCommand = TGAdminCommand.Invalid
        self.additionalData: typing.Optional[typing.Callable[[tgpdu.TGOutputStream], None]] = None


class AdminDescImpl(tgadmin.TGAdminDescription):

    def __init__(self):
        self._prime = None
        self._attrs: typing.List[tgmodel.TGAttributeDescriptor] = None
        self._idxs: typing.List[tgmodel.TGIndex] = None
        self._ft: tgmodel.TGNodeType = None
        self._tt: tgmodel.TGNodeType = None
        self._roles: typing.List[tgmodel.TGRole] = None

    @property
    def primary(self):
        return self._prime

    @property
    def attributes(self):
        return self._attrs

    @property
    def indices(self):
        return self._idxs

    @property
    def roles(self) -> typing.List[tgmodel.TGRole]:
        return self._roles

    @property
    def fromType(self):
        return self._ft

    @property
    def toType(self):
        return self._tt


class AdminResponse(pduimpl.AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.gmd: tgmodel.TGGraphMetadata = None
        self.__rawCommand: int = None
        self.__errorStr: str = None
        self.__storedProcedures: typing.List[tgmodel.TGStoredProcedure] = None
        self.__instream: pduimpl.ProtocolDataInputStream = None
        self.__command: TGAdminCommand = TGAdminCommand.Invalid
        self.__alreadyRead = False
        self.__users: typing.List[tgadmin.TGUserInfo] = None
        self.__roles: typing.List[tgmodel.TGRole] = None
        self.__connections: typing.List[tgadmin.TGConnectionInfo] = None
        self.__attrdescs: typing.List[tgmodel.TGAttributeDescriptor] = None
        self.__indices: typing.List[tgmodel.TGIndex] = None
        self.__info: tgadmin.TGServerInfo = None
        self.__except: tgexception.TGException = None
        self.__result: int = 0
        self.__description: AdminDescImpl = None
        self.__types: typing.List[tgmodel.TGEntityType] = None
        self.__numKilled: int = None

    @property
    def rawCommand(self) -> int:
        return self.__rawCommand

    @property
    def errorStr(self) -> str:
        return self.__errorStr

    @property
    def users(self) -> typing.List[tgadmin.TGUserInfo]:
        if self.__alreadyRead and not self.__users:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__users:
            self.__alreadyRead = True
            self.__users = []
            count = self.__instream.readInt()

            for i in range(count):
                self.__users.append(UserInfoImpl.createFromStream(self.__instream))
        return self.__users

    @property
    def roles(self) -> typing.List[tgmodel.TGRole]:
        if self.__alreadyRead and not self.__roles:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__roles:
            self.__alreadyRead = True
            self.__roles = []
            count = self.__instream.readInt()

            for i in range(count):
                toAdd = gmdimpl.RoleImpl()
                toAdd.readExternal(self.__instream)
                self.__roles.append(toAdd)
        return self.__roles

    @property
    def storedProcedures(self) -> typing.List[tgmodel.TGStoredProcedure]:
        if self.__alreadyRead and not self.__storedProcedures:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__storedProcedures:
            self.__alreadyRead = True
            self.__storedProcedures = []
            count = self.__instream.readInt()

            for i in range(count):
                toAdd = gmdimpl.StoredProcedureImpl()
                toAdd.readExternal(self.__instream)
                self.__storedProcedures.append(toAdd)
        return self.__storedProcedures

    @property
    def connections(self) -> typing.List[tgadmin.TGConnectionInfo]:
        if self.__alreadyRead and not self.__connections:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__connections:
            self.__alreadyRead = True
            self.__connections = []
            count = self.__instream.readLong()

            for i in range(count):
                self.__connections.append(ConnectionInfoImpl.readFromStream(self.__instream))
        return self.__connections

    @property
    def indices(self) -> typing.List[tgmodel.TGIndex]:
        if self.__alreadyRead and not self.__indices:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__indices:
            self.__alreadyRead = True
            self.__indices = []
            count = self.__instream.readInt()

            for i in range(count):
                idx = gmdimpl.IndexImpl(self.gmd.registry)
                idx.readExternal(self.__instream)
                self.__indices.append(idx)
        return self.__indices

    @property
    def attrdescs(self) -> typing.List[tgmodel.TGAttributeDescriptor]:
        if self.__alreadyRead and not self.__attrdescs:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__attrdescs:
            self.__alreadyRead = True
            self.__attrdescs = []
            count = self.__instream.readInt()

            for i in range(count):
                attrdesc = attrimpl.AttributeDescriptorImpl()
                attrdesc.readExternal(self.__instream)
                self.__attrdescs.append(attrdesc)
        return self.__attrdescs

    @property
    def info(self) -> tgadmin.TGServerInfo:
        if self.__alreadyRead and not self.__info:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__info:
            self.__alreadyRead = True
            self.__info = ServerInfoImpl.readFromStream(self.__instream)
        return self.__info

    @property
    def command(self) -> TGAdminCommand:
        return self.__command

    @property
    def exception(self) -> typing.Optional[tgexception.TGException]:
        return self.__except

    @property
    def result(self) -> int:
        return self.__result

    @property
    def numKilled(self) -> int:
        if self.__alreadyRead and self.__numKilled is None:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif self.__numKilled is None:
            self.__alreadyRead = True
            self.__numKilled = self.__instream.readUnsignedInt()
        return self.__numKilled

    def types(self, registry: gmdimpl.TypeRegistry) -> typing.List[tgmodel.TGEntityType]:
        if self.__alreadyRead and not self.__types:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__types:
            self.__alreadyRead = True
            self.__types = []
            numEnts = self.__instream.readInt()
            for i in range(numEnts):
                nodetype = gmdimpl.NodeTypeImpl(registry)
                nodetype.readExternal(self.__instream)
                self.__types.append(nodetype)
            numEnts = self.__instream.readInt()
            for i in range(numEnts):
                edgetype = gmdimpl.EdgeTypeImpl(registry)
                edgetype.readExternal(self.__instream)
                self.__types.append(edgetype)
        return self.__types

    def description(self, registry: gmdimpl.TypeRegistry) -> tgadmin.TGAdminDescription:
        if self.__alreadyRead and not self.__description:
            raise tgexception.TGException('Already read from this stream into another object!')
        elif not self.__description:
            sysType = tgmodel.TGSystemType.fromValue(self.__instream.readByte())
            ret = AdminDescImpl()
            self.__instream.__currpos__ -= 1                    # Go back one byte.
            cont: bool = True
            if sysType == tgmodel.TGSystemType.NodeType:
                primary = gmdimpl.NodeTypeImpl(registry)
                primary.readExternal(self.__instream)
                ret._prime = primary
                num_indices = self.__instream.readLong()
                toAdd = []
                for i in range(num_indices):
                    index = gmdimpl.IndexImpl(registry)
                    index.readExternal(self.__instream)
                    toAdd.append(index)
                ret._idxs = toAdd
            elif sysType == tgmodel.TGSystemType.EdgeType:
                primary = gmdimpl.EdgeTypeImpl(registry)
                primary.readExternal(self.__instream)
                ret._prime = primary
                readNodetype = self.__instream.readBoolean()
                fromType = None
                if readNodetype:
                    fromType = gmdimpl.NodeTypeImpl(registry)
                    fromType.readExternal(self.__instream)
                readNodetype = self.__instream.readBoolean()
                toType = None
                if readNodetype:
                    toType = gmdimpl.NodeTypeImpl(registry)
                    toType.readExternal(self.__instream)
                ret._ft = fromType
                ret._tt = toType
            elif sysType == tgmodel.TGSystemType.AttributeDescriptor:
                primary = attrimpl.AttributeDescriptorImpl()
                primary.readExternal(self.__instream)
                ret._prime = primary
                cont = False
            elif sysType == tgmodel.TGSystemType.Index:
                primary = gmdimpl.IndexImpl(registry)
                primary.readExternal(self.__instream)
                ret._prime = primary
            elif sysType == tgmodel.TGSystemType.StoredProcedure:
                primary = gmdimpl.StoredProcedureImpl()
                primary.readExternal(self.__instream)
                ret._prime = primary
                cont = False
            elif sysType == tgmodel.TGSystemType.Role:
                primary = gmdimpl.RoleImpl()
                primary.readExternal(self.__instream)
                ret._prime = primary
                cont = False
            elif sysType == tgmodel.TGSystemType.Principal:
                primary = UserInfoImpl.createFromStream(self.__instream)
                numRoles = self.__instream.readInt()
                roles = []
                for i in range(numRoles):
                    role = gmdimpl.RoleImpl()
                    role.readExternal(self.__instream)
                    roles.append(role)
                ret._roles = roles
                ret._prime = primary
                cont = False
            else:
                raise tgexception.TGException('Not allowed to describe anything but nodetypes, edgetypes, attribute des'
                                              'criptors, and indices.')

            if cont:
                num_attrs = self.__instream.readLong()
                toAdd = []
                for i in range(num_attrs):
                    attrdesc = attrimpl.AttributeDescriptorImpl()
                    attrdesc.readExternal(self.__instream)
                    toAdd.append(attrdesc)
                ret._attrs = toAdd

            self.__description = ret

        return self.__description

    def writePayload(self, os: TGOutputStream):
        pass

    def readPayload(self, instream: TGInputStream):
        self.__instream = instream
        result = instream.readInt()
        self.__rawCommand = instream.readInt()
        self.__command = TGAdminCommand.fromId(self.__rawCommand)
        if result != 0:
            self.__result = result
            self.__errorStr = 'Error processing admin request.'
            try:
                self.__errorStr = instream.readUTF()
            except:
                self.__errorStr = 'Error processing admin request.'
            self.__except = tgexception.TGException(self.__errorStr.strip(), errorcode=result)


def _sortSystemObjects(to_sort: typing.List[tgmodel.TGSystemObject]):

    def key(obj: tgmodel.TGSystemObject):
        return "%02d%s" % (obj.systemType.value, obj.name)

    to_sort.sort(key=key)


def getElipsisIfNecessary(name: str, length: int = 17) -> str:
    ret = name
    if len(name) > length:
        ret = ("%-." + str(length) + "s...") % name
    return ret


class AdminConnectionImpl(connimpl.ConnectionImpl, tgadmin.TGAdminConnection):

    def getInfo(self) -> tgadmin.TGServerInfo:
        req = self.__newRequest()
        req.command = TGAdminCommand.ShowInfo
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.ShowInfo:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        return res.info

    def getTypes(self) -> typing.List[tgmodel.TGEntityType]:
        req = self.__newRequest()
        req.command = TGAdminCommand.ShowTypes
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.ShowTypes:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        return res.types(self.graphMetadata.registry)

    def showTypes(self):
        types = self.getTypes()
        _sortSystemObjects(types)
        print(" %-20s  %-1s %-9s  %-20s" % ("Name", "T", "SysId", "#Entries"))
        for entType in types:
            if entType.name[0] == '@':
                continue
            kind = 'O'
            if entType.systemType == tgmodel.TGSystemType.NodeType:
                kind = 'N'
            elif entType.systemType == tgmodel.TGSystemType.EdgeType:
                kind = 'E'
            name = entType.name
            if name[0] == '$':
                if kind == 'N':
                    name = 'Default Nodetype'
                else:
                    edgetype: tgmodel.TGEdgeType = entType
                    if edgetype.directiontype == tgmodel.DirectionType.UnDirected:
                        name = "Default Undirected Edgetype"
                    elif edgetype.directiontype == tgmodel.DirectionType.Directed:
                        name = "Default Directed Edgetype"
                    elif edgetype.directiontype == tgmodel.DirectionType.BiDirectional:
                        name = "Default Bidriected Edgetype"
            print(" %-20s  %-1s %-9s  %-20s" % (getElipsisIfNecessary(name), kind, entType.id, entType.numberEntities))

        print(str(len(types)) + " type" + ("s" if len(types) > 1 else "") + " returned.")

    def fullImport(self, dir_path: str = "import") -> typing.List[tgadmin.TGImportDescriptor]:
        import os

        filenames = os.listdir(dir_path)
        file_batch_max_size = (1 << 15) - 4       # Changed this due to a weird SSL Bug.
        file_num_batches = {}

        total_batches = 0
        to_iterate = []
        for filename in filenames:
            if len(filename) < 5 or (filename[-4:] != '.csv' and filename[-5:] != '.conf') or filename[0] == '.':
                if filename != '.' and filename != '..':
                    tglog.gLogger.log(tglog.TGLevel.Warning, "Skipping file %s because it is not a CSV file or a Graph "
                                                             "Database configuration file!", filename)
                continue
            filesize = os.path.getsize(dir_path + os.sep + filename)
            overflow = filesize % file_batch_max_size
            base_batches = filesize // file_batch_max_size
            if overflow > 0:
                base_batches += 1
            file_num_batches[filename] = base_batches
            total_batches += base_batches
            to_iterate.append(filename)

        filenames = to_iterate
        waiter = self._genBCRWaiter()
        message: pduimpl.BeginImportSessionRequest =\
            pduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.BeginImportRequest, authtoken=self.__channel__.authtoken,
                                                   sessionid=self.__channel__.sessionid)

        message.isBatch = False
        message.numFiles = len(filenames)
        message.totalRequsts = total_batches

        response: pduimpl.BeginImportSessionResponse = self.__channel__.send(message, waiter)
        if response.error is not None:
            raise response.error

        totReqNum = 0
        descMap: typing.Dict[str, ImportDescImpl] = {}
        for filename in filenames:
            file_batches = file_num_batches[filename]
            filepath = dir_path + os.sep + filename
            with open(filepath, 'rb') as file:
                for i in range(file_batches):
                    buffer: bytes
                    if i + 1 < file_batches:
                        buffer = file.read(file_batch_max_size)
                    else:
                        buffer = file.read()
                    req: pduimpl.PartialImportRequest =\
                        pduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.PartialImportRequest, authtoken=self.__channel__.authtoken, sessionid=self.__channel__.authtoken)
                    req.isBatch = False
                    req.fileIdx = totReqNum
                    totReqNum += 1
                    req.fileTotalReq = file_batches
                    req.fileReqIdx = i
                    req.fileName = filename
                    req.data = buffer

                    response: pduimpl.PartialImportResponse = self.__channel__.send(req, self._genBCRWaiter())

                    if response.resultList is not None:
                        for desc in response.resultList:
                            if desc.typename in descMap:
                                descMap[desc.typename]._inc(desc.numinstances)
                            else:
                                descMap[desc.typename] = desc

                    if response.error is not None:
                        raise response.error

        self.__initMetadata__()
        return [descMap[descName] for descName in descMap]

    def getUsers(self) -> typing.List[tgadmin.TGUserInfo]:
        req = self.__newRequest()
        req.command = TGAdminCommand.ShowUsers
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.ShowUsers:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        return res.users

    def getRoles(self) -> typing.List[tgmodel.TGRole]:
        req = self.__newRequest()
        req.command = TGAdminCommand.ShowRoles
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.ShowRoles:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        return res.roles

    def getConnections(self) -> typing.List[tgadmin.TGConnectionInfo]:
        req = self.__newRequest()
        req.command = TGAdminCommand.ShowConnections
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.ShowConnections:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        return res.connections

    def stopServer(self):
        req = self.__newRequest()
        req.command = TGAdminCommand.StopServer
        self.__channel__.send(req)
        self.disconnect()

    def dumpServerStacktrace(self):
        sendMsg = pduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.DumpStacktraceRequest,
                                        authtoken=self.__channel__.authtoken, sessionid=self.__channel__.sessionid)
        self.__channel__.send(sendMsg)

    def checkpointServer(self):
        req = self.__newRequest()
        req.command = TGAdminCommand.CheckpointServer
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.CheckpointServer:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def killConnection(self, sessionid: typing.Optional[int] = None) -> int:
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUnsignedLong(sessionid if sessionid is not None else 0)
            outstream.writeBoolean(sessionid is None)
        req = self.__newRequest()
        req.command = TGAdminCommand.KillConnection
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.KillConnection:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        return res.numKilled

    def getAttributeDescriptors(self) -> typing.List[tgmodel.TGAttributeDescriptor]:
        req = self.__newRequest()
        req.command = TGAdminCommand.ShowAttrDescs
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.ShowAttrDescs:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        return res.attrdescs

    def getIndices(self) -> typing.List[tgmodel.TGIndex]:
        req = self.__newRequest()
        req.command = TGAdminCommand.ShowIndices
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.ShowIndices:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        res.gmd = self.graphMetadata
        return res.indices

    def setServerLogLevel(self, logcomponent: typing.Union[tglog.TGLogComponent, typing.List[tglog.TGLogComponent]],
                          level: tglog.TGLevel):
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeShort(level.value)
            if isinstance(logcomponent, list):
                outstream.writeUnsignedLong(tglog.TGLogComponent.aggregate_bits(logcomponent))
            else:
                outstream.writeUnsignedLong(logcomponent.bits)
        req = self.__newRequest()
        req.command = TGAdminCommand.SetLogLevel
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.SetLogLevel:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def describe(self, name: str) -> tgadmin.TGAdminDescription:
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(name)
        req = self.__newRequest()
        req.command = TGAdminCommand.Describe
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.exception is not None:
            raise res.exception
        if res.command != TGAdminCommand.Describe:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))
        return res.description(self.graphMetadata.registry)

    def createAttrDesc(self, name: str, ad_type: tgmodel.TGAttributeType, isArray: bool = False,
                       isEncrypted: bool = False, prec: int = None, scale: int = None):
        if ad_type != tgmodel.TGAttributeType.Number and (prec is not None or scale is not None):
            raise tgexception.TGException('Not allowed to set precision or scale on non-Number Attribute Type!')
        if prec is None:
            prec = 20
        if scale is None:
            scale = 5
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeBoolean(isArray)
            outstream.writeBoolean(isEncrypted)
            outstream.writeUTF(name)
            outstream.writeInt(ad_type.identifer)
            if ad_type == tgmodel.TGAttributeType.Number:
                outstream.writeUnsignedShort(prec)
                outstream.writeUnsignedShort(scale)
        req = self.__newRequest()
        req.command = TGAdminCommand.CreateAttrDesc
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            if res.result == 343:
                pass                    # Duplicate Attribute descriptor
            else:
                raise res.exception
        if res.command != TGAdminCommand.CreateAttrDesc:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

        self.__initMetadata__()         # Do this to refresh the metadata.

    def __validUserName(self, name: str):
        ret = True
        if len(name) < 1:
            ret = False
        elif len(name) > 255:
            ret = False
        elif '@' in name:
            ret = False
        return ret

    def createUser(self, name: str, password: str, roles: typing.List[str] = None):
        if roles is None:
            roles = []
        if not self.__validUserName(name):
            raise tgexception.TGException("Bad username: {0}".format(name))
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(name)
            outstream.writeUTF(password)
            #if roles is None or len(roles) == 0:
            #    outstream.writeUTF("users")
            #else:
            #    outstream.writeUTF(roles[0])
            outstream.writeInt(len(roles))
            for role in roles:
                outstream.writeUTF(role)
        req = self.__newRequest()
        req.command = TGAdminCommand.CreateUser
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            if res.result == 342:       # Duplicate System Object
                raise tgexception.TGException('Duplicated system object. Cannot have the username be the same as any ot'
                                              'her user, nor the same as a nodetype, edgetype, or index name.')
            elif res.result == 328:     # Invalid Role
                raise tgexception.TGException('Invalid role: {0}'.format(roles))
            else:
                raise res.exception
        if res.command != TGAdminCommand.CreateUser:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def __checkAttrDescs(self, attrDescs: typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]])\
                -> typing.List[tgmodel.TGAttributeDescriptor]:
        ad_list = []
        for ad in attrDescs:
            if isinstance(ad, tgmodel.TGEntityType):
                ad_list.append(ad)
            else:
                if ad in self.graphMetadata.attritubeDescriptors:
                    ad_list.append(self.graphMetadata.attritubeDescriptors[ad])
                else:
                    raise tgexception.TGException('Attritube Descriptor not in database: ' + str(ad))

        return ad_list

    def __checkTypes(self, types: typing.Optional[typing.List[typing.Union[tgmodel.TGEntityType, str]]],
                     allowNodetypes = True, allowEdgetypes = True) -> typing.List[tgmodel.TGEntityType]:
        type_list = []

        if types is not None:
            for ty in types:
                if isinstance(ty, tgmodel.TGEntityType):
                    type_list.append(ty)
                else:
                    if ty in self.graphMetadata.edgetypes and allowEdgetypes:
                        type_list.append(self.graphMetadata.edgetypes[ty])
                    elif ty in self.graphMetadata.nodetypes and allowNodetypes:
                        type_list.append(self.graphMetadata.nodetypes[ty])
                    else:
                        raise tgexception.TGException('Entity type not in database: ' + str(ty))
        return type_list

    def createIndex(self, name: str, attrDescs: typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]],
                    isUnique: bool = False,
                    types: typing.Optional[typing.List[typing.Union[tgmodel.TGEntityType, str]]] = None):
        # Error checking

        ad_list = self.__checkAttrDescs(attrDescs)

        type_list = self.__checkTypes(types)

        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeBoolean(isUnique)
            outstream.writeUTF(name)
            outstream.writeInt(len(ad_list))
            for ad in ad_list:
                outstream.writeChars(ad.name)
            outstream.writeInt(len(type_list))
            for ty in type_list:
                outstream.writeChars(ty.name)

        req = self.__newRequest()
        req.command = TGAdminCommand.CreateIndex
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            if res.result == 342:  # Duplicate System Object
                raise tgexception.TGException('Duplicated system object. Cannot have the index name be the same as any '
                                              'other index, nor the same as a nodetype, edgetype, or user name.')
            else:
                raise res.exception
        if res.command != TGAdminCommand.CreateIndex:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

        self.__initMetadata__()  # Do this to refresh the metadata.

    def createNodetype(self, name: str, attrDescs: typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]],
                    pageSize: int = None,
                    pkeys: typing.Optional[typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]]] = None):
        # Error checking

        ad_list = self.__checkAttrDescs(attrDescs)

        pkey_list = []
        if pkeys is not None:
            pkey_list = self.__checkAttrDescs(pkeys)

        ad_set = set()
        [ad_set.add(ad.name) for ad in ad_list]

        for pkey in pkey_list:
            if pkey.name not in ad_set:
                raise tgexception.TGException('Primary key attribute %s was not in list of attributes for this nodetype'
                                              '!' % pkey.name)

        if pageSize is None:
            pageSize = 512

        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(name)
            outstream.writeUnsignedInt(pageSize)
            outstream.writeInt(len(ad_list))
            for ad in ad_list:
                outstream.writeChars(ad.name)
            outstream.writeInt(len(pkey_list))
            for pkey in pkey_list:
                outstream.writeChars(pkey.name)

        req = self.__newRequest()
        req.command = TGAdminCommand.CreateNodeType
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            if res.result == 342:  # Duplicate System Object
                raise tgexception.TGException('Duplicated system object. Cannot have the nodetype name be the same as a'
                                              'ny other nodetype, nor the same as an index, edgetype, or user name.')
            else:
                raise res.exception
        if res.command != TGAdminCommand.CreateNodeType:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

        self.__initMetadata__()  # Do this to refresh the metadata.

    def __checkNodetypeValid(self, nodetype: typing.Optional[typing.Union[str, tgmodel.TGNodeType]] = None) -> str:
        ret: str
        if nodetype is None:
            ret = ""
        else:
            if isinstance(nodetype, str):
                if nodetype not in self.graphMetadata.nodetypes:
                    raise Exception
                ret = nodetype
            else:
                ret = nodetype.name

        return ret

    def createEdgetype(self, name: str, attrDescs: typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]],
                    dirType: typing.Optional[typing.Union[str, tgmodel.DirectionType]] = None,
                    fromType: typing.Optional[typing.Union[str, tgmodel.TGNodeType]] = None,
                    toType: typing.Optional[typing.Union[str, tgmodel.TGNodeType]] = None,
                    pkeys: typing.Optional[typing.List[typing.Union[tgmodel.TGAttributeDescriptor, str]]] = None):
        # Error checking

        ad_list = self.__checkAttrDescs(attrDescs)

        pkey_list = []
        if pkeys is not None:
            pkey_list = self.__checkAttrDescs(pkeys)

        dir_type: str = dirType
        if dirType is None:
            dir_type = 'directed'
        else:
            if isinstance(dirType, tgmodel.DirectionType):
                if dirType == tgmodel.DirectionType.Directed:
                    dir_type = 'directed'
                elif dirType == tgmodel.DirectionType.UnDirected:
                    dir_type = 'undirected'
                elif dirType == tgmodel.DirectionType.BiDirectional:
                    dir_type = 'bidirected'
                else:
                    raise tgexception.TGException('Unknown direction type: ' + str(dirType))
            else:
                if dirType not in ('directed', 'undirected', 'bidirected'):
                    raise tgexception.TGException('Unknown direction type as string: ' + dirType)

        ad_set = set()
        [ad_set.add(ad.name) for ad in ad_list]

        from_type = self.__checkNodetypeValid(fromType)
        to_type = self.__checkNodetypeValid(toType)

        for pkey in pkey_list:
            if pkey.name not in ad_set:
                raise tgexception.TGException('Primary key attribute {0} was not in list of attributes for this nodetyp'
                                              'e!'.format(pkey.name))

        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(name)
            outstream.writeUTF(from_type)
            outstream.writeUTF(to_type)
            outstream.writeUTF(dir_type)
            outstream.writeInt(len(ad_list))
            for ad in ad_list:
                outstream.writeChars(ad.name)
            outstream.writeInt(len(pkey_list))
            for pkey in pkey_list:
                outstream.writeChars(pkey.name)

        req = self.__newRequest()
        req.command = TGAdminCommand.CreateEdgeType
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            if res.result == 342:  # Duplicate System Object
                raise tgexception.TGException('Duplicated system object. Cannot have the edgetype name be the same as a'
                                              'ny other edgetype, nor the same as an index, nodetype, or user name.')
            else:
                raise res.exception
        if res.command != TGAdminCommand.CreateEdgeType:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

        self.__initMetadata__()  # Do this to refresh the metadata.

    def setSPDirectory(self, dir_path: str):

        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(dir_path)

        req = self.__newRequest()
        req.command = TGAdminCommand.SetSPDir
        req.additionalData = write
        res = self.__sendMessage(req)

        if res.result != 0:
            raise res.exception

        if res.command != TGAdminCommand.SetSPDir:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def getStoredProcedures(self) -> typing.List[tgmodel.TGStoredProcedure]:
        req = self.__newRequest()
        req.command = TGAdminCommand.ShowSP
        res = self.__sendMessage(req)

        if res.result != 0:
            raise res.exception

        if res.command != TGAdminCommand.ShowSP:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

        return res.storedProcedures

    def refreshStoredProcedures(self):
        req = self.__newRequest()
        req.command = TGAdminCommand.RefreshSPDir
        res = self.__sendMessage(req)

        if res.result != 0:
            raise res.exception

        if res.command != TGAdminCommand.RefreshSPDir:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def __convertToStreamStyle(self, permissions: typing.List[typing.Tuple[typing.Union[str, tgmodel.TGNodeType,
                  tgmodel.TGEdgeType, tgmodel.TGAttributeDescriptor], typing.Union[tgmodel.TGPermissionType, str]]])\
                  -> typing.List[typing.Tuple[str, int]]:
        """Converts a list of permissions into a list of system identifiers and byte-flag representation of the
        corresponding permission."""
        ret = []
        num = 0
        for perm_tuple in permissions:
            perm_type = perm_tuple[0]
            perm_role = perm_tuple[1]
            to_add_name: str
            to_add_role: int
            if isinstance(perm_type, str):
                to_add_name = perm_type
            else:
                to_add_name = perm_type.name
            if isinstance(perm_role, str):
                to_add_role = tgmodel.TGPermissionType.flagFromStr(perm_role)
            else:
                to_add_role = perm_role.flag
            ret.append((to_add_name, to_add_role))
            num += 1
        return ret

    def createRole(self, rolename: str, permissions: typing.List[typing.Tuple[typing.Union[str, tgmodel.TGNodeType,
                  tgmodel.TGEdgeType, tgmodel.TGAttributeDescriptor], typing.Union[tgmodel.TGPermissionType, str]]],
                   privileges: typing.Union[typing.Iterable[tgmodel.TGPrivilege], tgmodel.TGPrivilege] = tuple()):

        role_permissions = self.__convertToStreamStyle(permissions)

        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(rolename)
            outstream.writeUnsignedLong(tgmodel.TGPrivilege.privilegesToInt(privileges))

            num_perms = len(role_permissions)

            outstream.writeUnsignedInt(num_perms)

            for permission in role_permissions:
                name, perm = permission

                outstream.writeUTF(name)
                outstream.writeUnsignedLong(perm)
                outstream.writeBoolean(True)

        req = self.__newRequest()
        req.command = TGAdminCommand.CreateRole
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            if res.result == 342:  # Duplicate System Object
                raise tgexception.TGException('Duplicated system object. Cannot have the role name be the same as a'
                                              'ny other role, nor the same as an index, type, or user name.')
            else:
                raise res.exception
        if res.command != TGAdminCommand.CreateRole:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def grantPermissions(self, rolename: str, permissions: typing.List[typing.Tuple[
                    typing.Union[str, tgmodel.TGNodeType, tgmodel.TGEdgeType, tgmodel.TGAttributeDescriptor],
                    typing.Union[tgmodel.TGPermissionType, str]]]):
        role_permissions = self.__convertToStreamStyle(permissions)

        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(rolename)
            outstream.writeBoolean(False)   # Clear permissions
            outstream.writeBoolean(True)    # Is granting privileges.
            outstream.writeUnsignedLong(0)

            num_perms = len(role_permissions)

            outstream.writeUnsignedInt(num_perms)

            for permission in role_permissions:
                name, perm = permission

                outstream.writeUTF(name)
                outstream.writeUnsignedLong(perm)
                outstream.writeBoolean(True)

        req = self.__newRequest()
        req.command = TGAdminCommand.UpdateRole
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            raise res.exception
        if res.command != TGAdminCommand.UpdateRole:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def grantPrivileges(self, rolename: str,
                         privileges: typing.Union[typing.Iterable[tgmodel.TGPrivilege], tgmodel.TGPrivilege] = tuple()):
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(rolename)
            outstream.writeBoolean(False)   # Clear permissions
            outstream.writeBoolean(True)    # Is granting privileges.
            outstream.writeUnsignedLong(tgmodel.TGPrivilege.privilegesToInt(privileges))
            outstream.writeUnsignedInt(0)

        req = self.__newRequest()
        req.command = TGAdminCommand.UpdateRole
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            raise res.exception
        if res.command != TGAdminCommand.UpdateRole:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def revokePermissions(self, rolename: str, permissions: typing.List[typing.Tuple[
                    typing.Union[str, tgmodel.TGNodeType,tgmodel.TGEdgeType, tgmodel.TGAttributeDescriptor],
                    typing.Union[tgmodel.TGPermissionType, str]]]):
        role_permissions = self.__convertToStreamStyle(permissions)

        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(rolename)
            outstream.writeBoolean(False)   # Clear permissions
            outstream.writeBoolean(False)   # Is revoking privileges. NOT granting them.
            outstream.writeUnsignedLong(0)

            num_perms = len(role_permissions)

            outstream.writeUnsignedInt(num_perms)

            for permission in role_permissions:
                name, perm = permission

                outstream.writeUTF(name)
                outstream.writeUnsignedLong(perm)
                outstream.writeBoolean(False)

        req = self.__newRequest()
        req.command = TGAdminCommand.UpdateRole
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            raise res.exception
        if res.command != TGAdminCommand.UpdateRole:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def revokePrivileges(self, rolename: str, privileges: typing.Union[typing.Iterable[tgmodel.TGPrivilege],
                                                                       tgmodel.TGPrivilege] = tuple()):
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(rolename)
            outstream.writeBoolean(False)   # Clear permissions
            outstream.writeBoolean(False)   # Is revoking privileges. NOT granting them.
            outstream.writeUnsignedLong(tgmodel.TGPrivilege.privilegesToInt(privileges))
            outstream.writeUnsignedInt(0)

        req = self.__newRequest()
        req.command = TGAdminCommand.UpdateRole
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            raise res.exception
        if res.command != TGAdminCommand.UpdateRole:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def clearRole(self, rolename: str):
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(rolename)
            outstream.writeBoolean(True)
            outstream.writeBoolean(True)    # Just need to write a boolean here. Multipurpose
            outstream.writeUnsignedLong(0)

            outstream.writeUnsignedInt(0)   # The number of privileges. AKA none.

        req = self.__newRequest()
        req.command = TGAdminCommand.UpdateRole
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            raise res.exception
        if res.command != TGAdminCommand.UpdateRole:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def grantUser(self, username: str, roles: typing.List[str]):
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(username)
            outstream.writeBoolean(False)           # Will send password (set to False)
            outstream.writeBoolean(True)            # Will send roles to grant or revoke (set to True)
            outstream.writeInt(len(roles))
            for rolename in roles:
                outstream.writeBoolean(True)        # Is Granting role
                outstream.writeUTF(rolename)        # Rolename
        req = self.__newRequest()
        req.command = TGAdminCommand.UpdateUser
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            raise res.exception
        if res.command != TGAdminCommand.UpdateUser:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def revokeUser(self, username: str, roles: typing.List[str]):
        def write(outstream: tgpdu.TGOutputStream):
            outstream.writeUTF(username)
            outstream.writeBoolean(False)           # Will send password (set to False)
            outstream.writeBoolean(True)            # Will send roles to grant or revoke (set to True)
            outstream.writeInt(len(roles))
            for rolename in roles:
                outstream.writeBoolean(False)       # Is Granting role
                outstream.writeUTF(rolename)        # Rolename
        req = self.__newRequest()
        req.command = TGAdminCommand.UpdateUser
        req.additionalData = write
        res = self.__sendMessage(req)
        if res.result != 0:
            raise res.exception
        if res.command != TGAdminCommand.UpdateUser:
            raise tgexception.TGException("Did not expect {0} response.".format(str(res.command)))

    def __init__(self, url, username, password, dbName, env):
        connimpl.ConnectionImpl.__init__(self, url, username, password, dbName, env)

    def __sendMessage(self, request: AdminRequest) -> AdminResponse:
        waiter = self._genBCRWaiter()
        return self.__channel__.send(request, waiter)

    def __newRequest(self) -> AdminRequest:
        return pduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.AdminRequest, authtoken=self.__channel__.authtoken,
                                                      sessionid=self.__channel__.sessionid)
