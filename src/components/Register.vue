<template>
  <div class="register-container">
    <h2>Signup</h2>
    <form @submit.prevent="handleRegister">
      <div class="form-group">
        <label for="username">Account</label>
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
          <option value="client">client</option>
          <option value="sales">sales</option>
          <option value="trader">trader</option>
        </select>
      </div>

      <div class="form-group">
        <button type="submit">Signup</button>
      </div>
    </form>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';

export default {
  setup() {
    const username = ref('');
    const password = ref('');
    const confirmPassword = ref('');
    const role = ref('');
    const errorMessage = ref('');

    const handleRegister = async () => {
      // 检查密码是否一致
      if (password.value !== confirmPassword.value) {
        errorMessage.value = 'The passwords do not match';
        return;
      }

      try {
        // 发送 POST 请求到后端注册接口
        const response = await axios.post('/register', {
          username: username.value,
          password: password.value,
          role: role.value
        });

        // 处理成功响应
        if (response.status === 201) {
          alert('Registration successful! User ID: ' + response.data.user_id);
          console.log('Registration successful:', response.data);
        } else {
          // 处理其他状态码
          errorMessage.value = 'Registration failed: ' + response.data.message;
        }
      } catch (error) {
        // 处理错误响应
        if (error.response) {
          errorMessage.value = error.response.data.error || 'Registration failed';
        } else {
          errorMessage.value = 'Network error or invalid response';
        }
        console.error('Registration failed:', error);
      }
    };

    return {
      username,
      password,
      confirmPassword,
      role,
      errorMessage,
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

.error-message {
  color: #ff0000;
  margin-bottom: 15px;
  font-size: 14px;
}
</style>