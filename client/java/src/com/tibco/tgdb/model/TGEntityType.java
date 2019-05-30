
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
 * File name : TGEntityType.java
 * Created by: suresh
 * Created on: 01/24/2015 
 * SVN Id: $Id: TGEntityType.java 3142 2019-04-26 00:15:06Z nimish $
 */

package com.tibco.tgdb.model;

import java.util.Collection;

import com.tibco.tgdb.pdu.TGSerializable;


public interface TGEntityType extends TGSystemObject {

    /**
     * List of all the attribute descriptors defined for this entity desc.
     * @return a collection of attribute descriptors for this entity desc
     */
    Collection<TGAttributeDescriptor> getAttributeDescriptors();

    /**
     * Get the attribute descriptor for the name specified
     * @param attrName the attribute name whose attribute descriptor is sought after.
     * @return TGAttributeDescriptor for the entity desc. Can be NULL.
     */
    TGAttributeDescriptor getAttributeDescriptor(String attrName);

    /**
     * Get the parent EntityType
     * @return optional derived from
     */
    TGEntityType derivedFrom();
}
