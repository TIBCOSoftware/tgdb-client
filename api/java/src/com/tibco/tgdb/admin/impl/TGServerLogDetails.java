/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 *  File name : TGServerLogDetails.java
 *  Created on: 04/22/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: TGServerLogDetails.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.log.TGLogComponent;
import com.tibco.tgdb.log.TGLogLevel;


public class TGServerLogDetails {

	
	protected TGLogLevel logLevel;
	protected TGLogComponent logComponent;
	
	public TGServerLogDetails (TGLogComponent _logComponent, TGLogLevel _logLevel)
	{
		logComponent = _logComponent;
		logLevel = _logLevel;
	}
	
	public TGLogLevel getLogLevel() {
		return logLevel;
	}

	public TGLogComponent getLogComponent() {
		return logComponent;
	}	
	
}
