package com.tibco.tgdb.model;

import java.math.BigDecimal;

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
 *
 * File name :TGAttributeType
 * Created by: suresh

 * SVN Id: $Id: TGAttributeType.java 2344 2018-06-11 23:21:45Z ssubrama $
 */

/**
 * Attribute Type enumerations.
 * Do change the ID, as the server expects the attribute desc as the same value.
 */
public enum TGAttributeType {

    Invalid(0, null, null),
    Boolean(1, "Boolean", java.lang.Boolean.class),  //A single bit representing the truth value
    Byte(2, "Byte", java.lang.Byte.class),     //8bit octet
    Char(3, "Char", java.lang.Character.class),    //Fixed 8-Bit octet of N length
    Short(4, "Short", java.lang.Short.class),    //16bit
    Integer(5, "Integer", java.lang.Integer.class),      //32bit signed integer
    Long(6,"Long", java.lang.Long.class),     //64bit
    Float(7, "Float", java.lang.Float.class),    //32bit float
    Double(8, "Double", java.lang.Double.class),   //64bit float
    Number(9, "Number", BigDecimal.class),   //Number with precision
    String(10, "String", java.lang.String.class),   //Varying length String < 64K
    Date(11, "Date", java.util.Calendar.class),     //Only the Date part from the Calendar class is considered
    Time(12, "Time", java.util.Calendar.class), //Time hh:mm:ss.nnn with TIMEZONE
    TimeStamp(13, "Timestamp", java.util.Calendar.class), //Compatible with the SQL Timestamp field.
    Clob(14, "Clob", char[].class),     //Character -UTF-8 encoded string or large length > 64K
    Blob(15, "Blob", byte[].class)     //Binary object - a stream of octets (unsigned 8bit char) with length. A variation of such Blobs could
    ;

    TGAttributeType(int typeid, String typename, Class klazz) {
        this.typeid = typeid;
        this.typename = typename;
        this.klazz = klazz;
    }

    public static TGAttributeType fromTypeId(int typeid) {

        for (TGAttributeType attrType : TGAttributeType.values()) {
            if (attrType.typeid == typeid) return attrType;
        }
        return TGAttributeType.Invalid;
    }

    public static TGAttributeType fromTypeName(String typename) {

        for (TGAttributeType attrType : TGAttributeType.values()) {
            if (attrType.typename.equalsIgnoreCase(typename)) return attrType;
        }
        return TGAttributeType.Invalid;
    }

    public static TGAttributeType fromClass(Class klazz) {
    	// The null check is used to skip the first 'invalid' desc
        if (klazz.equals(java.util.Calendar.class)) return TimeStamp;  //We return the bigger container field.
        for (TGAttributeType attrType : TGAttributeType.values()) {
            if (attrType.klazz != null && attrType.klazz.equals(klazz)) return attrType;
        }
        return TGAttributeType.Invalid;
    }

    public int typeId() {
    	return typeid;
    }

    private int typeid;
    private String typename;
    private Class klazz;

}
