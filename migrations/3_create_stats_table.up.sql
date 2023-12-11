CREATE TABLE IF NOT EXISTS `stats` (
    id                INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    playerid          VARCHAR(128) NOT NULL,
    gamesplayed       INT,
    gamesstarted       INT,
    minutespergame    FLOAT,
    pointspergame     FLOAT,
    reboundspergame   FLOAT,
    assistspergame    FLOAT,
    stealspergame     FLOAT,
    fgpercentage      FLOAT,
    threeptpercentage FLOAT,
    ftpercentage      FLOAT,
    blockspergame     FLOAT,
    turnoverspergame  FLOAT,
    season VARCHAR(25), /* year of the season */
    position  VARCHAR(5),
    teamabbr VARCHAR(3),
    rookieseason BOOLEAN DEFAULT false,
    FOREIGN KEY(playerid) REFERENCES `players`(playerid),
    FOREIGN KEY(teamabbr) REFERENCES `teams`(teamabbr)
);