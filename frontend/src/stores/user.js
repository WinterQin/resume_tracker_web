import { defineStore } from 'pinia'
import { userApi } from '../api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.token,
    userProfile: (state) => state.user
  },
  
  actions: {
    async login(credentials) {
      try {
        const response = await userApi.login(credentials)
        this.token = response.token
        localStorage.setItem('token', response.token)
        localStorage.setItem('userInfo', JSON.stringify(response.user))
        return response
      } catch (error) {
        throw error
      }
    },
    
    async register(userData) {
      try {
        const response = await userApi.register(userData)
        return response
      } catch (error) {
        throw error
      }
    },
    
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
    }
  }
}) 