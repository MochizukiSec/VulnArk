-- 检查同名用户是否存在并软删除
UPDATE users SET deleted_at = NOW() 
WHERE username = 'admin' AND deleted_at IS NULL;

-- 创建新的管理员用户
INSERT INTO users (
  username, 
  password, 
  email, 
  real_name, 
  role, 
  active, 
  last_login, 
  created_at, 
  updated_at
) VALUES (
  'admin',
  'a0.ru7f10zHnrrOsEsAFTMKDo5QwY4VDJguC', 
  'admin@vulnark.com',
  '系统管理员',
  'admin',
  1,
  '2025-03-15 21:50:52',
  '2025-03-15 21:50:52',
  '2025-03-15 21:50:52'
);

-- 查询确认创建成功
SELECT id, username, email, role, created_at FROM users 
WHERE username = 'admin' AND deleted_at IS NULL;
