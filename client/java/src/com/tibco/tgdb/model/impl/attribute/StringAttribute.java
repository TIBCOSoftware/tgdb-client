/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : StringAttribute.${EXT}
 * Created on: 6/4/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.nio.CharBuffer;
import java.nio.charset.Charset;

class StringAttribute extends AbstractAttribute {
    final static int MAX_STRING_ATTR_LENGTH = 1000 - Short.BYTES - 1;

    StringAttribute(TGAttributeDescriptor desc)
    {
        super(desc);
    }

    @Override
    public void setValue(Object value) throws TGException {
        if (value == null) 
        {
        	this.value = value;
        	setModified();
        	return;
        }
        String s = value.toString();
        int len = utfLength(s);
        if (len > MAX_STRING_ATTR_LENGTH ) {
            throw new TGException(String.format("UTF length of String exceed the maximum string length supported for a String Attribute (%d > %d)",
                    len, MAX_STRING_ATTR_LENGTH));
        }
        this.value = s;
        setModified();
    }

    @Override
    public byte[] getAsBytes() {
        String str = (String) value;
        return str.getBytes();
    }

    @Override
    public char[] getAsChars() {
        String str = (String) value;
        CharBuffer cb = CharBuffer.wrap(str);
        return cb.array();

    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readUTF();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeUTF((String)this.value);
    }

    private int utfLength(String str) {
        int strlen = str.length();
        char c;
        int utflen = 0;
        for (int i = 0; i < strlen; i++) {
            c = str.charAt(i);
            if ((c >= 0x0001) && (c <= 0x007F)) {
                utflen++;
            } else if (c > 0x07FF) {
                utflen += 3;
            } else {
                utflen += 2;
            }
        }
        return utflen;
    }
}
