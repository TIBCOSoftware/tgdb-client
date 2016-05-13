/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not
 * use this file except in compliance with the License. A copy of the License is
 * included in the distribution package with this file. You also may obtain a
 * copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

exports.PrintUtility = {
	printEntityArray : function(entities) {
		console.log('{{{{{{{{{{{{{{{{{{{{ printEntityArray start }}}}}}}}}}}}}}}}}}}}}}'
						+ Object.keys(entities).length);
		for ( var key in entities) {
			console.log(type + ' - Entity : key = ' + key + ', vid = '
					+ entities[key].getVirtualId());
			var attributes = entities[key].getAttributes();
			for ( var index in attributes) {
				console.log('                  Attrinute : ' + index + ' - '
						+ attributes[index].getValue() + ' - ' /*
																 * +
																 * attributes[index].getAttributeDescriptor().getName()
																 */);
			}
		}
		console.log('{{{{{{{{{{{{{{{{{{{{ printEntityMap end  }}}}}}}}}}}}}}}}}}}}}}');
	},
	printEntityMap : function(entities, type) {
		console.log('{{{{{{{{{{{{{{{{{{{{ printEntityMap start }}}}}}}}}}}}}}}}}}}}}}'
						+ Object.keys(entities).length);
		for ( var key in entities) {
			console.log(type + ' - Entity : key = ' + key + ', vid = '
					+ entities[key].getVirtualId());
			var attributes = entities[key].getAttributes();
			for ( var index in attributes) {
				console.log('                  Attrinute : ' + index + ' - '
						+ attributes[index].getValue() + ' - ' /*
																 * +
																 * attributes[index].getAttributeDescriptor().getName()
																 */);
			}
		}
		console.log('{{{{{{{{{{{{{{{{{{{{ printEntityMap end  }}}}}}}}}}}}}}}}}}}}}}');
	}
}