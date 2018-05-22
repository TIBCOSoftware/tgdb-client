/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : TGChannelDisconnectedException.${EXT}
 * Created on: 3/14/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.exception;

public class TGChannelDisconnectedException extends TGException {

    public TGChannelDisconnectedException(Exception cause) {
        super(cause);
    }
    public TGChannelDisconnectedException(String msg) {
        super(msg);
    }
}
