FROM node:16.20.2-alpine as builder

WORKDIR /app

COPY ./ ./

RUN npm install 

RUN NEXT_PUBLIC_REST_ENDPOINT=REST_API_URL NEXT_PUBLIC_WS_ENDPOINT=WS_URL npm run build

RUN ls -la && echo "Listing directory contents done"

FROM node:16.20.2-alpine as runner

WORKDIR /app

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/build/entrypoint.sh ./entrypoint.sh
COPY --from=builder /app/public ./public


RUN npm install

RUN chmod 755 entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]

CMD [ "npm", "run", "start" ]
