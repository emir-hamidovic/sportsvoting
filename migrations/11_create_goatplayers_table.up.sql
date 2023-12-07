CREATE TABLE IF NOT EXISTS `goat_players` (
  playerid      VARCHAR(128) PRIMARY KEY NOT NULL,
  name          VARCHAR(128) NOT NULL,
  allstar       INT,
  allnba        INT,
  alldefense    INT,
  championships INT,
  dpoy          INT,
  sixman        INT,
  roy           INT,
  finalsmvp     INT,
  mvp           INT,
  isactive      BOOLEAN DEFAULT false
);