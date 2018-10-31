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
 * Exception thrown for general TG-related failures.
 * @author sbagi@tibco.com
 *
 */
public class TGGeneralException extends Exception {

	private static final long serialVersionUID = -6564871081756662355L;

	/**
	 * Create TG general exception
	 * @param message What is the exception about
	 */
	public TGGeneralException(String message) {
		super(message);
	}
}
