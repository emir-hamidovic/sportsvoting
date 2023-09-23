CREATE TABLE IF NOT EXISTS `users` (
  id            INT PRIMARY KEY AUTO_INCREMENT,
  username      VARCHAR(64),
  email         VARCHAR(64),
  password      TEXT,
  refresh_token TEXT,
  profile_pic   VARCHAR(255) DEFAULT ""
);