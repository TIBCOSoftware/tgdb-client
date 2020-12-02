package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.model.TGAttributeDescriptor;

import java.util.Collection;
import java.util.Iterator;

public class MutableNodeTypeImpl extends NodeTypeImpl {

	public MutableNodeTypeImpl (
		String _nodeTypeName,
		Collection<TGAttributeDescriptor> _listOfAttributes,
		Collection<TGAttributeDescriptor> _listOfPrimaryKeys
	) 
	{
		super(_nodeTypeName, null);
		
		attributes.clear();
		
		if (_listOfAttributes != null)
		{
			Iterator<TGAttributeDescriptor> itAttribs = _listOfAttributes.iterator();
			for (;itAttribs.hasNext();)
			{
				TGAttributeDescriptor attrib = itAttribs.next();
				String key = attrib.getName();
				attributes.put(key, attrib);
			}
		}
		
		pKeys.clear();
		if (_listOfPrimaryKeys != null)
		{
			pKeys.addAll(_listOfPrimaryKeys);
		}
	}

}