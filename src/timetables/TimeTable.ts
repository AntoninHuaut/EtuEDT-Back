// deno-lint-ignore-file no-explicit-any
import dayjs from 'dayjs';
import ICAL from 'npm:ical.js';

import { IEvent, ITimetable, ITimetableExtended, IUniv } from '/src/app.interface.ts';

export default class Timetable implements ITimetable {
    numUniv: number;
    nameUniv: string;
    adeResources: number;
    adeProjectId: number;
    descTT: string;
    numYearTT: number;
    nameTT: string;
    lastUpdate: Date;

    ics: string;
    private json: IEvent[];

    constructor(ttData: ITimetable & IUniv, lastUpdate: Date, ics: string) {
        this.numUniv = ttData.numUniv;
        this.nameUniv = ttData.nameUniv;
        this.adeResources = ttData.adeResources;
        this.adeProjectId = ttData.adeProjectId;
        this.descTT = ttData.descTT;
        this.numYearTT = ttData.numYearTT;
        this.nameTT = this.numYearTT + 'A ' + this.descTT;

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

    getAPIData(): ITimetableExtended {
        return {
            numUniv: this.numUniv,
            nameUniv: this.nameUniv,
            nameTT: this.nameTT,
            numYearTT: this.numYearTT,
            descTT: this.descTT,
            adeResources: this.adeResources,
            adeProjectId: this.adeProjectId,
            lastUpdate: this.lastUpdate,
        };
    }

    getICS(): string {
        return this.ics;
    }

    getJSON(): IEvent[] {
        return this.json;
    }

    setJSON() {
        this.json = this.toJson();
        this.regroupJson();
    }

    private toJson(): IEvent[] {
        try {
            if (!this.ics || this.ics.includes('HTTP ERROR')) {
                throw new Error(`${this.numUniv}#${this.adeResources} not available`);
            }

            const calParse: any = ICAL.parse(this.ics.trim());
            const eventComps = new ICAL.Component(calParse).getAllSubcomponents('vevent');

            return eventComps
                .map((item: any) => {
                    if (!hasValue(item)) return null;

                    const event: IEvent = {
                        title: getValue(item, 'summary'),
                        enseignant: getTeacher(getValue(item, 'description')),
                        description: formatDescription(getValue(item, 'description')),
                        start: getValue(item, 'dtstart').toJSDate(),
                        end: getValue(item, 'dtend').toJSDate(),
                        location: getValue(item, 'location'),
                    };

                    return event;
                })
                .filter((el: any) => el !== null && el !== undefined) as IEvent[];
        } catch (_) {
            return [];
        }
    }

    private regroupJson() {
        this.sortJson();

        const output: IEvent[] = [];

        this.json.forEach((event: IEvent) => {
            const eventStart = dayjs(event.start);
            const eventEnd = dayjs(event.end);

            const existing = output.filter((check: IEvent) => {
                const checkStart = dayjs(check.start);
                const checkEnd = dayjs(check.end);

                return (
                    event.title === check.title &&
                    event.enseignant === check.enseignant &&
                    event.location === check.location &&
                    (eventStart.isSame(checkEnd) || eventEnd.isSame(checkStart))
                );
            });

            if (existing.length) {
                const existingIndex = output.indexOf(existing[0]);

                if (eventStart.isBefore(dayjs(existing[0].end))) {
                    event.end = existing[0].end;
                } else {
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
        this.json = this.json.sort((a: IEvent, b: IEvent) => {
            const diffDate = dayjs(a.start).diff(dayjs(b.start));
            return diffDate !== 0 ? diffDate : a.title.localeCompare(b.title);
        });
    }
}

function formatDescription(str: string | undefined) {
    return str ? str.trim().replace(/^\(Exported.*\n?/m, '') : str;
}

function getTeacher(description: string | null): string {
    if (!description) return '?';

    const descSplit = description.split(/\n/);
    const indexSlice = 3;
    if (descSplit.length - indexSlice < 0) return '?';

    return descSplit[descSplit.length - indexSlice];
}

function getValue(item: any, value: string): any {
    return item.getFirstPropertyValue(value) || '?';
}

function hasProperty(item: any, value: string): boolean {
    return !!item.getFirstProperty(value);
}

function hasValue(item: any): boolean {
    return hasProperty(item, 'summary') && hasProperty(item, 'description') && hasProperty(item, 'location') && hasProperty(item, 'dtstart') && hasProperty(item, 'dtend');
}
