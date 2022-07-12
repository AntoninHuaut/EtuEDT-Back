import { IUniv } from '/model/UnivModel.ts';

export interface ITimetable {
    numUniv: number;
    adeResources: number;
    numYearTT: number;
    descTT: string;
    adeProjectId: number;
}

export type ITimeTableUniv = ITimetable & IUniv;

export interface ITimetableExtended extends ITimetable {
    nameTT: string;
    nameUniv: string;
    lastUpdate: Date;
}
