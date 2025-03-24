<template>
  <div class="password-container">
    <h3 class="password-title">Новый пароль</h3>
    <form @submit.prevent="saveNewPassword" class="password-form">
      <my-input 
        v-model="newPassword" 
        type="password" 
        placeholder="Введите новый пароль" 
        required 
        v-focus 
        class="password-input"
      />
      <my-button type="submit" class="password-button">Изменить</my-button>

      <div v-if="serverMessage" :class="{'success-message': isSuccess, 'error-message': !isSuccess}">
        {{ serverMessage }}
      </div>
    </form>
  </div>
</template>

<script>
import api from '@/services/api';

export default {
  name: 'NewPasswordForm',
  mounted() {
    // Извлечение токена из URL
    this.token = this.$route.query.token;
  },
  data() {
    return {
      newPassword: '',
      serverMessage: '',
      isSuccess: false,
    };
  },
  methods: {
    async saveNewPassword() {
      try {
        if (this.newPassword.length < 8){
          this.serverMessage='Пароль должен быть 8 или более символов'
          return
        }
        const response = await api.post('/users/reset-password', {
          token: this.token,
          password: this.newPassword
        });

        this.isSuccess = response.data.success;
        this.serverMessage = response.data.message;

        // Очищаем сообщение через 3 секунды
        setTimeout(() => {
          this.serverMessage = '';
        }, 3000);

      } catch (error) {
        this.isSuccess = false;
        this.serverMessage = error.response?.data?.message || 'Ошибка при смене пароля';
      }
    }
  }
}
</script>

<style>
.password-container {
  max-width: 400px;
  margin: 20px auto;
  padding: 15px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.password-title {
  color: teal;
  text-align: center;
  margin-bottom: 15px;
}

.password-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.password-input {
  width: 100%;
}

.password-button {
  width: 100%;
}

.success-message,
.error-message {
  text-align: center;
  padding: 10px;
  margin-top: 10px;
  border-radius: 5px;
}

.success-message {
  background-color: #d4edda;
  color: #155724;
}

.error-message {
  background-color: #f8d7da;
  color: #721c24;
}

/* Адаптивность */
@media (max-width: 480px) {
  .password-container {
    width: 90%;
    padding: 10px;
  }

  .password-form {
    gap: 8px;
  }

  .password-button {
    font-size: 14px;
  }
}
</style>
