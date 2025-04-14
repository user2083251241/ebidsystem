<template>
  <div class="table-container">
    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <table v-else>
      <thead>
        <tr>
          <th v-for="header in headers" :key="header">{{ header }}</th>
          <!-- <th>Operate</th> -->
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, rowIndex) in tableData" :key="rowIndex">
          <td v-for="(cell, cellIndex) in row" :key="cellIndex">{{ cell }}</td>
          <td>
            <button @click="openPurchaseModal(rowIndex)">Purchase</button>
          </td>
        </tr>
      </tbody>
    </table>
    <Purchase
      v-if="isPurchaseModalOpen"
      :order="selectedOrder"
      @close="closePurchaseModal"
      @submit-order="submitOrder"
    />
  </div>
</template>

<script>
import axios from 'axios';
import Purchase from './operation/Purchase.vue';

export default {
  name: 'ClientTableComponent',
  components: {
    Purchase
  },
  data() {
    return {
      headers: ['ID', 'Symbol', 'Quantity', 'Price', 'OrderType', 'Status', 'Operate'],
      tableData: [],
      isPurchaseModalOpen: false,
      selectedOrder: null,
      loading: true,
      error: null,
      TEMP: null
    };
  },
  async created() {
    await this.fetchUserProducts();
    console.log('clientTableComponent created');
  },
  methods: {
    async fetchUserProducts() {
      try {
        this.loading = true;
        this.error = null;

        // 从 localStorage 获取 token
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error('Please login first');
        }

        // 发送请求获取用户商品列表
        const response = await axios.get('/client/orders', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        console.log('Response data:', response.data);

        // 处理返回的数据
        if (response.data && Array.isArray(response.data)) {
          this.tableData = response.data.map(product => [
            product.ID,
            product.Symbol,
            product.Quantity,
            product.Price,
            product.OrderType,
            product.Status
          ]);
        } else {
          this.tableData = [];
        }
      } catch (error) {
        console.error('Fetch order error:', error);
        this.error = error.response?.data?.message || 'Fetch order failed';
      } finally {
        this.loading = false;
      }
    },
    openPurchaseModal(index) {
      this.selectedOrder = this.tableData[index];
      this.isPurchaseModalOpen = true;
      console.log(this.selectedOrder);
    },
    closePurchaseModal() {
      this.isPurchaseModalOpen = false;
    },
    async submitOrder(order) {
      try {
        const token = localStorage.getItem('token');console.log("235"+order);
        console.log("234"+order);
        if (!token) {
          throw new Error('Please login first');
        }
        console.log("233"+order);
        this.TEMP = order[0];
        console.log("TEMP"+TEMP);
        const response = await axios.post('/client/orders/buy', order, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (response.status === 200) {
          alert('Order submitted successfully' + (response.data.message || ''));
          this.closePurchaseModal();
        } else {
          console.log("250"+order);
          alert('Order submission failed' + (response.data.message || ''));
        }
      } catch (error) {
        console.log("251"+order);
        alert(`Order submission failed: ${error.response?.data?.message || 'Unknown error'}`);
      }
    }
  }
};
</script>
  
  <style scoped>
  .table-container {
    margin-top: 20px;
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
  }
  
  th, td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: center;
  }
  
  th {
    background-color: #f2f2f2;
    font-weight: bold;
  }
  
  tr:nth-child(even) {
    background-color: #f9f9f9;
  }
  
  tr:hover {
    background-color: #f1f1f1;
  }
  
  .loading {
    text-align: center;
    padding: 20px;
    color: #666;
  }
  
  button {
    margin: 0 5px;
    padding: 5px 10px;
    background-color: #4caf50;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #45a049;
  }
  
  .error {
    text-align: center;
    padding: 20px;
    color: #ff0000;
  }
  </style>