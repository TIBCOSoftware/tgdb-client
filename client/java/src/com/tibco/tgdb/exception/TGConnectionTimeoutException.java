/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : TGConnectionTimeoutException.${EXT}
 * Created on: 3/7/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.exception;

public class TGConnectionTimeoutException extends TGException {

    public TGConnectionTimeoutException(String errmsg) {
        super(errmsg);
    }

    public ExceptionType getExceptionType() {
        return ExceptionType.ConnectionTimeout;
    }
}
