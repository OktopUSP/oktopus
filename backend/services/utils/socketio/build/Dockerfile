FROM node:16.20.2-alpine

WORKDIR /app

COPY ../package*.json ./

RUN npm install

COPY ../ .

CMD ["node", "server.js"]