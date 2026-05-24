ALTER TABLE `problems` ADD COLUMN `time_limit` int NOT NULL DEFAULT 1000;
ALTER TABLE `problems` ADD COLUMN `memory_limit` int NOT NULL DEFAULT 256;