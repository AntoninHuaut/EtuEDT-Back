CREATE TABLE UNIV(
  numUniv INT AUTO_INCREMENT,
  nameUniv VARCHAR(64) NOT NULL,
  adeUniv VARCHAR(255) NOT NULL,
  PRIMARY KEY(numUniv)
);

CREATE TABLE TIMETABLE(
  numUniv INT,
  adeResources INT,
  numYearTT INT NOT NULL,
  descTT VARCHAR(64) NOT NULL,
  adeProjectId INT NOT NULL,
  PRIMARY KEY(numUniv, adeResources),
  FOREIGN KEY(numUniv) REFERENCES UNIV(numUniv)
);