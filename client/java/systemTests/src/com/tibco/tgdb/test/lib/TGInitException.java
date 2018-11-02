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
 * Exception thrown when TG server initialization fails.
 * @author sbagi@tibco.com
 *
 */
public class TGInitException extends TGGeneralException {

	private static final long serialVersionUID = 6236267002617589513L;

	private String output = "";
	
	/**
	 * Create an init exception
	 * @param message Reason for not initializing
	 */
	public TGInitException(String message) {
		super(message);
	}
	
	/**
	 * Create an init exception
	 * @param message Reason for not initializing
	 * @param output  Init output console for more info
	 */
	public TGInitException(String message, String output) {
		super(message);
		this.setOutput(output);
		
	}

	/**
	 * @return Init output console
	 */
	public String getOutput() {
		return this.output;
	}

	/**
	 * @param output Init output console
	 */
	public void setOutput(String output) {
		this.output = output;
	}

}
