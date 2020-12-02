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
 *  File name :TGMutableGraphMetadata.java
 *  Created on: 08/22/2019
 *  Created by: nimish
 *  
 *  <p>This interface allows users to create metadata (such as NodeType, EdgeType, AttributeType) 
 *  programmatically.
 *  
 *  SVN Id: $Id: $
 * 
 */

package com.tibco.tgdb.admin;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGNodeType;

import java.util.Collection;


public interface TGMutableGraphMetadata extends TGGraphMetadata {
	
	/**
	 * Allows the users to create a NodeType instance
	 * 
	 * @param typeName name of the nodeType instance
	 * @param attrDescs collection of attribute descriptors
	 * @param pkeyDescs collection of primary-key descriptors
	 * @throws TGException if any error occurs while performing the operation on server
	 */
    public void createNodeType(String typeName, Collection<TGAttributeDescriptor> attrDescs, Collection<TGAttributeDescriptor> pkeyDescs) throws TGException;

	
    /**
     * Allows the users to create EdgeType instance
     * 
     * @param edgeTypeName name of the edgeType instance
     * @param directionality to specify the direction (e.g. if the edge needs to be UnDirected, Directed, or BiDirectional) 
     * @param attrDescs collection of attribute descriptors
     * @param sourceNodeType name of the source nodeType
     * @param destinationNodeType name of the destination nodeType
     * @throws TGException if any error occurs while performing the operation on server
     */
    public void createEdgeType(String edgeTypeName, TGEdge.DirectionType directionality, Collection<TGAttributeDescriptor> attrDescs, TGNodeType sourceNodeType, TGNodeType destinationNodeType) throws TGException;


    /**
     * Allows the users to create an AttributeType instance
     *
     * @param attributeDescriptor attribute descriptor instance
     * @throws TGException if any error occurs while performing the operation on server
     */
	public void createAttributeType (TGAttributeDescriptor attributeDescriptor) throws TGException;

}
