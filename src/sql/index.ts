import postgres from 'postgres';
import { get } from '/config.ts';

const DB_HOST = get('POSTGRES_HOST');
const DB_PORT = get('POSTGRES_PORT') ?? '';
const DB_DATABASE = get('POSTGRES_DB');
const DB_USER = get('POSTGRES_USER');
const DB_PASSWORD = get('POSTGRES_PASSWORD');

if (!DB_HOST || isNaN(+DB_PORT) || !DB_DATABASE || !DB_USER || !DB_PASSWORD) {
    console.error('Invalid DB configuration');
    Deno.exit(3);
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

export { sql };
