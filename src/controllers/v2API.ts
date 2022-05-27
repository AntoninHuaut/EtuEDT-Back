import { oak } from "../../deps.ts";
import * as cacheManager from "../cache/cacheManager.ts";
import Timetable from "../cache/TimeTable.ts";

export async function handleUnivList(context: oak.RouterContext<any>) {
  context.response.body = cacheManager.getUnivList();
}

export async function handleTTList(context: oak.RouterContext<string, Record<string | number, string>, Record<string, any>>) {
  const univ_TTList = cacheManager.getTTListByUniv(context.params.numUniv);

  if (!univ_TTList) {
    context.response.body = {};
    context.response.status = 404;
  }
  else {
    context.response.body = univ_TTList;
  }
}

export async function handleTTFormat(context: oak.RouterContext<string, Record<string | number, string>, Record<string, any>>) {
  const dataTT: Timetable | undefined = cacheManager.getTTById(context.params.numUniv, context.params.adeResources);

  if (!dataTT) {
    context.response.body = {};
    context.response.status = 404;
  }
  else if (!context.params.format) {
    context.response.body = dataTT.getAPIData();
  }
  else {
    switch (context.params.format) {
      case "json":
        context.response.body = dataTT.getJSON();
        break;
      case "ics":
        context.response.body = dataTT.getICS();
        context.response.headers.set("Content-Type", "text/calendar");
        break;
      default:
        context.response.body = {};
        context.response.status = 404;
        break;
    }
  }
}
