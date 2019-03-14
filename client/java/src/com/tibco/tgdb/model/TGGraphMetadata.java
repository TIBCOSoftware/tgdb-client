package com.tibco.tgdb.model;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGSerializable;

import java.util.Collection;

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

 * File name : TGGraphMetadata.java
 * Created on: 1/22/15
 * Created by: suresh
 * Graph Metadata provides the catalogue data of the server.

 * SVN Id: $Id: TGGraphMetadata.java 2344 2018-06-11 23:21:45Z ssubrama $
 */

public interface TGGraphMetadata extends TGSerializable {

    /**
     * Return a set of Node Type defined in the System
     * @return Collection of TGNodeType. Can be Empty Collection
     */
    Collection<TGNodeType> getNodeTypes();

    /**
     * Get Node Type by name.
     * @param typeName a valid typename in the database catalogue
     * @return TGNodeType
     */
    TGNodeType getNodeType(String typeName);

    /**
     * Return a set of know edge Type
     * @return a Collection of TGEdgeType . Can be Empty Collection
     */
    Collection<TGEdgeType> getEdgeTypes();

    /**
     * Return the EdgeType by name
     * @param typeName a valid typename in the database catlogue
     * @return TGEdgeType
     */
    TGEdgeType getEdgeType(String typeName);

    /**
     * List of Attribute Descriptors
     * @return a collection of AttributeDescriptors
     */
    Collection<TGAttributeDescriptor> getAttributeDescriptors();


    /**
     * Get the Attribute Descriptor
     * @param attributeName the attribute nae for the Attribute Descriptor
     * @return A valid AttributeDescriptor
     */
    TGAttributeDescriptor getAttributeDescriptor(String attributeName);


    /**
     * Create Attribute Descriptor
     * @param attrName the attribute name for the descriptor
     * @param attrType the desc of the attribute from the enum
     * @param isArray is it an array
     * @return a a newly created TGAttributeDescriptor
     * @throws TGException if the attribute descriptor already exists
     */
    TGAttributeDescriptor createAttributeDescriptor(String attrName, TGAttributeType attrType, boolean isArray) throws TGException;

}
