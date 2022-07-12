import TimeTable from '/cache/TimeTable.ts';
import { getAllTT } from '/sql/timetable.ts';
import config from '/config/config.ts';

import checkConfig from '/config/checkConfig.ts';
import { now } from '/env.ts';
import dayjs from 'dayjs';

interface Request {
    status: number;
    body: string;
}

export default class TTCache {
    private cacheRefresh: [];

    private univObj_TTList;
    private univObj_TTObj;
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
        console.log(now, 'Refreshing Timetable...');

        getAllTT()
            .catch((err) => console.error(err))
            .then((ttList) => {
                if (!ttList || !Array.isArray(ttList)) return;

                const rqList = ttList.map((timetable: any) => requestTT(timetable));

                Promise.all(rqList)
                    .then((res: Request[]) => {
                        const bodyConverted: Request[] = res.map((subRes: Request) => convertBodyString(subRes));
                        this.updateTT(ttList, bodyConverted);
                    })
                    .catch((err) => console.error(now(), '[Catch]', err));
            });
    }

    updateTT(ttList: any[], res: Request[]) {
        const cacheRefresh: any = this.init ? this.cacheRefresh : [];
        const date = new Date();

        for (let i = 0; i < res.length; i++) {
            const item = cacheRefresh.find(
                (subItem: { numUniv: number; adeResources: number }) => subItem.numUniv == ttList[i].numUniv && subItem.adeResources == ttList[i].adeResources
            );
            const response = res[i];

            if (item) {
                if (response.status == 200) {
                    item.lastUpdate = date;
                    item.ics = response.body;
                    item.setJSON();
                }
            } else cacheRefresh.push(new TimeTable(ttList[i], date, response.body));
        }

        console.log(now(), 'Refreshed Timetables completed');

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

async function requestTT(timetableSql: any): Promise<Request> {
    const firstDate = dayjs().subtract('4', 'M').format('YYYY-MM-DD');
    const lastDate = dayjs().add('4', 'M').format('YYYY-MM-DD');

    const params = new URLSearchParams({
        resources: timetableSql.adeResources,
        projectId: timetableSql.adeProjectId,
        calType: 'ical',
        firstDate: firstDate,
        lastDate: lastDate,
    });

    const res = await fetch(timetableSql.adeUniv + '?' + params);
    const text = await res.text();
    return {
        status: res.status,
        body: text,
    };
}

function convertBodyString(request: Request): Request {
    if (!request.body) return request;

    const splited = request.body.split('\r');

    for (let i = 0; i < splited.length; i++) {
        splited[i] = convertStringSplit(splited, i, '_s');
        splited[i] = convertStringSplit(splited, i, '_1');
        splited[i] = convertStringSplit(splited, i, '_2');
    }

    request.body = splited.join('');
    return request;
}

function convertStringSplit(splited: { [x: string]: any }, i: number, strSplit: string) {
    if (splited[i].startsWith('\nSUMMARY:') && splited[i].toLowerCase().includes(strSplit))
        return splited[i].substring(0, splited[i].toLowerCase().indexOf(strSplit));

    return splited[i];
}
