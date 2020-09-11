import { config, now } from '../dps.ts';
import { getAllUniv } from '../sql/timetable.ts';

import checkConfig from '../config/checkConfig.ts';

export default class UnivCache {

    private cachedUniv: any;

    constructor() {
        checkConfig();

        setInterval(() => this.refresh(), config.refreshMinuts * 60 * 1000);

        this.cachedUniv = [];
        this.refresh();
    }

    getUnivList(): any {
        return this.cachedUniv;
    }

    refresh() {
        console.log(now(), "Refreshing Univ...");

        getAllUniv()
            .catch(err => console.error(err))
            .then(univList => this.cachedUniv = univList);
    }
}