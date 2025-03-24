<template>
  <div>
    <h1>Список пользователей</h1>
    <ul>
      <li v-for="note in users" :key="generateKey(note)">{{ note.title }}: {{ note.body }}</li>
    </ul>
    <button @click="fetchUsers">Загрузить пользователей</button>
  </div>
</template>

<script>
import axios from 'axios';
import { v4 as uuidv4 } from 'uuid'

export default {
  data() {
    return {
      users: []
    };
  },
  methods: {
    generateKey() {
      return uuidv4();
    },
    async fetchUsers() {
      try {
        const token = localStorage.getItem('token'); 
        const response = await axios.get('/api/users', {
          headers: {
            'Authorization': `Bearer ${token}` // Передаем токен в заголовке запроса с использованием схемы Bearer
          }
        });
        this.users = response.data;
      } catch (error) {
        console.error('Ошибка при получении данных:', error);
      }
    }
  }
};
</script>
