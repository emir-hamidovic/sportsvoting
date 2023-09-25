CREATE TABLE IF NOT EXISTS `player_votes` (
  id           INT PRIMARY KEY AUTO_INCREMENT,
  playerid     VARCHAR(128) NOT NULL,
  pollid       INT NOT NULL,
  userid       INT NOT NULL,
  votes_for     INT NOT NULL,
  FOREIGN KEY(playerid) REFERENCES `players`(playerid),
  FOREIGN KEY(pollid) REFERENCES `polls`(id),
  FOREIGN KEY(userid) REFERENCES `users`(id)
);