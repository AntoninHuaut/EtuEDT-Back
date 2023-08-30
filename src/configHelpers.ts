import config from '/config.ts';
import { ITimeTableUniv, IUniv } from '/src/app.interface.ts';

const minMinsRefresh = 15;

export function checkConfig() {
    if (!config.refreshMinuts || !Number.isInteger(config.refreshMinuts) || config.refreshMinuts < minMinsRefresh) {
        console.error(new Date(), `refreshMinuts (${config.refreshMinuts}) should be greater than or equal to ${minMinsRefresh}`);
        Deno.exit(1);
    } else console.log(new Date(), 'Timer OK');
}

export function getAllUniv(): IUniv[] {
    return config.univ.map((univ, index) => ({
        numUniv: index + 1,
        nameUniv: univ.nameUniv,
        adeUniv: univ.adeUniv,
    }));
}

export function getAllTT(): ITimeTableUniv[] {
    const response: ITimeTableUniv[] = [];

    config.univ.forEach((univ, index) => {
        univ.timetable.forEach((tt) => {
            response.push({
                numUniv: index + 1,
                nameUniv: univ.nameUniv,
                adeUniv: univ.adeUniv,
                adeResources: tt.adeResources,
                numYearTT: tt.numYearTT,
                descTT: tt.descTT,
                adeProjectId: tt.adeProjectId,
            });
        });
    });

    return response;
}
