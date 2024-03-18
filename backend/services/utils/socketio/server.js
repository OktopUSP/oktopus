const dotenv = require('dotenv')
dotenv.config();
dotenv.config({ path: `.env.local`, override: true });
const express = require('express');
const app = express();
const PORT = 5000;

const http = require('http').Server(app);
const cors = require('cors');
var allowedOrigins;
let allowedOriginsFromEnv = process.env.CORS_ALLOWED_ORIGINS.split(',')
if (allowedOriginsFromEnv.length > 1) {
  allowedOrigins = allowedOriginsFromEnv
}else{
  allowedOrigins = "*"
}
console.log("allowedOrigins:",allowedOrigins)

const io = require('socket.io')(http, {
    cors: {
        origin: allowedOrigins
    }
});

app.use(cors());

let users = []

io.on('connection', (socket) => {
    console.log(`🚀: ${socket.id} user just connected!`);

    socket.on("callUser", ({ userToCall, signalData, from }) => {
      console.log("user to call:",userToCall)
      let index = users.findIndex(x =>{ 
        return x.name === userToCall
      })
      console.log(index)
      if (index >= 0){
        console.log("calling user named "+ users[index].name+" and id "+users[index].id)
        io.to(users[index].id).emit("callUser", { signal: signalData, from });
      }else{
        console.log("There is no user named "+userToCall+" or he/she is offline")
      }
    });

    socket.on("answerCall", (data) => {
      io.to(data.to).emit("callAccepted", data.signal);
    });

    socket.on("newuser", (data) => {
      let index = users.findIndex(x =>{ x.name === data.name})
      if (index >=0){
        console.log("user already exists, but got connected with other id")
      }else{
        users.push(data)
      }
      console.log(data)
      console.log("total users: ", users)
      io.emit('users', users)
    })

    socket.on('disconnect', () => {
      console.log('🔥: A user disconnected');
      let index = users.findIndex(x => x.id === socket.id)
      if (index >= 0){
        let deletedEl = users.splice(index, 1)
        console.log("users deleted", deletedEl)
        console.log("users after disconection: ", users)
        io.emit('users', users)
      }else{
        console.log("couldn't find user with socket id:", socket.id)
      }
    });
});

http.listen(PORT, () => {
  console.log(`Server listening on ${PORT}`);
});