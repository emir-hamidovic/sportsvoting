CREATE TABLE IF NOT EXISTS `polls` (
  id           INT PRIMARY KEY AUTO_INCREMENT,
  name         VARCHAR(128) NOT NULL,
  description  VARCHAR(256) DEFAULT "",
  image        VARCHAR(500) DEFAULT "",
  selected_stats     VARCHAR(500) NOT NULL,
  season       CHAR(25) NOT NULL,
  userid       INT NOT NULL,
  FOREIGN KEY(userid) REFERENCES `users`(id)
);
