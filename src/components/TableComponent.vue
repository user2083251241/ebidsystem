<template>
  <div class="table-container">
    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <table v-else>
      <thead>
        <tr>
          <th v-for="header in headers" :key="header">{{ header }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, rowIndex) in tableData" :key="rowIndex">
          <td v-for="(cell, cellIndex) in row" :key="cellIndex">{{ cell }}</td>
          <td>
            <button @click="confirmDelete(rowIndex)">Delete</button>
            <button @click="showModifyModal(rowIndex)">Modify</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div class="create-button">
      <button @click="goToCreate">Create</button>
    </div>
    <ConfirmModal
      :isVisible="isConfirmModalVisible"
      title="Confirm Delete"
      message="Are you sure you want to delete this order?"
      @confirm="deleteOrder(selectedIndex)"
      @cancel="closeConfirmModal"
    />
    <ModifyModal
      :isVisible="isModifyModalVisible"
      title="Modify Order"
      :initialQuantity="selectedOrderQuantity"
      :initialPrice="selectedOrderPrice"
      @submit="modifyOrderSubmit"
      @close="closeModifyModal"
    />
  </div>
</template>

<script>
import axios from 'axios';
import ConfirmModal from './ConfirmModal.vue';
import ModifyModal from './ModifyModal.vue';

export default {
  name: 'TableComponent',
  components: {
    ConfirmModal,
    ModifyModal
  },
  data() {
    return {
      headers: ['ID', 'CreatedAt', 'UpdatedAt', 'DeletedAt', 'UserID', 'Symbol', 'Quantity', 'Price', 'OrderType', 'Direction', 'Status', 'Operate'],
      tableData: [],
      loading: true,
      error: null,
      isConfirmModalVisible: false,
      selectedIndex: -1,
      isModifyModalVisible: false,
      selectedOrderQuantity: null,
      selectedOrderPrice: null,
      selectedOrderId: null
    };
  },
  async created() {
    await this.fetchUserProducts();
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
        const response = await axios.get('/seller/orders', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        console.log('Response data:', response.data);

        // 处理返回的数据
        if (response.data && Array.isArray(response.data)) {
          this.tableData = response.data.map(product => [
            product.ID,
            this.formatDate(product.CreatedAt),
            this.formatDate(product.UpdatedAt),
            this.formatDate(product.DeletedAt),
            product.UserID,
            product.Symbol,
            product.Quantity,
            product.Price,
            product.OrderType,
            product.Direction,
            product.Status
          ]);
        } else {
          this.tableData = [];
        }
      } catch (error) {
        console.error('Fetch user products error:', error);
        this.error = error.response?.data?.message || 'Fetch user products failed';
      } finally {
        this.loading = false;
      }
    },
    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      });
    },
    confirmDelete(index) {
      this.selectedIndex = index;
      this.isConfirmModalVisible = true;
    },
    closeConfirmModal() {
      this.isConfirmModalVisible = false;
    },
    async deleteOrder(index) {
      try {
        // 获取订单的 ID
        const orderId = this.tableData[index][0]; // 假设订单 ID 是每行的第一个元素

        // 从 localStorage 获取 token
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error('Please login first');
        }

        // 发送 DELETE 请求
        const response = await axios.delete(`/seller/orders/${orderId}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          },
         // data: {
          //  orderId: orderId
         // }
        });

        // 检查响应状态码
        if (response.status === 200) {
          alert(response.data.message);
          //console.log('Order deleted successfully');
          // 从 tableData 中删除该订单
          //this.tableData.splice(index, 1);
        } else {
          alert('Failed to delete order');
          //console.error('Failed to delete order');
        }
      } catch (error) {
        alert(error);
        //console.error
        //console.error('Error deleting order:', error);
        console.error('Failed to delete order');
      } finally {
        this.closeConfirmModal();
      }
    },
    showModifyModal(index) {
      this.selectedIndex = index;
      const order = this.tableData[index];
      this.selectedOrderQuantity = order[6]; // 假设 Quantity 是第 7 列
      this.selectedOrderPrice = order[7]; // 假设 Price 是第 8 列
      this.selectedOrderId = order[0]; // 假设订单 ID 是第 1 列
      this.isModifyModalVisible = true;
    },
    closeModifyModal() {
      this.isModifyModalVisible = false;
    },
    async modifyOrderSubmit({ quantity, price }) {
      try {
        // 从 localStorage 获取 token
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error('Please login first');
        }

        // 获取用户 ID
        const userId = this.tableData[this.selectedIndex][4]; // 假设用户 ID 是第 5 列

        // 发送 POST 请求
        const response = await axios.put(`/seller/orders/${this.selectedOrderId}`, {
          quantity: quantity,
          price: price
        }, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        // 检查响应状态码
        if (response.status === 200) {
          alert('Modify successful! ' + response.data.message);
          //console.log('Order modified successfully');
          // 更新 tableData 中的订单信息
          this.tableData[this.selectedIndex][6] = quantity; // 更新 Quantity
          this.tableData[this.selectedIndex][7] = price; // 更新 Price
        } else {
          console.error('Failed to modify order');
          //alert(response.data.message);
        }
      } catch (error) {
        alert(error);
        //console.error('Error modifying order:', error);
        console.error('Failed to modify order');
      } finally {
        this.closeModifyModal();
      }
    },
    goToCreate() {
      this.$router.push('/create');
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