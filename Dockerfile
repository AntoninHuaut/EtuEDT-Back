FROM denoland/deno:alpine-1.25.1

EXPOSE 8080

WORKDIR /app

ADD . .

RUN deno cache --unstable --import-map=./import_map.json ./src/app.ts

CMD [ "task", "start" ]