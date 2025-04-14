import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// ����axiosĬ��ֵ
axios.defaults.baseURL = 'http://192.168.93.85:3000/api'//http://192.168.93.85:3000/api
axios.defaults.headers.common['Content-Type'] = 'application/json'

// ��localStorage��ȡtoken�����õ�axios headers
const token = localStorage.getItem('token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

// ����VueӦ��ʵ��
const app = createApp(App)

// ʹ��·��
app.use(router)

// ����Ӧ��
app.mount('#app') 