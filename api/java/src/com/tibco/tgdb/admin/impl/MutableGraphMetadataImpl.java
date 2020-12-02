package com.tibco.tgdb.admin.impl;

import java.util.Collection;
import java.util.Iterator;

import com.tibco.tgdb.admin.TGMutableGraphMetadata;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.model.impl.GraphMetadataImpl;
import com.tibco.tgdb.model.impl.GraphObjectFactoryImpl;
import com.tibco.tgdb.model.impl.MutableEdgeTypeImpl;
import com.tibco.tgdb.model.impl.MutableNodeTypeImpl;

public class MutableGraphMetadataImpl extends GraphMetadataImpl implements TGMutableGraphMetadata {

	public MutableGraphMetadataImpl(GraphObjectFactoryImpl gof) {
		super(gof);
	}

	@Override
	public void createNodeType(String typeName, Collection<TGAttributeDescriptor> attrDescs, Collection<TGAttributeDescriptor> pkeyDescs) throws TGException {
		
		runInitialValidationForParameters (typeName, attrDescs, pkeyDescs);
		
		MutableNodeTypeImpl nodeTypeImpl = new MutableNodeTypeImpl(typeName, attrDescs, pkeyDescs);
		AdminConnectionImpl connection = (AdminConnectionImpl) getConnection();
		AdminResponse adminResponse = (AdminResponse) connection.getCommandWithParameters(TGAdminCommand.CreateNodeType, nodeTypeImpl);
		if (adminResponse.createNodeTypeStatus.getResultId() == 0)
		{
		}
		else
		{
			TGException tgException = new TGException(adminResponse.createNodeTypeStatus.getErrorMessage());
			throw tgException;
		}
	}

	private void runInitialValidationForParameters(String typeName, Collection<TGAttributeDescriptor> attrDescs, Collection<TGAttributeDescriptor> pkeyDescs) throws TGException {
		// TODO-N: Error messages need to be externalized
		if (typeName == null)
		{
			throw new TGException ("NodeType name cannot be null"); 
		}
		
		if (attrDescs == null) {
			throw new TGException ("Attribute Descriptors can not be null");	
		}
		
		if (attrDescs.size() == 0) {
			throw new TGException ("Attribute Descriptors are not specified");	
		}

		if (pkeyDescs == null) return;
		
		Iterator<TGAttributeDescriptor> iteratorForPKeys = pkeyDescs.iterator();
		
		for (;iteratorForPKeys.hasNext();)
		{
			TGAttributeDescriptor currentPKey = iteratorForPKeys.next();
			if (!attrDescs.contains(currentPKey))
			{
				throw new TGException ("Error TGPrimaryKeyInvalid : Primary Key attributes must belong to the attributes of node type");
			}
		}
	}

	@Override
	public void createEdgeType(String edgeTypeName, TGEdge.DirectionType directionality, Collection<TGAttributeDescriptor> attrDescs, TGNodeType sourceNodeType, TGNodeType destinationNodeType) throws TGException
	{
		// TODO-N: Error messages need to be externalized
		if (edgeTypeName == null)
		{
			throw new TGException ("EdgeType name cannot be null"); 
		}
		
		if (sourceNodeType == null)
		{
			throw new TGException ("Source nodetype cannot be null");
		}
		
		if (destinationNodeType == null)
		{
			throw new TGException ("Destination nodetype cannot be null");
		}

		MutableEdgeTypeImpl edgeTypeImpl = new MutableEdgeTypeImpl(edgeTypeName, directionality, null, attrDescs, sourceNodeType, destinationNodeType);
		AdminConnectionImpl connection = (AdminConnectionImpl) getConnection();
		AdminResponse adminResponse = (AdminResponse) connection.getCommandWithParameters(TGAdminCommand.CreateEdgeType, edgeTypeImpl);
		
		if (adminResponse.createNodeTypeStatus.getResultId() == 0)
		{
		}
		else
		{
			TGException tgException = new TGException(adminResponse.createNodeTypeStatus.getErrorMessage());
			throw tgException;
		}
	}
	
	
	@Override
	public void createAttributeType(TGAttributeDescriptor attributeDescriptor) throws TGException {
		// TODO-N: Error messages need to be externalized
		if (attributeDescriptor == null)
		{
			throw new TGException ("AttributeDescriptor cannot be null");
		}
		
		AdminConnectionImpl connection = (AdminConnectionImpl) getConnection();
		AdminResponse adminResponse = (AdminResponse) connection.getCommandWithParameters(TGAdminCommand.CreateAttrDesc, attributeDescriptor);

		if (adminResponse.createNodeTypeStatus.getResultId() == 0)
		{
		}
		else
		{
			TGException tgException = new TGException(adminResponse.createNodeTypeStatus.getErrorMessage());
			throw tgException;
		}
	}
}
