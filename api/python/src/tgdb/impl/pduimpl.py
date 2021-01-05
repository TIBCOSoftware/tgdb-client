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
 *  File name :attr.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: pduimpl.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates all the Protocol Data Unit and its related classes
 """

import struct
import typing
import abc
import ssl
import datetime
import ctypes

import tgdb.log as tglog
from tgdb.pdu import *
from tgdb.impl.atomics import *
import tgdb.query as tgquery
import tgdb.impl.gmdimpl as tggmdimpl
import tgdb.model as tgmodel
import tgdb.impl.entityimpl as tgentimpl
import tgdb.exception as tgexception
import tgdb.impl.queryimpl as tgqueryimpl
import tgdb.impl.attrimpl as tgattrimpl
import tgdb.bulkio as tgbulk
import tgdb.utils as tgutils
import tgdb.channel as tgchan
import tgdb.admin as tgadm

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Data Input/Output Streams in PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class ProtocolDataOutputStream(TGOutputStream):

    def __init__(self, cipher = None):
        self._buffer = bytearray()
        self._count = 0
        self._pos = 0
        self._cipher = cipher

    @property
    def position(self):
        return self._pos

    @property
    def buffer(self):
        return self._buffer

    def skip(self, n):
        if (self._pos + n) >= len(self._buffer):
            raise tgexception.TGException("ArrayOutOfBounds exception")
        self._pos += n

    @property
    def length(self):
        return self.buffer.__len__()

    def writeBytes(self, buf: bytearray):
        _len = len(buf)
        self._buffer.extend(struct.pack('>i', _len))
        self._buffer.extend(buf)
        self._pos += (_len + 4)

    def writeUTF(self, s: str, encode: str = 'utf_8'):
        _buf = s.encode(encode)
        _len = len(_buf)
        self._buffer.extend(struct.pack('>h', _len))
        self._buffer.extend(_buf)
        self._pos += (_len + 2)

    def writeBoolean(self, b: bool):
        self._buffer.extend(struct.pack('?', b))
        self._pos += 1

    def writeByte(self, x):
        self._buffer.extend(struct.pack('b', x))
        self._pos += 1

    def writeUnsignedByte(self, x):
        self._buffer.extend(struct.pack('B', x))
        self._pos += 1

    # This Java Compliant Char : 2bytes - written as short
    def writeChar(self, c):
        self._buffer.extend(struct.pack('>h', ord(c)))
        self._pos += 2

    def writeChars(self, chars: str):
        length = len(chars)
        self.writeInt(length)
        for i in range(length):
            self.writeByte(ord(chars[i]))

    def writeShort(self, h):
        self._buffer.extend(struct.pack('>h', h))
        self._pos += 2

    def writeUnsignedShort(self, h):
        self._buffer.extend(struct.pack('>H', h))
        self._pos += 2

    def writeInt(self, i, littleEndian: bool = False):
        self._buffer.extend(struct.pack('<i' if littleEndian else '>i', i))
        self._pos += 4

    def writeUnsignedInt(self, i):
        self._buffer.extend(struct.pack('>I', i))
        self._pos += 4

    def writeLong(self, q, littleEndian: bool = False):
        self._buffer.extend(struct.pack('<q' if littleEndian else '>q', q))
        self._pos += 8

    def writeUnsignedLong(self, q, littleEndian: bool = False):
        self._buffer.extend(struct.pack('<Q' if littleEndian else '>Q', q))
        self._pos += 8

    def writeFloat(self, f):
        self._buffer.extend(struct.pack('>f', f))
        self._pos += 4

    def writeDouble(self, d):
        self._buffer.extend(struct.pack('>d', d))
        self._pos += 8

    def writeLongAsBytes(self, q):
        raise tgexception.TGException("Not implemented yet.")

    def writeIntAt(self, p, i):
        if p not in range(0, len(self._buffer)):
            raise tgexception.TGException('Invalid position :{p} specified.', format(p))
        struct.pack_into('>i', self._buffer, p, i)
        return

    def writeLongAt(self, p, q):
        if p not in range(0, len(self._buffer)):
            raise tgexception.TGException('Invalid position :{p} specified.', format(p))
        struct.pack_into('>q', self._buffer, p, q)
        return

    def writeHash64At(self, pos, startIdx, endIdx=-1, littleEndian: bool = False):
        hash: int = tgutils.tg_checkSum64(self._buffer[startIdx:endIdx if endIdx > startIdx else self._pos])
        struct.pack_into('<Q' if littleEndian else '>Q', self._buffer, pos, hash)

    def writeEncrypted(self, data: bytes):
        if self._cipher is None:
            raise tgexception.TGException("Trying to write encrypted data on an un-secure channel is not currently supp"
                                          "orted.")
        encrypted_data = self._cipher.encrypt(bytes(data))
        tglog.gLogger.log(tglog.TGLevel.DebugWire, str(encrypted_data))
        self.writeBytes(encrypted_data)


class ProtocolDataInputStream(TGInputStream):
    __buffer__: bytearray = None
    __buflen__: int = -1
    __currpos__: int = -1
    __startpos__: int = -1
    __refmap__: dict = None

    typecodeByteLength = {'b': 1, 'B': 1, 'c': 1, 'h': 2, 'H': 2, 'i': 4, 'I': 4, 'l': 4, 'L': 4, 'q': 8, 'Q': 8,
                          'e': 2, 'f': 4, 'd': 8}

    def __init__(self, buffer: bytearray, startpos=0, buflen=-1, cipher = None):
        if buffer is None:
            raise tgexception.TGException("Buffer is null")
        if buflen == -1:
            buflen = len(buffer)
        if (startpos + buflen) > len(buffer):
            raise tgexception.TGException("Start position exceeds the length")
        self.__buffer__ = buffer
        self._cipher = cipher
        self.__startpos__ = startpos
        self.__currpos__ = startpos
        self.__buflen__ = buflen
        self.__refmap__: dict = {}
        self.mark = startpos
        return

    def read(self) -> int:
        b = struct.unpack_from('b', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 1
        return b

    def skip(self, n):
        if (n < self.__startpos__) or (n >= self.__buflen__):
            raise tgexception.TGException("Invalid mark position specified")
        self.__currpos__ += n

    def readFully(self, buf: bytearray, startpos=0, buflen=-1):
        if buflen == -1:
            buflen = self.__buflen__ - self.__currpos__

        buf[startpos:startpos + buflen] = self.__buffer__[self.__currpos__:self.__currpos__ + buflen]
        # for i in range(0, buflen):
        #     buf[i+startpos] = self.readByte()
        return

    def readChars(self) -> str:
        length = self.readInt()
        ret: str = ''
        for i in range(length):
            ret += chr(self.readByte())
        return ret

    def readBytes(self) -> bytearray:
        _buflen = struct.unpack_from('>i', self.__buffer__, self.__currpos__)[0]
        _sp = self.__currpos__ + 4
        _ep = _sp + _buflen
        _buf = self.__buffer__[_sp:_ep]
        self.__currpos__ += (_buflen + 4)
        return _buf

    def readUTF(self) -> str:
        _buflen = struct.unpack_from('>h', self.__buffer__, self.__currpos__)[0]
        _sp = self.__currpos__ + 2
        _ep = _sp + _buflen
        _buf = self.__buffer__[_sp:_ep]
        self.__currpos__ += (_buflen + 2)
        s = str(_buf, "utf_8")
        return s

    def readBoolean(self):
        b = struct.unpack_from('?', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 1
        return b

    def readByte(self):
        b = struct.unpack_from('b', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 1
        return b

    def readUnsignedByte(self):
        b = struct.unpack_from('B', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 1
        return b

    def readChar(self):
        h = struct.unpack_from('>h', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 2
        return chr(h)

    def readShort(self):
        s = struct.unpack_from('>h', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 2
        return s

    def readInt(self, littleEndian: bool = False):
        i = struct.unpack_from('<i' if littleEndian else '>i', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 4
        return i

    def readUnsignedInt(self, littleEndian: bool = False):
        i = struct.unpack_from('<I' if littleEndian else '>I', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 4
        return i

    def readLong(self, littleEndian: bool = False):
        q = struct.unpack_from('<q' if littleEndian else '>q', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 8
        return q

    def readUnsignedLong(self, littleEndian: bool = False):
        q = struct.unpack_from('<Q' if littleEndian else '>Q', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 8
        return q

    def readFloat(self):
        f = struct.unpack_from('>f', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 4
        return f

    def readDouble(self):
        d = struct.unpack_from('>d', self.__buffer__, self.__currpos__)[0]
        self.__currpos__ += 8
        return d

    def readHash64From(self, startIdx, endIdx=-1) -> int:
        return tgutils.tg_checkSum64(self.__buffer__[startIdx:endIdx if endIdx > startIdx else self.__currpos__])

    def readStripedWith(self, listOfList: typing.List[typing.List], count: int,
                        typecodes: typing.Union[str, typing.List[str]]) -> bool:
        if isinstance(typecodes, typing.List) and len(typecodes) != len(listOfList):
            tglog.gLogger.log(tglog.TGLevel.Warning, "List of List's length and typecodes length don't match up!")
            return False
        for i in range(count):
            for j in range(len(listOfList)):
                lis: typing.List = listOfList[j]
                typecode: str = typecodes
                if isinstance(typecodes, typing.List):
                    typecode = typecodes[j]
                if typecode not in ProtocolDataInputStream.typecodeByteLength:
                    raise tgexception.TGException("Unknown or unimplemented typecode: %s", typecode)
                lis.append(struct.unpack_from('>' + typecode, self.__buffer__, self.__currpos__)[0])
                self.__currpos__ += ProtocolDataInputStream.typecodeByteLength[typecode]
        return True

    @property
    def available(self) -> int:
        return self.__buflen__ - self.__currpos__

    @property
    def position(self):
        return self.__currpos__

    @position.setter
    def position(self, pos):
        if (pos < self.__startpos__) or (pos >= self.__buflen__):
            raise tgexception.TGException("Invalid mark position specified")

        self.__currpos__ = pos

    @property
    def referencemap(self) -> dict:
        return self.__refmap__

    @referencemap.setter
    def referencemap(self, refmap):
        self.__refmap__ = refmap

    def readDecrypt(self) -> bytes:
        if self._cipher is None:
            raise tgexception.TGException("Trying to read encrypted data on an un-secure channel is not currently suppo"
                                          "rted.")
        length, data = self._cipher.deobfuscate(bytes(self.__buffer__[self.__currpos__:]))
        self.__currpos__ += length
        return data


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Data Input/Output Streams in PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of General Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class ProtocolMessageHeader:
    gSequenceNo = AtomicReference('Q', 0)

    def __init__(self, verbid=VerbId.InvalidMessage, authtoken=0, sessionid=0):
        self.__verbid__ = verbid
        self.__authtoken__ = authtoken
        self.__sessionid__ = sessionid
        self.__requestid__ = -1
        self.__sequenceno__ = ProtocolMessageHeader.gSequenceNo.increment()
        self.__timestamp__ = int(datetime.datetime.now().timestamp())
        self.__tenantid__ = 0
        self.__dataoffset__ = 0

    @property
    def verbid(self) -> VerbId:
        return self.__verbid__

    @property
    def sequenceno(self) -> int:
        return self.__sequenceno__

    @property
    def requestid(self) -> int:
        return self.__requestid__

    @requestid.setter
    def requestid(self, value: int):
        self.__requestid__ = value

    @property
    def timestamp(self) -> int:
        return self.__timestamp__

    @property
    def authtoken(self) -> int:
        return self.__authtoken__

    @property
    def sessionid(self) -> int:
        return self.__sessionid__

    @property
    def tenantid(self) -> int:
        return self.__tenantid__

    def updateSequenceAndTimeStamp(self):
        self.__sequenceno__ = ProtocolMessageHeader.gSequenceNo.increment()
        self.__timestamp__ = int(datetime.datetime.now().timestamp())

    def readHeader(self, instream: TGInputStream):
        size = instream.readInt()  # read the size and ignore
        magic = instream.readInt()
        if magic != TGProtocolVersion.magic():
            raise tgexception.TGBadMagic('Bad message magic')

        protver = instream.readShort()
        if not TGProtocolVersion.isCompatible(protver):
            raise tgexception.TGProtocolNotSupported('Unsupported protocol version')

        self.__verbid__ = VerbId.fromId(instream.readShort())
        self.__sequenceno__ = instream.readLong()
        self.__timestamp__ = instream.readLong()
        self.__requestid__ = instream.readLong()
        self.__authtoken__ = instream.readLong()
        self.__sessionid__ = instream.readLong()
        self.__tenantid__ = instream.readShort()
        self.__dataoffset__ = instream.readShort()

    def writeHeader(self, os: TGOutputStream):
        os.writeInt(0)  # length
        os.writeInt(TGProtocolVersion.magic())  # magic
        os.writeShort(TGProtocolVersion.getProtocolVersion())
        os.writeShort(self.verbid.value)
        os.writeLong(self.sequenceno)
        os.writeLong(self.timestamp)
        os.writeLong(self.requestid)
        os.writeLong(self.authtoken)
        os.writeLong(self.sessionid)
        os.writeShort(self.tenantid)
        os.writeShort(os.position + 2)  # Data offset


class AbstractProtocolMessage(TGMessage):

    def __new__(cls, *args, **kwargs):
        return super().__new__(cls)

    def __init__(self, verbid, authtoken=0, sessionid=0):
        self.__pduhdr__ = ProtocolMessageHeader(verbid, authtoken, sessionid)
        self.__buf__ = None

    @property
    def verbid(self) -> VerbId:
        return self.__pduhdr__.verbid

    @property
    def sequenceno(self) -> int:
        return self.__pduhdr__.sequenceno

    @property
    def requestid(self) -> int:
        return self.__pduhdr__.requestid

    @requestid.setter
    def requestid(self, value: int):
        self.__pduhdr__.requestid = value

    @property
    def timestamp(self) -> int:
        return self.__pduhdr__.timestamp

    @property
    def authtoken(self) -> int:
        return self.__pduhdr__.authtoken

    @property
    def sessionid(self) -> int:
        return self.__pduhdr__.sessionid

    @property
    def isUpdateable(self) -> bool:
        return True

    def updateSequenceAndTimeStamp(self):
        if self.isUpdateable:
            self.__pduhdr__.updateSequenceAndTimeStamp()
            self.__buf__ = None
        else:
            raise tgexception.TGException('Mutating a readonly message')

    def toBytes(self, cipher = None) -> bytearray:
        if self.__buf__ is None:
            os = ProtocolDataOutputStream(cipher)
            self.__pduhdr__.writeHeader(os)
            self.writePayload(os)
            os.writeIntAt(0, os.length)
            self.__buf__ = os.buffer

        return self.__buf__

    def fromBytes(self, buf: bytearray, cipher):
        instream: TGInputStream = ProtocolDataInputStream(buf, cipher=cipher)
        pos = instream.position
        size = instream.readInt()
        if size != len(buf):
            raise tgexception.TGInvalidMessageLength('buffer length mismatch, size was: {0}, buffer length: {1}'
                                                     .format(size, len(buf)))
        instream.position = pos
        self.__pduhdr__.readHeader(instream)
        self.readPayload(instream)

    def writeString(self, os, value: str):
        if value is None or len(value) == 0:
            os.writeBoolean(True)
        else:
            os.writeBoolean(False)
            os.writeUTF(value)
        return

    def readString(self, instream: TGInputStream) -> str:
        b = instream.readBoolean()
        return None if b else instream.readUTF()

    @abc.abstractmethod
    def writePayload(self, os: TGOutputStream):
        pass

    @abc.abstractmethod
    def readPayload(self, instream: TGInputStream):
        pass


class DefaultMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        return


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of General Purpose Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Handshake Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class HandshakeRequestType(enum.Enum):
    Invalid = 0
    Initiate = 1
    ChallengeAccepted = 2


class HandShakeRequestMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__sslmode__ = False
        self.__challenge__ = 0
        self.__requesttype__ = HandshakeRequestType.Invalid

    @property
    def sslmode(self):
        return self.__sslmode__

    @sslmode.setter
    def sslmode(self, v):
        self.__sslmode__ = v

    @property
    def challenge(self) -> int:
        return self.__challenge__

    @challenge.setter
    def challenge(self, v: int):
        self.__challenge__ = v

    @property
    def requesttype(self) -> HandshakeRequestType:
        return self.__requesttype__

    @requesttype.setter
    def requesttype(self, v: HandshakeRequestType):
        self.__requesttype__ = v

    @property
    def isUpdateable(self) -> bool:
        return True

    def writePayload(self, os: TGOutputStream):
        os.writeByte(self.__requesttype__.value)
        os.writeBoolean(self.__sslmode__)
        os.writeLong(self.__challenge__)

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Request message does not support read operation')


class HandshakeResponseStatus(enum.Enum):
    Invalid = 0
    AcceptChallenge = 1
    ProceedWithAuthentication = 2
    ChallengeFailed = 3

    @classmethod
    def fromId(cls, v):
        for status in HandshakeResponseStatus:
            if status.value == v:
                return status
        return HandshakeResponseStatus.Invalid


class HandShakeResponseMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__challenge__ = -1
        self.__status__ = HandshakeResponseStatus.Invalid
        self.__errormsg__ = ''

    @property
    def challenge(self):
        return self.__challenge__

    @property
    def status(self):
        return self.__status__

    @property
    def errormsg(self):
        return self.__errormsg__

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        self.__status__ = HandshakeResponseStatus.fromId(instream.readByte())
        self.__challenge__ = instream.readLong()
        if self.__status__ == HandshakeResponseStatus.ChallengeFailed:
            self.__errormsg__ = str(instream.readBytes())
        return


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Handshake Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Authenticate Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class AuthenticateRequestMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__clientid__ = ''
        self.__inboxaddr__ = ''
        self.__username__ = ''
        self.__dbName = ''
        self.__roles = None
        self.__password__ = []

    @property
    def clientid(self):
        return self.__clientid__

    @clientid.setter
    def clientid(self, value):
        self.__clientid__ = value

    @property
    def inboxaddr(self):
        return self.__inboxaddr__

    @inboxaddr.setter
    def inboxaddr(self, value):
        self.__inboxaddr__ = value

    @property
    def username(self):
        return self.__username__

    @username.setter
    def username(self, value):
        self.__username__ = value

    @property
    def dbname(self):
        return self.__dbname__

    @dbname.setter
    def dbname(self, value):
        self.__dbname__ = value

    @property
    def password(self):
        return self.__password__

    @password.setter
    def password(self, value: typing.Union[(bytes, str, bytearray)]):
        if isinstance(value, bytes):
            self.__password__ = value
        elif isinstance(value, str):
            self.__password__ = value.encode('utf_8')
        elif isinstance(value, bytearray):
            self.__password__ = value
        else:
            raise tgexception.TGException('Invalid type specified for password. Expected bytes, or str')

    @property
    def dbName(self) -> str:
        return self.__dbName

    @dbName.setter
    def dbName(self, val: str):
        self.__dbName = val

    @property
    def roles(self) -> typing.Optional[typing.List[str]]:
        return self.__roles

    @roles.setter
    def roles(self, val: typing.Optional[typing.List[str]]):
        self.__roles = val

    def writePayload(self, os: TGOutputStream):
        self.writeString(os, self.__dbName)
        self.writeString(os, self.__clientid__)
        self.writeString(os, self.__inboxaddr__)
        if self.__roles is None:
            os.writeInt(-1)
        else:
            os.writeInt(len(self.__roles))
            for roleName in self.__roles:
                os.writeUTF(roleName)
        self.writeString(os, self.__username__)
        os.writeBytes(self.__password__)
        return

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Request object can not be created')


class AuthenticateResponseMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__bsuccess__ = False
        self.__errorstatus__ = 0
        self.__authtoken__: int = 0
        self.__sessionid__: int = 0
        self.__certbuffer__ = []

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def errorStatus(self):
        return self.__errorstatus__

    @property
    def isSuccess(self):
        return self.__bsuccess__

    @property
    def serverCertificate(self):
        return self.__certbuffer__

    def writePayload(self, os: TGOutputStream):
        return

    @property
    def authtoken(self) -> int:
        return self.__authtoken__

    @property
    def sessionid(self) -> int:
        return self.__sessionid__

    def readPayload(self, instream: TGInputStream):
        self.__bsuccess__ = instream.readBoolean()
        if not self.__bsuccess__:
            self.__errorstatus__ = instream.readInt()
            return
        self.__authtoken__ = instream.readLong()
        self.__sessionid__ = instream.readLong()
        self.__certbuffer__ = instream.readBytes()


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Authenticate Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Connection Properties PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class ConnectionPropertiesMessage(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__props: tgutils.TGProperties = None

    @property
    def props(self) -> tgutils.TGProperties:
        return self.__props

    @props.setter
    def props(self, val: tgutils.TGProperties):
        self.__props = val

    def writePayload(self, os: TGOutputStream):
        if self.__props is None:
            raise tgexception.TGException("Can't write Null Properties!")
        props = dict()
        for key in self.__props:
            if key is not None:
                props[key] = self.__props[key]
        os.writeInt(len(props))
        for key in props:
            value = props[key]
            if isinstance(key, tgutils.ConfigName):
                os.writeUTF(key.propertyname)
            elif isinstance(key, str):
                os.writeUTF(key)
            else:
                raise tgexception.TGException("Unknown property key type: {0}!".format(str(type(key))))
            if value is None:
                os.writeBoolean(True)
            else:
                os.writeBoolean(False)
                os.writeUTF(str(value))

    def readPayload(self, instream: TGInputStream):
        pass


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Connection Properties PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Transaction Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class CommitTransactionRequestMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__addedEntities__: typing.Dict[int, tgmodel.TGEntity] = None
        self.__changeEntities__: typing.Dict[int, tgmodel.TGEntity] = None
        self.__removedEntities__: typing.Dict[int, tgmodel.TGEntity] = None
        self.__attrDiscSet__: typing.Set[tgmodel.TGAttributeDescriptor] = None

    def addCommitList(self, addedEntities: typing.Dict[int, tgentimpl.AbstractEntity],
                      changedEntities: typing.Dict[int, tgentimpl.AbstractEntity],
                      removedEntities: typing.Dict[int, tgentimpl.AbstractEntity],
                      attrDiscSet: typing.Set[tgmodel.TGAttributeDescriptor]):
        self.__addedEntities__ = addedEntities
        self.__changeEntities__ = changedEntities
        self.__removedEntities__ = removedEntities
        self.__attrDiscSet__ = attrDiscSet

    @property
    def isUpdateable(self) -> bool:
        return False

    def __checkListsDistinct(self):
        if self.__hasGT0Len(self.__addedEntities__):
            for id in self.__addedEntities__:
                if (self.__hasGT0Len(self.__changeEntities__) and id in self.__changeEntities__) or \
                        (self.__hasGT0Len(self.__removedEntities__) and id in self.__removedEntities__):
                    raise tgexception.TGTransactionGeneralErrorException(
                        "The same id appeared in multiple changelists! ID: {0}".format(id))
        if self.__hasGT0Len(self.__changeEntities__):
            for id in self.__changeEntities__:
                if (self.__hasGT0Len(self.__removedEntities__) and id in self.__removedEntities__):
                    raise tgexception.TGTransactionGeneralErrorException(
                        "The same id appeared in multiple changelists! ID: {0}".format(id))

    def __hasGT0Len(self, lis):
        return lis is not None or len(lis) > 0

    def __writeAndIgnore(self, os: TGOutputStream, ent: TGSerializable):
        try:
            ent.writeExternal(os)
        # Probably not a good idea to ignore these exceptions, but this is what the Java API does, so
        except tgexception.TGException as e:
            pass
        except IOError:
            pass

    def __writeList(self, os: TGOutputStream, lis: typing.Dict[int, tgmodel.TGEntity], code: int,
                    writeIfNew: bool = True):
        if self.__hasGT0Len(lis):
            os.writeShort(code)
            os.writeInt(sum(1 for id in lis if writeIfNew or not lis[id].isNew))
            [self.__writeAndIgnore(os, lis[id]) for id in lis if writeIfNew or not lis[id].isNew]

    '''
    // DH: TODO 3.0 introduces a transaction ID which is always -1
    '''

    def writePayload(self, os: TGOutputStream):
        startPos: int = os.position
        os.writeInt(0)  # commit buffer length
        os.writeInt(0)  # checksum - server does not validate it. DH TODO SHA 256 the byte array
        if tglog.gLogger.isEnabled(tglog.TGLevel.Debug):
            tglog.gLogger.log(tglog.TGLevel.Debug,
                        "Entering commit transaction request writePayload at output buffer position at : %d", startPos)

        # needs to be different because don't want to write over old attribute descriptors and have the transaction fail
        # because of it, unless we uglify the __writeList code above
        # Not allowed to commit attribute descriptors without admin privilege's
        # Allow various admin abilities in transaction, such as including node/edge types and attribute descriptors?
        self.__writeList(os, self.__addedEntities__, 0x1011)
        self.__writeList(os, self.__changeEntities__, 0x1012, False)
        self.__writeList(os, self.__removedEntities__, 0x1013, False)
        curPos: int = os.position
        length: int = curPos - startPos
        os.writeIntAt(startPos, length)
        os: ProtocolDataOutputStream
        hashCode = tgutils.tg_checkSum64(os.buffer[startPos + 8:])
        os.writeIntAt(startPos + 4, ctypes.c_int32((hashCode >> 32) ^ (hashCode & 0xFFFFFFFF)).value)
        return

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class CommitTransactionResponseMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__exception: tgexception.TGTransactionException = None
        self.__instream: TGInputStream = None

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    @property
    def exception(self):
        return self.__exception

    def __processTxnStatus(self, instream: TGInputStream, status: int):
        resp: tgexception.TGTransactionResponse = tgexception.TGTransactionResponse.fromId(status)
        ret = None
        if resp is not tgexception.TGTransactionResponse.TGTransactionSuccess:
            msg: str = None
            try:
                msg = instream.readUTF()
            except:
                msg = "Couldn't read error message."
            ret = tgexception.TGTransactionException.buildTransactionException(msg, resp)
        return ret

    def readPayload(self, instream: TGInputStream):
        instream.readInt()  # Buf length not needed, Python helps with that memory management
        instream.readInt()  # Unused checksum
        self.__exception: tgexception.TGTransactionException = self.__processTxnStatus(instream, instream.readInt())
        if self.__exception is None:
            self.__instream = instream

    def finishReadWith(self, addedEntities: typing.Dict[int, tgmodel.TGEntity],
                       changedEntities: typing.Dict[int, tgmodel.TGEntity],
                       removedEntities: typing.Dict[int, tgmodel.TGEntity],
                       typeReg: tggmdimpl.TypeRegistry):
        while 6 <= self.__instream.available:
            opCode: int = self.__instream.readShort()
            position: int = self.__instream.position
            count: int = self.__instream.readInt()
            vids: typing.List[int] = []
            ids: typing.List[int] = []
            vers: typing.List[int] = []
            if opCode == 0x1010:
                self.__instream.readStripedWith([vids, ids], count, 'i')
                for i in range(count):
                    if vids[i] in addedEntities:
                        typeReg[vids[i]].id = ids[i]
                        typeReg.addItem(typeReg[vids[i]])
                    else:
                        tglog.gLogger.log(tglog.TGLevel.Warning, "Could not find attributeDescriptor ID: %d", vids[i])
            elif opCode == 0x1011:
                self.__instream.readStripedWith([vids, ids, vers], count, 'q')
                for i in range(count):
                    if vids[i] in addedEntities:
                        addedEntities[vids[i]].entityId = ids[i]
                        addedEntities[vids[i]].version = vers[i]
                    else:
                        tglog.gLogger.log(tglog.TGLevel.Warning, "Could not find entity ID: %d in addedList", vids[i])
            elif opCode == 0x1012:
                self.__instream.readStripedWith([ids, vers], count, 'q')
                for i in range(count):
                    if ids[i] in changedEntities:
                        changedEntities[ids[i]].version = vers[i]
                    else:
                        tglog.gLogger.log(tglog.TGLevel.Warning, "Could not find entity ID: %d in changeList", ids[i])
            elif opCode == 0x1013:
                for i in range(count):
                    eid: int = self.__instream.readLong()
                    if eid in removedEntities:
                        removedEntities[eid].markDeleted()
                    else:
                        tglog.gLogger.log(tglog.TGLevel.Warning, "Could not find entity ID: %d in removedList", eid)
            elif opCode == 0x6789:
                self.__instream.position = position
                tglog.gLogger.log(tglog.TGLevel.Debug, "Received %d debug entities", count)
            else:
                tglog.gLogger.log(tglog.TGLevel.Warning, "Unknown OP Code: %x", opCode)


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Transaction Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''
'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Begin Bulk Import Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class BeginImportSessionRequest(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.isBatch: bool = True
        self.totalRequsts: int = None
        self.numFiles: int = None
        self.__loadopt: tgbulk.TGLoadOptions = None
        self.__erroropt: tgbulk.TGErrorOptions = None
        self.__dtformat: tgbulk.TGDateFormat = None

    @property
    def loadopt(self) -> tgbulk.TGLoadOptions:
        return self.__loadopt

    @loadopt.setter
    def loadopt(self, val: tgbulk.TGLoadOptions):
        self.__loadopt = val

    @property
    def erroropt(self) -> tgbulk.TGErrorOptions:
        return self.__erroropt

    @erroropt.setter
    def erroropt(self, val: tgbulk.TGErrorOptions):
        self.__erroropt = val

    @property
    def dtformat(self) -> tgbulk.TGDateFormat:
        return self.__dtformat

    @dtformat.setter
    def dtformat(self, val: tgbulk.TGDateFormat):
        self.__dtformat = val

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        os.writeBoolean(self.isBatch)    # Need to send this because this is only ever a batch import session, not file
                                         # based
        if self.isBatch:
            os.writeUTF(self.__erroropt.value)
            os.writeUTF(self.__loadopt.value)
            os.writeUTF(self.__dtformat.value)
        else:
            os.writeInt(self.totalRequsts)
            os.writeInt(self.numFiles)

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class BeginImportSessionResponse(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__error = None

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def error(self) -> typing.Optional[tgexception.TGException]:
        return self.__error

    def readPayload(self, instream: TGInputStream):
        result = instream.readInt()
        if result != 0:
            msg: str = "Error starting import"
            # No error message returned. Should there be one??? TODO Ask Katie
            """
            try:
                msg = instream.readUTF()
            except:
                msg = "Couldn't read error message."
            """
            self.__error = tgexception.TGImpExpException(msg, errorcode=result)

    def writePayload(self, os: TGOutputStream):
        return


class PartialImportRequest(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.isBatch = True
        self.fileIdx: int = None
        self.fileName: str = None
        self.fileTotalReq: int = None
        self.fileReqIdx: int = None
        self.__data: bytes = None
        self.__type: tgmodel.TGEntityType = None
        self.__reqIdx: int = None
        self.__totalRequestsForType: int = None
        self.__attrList: typing.List[str] = None        # Would prefer to have list of attr desc, but id, from, and to
                                                        # columns have no corresponding attribute descriptors
    @property
    def data(self) -> bytes:
        return self.__data

    @data.setter
    def data(self, val: typing.Union[bytes, str]):
        if isinstance(val, str):
            val = val.encode('utf_8')
        self.__data = val

    @property
    def type(self) -> tgmodel.TGEntityType:
        return self.__type

    @type.setter
    def type(self, val: tgmodel.TGEntityType):
        self.__type = val

    @property
    def reqIdx(self) -> int:
        return self.__reqIdx

    @reqIdx.setter
    def reqIdx(self, val: int):
        self.__reqIdx = val

    @property
    def totalRequestsForType(self) -> int:
        return self.__totalRequestsForType

    @totalRequestsForType.setter
    def totalRequestsForType(self, val: int):
        self.__totalRequestsForType = val

    @property
    def attrList(self) -> typing.List[str]:
        return self.__attrList

    @attrList.setter
    def attrList(self, val: typing.List[str]):
        self.__attrList = val

    @property
    def isUpdateable(self) -> bool:
        return False

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException("Write only message.")

    def writePayload(self, os: ProtocolDataOutputStream):
        os.writeBoolean(self.isBatch)       # Will always be batch-style and not file-style import session for Python
        buf_len = len(self.__data)
        os.writeUnsignedInt(buf_len)
        cs_pos = os.position
        os.writeLong(0)

        if self.isBatch:
            # Is a batch style

            os.writeUTF(self.__type.name)
            os.writeInt(self.__reqIdx)
            os.writeInt(self.__totalRequestsForType)

            if self.__reqIdx == 0:
                num_cols = len(self.attrList)
                os.writeInt(num_cols)
                for i in range(num_cols):
                    os.writeChars(self.attrList[i])

        else:
            # Is not a batch style
            os.writeUTF(self.fileName)
            os.writeInt(self.fileTotalReq)
            os.writeInt(self.fileReqIdx)

        data_start = os.position
        # Very bad, but fixes another bug. This mimics writeUTF(), but does it with bytes instead of strings.
        os.writeShort(buf_len)
        os._buffer.extend(bytes(self.__data))
        os._pos += buf_len
        os.writeHash64At(cs_pos, data_start + 2)            # Need to account for the offset caused by the length of the
                                                            # string.


class PartialImportResponse(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__resultList: typing.List[tgbulk.TGImportDescriptor] = None
        self.__error = None

    @property
    def error(self) -> typing.Optional[tgexception.TGException]:
        return self.__error

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def resultList(self) -> typing.List[tgadm.TGImportDescriptor]:
        return self.__resultList

    def readPayload(self, instream: TGInputStream):
        result = instream.readInt()
        if result != 0:
            msg: str = "Error running import"
            # No error message returned. Should there be one??? TODO Ask Katie
            """
            try:
                msg = instream.readUTF()
            except:
                msg = "Couldn't read error message."
            """
            self.__error = tgexception.TGImpExpException(msg, result)
        else:
            import tgdb.impl.adminimpl as tgadmimpl
            self.__resultList = []
            numResults = instream.readLong()
            typename: str
            isnode: bool
            numInstances: int
            for i in range(numResults):
                typename = instream.readUTF()               # Only unpacking them in this order to preserve
                isnode = instream.readBoolean()
                numInstances = instream.readLong()
                self.__resultList.append(tgadmimpl.ImportDescImpl(typename, isnode, numInstances))

    def writePayload(self, os: TGOutputStream):
        return


class EndBulkImportSessionRequest(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        # TODO write an actual payload request?
        pass

    def readPayload(self, instream: TGInputStream):
        return


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End Bulk Import Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''
'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Begin Bulk Export Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class BeginExportRequest(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__isBatch: bool = None
        self.__zipName: str = None
        self.__numEntities: int = None

    @property
    def isBatch(self) -> bool:
        return self.__isBatch

    @isBatch.setter
    def isBatch(self, val: bool):
        self.__isBatch = val

    @property
    def zipName(self) -> str:
        return self.__zipName

    @zipName.setter
    def zipName(self, val: str):
        self.__zipName = val

    @property
    def maxBatchEntities(self) -> int:
        return self.__numEntities

    @maxBatchEntities.setter
    def maxBatchEntities(self, val: int):
        self.__numEntities = val

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        os.writeBoolean(self.__isBatch)
        if self.__isBatch:
            os.writeInt(self.__numEntities)
        else:
            os.writeUTF("" if self.__zipName is None else self.__zipName)

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class BeginExportResponse(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__numRequests__: int = None
        self.__typeList__: typing.List[typing.Tuple[str, bool, int]] = None
        self.__error = None

    @property
    def error(self) -> typing.Optional[tgexception.TGException]:
        return self.__error

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def numRequests(self) -> int:
        return self.__numRequests__

    @property
    def typeList(self) -> typing.List[typing.Tuple[str, bool, int]]:
        return self.__typeList__

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        self.__numRequests__ = instream.readInt()
        result = instream.readInt()
        if result != 0:
            msg: str = "Error starting import."
            # No error message returned. Should there be one??? TODO Ask Katie
            """
            try:
                msg = instream.readUTF()
            except:
                msg = "Couldn't read error message."
            """
            self.__error = tgexception.TGImpExpException(msg, result)
        else:
            #  if bulk export returns success, but numTypes = 0, then the database was empty (don't send more requests)
            numTypes: int = instream.readLong()

            self.__typeList__ = []
            for i in range(numTypes):
                s = instream.readUTF()  # Type name
                b = instream.readBoolean()  # Is Node
                l = instream.readLong()  # Number of entities
                self.__typeList__.append((s, b, l,))


class PartialExportRequest(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__reqNum__: int = None

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def requestNum(self) -> int:
        return self.__reqNum__

    @requestNum.setter
    def requestNum(self, val: int):
        self.__reqNum__ = val

    def writePayload(self, os: TGOutputStream):
        os.writeInt(self.requestNum)

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class PartialExportResponse(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__fileName__: str = None
        self.__newType: bool = None
        self.__attrList: typing.List[str] = None
        self.__typeName: str = None
        self.__numEntities: int = None
        self.__hasMore: bool = None
        self.__data__: bytes = None

    @property
    def newType(self) -> bool:
        return self.__newType

    @property
    def attrList(self) -> typing.List[str]:
        return self.__attrList

    @property
    def typeName(self) -> str:
        return self.__typeName

    @property
    def numEntities(self) -> int:
        return self.__numEntities

    @property
    def hasMore(self) -> bool:
        return self.__hasMore

    @property
    def fileName(self) -> str:
        return self.__fileName__

    @property
    def data(self) -> bytes:
        return self.__data__

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        isBatch = instream.readBoolean()  # isBatch TODO DH
        if isBatch is False:
            self.__fileName__ = instream.readUTF()
        else:
            self.__newType = instream.readBoolean()
            if self.__newType:
                self.__typeName = instream.readUTF()
                numAttrs = instream.readInt()
                self.__attrList = []
                for i in range(numAttrs):
                    self.__attrList.append(instream.readChars())
            self.__numEntities = instream.readInt()
            self.__hasMore = instream.readBoolean()
        _ = instream.readInt()  # Buffer length
        _ = instream.readLong()  # Checksum
        self.__data__ = bytes(instream.readBytes())


"""
class BeginBulkExportSessionRequest(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class BeginBulkExportSessionResponse(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    @property
    def isUpdateable(self) -> bool:
        return False

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)

    def writePayload(self, os: TGOutputStream):
        raise tgexception.TGException('Read only response message.')

    def readPayload(self, instream: TGInputStream):
        result: tgexception.TGBulkImpExpResponse = tgexception.TGBulkImpExpResponse.fromId(instream.readInt())
        if result != tgexception.TGBulkImpExpResponse.TGBulkImpExpSuccess:
            msg: str
            try:
                msg = instream.readUTF()
            except:
                msg = "Couldn't read message!"
            raise tgexception.TGBulkImpExpException.buildBulkImportException(msg, result)


class BeginBatchExportEntityRequest(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__entKind__: tgmodel.TGEntityKind = None
        self.__entType__: tgmodel.TGEntityType = None
        self.__batchSize__: int = None

    @property
    def entKind(self) -> tgmodel.TGEntityKind:
        return self.__entKind__

    @entKind.setter
    def entKind(self, val: tgmodel.TGEntityKind):
        self.__entKind__ = val

    @property
    def entType(self) -> tgmodel.TGEntityType:
        return self.__entType__

    @entType.setter
    def entType(self, val: tgmodel.TGEntityType):
        self.__entType__ = val

    @property
    def batchSize(self) -> int:
        return self.__batchSize__

    @batchSize.setter
    def batchSize(self, val: int):
        self.__batchSize__ = val

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        os.writeByte(self.__entKind__.value)
        os.writeInt(self.__entType__.id)
        os.writeInt(self.__batchSize__)

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class BeginBatchExportEntityResponse(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__desc__: int = None
        self.__colLabels__: List[str] = None

    @property
    def descriptor(self) -> int:
        return self.__desc__

    @property
    def columnLabels(self) -> List[str]:
        return self.__colLabels__

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        raise tgexception.TGException('Read only response message.')

    def readPayload(self, instream: TGInputStream):
        result: tgexception.TGBulkImpExpResponse = tgexception.TGBulkImpExpResponse.fromId(instream.readInt())
        if result == tgexception.TGBulkImpExpResponse.TGBulkImpExpSuccess:
            self.__desc__ = instream.readInt()
            numLabels = instream.readInt()
            self.__colLabels__ = []
            for i in range(numLabels):
                self.__colLabels__.append(instream.readUTF())
        else:
            msg: str
            try:
                msg = instream.readUTF()
            except:
                msg = "Couldn't read message!"
            raise tgexception.TGBulkImpExpException.buildBulkImportException(msg, result)


class SingleBatchExportEntityRequest(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__desc__: int = None


    @property
    def descriptor(self) -> int:
        return self.__desc__

    @descriptor.setter
    def descriptor(self, val: int):
        self.__desc__ = val

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        os.writeInt(self.__desc__)

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class SingleBatchExportEntityResponse(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__numEnts__: int = None
        self.__data__: str = None
        self.__hasMore__: bool = None

    @property
    def data(self) -> str:
        return self.__data__

    @property
    def numEnts(self) -> int:
        return self.__numEnts__

    @property
    def hasMore(self) -> bool:
        return self.__hasMore__

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        raise tgexception.TGException('Read only response message.')

    def readPayload(self, instream: TGInputStream):
        result: tgexception.TGBulkImpExpResponse = tgexception.TGBulkImpExpResponse.fromId(instream.readInt())
        if result == tgexception.TGBulkImpExpResponse.TGBulkImpExpSuccess:
            self.__hasMore__ = instream.readBoolean()
            self.__numEnts__ = instream.readInt()
            self.__data__ = str(instream.readBytes(), "utf_8")
        else:
            msg: str
            try:
                msg = instream.readUTF()
            except:
                msg = "Couldn't read message!"
            raise tgexception.TGBulkImpExpException.buildBulkImportException(msg, result)


class EndBulkExportSessionRequest(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        # TODO find if there is anything else that we want to send to the server?
        return

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class EndBulkExportSessionResponse(AbstractProtocolMessage):
    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    @property
    def isUpdateable(self) -> bool:
        return False

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)

    def writePayload(self, os: TGOutputStream):
        raise tgexception.TGException('Read only response message.')

    def readPayload(self, instream: TGInputStream):
        result: tgexception.TGBulkImpExpResponse = tgexception.TGBulkImpExpResponse.fromId(instream.readInt())
        if result != tgexception.TGBulkImpExpResponse.TGBulkImpExpSuccess:
            msg: str
            try:
                msg = instream.readUTF()
            except:
                msg = "Couldn't read message!"
            raise tgexception.TGBulkImpExpException.buildBulkImportException(msg, result)
"""

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End Bulk Export Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''
'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Query Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class QueryRequestMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__command__: tgquery.TGQueryCommand = tgquery.TGQueryCommand.Invalid
        self.__option__: tgquery.TGQueryOption = tgquery.TGQueryOption()
        self.__queryHashId__: int = -1

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def command(self) -> tgquery.TGQueryCommand:
        return self.__command__

    @command.setter
    def command(self, value: tgquery.TGQueryCommand):
        self.__command__ = value

    @property
    def option(self) -> tgquery.TGQueryOption:
        return self.__option__

    @option.setter
    def option(self, opt: tgquery.TGQueryOption):
        self.__option__ = opt

    @property
    def query(self):
        return self.__option__.queryExpr

    @query.setter
    def query(self, value: str):
        self.__option__.queryExpr = value

    @property
    def queryHashId(self) -> int:
        return self.__queryHashId__

    @queryHashId.setter
    def queryHashId(self, value: int):
        self.__queryHashId__ = value

    def writePayload(self, os: TGOutputStream):
        startPos: int = os.position
        os.writeInt(0)  # Write num bytes later
        os.writeInt(0)  # Write checksum
        os.writeInt(self.__command__.value)
        if self.__command__ is tgquery.TGQueryCommand.Create or self.__command__ is tgquery.TGQueryCommand.Execute or \
                self.__command__ is tgquery.TGQueryCommand.ExecuteGremlin or self.__command__ is tgquery.TGQueryCommand.ExecuteGremlinStr:
            self.option.writeExternal(os, writeSort=True, writeQueryStr=True)
        elif self.__command__ is tgquery.TGQueryCommand.ExecuteID or self.__command__ is tgquery.TGQueryCommand.Close:
            self.option.writeExternal(os, writeSort=True)
            if self.__queryHashId__ < 0:
                raise tgexception.TGException("Must set a query hash ID if command is execute id or close!")
            os.writeLong(self.queryHashId)
        else:
            raise tgexception.TGException("Invalid command: {0}".format(self.__command__))
        os.writeIntAt(startPos, os.position - startPos)
        os: ProtocolDataOutputStream
        hashCode = tgutils.tg_checkSum64(os.buffer[startPos + 8:])
        hashCode = ((hashCode >> 32) ^ hashCode) & 0xFFFFFFFF
        hashCode = ctypes.c_int32(hashCode).value
        os.writeIntAt(startPos + 4, hashCode)

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class QueryResponseMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__result__: int = None
        self.__queryHashId__: int = None
        self.__resultCount__: int = None
        self.__totalCount__: int = None
        self.__entityStream__: ProtocolDataInputStream = None
        self.__fetchedEntites__: typing.Dict[int, tgmodel.TGEntity] = None
        self.__error: typing.Optional[tgquery.TGQueryException] = None
        self.__resultTypeAnnot: str = None

    @property
    def error(self) -> typing.Optional[tgquery.TGQueryException]:
        return self.__error

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def result(self) -> int:
        if self.__result__ is None:
            raise tgexception.TGException("No result read yet.")
        return self.__result__

    @property
    def resultCount(self) -> int:
        if self.__resultCount__ is None:
            raise tgexception.TGException("No result count read yet.")
        return self.__resultCount__

    @property
    def totalCount(self) -> int:
        if self.__totalCount__ is None and self.__resultCount__ is None:
            raise tgexception.TGException("No total count read yet.")
        elif self.__totalCount__ is None:
            return self.__resultCount__
        return self.__totalCount__

    @property
    def hasResult(self) -> bool:
        if self.__resultCount__ is None:
            return False
        return self.__resultCount__ > 0

    @property
    def queryHashId(self) -> int:
        if self.__queryHashId__ is None:
            raise tgexception.TGException("No query hash id read yet.")
        return self.__queryHashId__

    @property
    def resultTypeAnnot(self) -> str:
        return self.__resultTypeAnnot

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        count: int = instream.readInt()  # Unused count, could confirm that byte length is same
        checkSum: int = instream.readInt()  # Could check checksum
        instream: ProtocolDataInputStream
        actualCheckSum: int = tgutils.tg_checkSum64(instream.__buffer__[instream.position:instream.position + count])
        actualCheckSum = (actualCheckSum >> 32) ^ (actualCheckSum & 0xFFFFFFFF)
        # if checkSum != actualCheckSum:
        #     raise tgexception.TGException("Check sums do not match!")
        self.__result__ = instream.readInt()
        if self.__result__ != 0:
            queryError: typing.Optional[str] = None
            try:
                queryError = instream.readUTF()
                if queryError == '':
                    queryError = None
            except:
                pass
            self.__error = tgquery.TGQueryException(self.__result__, queryError)
            return
        self.__queryHashId__ = instream.readLong()
        syntax: int = instream.readByte()
        self.__resultTypeAnnot = instream.readUTF()
        self.__resultCount__ = instream.readInt()
        self.__entityStream__ = instream
        if syntax == 1:
            self.__totalCount__ = instream.readInt()
            tglog.gLogger.log(tglog.TGLevel.Debug, "Query has {0} result entities and {1} total entities".format(
                self.__resultCount__, self.__totalCount__))
        else:
            tglog.gLogger.log(tglog.TGLevel.Debug, "Query has {0} result count".format(self.__resultCount__))

    def finishReadWith(self, command: tgquery.TGQueryCommand, gof: tgmodel.TGGraphObjectFactory):
        ret = tgqueryimpl.ObjectResultSet(gof.connection, self.requestid)
        if self.__resultCount__ > 0:
            if command == tgquery.TGQueryCommand.ExecuteGremlinStr or command == tgquery.TGQueryCommand.ExecuteGremlin:
                elemType: tgqueryimpl.GremlinElementType = tgqueryimpl.GremlinElementType.fromId(
                    self.__entityStream__.readByte())
                if elemType is tgqueryimpl.GremlinElementType.List:
                    tglog.gLogger.log(tglog.TGLevel.DebugWire, "Finishing read query response...")
                    ret = tgqueryimpl.ObjectResultSet(gof.connection, self.requestid)
                    ret.setResultTypeAnnotation(self.__resultTypeAnnot)
                    ret.messageBytes = bytes(self.__entityStream__.__buffer__[self.__entityStream__.__currpos__:])
                else:
                    tglog.gLogger.log(tglog.TGLevel.Warning, "Unsupported Gremlin element type: %s", str(elemType))
            elif command == tgquery.TGQueryCommand.Execute:
                ret = tgqueryimpl.ResultSetImpl.parseEntityStream(self.__entityStream__, gof, self.requestid,
                                                                  self.totalCount, self.resultCount)
                self.__fetchedEntites__ = self.__entityStream__.referencemap
            else:
                raise tgexception.TGException("Unsupported command of finishReadWith: %s", str(command))
        return ret


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Query Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''
'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Metadata Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class MetadataRequestMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        return


class MetadataResponseMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__registry__ = tggmdimpl.TypeRegistry()

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        count = instream.readInt()
        while count > 0:
            systype = instream.readByte()
            typecount = instream.readInt()
            if systype == tgmodel.TGSystemType.AttributeDescriptor.value:
                for j in range(0, typecount):
                    attrdesc = tgattrimpl.AttributeDescriptorImpl()
                    attrdesc.readExternal(instream)
                    self.__registry__.addItem(attrdesc)
            elif systype == tgmodel.TGSystemType.NodeType.value:
                for j in range(0, typecount):
                    nodetype = tggmdimpl.NodeTypeImpl(self.__registry__)
                    nodetype.readExternal(instream)
                    self.__registry__.addItem(nodetype)
            elif systype == tgmodel.TGSystemType.EdgeType.value:
                for j in range(0, typecount):
                    edgetype = tggmdimpl.EdgeTypeImpl(self.__registry__)
                    edgetype.readExternal(instream)
                    self.__registry__.addItem(edgetype)
            else:
                tglog.gLogger.log(tglog.TGLevel.Error,
                            'Invalid type {0} encountered'.format(tgmodel.TGSystemType.fromValue(systype)))
            count -= typecount

    @property
    def typeregistry(self):
        return self.__registry__


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Metadata Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Get Entity Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class GetEntityCommand(enum.Enum):
    GetEntity = 0
    GetByID = 1
    GetMultiples = 2
    Continue = 10
    Close = 20


class GetEntityRequestMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__key__: tgmodel.TGKey = None
        self.__command__: GetEntityCommand = GetEntityCommand.GetEntity
        self.__option__: tgquery.TGQueryOption = tgquery.TGQueryOption()
        self.__resultId__: int = 0

    @property
    def resultId(self) -> int:
        return self.__resultId__

    @resultId.setter
    def resultId(self, value: int):
        self.__resultId__ = value

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def key(self) -> tgmodel.TGKey:
        return self.__key__

    @key.setter
    def key(self, value: tgmodel.TGKey):
        self.__key__ = value

    @property
    def command(self) -> GetEntityCommand:
        return self.__command__

    @command.setter
    def command(self, cmd):
        if not cmd in (GetEntityCommand.GetEntity, GetEntityCommand.GetByID, GetEntityCommand.GetMultiples):
            raise tgexception.TGException('Invalid command specified. It has to be one of value QueryCommand.GetEntity,'
                                          ' QueryCommand.CreateQuery, QueryCommand.ExecuteQuery ')
        self.__command__ = cmd

    @property
    def option(self):
        return self.__option__

    def __fixUpOptions(self):
        if self.__option__.batchSize < 10 or self.__option__.batchSize > 32767:
            self.__option__.batchSize = 50
        if self.__option__.edgeLimit < 0 or self.__option__.edgeLimit > 32767:
            self.__option__.edgeLimit = 1000
        if self.__option__.prefetchSize < 0:
            self.__option__.prefetchSize = 1000
        if self.__option__.traversalDepth < 1 or self.__option__.traversalDepth > 1000:
            self.__option__.traversalDepth = 3

    @option.setter
    def option(self, opt):
        self.__option__ = opt

    def writePayload(self, os: TGOutputStream):
        self.__fixUpOptions()
        os.writeShort(self.__command__.value)
        os.writeInt(self.__resultId__)  # self.resultId is always 0 - Not sure what Java is trying to say
        if self.__command__ == GetEntityCommand.GetEntity or self.__command__ == GetEntityCommand.GetByID or \
                self.__command__ == GetEntityCommand.GetMultiples:
            self.__option__.writeExternal(os)
        startPos: int = os.position
        os.writeInt(0) # key buffer length
        keyPos: int = os.position
        self.__key__.writeExternal(os)
        curPos: int = os.position
        length: int = curPos - keyPos
        os.writeIntAt(startPos, length)

    def readPayload(self, instream: TGInputStream):
        raise tgexception.TGException('Write only request message.')


class GetEntityResponseMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__entityStream__: TGInputStream = None
        self.__totalCount__: int = 0
        self.__resultCount__: int = 0  # Do we need resultCount???
        self.__resultId__: int = 0
        self.__fetchedEntities__: typing.Dict[int, tgmodel.TGEntity] = None

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def hasResult(self) -> bool:
        return self.__totalCount__ > 0

    @property
    def totalCount(self) -> int:
        return self.__totalCount__

    # Do we need resultCount?
    @property
    def resultCount(self) -> int:
        return self.__resultCount__

    @property
    def resultId(self) -> int:
        return self.__resultId__

    @property
    def entityStream(self) -> TGInputStream:
        return self.__entityStream__

    @property
    def fetchedEntities(self) -> typing.Dict[int, tgmodel.TGEntity]:
        return self.__fetchedEntities__

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        self.__entityStream__ = instream
        self.__resultId__ = instream.readInt()
        position = instream.position
        self.__totalCount__ = instream.readInt()

    def finishReadWith(self, gof: tgmodel.TGGraphObjectFactory):
        tgqueryimpl.ResultSetImpl.parseEntityStream(self.__entityStream__, gof, self.requestid, self.__totalCount__)
        self.__fetchedEntities__ = self.__entityStream__.referencemap


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Get Entity Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Get Large Object Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class GetLargeObjectRequestMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__entityId: int = None
        self.__decrypt: bool = None

    @property
    def isUpdateable(self) -> bool:
        return False

    @property
    def entityId(self) -> int:
        return self.__entityId

    @entityId.setter
    def entityId(self, val: int):
        self.__entityId = val

    @property
    def decrypt(self) -> bool:
        return self.__decrypt

    @decrypt.setter
    def decrypt(self, val: bool):
        self.__decrypt = val

    def writePayload(self, os: TGOutputStream):
        os.writeLong(self.__entityId)
        os.writeBoolean(self.__decrypt)

    def readPayload(self, instream: TGInputStream):
        return


class GetLargeObjectResponseMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)
        self.__entityId: int = None
        self.__data: bytearray = None
        self.__baseData: bytearray = None
        self.__cipher = None
        
    @property
    def entityId(self) -> int:
        return self.__entityId

    @property
    def data(self) -> bytes:
        return bytes(self.__data)

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        result = instream.readInt()
        if result == 0:                                 # Result is TGSuccess
            self.__entityId = instream.readLong()
            hasData = instream.readBoolean()
            if hasData:
                self.__data = instream.readBytes()
        else:
            raise tgexception.TGException("Could not retrieve large object from server.")

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Get Large Object Request/Response PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    Start of Session Exception PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''


class SessionExceptionMessage(AbstractProtocolMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid=VerbId.ExceptionMessage, authtoken=0, sessionid=0):
        super().__init__(verbid, authtoken, sessionid)
        self.__errormsg__: typing.Optional[str] = None
        self.__exceptionType: tgexception.ExceptionType = tgexception.ExceptionType.GeneralException

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        code = instream.readInt()
        if code != 0:
            isNull = instream.readBoolean()
            if not isNull:
                self.__errormsg__ = instream.readUTF()
            # TODO Map into an exception type.

    @property
    def exceptionType(self) -> tgexception.ExceptionType:
        return self.__exceptionType

    @exceptionType.setter
    def exceptionType(self, v: tgexception.ExceptionType):
        self.__exceptionType = v

    @property
    def errormsg(self) -> typing.Optional[str]:
        return self.__errormsg__

    @errormsg.setter
    def errormsg(self, v: typing.Optional[str]):
        self.__errormsg__ = v


class SessionForcefullyTerminatedMessage(SessionExceptionMessage):

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls, *args, **kwargs)
        instance.__init__(*args, **kwargs)
        return instance

    def __init__(self, verbid, authtoken, sessionid):
        super().__init__(verbid, authtoken, sessionid)

    @property
    def isUpdateable(self) -> bool:
        return False

    def writePayload(self, os: TGOutputStream):
        return

    def readPayload(self, instream: TGInputStream):
        return


'''
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    End of Session Exception PDU
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
'''

class TGMessageFactory:
    __binitialized__ = False
    __registry__ = None

    @classmethod
    def initialize(cls):
        import tgdb.impl.adminimpl as tgadmimpl

        if cls.__binitialized__:
            return

        registry = dict()

        registry[VerbId.PingMessage] = DefaultMessage

        registry[VerbId.HandShakeRequest] = HandShakeRequestMessage
        registry[VerbId.HandShakeResponse] = HandShakeResponseMessage

        registry[VerbId.AuthenticateRequest] = AuthenticateRequestMessage
        registry[VerbId.AuthenticateResponse] = AuthenticateResponseMessage

        registry[VerbId.BeginTransactionRequest] = DefaultMessage
        registry[VerbId.BeginTransactionResponse] = DefaultMessage
        registry[VerbId.CommitTransactionRequest] = CommitTransactionRequestMessage
        registry[VerbId.CommitTransactionResponse] = CommitTransactionResponseMessage
        registry[VerbId.RollbackTransactionRequest] = DefaultMessage
        registry[VerbId.RollbackTransactionResponse] = DefaultMessage

        registry[VerbId.QueryRequest] = QueryRequestMessage
        registry[VerbId.QueryResponse] = QueryResponseMessage

        registry[VerbId.TraverseRequest] = DefaultMessage
        registry[VerbId.TraverseResponse] = DefaultMessage

        registry[VerbId.AdminRequest] = tgadmimpl.AdminRequest
        registry[VerbId.AdminResponse] = tgadmimpl.AdminResponse

        registry[VerbId.MetadataRequest] = MetadataRequestMessage
        registry[VerbId.MetadataResponse] = MetadataResponseMessage

        registry[VerbId.GetEntityRequest] = GetEntityRequestMessage
        registry[VerbId.GetEntityResponse] = GetEntityResponseMessage

        registry[VerbId.GetLargeObjectRequest] = GetLargeObjectRequestMessage
        registry[VerbId.GetLargeObjectResponse] = GetLargeObjectResponseMessage

        registry[VerbId.BeginExportRequest] = BeginExportRequest
        registry[VerbId.BeginExportResponse] = BeginExportResponse
        registry[VerbId.PartialExportRequest] = PartialExportRequest
        registry[VerbId.PartialExportResponse] = PartialExportResponse
        registry[VerbId.CancelExportRequest] = DefaultMessage
        registry[VerbId.BeginImportRequest] = DefaultMessage
        registry[VerbId.BeginImportResponse] = DefaultMessage
        registry[VerbId.PartialImportRequest] = DefaultMessage
        registry[VerbId.PartialImportResponse] = DefaultMessage

        registry[VerbId.DumpStacktraceRequest] = DefaultMessage

        registry[VerbId.DisconnectChannelRequest] = DefaultMessage
        registry[VerbId.SessionForcefullyTerminated] = SessionForcefullyTerminatedMessage

        registry[VerbId.DecryptBufferRequest] = DefaultMessage
        registry[VerbId.DecryptBufferResponse] = DefaultMessage

        registry[VerbId.ExceptionMessage] = SessionExceptionMessage

        registry[VerbId.ConnectionPropertiesMessage] = ConnectionPropertiesMessage

        registry[VerbId.BeginImportRequest] = BeginImportSessionRequest
        registry[VerbId.BeginImportResponse] = BeginImportSessionResponse
        registry[VerbId.PartialImportRequest] = PartialImportRequest
        registry[VerbId.PartialImportResponse] = PartialImportResponse
        registry[VerbId.EndImportRequest] = EndBulkImportSessionRequest
        registry[VerbId.EndImportResponse] = PartialImportResponse
        """
        registry[VerbId.BeginBulkExportSessionRequest] = BeginBulkExportSessionRequest
        registry[VerbId.BeginBulkExportSessionResponse] = BeginBulkExportSessionResponse
        registry[VerbId.BeginBulkExportEntityRequest] = BeginBatchExportEntityRequest
        registry[VerbId.BeginBatchExportEntityRequest] = BeginBatchExportEntityResponse
        registry[VerbId.SingleBulkExportEntityRequest] = SingleBatchExportEntityRequest
        registry[VerbId.SingleBatchExportEntityRequest] = SingleBatchExportEntityResponse
        registry[VerbId.EndBulkExportSessionRequest] = EndBulkExportSessionRequest
        registry[VerbId.EndBulkExportSessionResponse] = EndBulkExportSessionResponse
        """
        cls.__registry__ = registry
        cls.__binitialized__ = True
        return

    @classmethod
    def createMessage(cls, verborbytes: typing.Union[VerbId, bytes], authtoken=0, sessionid=0, cipher=None, offset=0,
                      buflen=0) -> TGMessage:

        cls.initialize()

        if isinstance(verborbytes, VerbId):
            verbid: VerbId = verborbytes
            return cls.__createMessageFromVerb__(verbid, authtoken, sessionid)
        elif isinstance(verborbytes, bytes):
            buf: bytes = verborbytes
            blen = buflen if buflen > 0 else len(buf)
            return cls.__createMessageFromBuffer__(buf, offset, blen, cipher)
        else:
            raise TypeError('Invalid Argument specified for message creation')

    @classmethod
    def __createMessageFromVerb__(cls, verbid, authtoken, sessionid) -> TGMessage:

        cls.initialize()

        clz = cls.__registry__[verbid]
        if clz is None:
            raise tgexception.TGBadVerb('Bad Verb')

        obj = clz.__new__(clz, verbid, authtoken, sessionid)  # new will call __init__ internally.
        # obj.__init__(verbid, authtoken, sessionid)
        return obj

        """
        # Keep this for some time.
        deflist = [
            VerbId.PingMessage,
            VerbId.BeginTransactionRequest,
            VerbId.BeginTransactionResponse,
            VerbId.RollbackTransactionRequest,
            VerbId.RollbackTransactionResponse,
            VerbId.TraverseRequest,
            VerbId.TraverseResponse,
            VerbId.BeginExportRequest,
            VerbId.BeginExportResponse,
            VerbId.PartialExportRequest,
            VerbId.PartialExportResponse,
            VerbId.CancelExportRequest,
            VerbId.BeginImportRequest,
            VerbId.BeginImportResponse,
            VerbId.PartialImportRequest,
            VerbId.PartialImportResponse,
            VerbId.DumpStacktraceRequest,
            VerbId.DisconnectChannelRequest,
            VerbId. SessionForcefullyTerminated,
            VerbId.DecryptBufferRequest,
            VerbId.DecryptBufferResponse,
        ]
         
        if verbid == VerbId.HandShakeRequest:
            return HandShakeRequestMessage()

        elif verbid == VerbId.HandShakeResponse:
            return HandShakeResponseMessage()

        elif verbid == VerbId.AuthenticateRequest:
            return AuthenticateRequestMessage()

        elif verbid == VerbId.AuthenticateResponse:
            return AuthenticateResponseMessage()
        
        elif verbid in deflist:
            return DefaultMessage(verbid, authtoken, sessionid)
        
        else:
            raise tgexception.TGBadVerb('Bad Verb')
        """

    @classmethod
    def __createMessageFromBuffer__(cls, buf, offset, blen, cipher = None)\
            -> TGMessage:

        stream: ProtocolDataInputStream = ProtocolDataInputStream(buf, offset, blen, cipher)
        pduhdr: ProtocolMessageHeader = ProtocolMessageHeader()
        pduhdr.readHeader(stream)

        msg = cls.__createMessageFromVerb__(pduhdr.verbid, pduhdr.authtoken, pduhdr.sessionid)
        msg.fromBytes(buf[offset:offset + blen], cipher)

        return msg
