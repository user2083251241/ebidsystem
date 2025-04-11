<template>
    <div class="login-container">
      <h2>Bidsystem</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">Name</label>
          <input
            type="text"
            id="username"
            v-model="username"
            placeholder="Enter your name please"
            required
          />
        </div>
  
        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            placeholder="Enter your password please"
            required
          />
        </div>
  
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>
  
        <div class="form-group">
          <button type="submit">Login</button>
        </div>
      </form>
    </div>
  </template>
  
  <script>
  import { ref } from 'vue';
  import axios from 'axios';
  import { useRouter } from 'vue-router';
  
  export default {
    setup() {
      // 定义响应式数据
      const username = ref('');
      const password = ref('');
      const router = useRouter();
      const errorMessage = ref('');
  
      // 登录处理方法
      const handleLogin = async () => {
        try {
          errorMessage.value = '';
          const response = await axios.post('http://localhost:3000/api/login', {
            username: username.value,
            password: password.value
          });
  
          if (response.data.token) {
            // 保存 token 到 localStorage
            localStorage.setItem('token', response.data.token);
            // 保存用户信息
            localStorage.setItem('user', JSON.stringify(response.data.user));
            
            // 设置 axios 默认 headers
            axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.token}`;
            
            // 登录成功后跳转到首页
            router.push('/');
          }
        } catch (error) {
          console.error('登录失败:', error);
          errorMessage.value = error.response?.data?.message || 'Failed to login. Please check your credentials and try again.';
        }
      };
  
      return {
        username,
        password,
        errorMessage,
        handleLogin,
      };
    },
  };
  </script>
  
  <style scoped>
  .login-container {
    max-width: 300px;
    margin: 50px auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
    text-align: center;
  }
  
  h2 {
    margin-bottom: 20px;
  }
  
  .form-group {
    margin-bottom: 15px;
  }
  
  label {
    display: block;
    margin-bottom: 5px;
  }
  
  input {
    width: 100%;
    padding: 8px;
    box-sizing: border-box;
  }
  
  button {
    width: 100%;
    padding: 10px;
    background-color: #4caf50;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #45a049;
  }
  
  .error-message {
    color: #ff0000;
    margin-bottom: 15px;
    font-size: 14px;
  }
  </style>