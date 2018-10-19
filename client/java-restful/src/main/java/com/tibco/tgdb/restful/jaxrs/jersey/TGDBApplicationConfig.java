package com.tibco.tgdb.restful.jaxrs.jersey;

import java.util.logging.Logger;

import org.glassfish.jersey.server.ResourceConfig;

import com.tibco.tgdb.restful.TGDBRestService;

//@ApplicationPath("/api")
public class TGDBApplicationConfig extends ResourceConfig {
	
	public static final Logger LOGGER = Logger.getLogger(TGDBApplicationConfig.class.getName());

	public TGDBApplicationConfig() {
		LOGGER.info("[TGDBApplicationConfig::TGDBApplicationConfig] Entering .........");
		//System.out.println("[TGDBApplicationConfig::TGDBApplicationConfig] Entering .........");
        register(TGDBRestService.class);
        register(new TGDBApplicationBinder());
	}
}