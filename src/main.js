import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// 配置axios默认值
axios.defaults.baseURL = 'https://0125f20e-3482-4036-9bbd-f59f07ebd3f4.mock.pstmn.io'//http://192.168.81.85:3000/api
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