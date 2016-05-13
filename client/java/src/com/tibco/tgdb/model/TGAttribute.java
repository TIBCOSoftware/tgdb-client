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
 *		SVN Id: $Id: TGAttribute.java 771 2016-05-05 11:40:52Z vchung $
 *
 */

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGSerializable;

/**
 * An attribute is simple scalar value that is associated with an TGEntity.
 */
public interface TGAttribute extends TGSerializable {

    /**
     * Return the TGAttributeDescriptor for this attribute
     * @return TGAttributeDescriptor
     */
    TGAttributeDescriptor getAttributeType();

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
     * Set the value for this attribute. Appropriate data conversion to its attribute type will be performed
     * If the object is Null, then the object is explicitly set, but no value is provided.
     * @param value The value to be set
     * @throws TGException throws exception if conversion cannot be performed, or constraint violation exception occured
     */
    void setValue(Object value) throws TGException;

    /**
     * Return the value as a Boolean value
     * @return boolean value
     */
    boolean getAsBoolean();

    /**
     * Return as java primitive byte
     * @return byte value
     */
    byte getAsByte();

    /**
     * Return as java primitive char
     * @return  char value
     */
    char getAsChar();

    /**
     * Return as Short
     * @return short value
     */
    short getAsShort();

    /**
     * Return as int
     * @return int value
     */
    int getAsInt();

    /**
     * Return as java long
     * @return long value
     */
    long getAsLong();

    /**
     * Return as float
     * @return float
     */
    float getAsFloat();

    /**
     * Return as a double
     * @return as double
     */
    double getAsDouble();

    /**
     * Return the String value
     * @return string value
     */
    String getAsString();
}
