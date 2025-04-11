<template>
    <div class="register-container">
      <h2>Signup</h2>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="username">Accunt</label>
          <input
            type="text"
            id="username"
            v-model="username"
            placeholder="Input your account"
            required
          />
        </div>
  
        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            placeholder="Input your password"
            required
          />
        </div>
  
        <div class="form-group">
          <label for="confirm-password">Password Confirm</label>
          <input
            type="password"
            id="confirm-password"
            v-model="confirmPassword"
            placeholder="Confirm your password"
            required
          />
        </div>
  
        <div class="form-group">
          <label for="role">Role</label>
          <select id="role" v-model="role" required>
            <option value="buyer">buyer</option>
            <option value="seller">Seller</option>
            <option value="trader">trader</option>
          </select>
        </div>
  
        <div class="form-group">
          <button type="submit">Signup</button>
        </div>
      </form>
    </div>
  </template>
  
  <script>
  import { ref } from 'vue';
  
  export default {
    setup() {
      const username = ref('');
      const password = ref('');
      const confirmPassword = ref('');
      const role = ref('');
  
      const handleRegister = () => {
        if (password.value !== confirmPassword.value) {
          alert('The passwords do not match');
          return;
        }
  
        console.log('ÕËºÅ:', username.value);
        console.log('ÃÜÂë:', password.value);
        console.log('½ÇÉ«:', role.value);
        alert('Signup successful!');
      };
  
      return {
        username,
        password,
        confirmPassword,
        role,
        handleRegister,
      };
    },
  };
  </script>
  
  <style scoped>
  .register-container {
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
  
  input, select {
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