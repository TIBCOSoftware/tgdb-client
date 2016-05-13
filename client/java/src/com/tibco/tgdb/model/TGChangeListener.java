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
 * <p>
 * File name : TGChangeListener.${EXT}
 * Created on: 3/30/15
 * Created by: suresh
 * <p>
 * SVN Id: $Id: TGChangeListener.java 623 2016-03-19 21:41:13Z ssubrama $
 */

import com.tibco.tgdb.model.impl.TGGraph;

/**
 * For Internal Use only...
 */
public interface TGChangeListener {


    /**
     * New Entity Created.
     * @param entity the entity
     */
    void entityCreated(TGEntity entity);


    /**
     * Called when a Node is added to Graph
     * @param graph Node is added to Graph Object
     * @param node the node object
     */
    void nodeAdded(TGGraph graph, TGNode node);

    /**
     * Called when an attribute is Added to an entity.
     * @param attribute a new attribute is added
     * @param owner - The owner of the attribute
     */
    void attributeAdded(TGAttribute attribute, TGEntity owner);

    /**
     * Called when an attribute is set.
     * @param attribute - Attribute is changed
     * @param oldValue - the old value
     * @param newValue - the new value
     */
    void attributeChanged(TGAttribute attribute, Object oldValue, Object newValue);

    /**
     * Called when an attribute is removed from the entity.
     * @param attribute - The attribute being removed from the entity
     * @param owner - The owner of the attribute
     */
    void attributeRemoved(TGAttribute attribute, TGEntity owner);

    /**
     * Called when an node is from Graph
     * @param graph the Graph Object
     * @param node the node object
     */
    void nodeRemoved(TGGraph graph, TGNode node);


    /**
     * Called when the entity is deleted.
     * @param entity the entity about to be deleted
     */
    void entityDeleted(TGEntity entity);


}
