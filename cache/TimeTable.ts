import { moment } from '../dps.ts';
import ICAL from '../lib/ical.js';

export default class Timetable {

    private numUniv: number;
    private nameUniv: string;
    private adeResources: number;
    private adeProjectId: number;
    private descTT: string;
    private numYearTT: number;
    private nameTT: string;
    private lastUpdate: Date;

    private ics: string;
    private json: [];

    constructor(ttData: any, lastUpdate: Date, ics: string) {
        this.numUniv = ttData.numUniv;
        this.nameUniv = ttData.nameUniv;
        this.adeResources = ttData.adeResources;
        this.adeProjectId = ttData.adeProjectId;
        this.descTT = ttData.descTT;
        this.numYearTT = ttData.numYearTT;
        this.nameTT = this.numYearTT + "A " + this.descTT;

        this.lastUpdate = lastUpdate;
        this.ics = ics;
        this.json = [];

        this.setJSON();
    }

    getNumUniv(): number {
        return this.numUniv;
    }

    getAdeResources(): number {
        return this.adeResources;
    }

    getAPIData(): {} {
        return {
            numUniv: this.numUniv,
            nameUniv: this.nameUniv,
            nameTT: this.nameTT,
            numYearTT: this.numYearTT,
            descTT: this.descTT,
            adeResources: this.adeResources,
            adeProjectId: this.adeProjectId,
            lastUpdate: this.lastUpdate
        };
    }

    getICS(): string {
        return this.ics;
    }

    getJSON(): [] {
        return this.json;
    }

    setJSON() {
        this.json = this.toJson();
        this.regroupJson();
    }

    private toJson(): any {
        try {
            if (!this.ics || this.ics.includes('HTTP ERROR')) {
                throw new Error(`${this.numUniv}#${this.adeResources} not available`);
            }

            const calParse: any = ICAL.parse(this.ics.trim());
            let eventComps: any[] = new ICAL.Component(calParse).getAllSubcomponents("vevent");

            let events = eventComps.map((item: any) => {
                if (!hasValue(item)) return null;

                return {
                    "title": getValue(item, 'summary'),
                    "enseignant": getEnseignant(getValue(item, 'description')),
                    "description": removeExported(getValue(item, 'description')),
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

    private regroupJson(): any {
        this.sortJson();

        const output: any = [];

        this.json.forEach((event: any) => {
            const eventStart = moment(event.start);
            const eventEnd = moment(event.end);

            const existing = output.filter((check: any) => {
                const checkStart = moment(check.start);
                const checkEnd = moment(check.end);

                return event.title === check.title
                    && event.enseignant === check.enseignant
                    && event.location === check.location
                    && (eventStart.isSame(checkEnd) || eventEnd.isSame(checkStart))
            });

            if (existing.length) {
                const existingIndex = output.indexOf(existing[0]);

                if (eventStart.isBefore(moment(existing[0].end))) {
                    event.end = existing[0].end;
                }
                else {
                    event.start = existing[0].start;
                }

                output[existingIndex] = event;
            } else {
                output.push(event);
            }
        });

        this.json = output;
        this.sortJson();
    }

    private sortJson() {
        this.json = this.json.sort((a: any, b: any) => {
            const diffDate = moment(a.start).diff(moment(b.start));
            return diffDate !== 0 ? diffDate : a.title.localeCompare(b.title);
        });
    }
}

function removeExported(str: string | null) {
    return str ? str.replace(/^\(Exported.*\n?/m, '') : str;
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