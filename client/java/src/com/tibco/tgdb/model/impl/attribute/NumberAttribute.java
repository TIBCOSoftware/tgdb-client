/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : NumberAttribute.${EXT}
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
import java.math.BigDecimal;
import java.math.BigInteger;
import java.math.MathContext;
import java.math.RoundingMode;

class NumberAttribute extends AbstractAttribute {

    NumberAttribute(TGAttributeDescriptor desc) {
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
        if (value instanceof BigDecimal) {
            setBigDecimal(BigDecimal.class.cast(value));
        }
        else if (value instanceof BigInteger) {
            setBigDecimal(new BigDecimal((BigInteger)value, new MathContext(desc.getPrecision())));
        }
        else if (value instanceof Number) { //BigDecimal is also included in this.
            Number n = (Number) value;
            setBigDecimal(new BigDecimal(n.doubleValue(), new MathContext(desc.getPrecision())));
        }
        else if (value instanceof String) {
            String s = (String) value;
            setBigDecimal(new BigDecimal(s, new MathContext(desc.getPrecision())));
        }
        else if (value instanceof CharSequence) {
            CharSequence cs = (CharSequence) value;
            char cbuf[] = cs.toString().toCharArray();
            setBigDecimal(new BigDecimal(cbuf, new MathContext(desc.getPrecision())));
        }
        else if (value instanceof char[]) {
            char cbuf[] = (char[]) value;
            setBigDecimal(new BigDecimal(cbuf, new MathContext(desc.getPrecision())));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Number, value.getClass().getSimpleName());
        }

    }

    private void setBigDecimal(BigDecimal bd) {
        if (!isNull() && this.value.equals(bd)) return;
        this.value = bd;
        setModified();
    }

    @Override
    public long getAsLong() {
        BigDecimal bd = (BigDecimal) value;
        return bd.longValue();
    }

    @Override
    public float getAsFloat() {
        BigDecimal bd = (BigDecimal) value;
        return bd.floatValue();

    }

    @Override
    public double getAsDouble() {
        BigDecimal bd = (BigDecimal) value;
        return bd.doubleValue();
    }

    @Override
    public BigDecimal getAsNumber() {
        BigDecimal bd = (BigDecimal) value;
        return bd;
    }

    @Override
    public byte[] getAsBytes() throws TGException {
        return ConversionUtils.bigDecimal2ByteArray(BigDecimal.class.cast(this.value));
    }



    @Override
    public void readValue(TGInputStream is) throws TGException, IOException
    {
        short precision = is.readShort();
        short scale = is.readShort();
        String bdstr = is.readUTF();
        value = new BigDecimal(bdstr);
        setPrecisionAndScale(precision, scale);
    }

    private void setPrecisionAndScale(short precision, short scale)
    {
        BigDecimal bd = BigDecimal.class.cast(this.value);
        MathContext mc = new MathContext(precision, RoundingMode.HALF_UP);
        BigDecimal newbd = bd.round(mc);
        this.value = newbd.setScale(scale, BigDecimal.ROUND_HALF_UP);
    }


    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeShort(desc.getPrecision());
        os.writeShort(desc.getScale());
        BigDecimal bd = (BigDecimal) value;
        String s = bd.toPlainString();
        os.writeUTF(s);
    }
}
