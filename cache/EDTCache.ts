import { config, moment, now } from '../dps.ts';
import EDT from './EDT.ts';
import { getAllEDT } from '../sql/edt.ts';

import checkConfig from '../config/checkConfig.ts';

export default class EDTCache {

    private cached: [];

    private cacheAllEDT: any;
    private cachedEDT: any;
    private init: boolean;

    constructor() {
        checkConfig();

        setInterval(() => this.refresh(), config.refreshMinuts * 60 * 1000);

        this.cached = [];
        this.cacheAllEDT = [];
        this.cachedEDT = {};
        this.init = false;

        this.refresh();
    }

    getCacheAllEDT(): any {
        return this.cacheAllEDT;
    }

    getCacheEDT(): any {
        return this.cachedEDT;
    }

    refresh() {
        console.log(now(), "Refresh EDT...");

        getAllEDT()
            .catch(err => console.error(err))
            .then(edtList => {
                if (!edtList || !Array.isArray(edtList)) return;

                const rqList = edtList.map((edt: any) => requestEdt(edt));

                Promise.all(rqList)
                    .then((res: any) => {
                        res = res.map((subRes: any) => convertString(subRes));
                        this.updateEdt(edtList, res);
                    })
                    .catch(err => console.error(now(), "[Catch]", err));
            });
    }

    updateEdt(edtList: any[], res: string | any[]) {
        let cacheRefresh: any = this.init ? this.cached : [];
        const date = new Date();
        const ignoreOldEntries: any[] = [];

        for (let i = 0; i < res.length; i++) {
            ignoreOldEntries.push(edtList[i].resources);

            const item = cacheRefresh.find((subItem: { edtId: any; }) => subItem.edtId == edtList[i].resources);

            if (!!item) {
                if (!res[i].includes('HTTP ERROR')) {
                    item.lastUpdate = date;
                    item.edtIcs = res[i];
                    item.setJSON();
                }
            } else
                cacheRefresh.push(new EDT(edtList[i], date, res[i]));
        }

        console.log(now(), "Refresh EDT Terminé");

        cacheRefresh = cacheRefresh.filter((item: { edtId: any; }) => ignoreOldEntries.includes(item.edtId));

        this.cached = cacheRefresh;
        this.cacheAllEDT = cacheRefresh.map((item: EDT) => item.getAPIData());
        this.cached.forEach((item: EDT) => this.cachedEDT[item.getEdtId()] = item);

        this.init = true;
    }
}

function requestEdt(edtSql: any): Promise<string> {
    const startYear = moment().startOf('year');

    if (moment().isAfter(startYear.clone().add('6', 'M')))
        startYear.add('6', 'M');

    const firstDate = startYear.format('YYYY-MM-DD');
    const lastDate = startYear.add('6', 'M').format('YYYY-MM-DD');

    const params = new URLSearchParams({
        resources: edtSql.resources,
        projectId: edtSql.projectId,
        calType: 'ical',
        firstDate: firstDate,
        lastDate: lastDate
    });

    return fetch(edtSql.adeEta + '?' + params).then(res => res.text());
}

function convertString(str: string): string {
    if (!str)
        return str;

    const splited = str.split('\r');

    for (let i = 0; i < splited.length; i++) {
        splited[i] = convertStringSplit(splited, i, '_s');
        splited[i] = convertStringSplit(splited, i, '_1');
        splited[i] = convertStringSplit(splited, i, '_2');
    }

    str = splited.join('');

    return str;
}

function convertStringSplit(splited: { [x: string]: any; }, i: number, strSplit: string) {
    if (splited[i].startsWith('\nSUMMARY:') && splited[i].toLowerCase().includes(strSplit))
        return splited[i].substring(0, splited[i].toLowerCase().indexOf(strSplit));

    return splited[i];
}