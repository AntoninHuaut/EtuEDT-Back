import config from '/config.ts';
import { ITimetableUniv, IUniv } from '/src/app.interface.ts';

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

export function getAllTT(): ITimetableUniv[] {
    const response: ITimetableUniv[] = [];

    config.univ.forEach((univ, index) => {
        univ.timetable.forEach((tt) => {
            response.push({
                numUniv: index + 1,
                nameUniv: univ.nameUniv,
                adeUniv: univ.adeUniv,
                numYearTT: tt.numYearTT,
                descTT: tt.descTT,
                adeResources: tt.adeResources,
                adeProjectId: tt.adeProjectId,
            });
        });
    });

    return response;
}
