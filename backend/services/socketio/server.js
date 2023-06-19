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
    const sessionId = socket.id

    socket.on("newuser", (data) => {
      users.push(data)
      console.log(data)
      console.log("total users: ", users)
      io.emit('users', users)
    })

    socket.on('disconnect', () => {
      console.log('🔥: A user disconnected');
      users.splice(users.findIndex(x => x.id === sessionId), 1);
      console.log("users after disconection: ", users)
      io.emit('users', users)
    });
});

http.listen(PORT, () => {
  console.log(`Server listening on ${PORT}`);
});