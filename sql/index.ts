import {
    BufReader,
    config,
    mysql,
} from '../dps.ts';

export async function initTable() {
    const file = await Deno.open('./sql/table.sql');
    const bufReader = new BufReader(file);
    const textDecoder = new TextDecoder("utf-8");

    let connection = await getConnection();

    let statement: string | any;
    let currentLine: string = "";

    while ((statement = await bufReader.readLine()) != null) {
        currentLine += textDecoder.decode(statement.line);

        if (!currentLine.trim().endsWith(';')) continue;

        currentLine = currentLine.replace(/(\r\n|\n|\r)/gm, '');

        await connection.execute(currentLine).catch(() => { });
        currentLine = "";
    }

    await connection.close();
    file.close();
}

export async function getConnection() {
    return await new mysql.Client().connect(config.sql);
};