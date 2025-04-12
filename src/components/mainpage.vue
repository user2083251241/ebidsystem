<template>
    <div class="main-page">
      <h1>EBidsystem</h1>
      <div class="content">
        <Sidebar v-if="showSidebar"/>
        <div class="main-content">
          <p>Main area</p>
          <button class="toggle-sidebar-button" @click="toggleSidebar">&#9776;</button>
          <TableComponent />
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import Sidebar from './Sidebar.vue';
  import TableComponent from './TableComponent.vue';
  
  export default {
    components: {
      Sidebar,
      TableComponent
    },
    data() {
      return {
        showSidebar: true // 控制侧边栏的显示和隐藏
      };
    },
    methods: {
      toggleSidebar() {
        this.showSidebar = !this.showSidebar; // 切换侧边栏的显示状态
      }
    }
  };
  </script>
  
  <style scoped>
  .main-page {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-top: 50px;
  }
  
  .content {
    display: flex;
    width: 100%;
    max-width: 1200px;
    margin-top: 20px;
  }
  
  .sidebar {
    margin-right: 20px;
  }
  
  .main-content {
    flex: 1;
    background-color: #b7b4b4;
    padding: 20px;
    border-radius: 5px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: flex 0.3s ease;
    position: relative; /* 添加相对定位 */
  }
  
  .toggle-sidebar-button {
    position: absolute; /* 添加绝对定位 */
    top: 10px; /* 距离顶部 10px */
    left: 10px; /* 距离左侧 10px */
    background-color: #4caf50;
    color: white;
    border: none;
    border-radius: 5px;
    padding: 5px 10px;
    cursor: pointer;
  }
  
  .toggle-sidebar-button:hover {
    background-color: #45a049;
  }
  
  /* 响应式样式 */
  @media (max-width: 768px) {
    .content {
      flex-direction: column;
    }
  
    .sidebar {
      margin-right: 0;
      margin-bottom: 20px;
    }
  
    .main-content {
      width: 100%;
    }
  }
  </style>