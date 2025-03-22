-- 删除测试用户
DELETE FROM users WHERE username = 'admin123';
DELETE FROM users WHERE username = 'testadmin999';
DELETE FROM users WHERE username LIKE '%test%';
DELETE FROM users WHERE email LIKE '%test%';
DELETE FROM users WHERE real_name LIKE '%测试%';

-- 删除测试资产
DELETE FROM assets WHERE name LIKE '%测试%';
DELETE FROM assets WHERE identifier LIKE '%测试%';
DELETE FROM assets WHERE notes LIKE '%测试%';
DELETE FROM assets WHERE owner LIKE '%测试%';

-- 删除测试漏洞
DELETE FROM vulnerabilities WHERE title LIKE '%测试%';
DELETE FROM vulnerabilities WHERE description LIKE '%测试%';
DELETE FROM vulnerabilities WHERE notes LIKE '%测试%';

-- 删除测试知识库条目
DELETE FROM knowledge WHERE title LIKE '%测试%';
DELETE FROM knowledge WHERE content LIKE '%测试%';

-- 删除测试扫描任务
DELETE FROM scan_tasks WHERE name LIKE '%测试%';
DELETE FROM scan_tasks WHERE description LIKE '%测试%';

-- 删除测试扫描结果
DELETE FROM scan_results WHERE report LIKE '%测试%';

-- 确保admin账户仍然存在
INSERT INTO users (username, password, email, real_name, role, active, last_login, created_at, updated_at)
SELECT 'admin', '$2a$10$BhmVtbPkX/tSXbDVvM2aVO1qrSo2Fn8v4G2uh38EeHIcOEg6ftkDe', 'admin@vulnark.com', '系统管理员', 'admin', 1, NOW(), NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin' AND deleted_at IS NULL); 