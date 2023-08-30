import dayjs from 'dayjs';

import config from '/config.ts';
import { ITimetableExtended, ITimeTableUniv } from '/src/app.interface.ts';
import { checkConfig, getAllTT } from '/src/configHelpers.ts';
import TimeTable from '/src/timetables/TimeTable.ts';

type ResponseParsed = {
    bodyText: string;
} & Response;

interface UnivTTList {
    [numUniv: number]: ITimetableExtended[];
}

interface UnivTTObj {
    [numUniv: number]: {
        [adeResources: number]: TimeTable;
    };
}

export default class TTCache {
    private cacheRefresh: TimeTable[];

    private univObj_TTList: UnivTTList;
    private univObj_TTObj: UnivTTObj;

    constructor() {
        checkConfig();

        setInterval(() => this.refresh(), config.refreshMinuts * 60 * 1000);

        this.cacheRefresh = [];
        this.univObj_TTList = {};
        this.univObj_TTObj = {};

        this.refresh();
    }

    getUnivObj_TTList(): UnivTTList {
        return this.univObj_TTList;
    }

    getUnivObj_TTObj(): UnivTTObj {
        return this.univObj_TTObj;
    }

    refresh() {
        console.log(new Date(), 'Refreshing Timetable...');

        const ttList = getAllTT();
        const rqList = ttList.map((timetable) => requestTT(timetable));

        Promise.all(rqList)
            .then((res) => this.updateTT(ttList, res))
            .catch((err) => console.error(new Date(), '[Catch]', err));
    }

    updateTT(ttList: ITimeTableUniv[], res: ResponseParsed[]) {
        res = res.map((subRes) => convertBodyString(subRes));
        const cacheRefreshTmp = this.cacheRefresh;
        const date = new Date();

        for (let i = 0; i < res.length; i++) {
            const item = cacheRefreshTmp.find(
                (subItem: { numUniv: number; adeResources: number }) => subItem.numUniv == ttList[i].numUniv && subItem.adeResources == ttList[i].adeResources
            );
            const response = res[i];

            if (item) {
                if (response.status == 200) {
                    item.lastUpdate = date;
                    item.ics = response.bodyText;
                    item.setJSON();
                }
            } else {
                cacheRefreshTmp.push(new TimeTable(ttList[i], date, response.bodyText));
            }
        }
        this.cacheRefresh = cacheRefreshTmp;
        console.log(new Date(), 'Refreshed Timetables completed');

        const univObj_TTList_Tmp: UnivTTList = {};
        cacheRefreshTmp.forEach((item: TimeTable) => {
            if (!univObj_TTList_Tmp[item.getNumUniv()]) {
                univObj_TTList_Tmp[item.getNumUniv()] = [];
            }

            univObj_TTList_Tmp[item.getNumUniv()].push(item.getAPIData());
        });
        this.univObj_TTList = univObj_TTList_Tmp;

        const univObj_TTObj_Tmp: UnivTTObj = {};
        cacheRefreshTmp.forEach((item: TimeTable) => {
            if (!univObj_TTObj_Tmp[item.getNumUniv()]) {
                univObj_TTObj_Tmp[item.getNumUniv()] = {};
            }

            univObj_TTObj_Tmp[item.getNumUniv()][item.getAdeResources()] = item;
        });
        this.univObj_TTObj = univObj_TTObj_Tmp;
    }
}

async function requestTT(timetableSql: ITimeTableUniv): Promise<ResponseParsed> {
    const firstDate = dayjs().subtract('4', 'M').format('YYYY-MM-DD');
    const lastDate = dayjs().add('4', 'M').format('YYYY-MM-DD');

    const params = new URLSearchParams({
        resources: '' + timetableSql.adeResources,
        projectId: '' + timetableSql.adeProjectId,
        calType: 'ical',
        firstDate: firstDate,
        lastDate: lastDate,
    });

    const res = await fetch(timetableSql.adeUniv + '?' + params);
    const text = await res.text();
    return {
        ...res,
        bodyText: text,
    };
}

function convertBodyString(response: ResponseParsed): ResponseParsed {
    if (!response.bodyText) return response;

    const splited = response.bodyText.split('\r');

    for (let i = 0; i < splited.length; i++) {
        splited[i] = convertStringSplit(splited, i, '_s');
        splited[i] = convertStringSplit(splited, i, '_1');
        splited[i] = convertStringSplit(splited, i, '_2');
    }

    response.bodyText = splited.join('');
    return response;
}

function convertStringSplit(splited: string[], i: number, strSplit: string) {
    if (splited[i].startsWith('\nSUMMARY:') && splited[i].toLowerCase().includes(strSplit)) return splited[i].substring(0, splited[i].toLowerCase().indexOf(strSplit));

    return splited[i];
}
