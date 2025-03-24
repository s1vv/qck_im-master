<template>
  <div class="search-container">
    <h3 class="header">Информация по ссылке</h3>

    <div v-if="!searchResults" class="input-container">
      <my-input v-model="query" placeholder="Введите буквы после qck.im/" :maxlength="100" required class="input-field"/>
      <my-button @click="search" class="search-button">Поиск</my-button>
    </div>

    <!-- Индикатор загрузки -->
    <div v-if="isLoading" class="loading">Идет поиск...</div>

    <!-- Вывод ошибки -->
    <div v-if="errorInput" class="result-container" v-html="errorMarkdown"></div>

    <!-- Вывод результата -->
    <div v-if="searchResults" class="result-container" v-html="searchMarkdown"></div>

    <!-- Telegram-ссылка отображается только если нет searchResults -->
    <div v-if="!searchResults" class="telegram-container">
      <a href="https://t.me/qck_im" target="_blank" class="telegram-link">
        Telegram t.me/qck_im
      </a>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import { formatMarkdown } from '@/services/markdown';

export default {
  data() {
    return {
      query: '',
      searchResults: '',
      errorInput: '',
      isLoading: false
    };
  },
  computed: {
    errorMarkdown() {
      return this.errorInput ? formatMarkdown(this.errorInput) : '';
    },
    searchMarkdown() {
      return this.searchResults ? formatMarkdown(this.searchResults) : '';
    }
  },
  mounted() {
    if (this.$route.params.code) {
      this.query = this.$route.params.code;
      this.search();
    }
  },
  methods: {
    async search() {
      this.errorInput = ''; // Очистка ошибки перед запросом
      this.searchResults = ''; // Очистка предыдущих результатов
      
      if (this.query.length !== 8) {
        this.errorInput = 'Ссылка должна быть из 8 букв';
        return;
      }
      
      this.isLoading = true;
      try {
        const response = await axios.get('/api/qck/shared-data-link/', {
          params: { link: this.query },
        });
        this.searchResults = response.data.description || 'Нет данных';
      } catch (error) {
        console.error('Ошибка при поиске:', error);
        this.searchResults = 'Информация не найдена';
      } finally {
        this.isLoading = false;
      }
    }
  }
};
</script>


<style lang="scss" scoped>
.search-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  padding: 20px;
}

.header {
  color: teal;
  text-align: center;
  margin-bottom: 10px;
}

.input-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  max-width: 600px;
  width: 100%;
  margin-bottom: 20px;
}

.input-field {
  width: 100%;
  max-width: 400px;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-bottom: 10px;
}

.search-button {
  width: 100%;
  max-width: 100px;
  padding: 10px;
  font-size: 16px;
  margin-top: 10px;
  border-radius: 10px;
}

.loading {
  text-align: center;
  font-size: 16px;
  margin-top: 20px;
}

/* Контейнер для отображения результата */
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

/* Telegram-ссылка внизу страницы */
.telegram-container {
  position: absolute;
  bottom: 10px;
  width: 100%;
  text-align: center;
}

.telegram-link {
  font-size: 16px;
  color: #0088cc;
  text-decoration: none;
}

.telegram-link:hover {
  text-decoration: underline;
}
</style>
