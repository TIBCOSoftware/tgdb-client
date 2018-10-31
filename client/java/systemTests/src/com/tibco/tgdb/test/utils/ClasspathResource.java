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

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;

/**
 * Utility class to handle resources from classpath.
 * 
 * @author sbagi@tibco.com
 */
public class ClasspathResource {
	
	private ClasspathResource(){
		;
	}
	
	/**
	 * Make a copy of a resource file from the classpath to the file system.
	 * @param resourcePath Location of the resource inside the classpath
	 * @param targetPath Location of the file on the file system
	 * @return the target file
	 * @throws IOException I/O problems
	 */
	public static File getResourceAsFile(String resourcePath, String targetPath) throws IOException {
		
		 InputStream is = ClassLoader.getSystemClassLoader().getResourceAsStream(resourcePath);
		 if (is == null) {
			 throw new IOException("Resource " + resourcePath + " not found in classpath");
		 }

	     File targetFile = new File(targetPath);
	     if (!targetFile.exists()) {
	    	 targetFile.getParentFile().mkdirs();
	    	 if (!targetFile.createNewFile())
	    		 throw new IOException("Target file " + targetPath + " cannot be created");
	     }
	     FileOutputStream out = new FileOutputStream(targetFile);
	     byte[] buffer = new byte[1024];
	     int bytesRead;
	     while ((bytesRead = is.read(buffer)) != -1) {
	    	 out.write(buffer, 0, bytesRead);
	     }
	        
	     if (out != null)
	    	 out.close();
	     is.close();
	     
	     return targetFile;		
	}
	
}
