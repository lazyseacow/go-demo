-- 创建数据库
CREATE DATABASE IF NOT EXISTS go_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE go_demo;

-- 用户表（GORM 会自动创建，这里仅作为参考）
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT 1 COMMENT '状态 1正常 0禁用',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 插入测试用户（密码为 123456）
-- 注意: 实际使用时应该通过 API 注册，这里的密码已经过 bcrypt 加密
INSERT INTO `users` (`username`, `password`, `email`, `status`, `created_at`, `updated_at`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'admin@example.com', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `username` = `username`;

