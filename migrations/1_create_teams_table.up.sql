CREATE TABLE IF NOT EXISTS `teams` (
  id                INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  teamabbr          VARCHAR(5),
  name              VARCHAR(128) NOT NULL,
  logo              VARCHAR(256) NOT NULL,
  winlosspct        FLOAT,
  playoffs          INT,
  divisiontitles    INT,
  conferencetitles  INT,
  championships     INT,
  UNIQUE(teamabbr)
);