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
 *
 *  File name :TGEntity.java
 *  Created by: suresh
 *
 *		SVN Id: $Id: TGEntity.java 2473 2018-09-26 23:29:53Z vchung $
 *
 */
public interface TGEntity extends TGSerializable {

    enum TGEntityKind {
    	InvalidKind(0),
    	Entity(1), //FIXME: Should we call it something else? All the rest are entity also.
        Node(2),
        Edge(3),
        Graph(4),
        HyperEdge(5);

        private final int kind;

    	TGEntityKind(int kind) {
    		this.kind = kind;
    	}
    	
    	public int kind() {
    		return kind;
    	}

    	public static TGEntityKind fromValue(int kind) {
        	for (TGEntityKind ek : TGEntityKind.values()) {
            	if (kind == ek.kind) return ek;
        	}
        	return InvalidKind;
    	}
    }

    /**
     * Return the EntityKind
     * @return the kind of entity
     */
    TGEntityKind getEntityKind();

    /**
     * Check if this entity is currently being added to the system
     * @return boolean value indicating a newly created or fetched from the server
     */
    boolean isNew();

    /**
     * Check if this entity is already deleted in the system
     * @return boolean value indicating a deleted entity
     */
    boolean isDeleted();

    /**
     * Get the version of the Entity.
     * @return a int value as a version id.
     */
    int getVersion();

    /**
     * List of all the attributes set.
     * @return a collection of attributes for this entity
     */
    Collection<TGAttribute> getAttributes();

    /**
     * Get the attribute for the name specified
     * @param attrName the attribute name whose value and desc is sought after.
     * @return TGAttribute for the entity instance. Can be NULL.
     */
    TGAttribute getAttribute(String attrName);

    /**
     * Is TGAttribute set in the entity. It is check for the attributes existence and not if the value is NULL.
     * @param attrName the attribute name which is being tested as SET or not/
     * @return a boolean value indicating the same.
     */
    boolean isAttributeSet(String attrName);

    /**
     * Set the TGAttribute to this TGEntity
     * @param attr TGAttribute.
     */
    void setAttribute(TGAttribute attr);

    /**
     * Dynamically set the attribute to this entity.
     * If the AttributeDescriptor doesn't exist in the database, create a new one.
     * @param name - The attribute name
     * @param value - The value associated to it.
     * @throws com.tibco.tgdb.exception.TGException if can't set the value.
     */
    void setAttribute(String name, Object value) throws TGException;

    /**
     * Return the EntityType
     * @return the kind of entity
     */
    TGEntityType getEntityType();

}
