package com.tibco.tgdb.restful;

import java.util.HashMap;
import java.util.logging.Logger;

import javax.ws.rs.core.Context;
import javax.ws.rs.core.UriInfo;

import org.apache.commons.io.IOUtils;
import org.json.JSONObject;

public class TGDBServiceFactoryImpl extends TGDBServiceFactory{
	
	public static final Logger LOGGER = Logger.getLogger(TGDBServiceFactoryImpl.class.getName());

	private HashMap<String, TGDBService> services = new HashMap<String, TGDBService>();
	
    @Context
    UriInfo uriInfo;
    
    public TGDBServiceFactoryImpl() {
		try {
			LOGGER.info("[TGDBServiceFactory::TGDBServiceFactory] Creating ......... ");
			ClassLoader classLoader = getClass().getClassLoader();
		    String content = IOUtils.toString(classLoader.getResourceAsStream("tgdb.json"));
			JSONObject dbs = new JSONObject(content);
			for(Object key : dbs.names()) {
				JSONObject db = dbs.getJSONObject((String) key);
				TGDBService service = new TGDBServiceImpl();
				service.init(db.getString("url"), db.getString("user"), db.getString("password"));
				services.put((String) key, service);
			}
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
    }
    
    public TGDBService getService() {
        return services.get("default");
    }
    
    public void destroy(TGDBService service) {
        /* noop */
    }
	    
}