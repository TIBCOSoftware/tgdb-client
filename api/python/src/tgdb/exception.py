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
 *  File name :model.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: exception.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates all the exception defined the exception java package.
 """
import enum
from enum import *
from abc import *
import re


class ExceptionType(Enum):
    """The type of exception encountered."""
    BadVerb = 1
    InvalidMessageLength = 2
    BadMagic = 3
    ProtocolNotSupported = 4
    BadAuthentication = 5
    IOException = 6
    ConnectionTimeout = 7
    GeneralException = 8
    RetryIOException = 9
    DisconnectedException = 10


class TGResultNameLibrary:
    singleton = None

    def __init__(self, cdll):
        import ctypes
        self.__cdll: ctypes.CDLL = cdll

    def getResultName(self, result: int) -> str:
        try:
            import ctypes as ct
            function = getattr(self.__cdll, "tg_result_getName")
            function.restype = ct.c_char_p
            func_ret: bytes = function(result)
            return func_ret.decode('utf-8')
        except:
            return None

    @classmethod
    def initialize(cls):
        try:
            import ctypes as ct
            import re
            import os
            import sys
            libname: str
            if sys.platform == 'darwin':
                libname = "libtgcommon.dylib"
            elif sys.platform == 'linux':
                libname = "libtgcommon.so"
            else:
                libname = "libtgcommon.dll"

            re_format = "[^{0}]*\\.zip{0}".format(os.sep.replace("\\", "\\\\"))
            if not re.search(re_format, os.path.dirname(__file__)):
                raise Exception("Not where we thought we are.")
            base_dir = re.split(re_format.format(os.sep), os.path.dirname(__file__))[0]
            libpath = base_dir + libname
            library = ct.CDLL(libpath)
            cls.singleton = cls(library)
        except:
            cls.singleton = None


TGResultNameLibrary.initialize()


class TGException(Exception):
    """TIBCO Graph Database specific exceptions."""

    def __init__(self, reason, errorcode=None, cause=None):
        """Client code should never initialize any exceptions directly."""
        super().__init__(reason)
        if cause is not None and isinstance(cause, TGException):
            if errorcode is None:
                errorcode = cause.errorcode
        self._errorcode = errorcode
        self._cause = cause

    def __str__(self):
        to_ret = super().__str__() + self.errorname
        try:
            import re
            if not re.search("error code: \\d+", to_ret.lower()):
                to_ret += " Error Code: %d" % self.errorcode
        except:
            to_ret = super().__str__() + self.errorname
        return to_ret

    @classmethod
    def buildException(cls, reason, errorcode=None, cause=None):
        """Client code should never initialize any exceptions directly."""
        return TGException(reason, errorcode, cause)

    @property
    def errorname(self) -> str:
        if self.errorcode is None:
            return ""
        tg_resName = "Error Code: %d" % self.errorcode
        if TGResultNameLibrary.singleton is not None:
            tmpName = TGResultNameLibrary.singleton.getResultName(self.errorcode)
            if tmpName is not None:
                tg_resName = " Error Name: " + tmpName
        return tg_resName

    @property
    def exceptionType(self):
        """Gets the type of the exception."""
        return ExceptionType.GeneralException

    @property
    def errorcode(self):
        """Gets the exception's error code."""
        return self._errorcode

    @property
    def cause(self):
        """Represents the cause of this exception, if there was one.

        For example, this might occur if a KeyError gets raised when it should not, so the Python API handles that
        exception and tells the user what is the likely reason."""
        return self._cause


class TGAuthenticationException(TGException):
    """Exception for when the server refuses to authenticate the user because of a bad password or username."""
    @property
    def exceptionType(self):
        return ExceptionType.BadAuthentication


class TGBadMagic(TGException):
    """This exception could occur if the server and client are not on the same version of the Protocol Data Unit or if
    the server or clients messages are corrupted."""

    @property
    def exceptionType(self):
        return ExceptionType.BadMagic


class TGBadVerb(TGException):
    """This could occur if the client or server send a wrong verb identifier over."""

    @property
    @abstractmethod
    def exceptionType(self):
        return ExceptionType.BadVerb


class TGChannelDisconnectedException(TGException):
    """Occurs when the client loses connection with the server."""
    @property
    def exceptionType(self):
        return ExceptionType.DisconnectedException


class TGConnectionTimeoutException(TGException):
    """Occurs if the connection lasts too long.

    If you are running into this exception frequently, you might want to wrap your code around a while loop with
    reinitializing the connection object.
    """

    @property
    def exceptionType(self):
        return ExceptionType.ConnectionTimeout


class TGInvalidMessageLength(TGException):
    """Message received ended before the API expected."""

    @property
    def exceptionType(self):
        return ExceptionType.InvalidMessageLength


class TGProtocolNotSupported(TGException):
    """The server's protocol is incompatible with this API's protocol.

    Might want to update either the server or this API.
    """

    @property
    def exceptionType(self):
        return ExceptionType.ProtocolNotSupported


class TGSecurityException(TGException):
    """Occurs when a security related issue occurs.

    Contact your IT or security team to get this worked out. If only working """
    pass


class TGTransactionResponse(enum.Enum):
    """More thorough transaction response than just error or success."""

    TGTransactionInvalid = -1
    TGTransactionSuccess = 0
    TGTransactionAlreadyInProgress = 8001
    TGTransactionClientDisconnected = 8002
    TGTransactionMalFormed = 8003
    TGTransactionGeneralError = 8004
    TGTransactionVerificationError = 8005
    TGTransactionInBadState = 8006
    TGTransactionUniqueConstraintViolation = 8007
    TGTransactionOptimisticLockFailed = 8008
    TGTransactionResourceExceeded = 8009
    TGCurrentThreadNotinTransaction = 8010
    TGTransactionUniqueIndexKeyAttributeNullError = 8011

    @classmethod
    def fromId(cls, id: int):
        resp: TGTransactionResponse
        for resp in TGTransactionResponse:
            if resp.value == id:
                return resp
        return TGTransactionResponse.TGTransactionInvalid


class TGTransactionException(TGException, ABC):
    """Baseclass for more detailed transaction-related exceptions.

     Designed to catch specific instances of transaction failure. For example, if you want to upsert an entity, you can
     first try to insert it, which if it is already there, will cause a TGTransactionUniqueConstraintViolation, which
     you can catch with an except statement specifying that exception.
     """

    @property
    @abstractmethod
    def transactionException(self) -> TGTransactionResponse:
        pass

    @classmethod
    def buildTransactionException(cls, msg: str, resp: TGTransactionResponse):
        known_txn_errors = list((txn_error.name for txn_error in TGTransactionResponse))
        if msg is not None:
            to_match = re.compile("Transaction:\d+ failed. Root Error:(" + "|".join(known_txn_errors) + ")(.*)")
            matching = to_match.match(msg)
            while resp is TGTransactionResponse.TGTransactionGeneralError and matching:
                resp = getattr(TGTransactionResponse, matching.group(1),
                               TGTransactionResponse.TGTransactionGeneralError)
                msg = matching.group(2)
                matching = to_match.match(msg)
        if resp is TGTransactionResponse.TGTransactionAlreadyInProgress:
            return TGTransactionAlreadyInProgressException(msg)
        elif resp is TGTransactionResponse.TGTransactionMalFormed:
            return TGTransactionMalFormedException(msg)
        elif resp is TGTransactionResponse.TGTransactionInBadState:
            return TGTransactionInBadStateException(msg)
        elif resp is TGTransactionResponse.TGTransactionVerificationError:
            return TGTransactionVerificationErrorException(msg)
        elif resp is TGTransactionResponse.TGTransactionUniqueConstraintViolation:
            return TGTransactionUniqueConstraintViolationException(msg)
        elif resp is TGTransactionResponse.TGTransactionOptimisticLockFailed:
            return TGTransactionOptimisticLockFailedException(msg)
        elif resp is TGTransactionResponse.TGTransactionResourceExceeded:
            return TGTransactionResourceExceededException(msg)
        elif resp is TGTransactionResponse.TGTransactionUniqueIndexKeyAttributeNullError:
            return TGTransactionUniqueIndexKeyAttributeNullError(msg)
        else:
            return TGTransactionGeneralErrorException(msg)


class TGTransactionAlreadyInProgressException(TGTransactionException):
    @property
    def transactionException(self):
        return TGTransactionResponse.TGTransactionAlreadyInProgress


class TGTransactionMalFormedException(TGTransactionException):
    @property
    def transactionException(self) -> TGTransactionResponse:
        return TGTransactionResponse.TGTransactionMalFormed


class TGTransactionGeneralErrorException(TGTransactionException):
    @property
    def transactionException(self) -> TGTransactionResponse:
        return TGTransactionResponse.TGTransactionGeneralError


class TGTransactionInBadStateException(TGTransactionException):
    @property
    def transactionException(self) -> TGTransactionResponse:
        return TGTransactionResponse.TGTransactionInBadState


class TGTransactionVerificationErrorException(TGTransactionException):
    @property
    def transactionException(self) -> TGTransactionResponse:
        return TGTransactionResponse.TGTransactionVerificationError


class TGTransactionUniqueConstraintViolationException(TGTransactionException):
    @property
    def transactionException(self) -> TGTransactionResponse:
        return TGTransactionResponse.TGTransactionUniqueConstraintViolation


class TGTransactionOptimisticLockFailedException(TGTransactionException):
    @property
    def transactionException(self) -> TGTransactionResponse:
        return TGTransactionResponse.TGTransactionOptimisticLockFailed


class TGTransactionResourceExceededException(TGTransactionException):
    @property
    def transactionException(self) -> TGTransactionResponse:
        return TGTransactionResponse.TGTransactionResourceExceeded


class TGTransactionUniqueIndexKeyAttributeNullError(TGTransactionException):
    @property
    def transactionException(self) -> TGTransactionResponse:
        return TGTransactionResponse.TGTransactionUniqueIndexKeyAttributeNullError


class TGImpExpException(TGException):
    """Occurs when an import/export related error occurs."""


"""
class TGBulkImpExpResponse(enum.Enum):
     More thorough bulk loading response than just error or success. ""
    TGBulkImpExpInvalid = -1
    TGBulkImpExpSuccess = 0
    TGBulkImpExpAlreadyInProgress = 9001            # TODO Construct legitimate values... (these are just temporary)
    TGBulkImpExpAttributeNotFoundError = 9002
    TGBulkImpExpPrimaryKeyColumnMissing = 9003
    TGBulkImpExpPrimaryKeyAttributeMissing = 9004
    TGBulkImpExpUniqueConstraintViolation = 9005
    TGBulkImpExpVertexIDNotFound = 9006
    TGBulkImpExpGeneralError = 9007

    @classmethod
    def fromId(cls, id: int):
        resp: TGBulkImpExpResponse
        for resp in TGBulkImpExpResponse:
            if resp.value == id:
                return resp
        return TGBulkImpExpResponse.TGBulkImpExpGeneralError


class TGBulkImpExpException(TGException, ABC):
    ""Base class for better error-handling of specific bulk loading errors.""

    @property
    @abstractmethod
    def response(self):
        pass

    @property
    def rowAssociative(self) -> bool:
        return False

    @classmethod
    def buildBulkImportException(cls, msg: str, resp: TGBulkImpExpResponse):
        if resp is TGBulkImpExpResponse.TGBulkImpExpAlreadyInProgress:
            return TGBulkImpExpAlreadyInProgressException(msg)
        elif resp is TGBulkImpExpResponse.TGBulkImpExpAttributeNotFoundError:
            return TGBulkImpExpAttributeNotFoundError(msg)
        elif resp is TGBulkImpExpResponse.TGBulkImpExpPrimaryKeyColumnMissing:
            return TGBulkImpExpPrimaryKeyColumnMissing(msg)
        elif resp is TGBulkImpExpResponse.TGBulkImpExpPrimaryKeyAttributeMissing:
            return TGBulkImpExpPrimaryKeyAttributeMissing(msg)
        elif resp is TGBulkImpExpResponse.TGBulkImpExpUniqueConstraintViolation:
            return TGBulkImpExpUniqueConstraintViolation(msg)
        elif resp is TGBulkImpExpResponse.TGBulkImpExpVertexIDNotFound:
            return TGBulkImpExpVertexIDNotFoundError(msg)
        else:
            return TGBulkImpExpGeneralError(msg)


class TGBulkImpExpGeneralError(TGBulkImpExpException):

    @property
    def response(self):
        return TGBulkImpExpResponse.TGBulkImpExpGeneralError


class TGBulkImpExpAlreadyInProgressException(TGBulkImpExpException):

    @property
    def response(self):
        return TGBulkImpExpResponse.TGBulkImpExpAlreadyInProgress


class TGBulkImpExpAttributeNotFoundError(TGBulkImpExpException):

    @property
    def response(self):
        return TGBulkImpExpResponse.TGBulkImpExpAttributeNotFoundError


class TGBulkImpExpPrimaryKeyColumnMissing(TGBulkImpExpException):

    @property
    def response(self):
        return TGBulkImpExpResponse.TGBulkImpExpPrimaryKeyColumnMissing


class TGBulkImpExpPrimaryKeyAttributeMissing(TGBulkImpExpException):

    def __init__(self, reason, errorcode=None, cause=None):
        super().__init__(reason, errorcode, cause)
        self.__rowNum__: int = 0

    @property
    def rowAssociative(self) -> bool:
        return True

    @property
    def rowNum(self) -> int:
        return self.__rowNum__

    @rowNum.setter
    def rowNum(self, val: int):
        self.__rowNum__ = val

    @property
    def response(self):
        return TGBulkImpExpResponse.TGBulkImpExpPrimaryKeyAttributeMissing


class TGBulkImpExpUniqueConstraintViolation(TGBulkImpExpException):

    def __init__(self, reason, errorcode=None, cause=None):
        super().__init__(reason, errorcode, cause)
        self.__rowNum__: int = 0

    @property
    def rowAssociative(self) -> bool:
        return True

    @property
    def rowNum(self) -> int:
        return self.__rowNum__

    @rowNum.setter
    def rowNum(self, val: int):
        self.__rowNum__ = val

    @property
    def response(self):
        return TGBulkImpExpResponse.TGBulkImpExpUniqueConstraintViolation


class TGBulkImpExpVertexIDNotFoundError(TGBulkImpExpException):

    def __init__(self, reason, errorcode=None, cause=None):
        super().__init__( reason, errorcode, cause )
        self.__rowNum__: int = 0

    @property
    def rowAssociative(self) -> bool:
        return True

    @property
    def rowNum(self) -> int:
        return self.__rowNum__

    @rowNum.setter
    def rowNum(self, val: int):
        self.__rowNum__ = val

    @property
    def response(self):
        return TGBulkImpExpResponse.TGBulkImpExpVertexIDNotFound
"""

class TGTypeCoercionNotSupported(TGException):
    """Attempted to convert from one type to one that is incompatible.

    Think trying to convert a number type to a single character.

    If running into this issue, please make sure that the types are compatible (i.e. explicitly convert to the type you
    want it to convert to).
    """

    def __init__(self, fromtype, totype, isarray=False):
        super(TGTypeCoercionNotSupported, self).__init__(str.format("Cannot coerce value of desc:{0} to desc:{1}{2}"
                                                                    .format(fromtype.name, totype.name,
                                                                            {True: "[]", False: ""}[isarray])))


class TGTypeNotSupported(TGException):
    """That type is not supported yet."""
    pass


class TGVersionMismatchException(TGException):
    """The version between the server and client is out of sync."""
    pass


class TGAdminParseException(TGException):
    """Error with parsing the input to the administrative client terminal."""

    def __init__(self, reason, cause: Exception = None, position: int = 0, token: str = None):
        super(TGAdminParseException, self).__init__(reason, cause)
        if position < 0:
            position = 0
        self.__pos = position
        self.__token = token if isinstance(token, str) and token is not None else ""

    @property
    def token(self) -> str:
        return self.__token

    @property
    def position(self) -> int:
        """The position that the error occurred in."""
        return self.__pos

    @position.setter
    def position(self, pos: int):
        """The position that the error occurred in."""
        self.__pos = pos

    def __str__(self):
        return str(super().__str__()) + (" position: %d" % self.__pos)

