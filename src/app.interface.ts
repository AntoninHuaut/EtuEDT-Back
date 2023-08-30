export interface IEvent {
    title: string;
    enseignant: string;
    description?: string;
    start: Date;
    end: Date;
    location: string;
}

export interface IUniv {
    numUniv: number;
    nameUniv: string;
    adeUniv: string;
}

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
