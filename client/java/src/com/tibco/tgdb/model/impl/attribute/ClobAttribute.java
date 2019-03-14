/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : ClobAttribute.${EXT}
 * Created on: 6/4/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;

import java.nio.ByteBuffer;
import java.nio.CharBuffer;
import java.nio.charset.Charset;
import java.util.Iterator;
import java.util.Set;
import java.util.SortedMap;

class ClobAttribute extends BlobAttribute  {

    ClobAttribute(TGAttributeDescriptor desc)
    {
        super(desc);
    }

    public void setvalue(Object value) throws TGException
    {
    	if (value == null) 
        {
    		this.value = value;
    		setModified();
        	return;
        }
        else if (value instanceof char[]) {
            setCharBuffer(CharBuffer.wrap((char[])value));
        }
        else if (value instanceof CharBuffer) {
            setCharBuffer((CharBuffer) value);
        }
        else if (value instanceof CharSequence) {
            setCharBuffer(CharBuffer.wrap((CharSequence)value));
        }
        else {
            super.setValue(value);
        }
    }

    void setCharBuffer(CharBuffer cb)
    {
        Charset cs = Charset.forName("UTF-8");
        ByteBuffer bb = cs.encode(cb);
        this.value = bb.array();
        setModified();
    }

    @Override
    public char[] getAsChars() throws TGException {
        return getAsChars("UTF-8");

    }

    @Override
    public CharBuffer getAsCharBuffer() throws TGException {
        return getAsCharBuffer("UTF-8");
    }

    public char[] getAsChars(String encoding) throws TGException
    {
        CharBuffer cb = getAsCharBuffer(encoding);
        return cb.array();
    }

    public CharBuffer getAsCharBuffer(String encoding) throws TGException
    {
        ByteBuffer bb = getAsByteBuffer();
        Charset cs = Charset.forName(encoding);
        CharBuffer cb = cs.decode(bb);
        return cb;
    }

    public static void main(String[] args)
    {
        SortedMap charsets = Charset.availableCharsets();
        Set names = charsets.keySet();
        for (Iterator e = names.iterator(); e.hasNext();) {
            String name = (String) e.next();
            Charset charset = (Charset) charsets.get(name);
            System.out.println(charset);
            Set aliases = charset.aliases();
            for (Iterator ee = aliases.iterator(); ee.hasNext();) {
                System.out.println("    " + ee.next());
            }
        }
        Charset cs = Charset.forName("737");
        System.out.println("Charset is : " + cs);

    }




}
