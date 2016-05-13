/**
 * Created by aathalye on 1/9/15.
 */

//var ch = require('../channel/test/TESTChannel');
var conFactory = require('../connection/TGConnectionFactory');
var TGEdge = require('../model/TGEdge'),
    ProtocolDataInputStream  = require('../pdu/impl/ProtocolDataInputStream').ProtocolDataInputStream,
    StringUtils = require('../utils/StringUtils').StringUtils;

function test() {
    var connectionFactory = new conFactory.DefaultConnectionFactory();
    var linkURL = 'test://scott@192.168.1.18:8222';
    var connection = connectionFactory.createConnection(linkURL, 'scott', 'scott', null);
    var callback = function(connectionStatus) {
        if (connectionStatus) {
            console.log('Connection to server successful');
            createAndCommitTestObjects(connection);
        }
    };
    connection.connect(callback);
}

function createAndCommitTestObjects(connection) {
    var graphObjectFactory = connection.getGraphObjectFactory();
    var node1 = graphObjectFactory.createNode();

    node1.setAttribute('name', 'john doe');
    node1.setAttribute('multiple', 7);
    node1.setAttribute('rate', 3.3);

    //var node2 = graphObjectFactory.createNode();
    //node2.setAttribute("name", "julie");
    //node2.setAttribute("factor", "3.3");
    //
    //var edge1 = graphObjectFactory.createEdge(node1, node2, TGEdge.DirectionType.BIDIRECTIONAL);
    //edge1.setAttribute('name', 'spouse');
    //console.log('Creating edge between John and julie');

    connection.commit();
}
test();