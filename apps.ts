import { oak, config, env } from "./dps.ts";
import router from "./routes/api.ts";
import { initTable } from "./sql/index.ts";
import { critical } from "https://deno.land/x/std/log/mod.ts";
await initTable();

const app = new oak.Application();
app.use(async (ctx, next) => {
    try {
        await next();
    } catch (err) {
        ctx.response.status = 500;
    }
});

app.use((ctx, next) => {
    ctx.response.headers.set("Access-Control-Allow-Origin", "*");
    next();
});

app.use(router.routes());
app.use(router.allowedMethods());

const options = {
    port: config.port,
    secure: false,
    certFile: "",
    keyFile: ""
};

if (env["TYPE"] && env["TYPE"] == 'production') {
    options.secure = true;
    options.certFile = env["SSLPATH"] + "/fullchain.pem";
    options.keyFile = env["SSLPATH"] + "/privkey.pem";
}

console.log(`web start on http://localhost:${options.port}`);
await app.listen(options);