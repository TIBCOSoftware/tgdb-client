
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

 * File name :TGLogLevel.java
 * Created on: 04/24/2019 
 * Created by: nimish
 * SVN Id: $Id: TGLogLevel.java 3127 2019-04-25 22:56:22Z nimish $
 */

package com.tibco.tgdb.log;

public enum TGLogLevel {

	TGLL_Console(-2),
	TGLL_Invalid(-1),
	TGLL_Fatal(0),
	TGLL_Error(1),
	TGLL_Warn(2),
	TGLL_Info(3),
	TGLL_User(4),
	TGLL_Debug(5),
	TGLL_DebugFine(6),
	TGLL_DebugFiner(7),
	TGLL_MaxLogLevel(8);
	
	protected int ll;

	TGLogLevel(int _ll)
	{
		ll = _ll;
	}
	
	public int getLogLevel () 
	{
		return ll;
	}
}