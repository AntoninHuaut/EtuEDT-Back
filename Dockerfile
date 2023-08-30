FROM denoland/deno:alpine-1.36.3

EXPOSE 8080

WORKDIR /app

ADD . .

RUN deno cache ./src/app.ts

CMD [ "task", "start" ]