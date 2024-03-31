FROM node:16.20.2-alpine

WORKDIR /app

COPY ../ .

RUN npm install 

RUN npm run build

ENTRYPOINT [ "npm", "run", "start" ]