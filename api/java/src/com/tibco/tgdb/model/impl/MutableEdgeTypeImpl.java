package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGEdge.DirectionType;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGNodeType;

import java.util.Collection;
import java.util.Iterator;

public class MutableEdgeTypeImpl extends EdgeTypeImpl {

	public MutableEdgeTypeImpl(String name, DirectionType directionType, TGEdgeType parent, Collection<TGAttributeDescriptor> attrDescs, TGNodeType sourceNodeType, TGNodeType destinationNodeType) {
		super(name, directionType, parent);
		this.fromNodeType = sourceNodeType;
		this.toNodeType = destinationNodeType;
		
		attributes.clear();
		
		if (attrDescs != null)
		{
			Iterator<TGAttributeDescriptor> itAttribs = attrDescs.iterator();
			for (;itAttribs.hasNext();)
			{
				TGAttributeDescriptor attrib = itAttribs.next();
				String key = attrib.getName();
				attributes.put(key, attrib);
			}
		}
	}
	
	
}
