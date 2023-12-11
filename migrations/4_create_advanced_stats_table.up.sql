CREATE TABLE IF NOT EXISTS `advancedstats` (
    id                INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    playerid          VARCHAR(128) NOT NULL,
    per  FLOAT,
    tspct  FLOAT,
    usgpct  FLOAT,
    ows  FLOAT,
    dws  FLOAT,
    ws  FLOAT,
    obpm  FLOAT,
    dbpm  FLOAT,
    bpm  FLOAT,
    vorp  FLOAT,
    offrtg FLOAT,
    defrtg FLOAT,
    season VARCHAR(25), /* year of the season */
    teamabbr VARCHAR(3),
    FOREIGN KEY(playerid) REFERENCES `players`(playerid),
    FOREIGN KEY(teamabbr) REFERENCES `teams`(teamabbr)
);