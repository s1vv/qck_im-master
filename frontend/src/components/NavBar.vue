<template>
  <div class="navbar">
    <my-button @click="$router.push('/')">QCK.IM</my-button>
    <div class="navbar__btns">
      <!-- Показываем кнопку "Мои заметки", если есть токен -->
      <my-button v-if="isAuthenticated" @click="$router.push('/notes')">Мои заметки</my-button>
      <div class="navbar__spacer"></div>
      <my-button v-if="isAuthenticated" @click="logout">Выйти</my-button>
      <my-button v-else @click="$router.push('/login')" class="submit-button">Войти</my-button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      isAuthenticated: !!localStorage.getItem('token') // Проверка наличия токена
    };
  },
  methods: {
    async logout() {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          console.warn("Токен отсутствует в localStorage");
          return;
        }

        const response = await axios.post('/api/users/logout', {}, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        localStorage.removeItem('token');
        this.isAuthenticated = false;
        this.$router.push('/login');
      } catch (error) {
        console.error('Ошибка при выходе:', error);
      }
    }
  },
  watch: {
    '$route': function() {
      this.isAuthenticated = !!localStorage.getItem('token');
    }
  }
};
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  background-color: whitesmoke;
  box-shadow: 2px 2px 4px gray;
  display: flex;
  align-items: center;
  padding: 0 15px;
}
.navbar__btns {
  margin-left: auto;
  display: flex;
}
.navbar__spacer {
  width: 15px;
}
</style>
