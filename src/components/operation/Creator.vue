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
          <select id="orderType" v-model="orderType" required>
            <option value="market">Market</option>
            <option value="limit">Limit</option>
          </select>
        </div>
        <div class="form-group">
          <label for="direction">Direction</label>
          <select id="direction" v-model="direction" required>
            <option value="buy">Buy</option>
            <option value="sell">Sell</option>
          </select>
        </div>
        <button type="submit">Create Order</button>
      </form>
    </div>
  </template>
  
  <script>
  export default {
    name: 'Creator',
    data() {
      return {
        symbol: '',
        quantity: '',
        price: '',
        orderType: 'market',
        direction: 'buy'
      };
    },
    methods: {
      createOrder() {
        // 创建订单的逻辑
        console.log('Creating order:', {
          symbol: this.symbol,
          quantity: this.quantity,
          price: this.price,
          orderType: this.orderType,
          direction: this.direction
        });
        // 这里可以发送请求到后端创建订单
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