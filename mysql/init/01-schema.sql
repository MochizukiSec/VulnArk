-- 使用VulnArk数据库
USE vulnark;

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建settings表（如果不存在）
CREATE TABLE IF NOT EXISTS `settings` (
  `id` int(11) NOT NULL,
  `integrations` json DEFAULT NULL,
  `notifications` json DEFAULT NULL,
  `ai` json DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认设置
INSERT INTO `settings` (`id`, `integrations`, `notifications`, `ai`, `updated_at`, `updated_by`)
VALUES (1, '{}', '{}', '{}', NOW(), 1);

-- 创建setting_logs表（如果不存在）
CREATE TABLE IF NOT EXISTS `setting_logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `category` varchar(100) NOT NULL,
  `value` text DEFAULT NULL,
  `updated_by` int(11) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 重置外键检查
SET FOREIGN_KEY_CHECKS = 1; 