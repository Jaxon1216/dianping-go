import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { getUserInfo, logout as logoutApi } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  // State
  const token = ref(sessionStorage.getItem('token') || '')
  const userInfo = ref(null)
  const userDetail = ref(null)
  
  // Getters
  const isLoggedIn = computed(() => !!token.value)
  const userId = computed(() => userInfo.value?.id)
  
  // Actions
  const setToken = (newToken) => {
    token.value = newToken
    sessionStorage.setItem('token', newToken)
  }
  
  const clearToken = () => {
    token.value = ''
    userInfo.value = null
    userDetail.value = null
    sessionStorage.removeItem('token')
    sessionStorage.removeItem('userInfo')
  }
  
  const fetchUserInfo = async () => {
    try {
      const data = await getUserInfo()
      userInfo.value = data
      return data
    } catch (error) {
      clearToken()
      throw error
    }
  }
  
  const setUserDetail = (detail) => {
    userDetail.value = detail
    sessionStorage.setItem('userInfo', JSON.stringify(detail))
  }
  
  const logout = async () => {
    try {
      await logoutApi()
    } finally {
      clearToken()
    }
  }
  
  // 初始化时如果有token，获取用户信息
  const init = async () => {
    if (token.value && !userInfo.value) {
      try {
        await fetchUserInfo()
      } catch {
        // 获取失败，token可能已过期
      }
    }
  }
  
  return {
    token,
    userInfo,
    userDetail,
    isLoggedIn,
    userId,
    setToken,
    clearToken,
    fetchUserInfo,
    setUserDetail,
    logout,
    init
  }
})
