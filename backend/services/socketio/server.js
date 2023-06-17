const express = require('express');
const app = express();
const PORT = 5000;

const http = require('http').Server(app);
const cors = require('cors');

const io = require('socket.io')(http, {
    cors: {
        origin: "http://localhost:3000"
    }
});

app.use(cors());

let users = []

io.on('connection', (socket) => {
    console.log(`🚀: ${socket.id} user just connected!`);

    socket.on("newuser", (data) => {
      users.push(data)
      console.log(data)
    })

    socket.on("getusers", () => {
      socket.broadcast.emit('users', )
    });

    socket.on('disconnect', () => {
      console.log('🔥: A user disconnected');
    });
});

http.listen(PORT, () => {
  console.log(`Server listening on ${PORT}`);
});