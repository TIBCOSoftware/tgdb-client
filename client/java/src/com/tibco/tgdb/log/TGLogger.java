package com.tibco.tgdb.log;

/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * File name :TGLogger
 * Created on: 12/18/14
 * Created by: suresh
 *
 * SVN Id: $Id: TGLogger.java 623 2016-03-19 21:41:13Z ssubrama $
 */

/*
A Simple logging interface
 */
public interface TGLogger {

    /**
     * LogLevel
     */
    public enum TGLevel {
        None,
        Fatal,
        Error,
        Info,
        Warning,
        Debug,
        DebugWire
    }
    /**
     * log message
     * @param level  - The log level
     * @param format - The string format
     * @param args   - Args to the format
     */
    void log(TGLevel level, String format, Object ... args );

    /**
     * Log an Exception
     * @param msg A message to log
     * @param e The exception associated with the message
     */
    void logException(String msg, Exception e);


    /**
     * Is Log enabled for this level
     * @param level The level to check for the logger
     * @return boolean value indicating if it is enabled or not
     */
    boolean isEnabled(TGLevel level);

    /**
     * Set the Log Level
     * @param level the log level dynamically
     */
    void setLevel(TGLevel level);


}
