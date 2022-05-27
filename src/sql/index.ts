import {
    BufReader,
    mysql,
} from '../../deps.ts';

export async function initTable() {
    const file = await Deno.open('./src/sql/table.sql');
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
    return await new mysql.Client().connect({
        hostname: Deno.env.get('MYSQL_HOST') || Deno.exit(1),
        port: parseInt(Deno.env.get('MYSQL_PORT') ?? "") || Deno.exit(1),
        db: Deno.env.get('MYSQL_DATABASE') || Deno.exit(1),
        username: Deno.env.get('MYSQL_USER') || Deno.exit(1),
        password: Deno.env.get('MYSQL_PASSWORD') || Deno.exit(1),
    });
};