CREATE TABLE IF NOT EXISTS univ (
    "numUniv"  serial,
    "nameUniv" varchar(64)  not null,
    "adeUniv"  varchar(255) not null,
    PRIMARY KEY ("numUniv")
);

CREATE TABLE IF NOT EXISTS timetable (
    "numUniv"      integer     not null,
    "adeResources" integer     not null,
    "numYearTT"    integer     not null,
    "descTT"       varchar(64) not null,
    "adeProjectId" integer     not null,
    PRIMARY KEY("numUniv", "adeResources"),
    CONSTRAINT timetable_fk FOREIGN KEY("numUniv") REFERENCES univ("numUniv")
);

INSERT INTO univ ("numUniv", "nameUniv", "adeUniv") VALUES (1, 'IUT Grand Ouest Normandie', 'http://ade.unicaen.fr/jsp/custom/modules/plannings/anonymous_cal.jsp');

INSERT INTO timetable ("numUniv", "adeResources", "numYearTT", "descTT", "adeProjectId") 
VALUES  (1, 1177, 1, 'TP 1.1', 3),
        (1, 1179, 1, 'TP 1.2', 3),
        (1, 1185, 1, 'TP 2.1', 3),
        (1, 1186, 1, 'TP 2.2', 3),
        (1, 1189, 1, 'TP 3.1', 3),
        (1, 1191, 1, 'TP 3.2', 3),
        (1, 1200, 2, 'TP 1.1', 3),
        (1, 1201, 2, 'TP 1.2', 3),
        (1, 1204, 2, 'TP 2.1', 3),
        (1, 1205, 2, 'TP 2.2', 3),
        (1, 1208, 2, 'TP 3.1', 3),
        (1, 1209, 2, 'TP 3.2', 3);