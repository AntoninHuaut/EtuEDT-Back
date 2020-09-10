import ICAL from '../lib/ical.js';

export default class EDT {

    private numEta: number;
    private nomEta: string;
    private numAnnee: number;
    private numTP: string;
    private nomEDT: string;
    private edtName: string;
    private edtId: string;
    private lastUpdate: Date;

    private edtIcs: string;
    private edtJson: [];

    constructor(edtData: any, lastUpdate: Date, edtIcs: string) {
        this.numEta = edtData.numEta;
        this.nomEta = edtData.nomEta;
        this.numAnnee = edtData.numAnnee;
        this.numTP = edtData.numTP;
        this.nomEDT = edtData.nomEDT;
        this.edtName = this.numAnnee + "A TP " + this.numTP;
        this.edtId = edtData.resources;
        this.lastUpdate = lastUpdate;
        this.edtIcs = edtIcs;
        this.edtJson = [];

        this.setJSON();
    }

    getEdtId(): string {
        return this.edtId;
    }

    getAPIData(): {} {
        return {
            numEta: this.numEta,
            nomEta: this.nomEta,
            numAnnee: this.numAnnee,
            numTP: this.numTP,
            nomEDT: this.nomEDT,
            edtName: this.edtName,
            edtId: this.edtId,
            lastUpdate: this.lastUpdate
        };
    }

    getICS(): string {
        return this.edtIcs;
    }

    getJSON(): [] {
        return this.edtJson;
    }

    setJSON() {
        this.edtJson = this.toJson();
    }

    private toJson(): any {
        try {
            if (!this.edtIcs || this.edtIcs.includes('HTTP ERROR')) throw new Error(this.edtId + " not available");

            const calParse: any = ICAL.parse(this.edtIcs.trim());
            let eventComps: any[] = new ICAL.Component(calParse).getAllSubcomponents("vevent");

            let events = eventComps.map((item: any) => {
                if (!hasValue(item)) return null;

                return {
                    "title": getValue(item, 'summary'),
                    "enseignant": getEnseignant(getValue(item, 'description')),
                    "start": getValue(item, 'dtstart').toJSDate(),
                    "end": getValue(item, 'dtend').toJSDate(),
                    "location": getValue(item, 'location')
                };
            });

            return events.filter((el: any) => el != null);
        } catch (ex) {
            return [];
        }
    }
}

function getEnseignant(description: string): string {
    const descSplit = description.split(/\n/);
    const indexSlice = 3;
    if (descSplit.length - indexSlice < 0) return "?";

    return descSplit[descSplit.length - indexSlice];
}

function getValue(item: any, value: string): any {
    return item.getFirstPropertyValue(value) || "?";
}

function hasProperty(item: any, value: string): any {
    return !!item.getFirstProperty(value);
}

function hasValue(item: any): boolean {
    return hasProperty(item, 'summary') &&
        hasProperty(item, 'description') &&
        hasProperty(item, 'location') &&
        hasProperty(item, 'dtstart') &&
        hasProperty(item, 'dtend');
}