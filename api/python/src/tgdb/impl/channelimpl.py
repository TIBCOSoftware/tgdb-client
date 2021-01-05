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
 *  File name :channelimpl.py
 *  Created on: 5/30/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: channelimpl.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates all of the channel implementation
 """
import re
import abc
import sys
import socket
import select
import ssl
import traceback
import _thread
import struct
import threading
import typing
import enum
import ctypes as ct

import tgdb.pdu as tgpdu
import tgdb.impl.pduimpl as tgpduimpl
import tgdb.impl.atomics as tgatomics
import tgdb.log as tglog
import tgdb.channel as tgchannel
import tgdb.utils as tgutils

import tgdb.exception as tgexception
import tgdb.version as tgvers


class LinkUrl(tgchannel.TGChannelUrl):
    __defaulthost__ = 'localhost'
    __defaultport__ = 8222
    __defaultuser__ = 'admin'

    __protocol__: tgchannel.ProtocolType = tgchannel.ProtocolType.Tcp
    __host__ = 'localhost'
    __port__ = 8222
    __user__ = 'admin'
    __properties__: typing.Dict[str, str] = None

    def __init__(self, protocol=tgchannel.ProtocolType.Tcp, host='localhost', port=8222, user='admin', isipv6=False, props=None):
        self.__protocol__ = protocol
        self.__host__ = host
        self.__port__ = port
        self.__user__ = user
        self.__isipv6__ = isipv6
        self.__properties__ = tgutils.TGProperties() if props is None else props

    @property
    def protocol(self) -> tgchannel.ProtocolType:
        return self.__protocol__

    @property
    def host(self) -> str:
        return self.__host__

    @property
    def user(self) -> str:
        return self.__user__

    @property
    def port(self) -> int:
        return self.__port__

    @property
    def properties(self) -> dict:
        return self.__properties__

    @property
    def url(self) -> str:
        return '{0}://{1}@{2}{3}:{4}{5}'.format(
            self.protocol.name.lower(),
            self.user,
            '[' if self.__isipv6__ else '',
            self.host.lower(),
            self.port,
            ']' if self.__isipv6__ else ''
        )

    def __str__(self):
        return self.url

    @classmethod
    def parse(cls, url) -> tgchannel.TGChannelUrl:
        prot, host, port, user, isipv6, props = cls.__parseurl__(url)
        return LinkUrl(prot, host, port, user, isipv6, props)

    @classmethod
    def __parseurl__(cls, url) -> tuple:
        props = tgutils.TGProperties()
        urlpart, prot = cls.__parseprotocol__(url, props)
        urlpart, user = cls.__parseuser__(urlpart, props)
        urlpart, host, port, isipv6 = cls.__parsehostandport__(urlpart, props)
        props = cls.__parseproperties__(urlpart, props)
        return prot, host, port, user, isipv6, props

    @classmethod
    def __parseprotocol__(cls, url: str, props):

        if url.lower().startswith("tcp://"):
            prot = tgchannel.ProtocolType.Tcp
            pi = len("tcp://")
        elif url.lower().startswith("ssl://"):
            prot = tgchannel.ProtocolType.Ssl
            pi = len("ssl://")
        else:
            raise tgexception.TGException("Invalid protocol specified.")

        props[tgutils.ConfigName.ChannelProtocol] = prot

        return url[pi:], prot

    @classmethod
    def __parseuser__(cls, urlpart, props):
        pos = urlpart.find('@')
        if pos == -1:
            return urlpart, cls.__defaultuser__
        user = urlpart[:pos]
        props[tgutils.ConfigName.ChannelUserID] = user
        return urlpart[pos+1:], user


    @classmethod
    def __parsehostandport__(cls, urlpart: str, props):
        if len(urlpart) == 0:
            return urlpart, cls.__defaulthost__, cls.__defaultport__, False

        isipv6 = False
        propstart = urlpart.find('/')
        hostport = urlpart if propstart == -1 else urlpart[:propstart]
        stipv6 = hostport.find('[')
        edipv6 = hostport.find(']')
        lpos = hostport.rfind(':')

        if stipv6 != -1:
            if (edipv6 - stipv6) > 1:
                isipv6 = True
                lpos = hostport.rfind(':')
            else:
                raise tgexception.TGException('Invalid or missing host name')
        try:
            port = cls.__defaultport__ if lpos == -1 else int(hostport[lpos+1:])
        except(ValueError, OverflowError, TypeError):
            if isipv6:
                port = cls.__defaultport__
            else:
                raise tgexception.TGException("Bad port specified")
        host = urlpart if lpos == -1 else hostport[:lpos]
        props[tgutils.ConfigName.ChannelHost] = host
        props[tgutils.ConfigName.ChannelPort] = port
        if props[tgutils.ConfigName.TlsExpectedHostName] is None:
            props[tgutils.ConfigName.TlsExpectedHostName] = host + ":" + str(port)
        return '' if propstart == -1 else urlpart[propstart+1:], host, port, isipv6

    @classmethod
    def __parseproperties__(cls, urlpart, props):
        if len(urlpart) == 0:
            return props
        ll = len(urlpart)
        if urlpart[0] != '{' or urlpart[ll-1] !='}':
            raise tgexception.TGException('Invalid syntax - Property string must be embedded with {}.')

        urlpart = urlpart[1:ll-1]     # remove {}
        tokens = urlpart.split(';')     # split on ; assumption value doesn't contain ';'
        for token in tokens:
            key, value = token.split('=', 1)
            props[key] = value

        return props


class TrueVoidPtr(ct.c_void_p):
    pass


class SSLLibrary:

    singleton = None

    def __init__(self, library):
        import ctypes as ct
        self.__lib: ct.CDLL = library
        self.__lib.BIO_s_mem.restype = TrueVoidPtr
        self.__lib.BIO_new.restype = TrueVoidPtr
        self.__memBio = self.__lib.BIO_new(self.__lib.BIO_s_mem())

    def __del__(self):
        self.__lib.BIO_free(self.__memBio)

    def convertCertificate(self, cert_as_bytes: bytes) -> ct.c_void_p:
        str_buf = ct.create_string_buffer(cert_as_bytes, len(cert_as_bytes))
        buf_ptr = ct.pointer(str_buf)
        ret = ct.c_void_p()
        self.__run("d2i_X509", ct.byref(ret), ct.byref(buf_ptr), len(cert_as_bytes))
        if ret.value is None:
            raise tgexception.TGSecurityException("Could not parse X509 certificate.")
        return ret

    def getPubkeyContext(self, cert: ct.c_void_p) -> ct.c_void_p:
        pub_key = self.__run("X509_get0_pubkey", cert, ret_type=TrueVoidPtr)
        if pub_key is None:
            raise tgexception.TGSecurityException("Could not find public key in X509 certificate.")
        ret = self.__run("EVP_PKEY_CTX_new", pub_key, None, ret_type=TrueVoidPtr, num=2)
        if ret is None:
            raise tgexception.TGSecurityException("Could not find public key in X509 certificate.")
        return ret

    def getIssuerName(self, cert: ct.c_void_p):
        return self._getName("X509_get_issuer_name", cert)

    def getSubjectName(self, cert: ct.c_void_p):
        return self._getName("X509_get_subject_name", cert)

    def _getName(self, nameFunc: str, cert: ct.c_void_p) -> str:
        import ctypes as ct
        name = self.__run(nameFunc, cert, ret_type=TrueVoidPtr)
                        # Is the value of XN_FLAG_ONELINE for 1.1.1d \/ NEED TO UPDATE WHENEVER OPENSSL VERSION CHANGES!
        self.__run("X509_NAME_print_ex", self.__memBio, name, 0, 8520479, num=4, ret_type=ct.c_int)
        ret_buf = bytearray()
        num_recvd = 1
        while num_recvd > 0:
            read_buffer = ct.create_string_buffer(256)
            num_recvd = self.__run("BIO_gets", self.__memBio, read_buffer, 256, num=3, ret_type=ct.c_int)
            if num_recvd < 1:
                num_recvd = 0
            else:
                ret_buf.extend(read_buffer.raw[:num_recvd])
                if num_recvd < 255:
                    num_recvd = 0
        return ret_buf.decode('utf-8').split(',')[0]

    def encrypt(self, ctx: ct.c_void_p, data: bytes) -> bytes:
        if ctx is None or self.__run("EVP_PKEY_encrypt_init", ctx) <= 0:
            raise tgexception.TGSecurityException("Could not initialize public key context.")
        # The following corresponds to EVP_PKEY_CTX_set_rsa_padding(ctx, RSA_PKCS1_PADDING)
        # Need to make sure that the constants are correct whenever OpenSSL version changes.
        set_ret = self.__run("EVP_PKEY_CTX_ctrl_str", ctx, ct.c_char_p(b'rsa_padding_mode'), ct.c_char_p(b'pkcs1'),
                             num=3)
        if tglog.gLogger.isEnabled(tglog.TGLevel.DebugWire):
            num = ct.c_int()
            tglog.gLogger.log(tglog.TGLevel.DebugWire, "Return from EVP_PKEY_get_rsa_padding: %d",
                              self.__run("EVP_PKEY_CTX_ctrl", ctx, 6, -1, 4096+6, 0, ct.byref(num), num=6))
            tglog.gLogger.log(tglog.TGLevel.DebugWire, "Value of RSA padding: %d", num.value)
            tglog.gLogger.log(tglog.TGLevel.DebugWire, "Return from EVP_PKEY_set_rsa_padding: %d", set_ret)
            tglog.gLogger.log(tglog.TGLevel.DebugWire, "Return from EVP_PKEY_get_rsa_padding: %d",
                              self.__run("EVP_PKEY_CTX_ctrl", ctx, 6, -1, 4096+6, 0, ct.byref(num), num=6))
            tglog.gLogger.log(tglog.TGLevel.DebugWire, "Value of RSA padding: %d", num.value)
        outlen = ct.c_size_t(0)
        buffer = ct.create_string_buffer(data, len(data))
        if self.__run("EVP_PKEY_encrypt", ctx, None, ct.byref(outlen), buffer, len(data), num=5) <= 0:
            raise tgexception.TGSecurityException("Could not encrypt using context.")
        result = ct.create_string_buffer(outlen.value)
        if self.__run("EVP_PKEY_encrypt", ctx, result, ct.byref(outlen), buffer, len(data), num=5) <= 0:
            raise tgexception.TGSecurityException("Could not encrypt using context.")
        return result.raw

    def __run(self, functionName, one=None, two=None, three=None, four=None, five=None, six=None, ret_type=None,
              num=-1):
        try:
            function = getattr(self.__lib, functionName)
            if ret_type is not None:
                function.restype = ret_type
            if num == 0 or (one is None and num == -1):
                return function()
            elif num == 1 or (two is None and num == -1):
                return function(one)
            elif num == 2 or (three is None and num == -1):
                return function(one, two)
            elif num == 3 or (four is None and num == -1):
                return function(one, two, three)
            elif num == 4 or (five is None and num == -1):
                return function(one, two, three, four)
            elif num == 5 or (six is None and num == -1):
                return function(one, two, three, four, five)
            else:
                return function(one, two, three, four, five, six)
        except AttributeError as e:
            tglog.gLogger.log(tglog.TGLevel.Debug, "dir(self.__lib): %s", dir(self.__lib))
            tglog.gLogger.logException(tglog.TGLevel.Debug, str(e), e)
            raise tgexception.TGSecurityException("Could not find the correct SSL library function %s." % functionName,
                                                  cause=e)

    @classmethod
    def initialize(cls):
        import ctypes as ct
        import sys
        import os
        from ctypes.util import find_library
        libname: str
        if sys.platform == 'darwin':
            libname = "libcrypto.1.1.dylib"
        elif sys.platform == 'linux':
            libname = "libcrypto.so.1.1"
        else:
            libname = "libcrypto-1_1-x64.dll"

        library: ct.CDLL = None
        try:
            re_format = "[^{0}]*\\.zip{0}".format(os.sep.replace("\\", "\\\\"))
            if not re.search(re_format, os.path.dirname(__file__)):
                raise Exception("Not where we thought we are.")
            base_dir = re.split(re_format.format(os.sep), os.path.dirname(__file__))[0]
            libpath = base_dir + libname
            library = ct.CDLL(libpath)
        except Exception:
            tglog.gLogger.log(tglog.TGLevel.Warning, "Error finding SSL library, searching through local to find accept"
                  "able library.\nEven if successful, should not rely on this as a valid method for finding the SSL lib"
                  "rary, as that library may have different interfaces that can cause segfaults or worse.")
            libpath = ct.util.find_library("crypto")
            try:
                library = ct.CDLL(libpath)
            except:
                pass

        if library is not None:
            cls.singleton = cls(library)
        else:
            tglog.gLogger.log(tglog.TGLevel.Error, "Could not initialize SSL library. Do not encrypt/decrypt/use SSL ch"
                                                   "annel.")


SSLLibrary.initialize()


class Cipher(abc.ABC):

    @property
    @abc.abstractmethod
    def issuer_name(self) -> str:
        """Gets the certificate's issuer name."""

    @property
    @abc.abstractmethod
    def subject_name(self) -> str:
        """Gets the certificate's subject name."""

    @abc.abstractmethod
    def encrypt(self, buf: bytes) -> bytes:
        """Encrypts the data to send to the server."""

    @abc.abstractmethod
    def deobfuscate(self, buf: bytes) -> typing.Tuple[int, bytes]:
        """De-obfuscates the data from the server.

        NOTE: When receiving data from the server on a non-encrypted channel, the data IS NOT TRULY ENCRYPTED. This is
        to reduce the client's overhead. If it is truly vital that your data be encrypted from client to server and back
        again, then be sure to use SSL channels instead!

        :returns: a tuple containing first the number of bytes consumed from the buffer and second the de-obfuscated
                    data.
        """


class CipherImpl(Cipher):

    def __init__(self, cert: bytes, libraryLocation: typing.Optional[str] = None):
        self.__library: SSLLibrary = SSLLibrary.singleton
        if libraryLocation is not None:
            try:
                lib = ct.CDLL(libraryLocation)
                self.__library = SSLLibrary(lib)
            except:
                pass
        self.__cert: ct.c_void_p = None
        self.__pubKeyCtx: ct.c_void_p = None
        if self.__library is None:
            tglog.gLogger.log(tglog.TGLevel.Warning, "Could not find an SSL Library. Do not use encrypt, decrypt, or SS"
                                                     "L channels.")
            return
        try:
            self.__cert = self.__library.convertCertificate(cert)
            self.__pubKeyCtx = self.__library.getPubkeyContext(self.__cert)
        except tgexception.TGSecurityException as e:
            self.__library = None
            tglog.gLogger.logException(tglog.TGLevel.Debug, str(e), e)
            tglog.gLogger.log(tglog.TGLevel.Error, "Could not find the correct SSL library.")

    @property
    def subject_name(self) -> str:
        if self.__library is None:
            raise tgexception.TGSecurityException("Getting subject name of certificate when SSL Library was not loaded "
                                                  "properly!")
        return self.__library.getSubjectName(self.__cert)

    @property
    def issuer_name(self) -> str:
        if self.__library is None:
            raise tgexception.TGSecurityException("Getting subject name of certificate when SSL Library was not loaded "
                                                  "properly!")
        return self.__library.getIssuerName(self.__cert)

    def encrypt(self, buf: bytes) -> bytes:
        if self.__library is None:
            raise tgexception.TGSecurityException("Getting subject name of certificate when SSL Library was not loaded "
                                                  "properly!")
        if self.__pubKeyCtx is None:
            raise tgexception.TGSecurityException("Could not initialize the public key context successfully.")
        return self.__library.encrypt(self.__pubKeyCtx, buf)

    def deobfuscate(self, buf: bytes) -> typing.Tuple[int, bytes]:
        instream = tgpduimpl.ProtocolDataInputStream(bytearray(buf))
        ret = bytearray()
        rand = instream.readUnsignedLong()
        length = instream.readUnsignedLong()
        cnt = length // 8
        rem = length % 8

        for i in range(cnt):
            data = instream.readUnsignedLong() ^ rand
            ret.extend(data.to_bytes(8, byteorder='little'))

        for i in range(rem):
            ret.append(instream.readUnsignedByte())

        return instream.position, bytes(ret)


class BlockingChannelResponseWaiter(tgchannel.TGChannelResponseWaiter):

    def __init__(self, requestId, timeout=None):
        self.__requestid__ = requestId
        self.__timeout__ = timeout
        self.__status__ = tgchannel.TGChannelResponseWaiter.WaiterStatus.Waiting
        self.__reply__ = None
        self.__lock__ = threading.RLock()
        self.__cond__ = threading.Condition(self.__lock__)

    @property
    def isBlocking(self):
        return True

    @property
    def status(self):
        with self.__lock__:
            return self.__status__

    @property
    def requestid(self):
        return self.__requestid__

    def awaiting(self, status: tgchannel.TGChannelResponseWaiter.WaiterStatus):
        with self.__cond__:
            self.__cond__.wait_for(lambda: self.__status__ == status, timeout=self.__timeout__)

    @property
    def reply(self):
        with self.__lock__:
            return self.__reply__

    @reply.setter
    def reply(self, reply: tgpdu.TGMessage):
        with self.__cond__:
            self.__reply__ = reply
            self.__status__ = tgchannel.TGChannelResponseWaiter.WaiterStatus.Ok
            self.__cond__.notifyAll()


class SocketState(enum.Enum):

    Initialized = 0
    Connected = 1
    Closed = 2
    ErrorState = 3


class TcpSocket(tgchannel.TGSocket):

    # Java side socket type hierarchy doesn't make sense.
    _handle: socket.socket = None
    _props: tgutils.TGProperties = None
    _host: str = None
    _port: int = -1
    Read = 1
    Write = 2
    ReadWrite = 3

    def __init__(self, host='localhost', port=8222, timeout=5.000, props: tgutils.TGProperties=None):
        self._host = host
        self._port = port
        if timeout is not None and isinstance(timeout, str):
            timeout = float(timeout)
        self._timeout = timeout
        self._props = tgutils.TGProperties() if props is None else props
        try:
            if tglog.gLogger.isEnabled(tglog.TGLevel.DebugWire):
                for prop in self._props:
                    if prop is not None:
                        tglog.gLogger.log(tglog.TGLevel.DebugWire, "prop[%s] = %s", str(prop), str(self._props[prop]))
        except BaseException as e:
            tglog.gLogger.logException(tglog.TGLevel.Debug, "Unknown reason for exception while printing out props.", e)
            raise e
        self._handle = None
        self._state = SocketState.Initialized
        self._cipher: Cipher = None

    @property
    def cipher(self) -> Cipher:
        return self._cipher

    @cipher.setter
    def cipher(self, new_ciph: Cipher):
        self._cipher = new_ciph

    def _close(self, skt: socket.socket):
        if skt is not None:
            try:
                skt.shutdown(socket.SHUT_RDWR)
            except OSError:
                pass
            skt.close()

    def connect(self):
        for skt, inetaddr in self._createSocket():
            try:
                skt.connect(inetaddr)
                skt.setblocking(False)
                self._handle = skt
                self._state = SocketState.Connected
            except (socket.error, IOError, TimeoutError):
                self._close(skt)
                self._handle = None
        if self._handle is None:
            raise IOError('Invalid Address tuple ({0}, {1}) specified'.format(self._host, self._port))

    def close(self):
        self._close(self._handle)

    @property
    def handle(self) -> socket.socket:
        return self._handle

    @property
    def outboxaddr(self):
        return '{0}:{1}'.format(self._host, str(self._port))

    @property
    def inboxaddr(self):
        t: tuple = self._handle.getsockname()
        return '{0}:{1}'.format(t[0], t[1])

    def _isReady(self, type=Read, timeout=None) -> bool:
        # use select on the handle for read/write
        rlist, wlist, __ = select.select([self._handle] if type & TcpSocket.Read == TcpSocket.Read else [],
                                         [self._handle] if type & TcpSocket.Write == TcpSocket.Write else [],
                                         [],
                                         timeout)
        return self._handle in rlist or self._handle in wlist

    def _createSocket(self) -> typing.Generator[typing.Tuple[socket.socket, typing.Tuple], None, None]:
        # Acts like a Generator Function
        recvSize = self._props[tgutils.ConfigName.ChannelRecvSize]
        sendSize = self._props[tgutils.ConfigName.ChannelSendSize]
        rs = 32 if recvSize is None or len(recvSize) == 0 or int(recvSize) <= 0 else int(recvSize)
        ss = 32 if sendSize is None or len(sendSize) == 0 or int(sendSize) <= 0 else int(sendSize)
        lingerfmt = 'hh' if sys.platform == 'win32' else 'ii'
        try:
            tglog.gLogger.log(tglog.TGLevel.Debug, 'My host:port is %s:%d', self._host, self._port)
            addrinfos = socket.getaddrinfo(host=self._host, port=self._port, proto=socket.IPPROTO_TCP)
        except Exception as e:
            tglog.gLogger.logException(tglog.TGLevel.Debug, e)
            return

        for info in addrinfos:
            if self._state == SocketState.Connected:
                return
            skt = None
            skt = socket.socket(info[0], info[1])
            skt.setblocking(True)
            skt.settimeout(self._timeout)
            skt.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            if "SO_REUSEPORT" in dir(socket):
                tglog.gLogger.log(tglog.TGLevel.Debug, "Using SO_REUSEPORT")
                skt.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)
            elif "SO_REUSE_UNICASTPORT" in dir(socket):
                tglog.gLogger.log(tglog.TGLevel.Debug, "Using SO_REUSE_UNICASTPORT")
                skt.setsockopt(socket.SOL_SOCKET, socket.SO_REUSE_UNICASTPORT, 1)
            else:
                tglog.gLogger.log(tglog.TGLevel.Debug, "Not using any SO_REUSE\\w*PORT")
            skt.setsockopt(socket.SOL_SOCKET, socket.SO_SNDBUF, ss * 1024)
            skt.setsockopt(socket.SOL_SOCKET, socket.SO_RCVBUF, rs * 1024)
            skt.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)
            skt.setsockopt(socket.SOL_SOCKET, socket.SO_LINGER, struct.pack(lingerfmt, 1, 0))
            inetaddr = info[4]
            yield skt, inetaddr

    def send(self, msg: tgpdu.TGMessage) -> int:
        if self._handle is None:
            raise tgexception.TGException("Channel is closed")
        if not self._isReady(TcpSocket.Write, timeout=5.0):
            raise IOError('Peer terminated')

        if tglog.TGLogManager.getLogger().isEnabled(tglog.TGLevel.DebugWire):
            tglog.TGLogManager.getLogger().log(tglog.TGLevel.DebugWire,
                "----------------- outgoing message --------------------")
            tglog.TGLogManager.getLogger().log(tglog.TGLevel.DebugWire, str(msg))

        return self._sendMsg(msg)

    def _sendMsg(self, msg: tgpdu.TGMessage) -> int:
        buf = msg.toBytes(self._cipher)
        totalsent = 0
        size = len(buf)
        tglog.gLogger.log(tglog.TGLevel.Debug, "Sending message of size %d", size)
        while totalsent < size:
            sent = self._handle.send(buf[totalsent:])
            tglog.gLogger.log(tglog.TGLevel.TGLDebugFine, "Sent message of size %d", sent)
            if sent == 0:
                raise tgexception.TGException("Channel is disconnected")
            totalsent = totalsent + sent
        return totalsent

    def _readMsg(self, buf: bytearray, sock: ssl.SSLObject = None) -> tgpdu.TGMessage:
        if buf is None or len(buf) == 0:
            raise IOError('Peer terminated')
        try:
            size = struct.unpack('>i', buf)[0]
            if size <= 0:
                raise IOError('Peer terminated')
            while len(buf) < size:
                try:
                    tmp_buf = self._handle.recv(size - len(buf))
                    if tmp_buf is None or len(tmp_buf) <= 0:
                        raise IOError('Peer terminated')
                    buf.extend(tmp_buf)
                except BlockingIOError as e:    # Doing this because if the server takes a long time to respond it is
                    if e.errno != 35:           # does not always imply that the server is dead or connection lost.
                        raise IOError('Unknown error: ' + str(e))
            if tglog.gLogger.isEnabled(tglog.TGLevel.DebugWire):
                tglog.gLogger.log(tglog.TGLevel.DebugWire, tgutils.HexUtils.formatHex(buf))
        except BaseException as e:
            tglog.gLogger.logException(tglog.TGLevel.Debug, "Unknown exception while reading message!", e)
            raise e
        return tgpduimpl.TGMessageFactory.createMessage(bytes(buf), self._cipher)

    def recvMsg(self) -> tgpdu.TGMessage:
        while self._isReady(TcpSocket.Read):
            buf = bytearray()
            while len(buf) < 4:
                tmp_buf = self._handle.recv(4 - len(buf))
                if tmp_buf is None or len(tmp_buf) <= 0:
                    raise IOError('Peer terminated')
                buf.extend(tmp_buf)
            return self._readMsg(buf)

    def tryRead(self) -> tgpdu.TGMessage:
        if self._isReady(TcpSocket.Read, timeout=0.0):
            return self.recvMsg()
        return None

    @classmethod
    def test(cls):
        tcpsocket: TcpSocket = TcpSocket()
        tcpsocket.connect()
        print(tcpsocket.inboxaddr)


class SslSocket(TcpSocket):

    def __init__(self, host: str = 'localhost', port: int = 8222, timeout: float = 5.000,
                 props: tgutils.TGProperties = None, sslcontext: ssl.SSLContext = None):
        super().__init__(host, port, timeout, props)
        self.__sslcontext = sslcontext if sslcontext is not None else\
                        ssl.create_default_context(purpose=ssl.Purpose.SERVER_AUTH)
        self.__socket: socket.socket = None
        self.__sentLast = None

    def verifyDB(self, peercert: bytes):
        shouldVerify: bool = self._props[tgutils.ConfigName.TlsVerifyDatabaseName]
        dbName: str = self._props[tgutils.ConfigName.ConnectionDatabaseName]
        if shouldVerify:
            if dbName is None:
                raise tgexception.TGException("Should verify was set, but no database name was provided!")
            expected_subject = "CN = %s.ce$t" % dbName
            expected_issuer = "CN = %s.1nit" % dbName
            subj_name = self.cipher.subject_name
            issuer_name = self.cipher.issuer_name
            if subj_name is None:
                raise tgexception.TGException("Could not verify: No subject name discovered.")

            if issuer_name is None:
                raise tgexception.TGException("Could not verify: No subject name discovered.")

            if expected_issuer != issuer_name:
                raise tgexception.TGAuthenticationException("Server's issuer name: %s, expected: %s" %
                                                            (issuer_name, expected_issuer))

            if expected_subject != subj_name:
                raise tgexception.TGAuthenticationException("Server's subject name: %s, expected: %s" %
                                                            (subj_name, expected_subject))

    def connect(self):
        for skt, inetaddr in self._createSocket():
            self.__socket = skt
            try:
                self.__socket.connect(inetaddr)
                tglog.gLogger.log(tglog.TGLevel.Debug, "%s", self._props[tgutils.ConfigName.TlsExpectedHostName])
                skt = self.__sslcontext.wrap_socket(skt, server_side=False, do_handshake_on_connect=False,
                                                    server_hostname=self._props[tgutils.ConfigName.TlsExpectedHostName])
                handshake_failed = True
                while handshake_failed:
                    try:
                        skt.do_handshake()
                        handshake_failed = False
                    except ssl.SSLWantReadError:
                        select.select([skt], [], [])
                    except ssl.SSLWantWriteError:
                        select.select([], [skt], [])

                skt.setblocking(False)
                self._handle = skt
                self._state = SocketState.Connected
            except (socket.error, IOError, TimeoutError) as e:
                self._close(skt)
                self._handle = None
        if self._handle is None:
            raise IOError('Invalid Address tuple ({0}, {1}) specified'.format(self._host, self._port))

    def _sendMsg(self, msg: tgpdu.TGMessage):
        buf = msg.toBytes(self._cipher)
        totalsent = 0
        size = len(buf)
        tglog.gLogger.log(tglog.TGLevel.TGLDebug, "Sending message of size %d. Last bytes: %s", size, str(buf[-5:]))
        self.__sentLast = buf
        while totalsent < size:
            exception_raised = False
            sent: int = 0
            try:
                sent = self._handle.send(buf[totalsent:])
                tglog.gLogger.log(tglog.TGLevel.TGLDebugFine, "Sent message of size %d", size)
            except ssl.SSLWantWriteError:
                select.select([], [self._handle], [])
                exception_raised = True
            except ssl.SSLWantReadError:
                select.select([self._handle], [], [])
                exception_raised = True
            if sent == 0:
                raise IOError('Peer terminated')
            if not exception_raised:
                if sent == 0:
                    raise tgexception.TGException("Channel is disconnected")
                totalsent += sent
        return totalsent

    def _readMsg(self, buf: bytearray, sock: ssl.SSLObject = None) -> tgpdu.TGMessage:
        if buf is None or len(buf) == 0:
            raise IOError('Peer terminated')
        size = struct.unpack('>i', buf)[0]
        if size <= 0:
            raise IOError('Peer terminated')
        while len(buf) < size:
            try:
                tmp_buf = self._handle.recv(size - len(buf))
                if tmp_buf is None or len(tmp_buf) <= 0:
                    raise IOError('Peer terminated')
                buf.extend(tmp_buf)
            except ssl.SSLWantReadError:
                select.select([self._handle], [], [])
            except ssl.SSLWantWriteError:
                select.select([], [self._handle], [])
        if tglog.gLogger.isEnabled(tglog.TGLevel.DebugWire):
            tglog.gLogger.log(tglog.TGLevel.DebugWire, tgutils.HexUtils.formatHex(buf))
        return tgpduimpl.TGMessageFactory.createMessage(bytes(buf), cipher=self._cipher)


    def recvMsg(self) -> tgpdu.TGMessage:
        buf: bytearray
        try:
            buf = bytearray()
            while True:
                # Need to put into bytearray because SSLSockets do not support peeking at the data...
                try:
                    tmp_buf = self._handle.recv(4 - len(buf))
                    if tmp_buf is None or len(tmp_buf) <= 0:
                        raise IOError('Peer terminated')
                    buf.extend(tmp_buf)
                except ssl.SSLWantReadError:
                    select.select([self._handle], [], [], self._timeout)
                except ssl.SSLWantWriteError:
                    select.select([], [self._handle], [], self._timeout)
                if len(buf) == 4:
                    return self._readMsg(buf)
        except BaseException:
            if self.__sentLast is not None and tglog.gLogger.isEnabled(tglog.TGLevel.DebugWire):
                buf = self.__sentLast
                self.__sentLast = None
                for i in range(0, len(buf), 32):
                    toPrint = "%08x" % i
                    buffer = buf[i:min(len(buf), i + 32)]
                    for j in range(0, len(buffer), 2):
                        toPrint += " %02x" % buffer[j]
                        if j + 1 < len(buffer):
                            toPrint += "%02x" % buffer[j + 1]
                    tglog.gLogger.log(tglog.TGLevel.DebugWire, "%s", toPrint)
            raise

    def close(self):
        foundException = True
        skt: socket.socket = None
        if self._handle is None:
            return
        try:
            while foundException:
                try:
                    skt = self._handle.unwrap()
                    foundException = False
                except ssl.SSLWantWriteError:
                    select.select([], [self._handle], [])
                except ssl.SSLWantReadError:
                    select.select([self._handle], [], [])
                except OSError as e:
                    raise e
                except:
                    pass
            try:
                self._close(skt)
            except:
                pass
        except OSError as e:
            if e.errno != 0:        # For whatever reason, it seams that SSL sometimes causes an error to occur
                raise e


class ChannelReader(threading.Thread):

    _readerNum = tgatomics.AtomicReference('i', 0)

    def __init__(self, channel: tgchannel.TGChannel):
        super().__init__(name='TGLinkReader{0}.{1}'.format(channel.clientid, ChannelReader._readerNum.increment()))
        self.__channel__ = channel
        self.__isRunning__ = tgatomics.AtomicReference('i', 0)

    def start(self) -> None:
        self.__isRunning__.set(1)
        super().start()

    def stop(self):
        self.__isRunning__.set(0)

    def run(self) -> None:
        while self.__isRunning__.get() == 1:
            try:
                if not self.is_alive() or self.__channel__.isClosed or self.__isRunning__.get() == 0:
                    tglog.gLogger.log(tglog.TGLevel.Debug, 'Thread {0} interrupted or channel is closed. Stopping reade'
                                                           'r'.format(self.name))
                    return

                try:
                    # Added this in to prevent the extraneous output to the admin console. Will act like
                    # Connection.disconnect in terms of the end result.
                    msg = self.__channel__.readMessage()
                except IOError as e:
                    self.__channel__.disconnect()
                    self.__channel__.stop()
                    self.stop()
                    tglog.TGLogManager.getLogger().log(tglog.TGLevel.Error, "Channel closed. Reconnect or exit. More de"
                                                "tail: %s (of type %s)", str(e.args), str(type(e)))
                    break

                if msg is None or msg.verbid == tgpdu.VerbId.PingMessage:
                    continue

                if msg.verbid == tgpdu.VerbId.SessionForcefullyTerminated:
                    self.__channel__.stop(tgchannel.TGChannelStopMethod.RemoteKill, msg='Killed by Peer')
                    self.__isRunning__.set(0)
                    tglog.TGLogManager.getLogger().log(tglog.TGLevel.Error, "Connection killed by Peer.")
                    break

                if not self.is_alive() or self.__channel__.isClosed or self.__isRunning__.get() == 0:
                    tglog.TGLogManager.getLogger().log(tglog.TGLevel.Error, 'Thread {0} interrupted or channel is close'
                                                                            'd. Stopping reader'.format(self.name))
                    return

                self.__processMessage__(msg)
            except Exception as e:
                if self.__isRunning__.get() == 0:
                    break
                exres = self.__channel__.handleException(e)
                for requestid, waiter in self.__channel__.waiters.items():
                    msg = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.ExceptionMessage)
                    msg.errormsg = exres
                    waiter.reply = msg

                if exres != tgchannel.ExceptionHandleResult.RetryOperation:
                    tglog.TGLogManager.getLogger().logException(tglog.TGLevel.Debug, 'Exiting channel reader thread', e)
                    break

    def __processMessage__(self, msg: tgpdu.TGMessage):
        logger = tglog.TGLogManager.getLogger()
        waiter: tgchannel.TGChannelResponseWaiter =\
            self.__channel__.waiters[msg.requestid] if msg.requestid in self.__channel__.waiters else None
        if waiter is None:
            logger.log(tglog.TGLevel.Error, 'No waiter for the corresponding request:{0}'.format(msg.requestid))
            foundWaiter = False
            for _, waiter in self.__channel__.waiters.items():
                waiter.reply = msg
                foundWaiter = True
                break
            if not foundWaiter:
                raise tgexception.TGException("Could not find a corresponding waiter for requestid {0}"
                                               .format(msg.requestid))
        else:
            if logger.isEnabled(tglog.TGLevel.DebugWire):
                logger.log(tglog.TGLevel.DebugWire,
                           'Process msg: {0}'.format(tgutils.HexUtils.formatHex(msg.toBytes())))
            waiter.reply = msg


class TcpChannel(tgchannel.TGChannel):

    gRequestId = tgatomics.AtomicReference('q', 0)

    def __init__(self, url: tgchannel.TGChannelUrl, props: tgutils.TGProperties):
        self.__url__ = url
        self.__primaryurl__ = url
        self.__props__ = props
        self.__reader__ = ChannelReader(self)
        self.__waiters__: typing.Dict[int, tgchannel.TGChannelResponseWaiter] = dict()
        self.__authtoken__ = -1
        self.__sessionid__ = -1
        self.__socket__: TcpSocket = None
        self.__clientversion__ = 2.0
        self.__serverversion__ = 0.0
        self.__linkstate__: tgatomics.AtomicReference =\
            tgatomics.AtomicReference('i', tgchannel.LinkState.NotConnected.value)
        self.__sslmode__ = False
        self.__sendLock__ = threading.Lock()
        self.__servercert__ = []


    @property
    def linkstate(self) -> tgchannel.LinkState:
        return tgchannel.LinkState.fromId(self.__linkstate__.get())

    @linkstate.setter
    def linkstate(self, v: tgchannel.LinkState):
        self.__linkstate__.set(v.value)

    @property
    def outboxaddr(self):
        return self.__socket__.outboxaddr

    @property
    def properties(self) -> tgutils.TGProperties:
        return self.__props__

    @property
    def protocolversions(self) -> tuple:
        return self.__clientversion__, self.__serverversion__

    @property
    def authtoken(self) -> int:
        return self.__authtoken__

    @property
    def sessionid(self) -> int:
        return self.__sessionid__

    @property
    def clientid(self) -> str:
        return self.__props__[tgutils.ConfigName.ChannelClientId]

    def createSocket(self):     # override this method
        url = self.__primaryurl__
        tglog.gLogger.log(tglog.TGLevel.Warning, "Using insecure tcp connection instead of ssl. Should only use for tes"
                                                 "ting purposes or behind a secure firewall. This is not recommended!")
        return TcpSocket(host=url.host, port=url.port, timeout=self.__props__[tgutils.ConfigName.ChannelConnectTimeout],
                         props=self.properties)

    def connect(self):
        if self.isConnected:
            return
        if self.linkstate == tgchannel.LinkState.Reconnecting:
            return
        if self.inError():
            raise tgexception.TGException('Socket in error state.')

        if self.isClosed or self.linkstate == tgchannel.LinkState.NotConnected:
            try:
                self.__socket__ = self.createSocket()
                self.__socket__.connect()
                msg = self.__socket__.tryRead()
                if msg is not None and isinstance(msg, tgpduimpl.SessionForcefullyTerminatedMessage):
                    raise tgexception.TGChannelDisconnectedException('Channel disconnected')
                self.linkstate = tgchannel.LinkState.Connected
                self.__performHandShake__(self.__sslmode__)
                self.__authenticate__()
            except (tgexception.TGAuthenticationException, tgexception.TGChannelDisconnectedException, tgexception.TGVersionMismatchException) as e:
                raise
            except Exception as e:
                self.__socket__.close()
                raise

    def start(self):
        self.__reader__.start()

    def disconnect(self):
        if self.isConnected:
            for requestId in self.__waiters__:
                self.__waiters__[requestId].reply = None
            self.__linkstate__.set(tgchannel.LinkState.Closing.value)
            self.__reader__.stop()
            self.__socket__.close()
            self.__linkstate__.set(tgchannel.LinkState.Closed.value)

    def stop(self, stopmethod=tgchannel.TGChannelStopMethod.Graceful, msg=None):
        self.disconnect()

    def send(self, msg: tgpdu.TGMessage, waiter: tgchannel.TGChannelResponseWaiter = None) -> tgpdu.TGMessage:

        msg.requestid = waiter.requestid if waiter is not None else TcpChannel.gRequestId.decrement()

        # SS:TODO support retries and exception management the way java does.
        with self.__sendLock__:
            if not self.isConnected:
                raise tgexception.TGException('Channel Closed or not connected. Connect first and then try')
            if waiter is None or not waiter.isBlocking:
                self.__socket__.send(msg)
                return None
            self.__waiters__[waiter.requestid] = waiter
            self.__socket__.send(msg)
            waiter.awaiting(tgchannel.TGChannelResponseWaiter.WaiterStatus.Ok)
            self.__waiters__.pop(waiter.requestid)
            reply = waiter.reply
            if reply is None:
                raise tgexception.TGException('Expected a response, got none')
            if isinstance(reply, tgpduimpl.SessionExceptionMessage):
                raise tgexception.TGException(reply.errormsg)   # SS:TODO - support retry
            return reply

    def readMessage(self) -> tgpduimpl.TGMessage:
        return self.__socket__.recvMsg()

    def requestReply(self, request: tgpduimpl.TGMessage) -> tgpduimpl.TGMessage:
        if not self.isConnected:
            raise tgexception.TGException('Channel Closed or not connected. Connect first and then try')

        if self.__reader__.__isRunning__.get() == 0:
            with self.__sendLock__:
                try:
                    if tglog.gLogger.isEnabled(tglog.TGLevel.DebugWire):
                        tglog.gLogger.log(tglog.TGLevel.DebugWire, 'Send Message of Type:{0}:{1}'.format(type(request),
                                                                    tgutils.HexUtils.formatHex(request.toBytes())))
                    self.__socket__.send(request)
                    reply = self.__socket__.recvMsg()
                    return reply
                except Exception as e:
                    # SS:TODO - handle the exception

                    raise e
        else:
            timeout = self.__props__.get(tgutils.ConfigName.ConnectionOperationTimeoutSeconds, None)
            if timeout is not None and isinstance(timeout, str):
                timeout = float(timeout)
            bcr = BlockingChannelResponseWaiter(TcpChannel.gRequestId.decrement(), timeout=timeout)
            return self.send(request, bcr)

    @property
    def waiters(self) -> typing.Dict[int, tgchannel.TGChannelResponseWaiter]:
        return self.__waiters__

    def handleException(self, e: Exception) -> tgchannel.ExceptionHandleResult:
        return tgchannel.ExceptionHandleResult.RethrowException

    def __performHandShake__(self, sslMode):
        # TODO Move to using a ChannelResponseWaiter with an auto-decrementing request id?
        hsrequest: tgpduimpl.HandShakeRequestMessage = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.HandShakeRequest)
        hsrequest.requesttype = tgpduimpl.HandshakeRequestType.Initiate
        hsresponse: tgpduimpl.HandShakeResponseMessage = self.requestReply(hsrequest)
        self.__validatehandshakeResponse__(hsresponse, tgpduimpl.HandshakeResponseStatus.AcceptChallenge)

        hsrequest.updateSequenceAndTimeStamp()
        hsrequest.requesttype = tgpduimpl.HandshakeRequestType.ChallengeAccepted
        hsrequest.sslmode = self.__sslmode__
        hsrequest.challenge = tgvers.TGVersion.getVersion().toUInt64()
        hsresponse = self.requestReply(hsrequest)
        self.__validatehandshakeResponse__(hsresponse, tgpduimpl.HandshakeResponseStatus.ProceedWithAuthentication)

    def __validatehandshakeResponse__(self, response, status):
        if not isinstance(response, tgpduimpl.HandShakeResponseMessage):
            if isinstance(response, tgpduimpl.SessionForcefullyTerminatedMessage):
                raise tgexception.TGChannelDisconnectedException('Remote terminated')
            else:
                raise tgexception.TGException('Expected Handshake response')
        if response.status == status and status == tgpduimpl.HandshakeResponseStatus.AcceptChallenge:
            # match version
            return
        elif response.status == status and status == tgpduimpl.HandshakeResponseStatus.ProceedWithAuthentication:
            return
        else:
            raise tgexception.TGException('Handshake failed')

    def __authenticate__(self):
        # TODO Move to using a ChannelResponseWaiter with an auto-decrementing request id?
        authreq: tgpduimpl.AuthenticateRequestMessage = tgpduimpl.TGMessageFactory.createMessage(tgpdu.VerbId.AuthenticateRequest)
        authreq.clientid = self.clientid
        authreq.inboxaddr = self.__socket__.inboxaddr
        authreq.username = self.properties[tgutils.ConfigName.ChannelUserID]
        authreq.dbName = self.properties[tgutils.ConfigName.ConnectionDatabaseName]
        authreq.password = self.properties[tgutils.ConfigName.ChannelPassword]
        authreq.dbname = self.properties[tgutils.ConfigName.ConnectionDatabaseName]
        specifiedRoles = self.properties[tgutils.ConfigName.ConnectionSpecifiedRoles]
        if specifiedRoles is not None and isinstance(specifiedRoles, str):
            specifiedRoles = [] if specifiedRoles == '' else specifiedRoles.split(',')
        elif specifiedRoles is not None and isinstance(specifiedRoles, list):
            for rolename in specifiedRoles:
                if not isinstance(rolename, str):
                    specifiedRoles = None
                    break
        else:
            specifiedRoles = None
        if specifiedRoles is not None:
            for i in range(len(specifiedRoles)):
                specifiedRoles[i] = specifiedRoles[i].strip()
        if 'roles' in dir(authreq):
            authreq.roles = specifiedRoles
        authresp: tgpduimpl.AuthenticateResponseMessage = self.requestReply(authreq)
        if authresp.isSuccess:
            self.__authtoken__ = authresp.authtoken
            self.__sessionid__ = authresp.sessionid
            self.__servercert__ = authresp.serverCertificate
            self.__socket__.cipher = CipherImpl(bytes(self.__servercert__),
                                                self.properties[tgutils.ConfigName.TlsLibraryLocation])
            tglog.TGLogManager.getLogger().log(tglog.TGLevel.Info, 'Connected successfully to database %s hosted on ser'
                                                                   'ver at %s', authreq.dbName, self.__url__.url)
        else:
            raise tgexception.TGAuthenticationException('Bad user/password combination', errorcode=authresp.errorStatus)


class TGCipherSuites(enum.Enum):

    __supportedSuites__: typing.Set[str] = set()

    TLS_RSA_WITH_AES_128_CBC_SHA256 =\
    (0x3c, "AES128-SHA256", "RSA", "AES", "128")

    TLS_RSA_WITH_AES_256_CBC_SHA256 =\
    (0x3d, "AES256-SHA256", "RSA", "AES", "256")
    
    TLS_DHE_RSA_WITH_AES_128_CBC_SHA256 =\
    (0x67, "DHE-RSA-AES128-SHA256", "DH", "AES", "128")
    
    TLS_DHE_RSA_WITH_AES_256_CBC_SHA256 =\
    (0x6b, "DHE-RSA-AES256-SHA256", "DH", "AES", "256")
    
    TLS_RSA_WITH_AES_128_GCM_SHA256 =\
    (0x9c, "AES128-GCM-SHA256", "RSA", "AESGCM", "128")
    
    TLS_RSA_WITH_AES_256_GCM_SHA384 =\
    (0x9d, "AES256-GCM-SHA384", "RSA", "AESGCM", "256")
    
    TLS_DHE_RSA_WITH_AES_128_GCM_SHA256 =\
    (0x9e, "DHE-RSA-AES128-GCM-SHA256", "DH", "AESGCM", "128")
    
    TLS_DHE_RSA_WITH_AES_256_GCM_SHA384 =\
    (0x9f, "DHE-RSA-AES256-GCM-SHA384", "DH", "AESGCM", "256")
    
    TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256 =\
    (0xc023, "ECDHE-ECDSA-AES128-SHA256", "ECDH", "AES", "128")
    
    TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384 =\
    (0xc024, "ECDHE-ECDSA-AES256-SHA384", "ECDH", "AES", "256")
    
    TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256 =\
    (0xc027, "ECDHE-RSA-AES128-SHA256", "ECDH", "AES", "128")
    
    TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384 =\
    (0xc028, "ECDHE-RSA-AES256-SHA384", "ECDH", "AES", "256")
    
    TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 =\
    (0xc02b, "ECDHE-ECDSA-AES128-GCM-SHA256", "ECDH", "AESGCM", "128")
    
    TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384 =\
    (0xc02c, "ECDHE-ECDSA-AES256-GCM-SHA384", "ECDH", "AESGCM", "256")
    
    TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 =\
    (0xc02f, "ECDHE-RSA-AES128-GCM-SHA256", "ECDH", "AESGCM", "128")
    
    TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 =\
    (0xc030, "ECDHE-RSA-AES256-GCM-SHA384", "ECDH", "AESGCM", "256")
    
    TLS_INVALID_CIPHER = (0, None, None, None, None)

    def __init__(self, suiteID: int, openSSLName: str, keyExch: str, encr: str, bits: str):
        self.suiteID = suiteID
        self.openSSLName = openSSLName
        self.keyExch = keyExch
        self.encryption = encr
        self.bits = bits
        self.__class__.__supportedSuites__.add(openSSLName)

    @classmethod
    def filterSuites(cls, suites: typing.List[typing.Dict[str, typing.Any]]):

        for suite in suites:
            if suite['name'] not in cls.__supportedSuites__:
                suites.remove(suite)


class SslChannel(TcpChannel):

    def __init__(self, url: tgchannel.TGChannelUrl, props: tgutils.TGProperties):
        super().__init__(url, props)
        self.__sslmode__ = True

    def __updateProps(self, cfname: tgutils.ConfigName):
        props = self.properties
        url_props = self.__url__.properties
        val = url_props[cfname]
        if val is not None:
            props[cfname] = val

    def createSocket(self):     # override this method
        url = self.__primaryurl__
        ctx: ssl.SSLContext = self.__generateContext()
        self.__updateProps(tgutils.ConfigName.ConnectionDatabaseName)
        self.__updateProps(tgutils.ConfigName.TlsVerifyDatabaseName)
        return SslSocket(host=url.host, port=url.port, timeout=self.__props__[tgutils.ConfigName.ChannelConnectTimeout],
                         props=self.properties, sslcontext=ctx)

    def __authenticate__(self):
        super().__authenticate__()
        self.__socket__.verifyDB(self.__servercert__)

    def __generateContext(self) -> ssl.SSLContext:
        ctx = ssl.create_default_context()
        ctx.minimum_version = ssl.TLSVersion.TLSv1_2
        ctx.load_default_certs(purpose=ssl.Purpose.SERVER_AUTH)
        #hostname = self.properties[tgutils.ConfigName.TlsExpectedHostName]
        #if hostname is None:         TODO Make this an option?
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
        trustedCerts: str = self.__props__[tgutils.ConfigName.TlsTrustedCertificates]
        if trustedCerts is not None:
            strings: typing.List[str] = trustedCerts.split(";")
            for filepath in strings:
                ctx.load_verify_locations(cafile=filepath)

        return ctx

