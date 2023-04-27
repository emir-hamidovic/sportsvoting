CREATE TABLE IF NOT EXISTS `teams` (
  teamabbr          VARCHAR(2) PRIMARY KEY,
  name              VARCHAR(128) NOT NULL,
  logo              VARCHAR(256) NOT NULL,
  winlosspct        FLOAT,
  playoffs          INT,
  divisiontitles    INT,
  conferencetitles  INT,
  championships     INT
);