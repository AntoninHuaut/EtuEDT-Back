import { config, moment, now } from '../dps.ts';
import TimeTable from './TimeTable.ts';
import { getAllTT } from '../sql/timetable.ts';

import checkConfig from '../config/checkConfig.ts';

export default class TTCache {

    private cacheRefresh: [];

    private univObj_TTList: any;
    private univObj_TTObj: any;
    private init: boolean;

    constructor() {
        checkConfig();

        setInterval(() => this.refresh(), config.refreshMinuts * 60 * 1000);

        this.cacheRefresh = [];
        this.univObj_TTList = {};
        this.univObj_TTObj = {};
        this.init = false;

        this.refresh();
    }

    getUnivObj_TTList(): any {
        return this.univObj_TTList;
    }

    getUnivObj_TTObj(): any {
        return this.univObj_TTObj;
    }

    refresh() {
        console.log(now(), "Refreshing Timetable...");

        getAllTT()
            .catch(err => console.error(err))
            .then(ttList => {
                if (!ttList || !Array.isArray(ttList)) return;

                const rqList = ttList.map((timetable: any) => requestTT(timetable));

                Promise.all(rqList)
                    .then((res: any) => {
                        res = res.map((subRes: any) => convertString(subRes));
                        this.updateTT(ttList, res);
                    })
                    .catch(err => console.error(now(), "[Catch]", err));
            });
    }

    updateTT(ttList: any[], res: string | any[]) {
        let cacheRefresh: any = this.init ? this.cacheRefresh : [];
        const date = new Date();

        for (let i = 0; i < res.length; i++) {
            const item = cacheRefresh.find((subItem: { numUniv: any, adeResources: any }) => subItem.numUniv == ttList[i].numUniv && subItem.adeResources == ttList[i].adeResources);

            if (!!item) {
                if (!res[i].includes('HTTP ERROR')) {
                    item.lastUpdate = date;
                    item.ics = res[i];
                    item.setJSON();
                }
            } else
                cacheRefresh.push(new TimeTable(ttList[i], date, res[i]));
        }

        console.log(now(), "Refreshed Timetables completed");

        this.cacheRefresh = cacheRefresh;
        const tmp_UnivObj_TTList: any = {};
        const tmp_UnivObj_TTObj: any = {};

        cacheRefresh.forEach((item: TimeTable) => {
            if (!tmp_UnivObj_TTList[item.getNumUniv()]) {
                tmp_UnivObj_TTList[item.getNumUniv()] = [];
            }

            tmp_UnivObj_TTList[item.getNumUniv()].push(item.getAPIData());
        });
        this.univObj_TTList = tmp_UnivObj_TTList;

        this.cacheRefresh.forEach((item: TimeTable) => {
            if (!tmp_UnivObj_TTObj[item.getNumUniv()]) {
                tmp_UnivObj_TTObj[item.getNumUniv()] = {};
            }

            tmp_UnivObj_TTObj[item.getNumUniv()][item.getAdeResources()] = item;
        });
        this.univObj_TTObj = tmp_UnivObj_TTObj;

        this.init = true;
    }
}

function requestTT(timetableSql: any): Promise<string> {
    const firstDate = moment().subtract('4', 'M').format('YYYY-MM-DD');
    const lastDate = moment().add('4', 'M').format('YYYY-MM-DD');

    const params = new URLSearchParams({
        resources: timetableSql.adeResources,
        projectId: timetableSql.adeProjectId,
        calType: 'ical',
        firstDate: firstDate,
        lastDate: lastDate
    });

    return fetch(timetableSql.adeUniv + '?' + params).then(res => res.text());
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