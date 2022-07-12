import { moment } from 'https://deno.land/x/deno_moment@v1.1.2/moment.ts';
import { BufReader } from 'https://deno.land/std@0.120.0/io/buffer.ts';
import 'https://deno.land/x/dotenv/load.ts';

function now() {
    return moment().format('DD/MM/YYYY HH:mm');
}

export { BufReader, moment, now };
