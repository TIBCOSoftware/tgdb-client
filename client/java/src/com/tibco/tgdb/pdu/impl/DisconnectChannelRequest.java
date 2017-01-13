package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.pdu.VerbId;

public class DisconnectChannelRequest extends AuthenticatedMessage {

	@Override
	public boolean isUpdateable() {
		return false;
	}

	@Override
	public VerbId getVerbId() {
		return VerbId.DisconnectChannelRequest;
	}

}
