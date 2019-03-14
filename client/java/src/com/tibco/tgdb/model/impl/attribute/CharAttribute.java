/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : CharAttribute.${EXT}
 * Created on: 6/4/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;

class CharAttribute extends AbstractAttribute {

    CharAttribute (TGAttributeDescriptor desc)
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
        
        if (value instanceof Character) {
            setChar((Character)value);
        }
        else if (value instanceof Byte) {
            byte b = (Byte) value;
            setChar((char) b);
        }
        else if (value instanceof Number) {
            Number n = (Number) value;
            setChar((char) n.shortValue());
        }
        else if (value instanceof String) { //Unicode String value
            String s = value.toString();
            setChar(ConversionUtils.string2Character(s));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Char, value.getClass().getSimpleName());
        }

    }

    void setChar(Character c) {
        if (!isNull() && this.value.equals(c)) return;
        this.value = c;
        setModified();
    }

    @Override
    public char getAsChar() {
        return (char) this.value;
    }

    @Override
    public byte getAsByte() {
        return (byte)Character.getNumericValue((char)this.value);

    }

    @Override
    public short getAsShort() {
        return (short)Character.getNumericValue((char)this.value);
    }

    @Override
    public int getAsInt() {
        return Character.getNumericValue((char)this.value);
    }

    @Override
    public long getAsLong() {
        return Character.getNumericValue((char)this.value);
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readChar();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeChar((char)this.value);
    }
}
