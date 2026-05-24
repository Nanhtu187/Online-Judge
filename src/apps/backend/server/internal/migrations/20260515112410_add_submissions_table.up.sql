CREATE TABLE `submissions` (
  `id` char(36) NOT NULL,
  `problem_id` char(36) NOT NULL,
  `code_content` text NOT NULL,
  `status` varchar(20) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_submissions_problem_id` (`problem_id`),
  KEY `idx_submissions_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_submissions_problems` FOREIGN KEY (`problem_id`) REFERENCES `problems` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;