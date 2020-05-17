import { oak, config } from "./dps.ts";
import router from "./routes/api.ts";
import { initTable } from "./sql/index.ts";

await initTable();

const app = new oak.Application();
app.use(router.routes());
app.use(router.allowedMethods());

console.log(`web start on http://localhost:${config.port}`);
await app.listen({ port: config.port });
