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
var TGLogger      = require('./TGLogger').TGLogger;
var WinstonLogger = require('./WinstonLogger').WinstonLogger;
var TGLogLevel    = require('./TGLogger').TGLogLevel;

var TGLogManager = {
	logger : new WinstonLogger(),
	getLogger : function() {
		return this.logger;
	}
};

module.exports = TGLogManager;

function test() {
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.Info);
	logger.logInfo( 'Hello %s %s !!', 'Steven', 'Yang');
	logger.logError( 'Hello %s %s !!', 'Steven', 'Yang');
}

//test();
