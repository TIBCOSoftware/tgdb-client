/**
 * 
 */
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

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import bsh.EvalError;
import bsh.Interpreter;

/**
 * Utility class to convert piped data to TestNG data provider structure
 * <pre>{@code
 * Data structure looks like this:
 * [param1|param2|param3] <-- header
 * value11|value12|value13 <-- data
 * value21|value22|value23 <-- data
 * #value31|value32|value33 <-- commented out data
 * value41|{{return Math.random();}}|value43 <-- scripted data
 * ...						
 * }</pre>
 * @author sbagi@tibco.com
 */
public class PipedData {

	private PipedData() {
		;
	}
	
	/**
	 * Read a pipe-separated data file and return a TestNG-friendly data set.
	 * @param dataFile data file
	 * @return two-dimensional array - lines and parameter values
	 * @throws IOException Something wrong with reading the data file
	 * @throws EvalError Scripted data has syntax problem
	 */
	public static Object[][] read(File dataFile) throws IOException, EvalError {
		
		BufferedReader reader = new BufferedReader(new FileReader(dataFile));
		return read(reader);
	}

	/**
	 * Read a pipe-separated data stream and return a TestNG-friendly data set.
	 * @param dataStream data stream
	 * @return two-dimensional array - lines and parameter values
	 * @throws IOException Something wrong with reading the data file
	 * @throws EvalError Scripted data has syntax problem
	 */
	public static Object[][] read(InputStream dataStream) throws IOException, EvalError {
		BufferedReader reader = new BufferedReader(new InputStreamReader(dataStream));
		return read(reader);
	}
	
	/**
	 * Read a pipe-separated data stream and return a TestNG-friendly data set.
	 * @param dataStream data stream
	 * @param charset character set
	 * @return two-dimensional array - lines and parameter values
	 * @throws IOException Something wrong with reading the data file
	 * @throws EvalError Scripted data has syntax problem
	 */
	public static Object[][] read(InputStream dataStream, Charset charset) throws IOException, EvalError {
		BufferedReader reader = new BufferedReader(new InputStreamReader(dataStream, charset));
		return read(reader);
	}
	
	private static Object[][] read(BufferedReader dataReader) throws IOException, EvalError {
		
		Interpreter bsh = new Interpreter();
		
		String line = dataReader.readLine();
		List<Object[]> data = new ArrayList<Object[]>();
		boolean foundHeader = false;
		int nbHeaderParameters = 0;
		int nbLines = 0;
		while (line != null) {
			if (!line.startsWith("#")) {
				if (foundHeader) {
					String[] tmp = line.split("\\|");
					Object[] tmp2 = new Object[tmp.length];
					for (int i=0; i<tmp.length; i++) {
						if (tmp[i].matches("\\{\\{.*\\}\\}")) {
							tmp2[i] = bsh.eval(tmp[i].substring(2,(tmp[i].length()-2)));
						}
						else
							tmp2[i] = tmp[i];
					}
					data.add(Arrays.copyOf(tmp2,nbHeaderParameters));
					nbLines ++;
				 }
				 else {
					 if (line.matches("\\[.*\\]")) { 
						 foundHeader = true;
						 nbHeaderParameters = line.split("\\|").length;
					 }
				 }
			}
			line = dataReader.readLine();
		}
		dataReader.close();
		return (Object[][]) data.toArray(new Object[nbLines][nbHeaderParameters]);
	}
	
}
