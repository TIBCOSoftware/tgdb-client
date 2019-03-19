package types

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGChangeListener.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ChangeListener is an event listener that gets triggered when that event occurs
type TGChangeListener interface {
	// AttributeAdded gets called when an attribute is Added to an entity.
	AttributeAdded(attr TGAttribute, owner TGEntity)
	// AttributeChanged gets called when an attribute is set.
	AttributeChanged(attr TGAttribute, oldValue, newValue interface{})
	// AttributeRemoved gets called when an attribute is removed from the entity.
	AttributeRemoved(attr TGAttribute, owner TGEntity)
	// EntityCreated gets called when an entity is Added
	EntityCreated(entity TGEntity)
	// EntityDeleted gets called when the entity is deleted
	EntityDeleted(entity TGEntity)
	// NodeAdded gets called when a node is Added
	NodeAdded(graph TGGraph, node TGNode)
	// NodeRemoved gets called when a node is removed
	NodeRemoved(graph TGGraph, node TGNode)
}
