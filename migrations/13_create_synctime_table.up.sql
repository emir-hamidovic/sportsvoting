CREATE TABLE IF NOT EXISTS `sync_time` (
    name ENUM('Regular', 'GOAT') NOT NULL,
    last_sync_time INT DEFAULT NULL,
);