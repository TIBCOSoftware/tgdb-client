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
 *  File name :pdu.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: pdu.py 3255 2019-06-10 03:30:05Z ssubrama $
 *
 """

import abc
import enum
from tgdb.utils import *

@enum.unique
class VerbId(enum.Enum):
    """Server recognized verbs that describe what a message's purpose is and therefore how to interpret it.

    Should probably not be used by any client code.
    """

    # Ping Message - Heart  beats
    PingMessage     = 0

    # HandShake Request / Response protocol
    HandShakeRequest    = 1
    HandShakeResponse   = 2

    # Authenticate Request / Response protocol
    AuthenticateRequest = 3
    AuthenticateResponse    = 4

    # Transaction - begin / commit / rollback protocol verbs
    BeginTransactionRequest     = 5
    BeginTransactionResponse	= 6
    CommitTransactionRequest	= 7
    CommitTransactionResponse	= 8
    RollbackTransactionRequest	= 9
    RollbackTransactionResponse	= 10

    # Query Request / Response verbs
    QueryRequest	= 11
    QueryResponse	= 12

    # Graph Traversal verbs
    TraverseRequest	= 13
    TraverseResponse	= 14

    # Admin Request / Response verbs
    AdminRequest	= 15
    AdminResponse	= 16

    # Retrieve meta data
    MetadataRequest	    = 19
    MetadataResponse	= 20

    # Get entities
    GetEntityRequest	= 21
    GetEntityResponse	= 22

    # Get LargeObject
    GetLargeObjectRequest	= 23
    GetLargeObjectResponse	= 24

    # Import / Export verbs - They are admin request
    BeginExportRequest      = 25
    BeginExportResponse     = 26
    PartialExportRequest    = 27
    PartialExportResponse   = 28
    CancelExportRequest     = 29


    BeginImportRequest      = 31
    BeginImportResponse     = 32
    PartialImportRequest    = 33
    PartialImportResponse   = 34
    EndImportRequest = 35
    EndImportResponse = 36

    # Dump Stacktrace request verb
    DumpStacktraceRequest   = 39

    # Disconnect Request verbs
    DisconnectChannelRequest	= 40
    SessionForcefullyTerminated	= 41

    # Encrypt/Decrypt Request
    DecryptBufferRequest	= 44
    DecryptBufferResponse	= 45

    ConnectionPropertiesMessage = 48

    # New Export Stuff
    BeginBulkExportSessionRequest = 70
    BeginBulkExportSessionResponse = 71
    BeginBatchExportEntityRequest = 72
    BeginBatchExportEntityResponse = 73
    SingleBatchExportEntityRequest = 74
    SingleBatchExportEntityResponse = 75
    EndBulkExportSessionRequest = 76
    EndBulkExportSessionResponse = 77

    # Unknown Exception Message on the server.
    ExceptionMessage	    = 100

    # Bad Verb
    InvalidMessage	        = -1

    @classmethod
    def fromId(cls, id):
        """Get the verb corresponding to that id

        :param id: the identifier used to find the corresponding verb"""
        for verb in VerbId:
            if verb.value == id:
                return verb
        return VerbId.InvalidMessage


class TGOutputStream(abc.ABC):
    """A output writer for converting various Python objects to a uniform byte stream that the server can read.
    Should probably not be used by any client code.
    """

    @abc.abstractmethod
    def skip(self, n): pass

    @property
    @abc.abstractmethod
    def position(self): pass

    @property
    @abc.abstractmethod
    def buffer(self): pass

    @property
    @abc.abstractmethod
    def length(self): pass

    @abc.abstractmethod
    def writeBytes(self, buf): pass

    @abc.abstractmethod
    def writeChars(self, chars: str): pass

    @abc.abstractmethod
    def writeUTF(self, utfstr): pass

    @abc.abstractmethod
    def writeBoolean(self, b): pass

    @abc.abstractmethod
    def writeByte(self, x): pass

    @abc.abstractmethod
    def writeChar(self, c): pass

    @abc.abstractmethod
    def writeShort(self, h): pass

    @abc.abstractmethod
    def writeInt(self, i): pass

    @abc.abstractmethod
    def writeLong(self, q): pass

    @abc.abstractmethod
    def writeFloat(self, f): pass

    @abc.abstractmethod
    def writeDouble(self, d): pass

    @abc.abstractmethod
    def writeLongAsBytes(self, q): pass

    @abc.abstractmethod
    def writeIntAt(self, p, i): pass


class TGInputStream(abc.ABC):
    """An input reader for converting from a uniform byte stream from the server to various Python objects.

    Should probably not be used by any client code.
    """

    @abc.abstractmethod
    def read(self) -> int: pass

    @abc.abstractmethod
    def skip(self, n): pass

    @abc.abstractmethod
    def readChars(self) -> str: pass

    @abc.abstractmethod
    def readBytes(self) -> bytearray: pass

    @abc.abstractmethod
    def readFully(self, buf: bytearray, startpos=0, buflen=-1): pass
    
    @abc.abstractmethod
    def readUTF(self) -> str: pass

    @abc.abstractmethod
    def readBoolean(self): pass

    @abc.abstractmethod
    def readByte(self): pass

    @abc.abstractmethod
    def readChar(self): pass

    @abc.abstractmethod
    def readShort(self): pass

    @abc.abstractmethod
    def readInt(self): pass

    @abc.abstractmethod
    def readLong(self): pass

    @abc.abstractmethod
    def readFloat(self): pass

    @abc.abstractmethod
    def readDouble(self): pass

    @property
    @abc.abstractmethod
    def position(self): pass

    @position.setter
    @abc.abstractmethod
    def position(self, pos): pass


    """
    * A reference map for entity resolution 
    """
    @property
    @abc.abstractmethod
    def referencemap(self) -> dict: pass

    @referencemap.setter
    @abc.abstractmethod
    def referencemap(self, value): pass


class TGSerializable(abc.ABC):
    """Represents a serializable object.

    Should not be used by client code.
    """

    @abc.abstractmethod
    def writeExternal(self, ostream: TGOutputStream): pass

    @abc.abstractmethod
    def readExternal(self, istream: TGInputStream): pass


class TGMessage(abc.ABC):
    """Represents a message (either or to or from the server).

    Should not be used by client code.
    """

    @property
    @abc.abstractmethod
    def verbid(self) -> VerbId: pass

    @property
    @abc.abstractmethod
    def sequenceno(self) -> int: pass

    @property
    @abc.abstractmethod
    def requestid(self) -> int: pass

    @requestid.setter
    @abc.abstractmethod
    def requestid(self, value: int): pass

    @property
    @abc.abstractmethod
    def timestamp(self) -> int: pass

    @timestamp.setter
    @abc.abstractmethod
    def timestamp(self, value: int): pass

    @property
    @abc.abstractmethod
    def authtoken(self) -> int: pass

    @authtoken.setter
    @abc.abstractmethod
    def authtoken(self, value: int): pass

    @property
    @abc.abstractmethod
    def sessionid(self) -> int: pass

    @sessionid.setter
    @abc.abstractmethod
    def sessionid(self, value: int): pass

    @property
    @abc.abstractmethod
    def isUpdateable(self) -> bool: pass

    @abc.abstractmethod
    def updateSequenceAndTimeStamp(self): pass

    @abc.abstractmethod
    def toBytes(self, cipher) -> bytearray: pass

    @abc.abstractmethod
    def fromBytes(self, buf: bytearray): pass

    # def __str__(self):
    #     return HexUtils.formatHex(self.toBytes())
