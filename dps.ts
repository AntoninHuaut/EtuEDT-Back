import * as oak from "https://deno.land/x/oak/mod.ts";
import * as mysql from "https://deno.land/x/mysql/mod.ts";
import { moment } from "https://deno.land/x/moment/moment.ts";
import { BufReader } from "https://deno.land/std/io/bufio.ts";
import * as dotenv from "https://deno.land/x/dotenv/mod.ts";
import config from "./config/config.ts";

function now() {
  return moment().format("DD/MM/YYYY HH:mm");
}

const env = dotenv.config();

export { BufReader, config, env, moment, mysql, oak, now };
