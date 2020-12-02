package com.tibco.tgdb.admin.impl;

public class StatusForCreateEntityType {
	private int resultId;
	private String errorMessage;
	
	public StatusForCreateEntityType() {
	}
	
	public StatusForCreateEntityType(int resultId, String errorMessage) {
		this.resultId = resultId;
		this.errorMessage = errorMessage;
	}

	public int getResultId() {
		return resultId;
	}

	public void setResultId(int resultId) {
		this.resultId = resultId;
	}

	public String getErrorMessage() {
		return errorMessage;
	}

	public void setErrorMessage(String errorMessage) {
		this.errorMessage = errorMessage;
	}
	
	
}
