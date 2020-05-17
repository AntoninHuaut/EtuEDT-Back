import { getConnection } from './index.ts';

export async function getAllEDT(): Promise<any> {
    const connection = await getConnection();
    return await connection
        .query(`select * from EDT join Etablissement using(numEta) order by numEta, numAnnee, numTP`)
        .finally(() => connection.close());
}