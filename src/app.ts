import { Application, Context } from 'oak';

import router from '/routes/api.ts';

const app = new Application();
app.use(async (ctx: Context, next) => {
    try {
        await next();
    } catch (_err) {
        ctx.response.status = 500;
    }
});

app.use((ctx, next) => {
    ctx.response.headers.set('Access-Control-Allow-Origin', '*');
    next();
});

app.use(router.routes());
app.use(router.allowedMethods());

const options = {
    port: 8080,
};

console.log(`App start on http://localhost:${options.port}`);

app.listen(options);
