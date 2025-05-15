FROM node:18.18.0-alpine as builder

WORKDIR /app

COPY ../ ./
RUN rm .env && rm .env.local || true

RUN npm install 

RUN npm run build

FROM node:18.18.0-alpine as runner

WORKDIR /app

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/package-lock.json ./package-lock.json
COPY --from=builder /app/public ./public

RUN npm install

CMD [ "npm", "run", "start" ]
