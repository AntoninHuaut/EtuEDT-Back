import TimeTable from '/cache/TimeTable.ts';
import TTCache from '/cache/TTCache.ts';
import { ITimetableExtended } from '/model/TimeTableModel.ts';

const TT_CACHE = new TTCache();

export function refreshTT() {
    TT_CACHE.refresh();
}

export function getTTListByUniv(numUniv: number | undefined): ITimetableExtended[] {
    if (!numUniv) return [];

    return TT_CACHE.getUnivObj_TTList()[numUniv];
}

export function getTTById(numUniv: number | undefined, adeResources: number | undefined): TimeTable | undefined {
    if (!numUniv || !TT_CACHE.getUnivObj_TTObj()[numUniv] || !adeResources) return undefined;

    return TT_CACHE.getUnivObj_TTObj()[numUniv][adeResources];
}
