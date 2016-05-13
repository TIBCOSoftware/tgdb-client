/**
 * http://usejsdoc.org/
 */
/*
In the node.js intro tutorial (http://nodejs.org/), they show a basic tcp 
server, but for some reason omit a client connecting to it.  I added an 
example at the bottom.

Save the following server in example.js:
*/
var LOCAL_PORT  = 6512;

var net = require('net');

var buffer = new Buffer(0, 'hex');

var server = net.createServer(function(socket) {
	
    socket.on('data', function(data){
    	
        console.log('<< Request from client ', data);

    	buffer = Buffer.concat([buffer, new Buffer(data, 'hex')]);
    	
    });
	
	
	socket.write('13ho server\r\n');
	socket.pipe(socket);
});


server.listen(LOCAL_PORT, '127.0.0.1');

console.log("TCP GDB server accepting connection on port: " + LOCAL_PORT);
