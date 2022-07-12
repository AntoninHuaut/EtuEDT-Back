import { sql } from '/sql/index.ts';
import { ITimeTableUniv } from '/model/TimeTableModel.ts';
import { IUniv } from '/model/UnivModel.ts';

export async function getAllTT(): Promise<ITimeTableUniv[]> {
    return (await sql`SELECT * FROM "timetable" JOIN "univ" USING("numUniv") ORDER BY "numUniv", "numYearTT", "descTT"`) as unknown as ITimeTableUniv[];
}

export async function getAllUniv(): Promise<IUniv[]> {
    return (await sql`SELECT * FROM "univ" ORDER BY "numUniv"`) as unknown as IUniv[];
}
