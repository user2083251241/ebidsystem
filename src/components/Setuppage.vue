<template>
    <div class="setup-page">
      <h1>Bidsystem</h1>
      <div class="page-container">
        <div v-if="currentView === 'login'">
          <LoginPage />
          <button @click="switchView('register')">Switch to Register</button>
        </div>
        <div v-else>
          <Register />
          <button @click="switchView('login')">Switch to Login</button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import LoginPage from './LoginPage.vue';
  import Register from './Register.vue';
  
  export default {
    components: {
      LoginPage,
      Register
    },
    data() {
      return {
        currentView: 'login' // 默认显示登录页面
      };
    },
    methods: {
      switchView(view) {
        this.currentView = view; // 切换当前显示的页面
      }
    }
  };
  </script>
  
  <style scoped>
  .setup-page {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-top: 50px;
  }
  
  .page-container {
    margin-top: 20px;
  }
  
  button {
    margin-top: 10px;
    padding: 10px 20px;
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