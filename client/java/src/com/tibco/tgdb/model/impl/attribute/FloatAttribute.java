/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : FloatAttribute.${EXT}
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

class FloatAttribute extends AbstractAttribute {

    FloatAttribute(TGAttributeDescriptor desc)
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

        if (value instanceof Float) {
            setFloat((Float) value);
        }
        else if (value instanceof Number) {
            setFloat(((Number)value).floatValue());
        }
        else if (value instanceof String) {
            setFloat(ConversionUtils.string2Float(String.class.cast(value)));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Float, value.getClass().getSimpleName());
        }
    }

    void setFloat(Float f) {
        if (!isNull() && value.equals(f)) return;
        value = f;
        setModified();
    }

    @Override
    public boolean getAsBoolean() {
        if (isNull()) return false;
        float d = (Float) this.value;
        return d > 0;
    }


    @Override
    public short getAsShort() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.shortValue();
    }

    @Override
    public int getAsInt() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.intValue();
    }

    @Override
    public long getAsLong() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.longValue();
    }

    @Override
    public float getAsFloat() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.floatValue();
    }

    @Override
    public double getAsDouble() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.doubleValue();
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readFloat();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeFloat((Float)this.value);
    }
}
