import { getAllUniv } from '/sql/timetable.ts';
import config from '/config/config.ts';

import checkConfig from '/config/checkConfig.ts';
import { now } from '/env.ts';
import { IUniv } from '/model/UnivModel.ts';

export default class UnivCache {
    private cachedUniv: IUniv[];

    constructor() {
        checkConfig();

        setInterval(() => this.refresh(), config.refreshMinuts * 60 * 1000);

        this.cachedUniv = [];
        this.refresh();
    }

    getUnivList(): IUniv[] {
        return this.cachedUniv;
    }

    refresh() {
        console.log(now(), 'Refreshing Univ...');

        getAllUniv()
            .then((univList) => (this.cachedUniv = univList))
            .catch((err) => console.error(err));
    }
}
