// static file server for MessDFS

const express = require('express')
const app = express()
const listenPort = 3000

app.use(express.static("public"));
app.use("/index", express.static('index.html'))
app.use("/login", express.static('login.html'))
app.use("/signup", express.static('signup.html'))
app.use("/signout", express.static('signout.html'))

app.listen(listenPort, function() {
    console.log("Server Listening...")
})