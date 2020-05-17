import { oak } from "../dps.ts";
import * as edtManage from "../cache/edtManage.ts";
import EDT from "../cache/EDT.ts";

export async function handleEDTList(context: oak.RouterContext<Record<string | number, string | undefined>, Record<string, any>>) {
  context.response.body = edtManage.getCacheEDT();
}

export async function handleEDTDetails(context: oak.RouterContext<Record<string | number, string | undefined>, Record<string, any>>) {
  const dataEDT: EDT | undefined = edtManage.getCacheEDTById(context.params.edtId);

  if (!dataEDT)
    context.response.body = {};
  else if (!context.params.format)
    context.response.body = dataEDT.getAPIData();
  else
    switch (context.params.format) {
      case "json":
        context.response.body = dataEDT.getJSON();
        break;
      case "ics":
        context.response.body = dataEDT.getICS();
        context.response.headers.set("Content-Type", "text/calendar");
        break;
      default:
        context.response.body = {};
        break;
    }
}
