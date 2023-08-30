import { Context, helpers, httpErrors } from 'oak';
import * as cacheManager from '/cache/cacheManager.ts';
import Timetable from '/cache/TimeTable.ts';

export function handleUnivList(ctx: Context) {
    ctx.response.body = cacheManager.getUnivList();
}

export function handleTTList(ctx: Context) {
    const { numUniv } = helpers.getQuery(ctx, { mergeParams: true });
    if (isNaN(+numUniv)) {
        return new httpErrors.BadRequest('Invalid query parameter');
    }

    const univ_TTList = cacheManager.getTTListByUniv(+numUniv);
    if (univ_TTList.length === 0) {
        ctx.response.body = [];
        ctx.response.status = 404;
    } else {
        ctx.response.body = univ_TTList;
    }
}

export function handleTTFormat(ctx: Context) {
    const { adeResources, numUniv, format } = helpers.getQuery(ctx, { mergeParams: true });
    if (isNaN(+adeResources) || isNaN(+numUniv)) {
        return new httpErrors.BadRequest('Invalid query parameter');
    }

    const dataTT: Timetable | undefined = cacheManager.getTTById(+numUniv, +adeResources);

    if (!dataTT) {
        ctx.response.body = {};
        ctx.response.status = 404;
    } else if (!format) {
        ctx.response.body = dataTT.getAPIData();
    } else {
        switch (format) {
            case 'json':
                ctx.response.body = dataTT.getJSON();
                break;
            case 'ics':
                ctx.response.body = dataTT.getICS();
                ctx.response.headers.set('Content-Type', 'text/calendar');
                break;
            default:
                ctx.response.body = {};
                ctx.response.status = 404;
                break;
        }
    }
}
