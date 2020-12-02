/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name : NumberAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: NumberAttribute.java 3631 2019-12-11 01:12:03Z ssubrama $
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
        this.value = bd.setScale(this.desc.getScale(), BigDecimal.ROUND_HALF_UP);
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
        BigDecimal bd1 = new BigDecimal(bdstr, new MathContext(precision, RoundingMode.HALF_UP));
        this.value = bd1.setScale(scale, BigDecimal.ROUND_HALF_UP);
        //setPrecisionAndScale(precision, scale);
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

    public static void main(String[] args) {

        BigDecimal bd = new BigDecimal("1234567.56");
        MathContext mc = new MathContext(10, RoundingMode.HALF_UP);
        System.out.printf("BigDecimal = %s\n", bd.toPlainString());
        BigDecimal newbd = bd.round(mc);
        System.out.printf("BigDecimal = %s\n", newbd.toPlainString());
        BigDecimal newbd2 = newbd.setScale(3, BigDecimal.ROUND_HALF_UP);
        System.out.printf("BigDecimal = %s\n", newbd2.toPlainString());

        BigDecimal bd1 = new BigDecimal("123456789567.5696", new MathContext(10, RoundingMode.HALF_UP));
        System.out.printf("BigDecimal = %s\n", bd1.toPlainString());
        BigDecimal bd2 = bd1.setScale(3,RoundingMode.HALF_UP);
        System.out.printf("BigDecimal = %s\n", bd2.toPlainString());

    }
}
