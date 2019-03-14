/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : TGTypeNotSupported.${EXT}
 * Created on: 6/4/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.exception;

import com.tibco.tgdb.model.TGAttributeType;

public class TGTypeNotSupported extends TGException {

    public TGTypeNotSupported(TGAttributeType type)
    {
        super(String.format("Attribute desc:%s not supported", type.name()));
    }
}
