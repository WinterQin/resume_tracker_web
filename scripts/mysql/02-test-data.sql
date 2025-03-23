USE internship_manager;

-- 插入测试用户数据
INSERT INTO users (username, password, email) VALUES 
('test_user', '$2a$10$test_password_hash', 'test@example.com')
ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;

-- 插入测试申请数据
INSERT INTO applications (user_id, company, position, status, apply_date, event_link, salary, location) 
SELECT 
    (SELECT id FROM users WHERE username = 'test_user'),
    '测试公司',
    '测试职位',
    'submitted',
    NOW(),
    'http://example.com',
    '15k-20k',
    '北京'
WHERE EXISTS (SELECT 1 FROM users WHERE username = 'test_user'); 