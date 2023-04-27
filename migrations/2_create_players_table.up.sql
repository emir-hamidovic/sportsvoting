CREATE TABLE IF NOT EXISTS `players` (
  playerid   VARCHAR(128) PRIMARY KEY NOT NULL,
  name       VARCHAR(128) NOT NULL,
  teamabbr   VARCHAR(5),
  college    VARCHAR(128),
  height     VARCHAR(5) NOT NULL,
  weight     VARCHAR(5) NOT NULL,
  age        INT NOT NULL,
  FOREIGN KEY(teamabbr) REFERENCES `teams`(teamabbr)
);