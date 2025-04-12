<template>
    <div v-if="isVisible" class="modal-overlay">
      <div class="modal">
        <h3>{{ title }}</h3>
        <p>{{ message }}</p>
        <div class="modal-buttons">
          <button @click="confirm">Confirm</button>
          <button @click="cancel">Cancel</button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      isVisible: {
        type: Boolean,
        default: false
      },
      title: {
        type: String,
        default: 'Confirm'
      },
      message: {
        type: String,
        default: 'Are you sure?'
      }
    },
    methods: {
      confirm() {
        this.$emit('confirm');
      },
      cancel() {
        this.$emit('cancel');
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
    text-align: center;
  }
  
  .modal-buttons {
    margin-top: 20px;
  }
  
  button {
    margin: 0 10px;
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #45a049;
  }
  </style>