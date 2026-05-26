ALTER TABLE `submissions` ADD COLUMN `user_id` char(36) NOT NULL;
CREATE INDEX `idx_submissions_user_id` ON `submissions`(`user_id`);
