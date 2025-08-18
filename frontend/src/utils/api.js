import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 可以在这里添加token等认证信息
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

// 用户相关API
export const userApi = {
  // 用户注册
  register: (data) => api.post('/register', data),
  
  // 用户登录
  login: (data) => api.post('/login', data),
  
  // 获取用户信息
  getUserInfo: () => api.get('/user/info'),
  
  // 更新用户信息
  updateUserInfo: (data) => api.put('/user/info', data),
}

// 聊天相关API
export const chatApi = {
  // 获取会话列表
  getConversations: () => api.get('/conversations'),
  
  // 获取消息历史
  getMessages: (conversationId, page = 1) => 
    api.get(`/messages/${conversationId}?page=${page}`),
  
  // 发送消息
  sendMessage: (data) => api.post('/messages', data),
}

export default api
