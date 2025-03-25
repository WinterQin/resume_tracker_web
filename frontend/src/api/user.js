import api from './index'

export const userApi = {
  login(data) {
    return api.post('/auth/login', data)
  },
  register(data) {
    return api.post('/auth/register', data)
  },
  // 获取个人信息
  getProfile() {
    return api.get('/user/profile')
  },
  // 更新个人信息
  updateProfile(data) {
    return api.put('/user/profile', data)
  },
  // 更新密码
  updatePassword(data) {
    return api.put('/user/password', data)
  },
  // 删除账号
  deleteAccount() {
    return api.delete('/user/account')
  }
} 