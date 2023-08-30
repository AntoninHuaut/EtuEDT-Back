import { ITimetableExtended } from '/src/app.interface.ts';
import TimeTable from '/src/timetables/TimeTable.ts';
import TTCache from '/src/timetables/TTCache.ts';

const TT_CACHE = new TTCache();

export function getTTListByUniv(numUniv: number | undefined): ITimetableExtended[] {
    if (!numUniv) return [];

    return TT_CACHE.getUnivObj_TTList()[numUniv];
}

export function getTTById(numUniv: number | undefined, adeResources: number | undefined): TimeTable | undefined {
    if (!numUniv || !TT_CACHE.getUnivObj_TTObj()[numUniv] || !adeResources) return undefined;

    return TT_CACHE.getUnivObj_TTObj()[numUniv][adeResources];
}
