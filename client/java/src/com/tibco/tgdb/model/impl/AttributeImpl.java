/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * <p/>
 * File name : AttributeImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: AttributeImpl.java 2348 2018-06-22 16:34:26Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.math.BigDecimal;
import java.math.MathContext;
import java.math.RoundingMode;
import java.nio.ByteBuffer;
import java.nio.CharBuffer;
import java.util.Calendar;
import java.util.GregorianCalendar;

public class AttributeImpl implements TGAttribute {

    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();
    static final int DATE_ONLY = 0;
    static final int TIME_ONLY = 1;
    static final int TIMESTAMP = 2;
    static final int TGNoZone = -1;
    static final int TGZoneOffset = 0;
    static final int TGZoneId =1;
    static final int TGZoneName=2;


    AbstractEntity owner;

    TGAttributeDescriptor type;
    //addEdge does not trigger change notification
    Object value;
    boolean isModified = false;

    public AttributeImpl(AbstractEntity owner) {
        this.owner = owner;
    }

    public AttributeImpl(AbstractEntity owner, TGAttributeDescriptor type, Object value) {
        this.owner = owner;
        this.value = value;
        this.type = type;
    }

    @Override
    public TGEntity getOwner() {
        return owner;
    }

    @Override
    public TGAttributeDescriptor getAttributeType() {
        return type;
    }

    @Override
    public TGAttributeDescriptor getAttributeDescriptor() {
        return type;
    }

    @Override
    public boolean isNull() {
        return value==null;
    }

    @Override
    public Object getValue() {
        return value;
    }

    @Override
    public void setValue(Object value) throws TGException {
    	//FIXME: Need to match the desc of the attribute descriptor
        this.value = value;
        isModified = true;
        if (value == null) return;

        if (type.getType() == TGAttributeType.Number) {
            short precision = type.getPrecision();
            short scale = type.getScale();
            setPrecisionAndScale(precision, scale);
        }
    }

    @Override
    public boolean getAsBoolean() {
        return Boolean.class.cast(value);
    }

    @Override
    public byte getAsByte() {
        return Byte.class.cast(value);
    }

    @Override
    public char getAsChar() {
        return Character.class.cast(value);
    }

    @Override
    public short getAsShort() {
        return Short.class.cast(value);
    }

    @Override
    public int getAsInt() {
        return Integer.class.cast(value);
    }

    @Override
    public long getAsLong() {
        return Long.class.cast(value);
    }

    @Override
    public float getAsFloat() {
        return Float.class.cast(value);
    }

    @Override
    public double getAsDouble() {
        return Double.class.cast(value);
    }

    @Override
    public String getAsString() {
        return String.class.cast(value);
    }

    @Override
    public void writeExternal(TGOutputStream os) throws TGException, IOException {
       	int aid = type.getAttributeId();
        //null attribute is not allowed during entity creation
        os.writeInt(aid);
       	os.writeBoolean(isNull());
       	if (isNull()) {
       		return;
       	}
    	switch(type.getType()) {
    		case Boolean:
    			os.writeBoolean(Boolean.class.cast(value));
    			break;
    		case Byte:
    			os.writeByte(Byte.class.cast(value));
    			break;
    		case Char:
    			os.writeChar(Character.class.cast(value));
    			break;
    		case Short:
    			os.writeShort(Short.class.cast(value));
    			break;
    		case Integer:
    			os.writeInt(Integer.class.cast(value));
    			break;
    		case Long:
    			os.writeLong(Long.class.cast(value));
    			break;
    		case Float:
    			os.writeFloat(Float.class.cast(value));
    			break;
    		case Double:
    			os.writeDouble(Double.class.cast(value));
    			break;
    		case String:
    			os.writeUTF(String.class.cast(value));
    			break;
            case Date:
                writeTimestamp(os, DATE_ONLY);
                break;
            case Time:
                writeTimestamp(os, TIME_ONLY);
                break;
            case TimeStamp:
                writeTimestamp(os, TIMESTAMP);
                break;
            case Number:
                writeNumber(os);
                break;
    		default:
    			break;
    	}
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
    	int aid = is.readInt();
        TGAttributeDescriptor at = ((GraphMetadataImpl) owner.graphMetadata).getAttributeDescriptor(aid);
        this.type = at;
        if (at == null) {
        	//FIXME: retrieve entity desc together with the entity?
        	gLogger.log(TGLevel.Warning, "cannot lookup attribute descriptor %d from graph meta data cache", aid);
        }
        if (is.readByte() == 1) {
        	value = null;
        	return;
        }
    	switch(type.getType()) {
    		case Boolean:
    			byte b = is.readByte();
    			value = Boolean.valueOf(b != 0);
    			break;
    		case Byte:
    			value = is.readByte();
    			break;
    		case Char:
    			value = is.readChar();
    			break;
    		case Short:
    			value = is.readShort();
    			break;
    		case Integer:
    			value = is.readInt();
    			break;
    		case Long:
    			value = is.readLong();
    			break;
    		case Float:
    			value = is.readFloat();
    			break;
    		case Double:
    			value = is.readDouble();
    			break;
    		case String:
    			value = is.readUTF();
    			break;
            case Number:
            {
                readNumber(is);
                break;
            }
            case Date:
                readTimestamp(is, DATE_ONLY);
                break;
            case Time:
                readTimestamp(is, TIME_ONLY);
                break;
            case TimeStamp:
                readTimestamp(is, TIMESTAMP);
                break;

    		default:
    			break;
    	}
    }
    
    void resetIsModified() {
    	this.isModified = false;
    }
    
    @Override
    public boolean isModified() {
    	return isModified;
    }

    private void writeTimestamp(TGOutputStream os, int component2Write) throws TGException, IOException
    {
        Calendar cal = Calendar.class.cast(value);
        if (cal == null) throw new TGException("value is null");


        switch (component2Write) {
            case DATE_ONLY: {
                int era = cal.get(Calendar.ERA);
                os.writeBoolean(era == GregorianCalendar.AD);
                os.writeShort(cal.get(Calendar.YEAR));
                os.writeByte(cal.get(Calendar.MONTH) + 1); // Calendar starts January at 0, server starts January at 1
                os.writeByte(cal.get(Calendar.DAY_OF_MONTH));
                os.writeByte(0); //HR
                os.writeByte(0); //Min
                os.writeByte(0); //Sec
                os.writeShort(0); //msec
                os.writeByte(TGNoZone); //First to indicate we have no zone support
                //os.writeShort(TGNoZone); //This is for the zone ID
                break;
            }

            case TIME_ONLY: {
                os.writeBoolean(true);
                os.writeShort(0);
                os.writeByte(0);
                os.writeByte(0);
                os.writeByte(cal.get(Calendar.HOUR_OF_DAY)); //24 HR format
                os.writeByte(cal.get(Calendar.MINUTE)); //Min
                os.writeByte(cal.get(Calendar.SECOND)); //Sec
                os.writeShort(cal.get(Calendar.MILLISECOND)); //msec
                os.writeByte(TGNoZone); //First to indicate we have no zone support
                //os.writeShort(TGNoZone); //This is for the zone ID
                break;
            }

            case TIMESTAMP:
            {
                int era = cal.get(Calendar.ERA);
                os.writeBoolean(era == GregorianCalendar.AD);
                os.writeShort(cal.get(Calendar.YEAR));
                os.writeByte(cal.get(Calendar.MONTH) + 1); // Calendar starts January at 0, server starts January at 1
                os.writeByte(cal.get(Calendar.DAY_OF_MONTH));
                os.writeByte(cal.get(Calendar.HOUR_OF_DAY)); //24 HR format
                os.writeByte(cal.get(Calendar.MINUTE)); //Min
                os.writeByte(cal.get(Calendar.SECOND)); //Sec
                os.writeShort(cal.get(Calendar.MILLISECOND)); //msec
                os.writeByte(TGNoZone); //First to indicate we have no zone support
                //os.writeShort(TGNoZone); //This is for the zone ID
                break;
            }

            default:
                throw new TGException("Invalid spec provided to write the Calendar");
        }
    }

    //SS:TODO Support for Timezone - Post v1.0
    //SS:TODO Only support Gregorian Calendar - There is Japanese, Thai, ...
    private void readTimestamp(TGInputStream in, int component2read) throws TGException, IOException
    {
        boolean era;
        int year, mon, dom, hr, min, sec, ms, tztype, tzid;
        era     = in.readBoolean();
        year    = in.readShort();
        mon     = in.readByte();
        dom     = in.readByte();
        hr      = in.readByte();
        min     = in.readByte();
        sec     = in.readByte();
        ms      = in.readUnsignedShort();
        tztype  = in.readByte();

        if (tztype != -1) {
            tzid    = in.readShort();
        }

        --mon; // Java Calendar starts January at 0, server side starts January at 1

        switch (component2read) {
            case DATE_ONLY:
                value = new Calendar.Builder().setCalendarType("gregory")
                        .set(Calendar.ERA, era ? GregorianCalendar.AD : GregorianCalendar.BC)
                        .setDate(year,mon,dom)
                        .setTimeOfDay(0,0,0,0)
                        .build();

                break;
            case TIME_ONLY:
                value = new Calendar.Builder().setCalendarType("gregory")
                        .set(Calendar.ERA, 1) //1 is AD.
                        .setDate(0,0,0)
                        .setTimeOfDay(hr,min,sec,ms)
                        .build();

                break;
            case TIMESTAMP:
                value = new Calendar.Builder().setCalendarType("gregory")
                        .set(Calendar.ERA, era ? GregorianCalendar.AD : GregorianCalendar.BC)
                        .setDate(year,mon,dom)
                        .setTimeOfDay(hr,min,sec,ms)
                        .build();
                break;
            default:
                throw new TGException("Invalid spec provided to read the Calendar");

        }

    }

    private void writeNumber(TGOutputStream os) throws TGException, IOException
    {
        os.writeShort(type.getPrecision());
        os.writeShort(type.getScale());
        BigDecimal bd = (BigDecimal) value;
        String s = bd.toPlainString();
        os.writeUTF(s);
    }

    private void readNumber(TGInputStream in) throws TGException, IOException
    {
        short precision = in.readShort();
        short scale = in.readShort();
        String bdstr = in.readUTF();
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


}
