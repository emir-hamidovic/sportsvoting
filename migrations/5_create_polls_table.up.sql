CREATE TABLE IF NOT EXISTS `polls` (
  id           INT PRIMARY KEY AUTO_INCREMENT,
  name         VARCHAR(128) NOT NULL,
  description  VARCHAR(256) NOT NULL,
  image        VARCHAR(500) NOT NULL,
  endpoint     VARCHAR(128) NOT NULL
);
