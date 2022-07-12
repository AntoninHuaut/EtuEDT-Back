import { config as loadEnv } from 'dotenv';

loadEnv({ export: true });

export const get = (field: string) => Deno.env.get(field);
