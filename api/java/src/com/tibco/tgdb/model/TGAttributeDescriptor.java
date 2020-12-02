
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
 *  File name :TGAttributeDescriptor.java
 *  Created on: 3/18/14
 *  Created by: suresh
 *
 *		SVN Id: $Id: TGAttributeDescriptor.java 3631 2019-12-11 01:12:03Z ssubrama $
 *
 */

package com.tibco.tgdb.model;


import com.tibco.tgdb.pdu.TGSerializable;

/**
 * Basic definition of an TGAttribute
 */
public interface TGAttributeDescriptor extends TGSystemObject {

    /**
     * Return the AttributeId
     * @return the AttributeId
     */
    int getAttributeId();

    /**
     * Return the desc of TGAttribute
     * @return the attribute desc
     */
    TGAttributeType getType();

    /**
     * Is the AttributeType an array desc
     * @return boolean indicating the multiplicativeness of the desc
     */
    boolean isArray();

    /**
     * Is this field Encrypted?
     * @return
     */
    boolean isEncrypted();

    /**
     * For a Number desc - return the precision. The default Precision is 20
     * @return
     */
    short getPrecision();

    /**
     * For a Number desc - return the scale. The default Scale is 5
     * @return
     */
    short getScale();

}
