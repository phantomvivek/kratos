let cluster = require('cluster')
let WebSocketServer = require('uws').Server;
let express = require('express')
let http = require('http')
let connectCount = 0;

if (cluster.isMaster)
{
    let numCPUs = require("os").cpus().length;
	for (let i = 0; i < numCPUs; i++)
    {
        cluster.fork();
    }
    cluster.on('exit', function (worker)
    {
        console.error('Worker DISCONNECTED' + worker.process.pid + ' died. Replacing the died worker...');
        cluster.fork();
    });
}
else
{
    let app = express();
    let server = http.createServer(app);
    let wss = new WebSocketServer({ server: server, clientTracking: false });

    server.listen("8080", function ()
    {
        console.log("Server started")
    });
    wss.on('connection', Connect)
}

function Connect(ws) {
    connectCount++
    console.log("Open", connectCount)
    ws.on('message', (msg) => {
        //console.log("Msg", msg)
    });
    ws.on('close', () => {
        //console.log("Close")
        connectCount--;
    });
    ws.on('error', (err) => {
        console.log("Error", err)
        connectCount--
    });
}