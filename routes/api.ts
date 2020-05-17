import { oak } from "../dps.ts";
import { handleEDTList, handleEDTDetails } from "../controllers/api.ts";

const router = new oak.Router();

router
    .get('/', async (context) => await handleEDTList(context))
    .get('/:edtId/:format?', async (context) => await handleEDTDetails(context));

export default router;