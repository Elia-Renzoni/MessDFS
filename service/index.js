// static server

const express = require('express')
const app = express()
const listenPort = 3000

app.get("/", function (req, res) {
    res.send("Hello World!")
});

app.listen(listenPort, function() {
    console.log("Server Listening...")
})