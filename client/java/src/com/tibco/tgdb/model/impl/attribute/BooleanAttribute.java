/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : BooleanAttribute.${EXT}
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

class BooleanAttribute extends AbstractAttribute {

    BooleanAttribute(TGAttributeDescriptor type) {
        super(type);
        this.value = false;
    }

    @Override
    public void setValue(Object value) throws TGException {
    	if (value == null) 
        {
    		this.value = value;
    		setModified();
        	return;
        }
    	
        if (value instanceof Boolean) {
            setBoolean((Boolean)value);
        }
        else if (value instanceof Number) {
            int v = (int)value;
            setBoolean(v > 0);
        }
        else if (value instanceof String) {
            setBoolean(Boolean.valueOf(value.toString()));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Boolean, value.getClass().getSimpleName());
        }
    }

    void setBoolean(Boolean b) {
        setModified();
    	if (!isNull() && (this.value.equals(b))) return;
        this.value = b;
        setModified();
    }

    @Override
    public boolean getAsBoolean() {
        return Boolean.class.cast(this.value);
    }

    @Override
    public byte getAsByte() {
        Boolean b = (Boolean) this.value;
        return (byte) (b ? 1 : 0);
    }

    @Override
    public char getAsChar() {
        Boolean b = (Boolean) this.value;
        return (char) (b ? 1 : 0);
    }

    @Override
    public short getAsShort() {
        Boolean b = (Boolean) this.value;
        return (short) (b ? 1 : 0);
    }

    @Override
    public int getAsInt() {
        Boolean b = (Boolean) this.value;
        return (int) (b ? 1 : 0);
    }

    @Override
    public long getAsLong() {
        Boolean b = (Boolean) this.value;
        return (long) (b ? 1L : 0L);
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeBoolean(Boolean.class.cast(value));
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        value = is.readBoolean();
    }
}
