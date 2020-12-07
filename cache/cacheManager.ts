import TTCache from './TTCache.ts';
import UnivCache from './UnivCache.ts';
import Timetable from './TimeTable.ts';

const TT_CACHE = new TTCache();
const UNIV_CACHE = new UnivCache();

export function refreshTT() {
    TT_CACHE.refresh();
}

export function getUnivList(): any {
    return UNIV_CACHE.getUnivList();
}

export function getTTListByUniv(numUniv: string | undefined): any {
    if (!numUniv) return undefined;

    return TT_CACHE.getUnivObj_TTList()[numUniv];
}

export function getTTById(numUniv: string | undefined, adeResources: string | undefined): Timetable | undefined {
    if (!numUniv || !TT_CACHE.getUnivObj_TTObj()[numUniv] || !adeResources) return undefined;

    return TT_CACHE.getUnivObj_TTObj()[numUniv][adeResources];
}