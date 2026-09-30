import request from '@/utils/request'

const extractToken = (data) => data?.token ?? data

// 发送验证码
export const sendCode = (phone) => {
  return request.post(`/user/code?phone=${phone}`)
}

// 手机号+验证码登录
export const loginByCode = (data) => {
  return request.post('/user/login', data).then(extractToken)
}

// 手机号+密码登录
export const loginByPassword = (data) => {
  return request.post('/user/login', data).then(extractToken)
}

// 获取当前登录用户信息
export const getUserInfo = () => {
  return request.get('/user/me')
}

// 获取用户信息
export const getUserById = (id) => {
  return request.get(`/user/${id}`)
}

// 获取用户详情
export const getUserDetail = (id) => {
  return request.get(`/user/info/${id}`)
}

// 退出登录
export const logout = () => {
  return request.post('/user/logout')
}
