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

var conFactory = require('../lib/connection/TGConnectionFactory');
    
function createAndCommitTestObjects(conn) {
    
    var gof = conn.getGraphObjectFactory();    
  	console.log("Start transaction 1");
  	console.log("Create node1");
    var node1 = gof.createNode();
    node1.setStringAttribute("name", "john doe");
    node1.setIntegerAttribute("multiple", 7);
    node1.setDoubleAttribute("rate", 3.3);
    node1.setStringAttribute("nickname", "蠢豬");    
    conn.insertEntity(node1);

    console.log("Create node2");
    var node2 = gof.createNode();
    node2.setStringAttribute("name", "julie");
    node2.setDoubleAttribute("factor", 3.3);
    conn.insertEntity(node2);

    console.log("Create node3");
    var node3 = gof.createNode();
    node3.setStringAttribute("name", "margo");
    node3.setDoubleAttribute("factor", 2.3);
    node3.setBooleanAttribute("是超人嗎", false);
    conn.insertEntity(node3);
    
    console.log("Create edge1");
    var edge1 = gof.createBidirectionalEdge(node1, node2);
    edge1.setAttribute("name", "spouse");
    conn.insertEntity(edge1);

    console.log("Create edge2");
    var edge2 = gof.createDirectedEdge(node1, node3);
    edge2.setAttribute("name", "daughter");
    conn.insertEntity(edge2);

    console.log("Commit transaction 1");
    
    conn.commit(
    	function updateAndCommitTestObjects(status) {
    		
    		console.log("Start transaction 2");
    		// updates
    		console.log("Update node1");
    		node1.setIntegerAttribute("age", 40);
    		node1.setAttribute("multiple", 8);
    		conn.updateEntity(node1);

    		//console.log("Delete edge1");
    		//what is going to happen here?
    		//need to do proper cleanup
    		//conn.deleteEntity(edge1);

    		console.log("Update edge2");
    		//add a new attribute
    		edge2.setBooleanAttribute("extra", true);
    		//update existing one
    		edge2.setAttribute("name", "kid");
    		conn.updateEntity(edge2);
    
    		console.log("Create node4");
    		var node4 = gof.createNode();
    		node4.setAttribute("name", "Smith");
    		node4.setDoubleAttribute("level", 3.0);
    	    conn.insertEntity(node4);

    		console.log("Create edge3");
    		var edge3 = gof.createBidirectionalEdge(node1, node4);
    		edge3.setAttribute("name", "Tennis Partner");
    	    conn.insertEntity(edge3);

    		console.log("Commit transaction 2");
    		
    		conn.commit ( /* update John doe's value --------*/
    				function(status) {
    		    		conn.disconnect();
    		    		console.log("Disconnected.");
    				}
    		); 
    	}
    );
}

function test() {
    var connectionFactory = conFactory.getFactory();
    var linkURL = 'tcp://scott@192.168.1.6:8222';
    var conn = connectionFactory.createConnection(linkURL, 'scott', 'scott', null);
    
    conn.connect(function(connectionStatus) {
        if (connectionStatus) {
            console.log('Connection to server successful');
            createAndCommitTestObjects(conn);
        }
    });
}

test();