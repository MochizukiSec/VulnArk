import request from '@/utils/request'

// 用户登录
export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

// 更新用户信息
export function updateUserInfo(data) {
  return request({
    url: '/user/update',
    method: 'put',
    data
  })
}

// 获取所有用户列表 (仅管理员)
export function getUserList() {
  return request({
    url: '/admin/users',
    method: 'get'
  })
}

// 创建用户 (仅管理员)
export function createUser(data) {
  return request({
    url: '/admin/users',
    method: 'post',
    data
  })
}

// 删除用户 (仅管理员)
export function deleteUser(userId) {
  return request({
    url: `/admin/user/${userId}`,
    method: 'delete'
  })
}

// 更改用户角色 (仅管理员)
export function changeUserRole(userId, role) {
  return request({
    url: `/admin/user/${userId}/role`,
    method: 'put',
    data: { role }
  })
}

// 更新用户个人资料
export function updateProfile(data) {
  return request({
    url: '/user/profile',
    method: 'put',
    data
  })
}

// 更新用户密码
export function updatePassword(data) {
  return request({
    url: '/user/password',
    method: 'put',
    data
  })
} 