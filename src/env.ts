import { config as loadEnv } from 'dotenv';
import dayjs from 'dayjs';

loadEnv({ export: true, path: `.env.${Deno.env.get('ENV')}` });

export const get = (field: string) => Deno.env.get(field);
export const now = () => dayjs().format('DD/MM/YYYY HH:mm');
