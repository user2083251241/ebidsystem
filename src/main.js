import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// ����axiosĬ��ֵ
axios.defaults.baseURL = 'https://0125f20e-3482-4036-9bbd-f59f07ebd3f4.mock.pstmn.io'//http://192.168.81.85:3000/api
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