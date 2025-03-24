
<template>
    <div class="result-container">
      <h3 class="header">Результат поиска</h3>
      
      <div v-if="isLoading" class="loading">Идет загрузка...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="searchResults" class="result-content">
        <p>{{ searchResults }}</p>
      </div>
      <div v-else class="no-data">Нет данных</div>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  
  export default {
    props: ['query'],
    data() {
      return {
        searchResults: '',
        isLoading: false,
        error: null
      };
    },
    watch: {
      query: {
        immediate: true,
        handler(newQuery) {
          if (newQuery) {
            this.fetchData(newQuery);
          }
        }
      }
    },
    methods: {
      async fetchData(query) {
        this.isLoading = true;
        this.error = null;
        try {
          const response = await axios.get('/api/qck/shared-data-link/', {
            params: { link: query },
          });
          this.searchResults = response.data.description || 'Нет данных';
        } catch (err) {
          console.error('Ошибка при запросе:', err);
          this.error = 'Ошибка при загрузке данных';
        } finally {
          this.isLoading = false;
        }
      }
    }
  };
  </script>
  
  <style lang="scss" scoped>
  .result-container {
    margin-top: 20px;
    max-width: 600px;
    padding: 15px;
    background-color: #f9f9f9;
    border: 1px solid #ccc;
    border-radius: 5px;
    text-align: left;
    word-wrap: break-word;
    white-space: pre-wrap;
  }
  
  .header {
    color: teal;
    text-align: center;
    margin-bottom: 10px;
  }
  
  .loading, .error, .no-data {
    text-align: center;
    font-size: 16px;
    margin-top: 10px;
  }
  </style>
  