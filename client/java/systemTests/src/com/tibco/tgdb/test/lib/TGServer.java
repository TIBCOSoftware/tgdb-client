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

import java.io.BufferedReader;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.io.StringReader;
import java.net.Inet4Address;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.List;

import org.apache.commons.exec.CommandLine;
import org.apache.commons.exec.Executor;
import org.apache.commons.exec.OS;
import org.apache.commons.exec.DefaultExecutor;
import org.apache.commons.exec.DefaultExecuteResultHandler;
import org.apache.commons.exec.PumpStreamHandler;
import org.apache.commons.exec.util.StringUtils;
import org.apache.commons.exec.ExecuteWatchdog;
import org.apache.commons.exec.ExecuteException;

/**
 * TG Server to be lifecycled.
 * 
 * @author sbagi@tibco.com
 *
 */
public class TGServer {

	final private static String process = "tgdb";

	private File home;

	//private File workingDir;

	private File configFile;
	private File initFile;
	
	private File logFile;
	private String logFileBase;
	private File initLogFile;
	private String initLogFileBase;

	private int pid;

	final private String name = "tgdb"; // server name

	private String dbName;

	private File dbPath;
	
	private String systemUser;
	
	private String systemPwd;

	//private ByteArrayOutputStream outStream = new ByteArrayOutputStream();
	//private ByteArrayOutputStream errStream = new ByteArrayOutputStream();

	private String outIncrement = "";
	
	private boolean running = false;
	
	private String banner = "";

	/**
	 * Create a TG server
	 * 
	 * @param tgHome
	 *            TG server home
	 * @throws TGGeneralException
	 *             File path not found
	 */
	public TGServer(String tgHome) throws TGGeneralException {
		if (tgHome == null)
			throw new TGGeneralException("TGServer - TGDB home is not defined");
		File ftgHome = new File(tgHome);
		if (!ftgHome.exists())
			throw new TGGeneralException("TGServer - TGDB home '" + tgHome + "' does not exist");

		this.setHome(ftgHome);
		
	}

	/**
	 * Create a TG server
	 * 
	 * @param tgHome
	 *            TG server home
	 * @param tgConfig
	 *            TG server start-up config file
	 * @throws TGGeneralException
	 *             File path not found or mis-formatted file
	 */
	public TGServer(String tgHome, String tgConfig) throws TGGeneralException {
		this(tgHome);
		
		File fconfigFile = new File(tgConfig);
		if (!fconfigFile.exists())
			throw new TGGeneralException("TGServer - Config file '" + tgConfig + "' does not exist");
		this.setConfigFile(fconfigFile);
		this.setLogFile(this.name + "_" + this.getDbName());
		/*try { // this is done in setConfig()
			this.setDbName(fconfigFile);
			this.setDbPath(fconfigFile);
			this.setLogFile(fconfigFile);
			this.setNetListeners(fconfigFile);
		}
		catch(IOException ioe) {
			throw new TGGeneralException(ioe.getMessage());
		}*/
	}
	
	
	/**
	 * Create a TG server
	 * 
	 * @param tgHome
	 *            TG server home
	 * @param tgConfig
	 *            TG server start-up config file
	 * @param tgLog
	 *            TG server log filename (just filename - no path, no file extension). 
	 *            Log will be located in tgdb_home/bin/log folder
	 * @return 
	 * @throws TGGeneralException
	 *             File path not found or mis-formatted file
	 */
	public TGServer(String tgHome, String tgConfig, String tgLog) throws TGGeneralException {
		this(tgHome, tgConfig);
		this.setLogFile(tgLog);
	
	}
	
	/**
	 * Create a TG server
	 * 
	 * @param tgHome
	 *            TG server home
	 * @param tgInit
	 *            TG server init config file
	 * @param tgConfig
	 *            TG server start-up config file
	 * @throws TGGeneralException
	 *             File path not found or mis-formatted file
	 */
	/*
	public TGServer(String tgHome, String tgInit, String tgConfig) throws TGGeneralException {
		this(tgHome, tgConfig);
		
		File finitFile = new File(tgInit);
		if (!finitFile.exists())
			throw new TGGeneralException("TGServer - Init file '" + tgInit + "' does not exist");
		this.setInit(finitFile);
	}
	*/
	
	/**
	 * Network listener for a given TG server
	 * 
	 * @author sbagi@tibco.com
	 *
	 */
	public class NetListener {
		private String name;
		private String host;
		private String url;
		private int port;
		private int maxConnections;

		private NetListener() {
			;
		}

		private NetListener(BufferedReader br) throws NumberFormatException, IOException, TGGeneralException {
			this.setName(br);
			this.setHost(br);
			this.setPort(br);
			this.setMaxConnections(br);
			this.setUrl();
			
		}

		/**
		 * Get the listener name as defined in config file
		 * 
		 * @return listener name
		 */
		public String getName() {
			return this.name;
		}

		private void setName(BufferedReader br) throws IOException {
			String tmpName = "";
			String line = null;
			while ((line = br.readLine()) != null) {
				if (line.startsWith("name")) {
					if (line.indexOf("//") != -1)
						tmpName = line.substring(line.indexOf("=") + 1, line.indexOf("//")).trim();
					else
						tmpName = line.substring(line.indexOf("=") + 1).trim();
					break;
				}
			}
			this.name = tmpName;
		}

		/**
		 * Get the listener host as defined in config file
		 * 
		 * @return Listener host
		 */
		public String getHost() {
			return this.host;
		}

		/**
		 * Is the listener bound to an IPv6 address ?
		 * 
		 * @return true if listener bound to IPv6 address, false otherwise
		 * @throws TGGeneralException Cannot determine whether IPv6
		 */
		public boolean isIPv6() throws TGGeneralException {
			
			if (host != null && !host.equals("")) {
				InetAddress address;
				try {
					address = InetAddress.getByName(this.host);
					if (address instanceof Inet4Address) {
						return false;
					} else if (address instanceof Inet6Address) {
						return true;
					} else
						return false;
				}
				catch (UnknownHostException e) {
					throw new TGGeneralException("TGServer - " + e.getMessage());
				}
			} else
				throw new TGGeneralException("TGServer - Cannot whether determine IPv6 - Host is not set");
			
		}
		
		private void setHost(BufferedReader br) throws IOException {
			String tmpHost = "";
			String line = null;
			while ((line = br.readLine()) != null) {
				if (line.startsWith("host")) {
					if (line.indexOf("//") != -1)
						tmpHost = line.substring(line.indexOf("=") + 1, line.indexOf("//")).trim();
					else
						tmpHost = line.substring(line.indexOf("=") + 1).trim();
					break;
				}
			}
			if (tmpHost.equals("0.0.0.0")) {
				System.setProperty("java.net.preferIPv6Addresses", "false");
				if (Inet4Address.getLocalHost() instanceof Inet6Address)
					tmpHost = "localhost";
				else
					tmpHost = Inet4Address.getLocalHost().getHostAddress();
			}
			else if (tmpHost.equals("::")) {
				System.setProperty("java.net.preferIPv6Addresses", "true");
				tmpHost = InetAddress.getLocalHost().getHostName();
				InetAddress[] addr = InetAddress.getAllByName(tmpHost);
		    	for (InetAddress address : addr) {
		    		if (address instanceof Inet6Address) 
		    			tmpHost = address.getHostAddress().substring(0, (address.getHostAddress().contains("%")?address.getHostAddress().indexOf('%'):address.getHostAddress().length()));
		    	}
			}
			else 
				;
			this.host = tmpHost;
		}

		/**
		 * Get the listener port as defined in config file
		 * 
		 * @return listener port
		 */
		public int getPort() {
			return this.port;
		}
		
		/**
		 * Get the listener url
		 * 
		 * @return listener url
		 */
		public String getUrl() {
			return this.url;
		}

		private void setPort(BufferedReader br) throws NumberFormatException, IOException {
			int tmpPort = 0;
			String line = null;
			while ((line = br.readLine()) != null) {
				if (line.startsWith("port")) {
					if (line.indexOf("//") != -1)
						tmpPort = Integer.parseInt(line.substring(line.indexOf("=") + 1, line.indexOf("//")).trim());
					else
						tmpPort = Integer.parseInt(line.substring(line.indexOf("=") + 1).trim());
					break;
				}
			}
			this.port = tmpPort;
		}
		
		private void setUrl() throws TGGeneralException  {
			this.url = "tcp://" + (this.isIPv6()?"["+this.host+":"+this.port+"]":this.host+":"+this.port);
			
		}

		/**
		 * Get the max connection as defined in config file
		 * 
		 * @return listener max connection
		 */
		public int getMaxConnections() {
			return this.maxConnections;
		}

		private void setMaxConnections(BufferedReader br) throws NumberFormatException, IOException {
			int tmpConnections = 0;
			String line = null;
			while ((line = br.readLine()) != null) {
				if (line.startsWith("maxconnections")) {
					if (line.indexOf("//") != -1)
						tmpConnections = Integer
								.parseInt(line.substring(line.indexOf("=") + 1, line.indexOf("//")).trim());
					else
						tmpConnections = Integer.parseInt(line.substring(line.indexOf("=") + 1).trim());
					break;
				}
			}
			this.maxConnections = tmpConnections;
		}
	}

	private List<NetListener> netListeners = new ArrayList<NetListener>();

	

	/**
	 * Get the list of network listeners as defined in config file.
	 * 
	 * @return Array of listeners
	 */
	public NetListener[] getNetListeners() {
		return netListeners.toArray(new NetListener[0]);

	}
	
	/**
	 * Get the number of network listeners as defined in config file.
	 * 
	 * @return Number of listeners
	 */
	public int getNetListenersCount() {
		return this.getNetListeners().length;
	}

	private void setNetListeners(File configFile) throws Exception {

		BufferedReader br = new BufferedReader(new FileReader(configFile));
		this.netListeners.clear(); // reset
		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.startsWith("[netlistener]")) {
				this.netListeners.add(new NetListener(br));
			}
		}
		br.close();
	}

	/**
	 * Get the TG home folder
	 * 
	 * @return Full path of TG home
	 */
	public File getHome() {
		return home;
	}

	private void setHome(File tgHome) {
		this.home = tgHome;
	}

	/**
	 * Get TG server config file
	 * 
	 * @return full path of config file
	 */
	public File getConfigFile() {
		return this.configFile;
	}
	
	/**
	 * Get TG server init file
	 * 
	 * @return full path of init file
	 */
	public File getInitFile() {
		return this.initFile;
	}

	/**
	 * Set the config file TG server is running against
	 * 
	 * @param tgConfig
	 *            Full path of config file
	 * @throws TGGeneralException
	 *             File path not found or mis-formatted file
	 */
	public void setConfigFile(File tgConfig) throws TGGeneralException {
		
		if (this.running)
			throw new TGGeneralException("TGServer - Cannot set config file on a running server");
		if (!tgConfig.exists())
			throw new TGGeneralException("TGServer - Config file '" + tgConfig + "' does not exist");

		this.configFile = tgConfig;

		// (Re-)define the following parameters with this config file
		try {
			this.setDbName(tgConfig);
			this.setDbPath(tgConfig);
			//this.setLogFile(tgConfig);
			this.setNetListeners(tgConfig);
		}
		catch(Exception e) {
			throw new TGGeneralException(e.getMessage());
		}
	}
	
	/**
	 * Set the init files TG server is running against
	 * 
	 * @param tgInit
	 *            Full path of init file
	 * @throws TGGeneralException
	 *             File path not found or mis-formatted file
	 */
	private void setInit(File tgInit) throws TGGeneralException {
		
		if (this.running)
			throw new TGGeneralException("TGServer - Cannot set init file on a running server");
		if (!tgInit.exists())
			throw new TGGeneralException("TGServer - Init file '" + tgInit + "' does not exist");

		this.initFile = tgInit;

		try {
			// Init log file as well
			//this.setInitLogFile(tgInit);
			this.setInitLogFile("tgdb_initdb");
			// (Re-)define the following parameters with this init file
			this.setSystemUser(tgInit);
			this.setSystemPwd(tgInit);
		}
		catch(IOException ioe) {
			throw new TGGeneralException(ioe.getMessage());
		}
	}

	/**
	 * Get database path as defined in config file.
	 * 
	 * @return DB full path
	 */
	public File getDbPath() {
		return this.dbPath;
	}

	private void setDbPath(File configFile) throws IOException {
		String tmpPath = "";
		BufferedReader br = new BufferedReader(new FileReader(configFile));

		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.startsWith("[" + this.getName() + "]")) {
				String subLine = null;
				while ((subLine = br.readLine()) != null) {
					if (subLine.startsWith("dbpath")) {
						if (subLine.indexOf("//") != -1)
							tmpPath = subLine.substring(subLine.indexOf("=") + 1, subLine.indexOf("//")).trim();
						else
							tmpPath = subLine.substring(subLine.indexOf("=") + 1).trim();
						break;
					}
				}
				break;
			}
		}
		br.close();

		if (tmpPath.startsWith("."))
			tmpPath = this.home + "/bin/" + tmpPath;
		this.dbPath = new File(tmpPath);
	}

	/**
	 * Get database name a defined in config file
	 * 
	 * @return DB name
	 */
	public String getDbName() {
		return dbName;
	}

	private void setDbName(File configFile) throws IOException {
		String tmpName = "";
		BufferedReader br = new BufferedReader(new FileReader(configFile));

		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.startsWith("[" + this.getName() + "]")) {
				String subLine = null;
				while ((subLine = br.readLine()) != null) {
					if (subLine.startsWith("name")) {
						if (subLine.indexOf("//") != -1)
							tmpName = subLine.substring(subLine.indexOf("=") + 1, subLine.indexOf("//")).trim();
						else
							tmpName = subLine.substring(subLine.indexOf("=") + 1).trim();
						break;
					}
				}
				break;
			}
		}
		br.close();
		this.dbName = tmpName;
	}

	/**
	 * Get the engine name.
	 * 
	 * @return Engine name - "tgdb"
	 */
	public String getName() {
		return this.name;
	}

	/**
	 * Get the TG server PID
	 * 
	 * @return PID
	 * @throws TGGeneralException TG server does not have a PID probably due to start-up failure.
	 */
	public int getPid() throws TGGeneralException {
		if (this.pid == 0)
			throw new TGGeneralException("TGServer - Server does not have a PID - Probably due to start-up failure or crash");
		else
			return this.pid;
	}

	void setPid(int pid) {
		this.pid = pid;
	}

	/**
	 * Get the TG server log file as defined in the config file
	 * 
	 * @return Full path of log file
	 */
	public File getLogFile() {
		return this.logFile;
	}

	/**
	 * Set the log file based on what is in the config file
	 * @param configFile
	 * @throws IOException
	 */
	/*
	private void setLogFile(File configFile) throws IOException {

		String tmpLog = "";
		BufferedReader br = new BufferedReader(new FileReader(configFile));

		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.startsWith("[logger]")) {
				String subLine = null;
				while ((subLine = br.readLine()) != null) {
					if (subLine.startsWith("path")) {
						if (subLine.indexOf(";") != -1)
							tmpLog = subLine.substring(subLine.indexOf("=") + 1, subLine.indexOf(";")).trim();
						else
							tmpLog = subLine.substring(subLine.indexOf("=") + 1).trim();
						break;
					}
				}
				break;
			}
		}

		br.close();
		if (tmpLog.startsWith("."))
			tmpLog = this.home + "/bin/" + tmpLog;
		tmpLog = tmpLog + "/" + this.name + "_" + this.dbName + ".log";
		this.logFile = new File(tmpLog);

	}
	*/
	
	/**
	 * Set the log file based on the -l inline parameter
	 * @param logFile filename with no .log extension
	 */
	private void setLogFile(String fileName) {
		this.logFileBase = fileName;
		this.logFile = new File(this.home + "/bin/log/" + fileName + ".log");
	}
	
	/*
	private void setInitLogFile(File initFile) throws IOException {

		String tmpLog = "";
		BufferedReader br = new BufferedReader(new FileReader(initFile));

		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.startsWith("[logger]")) {
				String subLine = null;
				while ((subLine = br.readLine()) != null) {
					if (subLine.startsWith("path")) {
						if (subLine.indexOf(";") != -1)
							tmpLog = subLine.substring(subLine.indexOf("=") + 1, subLine.indexOf(";")).trim();
						else
							tmpLog = subLine.substring(subLine.indexOf("=") + 1).trim();
						break;
					}
				}
				break;
			}
		}
		br.close();
		
		if (tmpLog.equals("")) { // no logger defined in initdb.conf
			tmpLog = this.home + "/bin/log/tgdb_initdb.log";
		}
		else { // logger defined in initdb.conf
			if (tmpLog.startsWith("."))
				tmpLog = this.home + "/bin/" + tmpLog;
			tmpLog = tmpLog + "/" + this.name + "_" + this.dbName + ".log";
		}
		this.initLogFile = new File(tmpLog);

	}
	*/
	
	private void setInitLogFile(String fileName) {
		this.initLogFileBase = fileName;
		this.initLogFile = new File(this.home + "/bin/log/" + fileName + ".log");
	}
	
	/**
	 * Get the TG server init log file as defined in the init config file
	 * 
	 * @return Full path of log file
	 */
	public File getInitLogFile() {
		return this.initLogFile;
	}

	/**
	 * Get the TG server system user name as defined in the init file
	 * 
	 * @return System user name
	 */
	public String getSystemUser() {
		return this.systemUser;
	}
	
	private void setSystemUser(File initFile) throws IOException {

		String tmpUser = "";
		BufferedReader br = new BufferedReader(new FileReader(initFile));

		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.startsWith("[initdb]")) {
				String subLine = null;
				while ((subLine = br.readLine()) != null) {
					if (subLine.startsWith("sysuser")) {
						if (subLine.indexOf("//") != -1)
							tmpUser = subLine.substring(subLine.indexOf("=") + 1, subLine.indexOf("//")).trim();
						else
							tmpUser = subLine.substring(subLine.indexOf("=") + 1).trim();
						break;
					}
				}
				break;
			}
		}

		br.close();
		this.systemUser = tmpUser;
	}
	
	/**
	 * Get the TG server system user password as defined in the init file
	 * 
	 * @return System user password
	 */
	public String getSystemPwd() {
		return this.systemPwd;
	}
	
	private void setSystemPwd(File initFile) throws IOException {

		String tmpPwd = "";
		BufferedReader br = new BufferedReader(new FileReader(initFile));

		String line = null;
		while ((line = br.readLine()) != null) {
			if (line.startsWith("[initdb]")) {
				String subLine = null;
				while ((subLine = br.readLine()) != null) {
					if (subLine.startsWith("syspasswd")) {
						if (subLine.indexOf("//") != -1)
							tmpPwd = subLine.substring(subLine.indexOf("=") + 1, subLine.indexOf("//")).trim();
						else
							tmpPwd = subLine.substring(subLine.indexOf("=") + 1).trim();
						break;
					}
				}
				break;
			}
		}

		br.close();
		this.systemPwd = tmpPwd;
	}
	
	/**
	 * Is TG server up and running
	 * @return true if running, false otherwise
	 */
	public boolean isRunning() {
		return this.running;
	}

	void setRunning(boolean running) {
		this.running = running;
	}
	
	/**
	 * Get the full TG server standard output after start-up
	 * 
	 * @return Standard output console
	 * @throws TGGeneralException if output cannot be determined
	 */
	public String getOutput() throws TGGeneralException {
		//String outString = this.outStream.toString().replaceAll("\\x00", ""); // remove all Null characters !
		//this.outIncrement = outString.substring(outString.length());
		//return outString;
		// return this.outStream.toString();
		String outString = "";
		try {
			outString = new String(Files.readAllBytes(Paths.get(this.logFile.toURI())));
			this.outIncrement = outString.substring(outString.length());
			
		} catch (IOException ioe) {
			throw new TGGeneralException(ioe.getMessage());
		}
		return outString;
		
	}
	
	/**
	 * Get the full TG server banner after start-up
	 * 
	 * @param string 
	 * @return TG server banner
	 * @throws TGGeneralException If banner cannot be retrieved
	 */
	private void setBanner(String string) throws TGGeneralException {

		boolean bannerFound = false;
		String banner = "";
		try (BufferedReader reader = new BufferedReader(new StringReader(string))) {
            String line = reader.readLine();
            boolean printBanner = false; 
            while (line != null) {
                if (line.startsWith("**********") && !printBanner) {
                	printBanner = true;
                	banner = banner + "\n" + line;
                }
                else if (line.startsWith("**********") && printBanner) {
                	banner = banner + "\n" + line;
                	printBanner = false;
                }
                else if (printBanner) {
                	banner = banner + "\n" + line;
                	bannerFound = true;
                }
                line = reader.readLine();
            }
        } catch (IOException e) {
            throw new TGGeneralException(e.getMessage());
        }
		if (!bannerFound)
			throw new TGGeneralException("TGServer - TG server banner not found");
		this.banner= banner;
	}
	
	/**
	 * Get the full TG server banner
	 * 
	 * @return TG server banner
	 */
	public String getBanner() {
		return this.banner;
	}

	/**
	 * Get the incremented TG server standard output. This avoids to
	 * display the full output but just a part of it.
	 * 
	 * @return Part of the standard output that was never displayed so far
	 * @throws TGGeneralException if increment cannot be determined
	 */
	public String getOutputIncrement() throws TGGeneralException {
		String outString;
		if (this.outIncrement.equals(""))
			outString = this.getOutput();
		else
			outString = this.outIncrement;

		// Compute increment for next time
		this.outIncrement = outString
				.substring(this.getOutput().indexOf(outString) + outString.length());
		return outString;
	}

	/**
	 * Get the standard error of the TG server
	 * 
	 * @return Standard error console
	 */
	//public String getError() {
	//	return this.errStream.toString().replaceAll("\\x00", ""); // remove all Nul characters !
	//}
	
	/**
	 * Get the full name including version, build and edition
	 * 
	 * @return Full product name
	 */
	public String getFullName() {

		String fullName = "";
		try (BufferedReader reader = new BufferedReader(new StringReader(this.getBanner()))) {
            String line = reader.readLine();
            while (line != null) {
                if (line.matches("TIBCO\\(R\\) Graph Database [1-9]+\\.[0-9]+\\.[0-9]+ Build\\([1-9]+\\) (Enterprise|Community) Edition\\.")) {
                	fullName = line;
                }
                line = reader.readLine();
            }
        } catch (IOException e) {
            ;
        }
		return fullName;
	}

	/**
	 * Get the error statements found in the log file.
	 * 
	 * @return List of error statements
	 * @throws IOException
	 *             problem while reading log file
	 */
	public List<String> getErrorsInLog() throws IOException {

		List<String> errorStmts = new ArrayList<String>();
		BufferedReader reader = new BufferedReader(new FileReader(this.logFile));
		String line = reader.readLine();
		while (line != null) {
			if (line.matches("^.*[0-9][0-9] Error .*$")) { // matches last 2 digits of timestamp followed by Error
				errorStmts.add(line);
			}
			line = reader.readLine();
		}
		reader.close();
		return errorStmts;
	}

	/**
	 * Initialize the TG server synchronously. This Init operation blocks until
	 * it is completed.
	 * 
	 * @param initFile
	 *            TG server init config file
	 * @param forceCreation
	 * 			  Force creation. Delete all the data in the db directory first.
	 * @param timeout
	 *            Number of milliseconds allowed to initialize the server
	 * @return the output stream of init operation 
	 * @throws TGInitException
	 *             Init operation fails or timeout occurs 
	 */
	public String init(String initFile, boolean forceCreation, long timeout) throws TGInitException {

		File initF = new File(initFile);
		if (!initF.exists())
			throw new TGInitException("TGServer - Init file '" + initFile + "' does not exist");
		try {
			this.setInit(initF);
		}
		catch(TGGeneralException e) {
			throw new TGInitException(e.getMessage());
		}
		
		//ByteArrayOutputStream output = new ByteArrayOutputStream();
		
		PumpStreamHandler psh = new PumpStreamHandler(new ByteArrayOutputStream());
		Executor tgExec = new DefaultExecutor();
		tgExec.setStreamHandler(psh);
		tgExec.setWorkingDirectory(new File(this.home + "/bin"));
		CommandLine tgCL = new CommandLine((new File(this.home + "/bin/"+process)).getAbsolutePath());
		if (forceCreation)
			tgCL.addArguments(new String[] { "-i", "-f", "-Y", "-c", initFile, "-l", this.initLogFileBase });
		else
			tgCL.addArguments(new String[] { "-i", "-Y", "-c", initFile, "-l", this.initLogFileBase });

		ExecuteWatchdog tgWatch = new ExecuteWatchdog(timeout);
		tgExec.setWatchdog(tgWatch);
		System.out.println("TGServer - Initializing " + StringUtils.toString(tgCL.toStrings()," "));
		String output = "";
		try {
			tgExec.execute(tgCL);
			output = new String(Files.readAllBytes(Paths.get(this.getInitLogFile().toURI())));
		} catch (IOException ee) {
			if (tgWatch.killedProcess())
				throw new TGInitException("TGServer - Init did not complete within " + timeout + " ms");
			else {
				try {
					Thread.sleep(1000); // make sure output has time to fill up
					
				}
				catch(InterruptedException ie) {;}
				throw new TGInitException("TGServer - Init failed: " + ee.getMessage(), output);
			}
		}
		try {
			this.setBanner(output);
		}
		catch(TGGeneralException tge) {
			throw new TGInitException(tge.getMessage());
		}
		
		if (output.contains("TGSuccess")) {
			System.out.println("TGServer - Initialized successfully");
			return output;
		}
		else 
			throw new TGInitException("TGServer - Init failed", output);
	}

	/**
	 * Start the TG server synchronously.
	 * 
	 * @param timeout
	 *            Number of milliseconds allowed to start the server
	 * @throws TGStartException Start operation fails
	 */
	public void start(long timeout) throws TGStartException {

		if (this.configFile == null)
			throw new TGStartException("TGServer - Config file not set");
		
		if (this.logFile == null)
			this.setLogFile("tgdb_" + this.dbName);
		
		//this.outStream = new ByteArrayOutputStream(); // reset
		//this.errStream = new ByteArrayOutputStream(); // reset
		PumpStreamHandler psh = new PumpStreamHandler(new ByteArrayOutputStream());
		DefaultExecuteResultHandler resultHandler = new DefaultExecuteResultHandler();

		Executor tgExec = new DefaultExecutor();
		tgExec.setWorkingDirectory(new File(this.home + "/bin"));
		tgExec.setStreamHandler(psh);
		CommandLine tgCL = new CommandLine((new File(this.home + "/bin/"+process)).getAbsolutePath());
		tgCL.addArguments(new String[] { "-s", "-c", this.configFile.getAbsolutePath(), "-l", this.logFileBase });

		System.out.println("TGServer - Starting " + StringUtils.toString(tgCL.toStrings()," "));
		try {
			tgExec.execute(tgCL, resultHandler);
		} catch (IOException ioe) {
			try {
				Thread.sleep(1000); // Make sure output/error fill up
			}
			catch(InterruptedException ie) {;}
			throw new TGStartException(ioe.getMessage());
		}
		
		if (timeout > 0) {
			Calendar future = Calendar.getInstance();
			future.add(Calendar.MILLISECOND, (int) timeout);
			boolean started = false;
			boolean error = false;
			List<String> acceptedClients = new ArrayList<String>();
			String acceptedClient;
			while (!future.before(Calendar.getInstance())) {
				try {
					Thread.sleep(1000);
					BufferedReader reader = new BufferedReader(new StringReader(this.getOutput()));
					String line = reader.readLine();
					while (line != null) {
						if (line.contains("Process pid:"))
							this.setPid(Integer.parseInt(line.substring(line.lastIndexOf("Process pid:") + 12,
									line.indexOf(",", line.lastIndexOf("Process pid:") + 12))));
						if (line.contains("[Error]")) {
							error = true;
						}
						if (line.contains("Accepting clients on")) {
							started = true;
							acceptedClient = line.substring(line.indexOf("- Accepting clients on") + 2);
							if (!acceptedClients.contains(acceptedClient))
								acceptedClients.add(acceptedClient);
						}
						line = reader.readLine();
					}
					reader.close();
					if (started)
						break;
				} 
				catch (Exception e) {
					throw new TGStartException("TGServer - " + e.getMessage());
				}
			}
			if (!started)
				throw new TGStartException("TGServer - Did not start on time (after " + timeout
						+ " msec) - See log " + this.logFile);
			else {
				this.running = true;
				System.out.println("TGServer - Started successfully with pid " + this.pid + " :");
				System.out.println("\t\t- Log file: " + this.logFile);
				if (error)
					System.out.println("\t\t- With some error(s) - See log");
				for (String client : acceptedClients)
					System.out.println("\t\t- " + client);
				try {
					this.setBanner(this.getOutput());
				}
				catch(TGGeneralException tge) {
					throw new TGStartException(tge.getMessage());
				}
			}
		}
	}

	/**
	 * Not implemented
	 */
	public void restore() {
		;
	}

	/**
	 * Not implemented
	 */
	public void backup() {
		;
	}

	/**
	 * <pre>
	 * Kill the TG server.
	 * - taskkill on Windows.
	 * - kill -9 on Unix.
	 * </pre>
	 * 
	 * @throws Exception
	 *             Kill operation fails
	 */
	public void kill() throws Exception {

		if (this.pid == 0)
			throw new TGGeneralException("TG server does not have a PID - Probably due to a previous start-up failure");
		
		ByteArrayOutputStream output = new ByteArrayOutputStream();
		PumpStreamHandler psh = new PumpStreamHandler(output);
		DefaultExecutor executor = new DefaultExecutor();
		executor.setStreamHandler(psh);
		CommandLine cmdLine;
		if (OS.isFamilyWindows())
			cmdLine = CommandLine.parse("taskkill /f /pid " + this.getPid() + " /t");
		else
			cmdLine = CommandLine.parse("kill -9 " + this.getPid() + "");
		try {
			executor.execute(cmdLine);
		} catch (ExecuteException ee) {
			// System.out.println("TGServer with pid " + this.getPid() + " not killed :");
			// System.out.println("\t- " + output.toString().trim().replace("\n","\n\t- "));
			throw new ExecuteException(output.toString().trim(), 1); // re-throw with better message
		}
		System.out.println("TGServer - Server with pid " + this.getPid() + " successfully killed :");
		if (!output.toString().equals(""))
			System.out.println("\t\t- " + output.toString().trim().replace("\n", "\n\t\t- "));
		this.running = false;
		this.pid = 0;
	}

	/**
	 * <pre>
	 * Kill all the TG server processes.
	 * Note that this method blindly tries to kill all the servers of the machine 
	 * and do not update the running status of those servers.
	 * - taskkill on Windows.
	 * - kill -9 on Unix.
	 * </pre>
	 * 
	 * @throws Exception Kill operation fails
	 */
	public static void killAll() throws Exception {

		ByteArrayOutputStream output = new ByteArrayOutputStream();
		PumpStreamHandler psh = new PumpStreamHandler(output);
		DefaultExecutor executor = new DefaultExecutor();
		executor.setStreamHandler(psh);
		executor.setWorkingDirectory(new File(System.getProperty("java.io.tmpdir")));
		CommandLine cmdLine;
		if (OS.isFamilyWindows())
			cmdLine = CommandLine.parse("taskkill /f /im " + process + ".exe");
		else { // Unix
			File internalScriptFile = new File(
					ClassLoader.getSystemResource(TGServer.class.getPackage().getName().replace('.', '/') + "/TGKillProcessByName.sh").getFile());
			File finalScriptFile = new File(executor.getWorkingDirectory() + "/" + internalScriptFile.getName());
			Files.copy(internalScriptFile.toPath(), finalScriptFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
			cmdLine = CommandLine.parse("sh " + finalScriptFile.getAbsolutePath() + " " + process + "");
		}
		try {
			Thread.sleep(1000);
			executor.execute(cmdLine);
		} catch (ExecuteException ee) {
			if (output.toString().contains("ERROR: The process \"" + process + (OS.isFamilyWindows() ? ".exe\"" : "") + " not found"))
				return;
			throw new ExecuteException(output.toString().trim(), 1); // re-throw with better message
		}
		// Check one more thing :
		// On Windows when some processes do not get killed taskkill still
		// returns exit code
		if (OS.isFamilyWindows()) {
			if (output.toString().contains("ERROR")) 
				throw new ExecuteException(output.toString().trim(), 1);
		}
		else {
			if (output.toString().contains("ERROR: The process \"" + process + "\" not found"))
				return;
		}

		System.out.println("TGServer - Server(s) successfully killed :");
		if (!output.toString().equals(""))
			System.out.println("\t\t- " + output.toString().trim().replace("\n", "\n\t\t- "));
	}
	
	/**
	 * Returns a string representation of this TG server. 
	 * 
	 */
	public String toString() {
		String str = this.getBanner() + "\n";
		str = "TGServer - Description:\n";
		str = str + "\t\t- Home:\t\t" + this.getHome() + "\n";
		if (this.getConfigFile() == null) 
			str = str + "\t\t- Config:\tNot set\n";
		else {	
			str = str + "\t\t- Config:\t" + this.getConfigFile() + "\n";
			str = str + "\t\t- Log:\t\t" + this.getLogFile() + "\n";
			str = str + "\t\t- DB Name:\t" + this.getDbName() + "\n";
			str = str + "\t\t- DB Path:\t" + this.getDbPath() + "\n";
		}
		str = str + "\t\t- Running:\t" + this.isRunning() + "\n";
		if (this.isRunning()) {
			try {
				str = str + "\t\t- PID:\t\t" + this.getPid() + "\n";
			} catch (TGGeneralException e) {
				;
			}
			for (int i=0; i<this.getNetListeners().length; i++) {
				str = str + "\t\t- Listener:\t" + this.getNetListeners()[i].getName() + " - " + this.getNetListeners()[i].getHost() + ":" + this.getNetListeners()[i].getPort() +  "\n";
			}
		}
		return str;
	}
}
