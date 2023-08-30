import { Context, Router } from 'oak';

import { handleTTFormat, handleTTList, handleUnivList } from '/controllers/v2API.ts';

export default function getRouter(router: Router, path: string) {
    router
        .get(path, (context: Context) => handleUnivList(context))
        .get(path + '/:numUniv', (context: Context) => handleTTList(context))
        .get(path + '/:numUniv/:adeResources/:format?', (context: Context) => handleTTFormat(context));
}
