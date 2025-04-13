<template>
    <div v-if="isVisible" class="modal-overlay">
      <div class="modal">
        <h2>{{ title }}</h2>
        <form @submit.prevent="submit">
          <div class="form-group">
            <label for="quantity">Quantity</label>
            <input type="number" id="quantity" v-model="quantity" required />
          </div>
          <div class="form-group">
            <label for="price">Price</label>
            <input type="number" id="price" v-model="price" required />
          </div>
          <button type="submit">Submit</button>
          <button type="button" @click="close">Cancel</button>
        </form>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'ModifyModal',
    props: {
      isVisible: Boolean,
      title: String,
      initialQuantity: Number,
      initialPrice: Number
    },
    data() {
      return {
        quantity: this.initialQuantity,
        price: this.initialPrice
      };
    },
    methods: {
      submit() {
        this.$emit('submit', { quantity: this.quantity, price: this.price });
      },
      close() {
        this.$emit('close');
      }
    }
  };
  </script>
  
  <style scoped>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .modal {
    background: white;
    padding: 20px;
    border-radius: 5px;
    width: 300px;
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
    margin: 5px;
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