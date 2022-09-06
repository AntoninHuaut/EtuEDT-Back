import TTCache from '/cache/TTCache.ts';
import UnivCache from '/cache/UnivCache.ts';
import { IUniv } from '/model/UnivModel.ts';
import { ITimetableExtended } from '/model/TimeTableModel.ts';
import TimeTable from '/cache/TimeTable.ts';

const TT_CACHE = new TTCache();
const UNIV_CACHE = new UnivCache();

const isValidNumUniv = (numUniv: number | undefined): number | undefined => {
    if (!numUniv || !(numUniv in TT_CACHE.getUnivObj_TTList())) return;

    return numUniv;
};

const isValidAdeResources = (numUniv: number | undefined, adeResources: number | undefined): { numUniv: number; adeResources: number } | undefined => {
    numUniv = isValidNumUniv(numUniv);
    if (!numUniv || !adeResources || !(adeResources in TT_CACHE.getUnivObj_TTObj()[numUniv])) return;

    return { numUniv, adeResources };
};

export function refreshTT() {
    TT_CACHE.refresh();
}

export function getUnivList(): IUniv[] {
    return UNIV_CACHE.getUnivList();
}

export function getTTListByUniv(numUniv: number | undefined): ITimetableExtended[] {
    numUniv = isValidNumUniv(numUniv);
    if (!numUniv) return [];

    return TT_CACHE.getUnivObj_TTList()[numUniv];
}

export function getTTById(numUniv: number | undefined, adeResources: number | undefined): TimeTable | undefined {
    const res = isValidAdeResources(numUniv, adeResources);
    if (!res) return undefined;

    return TT_CACHE.getUnivObj_TTObj()[res.numUniv][res.adeResources];
}
