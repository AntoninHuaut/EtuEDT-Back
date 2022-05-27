import { oak } from "../../deps.ts";
import v2router from "./v2.ts";

const router = new oak.Router();

router.get('/', async (context) => {
    context.response.body = {
        "path": "v2"
    }
});

v2router(router, "/v2");

export default router;