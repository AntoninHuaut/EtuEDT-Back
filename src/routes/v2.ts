import { Context, Router } from 'oak';
import { handleUnivList, handleTTList, handleTTFormat } from '../controllers/v2API.ts';

export default function getRouter(router: Router, path: string) {
    router
        .get(path, async (context: Context) => await handleUnivList(context))
        .get(path + '/:numUniv', async (context: Context) => await handleTTList(context))
        .get(path + '/:numUniv/:adeResources/:format?', async (context: Context) => await handleTTFormat(context));
}
