CREATE TABLE IF NOT EXISTS `polls` (
  id           INT PRIMARY KEY AUTO_INCREMENT,
  name         VARCHAR(128) NOT NULL,
  description  VARCHAR(256) NOT NULL,
  image        VARCHAR(500) NOT NULL,
  selected_stats     VARCHAR(500) NOT NULL,
  season       CHAR(4) NOT NULL,
  userid       INT NOT NULL,
  FOREIGN KEY(userid) REFERENCES `users`(id)
);
