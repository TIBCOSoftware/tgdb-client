package com.tibco.tgdb.restful;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.logging.Logger;

import javax.inject.Inject;
import javax.ws.rs.Consumes;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.UriInfo;

import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.restful.jaxrs.jersey.TGDBApplicationConfig;
import com.tibco.tge.adapter.api.restful.common.model.APIError;
import com.tibco.tge.adapter.api.restful.common.model.APIResponse;

@Path( "/tgdb" ) 
public class TGDBRestService {
	
	public static final Logger LOGGER = Logger.getLogger(TGDBApplicationConfig.class.getName());

	private TGDBService tgdb;// = TGDBService.instance();
		
	private TGDBServiceFactory serviceFactory;
	
	@Inject 
	public TGDBRestService(TGDBServiceFactory serviceFactory) {
		this.tgdb = serviceFactory.getService();
	}
	
	//---------- Metadata API -----------
	
	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @GET
    @Path("/metadata")
    public Response getMetadata() {		
    	APIResponse response = new APIResponse();
    	
        try {
        	response.data = tgdb.getMetadata();
		} catch (Exception e) {
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();
    }
	
	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @GET
    @Path("/nodetypes")
    public Response getNodeTypes() {		
    	APIResponse response = new APIResponse();
    	
        try {
        	response.data = tgdb.getNodeTypes();
		} catch (Exception e) {
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();
    }

	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @GET
    @Path("/edgetypes")
    public Response getEdgeTypes() {		
    	APIResponse response = new APIResponse();
    	
        try {
        	response.data = tgdb.getEdgeTypes();
		} catch (Exception e) {
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();
    }
	
	//---------- Node API -----------
	
	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @GET
    @Path("/node/{type}")
    public Response getNode(
    		@Context final UriInfo uriInfo,
    		@PathParam("type") final String type) {		
    	APIResponse response = new APIResponse();
		System.out.println("[TGDBRestService::traverse] getNode ...............");

        try {
        	Map<String, List<String>> parameters = uriInfo.getQueryParameters();
        	TGEntity entity = tgdb.getNode(type, parameters);
    		Map<String, List<TGEntity>> tgResult = new HashMap<String, List<TGEntity>>();
    		tgResult.put("nodes", new ArrayList<TGEntity>());
    		tgResult.put("edges", new ArrayList<TGEntity>());
    		
    		int maxDepth = -1; 
    		if(null!=parameters.get("traversalDepth")) {
    			maxDepth = Integer.parseInt(parameters.get("traversalDepth").get(0));
    		}

    		traverse(tgResult, entity, maxDepth, 0);

        	response.data = buildResult(tgResult);
        	response.success = true;
		} catch (Exception e) {
			e.printStackTrace();
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();
    }
	
	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @PUT
    @Path("/node/{type}")
    public Response insertNode(
    		@PathParam("type") final String type,
            final Map<String, Object> data) {		
    	APIResponse response = new APIResponse();
		System.out.println("[TGDBRestService::traverse] insertNode ...............");
    	
        try {
        	TGNode node = tgdb.insertNode(type, data);
    		Map<String, List<TGEntity>> tgResult = new HashMap<String, List<TGEntity>>();
    		tgResult.put("nodes", new ArrayList<TGEntity>());
    		tgResult.put("edges", new ArrayList<TGEntity>());
        	traverse(tgResult, node, 0, 0);
        	
        	response.data = buildResult(tgResult);
        	response.success = true;
		} catch (Exception e) {
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();

    }
	
	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @POST
    @Path("/node/{type}")
    public Response updateNode(
    		@PathParam("type") final String type,
            final Map<String, Object> data) {
    	APIResponse response = new APIResponse();
		System.out.println("[TGDBRestService::traverse] updateNode ...............");
		
        try {
        	TGNode node = tgdb.updateNode(type, data);
    		Map<String, List<TGEntity>> tgResult = new HashMap<String, List<TGEntity>>();
    		tgResult.put("nodes", new ArrayList<TGEntity>());
    		tgResult.put("edges", new ArrayList<TGEntity>());
        	traverse(tgResult, node, 0, 0);
        	
        	response.data = buildResult(tgResult);
        	response.success = true;
		} catch (Exception e) {
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();

    }
	
	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @DELETE
    @Path("/node/{type}")
    public Response deleteNode(
    		@Context final UriInfo uriInfo,
    		@PathParam("type") final String type) {
    	APIResponse response = new APIResponse();
    	
        try {
           	Map<String, List<String>> parameters = uriInfo.getQueryParameters();
        	TGEntity entity = tgdb.deleteNode(type, parameters);
    		Map<String, List<TGEntity>> tgResult = new HashMap<String, List<TGEntity>>();
    		tgResult.put("nodes", new ArrayList<TGEntity>());
    		tgResult.put("edges", new ArrayList<TGEntity>());
    		traverse(tgResult, entity, 0, 0);

        	response.data = buildResult(tgResult);
        	response.success = true;
		} catch (Exception e) {
			e.printStackTrace();
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();
	}
	
	//---------- Edge API -----------

	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @PUT
    @Path("/edge/{type}")
    public Response insertEdge(
    		@PathParam("type") final String type,
            final Map<String, Object> data) {		
    	APIResponse response = new APIResponse();
    	
        try {        	
        	TGEntity entity = tgdb.insertEdge(type, data);
    		Map<String, List<TGEntity>> tgResult = new HashMap<String, List<TGEntity>>();
    		tgResult.put("nodes", new ArrayList<TGEntity>());
    		tgResult.put("edges", new ArrayList<TGEntity>());
    		traverse(tgResult, entity, 0, 0);

        	response.data = buildResult(tgResult);
        	response.success = true;        	
		} catch (Exception e) {
			e.printStackTrace();
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();
    }	

	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @DELETE
    @Path("/edge/{type}")
    public Response deleteEdge(
    		@PathParam("type") final String type,
            final Map<String, Object> data) {		
    	APIResponse response = new APIResponse();
    	
        try {        	
        	Collection<TGEdge> edges = tgdb.deleteEdge(type, data);
    		Map<String, List<TGEntity>> tgResult = new HashMap<String, List<TGEntity>>();
    		tgResult.put("nodes", new ArrayList<TGEntity>());
    		tgResult.put("edges", new ArrayList<TGEntity>());
    		for(TGEdge edge : edges) {
        		traverse(tgResult, edge, 0, 0);
    		}

        	response.data = buildResult(tgResult);
        	response.success = true;        	
		} catch (Exception e) {
			e.printStackTrace();
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();
    }		
	//---------- Search API -----------
	
	@Produces( { MediaType.APPLICATION_JSON } )
    @Consumes( { MediaType.APPLICATION_JSON } )
    @POST
    @Path("/search")
    public Response search(
    		final Map<String, Object> data) {
    	APIResponse response = new APIResponse();
    	
        try {
        	response.success = true;
    		Map<String, List<TGEntity>> tgResult = new HashMap<String, List<TGEntity>>();
    		tgResult.put("nodes", new ArrayList<TGEntity>());
    		tgResult.put("edges", new ArrayList<TGEntity>());
        	TGResultSet resultSet = tgdb.query(data);
        	if(null!=resultSet) {
                while (resultSet.hasNext()) {
                    TGNode node = (TGNode) resultSet.next();
            		traverse(tgResult, node, getTraversalDepth(data), 0);
                }
        	}
        	
        	response.data = buildResult(tgResult);
		} catch (Exception e) {
			e.printStackTrace();
 			APIError error = new APIError();
 			error.code = 1001;
 			error.message = e.getMessage();
 			response.error = error;
			response.success = false;
		}
        
        return Response.ok(response).build();
    }	
	
	//---------- Private -----------
	
	@SuppressWarnings("unchecked")
	private Map<String, Object> buildResult(Map<String, List<TGEntity>> tgResult) throws Exception {
		Map<String, Object> result = new HashMap<String, Object>();
		result.put("nodes", new ArrayList<Map<String, Object>>());
		result.put("edges", new ArrayList<Map<String, Object>>());
		
		for(TGEntity entity : tgResult.get("nodes")) {
			TGNode node = (TGNode)entity;
			if(null!=node&&null!=node.getType()) {
				Map<String, Object> aNode = new HashMap<String, Object>();
				((List<Map<String, Object>>)result.get("nodes")).add(aNode);
				aNode.put("_type", node.getType().getName());
				List<Object> id = getNodeKey(node);
				if(null!=id) {
					aNode.put("id", id);
				}
				for(TGAttribute attr : node.getAttributes()) {
					aNode.put(attr.getAttributeDescriptor().getName(), attr.getValue());
				}
			}
		}
		
		for(TGEntity edge : tgResult.get("edges")) {
			if(null!=edge.getType()) {
				List<Object> fromNodeId = getNodeKey(((TGEdge)edge).getVertices()[0]);
				List<Object> toNode = getNodeKey(((TGEdge)edge).getVertices()[1]);				
				if(null==fromNodeId||null==toNode) {
					continue;
				}
				
				Map<String, Object> anEdge = new HashMap<String, Object>();
				anEdge.put("_type", edge.getType().getName());
				anEdge.put("fromNode", fromNodeId);
				anEdge.put("toNode", toNode);
				for(TGAttribute attr : edge.getAttributes()) {
					anEdge.put(attr.getAttributeDescriptor().getName(), attr.getValue());
				}
				((List<Map<String, Object>>)result.get("edges")).add(anEdge);
			}
		}
		
		return result;
	}
	
	private void traverse(
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
	
	private List<Object> getNodeKey(TGNode node) throws Exception {
		List<String> keyFields = tgdb.getNodeKeyfields(node.getType());
		if(null==keyFields) {
			return null;
		}
		List<Object> key = new ArrayList<Object>();
		for(String keyField : keyFields) {
			key.add(node.getAttribute(keyField).getValue());
		}

		return key;
	}
	
	private int getTraversalDepth(Map<String, Object> para) {
		return null==para.get("traversalDepth")?TGQueryOption.DEFAULT_QUERY_OPTION.getTraversalDepth():(int)para.get("traversalDepth");
	}
}