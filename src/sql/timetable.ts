import { getConnection } from './index.ts';

export async function getAllTT(): Promise<any> {
    const connection = await getConnection();
    return await connection
        .query(`select * from TIMETABLE join UNIV using(numUniv) order by numUniv, numYearTT, descTT`)
        .finally(() => connection.close());
}

export async function getAllUniv(): Promise<any> {
    const connection = await getConnection();
    return await connection.query(`select * from UNIV order by numUniv`).finally(() => connection.close());
}