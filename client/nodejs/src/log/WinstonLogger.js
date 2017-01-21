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
 */
var winston    = require('winston'),
    util       = require('util'),
    TGLogger   = require('./TGLogger').TGLogger,
    TGLogLevel = require('./TGLogger').TGLogLevel;

function WinstonLogger() {
	this._LevelNames = [
	    'None',
	    'Fatal',
	    'Error',
	    'Warning',
	    'Info',
	    'Debug',
	    'DebugWire'
	];
	
	this._logger = new winston.Logger({
		levels : TGLogLevel,
		transports : [new(winston.transports.Console)({colorize: true})],
	    colors: {
	    	DebugWire : 'green',
	    	Debug : 'green',
	    	Info : 'green',
	        Warning : 'yellow',
	        Error : 'red',
	        Fatal : 'red',
	        None : 'green'
	    }
	});
		
	//this._logger.addColors(myCustomLevels.colors);
}

util.inherits(WinstonLogger, TGLogger);

WinstonLogger.prototype.logFatal = function(format, args){
	if(this.isFatal()) {
		this._logger.Fatal.apply(this, arguments);		
	}
};

WinstonLogger.prototype.logError = function(format, args){
	if(this.isError()) {
		this._logger.Error.apply(this, arguments);		
	}
};

WinstonLogger.prototype.logWarning = function(format, args){
	if(this.isWarning()) {
		this._logger.Warning.apply(this, arguments);		
	}
};

WinstonLogger.prototype.logInfo = function(format, args){
	if(this.isInfo()) {
		this._logger.Info.apply(this, arguments);		
	}
};

WinstonLogger.prototype.logDebug = function(format, args){
	if(this.isDebug()) {
		this._logger.Debug.apply(this, arguments);		
	}
};

WinstonLogger.prototype.logDebugWire = function(format, args){
	if(this.isDebugWire()) {
		this._logger.DebugWire.apply(this, arguments);		
	}
};

WinstonLogger.prototype.log = function(wLevel, format, args){
	//this._logger.log.apply(this, arguments);
};

WinstonLogger.prototype.setlevelImpl = function(wLevel){
	this._logger.level = this._LevelNames[wLevel];
};

exports.WinstonLogger = WinstonLogger;

function test() {
	var logger = new WinstonLogger();
	logger.setLevel(TGLogLevel.Info);
	logger.logDebugWire('Hello %s %s', 'Steven', '!!!');
	logger.logDebug('Hello %s %s', 'Steven', '!!!');
	logger.logInfo('Hello %s %s', 'Steven', '!!!');
	logger.logWarning('Hello %s %s', 'Steven', '!!!');
	logger.logError('Hello %s %s', 'Steven', '!!!');
	logger.logFatal('Hello %s %s', 'Steven', '!!!');
}

//test();