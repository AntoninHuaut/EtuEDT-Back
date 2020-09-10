import * as oak from "https://deno.land/x/oak/mod.ts";
import * as mysql from "https://deno.land/x/mysql/mod.ts";
import { moment } from "https://deno.land/x/deno_moment@1.0.0/moment.ts";
import { BufReader } from "https://deno.land/std/io/bufio.ts";
import config from "./config/config.ts";

function now() {
  return moment().format("DD/MM/YYYY HH:mm");
}

export { BufReader, config, moment, mysql, oak, now };
