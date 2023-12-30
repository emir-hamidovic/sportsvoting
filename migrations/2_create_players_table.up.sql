CREATE TABLE IF NOT EXISTS `players` (
  playerid   VARCHAR(128) PRIMARY KEY NOT NULL,
  name       VARCHAR(128) NOT NULL,
  teamabbr   VARCHAR(3) default "",
  college    VARCHAR(128) default "",
  height     VARCHAR(5) default "",
  weight     VARCHAR(3) default "",
  age        INT NOT NULL,
  FOREIGN KEY(teamabbr) REFERENCES `teams`(teamabbr)
);