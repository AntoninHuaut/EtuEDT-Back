import TTCache from '/cache/TTCache.ts';
import UnivCache from '/cache/UnivCache.ts';
import { IUniv } from '/model/UnivModel.ts';
import { ITimetableExtended } from '/model/TimeTableModel.ts';
import TimeTable from '/cache/TimeTable.ts';

const TT_CACHE = new TTCache();
const UNIV_CACHE = new UnivCache();

export function refreshTT() {
    TT_CACHE.refresh();
}

export function getUnivList(): IUniv[] {
    return UNIV_CACHE.getUnivList();
}

export function getTTListByUniv(numUniv: number | undefined): ITimetableExtended[] {
    if (!numUniv) return [];

    return TT_CACHE.getUnivObj_TTList()[numUniv];
}

export function getTTById(numUniv: number | undefined, adeResources: number | undefined): TimeTable | undefined {
    if (!numUniv || !TT_CACHE.getUnivObj_TTObj()[numUniv] || !adeResources) return undefined;

    return TT_CACHE.getUnivObj_TTObj()[numUniv][adeResources];
}
