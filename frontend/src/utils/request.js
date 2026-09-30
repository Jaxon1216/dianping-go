import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 创建axios实例
const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const data = response.data
    if (!data.success) {
      ElMessage.error(data.errorMsg || '请求失败')
      return Promise.reject(data.errorMsg || '请求失败')
    }
    return data.data
  },
  (error) => {
    console.error('请求错误:', error)
    
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.clearToken()
      ElMessage.error('请先登录')
      window.location.href = '/login'
      return Promise.reject('请先登录')
    }
    
    const message = error.response?.data?.errorMsg || '服务器异常'
    ElMessage.error(message)
    return Promise.reject(message)
  }
)

export default request
