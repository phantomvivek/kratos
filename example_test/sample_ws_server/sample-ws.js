let WebSocketServer = require('ws').Server;
let express = require('express')
let http = require('http')
let connectCount = 0;
let totalConnects = 1;

let app = express();
let server = http.createServer(app);
let wss = new WebSocketServer({ server: server, clientTracking: false });

server.listen("8080", function ()
{
    console.log("Server started")
});
wss.on('connection', Connect)

function Connect(ws) {
    let socketId = totalConnects;
    connectCount++;
    totalConnects++;
    console.log(`\nConnection opened: ${socketId}`)

    ws.on('message', (msg) => {
        console.log(`Message on socket ${socketId}: ${msg}`)
    });
    ws.on('close', () => {
        connectCount--;
    });
    ws.on('error', (err) => {
        console.log("Error", err)
        connectCount--
    });
}
