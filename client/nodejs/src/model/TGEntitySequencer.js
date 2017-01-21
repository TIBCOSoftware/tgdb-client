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

function TGEntitySequencer () {
  var _number = 0;
  function TGEntitySequencerImpl() {
	  this._instance = null;
	  this._number = 0;
      throw new Error('Singleton! Use TGEntitySequencer.getInstrace()');
  }

  // The private constructor
  var TGEntitySequencer = function () {
  };
  
  TGEntitySequencer.getNumber = function () {
		
  };
  
  TGEntitySequencer.prototype = TGEntitySequencerImpl.prototype;

  // The preferred smart constructor
  TGEntitySequencerImpl.getInstance = function () {
	  //if(!this._instance) {
	//	  this._instance = new TGEntitySequencer();
	  //}
	  return this;//._instance;
  };
  
  TGEntitySequencerImpl.getNumber = function () {
	  return _number++;
  };
  
  return TGEntitySequencerImpl;
}



exports.TGEntitySequencer = TGEntitySequencer;