CREATE TABLE `test_case_results` (
  `id` char(36) NOT NULL,
  `submission_id` char(36) NOT NULL,
  `test_case_id` char(36) NOT NULL,
  `status` varchar(20) NOT NULL,
  `actual_output` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_test_case_results_submission_id` (`submission_id`),
  KEY `idx_test_case_results_test_case_id` (`test_case_id`),
  KEY `idx_test_case_results_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_submissions_test_case_results` FOREIGN KEY (`submission_id`) REFERENCES `submissions` (`id`),
  CONSTRAINT `fk_test_cases_results` FOREIGN KEY (`test_case_id`) REFERENCES `test_cases` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;