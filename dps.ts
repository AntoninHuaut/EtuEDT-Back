import * as oak from "https://deno.land/x/oak@v10.1.0/mod.ts";
import * as mysql from "https://deno.land/x/mysql@v2.10.2/mod.ts";
import { moment } from "https://deno.land/x/deno_moment@v1.1.2/moment.ts";
import { BufReader } from "https://deno.land/std@0.120.0/io/buffer.ts";
import config from "./config/config.ts";

function now() {
  return moment().format("DD/MM/YYYY HH:mm");
}

export { BufReader, config, moment, mysql, oak, now };
