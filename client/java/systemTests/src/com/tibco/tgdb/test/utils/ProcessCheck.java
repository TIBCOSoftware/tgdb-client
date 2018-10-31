package com.tibco.tgdb.test.utils;

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

import java.io.IOException;

import org.apache.commons.exec.CommandLine;
import org.apache.commons.exec.DefaultExecutor;
import org.apache.commons.exec.OS;
import org.apache.commons.exec.PumpStreamHandler;

/**
 * Utility class to check on processes
 * 
 * @author sbagi@tibco.com
 */
public class ProcessCheck {
	
	private ProcessCheck() {
		;
	}
	
	/**
	 * Check whether a process is running or not
	 * @param pid pid of the process to monitor
	 * @return true if process is running
	 * @throws IOException IO exception
	 */
	public static boolean isProcessRunning(int pid) throws IOException {
	    String line;
	    if (OS.isFamilyWindows()) {
	        //tasklist exit code is always 0. Parse output
	        //findstr exit code 0 if found pid, 1 if it doesn't
	        line = "cmd /c \"tasklist /FI \"PID eq " + pid + "\" | findstr " + pid + "\"";
	    }
	    else {
	        //ps exit code 0 if process exists, 1 if it doesn't
	        line = "ps -p " + pid;
	    }
	    CommandLine cmdLine = CommandLine.parse(line);
	    DefaultExecutor executor = new DefaultExecutor();
	    executor.setStreamHandler(new PumpStreamHandler(null, null, null));
	    executor.setExitValues(new int[]{0, 1});
	    int exitValue = executor.execute(cmdLine);
	    if (exitValue == 0)
	    	return true;
	    else if (exitValue == 1)
	    	return false;
	    else // should never get to here in theory since execute would throw exception
	    	return true; 
	}
}
