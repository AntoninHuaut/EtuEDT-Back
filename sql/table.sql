CREATE TABLE TP(numTP VARCHAR(24) NOT NULL, PRIMARY KEY(numTP));

CREATE TABLE Annee(numAnnee INT, PRIMARY KEY(numAnnee));

CREATE TABLE Etablissement(
  numETa INT AUTO_INCREMENT,
  nomEta VARCHAR(64) NOT NULL,
  adeEta VARCHAR(255) NOT NULL,
  PRIMARY KEY(numEta)
);

CREATE TABLE EDT(
  numTP VARCHAR(24),
  numAnnee INT,
  numEta INT,
  resources VARCHAR(16) NOT NULL,
  projectId VARCHAR(16) NOT NULL,
  PRIMARY KEY(numTP, numAnnee, numEta),
  FOREIGN KEY(numTP) REFERENCES TP(numTP),
  FOREIGN KEY(numAnnee) REFERENCES Annee(numAnnee),
  FOREIGN KEY(numEta) REFERENCES Etablissement(numEta)
);