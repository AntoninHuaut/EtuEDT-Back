import { refreshTT } from "./cache/cacheManager.ts";

async function cli() {
    const buf = new Uint8Array(1024);

    await Deno.stdout.write(new TextEncoder().encode("\nCMD list:\n - refresh\nChoice: "));

    const n = <number>await Deno.stdin.read(buf);
    const answer = new TextDecoder().decode(buf.subarray(0, n)).trim();
    let response;

    switch (answer) {
        case 'refresh':
            refreshTT();
            response = "Refresh...";
            break;
        default:
            response = "Error: unknown cmd";
            break;
    }

    await Deno.stdout.write(new TextEncoder().encode("\n" + response));
    cli();
}

export default cli;