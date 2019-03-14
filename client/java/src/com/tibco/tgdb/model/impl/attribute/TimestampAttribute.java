/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : TimestampAttribute.${EXT}
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
import com.tibco.tgdb.utils.TGEnvironment;

import java.io.IOException;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.GregorianCalendar;

class TimestampAttribute extends AbstractAttribute {
    static final int TGNoZone = -1;
    static final int TGZoneOffset = 0;
    static final int TGZoneId =1;
    static final int TGZoneName=2;

    TimestampAttribute(TGAttributeDescriptor desc)
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

        if (value instanceof Calendar) {
            setCalendar((Calendar) value);
        }
        else if (value instanceof String) {
            setCalendar(ConversionUtils.string2Calendar(String.class.cast(value)));
        }
        else if (value instanceof Long) {
            setCalendar(ConversionUtils.long2Calendar(Long.class.cast(value)));
        }
        else {
            throw new TGTypeCoercionNotSupported(desc.getType(), value.getClass().getSimpleName());
        }
    }

    void setCalendar(Calendar cal) {
        if (!isNull() && this.value.equals(cal)) return;
        this.value = cal;
        setModified();
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        boolean era;
        int year, mon, dom, hr, min, sec, ms, tztype, tzid;
        era     = is.readBoolean();
        year    = is.readShort();
        mon     = is.readByte();--mon;
        dom     = is.readByte();
        hr      = is.readByte();
        min     = is.readByte();
        sec     = is.readByte();
        ms      = is.readUnsignedShort();
        tztype  = is.readByte();

        if (tztype != -1) {
            tzid    = is.readShort();
        }

        switch (this.desc.getType()) {
            case Date:
                value = new Calendar.Builder().setCalendarType("gregory")
                        .set(Calendar.ERA, era ? GregorianCalendar.AD : GregorianCalendar.BC)
                        .setDate(year,mon,dom)
                        .setTimeOfDay(0,0,0,0)
                        .build();

                break;
            case Time:
                value = new Calendar.Builder().setCalendarType("gregory")
                        .set(Calendar.ERA, 1) //1 is AD.
                        .setDate(0,0,0)
                        .setTimeOfDay(hr,min,sec,ms)
                        .build();

                break;
            case TimeStamp:
                value = new Calendar.Builder().setCalendarType("gregory")
                        .set(Calendar.ERA, era ? GregorianCalendar.AD : GregorianCalendar.BC)
                        .setDate(year,mon,dom)
                        .setTimeOfDay(hr,min,sec,ms)
                        .build();
                break;
            default:
                throw new TGException(String.format("Bad Descriptor :%s", this.desc.getType()));
        }
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        Calendar cal = Calendar.class.cast(value);
        int era = cal.get(Calendar.ERA);
        switch (this.desc.getType()) {
            case Date:
                os.writeBoolean(era == GregorianCalendar.AD);
                os.writeShort(cal.get(Calendar.YEAR));
                os.writeByte(cal.get(Calendar.MONTH) + 1); // Calendar starts January at 0, server starts January at 1
                os.writeByte(cal.get(Calendar.DAY_OF_MONTH));
                os.writeByte(0); //HR
                os.writeByte(0); //Min
                os.writeByte(0); //Sec
                os.writeShort(0); //msec
                os.writeByte(TGNoZone); //First to indicate we have no zone support
                break;

            case Time:
                os.writeBoolean(true);
                os.writeShort(0);
                os.writeByte(0);
                os.writeByte(0);
                os.writeByte(cal.get(Calendar.HOUR_OF_DAY)); //24 HR format
                os.writeByte(cal.get(Calendar.MINUTE)); //Min
                os.writeByte(cal.get(Calendar.SECOND)); //Sec
                os.writeShort(cal.get(Calendar.MILLISECOND)); //msec
                os.writeByte(TGNoZone); //First to indicate we have no zone support
                break;

            case TimeStamp:
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
            default:
                throw new TGException(String.format("Bad Descriptor :%s", this.desc.getType()));
        }
    }



}
