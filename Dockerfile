FROM denoland/deno:alpine-1.23.3

EXPOSE 8080

WORKDIR /app

ADD . .

RUN deno cache --import-map=./import_map.json ./src/app.ts

CMD [ "task", "start" ]