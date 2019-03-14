package com.tibco.tgdb.model;

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
 *  File name :TGAttribute.java
 *  Created on: 3/18/14
 *  Created by: suresh
 *
 *		SVN Id: $Id: TGAttribute.java 2348 2018-06-22 16:34:26Z ssubrama $
 *
 */

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.pdu.TGSerializable;

import java.math.BigDecimal;
import java.nio.ByteBuffer;
import java.nio.CharBuffer;
import java.util.Calendar;



/**
 * An attribute is simple scalar value that is associated with an TGEntity.
 */
public interface TGAttribute extends TGSerializable {

    /**
     * @deprecated
     * Return the TGAttributeDescriptor for this attribute
     * @return TGAttributeDescriptor
     */
    TGAttributeDescriptor getAttributeType();

    /**
     * Return the TGAttributeDescriptor for this attribute
     * @return TGAttributeDescriptor
     */
    TGAttributeDescriptor getAttributeDescriptor();

    /**
     * Get owner Entity.
     * @return owner Entity
     */
    TGEntity getOwner();

    /**
     * Is the attribute value null.
     * @return boolean true if it is null or false if it is not. If the attribute is not set, it will always be true
     */
    boolean isNull();

    /**
     * Is the attribute modified
     * @return boolean true if it is modified or false if it is not. 
     */
    boolean isModified();

    /**
     * Get the Value for this attribute as the most generic form
      * @return value of the object
     */
    Object getValue();

    /**
     * Set the value for this attribute. Appropriate data conversion to its attribute desc will be performed
     * If the object is Null, then the object is explicitly set, but no value is provided.
     * @param value The value to be set
     * @throws TGException throws exception if conversion cannot be performed, or constraint violation exception occured
     */
    void setValue(Object value) throws TGException;

    /**
     * Return the value as a Boolean value
     * @return boolean value
     */
    default boolean getAsBoolean() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), Boolean.class.getSimpleName());
    }

    /**
     * Return as java primitive byte
     * @return byte value
     */
    default byte getAsByte() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), Byte.class.getSimpleName());
    }

    /**
     * Return as java primitive char
     * @return  char value
     */
    default char getAsChar() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), Character.class.getSimpleName());
    }

    /**
     * Return as Short
     * @return short value
     */
    default short getAsShort() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), Short.class.getSimpleName());
    }

    /**
     * Return as int
     * @return int value
     */
    default int getAsInt() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), Integer.class.getSimpleName());
    }

    /**
     * Return as java long
     * @return long value
     */
    default long getAsLong() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), Long.class.getSimpleName());
    }

    /**
     * Return as float
     * @return float
     */
    default float getAsFloat() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), Float.class.getSimpleName());
    }

    /**
     * Return as a double
     * @return as double
     */

    default double getAsDouble() {throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), Double.class.getSimpleName()); }


    default String getAsString() {
        Object value = this.getValue();
        return value == null ? null : value.toString();
    }


    /**
     * Get as Date Calendar
     * @return
     */
    default Calendar getAsDate() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "Calendar.class");
    };

    /**
     * Get Time as Calendar.
     * @return
     */
    default Calendar getAsTime() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "Calendar.class");
    };

    /**
     * Get Timestamp (Datetime
     * @return
     */
    default Calendar getAsTimestamp() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "Calendar.class");
    }

    default BigDecimal getAsNumber() {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "BigDecimal.class");
    }


    /**
     * Get as Byte Array.
     * @return
     */
    default byte[] getAsBytes() throws TGException {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "byte array");
    }

    /**
     * Get as a Char Array
     * @return
     */
    default char[] getAsChars() throws TGException
    {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "char array");
    }

    default char[] getAsChars(String encoding) throws TGException
    {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "char array");
    }

    /**
     * Get As Blob
     * @return
     */
    default ByteBuffer getAsByteBuffer() throws TGException
    {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "ByteBuffer.class");
    }

    /**
     * Get as Clob
     * @return
     */
    default CharBuffer getAsCharBuffer() throws TGException
    {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "CharBuffer.class");
    }

    default CharBuffer getAsCharBuffer(String encoding) throws TGException
    {
        throw new TGTypeCoercionNotSupported(getAttributeDescriptor().getType(), "CharBuffer.class");
    }
}
