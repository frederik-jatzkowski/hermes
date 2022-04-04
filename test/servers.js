var http = require('http');

for(let i = 2; i < process.argv.length; i++){
    let port = parseInt(process.argv[i])
    //create a server object:
    http.createServer(function (req, res) {
        res.write(`Hello World from Port ${port}!`); //write a response to the client
        res.end(); //end the response
    }).listen(port); //the server object listens on port 8080
}
