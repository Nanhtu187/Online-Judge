CREATE TABLE `problems` (
  `id` char(36) NOT NULL,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `input_format` text,
  `output_format` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_problems_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `test_cases` (
  `id` char(36) NOT NULL,
  `problem_id` char(36) NOT NULL,
  `input` text,
  `expected_output` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_test_cases_problem_id` (`problem_id`),
  KEY `idx_test_cases_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_problems_test_cases` FOREIGN KEY (`problem_id`) REFERENCES `problems` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;