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
  
      // 登录处理方法
      const handleLogin = () => {
        // 仅在控制台打印用户名和密码
        console.log('用户名:', username.value);
        console.log('密码:', password.value);
        alert('Login successful!');
      };
  
      return {
        username,
        password,
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
  </style>