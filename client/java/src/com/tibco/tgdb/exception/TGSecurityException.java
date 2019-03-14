/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : TGSecurityException.${EXT}
 * Created on: 4/5/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.exception;

public class TGSecurityException extends TGException {

    public TGSecurityException(String msg)
    {
        super(msg);
    }

    public TGSecurityException(String reason, String errorCode)
    {
        super(reason, errorCode);
    }
}
