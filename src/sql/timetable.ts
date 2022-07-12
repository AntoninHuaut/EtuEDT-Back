import { sql } from './index.ts';

export async function getAllTT(): Promise<any> {
    return await sql`SELECT * FROM "timetable" JOIN "univ" USING("numUniv") ORDER BY "numUniv", "numYearTT", "descTT"`;
}

export async function getAllUniv(): Promise<any> {
    return await sql`SELECT * FROM "univ" ORDER BY"numUniv"`;
}
