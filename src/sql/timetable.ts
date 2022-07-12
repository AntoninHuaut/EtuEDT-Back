import { sql } from '/sql/index.ts';
import { ITimetable } from '/model/TimeTableModel.ts';
import { IUniv } from '/model/UnivModel.ts';

export async function getAllTT(): Promise<ITimetable[]> {
    return (await sql`SELECT * FROM "timetable" JOIN "univ" USING("numUniv") ORDER BY "numUniv", "numYearTT", "descTT"`) as unknown as ITimetable[];
}

export async function getAllUniv(): Promise<IUniv[]> {
    return (await sql`SELECT * FROM "univ" ORDER BY"numUniv"`) as unknown as IUniv[];
}
