CREATE TABLE IF NOT EXISTS `team_votes` (
  id           INT PRIMARY KEY AUTO_INCREMENT,
  teamabbr     VARCHAR(3),
  pollid       INT NOT NULL,
  votes_for    INT NOT NULL,
  FOREIGN KEY(teamabbr) REFERENCES `teams`(teamabbr),
  FOREIGN KEY(pollid) REFERENCES `polls`(id)
);