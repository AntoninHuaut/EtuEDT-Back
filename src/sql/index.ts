import postgres from 'postgres';
import { get } from '/env.ts';

const DB_HOST = get('POSTGRES_HOST');
const DB_PORT = get('POSTGRES_PORT') ?? '';
const DB_DATABASE = get('POSTGRES_DB');
const DB_USER = get('POSTGRES_USER');
const DB_PASSWORD = get('POSTGRES_PASSWORD');

if (!DB_HOST || isNaN(+DB_PORT) || !DB_DATABASE || !DB_USER || !DB_PASSWORD) {
    console.error('Invalid DB configuration');
    Deno.exit(2);
}

const sql = postgres({
    host: DB_HOST,
    port: +DB_PORT,
    database: DB_DATABASE,
    user: DB_USER,
    password: DB_PASSWORD,
    connection: {
        timezone: 'UTC',
    },
});

let currentTry = 0;
export async function connect() {
    if (currentTry > 5) {
        console.error('Failed to connect to DB');
        Deno.exit(3);
    }

    try {
        console.log('Try to connect to DB, try number: ' + ++currentTry);
        await sql`SELECT 1;`;
    } catch (err) {
        console.error(err);
        console.log('Trying to reconnect to DB in 5 seconds...');

        await new Promise<void>((resolve) =>
            setTimeout(async () => {
                await connect();
                resolve();
            }, 5000)
        );
    }
}

export { sql };
