// static file server for MessDFS

const express = require('express')
const app = express()
const listenPort = 3000

app.use("/index", express.static('public/index.html'))
app.use("/login", express.static('public/login.html'))
app.use("/signup", express.static('public/signup.html'))
app.use("/signout", express.static('public/signout.html'))

app.listen(listenPort, function() {
    console.log("Server Listening...")
})