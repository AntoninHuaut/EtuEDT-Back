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
    numYearTT: number;
    descTT: string;
    adeResources: number;
    adeProjectId: number;
}

export type ITimetableUniv = ITimetable & IUniv;

export type ITimetableAPI = ITimetable & {
    nameTT: string;
    nameUniv: string;
    lastUpdate: Date;
};
