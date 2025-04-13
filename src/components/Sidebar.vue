<template>
    <div class="sidebar">
      <h3>profile</h3>
      <div v-if="user">
        <ul>
        <li><a href="#">User:{{ user.username }}</a></li>
        <!-- <li><a href="#">More</a></li> -->
        <!-- <li><a href="#">Serve</a></li> -->
        <!-- <li><a href="#">Contact Us</a></li> -->
      </ul>
      </div>
  
    </div>
  </template>

  <script>
import { ref, onMounted } from 'vue';

  export default {
    name: 'Sidebar',
    // setup() {
    //   const user = inject('user');
    //   return { user }
    // }
    setup() {
    const user = ref(null);

    onMounted(() => {
      const storedUser = localStorage.getItem('user');
      if (storedUser) {
        user.value = JSON.parse(storedUser);
      }
    });
    
    return {
      user
    };
  }
  };
  </script>
  
  <style scoped>
  .sidebar {
    width: 200px;
    background-color: #696767;
    padding: 20px;
    border-radius: 5px;
    /* box-shadow: 0 2px 4px rgba(255, 255, 255, 0.1); */
  }
  
  h3 {
    margin-bottom: 15px;
  }
  
  ul {
    list-style-type: none;
    padding: 0;
  }
  
  li {
    margin-bottom: 10px;
  }
  
  a {
    text-decoration: none;
    color: #fcfcfc;
    font-weight: bold;
  }
  
  a:hover {
    text-decoration: underline;
  }
  </style>