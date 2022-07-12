import config from '/config/config.ts';
import { now } from '/env.ts';

const minMinsRefresh = 15;

export default function () {
    if (!config.refreshMinuts || !Number.isInteger(config.refreshMinuts) || config.refreshMinuts < minMinsRefresh) {
        console.error(now(), `refreshMinuts (${config.refreshMinuts}) should be greater than or equal to ${minMinsRefresh}`);
        Deno.exit(1);
    } else console.log(now(), 'Timer OK');
}
