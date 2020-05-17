import EDTCache from './EDTCache.ts';
import EDT from './EDT.ts';
const cache = new EDTCache();

export function getCacheEDT(): any {
    return cache.getCacheAllEDT();
}

export function getCacheEDTById(edtId: string | undefined): EDT | undefined {
    if (!edtId) return undefined;

    return cache.getCacheEDT()[edtId];
}