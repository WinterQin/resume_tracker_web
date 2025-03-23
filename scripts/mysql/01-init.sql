-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS internship_manager CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE internship_manager;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password VARCHAR(256) NOT NULL,
    email VARCHAR(128) NOT NULL UNIQUE,
    age INT,
    gender VARCHAR(10),
    phone VARCHAR(20),
    last_login_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建申请记录表
CREATE TABLE IF NOT EXISTS applications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    company VARCHAR(128) NOT NULL,
    position VARCHAR(128) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'submitted',
    apply_date DATETIME NOT NULL,
    next_event DATETIME NULL,
    event_type VARCHAR(32),
    event_link VARCHAR(256),
    notes TEXT,
    salary VARCHAR(64),
    location VARCHAR(128),
    contact_info VARCHAR(256),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 可以添加一些初始数据（可选）
INSERT INTO users (username, password, email) VALUES 
('admin', '$2a$10$your_hashed_password', 'admin@example.com')
ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP; 