import { oak } from "../../deps.ts";
import { handleUnivList, handleTTList, handleTTFormat } from "../controllers/v2API.ts";

export default function getRouter(router: oak.Router, path: string) {
    router
        .get(path, async (context) => await handleUnivList(context))
        .get(path + '/:numUniv', async (context) => await handleTTList(context))
        .get(path + '/:numUniv/:adeResources/:format?', async (context) => await handleTTFormat(context));
};