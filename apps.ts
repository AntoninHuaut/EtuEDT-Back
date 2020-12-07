import { oak, config } from "./dps.ts";
import router from "./routes/api.ts";
import { initTable } from "./sql/index.ts";
import cli from "./cli.ts";
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
    port: config.port
}

console.log(`web start on http://localhost:${options.port}`);

app.listen(options);
cli();