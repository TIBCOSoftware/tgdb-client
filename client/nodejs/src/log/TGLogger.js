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
var util    = require('util');
var TGLogLevel = {
        None      : 0,
        Fatal     : 1,
        Error     : 2,
        Warning   : 3,
        Info      : 4,
        Debug     : 5,
        DebugWire : 6
};

exports.TGLogLevel = TGLogLevel;

function TGLogger () {
	this._tgLogLevel = null;
}

/**
 * log message
 * @param level  - The log level
 * @param format - The string format
 * @param args   - Args to the format
 */
TGLogger.prototype.log = function(tgLevel, format, args ) {
	if(this._tgLogLevel===TGLogLevel.None) {
		return;
	}
	this.logImpl.apply(this, arguments);		
};

/**
 * Log an Exception
 * @param msg A message to log
 * @param e The exception associated with the message
 */
TGLogger.prototype.logException = function(msg, exception) {
    	
};


/**
 * Is Log enabled for this level
 * @param level The level to check for the logger
 * @return boolean value indicating if it is enabled or not
 */
TGLogger.prototype.isEnabled = function(tgLevel) {
    	
};

/**
 * Set the Log Level
 * @param level the log level dynamically
 */
TGLogger.prototype.setLevel = function(tgLevel) {
	this._tgLogLevel = tgLevel;
	this.setlevelImpl(tgLevel);
};

/**
 * Get the Log Level
 * @return level the log level dynamically
 */
TGLogger.prototype.getLevel = function() {
	return this._tgLogLevel;
};

TGLogger.prototype.isInfo = function() {
	return this._tgLogLevel >= TGLogLevel.Info;
};

TGLogger.prototype.isInfo = function() {
	return this._tgLogLevel >= TGLogLevel.Info;
};

TGLogger.prototype.isDebugWire = function() {
	return this._tgLogLevel >= TGLogLevel.DebugWire;
};

TGLogger.prototype.isDebug = function() {
	return this._tgLogLevel >= TGLogLevel.Debug;
};

TGLogger.prototype.isWarning = function() {
	return this._tgLogLevel >= TGLogLevel.Warning;
};

TGLogger.prototype.isInfo = function() {
	return this._tgLogLevel >= TGLogLevel.Info;
};

TGLogger.prototype.isError = function() {
	return this._tgLogLevel >= TGLogLevel.Error;
};

TGLogger.prototype.isFatal = function() {
	return this._tgLogLevel >= TGLogLevel.Fatal;
};

exports.TGLogger = TGLogger;

function test() {
}

//test();