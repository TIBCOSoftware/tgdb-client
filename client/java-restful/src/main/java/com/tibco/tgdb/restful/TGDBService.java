package com.tibco.tgdb.restful;

import java.util.Collection;
import java.util.List;
import java.util.Map;

import com.tibco.tgdb.connection.TGConnectionExceptionListener;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGEntityType;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.query.TGResultSet;

public interface TGDBService extends TGConnectionExceptionListener{
	
	public void init(String string, String string2, String string3);

	public Object getMetadata() throws Exception;

	public Object getNodeTypes() throws Exception;

	public Object getEdgeTypes() throws Exception;

	public TGEntity getNode(String type, Map<String, List<String>> parameters) throws Exception;
	
	public TGNode insertNode(String type, Map<String, Object> data) throws Exception;

	public TGNode updateNode(String type, Map<String, Object> data) throws Exception;

	public TGEntity deleteNode(String type, Map<String, List<String>> parameters) throws Exception;

	public TGEntity insertEdge(String type, Map<String, Object> data) throws Exception;

	public Collection<TGEdge> deleteEdge(String type, Map<String, Object> data) throws Exception;

	public TGResultSet query(Map<String, Object> data) throws TGException;

	public List<String> getNodeKeyfields(TGEntityType type) throws Exception;

}
