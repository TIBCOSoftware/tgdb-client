/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : ConversionUtils.${EXT}
 * Created on: 6/22/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.utils.TGEnvironment;

import java.io.ByteArrayOutputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.math.BigDecimal;
import java.sql.Timestamp;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Calendar;

public class ConversionUtils {

    public static byte[] inputStream2ByteArray(InputStream in) throws TGException
    {
        ByteArrayOutputStream bos = new ByteArrayOutputStream();
        try {
            int cnt = 0;
            while ((cnt = in.available()) > 0) {
                byte buf[] = new byte[cnt];
                int readcnt = in.read(buf, 0, cnt);
                bos.write(buf, 0, readcnt);
            }
            return bos.toByteArray();
        }

        catch (IOException ioe)
        {
            throw new TGException(ioe);
        }

    }

    public static byte[] bigDecimal2ByteArray(BigDecimal bd) throws TGException
    {
        try {
            ByteArrayOutputStream bos = new ByteArrayOutputStream();
            DataOutputStream dos = new DataOutputStream(bos);
            dos.write(bd.unscaledValue().toByteArray());
            dos.writeInt(bd.scale());
            return bos.toByteArray();
        }
        catch (IOException ioe) {
            throw new TGException(ioe);
        }
    }

    public static Calendar string2Calendar(String s) throws TGException
    {
        try {
            //SS:TODO Get the Datetime format from the environment.
            //SS:Should we take per connection?
            Calendar cal = Calendar.getInstance();
            SimpleDateFormat sdf = new SimpleDateFormat(TGEnvironment.getInstance().getDefaultDateTimeFormat());
            cal.setTime(sdf.parse(s));
            return cal;
        }
        catch (ParseException pe) {
            throw new TGException(pe);
        }
    }

    public static Calendar long2Calendar(long l) throws TGException
    {
        Calendar cal = Calendar.getInstance();
        cal.setTimeInMillis(l);
        return cal;
    }

    public static Double string2Double(String s) throws TGException
    {
        try {
            return Double.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a double", s),
                    "ConversionUtils.Error", nfe);
        }
    }

    public static Float string2Float(String s) throws TGException
    {
        try {
            return Float.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a float", s),
                    "ConversionUtils.Error", nfe);
        }
    }

    public static Integer string2Integer(String s) throws TGException
    {
        try {
            return Integer.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a integer", s),
                    "ConversionUtils.Error", nfe);
        }
    }

    public static Long string2Long(String s) throws TGException
    {
        try {
            return Long.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a long", s),
                    "ConversionUtils.Error", nfe);
        }
    }

    public static Short string2Short(String s) throws TGException
    {
        try {
            return Short.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a short", s),
                    "ConversionUtils.Error", nfe);
        }
    }

    public static Character string2Character(String s) throws TGException
    {
        if (s.length() == 1)  {
            return (char)s.charAt(0);
        }
        else { //let see if it is unicode point.
            try {
                int v = Integer.parseInt(s, 16); //The value is expected to be in "FFFF"
                return (char) v;
            }catch (NumberFormatException nfe) {
                throw TGException.buildException(String.format("Cannot convert String:%s to a character", s),
                        "ConversionUtils.Error", nfe);
            }
        }
    }



}
