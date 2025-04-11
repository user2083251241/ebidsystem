<template>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th v-for="header in headers" :key="header">{{ header }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, rowIndex) in tableData" :key="rowIndex">
            <td v-for="(cell, cellIndex) in row" :key="cellIndex">{{ cell }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script>
  export default {
    name: 'TableComponent',
    data() {
      return {
        headers: ['#', 'Cusip', '3', '4', '5', '6', '7', '8', '9', '10'],
        tableData: [
          [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
          [11, 12, 13, 14, 15, 16, 17, 18, 19, 20],
          [21, 22, 23, 24, 25, 26, 27, 28, 29, 30],
          [31, 32, 33, 34, 35, 36, 37, 38, 39, 40],
          [41, 42, 43, 44, 45, 46, 47, 48, 49, 50]
        ]
      };
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
  </style>