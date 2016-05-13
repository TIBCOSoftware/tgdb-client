/**
 * http://usejsdoc.org/
 */
var net = require('net');
 
var LOCAL_PORT  = 8222;
var REMOTE_PORT = 8222;
var REMOTE_ADDR = "192.168.1.18";
var serviceSocket = new net.Socket();

var server = net.createServer(function (socket) {
    socket.on('data', function (msg) {
        console.log('  ** START **');
        console.log('<< From client to proxy ', msg);

        
        serviceSocket.connect(parseInt(REMOTE_PORT), REMOTE_ADDR, function () {
            console.log('>> From proxy to remote', msg);
            serviceSocket.write(msg);
        });
        serviceSocket.on("data", function (data) {
            console.log('<< From remote to proxy', data);
            socket.write(data);
            console.log('>> From proxy to client', data);
        });
    });
});
 
server.listen(LOCAL_PORT);
console.log("TCP proxy server accepting connection on port: " + LOCAL_PORT);