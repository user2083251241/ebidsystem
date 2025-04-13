<template>
  <div class="creator-container">
    <h2>Create New Order</h2>
    <form @submit.prevent="createOrder">
      <div class="form-group">
        <label for="symbol">Symbol</label>
        <input type="text" id="symbol" v-model="symbol" required />
      </div>
      <div class="form-group">
        <label for="quantity">Quantity</label>
        <input type="number" id="quantity" v-model="quantity" required />
      </div>
      <div class="form-group">
        <label for="price">Price</label>
        <input type="number" id="price" v-model="price" required />
      </div>
      <div class="form-group">
        <label for="orderType">Order Type</label>
        <select id="orderType" v-model="type" required>
          <option value="market">Market</option>
          <option value="limit">Limit</option>
        </select>
      </div>
      <div class="form-group">
        <label for="direction">Direction</label>
        <select id="direction" v-model="direction" required>
          <!-- <option value="buy">Buy</option> -->
          <option value="sell">Sell</option>
        </select>
      </div>
      <button type="submit">Create Order</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Creator',
  data() {
    return {
      symbol: '',
      quantity: '',
      price: '',
      type: String(this.type), // 确保是字符串类型,
      direction: 'sell'
    };
  },
  methods: {
    async createOrder() {
      console.log(this.orderType); // 检查orderType的值
      try {
        // 从localStorage获取token
        const token = localStorage.getItem('token');
        if (!token) {
          alert('Please login first');
          return;
        }

        // 从localStorage获取用户信息
        const user = JSON.parse(localStorage.getItem('user'));
        if (!user) {
          alert('User information not found');
          return;
        }

        // 构造订单数据
        const orderData = {
          symbol: this.symbol,
          quantity: this.quantity,
          price: this.price,
          type: this.type,
          direction: this.direction,
          //userID: user.id // 假设用户信息中包含用户ID
        };

        // 发送POST请求到后端创建订单
        const response = await axios.post('/seller/orders', orderData, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        // 处理成功响应
        if (response.status === 200) {
          alert('Order created successfully');
          console.log('Order created:', response.data);
        } else {
          alert('Order creation failed: ' + response.data.message);
        }
      } catch (error) {
        console.error('Create order error:', error);
        alert('Order creation failed: ' + error.response?.data?.message || 'Please check your input and try again');
      }
      console.log(this.type); // 检查orderType的值
    }
  }
};
</script>

<style scoped>
.creator-container {
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