import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// 配置axios默认值
axios.defaults.baseURL = 'http://192.168.93.85:3000/api'//http://192.168.93.85:3000/api
axios.defaults.headers.common['Content-Type'] = 'application/json'

// 从localStorage获取token并设置到axios headers
const token = localStorage.getItem('token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

// 创建Vue应用实例
const app = createApp(App)

// 使用路由
app.use(router)

// 挂载应用
app.mount('#app') 