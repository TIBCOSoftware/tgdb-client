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

var conFactory               = require('../connection/TGConnectionFactory'),
    ProtocolDataInputStream  = require('../pdu/impl/ProtocolDataInputStream').ProtocolDataInputStream,
    StringUtils              = require('../utils/StringUtils').StringUtils,
    TGEdge                   = require('../model/TGEdge'),
    TGAttributeType          = require('../model/TGAttributeType').TGAttributeType;;

function test() {
    var connectionFactory = new conFactory.DefaultConnectionFactory();
    var linkURL = 'tcp://scott@192.168.1.18:8222';
    var conn = connectionFactory.createConnection(linkURL, 'scott', 'scott', null);

    var callback = function(connectionStatus) {
        if (connectionStatus) {
            console.log('Connection to server successful');
            createAndCommitTestObjects(conn);
        }
    };
    
    conn.connect(callback);
    	
}
    
function createAndCommitTestObjects(conn) {
    
    var gof = conn.getGraphObjectFactory();

    /*
    System.out.println("Start test query");
    TGQuery Query1 = conn.createQuery("testquery < X '5ef';");
    Query1.execute();
    conn.executeQuery("testquery < X '5ef';");
    */
    
  	console.log("Start transaction 1");
  	console.log("Create node1");
    var node1 = gof.createNode();
    node1.setAttribute("name", "john doe", TGAttributeType.STRING);
    node1.setAttribute("multiple", 7, TGAttributeType.INT);
    node1.setAttribute("rate", 3.3, TGAttributeType.DOUBLE);
    node1.setAttribute("nickname", "蠢豬", TGAttributeType.STRING);    
    
    console.log("Create node2");
    var node2 = gof.createNode();
    node2.setAttribute("name", "julie", TGAttributeType.STRING);
    node2.setAttribute("factor", 3.3, TGAttributeType.DOUBLE);

    console.log("Create node3");
    var node3 = gof.createNode();
    node3.setAttribute("name", "margo", TGAttributeType.STRING);
    node3.setAttribute("factor", 2.3, TGAttributeType.DOUBLE);
    node3.setAttribute("是超人嗎", false, TGAttributeType.BOOLEAN);
    
    console.log("Create edge1");
  	//FIXME: addEdge does not trigger change notification
    var edge1 = gof.createEdge(node1, node2, TGEdge.DirectionType.BIDIRECTIONAL);
    edge1.setAttribute("name", "spouse");

    console.log("Create edge2");
    var edge2 = gof.createEdge(node1, node3, TGEdge.DirectionType.DIRECTED);
    edge2.setAttribute("name", "daughter");

    console.log("Commit transaction 1");
    
    conn.commit(
    	function updateAndCommitTestObjects(status) {
    		
    		console.log("Start transaction 2");
    		// updates
    		console.log("Update node1");
    		node1.setAttribute("age", 40, TGAttributeType.INT);
    		node1.setAttribute("multiple", 8);
    		conn.updateEntity(node1);
    
    		console.log("Delete edge1");
    		//what is going to happen here?
    		//need to do proper cleanup
    		conn.deleteEntity(edge1);
    
    		console.log("Update edge2");
    		//add a new attribute
    		edge2.setAttribute("extra", true, TGAttributeType.BOOLEAN);
    		//update existing one
    		edge2.setAttribute("name", "kid");
    		conn.updateEntity(edge2);
    
    		console.log("Create node4");
    		var node4 = gof.createNode();
    		node4.setAttribute("name", "Smith");
    		node4.setAttribute("level", 3.0, TGAttributeType.DOUBLE);

    		console.log("Create edge3");
    		var edge3 = gof.createEdge(node1, node4, TGEdge.DirectionType.BIDIRECTIONAL);
    		edge3.setAttribute("name", "Tennis Partner");

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

test();