/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : IntegerAttribute.${EXT}
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

class IntegerAttribute extends AbstractAttribute {

    IntegerAttribute(TGAttributeDescriptor desc) {
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
        if (value instanceof Integer) {
            setInteger((Integer) value);
        }
        else if (value instanceof Number) {
            setInteger(((Number)value).intValue());
        }
        else if (value instanceof String) {
            setInteger(ConversionUtils.string2Integer(String.class.cast(value)));
        }
        else if (value instanceof Character) {
            char c = (char) value;
            setInteger(Character.getNumericValue(c));
        }
        else if (value instanceof Boolean) {
            boolean b = (boolean) value;
            setInteger(b ? 1 : 0);
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Integer, value.getClass().getSimpleName());
        }
    }

    void setInteger(Integer i) {
        if (!isNull() && this.value.equals(i)) return;
        this.value = i;
        setModified();
    }

    @Override
    public boolean getAsBoolean() {
        if (isNull()) return false;
        int d = (Integer) this.value;
        return d > 0;
    }


    @Override
    public short getAsShort() {
        if (isNull()) return 0;
        Integer d = (Integer) this.value;
        return d.shortValue();
    }

    @Override
    public int getAsInt() {
        if (isNull()) return 0;
        Integer d = (Integer) this.value;
        return d.intValue();
    }

    @Override
    public long getAsLong() {
        if (isNull()) return 0;
        Integer d = (Integer) this.value;
        return d.longValue();
    }

    @Override
    public float getAsFloat() {
        if (isNull()) return 0;
        Integer d = (Integer) this.value;
        return d.floatValue();
    }

    @Override
    public double getAsDouble() {
        if (isNull()) return 0;
        Integer d = (Integer) this.value;
        return d.doubleValue();
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readInt();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeInt((Integer)this.value);
    }
}
