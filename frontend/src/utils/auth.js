import Cookies from 'js-cookie'

const TokenKey = 'vulnark-token'

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  // 设置Cookie过期时间为7天，并添加路径参数确保在整个站点生效
  return Cookies.set(TokenKey, token, { 
    expires: 7, // 7天过期
    path: '/',  // 在整个站点生效
    sameSite: 'Lax' // 允许同站点请求携带Cookie
  })
}

export function removeToken() {
  return Cookies.remove(TokenKey, { path: '/' })
} 