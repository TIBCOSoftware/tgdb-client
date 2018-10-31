package com.tibco.tgdb.test.lib;

/**
 * Copyright 2018 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except 
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

/**
 * Exception thrown when TG server start-up fails.
 * @author sbagi@tibco.com
 *
 */
public class TGStartException extends TGGeneralException {

	private static final long serialVersionUID = 6387384148422834860L;
	private String output = "";
	private String error = "";

	/**
	 * Create a start exception
	 * @param message Reason for not starting up
	 */
	public TGStartException(String message) {
		super(message);
	}
	
	/**
	 * Create a start exception
	 * @param message Reason for not starting up
	 * @param output  Start output console for more info
	 */
	public TGStartException(String message, String output) {
		super(message);
		this.setOutput(output);
		
	}
	
	/**
	 * Create a start exception
	 * @param message Reason for not starting up
	 * @param output  Start output console for more info
	 * @param error  Start error console for more info
	 */
	public TGStartException(String message, String output, String error) {
		super(message);
		this.setOutput(output);
		this.setError(error);
	}
	
	/**
	 * @return Start output console
	 */
	public String getOutput() {
		return this.output;
	}

	/**
	 * @param output Start output console
	 */
	public void setOutput(String output) {
		this.output = output;
	}
	
	/**
	 * @return Start error console
	 */
	public String getError() {
		return this.error;
	}

	/**
	 * @param error Start error console
	 */
	public void setError(String error) {
		this.error = error;
	}

}
