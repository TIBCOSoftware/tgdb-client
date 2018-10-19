package com.tibco.tgdb.restful.jaxrs.jersey;

import javax.inject.Singleton;

import org.glassfish.hk2.utilities.binding.AbstractBinder;

import com.tibco.tgdb.restful.TGDBServiceFactory;
import com.tibco.tgdb.restful.TGDBServiceFactoryImpl;

public class TGDBApplicationBinder  extends AbstractBinder {
	@Override
	protected void configure() {
		bind(TGDBServiceFactoryImpl.class)
		.to(TGDBServiceFactory.class)
        .in(Singleton.class);
	}
}