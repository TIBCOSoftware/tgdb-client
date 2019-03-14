package com.tibco.tgdb.exception;

import com.tibco.tgdb.pdu.impl.ExceptionMessage;

import java.io.IOException;

/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * File name :TGException
 * Created by: suresh
 *
 * SVN Id: $Id: TGException.java 2344 2018-06-11 23:21:45Z ssubrama $
 */
public class TGException extends Exception {


    String errorCode;
    Exception linkedException;

    /**
     * Build TGException from a linked exception
     * @param reason    The reason for this exception being wrapped
     * @param errorCode An optional error code from the System
     * @param linkedException - The linked exception that is being wrapped
     * @return a new TGException
     */
    public static TGException buildException(String reason, String errorCode, Exception linkedException) {
        TGException ex = new TGException(reason, errorCode);
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
     * Create a new TGException from String
     * @param ex the descriptive message for the Exception
     */
    public TGException(String ex) {
        super(ex);
    }

    /**
     * Create a new TGException from Reason and Error code
     * @param reason - A descriptive reason for the error
     * @param errorCode - Optional Error code.
     */
    public TGException(String reason, String errorCode) {
        super(reason);
        this.errorCode = errorCode;
    }

    public TGException(Exception cause) {
        super(cause);
    }
    /**
     * @return Returns the error code for the exception
     */
    public String getErrorCode() { return errorCode; }

    /**
     * @return the LinkedException if it wrapped from one. Can be null.
     */
    public Exception getLinkedException() { return linkedException; }

    /**
     * @return A Exception desc from the the Exception.
     */
    public ExceptionType getExceptionType() {
        if (linkedException != null) {
            if (linkedException instanceof IOException) {
                return ExceptionType.IOException;
            }
            return ExceptionType.GeneralException;
        }
        return ExceptionType.GeneralException;
    }

    public enum ExceptionType {
        BadVerb,
        InvalidMessageLength,
        BadMagic,
        ProtocolNotSupported,
        BadAuthentication,
        IOException,
        ConnectionTimeout,
        GeneralException,
        RetryIOException,
        DisconnectedException
    }

    @Override
    public synchronized Throwable getCause() {
        if (linkedException != null) return linkedException;
        return super.getCause();
    }
}
