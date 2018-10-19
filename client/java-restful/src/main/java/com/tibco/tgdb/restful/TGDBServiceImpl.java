package com.tibco.tgdb.restful;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.logging.Logger;

import javax.servlet.ServletContext;
import javax.ws.rs.core.Context;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEdgeType;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGEntityType;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.query.TGResultSet;

public class TGDBServiceImpl implements TGDBService{
	
	public static final Logger LOGGER = Logger.getLogger(TGDBServiceImpl.class.getName());

	private TGConnection connection = null;
	private TGGraphObjectFactory gof = null;
	private TGGraphMetadata metadata = null;
	
	private HashMap<TGEntityType, List<String>> keyMap = new HashMap<TGEntityType, List<String>>();
	
	private String url = null;
	private String user = null;
    public String passwd = null;

	@Context
	private ServletContext application;
	
	public TGDBServiceImpl() {
	}

	public void init(String url, String user, String passwd) {
		this.url = url;
		this.user = user;
		this.passwd = passwd;
	}
		
	//---------- TGConnectionExceptionListener API -----------

	@Override
	public void onException(Exception e) {
		LOGGER.info("[TGDBService::onException] Exception Happens : " + e.getMessage());
		if(null!=this.connection) {
			this.connection.disconnect();
		}
		this.connection=null;
	}
	
	//---------- Metadata API -----------
	
	public TGGraphObjectFactory getFactory() throws Exception {
		ensureConnection();	
		return gof;
	}

	public TGGraphMetadata getMetadata() throws Exception {
		ensureConnection();	
		LOGGER.info("[TGDBService::getMetadata] force fetching metadata!!!");
		fetchMetadata();
		return this.metadata;
	}

	public Collection<TGNodeType> getNodeTypes() throws Exception {
		ensureConnection();	
		return metadata.getNodeTypes();
	}
	
	public TGEdgeType getEdgeType(String type) throws Exception {
		ensureConnection();	
		return metadata.getEdgeType(type);
	}
	
	public Collection<TGEdgeType> getEdgeTypes() throws Exception {
		ensureConnection();	
		return metadata.getEdgeTypes();
	}
	
	public List<String> getNodeKeyfields(String type) throws Exception {
		ensureConnection();	
        TGNodeType nodeType = metadata.getNodeType(type);
        if (null==nodeType) {
            throw new Exception("Node type not found");
        }
		List<String> key = new ArrayList<String>();
		for(TGAttributeDescriptor attrDesc : nodeType.getPKeyAttributeDescriptors()) {
			key.add(attrDesc.getName());
		}
		return key;
	}
	
	public List<String> getNodeKeyfields(TGEntityType type) throws Exception {
		return keyMap.get(type);
	}
	
	//---------- Node API -----------

	public TGNode insertNode(String type, Map<String, Object> data) throws Exception {
		
		ensureConnection();
		TGNode node = buildNode(type, data);
        connection.insertEntity(node);
        connection.commit();

		return node;
	}
	
	public TGNode updateNode(String type, Map<String, Object> data) throws Exception {
		TGNode node = buildNode(type, data);
		ensureConnection();		
        connection.updateEntity(node);
        connection.commit();

		return node;
	}

	private Object getAttributeValue(TGAttributeDescriptor keyElement, Map<String, List<String>> parameter) {
		TGAttributeType type = keyElement.getType();
		String strValue = parameter.get(keyElement.getName()).get(0);
		Object value = null;
		if(TGAttributeType.Int.equals(type)) {
			value = Integer.parseInt(strValue);
		} else if(TGAttributeType.Long.equals(type)) {
			value = Long.parseLong(strValue);
		} else if(TGAttributeType.Double.equals(type)) {
			value = Double.parseDouble(strValue);
		} else if(TGAttributeType.Float.equals(type)) {
			value = Float.parseFloat(strValue);
		} else if(TGAttributeType.Boolean.equals(type)) {
			value = Boolean.parseBoolean(strValue);
		} else {
			value = strValue;
		}
		return value;
	}
	
	
	public TGEntity getNode(String type, Map<String, List<String>> parameter) throws Exception {		
		ensureConnection();	
        TGNodeType nodeType = metadata.getNodeType(type);
        if (null==nodeType) {
            throw new Exception("Node type not found");
        }

        TGKey tgKey = gof.createCompositeKey(type);
        for(TGAttributeDescriptor attrDesc : nodeType.getPKeyAttributeDescriptors()) {
        	String name = attrDesc.getName();
        	tgKey.setAttribute(name, getAttributeValue(attrDesc, parameter));
    		LOGGER.info("[TGDBService::getNode] Key attribute set for " + name);
        }
        
        TGEntity entity = connection.getEntity(tgKey, getRequestParameterToOption(parameter));

		return entity;
	}
	
	public TGEntity deleteNode(String type, Map<String, List<String>> key) throws Exception {	
		
        TGEntity entity = getNode(type, key);
		ensureConnection();		
        connection.deleteEntity(entity);
        connection.commit();

		return entity;
	}
	
	//---------- Edge API -----------
	
	public TGEdge insertEdge(String type, Map<String, Object> data) throws Exception {
		
		TGEdge edge = buildEdge(type, data);
		ensureConnection();
        connection.insertEntity(edge);
		try {
	        connection.commit();
		} catch(Exception e) {
			e.printStackTrace();
		}

		return edge;
	}
	
	@SuppressWarnings("unchecked")
	public Collection<TGEdge> deleteEdge(String type, Map<String, Object> data) throws Exception {
		
        TGEdgeType edgeType = metadata.getEdgeType(type);
        if (null==edgeType) {
            throw new Exception("Edge type not found");
        }

		Map<String, Object> source = (Map<String, Object>)data.get("source");
		String sourceType = (String)source.get("type");
		Map<String, String> sourceKey = (Map<String, String>)source.get("key");
		TGNode sourceNode = fetchNode(sourceType, sourceKey);
		
		Map<String, Object> target = (Map<String, Object>)data.get("target");
		String targetType = (String)source.get("type");
		Map<String, String> targetKey = (Map<String, String>)target.get("key");
		TGNode targetNode = fetchNode(targetType, targetKey);
		
		//??????????????????
		Collection<TGEdge> edges = sourceNode.getEdges(edgeType, getEdgeDirection((int)data.get("direction")));
		for(TGEdge edge : edges) {
			for(TGNode node : edge.getVertices()) {
				if(node.equals(targetNode)) {
			        connection.deleteEntity(edge);
				}
			}
		}
		
		ensureConnection();
        connection.commit();

		return edges;
	}	
	//---------- Search API -----------

	@SuppressWarnings("unchecked")
	public TGResultSet query(Map<String, Object> para) throws TGException {				
		String queryString = null; 
		String edgeFilter = null;
		String traversalCondition = null;
		String endCondition = null;
		
		Map<String, Object> queryObj = null;
		if(null!=(queryObj=(Map<String, Object>)para.get("query"))) {
			queryString = (String)queryObj.get("queryString");
			edgeFilter = (String)queryObj.get("edgeFilter");
			traversalCondition = (String)queryObj.get("traversalCondition");
			endCondition = (String)queryObj.get("endCondition");
		}
		
		LOGGER.info("[TGDBService::query] queryString = " +  queryString);
		LOGGER.info("[TGDBService::query] edgeFilter = " +  edgeFilter);
		LOGGER.info("[TGDBService::query] traversalCondition = " +  traversalCondition);
		LOGGER.info("[TGDBService::query] endCondition = " +  endCondition);

		ensureConnection();
        TGResultSet resultSet = connection.executeQuery(
        		queryString, edgeFilter, traversalCondition, endCondition, buildQueryOption(para));
        
        return resultSet;
	}
	
	//---------- Private -----------
	
	private TGNode fetchNode(String type, Map<String, String> key) throws Exception {
		ensureConnection();	
        TGNodeType nodeType = metadata.getNodeType(type);
        if (null==nodeType) {
            throw new Exception("Node type not found");
        }

        TGKey tgKey = gof.createCompositeKey(type);
        
        for(TGAttributeDescriptor attrDesc : nodeType.getPKeyAttributeDescriptors()) {
        	String name = attrDesc.getName();
        	tgKey.setAttribute(name, key.get(name));
        	LOGGER.info("[TGDBService::createNode] Key attribute set for " + name);
        }
        
        TGQueryOption option = TGQueryOption.createQueryOption();
        return (TGNode)connection.getEntity(tgKey, option);
	}
	
	@SuppressWarnings("unchecked")
	private TGNode buildNode(String type, Map<String, Object> data) throws Exception {
		Map<String, Object> properties = (Map<String, Object>)data.get("properties");
		
        TGNodeType nodeType = metadata.getNodeType(type);
        if (null==nodeType) {
            throw new Exception("Node type not found");
        }

        for(TGAttributeDescriptor attrDesc : nodeType.getPKeyAttributeDescriptors()) {
        	String name = attrDesc.getName();
        	LOGGER.info("[TGDBService::createNode] Key attribute -> " + name);
        }
        
        TGNode node = gof.createNode(nodeType);
        for(TGAttributeDescriptor attrDesc : nodeType.getAttributeDescriptors()) {
        	String name = attrDesc.getName();
            node.setAttribute(name, properties.get(name));
            LOGGER.info("[TGDBService::createNode] attribute set -> " + name + " = " + node.getAttribute(name).getValue());
        }

		return node;
	}
	
	@SuppressWarnings("unchecked")
	private TGEdge buildEdge(String type, Map<String, Object> data) throws Exception {
		ensureConnection();	
        TGEdgeType edgeType = metadata.getEdgeType(type);
        if (null==edgeType) {
            throw new Exception("Edge type not found");
        }
		Map<String, Object> source = (Map<String, Object>)data.get("source");
		String sourceType = (String)source.get("type");
		Map<String, String> sourceKey = (Map<String, String>)source.get("key");
		TGNode sourceNode = fetchNode(sourceType, sourceKey);
		
		Map<String, Object> target = (Map<String, Object>)data.get("target");
		String targetType = (String)source.get("type");
		Map<String, String> targetKey = (Map<String, String>)target.get("key");
		TGNode targetNode = fetchNode(targetType, targetKey);
       
		Map<String, Object> properties = (Map<String, Object>)data.get("properties");
        TGEdge edge = gof.createEdge(sourceNode, targetNode, edgeType);
        //TGEdge edge = gof.createEdge(sourceNode, targetNode, TGEdge.DirectionType.BiDirectional);
        
        for(TGAttributeDescriptor attrDesc : edgeType.getAttributeDescriptors()) {
        	String name = attrDesc.getName();
        	if(null!=properties.get(name)) {
            	edge.setAttribute(name, properties.get(name));
				LOGGER.info("[TGDBService::buildEdge] attribute set -> " + name + " = " + edge.getAttribute(name).getValue());
        		//System.out.println("[TGDBService::buildEdge] attribute set -> " + name + " = " + edge.getAttribute(name).getValue());
        	}
        }

		return edge;
	}

	private void ensureConnection() throws TGException {
		if(null==connection||null==(gof=connection.getGraphObjectFactory())) {
			if(null!=this.connection) {
				LOGGER.info("[TGDBService::ensureConnection] Connection not null try to disconnect ..........");
		    	//System.out.println("[TGDBService::ensureConnection] Connection not null try to disconnect ..........");
				connection.disconnect();
			}
			LOGGER.info("[TGDBService::ensureConnection] Will try to connect ..........");
	    	//System.out.println("[TGDBService::ensureConnection] Will try to connect ..........");
	    	try {
	    		this.connection = TGConnectionFactory.getInstance().createConnection(url, user, passwd, null);
	    		this.connection.connect();
	    		fetchMetadata();
	    	} catch(TGException tge) {
	    		LOGGER.info("[TGDBService::ensureConnection] Unable to make connection !!! Will not connect ......");
	    		this.connection = null;
	    		throw tge;
	    	}
			gof = this.connection.getGraphObjectFactory();
			if (gof == null) {
				this.connection=null;
				throw new TGException("[TGDBService::ensureConnection] Graph Object Factory is null...exiting");
			}
			this.connection.setExceptionListener(this);
		} else {
			LOGGER.info("[TGDBService::ensureConnection] Connected !!! Will not reconnect ..........");
		}
	}
	
	private TGEdge.Direction getEdgeDirection(int edgeDirection) {
		switch(edgeDirection) {
		case 0:
			return TGEdge.Direction.Inbound;
		case 1:
			return TGEdge.Direction.Outbound;
		case 2:
			return TGEdge.Direction.Any;
		}
		return null;
	}
	
	@SuppressWarnings("unused")
	private TGEdge.DirectionType getEdgeDirectionType(int edgeDirectionType) {
		switch(edgeDirectionType) {
		case 0:
			return TGEdge.DirectionType.UnDirected;
		case 1:
			return TGEdge.DirectionType.Directed;
		case 2:
			return TGEdge.DirectionType.BiDirectional;
		}
		return null;
	}

	private void fetchMetadata() throws TGException {
		this.metadata = this.connection.getGraphMetadata(true);
		keyMap.clear();
		for(TGNodeType nodeType : this.metadata.getNodeTypes()) {
        	List<String> key = new ArrayList<String>();
	        for(TGAttributeDescriptor attrDesc : nodeType.getPKeyAttributeDescriptors()) {
	        	key.add(attrDesc.getName());
	        }
	        keyMap.put(nodeType, key);
		}		
	}
		
	private TGQueryOption getRequestParameterToOption(Map<String, List<String>> getRequestParameter) {
		Map<String, Object> para = new HashMap<String, Object>();
		if(null!=getRequestParameter.get("prefetchSize")) {
			para.put("prefetchSize", Integer.parseInt(getRequestParameter.get("prefetchSize").get(0)));
		}
		
		if(null!=getRequestParameter.get("edgeLimit")) {
			para.put("edgeLimit", Integer.parseInt(getRequestParameter.get("edgeLimit").get(0)));
		}
		
		if(null!=getRequestParameter.get("traversalDepth")) {
			para.put("traversalDepth", Integer.parseInt(getRequestParameter.get("traversalDepth").get(0)));
		}

		return buildQueryOption(para);
	}
	
	private TGQueryOption buildQueryOption(Map<String, Object> para) {
		TGQueryOption option = TGQueryOption.createQueryOption();
		
		option.setPrefetchSize(
				null==para.get("prefetchSize")?500:(int)para.get("prefetchSize"));
		option.setTraversalDepth(
				null==para.get("traversalDepth")?5:(int)para.get("traversalDepth"));
		option.setEdgeLimit(
				null==para.get("edgeLimit")?100:(int)para.get("edgeLimit"));
		
		LOGGER.info("[TGDBService::buildQueryOption] prefetchSize = " +  option.getPrefetchSize());
		LOGGER.info("[TGDBService::buildQueryOption] traversalDepth = " +  option.getTraversalDepth());
		LOGGER.info("[TGDBService::buildQueryOption] edgeLimit = " +  option.getEdgeLimit());

		return option;
	}
		
	@SuppressWarnings("unused")
	private static void traverse(
			Map<String, List<TGEntity>> result, 
			TGEntity entity, 
			int maxDepth, 
			int currDepth) {
		if(maxDepth<0||(maxDepth>=currDepth)) {
//			System.out.println("[TGDBRestService::traverse] maxDepth = " + maxDepth + ", currDepth = " + currDepth + ", result = " + result);
			if(entity instanceof TGNode) {
				TGNode node = (TGNode)entity;
				result.get("nodes").add(node);			
				for(TGEdge edge: node.getEdges()) {
					if(!result.get("edges").contains(edge)) {
						traverse(result, edge, maxDepth, ++currDepth);
					}
				}
				
			} else if(entity instanceof TGEdge) {
				TGEdge edge = (TGEdge)entity;
				result.get("edges").add(edge);			
				for(TGNode node: edge.getVertices()) {
					if(!result.get("nodes").contains(node)) {
						traverse(result, node, maxDepth, currDepth);
					}
				}
			}
		}
	}
}
