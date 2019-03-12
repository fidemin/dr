
DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `dispatch`;

CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `usertype` enum('PA','DR') COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `dispatch` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pa_id` int(11) NOT NULL,
  `dr_id` int(11) DEFAULT NULL,
  `address` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_complete` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `completed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_dispatch_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



