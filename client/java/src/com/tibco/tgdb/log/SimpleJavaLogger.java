package com.tibco.tgdb.log;

import java.util.logging.*;

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

 * File name :SimpleJavaLogger
 * Created by: suresh
 * SVN Id: $Id: SimpleJavaLogger.java 723 2016-04-16 19:21:18Z vchung $
 */
public class SimpleJavaLogger implements TGLogger {
    private static Logger logger = Logger.getLogger(SimpleJavaLogger.class.getPackage().getName());
    private static Level logLevel[] = null;
    private static TGLevel defaultLogLevel = TGLevel.Info;

    public SimpleJavaLogger() {
    	logLevel = new Level[TGLevel.DebugWire.ordinal() + 1];
    	logLevel[TGLevel.None.ordinal()] = Level.OFF;
        logLevel[TGLevel.Fatal.ordinal()] = Level.SEVERE;
        logLevel[TGLevel.Error.ordinal()] = Level.SEVERE;
        logLevel[TGLevel.Info.ordinal()] = Level.INFO;
        logLevel[TGLevel.Warning.ordinal()] = Level.WARNING;
        logLevel[TGLevel.Debug.ordinal()] = Level.FINE;
        logLevel[TGLevel.DebugWire.ordinal()] = Level.FINEST;
        setLevel(defaultLogLevel);
    }

    @Override
    public void log(TGLevel level, String format, Object... args) {
        logger.log(logLevel[level.ordinal()], String.format(format, args));
    }

    @Override
    public void logException(String msg, Exception e) {
        logger.log(Level.WARNING, msg, e);
    }

    @Override
    public boolean isEnabled(TGLevel level) {
    	return logger.isLoggable(logLevel[level.ordinal()]);
    }

    @Override
    public void setLevel(TGLevel level) {
    	logger.setLevel(logLevel[level.ordinal()]);
    }
}
