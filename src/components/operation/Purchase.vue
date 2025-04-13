<template>
  <div class="modal" v-if="isVisible">
    <div class="modal-content">
      <h2>Purchase {{ order[0] }}?</h2>
      <form @submit.prevent="submitOrder">
        <div class="form-group">
          <label for="quantity">Quantity:</label>
          <input
            type="number"
            id="quantity"
            v-model.number="quantity"
            :max="order[1]"
            required
          />
        </div>
        <div class="form-group">
          <label for="price">Price:</label>
          <input
            type="number"
            id="price"
            v-model.number="price"
            :disabled="isMarketOrder"
            :max="10000"
            required
          />
        </div>
        <button type="submit">Submit</button>
        <button @click="closeModal">Cancel</button>
      </form>
    </div>
  </div>
</template>
  
  <script>
  import axios from 'axios';
  
  export default {
    name: 'Purchase',
    props: {
      order: Object
    },
    data() {
      return {
        quantity: null,
        price: null,
        isVisible: true
      };
    },
    computed: {
      isMarketOrder() {
        return this.order[3] === 'market';
      }
    },
    // watch: {
    //   isMarketOrder(newVal) {
    //     if (newVal) {
    //       this.price = this.order[2];
    //     }
    //   }
    // },
    methods: {
      closeModal() {
        this.isVisible = false;
        this.$emit('close');
      },
      async submitOrder() {
        console.log(this.order);
        if (this.quantity > this.order[1]) {
          alert('Quantity exceeds available stock');
          return;
        }
  
        if (this.isMarketOrder && this.price !== this.order[2]) {
          alert('Price cannot be changed for market orders');
          return;
        }
  
        if (this.price > 10000) {
          alert('Price cannot exceed 10000');
          return;
        }
  
        const neworder = {
          Symbol: this.order[0],
          Quantity: this.quantity,
          Price: this.price
        };
  
        try {
          const token = localStorage.getItem('token');
          if (!token) {
            throw new Error('Please login first');
          }
  
          const response = await axios.post('/purchase', neworder, {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          });
  
          if (response.status === 201) {
            alert('Order submitted successfully' + (response.data.message || ''));
            this.closeModal();
          } else {
            alert('Order submission failed' + (response.data.message || ''));
          }
        } catch (error) {
          alert(`Order submission failed: ${error.response?.data?.message || 'Unknown error'}`);
        }
      }
    },
    mounted() {
      if (this.isMarketOrder) {
        this.price = this.order[2];
      }
    }
  };
  </script>
  
  <style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 5px;
  width: 300px; /* 设置一个固定宽度，使布局更美观 */
}

form {
  display: flex;
  flex-direction: column;
}

.form-group {
  display: flex;
  flex-direction: column;
  margin-bottom: 10px;
}

label {
  margin-bottom: 5px;
  font-weight: bold; /* 让标签更醒目 */
}

input {
  margin-bottom: 10px;
  padding: 8px; /* 增加内边距，让输入框更美观 */
  border: 1px solid #ccc; /* 添加边框 */
  border-radius: 4px; /* 添加圆角 */
}

button {
  margin: 5px 0;
  padding: 8px 16px; /* 增加按钮的内边距，使其更美观 */
  background-color: #007bff; /* 蓝色背景 */
  color: white; /* 白色文字 */
  border: none;
  border-radius: 4px; /* 圆角 */
  cursor: pointer; /* 鼠标悬停时显示手型 */
}

button:hover {
  background-color: #0056b3; /* 鼠标悬停时的背景颜色 */
}
</style>