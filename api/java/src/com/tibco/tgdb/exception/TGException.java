/**
 * Copyright (c) 2019 TIBCO Software Inc.
 * All rights reserved.
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
 * <p/>
 * File name: TGException.java
 * Created on: 2014-12-17
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGException.java 4619 2020-10-31 17:37:35Z vchung $
 */

package com.tibco.tgdb.exception;

import com.tibco.tgdb.pdu.impl.ExceptionMessage;
import java.io.IOException;

public class TGException extends Exception {
    public static int SvrCodeNotDefined = -1;

    private Exception linkedException;
    private TGExceptionType exceptionType = TGExceptionType.GeneralException;
    private int serverErrorCode = SvrCodeNotDefined;

    public enum TGExceptionType {
        GeneralException,
        BadVerb,
        BadMagic,
        HandshakeResponseError,
        InvalidMessageLength,
        ProtocolNotSupported,
        SSLInitError,
        BadAuthentication,
        ChannelSendError,
        ChannelError,
        IOException,
        ConnectionTimeout,
        RetryIOException,
        DisconnectedException,
        ThreadInterrupted,
        DuplicateSystemObject,
        BadAttributeDescriptor,
        TypeConversionError,
        UnexpectedOperation,
        TxnAlreadyInProgress,
        TxnClientDisconnected,
        TxnMalFormed,
        TxnGeneralError,
        TxnVerificationError,
        TxnInBadState,
        TxnUniqueConstraintViolation,
        TxnOptimisticLockFailed,
        TxnResourceExceeded,
        TxnCurrentThreadNotinTxn,
        TxnUniqueIndexKeyAttributeNullError,
        TxnInvalidPhase,
        TxnWriteAheadLogBusy,
        TxnWriteAheadLogFailed2Commit,
        QryError, //Generic query error
        QryProviderNotInitialized,
        QryParsingError,
        QryStepNotSupported,
        QryStepNotAllowed,
        QryStepArgMissing,
        QryStepArgNotSupported,
        QryStepMissing,
        QryNotDefined,
        QryAttrDescNotFound,
        QryEdgeTypeNotFound,
        QryNodeTypeNotFound,
        QryInternalDataMismatchError,
        QryStepSignatureNotSupported,
        QryInvalidDataType,
        QryExecSPFailure,
        QrySPNotFound,
        QrySPArgMissing,
        QryStepArgInvalid,
        QryStepModulationInvalid,
        QryAccessDenied;
    }

    /**
     * Create a new TGException from String
     * @param ex the descriptive message for the Exception
     */
    public TGException(String ex) {
        super(ex);
        exceptionType = deriveExceptionType();
    }

    /**
     * Create a new TGException from Reason and Error code
     * @param reason - A descriptive reason for the error
     * @param excpType - Optional Error code.
     */
    public TGException(String reason, TGExceptionType excpType) {
        super(reason);
        this.exceptionType = excpType;
    }

    /**
     * Create a new TGException from Reason and Error code
     * @param reason - A descriptive reason for the error
     * @param excpType - Optional Error code.
     * @param serverCode - Server side error code
     */
    public TGException(String reason, TGExceptionType excpType, int serverCode) {
        super(reason);
        this.exceptionType = excpType;
        this.serverErrorCode = serverCode;
    }

    /**
     * Create a new TGException from a Java exception
     * @param cause
     */
    public TGException(Exception cause) {
        super(cause);
        exceptionType = deriveExceptionType();
    }

    /**
     * @return the LinkedException if it wrapped from one. Can be null.
     */
    public Exception getLinkedException() { return linkedException; }

    /**
     * Determine the exception type based on the linkedException type
     * @return Exception type
     */
    TGExceptionType deriveExceptionType() {
        if (linkedException != null) {
            if (linkedException instanceof IOException) {
                return TGExceptionType.IOException;
            }
        }
        return TGExceptionType.GeneralException;
    }

    /**
     * Build TGException from a linked exception
     * @param reason    The reason for this exception being wrapped
     * @param excpType An optional error code from the System
     * @param svrErrCode Server error code if server error is mapped to a different excpType
     * @param linkedException - The linked exception that is being wrapped
     * @return a new TGException
     */
    public static TGException buildException(String reason, TGExceptionType excpType, int svrErrCode,
                                             Exception linkedException) {
        TGException ex = new TGException(reason, excpType);
        ex.linkedException = linkedException;
        ex.serverErrorCode = svrErrCode;
        return ex;
    }

    /**
     * Build TGException from a linked exception
     * @param reason    The reason for this exception being wrapped
     * @param excpType An optional error code from the System
     * @param linkedException - The linked exception that is being wrapped
     * @return a new TGException
     */
    public static TGException buildException(String reason, TGExceptionType excpType, Exception linkedException) {
        TGException ex = new TGException(reason, excpType);
        ex.linkedException = linkedException;
        return ex;
    }

    /**
     * Build TGException from a wire level exception message. This message is received from the server.
     * @param msg - A Exception message received from the Server
     * @return A new TGException from the wire.
     */
    public static TGException buildException(ExceptionMessage msg) {
        TGException ex = new TGException(msg.getMessage());
        return ex;
    }

    /**
     * @return TGExceptionType
     */
    public TGExceptionType getExceptionType() {
    	return this.exceptionType;
    }

    @Override
    public synchronized Throwable getCause() {
        if (linkedException != null) return linkedException;
        return super.getCause();
    }

    /**
     * @return Server error code
     */
    public int getServerErrorCode() {
        return serverErrorCode;
    }
}
