import { oak } from "../dps.ts";
import { handleOldAPI } from "../controllers/oldAPI.ts";
import v2router from "./v2.ts";

const router = new oak.Router();

router.get('/', async (context) => {
    context.response.body = {
        "path": "v2"
    }
});

v2router(router, "/v2");

// old
router.get('/:adeResources/ics', async (context) => handleOldAPI(context));

export default router;