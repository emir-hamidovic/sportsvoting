CREATE TABLE IF NOT EXISTS `users` (
  id           INT PRIMARY KEY AUTO_INCREMENT,
  username     VARCHAR(64),
  password     VARCHAR(256),
  refresh_token     VARCHAR(256),
  is_admin BOOLEAN
);